package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type PurchaseOrderRepository interface {
	Create(po *models.PurchaseOrder) (*models.PurchaseOrder, error)
	FindByID(id string) (*models.PurchaseOrder, error)
	FindAll(limit, offset int) ([]models.PurchaseOrder, int64, error)
	FindByVendor(vendorID uint, limit, offset int) ([]models.PurchaseOrder, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.PurchaseOrder, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.PurchaseOrder, int64, error)
	Update(id string, po *models.PurchaseOrder) (*models.PurchaseOrder, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetDB() *gorm.DB
}

type purchaseOrderRepository struct {
	db *gorm.DB
}

func NewPurchaseOrderRepository(db *gorm.DB) PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

func (r *purchaseOrderRepository) Create(po *models.PurchaseOrder) (*models.PurchaseOrder, error) {
	if err := r.db.Create(po).Error; err != nil {
		return nil, err
	}
	return po, nil
}

func (r *purchaseOrderRepository) FindByID(id string) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	if err := r.db.
		Preload("Vendor").
		Preload("Customer").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Where("id = ?", id).
		First(&po).Error; err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *purchaseOrderRepository) FindAll(limit, offset int) ([]models.PurchaseOrder, int64, error) {
	var pos []models.PurchaseOrder
	var total int64

	if err := r.db.Model(&models.PurchaseOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Preload("Vendor").
		Preload("Customer").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return pos, total, nil
}

func (r *purchaseOrderRepository) FindByVendor(vendorID uint, limit, offset int) ([]models.PurchaseOrder, int64, error) {
	var pos []models.PurchaseOrder
	var total int64

	if err := r.db.Model(&models.PurchaseOrder{}).
		Where("vendor_id = ?", vendorID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("vendor_id = ?", vendorID).
		Preload("Vendor").
		Preload("Customer").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return pos, total, nil
}

func (r *purchaseOrderRepository) FindByCustomer(customerID uint, limit, offset int) ([]models.PurchaseOrder, int64, error) {
	var pos []models.PurchaseOrder
	var total int64

	if err := r.db.Model(&models.PurchaseOrder{}).
		Where("customer_id = ?", customerID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("customer_id = ?", customerID).
		Preload("Vendor").
		Preload("Customer").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return pos, total, nil
}

func (r *purchaseOrderRepository) FindByStatus(status string, limit, offset int) ([]models.PurchaseOrder, int64, error) {
	var pos []models.PurchaseOrder
	var total int64

	if err := r.db.Model(&models.PurchaseOrder{}).
		Where("status = ?", status).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("status = ?", status).
		Preload("Vendor").
		Preload("Customer").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return pos, total, nil
}

func (r *purchaseOrderRepository) Update(id string, po *models.PurchaseOrder) (*models.PurchaseOrder, error) {
	if err := r.db.Model(&models.PurchaseOrder{}).Where("id = ?", id).Updates(po).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *purchaseOrderRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.PurchaseOrder{}).Error
}

func (r *purchaseOrderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.PurchaseOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *purchaseOrderRepository) GetDB() *gorm.DB {
	return r.db
}
