package input

import (
	"errors"
	"fmt"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/models"
)

type CreateItemInput struct {
	Name         string             `json:"name" validate:"required"`
	Type         domain.ItemType    `json:"type" validate:"required,oneof=goods service"`
	ItemDetails  ItemDetailsInput   `json:"item_details" validate:"required"`
	SalesInfo    SalesInfoInput     `json:"sales_info" validate:"required"`
	PurchaseInfo *PurchaseInfoInput `json:"purchase_info"`
	Inventory    *InventoryInput    `json:"inventory"`
	ReturnPolicy *ReturnPolicyInput `json:"return_policy"`
}

type ItemDetailsInput struct {
	Structure   domain.ItemStructure `json:"structure" validate:"required,oneof=single variants"`
	Unit        string               `json:"unit" validate:"required"`
	SKU         string               `json:"sku"`
	UPC         string               `json:"upc"`
	EAN         string               `json:"ean"`
	MPN         string               `json:"mpn"`
	ISBN        string               `json:"isbn"`
	Description string               `json:"description"`

	AttributeDefinitions []AttributeDefinitionInput `json:"attribute_definitions"`
	Variants             []VariantInput             `json:"variants"`
}

type AttributeDefinitionInput struct {
	Key     string   `json:"key" validate:"required"`
	Options []string `json:"options" validate:"required"`
}

type VariantInput struct {
	SKU           string            `json:"sku" validate:"required"`
	AttributeMap  map[string]string `json:"attribute_map" validate:"required"`
	SellingPrice  float64           `json:"selling_price" validate:"required,gt=0"`
	CostPrice     float64           `json:"cost_price" validate:"required,gt=0"`
	StockQuantity float64           `json:"stock_quantity" validate:"gte=0"`
}

type SalesInfoInput struct {
	Account      string  `json:"account" validate:"required"`
	SellingPrice float64 `json:"selling_price"`
	Currency     string  `json:"currency"`
	Description  string  `json:"description"`
}

type PurchaseInfoInput struct {
	Account           string  `json:"account" validate:"required"`
	CostPrice         float64 `json:"cost_price"`
	Currency          string  `json:"currency"`
	PreferredVendorID *uint   `json:"preferred_vendor_id"`
	Description       string  `json:"description"`
}

type InventoryInput struct {
	TrackInventory           bool   `json:"track_inventory"`
	InventoryAccount         string `json:"inventory_account"`
	InventoryValuationMethod string `json:"inventory_valuation_method"`
	ReorderPoint             int    `json:"reorder_point"`
}

type ReturnPolicyInput struct {
	Returnable bool `json:"returnable"`
}

type UpdateItemInput struct {
	Name         *string            `json:"name"`
	Type         *domain.ItemType   `json:"type"`
	ItemDetails  *ItemDetailsInput  `json:"item_details"`
	SalesInfo    *SalesInfoInput    `json:"sales_info"`
	PurchaseInfo *PurchaseInfoInput `json:"purchase_info"`
	Inventory    *InventoryInput    `json:"inventory"`
	ReturnPolicy *ReturnPolicyInput `json:"return_policy"`
}

type OpeningStockInput struct {
	OpeningStock            float64 `json:"opening_stock" validate:"gte=0"`
	OpeningStockRatePerUnit float64 `json:"opening_stock_rate_per_unit" validate:"gte=0"`
}

type VariantOpeningStockInput struct {
	VariantSKU              string  `json:"variant_sku" validate:"required"`
	OpeningStock            float64 `json:"opening_stock" validate:"gte=0"`
	OpeningStockRatePerUnit float64 `json:"opening_stock_rate_per_unit" validate:"gte=0"`
}

type UpdateVariantsOpeningStockInput struct {
	Variants []VariantOpeningStockInput `json:"variants" validate:"required,dive"`
}

