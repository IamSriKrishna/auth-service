package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type vendorPaymentRepository struct {
	db *gorm.DB
}

func NewVendorPaymentRepository(
	db *gorm.DB,
) VendorPaymentRepository {
	return &vendorPaymentRepository{db: db}
}

func (r *vendorPaymentRepository) vendorPaymentPreloads(
	db *gorm.DB,
) *gorm.DB {
	return db.
		Preload("PurchaseOrder").
		Preload("Vendor")
}

func vendorPaymentCompanyQuery(
	db *gorm.DB,
	companyID uint,
) *gorm.DB {
	return db.Where(
		"vendor_payments.created_by_company_id = ?",
		companyID,
	)
}

func (r *vendorPaymentRepository) Create(
	payment *models.VendorPayment,
) (*models.VendorPayment, error) {
	if payment == nil {
		return nil, gorm.ErrInvalidData
	}

	if err := r.db.Create(payment).Error; err != nil {
		return nil, err
	}

	return r.FindByID(payment.ID)
}

func (r *vendorPaymentRepository) FindByID(
	id uint,
) (*models.VendorPayment, error) {
	var payment models.VendorPayment

	err := r.vendorPaymentPreloads(r.db).
		Where("vendor_payments.id = ?", id).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *vendorPaymentRepository) FindByIDAndCompany(
	id uint,
	companyID uint,
) (*models.VendorPayment, error) {
	var payment models.VendorPayment

	err := r.vendorPaymentPreloads(
		vendorPaymentCompanyQuery(
			r.db.Model(&models.VendorPayment{}),
			companyID,
		),
	).
		Where("vendor_payments.id = ?", id).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *vendorPaymentRepository) FindByPaymentNumber(
	paymentNumber string,
) (*models.VendorPayment, error) {
	var payment models.VendorPayment

	err := r.vendorPaymentPreloads(r.db).
		Where(
			"vendor_payments.payment_number = ?",
			paymentNumber,
		).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *vendorPaymentRepository) FindByPaymentNumberAndCompany(
	paymentNumber string,
	companyID uint,
) (*models.VendorPayment, error) {
	var payment models.VendorPayment

	err := r.vendorPaymentPreloads(
		vendorPaymentCompanyQuery(
			r.db.Model(&models.VendorPayment{}),
			companyID,
		),
	).
		Where(
			"vendor_payments.payment_number = ?",
			paymentNumber,
		).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *vendorPaymentRepository) findList(
	query *gorm.DB,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	var payments []models.VendorPayment
	var total int64

	if err := query.
		Distinct("vendor_payments.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	err := r.vendorPaymentPreloads(query).
		Select("vendor_payments.*").
		Order("vendor_payments.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).
		Error

	return payments, total, err
}

func (r *vendorPaymentRepository) FindAll(
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.VendorPayment{}),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		vendorPaymentCompanyQuery(
			r.db.Model(&models.VendorPayment{}),
			companyID,
		),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindByPurchaseOrderID(
	purchaseOrderID string,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.VendorPayment{}).
			Where(
				"vendor_payments.purchase_order_id = ?",
				purchaseOrderID,
			),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindByPurchaseOrderIDAndCompany(
	purchaseOrderID string,
	companyID uint,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		vendorPaymentCompanyQuery(
			r.db.Model(&models.VendorPayment{}),
			companyID,
		).Where(
			"vendor_payments.purchase_order_id = ?",
			purchaseOrderID,
		),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindByVendorID(
	vendorID uint,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.VendorPayment{}).
			Where(
				"vendor_payments.vendor_id = ?",
				vendorID,
			),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindByVendorIDAndCompany(
	vendorID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		vendorPaymentCompanyQuery(
			r.db.Model(&models.VendorPayment{}),
			companyID,
		).Where(
			"vendor_payments.vendor_id = ?",
			vendorID,
		),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindByPaymentStatus(
	status string,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.VendorPayment{}).
			Where(
				"vendor_payments.payment_status = ?",
				status,
			),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) FindByPaymentStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]models.VendorPayment, int64, error) {
	return r.findList(
		vendorPaymentCompanyQuery(
			r.db.Model(&models.VendorPayment{}),
			companyID,
		).Where(
			"vendor_payments.payment_status = ?",
			status,
		),
		limit,
		offset,
	)
}

func (r *vendorPaymentRepository) Update(
	id uint,
	payment *models.VendorPayment,
) (*models.VendorPayment, error) {
	if payment == nil {
		return nil, gorm.ErrInvalidData
	}

	payment.ID = id

	if err := r.db.Save(payment).Error; err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *vendorPaymentRepository) UpdateByCompany(
	id uint,
	companyID uint,
	payment *models.VendorPayment,
) (*models.VendorPayment, error) {
	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, gorm.ErrInvalidData
	}

	payment.ID = id
	payment.CreatedByCompanyID = companyID

	if err := r.db.Save(payment).Error; err != nil {
		return nil, err
	}

	return r.FindByIDAndCompany(id, companyID)
}

func (r *vendorPaymentRepository) UpdatePaymentStatus(
	id uint,
	status string,
	paidAmount float64,
	remainingAmount float64,
) error {
	result := r.db.Model(&models.VendorPayment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status":   status,
			"paid_amount":      paidAmount,
			"remaining_amount": remainingAmount,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *vendorPaymentRepository) UpdatePaymentStatusByCompany(
	id uint,
	companyID uint,
	status string,
	paidAmount float64,
	remainingAmount float64,
) error {
	result := r.db.Model(&models.VendorPayment{}).
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Updates(map[string]interface{}{
			"payment_status":   status,
			"paid_amount":      paidAmount,
			"remaining_amount": remainingAmount,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *vendorPaymentRepository) Delete(
	id uint,
) error {
	result := r.db.Delete(
		&models.VendorPayment{},
		id,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *vendorPaymentRepository) DeleteByCompany(
	id uint,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Delete(&models.VendorPayment{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *vendorPaymentRepository) GetDB() *gorm.DB {
	return r.db
}
