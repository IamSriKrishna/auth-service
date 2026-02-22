package services

import (
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type OpeningStockService interface {
	UpdateOpeningStock(itemID string, input *input.OpeningStockInput, userID string) (*output.OpeningStockOutput, error)
	GetOpeningStock(itemID string) (*output.OpeningStockOutput, error)
	UpdateVariantsOpeningStock(itemID string, input *input.UpdateVariantsOpeningStockInput, userID string) ([]output.VariantOpeningStockOutput, error)
	GetVariantsOpeningStock(itemID string) ([]output.VariantOpeningStockOutput, error)
	GetStockSummary(itemID string) (*output.StockSummaryOutput, error)
}

type openingStockService struct {
	stockRepo     repo.OpeningStockRepository
	itemRepo      repo.ItemRepository
	inventoryRepo repo.InventoryBalanceRepository
}

func NewOpeningStockService(stockRepo repo.OpeningStockRepository, itemRepo repo.ItemRepository, inventoryRepo repo.InventoryBalanceRepository) OpeningStockService {
	return &openingStockService{
		stockRepo:     stockRepo,
		itemRepo:      itemRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *openingStockService) UpdateOpeningStock(itemID string, input *input.OpeningStockInput, userID string) (*output.OpeningStockOutput, error) {

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("item not found")
	}

	if item.ItemDetails.Structure != "single" {
		return nil, fmt.Errorf("opening stock for variant items must be set per variant")
	}

	err = s.stockRepo.CreateOrUpdateOpeningStock(itemID, input.OpeningStock, input.OpeningStockRatePerUnit)
	if err != nil {
		return nil, err
	}

	// Update inventory balance
	log.Printf("[OPEN_STOCK] Getting balance for item %s (single item)", itemID)
	balance, err := s.inventoryRepo.GetBalance(itemID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory balance: %v", err)
	}
	log.Printf("[OPEN_STOCK] Balance retrieved - ID: %d, Current Available: %.2f", balance.ID, balance.AvailableQuantity)

	balance.CurrentQuantity = input.OpeningStock
	balance.AvailableQuantity = input.OpeningStock
	balance.AverageRate = input.OpeningStockRatePerUnit
	balance.LastInventorySyncAt = time.Now()
	balance.UpdatedAt = time.Now()

	log.Printf("[OPEN_STOCK] Updating balance - ID: %d, Setting Available to: %.2f", balance.ID, input.OpeningStock)
	if err := s.inventoryRepo.UpdateBalance(balance); err != nil {
		return nil, fmt.Errorf("failed to update inventory balance: %v", err)
	}
	log.Printf("[OPEN_STOCK] Balance updated successfully - ID: %d", balance.ID)

	if input.OpeningStock > 0 {
		movement := &models.StockMovement{
			ItemID:        itemID,
			MovementType:  "opening_stock",
			Quantity:      input.OpeningStock,
			RatePerUnit:   input.OpeningStockRatePerUnit,
			ReferenceType: "adjustment",
			Notes:         "Opening stock adjustment",
			CreatedBy:     userID,
			CreatedAt:     time.Now(),
		}
		s.stockRepo.RecordStockMovement(movement)
	}

	stock, err := s.stockRepo.GetOpeningStock(itemID)
	if err != nil {
		return nil, err
	}

	return &output.OpeningStockOutput{
		OpeningStock:            stock.OpeningStock,
		OpeningStockRatePerUnit: stock.OpeningStockRatePerUnit,
		UpdatedAt:               stock.UpdatedAt,
	}, nil
}

func (s *openingStockService) GetOpeningStock(itemID string) (*output.OpeningStockOutput, error) {
	stock, err := s.stockRepo.GetOpeningStock(itemID)
	if err != nil {
		return nil, err
	}

	return &output.OpeningStockOutput{
		OpeningStock:            stock.OpeningStock,
		OpeningStockRatePerUnit: stock.OpeningStockRatePerUnit,
		UpdatedAt:               stock.UpdatedAt,
	}, nil
}

func (s *openingStockService) UpdateVariantsOpeningStock(itemID string, input *input.UpdateVariantsOpeningStockInput, userID string) ([]output.VariantOpeningStockOutput, error) {

	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("item not found")
	}

	if item.ItemDetails.Structure != "variants" {
		return nil, fmt.Errorf("this endpoint is only for variant items")
	}

	for _, variantInput := range input.Variants {
		log.Printf("[OPEN_STOCK] Processing variant %s for item %s", variantInput.VariantSKU, itemID)

		err := s.stockRepo.CreateOrUpdateVariantOpeningStock(
			variantInput.VariantSKU,
			variantInput.OpeningStock,
			variantInput.OpeningStockRatePerUnit,
		)
		if err != nil {
			return nil, err
		}

		// Update inventory balance for variant
		log.Printf("[OPEN_STOCK] Getting balance for variant %s of item %s", variantInput.VariantSKU, itemID)
		balance, err := s.inventoryRepo.GetBalance(itemID, &variantInput.VariantSKU)
		if err != nil {
			return nil, fmt.Errorf("failed to get inventory balance for variant %s: %v", variantInput.VariantSKU, err)
		}
		log.Printf("[OPEN_STOCK] Balance retrieved for variant %s - ID: %d, Current Available: %.2f", variantInput.VariantSKU, balance.ID, balance.AvailableQuantity)

		balance.CurrentQuantity = variantInput.OpeningStock
		balance.AvailableQuantity = variantInput.OpeningStock
		balance.AverageRate = variantInput.OpeningStockRatePerUnit
		balance.LastInventorySyncAt = time.Now()
		balance.UpdatedAt = time.Now()

		log.Printf("[OPEN_STOCK] Updating balance for variant %s - ID: %d, Setting Available to: %.2f", variantInput.VariantSKU, balance.ID, variantInput.OpeningStock)
		if err := s.inventoryRepo.UpdateBalance(balance); err != nil {
			return nil, fmt.Errorf("failed to update inventory balance for variant %s: %v", variantInput.VariantSKU, err)
		}
		log.Printf("[OPEN_STOCK] Balance updated for variant %s - ID: %d", variantInput.VariantSKU, balance.ID)

		if variantInput.OpeningStock > 0 {
			movement := &models.StockMovement{
				ItemID:        itemID,
				VariantSKU:    &variantInput.VariantSKU,
				MovementType:  "opening_stock",
				Quantity:      variantInput.OpeningStock,
				RatePerUnit:   variantInput.OpeningStockRatePerUnit,
				ReferenceType: "adjustment",
				Notes:         "Opening stock adjustment",
				CreatedBy:     userID,
				CreatedAt:     time.Now(),
			}
			s.stockRepo.RecordStockMovement(movement)
		}
	}

	return s.GetVariantsOpeningStock(itemID)
}

