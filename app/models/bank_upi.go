package models

import (
	"time"

	"gorm.io/gorm"
)

type CompanyBankDetail struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CompanyID         uint           `gorm:"not null;index" json:"company_id"`
	BankName          string         `gorm:"size:255;not null" json:"bank_name"`
	AccountHolderName string         `gorm:"size:255;not null" json:"account_holder_name"`
	AccountNumber     string         `gorm:"size:50;not null" json:"account_number"`
	IFSCCode          string         `gorm:"size:11;not null;index" json:"ifsc_code"`
	BranchName        string         `gorm:"size:255" json:"branch_name,omitempty"`
	IsPrimary         bool           `gorm:"default:false" json:"is_primary"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (CompanyBankDetail) TableName() string {
	return "company_bank_details"
}

type CompanyUPIDetail struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	UPIID     string    `gorm:"size:255;not null" json:"upi_id"`
	UPIQRURL  string    `gorm:"type:text" json:"upi_qr_url,omitempty"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (CompanyUPIDetail) TableName() string {
	return "company_upi_details"
}
