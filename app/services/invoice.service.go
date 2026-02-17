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
	CreateInvoice(input *input.CreateInvoiceInput, userID string) (*output.InvoiceOutput, error)
	GetInvoice(id string) (*output.InvoiceOutput, error)
	GetAllInvoices(limit, offset int) (*output.InvoiceListOutput, error)
	UpdateInvoice(id string, input *input.UpdateInvoiceInput, userID string) (*output.InvoiceOutput, error)
	DeleteInvoice(id string) error
	GetInvoicesByCustomer(customerID string, limit, offset int) (*output.InvoiceListOutput, error)
	GetInvoicesByStatus(status string, limit, offset int) (*output.InvoiceListOutput, error)
	UpdateInvoiceStatus(id string, status domain.InvoiceStatus) (*output.InvoiceOutput, error)
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
	CreatePayment(input *input.CreatePaymentInput, userID string) (*output.PaymentOutput, error)
	GetPayment(id uint) (*output.PaymentOutput, error)
	GetPaymentsByInvoice(invoiceID string) (*output.PaymentListOutput, error)
	DeletePayment(id uint) error
}

type invoiceService struct {
	invoiceRepo     repo.InvoiceRepository
	itemRepo        repo.ItemRepository
	customerRepo    repo.CustomerRepository
	salespersonRepo repo.SalespersonRepository
	taxRepo         repo.TaxRepository
	paymentRepo     repo.PaymentRepository
}

func NewInvoiceService(
	invoiceRepo repo.InvoiceRepository,
	itemRepo repo.ItemRepository,
	customerRepo repo.CustomerRepository,
	salespersonRepo repo.SalespersonRepository,
	taxRepo repo.TaxRepository,
	paymentRepo repo.PaymentRepository,
	pdfOutputDir string,
) InvoiceService {
	return &invoiceService{
		invoiceRepo:     invoiceRepo,
		itemRepo:        itemRepo,
		customerRepo:    customerRepo,
		salespersonRepo: salespersonRepo,
		taxRepo:         taxRepo,
		paymentRepo:     paymentRepo,
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
		item, err := s.itemRepo.FindByID(itemInput.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item %s not found", itemInput.ItemID)
		}

		if itemInput.VariantID != nil {
			variantFound := false
			for _, variant := range item.ItemDetails.Variants {
				if variant.ID == *itemInput.VariantID {
					variantFound = true
					break
				}
			}
			if !variantFound {
				return nil, fmt.Errorf("variant %d not found for item %s", *itemInput.VariantID, itemInput.ItemID)
			}
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		lineItems[i] = models.InvoiceLineItem{
			ItemID:         itemInput.ItemID,
			VariantID:      itemInput.VariantID,
			Description:    itemInput.Description,
			Quantity:       itemInput.Quantity,
			Rate:           itemInput.Rate,
			Amount:         amount,
			VariantDetails: itemInput.VariantDetails,
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
			item, err := s.itemRepo.FindByID(itemInput.ItemID)
			if err != nil {
				return nil, fmt.Errorf("item %s not found", itemInput.ItemID)
			}

			if itemInput.VariantID != nil {
				variantFound := false
				for _, variant := range item.ItemDetails.Variants {
					if variant.ID == *itemInput.VariantID {
						variantFound = true
						break
					}
				}
				if !variantFound {
					return nil, fmt.Errorf("variant %d not found for item %s", *itemInput.VariantID, itemInput.ItemID)
				}
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			lineItems[i] = models.InvoiceLineItem{
				ItemID:         itemInput.ItemID,
				VariantID:      itemInput.VariantID,
				Description:    itemInput.Description,
				Quantity:       itemInput.Quantity,
				Rate:           itemInput.Rate,
				Amount:         amount,
				VariantDetails: itemInput.VariantDetails,
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

	if invoice.Status == domain.InvoiceStatusSent || invoice.Status == domain.InvoiceStatusPaid {
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
