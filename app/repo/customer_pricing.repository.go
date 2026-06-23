package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

// CustomerPricingRepository defines database operations for customer pricing
type CustomerPricingRepository interface {
	Create(pricing *models.CustomerPricing) error
	Update(pricing *models.CustomerPricing) error
	Delete(id string) error
	GetByID(id string) (*models.CustomerPricing, error)
	GetByCustomerAndProduct(customerID uint, productID string) (*models.CustomerPricing, error)
	GetByCustomerID(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetAll(offset, limit int) ([]models.CustomerPricing, int64, error)
	GetActiveByCustomer(customerID uint) ([]models.CustomerPricing, error)
}

type customerPricingRepository struct {
	db *gorm.DB
}

// NewCustomerPricingRepository creates a new instance of customer pricing repository
func NewCustomerPricingRepository(db *gorm.DB) CustomerPricingRepository {
	return &customerPricingRepository{db: db}
}

// Create inserts a new customer pricing record
func (r *customerPricingRepository) Create(pricing *models.CustomerPricing) error {
	return r.db.Create(pricing).Error
}

// Update updates an existing customer pricing record
func (r *customerPricingRepository) Update(pricing *models.CustomerPricing) error {
	return r.db.Save(pricing).Error
}

// Delete soft deletes a customer pricing record
func (r *customerPricingRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.CustomerPricing{}).Error
}

// GetByID retrieves a customer pricing record by ID
func (r *customerPricingRepository) GetByID(id string) (*models.CustomerPricing, error) {
	var pricing models.CustomerPricing
	err := r.db.
		Joins("LEFT JOIN customers ON customers.id = customer_pricing.customer_id").
		Select("customer_pricing.*, customers.display_name as customer_name").
		Where("customer_pricing.id = ?", id).
		First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

// GetByCustomerAndProduct retrieves pricing for a specific customer and product
func (r *customerPricingRepository) GetByCustomerAndProduct(customerID uint, productID string) (*models.CustomerPricing, error) {
	var pricing models.CustomerPricing
	err := r.db.
		Joins("LEFT JOIN customers ON customers.id = customer_pricing.customer_id").
		Select("customer_pricing.*, customers.display_name as customer_name").
		Where("customer_pricing.customer_id = ? AND customer_pricing.product_id = ?", customerID, productID).
		First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

// GetByCustomerID retrieves all pricing records for a specific customer
func (r *customerPricingRepository) GetByCustomerID(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	var pricings []models.CustomerPricing
	var total int64

	err := r.db.
		Joins("LEFT JOIN customers ON customers.id = customer_pricing.customer_id").
		Select("customer_pricing.*, customers.display_name as customer_name").
		Where("customer_pricing.customer_id = ?", customerID).
		Offset(offset).
		Limit(limit).
		Find(&pricings).
		Offset(-1).
		Limit(-1).
		Count(&total).Error

	return pricings, total, err
}

// GetAll retrieves all customer pricing records with pagination
func (r *customerPricingRepository) GetAll(offset, limit int) ([]models.CustomerPricing, int64, error) {
	var pricings []models.CustomerPricing
	var total int64

	err := r.db.
		Joins("LEFT JOIN customers ON customers.id = customer_pricing.customer_id").
		Select("customer_pricing.*, customers.display_name as customer_name").
		Offset(offset).
		Limit(limit).
		Find(&pricings).
		Offset(-1).
		Limit(-1).
		Count(&total).Error

	return pricings, total, err
}

// GetActiveByCustomer retrieves all active pricing records for a specific customer
func (r *customerPricingRepository) GetActiveByCustomer(customerID uint) ([]models.CustomerPricing, error) {
	var pricings []models.CustomerPricing
	err := r.db.
		Joins("LEFT JOIN customers ON customers.id = customer_pricing.customer_id").
		Select("customer_pricing.*, customers.display_name as customer_name").
		Where("customer_pricing.customer_id = ? AND customer_pricing.is_active = ?", customerID, true).
		Find(&pricings).Error
	return pricings, err
}
