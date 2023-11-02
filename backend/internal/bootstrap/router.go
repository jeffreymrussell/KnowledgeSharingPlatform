package bootstrap

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/users"
	"github.com/gorilla/mux"
)

func Router(config internal.Config) *mux.Router {
	// Initialize UserUsecase with dependencies

	dbAdapter := &users.SQLiteAdapter{
		DB: config.DB,
	}
	userUsecase := &users.UserUsecase{
		Repository: dbAdapter,
	}

	// Initialize UserHandler with dependencies
	userHandler := &users.UserHandler{
		Usecase: userUsecase,
	}
	r := mux.NewRouter()
	r.HandleFunc("/register", userHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", userHandler.LogoutHandler).Methods("POST")
	return r
}
