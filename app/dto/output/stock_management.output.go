package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

// VariantStockOutput represents variant stock information in API responses
type VariantStockOutput struct {
	ID                string            `json:"id"`
	ProductID         string            `json:"product_id"`
	ProductName       string            `json:"product_name"`
	VariantSKU        string            `json:"variant_sku"`
	VariantName       string            `json:"variant_name,omitempty"`
	VariantAttributes map[string]string `json:"variant_attributes,omitempty"`
	CurrentStock      float64           `json:"current_stock"`
	PurchasedTotal    float64           `json:"purchased_total"`
	SoldTotal         float64           `json:"sold_total"`
	ReservedStock     float64           `json:"reserved_stock"`
	AvailableStock    float64           `json:"available_stock"`
	InTransitStock    float64           `json:"in_transit_stock"`
	AverageCost       float64           `json:"average_cost"`
	SellingPrice      float64           `json:"selling_price"`
	StockValue        float64           `json:"stock_value"`
	ReorderLevel      float64           `json:"reorder_level"`
	IsLowStock        bool              `json:"is_low_stock"`
	LastPurchasedDate *time.Time        `json:"last_purchased_date,omitempty"`
	LastSoldDate      *time.Time        `json:"last_sold_date,omitempty"`
	LastStockSyncAt   *time.Time        `json:"last_stock_sync_at,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// ProductStockOutput represents product stock information in API responses
type ProductStockOutput struct {
	ID                string     `json:"id"`
	ProductID         string     `json:"product_id"`
	ProductName       string     `json:"product_name"`
	SKU               string     `json:"sku"`
	CurrentStock      float64    `json:"current_stock"`
	PurchasedStock    float64    `json:"purchased_stock"`
	SoldStock         float64    `json:"sold_stock"`
	ReservedStock     float64    `json:"reserved_stock"`
	AvailableStock    float64    `json:"available_stock"`
	AverageCost       float64    `json:"average_cost"`
	StockValue        float64    `json:"stock_value"`
	LastPurchasedDate *time.Time `json:"last_purchased_date,omitempty"`
	LastSoldDate      *time.Time `json:"last_sold_date,omitempty"`
	LastStockSyncAt   time.Time  `json:"last_stock_sync_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	// Variant information if product has variants
	Variants []VariantStockOutput `json:"variants,omitempty"`
}

// StockManagementResponse represents the complete stock management response
type StockManagementResponse struct {
	Stocks          []interface{} `json:"stocks"` // Can be ProductStockOutput or VariantStockOutput
	TotalStockValue float64       `json:"total_stock_value"`
}

