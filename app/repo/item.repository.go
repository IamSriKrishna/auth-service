package repo

import (
	"fmt"
	"time"

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

// DeductStockQuantity reduces the stock quantity from an item or variant
// This checks the inventory_balance table (which includes opening stock, purchases, etc.)
// If variantSKU is provided, it deducts from the variant stock balance
// Otherwise, it deducts from the item stock balance
func (r *itemRepository) DeductStockQuantity(itemID string, variantSKU *string, quantity float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var balance models.InventoryBalance

		// Get inventory balance for the item or variant
		query := tx.Where("item_id = ?", itemID)
		if variantSKU != nil && *variantSKU != "" {
			query = query.Where("variant_sku = ?", *variantSKU)
		} else {
			query = query.Where("variant_sku IS NULL")
		}

		if err := query.First(&balance).Error; err != nil {
			if variantSKU != nil && *variantSKU != "" {
				return fmt.Errorf("inventory balance not found for variant %s of item %s", *variantSKU, itemID)
			}
			return fmt.Errorf("inventory balance not found for item %s", itemID)
		}

		// Check if enough stock is available
		if balance.AvailableQuantity < quantity {
			if variantSKU != nil && *variantSKU != "" {
				return fmt.Errorf("insufficient stock for variant %s: available=%f, required=%f", *variantSKU, balance.AvailableQuantity, quantity)
			}
			return fmt.Errorf("insufficient stock for item %s: available=%f, required=%f", itemID, balance.AvailableQuantity, quantity)
		}

		// Deduct from inventory balance
		balance.AvailableQuantity -= quantity
		balance.CurrentQuantity -= quantity
		balance.UpdatedAt = time.Now()

		if err := tx.Model(&balance).Updates(balance).Error; err != nil {
			return fmt.Errorf("failed to update inventory balance: %v", err)
		}

		return nil
	})
}

// CheckReorderPoint verifies if current stock is at or below reorder level
// Returns the variant with reorder point information
func (r *itemRepository) CheckReorderPoint(itemID string, variantSKU *string) (*models.Variant, error) {
	var item models.Item
	if err := r.db.
		Preload("ItemDetails.Variants").
		Preload("Inventory").
		Where("id = ?", itemID).
		First(&item).Error; err != nil {
		return nil, err
	}

	if variantSKU != nil && *variantSKU != "" {
		// Check variant reorder level
		for _, v := range item.ItemDetails.Variants {
			if v.SKU == *variantSKU {
				if v.StockQuantity <= v.ReorderLevel {
					return &v, nil // Returns variant at or below reorder point
				}
				return nil, nil // Stock is above reorder point
			}
		}
		return nil, fmt.Errorf("variant %s not found", *variantSKU)
	}

	return nil, nil
}

// GetVariantBySKU retrieves a variant by its SKU
func (r *itemRepository) GetVariantBySKU(sku string) (*models.Variant, error) {
	var variant models.Variant
	if err := r.db.Where("sku = ?", sku).First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

// UpdateVariantStock updates the stock quantity for a specific variant
func (r *itemRepository) UpdateVariantStock(variantID uint, newQuantity float64) error {
	return r.db.Model(&models.Variant{}).Where("id = ?", variantID).
		Update("stock_quantity", newQuantity).Error
}
