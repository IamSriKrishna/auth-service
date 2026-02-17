package repo
import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type brandRepository struct{
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{db: db}
}

func (r *brandRepository) Create(brand *models.Brand) error {
	return r.db.Create(brand).Error
}

func (r *brandRepository) FindByID(id uint) (*models.Brand, error) {
	var brand models.Brand
	err := r.db.First(&brand, id).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *brandRepository) FindAll(limit, offset int) ([]models.Brand, int64, error) {
	var brands []models.Brand
	var count int64
	err := r.db.Model(&models.Brand{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Limit(limit).Offset(offset).Find(&brands).Error
	if err != nil {
		return nil, 0, err
	}
	return brands, count, nil
}

func (r *brandRepository) Update(brand *models.Brand) error {
	return r.db.Save(brand).Error
}

func (r *brandRepository) Delete(id uint) error {
	return r.db.Delete(&models.Brand{}, id).Error
}
