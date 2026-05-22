package models

import (
	"time"
)

// ProductConversion represents a conversion rule from raw material to finished product
type ProductConversion struct {
	ID                   string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	RawProductID         string    `json:"raw_product_id" gorm:"type:varchar(255);not null;index"`
	RawProductName       string    `json:"raw_product_name" gorm:"type:varchar(255)"`
	RawProductSpec       string    `json:"raw_product_spec" gorm:"type:varchar(255)"` // e.g., "500 ml"
	FinishedProductID    string    `json:"finished_product_id" gorm:"type:varchar(255);not null;index"`
	FinishedProductName  string    `json:"finished_product_name" gorm:"type:varchar(255)"`
	FinishedProductSpec  string    `json:"finished_product_spec" gorm:"type:varchar(255)"`      // e.g., "300 ml"
	FinishedVariantSKU   string    `json:"finished_variant_sku" gorm:"type:varchar(255)"`       // Optional: If specified, add stock to this variant
	ConversionRatio      float64   `json:"conversion_ratio" gorm:"type:decimal(10,4);not null"` // How many units of raw material make 1 unit of finished product
	LossPercentage       float64   `json:"loss_percentage" gorm:"type:decimal(5,2);default:0"`  // Loss during conversion (%)
	IsActive             bool      `json:"is_active" gorm:"default:true;index"`
	Notes                string    `json:"notes" gorm:"type:text"`
	CreatedBy            string    `json:"created_by" gorm:"type:varchar(255)"`
	CreatedByUserName    string    `json:"created_by_user_name" gorm:"type:varchar(255)"`
	CreatedByCompanyID   uint      `json:"created_by_company_id"`
	CreatedByCompanyName string    `json:"created_by_company_name" gorm:"type:varchar(255)"`
	UpdatedBy            string    `json:"updated_by" gorm:"type:varchar(255)"`
	UpdatedByUserName    string    `json:"updated_by_user_name" gorm:"type:varchar(255)"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ProductConversionRecord represents a single conversion transaction
type ProductConversionRecord struct {
	ID                       string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	ConversionID             string    `json:"conversion_id" gorm:"type:varchar(255);not null;index"`
	RawProductID             string    `json:"raw_product_id" gorm:"type:varchar(255);not null;index"`
	RawProductName           string    `json:"raw_product_name" gorm:"type:varchar(255)"`
	RawQuantityUsed          float64   `json:"raw_quantity_used" gorm:"type:decimal(18,4);not null"`
	FinishedProductID        string    `json:"finished_product_id" gorm:"type:varchar(255);not null;index"`
	FinishedProductName      string    `json:"finished_product_name" gorm:"type:varchar(255)"`
	FinishedVariantSKU       string    `json:"finished_variant_sku" gorm:"type:varchar(255)"` // Optional: Which variant received the stock (if product has variants)
	FinishedQuantityProduced float64   `json:"finished_quantity_produced" gorm:"type:decimal(18,4);not null"`
	LossQuantity             float64   `json:"loss_quantity" gorm:"type:decimal(18,4);default:0"`
	ConversionDate           time.Time `json:"conversion_date"`
	Status                   string    `json:"status" gorm:"type:varchar(50);default:'COMPLETED'"` // PENDING, COMPLETED, FAILED
	Notes                    string    `json:"notes" gorm:"type:text"`
	CreatedBy                string    `json:"created_by" gorm:"type:varchar(255)"`
	CreatedByUserName        string    `json:"created_by_user_name" gorm:"type:varchar(255)"`
	CreatedByCompanyID       uint      `json:"created_by_company_id"`
	CreatedByCompanyName     string    `json:"created_by_company_name" gorm:"type:varchar(255)"`
	CreatedAt                time.Time `json:"created_at"`
}

// TableName specifies the table name for ProductConversion
func (ProductConversion) TableName() string {
	return "product_conversions"
}

// TableName specifies the table name for ProductConversionRecord
func (ProductConversionRecord) TableName() string {
	return "product_conversion_records"
}
