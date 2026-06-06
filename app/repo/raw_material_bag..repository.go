package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)


type rawMaterialBagRepository struct {
	db *gorm.DB
}

func NewRawMaterialBagRepository(db *gorm.DB) RawMaterialBagRepository {
	return &rawMaterialBagRepository{db: db}
}

func (r *rawMaterialBagRepository) CreateMany(bags []models.RawMaterialBag) error {
	return r.db.Create(&bags).Error
}

func (r *rawMaterialBagRepository) GetByID(id string) (*models.RawMaterialBag, error) {
	var bag models.RawMaterialBag
	err := r.db.Where("id = ?", id).First(&bag).Error
	if err != nil {
		return nil, err
	}
	return &bag, nil
}
func (r *rawMaterialBagRepository) GetAll(limit, offset int) ([]models.RawMaterialBag, int64, error) {
	var bags []models.RawMaterialBag
	var total int64

	err := r.db.Model(&models.RawMaterialBag{}).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&bags).Error

	return bags, total, err
}
func (r *rawMaterialBagRepository) GetByProductID(productID string) ([]models.RawMaterialBag, error) {
	var bags []models.RawMaterialBag
	err := r.db.
		Where("product_id = ? AND remaining_kg > 0", productID).
		Order("bag_number ASC").
		Find(&bags).Error
	return bags, err
}

func (r *rawMaterialBagRepository) GetByPurchaseOrderID(poID string) ([]models.RawMaterialBag, error) {
	var bags []models.RawMaterialBag
	err := r.db.
		Where("purchase_order_id = ?", poID).
		Order("bag_number ASC").
		Find(&bags).Error
	return bags, err
}

func (r *rawMaterialBagRepository) Update(bag *models.RawMaterialBag) error {
	return r.db.Save(bag).Error
}