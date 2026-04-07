package output

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type BillOutput struct {
	ID              string               `json:"id"`
	BillNumber      string               `json:"bill_number"`
	PurchaseOrderID *string              `json:"purchase_order_id,omitempty"`
	PurchaseOrderNo string               `json:"purchase_order_no,omitempty"`
	VendorID        uint                 `json:"vendor_id"`
	VendorName      string               `json:"vendor_name,omitempty"`
	Vendor          *VendorInfo          `json:"vendor,omitempty"`
	BillingAddress  string               `json:"billing_address,omitempty"`
	OrderNumber     string               `json:"order_number,omitempty"`
	BillDate        time.Time            `json:"bill_date"`
	DueDate         time.Time            `json:"due_date"`
	PaymentTerms    string               `json:"payment_terms"`
	Subject         string               `json:"subject,omitempty"`
	LineItems       []BillLineItemOutput `json:"line_items"`
	SubTotal        float64              `json:"sub_total"`
	Discount        float64              `json:"discount"`
	TaxType         *string              `json:"tax_type,omitempty"`
	TaxID           *uint                `json:"tax_id,omitempty"`
	Tax             *TaxInfo             `json:"tax,omitempty"`
	TaxAmount       float64              `json:"tax_amount"`
	Adjustment      float64              `json:"adjustment"`
	Total           float64              `json:"total"`
	Notes           string               `json:"notes,omitempty"`
	Status          string               `json:"status"`
	PaymentStatus   string               `json:"payment_status,omitempty"`
	Attachments     []string             `json:"attachments,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	UserID          int                  `json:"user_id,omitempty"`
	UserName        string               `json:"user_name,omitempty"`
	CompanyID       int                  `json:"company_id,omitempty"`
	CompanyName     string               `json:"company_name,omitempty"`
	UpdatedByName   string               `json:"updated_by_name,omitempty"`
}

type BillLineItemOutput struct {
	ID             uint              `json:"id"`
	ItemID         *string           `json:"item_id"`
	Item           *ItemInfo         `json:"item,omitempty"`
	VariantSKU     *string           `json:"variant_sku,omitempty"`
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
			VariantSKU:  item.VariantSKU,
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
		UserID:         parseUserID(bill.CreatedBy),
		UserName:       bill.CreatedByUserName,
		CompanyID:      int(bill.CreatedByCompanyID),
		CompanyName:    bill.CreatedByCompanyName,
		UpdatedByName:  bill.UpdatedByUserName,
	}, nil
}

func parseUserID(userIDStr string) int {
	if userIDStr == "" {
		return 0
	}
	var userID int
	_, _ = fmt.Sscanf(userIDStr, "%d", &userID)
	return userID
}

// CreateBillVariantOutput creates a bill output with variant-specific line items
func CreateBillVariantOutput(billNo string, purchaseOrderID *string, vendorID uint, vendorName string, lineItems []BillLineItemOutput, subTotal float64, taxAmount float64, total float64) *BillOutput {
	return &BillOutput{
		BillNumber:      billNo,
		PurchaseOrderID: purchaseOrderID,
		VendorID:        vendorID,
		VendorName:      vendorName,
		LineItems:       lineItems,
		SubTotal:        subTotal,
		TaxAmount:       taxAmount,
		Total:           total,
		Status:          "issued",
		PaymentStatus:   "unpaid",
		BillDate:        time.Now(),
		DueDate:         time.Now().AddDate(0, 1, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
