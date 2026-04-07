package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

// StockMovementType defines the type of stock movement for products
type ProductStockMovementType string

const (
	// Inbound movements
	StockMovementTypePurchaseOrder  ProductStockMovementType = "PURCHASE_ORDER"
	StockMovementTypePurchaseReturn ProductStockMovementType = "PURCHASE_RETURN"
	StockMovementTypeOpeningStock   ProductStockMovementType = "OPENING_STOCK"
	StockMovementTypeAdjustmentIn   ProductStockMovementType = "ADJUSTMENT_IN"

	// Outbound movements
	StockMovementTypeSalesOrder      ProductStockMovementType = "SALES_ORDER"
	StockMovementTypeSalesReturn     ProductStockMovementType = "SALES_RETURN"
	StockMovementTypeShipment        ProductStockMovementType = "SHIPMENT"
	StockMovementTypeAdjustmentOut   ProductStockMovementType = "ADJUSTMENT_OUT"
	StockMovementTypeProductionUsage ProductStockMovementType = "PRODUCTION_USAGE"
)

// StockManagementService handles product inventory operations
type StockManagementService interface {
	// Stock operations
	GetProductStock(productID string) (*models.ProductStock, error)
	GetAllProductStock(offset, limit int) ([]models.ProductStock, int64, error)
	GetLowStockProducts(threshold float64, offset, limit int) ([]models.ProductStock, int64, error)

	// Movement recording
	RecordInboundMovement(productID, referenceType, referenceID, referenceNo string, quantity, rate float64, notes, userID string) error
	RecordOutboundMovement(productID, referenceType, referenceID, referenceNo string, quantity float64, notes, userID string) error
	RecordStockAdjustment(productID string, quantity float64, adjustmentType string, reason, userID string) error

	// Stock history
	GetProductMovementHistory(productID string, offset, limit int) ([]models.StockLedger, int64, error)
	GetMovementsByReference(referenceID string) ([]models.StockLedger, error)
	GetMovementsByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.StockLedger, int64, error)

	// Stock summary
	GetStockSummary(productID string) (map[string]interface{}, error)
}

type stockManagementService struct {
	productStockRepo repo.ProductStockRepository
	stockLedgerRepo  repo.StockLedgerRepository
	productRepo      repo.ProductRepository
}

func NewStockManagementService(
	productStockRepo repo.ProductStockRepository,
	stockLedgerRepo repo.StockLedgerRepository,
	productRepo repo.ProductRepository,
) StockManagementService {
	return &stockManagementService{
		productStockRepo: productStockRepo,
		stockLedgerRepo:  stockLedgerRepo,
		productRepo:      productRepo,
	}
}

// GetProductStock retrieves the current stock for a product
func (s *stockManagementService) GetProductStock(productID string) (*models.ProductStock, error) {
	stock, err := s.productStockRepo.GetByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("product stock not found: %s", productID)
	}
	return stock, nil
}

// GetAllProductStock retrieves all product stocks with pagination
func (s *stockManagementService) GetAllProductStock(offset, limit int) ([]models.ProductStock, int64, error) {
	return s.productStockRepo.GetAll(offset, limit)
}

// GetLowStockProducts retrieves products with stock below threshold
func (s *stockManagementService) GetLowStockProducts(threshold float64, offset, limit int) ([]models.ProductStock, int64, error) {
	return s.productStockRepo.GetLowStockProducts(threshold, offset, limit)
}

