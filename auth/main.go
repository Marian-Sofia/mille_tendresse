package main

import auth_routes "github.com/Marian-Sofia/mille_tendresse/auth/routes"

func main () {
		// Genera las rutas
		router := auth_routes.AuthRoutes()
		router.Run(":5001")
}