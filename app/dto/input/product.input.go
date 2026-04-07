package input

import (
	"errors"
	"fmt"

	"github.com/bbapp-org/auth-service/app/models"
)

// CreateProductInput represents the input for creating a new product with variants.
// A product is a base item that can have multiple variants.
// Variants are specific configurations/options of the product (e.g., color, size).
//
// Example: "500ml Water Bottle" product can have variants like "Red 500ml", "Blue 500ml", etc.
// Each variant has its own SKU, pricing, and stock.
type CreateProductInput struct {
	Name           string              `json:"name" validate:"required" example:"500ml Water Bottle"`
	ProductDetails ProductDetailsInput `json:"product_details" validate:"required"`
	SalesInfo      SalesInfoInput      `json:"sales_info" validate:"required"`
	PurchaseInfo   *PurchaseInfoInput  `json:"purchase_info"`
	Inventory      *InventoryInput     `json:"inventory"`
	ReturnPolicy   *ReturnPolicyInput  `json:"return_policy"`
}

// ProductDetailsInput contains variant and attribute details for the product.
// Attributes define the possible variant options (e.g., Color: [Red, Blue, Green])
// Variants are specific combinations of these attributes
// For simple products (no variants), both AttributeDefinitions and Variants can be omitted,
// and a default variant will be auto-created from the base SKU and pricing.
type ProductDetailsInput struct {
	// Unit of measurement for this product (e.g., "piece", "kg", "liter")
	Unit string `json:"unit" validate:"required" example:"piece"`
	// Base SKU for the product (optional, used as prefix for variants)
	BaseSKU string `json:"base_sku" example:"WTR-BOT-500"`
	// Universal Product Code
	UPC string `json:"upc" example:"8904220500001"`
	// European Article Number
	EAN string `json:"ean" example:"8904220500001"`
	// Manufacturer Part Number
	MPN string `json:"mpn" example:"MPN-123"`
	// ISBN (for books)
	ISBN string `json:"isbn" example:"978-0-123456-78-9"`
	// Product description
	Description string `json:"description" example:"Premium 500ml drinking water bottle"`
	// Manufacturer ID reference
	ManufacturerID *uint `json:"manufacturer_id"`
	// Attribute definitions that define possible variant combinations
	// Example: Color: [Red, Blue, Green], Size: [500ml, 1L]
	// Optional: can be omitted for products without variants
	AttributeDefinitions []ProductAttributeDefinitionInput `json:"attribute_definitions" validate:"omitempty,dive"`
	// Array of variants for this product
	// Each variant must satisfy the attribute definitions
	// Optional: can be omitted for simple products without variants (a default variant will be auto-created)
	Variants []ProductVariantInput `json:"variants" validate:"omitempty,dive"`
}

// ProductAttributeDefinitionInput defines possible attribute values for product variants.
// Example:
//
//	Key: "Color", Options: ["Red", "Blue", "Green"]
//	Key: "Size", Options: ["Small", "Large"]
type ProductAttributeDefinitionInput struct {
	// Attribute name/key (e.g., "Color", "Size", "Capacity")
	Key string `json:"key" validate:"required" example:"Color"`
	// Possible values for this attribute
	Options []string `json:"options" validate:"required" example:"[\"Red\",\"Blue\",\"Green\"]"`
}

// ProductVariantInput represents a specific variant of the product.
// A variant is a unique combination of attribute values with its own SKU, pricing, and stock.
//
// Example variants for "500ml Water Bottle":
//   - Variant 1: Color=Red, Size=500ml, SKU=WTR-500-RED, Price=150, Stock=1000
//   - Variant 2: Color=Blue, Size=500ml, SKU=WTR-500-BLUE, Price=150, Stock=800
type ProductVariantInput struct {
	// Unique SKU for this variant (must be globally unique across all products)
	SKU string `json:"sku" validate:"required" example:"WTR-BOT-500-RED"`
	// Optional display name for the variant
	VariantName string `json:"variant_name" example:"500ml Red Water Bottle"`
	// Mapping of attribute keys to values for this variant
	// Example: {"Color": "Red", "Capacity": "500ml"}
	// Keys must match the attribute_definitions keys and values from the options
	AttributeMap map[string]string `json:"attribute_map" validate:"required" example:"{\"Color\":\"Red\",\"Capacity\":\"500ml\"}"`
	// Selling price for this variant
	SellingPrice float64 `json:"selling_price" validate:"required,gt=0" example:"150"`
	// Cost price for this variant
	CostPrice float64 `json:"cost_price" validate:"required,gt=0" example:"75"`
	// Initial stock quantity for this variant
	StockQuantity float64 `json:"stock_quantity" validate:"gte=0" example:"1000"`
	// Whether this variant is active and available for sale
	IsActive bool `json:"is_active" example:"true"`
}

