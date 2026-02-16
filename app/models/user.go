package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UserType represents the type of user
type UserType string

const (
	UserTypeMobile     UserType = "mobile_user"
	UserTypeSuperAdmin UserType = "superadmin"
	UserTypeAdmin      UserType = "admin"
	UserTypePartner    UserType = "partner"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
)

// StringArray is a custom type for handling JSON arrays in database
type StringArray []string

// Value implements the driver.Valuer interface for database writes
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for database reads
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}

	return json.Unmarshal(bytes, a)
}

// User represents the users table
type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Email             *string        `gorm:"unique;index" json:"email,omitempty"`
	Phone             *string        `gorm:"unique;index" json:"phone,omitempty"`
	Username          *string        `gorm:"unique;index" json:"username,omitempty"`
	PasswordHash      *string        `json:"-"`
	GoogleID          *string        `gorm:"unique;index" json:"google_id,omitempty"`
	AppleID           *string        `gorm:"unique;index" json:"apple_id,omitempty"`
	UserType          UserType       `gorm:"type:enum('mobile_user','superadmin','admin','partner');not null" json:"user_type"`
	RoleID            uint           `gorm:"not null;index" json:"role_id"`
	Role              Role           `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	EmailVerified     bool           `gorm:"default:false" json:"email_verified"`
	PhoneVerified     bool           `gorm:"default:false" json:"phone_verified"`
	Status            UserStatus     `gorm:"type:enum('active','inactive','pending');default:'active'" json:"status"`
	CreatedBy         *uint          `gorm:"index" json:"created_by,omitempty"`
	CreatedByUser     *User          `gorm:"foreignKey:CreatedBy;references:ID" json:"created_by_user,omitempty"`
	LastLoginAt       *time.Time     `json:"last_login_at,omitempty"`
	PasswordChangedAt *time.Time     `json:"password_changed_at,omitempty"`
}

// Role represents the roles table
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	RoleName    string         `gorm:"unique;not null" json:"role_name"`
	Permissions StringArray    `gorm:"type:json" json:"permissions"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
}

// RefreshToken represents the refresh tokens table
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TokenID   string    `gorm:"unique;not null;index" json:"token_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
}

// UserSession represents active user sessions
type UserSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	SessionID string    `gorm:"unique;not null;index" json:"session_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Validate that mobile users have at least one identity
	if u.UserType == UserTypeMobile {
		if u.Email == nil && u.Phone == nil && u.GoogleID == nil && u.AppleID == nil {
			return gorm.ErrInvalidData
		}
	}

	// Validate that admin/partner users have email
	if u.UserType == UserTypeAdmin || u.UserType == UserTypePartner || u.UserType == UserTypeSuperAdmin {
		if u.Email == nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// TableName returns the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// TableName returns the table name for RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// TableName returns the table name for UserSession model
func (UserSession) TableName() string {
	return "user_sessions"
}
