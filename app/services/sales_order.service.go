package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type SalesOrderService interface {
	// Basic CRUD Operations
	CreateSalesOrder(soInput *input.CreateSalesOrderInput, userID string) (*output.SalesOrderOutput, error)
	GetSalesOrder(id string) (*output.SalesOrderOutput, error)
	GetAllSalesOrders(limit, offset int) ([]output.SalesOrderOutput, int64, error)
	GetSalesOrdersByCustomer(customerID uint, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	GetSalesOrdersByStatus(status string, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	UpdateSalesOrder(id string, soInput *input.UpdateSalesOrderInput, userID string) (*output.SalesOrderOutput, error)

	// Step 4: Selling to Customers (Outbound Operations)
	// Update SO status and manage inventory reservations when customer commits to purchase
	UpdateSalesOrderStatus(id string, status string, userID string) (*output.SalesOrderOutput, error)
	DeleteSalesOrder(id string) error
}

type salesOrderService struct {
	soRepo             repo.SalesOrderRepository
	customerRepo       repo.CustomerRepository
	itemRepo           repo.ItemRepository
	taxRepo            repo.TaxRepository
	salespersonRepo    repo.SalespersonRepository
	inventoryRepo      repo.InventoryBalanceRepository
	stockMovementSvc   StockMovementService
	productStockRepo   repo.ProductStockRepository
	stockLedgerRepo    repo.StockLedgerRepository
	variantStockMgmt   VariantStockManagementService
	stockManagementSvc StockManagementService
}

func NewSalesOrderService(
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	itemRepo repo.ItemRepository,
	taxRepo repo.TaxRepository,
	salespersonRepo repo.SalespersonRepository,
	inventoryRepo repo.InventoryBalanceRepository,
	stockMovementSvc StockMovementService,
	productStockRepo repo.ProductStockRepository,
	stockLedgerRepo repo.StockLedgerRepository,
	variantStockMgmt VariantStockManagementService,
	stockManagementSvc StockManagementService,
) SalesOrderService {
	return &salesOrderService{
		soRepo:             soRepo,
		customerRepo:       customerRepo,
		itemRepo:           itemRepo,
		taxRepo:            taxRepo,
		salespersonRepo:    salespersonRepo,
		inventoryRepo:      inventoryRepo,
		stockMovementSvc:   stockMovementSvc,
		productStockRepo:   productStockRepo,
		stockLedgerRepo:    stockLedgerRepo,
		variantStockMgmt:   variantStockMgmt,
		stockManagementSvc: stockManagementSvc,
	}
}

func (s *salesOrderService) CreateSalesOrder(soInput *input.CreateSalesOrderInput, userID string) (*output.SalesOrderOutput, error) {
	lineItems := make([]models.SalesOrderLineItem, 0)
	subTotal := 0.0

	for _, itemInput := range soInput.LineItems {
		// Validate required fields
		if itemInput.ProductID == "" {
			return nil, errors.New("product_id is required for each line item")
		}
		if itemInput.ProductName == "" {
			return nil, errors.New("product_name is required for each line item")
		}

		// Calculate line item totals
		lineAmount := itemInput.Quantity * itemInput.Rate
		subTotal += lineAmount

		lineItem := models.SalesOrderLineItem{
			ProductID:      itemInput.ProductID,
			ProductName:    itemInput.ProductName,
			SKU:            itemInput.SKU,
			Account:        itemInput.Account,
			Quantity:       itemInput.Quantity,
			Rate:           itemInput.Rate,
			Amount:         lineAmount,
			VariantSKU:     itemInput.VariantSKU,
			VariantDetails: models.JSONB(itemInput.VariantDetails),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		lineItems = append(lineItems, lineItem)
	}

	// Calculate tax based on tax_rate
	taxAmount := (subTotal * soInput.TaxRate) / 100

	// Calculate total with shipping charges and adjustment
	total := subTotal + taxAmount + soInput.ShippingCharges + soInput.Adjustment

	// Generate SO number
	soNumber := fmt.Sprintf("SO-%s-%04d", time.Now().Format("2006"), s.generateSOSequence())

	so := &models.SalesOrder{
		ID:                   uuid.New().String(),
		SalesOrderNumber:     soNumber,
		CustomerID:           soInput.CustomerID,
		ReferenceNo:          soInput.ReferenceNo,
		Date:                 soInput.SalesOrderDate,
		ExpectedShipmentDate: soInput.ExpectedShipmentDate,
		DeliveryMethod:       soInput.DeliveryMethod,
		PaymentTerms:         domain.PaymentTerms(soInput.PaymentTerms),
		LineItems:            lineItems,
		SubTotal:             subTotal,
		ShippingCharges:      soInput.ShippingCharges,
		Adjustment:           soInput.Adjustment,
		TaxID:                soInput.TaxID,
		TaxRate:              soInput.TaxRate,
		TaxTotal:             taxAmount,
		Total:                total,
		CustomerNotes:        soInput.CustomerNotes,
		TermsAndConditions:   soInput.TermsAndConditions,
		SalespersonID:        soInput.SalespersonID,
		Status:               domain.SalesOrderStatusDraft,
		InventoryReserved:    false,
		InventoryDeducted:    false,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CreatedBy:            userID,
	}

	createdSO, err := s.soRepo.Create(so)
	if err != nil {
		return nil, errors.New("failed to create sales order: " + err.Error())
	}

	return output.ToSalesOrderOutput(createdSO)
}

func (s *salesOrderService) GetSalesOrder(id string) (*output.SalesOrderOutput, error) {
	so, err := s.soRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("sales order not found")
	}
	return output.ToSalesOrderOutput(so)
}

func (s *salesOrderService) GetAllSalesOrders(limit, offset int) ([]output.SalesOrderOutput, int64, error) {
	sos, total, err := s.soRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.SalesOrderOutput, len(sos))
	for i, so := range sos {
		out, _ := output.ToSalesOrderOutput(&so)
		outputs[i] = *out
	}

	return outputs, total, nil
}

func (s *salesOrderService) GetSalesOrdersByCustomer(customerID uint, limit, offset int) ([]output.SalesOrderOutput, int64, error) {
	sos, total, err := s.soRepo.FindByCustomer(customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.SalesOrderOutput, len(sos))
	for i, so := range sos {
		out, _ := output.ToSalesOrderOutput(&so)
		outputs[i] = *out
	}

	return outputs, total, nil
}

func (s *salesOrderService) GetSalesOrdersByStatus(status string, limit, offset int) ([]output.SalesOrderOutput, int64, error) {
	sos, total, err := s.soRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.SalesOrderOutput, len(sos))
	for i, so := range sos {
		out, _ := output.ToSalesOrderOutput(&so)
		outputs[i] = *out
	}

	return outputs, total, nil
}

