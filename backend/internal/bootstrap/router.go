package bootstrap

import (
	"KnowledgeSharingPlatform/internal/articles"
	"KnowledgeSharingPlatform/internal/middlewares"
	"KnowledgeSharingPlatform/internal/users"
	"database/sql"
	"github.com/gorilla/mux"
)

type Adapters struct {
	userRepositoryAdapter    users.UserRepository
	articleRepositoryAdapter articles.ArticleRepository
}

type Usecases struct {
	userUseCase    *users.UserUsecase
	articleUseCase *articles.ArticleUsecase
}

type Handlers struct {
	userHandler    *users.UserHandler
	articleHandler *articles.ArticleHandler
}

func SetupAdapters(db *sql.DB) Adapters {
	var adapters Adapters
	adapters.userRepositoryAdapter = &users.SQLiteAdapter{
		DB: db,
	}
	adapters.articleRepositoryAdapter = &articles.SQLiteArticleAdapter{
		DB: db,
	}
	return adapters
}
func SetupUseCases(adapters Adapters) Usecases {
	var usecases Usecases

	usecases.userUseCase = &users.UserUsecase{
		Repository: adapters.userRepositoryAdapter,
	}
	usecases.articleUseCase = &articles.ArticleUsecase{
		UserUsecase: usecases.userUseCase,
		ArticleRepo: adapters.articleRepositoryAdapter,
	}

	return usecases
}

func SetupHandlers(usecases Usecases) Handlers {
	var handlers Handlers

	// Initialize handlers with dependencies
	handlers.userHandler = &users.UserHandler{
		Usecase: usecases.userUseCase,
	}
	handlers.articleHandler = &articles.ArticleHandler{
		Usecase: usecases.articleUseCase,
	}
	return handlers
}
func SetupRouter(handlers Handlers) *mux.Router {

	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/register", handlers.userHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.userHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.userHandler.LogoutHandler).Methods("POST")

	// Article routes with authentication middleware
	s := r.PathPrefix("/articles").Subrouter()
	s.Use(middlewares.Authenticate)
	s.HandleFunc("", handlers.articleHandler.CreateArticle).Methods("POST")
	s.HandleFunc("", handlers.articleHandler.ListArticles).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", handlers.articleHandler.GetArticle).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", handlers.articleHandler.UpdateArticle).Methods("PUT")
	s.HandleFunc("/{id:[0-9]+}", handlers.articleHandler.DeleteArticle).Methods("DELETE")

	return r
}
