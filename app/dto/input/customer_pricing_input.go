package input

import "time"

// CustomerPricingLineItem represents a pricing line item tied to a product.
type CustomerPricingLineItem struct {
	ProductID   string  `json:"product_id" example:"prod_abc123"`
	ProductName string  `json:"product_name" example:"500ml Bottle"`
	Rate        float64 `json:"rate" binding:"required" example:"100.00"`
	Account     string  `json:"account" example:"SALES_REVENUE"`
	Description string  `json:"description" example:"Wholesale pricing"`
}

// CreateCustomerPricingDTO represents input for creating customer pricing with line items
type CreateCustomerPricingDTO struct {
	CustomerID uint                      `json:"customer_id" binding:"required" example:"1"`
	LineItems  []CustomerPricingLineItem `json:"line_items" binding:"required,dive" example:"[]"`
}

// UpdateCustomerPricingDTO represents input for updating customer pricing
type UpdateCustomerPricingDTO struct {
	Rate        float64 `json:"rate" binding:"required" example:"110.00"`
	Account     string  `json:"account" example:"SALES_REVENUE"`
	Description string  `json:"description" example:"Updated pricing"`
	IsActive    bool    `json:"is_active" example:"true"`
}

// SetDateRangeDTO represents input for setting effective date range
type SetDateRangeDTO struct {
	EffectiveFrom *time.Time `json:"effective_from" example:"2026-01-01T00:00:00Z"`
	EffectiveTo   *time.Time `json:"effective_to" example:"2026-12-31T23:59:59Z"`
}

// FilterCustomerPricingDTO represents input for filtering customer pricing
type FilterCustomerPricingDTO struct {
	CustomerID uint  `json:"customer_id" example:"1"`
	IsActive   *bool `json:"is_active" example:"true"`
	Offset     int   `json:"offset" example:"0"`
	Limit      int   `json:"limit" example:"10"`
}
