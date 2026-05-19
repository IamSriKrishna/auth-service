package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type customerPaymentRepository struct {
	db *gorm.DB
}

func NewCustomerPaymentRepository(db *gorm.DB) CustomerPaymentRepository {
	return &customerPaymentRepository{db: db}
}

func (r *customerPaymentRepository) Create(cp *models.CustomerPayment) (*models.CustomerPayment, error) {
	if err := r.db.Create(cp).Error; err != nil {
		return nil, err
	}
	return cp, nil
}

func (r *customerPaymentRepository) FindByID(id uint) (*models.CustomerPayment, error) {
	var cp models.CustomerPayment
	if err := r.db.
		Preload("SalesOrder").
		Preload("Customer").
		Where("id = ?", id).
		First(&cp).Error; err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *customerPaymentRepository) FindByPaymentNumber(paymentNumber string) (*models.CustomerPayment, error) {
	var cp models.CustomerPayment
	if err := r.db.
		Preload("SalesOrder").
		Preload("Customer").
		Where("payment_number = ?", paymentNumber).
		First(&cp).Error; err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *customerPaymentRepository) FindBySalesOrderID(soID string, limit, offset int) ([]models.CustomerPayment, int64, error) {
	var cps []models.CustomerPayment
	var total int64

	if err := r.db.Model(&models.CustomerPayment{}).
		Where("sales_order_id = ?", soID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("sales_order_id = ?", soID).
		Preload("SalesOrder").
		Preload("Customer").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&cps).Error; err != nil {
		return nil, 0, err
	}

	return cps, total, nil
}

func (r *customerPaymentRepository) FindByCustomerID(customerID uint, limit, offset int) ([]models.CustomerPayment, int64, error) {
	var cps []models.CustomerPayment
	var total int64

	if err := r.db.Model(&models.CustomerPayment{}).
		Where("customer_id = ?", customerID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("customer_id = ?", customerID).
		Preload("SalesOrder").
		Preload("Customer").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&cps).Error; err != nil {
		return nil, 0, err
	}

	return cps, total, nil
}

func (r *customerPaymentRepository) FindAll(limit, offset int) ([]models.CustomerPayment, int64, error) {
	var cps []models.CustomerPayment
	var total int64

	if err := r.db.Model(&models.CustomerPayment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Preload("SalesOrder").
		Preload("Customer").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&cps).Error; err != nil {
		return nil, 0, err
	}

	return cps, total, nil
}

func (r *customerPaymentRepository) FindByPaymentStatus(status string, limit, offset int) ([]models.CustomerPayment, int64, error) {
	var cps []models.CustomerPayment
	var total int64

	if err := r.db.Model(&models.CustomerPayment{}).
		Where("payment_status = ?", status).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("payment_status = ?", status).
		Preload("SalesOrder").
		Preload("Customer").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&cps).Error; err != nil {
		return nil, 0, err
	}

	return cps, total, nil
}

func (r *customerPaymentRepository) Update(id uint, cp *models.CustomerPayment) (*models.CustomerPayment, error) {
	if err := r.db.Model(&models.CustomerPayment{}).Where("id = ?", id).Updates(cp).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *customerPaymentRepository) UpdatePaymentStatus(id uint, status string, receivedAmount, remainingAmount float64) error {
	return r.db.Model(&models.CustomerPayment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"payment_status":   status,
		"received_amount":  receivedAmount,
		"remaining_amount": remainingAmount,
	}).Error
}

func (r *customerPaymentRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.CustomerPayment{}).Error
}

func (r *customerPaymentRepository) GetDB() *gorm.DB {
	return r.db
}
