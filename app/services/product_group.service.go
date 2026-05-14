package services

import (
	"fmt"
	"log"
	"math"
	"time"

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
	Reorder(id string, input *input.UpdateProductGroupInput) (*output.ProductGroupOutput, error)
	ReorderWithSummary(id string, input *input.UpdateProductGroupInput) (*output.ReorderProductGroupOutput, error)
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

// adjustStockForReorder adjusts stock based on quantity changes during product group reorder
// When component quantities decrease, you're creating MORE product groups from existing stock
// Both increase and decrease consume stock from inventory
func (s *productGroupService) adjustStockForReorder(productGroupID string, oldComponents []models.ProductGroupComponent, newComponents []input.ProductGroupComponentInput) error {
	log.Printf("[PG_REORDER_STOCK] ===== START: Adjusting stock for reorder of %s =====", productGroupID)

	// Build map of old quantities
	oldQtyMap := make(map[string]float64) // key: productID:variantSKU
	for _, comp := range oldComponents {
		variantSku := ""
		if comp.VariantSku != nil {
			variantSku = *comp.VariantSku
		}
		key := comp.ProductID + ":" + variantSku
		oldQtyMap[key] = comp.Quantity
	}

	// First, validate we have enough stock for the new assembly requirements
	for _, newComp := range newComponents {
		variantSku := ""
		if newComp.VariantSku == nil || *newComp.VariantSku == "" {
			log.Printf("[PG_REORDER_STOCK] Skipping variant SKU validation for %s (not specified)", newComp.ProductID)
			continue
		}
		variantSku = *newComp.VariantSku

		key := newComp.ProductID + ":" + variantSku
		oldQty := oldQtyMap[key]
		netConsumption := oldQty - newComp.Quantity // Additional stock consumed

		// When decreasing per-group quantity, we can make more groups, consuming additional stock
		// E.g., reducing from 90->10 means we can now make 9x per "old group", consuming 80 more
		if netConsumption > 0 {
			variantStock, err := s.variantStockMgmtService.GetVariantStockSummary(variantSku)
			if err != nil {
				return fmt.Errorf("variant %s stock not found: %w", variantSku, err)
			}

			if variantStock.AvailableStock < netConsumption {
				return fmt.Errorf("insufficient stock for variant %s during reorder: need %.0f more, have %.0f",
					variantSku, netConsumption, variantStock.AvailableStock)
			}
			log.Printf("[PG_REORDER_STOCK] Stock check passed for %s: need %.0f additional, available=%.0f",
				variantSku, netConsumption, variantStock.AvailableStock)
		}
	}

	// All checks passed, now adjust the stock
	// Logic: stock adjustment = ABS(oldQty - newQty)
	// Both increases and decreases consume stock from inventory
	for _, newComp := range newComponents {
		if newComp.VariantSku == nil || *newComp.VariantSku == "" {
			log.Printf("[PG_REORDER_STOCK] Skipping stock adjustment for %s (no variant SKU)", newComp.ProductID)
			continue
		}

		key := newComp.ProductID + ":" + *newComp.VariantSku
		oldQty := oldQtyMap[key]
		adjustment := oldQty - newComp.Quantity // Positive = more stock to consume, Negative = stock to release

		if adjustment == 0 {
			log.Printf("[PG_REORDER_STOCK] No change for %s, skipping adjustment", *newComp.VariantSku)
			continue
		}

		reason := fmt.Sprintf("Product Group Reorder: %s (%.0f->%.0f per unit)", productGroupID, oldQty, newComp.Quantity)

		if adjustment > 0 {
			// Decreasing component qty: can now make MORE groups - consume additional stock
			// E.g., 90->10 means making 9x more groups, consuming 80 more units
			log.Printf("[PG_REORDER_STOCK] Decreasing %s: %.0f -> %.0f (consume %.0f more stock for additional groups)",
				*newComp.VariantSku, oldQty, newComp.Quantity, adjustment)
			err := s.variantStockMgmtService.RecordStockAdjustment(
				*newComp.VariantSku,
				adjustment,
				"out",
				reason,
				"system",
			)
			if err != nil {
				return fmt.Errorf("failed to deduct stock for variant %s: %w", *newComp.VariantSku, err)
			}
			log.Printf("[PG_REORDER_STOCK] ✓ Deducted %.0f units from variant %s (additional PG assemblies)", adjustment, *newComp.VariantSku)
		} else {
			// Increasing component qty: need LESS stock, release the difference
			// E.g., 10->90 means each group needs 80 more, so release previously allocated stock
			absAdjustment := -adjustment
			log.Printf("[PG_REORDER_STOCK] Increasing %s: %.0f -> %.0f (release %.0f units)",
				*newComp.VariantSku, oldQty, newComp.Quantity, absAdjustment)
			err := s.variantStockMgmtService.RecordStockAdjustment(
				*newComp.VariantSku,
				absAdjustment,
				"in",
				reason,
				"system",
			)
			if err != nil {
				return fmt.Errorf("failed to release stock for variant %s: %w", *newComp.VariantSku, err)
			}
			log.Printf("[PG_REORDER_STOCK] ✓ Released %.0f units to variant %s", absAdjustment, *newComp.VariantSku)
		}
	}

	log.Printf("[PG_REORDER_STOCK] ===== END: Stock adjustments complete =====")
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

	// Also create ProductGroupInventory for the sales order stock deduction workflow
	if s.productGroupInventoryService != nil {
		if err := s.productGroupInventoryService.InitializeProductGroupInventory(productGroupID); err != nil {
			log.Printf("[PG_STOCK] Warning: Failed to initialize product group inventory: %v", err)
			// Don't return error - VariantStock is already created, this is just backup
		} else {
			log.Printf("[PG_STOCK] ✓ Created ProductGroupInventory record for %s", productGroupID)

			// Add initial stock to product group inventory
			if err := s.productGroupInventoryService.AddStock(
				productGroupID,
				assemblyQty,
				fmt.Sprintf("Initial stock from product group assembly: %s", input.Name),
				nil,
			); err != nil {
				log.Printf("[PG_STOCK] Warning: Failed to add initial stock to inventory: %v", err)
			}
		}
	}

	return nil
}

// updateProductGroupStock updates the stock entry for a product group after reordering
// Recalculates how many product groups can be made from available component stock
func (s *productGroupService) updateProductGroupStock(productGroupID string, productGroup *models.ProductGroup) error {
	log.Printf("[PG_STOCK_UPDATE] ===== START: Updating stock for PG %s =====", productGroupID)

	if s.variantStockMgmtService == nil {
		log.Printf("[PG_STOCK_UPDATE] Warning: variantStockMgmtService not available")
		return nil
	}

	// Calculate how many product groups can be made with current component stock
	// PG stock = min(componentA_available / componentA_per_unit, componentB_available / componentB_per_unit, ...)
	minAvailableGroups := math.MaxFloat64
	log.Printf("[PG_STOCK_UPDATE] Calculating available product groups from components:")

	for i, comp := range productGroup.Components {
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			log.Printf("[PG_STOCK_UPDATE]   [%d] No variant SKU, skipping", i)
			continue
		}

		variantStock, err := s.variantStockMgmtService.GetVariantStockSummary(*comp.VariantSku)
		if err != nil {
			log.Printf("[PG_STOCK_UPDATE]   [%d] Warning: Could not get stock for %s: %v", i, *comp.VariantSku, err)
			continue
		}

		if comp.Quantity > 0 {
			availableGroups := variantStock.AvailableStock / comp.Quantity
			log.Printf("[PG_STOCK_UPDATE]   [%d] %s: available=%.0f, per_unit=%.0f, groups=%.0f",
				i, *comp.VariantSku, variantStock.AvailableStock, comp.Quantity, availableGroups)

			if availableGroups < minAvailableGroups {
				minAvailableGroups = availableGroups
			}
		}
	}

	if minAvailableGroups == math.MaxFloat64 {
		minAvailableGroups = 0
	}

	// Round down to whole groups
	newPGStock := math.Floor(minAvailableGroups)
	log.Printf("[PG_STOCK_UPDATE] Calculated product groups available: %.0f", newPGStock)

	// Get or create variant stock for the product group
	pgStock, err := s.variantStockMgmtService.GetVariantStockSummary(productGroupID)
	if err != nil {
		log.Printf("[PG_STOCK_UPDATE] Product group stock entry not found, skipping update")
		return nil
	}

	// Get old quantity
	oldPGStock := pgStock.CurrentStock
	log.Printf("[PG_STOCK_UPDATE] Old PG stock: %.0f, New PG stock: %.0f", oldPGStock, newPGStock)

	// If stock changed, update via RecordStockAdjustment to maintain audit trail
	if newPGStock != oldPGStock {
		difference := newPGStock - oldPGStock

		if difference > 0 {
			// Stock increased - add the difference
			log.Printf("[PG_STOCK_UPDATE] Adding %.0f units to %s", difference, productGroupID)
			err := s.variantStockMgmtService.RecordStockAdjustment(
				productGroupID,
				difference,
				"in",
				fmt.Sprintf("Product Group Reorder: %s", productGroup.Name),
				"system",
			)
			if err != nil {
				log.Printf("[PG_STOCK_UPDATE] ✗ Error adding stock: %v", err)
				return fmt.Errorf("failed to add stock for product group: %w", err)
			}
			log.Printf("[PG_STOCK_UPDATE] ✓ Added %.0f units", difference)
		} else {
			// Stock decreased - deduct the difference
			absAdjustment := -difference
			log.Printf("[PG_STOCK_UPDATE] Deducting %.0f units from %s", absAdjustment, productGroupID)
			err := s.variantStockMgmtService.RecordStockAdjustment(
				productGroupID,
				absAdjustment,
				"out",
				fmt.Sprintf("Product Group Reorder: %s", productGroup.Name),
				"system",
			)
			if err != nil {
				log.Printf("[PG_STOCK_UPDATE] ✗ Error deducting stock: %v", err)
				return fmt.Errorf("failed to deduct stock for product group: %w", err)
			}
			log.Printf("[PG_STOCK_UPDATE] ✓ Deducted %.0f units", absAdjustment)
		}
	}

	log.Printf("[PG_STOCK_UPDATE] ===== END: Successfully updated stock for product group %s =====", productGroupID)
	return nil

	// Calculate total cost and selling price from components with new quantities
	totalCost := 0.0
	totalSellingPrice := 0.0

	for _, comp := range productGroup.Components {
		product, err := s.productRepo.FindByID(comp.ProductID)
		if err != nil {
			log.Printf("[PG_STOCK_UPDATE] Warning: Could not find product %s for cost calculation", comp.ProductID)
			continue
		}

		componentCost := product.PurchaseInfo.CostPrice * comp.Quantity
		componentPrice := product.SalesInfo.SellingPrice * comp.Quantity

		totalCost += componentCost
		totalSellingPrice += componentPrice

		log.Printf("[PG_STOCK_UPDATE] Component %s: Cost=%.2f, Selling=%.2f", comp.ProductID, componentCost, componentPrice)
	}

	log.Printf("[PG_STOCK_UPDATE] Total: Cost=%.2f, Selling=%.2f", totalCost, totalSellingPrice)

	// Get the existing variant stock entry for this product group
	existingStock, err := s.variantStockMgmtService.GetVariantStockSummary(productGroupID)
	if err != nil {
		log.Printf("[PG_STOCK_UPDATE] Creating new variant stock entry: %v", err)
		// If it doesn't exist, create it
		_, err := s.variantStockMgmtService.InitializeVariantStock(
			productGroupID,    // productID
			productGroupID,    // variantSKU
			productGroup.Name, // variantName
			productGroup.Name, // productName
			totalSellingPrice, // sellingPrice
			totalCost,         // costPrice
		)
		if err != nil {
			log.Printf("[PG_STOCK_UPDATE] ✗ Error initializing variant stock: %v", err)
			return fmt.Errorf("failed to initialize variant stock: %w", err)
		}
	} else {
		log.Printf("[PG_STOCK_UPDATE] Existing variant stock found: PurchasedStock=%.0f", existingStock.PurchasedStock)
		log.Printf("[PG_STOCK_UPDATE] Variant stock quantities are adjusted directly during reorder via RecordStockAdjustment()")
	}

	// Reinitialize product group inventory
	if s.productGroupInventoryService != nil {
		err = s.productGroupInventoryService.InitializeProductGroupInventory(productGroupID)
		if err != nil {
			log.Printf("[PG_STOCK_UPDATE] Warning: failed to reinitialize product group inventory: %v", err)
		}
	}

	log.Printf("[PG_STOCK_UPDATE] ===== END: Successfully updated stock for product group %s =====", productGroupID)
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

// Reorder updates an existing product group by changing quantities of existing components
// Can ONLY update quantities of existing products, not add/remove products
// Validates stock availability and consolidates stock from components
func (s *productGroupService) Reorder(id string, reorderInput *input.UpdateProductGroupInput) (*output.ProductGroupOutput, error) {
	// Get existing product group
	productGroup, err := s.productGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("product group not found")
	}

	// If products/components are provided in reorder request
	if len(reorderInput.Products) > 0 {
		// Build map of existing products by ID and SKU
		existingProductsMap := make(map[string]map[string]*models.ProductGroupComponent)
		for i := range productGroup.Components {
			productID := productGroup.Components[i].ProductID
			variantSku := ""
			if productGroup.Components[i].VariantSku != nil {
				variantSku = *productGroup.Components[i].VariantSku
			}

			if existingProductsMap[productID] == nil {
				existingProductsMap[productID] = make(map[string]*models.ProductGroupComponent)
			}
			existingProductsMap[productID][variantSku] = &productGroup.Components[i]
		}

		// Validate all new products already exist in the group
		newProductsMap := make(map[string]float64) // product_id:variant_sku -> new_quantity
		for i, comp := range reorderInput.Products {
			variantSku := ""
			if comp.VariantSku != nil {
				variantSku = *comp.VariantSku
			}

			key := comp.ProductID + ":" + variantSku

			// Check if this product+variant exists in the original group
			if existingProductsMap[comp.ProductID] == nil || existingProductsMap[comp.ProductID][variantSku] == nil {
				return nil, fmt.Errorf("product %d (%s - %s) does not exist in this product group. Reorder can only change quantities of existing products", i+1, comp.ProductID, variantSku)
			}

			// Validate quantity
			if comp.Quantity <= 0 {
				return nil, fmt.Errorf("quantity must be greater than 0 for product %s", comp.ProductID)
			}

			// Validate quantity is a whole number
			if comp.Quantity != float64(int64(comp.Quantity)) {
				return nil, fmt.Errorf("product %d quantity must be a whole number. Got: %f", i+1, comp.Quantity)
			}

			newProductsMap[key] = comp.Quantity
		}

		// Check that all existing products are being reordered
		for productID, variantMap := range existingProductsMap {
			for variantSku := range variantMap {
				key := productID + ":" + variantSku
				if _, exists := newProductsMap[key]; !exists {
					return nil, fmt.Errorf("all existing products must be reordered. Product %s (variant %s) is missing", productID, variantSku)
				}
			}
		}

		// Check stock availability before reorder
		stockCheckErrors := make([]string, 0)
		for _, comp := range reorderInput.Products {
			if comp.VariantSku != nil && *comp.VariantSku != "" {
				// Check if this variant has enough stock for the NEW quantity
				variant, err := s.variantStockMgmtService.GetVariantStockSummary(*comp.VariantSku)
				if err != nil {
					stockCheckErrors = append(stockCheckErrors, fmt.Sprintf("cannot check stock for %s: %v", *comp.VariantSku, err))
					continue
				}

				// Get old quantity for this product
				oldQuantity := 0.0
				for _, existing := range productGroup.Components {
					if existing.ProductID == comp.ProductID && ((comp.VariantSku == nil && existing.VariantSku == nil) || (existing.VariantSku != nil && *existing.VariantSku == *comp.VariantSku)) {
						oldQuantity = existing.Quantity
						break
					}
				}

				// Calculate additional stock needed
				additionalStock := comp.Quantity - oldQuantity

				// Only check if we need more stock
				if additionalStock > 0 {
					if variant == nil || variant.AvailableStock < additionalStock {
						availableStock := 0.0
						if variant != nil {
							availableStock = variant.AvailableStock
						}
						stockCheckErrors = append(stockCheckErrors, fmt.Sprintf("insufficient stock for %s: need +%d, have %.0f", *comp.VariantSku, int64(additionalStock), availableStock))
					}
				}
			}
		}

		if len(stockCheckErrors) > 0 {
			return nil, fmt.Errorf("stock check failed: %s", stockCheckErrors[0])
		}

		// Adjust variant stock based on component quantity changes
		// This handles both increases (deduct additional stock) and decreases (release stock back)
		err = s.adjustStockForReorder(id, productGroup.Components, reorderInput.Products)
		if err != nil {
			return nil, fmt.Errorf("failed to adjust stock for reorder: %w", err)
		}

		// Update quantities in existing components (don't replace, just update quantities)
		for i := range productGroup.Components {
			for _, comp := range reorderInput.Products {
				variantSku := ""
				if comp.VariantSku != nil {
					variantSku = *comp.VariantSku
				}

				if productGroup.Components[i].ProductID == comp.ProductID &&
					((comp.VariantSku == nil && productGroup.Components[i].VariantSku == nil) ||
						(productGroup.Components[i].VariantSku != nil && *productGroup.Components[i].VariantSku == variantSku)) {
					productGroup.Components[i].Quantity = comp.Quantity
					productGroup.Components[i].UpdatedAt = time.Now()
					break
				}
			}
		}

		// Reinitialize product group inventory
		if s.productGroupInventoryService != nil {
			err = s.productGroupInventoryService.InitializeProductGroupInventory(id)
			if err != nil {
				log.Printf("Warning: failed to reinitialize product group inventory: %v", err)
			}
		}
	}

	// Update name if provided
	if reorderInput.Name != "" {
		productGroup.Name = reorderInput.Name
	}

	// Update description if provided
	if reorderInput.Description != "" {
		productGroup.Description = reorderInput.Description
	}

	// Consolidate duplicate components (keep only one per product_id:variant_sku with latest UpdatedAt)
	consolidatedMap := make(map[string]*models.ProductGroupComponent)
	for i := range productGroup.Components {
		variantSku := ""
		if productGroup.Components[i].VariantSku != nil {
			variantSku = *productGroup.Components[i].VariantSku
		}
		key := productGroup.Components[i].ProductID + ":" + variantSku

		// Keep the component with the latest UpdatedAt timestamp
		if existing, exists := consolidatedMap[key]; !exists {
			consolidatedMap[key] = &productGroup.Components[i]
		} else if productGroup.Components[i].UpdatedAt.After(existing.UpdatedAt) {
			consolidatedMap[key] = &productGroup.Components[i]
		}
	}

	// Rebuild components array with consolidated ones
	consolidatedComponents := make([]models.ProductGroupComponent, 0, len(consolidatedMap))
	for _, comp := range consolidatedMap {
		consolidatedComponents = append(consolidatedComponents, *comp)
	}
	productGroup.Components = consolidatedComponents

	// Save updated product group
	err = s.productGroupRepo.Update(productGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to update product group: %v", err)
	}

	// Update the product group's own stock entry to reflect new component quantities
	if len(reorderInput.Products) > 0 {
		err = s.updateProductGroupStock(id, productGroup)
		if err != nil {
			log.Printf("Warning: failed to update product group stock: %v", err)
			// Don't fail - product group is already updated
		}
	}

	// Return updated product group
	return s.toOutput(productGroup)
}

