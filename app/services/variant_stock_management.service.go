package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VariantStockManagementService handles variant-level stock operations
type VariantStockManagementService interface {
	// Initialization
	InitializeVariantStock(productID, variantSKU, variantName, productName string, sellingPrice, costPrice float64) (*models.VariantStock, error)

	// Stock movements
	RecordPurchaseInbound(productID, variantSKU string, quantity, rate float64, referenceType, referenceID, referenceNumber string, userID string) error
	ReserveStockForSalesOrder(variantSKU string, quantity float64, salesOrderID, salesOrderNo string, userID string) (*models.StockReservation, error)
	ReleaseStockReservation(reservationID uint, reason string, userID string) error
	RecordPackageDeduction(reservationID uint, quantity float64, userID string) error
	RecordShipment(reservationID uint, userID string) error
	RecordInvoicing(reservationID uint, userID string) error
	RecordStockAdjustment(variantSKU string, quantity float64, adjustmentType, reason, userID string) error

	// Queries
	GetVariantStockSummary(variantSKU string) (*models.VariantStock, error)
	GetVariantMovementHistory(variantSKU string, offset, limit int) ([]models.VariantStockMovement, int64, error)
	GetReservationsForSalesOrder(salesOrderID string) ([]models.StockReservation, error)
	GetLowStockVariants(offset, limit int) ([]models.VariantStock, int64, error)
	GetAllVariantStocks(offset, limit int) ([]models.VariantStock, int64, error)
	GetAllVariantStocksByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error)
	GetAllRawMaterialVariantStocksByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error)

	// Damaged variant management
	MarkVariantAsDamaged(variantSKU string, quantity float64, reason, userID string) error
	GetDamagedVariants(offset, limit int) ([]models.VariantStock, int64, error)
	GetDamagedVariantsByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error)

	// Reconciliation
	SyncAggregateStock(productID string) error
}

type variantStockManagementService struct {
	variantStockRepo    repo.VariantStockRepository
	variantMovementRepo repo.VariantStockMovementRepository
	reservationRepo     repo.StockReservationRepository
	productStockRepo    repo.ProductStockRepository
	stockLedgerRepo     repo.StockLedgerRepository
	productRepo         repo.ProductRepository
	db                  *gorm.DB
}

func NewVariantStockManagementService(
	variantStockRepo repo.VariantStockRepository,
	variantMovementRepo repo.VariantStockMovementRepository,
	reservationRepo repo.StockReservationRepository,
	productStockRepo repo.ProductStockRepository,
	stockLedgerRepo repo.StockLedgerRepository,
	productRepo repo.ProductRepository,
	db *gorm.DB,
) VariantStockManagementService {
	return &variantStockManagementService{
		variantStockRepo:    variantStockRepo,
		variantMovementRepo: variantMovementRepo,
		reservationRepo:     reservationRepo,
		productStockRepo:    productStockRepo,
		stockLedgerRepo:     stockLedgerRepo,
		productRepo:         productRepo,
		db:                  db,
	}
}

func (s *variantStockManagementService) isRawMaterialVariant(productID string) bool {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return false
	}

	return product.IsRaw
}

// GetAllRawMaterialVariantStocksByUser returns raw-material variant stocks for a specific user
func (s *variantStockManagementService) GetAllRawMaterialVariantStocksByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error) {
	return s.variantStockRepo.GetAllByUserWithRawFilter(userID, offset, limit, true)
}

