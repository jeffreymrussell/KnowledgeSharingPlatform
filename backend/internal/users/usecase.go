package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepository interface {
	SaveUser(user User) error
	GetUserByUsername(username string) (User, error)
}
type UserUsecase struct {
	Repository UserRepository
}

func (usecase *UserUsecase) RegisterUser(registerUserDTO RegisterUserDTO) (User, error) {
	// Validate the User
	if err := registerUserDTO.Validate(); err != nil {
		return User{}, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUserDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		-1,
		registerUserDTO.Username,
		string(hashedPassword),
		registerUserDTO.Email,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
	}

	// Save the User to the database
	return user, usecase.Repository.SaveUser(user)
}

// Secret key to sign the JWT token
var jwtKey = []byte("your_secret_key_here")

// Claims struct to store the username
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// LoginUser function
func (usecase *UserUsecase) LoginUser(username, password string) (string, error) {
	// Retrieve the User from the database
	user, err := usecase.Repository.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// Validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("wrong password")
	}

	// Generate a token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error in generating token")
	}

	return tokenString, nil
}

// LogoutUser function
func (usecase *UserUsecase) LogoutUser(token string) (string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("invalid token signature")
		}
		return "", fmt.Errorf("could not parse token")
	}

	if !tkn.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Invalidate the token (For this example, just returning a confirmation)
	// In a real-world application, you might store invalidated tokens in a Client to check against during authentication
	return "Logged out successfully", nil
}
