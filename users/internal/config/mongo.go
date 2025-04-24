package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB es una variable global que contendrá la instancia de la base de datos MongoDB.
// Esta variable será accesible desde otras partes del proyecto.
var DB *mongo.Database

// ConnectDB establece la conexión con MongoDB.
func ConnectDB() {
	// Intenta obtener la URI de conexión desde la variable de entorno MONGO_URI.
	mongoURI := os.Getenv("MONGO_URI")

	// Si no se define MONGO_URI, se usa una URI por defecto (útil en desarrollo o Docker).
	if mongoURI == "" {
		mongoURI = "mongodb://mongo:27017"
	}

	// Se crea un contexto con un timeout de 10 segundos para evitar conexiones colgadas.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Se asegura de liberar los recursos del contexto al finalizar la función.

	// Se intenta conectar a MongoDB usando la URI.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		// Si ocurre un error, se imprime en consola y se detiene la ejecución.
		log.Fatal("❌ Error connecting to MongoDB:", err)
	}

	// Se selecciona la base de datos 'mille_tendresse' y se asigna a la variable global DB.
	DB = client.Database("mille_tendresse")

	// Se imprime un mensaje de éxito indicando que la conexión fue establecida.
	log.Println("✅ Connected to MongoDB at", mongoURI)
}