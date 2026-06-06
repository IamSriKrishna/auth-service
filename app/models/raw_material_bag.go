package models

import "time"

type RawMaterialBag struct {
	ID string `json:"id" gorm:"type:varchar(255);primaryKey"`

	PurchaseOrderID string `json:"purchase_order_id" gorm:"type:varchar(255);index;not null"`
	PurchaseOrderNo string `json:"purchase_order_no" gorm:"type:varchar(100)"`

	VendorID   uint   `json:"vendor_id" gorm:"index"`
	VendorName string `json:"vendor_name" gorm:"type:varchar(255)"`

	ProductID   string `json:"product_id" gorm:"type:varchar(255);index;not null"`
	ProductName string `json:"product_name" gorm:"type:varchar(255)"`

	BagNumber   int     `json:"bag_number"`
	ExpectedKg  float64 `json:"expected_kg"`
	ActualKg    float64 `json:"actual_kg"`
	RemainingKg float64 `json:"remaining_kg"`

	Status string `json:"status" gorm:"type:varchar(50);default:'available'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RawMaterialBag) TableName() string {
	return "raw_material_bags"
}