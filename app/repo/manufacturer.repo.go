package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type manufacturerRepository struct {
	db *gorm.DB
}

func NewManufacturerRepository(
	db *gorm.DB,
) ManufacturerRepository {
	return &manufacturerRepository{db: db}
}

func (r *manufacturerRepository) manufacturerPreloads(
	db *gorm.DB,
) *gorm.DB {
	return db.
		Preload("ProductGroup").
		Preload("ProductGroup.Components").
		Preload("ProductGroup.Components.Product")
}

func (r *manufacturerRepository) Create(
	manufacturer *models.Manufacturer,
) error {
	if manufacturer == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Create(manufacturer).Error
}

func (r *manufacturerRepository) FindByID(
	id uint,
) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer

	err := r.manufacturerPreloads(r.db).
		Where("manufacturers.id = ?", id).
		First(&manufacturer).
		Error
	if err != nil {
		return nil, err
	}

	return &manufacturer, nil
}

func (r *manufacturerRepository) FindByStringID(
	id string,
) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer

	err := r.manufacturerPreloads(r.db).
		Where("manufacturers.id = ?", id).
		First(&manufacturer).
		Error
	if err != nil {
		return nil, err
	}

	return &manufacturer, nil
}

func (r *manufacturerRepository) FindByStringIDAndCompany(
	id string,
	companyID uint,
) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer

	err := r.manufacturerPreloads(r.db).
		Where(
			"manufacturers.id = ? AND manufacturers.company_id = ?",
			id,
			companyID,
		).
		First(&manufacturer).
		Error
	if err != nil {
		return nil, err
	}

	return &manufacturer, nil
}

func (r *manufacturerRepository) findList(
	query *gorm.DB,
	limit int,
	offset int,
) ([]models.Manufacturer, int64, error) {
	var manufacturers []models.Manufacturer
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.manufacturerPreloads(query).
		Order("manufacturers.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&manufacturers).
		Error

	return manufacturers, total, err
}

func (r *manufacturerRepository) FindAll(
	limit int,
	offset int,
) ([]models.Manufacturer, int64, error) {
	return r.findList(
		r.db.Model(&models.Manufacturer{}),
		limit,
		offset,
	)
}

func (r *manufacturerRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.Manufacturer, int64, error) {
	return r.findList(
		r.db.Model(&models.Manufacturer{}).
			Where(
				"manufacturers.company_id = ?",
				companyID,
			),
		limit,
		offset,
	)
}

func (r *manufacturerRepository) FindAllWithFilter(
	limit int,
	offset int,
	companyID *uint,
	productGroupID *string,
) ([]models.Manufacturer, int64, error) {
	query := r.db.Model(&models.Manufacturer{})

	if companyID != nil && *companyID > 0 {
		query = query.Where(
			"manufacturers.company_id = ?",
			*companyID,
		)
	}

	if productGroupID != nil &&
		*productGroupID != "" {
		query = query.Where(
			"manufacturers.product_group_id = ?",
			*productGroupID,
		)
	}

	return r.findList(query, limit, offset)
}

func (r *manufacturerRepository) FindByProductGroupID(
	productGroupID string,
) ([]models.Manufacturer, error) {
	var manufacturers []models.Manufacturer

	err := r.manufacturerPreloads(r.db).
		Where(
			"manufacturers.product_group_id = ?",
			productGroupID,
		).
		Order("manufacturers.created_at DESC").
		Find(&manufacturers).
		Error

	return manufacturers, err
}

func (r *manufacturerRepository) FindByProductGroupIDAndCompany(
	productGroupID string,
	companyID uint,
) ([]models.Manufacturer, error) {
	var manufacturers []models.Manufacturer

	err := r.manufacturerPreloads(r.db).
		Where(
			"manufacturers.product_group_id = ? AND manufacturers.company_id = ?",
			productGroupID,
			companyID,
		).
		Order("manufacturers.created_at DESC").
		Find(&manufacturers).
		Error

	return manufacturers, err
}

func (r *manufacturerRepository) Update(
	manufacturer *models.Manufacturer,
) error {
	if manufacturer == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Save(manufacturer).Error
}

func (r *manufacturerRepository) UpdateByCompany(
	manufacturer *models.Manufacturer,
	companyID uint,
) error {
	if manufacturer == nil {
		return gorm.ErrInvalidData
	}

	result := r.db.Model(&models.Manufacturer{}).
		Where(
			"id = ? AND company_id = ?",
			manufacturer.ID,
			companyID,
		).
		Updates(manufacturer)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *manufacturerRepository) Delete(
	id uint,
) error {
	result := r.db.Delete(
		&models.Manufacturer{},
		id,
	)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *manufacturerRepository) DeleteByStringID(
	id string,
) error {
	result := r.db.
		Where("id = ?", id).
		Delete(&models.Manufacturer{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *manufacturerRepository) HasSalesOrderLineItems(
	id string,
) (bool, error) {
	var count int64

	result := r.db.
		Model(&models.SalesOrderLineItem{}).
		Where("manufacturer_id = ?", id).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (r *manufacturerRepository) DeleteByStringIDAndCompany(
	id string,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND company_id = ?",
			id,
			companyID,
		).
		Delete(&models.Manufacturer{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
