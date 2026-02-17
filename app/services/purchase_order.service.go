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

type PurchaseOrderService interface {
	CreatePurchaseOrder(poInput *input.CreatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error)
	GetPurchaseOrder(id string) (*output.PurchaseOrderOutput, error)
	GetAllPurchaseOrders(limit, offset int) (*output.PurchaseOrderListOutput, error)
	UpdatePurchaseOrder(id string, poInput *input.UpdatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error)
	DeletePurchaseOrder(id string) error
	GetPurchaseOrdersByVendor(vendorID uint, limit, offset int) (*output.PurchaseOrderListOutput, error)
	GetPurchaseOrdersByCustomer(customerID uint, limit, offset int) (*output.PurchaseOrderListOutput, error)
	GetPurchaseOrdersByStatus(status string, limit, offset int) (*output.PurchaseOrderListOutput, error)
	UpdatePurchaseOrderStatus(id string, status domain.PurchaseOrderStatus, userID string) (*output.PurchaseOrderOutput, error)
}

type purchaseOrderService struct {
	poRepo       repo.PurchaseOrderRepository
	vendorRepo   repo.VendorRepository
	customerRepo repo.CustomerRepository
	itemRepo     repo.ItemRepository
	taxRepo      repo.TaxRepository
}

func NewPurchaseOrderService(
	poRepo repo.PurchaseOrderRepository,
	vendorRepo repo.VendorRepository,
	customerRepo repo.CustomerRepository,
	itemRepo repo.ItemRepository,
	taxRepo repo.TaxRepository,
) PurchaseOrderService {
	return &purchaseOrderService{
		poRepo:       poRepo,
		vendorRepo:   vendorRepo,
		customerRepo: customerRepo,
		itemRepo:     itemRepo,
		taxRepo:      taxRepo,
	}
}

