package internal

import (
	"database/sql"
	"github.com/spf13/viper"
)

type DbConfig struct {
	DB         *sql.DB
	DbFilePath string
}

type Config struct {
	ServerPort  string
	JWTSecret   string
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(&config)
	validConfig := true
	errors := ""
	if config.ServerPort == "" {
		validConfig = false
		errors += "Server port not specified, please specify.\n"
	}
	if config.JWTSecret == "" {
		validConfig = false
		errors += "JWT Secret not specified, please specify.\n"
	}
	if config.DatabaseURL == "" {
		validConfig = false
		errors += "Database URL not specified, please specify.\n"
	}
	if !validConfig {
		panic(errors)
	}

	return &config, err
}
