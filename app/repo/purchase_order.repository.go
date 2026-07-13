package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

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
		Preload("LineItems.Product").
		Where("id = ?", id).
		First(&po).Error; err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *purchaseOrderRepository) purchaseOrderPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Vendor").
		Preload("Customer").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Product")
}

func (r *purchaseOrderRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder

	err := r.purchaseOrderPreloads(r.db).
		Where(
			"purchase_orders.id = ? AND purchase_orders.created_by_company_id = ?",
			id,
			companyID,
		).
		First(&purchaseOrder).
		Error

	if err != nil {
		return nil, err
	}

	return &purchaseOrder, nil
}

func (r *purchaseOrderRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.PurchaseOrder, int64, error) {
	var purchaseOrders []models.PurchaseOrder
	var total int64

	query := r.db.
		Model(&models.PurchaseOrder{}).
		Where("created_by_company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.purchaseOrderPreloads(r.db).
		Where("created_by_company_id = ?", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&purchaseOrders).
		Error

	return purchaseOrders, total, err
}

func (r *purchaseOrderRepository) FindByVendorAndCompany(
	vendorID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.PurchaseOrder, int64, error) {
	var purchaseOrders []models.PurchaseOrder
	var total int64

	query := r.db.
		Model(&models.PurchaseOrder{}).
		Where(
			"vendor_id = ? AND created_by_company_id = ?",
			vendorID,
			companyID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.purchaseOrderPreloads(r.db).
		Where(
			"vendor_id = ? AND created_by_company_id = ?",
			vendorID,
			companyID,
		).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&purchaseOrders).
		Error

	return purchaseOrders, total, err
}

func (r *purchaseOrderRepository) FindByCustomerAndCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.PurchaseOrder, int64, error) {
	var purchaseOrders []models.PurchaseOrder
	var total int64

	query := r.db.
		Model(&models.PurchaseOrder{}).
		Where(
			"customer_id = ? AND created_by_company_id = ?",
			customerID,
			companyID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.purchaseOrderPreloads(r.db).
		Where(
			"customer_id = ? AND created_by_company_id = ?",
			customerID,
			companyID,
		).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&purchaseOrders).
		Error

	return purchaseOrders, total, err
}

func (r *purchaseOrderRepository) FindByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]models.PurchaseOrder, int64, error) {
	var purchaseOrders []models.PurchaseOrder
	var total int64

	query := r.db.
		Model(&models.PurchaseOrder{}).
		Where(
			"status = ? AND created_by_company_id = ?",
			status,
			companyID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.purchaseOrderPreloads(r.db).
		Where(
			"status = ? AND created_by_company_id = ?",
			status,
			companyID,
		).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&purchaseOrders).
		Error

	return purchaseOrders, total, err
}

func (r *purchaseOrderRepository) UpdateByCompany(
	id string,
	companyID uint,
	purchaseOrder *models.PurchaseOrder,
) (*models.PurchaseOrder, error) {
	if purchaseOrder == nil {
		return nil, gorm.ErrInvalidData
	}

	var existing models.PurchaseOrder
	err := r.db.
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		First(&existing).
		Error
	if err != nil {
		return nil, err
	}

	updated, err := r.Update(id, purchaseOrder)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *purchaseOrderRepository) DeleteByCompany(
	id string,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Delete(&models.PurchaseOrder{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *purchaseOrderRepository) UpdateStatusByCompany(
	id string,
	companyID uint,
	status string,
) error {
	result := r.db.
		Model(&models.PurchaseOrder{}).
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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
		Preload("LineItems.Product").
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
		Preload("LineItems.Product").
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
		Preload("LineItems.Product").
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
		Preload("LineItems.Product").
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
