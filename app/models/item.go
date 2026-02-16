package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type Item struct {
	ID   string          `json:"id" gorm:"type:varchar(255);primaryKey"`
	Name string          `json:"name" gorm:"not null"`
	Type domain.ItemType `json:"type" gorm:"not null"`

	ItemDetails  ItemDetails  `json:"item_details" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	SalesInfo    SalesInfo    `json:"sales_info" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	PurchaseInfo PurchaseInfo `json:"purchase_info" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Inventory    Inventory    `json:"inventory" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	ReturnPolicy ReturnPolicy `json:"return_policy" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Item) TableName() string {
	return "items"
}

// AttributeDefinition represents a variant attribute type and its possible values
type AttributeDefinition struct {
	Key     string   `json:"key"`
	Options []string `json:"options"`
}

// AttributeDefinitions is a custom type for storing attribute definitions as JSON
type AttributeDefinitions []AttributeDefinition

func (a AttributeDefinitions) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *AttributeDefinitions) Scan(value interface{}) error {
	if value == nil {
		*a = []AttributeDefinition{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal AttributeDefinitions value")
	}
	return json.Unmarshal(bytes, a)
}

type ItemDetails struct {
	ID          uint                 `gorm:"primaryKey;autoIncrement"`
	ItemID      string               `gorm:"type:varchar(255);uniqueIndex"`
	Structure   domain.ItemStructure `json:"structure" gorm:"not null"`
	Unit        string               `json:"unit" gorm:"type:varchar(50);not null"`
	SKU         string               `json:"sku,omitempty" gorm:"type:varchar(255)"`
	UPC         string               `json:"upc,omitempty" gorm:"type:varchar(100)"`
	EAN         string               `json:"ean,omitempty" gorm:"type:varchar(100)"`
	MPN         string               `json:"mpn,omitempty" gorm:"type:varchar(100)"`
	ISBN        string               `json:"isbn,omitempty" gorm:"type:varchar(20)"`
	Description string               `json:"description,omitempty" gorm:"type:text"`

	// Store attribute definitions as JSON for variant items
	AttributeDefinitions AttributeDefinitions `json:"attribute_definitions,omitempty" gorm:"type:json"`

	Variants []Variant `json:"variants,omitempty" gorm:"foreignKey:ItemDetailsID;constraint:OnDelete:CASCADE"`
}

func (ItemDetails) TableName() string {
	return "item_details"
}

type Variant struct {
	ID            uint               `gorm:"primaryKey;autoIncrement"`
	ItemDetailsID uint               `gorm:"index;not null"`
	SKU           string             `json:"sku" gorm:"type:varchar(255);not null;uniqueIndex"`
	Attributes    []VariantAttribute `json:"attributes" gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE"`
	SellingPrice  float64            `json:"selling_price" gorm:"not null"`
	CostPrice     float64            `json:"cost_price" gorm:"not null"`
	StockQuantity float64            `json:"stock_quantity" gorm:"type:decimal(18,2);default:0"` // Changed to decimal for precision
	ReorderLevel  float64            `json:"reorder_level" gorm:"type:decimal(18,2);default:0"`  // Auto-reorder when below this
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

func (Variant) TableName() string {
	return "variants"
}

type VariantAttribute struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	VariantID uint   `gorm:"index;not null"`
	Key       string `json:"key" gorm:"type:varchar(100);not null"`
	Value     string `json:"value" gorm:"type:varchar(255);not null"`
}

func (VariantAttribute) TableName() string {
	return "variant_attributes"
}

type SalesInfo struct {
	ID           uint    `json:"-" gorm:"primaryKey;autoIncrement"`
	ItemID       string  `json:"item_id" gorm:"type:varchar(255);uniqueIndex;not null"`
	Account      string  `json:"account" gorm:"not null"`
	SellingPrice float64 `json:"selling_price,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Description  string  `json:"description,omitempty" gorm:"type:text"`
}

func (SalesInfo) TableName() string {
	return "sales_info"
}

type PurchaseInfo struct {
	ID                uint    `json:"-" gorm:"primaryKey;autoIncrement"`
	ItemID            string  `json:"item_id" gorm:"type:varchar(255);uniqueIndex;not null"`
	Account           string  `json:"account" gorm:"not null"`
	CostPrice         float64 `json:"cost_price,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	PreferredVendorID *uint   `json:"preferred_vendor_id,omitempty" gorm:"index"`
	PreferredVendor   *Vendor `json:"preferred_vendor,omitempty" gorm:"foreignKey:PreferredVendorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Description       string  `json:"description,omitempty" gorm:"type:text"`
}

func (PurchaseInfo) TableName() string {
	return "purchase_info"
}

type Inventory struct {
	ID                       uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	ItemID                   string `json:"item_id" gorm:"type:varchar(255);uniqueIndex;not null"`
	TrackInventory           bool   `json:"track_inventory"`
	InventoryAccount         string `json:"inventory_account" gorm:"type:varchar(255)"`
	InventoryValuationMethod string `json:"inventory_valuation_method" gorm:"type:varchar(50)"` // FIFO, LIFO, Weighted Average
	ReorderPoint             int    `json:"reorder_point" gorm:"default:0"`
}

func (Inventory) TableName() string {
	return "inventory"
}

type ReturnPolicy struct {
	ID         uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	ItemID     string `json:"item_id" gorm:"type:varchar(255);uniqueIndex;not null"`
	Returnable bool   `json:"returnable"`
}

func (ReturnPolicy) TableName() string {
	return "return_policy"
}

type OpeningStock struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	ItemID                  string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	OpeningStock            float64   `json:"opening_stock" gorm:"default:0"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit" gorm:"default:0"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (OpeningStock) TableName() string {
	return "opening_stock"
}

type VariantOpeningStock struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	VariantID               uint      `gorm:"uniqueIndex;not null"`
	OpeningStock            float64   `json:"opening_stock" gorm:"default:0"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit" gorm:"default:0"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (VariantOpeningStock) TableName() string {
	return "variant_opening_stock"
}

type StockMovement struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	ItemID        string    `gorm:"type:varchar(255);index;not null"`
	VariantID     *uint     `gorm:"index"`
	MovementType  string    `gorm:"type:varchar(50);not null"` // purchase_received, sales_reserved, sales_invoiced, manufactured, consumed, adjustment
	Quantity      float64   `gorm:"type:decimal(18,2);not null"`
	RatePerUnit   float64   `gorm:"not null"`
	ReferenceType string    `gorm:"type:varchar(50)"` // PurchaseOrder, SalesOrder, Invoice, ProductionOrder
	ReferenceID   string    `gorm:"type:varchar(255);index"`
	ReferenceNo   string    `gorm:"type:varchar(100)"` // PO-001, SO-001, INV-001, etc.
	Notes         string    `gorm:"type:text"`
	Status        string    `gorm:"type:varchar(50);default:'pending'"` // pending, completed, reversed
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `gorm:"type:varchar(255)"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `gorm:"type:varchar(255)"`
}

func (StockMovement) TableName() string {
	return "stock_movements"
}
