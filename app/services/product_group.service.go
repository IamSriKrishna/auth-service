package services

import (
	"fmt"
	"log"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ProductGroupService interface {
	// Basic CRUD Operations for bundled products
	Create(input *input.CreateProductGroupInput) (*output.ProductGroupOutput, error)
	FindByID(id string) (*output.ProductGroupOutput, error)
	FindAll(limit, offset int, search string) (*output.ProductGroupListOutput, error)
	Update(id string, input *input.UpdateProductGroupInput) (*output.ProductGroupOutput, error)
	Delete(id string) error

	// Extra operations
	FindByName(name string) (*output.ProductGroupOutput, error)
	ValidateStockAvailability(productGroupID string, quantityToCreate float64) ([]string, error)
}

type productGroupService struct {
	productGroupRepo             repo.ProductGroupRepository
	productRepo                  repo.ProductRepository
	variantStockMgmtService      VariantStockManagementService
	productGroupInventoryService ProductGroupInventoryService
	stockManagementService       StockManagementService
	productStockRepo             repo.ProductStockRepository
}

func NewProductGroupService(
	productGroupRepo repo.ProductGroupRepository,
	productRepo repo.ProductRepository,
) ProductGroupService {
	return &productGroupService{
		productGroupRepo: productGroupRepo,
		productRepo:      productRepo,
	}
}

// NewProductGroupServiceWithStockMgmt creates a new ProductGroupService with inventory management
func NewProductGroupServiceWithStockMgmt(
	productGroupRepo repo.ProductGroupRepository,
	productRepo repo.ProductRepository,
	variantStockMgmtService VariantStockManagementService,
	productGroupInventoryService ProductGroupInventoryService,
	stockManagementService StockManagementService,
	productStockRepo repo.ProductStockRepository,
) ProductGroupService {
	return &productGroupService{
		productGroupRepo:             productGroupRepo,
		productRepo:                  productRepo,
		variantStockMgmtService:      variantStockMgmtService,
		productGroupInventoryService: productGroupInventoryService,
		stockManagementService:       stockManagementService,
		productStockRepo:             productStockRepo,
	}
}

func (s *productGroupService) Create(input *input.CreateProductGroupInput) (*output.ProductGroupOutput, error) {
	// Validate all products exist and quantities are valid
	if len(input.Products) == 0 {
		return nil, fmt.Errorf("product group must have at least one product")
	}

	// Check that all component quantities are equal and whole numbers (integers)
	firstQuantity := input.Products[0].Quantity
	for i, comp := range input.Products {
		// Validate all quantities are whole numbers
		if comp.Quantity != float64(int64(comp.Quantity)) {
			return nil, fmt.Errorf("product %d quantity must be a whole number (no decimals). Got: %f", i+1, comp.Quantity)
		}

		if comp.Quantity != firstQuantity {
			return nil, fmt.Errorf("all product quantities must be equal. Product 1 has quantity %d, but product %d has quantity %d", int64(firstQuantity), i+1, int64(comp.Quantity))
		}

		if comp.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for product %s: got %d", comp.ProductID, int64(comp.Quantity))
		}

		product, err := s.productRepo.FindByID(comp.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", comp.ProductID)
		}

		// If variant is specified, validate it exists
		if comp.VariantSku != nil {
			variantFound := false
			if product.ProductDetails.ProductVariants != nil {
				for _, v := range product.ProductDetails.ProductVariants {
					if v.SKU == *comp.VariantSku {
						variantFound = true
						break
					}
				}
			}
			if !variantFound {
				return nil, fmt.Errorf("variant %s not found in product %s", *comp.VariantSku, comp.ProductID)
			}
		}
	}

	// Generate ProductGroup ID before stock deduction (in case we need it for reference)
	productGroupID := "pg_" + uuid.New().String()[:8]

	// ====== STOCK VALIDATION AND DEDUCTION ======
	// Check and deduct stock from component variants
	if s.variantStockMgmtService == nil {
		return nil, fmt.Errorf("stock management service not initialized - cannot create product group")
	}

	if err := s.checkAndDeductStock(productGroupID, input.Products); err != nil {
		return nil, fmt.Errorf("stock validation failed: %w", err)
	}

	// Create components
	components := make([]models.ProductGroupComponent, 0, len(input.Products))
	for _, comp := range input.Products {
		// Convert VariantDetails from interface{} to models.VariantDetails
		var variantDetails models.VariantDetails
		if comp.VariantDetails != nil {
			switch v := comp.VariantDetails.(type) {
			case map[string]interface{}:
				variantDetails = make(models.VariantDetails)
				for k, val := range v {
					if strVal, ok := val.(string); ok {
						variantDetails[k] = strVal
					}
				}
			case map[string]string:
				variantDetails = models.VariantDetails(v)
			}
		}

		component := models.ProductGroupComponent{
			ProductGroupID: productGroupID,
			ProductID:      comp.ProductID,
			VariantSku:     comp.VariantSku,
			Quantity:       comp.Quantity,
			Position:       comp.Position,
			VariantDetails: variantDetails,
		}
		components = append(components, component)
	}

	productGroup := &models.ProductGroup{
		ID:          productGroupID,
		Name:        input.Name,
		Description: input.Description,
		IsActive:    input.IsActive,
		Components:  components,
	}

	err := s.productGroupRepo.Create(productGroup)
	if err != nil {
		return nil, err
	}

	// Create stock entry for the product group itself
	if err := s.createProductGroupStock(productGroupID, input); err != nil {
		log.Printf("[PG_STOCK] Warning: Failed to create product group stock entry: %v", err)
		// Don't fail - product group is already created
	}

	return s.toOutput(productGroup)
}