// InitializeVariantStock creates new variant stock entry
func (s *variantStockManagementService) InitializeVariantStock(
	productID, variantSKU, variantName, productName string,
	sellingPrice, costPrice float64,
) (*models.VariantStock, error) {
	// Check if variant already exists
	existing, _ := s.variantStockRepo.GetBySKU(variantSKU)
	if existing != nil {
		return existing, nil
	}

	stock := &models.VariantStock{
		ID:              uuid.New().String(),
		ProductID:       productID,
		VariantSKU:      variantSKU,
		VariantName:     variantName,
		ProductName:     productName,
		CurrentStock:    0,
		PurchasedStock:  0,
		SoldStock:       0,
		ReservedStock:   0,
		AvailableStock:  0,
		InTransitStock:  0,
		AverageCost:     costPrice,
		SellingPrice:    sellingPrice,
		ReorderLevel:    10,
		ReorderQty:      50,
		IsLowStock:      false,
		LastStockSyncAt: timePtr(time.Now()),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.variantStockRepo.Create(stock); err != nil {
		log.Printf("[VARIANT_STOCK] Error creating variant stock: %v", err)
		return nil, fmt.Errorf("failed to initialize variant stock: %w", err)
	}

	log.Printf("[VARIANT_STOCK] Initialized: SKU=%s, Name=%s, Cost=%.2f, Selling=%.2f", variantSKU, variantName, costPrice, sellingPrice)
	return stock, nil
}

// RecordPurchaseInbound records purchase order received
func (s *variantStockManagementService) RecordPurchaseInbound(
	productID, variantSKU string,
	quantity, rate float64,
	referenceType, referenceID, referenceNumber string,
	userID string,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	log.Printf("[VARIANT_PURCHASE] Recording: ProductID=%s, SKU=%s, Qty=%.2f, Rate=%.2f, Ref=%s", productID, variantSKU, quantity, rate, referenceNumber)

	// Get or create variant stock
	stock, err := s.variantStockRepo.GetBySKU(variantSKU)
	if err != nil {
		// Variant doesn't exist, create it
		log.Printf("[VARIANT_PURCHASE] Variant not found, creating new variant: SKU=%s", variantSKU)
		product, err := s.productRepo.FindByID(productID)
		if err != nil {
			log.Printf("[VARIANT_PURCHASE] Error: Product not found: %s", productID)
			return fmt.Errorf("product not found: %s", productID)
		}

		// Initialize variant stock
		sellingPrice := 0.0
		if product.SalesInfo.SellingPrice > 0 {
			sellingPrice = product.SalesInfo.SellingPrice
		}
		stock, err = s.InitializeVariantStock(
			productID,
			variantSKU,
			variantSKU, // variant name defaults to SKU
			product.Name,
			sellingPrice,
			rate, // cost price from the purchase order
		)
		if err != nil {
			log.Printf("[VARIANT_PURCHASE] Error initializing variant: %v", err)
			return err
		}
	}

	// Update stock quantities
	prevQty := stock.PurchasedStock
	stock.PurchasedStock += quantity

	// Calculate weighted average cost
	if prevQty > 0 {
		stock.AverageCost = ((stock.AverageCost * prevQty) + (rate * quantity)) / stock.PurchasedStock
	} else {
		stock.AverageCost = rate
	}

	// CurrentStock = PurchasedStock - SoldStock
	stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.LastPurchasedDate = timePtr(time.Now())
	stock.LastStockSyncAt = timePtr(time.Now())
	stock.UpdatedAt = time.Now()

	if err := s.variantStockRepo.Update(stock); err != nil {
		log.Printf("[VARIANT_PURCHASE] Error updating stock: %v", err)
		return fmt.Errorf("failed to update variant stock: %w", err)
	}

	// Record movement
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       variantSKU,
		MovementType:     "PURCHASE_ORDER",
		Quantity:         quantity,
		Rate:             rate,
		Amount:           quantity * rate,
		ReferenceType:    referenceType,
		ReferenceID:      referenceID,
		ReferenceNumber:  referenceNumber,
		BalanceBeforeQty: prevQty,
		BalanceAfterQty:  stock.CurrentStock,
		Stage:            "confirmed",
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	if err := s.variantMovementRepo.Create(movement); err != nil {
		log.Printf("[VARIANT_PURCHASE] Warning: Failed to record movement: %v", err)
	}

	log.Printf("[VARIANT_PURCHASE] Success: New balance=%.2f, Avg Cost=%.2f", stock.CurrentStock, stock.AverageCost)
	return nil
}

// ReserveStockForSalesOrder reserves stock for sales order
func (s *variantStockManagementService) ReserveStockForSalesOrder(
	variantSKU string,
	quantity float64,
	salesOrderID, salesOrderNo string,
	userID string,
) (*models.StockReservation, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	log.Printf("[VARIANT_RESERVE] Reserving: SKU=%s, Qty=%.2f, SO=%s", variantSKU, quantity, salesOrderNo)

	// Get variant stock
	stock, err := s.variantStockRepo.GetBySKU(variantSKU)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %s", variantSKU)
	}

	// Check availability
	if stock.AvailableStock < quantity {
		log.Printf("[VARIANT_RESERVE] Insufficient stock: Available=%.2f, Requested=%.2f", stock.AvailableStock, quantity)
		return nil, fmt.Errorf("insufficient stock: available=%.2f, requested=%.2f", stock.AvailableStock, quantity)
	}

	// Update stock
	stock.ReservedStock += quantity
	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.UpdatedAt = time.Now()

	if err := s.variantStockRepo.Update(stock); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	// Create reservation
	reservation := &models.StockReservation{
		SalesOrderID:    salesOrderID,
		SalesOrderNo:    salesOrderNo,
		ProductID:       stock.ProductID,
		VariantSKU:      variantSKU,
		VariantStockID:  stock.ID,
		ReservedQty:     quantity,
		ShippedQty:      0,
		InvoicedQty:     0,
		Status:          "reserved",
		ReservationDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       userID,
	}

	if err := s.reservationRepo.Create(reservation); err != nil {
		log.Printf("[VARIANT_RESERVE] Error creating reservation: %v", err)
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Record movement
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       variantSKU,
		MovementType:     "SALES_ORDER",
		Quantity:         -quantity,
		Rate:             stock.AverageCost,
		Amount:           -quantity * stock.AverageCost,
		ReferenceType:    "sales_order",
		ReferenceID:      salesOrderID,
		ReferenceNumber:  salesOrderNo,
		BalanceBeforeQty: stock.CurrentStock + quantity - stock.ReservedStock + quantity,
		BalanceAfterQty:  stock.AvailableStock,
		Stage:            "confirmed",
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	if err := s.variantMovementRepo.Create(movement); err != nil {
		log.Printf("[VARIANT_RESERVE] Warning: Failed to record movement: %v", err)
	}

	log.Printf("[VARIANT_RESERVE] Success: Reserved=%.2f, Available=%.2f", stock.ReservedStock, stock.AvailableStock)
	return reservation, nil
}

// ReleaseStockReservation releases a stock reservation
func (s *variantStockManagementService) ReleaseStockReservation(reservationID uint, reason string, userID string) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	stock, err := s.variantStockRepo.GetBySKU(reservation.VariantSKU)
	if err != nil {
		return fmt.Errorf("variant not found: %s", reservation.VariantSKU)
	}

	// Release reserved stock
	stock.ReservedStock -= reservation.ReservedQty
	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.UpdatedAt = time.Now()

	if err := s.variantStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Update reservation status
	reservation.Status = "cancelled"
	reservation.UpdatedAt = time.Now()

	if err := s.reservationRepo.Update(reservation); err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	log.Printf("[VARIANT_RELEASE] Released: SKU=%s, Qty=%.2f, Reason=%s", reservation.VariantSKU, reservation.ReservedQty, reason)
	return nil
}

