package models

import (
	"time"
)

// CustomerPricing represents customer-specific pricing for products
type CustomerPricing struct {
	ID string `json:"id" gorm:"primaryKey;type:varchar(36)"`

	// Company
	CompanyID uint `json:"company_id" gorm:"not null;index"`

	// Customer
	CustomerID   uint   `json:"customer_id" gorm:"not null;index"`
	CustomerName string `json:"customer_name" gorm:"type:varchar(255);not null"`

	// Product
	ProductID   string `json:"product_id" gorm:"type:varchar(36);index"`
	ProductName string `json:"product_name" gorm:"type:varchar(255)"`

	// Pricing
	Rate        float64 `json:"rate" gorm:"type:decimal(10,2);not null"`
	Account     string  `json:"account" gorm:"type:varchar(100)"`
	Description string  `json:"description" gorm:"type:longtext"`
	Notes       string  `json:"notes" gorm:"type:longtext"`

	EffectiveFrom *time.Time `json:"effective_from" gorm:"type:datetime"`
	EffectiveTo   *time.Time `json:"effective_to" gorm:"type:datetime"`

	IsActive bool `json:"is_active" gorm:"default:true;index"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	// Audit
	CreatedBy string `json:"created_by" gorm:"type:varchar(36)"`
	UpdatedBy string `json:"updated_by" gorm:"type:varchar(36)"`
}

func (CustomerPricing) TableName() string {
	return "customer_pricing"
}
