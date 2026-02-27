package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (r *billRepository) Create(bill *models.Bill) (*models.Bill, error) {
	if err := r.db.Create(bill).Error; err != nil {
		return nil, err
	}
	return bill, nil
}

func (r *billRepository) FindByID(id string) (*models.Bill, error) {
	var bill models.Bill
	if err := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Where("id = ?", id).
		First(&bill).Error; err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *billRepository) FindAll(limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (r *billRepository) FindByVendor(vendorID uint, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("vendor_id = ?", vendorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("vendor_id = ?", vendorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (r *billRepository) FindByStatus(status string, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (r *billRepository) Update(id string, bill *models.Bill) (*models.Bill, error) {
	if err := r.db.Model(&models.Bill{}).Where("id = ?", id).Updates(bill).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *billRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Bill{}).Error
}

func (r *billRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Bill{}).Where("id = ?", id).Update("status", status).Error
}

func (r *billRepository) GetDB() *gorm.DB {
	return r.db
}
