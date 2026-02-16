
package services

import (
	"fmt"
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
	stockRepo repo.OpeningStockRepository
	itemRepo  repo.ItemRepository
}

func NewOpeningStockService(stockRepo repo.OpeningStockRepository, itemRepo repo.ItemRepository) OpeningStockService {
	return &openingStockService{
		stockRepo: stockRepo,
		itemRepo:  itemRepo,
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
		OpeningStock:          stock.OpeningStock,
		OpeningStockRatePerUnit: stock.OpeningStockRatePerUnit,
		UpdatedAt:             stock.UpdatedAt,
	}, nil
}

func (s *openingStockService) GetOpeningStock(itemID string) (*output.OpeningStockOutput, error) {
	stock, err := s.stockRepo.GetOpeningStock(itemID)
	if err != nil {
		return nil, err
	}
	
	return &output.OpeningStockOutput{
		OpeningStock:          stock.OpeningStock,
		OpeningStockRatePerUnit: stock.OpeningStockRatePerUnit,
		UpdatedAt:             stock.UpdatedAt,
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
		err := s.stockRepo.CreateOrUpdateVariantOpeningStock(
			variantInput.VariantID,
			variantInput.OpeningStock,
			variantInput.OpeningStockRatePerUnit,
		)
		if err != nil {
			return nil, err
		}
		
		if variantInput.OpeningStock > 0 {
			movement := &models.StockMovement{
				ItemID:        itemID,
				VariantID:     &variantInput.VariantID,
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
	
	stockMap := make(map[uint]*models.VariantOpeningStock)
	for i := range stocks {
		stockMap[stocks[i].VariantID] = &stocks[i]
	}
	
	result := make([]output.VariantOpeningStockOutput, len(item.ItemDetails.Variants))
	for i, variant := range item.ItemDetails.Variants {
		stock, exists := stockMap[variant.ID]
		if exists {
			result[i] = output.VariantOpeningStockOutput{
				VariantID:             variant.ID,
				VariantSKU:            variant.SKU,
				OpeningStock:          stock.OpeningStock,
				OpeningStockRatePerUnit: stock.OpeningStockRatePerUnit,
				UpdatedAt:             stock.UpdatedAt,
			}
		} else {
			result[i] = output.VariantOpeningStockOutput{
				VariantID:             variant.ID,
				VariantSKU:            variant.SKU,
				OpeningStock:          0,
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