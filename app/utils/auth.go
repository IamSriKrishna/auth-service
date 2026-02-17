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

var jwtSecretKey = []byte("your-secret-key-change-this-in-production")

func GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
		"exp":           time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iss":           "github.com/bbapp-org/auth-service",
		"sub":           fmt.Sprintf("%d", claims.UserID),
	})

	return token.SignedString(jwtSecretKey)
}

func GenerateRefreshToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24 * 90).Unix(),
		"iss":     "github.com/bbapp-org/auth-service",
		"sub":     fmt.Sprintf("%d", userID),
	})

	return token.SignedString(jwtSecretKey)
}

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

func IsValidEmail(email string) bool {
	return len(email) > 0 && len(email) <= 254
}

func IsValidPhone(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 15
}

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
