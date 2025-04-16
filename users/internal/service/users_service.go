package users_service

import (
	"fmt"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	users_repository "github.com/Marian-Sofia/mille_tendresse/users/internal/repository"
)

type usersService struct {
	repository users_interfaces.IUsersRepository
}

func NewUserService() users_interfaces.IUsersService {
	fmt.Println(config.DB)
	return &usersService{
		repository: users_repository.NewUserRepository(config.DB),
	}
}

func (srv *usersService) GetUsers() ([]users_model.User, error) {
	return srv.repository.Find()
}

func (srv *usersService) GetUserById(userId string) (users_model.User, error) {
	return srv.repository.FindById(userId)
}

func (srv *usersService) CreateUser(userModel users_model.CreateUser) (string, error) {
	if err := srv.Validatefields(userModel); err != nil {
		return "", err
	}

	user := users_model.User{
		Role:           2,
		Name:           userModel.Name,
		Email:          userModel.Email,
		Phone:          userModel.Phone,
		Identification: userModel.Identification,
		DateOfBirth:    userModel.DateOfBirth,
		Address:        userModel.Address,
	}

	return srv.repository.Create(user)
}

func (srv *usersService) UpdateUser(userId string, updates users_model.UpdateUser) (users_model.User, error) {
	updateMap := buildUpdateMap(updates)

	if len(updateMap) == 0 {
		return users_model.User{}, nil
	}
	return srv.repository.Update(userId, updateMap)
}

func (srv *usersService) DeleteUser(userId string) (string, error) {
	return srv.repository.Delete(userId)
}

func buildUpdateMap(update users_model.UpdateUser) map[string]interface{} {
	updates := make(map[string]interface{})

	if update.Name != nil {
		updates["name"] = *update.Name
	}
	if update.Email != nil {
		updates["email"] = *update.Email
	}
	if update.Phone != nil {
		updates["phone"] = *update.Phone
	}
	if update.Identification != nil {
		updates["identification"] = *update.Identification
	}
	if update.DateOfBirth != nil {
		updates["dateOfBirth"] = *update.DateOfBirth
	}
	if update.Address != nil {
		updates["address"] = *update.Address
	}

	return updates
}

func (srv *usersService) Validatefields(userModel users_model.CreateUser) error {

	if err := srv.repository.FindFields(userModel.Email, "email"); err != nil {
		return err
	}

	if err := srv.repository.FindFields(userModel.Phone, "phone"); err != nil {
		return err
	}

	if err := srv.repository.FindFields(userModel.Identification, "identification"); err != nil {
		return err
	}

	return nil
}
