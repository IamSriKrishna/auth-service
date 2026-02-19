package input

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type CreatePurchaseOrderInput struct {
	VendorID            uint                         `json:"vendor_id" validate:"required"`
	DeliveryAddressType string                       `json:"delivery_address_type" validate:"required,oneof=organization customer"`
	DeliveryAddressID   *uint                        `json:"delivery_address_id"`
	OrganizationName    string                       `json:"organization_name"`
	OrganizationAddress string                       `json:"organization_address"`
	CustomerID          *uint                        `json:"customer_id"`
	ReferenceNo         string                       `json:"reference_no"`
	Date                time.Time                    `json:"date" validate:"required"`
	DeliveryDate        time.Time                    `json:"delivery_date" validate:"required"`
	PaymentTerms        string                       `json:"payment_terms" validate:"required"`
	ShipmentPreference  string                       `json:"shipment_preference"`
	LineItems           []PurchaseOrderLineItemInput `json:"line_items" validate:"required,min=1,dive"`
	Discount            float64                      `json:"discount" validate:"gte=0"`
	DiscountType        string                       `json:"discount_type" validate:"omitempty,oneof=percentage amount"`
	TaxType             *string                      `json:"tax_type"`
	TaxID               *uint                        `json:"tax_id"`
	Adjustment          float64                      `json:"adjustment" validate:"gte=0"`
	Notes               string                       `json:"notes"`
	TermsAndConditions  string                       `json:"terms_and_conditions"`
	Attachments         []string                     `json:"attachments"`
}

type PurchaseOrderLineItemInput struct {
	ItemID         string            `json:"item_id" validate:"required"`
	VariantSKU     *string           `json:"variant_sku"`
	Account        string            `json:"account" validate:"required"`
	Quantity       float64           `json:"quantity" validate:"required,gt=0"`
	Rate           float64           `json:"rate" validate:"required,gt=0"`
	VariantDetails map[string]string `json:"variant_details"`
}

type UpdatePurchaseOrderInput struct {
	VendorID            *uint                        `json:"vendor_id"`
	DeliveryAddressType *string                      `json:"delivery_address_type" validate:"omitempty,oneof=organization customer"`
	DeliveryAddressID   *uint                        `json:"delivery_address_id"`
	OrganizationName    *string                      `json:"organization_name"`
	OrganizationAddress *string                      `json:"organization_address"`
	CustomerID          *uint                        `json:"customer_id"`
	ReferenceNo         *string                      `json:"reference_no"`
	Date                *time.Time                   `json:"date"`
	DeliveryDate        *time.Time                   `json:"delivery_date"`
	PaymentTerms        *string                      `json:"payment_terms"`
	ShipmentPreference  *string                      `json:"shipment_preference"`
	LineItems           []PurchaseOrderLineItemInput `json:"line_items" validate:"omitempty,dive"`
	Discount            *float64                     `json:"discount" validate:"omitempty,gte=0"`
	DiscountType        *string                      `json:"discount_type" validate:"omitempty,oneof=percentage amount"`
	TaxType             *string                      `json:"tax_type"`
	TaxID               *uint                        `json:"tax_id"`
	Adjustment          *float64                     `json:"adjustment" validate:"omitempty,gte=0"`
	Notes               *string                      `json:"notes"`
	TermsAndConditions  *string                      `json:"terms_and_conditions"`
	Attachments         []string                     `json:"attachments"`
}

type UpdatePurchaseOrderStatusInput struct {
	Status domain.PurchaseOrderStatus `json:"status" validate:"required,oneof=draft sent partially_received received cancelled"`
}
