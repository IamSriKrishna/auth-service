package repo

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productGroupInventoryRepository struct {
	db *gorm.DB
}

func NewProductGroupInventoryRepository(db *gorm.DB) *productGroupInventoryRepository {
	return &productGroupInventoryRepository{db: db}
}

func (r *productGroupInventoryRepository) Create(inventory *models.ProductGroupInventory) error {
	if err := r.db.Create(inventory).Error; err != nil {
		return fmt.Errorf("failed to create product group inventory: %w", err)
	}
	return nil
}

func (r *productGroupInventoryRepository) FindByID(id uint) (*models.ProductGroupInventory, error) {
	var inventory models.ProductGroupInventory
	if err := r.db.Where("id = ?", id).First(&inventory).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *productGroupInventoryRepository) FindByProductGroupID(productGroupID string) (*models.ProductGroupInventory, error) {
	var inventory models.ProductGroupInventory
	if err := r.db.Where("product_group_id = ?", productGroupID).First(&inventory).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *productGroupInventoryRepository) Update(inventory *models.ProductGroupInventory) error {
	if err := r.db.Save(inventory).Error; err != nil {
		return fmt.Errorf("failed to update product group inventory: %w", err)
	}
	return nil
}

func (r *productGroupInventoryRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.ProductGroupInventory{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete product group inventory: %w", err)
	}
	return nil
}

func (r *productGroupInventoryRepository) GetLowStockGroups(threshold float64) ([]models.ProductGroupInventory, error) {
	var inventories []models.ProductGroupInventory
	if err := r.db.Where("available_stock < ?", threshold).Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

// ComponentInventory Repository Implementation
type componentInventoryRepository struct {
	db *gorm.DB
}

func NewComponentInventoryRepository(db *gorm.DB) *componentInventoryRepository {
	return &componentInventoryRepository{db: db}
}

func (r *componentInventoryRepository) Create(inventory *models.ComponentInventory) error {
	if err := r.db.Create(inventory).Error; err != nil {
		return fmt.Errorf("failed to create component inventory: %w", err)
	}
	return nil
}

func (r *componentInventoryRepository) FindByID(id uint) (*models.ComponentInventory, error) {
	var inventory models.ComponentInventory
	if err := r.db.Where("id = ?", id).First(&inventory).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *componentInventoryRepository) FindByProductGroupID(productGroupID string) ([]models.ComponentInventory, error) {
	var inventories []models.ComponentInventory
	if err := r.db.Where("product_group_id = ?", productGroupID).
		Preload("ComponentProduct").
		Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

func (r *componentInventoryRepository) FindByComponentProductID(productID string) ([]models.ComponentInventory, error) {
	var inventories []models.ComponentInventory
	if err := r.db.Where("component_product_id = ?", productID).Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

func (r *componentInventoryRepository) Update(inventory *models.ComponentInventory) error {
	if err := r.db.Save(inventory).Error; err != nil {
		return fmt.Errorf("failed to update component inventory: %w", err)
	}
	return nil
}

func (r *componentInventoryRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.ComponentInventory{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete component inventory: %w", err)
	}
	return nil
}

func (r *componentInventoryRepository) UpdateBatch(inventories []models.ComponentInventory) error {
	if err := r.db.Save(inventories).Error; err != nil {
		return fmt.Errorf("failed to batch update component inventories: %w", err)
	}
	return nil
}

// ProductGroupTransaction Repository Implementation
type productGroupTransactionRepository struct {
	db *gorm.DB
}

func NewProductGroupTransactionRepository(db *gorm.DB) *productGroupTransactionRepository {
	return &productGroupTransactionRepository{db: db}
}

func (r *productGroupTransactionRepository) Create(transaction *models.ProductGroupTransaction) error {
	if err := r.db.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create product group transaction: %w", err)
	}
	return nil
}

func (r *productGroupTransactionRepository) FindByID(id uint) (*models.ProductGroupTransaction, error) {
	var transaction models.ProductGroupTransaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *productGroupTransactionRepository) FindByProductGroupID(productGroupID string, limit, offset int) ([]models.ProductGroupTransaction, int64, error) {
	var transactions []models.ProductGroupTransaction
	var total int64

	if err := r.db.Where("product_group_id = ?", productGroupID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *productGroupTransactionRepository) FindByReferenceID(referenceID string) ([]models.ProductGroupTransaction, error) {
	var transactions []models.ProductGroupTransaction
	if err := r.db.Where("purchase_order_id = ? OR sales_order_id = ? OR shipment_id = ?", referenceID, referenceID, referenceID).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *productGroupTransactionRepository) GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.ProductGroupTransaction, int64, error) {
	var transactions []models.ProductGroupTransaction
	var total int64

	if err := r.db.Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *productGroupTransactionRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.ProductGroupTransaction{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete product group transaction: %w", err)
	}
	return nil
}
