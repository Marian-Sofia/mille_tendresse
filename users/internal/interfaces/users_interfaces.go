package users_interfaces

import users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"

type IUsersService interface {
	GetUsers() ([]users_model.User, error)
	GetUserById(userId string) (users_model.User, error)
	CreateUser(users_model.CreateUser) (string, error)
	UpdateUser(userId string, userModel users_model.UpdateUser) (users_model.User, error)
	DeleteUser(userId string) (string, error)
	Validatefields(userModel users_model.CreateUser) error
}

type IUsersRepository interface {
	Find() ([]users_model.User, error)
	FindById(userId string) (users_model.User, error)
	Create(model users_model.User) (string, error)
	Update(userId string, userModel map[string]interface{}) (users_model.User, error)
	Delete(userId string) (string, error)
	FindFields(userModel string, field string) error
}
