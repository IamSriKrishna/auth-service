package models

import (
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	BankName   string         `gorm:"size:255;not null;uniqueIndex" json:"bank_name"`
	Address    string         `gorm:"type:text" json:"address,omitempty"`
	City       string         `gorm:"size:100" json:"city,omitempty"`
	State      string         `gorm:"size:100" json:"state,omitempty"`
	PostalCode string         `gorm:"size:10" json:"postal_code,omitempty"`
	Country    string         `gorm:"size:100" json:"country,omitempty"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	CompanyBankDetails []CompanyBankDetail `gorm:"foreignKey:BankID" json:"-"`
	VendorBankDetails  []VendorBankDetail  `gorm:"foreignKey:BankID" json:"-"`
}

func (Bank) TableName() string {
	return "banks"
}