func (s *openingStockService) GetVariantsOpeningStock(itemID string) ([]output.VariantOpeningStockOutput, error) {
	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return nil, err
	}

	stocks, err := s.stockRepo.GetAllVariantOpeningStocks(itemID)
	if err != nil {
		return nil, err
	}

	stockMap := make(map[string]*models.VariantOpeningStock)
	for i := range stocks {
		stockMap[stocks[i].VariantSKU] = &stocks[i]
	}

	result := make([]output.VariantOpeningStockOutput, len(item.ItemDetails.Variants))
	for i, variant := range item.ItemDetails.Variants {
		stock, exists := stockMap[variant.SKU]
		if exists {
			result[i] = output.VariantOpeningStockOutput{
				VariantID:               variant.ID,
				VariantSKU:              variant.SKU,
				OpeningStock:            stock.OpeningStock,
				OpeningStockRatePerUnit: stock.OpeningStockRatePerUnit,
				UpdatedAt:               stock.UpdatedAt,
			}
		} else {
			result[i] = output.VariantOpeningStockOutput{
				VariantID:               variant.ID,
				VariantSKU:              variant.SKU,
				OpeningStock:            0,
				OpeningStockRatePerUnit: 0,
			}
		}
	}

	return result, nil
}

func (s *openingStockService) GetStockSummary(itemID string) (*output.StockSummaryOutput, error) {
	stock, err := s.stockRepo.GetOpeningStock(itemID)
	if err != nil {
		return nil, err
	}

	return &output.StockSummaryOutput{
		StockOnHand:              stock.OpeningStock,
		CommittedStock:           0,
		AvailableForSale:         stock.OpeningStock,
		PhysicalStockOnHand:      stock.OpeningStock,
		PhysicalCommittedStock:   0,
		PhysicalAvailableForSale: stock.OpeningStock,
		ToBeInvoiced:             0,
		ToBeBilled:               0,
	}, nil
}
