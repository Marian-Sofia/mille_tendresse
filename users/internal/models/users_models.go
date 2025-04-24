package users_model

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User representa el modelo completo del usuario como se guarda en MongoDB
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// Rol del usuario, representado como número (por ejemplo: 1 = admin, 2 = cliente, etc.)
	Role           uint               `bson:"role" json:"role"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	// Contraseña. No se muestra en las respuestas JSON gracias al json:"-"
	Password       string             `bson:"password" json:"-"`
	Phone          string             `bson:"phone" json:"phone"`
	Identification string             `bson:"identification" json:"identification"`
	// Fecha de nacimiento en formato string (YYYY-MM-DD)
	DateOfBirth    string             `bson:"dateOfBirth" json:"dateOfBirth"`
	Address        string             `bson:"address" json:"address"`
}

// CreateUser es el modelo que se usa cuando se recibe un nuevo usuario por POST.
// Tiene validaciones incluidas usando etiquetas `validate`.
type CreateUser struct {
	Name           string `json:"name" validate:"required,min=3,max=50"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required"`
	Phone          string `json:"phone" validate:"required,numeric"`
	Identification string `json:"identification" validate:"required,numeric"`
	DateOfBirth    string `json:"dateOfBirth" validate:"required,datetime=2006-01-02"`
	Address        string `json:"address" validate:"required"`
}

// UpdateUser es para cuando se actualiza un usuario.
// Todos los campos son opcionales (punteros), así podés actualizar solo lo que cambió.
//  Los campos son punteros (*string) para poder saber si vinieron en el JSON o no. Si no se incluyen, quedan como nil, y no se actualizan.
type UpdateUser struct {
	Name           *string `json:"name"`
	Email          *string `json:"email"`
	// password
	Phone          *string `json:"phone"`
	Identification *string `json:"identification"`
	DateOfBirth    *string `json:"dateOfBirth"`
	Address        *string `json:"address"`
}

// Se crea una instancia global del validador
var validate = validator.New()

// Función que recibe un CreateUser y valida según las etiquetas `validate` del struct.
func ValidateUser(user CreateUser) error {
	return validate.Struct(user)
}
