package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/joho/godotenv"
)

type Config struct {
	Service  input.ServiceConfig
	Database input.DatabaseConfig
	Server   input.ServerConfig
	App      input.AppConfig
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load, using environment variables from ConfigMap/Secret")
	}

	config := &Config{
		Database: input.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "auth_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: input.ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8088),
		},
		App: input.AppConfig{
			Environment:    getEnv("ENV", "development"),
			JWTSecret:      getEnv("JWT_SECRET", ""),
			ServerPort:     getEnv("SERVER_PORT", "8088"),
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
		
	}

	// Validate required database environment variables
	log.Printf("MySQL database configuration loaded:")
	log.Printf("  DB_HOST: %s", config.Database.Host)
	log.Printf("  DB_PORT: %d", config.Database.Port)
	log.Printf("  DB_USER: %s", config.Database.User)
	log.Printf("  DB_NAME: %s", config.Database.DBName)
	log.Printf("  PASSWORD_SET: %t", config.Database.Password != "")

	return config
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
}

func (c *Config) GetReadReplicaDSN() string {
	// If read replica is not configured, return empty string
	if c.Database.ReadReplicaHost == "" {
		return ""
	}
	// Use read replica credentials if provided, otherwise fall back to primary credentials
	user := c.Database.ReadReplicaUser
	password := c.Database.ReadReplicaPassword
	if user == "" {
		user = c.Database.User
	}
	if password == "" {
		password = c.Database.Password
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=preferred",
		user,
		password,
		c.Database.ReadReplicaHost,
		c.Database.ReadReplicaPort,
		c.Database.DBName,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func getEnv(key, defaultValue string) string {
	// In Kubernetes, environment variables are injected from ConfigMaps and Secrets
	// Use os.Getenv to read them directly
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := getEnv(key, ""); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