// RecordInboundMovement records stock increase (purchase, return, adjustment)
func (s *stockManagementService) RecordInboundMovement(
	productID, referenceType, referenceID, referenceNo string,
	quantity, rate float64,
	notes, userID string,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive for inbound movement")
	}

	log.Printf("[STOCK_IN] Recording inbound: Product=%s, Type=%s, Qty=%.2f, Rate=%.2f, Ref=%s",
		productID, referenceType, quantity, rate, referenceNo)

	// Verify product exists
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("product not found: %s", productID)
	}

	// Get or create stock
	stock, err := s.productStockRepo.GetByProductID(productID)
	if err != nil {
		// Create new stock
		sku := ""
		if product.ProductDetails.BaseSKU != "" {
			sku = product.ProductDetails.BaseSKU
		}
		stock = &models.ProductStock{
			ID:             uuid.New().String(),
			ProductID:      productID,
			ProductName:    product.Name,
			SKU:            sku,
			CurrentStock:   quantity,
			PurchasedStock: quantity,
			AverageCost:    rate,
			AvailableStock: quantity,
		}
	} else {
		// Update existing stock
		prevQty := stock.CurrentStock
		stock.PurchasedStock += quantity

		// Update weighted average cost
		if prevQty > 0 {
			stock.AverageCost = ((stock.AverageCost * prevQty) + (rate * quantity)) / (stock.PurchasedStock)
		} else {
			stock.AverageCost = rate
		}

		// CurrentStock = PurchasedStock - SoldStock
		stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
		stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	}

	stock.LastPurchasedDate = timePtr(time.Now())
	stock.LastStockSyncAt = time.Now()
	stock.UpdatedAt = time.Now()

	if err := s.productStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Record ledger entry
	amount := quantity * rate
	ledger := &models.StockLedger{
		ProductID:        productID,
		MovementType:     string(StockMovementTypePurchaseOrder),
		Quantity:         quantity,
		Rate:             rate,
		Amount:           amount,
		ReferenceType:    referenceType,
		ReferenceID:      referenceID,
		ReferenceNumber:  referenceNo,
		BalanceBeforeQty: stock.CurrentStock - quantity,
		BalanceAfterQty:  stock.CurrentStock,
		CostBeforeAmount: (stock.CurrentStock - quantity) * (stock.AverageCost),
		CostAfterAmount:  stock.CurrentStock * stock.AverageCost,
		Notes:            notes,
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	if err := s.stockLedgerRepo.Create(ledger); err != nil {
		log.Printf("[STOCK_IN] Warning: Failed to create ledger entry: %v", err)
	}

	log.Printf("[STOCK_IN] Success: New balance=%.2f units, Avg Cost=%.2f", stock.CurrentStock, stock.AverageCost)
	return nil
}

// RecordOutboundMovement records stock decrease (sales, usage, etc.)
func (s *stockManagementService) RecordOutboundMovement(
	productID, referenceType, referenceID, referenceNo string,
	quantity float64,
	notes, userID string,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive for outbound movement")
	}

	log.Printf("[STOCK_OUT] Recording outbound: Product=%s, Type=%s, Qty=%.2f, Ref=%s",
		productID, referenceType, quantity, referenceNo)

	// Verify product exists
	if _, err := s.productRepo.FindByID(productID); err != nil {
		return fmt.Errorf("product not found: %s", productID)
	}

	// Get stock
	stock, err := s.productStockRepo.GetByProductID(productID)
	if err != nil {
		return fmt.Errorf("no stock found for product: %s", productID)
	}

	// Check available stock
	if stock.AvailableStock < quantity {
		return fmt.Errorf("insufficient stock: requested=%.2f, available=%.2f for product %s",
			quantity, stock.AvailableStock, productID)
	}

	// Deduct from stock - increment SoldStock
	stock.SoldStock += quantity
	// CurrentStock = PurchasedStock - SoldStock (automatic recalculation)
	stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.LastSoldDate = timePtr(time.Now())
	stock.LastStockSyncAt = time.Now()
	stock.UpdatedAt = time.Now()

	if err := s.productStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Record ledger entry
	amount := quantity * stock.AverageCost
	ledger := &models.StockLedger{
		ProductID:        productID,
		MovementType:     string(StockMovementTypeSalesOrder),
		Quantity:         -quantity, // Negative for outbound
		Rate:             stock.AverageCost,
		Amount:           -amount,
		ReferenceType:    referenceType,
		ReferenceID:      referenceID,
		ReferenceNumber:  referenceNo,
		BalanceBeforeQty: stock.CurrentStock + quantity,
		BalanceAfterQty:  stock.CurrentStock,
		CostBeforeAmount: (stock.CurrentStock + quantity) * stock.AverageCost,
		CostAfterAmount:  stock.CurrentStock * stock.AverageCost,
		Notes:            notes,
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	if err := s.stockLedgerRepo.Create(ledger); err != nil {
		log.Printf("[STOCK_OUT] Warning: Failed to create ledger entry: %v", err)
	}

	log.Printf("[STOCK_OUT] Success: New balance=%.2f units", stock.CurrentStock)
	return nil
}

