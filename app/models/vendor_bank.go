package models

import (
	"time"

	"gorm.io/gorm"
)

type VendorBankDetail struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	VendorID          uint           `gorm:"not null;index" json:"vendor_id"`
	AccountHolderName string         `gorm:"type:varchar(255)" json:"account_holder_name"`
	BankName          string         `gorm:"type:varchar(255)" json:"bank_name"`
	AccountNumber     string         `gorm:"type:varchar(50);not null" json:"account_number"`
	IFSC              string         `gorm:"type:varchar(11);not null" json:"ifsc"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (VendorBankDetail) TableName() string {
	return "vendor_bank_details"
}
