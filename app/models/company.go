package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CompanyName    string         `gorm:"size:255;not null" json:"company_name"`
	BusinessTypeID uint           `gorm:"not null;index" json:"business_type_id"`
	GSTNumber      string         `gorm:"size:15;uniqueIndex" json:"gst_number,omitempty"`
	PANNumber      string         `gorm:"size:10;uniqueIndex" json:"pan_number,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy      *uint          `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy      *uint          `gorm:"index" json:"updated_by,omitempty"`

	BusinessType     BusinessType            `gorm:"foreignKey:BusinessTypeID" json:"business_type,omitempty"`
	Contact          *CompanyContact         `gorm:"foreignKey:CompanyID" json:"contact,omitempty"`
	Address          *CompanyAddress         `gorm:"foreignKey:CompanyID" json:"address,omitempty"`
	BankDetails      []CompanyBankDetail     `gorm:"foreignKey:CompanyID" json:"bank_details,omitempty"`
	UPIDetails       *CompanyUPIDetail       `gorm:"foreignKey:CompanyID" json:"upi_details,omitempty"`
	InvoiceSettings  *CompanyInvoiceSetting  `gorm:"foreignKey:CompanyID" json:"invoice_settings,omitempty"`
	TaxSettings      *CompanyTaxSetting      `gorm:"foreignKey:CompanyID" json:"tax_settings,omitempty"`
	RegionalSettings *CompanyRegionalSetting `gorm:"foreignKey:CompanyID" json:"regional_settings,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}
