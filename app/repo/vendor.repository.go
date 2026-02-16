package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type vendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) VendorRepository {
	return &vendorRepository{db: db}
}

func (r *vendorRepository) Create(vendor *models.Vendor) error {
	return r.db.Create(vendor).Error
}

func (r *vendorRepository) Update(vendor *models.Vendor) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(vendor).Error
}

func (r *vendorRepository) FindByID(id uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Preload("OtherDetails").
		Preload("BillingAddress", "address_type = ?", "billing").
		Preload("ShippingAddress", "address_type = ?", "shipping").
		Preload("ContactPersons").
		Preload("BankDetails").
		Preload("Documents").
		First(&vendor, id).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (r *vendorRepository) FindAll(page, limit int) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.Vendor{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("OtherDetails").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&vendors).Error

	if err != nil {
		return nil, 0, err
	}

	return vendors, total, nil
}

func (r *vendorRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("vendor_id = ?", id).Delete(&models.EntityOtherDetails{}).Error; err != nil {
			return err
		}
		if err := tx.Where("vendor_id = ?", id).Delete(&models.EntityAddress{}).Error; err != nil {
			return err
		}
		if err := tx.Where("vendor_id = ?", id).Delete(&models.EntityContactPerson{}).Error; err != nil {
			return err
		}
		if err := tx.Where("vendor_id = ?", id).Delete(&models.VendorBankDetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("vendor_id = ?", id).Delete(&models.EntityDocument{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Vendor{}, id).Error
	})
}

func (r *vendorRepository) FindByMobile(mobile string) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Where("mobile = ?", mobile).First(&vendor).Error

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}
