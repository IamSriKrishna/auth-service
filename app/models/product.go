package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Product represents a product with variant support
type Product struct {
	ID   string `json:"id" gorm:"type:varchar(255);primaryKey"`
	Name string `json:"name" gorm:"not null"`

	ProductDetails ProductDetails `json:"product_details" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	SalesInfo      SalesInfo      `json:"sales_info" gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	PurchaseInfo   PurchaseInfo   `json:"purchase_info" gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	Inventory      Inventory      `json:"inventory" gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	ReturnPolicy   ReturnPolicy   `json:"return_policy" gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`

	CreatedBy            string    `json:"created_by" gorm:"type:varchar(255)"`
	CreatedByUserName    string    `json:"created_by_user_name" gorm:"type:varchar(255)"`
	CreatedByCompanyID   uint      `json:"created_by_company_id"`
	CreatedByCompanyName string    `json:"created_by_company_name" gorm:"type:varchar(255)"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

// ProductDetails contains detailed information about a product with its variants
type ProductDetails struct {
	ID             uint          `gorm:"primaryKey;autoIncrement"`
	ProductID      string        `gorm:"type:varchar(255);uniqueIndex"`
	Unit           string        `json:"unit" gorm:"type:varchar(50);not null"`
	BaseSKU        string        `json:"base_sku,omitempty" gorm:"type:varchar(255)"`
	UPC            string        `json:"upc,omitempty" gorm:"type:varchar(100)"`
	EAN            string        `json:"ean,omitempty" gorm:"type:varchar(100)"`
	MPN            string        `json:"mpn,omitempty" gorm:"type:varchar(100)"`
	ISBN           string        `json:"isbn,omitempty" gorm:"type:varchar(20)"`
	Description    string        `json:"description,omitempty" gorm:"type:text"`
	ManufacturerID *uint         `json:"manufacturer_id,omitempty" gorm:"index"`
	Manufacturer   *Manufacturer `json:"manufacturer,omitempty" gorm:"foreignKey:ManufacturerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	AttributeDefinitions ProductAttributeDefinitions `json:"attribute_definitions,omitempty" gorm:"type:json"`

	ProductVariants []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductDetailsID;constraint:OnDelete:CASCADE"`
}

func (ProductDetails) TableName() string {
	return "product_details"
}

// ProductAttributeDefinition defines possible attributes for product variants
type ProductAttributeDefinition struct {
	Key     string   `json:"key"`
	Options []string `json:"options"`
}

type ProductAttributeDefinitions []ProductAttributeDefinition

func (a ProductAttributeDefinitions) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *ProductAttributeDefinitions) Scan(value interface{}) error {
	if value == nil {
		*a = []ProductAttributeDefinition{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal ProductAttributeDefinitions value")
	}
	return json.Unmarshal(bytes, a)
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	ID               uint                      `gorm:"primaryKey;autoIncrement"`
	ProductDetailsID uint                      `gorm:"index;not null"`
	SKU              string                    `json:"sku" gorm:"type:varchar(255);not null;uniqueIndex"`
	VariantName      string                    `json:"variant_name,omitempty" gorm:"type:varchar(255)"`
	Attributes       []ProductVariantAttribute `json:"attributes" gorm:"foreignKey:ProductVariantID;constraint:OnDelete:CASCADE"`
	SellingPrice     float64                   `json:"selling_price" gorm:"not null"`
	CostPrice        float64                   `json:"cost_price" gorm:"not null"`
	StockQuantity    float64                   `json:"stock_quantity" gorm:"type:decimal(18,2);default:0"`
	ReorderLevel     float64                   `json:"reorder_level" gorm:"type:decimal(18,2);default:0"`
	IsActive         bool                      `json:"is_active" gorm:"default:true"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

// ProductVariantAttribute represents an attribute value for a product variant
type ProductVariantAttribute struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	ProductVariantID uint   `gorm:"index;not null"`
	Key              string `json:"key" gorm:"type:varchar(100);not null"`
	Value            string `json:"value" gorm:"type:varchar(255);not null"`
}

func (ProductVariantAttribute) TableName() string {
	return "product_variant_attributes"
}

// SalesInfo for Product (foreign key will be updated in migration)
// Using same SalesInfo model but with ProductID instead of ItemID

// PurchaseInfo for Product (foreign key will be updated in migration)
// Using same PurchaseInfo model but with ProductID instead of ItemID

// Inventory for Product (foreign key will be updated in migration)
// Using same Inventory model but with ProductID instead of ItemID

// ReturnPolicy for Product (foreign key will be updated in migration)
// Using same ReturnPolicy model but with ProductID instead of ItemID

// ProductOpeningStock tracks opening stock for products
type ProductOpeningStock struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	ProductID               string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	OpeningStock            float64   `json:"opening_stock" gorm:"default:0"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit" gorm:"default:0"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (ProductOpeningStock) TableName() string {
	return "product_opening_stock"
}

// ProductVariantOpeningStock tracks opening stock for product variants
type ProductVariantOpeningStock struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	ProductVariantSKU       string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	OpeningStock            float64   `json:"opening_stock" gorm:"default:0"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit" gorm:"default:0"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (ProductVariantOpeningStock) TableName() string {
	return "product_variant_opening_stock"
}
