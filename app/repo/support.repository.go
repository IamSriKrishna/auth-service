package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

// supportRepository implements SupportRepository interface
type supportRepository struct {
	db *gorm.DB // Database with dbresolver (handles read/write splitting automatically)
}

// NewSupportRepository creates a new support repository instance
func NewSupportRepository(db *gorm.DB) SupportRepository {
	return &supportRepository{db: db}
}

// Create creates a new support ticket
func (r *supportRepository) Create(support *models.Support) error {
	return r.db.Table(support.TableName()).Create(support).Error
}

// GetByID retrieves a support ticket by ID
func (r *supportRepository) GetByID(id uint) (*models.Support, error) {
	var support models.Support
	err := r.db.First(&support, id).Error
	if err != nil {
		return nil, err
	}
	return &support, nil
}

// List retrieves paginated support tickets
func (r *supportRepository) List(offset, limit int) ([]models.Support, int64, error) {
	var supports []models.Support
	var total int64

	// Count total records
	if err := r.db.Model(&models.Support{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&supports).Error
	if err != nil {
		return nil, 0, err
	}

	return supports, total, nil
}

// Update updates a support ticket
func (r *supportRepository) Update(support *models.Support) error {
	return r.db.Save(support).Error
}

// Delete soft deletes a support ticket
func (r *supportRepository) Delete(id uint) error {
	return r.db.Delete(&models.Support{}, id).Error
}
