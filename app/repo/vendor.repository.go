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

func (r *vendorRepository) Create(
	vendor *models.Vendor,
) error {
	return r.db.Create(vendor).Error
}

func (r *vendorRepository) Update(
	vendor *models.Vendor,
) error {
	return r.db.Save(vendor).Error
}

func (r *vendorRepository) FindByID(
	id uint,
) (*models.Vendor, error) {
	var vendor models.Vendor

	err := r.db.
		Where("id = ?", id).
		First(&vendor).
		Error

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (r *vendorRepository) FindAll(
	page int,
	limit int,
) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := r.db.Model(&models.Vendor{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&vendors).
		Error

	return vendors, total, err
}

func (r *vendorRepository) FindByUserID(
	userID uint,
	companyID uint,
	page int,
	limit int,
) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := r.db.
		Model(&models.Vendor{}).
		Where(
			"user_id = ? AND company_id = ?",
			userID,
			companyID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&vendors).
		Error

	return vendors, total, err
}

func (r *vendorRepository) FindByIDAndUser(
	id uint,
	userID uint,
) (*models.Vendor, error) {
	var vendor models.Vendor

	err := r.db.
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		First(&vendor).
		Error

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (r *vendorRepository) Delete(
	id uint,
) error {
	return r.db.Delete(&models.Vendor{}, id).Error
}

func (r *vendorRepository) FindByMobile(
	mobile string,
) (*models.Vendor, error) {
	var vendor models.Vendor

	err := r.db.
		Where("mobile = ?", mobile).
		First(&vendor).
		Error

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (r *vendorRepository) FindByIDAndCompany(
	id uint,
	companyID uint,
) (*models.Vendor, error) {
	var vendor models.Vendor

	err := r.db.
		Preload("User").
		Preload("Company").
		Preload("OtherDetails").
		Preload("BillingAddress", "address_type = ?", "billing").
		Preload("ShippingAddress", "address_type = ?", "shipping").
		Preload("ContactPersons").
		Preload("BankDetails").
		Preload("BankDetails.Bank").
		Preload("Documents").
		Where(
			"vendors.id = ? AND vendors.company_id = ?",
			id,
			companyID,
		).
		First(&vendor).
		Error

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (r *vendorRepository) FindByCompanyID(
	companyID uint,
	page int,
	limit int,
) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := r.db.
		Model(&models.Vendor{}).
		Where("company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&vendors).
		Error

	return vendors, total, err
}

func (r *vendorRepository) UpdateByCompanyID(
	vendor *models.Vendor,
	companyID uint,
) error {
	if vendor == nil {
		return gorm.ErrInvalidData
	}

	result := r.db.
		Model(&models.Vendor{}).
		Where(
			"id = ? AND company_id = ?",
			vendor.ID,
			companyID,
		).
		Updates(vendor)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *vendorRepository) DeleteByIDAndCompany(
	id uint,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND company_id = ?",
			id,
			companyID,
		).
		Delete(&models.Vendor{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *vendorRepository) FindByMobileAndCompany(
	mobile string,
	companyID uint,
) (*models.Vendor, error) {
	var vendor models.Vendor

	err := r.db.
		Where(
			"mobile = ? AND company_id = ?",
			mobile,
			companyID,
		).
		First(&vendor).
		Error

	if err != nil {
		return nil, err
	}

	return &vendor, nil
}
