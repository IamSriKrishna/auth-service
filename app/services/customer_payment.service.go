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
	// Existing methods retained.
	CreateCustomerPayment(
		paymentInput *input.CreateCustomerPaymentInput,
		userID string,
		userName string,
		companyName string,
		companyID uint,
	) (*output.CustomerPaymentOutput, error)

	GetCustomerPayment(
		id uint,
	) (*output.CustomerPaymentOutput, error)

	GetAllCustomerPayments(
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	GetCustomerPaymentsBySalesOrder(
		salesOrderID string,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	GetCustomerPaymentsByCustomer(
		customerID uint,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	GetCustomerPaymentsByStatus(
		status string,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	UpdateCustomerPayment(
		id uint,
		paymentInput *input.UpdateCustomerPaymentInput,
		userID string,
	) (*output.CustomerPaymentOutput, error)

	RecordPayment(
		id uint,
		paymentInput *input.RecordCustomerPaymentInput,
		userID string,
		userName string,
		companyName string,
		companyID uint,
	) (*output.CustomerPaymentOutput, error)

	DeleteCustomerPayment(
		id uint,
	) error

	// Company-scoped methods used by protected routes.
	CreateCustomerPaymentForCompany(
		paymentInput *input.CreateCustomerPaymentInput,
		userID string,
		userName string,
		companyName string,
		companyID uint,
	) (*output.CustomerPaymentOutput, error)

	GetCustomerPaymentByCompany(
		id uint,
		companyID uint,
	) (*output.CustomerPaymentOutput, error)

	GetAllCustomerPaymentsByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	GetCustomerPaymentsBySalesOrderAndCompany(
		salesOrderID string,
		companyID uint,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	GetCustomerPaymentsByCustomerAndCompany(
		customerID uint,
		companyID uint,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	GetCustomerPaymentsByStatusAndCompany(
		status string,
		companyID uint,
		limit int,
		offset int,
	) ([]output.CustomerPaymentOutput, int64, error)

	UpdateCustomerPaymentForCompany(
		id uint,
		paymentInput *input.UpdateCustomerPaymentInput,
		userID string,
		companyID uint,
	) (*output.CustomerPaymentOutput, error)

	RecordPaymentForCompany(
		id uint,
		paymentInput *input.RecordCustomerPaymentInput,
		userID string,
		userName string,
		companyName string,
		companyID uint,
	) (*output.CustomerPaymentOutput, error)

	DeleteCustomerPaymentForCompany(
		id uint,
		companyID uint,
	) error
}

type customerPaymentService struct {
	customerPaymentRepo repo.CustomerPaymentRepository
	soRepo              repo.SalesOrderRepository
	customerRepo        repo.CustomerRepository
	userRepo            repo.UserRepository
}

func NewCustomerPaymentService(
	customerPaymentRepo repo.CustomerPaymentRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	userRepo repo.UserRepository,
) CustomerPaymentService {
	return &customerPaymentService{
		customerPaymentRepo: customerPaymentRepo,
		soRepo:              soRepo,
		customerRepo:        customerRepo,
		userRepo:            userRepo,
	}
}

func (s *customerPaymentService) CreateCustomerPayment(
	paymentInput *input.CreateCustomerPaymentInput,
	userID string,
	userName string,
	companyName string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	salesOrder, err := s.soRepo.FindByID(
		paymentInput.SalesOrderID,
	)
	if err != nil {
		return nil, errors.New("sales order not found")
	}

	customer, err := s.customerRepo.FindByID(
		paymentInput.CustomerID,
	)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	if customer.ID != salesOrder.CustomerID {
		return nil, errors.New(
			"customer does not match sales order customer",
		)
	}

	existingPayments, _, err :=
		s.customerPaymentRepo.FindBySalesOrderID(
			paymentInput.SalesOrderID,
			1000,
			0,
		)
	if err != nil {
		return nil, err
	}

	if paymentInput.Amount > salesOrder.Total {
		return nil, fmt.Errorf(
			"payment amount exceeds SO total. Total SO: %.2f, Attempting to add: %.2f",
			salesOrder.Total,
			paymentInput.Amount,
		)
	}

	paymentNumber :=
		"CP-" +
			time.Now().Format("20060102") +
			"-" +
			uuid.New().String()[:8]

	totalReceivedAmount := 0.0
	for _, payment := range existingPayments {
		totalReceivedAmount += payment.ReceivedAmount
	}
	totalReceivedAmount += paymentInput.Amount

	newRemainingAmount :=
		salesOrder.Total - totalReceivedAmount
	if newRemainingAmount < 0 {
		newRemainingAmount = 0
	}

	paymentStatus := domain.PaymentStatusPending
	if newRemainingAmount <= 0 {
		paymentStatus = domain.PaymentStatusCompleted
	} else if totalReceivedAmount > 0 {
		paymentStatus = domain.PaymentStatusPartial
	}

	customerPayment := models.CustomerPayment{
		PaymentNumber:        paymentNumber,
		SalesOrderID:         paymentInput.SalesOrderID,
		CustomerID:           paymentInput.CustomerID,
		PaymentMode:          domain.PaymentMode(paymentInput.PaymentMode),
		Amount:               paymentInput.Amount,
		ReceivedAmount:       paymentInput.Amount,
		RemainingAmount:      newRemainingAmount,
		PaymentStatus:        paymentStatus,
		PaymentDate:          paymentInput.PaymentDate,
		ReferenceNumber:      paymentInput.ReferenceNumber,
		Notes:                paymentInput.Notes,
		CreatedByUserID:      userID,
		CreatedByUserName:    userName,
		CreatedByCompanyID:   companyID,
		CreatedByCompanyName: companyName,
	}

	createdPayment, err :=
		s.customerPaymentRepo.Create(&customerPayment)
	if err != nil {
		return nil, err
	}

	salesOrder.PaidAmount = totalReceivedAmount
	salesOrder.RemainingAmount = newRemainingAmount

	if newRemainingAmount <= 0 {
		salesOrder.Status = domain.SalesOrderStatusPaid
		salesOrder.RemainingAmount = 0
	}

	if _, err := s.soRepo.Update(
		salesOrder.ID,
		salesOrder,
	); err != nil {
		return nil, err
	}

	return output.ConvertCustomerPaymentToOutput(
		createdPayment,
	), nil
}

func (s *customerPaymentService) GetCustomerPayment(
	id uint,
) (*output.CustomerPaymentOutput, error) {
	payment, err := s.customerPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(
			"customer payment not found",
		)
	}

	return output.ConvertCustomerPaymentToOutput(
		payment,
	), nil
}

func (s *customerPaymentService) GetAllCustomerPayments(
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err :=
		s.customerPaymentRepo.FindAll(
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsBySalesOrder(
	salesOrderID string,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err :=
		s.customerPaymentRepo.FindBySalesOrderID(
			salesOrderID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsByCustomer(
	customerID uint,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err :=
		s.customerPaymentRepo.FindByCustomerID(
			customerID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsByStatus(
	status string,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err :=
		s.customerPaymentRepo.FindByPaymentStatus(
			status,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) UpdateCustomerPayment(
	id uint,
	paymentInput *input.UpdateCustomerPaymentInput,
	userID string,
) (*output.CustomerPaymentOutput, error) {
	payment, err := s.customerPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(
			"customer payment not found",
		)
	}

	if paymentInput.PaymentMode != nil {
		payment.PaymentMode =
			domain.PaymentMode(*paymentInput.PaymentMode)
	}
	if paymentInput.Amount != nil {
		payment.Amount = *paymentInput.Amount
	}
	if paymentInput.PaymentDate != nil {
		payment.PaymentDate = *paymentInput.PaymentDate
	}
	if paymentInput.ReferenceNumber != nil {
		payment.ReferenceNumber =
			*paymentInput.ReferenceNumber
	}
	if paymentInput.Notes != nil {
		payment.Notes = *paymentInput.Notes
	}

	updatedPayment, err :=
		s.customerPaymentRepo.Update(
			id,
			payment,
		)
	if err != nil {
		return nil, err
	}

	return output.ConvertCustomerPaymentToOutput(
		updatedPayment,
	), nil
}

func (s *customerPaymentService) RecordPayment(
	id uint,
	paymentInput *input.RecordCustomerPaymentInput,
	userID string,
	userName string,
	companyName string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	payment, err := s.customerPaymentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(
			"customer payment not found",
		)
	}

	salesOrder, err := s.soRepo.FindByID(
		payment.SalesOrderID,
	)
	if err != nil {
		return nil, errors.New("sales order not found")
	}

	payment.PaymentMode =
		domain.PaymentMode(paymentInput.PaymentMode)
	payment.ReferenceNumber =
		paymentInput.ReferenceNumber
	payment.Notes = paymentInput.Notes

	updatedPayment, err :=
		s.customerPaymentRepo.Update(
			id,
			payment,
		)
	if err != nil {
		return nil, err
	}

	if updatedPayment.PaymentStatus ==
		domain.PaymentStatusPartial ||
		updatedPayment.PaymentStatus ==
			domain.PaymentStatusCompleted {
		allPayments, _, findErr :=
			s.customerPaymentRepo.FindBySalesOrderID(
				payment.SalesOrderID,
				1000,
				0,
			)
		if findErr == nil && len(allPayments) > 0 {
			totalReceivedAmount := 0.0
			for _, currentPayment := range allPayments {
				totalReceivedAmount +=
					currentPayment.ReceivedAmount
			}

			salesOrder.PaidAmount = totalReceivedAmount
			salesOrder.RemainingAmount =
				salesOrder.Total - totalReceivedAmount
			if salesOrder.RemainingAmount < 0 {
				salesOrder.RemainingAmount = 0
			}

			if salesOrder.RemainingAmount <= 0 {
				salesOrder.Status =
					domain.SalesOrderStatusPaid
				salesOrder.RemainingAmount = 0
			}

			if _, updateErr := s.soRepo.Update(
				salesOrder.ID,
				salesOrder,
			); updateErr != nil {
				fmt.Printf(
					"warning: failed to update SO payment status: %v\n",
					updateErr,
				)
			}
		}
	}

	return output.ConvertCustomerPaymentToOutput(
		updatedPayment,
	), nil
}

func (s *customerPaymentService) DeleteCustomerPayment(
	id uint,
) error {
	return s.customerPaymentRepo.Delete(id)
}

func (s *customerPaymentService) validateCustomerPaymentUserCompany(
	userID string,
	companyID uint,
) error {
	if companyID == 0 {
		return errors.New("invalid company")
	}

	var userIDUint uint
	if _, err := fmt.Sscanf(
		userID,
		"%d",
		&userIDUint,
	); err != nil || userIDUint == 0 {
		return errors.New("invalid authenticated user")
	}

	user, err := s.userRepo.GetByIDAndCompanyID(
		userIDUint,
		companyID,
	)
	if err != nil || user == nil {
		return errors.New(
			"user does not belong to the company",
		)
	}

	return nil
}

func (s *customerPaymentService) CreateCustomerPaymentForCompany(
	paymentInput *input.CreateCustomerPaymentInput,
	userID string,
	userName string,
	companyName string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	if paymentInput == nil {
		return nil, errors.New("input cannot be nil")
	}

	if err := s.validateCustomerPaymentUserCompany(
		userID,
		companyID,
	); err != nil {
		return nil, err
	}

	salesOrder, err :=
		s.soRepo.FindByIDAndCompany(
			paymentInput.SalesOrderID,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"sales order not found in your company",
		)
	}

	customer, err :=
		s.customerRepo.FindByIDAndCompany(
			paymentInput.CustomerID,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"customer not found in your company",
		)
	}

	if customer.ID != salesOrder.CustomerID {
		return nil, errors.New(
			"customer does not match sales order customer",
		)
	}

	existingPayments, _, err :=
		s.customerPaymentRepo.FindBySalesOrderIDAndCompany(
			paymentInput.SalesOrderID,
			companyID,
			1000,
			0,
		)
	if err != nil {
		return nil, err
	}

	if paymentInput.Amount > salesOrder.Total {
		return nil, fmt.Errorf(
			"payment amount exceeds SO total. Total SO: %.2f, Attempting to add: %.2f",
			salesOrder.Total,
			paymentInput.Amount,
		)
	}

	totalReceivedAmount := 0.0
	for _, payment := range existingPayments {
		totalReceivedAmount += payment.ReceivedAmount
	}
	totalReceivedAmount += paymentInput.Amount

	if totalReceivedAmount > salesOrder.Total {
		return nil, fmt.Errorf(
			"total received amount exceeds sales order total. SO total: %.2f, total received: %.2f",
			salesOrder.Total,
			totalReceivedAmount,
		)
	}

	createdPayment, err := s.CreateCustomerPayment(
		paymentInput,
		userID,
		userName,
		companyName,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	return s.GetCustomerPaymentByCompany(
		createdPayment.ID,
		companyID,
	)
}

func (s *customerPaymentService) GetCustomerPaymentByCompany(
	id uint,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	payment, err :=
		s.customerPaymentRepo.FindByIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"customer payment not found",
		)
	}

	return output.ConvertCustomerPaymentToOutput(
		payment,
	), nil
}

func (s *customerPaymentService) GetAllCustomerPaymentsByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err :=
		s.customerPaymentRepo.FindAllByCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsBySalesOrderAndCompany(
	salesOrderID string,
	companyID uint,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	if _, err := s.soRepo.FindByIDAndCompany(
		salesOrderID,
		companyID,
	); err != nil {
		return nil, 0, errors.New(
			"sales order not found in your company",
		)
	}

	payments, total, err :=
		s.customerPaymentRepo.FindBySalesOrderIDAndCompany(
			salesOrderID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsByCustomerAndCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	if _, err := s.customerRepo.FindByIDAndCompany(
		customerID,
		companyID,
	); err != nil {
		return nil, 0, errors.New(
			"customer not found in your company",
		)
	}

	payments, total, err :=
		s.customerPaymentRepo.FindByCustomerIDAndCompany(
			customerID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) GetCustomerPaymentsByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]output.CustomerPaymentOutput, int64, error) {
	payments, total, err :=
		s.customerPaymentRepo.FindByPaymentStatusAndCompany(
			status,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	return output.ConvertCustomerPaymentsToOutput(
		payments,
	), total, nil
}

func (s *customerPaymentService) UpdateCustomerPaymentForCompany(
	id uint,
	paymentInput *input.UpdateCustomerPaymentInput,
	userID string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	if paymentInput == nil {
		return nil, errors.New("input cannot be nil")
	}

	if err := s.validateCustomerPaymentUserCompany(
		userID,
		companyID,
	); err != nil {
		return nil, err
	}

	payment, err :=
		s.customerPaymentRepo.FindByIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"customer payment not found",
		)
	}

	if paymentInput.PaymentMode != nil {
		payment.PaymentMode =
			domain.PaymentMode(*paymentInput.PaymentMode)
	}
	if paymentInput.Amount != nil {
		if *paymentInput.Amount < 0 {
			return nil, errors.New(
				"payment amount cannot be negative",
			)
		}
		payment.Amount = *paymentInput.Amount
	}
	if paymentInput.PaymentDate != nil {
		payment.PaymentDate = *paymentInput.PaymentDate
	}
	if paymentInput.ReferenceNumber != nil {
		payment.ReferenceNumber =
			*paymentInput.ReferenceNumber
	}
	if paymentInput.Notes != nil {
		payment.Notes = *paymentInput.Notes
	}

	updatedPayment, err :=
		s.customerPaymentRepo.UpdateByCompany(
			id,
			companyID,
			payment,
		)
	if err != nil {
		return nil, err
	}

	return output.ConvertCustomerPaymentToOutput(
		updatedPayment,
	), nil
}

func (s *customerPaymentService) RecordPaymentForCompany(
	id uint,
	paymentInput *input.RecordCustomerPaymentInput,
	userID string,
	userName string,
	companyName string,
	companyID uint,
) (*output.CustomerPaymentOutput, error) {
	if paymentInput == nil {
		return nil, errors.New("input cannot be nil")
	}

	if err := s.validateCustomerPaymentUserCompany(
		userID,
		companyID,
	); err != nil {
		return nil, err
	}

	payment, err :=
		s.customerPaymentRepo.FindByIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"customer payment not found",
		)
	}

	salesOrder, err :=
		s.soRepo.FindByIDAndCompany(
			payment.SalesOrderID,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"sales order not found in your company",
		)
	}

	if _, err := s.customerRepo.FindByIDAndCompany(
		payment.CustomerID,
		companyID,
	); err != nil {
		return nil, errors.New(
			"customer not found in your company",
		)
	}

	payment.PaymentMode =
		domain.PaymentMode(paymentInput.PaymentMode)
	payment.ReferenceNumber =
		paymentInput.ReferenceNumber
	payment.Notes = paymentInput.Notes
	payment.CreatedByCompanyID = companyID

	updatedPayment, err :=
		s.customerPaymentRepo.UpdateByCompany(
			id,
			companyID,
			payment,
		)
	if err != nil {
		return nil, err
	}

	if updatedPayment.PaymentStatus ==
		domain.PaymentStatusPartial ||
		updatedPayment.PaymentStatus ==
			domain.PaymentStatusCompleted {
		allPayments, _, findErr :=
			s.customerPaymentRepo.FindBySalesOrderIDAndCompany(
				payment.SalesOrderID,
				companyID,
				1000,
				0,
			)
		if findErr == nil && len(allPayments) > 0 {
			totalReceivedAmount := 0.0
			for _, currentPayment := range allPayments {
				totalReceivedAmount +=
					currentPayment.ReceivedAmount
			}

			salesOrder.PaidAmount = totalReceivedAmount
			salesOrder.RemainingAmount =
				salesOrder.Total - totalReceivedAmount
			if salesOrder.RemainingAmount < 0 {
				salesOrder.RemainingAmount = 0
			}

			if salesOrder.RemainingAmount <= 0 {
				salesOrder.Status =
					domain.SalesOrderStatusPaid
				salesOrder.RemainingAmount = 0
			}

			if _, updateErr :=
				s.soRepo.UpdateByCompany(
					salesOrder.ID,
					companyID,
					salesOrder,
				); updateErr != nil {
				return nil, fmt.Errorf(
					"failed to update sales order payment status: %w",
					updateErr,
				)
			}
		}
	}

	return output.ConvertCustomerPaymentToOutput(
		updatedPayment,
	), nil
}

func (s *customerPaymentService) DeleteCustomerPaymentForCompany(
	id uint,
	companyID uint,
) error {
	if _, err :=
		s.customerPaymentRepo.FindByIDAndCompany(
			id,
			companyID,
		); err != nil {
		return errors.New(
			"customer payment not found",
		)
	}

	return s.customerPaymentRepo.DeleteByCompany(
		id,
		companyID,
	)
}
