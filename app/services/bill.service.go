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
	// Existing methods retained for compatibility.
	CreateBill(
		billInput *input.CreateBillInput,
		userID string,
	) (*output.BillOutput, error)

	GetBill(id string) (*output.BillOutput, error)

	GetAllBills(
		limit int,
		offset int,
		createdBy string,
	) ([]output.BillOutput, int64, error)

	GetBillsByVendor(
		vendorID uint,
		limit int,
		offset int,
		createdBy string,
	) ([]output.BillOutput, int64, error)

	GetBillsByStatus(
		status string,
		limit int,
		offset int,
		createdBy string,
	) ([]output.BillOutput, int64, error)

	UpdateBill(
		id string,
		billInput *input.UpdateBillInput,
		userID string,
	) (*output.BillOutput, error)

	UpdateBillStatus(
		id string,
		status string,
		userID string,
	) (*output.BillOutput, error)

	DeleteBill(
		id string,
		createdBy string,
	) error

	// Company-scoped methods.
	CreateBillForCompany(
		billInput *input.CreateBillInput,
		userID string,
		companyID uint,
	) (*output.BillOutput, error)

	GetBillByCompany(
		id string,
		companyID uint,
	) (*output.BillOutput, error)

	GetAllBillsByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]output.BillOutput, int64, error)

	GetBillsByVendorAndCompany(
		vendorID uint,
		companyID uint,
		limit int,
		offset int,
	) ([]output.BillOutput, int64, error)

	GetBillsByStatusAndCompany(
		status string,
		companyID uint,
		limit int,
		offset int,
	) ([]output.BillOutput, int64, error)

	UpdateBillForCompany(
		id string,
		billInput *input.UpdateBillInput,
		userID string,
		companyID uint,
	) (*output.BillOutput, error)

	UpdateBillStatusForCompany(
		id string,
		status string,
		userID string,
		companyID uint,
	) (*output.BillOutput, error)

	DeleteBillForCompany(
		id string,
		companyID uint,
	) error

	ValidateBillInputForCompany(
		billInput *input.CreateBillInput,
		companyID uint,
	) error
}

type billService struct {
	billRepo    repo.BillRepository
	vendorRepo  repo.VendorRepository
	productRepo repo.ProductRepository
	taxRepo     repo.TaxRepository
	userRepo    repo.UserRepository
}

func NewBillService(
	billRepo repo.BillRepository,
	vendorRepo repo.VendorRepository,
	productRepo repo.ProductRepository,
	taxRepo repo.TaxRepository,
	userRepo repo.UserRepository,
) BillService {
	return &billService{
		billRepo:    billRepo,
		vendorRepo:  vendorRepo,
		productRepo: productRepo,
		taxRepo:     taxRepo,
		userRepo:    userRepo,
	}
}

func (s *billService) CreateBill(
	billInput *input.CreateBillInput,
	userID string,
) (*output.BillOutput, error) {
	if billInput == nil {
		return nil, errors.New("input cannot be nil")
	}

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

	lineItems := make([]models.BillLineItem, 0, len(billInput.LineItems))
	subTotal := 0.0

	for _, itemInput := range billInput.LineItems {
		if itemInput.ProductID == nil || *itemInput.ProductID == "" {
			return nil, errors.New("product_id is required for each line item")
		}

		product, err := s.productRepo.FindByID(*itemInput.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + *itemInput.ProductID)
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		lineItems = append(lineItems, models.BillLineItem{
			ProductID:   itemInput.ProductID,
			Product:     product,
			ProductName: itemInput.ProductName,
			Description: itemInput.ProductName,
			Account:     itemInput.Account,
			Quantity:    itemInput.Quantity,
			Rate:        itemInput.Rate,
			Amount:      amount,
		})
	}

	taxAmount := 0.0
	if tax != nil {
		taxAmount = (subTotal - billInput.Discount) * (tax.Rate / 100)
	}

	total := subTotal -
		billInput.Discount +
		taxAmount +
		billInput.Adjustment

	bill := &models.Bill{
		ID:              uuid.New().String(),
		BillNumber:      billInput.BillNumber,
		VendorID:        billInput.VendorID,
		Vendor:          vendor,
		PurchaseOrderID: billInput.PurchaseOrderID,
		BillingAddress:  billInput.BillingAddress,
		OrderNumber:     billInput.OrderNumber,
		BillDate:        billInput.BillDate,
		DueDate:         billInput.DueDate,
		PaymentTerms:    domain.PaymentTerms(billInput.PaymentTerms),
		Subject:         billInput.Subject,
		LineItems:       lineItems,
		SubTotal:        subTotal,
		Discount:        billInput.Discount,
		TaxType:         (*domain.TaxType)(billInput.TaxType),
		TaxID:           billInput.TaxID,
		Tax:             tax,
		TaxAmount:       taxAmount,
		Adjustment:      billInput.Adjustment,
		Total:           total,
		Notes:           billInput.Notes,
		Status:          domain.BillStatusDraft,
		Attachments:     billInput.Attachments,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       userID,
		UpdatedBy:       userID,
	}

	savedBill, err := s.billRepo.Create(bill)
	if err != nil {
		return nil, err
	}

	reloadedBill, err := s.billRepo.FindByID(savedBill.ID)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(reloadedBill)
}

