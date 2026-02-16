package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type taxTypeRepository struct {
	db *gorm.DB
}

func NewTaxTypeRepository(db *gorm.DB) TaxTypeRepository {
	return &taxTypeRepository{db: db}
}

func (r *taxTypeRepository) FindAll() ([]models.TaxType, error) {
	var taxTypes []models.TaxType
	err := r.db.Where("is_active = ?", true).Order("tax_name ASC").Find(&taxTypes).Error
	return taxTypes, err
}

func (r *taxTypeRepository) FindByID(id uint) (*models.TaxType, error) {
	var taxType models.TaxType
	err := r.db.First(&taxType, id).Error
	return &taxType, err
}

func (r *taxTypeRepository) Create(taxType *models.TaxType) error {
	return r.db.Create(taxType).Error
}

func (r *taxTypeRepository) Update(taxType *models.TaxType) error {
	return r.db.Save(taxType).Error
}

func (r *taxTypeRepository) Delete(id uint) error {
	return r.db.Delete(&models.TaxType{}, id).Error
}
