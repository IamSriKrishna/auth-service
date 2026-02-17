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

type BillService interface {
	CreateBill(billInput *input.CreateBillInput, userID string) (*output.BillOutput, error)
	GetBill(id string) (*output.BillOutput, error)
	GetAllBills(limit, offset int) ([]output.BillOutput, int64, error)
	GetBillsByVendor(vendorID uint, limit, offset int) ([]output.BillOutput, int64, error)
	GetBillsByStatus(status string, limit, offset int) ([]output.BillOutput, int64, error)
	UpdateBill(id string, billInput *input.UpdateBillInput, userID string) (*output.BillOutput, error)
	UpdateBillStatus(id string, status string, userID string) (*output.BillOutput, error)
	DeleteBill(id string) error
}

type billService struct {
	billRepo   repo.BillRepository
	vendorRepo repo.VendorRepository
	itemRepo   repo.ItemRepository
	taxRepo    repo.TaxRepository
}

func NewBillService(
	billRepo repo.BillRepository,
	vendorRepo repo.VendorRepository,
	itemRepo repo.ItemRepository,
	taxRepo repo.TaxRepository,
) BillService {
	return &billService{
		billRepo:   billRepo,
		vendorRepo: vendorRepo,
		itemRepo:   itemRepo,
		taxRepo:    taxRepo,
	}
}

func (s *billService) CreateBill(billInput *input.CreateBillInput, userID string) (*output.BillOutput, error) {
	vendor, err := s.vendorRepo.FindByID(billInput.VendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	var tax *models.Tax
	if billInput.TaxID != nil {
		tax, err = s.taxRepo.FindByID(*billInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
	}

	lineItems := make([]models.BillLineItem, 0)
	subTotal := 0.0

	for _, itemInput := range billInput.LineItems {
		item, err := s.itemRepo.FindByID(itemInput.ItemID)
		if err != nil {
			return nil, errors.New("item not found: " + itemInput.ItemID)
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		lineItem := models.BillLineItem{
			ItemID:      itemInput.ItemID,
			Item:        item,
			VariantID:   itemInput.VariantID,
			Description: itemInput.Description,
			Account:     itemInput.Account,
			Quantity:    itemInput.Quantity,
			Rate:        itemInput.Rate,
			Amount:      amount,
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
		taxAmount = (subTotal - billInput.Discount) * (tax.Rate / 100)
	}

	total := subTotal - billInput.Discount + taxAmount + billInput.Adjustment

	bill := &models.Bill{
		ID:             uuid.New().String(),
		BillNumber:     fmt.Sprintf("BILL-%d", time.Now().Unix()),
		VendorID:       billInput.VendorID,
		Vendor:         vendor,
		BillingAddress: billInput.BillingAddress,
		OrderNumber:    billInput.OrderNumber,
		BillDate:       billInput.BillDate,
		DueDate:        billInput.DueDate,
		PaymentTerms:   domain.PaymentTerms(billInput.PaymentTerms),
		Subject:        billInput.Subject,
		LineItems:      lineItems,
		SubTotal:       subTotal,
		Discount:       billInput.Discount,
		TaxType:        (*domain.TaxType)(billInput.TaxType),
		TaxID:          billInput.TaxID,
		Tax:            tax,
		TaxAmount:      taxAmount,
		Adjustment:     billInput.Adjustment,
		Total:          total,
		Notes:          billInput.Notes,
		Status:         domain.BillStatusDraft,
		Attachments:    billInput.Attachments,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userID,
		UpdatedBy:      userID,
	}

	savedBill, err := s.billRepo.Create(bill)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(savedBill)
}

func (s *billService) GetBill(id string) (*output.BillOutput, error) {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return output.ToBillOutput(bill)
}

func (s *billService) GetAllBills(limit, offset int) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.BillOutput, len(bills))
	for i, bill := range bills {
		billOut, _ := output.ToBillOutput(&bill)
		outputs[i] = *billOut
	}

	return outputs, total, nil
}

func (s *billService) GetBillsByVendor(vendorID uint, limit, offset int) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindByVendor(vendorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.BillOutput, len(bills))
	for i, bill := range bills {
		billOut, _ := output.ToBillOutput(&bill)
		outputs[i] = *billOut
	}

	return outputs, total, nil
}

func (s *billService) GetBillsByStatus(status string, limit, offset int) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.BillOutput, len(bills))
	for i, bill := range bills {
		billOut, _ := output.ToBillOutput(&bill)
		outputs[i] = *billOut
	}

	return outputs, total, nil
}

