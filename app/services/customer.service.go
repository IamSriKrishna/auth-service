package services

import (
	"errors"
	"fmt"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/helper"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"gorm.io/gorm"
)

type CustomerService interface {
	// Existing methods retained for compatibility.
	CreateCustomer(input *input.CreateCustomerInput) (*output.CustomerOutput, error)
	UpdateCustomer(id uint, input *input.UpdateCustomerInput) (*output.CustomerOutput, error)
	GetCustomerByID(id uint) (*output.CustomerOutput, error)
	GetAllCustomers(page, limit int) ([]output.CustomerListOutput, int64, error)
	DeleteCustomer(customer *models.Customer) error

	CreateCustomerForUser(
		userID uint,
		companyID uint,
		input *input.CreateCustomerInput,
	) (*output.CustomerOutput, error)

	UpdateCustomerForUser(
		id uint,
		userID uint,
		input *input.UpdateCustomerInput,
	) (*output.CustomerOutput, error)

	GetCustomerByIDAndUser(
		id uint,
		userID uint,
	) (*output.CustomerOutput, error)

	GetCustomersByUser(
		userID uint,
		companyID uint,
		page int,
		limit int,
	) ([]output.CustomerListOutput, int64, error)

	DeleteCustomerForUser(
		id uint,
		userID uint,
	) error

	// Company-scoped methods added.
	GetCustomerByIDAndCompany(
		id uint,
		companyID uint,
	) (*output.CustomerOutput, error)

	GetCustomersByCompany(
		companyID uint,
		page int,
		limit int,
	) ([]output.CustomerListOutput, int64, error)

	UpdateCustomerForCompany(
		id uint,
		companyID uint,
		input *input.UpdateCustomerInput,
	) (*output.CustomerOutput, error)

	DeleteCustomerForCompany(
		id uint,
		companyID uint,
	) error
}

type customerService struct {
	repo repo.CustomerRepository
}

func NewCustomerService(repo repo.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func customerListOutputs(
	customers []models.Customer,
) []output.CustomerListOutput {
	outputs := make([]output.CustomerListOutput, 0, len(customers))

	for _, customer := range customers {
		outputs = append(outputs, output.CustomerListOutput{
			ID:               customer.ID,
			DisplayName:      customer.DisplayName,
			EmailAddress:     customer.EmailAddress,
			WorkPhone:        customer.WorkPhone,
			Mobile:           customer.Mobile,
			CustomerLanguage: customer.CustomerLanguage,
			CreatedAt:        customer.CreatedAt,
			UpdatedAt:        customer.UpdatedAt,
		})
	}

	return outputs
}

func (s *customerService) CreateCustomer(
	input *input.CreateCustomerInput,
) (*output.CustomerOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if input.Mobile != "" {
		existingCustomer, err := s.repo.FindByMobile(input.Mobile)
		if err == nil && existingCustomer != nil {
			return nil, fmt.Errorf(
				"mobile number %s already exists with customer: %s",
				input.Mobile,
				existingCustomer.DisplayName,
			)
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	customer := helper.MapCreateCustomerInput(input)

	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}

	return s.GetCustomerByID(customer.ID)
}

func (s *customerService) UpdateCustomer(
	id uint,
	input *input.UpdateCustomerInput,
) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	if input.Mobile != nil &&
		*input.Mobile != "" &&
		*input.Mobile != customer.Mobile {
		existingCustomer, err := s.repo.FindByMobile(*input.Mobile)
		if err == nil &&
			existingCustomer != nil &&
			existingCustomer.ID != id {
			return nil, fmt.Errorf(
				"mobile number %s already exists with customer: %s",
				*input.Mobile,
				existingCustomer.DisplayName,
			)
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	helper.ApplyUpdateCustomerInput(customer, input)

	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}

	return s.GetCustomerByID(customer.ID)
}

func (s *customerService) GetCustomerByID(
	id uint,
) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return helper.MapCustomerToOutput(customer), nil
}

func (s *customerService) GetAllCustomers(
	page int,
	limit int,
) ([]output.CustomerListOutput, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	customers, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return customerListOutputs(customers), total, nil
}

func (s *customerService) DeleteCustomer(
	customer *models.Customer,
) error {
	return s.repo.Delete(customer)
}

func (s *customerService) CreateCustomerForUser(
	userID uint,
	companyID uint,
	input *input.CreateCustomerInput,
) (*output.CustomerOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if userID == 0 {
		return nil, errors.New("invalid user")
	}

	if companyID == 0 {
		return nil, errors.New("invalid company")
	}

	if input.Mobile != "" {
		existingCustomer, err := s.repo.FindByMobileAndCompany(
			input.Mobile,
			companyID,
		)

		if err == nil && existingCustomer != nil {
			return nil, fmt.Errorf(
				"mobile number %s already exists with customer: %s",
				input.Mobile,
				existingCustomer.DisplayName,
			)
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	customer := helper.MapCreateCustomerInput(input)
	customer.UserID = userID
	customer.CompanyID = companyID

	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}

	return s.GetCustomerByIDAndCompany(
		customer.ID,
		companyID,
	)
}

// Existing compatibility method retained.
func (s *customerService) GetCustomerByIDAndUser(
	id uint,
	userID uint,
) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	return helper.MapCustomerToOutput(customer), nil
}

// Existing compatibility method retained.
func (s *customerService) UpdateCustomerForUser(
	id uint,
	userID uint,
	input *input.UpdateCustomerInput,
) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	helper.ApplyUpdateCustomerInput(customer, input)

	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}

	return helper.MapCustomerToOutput(customer), nil
}

