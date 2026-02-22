package models

import (
	"time"

	"gorm.io/gorm"
)

type VendorBankDetail struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	VendorID          uint           `gorm:"not null;index" json:"vendor_id"`
	BankID            uint           `gorm:"not null;index" json:"bank_id"`
	AccountHolderName string         `gorm:"type:varchar(255)" json:"account_holder_name"`
	AccountNumber     string         `gorm:"type:varchar(50);not null" json:"account_number"`
	IFSCCode          string         `gorm:"size:11" json:"ifsc_code,omitempty"`
	BranchName        string         `gorm:"size:255" json:"branch_name,omitempty"`
	IsPrimary         bool           `gorm:"default:false" json:"is_primary"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	Vendor Vendor `gorm:"foreignKey:VendorID" json:"-"`
	Bank   Bank   `gorm:"foreignKey:BankID" json:"bank,omitempty"`
}

func (VendorBankDetail) TableName() string {
	return "vendor_bank_details"
}
