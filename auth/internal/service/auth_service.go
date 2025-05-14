package auth_service

import (
	auth_interfaces "github.com/Marian-Sofia/mille_tendresse/auth/internal/interfaces"
	auth_models "github.com/Marian-Sofia/mille_tendresse/auth/internal/models"
	auth_repository "github.com/Marian-Sofia/mille_tendresse/auth/internal/repository"
)

type authService struct {
	repository auth_interfaces.IAuthRepository
}

func NewAuthService() auth_interfaces.IAuthService {
	return &authService {
		repository: auth_repository.NewAuthRepository(),
	}
}

func (srv *authService) Login(userLogin auth_models.AuthLogin) error {
	return srv.repository.Login(userLogin)
}


func (srv *authService) Register(userRegister auth_models.AuthRegister) error {
	return srv.repository.Register(userRegister)
}