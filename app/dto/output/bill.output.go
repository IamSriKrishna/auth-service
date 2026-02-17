package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type BillOutput struct {
	ID             string               `json:"id"`
	BillNumber     string               `json:"bill_number"`
	VendorID       uint                 `json:"vendor_id"`
	Vendor         *VendorInfo          `json:"vendor,omitempty"`
	BillingAddress string               `json:"billing_address,omitempty"`
	OrderNumber    string               `json:"order_number,omitempty"`
	BillDate       time.Time            `json:"bill_date"`
	DueDate        time.Time            `json:"due_date"`
	PaymentTerms   string               `json:"payment_terms"`
	Subject        string               `json:"subject,omitempty"`
	LineItems      []BillLineItemOutput `json:"line_items"`
	SubTotal       float64              `json:"sub_total"`
	Discount       float64              `json:"discount"`
	TaxType        *string              `json:"tax_type,omitempty"`
	TaxID          *uint                `json:"tax_id,omitempty"`
	Tax            *TaxInfo             `json:"tax,omitempty"`
	TaxAmount      float64              `json:"tax_amount"`
	Adjustment     float64              `json:"adjustment"`
	Total          float64              `json:"total"`
	Notes          string               `json:"notes,omitempty"`
	Status         string               `json:"status"`
	Attachments    []string             `json:"attachments,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	CreatedBy      string               `json:"created_by,omitempty"`
	UpdatedBy      string               `json:"updated_by,omitempty"`
}

type BillLineItemOutput struct {
	ID             uint              `json:"id"`
	ItemID         string            `json:"item_id"`
	Item           *ItemInfo         `json:"item,omitempty"`
	VariantID      *uint             `json:"variant_id,omitempty"`
	Variant        *VariantInfo      `json:"variant,omitempty"`
	Description    string            `json:"description,omitempty"`
	Account        string            `json:"account,omitempty"`
	Quantity       float64           `json:"quantity"`
	Rate           float64           `json:"rate"`
	Amount         float64           `json:"amount"`
	VariantDetails map[string]string `json:"variant_details,omitempty"`
}

func ToBillOutput(bill *models.Bill) (*BillOutput, error) {
	lineItems := make([]BillLineItemOutput, 0)

	for _, item := range bill.LineItems {
		lineItemOutput := BillLineItemOutput{
			ID:          item.ID,
			ItemID:      item.ItemID,
			VariantID:   item.VariantID,
			Description: item.Description,
			Account:     item.Account,
			Quantity:    item.Quantity,
			Rate:        item.Rate,
			Amount:      item.Amount,
		}

		if item.Item != nil {
			lineItemOutput.Item = &ItemInfo{
				ID:   item.Item.ID,
				Name: item.Item.Name,
				SKU:  item.Item.ItemDetails.SKU,
			}
		}

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

		lineItems = append(lineItems, lineItemOutput)
	}

	var vendor *VendorInfo
	if bill.Vendor != nil {
		vendor = &VendorInfo{
			ID:           bill.Vendor.ID,
			DisplayName:  bill.Vendor.DisplayName,
			CompanyName:  bill.Vendor.CompanyName,
			EmailAddress: bill.Vendor.EmailAddress,
			WorkPhone:    bill.Vendor.WorkPhone,
		}
	}

	var tax *TaxInfo
	if bill.Tax != nil {
		tax = &TaxInfo{
			ID:   bill.Tax.ID,
			Name: bill.Tax.Name,
			Rate: bill.Tax.Rate,
		}
	}

	return &BillOutput{
		ID:             bill.ID,
		BillNumber:     bill.BillNumber,
		VendorID:       bill.VendorID,
		Vendor:         vendor,
		BillingAddress: bill.BillingAddress,
		OrderNumber:    bill.OrderNumber,
		BillDate:       bill.BillDate,
		DueDate:        bill.DueDate,
		PaymentTerms:   string(bill.PaymentTerms),
		Subject:        bill.Subject,
		LineItems:      lineItems,
		SubTotal:       bill.SubTotal,
		Discount:       bill.Discount,
		TaxType:        (*string)(bill.TaxType),
		TaxID:          bill.TaxID,
		Tax:            tax,
		TaxAmount:      bill.TaxAmount,
		Adjustment:     bill.Adjustment,
		Total:          bill.Total,
		Notes:          bill.Notes,
		Status:         string(bill.Status),
		Attachments:    bill.Attachments,
		CreatedAt:      bill.CreatedAt,
		UpdatedAt:      bill.UpdatedAt,
		CreatedBy:      bill.CreatedBy,
		UpdatedBy:      bill.UpdatedBy,
	}, nil
}
