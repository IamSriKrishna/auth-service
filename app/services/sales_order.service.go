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

type SalesOrderService interface {
	CreateSalesOrder(soInput *input.CreateSalesOrderInput, userID string) (*output.SalesOrderOutput, error)
	GetSalesOrder(id string) (*output.SalesOrderOutput, error)
	GetAllSalesOrders(limit, offset int) ([]output.SalesOrderOutput, int64, error)
	GetSalesOrdersByCustomer(customerID uint, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	GetSalesOrdersByStatus(status string, limit, offset int) ([]output.SalesOrderOutput, int64, error)
	UpdateSalesOrder(id string, soInput *input.UpdateSalesOrderInput, userID string) (*output.SalesOrderOutput, error)
	UpdateSalesOrderStatus(id string, status string, userID string) (*output.SalesOrderOutput, error)
	DeleteSalesOrder(id string) error
}

type salesOrderService struct {
	soRepo          repo.SalesOrderRepository
	customerRepo    repo.CustomerRepository
	itemRepo        repo.ItemRepository
	taxRepo         repo.TaxRepository
	salespersonRepo repo.SalespersonRepository
	inventoryRepo   repo.InventoryBalanceRepository
}

func NewSalesOrderService(
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	itemRepo repo.ItemRepository,
	taxRepo repo.TaxRepository,
	salespersonRepo repo.SalespersonRepository,
	inventoryRepo repo.InventoryBalanceRepository,
) SalesOrderService {
	return &salesOrderService{
		soRepo:          soRepo,
		customerRepo:    customerRepo,
		itemRepo:        itemRepo,
		taxRepo:         taxRepo,
		salespersonRepo: salespersonRepo,
		inventoryRepo:   inventoryRepo,
	}
}

func (s *salesOrderService) CreateSalesOrder(soInput *input.CreateSalesOrderInput, userID string) (*output.SalesOrderOutput, error) {
	customer, err := s.customerRepo.FindByID(soInput.CustomerID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	var salesperson *models.Salesperson
	if soInput.SalespersonID != nil {
		salesperson, err = s.salespersonRepo.FindByID(*soInput.SalespersonID)
		if err != nil {
			return nil, errors.New("salesperson not found")
		}
	}

	var tax *models.Tax
	if soInput.TaxID != nil {
		tax, err = s.taxRepo.FindByID(*soInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
	}

	lineItems := make([]models.SalesOrderLineItem, 0)
	subTotal := 0.0

	for _, itemInput := range soInput.LineItems {
		item, err := s.itemRepo.FindByID(itemInput.ItemID)
		if err != nil {
			return nil, errors.New("item not found: " + itemInput.ItemID)
		}

		// Check inventory availability
		inventoryBalance, err := s.inventoryRepo.GetBalance(itemInput.ItemID, itemInput.VariantID)
		if err != nil {
			return nil, fmt.Errorf("failed to check inventory for item %s: %v", itemInput.ItemID, err)
		}

		if inventoryBalance.AvailableQuantity < itemInput.Quantity {
			variantName := "N/A"
			if itemInput.VariantID != nil && item.ItemDetails.Variants != nil {
				for _, v := range item.ItemDetails.Variants {
					if v.ID == *itemInput.VariantID {
						variantName = v.SKU
						break
					}
				}
			}
			return nil, fmt.Errorf("insufficient inventory for %s (%s). Required: %f units, Available: %f units",
				item.Name, variantName, itemInput.Quantity, inventoryBalance.AvailableQuantity)
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		lineItem := models.SalesOrderLineItem{
			ItemID:    itemInput.ItemID,
			Item:      item,
			VariantID: itemInput.VariantID,
			Quantity:  itemInput.Quantity,
			Rate:      itemInput.Rate,
			Amount:    amount,
		}

		if itemInput.VariantDetails != nil {
			variantDetails := make(models.VariantDetails)
			for k, v := range itemInput.VariantDetails {
				variantDetails[k] = v
			}
			lineItem.VariantDetails = variantDetails
		}

		lineItems = append(lineItems, lineItem)
	}

	taxAmount := 0.0
	if tax != nil {
		taxAmount = ((subTotal + soInput.ShippingCharges) * tax.Rate) / 100
	}

	total := subTotal + soInput.ShippingCharges + taxAmount + soInput.Adjustment

	soNumber := fmt.Sprintf("SO-%s-%04d", time.Now().Format("20060102"), s.generateSOSequence())

	var taxType *domain.TaxType
	if soInput.TaxType != nil {
		tt := domain.TaxType(*soInput.TaxType)
		taxType = &tt
	}

	so := &models.SalesOrder{
		ID:                   uuid.New().String(),
		SalesOrderNumber:     soNumber,
		CustomerID:           soInput.CustomerID,
		Customer:             customer,
		SalespersonID:        soInput.SalespersonID,
		Salesperson:          salesperson,
		ReferenceNo:          soInput.ReferenceNo,
		SODate:               soInput.SODate,
		ExpectedShipmentDate: soInput.ExpectedShipmentDate,
		PaymentTerms:         domain.PaymentTerms(soInput.PaymentTerms),
		DeliveryMethod:       soInput.DeliveryMethod,
		LineItems:            lineItems,
		SubTotal:             subTotal,
		ShippingCharges:      soInput.ShippingCharges,
		TaxType:              taxType,
		TaxID:                soInput.TaxID,
		Tax:                  tax,
		TaxAmount:            taxAmount,
		Adjustment:           soInput.Adjustment,
		Total:                total,
		CustomerNotes:        soInput.CustomerNotes,
		TermsAndConditions:   soInput.TermsAndConditions,
		Status:               domain.SalesOrderStatusDraft,
		Attachments:          soInput.Attachments,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CreatedBy:            userID,
		UpdatedBy:            userID,
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
		customer, err := s.customerRepo.FindByID(*soInput.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
		so.CustomerID = *soInput.CustomerID
		so.Customer = customer
	}

	if soInput.SalespersonID != nil {
		salesperson, err := s.salespersonRepo.FindByID(*soInput.SalespersonID)
		if err != nil {
			return nil, errors.New("salesperson not found")
		}
		so.SalespersonID = soInput.SalespersonID
		so.Salesperson = salesperson
	}

	if soInput.ReferenceNo != nil {
		so.ReferenceNo = *soInput.ReferenceNo
	}

	if soInput.SODate != nil {
		so.SODate = *soInput.SODate
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

	if len(soInput.LineItems) > 0 {
		lineItems := make([]models.SalesOrderLineItem, 0)
		subTotal := 0.0

		for _, itemInput := range soInput.LineItems {
			item, err := s.itemRepo.FindByID(itemInput.ItemID)
			if err != nil {
				return nil, errors.New("item not found")
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			lineItem := models.SalesOrderLineItem{
				ItemID:    itemInput.ItemID,
				Item:      item,
				VariantID: itemInput.VariantID,
				Quantity:  itemInput.Quantity,
				Rate:      itemInput.Rate,
				Amount:    amount,
			}

			if itemInput.VariantDetails != nil {
				variantDetails := make(models.VariantDetails)
				for k, v := range itemInput.VariantDetails {
					variantDetails[k] = v
				}
				lineItem.VariantDetails = variantDetails
			}

			lineItems = append(lineItems, lineItem)
		}

		so.LineItems = lineItems
		so.SubTotal = subTotal

		if so.Tax != nil {
			so.TaxAmount = ((so.SubTotal + so.ShippingCharges) * so.Tax.Rate) / 100
		}

		so.Total = so.SubTotal + so.ShippingCharges + so.TaxAmount + so.Adjustment
	}

	if soInput.ShippingCharges != nil {
		so.ShippingCharges = *soInput.ShippingCharges

		if so.Tax != nil {
			so.TaxAmount = ((so.SubTotal + so.ShippingCharges) * so.Tax.Rate) / 100
		}

		so.Total = so.SubTotal + so.ShippingCharges + so.TaxAmount + so.Adjustment
	}

	if soInput.TaxID != nil {
		tax, err := s.taxRepo.FindByID(*soInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
		so.TaxID = soInput.TaxID
		so.Tax = tax

		so.TaxAmount = ((so.SubTotal + so.ShippingCharges) * tax.Rate) / 100
		so.Total = so.SubTotal + so.ShippingCharges + so.TaxAmount + so.Adjustment
	}

	if soInput.Adjustment != nil {
		so.Adjustment = *soInput.Adjustment
		so.Total = so.SubTotal + so.ShippingCharges + so.TaxAmount + *soInput.Adjustment
	}

	if soInput.CustomerNotes != nil {
		so.CustomerNotes = *soInput.CustomerNotes
	}

	if soInput.TermsAndConditions != nil {
		so.TermsAndConditions = *soInput.TermsAndConditions
	}

	if len(soInput.Attachments) > 0 {
		so.Attachments = soInput.Attachments
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
	err := s.soRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, errors.New("failed to update status: " + err.Error())
	}

	return s.GetSalesOrder(id)
}

func (s *salesOrderService) DeleteSalesOrder(id string) error {
	return s.soRepo.Delete(id)
}

func (s *salesOrderService) generateSOSequence() int {
	var count int64
	today := time.Now().Format("2006-01-02")

	s.soRepo.GetDB().Where("DATE(created_at) = ?", today).Model(&models.SalesOrder{}).Count(&count)

	return int(count) + 1
}