// RecordStockAdjustment records manual stock adjustments
func (s *stockManagementService) RecordStockAdjustment(
	productID string,
	quantity float64,
	adjustmentType string,
	reason, userID string,
) error {
	stock, err := s.productStockRepo.GetByProductID(productID)
	if err != nil {
		return fmt.Errorf("product stock not found: %s", productID)
	}

	if adjustmentType == "in" {
		stock.PurchasedStock += quantity
		// CurrentStock = PurchasedStock - SoldStock
		stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
	} else if adjustmentType == "out" {
		if stock.AvailableStock < quantity {
			return errors.New("insufficient stock for adjustment")
		}
		stock.SoldStock += quantity
		// CurrentStock = PurchasedStock - SoldStock
		stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
	} else {
		return errors.New("invalid adjustment type")
	}

	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.UpdatedAt = time.Now()

	if err := s.productStockRepo.Update(stock); err != nil {
		return err
	}

	// Record ledger
	var movementQty float64
	if adjustmentType == "in" {
		movementQty = quantity
	} else {
		movementQty = -quantity
	}

	ledger := &models.StockLedger{
		ProductID:        productID,
		MovementType:     string(StockMovementTypeAdjustmentIn),
		Quantity:         movementQty,
		Rate:             stock.AverageCost,
		Amount:           movementQty * stock.AverageCost,
		ReferenceType:    "ADJUSTMENT",
		ReferenceNumber:  reason,
		BalanceBeforeQty: stock.CurrentStock - movementQty,
		BalanceAfterQty:  stock.CurrentStock,
		Notes:            reason,
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	s.stockLedgerRepo.Create(ledger)
	return nil
}

// GetProductMovementHistory retrieves movement history for a product
func (s *stockManagementService) GetProductMovementHistory(productID string, offset, limit int) ([]models.StockLedger, int64, error) {
	return s.stockLedgerRepo.GetProductMovementHistory(productID, offset, limit)
}

// GetMovementsByReference retrieves all movements for a reference document
func (s *stockManagementService) GetMovementsByReference(referenceID string) ([]models.StockLedger, error) {
	return s.stockLedgerRepo.GetByReferenceID(referenceID)
}

// GetMovementsByDateRange retrieves movements within a date range
func (s *stockManagementService) GetMovementsByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.StockLedger, int64, error) {
	return s.stockLedgerRepo.GetByDateRange(fromDate, toDate, offset, limit)
}

// GetStockSummary returns a comprehensive stock summary for a product
func (s *stockManagementService) GetStockSummary(productID string) (map[string]interface{}, error) {
	stock, err := s.productStockRepo.GetByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("product stock not found: %s", productID)
	}

	return map[string]interface{}{
		"product_id":      stock.ProductID,
		"product_name":    stock.ProductName,
		"sku":             stock.SKU,
		"current_stock":   stock.CurrentStock,
		"purchased_total": stock.PurchasedStock,
		"sold_total":      stock.SoldStock,
		"reserved_stock":  stock.ReservedStock,
		"available_stock": stock.AvailableStock,
		"average_cost":    stock.AverageCost,
		"stock_value":     stock.CurrentStock * stock.AverageCost,
		"last_purchased":  stock.LastPurchasedDate,
		"last_sold":       stock.LastSoldDate,
		"last_sync":       stock.LastStockSyncAt,
	}, nil
}
