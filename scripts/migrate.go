package main

import (
	"log"
	"path/filepath"

	"github.com/bbapp-org/auth-service/app/config"
	"github.com/bbapp-org/auth-service/app/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(filepath.Join(".env")); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	config.InitDB()

	db := config.GetDB()

	// Auto-migrate schemas
	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
		&models.UserSession{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")
}