// RecordPackageDeduction moves from reserved to in_transit
func (s *variantStockManagementService) RecordPackageDeduction(reservationID uint, quantity float64, userID string) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	stock, err := s.variantStockRepo.GetBySKU(reservation.VariantSKU)
	if err != nil {
		return fmt.Errorf("variant not found: %s", reservation.VariantSKU)
	}

	// Update quantities
	stock.ReservedStock -= quantity
	stock.InTransitStock += quantity
	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.UpdatedAt = time.Now()

	if err := s.variantStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	reservation.ShippedQty += quantity
	reservation.Status = "partial_shipped"
	if reservation.ShippedQty >= reservation.ReservedQty {
		reservation.Status = "fully_shipped"
	}
	reservation.UpdatedAt = time.Now()

	if err := s.reservationRepo.Update(reservation); err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	// Record movement
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       reservation.VariantSKU,
		MovementType:     "SHIPMENT",
		Quantity:         quantity,
		Rate:             stock.AverageCost,
		Amount:           quantity * stock.AverageCost,
		ReferenceType:    "package",
		ReferenceID:      fmt.Sprintf("%d", reservationID),
		BalanceBeforeQty: stock.InTransitStock - quantity,
		BalanceAfterQty:  stock.InTransitStock,
		Stage:            "packed",
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	s.variantMovementRepo.Create(movement)
	log.Printf("[VARIANT_PACKAGE] Packed: SKU=%s, Qty=%.2f, InTransit=%.2f", reservation.VariantSKU, quantity, stock.InTransitStock)
	return nil
}