// ReorderWithSummary returns a reorder summary response instead of full product group
func (s *productGroupService) ReorderWithSummary(id string, reorderInput *input.UpdateProductGroupInput) (*output.ReorderProductGroupOutput, error) {
	// First, get the old product group to track changes
	oldProductGroup, err := s.productGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("product group not found")
	}

	// Build old quantities map
	oldQuantitiesMap := make(map[string]float64) // product_id:variant_sku -> old_quantity
	for _, comp := range oldProductGroup.Components {
		variantSku := ""
		if comp.VariantSku != nil {
			variantSku = *comp.VariantSku
		}
		key := comp.ProductID + ":" + variantSku
		oldQuantitiesMap[key] = comp.Quantity
	}

	// Call Reorder to do the actual update
	_, err = s.Reorder(id, reorderInput)
	if err != nil {
		return nil, err
	}

	// Get the updated product group
	updatedProductGroup, err := s.productGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated product group")
	}

	// Build updates summary
	updates := make([]output.ReorderUpdateOutput, 0, len(reorderInput.Products))
	for _, comp := range reorderInput.Products {
		variantSku := ""
		if comp.VariantSku != nil {
			variantSku = *comp.VariantSku
		}
		key := comp.ProductID + ":" + variantSku

		oldQuantity := oldQuantitiesMap[key]
		newQuantity := comp.Quantity
		stockAdjusted := newQuantity - oldQuantity

		updates = append(updates, output.ReorderUpdateOutput{
			VariantSku:    variantSku,
			OldQuantity:   oldQuantity,
			NewQuantity:   newQuantity,
			StockAdjusted: stockAdjusted,
		})
	}

	// Return summary response
	return &output.ReorderProductGroupOutput{
		ID: id,
		ReorderSummary: &output.ReorderSummaryOutput{
			TotalProducts: len(updates),
			Updates:       updates,
		},
		UpdatedAt: updatedProductGroup.UpdatedAt,
	}, nil
}