func (s *billService) GetBill(
	id string,
) (*output.BillOutput, error) {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(bill)
}

func (s *billService) GetAllBills(
	limit int,
	offset int,
	createdBy string,
) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindByCreatedBy(
		createdBy,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	results, err := billOutputs(bills)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *billService) GetBillsByVendor(
	vendorID uint,
	limit int,
	offset int,
	createdBy string,
) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindByVendorAndCreatedBy(
		vendorID,
		createdBy,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	results, err := billOutputs(bills)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *billService) GetBillsByStatus(
	status string,
	limit int,
	offset int,
	createdBy string,
) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindByStatusAndCreatedBy(
		status,
		createdBy,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	results, err := billOutputs(bills)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *billService) UpdateBill(
	id string,
	billInput *input.UpdateBillInput,
	userID string,
) (*output.BillOutput, error) {
	if billInput == nil {
		return nil, errors.New("input cannot be nil")
	}

	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	if err := s.applyBillUpdate(
		bill,
		billInput,
		userID,
		0,
		false,
	); err != nil {
		return nil, err
	}

	updatedBill, err := s.billRepo.Update(id, bill)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(updatedBill)
}

func (s *billService) UpdateBillStatus(
	id string,
	status string,
	userID string,
) (*output.BillOutput, error) {
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

func (s *billService) DeleteBill(
	id string,
	createdBy string,
) error {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return err
	}

	if bill.CreatedBy != createdBy {
		return errors.New(
			"unauthorized: you can only delete bills you created",
		)
	}

	return s.billRepo.Delete(id)
}

func (s *billService) ValidateBillInputForCompany(
	billInput *input.CreateBillInput,
	companyID uint,
) error {
	if billInput == nil {
		return errors.New("input cannot be nil")
	}

	if companyID == 0 {
		return errors.New("invalid company")
	}

	if _, err := s.vendorRepo.FindByIDAndCompany(
		billInput.VendorID,
		companyID,
	); err != nil {
		return errors.New("vendor not found in your company")
	}

	for _, item := range billInput.LineItems {
		if item.ProductID == nil || *item.ProductID == "" {
			return errors.New(
				"product_id is required for each line item",
			)
		}

		if _, err := s.productRepo.FindByIDAndCompany(
			*item.ProductID,
			companyID,
		); err != nil {
			return fmt.Errorf(
				"product %s not found in your company",
				*item.ProductID,
			)
		}
	}

	return nil
}

func (s *billService) validateBillUpdateForCompany(
	billInput *input.UpdateBillInput,
	companyID uint,
) error {
	if billInput == nil {
		return errors.New("input cannot be nil")
	}

	if billInput.VendorID != nil {
		if _, err := s.vendorRepo.FindByIDAndCompany(
			*billInput.VendorID,
			companyID,
		); err != nil {
			return errors.New(
				"vendor not found in your company",
			)
		}
	}

	for _, item := range billInput.LineItems {
		if item.ProductID == nil || *item.ProductID == "" {
			return errors.New(
				"product_id is required for each line item",
			)
		}

		if _, err := s.productRepo.FindByIDAndCompany(
			*item.ProductID,
			companyID,
		); err != nil {
			return fmt.Errorf(
				"product %s not found in your company",
				*item.ProductID,
			)
		}
	}

	return nil
}

