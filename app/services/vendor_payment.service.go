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

type VendorPaymentService interface {
	// Existing methods retained.
	CreateVendorPayment(
		paymentInput *input.CreateVendorPaymentInput,
		userID, userName, companyName string,
		companyID uint,
	) (*output.VendorPaymentOutput, error)
	GetVendorPayment(id uint) (*output.VendorPaymentOutput, error)
	GetAllVendorPayments(limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByPurchaseOrder(poID string, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByVendor(vendorID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByStatus(status string, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	UpdateVendorPayment(id uint, paymentInput *input.UpdateVendorPaymentInput, userID string) (*output.VendorPaymentOutput, error)
	RecordPayment(
		id uint,
		paymentInput *input.RecordVendorPaymentInput,
		userID, userName, companyName string,
		companyID uint,
	) (*output.VendorPaymentOutput, error)
	DeleteVendorPayment(id uint) error

	// Company-scoped methods used by authenticated routes.
	CreateVendorPaymentForCompany(
		paymentInput *input.CreateVendorPaymentInput,
		userID, userName, companyName string,
		companyID uint,
	) (*output.VendorPaymentOutput, error)
	GetVendorPaymentByCompany(id uint, companyID uint) (*output.VendorPaymentOutput, error)
	GetAllVendorPaymentsByCompany(companyID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByPurchaseOrderAndCompany(poID string, companyID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByVendorAndCompany(vendorID, companyID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByStatusAndCompany(status string, companyID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	UpdateVendorPaymentForCompany(
		id uint,
		paymentInput *input.UpdateVendorPaymentInput,
		userID string,
		companyID uint,
	) (*output.VendorPaymentOutput, error)
	RecordPaymentForCompany(
		id uint,
		paymentInput *input.RecordVendorPaymentInput,
		userID, userName, companyName string,
		companyID uint,
	) (*output.VendorPaymentOutput, error)
	DeleteVendorPaymentForCompany(id uint, companyID uint) error
}

type vendorPaymentService struct {
	vendorPaymentRepo repo.VendorPaymentRepository
	poRepo            repo.PurchaseOrderRepository
	vendorRepo        repo.VendorRepository
	stockMgmtSvc      StockManagementService
	variantStockMgmt  VariantStockManagementService
	userRepo          repo.UserRepository
}

func NewVendorPaymentService(
	vendorPaymentRepo repo.VendorPaymentRepository,
	poRepo repo.PurchaseOrderRepository,
	vendorRepo repo.VendorRepository,
	stockMgmtSvc StockManagementService,
	variantStockMgmt VariantStockManagementService,
	userRepo repo.UserRepository,
) VendorPaymentService {
	return &vendorPaymentService{
		vendorPaymentRepo: vendorPaymentRepo,
		poRepo:            poRepo,
		vendorRepo:        vendorRepo,
		stockMgmtSvc:      stockMgmtSvc,
		variantStockMgmt:  variantStockMgmt,
		userRepo:          userRepo,
	}
}

func (s *vendorPaymentService) CreateVendorPayment(
	paymentInput *input.CreateVendorPaymentInput,
	userID, userName, companyName string,
	companyID uint,
) (*output.VendorPaymentOutput, error) {
	// Validate PurchaseOrder exists
	po, err := s.poRepo.FindByID(paymentInput.PurchaseOrderID)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	// Validate Vendor exists
	vendor, err := s.vendorRepo.FindByID(paymentInput.VendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	// Verify vendor matches purchase order vendor
	if vendor.ID != po.VendorID {
		return nil, errors.New("vendor does not match purchase order vendor")
	}

	// Get all existing payments for this PO to check for duplicates
	existingPayments, _, err := s.vendorPaymentRepo.FindByPurchaseOrderID(paymentInput.PurchaseOrderID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// Check if new payment amount is valid (just a sanity check, not a hard requirement)
	if paymentInput.Amount > po.Total {
		return nil, fmt.Errorf("payment amount exceeds PO total. Total PO: %.2f, Attempting to add: %.2f",
			po.Total, paymentInput.Amount)
	}

	paymentNumber := "VP-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:8]

	// Calculate remaining amount for the PO after this payment
	totalPaidAmount := 0.0
	for _, payment := range existingPayments {
		totalPaidAmount += payment.PaidAmount
	}
	totalPaidAmount += paymentInput.Amount
	newRemainingAmount := po.Total - totalPaidAmount
	if newRemainingAmount < 0 {
		newRemainingAmount = 0
	}

	// Determine payment status based on remaining amount (PO payment status)
	paymentStatus := domain.PaymentStatusPending
	if newRemainingAmount <= 0 {
		paymentStatus = domain.PaymentStatusCompleted
	} else if totalPaidAmount > 0 {
		paymentStatus = domain.PaymentStatusPartial
	}

	// Create payment record
	vendorPayment := models.VendorPayment{
		PaymentNumber:        paymentNumber,
		PurchaseOrderID:      paymentInput.PurchaseOrderID,
		VendorID:             paymentInput.VendorID,
		PaymentMode:          domain.PaymentMode(paymentInput.PaymentMode),
		Amount:               paymentInput.Amount,
		PaidAmount:           paymentInput.Amount, // Amount paid in this payment
		RemainingAmount:      newRemainingAmount,  // What's left on PO after this payment
		PaymentStatus:        paymentStatus,       // PO payment status: pending/partial/completed
		PaymentDate:          paymentInput.PaymentDate,
		ReferenceNumber:      paymentInput.ReferenceNumber,
		Notes:                paymentInput.Notes,
		CreatedByUserID:      userID,
		CreatedByUserName:    userName,
		CreatedByCompanyID:   companyID,
		CreatedByCompanyName: companyName,
	}

	createdPayment, err := s.vendorPaymentRepo.Create(&vendorPayment)
	if err != nil {
		return nil, err
	}

	// Update PO's payment tracking fields
	po.PaidAmount = totalPaidAmount
	po.RemainingAmount = newRemainingAmount
	if newRemainingAmount <= 0 {
		po.POPaymentStatus = domain.PaymentStatusCompleted
		po.RemainingAmount = 0
	} else if totalPaidAmount > 0 {
		po.POPaymentStatus = domain.PaymentStatusPartial
	}

	// Automatically update PO status to "received" when fully paid
	if po.POPaymentStatus == domain.PaymentStatusCompleted && po.Status == "draft" {
		po.Status = "received"
	}

	_, err = s.poRepo.Update(po.ID, po)
	if err != nil {
		return nil, err
	}

	// When PO is automatically transitioned to "received" after full payment, record stock once
	// Stock should be recorded exactly once: when PO reaches "received" status, not on every payment
	if po.Status == "received" && !po.InventorySynced {
		for _, lineItem := range po.LineItems {
			if lineItem.ProductID == nil || *lineItem.ProductID == "" {
				continue
			}

			stockQty := lineItem.Quantity
			stockRate := lineItem.Rate

			switch lineItem.PurchaseUnit {
			case "kg":
				stockQty = lineItem.Quantity * 1000
				stockRate = lineItem.Rate / 1000

			case "pieces":
				stockQty = lineItem.Quantity
				stockRate = lineItem.Rate

			default:
				stockQty = lineItem.Quantity
				stockRate = lineItem.Rate
			}

			if lineItem.SKU != "" {
				err := s.variantStockMgmt.RecordPurchaseInbound(
					*lineItem.ProductID,
					lineItem.SKU,
					stockQty,
					stockRate,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					userID,
				)

				if err != nil {
					fmt.Printf("[VP_CREATE] Error recording variant stock for SKU %s: %v\n", lineItem.SKU, err)
				}
			} else {
				err := s.stockMgmtSvc.RecordInboundMovement(
					*lineItem.ProductID,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					stockQty,
					stockRate,
					fmt.Sprintf("Received from vendor %s", po.Vendor.DisplayName),
					userID,
				)

				if err != nil {
					fmt.Printf("[VP_CREATE] Error recording stock for product %s: %v\n", *lineItem.ProductID, err)
				}
			}
		}

		po.InventorySynced = true
		now := time.Now()
		po.InventorySyncDate = &now

		_, err = s.poRepo.Update(po.ID, po)
		if err != nil {
			fmt.Printf("[VP_CREATE] Warning: Failed to mark inventory as synced: %v\n", err)
		}
	}

	return output.ConvertVendorPaymentToOutput(createdPayment), nil
}

func (s *vendorPaymentService) GetVendorPayment(id uint) (*output.VendorPaymentOutput, error) {
	payment, err := s.vendorPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("vendor payment not found")
	}
	return output.ConvertVendorPaymentToOutput(payment), nil
}

func (s *vendorPaymentService) GetAllVendorPayments(limit, offset int) ([]output.VendorPaymentOutput, int64, error) {
	payments, total, err := s.vendorPaymentRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) GetVendorPaymentsByPurchaseOrder(poID string, limit, offset int) ([]output.VendorPaymentOutput, int64, error) {
	payments, total, err := s.vendorPaymentRepo.FindByPurchaseOrderID(poID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) GetVendorPaymentsByVendor(vendorID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error) {
	payments, total, err := s.vendorPaymentRepo.FindByVendorID(vendorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) GetVendorPaymentsByStatus(status string, limit, offset int) ([]output.VendorPaymentOutput, int64, error) {
	payments, total, err := s.vendorPaymentRepo.FindByPaymentStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) UpdateVendorPayment(
	id uint,
	paymentInput *input.UpdateVendorPaymentInput,
	userID string,
) (*output.VendorPaymentOutput, error) {
	payment, err := s.vendorPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("vendor payment not found")
	}

	if paymentInput.PaymentMode != nil {
		payment.PaymentMode = domain.PaymentMode(*paymentInput.PaymentMode)
	}
	if paymentInput.Amount != nil {
		payment.Amount = *paymentInput.Amount
	}
	if paymentInput.PaymentDate != nil {
		payment.PaymentDate = *paymentInput.PaymentDate
	}
	if paymentInput.ReferenceNumber != nil {
		payment.ReferenceNumber = *paymentInput.ReferenceNumber
	}
	if paymentInput.Notes != nil {
		payment.Notes = *paymentInput.Notes
	}

	updatedPayment, err := s.vendorPaymentRepo.Update(id, payment)
	if err != nil {
		return nil, err
	}

	return output.ConvertVendorPaymentToOutput(updatedPayment), nil
}

func (s *vendorPaymentService) RecordPayment(
	id uint,
	paymentInput *input.RecordVendorPaymentInput,
	userID, userName, companyName string,
	companyID uint,
) (*output.VendorPaymentOutput, error) {
	payment, err := s.vendorPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("vendor payment not found")
	}

	// Fetch PO to check line items
	po, err := s.poRepo.FindByID(payment.PurchaseOrderID)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	// Update payment record details
	payment.PaymentMode = domain.PaymentMode(paymentInput.PaymentMode)
	payment.ReferenceNumber = paymentInput.ReferenceNumber
	payment.Notes = paymentInput.Notes

	updatedPayment, err := s.vendorPaymentRepo.Update(id, payment)
	if err != nil {
		return nil, err
	}

	// If payment is being made, update PO payment tracking
	if updatedPayment.PaymentStatus == domain.PaymentStatusPartial || updatedPayment.PaymentStatus == domain.PaymentStatusCompleted {
		// Get all payments for this PO
		allPayments, _, err := s.vendorPaymentRepo.FindByPurchaseOrderID(payment.PurchaseOrderID, 1000, 0)
		if err == nil && len(allPayments) > 0 {
			// Calculate total paid across all payments
			totalPaidAmount := 0.0
			for _, p := range allPayments {
				totalPaidAmount += p.PaidAmount
			}

			// Update PO's payment tracking
			po.PaidAmount = totalPaidAmount
			po.RemainingAmount = po.Total - totalPaidAmount
			if po.RemainingAmount < 0 {
				po.RemainingAmount = 0
			}

			// Determine PO payment status
			if po.RemainingAmount <= 0 {
				po.POPaymentStatus = domain.PaymentStatusCompleted
				po.RemainingAmount = 0
			} else if po.PaidAmount > 0 {
				po.POPaymentStatus = domain.PaymentStatusPartial
			}

			// Automatically update PO status to "received" when fully paid
			if po.POPaymentStatus == domain.PaymentStatusCompleted && po.Status == "draft" {
				po.Status = "received"
			}

			// Update the PO
			_, err = s.poRepo.Update(po.ID, po)
			if err != nil {
				fmt.Printf("warning: failed to update PO payment status: %v\n", err)
			}
		}
	}

	if po.Status == "received" && !po.InventorySynced {
		for _, lineItem := range po.LineItems {
			if lineItem.ProductID == nil || *lineItem.ProductID == "" {
				continue
			}

			stockQty := lineItem.Quantity
			stockRate := lineItem.Rate

			// Convert kg to grams
			if lineItem.PurchaseUnit == "kg" {
				stockQty = lineItem.Quantity * 1000
				stockRate = lineItem.Rate / 1000
			}

			if lineItem.SKU != "" {
				err := s.variantStockMgmt.RecordPurchaseInbound(
					*lineItem.ProductID,
					lineItem.SKU,
					stockQty,  // 50 kg becomes 50000 grams
					stockRate, // 150/kg becomes 0.15/gram
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					userID,
				)

				if err != nil {
					fmt.Printf("[VP_CREATE] Error recording variant stock for SKU %s: %v\n", lineItem.SKU, err)
				}
			} else {
				err := s.stockMgmtSvc.RecordInboundMovement(
					*lineItem.ProductID,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					stockQty,
					stockRate,
					fmt.Sprintf("Received from vendor %s", po.Vendor.DisplayName),
					userID,
				)

				if err != nil {
					fmt.Printf("[VP_CREATE] Error recording stock for product %s: %v\n", *lineItem.ProductID, err)
				}
			}
		}

		po.InventorySynced = true
		now := time.Now()
		po.InventorySyncDate = &now

		_, err = s.poRepo.Update(po.ID, po)
		if err != nil {
			fmt.Printf("[VP_CREATE] Warning: Failed to mark inventory as synced: %v\n", err)
		}
	}

	return output.ConvertVendorPaymentToOutput(updatedPayment), nil
}

func (s *vendorPaymentService) DeleteVendorPayment(id uint) error {
	return s.vendorPaymentRepo.Delete(id)
}

func (s *vendorPaymentService) validateVendorPaymentUserCompany(
	userID string,
	companyID uint,
) error {
	if companyID == 0 {
		return errors.New("invalid company")
	}

	var parsedUserID uint
	if _, err := fmt.Sscanf(userID, "%d", &parsedUserID); err != nil || parsedUserID == 0 {
		return errors.New("invalid authenticated user")
	}

	user, err := s.userRepo.GetByIDAndCompanyID(parsedUserID, companyID)
	if err != nil || user == nil {
		return errors.New("user does not belong to the company")
	}

	return nil
}

func (s *vendorPaymentService) validateVendorPaymentReferences(
	paymentInput *input.CreateVendorPaymentInput,
	companyID uint,
) error {
	if paymentInput == nil {
		return errors.New("input cannot be nil")
	}

	po, err := s.poRepo.FindByIDAndCompany(
		paymentInput.PurchaseOrderID,
		companyID,
	)
	if err != nil {
		return errors.New("purchase order not found in your company")
	}

	vendor, err := s.vendorRepo.FindByIDAndCompany(
		paymentInput.VendorID,
		companyID,
	)
	if err != nil {
		return errors.New("vendor not found in your company")
	}

	if po.VendorID != vendor.ID {
		return errors.New("vendor does not match purchase order vendor")
	}

	return nil
}

func (s *vendorPaymentService) CreateVendorPaymentForCompany(
	paymentInput *input.CreateVendorPaymentInput,
	userID, userName, companyName string,
	companyID uint,
) (*output.VendorPaymentOutput, error) {
	if err := s.validateVendorPaymentUserCompany(userID, companyID); err != nil {
		return nil, err
	}

	if err := s.validateVendorPaymentReferences(paymentInput, companyID); err != nil {
		return nil, err
	}

	payment, err := s.CreateVendorPayment(
		paymentInput,
		userID,
		userName,
		companyName,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	return s.GetVendorPaymentByCompany(payment.ID, companyID)
}

func (s *vendorPaymentService) GetVendorPaymentByCompany(
	id uint,
	companyID uint,
) (*output.VendorPaymentOutput, error) {
	payment, err := s.vendorPaymentRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("vendor payment not found")
	}

	return output.ConvertVendorPaymentToOutput(payment), nil
}

func (s *vendorPaymentService) GetAllVendorPaymentsByCompany(
	companyID uint,
	limit, offset int,
) ([]output.VendorPaymentOutput, int64, error) {
	payments, total, err := s.vendorPaymentRepo.FindAllByCompany(
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) GetVendorPaymentsByPurchaseOrderAndCompany(
	poID string,
	companyID uint,
	limit, offset int,
) ([]output.VendorPaymentOutput, int64, error) {
	if _, err := s.poRepo.FindByIDAndCompany(poID, companyID); err != nil {
		return nil, 0, errors.New("purchase order not found in your company")
	}

	payments, total, err := s.vendorPaymentRepo.FindByPurchaseOrderIDAndCompany(
		poID,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) GetVendorPaymentsByVendorAndCompany(
	vendorID, companyID uint,
	limit, offset int,
) ([]output.VendorPaymentOutput, int64, error) {
	if _, err := s.vendorRepo.FindByIDAndCompany(vendorID, companyID); err != nil {
		return nil, 0, errors.New("vendor not found in your company")
	}

	payments, total, err := s.vendorPaymentRepo.FindByVendorIDAndCompany(
		vendorID,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) GetVendorPaymentsByStatusAndCompany(
	status string,
	companyID uint,
	limit, offset int,
) ([]output.VendorPaymentOutput, int64, error) {
	payments, total, err := s.vendorPaymentRepo.FindByPaymentStatusAndCompany(
		status,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertVendorPaymentsToOutput(payments), total, nil
}

func (s *vendorPaymentService) UpdateVendorPaymentForCompany(
	id uint,
	paymentInput *input.UpdateVendorPaymentInput,
	userID string,
	companyID uint,
) (*output.VendorPaymentOutput, error) {
	if err := s.validateVendorPaymentUserCompany(userID, companyID); err != nil {
		return nil, err
	}

	payment, err := s.vendorPaymentRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("vendor payment not found")
	}

	if paymentInput == nil {
		return nil, errors.New("input cannot be nil")
	}

	if paymentInput.PaymentMode != nil {
		payment.PaymentMode = domain.PaymentMode(*paymentInput.PaymentMode)
	}
	if paymentInput.Amount != nil {
		payment.Amount = *paymentInput.Amount
	}
	if paymentInput.PaymentDate != nil {
		payment.PaymentDate = *paymentInput.PaymentDate
	}
	if paymentInput.ReferenceNumber != nil {
		payment.ReferenceNumber = *paymentInput.ReferenceNumber
	}
	if paymentInput.Notes != nil {
		payment.Notes = *paymentInput.Notes
	}

	updated, err := s.vendorPaymentRepo.UpdateByCompany(
		id,
		companyID,
		payment,
	)
	if err != nil {
		return nil, err
	}

	return output.ConvertVendorPaymentToOutput(updated), nil
}

func (s *vendorPaymentService) RecordPaymentForCompany(
	id uint,
	paymentInput *input.RecordVendorPaymentInput,
	userID, userName, companyName string,
	companyID uint,
) (*output.VendorPaymentOutput, error) {
	if err := s.validateVendorPaymentUserCompany(userID, companyID); err != nil {
		return nil, err
	}

	payment, err := s.vendorPaymentRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("vendor payment not found")
	}

	if _, err := s.poRepo.FindByIDAndCompany(
		payment.PurchaseOrderID,
		companyID,
	); err != nil {
		return nil, errors.New("purchase order not found in your company")
	}

	if _, err := s.vendorRepo.FindByIDAndCompany(
		payment.VendorID,
		companyID,
	); err != nil {
		return nil, errors.New("vendor not found in your company")
	}

	if _, err := s.RecordPayment(
		id,
		paymentInput,
		userID,
		userName,
		companyName,
		companyID,
	); err != nil {
		return nil, err
	}

	return s.GetVendorPaymentByCompany(id, companyID)
}

func (s *vendorPaymentService) DeleteVendorPaymentForCompany(
	id uint,
	companyID uint,
) error {
	if _, err := s.vendorPaymentRepo.FindByIDAndCompany(id, companyID); err != nil {
		return errors.New("vendor payment not found")
	}

	return s.vendorPaymentRepo.DeleteByCompany(id, companyID)
}
