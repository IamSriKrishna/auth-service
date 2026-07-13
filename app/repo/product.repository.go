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

func productPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("ProductDetails.ProductVariants.Attributes").
		Preload("SalesInfo").
		Preload("PurchaseInfo.PreferredVendor").
		Preload("Inventory").
		Preload("ReturnPolicy")
}

func (r *productRepository) Create(
	product *models.Product,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Omit(
				"ProductDetails",
				"SalesInfo",
				"PurchaseInfo",
				"Inventory",
				"ReturnPolicy",
			).
			Create(product).
			Error; err != nil {
			return err
		}

		product.ProductDetails.ProductID = product.ID
		productIDPointer := product.ID
		product.SalesInfo.ProductID = &productIDPointer
		product.Inventory.ProductID = &productIDPointer
		product.ReturnPolicy.ProductID = &productIDPointer

		if product.PurchaseInfo.Account != "" {
			product.PurchaseInfo.ProductID = &productIDPointer
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

		if err := tx.
			Omit("ProductVariants").
			Create(&product.ProductDetails).
			Error; err != nil {
			return err
		}

		if len(product.ProductDetails.ProductVariants) > 0 {
			for index := range product.ProductDetails.ProductVariants {
				product.ProductDetails.ProductVariants[index].
					ProductDetailsID = product.ProductDetails.ID
			}

			if err := tx.
				Omit("Attributes").
				Create(&product.ProductDetails.ProductVariants).
				Error; err != nil {
				return err
			}

			var attributes []models.ProductVariantAttribute

			for variantIndex := range product.ProductDetails.ProductVariants {
				for attributeIndex := range product.ProductDetails.
					ProductVariants[variantIndex].
					Attributes {
					product.ProductDetails.
						ProductVariants[variantIndex].
						Attributes[attributeIndex].
						ProductVariantID = product.
						ProductDetails.
						ProductVariants[variantIndex].
						ID

					attributes = append(
						attributes,
						product.ProductDetails.
							ProductVariants[variantIndex].
							Attributes[attributeIndex],
					)
				}
			}

			if len(attributes) > 0 {
				if err := tx.Create(&attributes).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *productRepository) FindByID(
	id string,
) (*models.Product, error) {
	var product models.Product

	err := productPreloads(r.db).
		Where("products.id = ?", id).
		First(&product).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.Product, error) {
	var product models.Product

	err := productPreloads(r.db).
		Where(
			"products.id = ? AND products.created_by_company_id = ?",
			id,
			companyID,
		).
		First(&product).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindAll(
	limit int,
	offset int,
) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if err := r.db.
		Model(&models.Product{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	err := productPreloads(r.db).
		Limit(limit).
		Offset(offset).
		Order("products.created_at DESC").
		Find(&products).
		Error

	return products, total, err
}

func (r *productRepository) FindByCreatedBy(
	createdBy string,
	limit int,
	offset int,
) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.
		Model(&models.Product{}).
		Where("created_by = ?", createdBy)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := productPreloads(r.db).
		Where("created_by = ?", createdBy).
		Limit(limit).
		Offset(offset).
		Order("products.created_at DESC").
		Find(&products).
		Error

	return products, total, err
}

func (r *productRepository) FindByCreatedByAndCompany(
	createdBy string,
	companyID uint,
	limit int,
	offset int,
) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.
		Model(&models.Product{}).
		Where(
			"created_by = ? AND created_by_company_id = ?",
			createdBy,
			companyID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := productPreloads(r.db).
		Where(
			"created_by = ? AND created_by_company_id = ?",
			createdBy,
			companyID,
		).
		Limit(limit).
		Offset(offset).
		Order("products.created_at DESC").
		Find(&products).
		Error

	return products, total, err
}

func (r *productRepository) FindByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.
		Model(&models.Product{}).
		Where("created_by_company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := productPreloads(r.db).
		Where("created_by_company_id = ?", companyID).
		Limit(limit).
		Offset(offset).
		Order("products.created_at DESC").
		Find(&products).
		Error

	return products, total, err
}

func (r *productRepository) Update(
	product *models.Product,
) error {
	return r.updateProductTransaction(
		r.db,
		product,
	)
}

func (r *productRepository) UpdateByCompany(
	product *models.Product,
	companyID uint,
) error {
	if product == nil {
		return gorm.ErrInvalidData
	}

	var count int64
	if err := r.db.
		Model(&models.Product{}).
		Where(
			"id = ? AND created_by_company_id = ?",
			product.ID,
			companyID,
		).
		Count(&count).
		Error; err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.updateProductTransaction(
		r.db,
		product,
	)
}

func (r *productRepository) updateProductTransaction(
	db *gorm.DB,
	product *models.Product,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(product).
			Omit(
				"ProductDetails",
				"SalesInfo",
				"PurchaseInfo",
				"Inventory",
				"ReturnPolicy",
			).
			Updates(product).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("product_id = ?", product.ID).
			Updates(&product.SalesInfo).
			Error; err != nil {
			return err
		}

		if product.PurchaseInfo.Account != "" {
			if err := tx.
				Where("product_id = ?", product.ID).
				Save(&product.PurchaseInfo).
				Error; err != nil {
				return err
			}
		}

		if err := tx.
			Where("product_id = ?", product.ID).
			Updates(&product.Inventory).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("product_id = ?", product.ID).
			Updates(&product.ReturnPolicy).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("product_id = ?", product.ID).
			Omit("ProductVariants").
			Updates(&product.ProductDetails).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where(
				"product_details_id = ?",
				product.ProductDetails.ID,
			).
			Delete(&models.ProductVariant{}).
			Error; err != nil {
			return err
		}

		if len(product.ProductDetails.ProductVariants) > 0 {
			for index := range product.ProductDetails.ProductVariants {
				product.ProductDetails.ProductVariants[index].
					ProductDetailsID = product.ProductDetails.ID
			}

			if err := tx.
				Omit("Attributes").
				Create(&product.ProductDetails.ProductVariants).
				Error; err != nil {
				return err
			}

			var attributes []models.ProductVariantAttribute

			for variantIndex := range product.ProductDetails.ProductVariants {
				for attributeIndex := range product.ProductDetails.
					ProductVariants[variantIndex].
					Attributes {
					product.ProductDetails.
						ProductVariants[variantIndex].
						Attributes[attributeIndex].
						ProductVariantID = product.
						ProductDetails.
						ProductVariants[variantIndex].
						ID

					attributes = append(
						attributes,
						product.ProductDetails.
							ProductVariants[variantIndex].
							Attributes[attributeIndex],
					)
				}
			}

			if len(attributes) > 0 {
				if err := tx.Create(&attributes).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *productRepository) Delete(
	id string,
) error {
	return r.deleteProductTransaction(
		r.db,
		id,
	)
}

func (r *productRepository) DeleteByCreatedBy(
	id string,
	createdBy string,
) error {
	var count int64

	if err := r.db.
		Model(&models.Product{}).
		Where(
			"id = ? AND created_by = ?",
			id,
			createdBy,
		).
		Count(&count).
		Error; err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.deleteProductTransaction(
		r.db,
		id,
	)
}

func (r *productRepository) DeleteByCompany(
	id string,
	companyID uint,
) error {
	var count int64

	if err := r.db.
		Model(&models.Product{}).
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Count(&count).
		Error; err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.deleteProductTransaction(
		r.db,
		id,
	)
}

func (r *productRepository) deleteProductTransaction(
	db *gorm.DB,
	id string,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("product_id = ?", id).
			Delete(&models.StockLedger{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("product_id = ?", id).
			Delete(&models.ProductStock{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			Delete(
				&models.ProductGroupComponent{},
				"product_id = ?",
				id,
			).
			Error; err != nil {
			return err
		}

		return tx.
			Delete(&models.Product{}, "id = ?", id).
			Error
	})
}

func (r *productRepository) DeductProductVariantStock(
	productID string,
	variantSKU string,
	quantity float64,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var balance models.InventoryBalance

		if err := tx.
			Where(
				"item_id = ? AND variant_sku = ?",
				productID,
				variantSKU,
			).
			First(&balance).
			Error; err != nil {
			return fmt.Errorf(
				"inventory balance not found for product variant %s of product %s",
				variantSKU,
				productID,
			)
		}

		if balance.AvailableQuantity < quantity {
			return fmt.Errorf(
				"insufficient stock for product variant %s: available=%f, required=%f",
				variantSKU,
				balance.AvailableQuantity,
				quantity,
			)
		}

		balance.AvailableQuantity -= quantity
		balance.CurrentQuantity -= quantity
		balance.UpdatedAt = time.Now()

		return tx.
			Model(&balance).
			Updates(balance).
			Error
	})
}

func (r *productRepository) CheckProductVariantReorderPoint(
	productID string,
	variantSKU string,
) (*models.ProductVariant, error) {
	var product models.Product

	if err := r.db.
		Preload("ProductDetails.ProductVariants").
		Preload("Inventory").
		Where("id = ?", productID).
		First(&product).
		Error; err != nil {
		return nil, err
	}

	for _, variant := range product.ProductDetails.ProductVariants {
		if variant.SKU == variantSKU {
			if variant.StockQuantity <= variant.ReorderLevel {
				return &variant, nil
			}

			return nil, nil
		}
	}

	return nil, fmt.Errorf("variant %s not found", variantSKU)
}

func (r *productRepository) GetProductVariantBySKU(
	sku string,
) (*models.ProductVariant, error) {
	var variant models.ProductVariant

	if err := r.db.
		Where("sku = ?", sku).
		First(&variant).
		Error; err != nil {
		return nil, err
	}

	return &variant, nil
}

func (r *productRepository) UpdateProductVariantStock(
	variantID uint,
	newQuantity float64,
) error {
	return r.db.
		Model(&models.ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock_quantity", newQuantity).
		Error
}

func (r *productRepository) GetProductVariantsByProductID(
	productID string,
) ([]models.ProductVariant, error) {
	var product models.Product

	if err := r.db.
		Preload("ProductDetails.ProductVariants.Attributes").
		Where("id = ?", productID).
		First(&product).
		Error; err != nil {
		return nil, err
	}

	return product.ProductDetails.ProductVariants, nil
}

func (r *productRepository) GetProductVariantsByProductIDAndCompany(
	productID string,
	companyID uint,
) ([]models.ProductVariant, error) {
	var product models.Product

	if err := r.db.
		Preload("ProductDetails.ProductVariants.Attributes").
		Where(
			"id = ? AND created_by_company_id = ?",
			productID,
			companyID,
		).
		First(&product).
		Error; err != nil {
		return nil, err
	}

	return product.ProductDetails.ProductVariants, nil
}