func (s *billService) CreateBillForCompany(
	billInput *input.CreateBillInput,
	userID string,
	companyID uint,
) (*output.BillOutput, error) {
	if err := s.ValidateBillInputForCompany(
		billInput,
		companyID,
	); err != nil {
		return nil, err
	}

	userIDUint, err := parseBillUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByIDAndCompanyID(
		userIDUint,
		companyID,
	)
	if err != nil || user == nil {
		return nil, errors.New(
			"user does not belong to the company",
		)
	}

	createdBill, err := s.CreateBill(
		billInput,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return s.GetBillByCompany(
		createdBill.ID,
		companyID,
	)
}

func (s *billService) GetBillByCompany(
	id string,
	companyID uint,
) (*output.BillOutput, error) {
	bill, err := s.billRepo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	return output.ToBillOutput(bill)
}

func (s *billService) GetAllBillsByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]output.BillOutput, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	bills, total, err := s.billRepo.FindAllByCompany(
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	results, err := billOutputs(bills)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *billService) GetBillsByVendorAndCompany(
	vendorID uint,
	companyID uint,
	limit int,
	offset int,
) ([]output.BillOutput, int64, error) {
	if _, err := s.vendorRepo.FindByIDAndCompany(
		vendorID,
		companyID,
	); err != nil {
		return nil, 0, errors.New(
			"vendor not found in your company",
		)
	}

	bills, total, err := s.billRepo.FindByVendorAndCompany(
		vendorID,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	results, err := billOutputs(bills)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *billService) GetBillsByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]output.BillOutput, int64, error) {
	bills, total, err := s.billRepo.FindByStatusAndCompany(
		status,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	results, err := billOutputs(bills)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *billService) UpdateBillForCompany(
	id string,
	billInput *input.UpdateBillInput,
	userID string,
	companyID uint,
) (*output.BillOutput, error) {
	if billInput == nil {
		return nil, errors.New("input cannot be nil")
	}

	bill, err := s.billRepo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	if err := s.validateBillUpdateForCompany(
		billInput,
		companyID,
	); err != nil {
		return nil, err
	}

	if err := s.applyBillUpdate(
		bill,
		billInput,
		userID,
		companyID,
		true,
	); err != nil {
		return nil, err
	}

	updatedBill, err := s.billRepo.UpdateByCompany(
		id,
		companyID,
		bill,
	)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(updatedBill)
}

func (s *billService) UpdateBillStatusForCompany(
	id string,
	status string,
	userID string,
	companyID uint,
) (*output.BillOutput, error) {
	bill, err := s.billRepo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	bill.Status = domain.BillStatus(status)
	bill.UpdatedAt = time.Now()
	bill.UpdatedBy = userID

	updatedBill, err := s.billRepo.UpdateByCompany(
		id,
		companyID,
		bill,
	)
	if err != nil {
		return nil, err
	}

	return output.ToBillOutput(updatedBill)
}

func (s *billService) DeleteBillForCompany(
	id string,
	companyID uint,
) error {
	return s.billRepo.DeleteByCompany(
		id,
		companyID,
	)
}

func (s *billService) applyBillUpdate(
	bill *models.Bill,
	billInput *input.UpdateBillInput,
	userID string,
	companyID uint,
	useCompanyFilter bool,
) error {
	if billInput.VendorID != nil {
		var vendor *models.Vendor
		var err error

		if useCompanyFilter {
			vendor, err = s.vendorRepo.FindByIDAndCompany(
				*billInput.VendorID,
				companyID,
			)
		} else {
			vendor, err = s.vendorRepo.FindByID(
				*billInput.VendorID,
			)
		}

		if err != nil {
			return errors.New("vendor not found")
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
		bill.PaymentTerms = domain.PaymentTerms(
			*billInput.PaymentTerms,
		)
	}

	if billInput.Subject != nil {
		bill.Subject = *billInput.Subject
	}

	if len(billInput.LineItems) > 0 {
		lineItems := make(
			[]models.BillLineItem,
			0,
			len(billInput.LineItems),
		)
		subTotal := 0.0

		for _, itemInput := range billInput.LineItems {
			if itemInput.ProductID == nil ||
				*itemInput.ProductID == "" {
				return errors.New(
					"product_id is required for each line item",
				)
			}

			var product *models.Product
			var err error

			if useCompanyFilter {
				product, err = s.productRepo.FindByIDAndCompany(
					*itemInput.ProductID,
					companyID,
				)
			} else {
				product, err = s.productRepo.FindByID(
					*itemInput.ProductID,
				)
			}

			if err != nil {
				return errors.New(
					"product not found: " +
						*itemInput.ProductID,
				)
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			lineItems = append(
				lineItems,
				models.BillLineItem{
					ProductID:   itemInput.ProductID,
					Product:     product,
					ProductName: itemInput.ProductName,
					Description: itemInput.ProductName,
					Account:     itemInput.Account,
					Quantity:    itemInput.Quantity,
					Rate:        itemInput.Rate,
					Amount:      amount,
				},
			)
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
			return errors.New("tax not found")
		}

		bill.TaxID = billInput.TaxID
		bill.Tax = tax
	}

	if billInput.TaxType != nil {
		bill.TaxType = (*domain.TaxType)(
			billInput.TaxType,
		)
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
		bill.TaxAmount = (bill.SubTotal - bill.Discount) *
			(bill.Tax.Rate / 100)
	} else {
		bill.TaxAmount = 0
	}

	bill.Total = bill.SubTotal -
		bill.Discount +
		bill.TaxAmount +
		bill.Adjustment

	bill.UpdatedAt = time.Now()
	bill.UpdatedBy = userID

	return nil
}

func parseBillUserID(userID string) (uint, error) {
	var parsed uint

	if _, err := fmt.Sscanf(userID, "%d", &parsed); err != nil ||
		parsed == 0 {
		return 0, errors.New("invalid authenticated user")
	}

	return parsed, nil
}

func billOutputs(
	bills []models.Bill,
) ([]output.BillOutput, error) {
	results := make([]output.BillOutput, len(bills))

	for index := range bills {
		billOutput, err := output.ToBillOutput(&bills[index])
		if err != nil {
			return nil, err
		}

		results[index] = *billOutput
	}

	return results, nil
}
