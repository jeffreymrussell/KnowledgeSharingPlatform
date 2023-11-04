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
	err := internal.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := bootstrap.SetupDatabase(internal.GlobalConfig.DatabaseURL)

	router := bootstrap.SetupRouter(bootstrap.SetupHandlers(bootstrap.SetupUseCases(bootstrap.SetupAdapters(db))))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + internal.GlobalConfig.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on :" + internal.GlobalConfig.ServerPort)
	err = srv.ListenAndServe()
	if err != nil {
		panic("Failed to start server")
	}

	err = db.Close()
	if err != nil {
		return
	}
}
