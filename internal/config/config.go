package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DBUrlMigration string
	SecreetJwt string

	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	// Try current directory first, then parent directory (for running from cmd/)
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	fmt.Println("Environment variables loaded successfully")
	return &Config{
		Port:           os.Getenv("PORT"),
		DBUrlMigration: os.Getenv("DATABASE_URL"),
		SecreetJwt:     os.Getenv("JWT_SECRET"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
	}, nil

}