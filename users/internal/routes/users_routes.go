package users_routes

import (
	users_controller "github.com/Marian-Sofia/mille_tendresse/users/internal/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes () *gin.Engine {
	c := users_controller.NewUserController()

	r := gin.Default()

	api := r.Group("/users")

	{
		api.GET("/", c.GetUsers)
		api.GET("/:userId", c.GetUserById)
		api.POST("/create", c.CreateUser)
		api.PATCH("/update/:userId", c.UpdateUser)
		api.DELETE("/delete/:userId", c.DeleteUser)
	}

	return r
}