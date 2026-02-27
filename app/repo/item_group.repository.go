package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)


type itemGroupRepository struct {
	db *gorm.DB
}

func NewItemGroupRepository(db *gorm.DB) ItemGroupRepository {
	return &itemGroupRepository{db: db}
}

func (r *itemGroupRepository) Create(itemGroup *models.ItemGroup) error {
	return r.db.Create(itemGroup).Error
}

func (r *itemGroupRepository) FindByID(id string) (*models.ItemGroup, error) {
	var itemGroup models.ItemGroup
	err := r.db.Preload("Components.Item").First(&itemGroup, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &itemGroup, nil
}

func (r *itemGroupRepository) FindAll(limit, offset int, search string) ([]models.ItemGroup, int64, error) {
	var itemGroups []models.ItemGroup
	var count int64

	query := r.db.Model(&models.ItemGroup{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Components.Item").Limit(limit).Offset(offset).Order("created_at DESC").Find(&itemGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return itemGroups, count, nil
}

func (r *itemGroupRepository) Update(itemGroup *models.ItemGroup) error {
	return r.db.Save(itemGroup).Error
}

func (r *itemGroupRepository) Delete(id string) error {
	return r.db.Delete(&models.ItemGroup{}, "id = ?", id).Error
}

func (r *itemGroupRepository) FindByName(name string) (*models.ItemGroup, error) {
	var itemGroup models.ItemGroup
	err := r.db.Preload("Components.Item").First(&itemGroup, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &itemGroup, nil
}

func (r *itemGroupRepository) FindActiveGroups(limit, offset int) ([]models.ItemGroup, int64, error) {
	var itemGroups []models.ItemGroup
	var count int64

	query := r.db.Where("is_active = ?", true)

	err := query.Model(&models.ItemGroup{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Components.Item").Limit(limit).Offset(offset).Order("created_at DESC").Find(&itemGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return itemGroups, count, nil
}