// Existing compatibility method retained.
// It now returns company-level data because company sharing is required.
func (s *customerService) GetCustomersByUser(
	userID uint,
	companyID uint,
	page int,
	limit int,
) ([]output.CustomerListOutput, int64, error) {
	_ = userID

	return s.GetCustomersByCompany(
		companyID,
		page,
		limit,
	)
}

// Existing compatibility method retained.
func (s *customerService) DeleteCustomerForUser(
	id uint,
	userID uint,
) error {
	customer, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return fmt.Errorf("customer not found")
	}

	return s.repo.Delete(customer)
}

func (s *customerService) GetCustomerByIDAndCompany(
	id uint,
	companyID uint,
) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	return helper.MapCustomerToOutput(customer), nil
}

func (s *customerService) GetCustomersByCompany(
	companyID uint,
	page int,
	limit int,
) ([]output.CustomerListOutput, int64, error) {
	if companyID == 0 {
		return nil, 0, errors.New("invalid company")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	customers, total, err := s.repo.FindByCompanyID(
		companyID,
		page,
		limit,
	)
	if err != nil {
		return nil, 0, err
	}

	return customerListOutputs(customers), total, nil
}

func (s *customerService) UpdateCustomerForCompany(
	id uint,
	companyID uint,
	input *input.UpdateCustomerInput,
) (*output.CustomerOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	customer, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	if input.Mobile != nil &&
		*input.Mobile != "" &&
		*input.Mobile != customer.Mobile {
		existingCustomer, err := s.repo.FindByMobileAndCompany(
			*input.Mobile,
			companyID,
		)

		if err == nil &&
			existingCustomer != nil &&
			existingCustomer.ID != id {
			return nil, fmt.Errorf(
				"mobile number %s already exists with customer: %s",
				*input.Mobile,
				existingCustomer.DisplayName,
			)
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	helper.ApplyUpdateCustomerInput(customer, input)

	if err := s.repo.UpdateByCompanyID(
		customer,
		companyID,
	); err != nil {
		return nil, err
	}

	return s.GetCustomerByIDAndCompany(
		customer.ID,
		companyID,
	)
}

func (s *customerService) DeleteCustomerForCompany(
	id uint,
	companyID uint,
) error {
	return s.repo.DeleteByIDAndCompany(
		id,
		companyID,
	)
}