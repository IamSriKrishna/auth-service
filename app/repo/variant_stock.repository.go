package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type variantStockRepository struct {
	db *gorm.DB
}

func NewVariantStockRepository(db *gorm.DB) VariantStockRepository {
	return &variantStockRepository{db: db}
}

func (r *variantStockRepository) Create(stock *models.VariantStock) error {
	return r.db.Create(stock).Error
}

func (r *variantStockRepository) GetByID(id string) (*models.VariantStock, error) {
	var stock models.VariantStock
	err := r.db.Preload("Product").Where("id = ?", id).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *variantStockRepository) GetBySKU(sku string) (*models.VariantStock, error) {
	var stock models.VariantStock
	err := r.db.Preload("Product").Where("variant_sku = ?", sku).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *variantStockRepository) GetByProductID(productID string, offset, limit int) ([]models.VariantStock, int64, error) {
	var stocks []models.VariantStock
	var total int64

	err := r.db.Model(&models.VariantStock{}).
		Where("product_id = ?", productID).
		Count(&total).
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

func (r *variantStockRepository) Update(stock *models.VariantStock) error {
	return r.db.Save(stock).Error
}

func (r *variantStockRepository) Delete(id string) error {
	return r.db.Delete(&models.VariantStock{}, "id = ?", id).Error
}

func (r *variantStockRepository) GetAll(offset, limit int) ([]models.VariantStock, int64, error) {
	var stocks []models.VariantStock
	var total int64

	err := r.db.Model(&models.VariantStock{}).
		Count(&total).
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

func (r *variantStockRepository) GetAllByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error) {
	var stocks []models.VariantStock
	var total int64

	err := r.db.Model(&models.VariantStock{}).
		Joins("JOIN products ON variant_stocks.product_id = products.id").
		Where("products.created_by = ?", userID).
		Count(&total).
		Preload("Product").
		Offset(offset).
		Limit(limit).
		Order("variant_stocks.created_at DESC").
		Find(&stocks).Error

	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (r *variantStockRepository) GetBySKUs(skus []string) ([]models.VariantStock, error) {
	var stocks []models.VariantStock
	err := r.db.Preload("Product").
		Where("variant_sku IN ?", skus).
		Find(&stocks).Error
	return stocks, err
}

func (r *variantStockRepository) GetLowStockVariants(threshold float64, offset, limit int) ([]models.VariantStock, int64, error) {
	var stocks []models.VariantStock
	var total int64

	query := r.db.Model(&models.VariantStock{}).
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

// VariantStockMovement Repository

type variantStockMovementRepository struct {
	db *gorm.DB
}

func NewVariantStockMovementRepository(db *gorm.DB) VariantStockMovementRepository {
	return &variantStockMovementRepository{db: db}
}

func (r *variantStockMovementRepository) Create(movement *models.VariantStockMovement) error {
	return r.db.Create(movement).Error
}

func (r *variantStockMovementRepository) GetByID(id uint) (*models.VariantStockMovement, error) {
	var movement models.VariantStockMovement
	err := r.db.Where("id = ?", id).First(&movement).Error
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

func (r *variantStockMovementRepository) GetByVariantSKU(sku string, offset, limit int) ([]models.VariantStockMovement, int64, error) {
	var movements []models.VariantStockMovement
	var total int64

	err := r.db.Model(&models.VariantStockMovement{}).
		Where("variant_sku = ?", sku).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&movements).Error

	if err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}

func (r *variantStockMovementRepository) GetByReferenceID(referenceID string) ([]models.VariantStockMovement, error) {
	var movements []models.VariantStockMovement
	err := r.db.Where("reference_id = ?", referenceID).
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}

func (r *variantStockMovementRepository) GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.VariantStockMovement, int64, error) {
	var movements []models.VariantStockMovement
	var total int64

	err := r.db.Model(&models.VariantStockMovement{}).
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&movements).Error

	if err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}

func (r *variantStockMovementRepository) DeleteByReferenceID(referenceID string) error {
	return r.db.Where("reference_id = ?", referenceID).Delete(&models.VariantStockMovement{}).Error
}

// StockReservation Repository

type stockReservationRepository struct {
	db *gorm.DB
}

func NewStockReservationRepository(db *gorm.DB) StockReservationRepository {
	return &stockReservationRepository{db: db}
}

func (r *stockReservationRepository) Create(reservation *models.StockReservation) error {
	return r.db.Create(reservation).Error
}

func (r *stockReservationRepository) GetByID(id uint) (*models.StockReservation, error) {
	var reservation models.StockReservation
	err := r.db.Where("id = ?", id).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *stockReservationRepository) GetBySalesOrderID(salesOrderID string) ([]models.StockReservation, error) {
	var reservations []models.StockReservation
	err := r.db.Where("sales_order_id = ?", salesOrderID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

func (r *stockReservationRepository) GetByVariantSKU(sku string, offset, limit int) ([]models.StockReservation, int64, error) {
	var reservations []models.StockReservation
	var total int64

	err := r.db.Model(&models.StockReservation{}).
		Where("variant_sku = ?", sku).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&reservations).Error

	if err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

func (r *stockReservationRepository) GetByStatus(status string, offset, limit int) ([]models.StockReservation, int64, error) {
	var reservations []models.StockReservation
	var total int64

	err := r.db.Model(&models.StockReservation{}).
		Where("status = ?", status).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&reservations).Error

	if err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

func (r *stockReservationRepository) Update(reservation *models.StockReservation) error {
	return r.db.Save(reservation).Error
}

func (r *stockReservationRepository) Delete(id uint) error {
	return r.db.Delete(&models.StockReservation{}, "id = ?", id).Error
}

func (r *stockReservationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.StockReservation{}).Where("id = ?", id).Update("status", status).Error
}
