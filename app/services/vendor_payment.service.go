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
	// Basic CRUD Operations
	CreateVendorPayment(paymentInput *input.CreateVendorPaymentInput, userID, userName, companyName string, companyID uint) (*output.VendorPaymentOutput, error)
	GetVendorPayment(id uint) (*output.VendorPaymentOutput, error)
	GetAllVendorPayments(limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByPurchaseOrder(poID string, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByVendor(vendorID uint, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	GetVendorPaymentsByStatus(status string, limit, offset int) ([]output.VendorPaymentOutput, int64, error)
	UpdateVendorPayment(id uint, paymentInput *input.UpdateVendorPaymentInput, userID string) (*output.VendorPaymentOutput, error)

	// Payment recording and status management
	RecordPayment(id uint, paymentInput *input.RecordVendorPaymentInput, userID, userName, companyName string, companyID uint) (*output.VendorPaymentOutput, error)
	DeleteVendorPayment(id uint) error
}

type vendorPaymentService struct {
	vendorPaymentRepo repo.VendorPaymentRepository
	poRepo            repo.PurchaseOrderRepository
	vendorRepo        repo.VendorRepository
	stockMgmtSvc      StockManagementService
	variantStockMgmt  VariantStockManagementService
}

func NewVendorPaymentService(
	vendorPaymentRepo repo.VendorPaymentRepository,
	poRepo repo.PurchaseOrderRepository,
	vendorRepo repo.VendorRepository,
	stockMgmtSvc StockManagementService,
	variantStockMgmt VariantStockManagementService,
) VendorPaymentService {
	return &vendorPaymentService{
		vendorPaymentRepo: vendorPaymentRepo,
		poRepo:            poRepo,
		vendorRepo:        vendorRepo,
		stockMgmtSvc:      stockMgmtSvc,
		variantStockMgmt:  variantStockMgmt,
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

	// Check for duplicate payment (same amount, vendor, and date within 1 hour)
	for _, payment := range existingPayments {
		timeDiff := paymentInput.PaymentDate.Sub(payment.PaymentDate).Abs()
		if payment.Amount == paymentInput.Amount && timeDiff.Hours() < 1 {
			return nil, fmt.Errorf("duplicate payment detected. Payment of %.2f already exists for this PO", paymentInput.Amount)
		}
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

			if lineItem.SKU != "" {
				// Variant purchase - record at variant level only
				err := s.variantStockMgmt.RecordPurchaseInbound(
					*lineItem.ProductID,
					lineItem.SKU,
					lineItem.Quantity,
					lineItem.Rate,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					userID,
				)
				if err != nil {
					fmt.Printf("[VP_CREATE] Error recording variant stock for SKU %s: %v\n", lineItem.SKU, err)
				}
			} else {
				// Base product (no SKU) - record at product level only
				err := s.stockMgmtSvc.RecordInboundMovement(
					*lineItem.ProductID,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					lineItem.Quantity,
					lineItem.Rate,
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

	// When PO is automatically transitioned to "received" after full payment, record stock once
	// Stock should be recorded exactly once: when PO reaches "received" status, not on every payment
	if po.Status == "received" && !po.InventorySynced {
		for _, lineItem := range po.LineItems {
			if lineItem.ProductID == nil || *lineItem.ProductID == "" {
				continue
			}

			if lineItem.SKU != "" {
				// Variant purchase - record at variant level only
				err := s.variantStockMgmt.RecordPurchaseInbound(
					*lineItem.ProductID,
					lineItem.SKU,
					lineItem.Quantity,
					lineItem.Rate,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					userID,
				)
				if err != nil {
					fmt.Printf("[VP_RECORD] Error recording variant stock for SKU %s: %v\n", lineItem.SKU, err)
				}
			} else {
				// Base product (no SKU) - record at product level only
				err := s.stockMgmtSvc.RecordInboundMovement(
					*lineItem.ProductID,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					lineItem.Quantity,
					lineItem.Rate,
					fmt.Sprintf("Received from vendor %s", po.Vendor.DisplayName),
					userID,
				)
				if err != nil {
					fmt.Printf("[VP_RECORD] Error recording stock for product %s: %v\n", *lineItem.ProductID, err)
				}
			}
		}
		po.InventorySynced = true
		now := time.Now()
		po.InventorySyncDate = &now
		_, err = s.poRepo.Update(po.ID, po)
		if err != nil {
			fmt.Printf("[VP_RECORD] Warning: Failed to mark inventory as synced: %v\n", err)
		}
	}

	return output.ConvertVendorPaymentToOutput(updatedPayment), nil
}

func (s *vendorPaymentService) DeleteVendorPayment(id uint) error {
	return s.vendorPaymentRepo.Delete(id)
}
