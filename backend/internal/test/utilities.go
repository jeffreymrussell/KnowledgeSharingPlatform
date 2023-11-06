package test

import (
	"KnowledgeSharingPlatform/internal/users"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
)

type KSPTestSuite struct {
	suite.Suite
	registeredUser users.LoginUserDTO
	Server         *httptest.Server
}

// AuthResponse represents the response object containing the authentication token.
type AuthResponse struct {
	Token string `json:"token"`
}

func GetBodyAsString(res *http.Response) []byte {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	data, _ := io.ReadAll(res.Body)
	return data
}

func RegisterUserAndLogin(t require.TestingT, url string) (users.LoginUserDTO, string) {
	registeredUser := users.LoginUserDTO{
		Username: "registereduser",
		Password: "Registeredpassword1!",
	}
	payload := "{\"username\": \"registereduser\", \"password\": \"Registeredpassword1!\", \"email\": \"test@example.com\"}"

	registerURL := fmt.Sprintf("%s/register", url)

	req, _ := http.NewRequest(http.MethodPost, registerURL, bytes.NewBuffer([]byte(payload)))
	res, err := http.DefaultClient.Do(req)
	data := GetBodyAsString(res)
	require.NoError(t, err, "Shouldn't have an error running the registration")
	if !(res.StatusCode == 400 && string(data) == "UNIQUE constraint failed: users.email\n") {
		require.Equal(t, 201, res.StatusCode, fmt.Sprintf("201 for user registration: %s", data))
	}

	authPayload, _ := json.Marshal(registeredUser)
	loginURL := fmt.Sprintf("%s/login", url)
	req, _ = http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(authPayload))
	res, err = http.DefaultClient.Do(req)
	data = GetBodyAsString(res)
	require.NoError(t, err, fmt.Sprintf("There shouldn't be an error from login: %s", data))

	require.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("Login: Expected status code 200. %s", data))

	var authResponse AuthResponse
	json.Unmarshal(data, &authResponse)

	return registeredUser, authResponse.Token
}
