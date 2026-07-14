package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type InvoiceService interface {
	// Existing methods retained.
	CreateInvoice(input *input.CreateInvoiceInput, userID string) (*output.InvoiceOutput, error)
	GetInvoice(id string) (*output.InvoiceOutput, error)
	GetAllInvoices(limit, offset int) (*output.InvoiceListOutput, error)
	UpdateInvoice(id string, input *input.UpdateInvoiceInput, userID string) (*output.InvoiceOutput, error)
	DeleteInvoice(id string) error
	GetInvoicesByCustomer(customerID string, limit, offset int) (*output.InvoiceListOutput, error)
	GetInvoicesByStatus(status string, limit, offset int) (*output.InvoiceListOutput, error)
	UpdateInvoiceStatus(id string, status domain.InvoiceStatus) (*output.InvoiceOutput, error)

	// Company-scoped methods.
	CreateInvoiceForCompany(
		input *input.CreateInvoiceInput,
		userID string,
		companyID uint,
	) (*output.InvoiceOutput, error)

	GetInvoiceByCompany(
		id string,
		companyID uint,
	) (*output.InvoiceOutput, error)

	GetAllInvoicesByCompany(
		companyID uint,
		limit int,
		offset int,
	) (*output.InvoiceListOutput, error)

	UpdateInvoiceForCompany(
		id string,
		input *input.UpdateInvoiceInput,
		userID string,
		companyID uint,
	) (*output.InvoiceOutput, error)

	DeleteInvoiceForCompany(
		id string,
		companyID uint,
	) error

	GetInvoicesByCustomerAndCompany(
		customerID uint,
		companyID uint,
		limit int,
		offset int,
	) (*output.InvoiceListOutput, error)

	GetInvoicesByStatusAndCompany(
		status string,
		companyID uint,
		limit int,
		offset int,
	) (*output.InvoiceListOutput, error)

	UpdateInvoiceStatusForCompany(
		id string,
		status domain.InvoiceStatus,
		userID string,
		companyID uint,
	) (*output.InvoiceOutput, error)

	ValidateInvoiceInputForCompany(
		input *input.CreateInvoiceInput,
		companyID uint,
	) error
}

type SalespersonService interface {
	CreateSalesperson(input *input.CreateSalespersonInput) (*output.SalespersonOutput, error)
	GetSalesperson(id uint) (*output.SalespersonOutput, error)
	GetAllSalespersons(limit, offset int) (*output.SalespersonListOutput, error)
	UpdateSalesperson(id uint, input *input.UpdateSalespersonInput) (*output.SalespersonOutput, error)
	DeleteSalesperson(id uint) error
}

type TaxService interface {
	CreateTax(input *input.CreateTaxInput) (*output.TaxOutput, error)
	GetTax(id uint) (*output.TaxOutput, error)
	GetAllTaxes(limit, offset int) (*output.TaxListOutput, error)
	UpdateTax(id uint, input *input.UpdateTaxInput) (*output.TaxOutput, error)
	DeleteTax(id uint) error
}

type PaymentService interface {
	// Existing methods retained.
	CreatePayment(input *input.CreatePaymentInput, userID string) (*output.PaymentOutput, error)
	GetPayment(id uint) (*output.PaymentOutput, error)
	GetPaymentsByInvoice(invoiceID string) (*output.PaymentListOutput, error)
	DeletePayment(id uint) error

	// Company-scoped methods.
	CreatePaymentForCompany(
		input *input.CreatePaymentInput,
		userID string,
		companyID uint,
	) (*output.PaymentOutput, error)

	GetPaymentForCompany(
		id uint,
		companyID uint,
	) (*output.PaymentOutput, error)

	GetPaymentsByInvoiceForCompany(
		invoiceID string,
		companyID uint,
	) (*output.PaymentListOutput, error)

	DeletePaymentForCompany(
		id uint,
		companyID uint,
	) error
}

type invoiceService struct {
	invoiceRepo      repo.InvoiceRepository
	itemRepo         repo.ItemRepository
	customerRepo     repo.CustomerRepository
	salespersonRepo  repo.SalespersonRepository
	taxRepo          repo.TaxRepository
	paymentRepo      repo.PaymentRepository
	productRepo      repo.ProductRepository
	productStockRepo repo.ProductStockRepository
	stockLedgerRepo  repo.StockLedgerRepository
	userRepo         repo.UserRepository
}

