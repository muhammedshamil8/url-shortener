package config

import (
	"os"
	"strings"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	Env    string
}

type JWTConfig struct {
	AccessTokenSecret  string
	AccessTokenExpiry  string
	RefreshTokenSecret string
	RefreshTokenExpiry string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port           string
	BaseURL        string
	AllowedOrigins []string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port:           os.Getenv("APP_PORT"),
			BaseURL:        os.Getenv("BASE_URL"),
			AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", ""), ","),
		},
		JWT: JWTConfig{
			AccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
			AccessTokenExpiry:  os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"),
			RefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
			RefreshTokenExpiry: os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"),
		},
		Env: getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
