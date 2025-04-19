package users_service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Marian-Sofia/mille_tendresse/users/internal/config"
	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	users_repository "github.com/Marian-Sofia/mille_tendresse/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	repository      users_interfaces.IUsersRepository
	redisRepository users_repository.RedisRepository
}

func NewUserService() users_interfaces.IUsersService {
	return &usersService{
		repository:      users_repository.NewUserRepository(config.DB),
		redisRepository: *users_repository.NewRedisRepository(config.Client),
	}
}

func (srv *usersService) GetUsers() ([]users_model.User, error) {
	// Buscando los usuarios en cache
	usersRedis, err := srv.redisRepository.GetCache("users")
	if err != nil {
		log.Fatal(err)
	}

	// Parseando los usuarios que estan guardados en cache que vienen en formato json a type users_model.User
	var usersCache []users_model.User
	err = json.Unmarshal([]byte(usersRedis), &usersCache)
	if err != nil {
		fmt.Println(err)
	}

	// Valida si los datos en cache existen, si exiten los retorna
	if len(usersCache) == 0 {
		fmt.Println("No data in cache")
	} else {
		return usersCache, nil
	}

	// Esta buscando los usuarios en la DB, para la primera consulta
	usersDB, err := srv.repository.Find()
	if err != nil {
		return nil, err
	}

	// Esta parseando los users de la DB a Json para guardarlon en Redis
	jsonData, err := json.Marshal(usersDB)
	if err != nil {
		log.Fatal(err)
	}

	// Esta guardando los usuarios en Redis para la segunda consulta
	err = srv.redisRepository.SetCache("users", jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// Retorna los usuarios de la DB, primera consulta. Luego de la primera consulta retorna los de Redis
	return usersDB, nil
}

func (srv *usersService) GetUserById(userId string) (users_model.User, error) {
	// Buscar el usuario en Redis
	userRedis, err := srv.redisRepository.GetCache(userId)
	fmt.Println("userRedis", userRedis)
	if err != nil {
		fmt.Println("error cache", err)
	}

	// Parseando el usuario que esta guardado en cache en formato json a type users_model.User
	var userCache users_model.User
	err = json.Unmarshal([]byte(userRedis), &userCache)
	if err != nil {
		fmt.Println("Error Unmarshal", err)
	}
	
	// Valida si los datos en cache existen, si exiten los retorna
	if userCache.ID.IsZero() {
		fmt.Println("no data in cache")
	} else {
		fmt.Println("CacheUserId", userCache)
		return userCache, nil
	}

	// Busca el user en la DB
	userDB, err := srv.repository.FindById(userId)
	if err != nil {
		return userDB, err
	}

	// Parsea los datos de la DB a Json
	jsonData, err := json.Marshal(userDB)
	if err != nil {
		fmt.Println("Error JsonData", err)
	}

	// Guarda el usuario en cache
	err = srv.redisRepository.SetCache(userId, jsonData)
	if err != nil {
		fmt.Println("Error set cache", err)
	}

	// Retorna el usuario desde la DB
	return userDB, nil
}

func (srv *usersService) CreateUser(userModel users_model.CreateUser) (string, error) {
	// Trae el metodo para validar los campos que llegan
	if err := srv.Validatefields(userModel); err != nil {
		return "", err
	}

	// hasea la contraseña que llega
	hash, _ := HashPassword(userModel.Password)
	fmt.Println("Hash:", hash)

	// Revisa si la contraseña que llega es la misma que esta guardada (Este luego pasa para el Auth)
	match := CheckPasswordHash(userModel.Password, hash)
	fmt.Println("Match:", match)

	// Esto convertia el modelo CreateUser a el modelo User
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

	// Se crea el usuario en la DB
	msg, err := srv.repository.Create(user)
	if err != nil {
		fmt.Println(err)
	}

	// Limpiar cache de Redis
	if err := srv.redisRepository.CleanCache(); err != nil {
		fmt.Println(err)
	}

	// Retorna el string y el error si hay
	return msg, err
}

func (srv *usersService) UpdateUser(userId string, updates users_model.UpdateUser) (users_model.User, error) {
	// Se Contruye el modelo update por campos
	updateMap := buildUpdateMap(updates)

	// Se verifica si llega algo para actualizar
	if len(updateMap) == 0 {
		return users_model.User{}, nil
	}

	// Se actualiza el user en la DB
	msg, err := srv.repository.Update(userId, updateMap)
	if err != nil {
		fmt.Println(err)
	}

	// Se limpia el cache de Redis
	if err := srv.redisRepository.CleanCache(); err != nil {
		fmt.Println(err)
	}

	// Se retorna un mensaje y el error
	return msg, err
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
