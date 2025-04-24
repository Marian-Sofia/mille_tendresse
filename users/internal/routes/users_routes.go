package users_routes

import (
	users_controller "github.com/Marian-Sofia/mille_tendresse/users/internal/controller"
	"github.com/gin-gonic/gin"
)

// UserRoutes configura las rutas de la API para las operaciones de usuario.
func UserRoutes() *gin.Engine {
	// Creamos un nuevo controlador de usuarios
	c := users_controller.NewUserController()

	// Inicializamos Gin, que es un framework para construir aplicaciones web.
	r := gin.Default()

	// Definimos un grupo de rutas para los endpoints relacionados con usuarios
	api := r.Group("/users") // Prefijo para todas las rutas que sigan
	{
		// Ruta GET para obtener todos los usuarios
		api.GET("/", c.GetUsers)

		// Ruta GET para obtener un usuario por su ID
		api.GET("/:userId", c.GetUserById)

		// Ruta POST para crear un nuevo usuario
		api.POST("/create", c.CreateUser)

		// Ruta PATCH para actualizar un usuario existente
		api.PATCH("/update/:userId", c.UpdateUser)

		// Ruta DELETE para eliminar un usuario
		api.DELETE("/delete/:userId", c.DeleteUser)
	}

	// Retorna el servidor con las rutas definidas
	return r
}
