package input

import (
	"time"
)

type CreateSalesOrderInput struct {
	CustomerID           uint                      `json:"customer_id" validate:"required"`
	ReferenceNo          string                    `json:"reference_no"`
	SODate               time.Time                 `json:"sales_order_date" validate:"required"`
	ExpectedShipmentDate time.Time                 `json:"expected_shipment_date" validate:"required"`
	PaymentTerms         string                    `json:"payment_terms" validate:"required"`
	DeliveryMethod       string                    `json:"delivery_method"`
	SalespersonID        *uint                     `json:"salesperson_id"`
	LineItems            []SalesOrderLineItemInput `json:"line_items" validate:"required,min=1,dive"`
	ShippingCharges      float64                   `json:"shipping_charges" validate:"gte=0"`
	TaxType              *string                   `json:"tax_type"`
	TaxID                *uint                     `json:"tax_id"`
	Adjustment           float64                   `json:"adjustment" validate:"gte=0"`
	CustomerNotes        string                    `json:"customer_notes"`
	TermsAndConditions   string                    `json:"terms_and_conditions"`
	Attachments          []string                  `json:"attachments"`
}

type SalesOrderLineItemInput struct {
	ItemID         string            `json:"item_id" validate:"required"`
	VariantID      *uint             `json:"variant_id"`
	Quantity       float64           `json:"quantity" validate:"required,gt=0"`
	Rate           float64           `json:"rate" validate:"required,gt=0"`
	VariantDetails map[string]string `json:"variant_details"`
}

type UpdateSalesOrderInput struct {
	CustomerID           *uint                     `json:"customer_id"`
	ReferenceNo          *string                   `json:"reference_no"`
	SODate               *time.Time                `json:"sales_order_date"`
	ExpectedShipmentDate *time.Time                `json:"expected_shipment_date"`
	PaymentTerms         *string                   `json:"payment_terms"`
	DeliveryMethod       *string                   `json:"delivery_method"`
	SalespersonID        *uint                     `json:"salesperson_id"`
	LineItems            []SalesOrderLineItemInput `json:"line_items" validate:"omitempty,dive"`
	ShippingCharges      *float64                  `json:"shipping_charges" validate:"omitempty,gte=0"`
	TaxType              *string                   `json:"tax_type"`
	TaxID                *uint                     `json:"tax_id"`
	Adjustment           *float64                  `json:"adjustment" validate:"omitempty,gte=0"`
	CustomerNotes        *string                   `json:"customer_notes"`
	TermsAndConditions   *string                   `json:"terms_and_conditions"`
	Attachments          []string                  `json:"attachments"`
}

type UpdateSalesOrderStatusInput struct {
	Status string `json:"status" validate:"required,oneof=draft sent confirmed partial_shipped shipped delivered cancelled"`
}
