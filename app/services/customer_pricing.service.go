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

// CustomerPricingService handles customer pricing operations
type CustomerPricingService interface {
	CreatePricing(customerID uint, lineItems []input.CustomerPricingLineItem) ([]models.CustomerPricing, error)
	UpdatePricing(id string, rate float64, account, description string, isActive bool) (*models.CustomerPricing, error)
	DeletePricing(id string) error
	GetPricingByID(id string) (*models.CustomerPricing, error)
	GetPricingByCustomerAndManufacturer(customerID uint, manufacturerID string) (*models.CustomerPricing, error)
	GetPricingByCustomer(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetPricingByManufacturer(manufacturerID string, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetAllPricing(offset, limit int) ([]models.CustomerPricing, int64, error)
	GetActivePricingByCustomer(customerID uint) ([]models.CustomerPricing, error)
	GetActivePricingByManufacturer(manufacturerID string) ([]models.CustomerPricing, error)
	SetEffectiveDateRange(id string, fromDate, toDate *time.Time) error
}

type customerPricingService struct {
	customerPricingRepo repo.CustomerPricingRepository
}

// NewCustomerPricingService creates a new instance of customer pricing service
func NewCustomerPricingService(
	customerPricingRepo repo.CustomerPricingRepository,
) CustomerPricingService {
	return &customerPricingService{
		customerPricingRepo: customerPricingRepo,
	}
}

// CreatePricing creates multiple customer pricing records from line items
func (s *customerPricingService) CreatePricing(
	customerID uint,
	lineItems []input.CustomerPricingLineItem,
) ([]models.CustomerPricing, error) {
	if customerID == 0 {
		return nil, errors.New("customer ID is required")
	}

	if len(lineItems) == 0 {
		return nil, errors.New("at least one line item is required")
	}

	var createdPricings []models.CustomerPricing

	for _, item := range lineItems {
		if item.ManufacturerID == "" {
			return nil, errors.New("manufacturer ID is required for all line items")
		}

		if item.Rate < 0 {
			return nil, errors.New("rate cannot be negative")
		}

		// Check if pricing already exists for this customer-manufacturer combination
		existing, err := s.customerPricingRepo.GetByCustomerAndManufacturer(customerID, item.ManufacturerID)
		if err == nil && existing != nil {
			log.Printf("[PRICING] Pricing already exists for customer %d and manufacturer %s, updating instead", customerID, item.ManufacturerID)
			// Update instead of creating
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
			ID:               uuid.New().String(),
			CustomerID:       customerID,
			ManufacturerID:   item.ManufacturerID,
			ManufacturerName: item.ManufacturerName,
			Rate:             item.Rate,
			Account:          item.Account,
			Description:      item.Description,
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := s.customerPricingRepo.Create(pricing); err != nil {
			return nil, fmt.Errorf("failed to create customer pricing: %w", err)
		}

		createdPricings = append(createdPricings, *pricing)
	}

	return createdPricings, nil
}

// UpdatePricing updates an existing customer pricing record
func (s *customerPricingService) UpdatePricing(
	id string,
	rate float64,
	account, description string,
	isActive bool,
) (*models.CustomerPricing, error) {
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

// DeletePricing deletes a customer pricing record
func (s *customerPricingService) DeletePricing(id string) error {
	if id == "" {
		return errors.New("pricing ID is required")
	}

	if _, err := s.customerPricingRepo.GetByID(id); err != nil {
		return fmt.Errorf("pricing not found: %s", id)
	}

	return s.customerPricingRepo.Delete(id)
}

// GetPricingByID retrieves a pricing record by ID
func (s *customerPricingService) GetPricingByID(id string) (*models.CustomerPricing, error) {
	if id == "" {
		return nil, errors.New("pricing ID is required")
	}

	return s.customerPricingRepo.GetByID(id)
}

// GetPricingByCustomerAndManufacturer retrieves pricing for a specific customer-manufacturer combination
func (s *customerPricingService) GetPricingByCustomerAndManufacturer(customerID uint, manufacturerID string) (*models.CustomerPricing, error) {
	if customerID == 0 || manufacturerID == "" {
		return nil, errors.New("customer ID and manufacturer ID are required")
	}

	return s.customerPricingRepo.GetByCustomerAndManufacturer(customerID, manufacturerID)
}

// GetPricingByCustomer retrieves all pricing records for a specific customer
func (s *customerPricingService) GetPricingByCustomer(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	if customerID == 0 {
		return nil, 0, errors.New("customer ID is required")
	}

	return s.customerPricingRepo.GetByCustomerID(customerID, offset, limit)
}

// GetPricingByManufacturer retrieves all pricing records for a specific manufacturer
func (s *customerPricingService) GetPricingByManufacturer(manufacturerID string, offset, limit int) ([]models.CustomerPricing, int64, error) {
	if manufacturerID == "" {
		return nil, 0, errors.New("manufacturer ID is required")
	}

	return s.customerPricingRepo.GetByManufacturerID(manufacturerID, offset, limit)
}

// GetAllPricing retrieves all customer pricing records
func (s *customerPricingService) GetAllPricing(offset, limit int) ([]models.CustomerPricing, int64, error) {
	return s.customerPricingRepo.GetAll(offset, limit)
}

// GetActivePricingByCustomer retrieves all active pricing records for a customer
func (s *customerPricingService) GetActivePricingByCustomer(customerID uint) ([]models.CustomerPricing, error) {
	if customerID == 0 {
		return nil, errors.New("customer ID is required")
	}

	return s.customerPricingRepo.GetActiveByCustomer(customerID)
}

// GetActivePricingByManufacturer retrieves all active pricing records for a manufacturer
func (s *customerPricingService) GetActivePricingByManufacturer(manufacturerID string) ([]models.CustomerPricing, error) {
	if manufacturerID == "" {
		return nil, errors.New("manufacturer ID is required")
	}

	return s.customerPricingRepo.GetActiveByManufacturer(manufacturerID)
}

// SetEffectiveDateRange sets the effective date range for a pricing record
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
