package auth_controller

import (
	"fmt"
	"net/http"

	auth_interfaces "github.com/Marian-Sofia/mille_tendresse/auth/internal/interfaces"
	auth_models "github.com/Marian-Sofia/mille_tendresse/auth/internal/models"
	auth_service "github.com/Marian-Sofia/mille_tendresse/auth/internal/service"
	"github.com/gin-gonic/gin"
)

type authController struct {
	service auth_interfaces.IAuthService
}

func NewAuthController() *authController {
	return &authController{
		service: auth_service.NewAuthService(),
	}
}

func (ctrl *authController) Login (c *gin.Context) {
	var dataUser auth_models.AuthLogin
	
	if err := c.ShouldBind(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println("controller", dataUser)

	if err := ctrl.service.Login(dataUser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, "ok")
}