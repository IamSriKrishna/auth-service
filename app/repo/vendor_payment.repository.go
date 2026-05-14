package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type vendorPaymentRepository struct {
	db *gorm.DB
}

func NewVendorPaymentRepository(db *gorm.DB) VendorPaymentRepository {
	return &vendorPaymentRepository{db: db}
}

func (r *vendorPaymentRepository) Create(vp *models.VendorPayment) (*models.VendorPayment, error) {
	if err := r.db.Create(vp).Error; err != nil {
		return nil, err
	}
	return vp, nil
}

func (r *vendorPaymentRepository) FindByID(id uint) (*models.VendorPayment, error) {
	var vp models.VendorPayment
	if err := r.db.
		Preload("PurchaseOrder").
		Preload("Vendor").
		Where("id = ?", id).
		First(&vp).Error; err != nil {
		return nil, err
	}
	return &vp, nil
}

func (r *vendorPaymentRepository) FindByPaymentNumber(paymentNumber string) (*models.VendorPayment, error) {
	var vp models.VendorPayment
	if err := r.db.
		Preload("PurchaseOrder").
		Preload("Vendor").
		Where("payment_number = ?", paymentNumber).
		First(&vp).Error; err != nil {
		return nil, err
	}
	return &vp, nil
}

func (r *vendorPaymentRepository) FindByPurchaseOrderID(poID string, limit, offset int) ([]models.VendorPayment, int64, error) {
	var vps []models.VendorPayment
	var total int64

	if err := r.db.Model(&models.VendorPayment{}).
		Where("purchase_order_id = ?", poID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("purchase_order_id = ?", poID).
		Preload("PurchaseOrder").
		Preload("Vendor").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&vps).Error; err != nil {
		return nil, 0, err
	}

	return vps, total, nil
}

func (r *vendorPaymentRepository) FindByVendorID(vendorID uint, limit, offset int) ([]models.VendorPayment, int64, error) {
	var vps []models.VendorPayment
	var total int64

	if err := r.db.Model(&models.VendorPayment{}).
		Where("vendor_id = ?", vendorID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("vendor_id = ?", vendorID).
		Preload("PurchaseOrder").
		Preload("Vendor").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&vps).Error; err != nil {
		return nil, 0, err
	}

	return vps, total, nil
}

func (r *vendorPaymentRepository) FindAll(limit, offset int) ([]models.VendorPayment, int64, error) {
	var vps []models.VendorPayment
	var total int64

	if err := r.db.Model(&models.VendorPayment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Preload("PurchaseOrder").
		Preload("Vendor").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&vps).Error; err != nil {
		return nil, 0, err
	}

	return vps, total, nil
}

func (r *vendorPaymentRepository) FindByPaymentStatus(status string, limit, offset int) ([]models.VendorPayment, int64, error) {
	var vps []models.VendorPayment
	var total int64

	if err := r.db.Model(&models.VendorPayment{}).
		Where("payment_status = ?", status).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("payment_status = ?", status).
		Preload("PurchaseOrder").
		Preload("Vendor").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&vps).Error; err != nil {
		return nil, 0, err
	}

	return vps, total, nil
}

func (r *vendorPaymentRepository) Update(id uint, vp *models.VendorPayment) (*models.VendorPayment, error) {
	if err := r.db.Model(&models.VendorPayment{}).Where("id = ?", id).Updates(vp).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *vendorPaymentRepository) UpdatePaymentStatus(id uint, status string, paidAmount, remainingAmount float64) error {
	return r.db.Model(&models.VendorPayment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"payment_status":   status,
		"paid_amount":      paidAmount,
		"remaining_amount": remainingAmount,
	}).Error
}

func (r *vendorPaymentRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.VendorPayment{}).Error
}

func (r *vendorPaymentRepository) GetDB() *gorm.DB {
	return r.db
}
