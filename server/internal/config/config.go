package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Host     string
	Port     string
	BasePath string // Base path for all routes (e.g., "/app" or "/")
}

type DatabaseConfig struct {
	Driver   string // "sqlite" or "postgres"
	DSN      string
	MaxConns int
	MinConns int
}

type AuthConfig struct {
	JWTSecret            string
	AccessTokenExpiry    time.Duration
	RefreshExpiry        time.Duration
	AllowRegistration    bool
	SecureCookies        bool   // true for HTTPS (production)
	CookieDomain         string // cookie domain (optional)
	InitialAdminUsername string // username for initial admin (optional)
	InitialAdminPassword string // password for initial admin (optional)
}

func Load() *Config {
	basePath := getEnv("BASE_PATH", "/")
	// Normalize base path
	if basePath != "/" {
		basePath = "/" + strings.Trim(basePath, "/")
	}
	return &Config{
		Server: ServerConfig{
			Host:     getEnv("SERVER_HOST", "0.0.0.0"),
			Port:     getEnv("SERVER_PORT", "8080"),
			BasePath: basePath,
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			DSN:      getEnv("DB_DSN", "file:data/lovely_eye.db?cache=shared&mode=rwc"),
			MaxConns: getEnvInt("DB_MAX_CONNS", 10),
			MinConns: getEnvInt("DB_MIN_CONNS", 1),
		},
		Auth: AuthConfig{
			JWTSecret:            getJWTSecret(),
			AccessTokenExpiry:    time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 15)) * time.Minute,
			RefreshExpiry:        time.Duration(getEnvInt("JWT_REFRESH_DAYS", 7)) * 24 * time.Hour,
			AllowRegistration:    getEnvBool("ALLOW_REGISTRATION", false),
			SecureCookies:        getEnvBool("SECURE_COOKIES", false),
			CookieDomain:         getEnv("COOKIE_DOMAIN", ""),
			InitialAdminUsername: getEnv("INITIAL_ADMIN_USERNAME", ""),
			InitialAdminPassword: getEnv("INITIAL_ADMIN_PASSWORD", ""),
		},
	}
}

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		if len(secret) < 32 {
			log.Fatal("JWT_SECRET must be at least 32 characters")
		}
		return secret
	}
	// Generate random secret for development (tokens won't survive restarts)
	log.Println("WARNING: JWT_SECRET not set, generating random secret. Tokens will not survive server restarts.")
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate random JWT secret")
	}
	return hex.EncodeToString(bytes)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
