package models

import "time"

type PurchaseDispense struct {
	ID                  string `gorm:"type:varchar(255);primaryKey" json:"id"`
	PurchaseClaimID     string `gorm:"type:varchar(255);not null;index" json:"purchase_claim_id"`
	PurchaseClaimItemID uint   `gorm:"not null;index" json:"purchase_claim_item_id"`
	PurchaseOrderID     string `gorm:"type:varchar(255);not null;index" json:"purchase_order_id"`
	ProductID           string `gorm:"type:varchar(255);not null;index" json:"product_id"`
	ProductName         string `gorm:"type:varchar(255)" json:"product_name"`
	IsRawMaterial       bool   `gorm:"not null;default:false" json:"is_raw_material"`

	Quantity     float64 `gorm:"type:decimal(18,6);not null" json:"quantity"`
	Unit         string  `gorm:"type:varchar(50);not null" json:"unit"`
	BaseQuantity float64 `gorm:"type:decimal(18,6);not null" json:"base_quantity"`
	BaseUnit     string  `gorm:"type:varchar(50);not null" json:"base_unit"`

	DispenseDate time.Time `gorm:"not null;index" json:"dispense_date"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedBy    string    `gorm:"type:varchar(255);not null" json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

func (PurchaseDispense) TableName() string {
	return "purchase_dispenses"
}