func (s *salesOrderService) UpdateSalesOrder(id string, soInput *input.UpdateSalesOrderInput, userID string) (*output.SalesOrderOutput, error) {
	so, err := s.soRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("sales order not found")
	}

	if soInput.CustomerID != nil {
		so.CustomerID = *soInput.CustomerID
	}

	if soInput.ReferenceNo != nil {
		so.ReferenceNo = *soInput.ReferenceNo
	}

	if soInput.SalesOrderDate != nil {
		so.Date = *soInput.SalesOrderDate
	}

	if soInput.ExpectedShipmentDate != nil {
		so.ExpectedShipmentDate = *soInput.ExpectedShipmentDate
	}

	if soInput.PaymentTerms != nil {
		so.PaymentTerms = domain.PaymentTerms(*soInput.PaymentTerms)
	}

	if soInput.DeliveryMethod != nil {
		so.DeliveryMethod = *soInput.DeliveryMethod
	}

	if soInput.ShippingCharges != nil {
		so.ShippingCharges = *soInput.ShippingCharges
	}

	if soInput.Adjustment != nil {
		so.Adjustment = *soInput.Adjustment
	}

	if soInput.TaxID != nil {
		so.TaxID = soInput.TaxID
	}

	if soInput.TaxRate != nil {
		so.TaxRate = *soInput.TaxRate
	}

	if soInput.CustomerNotes != nil {
		so.CustomerNotes = *soInput.CustomerNotes
	}

	if soInput.TermsAndConditions != nil {
		so.TermsAndConditions = *soInput.TermsAndConditions
	}

	if soInput.SalespersonID != nil {
		so.SalespersonID = soInput.SalespersonID
	}

	if len(soInput.LineItems) > 0 {
		lineItems := make([]models.SalesOrderLineItem, 0)
		subTotal := 0.0

		for _, itemInput := range soInput.LineItems {
			if itemInput.ProductID == "" {
				return nil, errors.New("product_id is required for each line item")
			}

			lineAmount := itemInput.Quantity * itemInput.Rate
			subTotal += lineAmount

			lineItem := models.SalesOrderLineItem{
				ProductID:      itemInput.ProductID,
				ProductName:    itemInput.ProductName,
				SKU:            itemInput.SKU,
				Account:        itemInput.Account,
				Quantity:       itemInput.Quantity,
				Rate:           itemInput.Rate,
				Amount:         lineAmount,
				VariantSKU:     itemInput.VariantSKU,
				VariantDetails: models.JSONB(itemInput.VariantDetails),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			lineItems = append(lineItems, lineItem)
		}

		so.LineItems = lineItems
		so.SubTotal = subTotal

		// Recalculate tax with the current tax rate
		taxAmount := (subTotal * so.TaxRate) / 100
		so.TaxTotal = taxAmount
		so.Total = subTotal + taxAmount + so.ShippingCharges + so.Adjustment
	}

	so.UpdatedAt = time.Now()
	so.UpdatedBy = userID

	updatedSO, err := s.soRepo.Update(id, so)
	if err != nil {
		return nil, errors.New("failed to update sales order: " + err.Error())
	}

	return output.ToSalesOrderOutput(updatedSO)
}

