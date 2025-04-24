package repository

import (
	"time"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	"github.com/redis/go-redis/v9"
)

// RedisRepository es una estructura que maneja el acceso a Redis.
// Tiene un cliente Redis ya inicializado.
type RedisRepository struct {
	cache *redis.Client // Cliente que permite conectarse y operar con Redis
}

// NewRedisRepository recibe un cliente Redis y devuelve una nueva instancia del repositorio.
func NewRedisRepository (cache *redis.Client) *RedisRepository {
	return &RedisRepository{
		cache: cache,
	}
}

// GetCache obtiene un valor de Redis usando una clave.
func (st *RedisRepository) GetCache (key string) (string, error) {
	value, err := st.cache.Get(config.Ctx, key).Result()
	if err == redis.Nil { // redis.Nil es un error especial que significa "clave no encontrada"
		// Si la clave no existe, no es un error grave. Devolvemos vacío.
		return "", nil
	}

	return value, err // Si hubo otro error (conexión, etc.), lo devolvemos.
}

// SetCache guarda un valor en Redis con una clave.
// El TTL está en 0, lo que significa que **no expira** automáticamente.
func (st *RedisRepository) SetCache(key string, value []byte) error {
	err := st.cache.Set(config.Ctx, key, value, time.Duration(0)).Err()
	return err
}

// CleanCache borra toda la base de datos de Redis.
// ¡Esto elimina tod el cache! Usar con cuidado.
func (st *RedisRepository) CleanCache() error {
	err := st.cache.FlushDB(config.Ctx).Err()
	if err != nil {
		return err
	}
	return nil
}