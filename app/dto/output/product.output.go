package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type ProductOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Resource fields
	IsResource          bool    `json:"is_resource"`
	ResourceName        string  `json:"resource_name,omitempty"`
	ResourceUnit        string  `json:"resource_unit,omitempty"`
	ResourceCostPerUnit float64 `json:"resource_cost_per_unit,omitempty"`

	ProductDetails ProductDetailsOutput `json:"product_details,omitempty"`
	SalesInfo      SalesInfoOutput      `json:"sales_info,omitempty"`
	PurchaseInfo   PurchaseInfoOutput   `json:"purchase_info,omitempty"`
	Inventory      InventoryOutput      `json:"inventory,omitempty"`
	ReturnPolicy   ReturnPolicyOutput   `json:"return_policy,omitempty"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id,omitempty"`
	UserName    string    `json:"user_name,omitempty"`
	CompanyID   uint      `json:"company_id,omitempty"`
	CompanyName string    `json:"company_name,omitempty"`
}

type ProductAttributeDefinitionOutput struct {
	Key     string   `json:"key"`
	Options []string `json:"options"`
}

type ProductDetailsOutput struct {
	Unit                 string                             `json:"unit"`
	BaseSKU              string                             `json:"base_sku,omitempty"`
	UPC                  string                             `json:"upc,omitempty"`
	EAN                  string                             `json:"ean,omitempty"`
	MPN                  string                             `json:"mpn,omitempty"`
	ISBN                 string                             `json:"isbn,omitempty"`
	Description          string                             `json:"description,omitempty"`
	AttributeDefinitions []ProductAttributeDefinitionOutput `json:"attribute_definitions,omitempty"`
	Variants             []ProductVariantOutput             `json:"variants,omitempty"`
}

type ProductVariantOutput struct {
	SKU           string            `json:"sku"`
	VariantName   string            `json:"variant_name,omitempty"`
	AttributeMap  map[string]string `json:"attribute_map"`
	SellingPrice  float64           `json:"selling_price"`
	CostPrice     float64           `json:"cost_price"`
	StockQuantity float64           `json:"stock_quantity"`
	ReorderLevel  float64           `json:"reorder_level"`
	IsActive      bool              `json:"is_active"`
}

type ProductListOutput struct {
	Products []ProductOutput `json:"products"`
	Total    int64           `json:"total"`
}

type ProductOpeningStockOutput struct {
	OpeningStock            float64   `json:"opening_stock"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type ProductVariantOpeningStockOutput struct {
	VariantSKU              string    `json:"variant_sku"`
	OpeningStock            float64   `json:"opening_stock"`
	OpeningStockRatePerUnit float64   `json:"opening_stock_rate_per_unit"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// Conversion functions

func ToProductOutput(product *models.Product) (*ProductOutput, error) {
	output := &ProductOutput{
		ID:                  product.ID,
		Name:                product.Name,
		IsResource:          product.IsResource,
		ResourceName:        product.ResourceName,
		ResourceUnit:        product.ResourceUnit,
		ResourceCostPerUnit: product.ResourceCostPerUnit,
		CreatedAt:           product.CreatedAt,
		UpdatedAt:           product.UpdatedAt,
		UserID:              product.CreatedBy,
		UserName:            product.CreatedByUserName,
		CompanyID:           product.CreatedByCompanyID,
		CompanyName:         product.CreatedByCompanyName,
	}

	// Only populate product details, sales info, purchase info, inventory, and return policy for non-resource products
	if !product.IsResource {
		attributeDefs := make([]ProductAttributeDefinitionOutput, len(product.ProductDetails.AttributeDefinitions))
		for i, def := range product.ProductDetails.AttributeDefinitions {
			attributeDefs[i] = ProductAttributeDefinitionOutput{
				Key:     def.Key,
				Options: def.Options,
			}
		}

		variants := make([]ProductVariantOutput, len(product.ProductDetails.ProductVariants))
		for i, v := range product.ProductDetails.ProductVariants {
			attributeMap := make(map[string]string)
			for _, attr := range v.Attributes {
				attributeMap[attr.Key] = attr.Value
			}

			variants[i] = ProductVariantOutput{
				SKU:           v.SKU,
				VariantName:   v.VariantName,
				AttributeMap:  attributeMap,
				SellingPrice:  v.SellingPrice,
				CostPrice:     v.CostPrice,
				StockQuantity: v.StockQuantity,
				ReorderLevel:  v.ReorderLevel,
				IsActive:      v.IsActive,
			}
		}

		purchaseInfo := PurchaseInfoOutput{
			Account:     product.PurchaseInfo.Account,
			CostPrice:   product.PurchaseInfo.CostPrice,
			Currency:    product.PurchaseInfo.Currency,
			Description: product.PurchaseInfo.Description,
		}

		productDetails := ProductDetailsOutput{
			Unit:                 product.ProductDetails.Unit,
			BaseSKU:              product.ProductDetails.BaseSKU,
			UPC:                  product.ProductDetails.UPC,
			EAN:                  product.ProductDetails.EAN,
			MPN:                  product.ProductDetails.MPN,
			ISBN:                 product.ProductDetails.ISBN,
			Description:          product.ProductDetails.Description,
			AttributeDefinitions: attributeDefs,
			Variants:             variants,
		}

		output.ProductDetails = productDetails
		output.SalesInfo = SalesInfoOutput{
			Account:      product.SalesInfo.Account,
			SellingPrice: product.SalesInfo.SellingPrice,
			Currency:     product.SalesInfo.Currency,
			Description:  product.SalesInfo.Description,
		}
		output.PurchaseInfo = purchaseInfo
		output.Inventory = InventoryOutput{
			TrackInventory:           product.Inventory.TrackInventory,
			InventoryAccount:         product.Inventory.InventoryAccount,
			InventoryValuationMethod: product.Inventory.InventoryValuationMethod,
			ReorderPoint:             product.Inventory.ReorderPoint,
		}
		output.ReturnPolicy = ReturnPolicyOutput{
			Returnable: product.ReturnPolicy.Returnable,
		}
	}

	return output, nil
}

func ToProductListOutput(products []models.Product, total int64) (*ProductListOutput, error) {
	outputs := make([]ProductOutput, len(products))
	for i, product := range products {
		output, err := ToProductOutput(&product)
		if err != nil {
			return nil, err
		}
		outputs[i] = *output
	}

	return &ProductListOutput{
		Products: outputs,
		Total:    total,
	}, nil
}
