package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

// Shipment represents a shipment for delivery
type Shipment struct {
	ID              string                `json:"id" gorm:"type:varchar(255);primaryKey"`
	ShipmentNo      string                `json:"shipment_no" gorm:"column:shipment_no;type:varchar(100);uniqueIndex;not null"`
	PackageID       string                `json:"package_id" gorm:"type:varchar(255);not null;index"`
	Package         *Package              `json:"package,omitempty" gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SalesOrderID    string                `json:"sales_order_id" gorm:"type:varchar(255);not null;index"`
	SalesOrder      *SalesOrder           `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CustomerID      uint                  `json:"customer_id" gorm:"not null;index"`
	Customer        *Customer             `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ShipDate        time.Time             `json:"ship_date" gorm:"not null"`
	Carrier         string                `json:"carrier" gorm:"type:varchar(255)"`
	TrackingNo      string                `json:"tracking_no" gorm:"type:varchar(100)"`
	TrackingURL     string                `json:"tracking_url" gorm:"type:varchar(500)"`
	ShippingCharges float64               `json:"shipping_charges" gorm:"default:0"`
	Status          domain.ShipmentStatus `json:"status" gorm:"type:varchar(50);not null;default:'created'"`
	Notes           string                `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	CreatedBy       string                `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy       string                `json:"updated_by" gorm:"type:varchar(255)"`
}

// TableName specifies the table name for Shipment
func (Shipment) TableName() string {
	return "shipments"
}
