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
	CreateCustomer(input *input.CreateCustomerInput) (*output.CustomerOutput, error)
	UpdateCustomer(id uint, input *input.UpdateCustomerInput) (*output.CustomerOutput, error)
	GetCustomerByID(id uint) (*output.CustomerOutput, error)
	GetAllCustomers(page, limit int) ([]output.CustomerListOutput, int64, error)
	DeleteCustomer(customer *models.Customer) error
}

type customerService struct {
	repo repo.CustomerRepository
}

func NewCustomerService(repo repo.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(input *input.CreateCustomerInput) (*output.CustomerOutput, error) {
	// Check if mobile number already exists
	if input.Mobile != "" {
		existingCustomer, err := s.repo.FindByMobile(input.Mobile)
		if err == nil && existingCustomer != nil {
			return nil, fmt.Errorf("mobile number %s already exists with customer: %s", 
				input.Mobile, existingCustomer.DisplayName)
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

func (s *customerService) UpdateCustomer(id uint, input *input.UpdateCustomerInput) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	// Check if mobile number is being updated and if it already exists
	if input.Mobile != nil && *input.Mobile != "" && *input.Mobile != customer.Mobile {
		existingCustomer, err := s.repo.FindByMobile(*input.Mobile)
		if err == nil && existingCustomer != nil && existingCustomer.ID != id {
			return nil, fmt.Errorf("mobile number %s already exists with customer: %s", 
				*input.Mobile, existingCustomer.DisplayName)
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

func (s *customerService) GetCustomerByID(id uint) (*output.CustomerOutput, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return helper.MapCustomerToOutput(customer), nil
}

func (s *customerService) GetAllCustomers(page, limit int) ([]output.CustomerListOutput, int64, error) {
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

	var outputs []output.CustomerListOutput
	for _, customer := range customers {
		outputs = append(outputs, output.CustomerListOutput{
			ID:               customer.ID,
			DisplayName:      customer.DisplayName,
			CompanyName:      customer.CompanyName,
			EmailAddress:     customer.EmailAddress,
			WorkPhone:        customer.WorkPhone,
			Mobile:           customer.Mobile,
			CustomerLanguage: customer.CustomerLanguage,
			CreatedAt:        customer.CreatedAt,
			UpdatedAt:        customer.UpdatedAt,
		})
	}

	return outputs, total, nil
}

func (s *customerService) DeleteCustomer(customer *models.Customer) error {
	return s.repo.Delete(customer)
}