package models

// Credentials represents the login credentials.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
