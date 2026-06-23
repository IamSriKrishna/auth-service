package models

import (
	"time"
)

// CustomerPricing represents customer-specific pricing for products
type CustomerPricing struct {
	ID            string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CustomerID    uint       `json:"customer_id" gorm:"type:bigint;not null;index"`
	CustomerName  string     `json:"customer_name" gorm:"type:varchar(255);not null"`
	ProductID     string     `json:"product_id" gorm:"type:varchar(36);index"`
	ProductName   string     `json:"product_name" gorm:"type:varchar(255)"`
	Rate          float64    `json:"rate" gorm:"type:decimal(10,2);not null"`
	Account       string     `json:"account" gorm:"type:varchar(100)"`
	Description   string     `json:"description" gorm:"type:longtext"`
	EffectiveFrom *time.Time `json:"effective_from" gorm:"type:datetime"`
	EffectiveTo   *time.Time `json:"effective_to" gorm:"type:datetime"`
	IsActive      bool       `json:"is_active" gorm:"type:boolean;default:true;index"`
	Notes         string     `json:"notes" gorm:"type:longtext"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
	CreatedBy     string     `json:"created_by" gorm:"type:varchar(36)"`
	UpdatedBy     string     `json:"updated_by" gorm:"type:varchar(36)"`
}

// TableName specifies the table name for CustomerPricing model
func (cp *CustomerPricing) TableName() string {
	return "customer_pricing"
}