func (s *productGroupService) FindByID(id string) (*output.ProductGroupOutput, error) {
	productGroup, err := s.productGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("product group not found")
	}

	return s.toOutput(productGroup)
}

func (s *productGroupService) FindAll(limit, offset int, search string) (*output.ProductGroupListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	productGroups, total, err := s.productGroupRepo.FindAll(limit, offset, search)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.ProductGroupOutput, len(productGroups))
	for i, pg := range productGroups {
		out, _ := s.toOutput(&pg)
		outputs[i] = *out
	}

	return &output.ProductGroupListOutput{
		ProductGroups: outputs,
		Total:         total,
	}, nil
}

func (s *productGroupService) Update(id string, input *input.UpdateProductGroupInput) (*output.ProductGroupOutput, error) {
	productGroup, err := s.productGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("product group not found")
	}

	if input.Name != "" {
		productGroup.Name = input.Name
	}

	if input.Description != "" {
		productGroup.Description = input.Description
	}

	if input.IsActive != nil {
		productGroup.IsActive = *input.IsActive
	}

	// Update products if provided
	if len(input.Products) > 0 {
		// Validate all products exist
		for _, comp := range input.Products {
			product, err := s.productRepo.FindByID(comp.ProductID)
			if err != nil {
				return nil, fmt.Errorf("product %s not found", comp.ProductID)
			}

			if comp.VariantSku != nil {
				variantFound := false
				if product.ProductDetails.ProductVariants != nil {
					for _, v := range product.ProductDetails.ProductVariants {
						if v.SKU == *comp.VariantSku {
							variantFound = true
							break
						}
					}
				}
				if !variantFound {
					return nil, fmt.Errorf("variant %s not found in product %s", *comp.VariantSku, comp.ProductID)
				}
			}
		}

		productGroup.Components = make([]models.ProductGroupComponent, 0, len(input.Products))
		for _, comp := range input.Products {
			// Convert VariantDetails from interface{} to models.VariantDetails
			var variantDetails models.VariantDetails
			if comp.VariantDetails != nil {
				switch v := comp.VariantDetails.(type) {
				case map[string]interface{}:
					variantDetails = make(models.VariantDetails)
					for k, val := range v {
						if strVal, ok := val.(string); ok {
							variantDetails[k] = strVal
						}
					}
				case map[string]string:
					variantDetails = models.VariantDetails(v)
				}
			}

			component := models.ProductGroupComponent{
				ProductGroupID: id,
				ProductID:      comp.ProductID,
				VariantSku:     comp.VariantSku,
				Quantity:       comp.Quantity,
				Position:       comp.Position,
				VariantDetails: variantDetails,
			}
			productGroup.Components = append(productGroup.Components, component)
		}
	}

	err = s.productGroupRepo.Update(productGroup)
	if err != nil {
		return nil, err
	}

	return s.toOutput(productGroup)
}

func (s *productGroupService) Delete(id string) error {
	return s.productGroupRepo.Delete(id)
}

func (s *productGroupService) FindByName(name string) (*output.ProductGroupOutput, error) {
	productGroup, err := s.productGroupRepo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("product group not found")
	}

	return s.toOutput(productGroup)
}

