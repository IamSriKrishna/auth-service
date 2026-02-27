package input

import (
	"time"
)

type CreatePackageInput struct {
	SalesOrderID  string                 `json:"sales_order_id" validate:"required"`
	CustomerID    uint                   `json:"customer_id" validate:"required"`
	PackageDate   time.Time              `json:"package_date" validate:"required"`
	PackageSlipNo *string                `json:"package_slip_no"` // Optional, will be auto-generated if not provided
	Items         []PackageLineItemInput `json:"items" validate:"omitempty,dive"`
	InternalNotes string                 `json:"internal_notes"`
}

// PackageLineItemInput represents items to be packed
// If items array is empty, all sales order items will be auto-populated
type PackageLineItemInput struct {
	SalesOrderItemID uint    `json:"sales_order_item_id" validate:"required"`
	PackedQty        float64 `json:"packed_qty" validate:"required,gte=0"`
}

type UpdatePackageInput struct {
	PackageDate   *time.Time             `json:"package_date"`
	Items         []PackageLineItemInput `json:"items" validate:"omitempty,dive"`
	InternalNotes *string                `json:"internal_notes"`
	Status        *string                `json:"status" validate:"omitempty,oneof=created packed shipped delivered cancelled"`
}

type UpdatePackageStatusInput struct {
	Status string `json:"status" validate:"required,oneof=created packed shipped delivered cancelled"`
}
