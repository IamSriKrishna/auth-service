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

type CustomerPaymentService interface {
	// Basic CRUD Operations
	CreateCustomerPayment(paymentInput *input.CreateCustomerPaymentInput, userID, userName, companyName string, companyID uint) (*output.CustomerPaymentOutput, error)
	GetCustomerPayment(id uint) (*output.CustomerPaymentOutput, error)
	GetAllCustomerPayments(limit, offset int) ([]output.CustomerPaymentOutput, int64, error)
	GetCustomerPaymentsBySalesOrder(soID string, limit, offset int) ([]output.CustomerPaymentOutput, int64, error)
	GetCustomerPaymentsByCustomer(customerID uint, limit, offset int) ([]output.CustomerPaymentOutput, int64, error)
	GetCustomerPaymentsByStatus(status string, limit, offset int) ([]output.CustomerPaymentOutput, int64, error)
	UpdateCustomerPayment(id uint, paymentInput *input.UpdateCustomerPaymentInput, userID string) (*output.CustomerPaymentOutput, error)

	// Payment recording and status management
	RecordPayment(id uint, paymentInput *input.RecordCustomerPaymentInput, userID, userName, companyName string, companyID uint) (*output.CustomerPaymentOutput, error)
	DeleteCustomerPayment(id uint) error
}

type customerPaymentService struct {
	customerPaymentRepo repo.CustomerPaymentRepository
	soRepo              repo.SalesOrderRepository
	customerRepo        repo.CustomerRepository
}

func NewCustomerPaymentService(
	customerPaymentRepo repo.CustomerPaymentRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
) CustomerPaymentService {
	return &customerPaymentService{
		customerPaymentRepo: customerPaymentRepo,
		soRepo:              soRepo,
		customerRepo:        customerRepo,
	}
}

