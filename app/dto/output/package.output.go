package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type PackageOutput struct {
	ID            string              `json:"id"`
	PackageSlipNo string              `json:"package_slip_no"`
	SalesOrderID  string              `json:"sales_order_id"`
	SalesOrder    *SalesOrderInfo     `json:"sales_order,omitempty"`
	CustomerID    uint                `json:"customer_id"`
	Customer      *CustomerInfo       `json:"customer,omitempty"`
	PackageDate   time.Time           `json:"package_date"`
	Items         []PackageItemOutput `json:"items"`
	Status        string              `json:"status"`
	InternalNotes string              `json:"internal_notes,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	CreatedBy     string              `json:"created_by,omitempty"`
	UpdatedBy     string              `json:"updated_by,omitempty"`
}

type PackageItemOutput struct {
	ID             uint              `json:"id"`
	ItemID         string            `json:"item_id"`
	Item           *ItemInfo         `json:"item,omitempty"`
	VariantSKU     *string           `json:"variant_sku,omitempty"`
	Variant        *VariantInfo      `json:"variant,omitempty"`
	OrderedQty     float64           `json:"ordered_qty"`
	PackedQty      float64           `json:"packed_qty"`
	VariantDetails map[string]string `json:"variant_details,omitempty"`
}

type SalesOrderInfo struct {
	ID               string    `json:"id"`
	SalesOrderNo     string    `json:"sales_order_no"`
	CustomerID       uint      `json:"customer_id"`
	ReferenceNo      string    `json:"reference_no,omitempty"`
	SODate           time.Time `json:"sales_order_date"`
	ExpectedShipDate time.Time `json:"expected_shipment_date"`
	Status           string    `json:"status"`
}

func ToPackageOutput(pkg *models.Package) (*PackageOutput, error) {
	items := make([]PackageItemOutput, 0)

	for _, item := range pkg.Items {
		packageItemOutput := PackageItemOutput{
			ID:         item.ID,
			ItemID:     item.ItemID,
			VariantSKU: item.VariantSKU,
			OrderedQty: item.OrderedQty,
			PackedQty:  item.PackedQty,
		}

		if item.Item != nil {
			packageItemOutput.Item = &ItemInfo{
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
			packageItemOutput.Variant = &VariantInfo{
				ID:           item.Variant.ID,
				SKU:          item.Variant.SKU,
				AttributeMap: attributeMap,
			}
		}

		if item.VariantDetails != nil {
			packageItemOutput.VariantDetails = convertVariantDetails(item.VariantDetails)
		}

		items = append(items, packageItemOutput)
	}

	var soInfo *SalesOrderInfo
	if pkg.SalesOrder != nil {
		soInfo = &SalesOrderInfo{
			ID:               pkg.SalesOrder.ID,
			SalesOrderNo:     pkg.SalesOrder.SalesOrderNumber,
			CustomerID:       pkg.SalesOrder.CustomerID,
			ReferenceNo:      pkg.SalesOrder.ReferenceNo,
			SODate:           pkg.SalesOrder.SODate,
			ExpectedShipDate: pkg.SalesOrder.ExpectedShipmentDate,
			Status:           string(pkg.SalesOrder.Status),
		}
	}

	var customerInfo *CustomerInfo
	if pkg.Customer != nil {
		customerInfo = &CustomerInfo{
			ID:          pkg.Customer.ID,
			DisplayName: pkg.Customer.DisplayName,
			CompanyName: pkg.Customer.CompanyName,
			Email:       pkg.Customer.EmailAddress,
			Phone:       pkg.Customer.Mobile,
		}
	}

	output := &PackageOutput{
		ID:            pkg.ID,
		PackageSlipNo: pkg.PackageSlipNo,
		SalesOrderID:  pkg.SalesOrderID,
		SalesOrder:    soInfo,
		CustomerID:    pkg.CustomerID,
		Customer:      customerInfo,
		PackageDate:   pkg.PackageDate,
		Items:         items,
		Status:        string(pkg.Status),
		InternalNotes: pkg.InternalNotes,
		CreatedAt:     pkg.CreatedAt,
		UpdatedAt:     pkg.UpdatedAt,
		CreatedBy:     pkg.CreatedBy,
		UpdatedBy:     pkg.UpdatedBy,
	}

	return output, nil
}
