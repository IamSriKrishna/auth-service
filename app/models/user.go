package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserType string

const (
	UserTypeMobile     UserType = "mobile_user"
	UserTypeSuperAdmin UserType = "superadmin"
	UserTypeAdmin      UserType = "admin"
	UserTypePartner    UserType = "partner"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserType == UserTypeMobile {
		if u.Email == nil && u.Phone == nil && u.GoogleID == nil && u.AppleID == nil {
			return gorm.ErrInvalidData
		}
	}

	if u.UserType == UserTypeAdmin || u.UserType == UserTypePartner || u.UserType == UserTypeSuperAdmin {
		if u.Email == nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (UserSession) TableName() string {
	return "user_sessions"
}
