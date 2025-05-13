package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"
	"unicode"

	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// usersRepository representa el repositorio para acceder a la colección de usuarios en MongoDB.
type usersRepository struct {
	collection *mongo.Collection // Aquí guardamos la colección 'users' de MongoDB
}

// NewUserRepository es el constructor del repositorio de usuarios.
// Recibe una base de datos MongoDB y retorna una nueva instancia de usersRepository.
func NewUserRepository(db *mongo.Database) users_interfaces.IUsersRepository {
	return &usersRepository{
		collection: db.Collection("users"), // Asigna la colección 'users' a la propiedad
	}
}
// Inyección de dependencias: El repositorio recibe la base de datos para poder obtener la colección de usuarios.

// Find devuelve todos los usuarios almacenados en la colección de MongoDB.
func (rpt *usersRepository) Find() ([]users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel() // Asegura que se cancele el contexto al terminar

	var users []users_model.User // Slice para almacenar los usuarios recuperados

	// Realiza una consulta para obtener todos los usuarios
	cursor, err := rpt.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err // Si hay error, lo devuelve
	}
	defer cursor.Close(ctx) // Cierra el cursor cuando termine

	// Decodifica los documentos del cursor en el slice users
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err // Si hay error al leer, lo devuelve
	}
	return users, err // Devuelve la lista de usuarios
}

// FindById busca un usuario por su ID.
func (rpt *usersRepository) FindById(userId string) (users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel()

	var user users_model.User // Variable para almacenar el usuario encontrado

	// Convierte el string del ID a un ObjectID de MongoDB
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return user, err // Si hay error en la conversión, lo devuelve
	}

	// Realiza la búsqueda del usuario por su ObjectID
	err = rpt.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		// Si no encuentra el usuario, devuelve un error personalizado
		if err == mongo.ErrNoDocuments {
			return user, errors.New(" User does not exist")
		}
		return user, err
	}

	return user, nil // Devuelve el usuario encontrado
}

// Create inserta un nuevo usuario en la colección de MongoDB.
func (rpt *usersRepository) Create(userModel users_model.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel()

	// Inserta un nuevo documento (usuario) en la colección
	_, err := rpt.collection.InsertOne(ctx, userModel)
	if err != nil {
		return "", err // Si hay error, lo devuelve
	}

	return "User is created", nil // Devuelve un mensaje de éxito
}

// Update actualiza un usuario existente en la colección MongoDB.
func (rpt *usersRepository) Update(userId string, updates map[string]interface{}) (users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel()

	var updateUser users_model.User // Variable para almacenar el usuario actualizado

	// Convierte el string del ID a un ObjectID de MongoDB
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return users_model.User{}, err // Si hay error en la conversión, lo devuelve
	}

	// Definimos los cambios a aplicar (en formato BSON)
	update := bson.M{
		"$set": updates, // Actualiza solo los campos que se pasen
	}

	// Realiza la actualización y devuelve el documento actualizado
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // Retorna el documento actualizado
	err = rpt.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updateUser)

	if err != nil {
		return updateUser, err // Si hay error, lo devuelve
	}

	return updateUser, nil // Devuelve el usuario actualizado
}

// Delete elimina un usuario de la colección MongoDB por su ID.
func (rpt *usersRepository) Delete(userId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel()

	// Convierte el string del ID a un ObjectID de MongoDB
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err // Si hay error en la conversión, lo devuelve
	}

	// Elimina el usuario por su ObjectID
	_, err = rpt.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		// Si no hay documento para eliminar, lo maneja como un error de "no encontrado"
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}

	return "User has been deleted", err // Devuelve un mensaje de éxito
}


// FindFields busca si un campo específico de un usuario (por ejemplo, email) ya existe en la base de datos.
func (rpt *usersRepository) FindFields(userModel string, field string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel()

	var user users_model.User // Variable para almacenar el usuario encontrado

	// Busca un documento que tenga el valor del campo proporcionado
	errUser := rpt.collection.FindOne(ctx, bson.M{field: userModel}).Decode(&user)
	if errUser != nil {
		// Si no lo encuentra, devuelve nil
		if errUser == mongo.ErrNoDocuments {
			return nil
		}
		return errUser // Si hubo un error, lo devuelve
	}

	// Usamos reflexión para verificar si el valor del campo coincide
	val := reflect.ValueOf(user)
	f := val.FieldByName(string(unicode.ToUpper(rune(field[0]))) + field[1:])

	if f.IsValid() && f.String() == userModel {
		return fmt.Errorf("%s already exists", field) // Si el valor existe, devuelve un error
	}
	return nil // Si tod está bien, retorna nil (sin errores)
}

func (rpt *usersRepository) AuthUser (email string) (users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Establece un tiempo de espera
	defer cancel()

	var dataUser users_model.User

	err := rpt.collection.FindOne(ctx, bson.M{"email": email}).Decode(&dataUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return dataUser, errors.New(" User does not exist")
		}
		return dataUser, err
	}

	return dataUser, nil
}