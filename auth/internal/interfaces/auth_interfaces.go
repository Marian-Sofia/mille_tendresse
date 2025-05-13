package auth_interfaces

import "github.com/Marian-Sofia/mille_tendresse/auth/internal/models"

type IAuthService interface {
	Login(userLogin auth_models.AuthLogin) error
	Register(userRegister auth_models.AuthRegister) error
}

type IAuthRepository interface {
	Login(userLogin auth_models.AuthLogin) error
	Register(userRegister auth_models.AuthRegister) error
}