// RecordShipment confirms shipment
func (s *variantStockManagementService) RecordShipment(reservationID uint, userID string) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	stock, err := s.variantStockRepo.GetBySKU(reservation.VariantSKU)
	if err != nil {
		return fmt.Errorf("variant not found: %s", reservation.VariantSKU)
	}

	// Record movement
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       reservation.VariantSKU,
		MovementType:     "SHIPMENT",
		Quantity:         reservation.ShippedQty,
		Rate:             stock.AverageCost,
		Amount:           reservation.ShippedQty * stock.AverageCost,
		ReferenceType:    "shipment",
		ReferenceID:      fmt.Sprintf("%d", reservationID),
		BalanceBeforeQty: stock.InTransitStock,
		BalanceAfterQty:  stock.InTransitStock,
		Stage:            "shipped",
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	s.variantMovementRepo.Create(movement)
	log.Printf("[VARIANT_SHIPMENT] Shipped: SKU=%s, Qty=%.2f", reservation.VariantSKU, reservation.ShippedQty)
	return nil
}

// RecordInvoicing finalizes sale
func (s *variantStockManagementService) RecordInvoicing(reservationID uint, userID string) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	stock, err := s.variantStockRepo.GetBySKU(reservation.VariantSKU)
	if err != nil {
		return fmt.Errorf("variant not found: %s", reservation.VariantSKU)
	}

	// Final deduction: increment sold stock
	stock.InTransitStock -= reservation.ShippedQty
	stock.SoldStock += reservation.ShippedQty
	// CurrentStock = PurchasedStock - SoldStock (automatic recalculation)
	stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.LastSoldDate = timePtr(time.Now())
	stock.UpdatedAt = time.Now()

	if err := s.variantStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Update reservation
	reservation.InvoicedQty = reservation.ShippedQty
	reservation.Status = "invoiced"
	reservation.UpdatedAt = time.Now()

	if err := s.reservationRepo.Update(reservation); err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	// Record movement
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       reservation.VariantSKU,
		MovementType:     "INVOICE",
		Quantity:         -reservation.ShippedQty,
		Rate:             stock.AverageCost,
		Amount:           -reservation.ShippedQty * stock.AverageCost,
		ReferenceType:    "invoice",
		ReferenceID:      reservation.SalesOrderID,
		ReferenceNumber:  reservation.SalesOrderNo,
		BalanceBeforeQty: stock.CurrentStock + reservation.ShippedQty,
		BalanceAfterQty:  stock.CurrentStock,
		Stage:            "invoiced",
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	s.variantMovementRepo.Create(movement)
	log.Printf("[VARIANT_INVOICE] Invoiced: SKU=%s, Qty=%.2f, Final Balance=%.2f", reservation.VariantSKU, reservation.ShippedQty, stock.CurrentStock)
	return nil
}

