package input

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type CreateInvoiceInput struct {
	CustomerID         uint                   `json:"customer_id" validate:"required"`
	OrderNumber        string                 `json:"order_number"`
	InvoiceDate        time.Time              `json:"invoice_date" validate:"required"`
	Terms              string                 `json:"terms" validate:"required"`
	DueDate            time.Time              `json:"due_date" validate:"required"`
	SalespersonID      *uint                  `json:"salesperson_id"`
	Subject            string                 `json:"subject"`
	LineItems          []InvoiceLineItemInput `json:"line_items" validate:"required,min=1,dive"`
	ShippingCharges    float64                `json:"shipping_charges" validate:"gte=0"`
	TaxType            *string                `json:"tax_type"`
	TaxID              *uint                  `json:"tax_id"`
	Adjustment         float64                `json:"adjustment"`
	CustomerNotes      string                 `json:"customer_notes"`
	TermsAndConditions string                 `json:"terms_and_conditions"`
	Attachments        []string               `json:"attachments"`
	PaymentReceived    bool                   `json:"payment_received"`
	PaymentSplits      []PaymentSplitInput    `json:"payment_splits"`
	EmailRecipients    []string               `json:"email_recipients"`
}

type InvoiceLineItemInput struct {
	ItemID         string            `json:"item_id" validate:"required"`
	VariantID      *uint             `json:"variant_id"`
	Description    string            `json:"description"`
	Quantity       float64           `json:"quantity" validate:"required,gt=0"`
	Rate           float64           `json:"rate" validate:"required,gt=0"`
	VariantDetails map[string]string `json:"variant_details"`
}

type UpdateInvoiceInput struct {
	CustomerID         *uint                  `json:"customer_id"`
	OrderNumber        *string                `json:"order_number"`
	InvoiceDate        *time.Time             `json:"invoice_date"`
	Terms              *string                `json:"terms"`
	DueDate            *time.Time             `json:"due_date"`
	SalespersonID      *uint                  `json:"salesperson_id"`
	Subject            *string                `json:"subject"`
	LineItems          []InvoiceLineItemInput `json:"line_items" validate:"omitempty,dive"`
	ShippingCharges    *float64               `json:"shipping_charges" validate:"omitempty,gte=0"`
	TaxType            *string                `json:"tax_type"`
	TaxID              *uint                  `json:"tax_id"`
	Adjustment         *float64               `json:"adjustment"`
	CustomerNotes      *string                `json:"customer_notes"`
	TermsAndConditions *string                `json:"terms_and_conditions"`
	Attachments        []string               `json:"attachments"`
	PaymentReceived    *bool                  `json:"payment_received"`
	PaymentSplits      []PaymentSplitInput    `json:"payment_splits"`
	EmailRecipients    []string               `json:"email_recipients"`
}

type CreateSalespersonInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateSalespersonInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email" validate:"omitempty,email"`
}

type CreateTaxInput struct {
	Name    string  `json:"name" validate:"required"`
	TaxType string  `json:"tax_type" validate:"required"`
	Rate    float64 `json:"rate" validate:"required,gte=0,lte=100"`
}

type UpdateTaxInput struct {
	Name    *string  `json:"name"`
	TaxType *string  `json:"tax_type"`
	Rate    *float64 `json:"rate" validate:"omitempty,gte=0,lte=100"`
}

type PaymentSplitInput struct {
	PaymentMode    string  `json:"payment_mode" validate:"required"`
	DepositTo      string  `json:"deposit_to"`
	AmountReceived float64 `json:"amount_received" validate:"gte=0"`
}

type CreatePaymentInput struct {
	InvoiceID   string    `json:"invoice_id" validate:"required"`
	PaymentDate time.Time `json:"payment_date" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	PaymentMode string    `json:"payment_mode" validate:"required"`
	Reference   string    `json:"reference"`
	Notes       string    `json:"notes"`
}

type InvoiceStatusUpdateInput struct {
	Status domain.InvoiceStatus `json:"status" validate:"required,oneof=draft sent partial paid overdue void"`
}