func (s *salesOrderService) UpdateSalesOrderStatus(id string, status string, userID string) (*output.SalesOrderOutput, error) {
	so, err := s.soRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("sales order not found")
	}

	so.Status = domain.SalesOrderStatus(status)
	so.UpdatedAt = time.Now()
	so.UpdatedBy = userID

	// Deduct inventory when status is "paid" or "delivered" and not already deducted
	if (status == string(domain.SalesOrderStatusPaid) || status == string(domain.SalesOrderStatusDelivered)) && !so.InventoryDeducted {
		// Deduct stock for each line item
		for _, lineItem := range so.LineItems {
			deducted := false

			// Check if this is a variant-based product (has VariantSKU)
			if lineItem.VariantSKU != "" {
				// For variant-based products, try to deduct from variant stock first
				variantStock, err := s.variantStockMgmt.GetVariantStockSummary(lineItem.VariantSKU)
				if err == nil && variantStock != nil && variantStock.AvailableStock >= lineItem.Quantity {
					// Variant stock exists and has sufficient quantity - perform the actual deduction
					adjustErr := s.variantStockMgmt.RecordStockAdjustment(
						lineItem.VariantSKU,
						lineItem.Quantity,
						"out",
						fmt.Sprintf("Stock deducted for Sales Order: %s", so.SalesOrderNumber),
						userID,
					)
					if adjustErr != nil {
						log.Printf("[SO_STOCK_DEDUCTION] Error deducting variant stock for SKU=%s: %v", lineItem.VariantSKU, adjustErr)
						return nil, fmt.Errorf("failed to deduct variant stock for SKU %s: %v", lineItem.VariantSKU, adjustErr)
					}
					log.Printf("[SO_STOCK_DEDUCTION] Successfully deducted %.2f units of variant %s", lineItem.Quantity, lineItem.VariantSKU)
					deducted = true
				} else if err != nil || variantStock == nil {
					log.Printf("[SO_STOCK_DEDUCTION] Warning: Could not find variant stock for SKU=%s, will attempt product-level deduction", lineItem.VariantSKU)
				} else if variantStock.AvailableStock < lineItem.Quantity {
					log.Printf("[SO_STOCK_DEDUCTION] Insufficient variant stock: SKU=%s, Available=%.2f, Required=%.2f",
						lineItem.VariantSKU, variantStock.AvailableStock, lineItem.Quantity)
					return nil, fmt.Errorf("insufficient stock for variant %s: available=%.2f, required=%.2f",
						lineItem.VariantSKU, variantStock.AvailableStock, lineItem.Quantity)
				}
			}

			// If not deducted at variant level, try product-level stock deduction
			if !deducted {
				err := s.stockManagementSvc.RecordOutboundMovement(
					lineItem.ProductID,
					"SALES_ORDER",
					id,
					so.SalesOrderNumber,
					lineItem.Quantity,
					fmt.Sprintf("Stock deducted for Sales Order: %s", so.SalesOrderNumber),
					userID,
				)
				if err != nil {
					// If this is a variant-based product and product stock doesn't exist, that's OK
					if lineItem.VariantSKU != "" && err.Error() == fmt.Sprintf("no stock found for product: %s", lineItem.ProductID) {
						log.Printf("[SO_STOCK_DEDUCTION] Note: Variant-based product %s has no aggregated ProductStock, stock deduction handled at variant level", lineItem.ProductID)
					} else {
						log.Printf("[SO_STOCK_DEDUCTION] Error deducting stock for product %s: %v", lineItem.ProductID, err)
						return nil, fmt.Errorf("failed to deduct stock for product %s: %v", lineItem.ProductID, err)
					}
				} else {
					log.Printf("[SO_STOCK_DEDUCTION] Successfully deducted %.2f units of product %s", lineItem.Quantity, lineItem.ProductID)
				}
			}
		}

		// Mark inventory as deducted
		now := time.Now()
		so.InventoryDeducted = true
		so.DeductedDate = &now
		log.Printf("[SO_STOCK_DEDUCTION] Inventory deducted for Sales Order: %s", so.SalesOrderNumber)
	}

	err = s.soRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, errors.New("failed to update status: " + err.Error())
	}

	return s.GetSalesOrder(id)
}

