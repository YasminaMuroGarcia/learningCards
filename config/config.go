package config

import (
	"os"
)

type DBConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	SSLMode  string
}

type AppConfig struct {
	Hostname string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Hostname: getEnv("HOSTNAME", "localhost"),
	}
}
func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "defaultuser"),
		Password: getEnv("DB_PASSWORD", "defaultpassword"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
