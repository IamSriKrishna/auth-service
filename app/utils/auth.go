package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/output"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key - should be loaded from environment in production
var jwtSecretKey = []byte("your-secret-key-change-this-in-production")

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token
func GenerateJWT(claims output.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":       claims.UserID,
		"user_type":     claims.UserType,
		"role":          claims.Role,
		"email":         claims.Email,
		"phone":         claims.Phone,
		"google_id":     claims.GoogleID,
		"identity_type": claims.IdentityType,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iss":           "github.com/bbapp-org/auth-service",
		"sub":           fmt.Sprintf("%d", claims.UserID),
	})

	return token.SignedString(jwtSecretKey)
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24 * 90).Unix(), // 90 days
		"iss":     "github.com/bbapp-org/auth-service",
		"sub":     fmt.Sprintf("%d", userID),
	})

	return token.SignedString(jwtSecretKey)
}

// ValidateJWT validates a JWT token and returns claims
func ValidateJWT(tokenString string) (*output.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid user_id in token")
		}

		userType, ok := claims["user_type"].(string)
		if !ok {
			return nil, errors.New("invalid user_type in token")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, errors.New("invalid role in token")
		}

		email, _ := claims["email"].(string)
		phone, _ := claims["phone"].(string)
		googleID, _ := claims["google_id"].(string)
		identityType, _ := claims["identity_type"].(string)

		return &output.Claims{
			UserID:       uint(userID),
			UserType:     userType,
			Role:         role,
			Email:        email,
			Phone:        phone,
			GoogleID:     googleID,
			IdentityType: identityType,
		}, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates a refresh token
func ValidateRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			return 0, errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid user_id in token")
		}

		return uint(userID), nil
	}

	return 0, errors.New("invalid refresh token")
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// IsValidEmail checks if email format is valid
func IsValidEmail(email string) bool {
	// Simple email validation - in production, use a proper email validation library
	return len(email) > 0 && len(email) <= 254 // Basic check
}

// IsValidPhone checks if phone format is valid
func IsValidPhone(phone string) bool {
	// Simple phone validation - in production, use a proper phone validation library
	return len(phone) >= 10 && len(phone) <= 15 // Basic check
}

// GetIdentityType determines the identity type based on the request
func GetIdentityType(email, phone, googleID string) string {
	if googleID != "" {
		return "google_oidc"
	}
	if phone != "" {
		return "phone"
	}
	if email != "" {
		return "email"
	}
	return ""
}
