package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	UserName string
	Password string
	Database string
}

func GetEnvFromFile() *DBConfig {
	_ = godotenv.Load("../.env")

	return &DBConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		UserName: getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		Database: getEnv("POSTGRES_DB", "events_db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
