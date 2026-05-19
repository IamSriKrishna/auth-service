package input

// CreateProductGroupInput represents the input for creating a product group/bundle.
// A product group combines multiple individual products and their specific variants
// into a single bundle/kit/combo that can be sold together.
//
// Example: "500ml Water Bottle Kit" = 1x Bottle (specific color) + 1x Cap (specific size) + 1x Label
type CreateProductGroupInput struct {
	// Name of the product group
	Name string `json:"name" validate:"required" example:"500ml Water Bottle Complete Kit"`
	// Description of what this product group contains
	Description string `json:"description" example:"Complete water bottle kit with bottle, cap, and label"`
	// Status of the product group (e.g., "active", "inactive")
	Status string `json:"status,omitempty" example:"active"`
	// Whether this product group is active
	IsActive bool `json:"is_active" example:"true"`
	// Products/components of the product group - the products that make up this bundle
	Products []ProductGroupComponentInput `json:"products" validate:"required,min=1,dive,required"`
}

// ProductGroupComponentInput represents a single product or product variant
// that is part of the product group.
//
// Each component references:
// - A specific product (by product_id)
// - Optionally a specific variant of that product (by variant_sku)
// - The quantity of this component per product group
type ProductGroupComponentInput struct {
	// ID of the product to include in this group
	ProductID string `json:"product_id" validate:"required" example:"prod_abc123"`
	// Optional: SKU of the specific variant to use
	// If not specified, the product group can use any variant of the product
	VariantSku *string `json:"variant_sku,omitempty" example:"WTR-BOT-500-RED"`
	// Quantity of this component per product group
	// All components must have the same quantity
	Quantity float64 `json:"quantity" validate:"required,gt=0" example:"1"`
	// Optional: Position/order of this component in the product group
	Position int `json:"position,omitempty" example:"1"`
	// Optional: Metadata about the specific variant selected
	// Useful for storing variant attribute details
	VariantDetails interface{} `json:"variant_details,omitempty"`
}

// UpdateProductGroupInput represents updates to an existing product group
type UpdateProductGroupInput struct {
	// Name of the product group
	Name string `json:"name"`
	// Description of what this product group contains
	Description string `json:"description"`
	// Status of the product group
	Status string `json:"status"`
	// Whether this product group is active
	IsActive *bool `json:"is_active"`
	// Products/components of the product group
	Products []ProductGroupComponentInput `json:"products"`
}