// RecordStockAdjustment records manual adjustments
func (s *variantStockManagementService) RecordStockAdjustment(
	variantSKU string,
	quantity float64,
	adjustmentType, reason, userID string,
) error {
	stock, err := s.variantStockRepo.GetBySKU(variantSKU)
	if err != nil {
		return fmt.Errorf("variant not found: %s", variantSKU)
	}

	var movementQty float64
	now := time.Now()
	if adjustmentType == "in" {
		stock.PurchasedStock += quantity
		// CurrentStock = PurchasedStock - SoldStock
		stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
		stock.LastPurchasedDate = &now
		movementQty = quantity
	} else if adjustmentType == "out" {
		if stock.AvailableStock < quantity {
			return errors.New("insufficient stock for adjustment out")
		}
		stock.SoldStock += quantity
		// CurrentStock = PurchasedStock - SoldStock
		stock.CurrentStock = stock.PurchasedStock - stock.SoldStock
		stock.LastSoldDate = &now
		movementQty = -quantity
	} else {
		return errors.New("invalid adjustment type")
	}

	stock.AvailableStock = stock.CurrentStock - stock.ReservedStock
	stock.UpdatedAt = time.Now()

	if err := s.variantStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Record movement
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       variantSKU,
		MovementType:     "ADJUSTMENT",
		Quantity:         movementQty,
		Rate:             stock.AverageCost,
		Amount:           movementQty * stock.AverageCost,
		ReferenceType:    "adjustment",
		ReferenceNumber:  reason,
		BalanceBeforeQty: stock.CurrentStock - movementQty,
		BalanceAfterQty:  stock.CurrentStock,
		Stage:            "confirmed",
		Notes:            reason,
		CreatedAt:        time.Now(),
		CreatedBy:        userID,
	}

	s.variantMovementRepo.Create(movement)
	log.Printf("[VARIANT_ADJUSTMENT] %s: SKU=%s, Qty=%.2f, Reason=%s", adjustmentType, variantSKU, quantity, reason)
	return nil
}

// GetVariantStockSummary returns current stock summary
func (s *variantStockManagementService) GetVariantStockSummary(variantSKU string) (*models.VariantStock, error) {
	return s.variantStockRepo.GetBySKU(variantSKU)
}

// GetVariantMovementHistory returns movement history
func (s *variantStockManagementService) GetVariantMovementHistory(variantSKU string, offset, limit int) ([]models.VariantStockMovement, int64, error) {
	return s.variantMovementRepo.GetByVariantSKU(variantSKU, offset, limit)
}

// GetReservationsForSalesOrder returns reservations for SO
func (s *variantStockManagementService) GetReservationsForSalesOrder(salesOrderID string) ([]models.StockReservation, error) {
	return s.reservationRepo.GetBySalesOrderID(salesOrderID)
}

// GetLowStockVariants returns low stock variants
func (s *variantStockManagementService) GetLowStockVariants(offset, limit int) ([]models.VariantStock, int64, error) {
	return s.variantStockRepo.GetLowStockVariants(50, offset, limit)
}

// GetAllVariantStocks returns all variant stocks with pagination
func (s *variantStockManagementService) GetAllVariantStocks(offset, limit int) ([]models.VariantStock, int64, error) {
	return s.variantStockRepo.GetAll(offset, limit)
}

// GetAllVariantStocksByUser returns all variant stocks for a specific user
func (s *variantStockManagementService) GetAllVariantStocksByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error) {
	return s.variantStockRepo.GetAllByUserWithRawFilter(userID, offset, limit, false)
}

// SyncAggregateStock syncs variant stock to product level
func (s *variantStockManagementService) SyncAggregateStock(productID string) error {
	// Get all variants for product
	variants, _, err := s.variantStockRepo.GetByProductID(productID, 0, 1000)
	if err != nil {
		return fmt.Errorf("failed to get variants: %w", err)
	}

	// Aggregate totals
	var totalCurrent, totalPurchased, totalSold, totalReserved, totalInTransit, totalAvailable float64
	var sumCost float64
	count := 0

	for _, v := range variants {
		totalCurrent += v.CurrentStock
		totalPurchased += v.PurchasedStock
		totalSold += v.SoldStock
		totalReserved += v.ReservedStock
		totalInTransit += v.InTransitStock
		totalAvailable += v.AvailableStock
		sumCost += v.AverageCost
		count++
	}

	var avgCost float64
	if count > 0 {
		avgCost = sumCost / float64(count)
	}

	// Get or create product stock
	productStock, _ := s.productStockRepo.GetByProductID(productID)
	if productStock == nil {
		productStock = &models.ProductStock{
			ID:        uuid.New().String(),
			ProductID: productID,
		}
	}

	productStock.CurrentStock = totalCurrent
	productStock.PurchasedStock = totalPurchased
	productStock.SoldStock = totalSold
	productStock.ReservedStock = totalReserved
	productStock.AvailableStock = totalAvailable
	productStock.AverageCost = avgCost
	productStock.UpdatedAt = time.Now()

	if err := s.productStockRepo.Update(productStock); err != nil {
		return fmt.Errorf("failed to sync product stock: %w", err)
	}

	log.Printf("[VARIANT_SYNC] Synced ProductID=%s, Variants=%d, Total=%.2f", productID, count, totalCurrent)
	return nil
}

