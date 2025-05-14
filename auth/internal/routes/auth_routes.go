package auth_routes

import (
	auth_controller "github.com/Marian-Sofia/mille_tendresse/auth/internal/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes () *gin.Engine {
	c := auth_controller.NewAuthController()

	r := gin.Default()

	api := r.Group("/auth")
	{
		api.POST("/login", c.Login)
	}

	return r
}