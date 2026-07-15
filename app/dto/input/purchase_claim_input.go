package input

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type CreatePurchaseClaimInput struct {
	PurchaseOrderID string                         `json:"purchase_order_id" validate:"required"`
	Date            time.Time                      `json:"date" validate:"required"`
	Notes           string                         `json:"notes"`
	Items           []CreatePurchaseClaimItemInput `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseClaimItemInput struct {
	PurchaseOrderItemID uint                       `json:"purchase_order_item_id" validate:"required"`
	Type                models.PurchaseClaimType   `json:"type" validate:"required,oneof=missing damaged"`
	Quantity            float64                    `json:"quantity" validate:"required,gt=0"`
	Unit                string                     `json:"unit" validate:"required"`
	Reason              string                     `json:"reason" validate:"required"`
	Action              models.PurchaseClaimAction `json:"action" validate:"required,oneof=replacement credit_note return_to_vendor scrap adjustment_only"`
}

type ReceivePurchaseClaimReplacementInput struct {
	PurchaseClaimItemID uint      `json:"purchase_claim_item_id" validate:"required"`
	Quantity            float64   `json:"quantity" validate:"required,gt=0"`
	Unit                string    `json:"unit" validate:"required"`
	ReceivedDate        time.Time `json:"received_date" validate:"required"`
	Notes               string    `json:"notes"`
}
