package models

import "time"

// ConversionRecordBagUsage tracks which bags were used for each conversion record
type ConversionRecordBagUsage struct {
	ID                 string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	ConversionRecordID string    `json:"conversion_record_id" gorm:"type:varchar(255);not null;index"`
	BagID              string    `json:"bag_id" gorm:"type:varchar(255);not null;index"`
	BagNumber          int       `json:"bag_number"`
	ProductID          string    `json:"product_id" gorm:"type:varchar(255);not null"`
	ProductName        string    `json:"product_name" gorm:"type:varchar(255)"`
	QuantityUsedKg     float64   `json:"quantity_used_kg" gorm:"type:decimal(18,4);not null"`
	CreatedAt          time.Time `json:"created_at"`
}

// TableName specifies the table name for ConversionRecordBagUsage
func (ConversionRecordBagUsage) TableName() string {
	return "conversion_record_bag_usages"
}
