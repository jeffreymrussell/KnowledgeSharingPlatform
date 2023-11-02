package main

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/bootstrap"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := bootstrap.SetupDatabase(config.DatabaseURL)
	dbConfig := internal.DbConfig{
		DB:         db,
		DbFilePath: config.DatabaseURL,
	}
	router := bootstrap.SetupRouter(dbConfig)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on :8080")
	err = srv.ListenAndServe()
	if err != nil {
		panic("Failed to start server")
	}

	err = db.Close()
	if err != nil {
		return
	}
}
