package input

import "time"

type CreateVendorPaymentInput struct {
	PurchaseOrderID string    `json:"purchase_order_id" validate:"required"`
	VendorID        uint      `json:"vendor_id" validate:"required"`
	PaymentMode     string    `json:"payment_mode" validate:"required,oneof=cash online"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	PaymentDate     time.Time `json:"payment_date" validate:"required"`
	ReferenceNumber string    `json:"reference_number"`
	Notes           string    `json:"notes"`
}

type UpdateVendorPaymentInput struct {
	PaymentMode     *string    `json:"payment_mode" validate:"omitempty,oneof=cash online"`
	Amount          *float64   `json:"amount" validate:"omitempty,gt=0"`
	PaymentDate     *time.Time `json:"payment_date"`
	ReferenceNumber *string    `json:"reference_number"`
	Notes           *string    `json:"notes"`
}

type RecordVendorPaymentInput struct {
	PaidAmount      float64 `json:"paid_amount" validate:"required,gt=0"`
	PaymentMode     string  `json:"payment_mode" validate:"required,oneof=cash online"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
}
