package auth_repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	fmt.Println("repository", userLogin)
	url := fmt.Sprintf("%s/users/validate", rpt.BaseURL)

	userData, err :=	json.Marshal(userLogin)
	if err != nil {
		return errors.New("the body could not be serialized")
	}
	fmt.Println("repository", userData)


	resp, err := http.Post(url, "application/json", bytes.NewBuffer(userData))
	if err != nil {
		return err
	}
	fmt.Println("repository", resp)

	if resp.StatusCode != http.StatusOK {
		 defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("request failed %v", err)
    }

    return fmt.Errorf("request failed %s", string(bodyBytes))
	}

	defer resp.Body.Close()

	return nil
}


func (r *userHTTPRepository) Register(userRegister auth_models.AuthRegister) error {
	return nil
}