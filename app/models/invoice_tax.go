package models

import (
	"time"
)

type TaxType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaxName     string    `gorm:"size:50;not null;uniqueIndex" json:"tax_name"`
	TaxCode     string    `gorm:"size:10;not null;uniqueIndex" json:"tax_code"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (TaxType) TableName() string {
	return "tax_types"
}

type CompanyInvoiceSetting struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	CompanyID            uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	InvoicePrefix        string    `gorm:"size:10;not null;default:'INV'" json:"invoice_prefix"`
	InvoiceStartNumber   int       `gorm:"not null;default:1" json:"invoice_start_number"`
	CurrentInvoiceNumber int       `gorm:"not null;default:1" json:"current_invoice_number"`
	ShowLogo             bool      `gorm:"default:true" json:"show_logo"`
	ShowSignature        bool      `gorm:"default:false" json:"show_signature"`
	RoundOffTotal        bool      `gorm:"default:true" json:"round_off_total"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (CompanyInvoiceSetting) TableName() string {
	return "company_invoice_settings"
}

type CompanyTaxSetting struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CompanyID  uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	GSTEnabled bool      `gorm:"default:true" json:"gst_enabled"`
	TaxTypeID  uint      `gorm:"not null;index" json:"tax_type_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
	TaxType TaxType `gorm:"foreignKey:TaxTypeID" json:"tax_type,omitempty"`
}

func (CompanyTaxSetting) TableName() string {
	return "company_tax_settings"
}
