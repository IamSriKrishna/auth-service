package repo

import (
	"errors"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productConversionRepository struct {
	db *gorm.DB
}

func NewProductConversionRepository(db *gorm.DB) ProductConversionRepository {
	return &productConversionRepository{db: db}
}

// Create creates a new product conversion rule
func (r *productConversionRepository) Create(conversion *models.ProductConversion) error {
	return r.db.Create(conversion).Error
}

// GetByID retrieves a conversion rule by ID
func (r *productConversionRepository) GetByID(id string) (*models.ProductConversion, error) {
	var conversion models.ProductConversion
	if err := r.db.Where("id = ?", id).First(&conversion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conversion, nil
}

// GetAll retrieves all conversion rules with pagination
func (r *productConversionRepository) GetAll(offset, limit int) ([]models.ProductConversion, int64, error) {
	var conversions []models.ProductConversion
	var total int64

	if err := r.db.Model(&models.ProductConversion{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&conversions).Error; err != nil {
		return nil, 0, err
	}

	return conversions, total, nil
}

// GetActiveConversions retrieves all active conversion rules
func (r *productConversionRepository) GetActiveConversions(offset, limit int) ([]models.ProductConversion, int64, error) {
	var conversions []models.ProductConversion
	var total int64

	if err := r.db.Model(&models.ProductConversion{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("is_active = ?", true).Offset(offset).Limit(limit).Find(&conversions).Error; err != nil {
		return nil, 0, err
	}

	return conversions, total, nil
}

// GetByRawProductID retrieves conversion rules for a raw product
func (r *productConversionRepository) GetByRawProductID(rawProductID string, offset, limit int) ([]models.ProductConversion, int64, error) {
	var conversions []models.ProductConversion
	var total int64

	if err := r.db.Model(&models.ProductConversion{}).Where("raw_product_id = ?", rawProductID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("raw_product_id = ?", rawProductID).Offset(offset).Limit(limit).Find(&conversions).Error; err != nil {
		return nil, 0, err
	}

	return conversions, total, nil
}

// GetByFinishedProductID retrieves conversion rules for a finished product
func (r *productConversionRepository) GetByFinishedProductID(finishedProductID string, offset, limit int) ([]models.ProductConversion, int64, error) {
	var conversions []models.ProductConversion
	var total int64

	if err := r.db.Model(&models.ProductConversion{}).Where("finished_product_id = ?", finishedProductID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("finished_product_id = ?", finishedProductID).Offset(offset).Limit(limit).Find(&conversions).Error; err != nil {
		return nil, 0, err
	}

	return conversions, total, nil
}

// GetByProductPair retrieves a conversion rule by raw and finished product IDs
func (r *productConversionRepository) GetByProductPair(rawProductID, finishedProductID string) (*models.ProductConversion, error) {
	var conversion models.ProductConversion
	if err := r.db.Where("raw_product_id = ? AND finished_product_id = ?", rawProductID, finishedProductID).First(&conversion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conversion, nil
}

// Update updates a conversion rule
func (r *productConversionRepository) Update(conversion *models.ProductConversion) error {
	return r.db.Model(conversion).Updates(conversion).Error
}

// Delete deletes a conversion rule
func (r *productConversionRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductConversion{}, "id = ?", id).Error
}

// GetDB returns the database connection
func (r *productConversionRepository) GetDB() *gorm.DB {
	return r.db
}

// --- ProductConversionRecord Repository ---

type productConversionRecordRepository struct {
	db *gorm.DB
}

func NewProductConversionRecordRepository(db *gorm.DB) ProductConversionRecordRepository {
	return &productConversionRecordRepository{db: db}
}

// Create creates a new conversion record
func (r *productConversionRecordRepository) Create(record *models.ProductConversionRecord) error {
	return r.db.Create(record).Error
}

// GetByID retrieves a conversion record by ID
func (r *productConversionRecordRepository) GetByID(id string) (*models.ProductConversionRecord, error) {
	var record models.ProductConversionRecord
	if err := r.db.Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetByConversionID retrieves conversion records for a specific conversion rule
func (r *productConversionRecordRepository) GetByConversionID(conversionID string, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	var records []models.ProductConversionRecord
	var total int64

	if err := r.db.Model(&models.ProductConversionRecord{}).Where("conversion_id = ?", conversionID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("conversion_id = ?", conversionID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetAll retrieves all conversion records with pagination
func (r *productConversionRecordRepository) GetAll(offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	var records []models.ProductConversionRecord
	var total int64

	if err := r.db.Model(&models.ProductConversionRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetByDateRange retrieves conversion records within a date range
func (r *productConversionRecordRepository) GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	var records []models.ProductConversionRecord
	var total int64

	query := r.db.Model(&models.ProductConversionRecord{}).Where("conversion_date BETWEEN ? AND ?", fromDate, toDate)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetByStatus retrieves conversion records by status
func (r *productConversionRecordRepository) GetByStatus(status string, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	var records []models.ProductConversionRecord
	var total int64

	if err := r.db.Model(&models.ProductConversionRecord{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("status = ?", status).Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// Update updates a conversion record
func (r *productConversionRecordRepository) Update(record *models.ProductConversionRecord) error {
	return r.db.Model(record).Updates(record).Error
}

// Delete deletes a conversion record
func (r *productConversionRecordRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductConversionRecord{}, "id = ?", id).Error
}

// GetDB returns the database connection
func (r *productConversionRecordRepository) GetDB() *gorm.DB {
	return r.db
}

// --- ConversionRecordBagUsage Repository ---

type conversionRecordBagUsageRepository struct {
	db *gorm.DB
}

func NewConversionRecordBagUsageRepository(db *gorm.DB) ConversionRecordBagUsageRepository {
	return &conversionRecordBagUsageRepository{db: db}
}

// Create creates a new bag usage record for a conversion record
func (r *conversionRecordBagUsageRepository) Create(bagUsage *models.ConversionRecordBagUsage) error {
	return r.db.Create(bagUsage).Error
}

// GetByConversionRecordID retrieves all bags used for a specific conversion record
func (r *conversionRecordBagUsageRepository) GetByConversionRecordID(recordID string) ([]models.ConversionRecordBagUsage, error) {
	var usages []models.ConversionRecordBagUsage
	if err := r.db.Where("conversion_record_id = ?", recordID).Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

// GetByBagID retrieves all conversion records that used a specific bag
func (r *conversionRecordBagUsageRepository) GetByBagID(bagID string) ([]models.ConversionRecordBagUsage, error) {
	var usages []models.ConversionRecordBagUsage
	if err := r.db.Where("bag_id = ?", bagID).Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

// GetDB returns the database connection
func (r *conversionRecordBagUsageRepository) GetDB() *gorm.DB {
	return r.db
}
