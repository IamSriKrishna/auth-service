package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type CustomerPricingRepository interface {
	Create(pricing *models.CustomerPricing) error
	Update(pricing *models.CustomerPricing) error
	Delete(id string) error
	GetByID(id string) (*models.CustomerPricing, error)
	GetByCustomerAndProduct(customerID uint, productID string) (*models.CustomerPricing, error)
	GetByCustomerID(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetAll(offset, limit int) ([]models.CustomerPricing, int64, error)
	GetActiveByCustomer(customerID uint) ([]models.CustomerPricing, error)

	CreateForCompany(pricing *models.CustomerPricing, companyID uint) error
	UpdateByCompany(pricing *models.CustomerPricing, companyID uint) error
	DeleteByCompany(id string, companyID uint) error
	GetByIDAndCompany(id string, companyID uint) (*models.CustomerPricing, error)
	GetByCustomerAndProductAndCompany(customerID uint, productID string, companyID uint) (*models.CustomerPricing, error)
	GetByCustomerIDAndCompany(customerID uint, companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetAllByCompany(companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error)
	GetActiveByCustomerAndCompany(customerID uint, companyID uint) ([]models.CustomerPricing, error)
}

type customerPricingRepository struct {
	db *gorm.DB
}

func NewCustomerPricingRepository(db *gorm.DB) CustomerPricingRepository {
	return &customerPricingRepository{db: db}
}

func (r *customerPricingRepository) pricingPreloads(db *gorm.DB) *gorm.DB {
	// CustomerPricing stores the customer and product names directly on the record,
	// so no relation preloading is required here.
	return db
}

func (r *customerPricingRepository) Create(pricing *models.CustomerPricing) error {
	if pricing == nil {
		return gorm.ErrInvalidData
	}
	return r.db.Create(pricing).Error
}

func (r *customerPricingRepository) CreateForCompany(pricing *models.CustomerPricing, companyID uint) error {
	if pricing == nil {
		return gorm.ErrInvalidData
	}
	pricing.CompanyID = companyID
	return r.db.Create(pricing).Error
}

func (r *customerPricingRepository) Update(pricing *models.CustomerPricing) error {
	if pricing == nil {
		return gorm.ErrInvalidData
	}
	return r.db.Save(pricing).Error
}

func (r *customerPricingRepository) UpdateByCompany(pricing *models.CustomerPricing, companyID uint) error {
	if pricing == nil {
		return gorm.ErrInvalidData
	}

	result := r.db.Model(&models.CustomerPricing{}).
		Where("id = ? AND company_id = ?", pricing.ID, companyID).
		Updates(pricing)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *customerPricingRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.CustomerPricing{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *customerPricingRepository) DeleteByCompany(id string, companyID uint) error {
	result := r.db.Where("id = ? AND company_id = ?", id, companyID).Delete(&models.CustomerPricing{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *customerPricingRepository) GetByID(id string) (*models.CustomerPricing, error) {
	var pricing models.CustomerPricing
	err := r.pricingPreloads(r.db).
		Where("customer_pricing.id = ?", id).
		First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *customerPricingRepository) GetByIDAndCompany(id string, companyID uint) (*models.CustomerPricing, error) {
	var pricing models.CustomerPricing
	err := r.pricingPreloads(r.db).
		Where("customer_pricing.id = ? AND customer_pricing.company_id = ?", id, companyID).
		First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *customerPricingRepository) GetByCustomerAndProduct(customerID uint, productID string) (*models.CustomerPricing, error) {
	var pricing models.CustomerPricing
	err := r.pricingPreloads(r.db).
		Where("customer_id = ? AND product_id = ?", customerID, productID).
		First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *customerPricingRepository) GetByCustomerAndProductAndCompany(customerID uint, productID string, companyID uint) (*models.CustomerPricing, error) {
	var pricing models.CustomerPricing
	err := r.pricingPreloads(r.db).
		Where("customer_id = ? AND product_id = ? AND company_id = ?", customerID, productID, companyID).
		First(&pricing).Error
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *customerPricingRepository) list(query *gorm.DB, offset, limit int) ([]models.CustomerPricing, int64, error) {
	var pricings []models.CustomerPricing
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.pricingPreloads(query).
		Order("customer_pricing.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&pricings).Error

	return pricings, total, err
}

func (r *customerPricingRepository) GetByCustomerID(customerID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	return r.list(
		r.db.Model(&models.CustomerPricing{}).Where("customer_id = ?", customerID),
		offset,
		limit,
	)
}

func (r *customerPricingRepository) GetByCustomerIDAndCompany(customerID uint, companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	return r.list(
		r.db.Model(&models.CustomerPricing{}).
			Where("customer_id = ? AND company_id = ?", customerID, companyID),
		offset,
		limit,
	)
}

func (r *customerPricingRepository) GetAll(offset, limit int) ([]models.CustomerPricing, int64, error) {
	return r.list(r.db.Model(&models.CustomerPricing{}), offset, limit)
}

func (r *customerPricingRepository) GetAllByCompany(companyID uint, offset, limit int) ([]models.CustomerPricing, int64, error) {
	return r.list(
		r.db.Model(&models.CustomerPricing{}).Where("company_id = ?", companyID),
		offset,
		limit,
	)
}

func (r *customerPricingRepository) GetActiveByCustomer(customerID uint) ([]models.CustomerPricing, error) {
	var pricings []models.CustomerPricing
	err := r.pricingPreloads(r.db).
		Where("customer_id = ? AND is_active = ?", customerID, true).
		Order("customer_pricing.created_at DESC").
		Find(&pricings).Error
	return pricings, err
}

func (r *customerPricingRepository) GetActiveByCustomerAndCompany(customerID uint, companyID uint) ([]models.CustomerPricing, error) {
	var pricings []models.CustomerPricing
	err := r.pricingPreloads(r.db).
		Where("customer_id = ? AND company_id = ? AND is_active = ?", customerID, companyID, true).
		Order("customer_pricing.created_at DESC").
		Find(&pricings).Error
	return pricings, err
}
