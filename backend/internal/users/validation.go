package users

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// ValidateRegisterUserDTO enhances the validation for registering a user.
func ValidateRegisterUserDTO(user RegisterUserDTO) error {
	// Trim whitespace
	user.Email = strings.TrimSpace(user.Email)
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	// Email validation
	if user.Email == "" {
		return errors.New("missing email")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Username validation
	if user.Username == "" {
		return errors.New("missing username")
	}
	if len(user.Username) < 3 || len(user.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !isAlphanumeric(user.Username) {
		return errors.New("username must be alphanumeric")
	}

	// Password validation
	if user.Password == "" {
		return errors.New("missing password")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasPasswordRequirements(user.Password) {
		return errors.New("password must contain at least one number, one uppercase letter, one lowercase letter, and one special character")
	}

	return nil
}

// isValidEmail checks if the email address has a valid format.
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// isAlphanumeric checks if the string contains only alphanumeric characters.
func isAlphanumeric(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(str)
}

// hasPasswordRequirements checks the password for digits, upper & lower case, and special chars.
func hasPasswordRequirements(password string) bool {
	var (
		hasMinLen      = false
		hasUpper       = false
		hasLower       = false
		hasNumber      = false
		hasSpecialChar = false
	)
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecialChar
}
