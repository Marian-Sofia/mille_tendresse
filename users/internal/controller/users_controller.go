package users_controller

import (
	"net/http"

	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	users_service "github.com/Marian-Sofia/mille_tendresse/users/internal/service"
	"github.com/gin-gonic/gin"
)

type usersController struct {
	service users_interfaces.IUsersService
}

func NewUserController () *usersController {
	return &usersController{
		service: users_service.NewUserService(),
	}
}

func (ctrl *usersController) GetUsers (c *gin.Context){
	users, err := ctrl.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl *usersController) GetUserById (c *gin.Context){
	userId := c.Param("userId")
	user, err := ctrl.service.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *usersController) CreateUser (c *gin.Context){
	var userModel  users_model.CreateUser
	c.ShouldBind(&userModel)
	user, err := ctrl.service.CreateUser(userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *usersController) UpdateUser (c *gin.Context){
	userId := c.Param("userId")
	var userModel users_model.UpdateUser
	c.ShouldBind(&userModel)
	user, err := ctrl.service.UpdateUser(userId, userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *usersController) DeleteUser (c *gin.Context){
	userId := c.Param("userId")
	user, err := ctrl.service.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}