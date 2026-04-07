package repo

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("ProductDetails", "SalesInfo", "PurchaseInfo", "Inventory", "ReturnPolicy").
			Create(product).Error; err != nil {
			return err
		}

		product.ProductDetails.ProductID = product.ID
		productIDPtr := product.ID
		product.SalesInfo.ProductID = &productIDPtr
		product.Inventory.ProductID = &productIDPtr
		product.ReturnPolicy.ProductID = &productIDPtr

		if product.PurchaseInfo.Account != "" {
			product.PurchaseInfo.ProductID = &productIDPtr
		}

		if err := tx.Create(&product.SalesInfo).Error; err != nil {
			return err
		}

		if product.PurchaseInfo.Account != "" {
			if err := tx.Create(&product.PurchaseInfo).Error; err != nil {
				return err
			}
		}

		if product.Inventory.TrackInventory {
			if err := tx.Create(&product.Inventory).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&product.ReturnPolicy).Error; err != nil {
			return err
		}

		if err := tx.Omit("ProductVariants").Create(&product.ProductDetails).Error; err != nil {
			return err
		}

		if len(product.ProductDetails.ProductVariants) > 0 {
			for i := range product.ProductDetails.ProductVariants {
				product.ProductDetails.ProductVariants[i].ProductDetailsID = product.ProductDetails.ID
			}

			if err := tx.Omit("Attributes").Create(&product.ProductDetails.ProductVariants).Error; err != nil {
				return err
			}

			var allAttributes []models.ProductVariantAttribute
			for i := range product.ProductDetails.ProductVariants {
				for j := range product.ProductDetails.ProductVariants[i].Attributes {
					product.ProductDetails.ProductVariants[i].Attributes[j].ProductVariantID = product.ProductDetails.ProductVariants[i].ID
					allAttributes = append(allAttributes, product.ProductDetails.ProductVariants[i].Attributes[j])
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

func (r *productRepository) FindByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("ProductDetails.ProductVariants.Attributes").
		Preload("ProductDetails.Manufacturer").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Preload("ProductDetails.ProductVariants.Attributes").
		Preload("ProductDetails.Manufacturer").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(product).
			Omit("ProductDetails", "SalesInfo", "PurchaseInfo", "Inventory", "ReturnPolicy").
			Updates(product).Error; err != nil {
			return err
		}

		if err := tx.Where("item_id = ?", product.ID).Updates(&product.SalesInfo).Error; err != nil {
			return err
		}

		if product.PurchaseInfo.Account != "" {
			if err := tx.Where("item_id = ?", product.ID).Save(&product.PurchaseInfo).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("item_id = ?", product.ID).Updates(&product.Inventory).Error; err != nil {
			return err
		}

		if err := tx.Where("item_id = ?", product.ID).Updates(&product.ReturnPolicy).Error; err != nil {
			return err
		}

		if err := tx.Where("item_id = ?", product.ID).
			Omit("ProductVariants").
			Updates(&product.ProductDetails).Error; err != nil {
			return err
		}

		if err := tx.Where("product_details_id = ?", product.ProductDetails.ID).
			Delete(&models.ProductVariant{}).Error; err != nil {
			return err
		}

		if len(product.ProductDetails.ProductVariants) > 0 {
			for i := range product.ProductDetails.ProductVariants {
				product.ProductDetails.ProductVariants[i].ProductDetailsID = product.ProductDetails.ID
			}

			if err := tx.Omit("Attributes").Create(&product.ProductDetails.ProductVariants).Error; err != nil {
				return err
			}

			var allAttributes []models.ProductVariantAttribute
			for i := range product.ProductDetails.ProductVariants {
				for j := range product.ProductDetails.ProductVariants[i].Attributes {
					product.ProductDetails.ProductVariants[i].Attributes[j].ProductVariantID = product.ProductDetails.ProductVariants[i].ID
					allAttributes = append(allAttributes, product.ProductDetails.ProductVariants[i].Attributes[j])
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

func (r *productRepository) Delete(id string) error {
	// Start transaction
	tx := r.db.Begin()

	// First, delete all stock ledgers for this product
	if err := tx.Where("product_id = ?", id).Delete(&models.StockLedger{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Then, delete all product stocks for this product
	if err := tx.Where("product_id = ?", id).Delete(&models.ProductStock{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Finally, delete the product itself
	if err := tx.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

func (r *productRepository) DeleteByCreatedBy(id string, createdBy string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, delete all stock ledgers for this product
		if err := tx.Where("product_id = ?", id).Delete(&models.StockLedger{}).Error; err != nil {
			return err
		}

		// Delete all product stocks for this product
		if err := tx.Where("product_id = ?", id).Delete(&models.ProductStock{}).Error; err != nil {
			return err
		}

		// Delete all product_group_components that reference this product
		if err := tx.Delete(&models.ProductGroupComponent{}, "product_id = ?", id).Error; err != nil {
			return err
		}

		// Then delete the product itself
		if err := tx.Delete(&models.Product{}, "id = ? AND created_by = ?", id, createdBy).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *productRepository) FindByCreatedBy(createdBy string, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Count total products for this user
	countResult := r.db.Model(&models.Product{}).Where("created_by = ?", createdBy).Count(&total)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	// Fetch products with all relationships
	result := r.db.
		Where("created_by = ?", createdBy).
		Preload("ProductDetails.ProductVariants.Attributes").
		Preload("ProductDetails.Manufacturer").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, total, nil
}

// DeductProductVariantStock reduces the stock quantity from a product variant
func (r *productRepository) DeductProductVariantStock(productID string, variantSKU string, quantity float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var balance models.InventoryBalance

		query := tx.Where("item_id = ? AND variant_sku = ?", productID, variantSKU)

		if err := query.First(&balance).Error; err != nil {
			return fmt.Errorf("inventory balance not found for product variant %s of product %s", variantSKU, productID)
		}

		// Check if enough stock is available
		if balance.AvailableQuantity < quantity {
			return fmt.Errorf("insufficient stock for product variant %s: available=%f, required=%f", variantSKU, balance.AvailableQuantity, quantity)
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

// CheckProductVariantReorderPoint verifies if current stock is at or below reorder level
func (r *productRepository) CheckProductVariantReorderPoint(productID string, variantSKU string) (*models.ProductVariant, error) {
	var product models.Product
	if err := r.db.
		Preload("ProductDetails.ProductVariants").
		Preload("Inventory").
		Where("id = ?", productID).
		First(&product).Error; err != nil {
		return nil, err
	}

	// Check variant reorder level
	for _, v := range product.ProductDetails.ProductVariants {
		if v.SKU == variantSKU {
			if v.StockQuantity <= v.ReorderLevel {
				return &v, nil // Returns variant at or below reorder point
			}
			return nil, nil // Stock is above reorder point
		}
	}
	return nil, fmt.Errorf("variant %s not found", variantSKU)
}

// GetProductVariantBySKU retrieves a product variant by its SKU
func (r *productRepository) GetProductVariantBySKU(sku string) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	if err := r.db.Where("sku = ?", sku).First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

// UpdateProductVariantStock updates the stock quantity for a specific product variant
func (r *productRepository) UpdateProductVariantStock(variantID uint, newQuantity float64) error {
	return r.db.Model(&models.ProductVariant{}).Where("id = ?", variantID).
		Update("stock_quantity", newQuantity).Error
}

// GetProductVariantsByProductID retrieves all variants for a product
func (r *productRepository) GetProductVariantsByProductID(productID string) ([]models.ProductVariant, error) {
	var product models.Product
	if err := r.db.
		Preload("ProductDetails.ProductVariants.Attributes").
		Where("id = ?", productID).
		First(&product).Error; err != nil {
		return nil, err
	}
	return product.ProductDetails.ProductVariants, nil
}
