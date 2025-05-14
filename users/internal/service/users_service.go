package users_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	users_repository "github.com/Marian-Sofia/mille_tendresse/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Definición del servicio de usuarios
type usersService struct {
	repository      users_interfaces.IUsersRepository // Repositorio para manejar la base de datos
	redisRepository users_repository.RedisRepository   // Repositorio para manejar el cache de Redis
}

// Constructor de la estructura usersService
func NewUserService() users_interfaces.IUsersService {
	// Retorna una nueva instancia de usersService con los repositorios necesarios
	return &usersService{
		repository:      users_repository.NewUserRepository(config.DB),    // Inicializa el repositorio de usuarios con la base de datos
		redisRepository: *users_repository.NewRedisRepository(config.Client), // Inicializa el repositorio de Redis
	}
}

// Método para obtener todos los usuarios
func (srv *usersService) GetUsers() ([]users_model.User, error) {
	// Primero intenta obtener los usuarios desde Redis (cache)
	usersRedis, err := srv.redisRepository.GetCache("users")
	if err != nil {
		log.Fatal(err) // En caso de error, lo maneja fatalmente
	}

	// Si los usuarios están en cache, los parsea y los retorna
	var usersCache []users_model.User
	err = json.Unmarshal([]byte(usersRedis), &usersCache)
	if err != nil {
		fmt.Println(err) // Error en la deserialización
	}

	// Si hay datos en cache, los devuelve
	if len(usersCache) == 0 {
		fmt.Println("No data in cache") // Si no hay datos en cache, pasa al siguiente paso
	} else {
		return usersCache, nil // Retorna los datos cacheados
	}

	// Si no se encuentran en cache, busca los usuarios en la base de datos
	usersDB, err := srv.repository.Find()
	if err != nil {
		return nil, err // Si ocurre un error, lo devuelve
	}

	// Si se encuentra en la base de datos, lo guarda en cache (para futuras consultas)
	jsonData, err := json.Marshal(usersDB)
	if err != nil {
		log.Fatal(err) // Error al convertir la data a JSON
	}

	err = srv.redisRepository.SetCache("users", jsonData)
	if err != nil {
		log.Fatal(err) // Error al guardar en el cache de Redis
	}

	// Retorna los usuarios desde la base de datos
	return usersDB, nil
}

// Método para obtener un usuario por ID
func (srv *usersService) GetUserById(userId string) (users_model.User, error) {
	// Intenta obtener el usuario desde Redis (cache)
	userRedis, err := srv.redisRepository.GetCache(userId)
	if err != nil {
		fmt.Println("error cache", err) // Error al obtener desde el cache
	}

	// Convierte los datos de JSON a un objeto de tipo User
	var userCache users_model.User
	err = json.Unmarshal([]byte(userRedis), &userCache)
	if err != nil {
		fmt.Println("Error Unmarshal", err) // Error en la deserialización
	}

	// Si encuentra el usuario en el cache, lo retorna
	if userCache.ID.IsZero() {
		fmt.Println("no data in cache") // Si no está en cache, pasa a la DB
	} else {
		fmt.Println("CacheUserId", userCache)
		return userCache, nil // Retorna el usuario desde el cache
	}

	// Si no está en cache, busca el usuario en la base de datos
	userDB, err := srv.repository.FindById(userId)
	if err != nil {
		return userDB, err // Si ocurre un error, lo retorna
	}

	// Si el usuario se encuentra en la DB, lo guarda en el cache para futuras consultas
	jsonData, err := json.Marshal(userDB)
	if err != nil {
		fmt.Println("Error JsonData", err) // Error al convertir a JSON
	}

	err = srv.redisRepository.SetCache(userId, jsonData)
	if err != nil {
		fmt.Println("Error set cache", err) // Error al guardar en el cache
	}

	// Retorna el usuario desde la base de datos
	return userDB, nil
}

