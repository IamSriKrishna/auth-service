package input

import (
	"time"
)

type CreatePackageInput struct {
	SalesOrderID  string             `json:"sales_order_id" validate:"required"`
	CustomerID    uint               `json:"customer_id" validate:"required"`
	PackageDate   time.Time          `json:"package_date" validate:"required"`
	Items         []PackageItemInput `json:"items" validate:"required,min=1,dive"`
	InternalNotes string             `json:"internal_notes"`
}

type PackageItemInput struct {
	SalesOrderItemID uint              `json:"sales_order_item_id" validate:"required"`
	ItemID           string            `json:"item_id" validate:"required"`
	VariantID        *uint             `json:"variant_id"`
	OrderedQty       float64           `json:"ordered_qty" validate:"required,gt=0"`
	PackedQty        float64           `json:"packed_qty" validate:"required,gte=0"`
	VariantDetails   map[string]string `json:"variant_details"`
}

type UpdatePackageInput struct {
	PackageDate   *time.Time         `json:"package_date"`
	Items         []PackageItemInput `json:"items" validate:"omitempty,dive"`
	InternalNotes *string            `json:"internal_notes"`
	Status        *string            `json:"status" validate:"omitempty,oneof=created packed shipped delivered cancelled"`
}

type UpdatePackageStatusInput struct {
	Status string `json:"status" validate:"required,oneof=created packed shipped delivered cancelled"`
}
