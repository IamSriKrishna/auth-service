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

	// Company-scoped operations used by authenticated routes.
	CreateSalesOrderForCompany(soInput *input.CreateSalesOrderInput, userID string, companyID uint) (*output.SalesOrderOutput, error)
	GetSalesOrderForCompany(id string, companyID uint) (*output.SalesOrderOutput, error)
	GetAllSalesOrdersForCompany(companyID uint, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	GetSalesOrdersByCustomerForCompany(customerID, companyID uint, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	GetSalesOrdersByStatusForCompany(status string, companyID uint, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	UpdateSalesOrderForCompany(id string, soInput *input.UpdateSalesOrderInput, userID string, companyID uint) (*output.SalesOrderOutput, error)
	UpdateSalesOrderStatusForCompany(id, status, userID string, companyID uint) (*output.SalesOrderOutput, error)
	DeleteSalesOrderForCompany(id string, companyID uint) error
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
	pgInventoryService ProductGroupInventoryService
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
	pgInventoryService ProductGroupInventoryService,
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
		pgInventoryService: pgInventoryService,
	}
}

func (s *salesOrderService) CreateSalesOrder(soInput *input.CreateSalesOrderInput, userID string) (*output.SalesOrderOutput, error) {
	lineItems := make([]models.SalesOrderLineItem, 0)
	subTotal := 0.0

	for _, itemInput := range soInput.LineItems {
		// Validate required fields
		if itemInput.ManufacturerID == "" {
			return nil, errors.New("manufacturer_id is required for each line item")
		}
		if itemInput.ManufacturerName == "" {
			return nil, errors.New("manufacturer_name is required for each line item")
		}

		// Calculate line item totals
		lineAmount := itemInput.Quantity * itemInput.Rate
		subTotal += lineAmount

		lineItem := models.SalesOrderLineItem{
			ManufacturerID:   itemInput.ManufacturerID,
			ManufacturerName: itemInput.ManufacturerName,
			Account:          itemInput.Account,
			Quantity:         itemInput.Quantity,
			Rate:             itemInput.Rate,
			Amount:           lineAmount,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
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

	// Stock is NOT deducted at creation time
	// Stock will be deducted when the sales order status is changed to "paid"
	// This ensures stock is only committed when payment is confirmed

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
			if itemInput.ManufacturerID == "" {
				return nil, errors.New("manufacturer_id is required for each line item")
			}

			lineAmount := itemInput.Quantity * itemInput.Rate
			subTotal += lineAmount

			lineItem := models.SalesOrderLineItem{
				ManufacturerID:   itemInput.ManufacturerID,
				ManufacturerName: itemInput.ManufacturerName,
				Account:          itemInput.Account,
				Quantity:         itemInput.Quantity,
				Rate:             itemInput.Rate,
				Amount:           lineAmount,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
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

	oldStatus := so.Status
	so.Status = domain.SalesOrderStatus(status)
	so.UpdatedAt = time.Now()
	so.UpdatedBy = userID

	// When status changes to "paid", deduct stock from inventory
	if status == "paid" && oldStatus != domain.SalesOrderStatusPaid {
		// Only deduct once
		if !so.InventoryDeducted {
			log.Printf("[SO_PAID] Processing stock deduction for sales order: %s", so.SalesOrderNumber)

			// Deduct stock for each manufacturer in the sales order
			for _, lineItem := range so.LineItems {
				if s.variantStockMgmt != nil {
					err := s.variantStockMgmt.RecordStockAdjustment(
						lineItem.ManufacturerID, // Use manufacturer ID as SKU
						lineItem.Quantity,       // Positive quantity for "out"
						"out",
						fmt.Sprintf("Sold via Sales Order: %s", so.SalesOrderNumber),
						userID,
					)
					if err != nil {
						log.Printf("[SO_PAID] Warning: Failed to deduct stock for manufacturer %s: %v", lineItem.ManufacturerID, err)
						// Don't fail the status update if stock deduction fails - stock may have already been adjusted
						// This ensures payment processing continues even if there's a stock sync issue
					} else {
						log.Printf("[SO_PAID] Deducted stock: %s -%.2f units", lineItem.ManufacturerID, lineItem.Quantity)
					}
				}
			}

			// Mark inventory as deducted
			so.InventoryDeducted = true
			now := time.Now()
			so.DeductedDate = &now
		}
	}

	err = s.soRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, errors.New("failed to update status: " + err.Error())
	}

	// Update InventoryDeducted flag if it changed
	if so.InventoryDeducted {
		_, err = s.soRepo.Update(id, so)
		if err != nil {
			log.Printf("[SO_PAID] Warning: Failed to update InventoryDeducted flag: %v", err)
		}
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

func (s *salesOrderService) validateSalesOrderCompanyInput(
	soInput *input.CreateSalesOrderInput,
	companyID uint,
) error {
	if soInput == nil {
		return errors.New("input cannot be nil")
	}
	if companyID == 0 {
		return errors.New("invalid company")
	}

	if _, err := s.customerRepo.FindByIDAndCompany(soInput.CustomerID, companyID); err != nil {
		return errors.New("customer not found in your company")
	}

	for _, lineItem := range soInput.LineItems {
		if lineItem.ManufacturerID == "" {
			return errors.New("manufacturer_id is required for each line item")
		}

		var count int64
		if err := s.soRepo.GetDB().Model(&models.Manufacturer{}).
			Where("id = ? AND company_id = ?", lineItem.ManufacturerID, companyID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to validate manufacturer %s: %w", lineItem.ManufacturerID, err)
		}
		if count == 0 {
			return fmt.Errorf("manufacturer %s not found in your company", lineItem.ManufacturerID)
		}
	}

	if soInput.SalespersonID != nil {
		var count int64
		query := s.soRepo.GetDB().Model(&models.Salesperson{}).Where("id = ?", *soInput.SalespersonID)
		// Apply company filter only when the salespersons table has company_id.
		if s.soRepo.GetDB().Migrator().HasColumn(&models.Salesperson{}, "company_id") {
			query = query.Where("company_id = ?", companyID)
		}
		if err := query.Count(&count).Error; err != nil {
			return fmt.Errorf("failed to validate salesperson: %w", err)
		}
		if count == 0 {
			return errors.New("salesperson not found in your company")
		}
	}

	return nil
}

func (s *salesOrderService) validateSalesOrderCompanyUpdate(
	soInput *input.UpdateSalesOrderInput,
	companyID uint,
) error {
	if soInput == nil {
		return errors.New("input cannot be nil")
	}

	if soInput.CustomerID != nil {
		if _, err := s.customerRepo.FindByIDAndCompany(*soInput.CustomerID, companyID); err != nil {
			return errors.New("customer not found in your company")
		}
	}

	for _, lineItem := range soInput.LineItems {
		if lineItem.ManufacturerID == "" {
			return errors.New("manufacturer_id is required for each line item")
		}

		var count int64
		if err := s.soRepo.GetDB().Model(&models.Manufacturer{}).
			Where("id = ? AND company_id = ?", lineItem.ManufacturerID, companyID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to validate manufacturer %s: %w", lineItem.ManufacturerID, err)
		}
		if count == 0 {
			return fmt.Errorf("manufacturer %s not found in your company", lineItem.ManufacturerID)
		}
	}

	return nil
}

func (s *salesOrderService) CreateSalesOrderForCompany(
	soInput *input.CreateSalesOrderInput,
	userID string,
	companyID uint,
) (*output.SalesOrderOutput, error) {
	if err := s.validateSalesOrderCompanyInput(soInput, companyID); err != nil {
		return nil, err
	}
	return s.CreateSalesOrder(soInput, userID)
}

func (s *salesOrderService) GetSalesOrderForCompany(
	id string,
	companyID uint,
) (*output.SalesOrderOutput, error) {
	so, err := s.soRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("sales order not found")
	}
	return output.ToSalesOrderOutput(so)
}

func (s *salesOrderService) GetAllSalesOrdersForCompany(
	companyID uint,
	limit, offset int,
) ([]output.SalesOrderOutput, int64, error) {
	sos, total, err := s.soRepo.FindAllByCompany(companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return salesOrdersToOutput(sos, total)
}

func (s *salesOrderService) GetSalesOrdersByCustomerForCompany(
	customerID, companyID uint,
	limit, offset int,
) ([]output.SalesOrderOutput, int64, error) {
	if _, err := s.customerRepo.FindByIDAndCompany(customerID, companyID); err != nil {
		return nil, 0, errors.New("customer not found in your company")
	}

	sos, total, err := s.soRepo.FindByCustomerAndCompany(customerID, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return salesOrdersToOutput(sos, total)
}

func (s *salesOrderService) GetSalesOrdersByStatusForCompany(
	status string,
	companyID uint,
	limit, offset int,
) ([]output.SalesOrderOutput, int64, error) {
	sos, total, err := s.soRepo.FindByStatusAndCompany(status, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return salesOrdersToOutput(sos, total)
}

func (s *salesOrderService) UpdateSalesOrderForCompany(
	id string,
	soInput *input.UpdateSalesOrderInput,
	userID string,
	companyID uint,
) (*output.SalesOrderOutput, error) {
	if _, err := s.soRepo.FindByIDAndCompany(id, companyID); err != nil {
		return nil, errors.New("sales order not found")
	}
	if err := s.validateSalesOrderCompanyUpdate(soInput, companyID); err != nil {
		return nil, err
	}

	// Existing update logic is retained after the company ownership check.
	return s.UpdateSalesOrder(id, soInput, userID)
}

func (s *salesOrderService) UpdateSalesOrderStatusForCompany(
	id, status, userID string,
	companyID uint,
) (*output.SalesOrderOutput, error) {
	if _, err := s.soRepo.FindByIDAndCompany(id, companyID); err != nil {
		return nil, errors.New("sales order not found")
	}

	// Existing status and stock deduction logic is retained.
	return s.UpdateSalesOrderStatus(id, status, userID)
}

func (s *salesOrderService) DeleteSalesOrderForCompany(
	id string,
	companyID uint,
) error {
	if _, err := s.soRepo.FindByIDAndCompany(id, companyID); err != nil {
		return errors.New("sales order not found")
	}
	return s.soRepo.DeleteByCompany(id, companyID)
}

func salesOrdersToOutput(
	salesOrders []models.SalesOrder,
	total int64,
) ([]output.SalesOrderOutput, int64, error) {
	outputs := make([]output.SalesOrderOutput, len(salesOrders))
	for i := range salesOrders {
		out, err := output.ToSalesOrderOutput(&salesOrders[i])
		if err != nil {
			return nil, 0, err
		}
		outputs[i] = *out
	}
	return outputs, total, nil
}