// checkAndDeductStock checks availability and deducts stock for all components
// This is called during product group creation to ensure components are reserved
func (s *productGroupService) checkAndDeductStock(productGroupID string, components []input.ProductGroupComponentInput) error {
	// First, check all components have sufficient stock available
	for _, comp := range components {
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			log.Printf("[PG_STOCK] Warning: Variant SKU not specified for product %s, skipping stock check", comp.ProductID)
			continue
		}

		// Get current stock for variant
		variantStock, err := s.variantStockMgmtService.GetVariantStockSummary(*comp.VariantSku)
		if err != nil {
			return fmt.Errorf("variant %s stock not found for product %s", *comp.VariantSku, comp.ProductID)
		}

		// Check if enough stock is available
		if variantStock.AvailableStock < comp.Quantity {
			return fmt.Errorf("insufficient stock for variant %s: required=%.0f, available=%.2f",
				*comp.VariantSku, comp.Quantity, variantStock.AvailableStock)
		}

		log.Printf("[PG_STOCK] Stock check passed for %s: required=%.0f, available=%.2f",
			*comp.VariantSku, comp.Quantity, variantStock.AvailableStock)
	}

	// All checks passed, now deduct the stock
	for _, comp := range components {
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			continue
		}

		reason := fmt.Sprintf("Product Group Assembly: PG-%s", productGroupID)

		// Use RecordStockAdjustment with "out" type to deduct stock
		err := s.variantStockMgmtService.RecordStockAdjustment(
			*comp.VariantSku,
			comp.Quantity,
			"out",
			reason,
			"system", // userID - could be enhanced to get from context
		)
		if err != nil {
			return fmt.Errorf("failed to deduct stock for variant %s: %w", *comp.VariantSku, err)
		}

		log.Printf("[PG_STOCK] Deducted %.0f units from variant %s for product group %s",
			comp.Quantity, *comp.VariantSku, productGroupID)
	}

	return nil
}

// createProductGroupStock creates a stock entry for the product group itself
// Records in BOTH ProductGroupInventory (for internal tracking) AND ProductStock (for stock APIs)
// This allows the product group to appear in /api/stock/summary, /api/dashboard/stock
func (s *productGroupService) createProductGroupStock(productGroupID string, input *input.CreateProductGroupInput) error {
	log.Printf("[PG_STOCK] ===== START: Creating stock for PG %s (%s) =====", productGroupID, input.Name)

	// Get the initial stock quantity from the first component
	// (all components have equal quantity per product group design)
	assemblyQty := input.Products[0].Quantity
	log.Printf("[PG_STOCK] Assembly quantity: %.0f", assemblyQty)

	// Calculate total cost and selling price from components
	totalCost := 0.0
	totalSellingPrice := 0.0

	for _, comp := range input.Products {
		product, err := s.productRepo.FindByID(comp.ProductID)
		if err != nil {
			log.Printf("[PG_STOCK] Warning: Could not find product %s for cost calculation", comp.ProductID)
			continue
		}

		componentCost := product.PurchaseInfo.CostPrice * comp.Quantity
		componentPrice := product.SalesInfo.SellingPrice * comp.Quantity

		totalCost += componentCost
		totalSellingPrice += componentPrice

		log.Printf("[PG_STOCK] Component %s: Cost=%.2f, Selling=%.2f", comp.ProductID, componentCost, componentPrice)
	}

	log.Printf("[PG_STOCK] Total: Cost=%.2f, Selling=%.2f", totalCost, totalSellingPrice)

	// IMPORTANT: We create VariantStock entry for the product group
	// The dashboard stock API (/dashboard/stock) queries from VariantStock table, not ProductStock
	// VariantStock doesn't have foreign key constraints, so we can use pg_xxx IDs directly
	if s.variantStockMgmtService == nil {
		log.Printf("[PG_STOCK] ✗ variantStockMgmtService is nil, cannot create VariantStock entry")
		return fmt.Errorf("variant stock management service not available")
	}

	log.Printf("[PG_STOCK] Starting VariantStock initialization for %s", productGroupID)

	// Initialize variant stock for the product group
	variantStock, err := s.variantStockMgmtService.InitializeVariantStock(
		productGroupID,    // productID - use PG ID
		productGroupID,    // variantSKU - use PG ID as SKU
		input.Name,        // variantName - use PG name
		input.Name,        // productName - use PG name
		totalSellingPrice, // sellingPrice
		totalCost,         // costPrice (total cost)
	)
	if err != nil {
		log.Printf("[PG_STOCK] ✗ Error initializing variant stock: %v", err)
		return fmt.Errorf("failed to initialize variant stock: %w", err)
	}
	log.Printf("[PG_STOCK] ✓ Initialized VariantStock: ID=%s, SKU=%s, ProductName=%s", variantStock.ID, variantStock.VariantSKU, variantStock.ProductName)

	// Record initial stock in variant stock
	pgReason := fmt.Sprintf("Product Group Assembly: %s", input.Name)
	log.Printf("[PG_STOCK] Recording stock adjustment: SKU=%s, Qty=%.0f, Type=in", productGroupID, assemblyQty)

	err = s.variantStockMgmtService.RecordStockAdjustment(
		productGroupID,
		assemblyQty,
		"in",
		pgReason,
		"system",
	)
	if err != nil {
		log.Printf("[PG_STOCK] ✗ Error recording variant stock adjustment: %v", err)
		return fmt.Errorf("failed to record stock adjustment: %w", err)
	}
	log.Printf("[PG_STOCK] ✓ Recorded stock adjustment: %.0f units added to %s", assemblyQty, productGroupID)

	log.Printf("[PG_STOCK] ===== END: Successfully created VariantStock entry for product group %s (%s) with %.0f units =====", productGroupID, input.Name, assemblyQty)
	return nil
}