// UpdateProductInput represents updates to an existing product
type UpdateProductInput struct {
	Name           *string              `json:"name"`
	ProductDetails *ProductDetailsInput `json:"product_details"`
	SalesInfo      *SalesInfoInput      `json:"sales_info"`
	PurchaseInfo   *PurchaseInfoInput   `json:"purchase_info"`
	Inventory      *InventoryInput      `json:"inventory"`
	ReturnPolicy   *ReturnPolicyInput   `json:"return_policy"`
}

// ProductVariantOpeningStockInput for setting opening stock for a variant
type ProductVariantOpeningStockInput struct {
	VariantSKU              string  `json:"variant_sku" validate:"required"`
	OpeningStock            float64 `json:"opening_stock" validate:"gte=0"`
	OpeningStockRatePerUnit float64 `json:"opening_stock_rate_per_unit" validate:"gte=0"`
}

// UpdateProductVariantsOpeningStockInput for updating opening stock for multiple variants
type UpdateProductVariantsOpeningStockInput struct {
	Variants []ProductVariantOpeningStockInput `json:"variants" validate:"required,dive"`
}

// ProductOpeningStockInput for setting opening stock for the base product
type ProductOpeningStockInput struct {
	OpeningStock            float64 `json:"opening_stock" validate:"gte=0"`
	OpeningStockRatePerUnit float64 `json:"opening_stock_rate_per_unit" validate:"gte=0"`
}

// Conversion functions

// ToProductDetails converts ProductDetailsInput to models.ProductDetails
func (p *CreateProductInput) ToProductDetails() models.ProductDetails {
	productDetails := models.ProductDetails{
		Unit:           p.ProductDetails.Unit,
		BaseSKU:        p.ProductDetails.BaseSKU,
		UPC:            p.ProductDetails.UPC,
		EAN:            p.ProductDetails.EAN,
		MPN:            p.ProductDetails.MPN,
		ISBN:           p.ProductDetails.ISBN,
		Description:    p.ProductDetails.Description,
		ManufacturerID: p.ProductDetails.ManufacturerID,
	}

	if len(p.ProductDetails.AttributeDefinitions) > 0 {
		productDetails.AttributeDefinitions = make([]models.ProductAttributeDefinition, len(p.ProductDetails.AttributeDefinitions))
		for i, attr := range p.ProductDetails.AttributeDefinitions {
			productDetails.AttributeDefinitions[i] = models.ProductAttributeDefinition{
				Key:     attr.Key,
				Options: attr.Options,
			}
		}
	}

	if len(p.ProductDetails.Variants) > 0 {
		productDetails.ProductVariants = make([]models.ProductVariant, len(p.ProductDetails.Variants))
		for i, v := range p.ProductDetails.Variants {
			attributes := make([]models.ProductVariantAttribute, 0, len(v.AttributeMap))
			for key, value := range v.AttributeMap {
				attributes = append(attributes, models.ProductVariantAttribute{
					Key:   key,
					Value: value,
				})
			}

			productDetails.ProductVariants[i] = models.ProductVariant{
				SKU:           v.SKU,
				VariantName:   v.VariantName,
				Attributes:    attributes,
				SellingPrice:  v.SellingPrice,
				CostPrice:     v.CostPrice,
				StockQuantity: v.StockQuantity,
				IsActive:      v.IsActive,
			}
		}
	}

	return productDetails
}

// ToSalesInfo converts to SalesInfo model
func (p *CreateProductInput) ToSalesInfo() models.SalesInfo {
	return models.SalesInfo{
		Account:      p.SalesInfo.Account,
		SellingPrice: p.SalesInfo.SellingPrice,
		Currency:     p.SalesInfo.Currency,
		Description:  p.SalesInfo.Description,
	}
}

// ToPurchaseInfo converts to PurchaseInfo model
func (p *CreateProductInput) ToPurchaseInfo() models.PurchaseInfo {
	if p.PurchaseInfo == nil {
		return models.PurchaseInfo{}
	}
	return models.PurchaseInfo{
		Account:           p.PurchaseInfo.Account,
		CostPrice:         p.PurchaseInfo.CostPrice,
		Currency:          p.PurchaseInfo.Currency,
		PreferredVendorID: p.PurchaseInfo.PreferredVendorID,
		Description:       p.PurchaseInfo.Description,
	}
}

