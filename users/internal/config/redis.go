package config

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/redis/go-redis/v9"
)

// Client es la instancia global del cliente Redis que se usará en toda la aplicacion.
var Client *redis.Client

// Ctx es un contexto global (sin timeout) que se usará para todas las operaciones Redis.
// También podrías crear contextos con timeout por operación si querés mayor control.
var Ctx = context.Background()

// InitRedis inicializa la conexión a Redis.
func InitRedis() {
	// Se imprimen las variables de entorno para debug (podés quitar esto en producción).
	fmt.Println("REDIS_ADDR", os.Getenv("REDIS_ADDR"), "REDIS_PASSWORD", os.Getenv("REDIS_PASSWORD"))

	// Se crea un nuevo cliente Redis usando los valores de las variables de entorno.
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // Dirección del servidor Redis, ej: "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // Contraseña si está configurada
		DB:       0,                           // Número de base de datos (por defecto es 0)
	})

	// Se hace un ping para verificar que la conexión a Redis funcione correctamente.
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		// Si falla la conexión, se termina la ejecución y se imprime el error.
		log.Fatalf("❌ Unable to connect to Redis: %v", err)
	}

	// Mensaje de éxito si la conexión se realizó correctamente.
	log.Println("✅ Connected to Redis correctly")
}