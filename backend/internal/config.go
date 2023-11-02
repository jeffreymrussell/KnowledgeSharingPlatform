package internal

import (
	"database/sql"
)

type Config struct {
	DB         *sql.DB
	DbFilePath string
}
