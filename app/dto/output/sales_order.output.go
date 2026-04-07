package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type SalesOrderOutput struct {
	ID                   string                     `json:"id"`
	SalesOrderNo         string                     `json:"sales_order_no"`
	CustomerID           uint                       `json:"customer_id"`
	ReferenceNo          string                     `json:"reference_no,omitempty"`
	Status               string                     `json:"status"`
	Date                 time.Time                  `json:"date"`
	ExpectedShipmentDate time.Time                  `json:"expected_shipment_date"`
	DeliveryMethod       string                     `json:"delivery_method,omitempty"`
	PaymentTerms         string                     `json:"payment_terms"`
	LineItems            []SalesOrderLineItemOutput `json:"line_items"`
	SubTotal             float64                    `json:"sub_total"`
	ShippingCharges      float64                    `json:"shipping_charges"`
	Adjustment           float64                    `json:"adjustment"`
	TaxRate              float64                    `json:"tax_rate"`
	TaxTotal             float64                    `json:"tax_total"`
	Total                float64                    `json:"total"`
	CustomerNotes        string                     `json:"customer_notes,omitempty"`
	TermsAndConditions   string                     `json:"terms_and_conditions,omitempty"`
	SalespersonID        *uint                      `json:"salesperson_id,omitempty"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
}

type SalesOrderLineItemOutput struct {
	ID                uint                   `json:"id"`
	ProductID         string                 `json:"product_id"`
	ProductName       string                 `json:"product_name"`
	SKU               string                 `json:"sku"`
	Account           string                 `json:"account"`
	Quantity          float64                `json:"quantity"`
	DeliveredQuantity float64                `json:"delivered_quantity"`
	Rate              float64                `json:"rate"`
	Amount            float64                `json:"amount"`
	VariantSKU        string                 `json:"variant_sku,omitempty"`
	VariantDetails    map[string]interface{} `json:"variant_details,omitempty"`
}

func ToSalesOrderOutput(so *models.SalesOrder) (*SalesOrderOutput, error) {
	lineItems := make([]SalesOrderLineItemOutput, 0)

	for _, item := range so.LineItems {
		variantDetails := make(map[string]interface{})
		if item.VariantDetails != nil {
			variantDetails = item.VariantDetails
		}

		lineItemOutput := SalesOrderLineItemOutput{
			ID:                item.ID,
			ProductID:         item.ProductID,
			ProductName:       item.ProductName,
			SKU:               item.SKU,
			Account:           item.Account,
			Quantity:          item.Quantity,
			DeliveredQuantity: item.DeliveredQuantity,
			Rate:              item.Rate,
			Amount:            item.Amount,
			VariantSKU:        item.VariantSKU,
			VariantDetails:    variantDetails,
		}
		lineItems = append(lineItems, lineItemOutput)
	}

	output := &SalesOrderOutput{
		ID:                   so.ID,
		SalesOrderNo:         so.SalesOrderNumber,
		CustomerID:           so.CustomerID,
		ReferenceNo:          so.ReferenceNo,
		Status:               string(so.Status),
		Date:                 so.Date,
		ExpectedShipmentDate: so.ExpectedShipmentDate,
		DeliveryMethod:       so.DeliveryMethod,
		PaymentTerms:         string(so.PaymentTerms),
		LineItems:            lineItems,
		SubTotal:             so.SubTotal,
		ShippingCharges:      so.ShippingCharges,
		Adjustment:           so.Adjustment,
		TaxRate:              so.TaxRate,
		TaxTotal:             so.TaxTotal,
		Total:                so.Total,
		CustomerNotes:        so.CustomerNotes,
		TermsAndConditions:   so.TermsAndConditions,
		SalespersonID:        so.SalespersonID,
		CreatedAt:            so.CreatedAt,
		UpdatedAt:            so.UpdatedAt,
	}

	return output, nil
}
