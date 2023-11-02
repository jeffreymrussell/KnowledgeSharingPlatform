package test

import (
	"KnowledgeSharingPlatform/internal"
	"database/sql"
	"os"
)

func InitializeDatabase(db *sql.DB) error {
	// Execute table creation queries
	tableCreationQueries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE TABLE IF NOT EXISTS articles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            author_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (author_id) REFERENCES users (id)
        );`,
		// ... Add other table creation queries here
	}

	for _, query := range tableCreationQueries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTable(dbConfig internal.DbConfig) {

	err := dbConfig.DB.Close()
	if err != nil {
		panic("Failed to close DB")
	}
	err = os.Remove(dbConfig.DbFilePath)
	if err != nil {
		panic("failed to remove db")
	}
}
