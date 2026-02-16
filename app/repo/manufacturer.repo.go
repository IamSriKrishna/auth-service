package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type manufacturerRepository struct{
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

func (r *manufacturerRepository) FindAll(limit, offset int) ([]models.Manufacturer, int64, error) {
	var manufacturers []models.Manufacturer
	var count int64
	err := r.db.Model(&models.Manufacturer{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Limit(limit).Offset(offset).Find(&manufacturers).Error
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