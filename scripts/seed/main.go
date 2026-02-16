package main

import (
	"log"
	"path/filepath"

	"github.com/bbapp-org/auth-service/app/config"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/utils"

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

	// Seed roles
	roles := []models.Role{
		{
			RoleName:    "mobile_user",
			Permissions: []string{"mobile_app_access", "profile_read", "profile_update"},
			Description: "Standard mobile application user",
			IsActive:    true,
		},
		{
			RoleName:    "superadmin",
			Permissions: []string{"*"},
			Description: "Super administrator with full system access",
			IsActive:    true,
		},
		{
			RoleName:    "admin",
			Permissions: []string{"admin_panel_access", "user_management", "reports_read"},
			Description: "Administrative user with management capabilities",
			IsActive:    true,
		},
		{
			RoleName:    "partner",
			Permissions: []string{"partner_tools_access", "data_operations", "system_monitoring", "revenue_dashboard"},
			Description: "Business partner with operational access and revenue sharing",
			IsActive:    true,
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Error creating role %s: %v", role.RoleName, err)
			} else {
				log.Printf("Created role: %s", role.RoleName)
			}
		} else {
			log.Printf("Role %s already exists", role.RoleName)
		}
	}

	// Create super admin user
	var superAdminRole models.Role
	if err := db.Where("role_name = ?", "superadmin").First(&superAdminRole).Error; err != nil {
		log.Fatal("Super admin role not found")
	}

	var existingUser models.User
	superAdminEmail := "admin@example.com"
	if err := db.Where("email = ?", superAdminEmail).First(&existingUser).Error; err != nil {
		// Super admin doesn't exist, create it
		passwordHash, err := utils.HashPassword("admin123")
		if err != nil {
			log.Fatal("Error hashing password:", err)
		}

		superAdmin := models.User{
			Email:        &superAdminEmail,
			Username:     stringPtr("superadmin"),
			PasswordHash: &passwordHash,
			UserType:     models.UserTypeSuperAdmin,
			RoleID:       superAdminRole.ID,
			Status:       models.UserStatusActive,
		}

		if err := db.Create(&superAdmin).Error; err != nil {
			log.Fatal("Error creating super admin user:", err)
		} else {
			log.Printf("Created super admin user: %s", superAdminEmail)
		}
	} else {
		log.Printf("Super admin user %s already exists", superAdminEmail)
	}

	log.Println("Database seeding completed successfully")
}

func stringPtr(s string) *string {
	return &s
}
