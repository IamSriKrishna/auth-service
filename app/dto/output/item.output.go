package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/models"
)

type ItemOutput struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           domain.ItemType `json:"type"`
	Brand          string          `json:"brand,omitempty"`
	ManufacturerID *uint           `json:"manufacturer_id,omitempty"`
	Manufacturer   *ManufacturerInfo `json:"manufacturer,omitempty"`

	ItemDetails  ItemDetailsOutput  `json:"item_details"`
	SalesInfo    SalesInfoOutput    `json:"sales_info"`
	PurchaseInfo PurchaseInfoOutput `json:"purchase_info,omitempty"`
	Inventory    InventoryOutput    `json:"inventory,omitempty"`
	ReturnPolicy ReturnPolicyOutput `json:"return_policy,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ManufacturerInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AttributeDefinitionOutput struct {
	Key     string   `json:"key"`
	Options []string `json:"options"`
}

type ItemDetailsOutput struct {
	Structure            domain.ItemStructure        `json:"structure"`
	Unit                 string                      `json:"unit"`
	SKU                  string                      `json:"sku,omitempty"`
	UPC                  string                      `json:"upc,omitempty"`
	EAN                  string                      `json:"ean,omitempty"`
	MPN                  string                      `json:"mpn,omitempty"`
	ISBN                 string                      `json:"isbn,omitempty"`
	Description          string                      `json:"description,omitempty"`
	AttributeDefinitions []AttributeDefinitionOutput `json:"attribute_definitions,omitempty"`
	Variants             []VariantOutput             `json:"variants,omitempty"`
}

type VariantOutput struct {
	SKU           string            `json:"sku"`
	AttributeMap  map[string]string `json:"attribute_map"`
	SellingPrice  float64           `json:"selling_price"`
	CostPrice     float64           `json:"cost_price"`
	StockQuantity int               `json:"stock_quantity"`
}

type SalesInfoOutput struct {
	Account      string  `json:"account"`
	SellingPrice float64 `json:"selling_price,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Description  string  `json:"description,omitempty"`
}

type PurchaseInfoOutput struct {
	Account           string               `json:"account"`
	CostPrice         float64              `json:"cost_price,omitempty"`
	Currency          string               `json:"currency,omitempty"`
	PreferredVendorID *uint                `json:"preferred_vendor_id,omitempty"`
	PreferredVendor   *PreferredVendorInfo `json:"preferred_vendor,omitempty"`
	Description       string               `json:"description,omitempty"`
}

