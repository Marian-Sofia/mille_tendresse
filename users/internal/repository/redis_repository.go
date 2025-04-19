package repository

import (
	"time"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	cache *redis.Client
}

func NewRedisRepository (cache *redis.Client) *RedisRepository {
	return &RedisRepository{
		cache: cache,
	}
}

// Metodo para traer datos de cache
func (st *RedisRepository) GetCache (key string) (string, error) {
	value, err := st.cache.Get(config.Ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	return value, err
}

// Metodo para añadir datos al cache
func (st *RedisRepository) SetCache (key string, value []byte) error {
	// luego hay que cambiar el TTLs
	err  := st.cache.Set(config.Ctx, key, value, time.Duration(0)).Err()
	return err
}

// Metodo para borrar datos de cache
func (st *RedisRepository) CleanCache() error {
	err := st.cache.FlushDB(config.Ctx).Err()
	if err !=  nil {
		return err
	}
	return nil
}