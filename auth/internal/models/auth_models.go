package auth_models

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRegister struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Phone          string `bson:"phone" json:"phone"`
	Identification string `json:"identification"`
	DateOfBirth    string `json:"dateOfBirth"`
	Address        string `json:"address"`
}
