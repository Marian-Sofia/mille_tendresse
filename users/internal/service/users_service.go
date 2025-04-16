package users_service

import (
	"fmt"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	users_repository "github.com/Marian-Sofia/mille_tendresse/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
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

	hash, _ := HashPassword(userModel.Password)
	fmt.Println("Hash:", hash)

	// Este luego pasa para el Auth
	match := CheckPasswordHash(userModel.Password, hash)
	fmt.Println("Match:", match)

	user := users_model.User{
		Role:           2,
		Name:           userModel.Name,
		Email:          userModel.Email,
		Password:       hash,
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

// Funcion para comparar los campos a actualizar
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

// Funcion para validar los campos
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

// Funcion para hashear password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Funcion para verificar la password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
