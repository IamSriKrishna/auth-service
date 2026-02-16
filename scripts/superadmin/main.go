package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	email := "admin@bbcloud.app"
	password := "admin123"

	// Hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	passwordHash := string(bytes)

	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Password Hash: %s\n", passwordHash)

	// Generate the SQL insert statement
	sql := fmt.Sprintf(`INSERT INTO users (email, username, password_hash, user_type, role_id, status, email_verified, created_at, updated_at) 
VALUES ('%s', 'superadmin', '%s', 'superadmin', 3, 'active', 1, NOW(), NOW());`, email, passwordHash)

	fmt.Printf("\nSQL Insert Statement:\n%s\n", sql)
}