// Método para crear un nuevo usuario
func (srv *usersService) CreateUser(userModel users_model.CreateUser) (string, error) {
	// Valida los campos del usuario (por ejemplo, que no estén vacíos o mal formateados)
	if err := srv.Validatefields(userModel); err != nil {
		return "", err // Si hay error de validación, lo retorna
	}

	// Hashea la contraseña que el usuario ingresa
	hash, _ := HashPassword(userModel.Password)
	fmt.Println("Hash:", hash)

	// Verifica si la contraseña ingresada coincide con la versión hasheada (esto se usa para autenticación)
	match := CheckPasswordHash(userModel.Password, hash)
	fmt.Println("Match:", match)

	// Convierte el modelo CreateUser al modelo User para insertar en la base de datos
	user := users_model.User{
		Role:           2,
		Name:           userModel.Name,
		Email:          userModel.Email,
		Password:       hash, // La contraseña se guarda hasheada
		Phone:          userModel.Phone,
		Identification: userModel.Identification,
		DateOfBirth:    userModel.DateOfBirth,
		Address:        userModel.Address,
	}

	// Inserta el nuevo usuario en la base de datos
	msg, err := srv.repository.Create(user)
	if err != nil {
		fmt.Println(err) // Si ocurre error al crear el usuario, lo muestra
	}

	// Limpia el cache de Redis para asegurar que los datos estén actualizados
	if err := srv.redisRepository.CleanCache(); err != nil {
		fmt.Println(err) // Error al limpiar el cache
	}

	// Retorna el mensaje de éxito o error
	return msg, err
}

// Método para actualizar un usuario
func (srv *usersService) UpdateUser(userId string, updates users_model.UpdateUser) (users_model.User, error) {
	// Convierte el modelo de actualización a un mapa de campos para actualizar
	updateMap := buildUpdateMap(updates)

	// Si no hay campos para actualizar, retorna un error vacío
	if len(updateMap) == 0 {
		return users_model.User{}, nil
	}

	// Realiza la actualización en la base de datos
	msg, err := srv.repository.Update(userId, updateMap)
	if err != nil {
		fmt.Println(err) // Si ocurre error al actualizar el usuario, lo muestra
	}

	// Limpia el cache de Redis
	if err := srv.redisRepository.CleanCache(); err != nil {
		fmt.Println(err) // Error al limpiar el cache
	}

	// Retorna el usuario actualizado
	return msg, err
}

// Método para eliminar un usuario
func (srv *usersService) DeleteUser(userId string) (string, error) {
	// Elimina el usuario de la base de datos
	msg, err := srv.repository.Delete(userId)
	if err != nil {
		fmt.Println(err) // Error al eliminar el usuario
	}

	// Limpia el cache de Redis
	if err := srv.redisRepository.CleanCache(); err != nil {
		fmt.Println(err) // Error al limpiar el cache
	}

	// Retorna el mensaje de éxito o error
	return msg, err
}

// Función para comparar los campos a actualizar
func buildUpdateMap(update users_model.UpdateUser) map[string]interface{} {
	updates := make(map[string]interface{})

	// Revisa cada campo para ver si se tiene un valor que debe actualizarse
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

// Método para validar los campos de entrada de un usuario
func (srv *usersService) Validatefields(userModel users_model.CreateUser) error {
	// Valida que no existan campos duplicados en la base de datos (como email, teléfono, etc.)
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

// Función para hashear la contraseña del usuario
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Hashea la contraseña
	return string(bytes), err
}

// Función para comparar si la contraseña ingresada coincide con el hash almacenado
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) // Compara las contraseñas
	return err == nil
}

func (srv *usersService) AuthUser (userModel users_model.AuthUser) error {
	userData, err := srv.repository.AuthUser(userModel.Email)
	if err != nil {
		return err
	}

	if !CheckPasswordHash(userModel.Password, userData.Password) {
		return errors.New("invalid password") 
	}

	return nil
}