func (s *customerPaymentService) CreateCustomerPayment(
	paymentInput *input.CreateCustomerPaymentInput,
	userID, userName, companyName string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	// Validate SalesOrder exists
	so, err := s.soRepo.FindByID(paymentInput.SalesOrderID)
	if err != nil {
		return nil, errors.New("sales order not found")
	}

	// Validate Customer exists
	customer, err := s.customerRepo.FindByID(paymentInput.CustomerID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	// Verify customer matches sales order customer
	if customer.ID != so.CustomerID {
		return nil, errors.New("customer does not match sales order customer")
	}

	// Get all existing payments for this SO to check for duplicates
	existingPayments, _, err := s.customerPaymentRepo.FindBySalesOrderID(paymentInput.SalesOrderID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// Check if new payment amount is valid (just a sanity check, not a hard requirement)
	if paymentInput.Amount > so.Total {
		return nil, fmt.Errorf("payment amount exceeds SO total. Total SO: %.2f, Attempting to add: %.2f",
			so.Total, paymentInput.Amount)
	}

	paymentNumber := "CP-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:8]

	// Calculate remaining amount for the SO after this payment
	totalReceivedAmount := 0.0
	for _, payment := range existingPayments {
		totalReceivedAmount += payment.ReceivedAmount
	}
	totalReceivedAmount += paymentInput.Amount
	newRemainingAmount := so.Total - totalReceivedAmount
	if newRemainingAmount < 0 {
		newRemainingAmount = 0
	}

	// Determine payment status based on remaining amount (SO payment status)
	paymentStatus := domain.PaymentStatusPending
	if newRemainingAmount <= 0 {
		paymentStatus = domain.PaymentStatusCompleted
	} else if totalReceivedAmount > 0 {
		paymentStatus = domain.PaymentStatusPartial
	}

	// Create payment record
	customerPayment := models.CustomerPayment{
		PaymentNumber:        paymentNumber,
		SalesOrderID:         paymentInput.SalesOrderID,
		CustomerID:           paymentInput.CustomerID,
		PaymentMode:          domain.PaymentMode(paymentInput.PaymentMode),
		Amount:               paymentInput.Amount,
		ReceivedAmount:       paymentInput.Amount, // Amount received in this payment
		RemainingAmount:      newRemainingAmount,  // What's left on SO after this payment
		PaymentStatus:        paymentStatus,       // SO payment status: pending/partial/completed
		PaymentDate:          paymentInput.PaymentDate,
		ReferenceNumber:      paymentInput.ReferenceNumber,
		Notes:                paymentInput.Notes,
		CreatedByUserID:      userID,
		CreatedByUserName:    userName,
		CreatedByCompanyID:   companyID,
		CreatedByCompanyName: companyName,
	}

	createdPayment, err := s.customerPaymentRepo.Create(&customerPayment)
	if err != nil {
		return nil, err
	}

	// Update SO's payment tracking fields
	so.PaidAmount = totalReceivedAmount
	so.RemainingAmount = newRemainingAmount
	if newRemainingAmount <= 0 {
		// Update SO status to "paid" when fully paid
		so.Status = domain.SalesOrderStatusPaid
		so.RemainingAmount = 0
	}

	_, err = s.soRepo.Update(so.ID, so)
	if err != nil {
		return nil, err
	}

	return output.ConvertCustomerPaymentToOutput(createdPayment), nil
}

func (s *customerPaymentService) GetCustomerPayment(id uint) (*output.CustomerPaymentOutput, error) {
	payment, err := s.customerPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("customer payment not found")
	}
	return output.ConvertCustomerPaymentToOutput(payment), nil
}

func (s *customerPaymentService) GetAllCustomerPayments(limit, offset int) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err := s.customerPaymentRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertCustomerPaymentsToOutput(payments), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsBySalesOrder(soID string, limit, offset int) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err := s.customerPaymentRepo.FindBySalesOrderID(soID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertCustomerPaymentsToOutput(payments), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsByCustomer(customerID uint, limit, offset int) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err := s.customerPaymentRepo.FindByCustomerID(customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertCustomerPaymentsToOutput(payments), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsByStatus(status string, limit, offset int) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err := s.customerPaymentRepo.FindByPaymentStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return output.ConvertCustomerPaymentsToOutput(payments), total, nil
}

func (s *customerPaymentService) UpdateCustomerPayment(
	id uint,
	paymentInput *input.UpdateCustomerPaymentInput,
	userID string,
) (*output.CustomerPaymentOutput, error) {
	payment, err := s.customerPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("customer payment not found")
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

	updatedPayment, err := s.customerPaymentRepo.Update(id, payment)
	if err != nil {
		return nil, err
	}

	return output.ConvertCustomerPaymentToOutput(updatedPayment), nil
}

func (s *customerPaymentService) RecordPayment(
	id uint,
	paymentInput *input.RecordCustomerPaymentInput,
	userID, userName, companyName string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	payment, err := s.customerPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("customer payment not found")
	}

	// Fetch SO to verify it exists
	so, err := s.soRepo.FindByID(payment.SalesOrderID)
	if err != nil {
		return nil, errors.New("sales order not found")
	}

	// Update payment record details
	payment.PaymentMode = domain.PaymentMode(paymentInput.PaymentMode)
	payment.ReferenceNumber = paymentInput.ReferenceNumber
	payment.Notes = paymentInput.Notes

	updatedPayment, err := s.customerPaymentRepo.Update(id, payment)
	if err != nil {
		return nil, err
	}

	// If payment is being recorded, update SO payment tracking
	if updatedPayment.PaymentStatus == domain.PaymentStatusPartial || updatedPayment.PaymentStatus == domain.PaymentStatusCompleted {
		// Get all payments for this SO
		allPayments, _, err := s.customerPaymentRepo.FindBySalesOrderID(payment.SalesOrderID, 1000, 0)
		if err == nil && len(allPayments) > 0 {
			// Calculate total received across all payments
			totalReceivedAmount := 0.0
			for _, p := range allPayments {
				totalReceivedAmount += p.ReceivedAmount
			}

			// Update SO's payment tracking
			so.PaidAmount = totalReceivedAmount
			so.RemainingAmount = so.Total - totalReceivedAmount
			if so.RemainingAmount < 0 {
				so.RemainingAmount = 0
			}

			// Automatically update SO status to "paid" when fully paid
			if so.RemainingAmount <= 0 {
				so.Status = domain.SalesOrderStatusPaid
				so.RemainingAmount = 0
			}

			// Update the SO
			_, err = s.soRepo.Update(so.ID, so)
			if err != nil {
				fmt.Printf("warning: failed to update SO payment status: %v\n", err)
			}
		}
	}

	return output.ConvertCustomerPaymentToOutput(updatedPayment), nil
}

func (s *customerPaymentService) DeleteCustomerPayment(id uint) error {
	return s.customerPaymentRepo.Delete(id)
}
