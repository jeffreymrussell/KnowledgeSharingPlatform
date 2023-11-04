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
	JWTSecret   []byte
	DatabaseURL string
}

var GlobalConfig *Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	validConfig := true
	errors := ""
	if config.ServerPort == "" {
		validConfig = false
		errors += "Server port not specified, please specify.\n"
	}
	if string(config.JWTSecret) == "" {
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

	GlobalConfig = &config

	return nil
}
