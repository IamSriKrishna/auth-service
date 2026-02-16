package repo

import (
	"fmt"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type ShipmentRepository interface {
	Create(shipment *models.Shipment) (*models.Shipment, error)
	FindByID(id string) (*models.Shipment, error)
	FindAll(limit, offset int) ([]models.Shipment, int64, error)
	FindByPackage(packageID string, limit, offset int) ([]models.Shipment, int64, error)
	FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Shipment, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.Shipment, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.Shipment, int64, error)
	Update(id string, shipment *models.Shipment) (*models.Shipment, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetNextShipmentNo() (string, error)
	GetDB() *gorm.DB
}

type shipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(shipment *models.Shipment) (*models.Shipment, error) {
	if err := r.db.Create(shipment).Error; err != nil {
		return nil, err
	}
	return shipment, nil
}

func (r *shipmentRepository) FindByID(id string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.
		Preload("Package").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Where("id = ?", id).
		First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) FindAll(limit, offset int) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	query := r.db.
		Preload("Package").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer")

	if err := query.Model(&models.Shipment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) FindByPackage(packageID string, limit, offset int) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	query := r.db.
		Preload("Package").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Where("package_id = ?", packageID)

	if err := query.Model(&models.Shipment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	query := r.db.
		Preload("Package").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Where("sales_order_id = ?", salesOrderID)

	if err := query.Model(&models.Shipment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) FindByCustomer(customerID uint, limit, offset int) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	query := r.db.
		Preload("Package").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Where("customer_id = ?", customerID)

	if err := query.Model(&models.Shipment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) FindByStatus(status string, limit, offset int) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	query := r.db.
		Preload("Package").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Where("status = ?", status)

	if err := query.Model(&models.Shipment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) Update(id string, shipment *models.Shipment) (*models.Shipment, error) {
	if err := r.db.Model(&models.Shipment{}).Where("id = ?", id).Updates(shipment).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *shipmentRepository) Delete(id string) error {
	return r.db.Delete(&models.Shipment{}, "id = ?", id).Error
}

func (r *shipmentRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Shipment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *shipmentRepository) GetNextShipmentNo() (string, error) {
	var count int64
	if err := r.db.Model(&models.Shipment{}).Count(&count).Error; err != nil {
		return "", err
	}

	shipmentNo := fmt.Sprintf("SHP-%05d", count+1)
	return shipmentNo, nil
}

func (r *shipmentRepository) GetDB() *gorm.DB {
	return r.db
}