func NewInvoiceService(
	invoiceRepo repo.InvoiceRepository,
	itemRepo repo.ItemRepository,
	customerRepo repo.CustomerRepository,
	salespersonRepo repo.SalespersonRepository,
	taxRepo repo.TaxRepository,
	paymentRepo repo.PaymentRepository,
	productRepo repo.ProductRepository,
	productStockRepo repo.ProductStockRepository,
	stockLedgerRepo repo.StockLedgerRepository,
	userRepo repo.UserRepository,
	pdfOutputDir string,
) InvoiceService {
	return &invoiceService{
		invoiceRepo:      invoiceRepo,
		itemRepo:         itemRepo,
		customerRepo:     customerRepo,
		salespersonRepo:  salespersonRepo,
		taxRepo:          taxRepo,
		paymentRepo:      paymentRepo,
		productRepo:      productRepo,
		productStockRepo: productStockRepo,
		stockLedgerRepo:  stockLedgerRepo,
		userRepo:         userRepo,
	}
}

func (s *invoiceService) CreateInvoice(input *input.CreateInvoiceInput, userID string) (*output.InvoiceOutput, error) {
	_, err := s.customerRepo.FindByID(input.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	if input.SalespersonID != nil {
		_, err := s.salespersonRepo.FindByID(*input.SalespersonID)
		if err != nil {
			return nil, fmt.Errorf("salesperson not found")
		}
	}

	var tax *models.Tax
	if input.TaxID != nil {
		tax, err = s.taxRepo.FindByID(*input.TaxID)
		if err != nil {
			return nil, fmt.Errorf("tax not found")
		}
	}

	id := fmt.Sprintf("inv_%s", uuid.New().String()[:8])
	invoiceNumber, err := s.invoiceRepo.GetNextInvoiceNumber()
	if err != nil {
		return nil, err
	}

	lineItems := make([]models.InvoiceLineItem, len(input.LineItems))
	var subTotal float64

	for i, itemInput := range input.LineItems {
		lineItemProductID := itemInput.ProductID
		if lineItemProductID != nil && *lineItemProductID != "" {
			if _, err := s.productRepo.FindByID(*lineItemProductID); err != nil {
				lineItemProductID = nil
			}
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		variantDetails := models.VariantDetails{
			"sku":     itemInput.SKU,
			"account": itemInput.Account,
		}

		lineItems[i] = models.InvoiceLineItem{
			ProductID:      lineItemProductID,
			ProductName:    itemInput.ProductName,
			Description:    itemInput.ProductName,
			Quantity:       itemInput.Quantity,
			Rate:           itemInput.Rate,
			Amount:         amount,
			VariantDetails: variantDetails,
		}
	}

	var taxAmount float64
	if tax != nil {
		taxAmount = (subTotal + input.ShippingCharges) * tax.Rate / 100
	}

	total := subTotal + input.ShippingCharges + taxAmount + input.Adjustment

	invoice := &models.Invoice{
		ID:                 id,
		InvoiceNumber:      invoiceNumber,
		SalesOrderID:       input.SalesOrderID,
		CustomerID:         input.CustomerID,
		OrderNumber:        input.OrderNumber,
		InvoiceDate:        input.InvoiceDate,
		Terms:              domain.PaymentTerms(input.Terms),
		DueDate:            input.DueDate,
		SalespersonID:      input.SalespersonID,
		Subject:            input.Subject,
		LineItems:          lineItems,
		SubTotal:           subTotal,
		ShippingCharges:    input.ShippingCharges,
		TaxID:              input.TaxID,
		TaxAmount:          taxAmount,
		Adjustment:         input.Adjustment,
		Total:              total,
		CustomerNotes:      input.CustomerNotes,
		TermsAndConditions: input.TermsAndConditions,
		Status:             domain.InvoiceStatusDraft,
		Attachments:        input.Attachments,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          userID,
	}

	if input.TaxType != nil {
		invoice.TaxType = domain.TaxType(*input.TaxType)
	}

	if err := s.invoiceRepo.Create(invoice); err != nil {
		return nil, err
	}

	createdInvoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	invoiceOutput, err := output.ToInvoiceOutput(createdInvoice)
	if err != nil {
		return nil, err
	}

	return invoiceOutput, nil
}

func (s *invoiceService) GetInvoice(id string) (*output.InvoiceOutput, error) {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceOutput(invoice)
}

func (s *invoiceService) GetAllInvoices(limit, offset int) (*output.InvoiceListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	invoices, total, err := s.invoiceRepo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceListOutput(invoices, total)
}

func (s *invoiceService) UpdateInvoice(id string, input *input.UpdateInvoiceInput, userID string) (*output.InvoiceOutput, error) {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if invoice.Status == domain.InvoiceStatusSent || invoice.Status == domain.InvoiceStatusPaid {
		return nil, fmt.Errorf("cannot update invoice with status %s", invoice.Status)
	}

	if input.CustomerID != nil {
		_, err := s.customerRepo.FindByID(*input.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("customer not found")
		}
		invoice.CustomerID = *input.CustomerID
	}

	if input.SalespersonID != nil {
		_, err := s.salespersonRepo.FindByID(*input.SalespersonID)
		if err != nil {
			return nil, fmt.Errorf("salesperson not found")
		}
		invoice.SalespersonID = input.SalespersonID
	}

	if input.SalesOrderID != nil {
		invoice.SalesOrderID = input.SalesOrderID
	}

	if input.OrderNumber != nil {
		invoice.OrderNumber = *input.OrderNumber
	}

	if input.InvoiceDate != nil {
		invoice.InvoiceDate = *input.InvoiceDate
	}

	if input.Terms != nil {
		invoice.Terms = domain.PaymentTerms(*input.Terms)
	}

	if input.DueDate != nil {
		invoice.DueDate = *input.DueDate
	}

	if input.Subject != nil {
		invoice.Subject = *input.Subject
	}

	if input.CustomerNotes != nil {
		invoice.CustomerNotes = *input.CustomerNotes
	}

	if input.TermsAndConditions != nil {
		invoice.TermsAndConditions = *input.TermsAndConditions
	}

	if input.Attachments != nil {
		invoice.Attachments = input.Attachments
	}

	if len(input.LineItems) > 0 {
		lineItems := make([]models.InvoiceLineItem, len(input.LineItems))
		var subTotal float64

		for i, itemInput := range input.LineItems {
			// Fetch product
			_, err := s.productRepo.FindByID(*itemInput.ProductID)
			if err != nil {
				return nil, fmt.Errorf("product %s not found", *itemInput.ProductID)
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			variantDetails := models.VariantDetails{
				"sku":     itemInput.SKU,
				"account": itemInput.Account,
			}

			lineItems[i] = models.InvoiceLineItem{
				ProductID:      itemInput.ProductID,
				ProductName:    itemInput.ProductName,
				Description:    itemInput.ProductName,
				Quantity:       itemInput.Quantity,
				Rate:           itemInput.Rate,
				Amount:         amount,
				VariantDetails: variantDetails,
			}
		}

		invoice.LineItems = lineItems
		invoice.SubTotal = subTotal
	}

	if input.ShippingCharges != nil {
		invoice.ShippingCharges = *input.ShippingCharges
	}

	if input.TaxID != nil {
		tax, err := s.taxRepo.FindByID(*input.TaxID)
		if err != nil {
			return nil, fmt.Errorf("tax not found")
		}
		invoice.TaxID = input.TaxID
		invoice.TaxAmount = (invoice.SubTotal + invoice.ShippingCharges) * tax.Rate / 100
	}

	if input.TaxType != nil {
		invoice.TaxType = domain.TaxType(*input.TaxType)
	}

	if input.Adjustment != nil {
		invoice.Adjustment = *input.Adjustment
	}

	invoice.Total = invoice.SubTotal + invoice.ShippingCharges + invoice.TaxAmount + invoice.Adjustment
	invoice.UpdatedAt = time.Now()
	invoice.UpdatedBy = userID

	if err := s.invoiceRepo.Update(invoice); err != nil {
		return nil, err
	}

	updatedInvoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceOutput(updatedInvoice)
}

func (s *invoiceService) DeleteInvoice(id string) error {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return errors.New("invoice not found")
	}

	if invoice.Status == domain.InvoiceStatusIssued || invoice.Status == domain.InvoiceStatusSent || invoice.Status == domain.InvoiceStatusPaid {
		return fmt.Errorf("cannot delete invoice with status %s", invoice.Status)
	}

	return s.invoiceRepo.Delete(id)
}

func (s *invoiceService) GetInvoicesByCustomer(customerID string, limit, offset int) (*output.InvoiceListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	invoices, total, err := s.invoiceRepo.FindByCustomerID(customerID, limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceListOutput(invoices, total)
}

func (s *invoiceService) GetInvoicesByStatus(status string, limit, offset int) (*output.InvoiceListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	invoices, total, err := s.invoiceRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceListOutput(invoices, total)
}

func (s *invoiceService) UpdateInvoiceStatus(id string, status domain.InvoiceStatus) (*output.InvoiceOutput, error) {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Only handle stock deduction when transitioning to "issued" status
	if status == domain.InvoiceStatusIssued && invoice.Status != domain.InvoiceStatusIssued {
		// Deduct stock for each line item with product_id
		for _, lineItem := range invoice.LineItems {
			if lineItem.ProductID == nil {
				continue // Skip items without product reference
			}

			// Get current product stock
			productStock, err := s.productStockRepo.GetByProductID(*lineItem.ProductID)
			if err != nil {
				return nil, fmt.Errorf("failed to get stock for product %s: %w", *lineItem.ProductID, err)
			}

			// Calculate new stock levels
			newCurrentStock := productStock.CurrentStock - lineItem.Quantity
			newAvailableStock := productStock.AvailableStock - lineItem.Quantity

			if newCurrentStock < 0 {
				return nil, fmt.Errorf("insufficient stock for product %s. Required: %f, Available: %f", *lineItem.ProductID, lineItem.Quantity, productStock.CurrentStock)
			}

			// Update product stock
			productStock.CurrentStock = newCurrentStock
			productStock.AvailableStock = newAvailableStock
			productStock.SoldStock = productStock.SoldStock + lineItem.Quantity
			productStock.LastSoldDate = &time.Time{}
			*productStock.LastSoldDate = time.Now()
			productStock.UpdatedAt = time.Now()

			if err := s.productStockRepo.Update(productStock); err != nil {
				return nil, fmt.Errorf("failed to update stock for product %s: %w", *lineItem.ProductID, err)
			}

			// Create stock ledger entry using invoice_number as reference
			ledgerEntry := &models.StockLedger{
				ProductID:        *lineItem.ProductID,
				MovementType:     "SALES_INVOICE",
				Quantity:         -lineItem.Quantity, // Negative for outbound
				Rate:             productStock.AverageCost,
				Amount:           -lineItem.Quantity * productStock.AverageCost,
				ReferenceType:    "INVOICE",
				ReferenceID:      invoice.ID,
				ReferenceNumber:  invoice.InvoiceNumber,
				BalanceBeforeQty: productStock.CurrentStock + lineItem.Quantity, // Before deduction
				BalanceAfterQty:  productStock.CurrentStock,                     // After deduction
				Notes:            fmt.Sprintf("Invoice %s issued", invoice.InvoiceNumber),
				CreatedAt:        time.Now(),
				CreatedBy:        invoice.UpdatedBy,
			}

			if err := s.stockLedgerRepo.Create(ledgerEntry); err != nil {
				return nil, fmt.Errorf("failed to create stock ledger entry: %w", err)
			}
		}
	}

	// Update invoice status
	invoice.Status = status
	invoice.UpdatedAt = time.Now()

	if err := s.invoiceRepo.Update(invoice); err != nil {
		return nil, err
	}

	updatedInvoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceOutput(updatedInvoice)
}

func parseInvoiceUserID(userID string) (uint, error) {
	var parsed uint
	if _, err := fmt.Sscanf(userID, "%d", &parsed); err != nil || parsed == 0 {
		return 0, errors.New("invalid authenticated user")
	}
	return parsed, nil
}

func (s *invoiceService) ValidateInvoiceInputForCompany(
	invoiceInput *input.CreateInvoiceInput,
	companyID uint,
) error {
	if invoiceInput == nil {
		return errors.New("input cannot be nil")
	}
	if companyID == 0 {
		return errors.New("invalid company")
	}

	if _, err := s.customerRepo.FindByIDAndCompany(
		invoiceInput.CustomerID,
		companyID,
	); err != nil {
		return errors.New("customer not found in your company")
	}

	for _, lineItem := range invoiceInput.LineItems {
		if lineItem.ProductID == nil || *lineItem.ProductID == "" {
			return errors.New("product_id is required for each line item")
		}
		if _, err := s.productRepo.FindByIDAndCompany(
			*lineItem.ProductID,
			companyID,
		); err != nil {
			_ = err
		}
	}

	return nil
}

func (s *invoiceService) validateInvoiceUpdateForCompany(
	invoiceInput *input.UpdateInvoiceInput,
	companyID uint,
) error {
	if invoiceInput == nil {
		return errors.New("input cannot be nil")
	}

	if invoiceInput.CustomerID != nil {
		if _, err := s.customerRepo.FindByIDAndCompany(
			*invoiceInput.CustomerID,
			companyID,
		); err != nil {
			return errors.New("customer not found in your company")
		}
	}

	for _, lineItem := range invoiceInput.LineItems {
		if lineItem.ProductID == nil || *lineItem.ProductID == "" {
			return errors.New("product_id is required for each line item")
		}
		if _, err := s.productRepo.FindByIDAndCompany(
			*lineItem.ProductID,
			companyID,
		); err != nil {
			_ = err
		}
	}

	return nil
}

func (s *invoiceService) CreateInvoiceForCompany(
	invoiceInput *input.CreateInvoiceInput,
	userID string,
	companyID uint,
) (*output.InvoiceOutput, error) {
	if err := s.ValidateInvoiceInputForCompany(
		invoiceInput,
		companyID,
	); err != nil {
		return nil, err
	}

	userIDUint, err := parseInvoiceUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByIDAndCompanyID(
		userIDUint,
		companyID,
	)
	if err != nil || user == nil {
		return nil, errors.New("user does not belong to the company")
	}

	createdInvoice, err := s.CreateInvoice(invoiceInput, userID)
	if err != nil {
		return nil, err
	}

	return s.GetInvoiceByCompany(
		createdInvoice.ID,
		companyID,
	)
}

func (s *invoiceService) GetInvoiceByCompany(
	id string,
	companyID uint,
) (*output.InvoiceOutput, error) {
	invoice, err := s.invoiceRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	return output.ToInvoiceOutput(invoice)
}

func (s *invoiceService) GetAllInvoicesByCompany(
	companyID uint,
	limit int,
	offset int,
) (*output.InvoiceListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	invoices, total, err := s.invoiceRepo.FindAllByCompany(
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceListOutput(invoices, total)
}

func (s *invoiceService) UpdateInvoiceForCompany(
	id string,
	invoiceInput *input.UpdateInvoiceInput,
	userID string,
	companyID uint,
) (*output.InvoiceOutput, error) {
	if _, err := s.invoiceRepo.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return nil, errors.New("invoice not found")
	}

	if err := s.validateInvoiceUpdateForCompany(
		invoiceInput,
		companyID,
	); err != nil {
		return nil, err
	}

	updatedInvoice, err := s.UpdateInvoice(
		id,
		invoiceInput,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return s.GetInvoiceByCompany(
		updatedInvoice.ID,
		companyID,
	)
}

func (s *invoiceService) DeleteInvoiceForCompany(
	id string,
	companyID uint,
) error {
	invoice, err := s.invoiceRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return errors.New("invoice not found")
	}

	if invoice.Status == domain.InvoiceStatusIssued ||
		invoice.Status == domain.InvoiceStatusSent ||
		invoice.Status == domain.InvoiceStatusPaid {
		return fmt.Errorf(
			"cannot delete invoice with status %s",
			invoice.Status,
		)
	}

	return s.invoiceRepo.DeleteByCompany(id, companyID)
}

func (s *invoiceService) GetInvoicesByCustomerAndCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) (*output.InvoiceListOutput, error) {
	if _, err := s.customerRepo.FindByIDAndCompany(
		customerID,
		companyID,
	); err != nil {
		return nil, errors.New("customer not found in your company")
	}

	invoices, total, err := s.invoiceRepo.FindByCustomerIDAndCompany(
		customerID,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceListOutput(invoices, total)
}

func (s *invoiceService) GetInvoicesByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) (*output.InvoiceListOutput, error) {
	invoices, total, err := s.invoiceRepo.FindByStatusAndCompany(
		status,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return output.ToInvoiceListOutput(invoices, total)
}

func (s *invoiceService) UpdateInvoiceStatusForCompany(
	id string,
	status domain.InvoiceStatus,
	userID string,
	companyID uint,
) (*output.InvoiceOutput, error) {
	invoice, err := s.invoiceRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	shouldFilter := companyID > 0

	if status == domain.InvoiceStatusIssued &&
		invoice.Status != domain.InvoiceStatusIssued {
		for _, lineItem := range invoice.LineItems {
			if lineItem.ProductID == nil {
				continue
			}

			if _, err := s.productRepo.FindByIDAndCompany(
				*lineItem.ProductID,
				companyID,
			); err != nil {
				return nil, fmt.Errorf(
					"product %s not found in your company",
					*lineItem.ProductID,
				)
			}

			productStock, err :=
				s.productStockRepo.GetByProductIDAndCompany(
					*lineItem.ProductID,
					companyID,
					shouldFilter,
				)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to get stock for product %s: %w",
					*lineItem.ProductID,
					err,
				)
			}

			if productStock.CurrentStock < lineItem.Quantity {
				return nil, fmt.Errorf(
					"insufficient stock for product %s. Required: %f, Available: %f",
					*lineItem.ProductID,
					lineItem.Quantity,
					productStock.CurrentStock,
				)
			}

			beforeQuantity := productStock.CurrentStock
			productStock.CurrentStock -= lineItem.Quantity
			productStock.AvailableStock -= lineItem.Quantity
			productStock.SoldStock += lineItem.Quantity
			now := time.Now()
			productStock.LastSoldDate = &now
			productStock.UpdatedAt = now

			if err := s.productStockRepo.UpdateByCompany(
				productStock,
				companyID,
				shouldFilter,
			); err != nil {
				return nil, fmt.Errorf(
					"failed to update stock for product %s: %w",
					*lineItem.ProductID,
					err,
				)
			}

			ledgerEntry := &models.StockLedger{
				ProductID:        *lineItem.ProductID,
				MovementType:     "SALES_INVOICE",
				Quantity:         -lineItem.Quantity,
				Rate:             productStock.AverageCost,
				Amount:           -lineItem.Quantity * productStock.AverageCost,
				ReferenceType:    "INVOICE",
				ReferenceID:      invoice.ID,
				ReferenceNumber:  invoice.InvoiceNumber,
				BalanceBeforeQty: beforeQuantity,
				BalanceAfterQty:  productStock.CurrentStock,
				Notes:            fmt.Sprintf("Invoice %s issued", invoice.InvoiceNumber),
				CreatedAt:        now,
				CreatedBy:        userID,
			}

			if err := s.stockLedgerRepo.CreateForCompany(
				ledgerEntry,
				companyID,
				shouldFilter,
			); err != nil {
				return nil, fmt.Errorf(
					"failed to create stock ledger entry: %w",
					err,
				)
			}
		}
	}

	invoice.Status = status
	invoice.UpdatedAt = time.Now()
	invoice.UpdatedBy = userID

	if err := s.invoiceRepo.UpdateByCompany(
		invoice,
		companyID,
	); err != nil {
		return nil, err
	}

	return s.GetInvoiceByCompany(id, companyID)
}

type salespersonService struct {
	repo repo.SalespersonRepository
}

func NewSalespersonService(repo repo.SalespersonRepository) SalespersonService {
	return &salespersonService{repo: repo}
}

func (s *salespersonService) CreateSalesperson(input *input.CreateSalespersonInput) (*output.SalespersonOutput, error) {
	existing, _ := s.repo.FindByEmail(input.Email)
	if existing != nil {
		return nil, fmt.Errorf("salesperson with email %s already exists", input.Email)
	}

	salesperson := &models.Salesperson{
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(salesperson); err != nil {
		return nil, err
	}

	return output.ToSalespersonOutput(salesperson), nil
}

func (s *salespersonService) GetSalesperson(id uint) (*output.SalespersonOutput, error) {
	salesperson, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToSalespersonOutput(salesperson), nil
}

func (s *salespersonService) GetAllSalespersons(limit, offset int) (*output.SalespersonListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	salespersons, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToSalespersonListOutput(salespersons, total), nil
}

func (s *salespersonService) UpdateSalesperson(id uint, input *input.UpdateSalespersonInput) (*output.SalespersonOutput, error) {
	salesperson, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		salesperson.Name = *input.Name
	}

	if input.Email != nil {
		existing, _ := s.repo.FindByEmail(*input.Email)
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("email %s is already taken", *input.Email)
		}
		salesperson.Email = *input.Email
	}

	salesperson.UpdatedAt = time.Now()

	if err := s.repo.Update(salesperson); err != nil {
		return nil, err
	}

	return output.ToSalespersonOutput(salesperson), nil
}

