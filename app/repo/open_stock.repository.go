package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type openingStockRepository struct {
	db *gorm.DB
}

func NewOpeningStockRepository(db *gorm.DB) OpeningStockRepository {
	return &openingStockRepository{db: db}
}

func (r *openingStockRepository) CreateOrUpdateOpeningStock(itemID string, openingStock, ratePerUnit float64) error {
	var existing models.OpeningStock
	err := r.db.Where("item_id = ?", itemID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&models.OpeningStock{
			ItemID:                  itemID,
			OpeningStock:            openingStock,
			OpeningStockRatePerUnit: ratePerUnit,
		}).Error
	}

	if err != nil {
		return err
	}

	return r.db.Model(&existing).Updates(map[string]interface{}{
		"opening_stock":               openingStock,
		"opening_stock_rate_per_unit": ratePerUnit,
	}).Error
}

func (r *openingStockRepository) GetOpeningStock(itemID string) (*models.OpeningStock, error) {
	var stock models.OpeningStock
	err := r.db.Where("item_id = ?", itemID).First(&stock).Error
	if err == gorm.ErrRecordNotFound {
		return &models.OpeningStock{
			ItemID:                  itemID,
			OpeningStock:            0,
			OpeningStockRatePerUnit: 0,
		}, nil
	}
	return &stock, err
}

func (r *openingStockRepository) CreateOrUpdateVariantOpeningStock(variantID uint, openingStock, ratePerUnit float64) error {
	var existing models.VariantOpeningStock
	err := r.db.Where("variant_id = ?", variantID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(&models.VariantOpeningStock{
			VariantID:               variantID,
			OpeningStock:            openingStock,
			OpeningStockRatePerUnit: ratePerUnit,
		}).Error
	}

	if err != nil {
		return err
	}

	return r.db.Model(&existing).Updates(map[string]interface{}{
		"opening_stock":               openingStock,
		"opening_stock_rate_per_unit": ratePerUnit,
	}).Error
}

func (r *openingStockRepository) GetVariantOpeningStock(variantID uint) (*models.VariantOpeningStock, error) {
	var stock models.VariantOpeningStock
	err := r.db.Where("variant_id = ?", variantID).First(&stock).Error
	if err == gorm.ErrRecordNotFound {
		return &models.VariantOpeningStock{
			VariantID:               variantID,
			OpeningStock:            0,
			OpeningStockRatePerUnit: 0,
		}, nil
	}
	return &stock, err
}

func (r *openingStockRepository) GetAllVariantOpeningStocks(itemID string) ([]models.VariantOpeningStock, error) {
	var stocks []models.VariantOpeningStock
	err := r.db.Raw(`
		SELECT vos.* 
		FROM variant_opening_stock vos
		INNER JOIN variants v ON v.id = vos.variant_id
		INNER JOIN item_details id ON id.id = v.item_details_id
		WHERE id.item_id = ?
	`, itemID).Scan(&stocks).Error

	return stocks, err
}

func (r *openingStockRepository) RecordStockMovement(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}

func (r *openingStockRepository) GetStockMovements(itemID string) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Where("item_id = ?", itemID).
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}
