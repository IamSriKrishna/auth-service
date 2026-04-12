package input

import (
	"time"
)

type CreateBillInput struct {
	VendorID        uint                `json:"vendor_id" validate:"required"`
	BillNumber      string              `json:"bill_number" validate:"required"`
	PurchaseOrderID *string             `json:"purchase_order_id"`
	BillingAddress  string              `json:"billing_address"`
	OrderNumber     string              `json:"order_number"`
	BillDate        time.Time           `json:"bill_date" validate:"required"`
	DueDate         time.Time           `json:"due_date" validate:"required"`
	PaymentTerms    string              `json:"payment_terms" validate:"required"`
	Subject         string              `json:"subject"`
	LineItems       []BillLineItemInput `json:"line_items" validate:"required,min=1,dive"`
	Discount        float64             `json:"discount" validate:"gte=0"`
	TaxType         *string             `json:"tax_type"`
	TaxPercentage   *float64            `json:"tax_percentage"`
	TaxID           *uint               `json:"tax_id"`
	Adjustment      float64             `json:"adjustment" validate:"gte=0"`
	Notes           string              `json:"notes"`
	Attachments     []string            `json:"attachments"`
}

type BillLineItemInput struct {
	ProductID   *string `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	SKU         string  `json:"sku"`
	Account     string  `json:"account" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required,gt=0"`
	Rate        float64 `json:"rate" validate:"required,gt=0"`
}

type UpdateBillInput struct {
	VendorID        *uint               `json:"vendor_id"`
	BillNumber      *string             `json:"bill_number"`
	PurchaseOrderID *string             `json:"purchase_order_id"`
	BillingAddress  *string             `json:"billing_address"`
	OrderNumber     *string             `json:"order_number"`
	BillDate        *time.Time          `json:"bill_date"`
	DueDate         *time.Time          `json:"due_date"`
	PaymentTerms    *string             `json:"payment_terms"`
	Subject         *string             `json:"subject"`
	LineItems       []BillLineItemInput `json:"line_items" validate:"omitempty,dive"`
	Discount        *float64            `json:"discount" validate:"omitempty,gte=0"`
	TaxType         *string             `json:"tax_type"`
	TaxID           *uint               `json:"tax_id"`
	Adjustment      *float64            `json:"adjustment" validate:"omitempty,gte=0"`
	Notes           *string             `json:"notes"`
	Attachments     []string            `json:"attachments"`
}

type UpdateBillStatusInput struct {
	Status string `json:"status" validate:"required,oneof=draft sent partial paid overdue void"`
}