func (s *productGroupService) toOutput(productGroup *models.ProductGroup) (*output.ProductGroupOutput, error) {
	return output.ToProductGroupOutput(productGroup)
}

// ValidateStockAvailability checks if there is enough stock to fulfill product group usage
// quantityToCreate: number of product groups to create/consume
func (s *productGroupService) ValidateStockAvailability(productGroupID string, quantityToCreate float64) ([]string, error) {
	productGroup, err := s.productGroupRepo.FindByID(productGroupID)
	if err != nil {
		return nil, fmt.Errorf("product group not found: %v", err)
	}

	warnings := []string{}

	for _, comp := range productGroup.Components {
		totalRequired := comp.Quantity * quantityToCreate

		product, err := s.productRepo.FindByID(comp.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %v", comp.ProductID, err)
		}

		// For product groups, variant SKU is required
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			return nil, fmt.Errorf("variant SKU required for product %s in product group", comp.ProductID)
		}

		// Check variant stock
		variant, err := s.productRepo.GetProductVariantBySKU(*comp.VariantSku)
		if err != nil {
			return nil, fmt.Errorf("variant %s not found: %v", *comp.VariantSku, err)
		}

		available := variant.StockQuantity

		// Check if would go below reorder level
		if variant.ReorderLevel > 0 && (available-totalRequired) <= variant.ReorderLevel {
			warnings = append(warnings,
				fmt.Sprintf("WARNING: %s (variant: %s) stock would reach reorder level. Current: %f, Required: %f, Reorder Level: %f",
					product.Name, *comp.VariantSku, available, totalRequired, variant.ReorderLevel))
		}

		if available < totalRequired {
			return nil, fmt.Errorf("insufficient stock for %s (variant: %s): available=%f, required=%f",
				product.Name, *comp.VariantSku, available, totalRequired)
		}
	}

	return warnings, nil
}

// ========================
// Helper Methods
// ========================

// GetProductGroupComponents returns detailed information about all components in a product group
// This is useful for displaying what's inside the product group
func (s *productGroupService) GetProductGroupComponents(productGroupID string) (map[string]interface{}, error) {
	productGroup, err := s.productGroupRepo.FindByID(productGroupID)
	if err != nil {
		return nil, fmt.Errorf("product group not found")
	}

	components := make([]map[string]interface{}, 0, len(productGroup.Components))

	for _, comp := range productGroup.Components {
		product, err := s.productRepo.FindByID(comp.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %v", comp.ProductID, err)
		}

		componentInfo := map[string]interface{}{
			"product_id":   comp.ProductID,
			"product_name": product.Name,
			"quantity":     comp.Quantity,
			"variant_sku":  comp.VariantSku,
		}

		// If variant SKU is specified, get variant details
		if comp.VariantSku != nil && *comp.VariantSku != "" {
			variant, err := s.productRepo.GetProductVariantBySKU(*comp.VariantSku)
			if err == nil && variant != nil {
				variantAttrs := make(map[string]string)
				for _, attr := range variant.Attributes {
					variantAttrs[attr.Key] = attr.Value
				}

				componentInfo["variant_name"] = variant.VariantName
				componentInfo["variant_attributes"] = variantAttrs
				componentInfo["variant_price"] = variant.SellingPrice
				componentInfo["variant_stock"] = variant.StockQuantity
			}
		}

		components = append(components, componentInfo)
	}

	return map[string]interface{}{
		"product_group_id":   productGroup.ID,
		"product_group_name": productGroup.Name,
		"is_active":          productGroup.IsActive,
		"components":         components,
		"total_components":   len(components),
	}, nil
}

// GetTotalPackageValue calculates the total cost/selling price of a product group
func (s *productGroupService) GetTotalPackageValue(productGroupID string, valueType string) (float64, error) {
	productGroup, err := s.productGroupRepo.FindByID(productGroupID)
	if err != nil {
		return 0, fmt.Errorf("product group not found")
	}

	totalValue := 0.0

	for _, comp := range productGroup.Components {
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			return 0, fmt.Errorf("variant SKU required for calculating package value")
		}

		variant, err := s.productRepo.GetProductVariantBySKU(*comp.VariantSku)
		if err != nil {
			return 0, fmt.Errorf("variant %s not found: %v", *comp.VariantSku, err)
		}

		var price float64
		if valueType == "cost" {
			price = variant.CostPrice
		} else {
			price = variant.SellingPrice
		}

		totalValue += price * comp.Quantity
	}

	return totalValue, nil
}
