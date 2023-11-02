package users_test

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/bootstrap"
	"KnowledgeSharingPlatform/internal/test"
	"KnowledgeSharingPlatform/internal/users"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUser represents a user object with a username and password.
type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response object containing the authentication token.
type AuthResponse struct {
	Token string `json:"token"`
}

type UserTestSuite struct {
	suite.Suite
	Server         *httptest.Server
	config         internal.DbConfig
	registeredUser users.LoginUserDTO
}

func (suite *UserTestSuite) SetupTest() {
	db := bootstrap.SetupDatabase("test_db.sqlite")
	err := test.InitializeDatabase(db)
	if err != nil {
		panic("Failed to open sqlite")
	}

	suite.config = internal.DbConfig{
		DB:         db,
		DbFilePath: "test_db.sqlite",
	}
	router := bootstrap.SetupRouter(suite.config)
	suite.Server = httptest.NewServer(router)
	if err != nil {
		return
	}

	suite.registeredUser = users.LoginUserDTO{
		Username: "registereduser",
		Password: "Registeredpassword1!",
	}
	payload := "{\"username\": \"registereduser\", \"password\": \"Registeredpassword1!\", \"email\": \"test@example.com\"}"

	registerURL := fmt.Sprintf("%s/register", suite.Server.URL)

	req, _ := http.NewRequest(http.MethodPost, registerURL, bytes.NewBuffer([]byte(payload)))
	res, err := http.DefaultClient.Do(req)
	data := getBodyAsString(res)
	assert.NoError(suite.T(), err, "Shouldn't have an error running the registration")
	assert.Equal(suite.T(), 201, res.StatusCode, fmt.Sprintf("201 for user registration: %s", data))
}

func (suite *UserTestSuite) AfterTest(_ string, _ string) {
	test.DeleteTable(suite.config)
}
func Test_UserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{})

}
func (suite *UserTestSuite) TestUserAuthentication() {
	suite.T().Run("LoginUser and LogoutUser", func(t *testing.T) {
		// Login
		payload, _ := json.Marshal(suite.registeredUser)
		loginURL := fmt.Sprintf("%s/login", suite.Server.URL)
		req, _ := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(payload))
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "There shouldn't be an error from login")

		data := getBodyAsString(res)
		require.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("Login: Expected status code 200. %s", data))

		// Extract token from login response
		var authResponse AuthResponse
		json.Unmarshal(data, &authResponse)
		token := authResponse.Token

		// Logout
		logoutURL := fmt.Sprintf("%s/logout", suite.Server.URL)
		req, _ = http.NewRequest("POST", logoutURL, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err, "No error from logout")
		data = getBodyAsString(res)
		assert.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("Logout: Expected status code 200. %s", data))
	})

	// LoginUserWithIncorrectCredentials Test
	suite.T().Run("LoginUserWithIncorrectCredentials", func(t *testing.T) {
		user := TestUser{
			Username: "wronguser",
			Password: "wrongpassword",
		}
		payload, _ := json.Marshal(user)
		loginURL := fmt.Sprintf("%s/login", suite.Server.URL)
		req, _ := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(payload))
		res, _ := http.DefaultClient.Do(req)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "Expected status code 401")
	})
}

func getBodyAsString(res *http.Response) []byte {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	data, _ := io.ReadAll(res.Body)
	return data
}
