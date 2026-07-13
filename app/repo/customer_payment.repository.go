package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type customerPaymentRepository struct {
	db *gorm.DB
}

func NewCustomerPaymentRepository(
	db *gorm.DB,
) CustomerPaymentRepository {
	return &customerPaymentRepository{db: db}
}

func (r *customerPaymentRepository) customerPaymentPreloads(
	db *gorm.DB,
) *gorm.DB {
	return db.
		Preload("SalesOrder").
		Preload("Customer")
}

func customerPaymentCompanyQuery(
	db *gorm.DB,
	companyID uint,
) *gorm.DB {
	return db.Where(
		"customer_payments.created_by_company_id = ?",
		companyID,
	)
}

func (r *customerPaymentRepository) Create(
	payment *models.CustomerPayment,
) (*models.CustomerPayment, error) {
	if payment == nil {
		return nil, gorm.ErrInvalidData
	}

	if err := r.db.Create(payment).Error; err != nil {
		return nil, err
	}

	return r.FindByID(payment.ID)
}

func (r *customerPaymentRepository) FindByID(
	id uint,
) (*models.CustomerPayment, error) {
	var payment models.CustomerPayment

	err := r.customerPaymentPreloads(r.db).
		Where("customer_payments.id = ?", id).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *customerPaymentRepository) FindByIDAndCompany(
	id uint,
	companyID uint,
) (*models.CustomerPayment, error) {
	var payment models.CustomerPayment

	err := r.customerPaymentPreloads(
		customerPaymentCompanyQuery(
			r.db.Model(&models.CustomerPayment{}),
			companyID,
		),
	).
		Where("customer_payments.id = ?", id).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *customerPaymentRepository) FindByPaymentNumber(
	paymentNumber string,
) (*models.CustomerPayment, error) {
	var payment models.CustomerPayment

	err := r.customerPaymentPreloads(r.db).
		Where(
			"customer_payments.payment_number = ?",
			paymentNumber,
		).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *customerPaymentRepository) FindByPaymentNumberAndCompany(
	paymentNumber string,
	companyID uint,
) (*models.CustomerPayment, error) {
	var payment models.CustomerPayment

	err := r.customerPaymentPreloads(
		customerPaymentCompanyQuery(
			r.db.Model(&models.CustomerPayment{}),
			companyID,
		),
	).
		Where(
			"customer_payments.payment_number = ?",
			paymentNumber,
		).
		First(&payment).
		Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *customerPaymentRepository) findList(
	query *gorm.DB,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	var payments []models.CustomerPayment
	var total int64

	if err := query.
		Distinct("customer_payments.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	err := r.customerPaymentPreloads(query).
		Select("customer_payments.*").
		Order("customer_payments.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).
		Error

	return payments, total, err
}

func (r *customerPaymentRepository) FindAll(
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.CustomerPayment{}),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		customerPaymentCompanyQuery(
			r.db.Model(&models.CustomerPayment{}),
			companyID,
		),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindBySalesOrderID(
	salesOrderID string,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.CustomerPayment{}).
			Where(
				"customer_payments.sales_order_id = ?",
				salesOrderID,
			),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindBySalesOrderIDAndCompany(
	salesOrderID string,
	companyID uint,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		customerPaymentCompanyQuery(
			r.db.Model(&models.CustomerPayment{}),
			companyID,
		).Where(
			"customer_payments.sales_order_id = ?",
			salesOrderID,
		),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindByCustomerID(
	customerID uint,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.CustomerPayment{}).
			Where(
				"customer_payments.customer_id = ?",
				customerID,
			),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindByCustomerIDAndCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		customerPaymentCompanyQuery(
			r.db.Model(&models.CustomerPayment{}),
			companyID,
		).Where(
			"customer_payments.customer_id = ?",
			customerID,
		),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindByPaymentStatus(
	status string,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		r.db.Model(&models.CustomerPayment{}).
			Where(
				"customer_payments.payment_status = ?",
				status,
			),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) FindByPaymentStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]models.CustomerPayment, int64, error) {
	return r.findList(
		customerPaymentCompanyQuery(
			r.db.Model(&models.CustomerPayment{}),
			companyID,
		).Where(
			"customer_payments.payment_status = ?",
			status,
		),
		limit,
		offset,
	)
}

func (r *customerPaymentRepository) Update(
	id uint,
	payment *models.CustomerPayment,
) (*models.CustomerPayment, error) {
	if payment == nil {
		return nil, gorm.ErrInvalidData
	}

	payment.ID = id

	if err := r.db.Save(payment).Error; err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *customerPaymentRepository) UpdateByCompany(
	id uint,
	companyID uint,
	payment *models.CustomerPayment,
) (*models.CustomerPayment, error) {
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

func (r *customerPaymentRepository) UpdatePaymentStatus(
	id uint,
	status string,
	receivedAmount float64,
	remainingAmount float64,
) error {
	result := r.db.Model(&models.CustomerPayment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status":   status,
			"received_amount":  receivedAmount,
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

func (r *customerPaymentRepository) UpdatePaymentStatusByCompany(
	id uint,
	companyID uint,
	status string,
	receivedAmount float64,
	remainingAmount float64,
) error {
	result := r.db.Model(&models.CustomerPayment{}).
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Updates(map[string]interface{}{
			"payment_status":   status,
			"received_amount":  receivedAmount,
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

func (r *customerPaymentRepository) Delete(
	id uint,
) error {
	result := r.db.Delete(
		&models.CustomerPayment{},
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

func (r *customerPaymentRepository) DeleteByCompany(
	id uint,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND created_by_company_id = ?",
			id,
			companyID,
		).
		Delete(&models.CustomerPayment{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *customerPaymentRepository) GetDB() *gorm.DB {
	return r.db
}