// StockLedgerOutput represents a stock movement entry in API responses
type StockLedgerOutput struct {
	ID               uint      `json:"id"`
	ProductID        string    `json:"product_id"`
	ProductName      string    `json:"product_name,omitempty"`
	MovementType     string    `json:"movement_type"`
	Quantity         float64   `json:"quantity"`
	Rate             float64   `json:"rate"`
	Amount           float64   `json:"amount"`
	ReferenceType    string    `json:"reference_type"`
	ReferenceID      string    `json:"reference_id"`
	ReferenceNumber  string    `json:"reference_number"`
	BalanceBeforeQty float64   `json:"balance_before_qty"`
	BalanceAfterQty  float64   `json:"balance_after_qty"`
	CostBeforeAmount float64   `json:"cost_before_amount"`
	CostAfterAmount  float64   `json:"cost_after_amount"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
}

// ToVariantStockOutput converts a variant stock model to output DTO
func ToVariantStockOutput(stock *models.VariantStock) *VariantStockOutput {
	variantAttributes := make(map[string]string)
	// Parse variant attributes if they exist in the model
	// This would depend on how attributes are stored in your VariantStock model

	stockValue := stock.CurrentStock * stock.AverageCost
	if stockValue < 0 {
		stockValue = 0
	}

	return &VariantStockOutput{
		ID:                stock.ID,
		ProductID:         stock.ProductID,
		ProductName:       stock.ProductName,
		VariantSKU:        stock.VariantSKU,
		VariantName:       stock.VariantName,
		VariantAttributes: variantAttributes,
		CurrentStock:      stock.CurrentStock,
		PurchasedTotal:    stock.PurchasedStock,
		SoldTotal:         stock.SoldStock,
		ReservedStock:     stock.ReservedStock,
		AvailableStock:    stock.AvailableStock,
		InTransitStock:    stock.InTransitStock,
		AverageCost:       stock.AverageCost,
		SellingPrice:      stock.SellingPrice,
		StockValue:        stockValue,
		ReorderLevel:      stock.ReorderLevel,
		IsLowStock:        stock.IsLowStock,
		LastPurchasedDate: stock.LastPurchasedDate,
		LastSoldDate:      stock.LastSoldDate,
		LastStockSyncAt:   stock.LastStockSyncAt,
		CreatedAt:         stock.CreatedAt,
		UpdatedAt:         stock.UpdatedAt,
	}
}

// ToProductStockOutput converts a model to output DTO
// If product has variants, includes them in the response
func ToProductStockOutput(stock *models.ProductStock) *ProductStockOutput {
	return &ProductStockOutput{
		ID:                stock.ID,
		ProductID:         stock.ProductID,
		ProductName:       stock.ProductName,
		SKU:               stock.SKU,
		CurrentStock:      stock.CurrentStock,
		PurchasedStock:    stock.PurchasedStock,
		SoldStock:         stock.SoldStock,
		ReservedStock:     stock.ReservedStock,
		AvailableStock:    stock.AvailableStock,
		AverageCost:       stock.AverageCost,
		StockValue:        stock.CurrentStock * stock.AverageCost,
		LastPurchasedDate: stock.LastPurchasedDate,
		LastSoldDate:      stock.LastSoldDate,
		LastStockSyncAt:   stock.LastStockSyncAt,
		CreatedAt:         stock.CreatedAt,
		UpdatedAt:         stock.UpdatedAt,
	}
}

// ToStockLedgerOutput converts a ledger entry to output DTO
func ToStockLedgerOutput(ledger *models.StockLedger) *StockLedgerOutput {
	productName := ""
	if ledger.Product != nil {
		productName = ledger.Product.Name
	}

	return &StockLedgerOutput{
		ID:               ledger.ID,
		ProductID:        ledger.ProductID,
		ProductName:      productName,
		MovementType:     ledger.MovementType,
		Quantity:         ledger.Quantity,
		Rate:             ledger.Rate,
		Amount:           ledger.Amount,
		ReferenceType:    ledger.ReferenceType,
		ReferenceID:      ledger.ReferenceID,
		ReferenceNumber:  ledger.ReferenceNumber,
		BalanceBeforeQty: ledger.BalanceBeforeQty,
		BalanceAfterQty:  ledger.BalanceAfterQty,
		CostBeforeAmount: ledger.CostBeforeAmount,
		CostAfterAmount:  ledger.CostAfterAmount,
		Notes:            ledger.Notes,
		CreatedAt:        ledger.CreatedAt,
		CreatedBy:        ledger.CreatedBy,
	}
}

// BuildVariantStockResponse builds a stock management response for variant stocks
// Shows variant details if product has variants, otherwise shows product information
func BuildVariantStockResponse(variantStocks []models.VariantStock) *StockManagementResponse {
	response := &StockManagementResponse{
		Stocks: make([]interface{}, 0),
	}

	totalValue := 0.0
	for i := range variantStocks {
		variantOutput := ToVariantStockOutput(&variantStocks[i])
		response.Stocks = append(response.Stocks, variantOutput)
		totalValue += variantOutput.StockValue
	}

	response.TotalStockValue = totalValue
	return response
}

// BuildProductStockResponse builds a stock management response for product stocks
// If a product has variants, includes variant details in the response
func BuildProductStockResponse(productStocks []models.ProductStock, variantsByProduct map[string][]models.VariantStock) *StockManagementResponse {
	response := &StockManagementResponse{
		Stocks: make([]interface{}, 0),
	}

	totalValue := 0.0
	for i := range productStocks {
		productOutput := ToProductStockOutput(&productStocks[i])

		// If product has variants, include them
		if variants, hasVariants := variantsByProduct[productStocks[i].ProductID]; hasVariants && len(variants) > 0 {
			productOutput.Variants = make([]VariantStockOutput, 0)
			for j := range variants {
				variantOutput := ToVariantStockOutput(&variants[j])
				productOutput.Variants = append(productOutput.Variants, *variantOutput)
				totalValue += variantOutput.StockValue
			}
		} else {
			// No variants, just add product stock value
			totalValue += productOutput.StockValue
		}

		response.Stocks = append(response.Stocks, productOutput)
	}

	response.TotalStockValue = totalValue
	return response
}

// BuildCombinedStockResponse builds a combined stock response showing either variants or products
// Perfect for stock management - shows variants if available, product alone if no variants
func BuildCombinedStockResponse(variantStocks []models.VariantStock, productStocks []models.ProductStock) *StockManagementResponse {
	response := &StockManagementResponse{
		Stocks: make([]interface{}, 0),
	}

	totalValue := 0.0

	// Add variant stocks (if variants exist, show them)
	processedProductIDs := make(map[string]bool)
	for i := range variantStocks {
		variantOutput := ToVariantStockOutput(&variantStocks[i])
		response.Stocks = append(response.Stocks, variantOutput)
		totalValue += variantOutput.StockValue
		processedProductIDs[variantStocks[i].ProductID] = true
	}

	// Add product stocks that don't have variants
	for i := range productStocks {
		if !processedProductIDs[productStocks[i].ProductID] {
			productOutput := ToProductStockOutput(&productStocks[i])
			response.Stocks = append(response.Stocks, productOutput)
			totalValue += productOutput.StockValue
		}
	}

	response.TotalStockValue = totalValue
	return response
}
