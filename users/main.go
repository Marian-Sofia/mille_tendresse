package main

import (
	"fmt"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	users_routes "github.com/Marian-Sofia/mille_tendresse/users/internal/routes"
)

func main (){
	// Se conecta a MongoDB
	config.ConnectDB()
	if config.DB == nil {
		fmt.Errorf("❌ No se pudo conectar a la base de datos")
	}
	fmt.Println("✅ Conectado a MongoDB")

	// Se conecta a Redis
	config.InitRedis()

	// Genera las rutas
	router := users_routes.UserRoutes()
	router.Run(":5001")
}