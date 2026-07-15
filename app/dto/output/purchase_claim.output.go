package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type PurchaseOrderClaimSourceOutput struct {
	PurchaseOrderID     string                               `json:"purchase_order_id"`
	PurchaseOrderNumber string                               `json:"purchase_order_number"`
	VendorID            uint                                 `json:"vendor_id"`
	VendorName          string                               `json:"vendor_name"`
	Status              string                               `json:"status"`
	InventorySynced     bool                                 `json:"inventory_synced"`
	Items               []PurchaseOrderClaimSourceItemOutput `json:"items"`
}

type PurchaseOrderClaimSourceItemOutput struct {
	PurchaseOrderItemID uint   `json:"purchase_order_item_id"`
	ProductID           string `json:"product_id"`
	ProductName         string `json:"product_name"`
	SKU                 string `json:"sku"`
	IsRawMaterial       bool   `json:"is_raw_material"`

	OrderedQuantity     float64 `json:"ordered_quantity"`
	OrderedUnit         string  `json:"ordered_unit"`
	OrderedBaseQuantity float64 `json:"ordered_base_quantity"`
	BaseUnit            string  `json:"base_unit"`

	ReceivedQuantity     float64 `json:"received_quantity"`
	ReceivedBaseQuantity float64 `json:"received_base_quantity"`

	MissingReportedBase float64 `json:"missing_reported_base"`
	DamagedReportedBase float64 `json:"damaged_reported_base"`

	MissingRemainingBase float64 `json:"missing_remaining_base"`
	DamagedRemainingBase float64 `json:"damaged_remaining_base"`

	ReplacementPendingBase float64 `json:"replacement_pending_base"`

	NumberOfPacks   float64 `json:"number_of_packs,omitempty"`
	QuantityPerPack float64 `json:"quantity_per_pack,omitempty"`
	ReceivedPacks   float64 `json:"received_packs,omitempty"`
	Rate            float64 `json:"rate"`
}

type PurchaseClaimOutput struct {
	ID                  string                     `json:"id"`
	ClaimNumber         string                     `json:"claim_number"`
	PurchaseOrderID     string                     `json:"purchase_order_id"`
	PurchaseOrderNumber string                     `json:"purchase_order_number"`
	VendorID            uint                       `json:"vendor_id"`
	CompanyID           uint                       `json:"company_id"`
	Date                time.Time                  `json:"date"`
	Status              models.PurchaseClaimStatus `json:"status"`
	Notes               string                     `json:"notes"`
	CreatedBy           string                     `json:"created_by"`
	Items               []models.PurchaseClaimItem `json:"items"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}

func ToPurchaseClaimOutput(value *models.PurchaseClaim) *PurchaseClaimOutput {
	return &PurchaseClaimOutput{
		ID:                  value.ID,
		ClaimNumber:         value.ClaimNumber,
		PurchaseOrderID:     value.PurchaseOrderID,
		PurchaseOrderNumber: value.PurchaseOrderNumber,
		VendorID:            value.VendorID,
		CompanyID:           value.CompanyID,
		Date:                value.ClaimDate,
		Status:              value.Status,
		Notes:               value.Notes,
		CreatedBy:           value.CreatedBy,
		Items:               value.Items,
		CreatedAt:           value.CreatedAt,
		UpdatedAt:           value.UpdatedAt,
	}
}
