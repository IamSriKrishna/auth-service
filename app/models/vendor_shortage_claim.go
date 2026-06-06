package models

import "time"

type VendorShortageClaim struct {
	ID string `json:"id" gorm:"type:varchar(255);primaryKey"`

	PurchaseOrderID string `json:"purchase_order_id" gorm:"type:varchar(255);index"`
	PurchaseOrderNo string `json:"purchase_order_no"`

	VendorID   uint   `json:"vendor_id"`
	VendorName string `json:"vendor_name"`

	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`

	ExpectedKg    float64 `json:"expected_kg"`
	ReceivedKg    float64 `json:"received_kg"`
	ShortageKg    float64 `json:"shortage_kg"`
	ShortageGrams float64 `json:"shortage_grams"`

	RatePerKg       float64 `json:"rate_per_kg"`
	ClaimAmount     float64 `json:"claim_amount"`
	RecoveredAmount float64 `json:"recovered_amount"`

	Status string `json:"status" gorm:"type:varchar(50);default:'pending'"`
	Notes  string `json:"notes" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (VendorShortageClaim) TableName() string {
	return "vendor_shortage_claims"
}