func (s *billService) UpdateBill(id string, billInput *input.UpdateBillInput, userID string) (*output.BillOutput, error) {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	if billInput.VendorID != nil {
		vendor, err := s.vendorRepo.FindByID(*billInput.VendorID)
		if err != nil {
			return nil, errors.New("vendor not found")
		}
		bill.VendorID = *billInput.VendorID
		bill.Vendor = vendor
	}

	if billInput.BillingAddress != nil {
		bill.BillingAddress = *billInput.BillingAddress
	}

	if billInput.OrderNumber != nil {
		bill.OrderNumber = *billInput.OrderNumber
	}

	if billInput.BillDate != nil {
		bill.BillDate = *billInput.BillDate
	}

	if billInput.DueDate != nil {
		bill.DueDate = *billInput.DueDate
	}

	if billInput.PaymentTerms != nil {
		bill.PaymentTerms = domain.PaymentTerms(*billInput.PaymentTerms)
	}

	if billInput.Subject != nil {
		bill.Subject = *billInput.Subject
	}

	if len(billInput.LineItems) > 0 {
		lineItems := make([]models.BillLineItem, 0)
		subTotal := 0.0

		for _, itemInput := range billInput.LineItems {
			item, err := s.itemRepo.FindByID(itemInput.ItemID)
			if err != nil {
				return nil, errors.New("item not found: " + itemInput.ItemID)
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			lineItem := models.BillLineItem{
				ItemID:      itemInput.ItemID,
				Item:        item,
				VariantID:   itemInput.VariantID,
				Description: itemInput.Description,
				Account:     itemInput.Account,
				Quantity:    itemInput.Quantity,
				Rate:        itemInput.Rate,
				Amount:      amount,
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

		bill.LineItems = lineItems
		bill.SubTotal = subTotal
	}

	if billInput.Discount != nil {
		bill.Discount = *billInput.Discount
	}

	if billInput.TaxID != nil {
		tax, err := s.taxRepo.FindByID(*billInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
		bill.TaxID = billInput.TaxID
		bill.Tax = tax
	}

	if billInput.TaxType != nil {
		bill.TaxType = (*domain.TaxType)(billInput.TaxType)
	}

	if billInput.Adjustment != nil {
		bill.Adjustment = *billInput.Adjustment
	}

	if billInput.Notes != nil {
		bill.Notes = *billInput.Notes
	}

	if billInput.Attachments != nil {
		bill.Attachments = billInput.Attachments
	}

	if bill.Tax != nil {
		bill.TaxAmount = (bill.SubTotal - bill.Discount) * (bill.Tax.Rate / 100)
	} else {
		bill.TaxAmount = 0
	}

	bill.Total = bill.SubTotal - bill.Discount + bill.TaxAmount + bill.Adjustment
	bill.UpdatedAt = time.Now()
	bill.UpdatedBy = userID

	updatedBill, err := s.billRepo.Update(id, bill)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(updatedBill)
}

func (s *billService) UpdateBillStatus(id string, status string, userID string) (*output.BillOutput, error) {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	bill.Status = domain.BillStatus(status)
	bill.UpdatedAt = time.Now()
	bill.UpdatedBy = userID

	updatedBill, err := s.billRepo.Update(id, bill)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(updatedBill)
}

func (s *billService) DeleteBill(id string) error {
	return s.billRepo.Delete(id)
}
