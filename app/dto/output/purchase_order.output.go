package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type PurchaseOrderOutput struct {
	ID                  string                        `json:"id"`
	PurchaseOrderNo     string                        `json:"purchase_order_no"`
	VendorID            uint                          `json:"vendor_id"`
	Vendor              *VendorInfo                   `json:"vendor,omitempty"`
	DeliveryAddressType string                        `json:"delivery_address_type"`
	DeliveryAddressID   *uint                         `json:"delivery_address_id,omitempty"`
	OrganizationName    string                        `json:"organization_name,omitempty"`
	OrganizationAddress string                        `json:"organization_address,omitempty"`
	CustomerID          *uint                         `json:"customer_id,omitempty"`
	Customer            *CustomerInfo                 `json:"customer,omitempty"`
	ReferenceNo         string                        `json:"reference_no,omitempty"`
	Date                time.Time                     `json:"date"`
	DeliveryDate        time.Time                     `json:"delivery_date"`
	PaymentTerms        string                        `json:"payment_terms"`
	ShipmentPreference  string                        `json:"shipment_preference,omitempty"`
	LineItems           []PurchaseOrderLineItemOutput `json:"line_items"`
	SubTotal            float64                       `json:"sub_total"`
	Discount            float64                       `json:"discount"`
	DiscountType        string                        `json:"discount_type,omitempty"`
	TaxType             *string                       `json:"tax_type,omitempty"`
	TaxID               *uint                         `json:"tax_id,omitempty"`
	Tax                 *TaxInfo                      `json:"tax,omitempty"`
	TaxAmount           float64                       `json:"tax_amount"`
	Adjustment          float64                       `json:"adjustment"`
	Total               float64                       `json:"total"`
	Notes               string                        `json:"notes,omitempty"`
	TermsAndConditions  string                        `json:"terms_and_conditions,omitempty"`
	Status              string                        `json:"status"`
	Attachments         []string                      `json:"attachments,omitempty"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`
	CreatedBy           string                        `json:"created_by,omitempty"`
	UpdatedBy           string                        `json:"updated_by,omitempty"`
}

type PurchaseOrderLineItemOutput struct {
	ID             uint              `json:"id"`
	ItemID         string            `json:"item_id"`
	Item           *ItemInfo         `json:"item,omitempty"`
	VariantID      *uint             `json:"variant_id,omitempty"`
	Variant        *VariantInfo      `json:"variant,omitempty"`
	Account        string            `json:"account"`
	Quantity       float64           `json:"quantity"`
	Rate           float64           `json:"rate"`
	Amount         float64           `json:"amount"`
	VariantDetails map[string]string `json:"variant_details,omitempty"`
}

type PurchaseOrderListOutput struct {
	PurchaseOrders []PurchaseOrderOutput `json:"purchase_orders"`
	Total          int64                 `json:"total"`
}

func ToPurchaseOrderOutput(po *models.PurchaseOrder) (*PurchaseOrderOutput, error) {
	lineItems := make([]PurchaseOrderLineItemOutput, len(po.LineItems))
	for i, item := range po.LineItems {
		lineItemOutput := PurchaseOrderLineItemOutput{
			ID:       item.ID,
			ItemID:   item.ItemID,
			Account:  item.Account,
			Quantity: item.Quantity,
			Rate:     item.Rate,
			Amount:   item.Amount,
		}

		if item.VariantID != nil {
			lineItemOutput.VariantID = item.VariantID
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

		if item.VariantDetails != nil {
			lineItemOutput.VariantDetails = convertVariantDetails(item.VariantDetails)
		}

		lineItems[i] = lineItemOutput
	}

	attachments := po.Attachments
	if attachments == nil {
		attachments = []string{}
	}

	output := &PurchaseOrderOutput{
		ID:                  po.ID,
		PurchaseOrderNo:     po.PurchaseOrderNumber,
		VendorID:            po.VendorID,
		DeliveryAddressType: po.DeliveryAddressType,
		DeliveryAddressID:   po.DeliveryAddressID,
		OrganizationName:    po.OrganizationName,
		OrganizationAddress: po.OrganizationAddress,
		CustomerID:          po.CustomerID,
		ReferenceNo:         po.ReferenceNo,
		Date:                po.PODate,
		DeliveryDate:        po.DeliveryDate,
		PaymentTerms:        string(po.PaymentTerms),
		ShipmentPreference:  po.ShipmentPreference,
		LineItems:           lineItems,
		SubTotal:            po.SubTotal,
		Discount:            po.Discount,
		DiscountType:        po.DiscountType,
		TaxAmount:           po.TaxAmount,
		Adjustment:          po.Adjustment,
		Total:               po.Total,
		Notes:               po.Notes,
		TermsAndConditions:  po.TermsAndConditions,
		Status:              string(po.Status),
		Attachments:         attachments,
		CreatedAt:           po.CreatedAt,
		UpdatedAt:           po.UpdatedAt,
		CreatedBy:           po.CreatedBy,
		UpdatedBy:           po.UpdatedBy,
	}

	if po.Vendor != nil {
		output.Vendor = &VendorInfo{
			ID:           po.Vendor.ID,
			DisplayName:  po.Vendor.DisplayName,
			CompanyName:  po.Vendor.CompanyName,
			EmailAddress: po.Vendor.EmailAddress,
			WorkPhone:    po.Vendor.WorkPhone,
		}
	}

	if po.Customer != nil {
		output.Customer = &CustomerInfo{
			ID:          po.Customer.ID,
			DisplayName: po.Customer.DisplayName,
			CompanyName: po.Customer.CompanyName,
			Email:       po.Customer.EmailAddress,
			Phone:       po.Customer.WorkPhone,
		}
	}

	if po.Tax != nil {
		taxTypeStr := string(*po.TaxType)
		output.TaxType = &taxTypeStr
		output.Tax = &TaxInfo{
			ID:      po.Tax.ID,
			Name:    po.Tax.Name,
			TaxType: string(po.Tax.TaxType),
			Rate:    po.Tax.Rate,
		}
	}

	return output, nil
}
