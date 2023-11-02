package users

import (
	"testing"
)

func TestValidateRegisterUserDTO(t *testing.T) {
	tests := []struct {
		name    string
		user    RegisterUserDTO
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "Password123!",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			user: RegisterUserDTO{
				Username: "testuser",
				Password: "Password123!",
			},
			wantErr: true,
			errMsg:  "missing email",
		},
		{
			name: "invalid email format",
			user: RegisterUserDTO{
				Email:    "test",
				Username: "testuser",
				Password: "Password123!",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "missing username",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			wantErr: true,
			errMsg:  "missing username",
		},
		{
			name: "username too short",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "tu",
				Password: "Password123!",
			},
			wantErr: true,
			errMsg:  "username must be between 3 and 20 characters",
		},
		{
			name: "username too long",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "thisisaverylongusernamethatexceedsthelimit",
				Password: "Password123!",
			},
			wantErr: true,
			errMsg:  "username must be between 3 and 20 characters",
		},
		{
			name: "username with special characters",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "test$user",
				Password: "Password123!",
			},
			wantErr: true,
			errMsg:  "username must be alphanumeric",
		},
		{
			name: "missing password",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "testuser",
			},
			wantErr: true,
			errMsg:  "missing password",
		},
		{
			name: "weak password",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "pass",
			},
			wantErr: true,
			errMsg:  "password must be at least 8 characters long",
		},
		{
			name: "password without special character",
			user: RegisterUserDTO{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "Password123",
			},
			wantErr: true,
			errMsg:  "password must contain at least one number, one uppercase letter, one lowercase letter, and one special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisterUserDTO(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegisterUserDTO() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ValidateRegisterUserDTO() gotErr = %v, wantErr %v", err.Error(), tt.errMsg)
			}
		})
	}
}
