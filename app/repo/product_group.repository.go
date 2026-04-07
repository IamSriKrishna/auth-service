package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productGroupRepository struct {
	db *gorm.DB
}

func NewProductGroupRepository(db *gorm.DB) ProductGroupRepository {
	return &productGroupRepository{db: db}
}

func (r *productGroupRepository) Create(productGroup *models.ProductGroup) error {
	return r.db.Create(productGroup).Error
}

func (r *productGroupRepository) FindByID(id string) (*models.ProductGroup, error) {
	var productGroup models.ProductGroup
	err := r.db.Preload("Components.Product.ProductDetails").
		Preload("Components.Product.SalesInfo").
		Preload("Components.Product.PurchaseInfo").
		First(&productGroup, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &productGroup, nil
}

func (r *productGroupRepository) FindAll(limit, offset int, search string) ([]models.ProductGroup, int64, error) {
	var productGroups []models.ProductGroup
	var count int64

	query := r.db.Model(&models.ProductGroup{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Components.Product.ProductDetails").
		Preload("Components.Product.SalesInfo").
		Preload("Components.Product.PurchaseInfo").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&productGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return productGroups, count, nil
}

func (r *productGroupRepository) Update(productGroup *models.ProductGroup) error {
	return r.db.Save(productGroup).Error
}

func (r *productGroupRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductGroup{}, "id = ?", id).Error
}

func (r *productGroupRepository) FindByName(name string) (*models.ProductGroup, error) {
	var productGroup models.ProductGroup
	err := r.db.Preload("Components.Product.ProductDetails").
		Preload("Components.Product.SalesInfo").
		Preload("Components.Product.PurchaseInfo").
		First(&productGroup, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &productGroup, nil
}

func (r *productGroupRepository) FindActiveGroups(limit, offset int) ([]models.ProductGroup, int64, error) {
	var productGroups []models.ProductGroup
	var count int64

	query := r.db.Where("is_active = ?", true)

	err := query.Model(&models.ProductGroup{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Components.Product.ProductDetails").
		Preload("Components.Product.SalesInfo").
		Preload("Components.Product.PurchaseInfo").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&productGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return productGroups, count, nil
}
