package input

import "time"

type CreatePurchaseDispenseInput struct {
	PurchaseClaimItemID uint      `json:"purchase_claim_item_id" validate:"required"`
	Quantity            float64   `json:"quantity" validate:"required,gt=0"`
	Unit                string    `json:"unit" validate:"required"`
	DispenseDate        time.Time `json:"dispense_date" validate:"required"`
	Notes               string    `json:"notes"`
}
