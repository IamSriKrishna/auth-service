package output

import "time"

// CustomerPricingDTO represents output for customer pricing
type CustomerPricingDTO struct {
	ID               string     `json:"id" example:"pricing-789"`
	CustomerID       string     `json:"customer_id" example:"cust-123"`
	CustomerName     string     `json:"customer_name" example:"Acme Corp"`
	ManufacturerID   string     `json:"manufacturer_id" example:"mfg-456"`
	ManufacturerRate float64    `json:"manufacturer_rate" example:"150.50"`
	EffectiveFrom    *time.Time `json:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to"`
	IsActive         bool       `json:"is_active" example:"true"`
	Notes            string     `json:"notes" example:"Special pricing"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        string     `json:"created_by" example:"user-001"`
	UpdatedBy        string     `json:"updated_by" example:"user-002"`
}

// CustomerPricingListDTO represents a list of customer pricing records
type CustomerPricingListDTO struct {
	Data  []CustomerPricingDTO `json:"data"`
	Total int64                `json:"total" example:"100"`
}

// CustomerPricingDetailDTO represents detailed customer pricing information
type CustomerPricingDetailDTO struct {
	ID               string     `json:"id"`
	CustomerID       string     `json:"customer_id"`
	CustomerName     string     `json:"customer_name"`
	ManufacturerID   string     `json:"manufacturer_id"`
	ManufacturerRate float64    `json:"manufacturer_rate"`
	EffectiveFrom    *time.Time `json:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to"`
	IsActive         bool       `json:"is_active"`
	Notes            string     `json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        string     `json:"created_by"`
	UpdatedBy        string     `json:"updated_by"`
}
