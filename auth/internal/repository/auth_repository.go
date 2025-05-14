package auth_repository

import (
	"bytes"
	"encoding/json"
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
	return httpRequest(rpt.BaseURL, "/users/validate", userLogin)
}

func (rpt *userHTTPRepository) Register(userRegister auth_models.AuthRegister) error {
	return httpRequest(rpt.BaseURL, "/users/create", userRegister)
}

func httpRequest (BaseURL, endpoint string, payload interface{}) error {
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)

	// Serializar el payload a JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize request body: %w", err)
	}

	// Hacer la peticion POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Manejo de errores de respuesta HTTP
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("request failed with status %d and unreadable body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}