// MarkVariantAsDamaged marks a variant stock as damaged and reduces available stock
func (s *variantStockManagementService) MarkVariantAsDamaged(variantSKU string, quantity float64, reason, userID string) error {
	if quantity <= 0 {
		return errors.New("damage quantity must be positive")
	}

	if reason == "" {
		return errors.New("damage reason is required")
	}

	log.Printf("[DAMAGE_TRACKING] Marking variant as damaged: SKU=%s, Qty=%.2f, Reason=%s",
		variantSKU, quantity, reason)

	// Get variant stock
	stock, err := s.variantStockRepo.GetBySKU(variantSKU)
	if err != nil {
		return fmt.Errorf("variant stock not found: %s", variantSKU)
	}

	// Check if available stock is sufficient
	if stock.AvailableStock < quantity {
		return fmt.Errorf("insufficient available stock: have %.2f, trying to mark %.2f as damaged",
			stock.AvailableStock, quantity)
	}

	// Update stock
	prevDamagedStock := stock.DamagedStock
	stock.DamagedStock += quantity
	stock.AvailableStock -= quantity
	stock.DamageReason = reason
	now := time.Now()
	stock.DamagedAt = &now
	stock.DamagedBy = userID
	stock.UpdatedAt = now

	if err := s.variantStockRepo.Update(stock); err != nil {
		return fmt.Errorf("failed to update variant stock with damage: %w", err)
	}

	// Record movement entry for damage
	amount := quantity * stock.AverageCost
	movement := &models.VariantStockMovement{
		VariantID:        stock.ID,
		ProductID:        stock.ProductID,
		VariantSKU:       stock.VariantSKU,
		MovementType:     "DAMAGE",
		Quantity:         -quantity, // Negative to indicate reduction
		Rate:             stock.AverageCost,
		Amount:           -amount, // Negative to indicate value reduction
		ReferenceType:    "damage_record",
		ReferenceID:      uuid.New().String(),
		ReferenceNumber:  fmt.Sprintf("DMG-%d", time.Now().Unix()),
		BalanceBeforeQty: stock.AvailableStock + quantity,
		BalanceAfterQty:  stock.AvailableStock,
		Notes:            fmt.Sprintf("Damage reason: %s", reason),
		CreatedAt:        now,
		CreatedBy:        userID,
	}

	if err := s.variantMovementRepo.Create(movement); err != nil {
		log.Printf("[DAMAGE_TRACKING] Warning: Failed to create movement entry: %v", err)
	}

	log.Printf("[DAMAGE_TRACKING] Success: Variant damaged stock updated from %.2f to %.2f units",
		prevDamagedStock, stock.DamagedStock)
	return nil
}

// GetDamagedVariants retrieves all variants with damaged stock
func (s *variantStockManagementService) GetDamagedVariants(offset, limit int) ([]models.VariantStock, int64, error) {
	return s.variantStockRepo.GetDamagedVariants(offset, limit)
}

// GetDamagedVariantsByUser retrieves damaged variants for a specific user
func (s *variantStockManagementService) GetDamagedVariantsByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error) {
	return s.variantStockRepo.GetDamagedVariantsByUser(userID, offset, limit)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