func (s *purchaseOrderService) CreatePurchaseOrder(poInput *input.CreatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error) {
	vendor, err := s.vendorRepo.FindByID(poInput.VendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	if poInput.DeliveryAddressType != "organization" && poInput.DeliveryAddressType != "customer" {
		return nil, errors.New("invalid delivery address type")
	}

	var customer *models.Customer
	if poInput.DeliveryAddressType == "customer" {
		if poInput.CustomerID == nil {
			return nil, errors.New("customer_id is required for customer delivery")
		}
		var err error
		customer, err = s.customerRepo.FindByID(*poInput.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}

	var tax *models.Tax
	if poInput.TaxID != nil {
		tax, err = s.taxRepo.FindByID(*poInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
	}

	lineItems := make([]models.PurchaseOrderLineItem, 0)
	subTotal := 0.0

	for _, itemInput := range poInput.LineItems {
		item, err := s.itemRepo.FindByID(itemInput.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item %s not found", itemInput.ItemID)
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		lineItem := models.PurchaseOrderLineItem{
			ItemID:    itemInput.ItemID,
			Item:      item,
			VariantID: itemInput.VariantID,
			Account:   itemInput.Account,
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

	discount := poInput.Discount
	if poInput.DiscountType == "percentage" {
		discount = (subTotal * poInput.Discount) / 100
	}

	taxAmount := 0.0
	if tax != nil {
		taxAmount = ((subTotal - discount) * tax.Rate) / 100
	}

	total := subTotal - discount + taxAmount + poInput.Adjustment

	poNumber := fmt.Sprintf("PO-%s-%04d", time.Now().Format("20060102"), s.generatePOSequence())

	var taxType *domain.TaxType
	if poInput.TaxType != nil {
		tt := domain.TaxType(*poInput.TaxType)
		taxType = &tt
	}

	po := &models.PurchaseOrder{
		ID:                  uuid.New().String(),
		PurchaseOrderNumber: poNumber,
		VendorID:            poInput.VendorID,
		Vendor:              vendor,
		DeliveryAddressType: poInput.DeliveryAddressType,
		DeliveryAddressID:   poInput.DeliveryAddressID,
		OrganizationName:    poInput.OrganizationName,
		OrganizationAddress: poInput.OrganizationAddress,
		CustomerID:          poInput.CustomerID,
		Customer:            customer,
		ReferenceNo:         poInput.ReferenceNo,
		PODate:              poInput.Date,
		DeliveryDate:        poInput.DeliveryDate,
		PaymentTerms:        domain.PaymentTerms(poInput.PaymentTerms),
		ShipmentPreference:  poInput.ShipmentPreference,
		LineItems:           lineItems,
		SubTotal:            subTotal,
		Discount:            discount,
		DiscountType:        poInput.DiscountType,
		TaxType:             taxType,
		TaxID:               poInput.TaxID,
		Tax:                 tax,
		TaxAmount:           taxAmount,
		Adjustment:          poInput.Adjustment,
		Total:               total,
		Notes:               poInput.Notes,
		TermsAndConditions:  poInput.TermsAndConditions,
		Status:              domain.PurchaseOrderStatusDraft,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CreatedBy:           userID,
		UpdatedBy:           userID,
	}

	if len(poInput.Attachments) > 0 {
		po.Attachments = poInput.Attachments
	}

	createdPO, err := s.poRepo.Create(po)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return output.ToPurchaseOrderOutput(createdPO)
}

func (s *purchaseOrderService) GetPurchaseOrder(id string) (*output.PurchaseOrderOutput, error) {
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	return output.ToPurchaseOrderOutput(po)
}

func (s *purchaseOrderService) GetAllPurchaseOrders(limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) UpdatePurchaseOrder(id string, poInput *input.UpdatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error) {
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	if poInput.VendorID != nil {
		vendor, err := s.vendorRepo.FindByID(*poInput.VendorID)
		if err != nil {
			return nil, errors.New("vendor not found")
		}
		po.VendorID = *poInput.VendorID
		po.Vendor = vendor
	}

	if poInput.DeliveryAddressType != nil {
		po.DeliveryAddressType = *poInput.DeliveryAddressType
	}

	if poInput.OrganizationName != nil {
		po.OrganizationName = *poInput.OrganizationName
	}

	if poInput.OrganizationAddress != nil {
		po.OrganizationAddress = *poInput.OrganizationAddress
	}

	if poInput.CustomerID != nil {
		customer, err := s.customerRepo.FindByID(*poInput.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
		po.CustomerID = poInput.CustomerID
		po.Customer = customer
	}

	if poInput.ReferenceNo != nil {
		po.ReferenceNo = *poInput.ReferenceNo
	}

	if poInput.Date != nil {
		po.PODate = *poInput.Date
	}

	if poInput.DeliveryDate != nil {
		po.DeliveryDate = *poInput.DeliveryDate
	}

	if poInput.PaymentTerms != nil {
		po.PaymentTerms = domain.PaymentTerms(*poInput.PaymentTerms)
	}

	if poInput.ShipmentPreference != nil {
		po.ShipmentPreference = *poInput.ShipmentPreference
	}

	if len(poInput.LineItems) > 0 {
		lineItems := make([]models.PurchaseOrderLineItem, 0)
		subTotal := 0.0

		for _, itemInput := range poInput.LineItems {
			item, err := s.itemRepo.FindByID(itemInput.ItemID)
			if err != nil {
				return nil, fmt.Errorf("item %s not found", itemInput.ItemID)
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			lineItem := models.PurchaseOrderLineItem{
				ItemID:    itemInput.ItemID,
				Item:      item,
				VariantID: itemInput.VariantID,
				Account:   itemInput.Account,
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

		po.LineItems = lineItems
		po.SubTotal = subTotal
	}

	if poInput.Discount != nil {
		po.Discount = *poInput.Discount
	}

	if poInput.DiscountType != nil {
		po.DiscountType = *poInput.DiscountType
	}

	if poInput.TaxID != nil {
		tax, err := s.taxRepo.FindByID(*poInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
		po.TaxID = poInput.TaxID
		po.Tax = tax
	}

	if poInput.TaxType != nil {
		taxType := domain.TaxType(*poInput.TaxType)
		po.TaxType = &taxType
	}

	if poInput.Adjustment != nil {
		po.Adjustment = *poInput.Adjustment
	}

	discount := po.Discount
	if po.DiscountType == "percentage" {
		discount = (po.SubTotal * po.Discount) / 100
	}

	taxAmount := 0.0
	if po.Tax != nil {
		taxAmount = ((po.SubTotal - discount) * po.Tax.Rate) / 100
	}

	po.TaxAmount = taxAmount
	po.Total = po.SubTotal - discount + taxAmount + po.Adjustment

	if poInput.Notes != nil {
		po.Notes = *poInput.Notes
	}

	if poInput.TermsAndConditions != nil {
		po.TermsAndConditions = *poInput.TermsAndConditions
	}

	if len(poInput.Attachments) > 0 {
		po.Attachments = poInput.Attachments
	}

	updatedPO, err := s.poRepo.Update(id, po)
	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return output.ToPurchaseOrderOutput(updatedPO)
}

func (s *purchaseOrderService) DeletePurchaseOrder(id string) error {
	return s.poRepo.Delete(id)
}

func (s *purchaseOrderService) GetPurchaseOrdersByVendor(vendorID uint, limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindByVendor(vendorID, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) GetPurchaseOrdersByCustomer(customerID uint, limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindByCustomer(customerID, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) GetPurchaseOrdersByStatus(status string, limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) UpdatePurchaseOrderStatus(id string, status domain.PurchaseOrderStatus, userID string) (*output.PurchaseOrderOutput, error) {
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	po.Status = status
	po.UpdatedAt = time.Now()
	po.UpdatedBy = userID

	err = s.poRepo.UpdateStatus(id, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order status: %w", err)
	}

	return s.GetPurchaseOrder(id)
}

func (s *purchaseOrderService) generatePOSequence() int {
	var count int64
	today := time.Now().Format("2006-01-02")

	s.poRepo.GetDB().Where("DATE(created_at) = ?", today).Model(&models.PurchaseOrder{}).Count(&count)

	return int(count) + 1
}
