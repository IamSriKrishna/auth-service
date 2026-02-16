package input

import (
	"time"
)

type CreateShipmentInput struct {
	PackageID       string    `json:"package_id" validate:"required"`
	SalesOrderID    string    `json:"sales_order_id" validate:"required"`
	CustomerID      uint      `json:"customer_id" validate:"required"`
	ShipDate        time.Time `json:"ship_date" validate:"required"`
	Carrier         string    `json:"carrier"`
	TrackingNo      string    `json:"tracking_no"`
	TrackingURL     string    `json:"tracking_url"`
	ShippingCharges float64   `json:"shipping_charges" validate:"gte=0"`
	Notes           string    `json:"notes"`
}

type UpdateShipmentInput struct {
	ShipDate        *time.Time `json:"ship_date"`
	Carrier         *string    `json:"carrier"`
	TrackingNo      *string    `json:"tracking_no"`
	TrackingURL     *string    `json:"tracking_url"`
	ShippingCharges *float64   `json:"shipping_charges" validate:"omitempty,gte=0"`
	Notes           *string    `json:"notes"`
	Status          *string    `json:"status" validate:"omitempty,oneof=created shipped in_transit delivered cancelled"`
}

type UpdateShipmentStatusInput struct {
	Status string `json:"status" validate:"required,oneof=created shipped in_transit delivered cancelled"`
}
