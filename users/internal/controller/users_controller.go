package users_controller

import (
	"net/http"

	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	users_service "github.com/Marian-Sofia/mille_tendresse/users/internal/service"
	"github.com/gin-gonic/gin"
)

// usersController es la estructura que manejará las rutas relacionadas a usuarios.
type usersController struct {
	// service es una interfaz que define lo que el controlador puede hacer (crear, obtener, etc.).
	service users_interfaces.IUsersService
}

// NewUserController es un constructor que devuelve una instancia del controlador de usuarios.
func NewUserController() *usersController {
	return &usersController{
		service: users_service.NewUserService(), // Inyecta el servicio (por ahora hardcodeado, se puede mejorar con DI).
	}
}

// GetUsers maneja GET /users y devuelve una lista de todos los usuarios.
func (ctrl *usersController) GetUsers(c *gin.Context) {
	users, err := ctrl.service.GetUsers()
	if err != nil {
		// Si hay error, se responde con 400 y el mensaje de error.
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// Si tod va bien, responde con 200 y la lista de usuarios.
	c.JSON(http.StatusOK, users)
}

// GetUserById maneja GET /users/:userId y devuelve un usuario por su ID.
func (ctrl *usersController) GetUserById(c *gin.Context) {
	// Obtiene el parámetro userId de la URL.
	userId := c.Param("userId")

	user, err := ctrl.service.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser maneja POST /users y crea un nuevo usuario.
func (ctrl *usersController) CreateUser(c *gin.Context) {
	var userModel users_model.CreateUser

	// Intenta bindear el JSON del body del request al struct CreateUser.
	if err := c.ShouldBind(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	// Valida los datos antes de crear el usuario.
	if err := users_model.ValidateUser(userModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llama al servicio para crear el usuario.
	user, err := ctrl.service.CreateUser(userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser maneja PUT /users/:userId y actualiza un usuario existente.
func (ctrl *usersController) UpdateUser(c *gin.Context) {
	// Obtiene el ID del usuario desde la URL.
	userId := c.Param("userId")
	var userModel users_model.UpdateUser

	// Bindea los nuevos datos del usuario. (bindear es llenar un struct con los datos que vienen en una petición)
	c.ShouldBind(&userModel)

	// Llama al servicio para actualizar al usuario.
	user, err := ctrl.service.UpdateUser(userId, userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser maneja DELETE /users/:userId y elimina un usuario por su ID.
func (ctrl *usersController) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")

	user, err := ctrl.service.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// AuthUser maneja POST /users/email para autentificar un usuario por email y contraseña
func (ctrl *usersController) AuthUser (c *gin.Context) {
	var userData users_model.AuthUser 

	if err := c.ShouldBind(&userData); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, nil)
}