func (s *salesOrderService) DeleteSalesOrder(id string) error {
	so, err := s.soRepo.FindByID(id)
	if err != nil {
		return errors.New("sales order not found")
	}

	log.Printf("[SO_DELETE] Deleted SO: %s", so.SalesOrderNumber)
	return s.soRepo.Delete(id)
}

func (s *salesOrderService) generateSOSequence() int {
	year := time.Now().Format("2006")
	var count int64
	db := s.soRepo.GetDB()
	lockName := fmt.Sprintf("so_seq_%s", year)

	// Acquire a named lock to ensure only one goroutine increments the counter at a time
	var lockResult int
	result := db.Raw("SELECT GET_LOCK(?, 30)", lockName).Scan(&lockResult)
	if result.Error != nil || lockResult != 1 {
		log.Printf("[SO_SEQUENCE] Warning: Failed to acquire lock for sequence generation: %v, lockResult=%d", result.Error, lockResult)
	}

	// Count ALL SOs for this year (not just today) to ensure globally unique sequence numbers
	// This prevents duplicate numbers when no SOs are created on some days
	db.Where("YEAR(created_at) = ?", year).Model(&models.SalesOrder{}).Count(&count)

	sequence := int(count) + 1

	// Release the lock
	db.Raw("SELECT RELEASE_LOCK(?)", lockName).Scan(&lockResult)

	log.Printf("[SO_SEQUENCE] Generated SO sequence number %d for year %s", sequence, year)
	return sequence
}
