package main

import (
	"KnowledgeSharingPlatform/internal"
	"KnowledgeSharingPlatform/internal/bootstrap"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic("Failed to open sqlite")
	}
	config := internal.Config{
		DB:         db,
		DbFilePath: "db.sqlite",
	}
	r := bootstrap.Router(config)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
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