func (s *salespersonService) DeleteSalesperson(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("salesperson not found")
	}

	return s.repo.Delete(id)
}

type taxService struct {
	repo repo.TaxRepository
}

func NewTaxService(repo repo.TaxRepository) TaxService {
	return &taxService{repo: repo}
}

func (s *taxService) CreateTax(input *input.CreateTaxInput) (*output.TaxOutput, error) {
	tax := &models.Tax{
		Name:      input.Name,
		TaxType:   input.TaxType,
		Rate:      input.Rate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(tax); err != nil {
		return nil, err
	}

	return output.ToTaxOutput(tax), nil
}

func (s *taxService) GetTax(id uint) (*output.TaxOutput, error) {
	tax, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToTaxOutput(tax), nil
}

func (s *taxService) GetAllTaxes(limit, offset int) (*output.TaxListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	taxes, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToTaxListOutput(taxes, total), nil
}

func (s *taxService) UpdateTax(id uint, input *input.UpdateTaxInput) (*output.TaxOutput, error) {
	tax, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		tax.Name = *input.Name
	}

	if input.TaxType != nil {
		tax.TaxType = *input.TaxType
	}

	if input.Rate != nil {
		tax.Rate = *input.Rate
	}

	tax.UpdatedAt = time.Now()

	if err := s.repo.Update(tax); err != nil {
		return nil, err
	}

	return output.ToTaxOutput(tax), nil
}

