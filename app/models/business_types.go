package models

import (
	"time"

	"gorm.io/gorm"
)

type BusinessType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TypeName    string         `gorm:"size:100;not null;uniqueIndex" json:"type_name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BusinessType) TableName() string {
	return "business_types"
}
