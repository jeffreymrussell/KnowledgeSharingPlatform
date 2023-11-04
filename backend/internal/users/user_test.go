package users_test

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/bootstrap"
	"KnowledgeSharingPlatform/internal/test"
	"KnowledgeSharingPlatform/internal/users"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUser represents a user object with a username and password.
type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	router := bootstrap.SetupRouter(bootstrap.SetupHandlers(bootstrap.SetupUseCases(bootstrap.SetupAdapters(db))))
	suite.Server = httptest.NewServer(router)
	if err != nil {
		return
	}

	suite.registeredUser, _ = test.RegisterUserAndLogin(suite.T(), suite.Server.URL)
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
		require.NoError(t, err, "There shouldn't be an error from login")

		data := test.GetBodyAsString(res)
		require.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("Login: Expected status code 200. %s", data))

		// Extract token from login response
		var authResponse test.AuthResponse
		json.Unmarshal(data, &authResponse)
		token := authResponse.Token

		// Logout
		logoutURL := fmt.Sprintf("%s/logout", suite.Server.URL)
		req, _ = http.NewRequest("POST", logoutURL, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err = http.DefaultClient.Do(req)
		require.NoError(t, err, "No error from logout")
		data = test.GetBodyAsString(res)
		require.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("Logout: Expected status code 200. %s", data))
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
		require.Equal(t, http.StatusUnauthorized, res.StatusCode, "Expected status code 401")
	})
}