func (s *taxService) DeleteTax(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("tax not found")
	}

	return s.repo.Delete(id)
}

type paymentService struct {
	paymentRepo repo.PaymentRepository
	invoiceRepo repo.InvoiceRepository
}

func NewPaymentService(paymentRepo repo.PaymentRepository, invoiceRepo repo.InvoiceRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
	}
}

func (s *paymentService) CreatePayment(input *input.CreatePaymentInput, userID string) (*output.PaymentOutput, error) {
	invoice, err := s.invoiceRepo.FindByID(input.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found")
	}

	existingPayments, err := s.paymentRepo.FindByInvoiceID(input.InvoiceID)
	if err != nil {
		return nil, err
	}

	var totalPaid float64
	for _, p := range existingPayments {
		totalPaid += p.Amount
	}

	if totalPaid+input.Amount > invoice.Total {
		return nil, fmt.Errorf("payment amount exceeds remaining invoice balance")
	}

	payment := &models.Payment{
		InvoiceID:   input.InvoiceID,
		PaymentDate: input.PaymentDate,
		Amount:      input.Amount,
		PaymentMode: input.PaymentMode,
		Reference:   input.Reference,
		Notes:       input.Notes,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	totalPaid += input.Amount
	if totalPaid >= invoice.Total {
		invoice.Status = domain.InvoiceStatusPaid
	} else {
		invoice.Status = domain.InvoiceStatusPartial
	}
	invoice.UpdatedAt = time.Now()

	if err := s.invoiceRepo.Update(invoice); err != nil {
		return nil, err
	}

	return output.ToPaymentOutput(payment), nil
}

func (s *paymentService) GetPayment(id uint) (*output.PaymentOutput, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToPaymentOutput(payment), nil
}

func (s *paymentService) GetPaymentsByInvoice(invoiceID string) (*output.PaymentListOutput, error) {
	payments, err := s.paymentRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return nil, err
	}

	return output.ToPaymentListOutput(payments, int64(len(payments))), nil
}

