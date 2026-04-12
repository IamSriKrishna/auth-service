package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productStockRepository struct {
	db *gorm.DB
}

func NewProductStockRepository(db *gorm.DB) ProductStockRepository {
	return &productStockRepository{db: db}
}

func (r *productStockRepository) Create(stock *models.ProductStock) error {
	return r.db.Create(stock).Error
}

func (r *productStockRepository) GetByID(id string) (*models.ProductStock, error) {
	var stock models.ProductStock
	err := r.db.Preload("Product").Where("id = ?", id).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *productStockRepository) GetByProductID(productID string) (*models.ProductStock, error) {
	var stock models.ProductStock
	err := r.db.Preload("Product").Where("product_id = ?", productID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *productStockRepository) Update(stock *models.ProductStock) error {
	return r.db.Save(stock).Error
}

func (r *productStockRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductStock{}, "id = ?", id).Error
}

func (r *productStockRepository) GetAll(offset, limit int) ([]models.ProductStock, int64, error) {
	var stocks []models.ProductStock
	var total int64

	err := r.db.Model(&models.ProductStock{}).Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&stocks).Error

	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (r *productStockRepository) GetByProductIDs(productIDs []string) ([]models.ProductStock, error) {
	var stocks []models.ProductStock
	err := r.db.Preload("Product").
		Where("product_id IN ?", productIDs).
		Find(&stocks).Error
	return stocks, err
}

func (r *productStockRepository) GetLowStockProducts(threshold float64, offset, limit int) ([]models.ProductStock, int64, error) {
	var stocks []models.ProductStock
	var total int64

	query := r.db.Model(&models.ProductStock{}).
		Where("available_stock <= ?", threshold)

	err := query.Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("available_stock ASC").
		Find(&stocks).Error

	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

// GetAllByUser retrieves all product stocks for a specific user
func (r *productStockRepository) GetAllByUser(userID uint, offset, limit int) ([]models.ProductStock, int64, error) {
	var stocks []models.ProductStock
	var total int64

	err := r.db.Model(&models.ProductStock{}).
		Joins("JOIN products ON product_stocks.product_id = products.id").
		Where("products.created_by = ?", userID).
		Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("product_stocks.created_at DESC").
		Find(&stocks).Error

	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

// GetLowStockProductsByUser retrieves low stock products for a specific user
func (r *productStockRepository) GetLowStockProductsByUser(userID uint, threshold float64, offset, limit int) ([]models.ProductStock, int64, error) {
	var stocks []models.ProductStock
	var total int64

	query := r.db.Model(&models.ProductStock{}).
		Joins("JOIN products ON product_stocks.product_id = products.id").
		Where("products.created_by = ? AND product_stocks.available_stock <= ?", userID, threshold)

	err := query.Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("product_stocks.available_stock ASC").
		Find(&stocks).Error

	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

type stockLedgerRepository struct {
	db *gorm.DB
}

func NewStockLedgerRepository(db *gorm.DB) StockLedgerRepository {
	return &stockLedgerRepository{db: db}
}

func (r *stockLedgerRepository) Create(ledger *models.StockLedger) error {
	return r.db.Create(ledger).Error
}

func (r *stockLedgerRepository) GetByID(id uint) (*models.StockLedger, error) {
	var ledger models.StockLedger
	err := r.db.Preload("Product").Where("id = ?", id).First(&ledger).Error
	if err != nil {
		return nil, err
	}
	return &ledger, nil
}

func (r *stockLedgerRepository) GetByProductID(productID string, offset, limit int) ([]models.StockLedger, int64, error) {
	var ledgers []models.StockLedger
	var total int64

	err := r.db.Model(&models.StockLedger{}).
		Where("product_id = ?", productID).
		Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ledgers).Error

	if err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}

func (r *stockLedgerRepository) GetByReferenceID(referenceID string) ([]models.StockLedger, error) {
	var ledgers []models.StockLedger
	err := r.db.Preload("Product").
		Where("reference_id = ?", referenceID).
		Order("created_at DESC").
		Find(&ledgers).Error
	return ledgers, err
}

func (r *stockLedgerRepository) DeleteByReferenceID(referenceID string) error {
	return r.db.Where("reference_id = ?", referenceID).Delete(&models.StockLedger{}).Error
}

func (r *stockLedgerRepository) GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.StockLedger, int64, error) {
	var ledgers []models.StockLedger
	var total int64

	err := r.db.Model(&models.StockLedger{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ledgers).Error

	if err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}

func (r *stockLedgerRepository) GetProductMovementHistory(productID string, offset, limit int) ([]models.StockLedger, int64, error) {
	var ledgers []models.StockLedger
	var total int64

	err := r.db.Model(&models.StockLedger{}).
		Where("product_id = ?", productID).
		Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ledgers).Error

	if err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}
