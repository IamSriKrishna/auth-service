package services

import (
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type InventoryService interface {
	// ConsumeItemGroupComponents deducts stock for all components in an item group
	// quantity is the number of item groups to consume (e.g., 10 bottles)
	ConsumeItemGroupComponents(itemGroup *models.ItemGroup, quantity float64) ([]InventoryWarning, error)

	// CheckItemGroupAvailability verifies if enough stock exists for all components
	CheckItemGroupAvailability(itemGroup *models.ItemGroup, quantity float64) (bool, []InventoryIssue, error)

	// CheckReorderPointsForItem checks if using the item will breach reorder point
	CheckReorderPointsForItem(itemID string, variantSKU *string, quantityToUse float64) (bool, string, error)

	// SyncOpeningStockToInventoryBalance syncs opening stock to inventory balance for an item
	SyncOpeningStockToInventoryBalance(itemID string) error
}

type InventoryWarning struct {
	ItemID       string
	ItemName     string
	VariantSKU   *string
	CurrentStock float64
	ReorderLevel float64
	Message      string
}

type InventoryIssue struct {
	ItemID         string
	ItemName       string
	VariantSKU     *string
	StockRequired  float64
	StockAvailable float64
	Message        string
}

type inventoryService struct {
	itemRepo             repo.ItemRepository
	itemGroupRepo        repo.ItemGroupRepository
	inventoryBalanceRepo repo.InventoryBalanceRepository
	openingStockRepo     repo.OpeningStockRepository
}

func NewInventoryService(itemRepo repo.ItemRepository, itemGroupRepo repo.ItemGroupRepository, inventoryBalanceRepo repo.InventoryBalanceRepository, openingStockRepo repo.OpeningStockRepository) InventoryService {
	return &inventoryService{
		itemRepo:             itemRepo,
		itemGroupRepo:        itemGroupRepo,
		inventoryBalanceRepo: inventoryBalanceRepo,
		openingStockRepo:     openingStockRepo,
	}
}

// ConsumeItemGroupComponents deducts stock for all components in an item group
func (s *inventoryService) ConsumeItemGroupComponents(itemGroup *models.ItemGroup, quantity float64) ([]InventoryWarning, error) {
	warnings := []InventoryWarning{}

	for _, component := range itemGroup.Components {
		totalQuantityNeeded := component.Quantity * quantity

		// Only deduct stock for variant items
		if component.VariantSku == nil || *component.VariantSku == "" {
			continue
		}

		// First check reorder point
		atReorderPoint, msg, err := s.CheckReorderPointsForItem(component.ItemID, component.VariantSku, totalQuantityNeeded)
		if err != nil {
			return nil, fmt.Errorf("error checking reorder point for %s: %v", component.ItemID, err)
		}

		// Deduct the stock
		if err := s.itemRepo.DeductStockQuantity(component.ItemID, component.VariantSku, totalQuantityNeeded); err != nil {
			return nil, fmt.Errorf("failed to deduct stock for component %s: %v", component.ItemID, err)
		}

		// Add warning if at reorder point
		if atReorderPoint {
			warnings = append(warnings, InventoryWarning{
				ItemID:     component.ItemID,
				VariantSKU: component.VariantSku,
				Message:    msg,
			})
		}
	}

	return warnings, nil
}

// CheckItemGroupAvailability verifies if enough stock exists for all components
// Uses InventoryBalance table which tracks opening stock, purchases, and consumption
// Falls back to opening_stock tables if inventory_balance shows 0
func (s *inventoryService) CheckItemGroupAvailability(itemGroup *models.ItemGroup, quantity float64) (bool, []InventoryIssue, error) {
	issues := []InventoryIssue{}
	available := true

	if len(itemGroup.Components) == 0 {
		return false, nil, fmt.Errorf("item group has no components")
	}

	// Get the base quantity from first component to calculate per-unit requirements
	baseQuantity := itemGroup.Components[0].Quantity

	for _, component := range itemGroup.Components {
		// Calculate per-unit requirement: (component_qty / base_qty) * units_to_manufacture
		// Example: if component qty=500 and base qty=500, ratio is 1:1, so for 100 units we need 100 of that component
		totalQuantityNeeded := (component.Quantity / baseQuantity) * quantity

		var item *models.Item
		var err error

		item, err = s.itemRepo.FindByID(component.ItemID)
		if err != nil {
			return false, nil, fmt.Errorf("item %s not found", component.ItemID)
		}

		var currentStock float64
		var balance *models.InventoryBalance

		// For variant items, check variant SKU is provided and use inventory balance
		if item.ItemDetails.Structure == "variants" {
			if component.VariantSku == nil || *component.VariantSku == "" {
				return false, nil, fmt.Errorf("variant SKU required for variant item %s in item group", component.ItemID)
			}

			// Check inventory balance for the variant
			log.Printf("[INVENTORY_CHECK] Checking variant item %s variant %s - Need %.0f units", component.ItemID, *component.VariantSku, totalQuantityNeeded)
			balance, err = s.inventoryBalanceRepo.GetBalance(component.ItemID, component.VariantSku)
			if err != nil {
				return false, nil, fmt.Errorf("failed to check inventory balance for variant %s: %v", *component.VariantSku, err)
			}
			currentStock = balance.AvailableQuantity
			log.Printf("[INVENTORY_CHECK] Variant item %s variant %s - Available: %.0f", component.ItemID, *component.VariantSku, currentStock)

			// Fallback to opening_stock if available quantity is 0
			if currentStock == 0 {
				variantStock, err := s.openingStockRepo.GetVariantOpeningStock(*component.VariantSku)
				if err == nil && variantStock != nil && variantStock.OpeningStock > 0 {
					log.Printf("[INVENTORY_CHECK] Variant item %s variant %s - Fallback to opening stock: %.0f", component.ItemID, *component.VariantSku, variantStock.OpeningStock)
					currentStock = variantStock.OpeningStock
				}
			}
		} else {
			// For single-structure items, check inventory balance (opening stock is tracked here)
			log.Printf("[INVENTORY_CHECK] Checking single item %s - Need %.0f units", component.ItemID, totalQuantityNeeded)
			balance, err = s.inventoryBalanceRepo.GetBalance(component.ItemID, nil)
			if err != nil {
				return false, nil, fmt.Errorf("failed to check inventory balance for item %s: %v", component.ItemID, err)
			}
			currentStock = balance.AvailableQuantity
			log.Printf("[INVENTORY_CHECK] Single item %s - Available: %.0f", component.ItemID, currentStock)

			// Fallback to opening_stock if available quantity is 0
			if currentStock == 0 {
				openingStock, err := s.openingStockRepo.GetOpeningStock(component.ItemID)
				if err == nil && openingStock != nil && openingStock.OpeningStock > 0 {
					log.Printf("[INVENTORY_CHECK] Single item %s - Fallback to opening stock: %.0f", component.ItemID, openingStock.OpeningStock)
					currentStock = openingStock.OpeningStock
				}
			}

			// If variant_sku is provided but item has no variants, try to get variant balance anyway
			if currentStock == 0 && component.VariantSku != nil && *component.VariantSku != "" {
				variantBalance, err := s.inventoryBalanceRepo.GetBalance(component.ItemID, component.VariantSku)
				if err == nil && variantBalance != nil {
					currentStock = variantBalance.AvailableQuantity
					log.Printf("[INVENTORY_CHECK] Single item %s with variant sku %s found - Available: %.0f", component.ItemID, *component.VariantSku, currentStock)
				}

				// Fallback to opening_stock for variant if still 0
				if currentStock == 0 {
					variantStock, err := s.openingStockRepo.GetVariantOpeningStock(*component.VariantSku)
					if err == nil && variantStock != nil && variantStock.OpeningStock > 0 {
						log.Printf("[INVENTORY_CHECK] Single item %s variant %s - Fallback to opening stock: %.0f", component.ItemID, *component.VariantSku, variantStock.OpeningStock)
						currentStock = variantStock.OpeningStock
					}
				}
			}
		}

		if currentStock < totalQuantityNeeded {
			available = false
			issues = append(issues, InventoryIssue{
				ItemID:         component.ItemID,
				ItemName:       item.Name,
				VariantSKU:     component.VariantSku,
				StockRequired:  totalQuantityNeeded,
				StockAvailable: currentStock,
				Message: fmt.Sprintf("Insufficient stock: %s (required: %.0f, available: %.0f). Please set opening stock first using PUT /items/{id}/opening-stock or PUT /items/{id}/variants/opening-stock",
					item.Name, totalQuantityNeeded, currentStock),
			})
		}
	}

	return available, issues, nil
}

// CheckReorderPointsForItem checks if using the item will breach reorder point
func (s *inventoryService) CheckReorderPointsForItem(itemID string, variantSKU *string, quantityToUse float64) (bool, string, error) {
	variant, err := s.itemRepo.CheckReorderPoint(itemID, variantSKU)
	if err != nil {
		return false, "", err
	}

	if variant != nil {
		// Stock is at or below reorder level
		msg := fmt.Sprintf("WARNING: %s stock is at reorder level (%f). Current stock: %f, Reorder level: %f",
			variant.SKU, variant.StockQuantity, variant.StockQuantity, variant.ReorderLevel)
		return true, msg, nil
	}

	return false, "", nil
}

// SyncOpeningStockToInventoryBalance syncs opening stock data to inventory balance table for an item
func (s *inventoryService) SyncOpeningStockToInventoryBalance(itemID string) error {
	log.Printf("[SYNC] Starting sync for item %s", itemID)

	// Get the item to check its structure
	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return fmt.Errorf("item %s not found: %v", itemID, err)
	}

	// For single-structure items
	if item.ItemDetails.Structure == "single" {
		openingStock, err := s.openingStockRepo.GetOpeningStock(itemID)
		if err != nil {
			return fmt.Errorf("failed to get opening stock for item %s: %v", itemID, err)
		}

		if openingStock.OpeningStock > 0 {
			balance, err := s.inventoryBalanceRepo.GetBalance(itemID, nil)
			if err != nil {
				return fmt.Errorf("failed to get inventory balance for item %s: %v", itemID, err)
			}

			balance.CurrentQuantity = openingStock.OpeningStock
			balance.AvailableQuantity = openingStock.OpeningStock
			balance.AverageRate = openingStock.OpeningStockRatePerUnit
			balance.LastInventorySyncAt = time.Now()
			balance.UpdatedAt = time.Now()

			if err := s.inventoryBalanceRepo.UpdateBalance(balance); err != nil {
				return fmt.Errorf("failed to update inventory balance for item %s: %v", itemID, err)
			}
			log.Printf("[SYNC] Synced single item %s - Stock: %.0f", itemID, openingStock.OpeningStock)
		}
	} else if item.ItemDetails.Structure == "variants" {
		// For variant-structure items, sync each variant
		variantStocks, err := s.openingStockRepo.GetAllVariantOpeningStocks(itemID)
		if err != nil {
			return fmt.Errorf("failed to get variant opening stocks for item %s: %v", itemID, err)
		}

		for _, variantStock := range variantStocks {
			if variantStock.OpeningStock > 0 {
				balance, err := s.inventoryBalanceRepo.GetBalance(itemID, &variantStock.VariantSKU)
				if err != nil {
					return fmt.Errorf("failed to get inventory balance for variant %s: %v", variantStock.VariantSKU, err)
				}

				balance.CurrentQuantity = variantStock.OpeningStock
				balance.AvailableQuantity = variantStock.OpeningStock
				balance.AverageRate = variantStock.OpeningStockRatePerUnit
				balance.LastInventorySyncAt = time.Now()
				balance.UpdatedAt = time.Now()

				if err := s.inventoryBalanceRepo.UpdateBalance(balance); err != nil {
					return fmt.Errorf("failed to update inventory balance for variant %s: %v", variantStock.VariantSKU, err)
				}
				log.Printf("[SYNC] Synced variant %s of item %s - Stock: %.0f", variantStock.VariantSKU, itemID, variantStock.OpeningStock)
			}
		}
	}

	log.Printf("[SYNC] Completed sync for item %s", itemID)
	return nil
}
