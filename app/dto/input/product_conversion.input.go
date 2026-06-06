package input

import "time"

// CreateProductConversionInput for creating a new conversion rule
type CreateProductConversionInput struct {
	RawProductID       string  `json:"raw_product_id" validate:"required"`
	FinishedProductID  string  `json:"finished_product_id" validate:"required"`
	FinishedVariantSKU string  `json:"finished_variant_sku"`                      // Optional: If product has variants, specify SKU to add stock to that variant
	ConversionRatio    float64 `json:"conversion_ratio" validate:"required,gt=0"` // How many raw units for 1 finished unit
	LossPercentage     float64 `json:"loss_percentage" validate:"min=0,max=100"`
	IsActive           bool    `json:"is_active"`
	Notes              string  `json:"notes"`
}

// UpdateProductConversionInput for updating a conversion rule
type UpdateProductConversionInput struct {
	ConversionRatio    *float64 `json:"conversion_ratio" validate:"omitempty,gt=0"`
	LossPercentage     *float64 `json:"loss_percentage" validate:"omitempty,min=0,max=100"`
	IsActive           *bool    `json:"is_active"`
	FinishedVariantSKU *string  `json:"finished_variant_sku"` // Optional: Update which variant receives stock
	Notes              *string  `json:"notes"`
}

// CreateProductConversionRecordInput for performing a conversion
type CreateProductConversionRecordInput struct {
	ConversionID       string     `json:"conversion_id" validate:"required"`
	RawQuantityUsed    float64    `json:"raw_quantity_used" validate:"omitempty,gt=0"`
	ConversionDate     *time.Time `json:"conversion_date"`
	Notes              string     `json:"notes"`
	ExecuteConversion  bool       `json:"execute_conversion"`
	FinishedVariantSKU string     `json:"finished_variant_sku"`

	RawMaterialBags []UseRawMaterialBagInput `json:"raw_material_bags"`
}

// ListConversionsQuery for filtering conversions
type ListConversionsQuery struct {
	Page              int    `json:"page" validate:"min=1"`
	Limit             int    `json:"limit" validate:"min=1,max=100"`
	IsActive          *bool  `json:"is_active"`
	RawProductID      string `json:"raw_product_id"`
	FinishedProductID string `json:"finished_product_id"`
	Search            string `json:"search"`
}
