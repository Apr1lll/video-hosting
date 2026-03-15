package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	Env    string
}

type DBConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConnections int32
	IdleTimeout    time.Duration
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type JWTConfig struct {
	PrivateKeyPath  string
	PublicKeyPath   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

func Load() (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		// Не фатально, просто используем системные переменные
		fmt.Println("⚠️  .env file not found, using system environment variables")
	}

	cfg := &Config{
		Env: getEnv("APP_ENV", "development"),
		DB: DBConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnvAsInt("DB_PORT", 5432),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", "postgres"),
			Name:           getEnv("DB_NAME", "postgres"),
			SSLMode:        getEnv("DB_SSL_MODE", "disable"),
			MaxConnections: getEnvAsInt32("DB_MAX_CONNECTIONS", 20),
			IdleTimeout:    getEnvAsDuration("DB_IDLE_TIMEOUT", 5*time.Minute),
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		JWT: JWTConfig{
			PrivateKeyPath:  getEnv("JWT_PRIVATE_KEY_PATH", "keys/private.pem"),
			PublicKeyPath:   getEnv("JWT_PUBLIC_KEY_PATH", "keys/public.pem"),
			AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
			Issuer:          getEnv("JWT_ISSUER", "auth-service"),
		},
	}

	return cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsInt32(key string, defaultValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(intVal)
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
