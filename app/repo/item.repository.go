package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item *models.Item) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("ItemDetails", "SalesInfo", "PurchaseInfo", "Inventory", "ReturnPolicy").
			Create(item).Error; err != nil {
			return err
		}

		item.ItemDetails.ItemID = item.ID
		item.SalesInfo.ItemID = item.ID
		item.Inventory.ItemID = item.ID
		item.ReturnPolicy.ItemID = item.ID
		
		if item.PurchaseInfo.Account != "" {
			item.PurchaseInfo.ItemID = item.ID
		}

		if err := tx.Create(&item.SalesInfo).Error; err != nil {
			return err
		}

		if item.PurchaseInfo.Account != "" {
			if err := tx.Create(&item.PurchaseInfo).Error; err != nil {
				return err
			}
		}

		if item.Inventory.TrackInventory || item.ItemDetails.Structure == "single" {
			if err := tx.Create(&item.Inventory).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&item.ReturnPolicy).Error; err != nil {
			return err
		}

		if err := tx.Omit("Variants").Create(&item.ItemDetails).Error; err != nil {
			return err
		}

		if item.ItemDetails.Structure == "variants" && len(item.ItemDetails.Variants) > 0 {
			for i := range item.ItemDetails.Variants {
				item.ItemDetails.Variants[i].ItemDetailsID = item.ItemDetails.ID
			}

			if err := tx.Omit("Attributes").Create(&item.ItemDetails.Variants).Error; err != nil {
				return err
			}

			var allAttributes []models.VariantAttribute
			for i := range item.ItemDetails.Variants {
				for j := range item.ItemDetails.Variants[i].Attributes {
					item.ItemDetails.Variants[i].Attributes[j].VariantID = item.ItemDetails.Variants[i].ID
					allAttributes = append(allAttributes, item.ItemDetails.Variants[i].Attributes[j])
				}
			}

			if len(allAttributes) > 0 {
				if err := tx.Create(&allAttributes).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *itemRepository) FindByID(id string) (*models.Item, error) {
	var item models.Item
	err := r.db.
		Preload("ItemDetails.Variants.Attributes").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy").
		Where("id = ?", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) FindAll(limit, offset int) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	if err := r.db.Model(&models.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.

		Preload("ItemDetails.Variants.Attributes").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *itemRepository) Update(item *models.Item) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(item).
			Omit("ItemDetails", "SalesInfo", "PurchaseInfo", "Inventory", "ReturnPolicy").
			Updates(item).Error; err != nil {
			return err
		}

		if err := tx.Where("item_id = ?", item.ID).Updates(&item.SalesInfo).Error; err != nil {
			return err
		}

		if item.PurchaseInfo.Account != "" {
			if err := tx.Where("item_id = ?", item.ID).Save(&item.PurchaseInfo).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("item_id = ?", item.ID).Updates(&item.Inventory).Error; err != nil {
			return err
		}

		if err := tx.Where("item_id = ?", item.ID).Updates(&item.ReturnPolicy).Error; err != nil {
			return err
		}

		if err := tx.Where("item_id = ?", item.ID).
			Omit("Variants").
			Updates(&item.ItemDetails).Error; err != nil {
			return err
		}

		if item.ItemDetails.Structure == "variants" {
			if err := tx.Where("item_details_id = ?", item.ItemDetails.ID).
				Delete(&models.Variant{}).Error; err != nil {
				return err
			}

			if len(item.ItemDetails.Variants) > 0 {
				for i := range item.ItemDetails.Variants {
					item.ItemDetails.Variants[i].ItemDetailsID = item.ItemDetails.ID
				}

				if err := tx.Omit("Attributes").Create(&item.ItemDetails.Variants).Error; err != nil {
					return err
				}

				var allAttributes []models.VariantAttribute
				for i := range item.ItemDetails.Variants {
					for j := range item.ItemDetails.Variants[i].Attributes {
						item.ItemDetails.Variants[i].Attributes[j].VariantID = item.ItemDetails.Variants[i].ID
						allAttributes = append(allAttributes, item.ItemDetails.Variants[i].Attributes[j])
					}
				}

				if len(allAttributes) > 0 {
					if err := tx.Create(&allAttributes).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func (r *itemRepository) Delete(id string) error {
	return r.db.Delete(&models.Item{}, "id = ?", id).Error
}

func (r *itemRepository) FindByType(itemType string, limit, offset int) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := r.db.Model(&models.Item{}).Where("type = ?", itemType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("ItemDetails.Variants.Attributes").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
