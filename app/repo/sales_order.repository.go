package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type salesOrderRepository struct {
	db *gorm.DB
}

func NewSalesOrderRepository(db *gorm.DB) SalesOrderRepository {
	return &salesOrderRepository{db: db}
}

func (r *salesOrderRepository) Create(so *models.SalesOrder) (*models.SalesOrder, error) {
	if err := r.db.Create(so).Error; err != nil {
		return nil, err
	}
	// Reload the sales order with all relationships to ensure everything is populated
	return r.FindByID(so.ID)
}

func (r *salesOrderRepository) FindByID(id string) (*models.SalesOrder, error) {
	var so models.SalesOrder
	if err := r.db.
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Item.ItemDetails").
		Preload("LineItems.Variant").
		Preload("LineItems.Variant.Attributes").
		Where("id = ?", id).
		First(&so).Error; err != nil {
		return nil, err
	}
	return &so, nil
}

func (r *salesOrderRepository) FindAll(limit, offset int) ([]models.SalesOrder, int64, error) {
	var sos []models.SalesOrder
	var total int64

	query := r.db.
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Item.ItemDetails").
		Preload("LineItems.Variant").
		Preload("LineItems.Variant.Attributes")

	if err := query.Model(&models.SalesOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sos).Error; err != nil {
		return nil, 0, err
	}

	return sos, total, nil
}

func (r *salesOrderRepository) FindByCustomer(customerID uint, limit, offset int) ([]models.SalesOrder, int64, error) {
	var sos []models.SalesOrder
	var total int64

	query := r.db.
		Where("customer_id = ?", customerID).
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.SalesOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sos).Error; err != nil {
		return nil, 0, err
	}

	return sos, total, nil
}

func (r *salesOrderRepository) FindByStatus(status string, limit, offset int) ([]models.SalesOrder, int64, error) {
	var sos []models.SalesOrder
	var total int64

	query := r.db.
		Where("status = ?", status).
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.SalesOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sos).Error; err != nil {
		return nil, 0, err
	}

	return sos, total, nil
}

func (r *salesOrderRepository) Update(id string, so *models.SalesOrder) (*models.SalesOrder, error) {
	if err := r.db.Model(&models.SalesOrder{}).Where("id = ?", id).
		Updates(so).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *salesOrderRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.SalesOrder{}).Error
}

func (r *salesOrderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.SalesOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *salesOrderRepository) GetDB() *gorm.DB {
	return r.db
}
