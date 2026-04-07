package services

import (
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

// StockMovementType defines the type of stock movement
type StockMovementType string

const (
	// Inbound movements
	StockMovementPurchaseOrder  StockMovementType = "PURCHASE_ORDER"
	StockMovementPurchaseReturn StockMovementType = "PURCHASE_RETURN"
	StockMovementOpeningStock   StockMovementType = "OPENING_STOCK"
	StockMovementAdjustmentIn   StockMovementType = "ADJUSTMENT_IN"

	// Outbound movements
	StockMovementSalesOrder      StockMovementType = "SALES_ORDER"
	StockMovementSalesReturn     StockMovementType = "SALES_RETURN"
	StockMovementShipment        StockMovementType = "SHIPMENT"
	StockMovementAdjustmentOut   StockMovementType = "ADJUSTMENT_OUT"
	StockMovementProductionUsage StockMovementType = "PRODUCTION_USAGE"
)

// StockMovementRequest holds information for stock movement transaction
type StockMovementRequest struct {
	ItemID       string
	VariantSKU   *string
	Quantity     float64
	Rate         float64
	MovementType StockMovementType
	ReferenceID  string
	ReferenceNo  string
	Notes        string
	CreatedBy    string
	Status       string // For tracking movement status
	SourceType   string // purchase_order, sales_order, shipment, package, etc.
	SourceID     string // ID of the source document
}

// StockMovementResponse holds the result of a stock movement
type StockMovementResponse struct {
	InventoryJournalID uint64
	PreviousQuantity   float64
	NewQuantity        float64
	MovementQuantity   float64
	AverageRate        float64
	Message            string
}

// StockMovementService handles all stock movements across the system
type StockMovementService interface {
	// Record stock movements
	RecordInboundMovement(req *StockMovementRequest) (*StockMovementResponse, error)
	RecordOutboundMovement(req *StockMovementRequest) (*StockMovementResponse, error)

	// Get movement history
	GetInventoryJournal(itemID string, variantSKU *string) ([]models.InventoryJournal, error)
	GetItemInventoryBalance(itemID string, variantSKU *string) (*models.InventoryBalance, error)
}

type stockMovementService struct {
	inventoryBalanceRepo repo.InventoryBalanceRepository
	itemRepo             repo.ItemRepository
}

// NewStockMovementService creates a new stock movement service
func NewStockMovementService(
	inventoryBalanceRepo repo.InventoryBalanceRepository,
	itemRepo repo.ItemRepository,
) StockMovementService {
	return &stockMovementService{
		inventoryBalanceRepo: inventoryBalanceRepo,
		itemRepo:             itemRepo,
	}
}

// RecordInboundMovement records a stock increase (purchase, return, etc.)
func (s *stockMovementService) RecordInboundMovement(req *StockMovementRequest) (*StockMovementResponse, error) {
	log.Printf("[STOCK_IN] Recording inbound movement: ItemID=%s, Type=%s, Qty=%.2f, Rate=%.2f",
		req.ItemID, req.MovementType, req.Quantity, req.Rate)

	// Validate item exists
	item, err := s.itemRepo.FindByID(req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("item not found: %s", req.ItemID)
	}

	// Get or create inventory balance
	balance, err := s.inventoryBalanceRepo.GetBalance(req.ItemID, req.VariantSKU)
	if err != nil {
		// Create new balance if not exists
		balance = &models.InventoryBalance{
			ItemID:            req.ItemID,
			VariantSKU:        req.VariantSKU,
			CurrentQuantity:   0,
			ReservedQuantity:  0,
			AvailableQuantity: 0,
			InTransitQuantity: 0,
			AverageRate:       req.Rate,
		}
	}

	previousQty := balance.CurrentQuantity
	newQty := previousQty + req.Quantity

	// Update average rate (weighted average)
	if previousQty > 0 {
		balance.AverageRate = ((balance.AverageRate * previousQty) + (req.Rate * req.Quantity)) / newQty
	} else {
		balance.AverageRate = req.Rate
	}

	// Update quantities
	balance.CurrentQuantity = newQty
	balance.AvailableQuantity = newQty - balance.ReservedQuantity - balance.InTransitQuantity
	balance.LastReceivedDate = timePtr(time.Now())
	balance.LastInventorySyncAt = time.Now()
	balance.UpdatedAt = time.Now()

	// Save inventory balance
	if err := s.inventoryBalanceRepo.UpdateBalance(balance); err != nil {
		return nil, fmt.Errorf("failed to save inventory balance: %v", err)
	}

	// Record journal entry
	journal := &models.InventoryJournal{
		ItemID:          req.ItemID,
		VariantSKU:      req.VariantSKU,
		TransactionType: string(req.MovementType),
		Quantity:        req.Quantity,
		ReferenceType:   req.SourceType,
		ReferenceID:     req.SourceID,
		ReferenceNo:     req.ReferenceNo,
		Notes:           fmt.Sprintf("%s - %s", req.Notes, item.Name),
		CreatedAt:       time.Now(),
		CreatedBy:       req.CreatedBy,
	}

	if err := s.inventoryBalanceRepo.CreateJournalEntry(journal); err != nil {
		log.Printf("[STOCK_IN] Error saving journal: %v", err)
		// Don't return error, journal is not critical
	}

	log.Printf("[STOCK_IN] Movement recorded successfully: PreviousQty=%.2f, NewQty=%.2f, AvgRate=%.2f",
		previousQty, newQty, balance.AverageRate)

	return &StockMovementResponse{
		PreviousQuantity: previousQty,
		NewQuantity:      newQty,
		MovementQuantity: req.Quantity,
		AverageRate:      balance.AverageRate,
		Message:          fmt.Sprintf("Stock increased from %.2f to %.2f units", previousQty, newQty),
	}, nil
}

// RecordOutboundMovement records a stock decrease (sales, usage, etc.)
func (s *stockMovementService) RecordOutboundMovement(req *StockMovementRequest) (*StockMovementResponse, error) {
	log.Printf("[STOCK_OUT] Recording outbound movement: ItemID=%s, Type=%s, Qty=%.2f",
		req.ItemID, req.MovementType, req.Quantity)

	// Validate item exists
	item, err := s.itemRepo.FindByID(req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("item not found: %s", req.ItemID)
	}

	// Get inventory balance
	balance, err := s.inventoryBalanceRepo.GetBalance(req.ItemID, req.VariantSKU)
	if err != nil {
		return nil, fmt.Errorf("no inventory balance found for item %s", req.ItemID)
	}

	previousQty := balance.AvailableQuantity

	// Check if sufficient stock available
	if balance.AvailableQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: requested=%.2f, available=%.2f for item %s",
			req.Quantity, balance.AvailableQuantity, item.Name)
	}

	// Update quantities
	balance.CurrentQuantity = balance.CurrentQuantity - req.Quantity
	balance.AvailableQuantity = balance.AvailableQuantity - req.Quantity
	balance.LastConsumedDate = timePtr(time.Now())
	balance.LastInventorySyncAt = time.Now()
	balance.UpdatedAt = time.Now()

	// Save inventory balance
	if err := s.inventoryBalanceRepo.UpdateBalance(balance); err != nil {
		return nil, fmt.Errorf("failed to save inventory balance: %v", err)
	}

	// Record journal entry
	journal := &models.InventoryJournal{
		ItemID:          req.ItemID,
		VariantSKU:      req.VariantSKU,
		TransactionType: string(req.MovementType),
		Quantity:        -req.Quantity, // Negative for outbound
		ReferenceType:   req.SourceType,
		ReferenceID:     req.SourceID,
		ReferenceNo:     req.ReferenceNo,
		Notes:           fmt.Sprintf("%s - %s", req.Notes, item.Name),
		CreatedAt:       time.Now(),
		CreatedBy:       req.CreatedBy,
	}

	if err := s.inventoryBalanceRepo.CreateJournalEntry(journal); err != nil {
		log.Printf("[STOCK_OUT] Error saving journal: %v", err)
		// Don't return error, journal is not critical
	}

	log.Printf("[STOCK_OUT] Movement recorded successfully: PreviousQty=%.2f, NewQty=%.2f",
		previousQty, balance.AvailableQuantity)

	return &StockMovementResponse{
		PreviousQuantity: previousQty,
		NewQuantity:      balance.AvailableQuantity,
		MovementQuantity: req.Quantity,
		AverageRate:      balance.AverageRate,
		Message:          fmt.Sprintf("Stock decreased from %.2f to %.2f units", previousQty, balance.AvailableQuantity),
	}, nil
}

// GetInventoryJournal retrieves the transaction history for an item
func (s *stockMovementService) GetInventoryJournal(itemID string, variantSKU *string) ([]models.InventoryJournal, error) {
	_, journalCount, err := s.inventoryBalanceRepo.GetJournalEntries(itemID, 0, 100)
	if err != nil {
		return []models.InventoryJournal{}, fmt.Errorf("failed to get journals: %v", err)
	}
	journals, _, err := s.inventoryBalanceRepo.GetJournalEntries(itemID, 0, int(journalCount))
	if err != nil {
		return []models.InventoryJournal{}, err
	}
	return journals, nil
}

// GetItemInventoryBalance retrieves the current balance for an item
func (s *stockMovementService) GetItemInventoryBalance(itemID string, variantSKU *string) (*models.InventoryBalance, error) {
	return s.inventoryBalanceRepo.GetBalance(itemID, variantSKU)
}