// ToInventory converts to Inventory model
func (p *CreateProductInput) ToInventory() models.Inventory {
	if p.Inventory == nil {
		return models.Inventory{
			TrackInventory: false,
		}
	}
	return models.Inventory{
		TrackInventory:           p.Inventory.TrackInventory,
		InventoryAccount:         p.Inventory.InventoryAccount,
		InventoryValuationMethod: p.Inventory.InventoryValuationMethod,
		ReorderPoint:             p.Inventory.ReorderPoint,
	}
}

// ToReturnPolicy converts to ReturnPolicy model
func (p *CreateProductInput) ToReturnPolicy() models.ReturnPolicy {
	if p.ReturnPolicy == nil {
		return models.ReturnPolicy{}
	}
	return models.ReturnPolicy{
		Returnable: p.ReturnPolicy.Returnable,
	}
}

// ValidateVariantAttributes validates that variants have proper attribute mappings
// For products with attribute definitions, all variants must have matching attributes
// For products without attributes (single variant), attribute_map can be empty {}
// For products without any variants, this is a no-op (default variant will be auto-created)
func (p *CreateProductInput) ValidateVariantAttributes() error {
	// Case 0: No variants provided - a default variant will be auto-created
	if len(p.ProductDetails.Variants) == 0 {
		// If no variants, there should be no attribute definitions either
		if len(p.ProductDetails.AttributeDefinitions) > 0 {
			return errors.New("attribute_definitions cannot be provided without variants")
		}
		return nil
	}

	// Case 1: No attribute definitions (single product/default variant)
	// In this case, all variants should have empty attribute_map
	if len(p.ProductDetails.AttributeDefinitions) == 0 {
		for i, variant := range p.ProductDetails.Variants {
			if len(variant.AttributeMap) > 0 {
				return fmt.Errorf("variant %d (%s) has attribute_map but no attributes are defined for this product", i+1, variant.SKU)
			}
		}
		return nil
	}

	// Case 2: Attribute definitions exist
	// Create a map of valid attribute keys
	validKeys := make(map[string][]string)
	for _, attr := range p.ProductDetails.AttributeDefinitions {
		validKeys[attr.Key] = attr.Options
	}

	// Validate each variant
	for i, variant := range p.ProductDetails.Variants {
		if len(variant.AttributeMap) != len(validKeys) {
			return fmt.Errorf("variant %d (%s) has %d attributes but %d are defined", i+1, variant.SKU, len(variant.AttributeMap), len(validKeys))
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

// ========================
// Helper Builder Methods
// ========================

// NewSingleVariantProduct creates a product with a single default variant (no attributes).
// Use this for products without variants like labels, caps, etc.
//
// Example: Creating a "Water Bottle Label" with just one SKU
func NewSingleVariantProduct(
	name string,
	defaultVariantSKU string,
	unit string,
	sellingPrice, costPrice float64,
) *CreateProductInput {
	return &CreateProductInput{
		Name: name,
		ProductDetails: ProductDetailsInput{
			Unit:                 unit,
			BaseSKU:              defaultVariantSKU,
			AttributeDefinitions: []ProductAttributeDefinitionInput{}, // No attributes
			Variants: []ProductVariantInput{
				{
					SKU:           defaultVariantSKU,
					VariantName:   name,
					AttributeMap:  map[string]string{}, // Empty for single variant
					SellingPrice:  sellingPrice,
					CostPrice:     costPrice,
					StockQuantity: 0,
					IsActive:      true,
				},
			},
		},
	}
}

// AddVariant adds a new variant to the product.
// The variant must match the attribute definitions already defined.
func (p *CreateProductInput) AddVariant(variant ProductVariantInput) {
	p.ProductDetails.Variants = append(p.ProductDetails.Variants, variant)
}

// AddAttribute defines a new attribute for this product's variants.
// All variants must have values for all defined attributes.
func (p *CreateProductInput) AddAttribute(key string, options []string) {
	p.ProductDetails.AttributeDefinitions = append(
		p.ProductDetails.AttributeDefinitions,
		ProductAttributeDefinitionInput{
			Key:     key,
			Options: options,
		},
	)
}

// GetVariantCount returns the number of variants for this product
func (p *CreateProductInput) GetVariantCount() int {
	return len(p.ProductDetails.Variants)
}

// HasVariantAttributes returns true if the product has variant attributes defined
func (p *CreateProductInput) HasVariantAttributes() bool {
	return len(p.ProductDetails.AttributeDefinitions) > 0
}

// ValidateSKUUniqueness checks that all variant SKUs are unique
func (p *CreateProductInput) ValidateSKUUniqueness() error {
	skus := make(map[string]bool)
	for i, variant := range p.ProductDetails.Variants {
		if _, exists := skus[variant.SKU]; exists {
			return fmt.Errorf("variant %d has duplicate SKU '%s'", i+1, variant.SKU)
		}
		skus[variant.SKU] = true
	}
	return nil
}
