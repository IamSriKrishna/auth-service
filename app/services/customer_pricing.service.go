package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type CustomerPricingService interface {
	CreatePricing(customerID uint, lineItems []input.CustomerPricingLineItem) ([]models.CustomerPricing, error)
	UpdatePricing(id string, rate float64, account, description string, isActive bool) (*models.CustomerPricing, error)
	DeletePricing(id string) error
	GetPricingByID(id string) (*models.CustomerPricing, error)
	GetPricingByCustomer(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetAllPricing(offset, limit int) ([]models.CustomerPricing, int64, error)
	GetActivePricingByCustomer(customerID uint) ([]models.CustomerPricing, error)
	SetEffectiveDateRange(id string, fromDate, toDate *time.Time) error

	CreatePricingForCompany(customerID uint, lineItems []input.CustomerPricingLineItem, userID, companyID uint) ([]models.CustomerPricing, error)
	UpdatePricingForCompany(id string, rate float64, account, description string, isActive bool, userID, companyID uint) (*models.CustomerPricing, error)
	DeletePricingForCompany(id string, userID, companyID uint) error
	GetPricingByIDForCompany(id string, companyID uint) (*models.CustomerPricing, error)
	GetPricingByCustomerForCompany(customerID, companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetAllPricingForCompany(companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetActivePricingByCustomerForCompany(customerID, companyID uint) ([]models.CustomerPricing, error)
	SetEffectiveDateRangeForCompany(id string, fromDate, toDate *time.Time, userID, companyID uint) error
}

type customerPricingService struct {
	customerPricingRepo repo.CustomerPricingRepository
	customerRepo        repo.CustomerRepository
	productRepo         repo.ProductRepository
	userRepo            repo.UserRepository
}

func NewCustomerPricingService(customerPricingRepo repo.CustomerPricingRepository) CustomerPricingService {
	return &customerPricingService{customerPricingRepo: customerPricingRepo}
}

func NewCustomerPricingServiceWithDependencies(
	customerPricingRepo repo.CustomerPricingRepository,
	customerRepo repo.CustomerRepository,
	productRepo repo.ProductRepository,
	userRepo repo.UserRepository,
) CustomerPricingService {
	return &customerPricingService{
		customerPricingRepo: customerPricingRepo,
		customerRepo:        customerRepo,
		productRepo:         productRepo,
		userRepo:            userRepo,
	}
}

func (s *customerPricingService) validateUserCompany(userID, companyID uint) error {
	if userID == 0 {
		return errors.New("invalid authenticated user")
	}
	if companyID == 0 {
		return errors.New("invalid company")
	}
	if s.userRepo == nil {
		return errors.New("user repository is required")
	}

	user, err := s.userRepo.GetByIDAndCompanyID(userID, companyID)
	if err != nil || user == nil {
		return errors.New("user does not belong to the company")
	}
	return nil
}

func (s *customerPricingService) validateCustomerCompany(customerID, companyID uint) error {
	if customerID == 0 {
		return errors.New("customer ID is required")
	}
	if s.customerRepo == nil {
		return errors.New("customer repository is required")
	}

	customer, err := s.customerRepo.FindByIDAndCompany(customerID, companyID)
	if err != nil || customer == nil {
		return errors.New("customer not found in your company")
	}
	return nil
}

func (s *customerPricingService) validateProductCompany(productID string, companyID uint) (*models.Product, error) {
	if productID == "" {
		return nil, errors.New("product ID is required")
	}
	if s.productRepo == nil {
		return nil, errors.New("product repository is required")
	}

	product, err := s.productRepo.FindByIDAndCompany(productID, companyID)
	if err != nil || product == nil {
		return nil, fmt.Errorf("product %s not found in your company", productID)
	}
	return product, nil
}

func (s *customerPricingService) CreatePricing(customerID uint, lineItems []input.CustomerPricingLineItem) ([]models.CustomerPricing, error) {
	if customerID == 0 {
		return nil, errors.New("customer ID is required")
	}
	if len(lineItems) == 0 {
		return nil, errors.New("at least one line item is required")
	}

	createdPricings := make([]models.CustomerPricing, 0, len(lineItems))
	for _, item := range lineItems {
		if item.Rate < 0 {
			return nil, errors.New("rate cannot be negative")
		}
		if item.ProductID == "" {
			return nil, errors.New("product ID is required for all line items")
		}

		existing, err := s.customerPricingRepo.GetByCustomerAndProduct(customerID, item.ProductID)
		if err == nil && existing != nil {
			existing.Rate = item.Rate
			existing.Account = item.Account
			existing.Description = item.Description
			existing.IsActive = true
			existing.UpdatedAt = time.Now()
			if err := s.customerPricingRepo.Update(existing); err != nil {
				return nil, fmt.Errorf("failed to update existing pricing: %w", err)
			}
			createdPricings = append(createdPricings, *existing)
			continue
		}

		pricing := &models.CustomerPricing{
			ID: uuid.New().String(), CustomerID: customerID, ProductID: item.ProductID,
			ProductName: item.ProductName, Rate: item.Rate, Account: item.Account,
			Description: item.Description, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		if err := s.customerPricingRepo.Create(pricing); err != nil {
			return nil, fmt.Errorf("failed to create customer pricing: %w", err)
		}
		createdPricings = append(createdPricings, *pricing)
	}
	return createdPricings, nil
}

func (s *customerPricingService) UpdatePricing(id string, rate float64, account, description string, isActive bool) (*models.CustomerPricing, error) {
	if id == "" {
		return nil, errors.New("pricing ID is required")
	}
	if rate < 0 {
		return nil, errors.New("rate cannot be negative")
	}

	pricing, err := s.customerPricingRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("pricing not found: %s", id)
	}

	pricing.Rate = rate
	pricing.Account = account
	pricing.Description = description
	pricing.IsActive = isActive
	pricing.UpdatedAt = time.Now()

	if err := s.customerPricingRepo.Update(pricing); err != nil {
		return nil, fmt.Errorf("failed to update customer pricing: %w", err)
	}
	return pricing, nil
}

func (s *customerPricingService) DeletePricing(id string) error {
	if id == "" {
		return errors.New("pricing ID is required")
	}
	if _, err := s.customerPricingRepo.GetByID(id); err != nil {
		return fmt.Errorf("pricing not found: %s", id)
	}
	return s.customerPricingRepo.Delete(id)
}

func (s *customerPricingService) GetPricingByID(id string) (*models.CustomerPricing, error) {
	if id == "" {
		return nil, errors.New("pricing ID is required")
	}
	return s.customerPricingRepo.GetByID(id)
}

func (s *customerPricingService) GetPricingByCustomer(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	if customerID == 0 {
		return nil, 0, errors.New("customer ID is required")
	}
	return s.customerPricingRepo.GetByCustomerID(customerID, offset, limit)
}

func (s *customerPricingService) GetAllPricing(offset, limit int) ([]models.CustomerPricing, int64, error) {
	return s.customerPricingRepo.GetAll(offset, limit)
}

func (s *customerPricingService) GetActivePricingByCustomer(customerID uint) ([]models.CustomerPricing, error) {
	if customerID == 0 {
		return nil, errors.New("customer ID is required")
	}
	return s.customerPricingRepo.GetActiveByCustomer(customerID)
}

func (s *customerPricingService) SetEffectiveDateRange(id string, fromDate, toDate *time.Time) error {
	if id == "" {
		return errors.New("pricing ID is required")
	}
	pricing, err := s.customerPricingRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("pricing not found: %s", id)
	}
	if fromDate != nil && toDate != nil && fromDate.After(*toDate) {
		return errors.New("effective from date cannot be after effective to date")
	}
	pricing.EffectiveFrom = fromDate
	pricing.EffectiveTo = toDate
	pricing.UpdatedAt = time.Now()
	return s.customerPricingRepo.Update(pricing)
}

func (s *customerPricingService) CreatePricingForCompany(customerID uint, lineItems []input.CustomerPricingLineItem, userID, companyID uint) ([]models.CustomerPricing, error) {
	if err := s.validateUserCompany(userID, companyID); err != nil {
		return nil, err
	}
	if err := s.validateCustomerCompany(customerID, companyID); err != nil {
		return nil, err
	}
	if len(lineItems) == 0 {
		return nil, errors.New("at least one line item is required")
	}

	createdPricings := make([]models.CustomerPricing, 0, len(lineItems))
	for _, item := range lineItems {
		if item.Rate < 0 {
			return nil, errors.New("rate cannot be negative")
		}

		product, err := s.validateProductCompany(item.ProductID, companyID)
		if err != nil {
			return nil, err
		}

		existing, err := s.customerPricingRepo.GetByCustomerAndProductAndCompany(customerID, item.ProductID, companyID)
		if err == nil && existing != nil {
			log.Printf("[PRICING] Updating company pricing for customer %d and product %s", customerID, item.ProductID)
			existing.Rate = item.Rate
			existing.Account = item.Account
			existing.Description = item.Description
			existing.IsActive = true
			existing.ProductName = product.Name
			existing.UpdatedBy = string(userID)
			existing.UpdatedAt = time.Now()
			if err := s.customerPricingRepo.UpdateByCompany(existing, companyID); err != nil {
				return nil, fmt.Errorf("failed to update existing pricing: %w", err)
			}
			createdPricings = append(createdPricings, *existing)
			continue
		}

		pricing := &models.CustomerPricing{
			ID: uuid.New().String(), CustomerID: customerID, ProductID: item.ProductID,
			ProductName: product.Name, Rate: item.Rate, Account: item.Account,
			Description: item.Description, IsActive: true, CompanyID: companyID,
			CreatedBy: string(userID), UpdatedBy: string(userID), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		if err := s.customerPricingRepo.CreateForCompany(pricing, companyID); err != nil {
			return nil, fmt.Errorf("failed to create customer pricing: %w", err)
		}
		createdPricings = append(createdPricings, *pricing)
	}
	return createdPricings, nil
}

func (s *customerPricingService) UpdatePricingForCompany(id string, rate float64, account, description string, isActive bool, userID, companyID uint) (*models.CustomerPricing, error) {
	if err := s.validateUserCompany(userID, companyID); err != nil {
		return nil, err
	}
	if id == "" {
		return nil, errors.New("pricing ID is required")
	}
	if rate < 0 {
		return nil, errors.New("rate cannot be negative")
	}

	pricing, err := s.customerPricingRepo.GetByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("pricing not found")
	}

	pricing.Rate = rate
	pricing.Account = account
	pricing.Description = description
	pricing.IsActive = isActive
	pricing.UpdatedBy = string(userID)
	pricing.UpdatedAt = time.Now()

	if err := s.customerPricingRepo.UpdateByCompany(pricing, companyID); err != nil {
		return nil, fmt.Errorf("failed to update customer pricing: %w", err)
	}
	return pricing, nil
}

func (s *customerPricingService) DeletePricingForCompany(id string, userID, companyID uint) error {
	if err := s.validateUserCompany(userID, companyID); err != nil {
		return err
	}
	if id == "" {
		return errors.New("pricing ID is required")
	}
	if _, err := s.customerPricingRepo.GetByIDAndCompany(id, companyID); err != nil {
		return errors.New("pricing not found")
	}
	return s.customerPricingRepo.DeleteByCompany(id, companyID)
}

func (s *customerPricingService) GetPricingByIDForCompany(id string, companyID uint) (*models.CustomerPricing, error) {
	if id == "" {
		return nil, errors.New("pricing ID is required")
	}
	return s.customerPricingRepo.GetByIDAndCompany(id, companyID)
}

func (s *customerPricingService) GetPricingByCustomerForCompany(customerID, companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	if err := s.validateCustomerCompany(customerID, companyID); err != nil {
		return nil, 0, err
	}
	return s.customerPricingRepo.GetByCustomerIDAndCompany(customerID, companyID, offset, limit)
}

func (s *customerPricingService) GetAllPricingForCompany(companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	if companyID == 0 {
		return nil, 0, errors.New("invalid company")
	}
	return s.customerPricingRepo.GetAllByCompany(companyID, offset, limit)
}

func (s *customerPricingService) GetActivePricingByCustomerForCompany(customerID, companyID uint) ([]models.CustomerPricing, error) {
	if err := s.validateCustomerCompany(customerID, companyID); err != nil {
		return nil, err
	}
	return s.customerPricingRepo.GetActiveByCustomerAndCompany(customerID, companyID)
}

func (s *customerPricingService) SetEffectiveDateRangeForCompany(id string, fromDate, toDate *time.Time, userID, companyID uint) error {
	if err := s.validateUserCompany(userID, companyID); err != nil {
		return err
	}
	if id == "" {
		return errors.New("pricing ID is required")
	}
	if fromDate != nil && toDate != nil && fromDate.After(*toDate) {
		return errors.New("effective from date cannot be after effective to date")
	}

	pricing, err := s.customerPricingRepo.GetByIDAndCompany(id, companyID)
	if err != nil {
		return errors.New("pricing not found")
	}

	pricing.EffectiveFrom = fromDate
	pricing.EffectiveTo = toDate
	pricing.UpdatedBy = string(userID)
	pricing.UpdatedAt = time.Now()
	return s.customerPricingRepo.UpdateByCompany(pricing, companyID)
}
