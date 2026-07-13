package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productConversionRepository struct{ db *gorm.DB }

func NewProductConversionRepository(db *gorm.DB) ProductConversionRepository {
	return &productConversionRepository{db: db}
}

func conversionCompanyQuery(db *gorm.DB, companyID uint) *gorm.DB {
	return db.Joins(`
		JOIN users conversion_creator
			ON conversion_creator.id =
				CAST(product_conversions.created_by AS UNSIGNED)
	`).Where("conversion_creator.company_id = ?", companyID)
}

func (r *productConversionRepository) Create(c *models.ProductConversion) error {
	return r.db.Create(c).Error
}

func (r *productConversionRepository) CreateForCompany(c *models.ProductConversion, companyID uint) error {
	if c == nil || c.CreatedBy == "" {
		return gorm.ErrInvalidData
	}
	var count int64
	if err := r.db.Table("users").
		Where("id = CAST(? AS UNSIGNED) AND company_id = ?", c.CreatedBy, companyID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.Create(c).Error
}

func (r *productConversionRepository) GetByID(id string) (*models.ProductConversion, error) {
	var c models.ProductConversion
	if err := r.db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *productConversionRepository) GetByIDAndCompany(id string, companyID uint) (*models.ProductConversion, error) {
	var c models.ProductConversion
	q := conversionCompanyQuery(r.db.Model(&models.ProductConversion{}), companyID)
	if err := q.Select("product_conversions.*").
		Where("product_conversions.id = ?", id).
		First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *productConversionRepository) list(query *gorm.DB, offset, limit int) ([]models.ProductConversion, int64, error) {
	var rows []models.ProductConversion
	var total int64
	if err := query.Distinct("product_conversions.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Select("product_conversions.*").
		Order("product_conversions.created_at DESC").
		Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *productConversionRepository) GetAll(offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(r.db.Model(&models.ProductConversion{}), offset, limit)
}
func (r *productConversionRepository) GetAllByCompany(companyID uint, offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(conversionCompanyQuery(r.db.Model(&models.ProductConversion{}), companyID), offset, limit)
}
func (r *productConversionRepository) GetActiveConversions(offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(r.db.Model(&models.ProductConversion{}).Where("is_active = ?", true), offset, limit)
}
func (r *productConversionRepository) GetActiveConversionsByCompany(companyID uint, offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(conversionCompanyQuery(r.db.Model(&models.ProductConversion{}), companyID).
		Where("product_conversions.is_active = ?", true), offset, limit)
}
func (r *productConversionRepository) GetByRawProductID(id string, offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(r.db.Model(&models.ProductConversion{}).Where("raw_product_id = ?", id), offset, limit)
}
func (r *productConversionRepository) GetByRawProductIDAndCompany(id string, companyID uint, offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(conversionCompanyQuery(r.db.Model(&models.ProductConversion{}), companyID).
		Where("product_conversions.raw_product_id = ?", id), offset, limit)
}
func (r *productConversionRepository) GetByFinishedProductID(id string, offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(r.db.Model(&models.ProductConversion{}).Where("finished_product_id = ?", id), offset, limit)
}
func (r *productConversionRepository) GetByFinishedProductIDAndCompany(id string, companyID uint, offset, limit int) ([]models.ProductConversion, int64, error) {
	return r.list(conversionCompanyQuery(r.db.Model(&models.ProductConversion{}), companyID).
		Where("product_conversions.finished_product_id = ?", id), offset, limit)
}
func (r *productConversionRepository) GetByProductPair(rawID, finishedID string) (*models.ProductConversion, error) {
	var c models.ProductConversion
	err := r.db.Where("raw_product_id = ? AND finished_product_id = ?", rawID, finishedID).First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &c, err
}
func (r *productConversionRepository) GetByProductPairAndCompany(rawID, finishedID string, companyID uint) (*models.ProductConversion, error) {
	var c models.ProductConversion
	q := conversionCompanyQuery(r.db.Model(&models.ProductConversion{}), companyID)
	err := q.Select("product_conversions.*").
		Where("product_conversions.raw_product_id = ? AND product_conversions.finished_product_id = ?", rawID, finishedID).
		First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &c, err
}
func (r *productConversionRepository) Update(c *models.ProductConversion) error {
	return r.db.Save(c).Error
}
func (r *productConversionRepository) UpdateByCompany(c *models.ProductConversion, companyID uint) error {
	if c == nil {
		return gorm.ErrInvalidData
	}
	if _, err := r.GetByIDAndCompany(c.ID, companyID); err != nil {
		return err
	}
	return r.db.Save(c).Error
}
func (r *productConversionRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductConversion{}, "id = ?", id).Error
}
func (r *productConversionRepository) DeleteByCompany(id string, companyID uint) error {
	if _, err := r.GetByIDAndCompany(id, companyID); err != nil {
		return err
	}
	result := r.db.Delete(&models.ProductConversion{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *productConversionRepository) GetDB() *gorm.DB { return r.db }

// --- ProductConversionRecord Repository ---

type productConversionRecordRepository struct{ db *gorm.DB }

func NewProductConversionRecordRepository(db *gorm.DB) ProductConversionRecordRepository {
	return &productConversionRecordRepository{db: db}
}

func conversionRecordCompanyQuery(db *gorm.DB, companyID uint) *gorm.DB {
	return db.Joins(`
		JOIN users record_creator
			ON record_creator.id =
				CAST(product_conversion_records.created_by AS UNSIGNED)
	`).Where("record_creator.company_id = ?", companyID)
}

func (r *productConversionRecordRepository) Create(x *models.ProductConversionRecord) error {
	return r.db.Create(x).Error
}
func (r *productConversionRecordRepository) CreateForCompany(x *models.ProductConversionRecord, companyID uint) error {
	if x == nil || x.CreatedBy == "" {
		return gorm.ErrInvalidData
	}
	var count int64
	if err := r.db.Table("users").
		Where("id = CAST(? AS UNSIGNED) AND company_id = ?", x.CreatedBy, companyID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.Create(x).Error
}
func (r *productConversionRecordRepository) GetByID(id string) (*models.ProductConversionRecord, error) {
	var x models.ProductConversionRecord
	if err := r.db.Where("id = ?", id).First(&x).Error; err != nil {
		return nil, err
	}
	return &x, nil
}
func (r *productConversionRecordRepository) GetByIDAndCompany(id string, companyID uint) (*models.ProductConversionRecord, error) {
	var x models.ProductConversionRecord
	q := conversionRecordCompanyQuery(r.db.Model(&models.ProductConversionRecord{}), companyID)
	if err := q.Select("product_conversion_records.*").
		Where("product_conversion_records.id = ?", id).First(&x).Error; err != nil {
		return nil, err
	}
	return &x, nil
}
func (r *productConversionRecordRepository) list(q *gorm.DB, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	var rows []models.ProductConversionRecord
	var total int64
	if err := q.Distinct("product_conversion_records.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Select("product_conversion_records.*").Order("product_conversion_records.created_at DESC").
		Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
func (r *productConversionRecordRepository) GetAll(offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(r.db.Model(&models.ProductConversionRecord{}), offset, limit)
}
func (r *productConversionRecordRepository) GetAllByCompany(companyID uint, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(conversionRecordCompanyQuery(r.db.Model(&models.ProductConversionRecord{}), companyID), offset, limit)
}
func (r *productConversionRecordRepository) GetByConversionID(id string, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(r.db.Model(&models.ProductConversionRecord{}).Where("conversion_id = ?", id), offset, limit)
}
func (r *productConversionRecordRepository) GetByConversionIDAndCompany(id string, companyID uint, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(conversionRecordCompanyQuery(r.db.Model(&models.ProductConversionRecord{}), companyID).
		Where("product_conversion_records.conversion_id = ?", id), offset, limit)
}
func (r *productConversionRecordRepository) GetByDateRange(from, to time.Time, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(r.db.Model(&models.ProductConversionRecord{}).Where("conversion_date BETWEEN ? AND ?", from, to), offset, limit)
}
func (r *productConversionRecordRepository) GetByDateRangeAndCompany(from, to time.Time, companyID uint, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(conversionRecordCompanyQuery(r.db.Model(&models.ProductConversionRecord{}), companyID).
		Where("product_conversion_records.conversion_date BETWEEN ? AND ?", from, to), offset, limit)
}
func (r *productConversionRecordRepository) GetByStatus(status string, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(r.db.Model(&models.ProductConversionRecord{}).Where("status = ?", status), offset, limit)
}
func (r *productConversionRecordRepository) GetByStatusAndCompany(status string, companyID uint, offset, limit int) ([]models.ProductConversionRecord, int64, error) {
	return r.list(conversionRecordCompanyQuery(r.db.Model(&models.ProductConversionRecord{}), companyID).
		Where("product_conversion_records.status = ?", status), offset, limit)
}
func (r *productConversionRecordRepository) Update(x *models.ProductConversionRecord) error {
	return r.db.Save(x).Error
}
func (r *productConversionRecordRepository) UpdateByCompany(x *models.ProductConversionRecord, companyID uint) error {
	if x == nil {
		return gorm.ErrInvalidData
	}
	if _, err := r.GetByIDAndCompany(x.ID, companyID); err != nil {
		return err
	}
	return r.db.Save(x).Error
}
func (r *productConversionRecordRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductConversionRecord{}, "id = ?", id).Error
}
func (r *productConversionRecordRepository) DeleteByCompany(id string, companyID uint) error {
	if _, err := r.GetByIDAndCompany(id, companyID); err != nil {
		return err
	}
	return r.db.Delete(&models.ProductConversionRecord{}, "id = ?", id).Error
}
func (r *productConversionRecordRepository) GetDB() *gorm.DB { return r.db }

// --- ConversionRecordBagUsage Repository ---

type conversionRecordBagUsageRepository struct{ db *gorm.DB }

func NewConversionRecordBagUsageRepository(db *gorm.DB) ConversionRecordBagUsageRepository {
	return &conversionRecordBagUsageRepository{db: db}
}

func (r *conversionRecordBagUsageRepository) Create(x *models.ConversionRecordBagUsage) error {
	return r.db.Create(x).Error
}
func (r *conversionRecordBagUsageRepository) CreateForCompany(x *models.ConversionRecordBagUsage, companyID uint) error {
	if x == nil {
		return gorm.ErrInvalidData
	}
	var count int64
	if err := conversionRecordCompanyQuery(r.db.Model(&models.ProductConversionRecord{}), companyID).
		Where("product_conversion_records.id = ?", x.ConversionRecordID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.Create(x).Error
}
func (r *conversionRecordBagUsageRepository) GetByConversionRecordID(id string) ([]models.ConversionRecordBagUsage, error) {
	var rows []models.ConversionRecordBagUsage
	err := r.db.Where("conversion_record_id = ?", id).Order("created_at ASC").Find(&rows).Error
	return rows, err
}
func (r *conversionRecordBagUsageRepository) GetByConversionRecordIDAndCompany(id string, companyID uint) ([]models.ConversionRecordBagUsage, error) {
	var rows []models.ConversionRecordBagUsage
	err := r.db.Model(&models.ConversionRecordBagUsage{}).
		Joins("JOIN product_conversion_records ON product_conversion_records.id = conversion_record_bag_usages.conversion_record_id").
		Joins(`JOIN users record_creator ON record_creator.id = CAST(product_conversion_records.created_by AS UNSIGNED)`).
		Where("conversion_record_bag_usages.conversion_record_id = ? AND record_creator.company_id = ?", id, companyID).
		Select("conversion_record_bag_usages.*").Order("conversion_record_bag_usages.created_at ASC").
		Find(&rows).Error
	return rows, err
}
func (r *conversionRecordBagUsageRepository) GetByBagID(id string) ([]models.ConversionRecordBagUsage, error) {
	var rows []models.ConversionRecordBagUsage
	err := r.db.Where("bag_id = ?", id).Order("created_at DESC").Find(&rows).Error
	return rows, err
}
func (r *conversionRecordBagUsageRepository) GetByBagIDAndCompany(id string, companyID uint) ([]models.ConversionRecordBagUsage, error) {
	var rows []models.ConversionRecordBagUsage
	err := r.db.Model(&models.ConversionRecordBagUsage{}).
		Joins("JOIN product_conversion_records ON product_conversion_records.id = conversion_record_bag_usages.conversion_record_id").
		Joins(`JOIN users record_creator ON record_creator.id = CAST(product_conversion_records.created_by AS UNSIGNED)`).
		Where("conversion_record_bag_usages.bag_id = ? AND record_creator.company_id = ?", id, companyID).
		Select("conversion_record_bag_usages.*").Order("conversion_record_bag_usages.created_at DESC").
		Find(&rows).Error
	return rows, err
}
func (r *conversionRecordBagUsageRepository) GetDB() *gorm.DB { return r.db }