func (c *CreateItemInput) ToItemDetails() models.ItemDetails {
	itemDetails := models.ItemDetails{
		Structure:   c.ItemDetails.Structure,
		Unit:        c.ItemDetails.Unit,
		SKU:         c.ItemDetails.SKU,
		UPC:         c.ItemDetails.UPC,
		EAN:         c.ItemDetails.EAN,
		MPN:         c.ItemDetails.MPN,
		ISBN:        c.ItemDetails.ISBN,
		Description: c.ItemDetails.Description,
	}

	if len(c.ItemDetails.AttributeDefinitions) > 0 {
		itemDetails.AttributeDefinitions = make([]models.AttributeDefinition, len(c.ItemDetails.AttributeDefinitions))
		for i, attr := range c.ItemDetails.AttributeDefinitions {
			itemDetails.AttributeDefinitions[i] = models.AttributeDefinition{
				Key:     attr.Key,
				Options: attr.Options,
			}
		}
	}

	if len(c.ItemDetails.Variants) > 0 {
		itemDetails.Variants = make([]models.Variant, len(c.ItemDetails.Variants))
		for i, v := range c.ItemDetails.Variants {
			attributes := make([]models.VariantAttribute, 0, len(v.AttributeMap))
			for key, value := range v.AttributeMap {
				attributes = append(attributes, models.VariantAttribute{
					Key:   key,
					Value: value,
				})
			}

			itemDetails.Variants[i] = models.Variant{
				SKU:           v.SKU,
				Attributes:    attributes,
				SellingPrice:  v.SellingPrice,
				CostPrice:     v.CostPrice,
				StockQuantity: v.StockQuantity,
			}
		}
	}

	return itemDetails
}

func (c *CreateItemInput) ToSalesInfo() models.SalesInfo {
	return models.SalesInfo{
		Account:      c.SalesInfo.Account,
		SellingPrice: c.SalesInfo.SellingPrice,
		Currency:     c.SalesInfo.Currency,
		Description:  c.SalesInfo.Description,
	}
}

func (c *CreateItemInput) ToPurchaseInfo() models.PurchaseInfo {
	if c.PurchaseInfo == nil {
		return models.PurchaseInfo{}
	}
	return models.PurchaseInfo{
		Account:           c.PurchaseInfo.Account,
		CostPrice:         c.PurchaseInfo.CostPrice,
		Currency:          c.PurchaseInfo.Currency,
		PreferredVendorID: c.PurchaseInfo.PreferredVendorID,
		Description:       c.PurchaseInfo.Description,
	}
}

func (c *CreateItemInput) ToInventory() models.Inventory {
	if c.Inventory == nil {
		return models.Inventory{
			TrackInventory: false,
		}
	}
	return models.Inventory{
		TrackInventory:           c.Inventory.TrackInventory,
		InventoryAccount:         c.Inventory.InventoryAccount,
		InventoryValuationMethod: c.Inventory.InventoryValuationMethod,
		ReorderPoint:             c.Inventory.ReorderPoint,
	}
}

func (c *CreateItemInput) ToReturnPolicy() models.ReturnPolicy {
	if c.ReturnPolicy == nil {
		return models.ReturnPolicy{}
	}
	return models.ReturnPolicy{
		Returnable: c.ReturnPolicy.Returnable,
	}
}

// ValidateVariantAttributes validates that variants have proper attribute mappings
func (c *CreateItemInput) ValidateVariantAttributes() error {
	if c.ItemDetails.Structure != "variants" {
		return nil
	}

	if len(c.ItemDetails.AttributeDefinitions) == 0 {
		return errors.New("variant items must define attributes")
	}

	if len(c.ItemDetails.Variants) == 0 {
		return errors.New("variant items must have at least one variant")
	}

	// Create a map of valid attribute keys
	validKeys := make(map[string][]string)
	for _, attr := range c.ItemDetails.AttributeDefinitions {
		validKeys[attr.Key] = attr.Options
	}

	// Validate each variant
	for i, variant := range c.ItemDetails.Variants {
		if len(variant.AttributeMap) == 0 {
			return fmt.Errorf("variant %d (%s) must have attribute_map defined", i+1, variant.SKU)
		}

		// Check that all defined attributes are in the variant's attribute_map
		for attrKey, validOptions := range validKeys {
			variantValue, exists := variant.AttributeMap[attrKey]
			if !exists {
				return fmt.Errorf("variant %d (%s) missing required attribute '%s'", i+1, variant.SKU, attrKey)
			}

			// Check that the value is one of the valid options
			valid := false
			for _, option := range validOptions {
				if variantValue == option {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("variant %d (%s) has invalid value '%s' for attribute '%s'. Valid options: %v", i+1, variant.SKU, variantValue, attrKey, validOptions)
			}
		}

		// Check for extra attributes not defined
		for variantKey := range variant.AttributeMap {
			if _, exists := validKeys[variantKey]; !exists {
				return fmt.Errorf("variant %d (%s) has undefined attribute '%s'", i+1, variant.SKU, variantKey)
			}
		}
	}

	return nil
}
