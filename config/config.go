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
	Hostname   string
	HostnameIP string
	FrontendIP string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Hostname:   getEnv("HOSTNAME", "localhost"),
		HostnameIP: getEnv("HOSTNAME_IP", "192.168.0.1"),
		FrontendIP: getEnv("FRONTEND_IP", "192.168.0.2"),
	}
}
func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "defaultuser"),
		Password: getEnv("DB_PASSWORD", ""),
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
