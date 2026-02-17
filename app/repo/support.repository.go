package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type supportRepository struct {
	db *gorm.DB
}

func NewSupportRepository(db *gorm.DB) SupportRepository {
	return &supportRepository{db: db}
}

func (r *supportRepository) Create(support *models.Support) error {
	return r.db.Table(support.TableName()).Create(support).Error
}

func (r *supportRepository) GetByID(id uint) (*models.Support, error) {
	var support models.Support
	err := r.db.First(&support, id).Error
	if err != nil {
		return nil, err
	}
	return &support, nil
}

func (r *supportRepository) List(offset, limit int) ([]models.Support, int64, error) {
	var supports []models.Support
	var total int64

	if err := r.db.Model(&models.Support{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&supports).Error
	if err != nil {
		return nil, 0, err
	}

	return supports, total, nil
}

func (r *supportRepository) Update(support *models.Support) error {
	return r.db.Save(support).Error
}

func (r *supportRepository) Delete(id uint) error {
	return r.db.Delete(&models.Support{}, id).Error
}