type PreferredVendorInfo struct {
	ID           uint   `json:"id"`
	DisplayName  string `json:"display_name"`
	CompanyName  string `json:"company_name,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	WorkPhone    string `json:"work_phone,omitempty"`
}

type InventoryOutput struct {
	TrackInventory           bool   `json:"track_inventory"`
	InventoryAccount         string `json:"inventory_account,omitempty"`
	InventoryValuationMethod string `json:"inventory_valuation_method,omitempty"`
	ReorderPoint             int    `json:"reorder_point,omitempty"`
}

type ReturnPolicyOutput struct {
	Returnable bool `json:"returnable"`
}

type ItemListOutput struct {
	Items []ItemOutput `json:"items"`
	Total int64        `json:"total"`
}
type OpeningStockOutput struct {
	OpeningStock            float64   `json:"opening_stock"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type VariantOpeningStockOutput struct {
	VariantID               uint      `json:"variant_id"`
	VariantSKU              string    `json:"variant_sku"`
	OpeningStock            float64   `json:"opening_stock"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type StockSummaryOutput struct {
	StockOnHand      float64 `json:"stock_on_hand"`
	CommittedStock   float64 `json:"committed_stock"`
	AvailableForSale float64 `json:"available_for_sale"`

	PhysicalStockOnHand      float64 `json:"physical_stock_on_hand"`
	PhysicalCommittedStock   float64 `json:"physical_committed_stock"`
	PhysicalAvailableForSale float64 `json:"physical_available_for_sale"`

	ToBeInvoiced float64 `json:"to_be_invoiced"`
	ToBeBilled   float64 `json:"to_be_billed"`
}

func ToItemOutput(item *models.Item) (*ItemOutput, error) {
	// Convert attribute definitions
	attributeDefs := make([]AttributeDefinitionOutput, len(item.ItemDetails.AttributeDefinitions))
	for i, def := range item.ItemDetails.AttributeDefinitions {
		attributeDefs[i] = AttributeDefinitionOutput{
			Key:     def.Key,
			Options: def.Options,
		}
	}

	// Convert variants
	variants := make([]VariantOutput, len(item.ItemDetails.Variants))
	for i, v := range item.ItemDetails.Variants {
		// Convert attributes array to map
		attributeMap := make(map[string]string)
		for _, attr := range v.Attributes {
			attributeMap[attr.Key] = attr.Value
		}

		variants[i] = VariantOutput{
			SKU:           v.SKU,
			AttributeMap:  attributeMap,
			SellingPrice:  v.SellingPrice,
			CostPrice:     v.CostPrice,
			StockQuantity: v.StockQuantity,
		}
	}

	// Build purchase info (might be empty for services)
	purchaseInfo := PurchaseInfoOutput{
		Account:           item.PurchaseInfo.Account,
		CostPrice:         item.PurchaseInfo.CostPrice,
		Currency:          item.PurchaseInfo.Currency,
		PreferredVendorID: item.PurchaseInfo.PreferredVendorID,
		Description:       item.PurchaseInfo.Description,
	}

	// Add vendor details if available
	if item.PurchaseInfo.PreferredVendor != nil {
		purchaseInfo.PreferredVendor = &PreferredVendorInfo{
			ID:           item.PurchaseInfo.PreferredVendor.ID,
			DisplayName:  item.PurchaseInfo.PreferredVendor.DisplayName,
			CompanyName:  item.PurchaseInfo.PreferredVendor.CompanyName,
			EmailAddress: item.PurchaseInfo.PreferredVendor.EmailAddress,
			WorkPhone:    item.PurchaseInfo.PreferredVendor.WorkPhone,
		}
	}

	var manufacturerInfo *ManufacturerInfo
	if item.Manufacturer != nil {
		manufacturerInfo = &ManufacturerInfo{
			ID:   item.Manufacturer.ID,
			Name: item.Manufacturer.Name,
		}
	}

	output := &ItemOutput{
		ID:           item.ID,
		Name:         item.Name,
		Type:         item.Type,
		Brand:        item.Brand,
		Manufacturer: manufacturerInfo,
		ManufacturerID: item.ManufacturerID,
		ItemDetails: ItemDetailsOutput{
			Structure:            item.ItemDetails.Structure,
			Unit:                 item.ItemDetails.Unit,
			SKU:                  item.ItemDetails.SKU,
			UPC:                  item.ItemDetails.UPC,
			EAN:                  item.ItemDetails.EAN,
			MPN:                  item.ItemDetails.MPN,
			ISBN:                 item.ItemDetails.ISBN,
			Description:          item.ItemDetails.Description,
			AttributeDefinitions: attributeDefs,
			Variants:             variants,
		},
		SalesInfo: SalesInfoOutput{
			Account:      item.SalesInfo.Account,
			SellingPrice: item.SalesInfo.SellingPrice,
			Currency:     item.SalesInfo.Currency,
			Description:  item.SalesInfo.Description,
		},
		PurchaseInfo: purchaseInfo,
		Inventory: InventoryOutput{
			TrackInventory:           item.Inventory.TrackInventory,
			InventoryAccount:         item.Inventory.InventoryAccount,
			InventoryValuationMethod: item.Inventory.InventoryValuationMethod,
			ReorderPoint:             item.Inventory.ReorderPoint,
		},
		ReturnPolicy: ReturnPolicyOutput{
			Returnable: item.ReturnPolicy.Returnable,
		},
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	return output, nil
}

func ToItemListOutput(items []models.Item, total int64) (*ItemListOutput, error) {
	outputs := make([]ItemOutput, len(items))
	for i, item := range items {
		output, err := ToItemOutput(&item)
		if err != nil {
			return nil, err
		}
		outputs[i] = *output
	}

	return &ItemListOutput{
		Items: outputs,
		Total: total,
	}, nil
}
