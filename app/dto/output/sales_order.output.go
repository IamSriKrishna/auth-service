package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type SalesOrderOutput struct {
	ID                   string                     `json:"id"`
	SalesOrderNo         string                     `json:"sales_order_no"`
	CustomerID           uint                       `json:"customer_id"`
	Customer             *CustomerInfo              `json:"customer,omitempty"`
	SalespersonID        *uint                      `json:"salesperson_id,omitempty"`
	Salesperson          *SalespersonInfo           `json:"salesperson,omitempty"`
	ReferenceNo          string                     `json:"reference_no,omitempty"`
	SODate               time.Time                  `json:"sales_order_date"`
	ExpectedShipmentDate time.Time                  `json:"expected_shipment_date"`
	PaymentTerms         string                     `json:"payment_terms"`
	DeliveryMethod       string                     `json:"delivery_method,omitempty"`
	LineItems            []SalesOrderLineItemOutput `json:"line_items"`
	SubTotal             float64                    `json:"sub_total"`
	ShippingCharges      float64                    `json:"shipping_charges"`
	TaxType              *string                    `json:"tax_type,omitempty"`
	TaxID                *uint                      `json:"tax_id,omitempty"`
	Tax                  *TaxInfo                   `json:"tax,omitempty"`
	TaxAmount            float64                    `json:"tax_amount"`
	Adjustment           float64                    `json:"adjustment"`
	Total                float64                    `json:"total"`
	CustomerNotes        string                     `json:"customer_notes,omitempty"`
	TermsAndConditions   string                     `json:"terms_and_conditions,omitempty"`
	Status               string                     `json:"status"`
	Attachments          []string                   `json:"attachments,omitempty"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
	CreatedBy            string                     `json:"created_by,omitempty"`
	UpdatedBy            string                     `json:"updated_by,omitempty"`
}

type SalesOrderLineItemOutput struct {
	ID             uint              `json:"id"`
	ItemID         string            `json:"item_id"`
	Item           *ItemInfo         `json:"item,omitempty"`
	VariantID      *uint             `json:"variant_id,omitempty"`
	Variant        *VariantInfo      `json:"variant,omitempty"`
	Quantity       float64           `json:"quantity"`
	Rate           float64           `json:"rate"`
	Amount         float64           `json:"amount"`
	VariantDetails map[string]string `json:"variant_details,omitempty"`
}

func ToSalesOrderOutput(so *models.SalesOrder) (*SalesOrderOutput, error) {
	lineItems := make([]SalesOrderLineItemOutput, 0)

	for _, item := range so.LineItems {
		lineItemOutput := SalesOrderLineItemOutput{
			ID:        item.ID,
			ItemID:    item.ItemID,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
			Rate:      item.Rate,
			Amount:    item.Amount,
		}

		// Add item info
		if item.Item != nil {
			lineItemOutput.Item = &ItemInfo{
				ID:   item.Item.ID,
				Name: item.Item.Name,
				SKU:  item.Item.ItemDetails.SKU,
			}
		}

		// Add variant info if available
		if item.Variant != nil {
			attributeMap := make(map[string]string)
			for _, attr := range item.Variant.Attributes {
				attributeMap[attr.Key] = attr.Value
			}
			lineItemOutput.Variant = &VariantInfo{
				ID:           item.Variant.ID,
				SKU:          item.Variant.SKU,
				AttributeMap: attributeMap,
			}
		}

		// Convert variant details
		if item.VariantDetails != nil {
			lineItemOutput.VariantDetails = convertVariantDetails(item.VariantDetails)
		}

		lineItems = append(lineItems, lineItemOutput)
	}

	// Convert attachments
	attachments := so.Attachments
	if attachments == nil {
		attachments = []string{}
	}

	// Build output
	output := &SalesOrderOutput{
		ID:                   so.ID,
		SalesOrderNo:         so.SalesOrderNumber,
		CustomerID:           so.CustomerID,
		SalespersonID:        so.SalespersonID,
		ReferenceNo:          so.ReferenceNo,
		SODate:               so.SODate,
		ExpectedShipmentDate: so.ExpectedShipmentDate,
		PaymentTerms:         string(so.PaymentTerms),
		DeliveryMethod:       so.DeliveryMethod,
		LineItems:            lineItems,
		SubTotal:             so.SubTotal,
		ShippingCharges:      so.ShippingCharges,
		TaxAmount:            so.TaxAmount,
		Adjustment:           so.Adjustment,
		Total:                so.Total,
		CustomerNotes:        so.CustomerNotes,
		TermsAndConditions:   so.TermsAndConditions,
		Status:               string(so.Status),
		Attachments:          attachments,
		CreatedAt:            so.CreatedAt,
		UpdatedAt:            so.UpdatedAt,
		CreatedBy:            so.CreatedBy,
		UpdatedBy:            so.UpdatedBy,
	}

	// Add customer info
	if so.Customer != nil {
		output.Customer = &CustomerInfo{
			ID:          so.Customer.ID,
			DisplayName: so.Customer.DisplayName,
			CompanyName: so.Customer.CompanyName,
			Email:       so.Customer.EmailAddress,
			Phone:       so.Customer.WorkPhone,
		}
	}

	// Add salesperson info
	if so.Salesperson != nil {
		output.Salesperson = &SalespersonInfo{
			ID:   so.Salesperson.ID,
			Name: so.Salesperson.Name,
		}
	}

	// Add tax info if available
	if so.Tax != nil {
		taxTypeStr := string(*so.TaxType)
		output.TaxType = &taxTypeStr
		output.Tax = &TaxInfo{
			ID:      so.Tax.ID,
			Name:    so.Tax.Name,
			TaxType: string(so.Tax.TaxType),
			Rate:    so.Tax.Rate,
		}
	}

	return output, nil
}
