package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type businessTypeRepository struct {
	db *gorm.DB
}

func NewBusinessTypeRepository(db *gorm.DB) BusinessTypeRepository {
	return &businessTypeRepository{db: db}
}

func (r *businessTypeRepository) FindAll() ([]models.BusinessType, error) {
	var businessTypes []models.BusinessType
	err := r.db.Where("is_active = ?", true).Order("type_name ASC").Find(&businessTypes).Error
	return businessTypes, err
}

func (r *businessTypeRepository) FindByID(id uint) (*models.BusinessType, error) {
	var businessType models.BusinessType
	err := r.db.First(&businessType, id).Error
	return &businessType, err
}

func (r *businessTypeRepository) Create(businessType *models.BusinessType) error {
	return r.db.Create(businessType).Error
}

func (r *businessTypeRepository) Update(businessType *models.BusinessType) error {
	return r.db.Save(businessType).Error
}

func (r *businessTypeRepository) Delete(id uint) error {
	return r.db.Delete(&models.BusinessType{}, id).Error
}
