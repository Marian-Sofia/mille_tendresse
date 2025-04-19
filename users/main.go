package main

import (
	"fmt"
	"log"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	users_routes "github.com/Marian-Sofia/mille_tendresse/users/internal/routes"
)

func main (){
	config.ConnectDB()
	if config.DB == nil {
		log.Fatal("❌ No se pudo conectar a la base de datos")
	}
	fmt.Println("✅ Conectado a MongoDB")

	config.InitRedis()


	router := users_routes.UserRoutes()
	router.Run(":5001")
}