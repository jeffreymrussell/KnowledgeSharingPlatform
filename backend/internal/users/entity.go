package users

import "fmt"

// RegisterUserDTO represents the data transfer object for user registration.
type RegisterUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginUserDTO represents the data transfer object for user login.
type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Username == "" || u.PasswordHash == "" {
		return fmt.Errorf("username and password must not be empty")
	}
	return nil
}

func (u *RegisterUserDTO) Validate() error {
	if u.Username == "" || u.Password == "" {
		return fmt.Errorf("username and password must not be empty")
	}
	return nil
}
