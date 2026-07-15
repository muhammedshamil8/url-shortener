package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DB     DBConfig
	Redis  RedisConfig
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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
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
	raw := strings.Split(getEnv("ALLOWED_ORIGINS", ""), ",")

	origins := make([]string, 0, len(raw))
	for _, origin := range raw {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
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
			AllowedOrigins: origins,
		},
		JWT: JWTConfig{
			AccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
			AccessTokenExpiry:  os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"),
			RefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
			RefreshTokenExpiry: os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       getIntEnv("REDIS_DB", 0),
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

func getIntEnv(key string, fallback int) int {
	value := getEnv(key, strconv.Itoa(fallback))
	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return result
}
