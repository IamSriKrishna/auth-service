package repo

import (
	"errors"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type manufacturerRepository struct {
	db *gorm.DB
}

func NewManufacturerRepository(db *gorm.DB) ManufacturerRepository {
	return &manufacturerRepository{db: db}
}

func (r *manufacturerRepository) Create(manufacturer *models.Manufacturer) error {
	return r.db.Create(manufacturer).Error
}

func (r *manufacturerRepository) FindByID(id uint) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	err := r.db.First(&manufacturer, id).Error
	if err != nil {
		return nil, err
	}
	return &manufacturer, nil
}

// FindByStringID finds manufacturer by string ID
func (r *manufacturerRepository) FindByStringID(id string) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	err := r.db.
		Preload("ProductGroup.Components.Product").
		Preload("Company").
		Preload("User").
		Where("id = ?", id).
		First(&manufacturer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("manufacturer not found")
		}
		return nil, err
	}
	return &manufacturer, nil
}

// FindByProductGroupID finds all manufacturers for a specific product group
func (r *manufacturerRepository) FindByProductGroupID(productGroupID string) ([]models.Manufacturer, error) {
	var manufacturers []models.Manufacturer
	err := r.db.
		Where("product_group_id = ?", productGroupID).
		Preload("ProductGroup.Components.Product").
		Find(&manufacturers).Error
	return manufacturers, err
}

func (r *manufacturerRepository) FindAll(limit, offset int) ([]models.Manufacturer, int64, error) {
	var manufacturers []models.Manufacturer
	var count int64
	err := r.db.Model(&models.Manufacturer{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.
		Preload("ProductGroup.Components.Product").
		Preload("Company").
		Preload("User").
		Limit(limit).
		Offset(offset).
		Find(&manufacturers).Error
	if err != nil {
		return nil, 0, err
	}
	return manufacturers, count, nil
}

// FindAllWithFilter finds all manufacturers with optional filters for company and product group
func (r *manufacturerRepository) FindAllWithFilter(limit, offset int, companyID *uint, productGroupID *string) ([]models.Manufacturer, int64, error) {
	var manufacturers []models.Manufacturer
	var count int64
	query := r.db.Model(&models.Manufacturer{})

	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	}
	if productGroupID != nil {
		query = query.Where("product_group_id = ?", *productGroupID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Preload("ProductGroup.Components.Product").
		Preload("Company").
		Preload("User").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&manufacturers).Error
	if err != nil {
		return nil, 0, err
	}
	return manufacturers, count, nil
}

func (r *manufacturerRepository) Update(manufacturer *models.Manufacturer) error {
	return r.db.Save(manufacturer).Error
}

func (r *manufacturerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Manufacturer{}, id).Error
}

// DeleteByStringID deletes a manufacturer by string ID
func (r *manufacturerRepository) DeleteByStringID(id string) error {
	return r.db.Delete(&models.Manufacturer{}, "id = ?", id).Error
}
