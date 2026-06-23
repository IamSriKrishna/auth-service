package output

import "time"

// ProductConversionOutput for conversion rule response
type ProductConversionOutput struct {
	ID                   string    `json:"id"`
	RawProductID         string    `json:"raw_product_id"`
	RawProductName       string    `json:"raw_product_name"`
	RawProductSpec       string    `json:"raw_product_spec"`
	FinishedProductID    string    `json:"finished_product_id"`
	FinishedProductName  string    `json:"finished_product_name"`
	FinishedProductSpec  string    `json:"finished_product_spec"`
	FinishedVariantSKU   string    `json:"finished_variant_sku"` // Optional: Variant SKU to receive converted stock
	ConversionRatio      float64   `json:"conversion_ratio"`
	LossPercentage       float64   `json:"loss_percentage"`
	IsActive             bool      `json:"is_active"`
	Notes                string    `json:"notes"`
	CreatedBy            string    `json:"created_by"`
	CreatedByUserName    string    `json:"created_by_user_name"`
	CreatedByCompanyID   uint      `json:"created_by_company_id"`
	CreatedByCompanyName string    `json:"created_by_company_name"`
	UpdatedBy            string    `json:"updated_by"`
	UpdatedByUserName    string    `json:"updated_by_user_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ProductConversionRecordOutput for conversion record response
type ProductConversionRecordOutput struct {
	ID                       string                           `json:"id"`
	ConversionID             string                           `json:"conversion_id"`
	RawProductID             string                           `json:"raw_product_id"`
	RawProductName           string                           `json:"raw_product_name"`
	RawQuantityUsed          float64                          `json:"raw_quantity_used"`
	FinishedProductID        string                           `json:"finished_product_id"`
	FinishedProductName      string                           `json:"finished_product_name"`
	FinishedVariantSKU       string                           `json:"finished_variant_sku"` // Variant SKU that received the stock (if applicable)
	FinishedQuantityProduced float64                          `json:"finished_quantity_produced"`
	LossQuantity             float64                          `json:"loss_quantity"`
	ConversionDate           time.Time                        `json:"conversion_date"`
	Status                   string                           `json:"status"`
	Notes                    string                           `json:"notes"`
	CreatedBy                string                           `json:"created_by"`
	CreatedByUserName        string                           `json:"created_by_user_name"`
	CreatedByCompanyID       uint                             `json:"created_by_company_id"`
	CreatedByCompanyName     string                           `json:"created_by_company_name"`
	BagsUsed                 []ConversionRecordBagUsageOutput `json:"bags_used"` // Bags used for this conversion
	CreatedAt                time.Time                        `json:"created_at"`
}

// ConversionRecordBagUsageOutput for bag usage details
type ConversionRecordBagUsageOutput struct {
	ID             string    `json:"id"`
	BagID          string    `json:"bag_id"`
	BagNumber      int       `json:"bag_number"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	QuantityUsedKg float64   `json:"quantity_used_kg"`
	CreatedAt      time.Time `json:"created_at"`
}

// ProductConversionListOutput for list of conversions
type ProductConversionListOutput struct {
	Conversions []ProductConversionOutput `json:"conversions"`
	Total       int64                     `json:"total"`
	Page        int                       `json:"page"`
	Limit       int                       `json:"limit"`
}

// ProductConversionRecordListOutput for list of conversion records
type ProductConversionRecordListOutput struct {
	Records []ProductConversionRecordOutput `json:"records"`
	Total   int64                           `json:"total"`
	Page    int                             `json:"page"`
	Limit   int                             `json:"limit"`
}

// ConversionExecutionOutput for conversion execution response
type ConversionExecutionOutput struct {
	RecordID                 string  `json:"record_id"`
	Status                   string  `json:"status"`
	RawProductName           string  `json:"raw_product_name"`
	RawQuantityUsed          float64 `json:"raw_quantity_used"`
	FinishedProductName      string  `json:"finished_product_name"`
	FinishedVariantSKU       string  `json:"finished_variant_sku"` // Variant that received stock (if applicable)
	FinishedQuantityProduced float64 `json:"finished_quantity_produced"`
	LossQuantity             float64 `json:"loss_quantity"`
	Message                  string  `json:"message"`
}