func (s *paymentService) DeletePayment(id uint) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return errors.New("payment not found")
	}

	invoice, err := s.invoiceRepo.FindByID(payment.InvoiceID)
	if err != nil {
		return err
	}

	if err := s.paymentRepo.Delete(id); err != nil {
		return err
	}

	payments, err := s.paymentRepo.FindByInvoiceID(payment.InvoiceID)
	if err != nil {
		return err
	}

	var totalPaid float64
	for _, p := range payments {
		totalPaid += p.Amount
	}

	if totalPaid == 0 {
		invoice.Status = domain.InvoiceStatusSent
	} else if totalPaid < invoice.Total {
		invoice.Status = domain.InvoiceStatusPartial
	}
	invoice.UpdatedAt = time.Now()

	return s.invoiceRepo.Update(invoice)
}

func (s *paymentService) CreatePaymentForCompany(
	paymentInput *input.CreatePaymentInput,
	userID string,
	companyID uint,
) (*output.PaymentOutput, error) {
	invoice, err := s.invoiceRepo.FindByIDAndCompany(
		paymentInput.InvoiceID,
		companyID,
	)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	existingPayments, err := s.paymentRepo.FindByInvoiceID(
		paymentInput.InvoiceID,
	)
	if err != nil {
		return nil, err
	}

	var totalPaid float64
	for _, payment := range existingPayments {
		totalPaid += payment.Amount
	}

	if totalPaid+paymentInput.Amount > invoice.Total {
		return nil, errors.New(
			"payment amount exceeds remaining invoice balance",
		)
	}

	payment := &models.Payment{
		InvoiceID:   paymentInput.InvoiceID,
		PaymentDate: paymentInput.PaymentDate,
		Amount:      paymentInput.Amount,
		PaymentMode: paymentInput.PaymentMode,
		Reference:   paymentInput.Reference,
		Notes:       paymentInput.Notes,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	totalPaid += paymentInput.Amount
	if totalPaid >= invoice.Total {
		invoice.Status = domain.InvoiceStatusPaid
	} else {
		invoice.Status = domain.InvoiceStatusPartial
	}
	invoice.UpdatedAt = time.Now()
	invoice.UpdatedBy = userID

	if err := s.invoiceRepo.UpdateByCompany(
		invoice,
		companyID,
	); err != nil {
		return nil, err
	}

	return output.ToPaymentOutput(payment), nil
}

func (s *paymentService) GetPaymentForCompany(
	id uint,
	companyID uint,
) (*output.PaymentOutput, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if _, err := s.invoiceRepo.FindByIDAndCompany(
		payment.InvoiceID,
		companyID,
	); err != nil {
		return nil, errors.New("payment not found")
	}

	return output.ToPaymentOutput(payment), nil
}

func (s *paymentService) GetPaymentsByInvoiceForCompany(
	invoiceID string,
	companyID uint,
) (*output.PaymentListOutput, error) {
	if _, err := s.invoiceRepo.FindByIDAndCompany(
		invoiceID,
		companyID,
	); err != nil {
		return nil, errors.New("invoice not found")
	}

	payments, err := s.paymentRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return nil, err
	}

	return output.ToPaymentListOutput(
		payments,
		int64(len(payments)),
	), nil
}

func (s *paymentService) DeletePaymentForCompany(
	id uint,
	companyID uint,
) error {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return errors.New("payment not found")
	}

	invoice, err := s.invoiceRepo.FindByIDAndCompany(
		payment.InvoiceID,
		companyID,
	)
	if err != nil {
		return errors.New("payment not found")
	}

	if err := s.paymentRepo.Delete(id); err != nil {
		return err
	}

	payments, err := s.paymentRepo.FindByInvoiceID(
		payment.InvoiceID,
	)
	if err != nil {
		return err
	}

	var totalPaid float64
	for _, existingPayment := range payments {
		totalPaid += existingPayment.Amount
	}

	if totalPaid == 0 {
		invoice.Status = domain.InvoiceStatusSent
	} else if totalPaid < invoice.Total {
		invoice.Status = domain.InvoiceStatusPartial
	}

	invoice.UpdatedAt = time.Now()

	return s.invoiceRepo.UpdateByCompany(
		invoice,
		companyID,
	)
}
