package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type ProductGroupOutput struct {
	ID           string                        `json:"id"`
	Name         string                        `json:"name"`
	Description  string                        `json:"description"`
	IsActive     bool                          `json:"is_active"`
	Cost         float64                       `json:"cost"`
	SellingPrice float64                       `json:"selling_price"`
	Profit       float64                       `json:"profit"`
	Components   []ProductGroupComponentOutput `json:"components"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

type ProductGroupComponentOutput struct {
	ID             uint              `json:"id"`
	ProductID      string            `json:"product_id"`
	Product        *ProductOutput    `json:"product,omitempty"`
	VariantSku     *string           `json:"variant_sku,omitempty"`
	Quantity       float64           `json:"quantity"`
	Position       int               `json:"position,omitempty"`
	VariantDetails map[string]string `json:"variant_details,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type ProductGroupListOutput struct {
	ProductGroups []ProductGroupOutput `json:"product_groups"`
	Total         int64                `json:"total"`
}

func ToProductGroupOutput(pg *models.ProductGroup) (*ProductGroupOutput, error) {
	// First, consolidate duplicate components (keep only the latest one per product_id:variant_sku)
	consolidatedMap := make(map[string]*models.ProductGroupComponent)

	for i := range pg.Components {
		variantSku := ""
		if pg.Components[i].VariantSku != nil {
			variantSku = *pg.Components[i].VariantSku
		}
		key := pg.Components[i].ProductID + ":" + variantSku

		// Keep the component with the latest UpdatedAt timestamp
		if existing, exists := consolidatedMap[key]; !exists {
			consolidatedMap[key] = &pg.Components[i]
		} else if pg.Components[i].UpdatedAt.After(existing.UpdatedAt) {
			consolidatedMap[key] = &pg.Components[i]
		}
	}

	components := make([]ProductGroupComponentOutput, len(consolidatedMap))
	totalCost := 0.0
	totalSelling := 0.0
	idx := 0

	for _, comp := range consolidatedMap {
		componentOutput := ProductGroupComponentOutput{
			ID:         comp.ID,
			ProductID:  comp.ProductID,
			VariantSku: comp.VariantSku,
			Quantity:   comp.Quantity,
			Position:   comp.Position,
			CreatedAt:  comp.CreatedAt,
			UpdatedAt:  comp.UpdatedAt,
		}

		// Convert VariantDetails to output format
		if len(comp.VariantDetails) > 0 {
			componentOutput.VariantDetails = make(map[string]string)
			for k, v := range comp.VariantDetails {
				componentOutput.VariantDetails[k] = v
			}
		}

		// Include product details if available
		if comp.Product != nil {
			productOutput, err := ToProductOutput(comp.Product)
			if err == nil {
				componentOutput.Product = productOutput

				// Calculate cost and selling price for this component
				// Get the cost and selling price from the product's sales and purchase info
				costPerUnit := comp.Product.PurchaseInfo.CostPrice
				sellingPerUnit := comp.Product.SalesInfo.SellingPrice

				totalCost += costPerUnit * comp.Quantity
				totalSelling += sellingPerUnit * comp.Quantity
			}
		}

		components[idx] = componentOutput
		idx++
	}

	profit := totalSelling - totalCost

	return &ProductGroupOutput{
		ID:           pg.ID,
		Name:         pg.Name,
		Description:  pg.Description,
		IsActive:     pg.IsActive,
		Cost:         totalCost,
		SellingPrice: totalSelling,
		Profit:       profit,
		Components:   components,
		CreatedAt:    pg.CreatedAt,
		UpdatedAt:    pg.UpdatedAt,
	}, nil
}

func ToProductGroupListOutput(pgs []models.ProductGroup, total int64) (*ProductGroupListOutput, error) {
	outputs := make([]ProductGroupOutput, len(pgs))
	for i, pg := range pgs {
		output, err := ToProductGroupOutput(&pg)
		if err != nil {
			return nil, err
		}
		outputs[i] = *output
	}

	return &ProductGroupListOutput{
		ProductGroups: outputs,
		Total:         total,
	}, nil
}

type ReorderProductGroupOutput struct {
	ID             string                `json:"id"`
	ReorderSummary *ReorderSummaryOutput `json:"reorder_summary"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

type ReorderSummaryOutput struct {
	TotalProducts int                   `json:"total_products"`
	Updates       []ReorderUpdateOutput `json:"updates"`
}

type ReorderUpdateOutput struct {
	VariantSku    string  `json:"variant_sku"`
	OldQuantity   float64 `json:"old_quantity"`
	NewQuantity   float64 `json:"new_quantity"`
	StockAdjusted float64 `json:"stock_adjusted"`
}
