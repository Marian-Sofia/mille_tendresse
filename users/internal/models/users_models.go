package users_model

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Role           uint               `bson:"role" json:"role"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	Password       string             `bson:"password" json:"-"`
	Phone          string             `bson:"phone" json:"phone"`
	Identification string             `bson:"identification" json:"identification"`
	DateOfBirth    string             `bson:"dateOfBirth" json:"dateOfBirth"`
	Address        string             `bson:"address" json:"address"`
}

type CreateUser struct {
	Name           string `json:"name" validate:"required,min=3,max=50"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required"`
	Phone          string `json:"phone" validate:"required,numeric"`
	Identification string `json:"identification" validate:"required,numeric"`
	DateOfBirth    string `json:"dateOfBirth" validate:"required,datetime=2006-01-02"`
	Address        string `json:"address" validate:"required"`
}

type UpdateUser struct {
	Name           *string `json:"name"`
	Email          *string `json:"email"`
	// password
	Phone          *string `json:"phone"`
	Identification *string `json:"identification"`
	DateOfBirth    *string `json:"dateOfBirth"`
	Address        *string `json:"address"`
}

var validate = validator.New()

func ValidateUser(user CreateUser) error {
	return validate.Struct(user)
}
