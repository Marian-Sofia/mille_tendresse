package auth_repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	auth_interfaces "github.com/Marian-Sofia/mille_tendresse/auth/internal/interfaces"
	auth_models "github.com/Marian-Sofia/mille_tendresse/auth/internal/models"
)

type userHTTPRepository struct {
	BaseURL string 
}

func NewAuthRepository () auth_interfaces.IAuthRepository {
	return &userHTTPRepository{
		BaseURL: os.Getenv("BASE_URL"), 
	}
}

func (rpt *userHTTPRepository) Login(userLogin auth_models.AuthLogin) error {
	url := fmt.Sprintf("%s/users/email", rpt.BaseURL)

	userData, err :=	json.Marshal(userLogin)
	if err != nil {
		return errors.New("the body could not be serialized")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userData))
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	defer resp.Body.Close()

	return nil
}


func (r *userHTTPRepository) Register(userRegister auth_models.AuthRegister) error {
	return nil
}