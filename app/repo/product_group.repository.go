package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productGroupRepository struct{ db *gorm.DB }

func NewProductGroupRepository(db *gorm.DB) ProductGroupRepository {
	return &productGroupRepository{db: db}
}

func (r *productGroupRepository) productGroupPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Components.Product.ProductDetails.ProductVariants.Attributes").
		Preload("Components.Product.SalesInfo").
		Preload("Components.Product.PurchaseInfo").
		Preload("Resources")
}

func (r *productGroupRepository) Create(productGroup *models.ProductGroup) error {
	return r.db.Create(productGroup).Error
}

func (r *productGroupRepository) CreateForCompany(productGroup *models.ProductGroup, companyID, userID uint) error {
	productGroup.CompanyID = companyID
	productGroup.CreatedBy = userID
	productGroup.UpdatedBy = userID
	return r.db.Create(productGroup).Error
}

func (r *productGroupRepository) FindByID(id string) (*models.ProductGroup, error) {
	var productGroup models.ProductGroup
	err := r.productGroupPreloads(r.db).
		Where("product_groups.id = ?", id).
		First(&productGroup).Error
	if err != nil {
		return nil, err
	}
	return &productGroup, nil
}

func (r *productGroupRepository) FindByIDAndCompany(id string, companyID uint) (*models.ProductGroup, error) {
	var productGroup models.ProductGroup
	err := r.productGroupPreloads(r.db).
		Where("product_groups.id = ? AND product_groups.company_id = ?", id, companyID).
		First(&productGroup).Error
	if err != nil {
		return nil, err
	}
	return &productGroup, nil
}

func (r *productGroupRepository) FindAll(limit, offset int, search string) ([]models.ProductGroup, int64, error) {
	query := r.db.Model(&models.ProductGroup{})
	if search != "" {
		query = query.Where("product_groups.name LIKE ?", "%"+search+"%")
	}
	return r.list(query, limit, offset)
}

func (r *productGroupRepository) FindAllByCompany(companyID uint, limit, offset int, search string) ([]models.ProductGroup, int64, error) {
	query := r.db.Model(&models.ProductGroup{}).Where("product_groups.company_id = ?", companyID)
	if search != "" {
		query = query.Where("product_groups.name LIKE ?", "%"+search+"%")
	}
	return r.list(query, limit, offset)
}

func (r *productGroupRepository) list(query *gorm.DB, limit, offset int) ([]models.ProductGroup, int64, error) {
	var productGroups []models.ProductGroup
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.productGroupPreloads(query).
		Order("product_groups.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&productGroups).Error
	if err != nil {
		return nil, 0, err
	}
	return productGroups, total, nil
}

func (r *productGroupRepository) Update(productGroup *models.ProductGroup) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ProductGroup{}).
			Where("id = ?", productGroup.ID).
			Omit("Components", "Resources").
			Updates(productGroup).Error; err != nil {
			return err
		}
		if productGroup.Components != nil {
			if err := tx.Where("product_group_id = ?", productGroup.ID).Delete(&models.ProductGroupComponent{}).Error; err != nil {
				return err
			}
			for i := range productGroup.Components {
				productGroup.Components[i].ID = 0
				productGroup.Components[i].ProductGroupID = productGroup.ID
			}
			if len(productGroup.Components) > 0 {
				if err := tx.Create(&productGroup.Components).Error; err != nil {
					return err
				}
			}
		}
		if productGroup.Resources != nil {
			if err := tx.Where("product_group_id = ?", productGroup.ID).Delete(&models.ProductGroupResource{}).Error; err != nil {
				return err
			}
			for i := range productGroup.Resources {
				productGroup.Resources[i].ID = 0
				productGroup.Resources[i].ProductGroupID = productGroup.ID
			}
			if len(productGroup.Resources) > 0 {
				if err := tx.Create(&productGroup.Resources).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *productGroupRepository) UpdateByCompany(productGroup *models.ProductGroup, companyID uint) error {
	if productGroup == nil {
		return gorm.ErrInvalidData
	}
	if productGroup.CompanyID != companyID {
		return gorm.ErrRecordNotFound
	}
	var count int64
	if err := r.db.Model(&models.ProductGroup{}).
		Where("id = ? AND company_id = ?", productGroup.ID, companyID).
		Count(&count).Error; err != nil {
		return err
	}
	// Supports a just-created legacy row whose company_id was assigned in memory.
	if count == 0 {
		if err := r.db.Model(&models.ProductGroup{}).Where("id = ?", productGroup.ID).
			Updates(map[string]interface{}{"company_id": companyID, "created_by": productGroup.CreatedBy, "updated_by": productGroup.UpdatedBy}).Error; err != nil {
			return err
		}
	}
	return r.Update(productGroup)
}

func (r *productGroupRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductGroup{}, "id = ?", id).Error
}

func (r *productGroupRepository) DeleteByCompany(id string, companyID uint) error {
	result := r.db.Where("id = ? AND company_id = ?", id, companyID).Delete(&models.ProductGroup{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *productGroupRepository) FindByName(name string) (*models.ProductGroup, error) {
	var productGroup models.ProductGroup
	err := r.productGroupPreloads(r.db).Where("product_groups.name = ?", name).First(&productGroup).Error
	if err != nil {
		return nil, err
	}
	return &productGroup, nil
}

func (r *productGroupRepository) FindByNameAndCompany(name string, companyID uint) (*models.ProductGroup, error) {
	var productGroup models.ProductGroup
	err := r.productGroupPreloads(r.db).
		Where("product_groups.name = ? AND product_groups.company_id = ?", name, companyID).
		First(&productGroup).Error
	if err != nil {
		return nil, err
	}
	return &productGroup, nil
}

func (r *productGroupRepository) FindActiveGroups(limit, offset int) ([]models.ProductGroup, int64, error) {
	return r.list(r.db.Model(&models.ProductGroup{}).Where("product_groups.is_active = ?", true), limit, offset)
}

func (r *productGroupRepository) FindActiveGroupsByCompany(companyID uint, limit, offset int) ([]models.ProductGroup, int64, error) {
	return r.list(r.db.Model(&models.ProductGroup{}).
		Where("product_groups.company_id = ? AND product_groups.is_active = ?", companyID, true), limit, offset)
}
