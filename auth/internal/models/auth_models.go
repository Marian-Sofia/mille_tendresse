package auth_models

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

type AuthRegister struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	Phone          string `bson:"phone" json:"phone"`
	Identification string `json:"identification"`
	DateOfBirth    string `json:"dateOfBirth"`
	Address        string `json:"address"`
}
