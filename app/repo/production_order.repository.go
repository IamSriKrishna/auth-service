package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type productionOrderRepository struct {
	db *gorm.DB
}

func NewProductionOrderRepository(db *gorm.DB) ProductionOrderRepository {
	return &productionOrderRepository{db: db}
}

func (r *productionOrderRepository) Create(order *models.ProductionOrder) error {
	return r.db.Create(order).Error
}

func (r *productionOrderRepository) FindByID(id string) (*models.ProductionOrder, error) {
	var order models.ProductionOrder
	err := r.db.
		Preload("ItemGroup").
		Preload("ProductionOrderItems").
		Preload("ProductionOrderItems.ItemGroupComponent").
		Preload("ProductionOrderItems.ItemGroupComponent.Item").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *productionOrderRepository) FindAll(limit, offset int) ([]models.ProductionOrder, int64, error) {
	var orders []models.ProductionOrder
	var count int64

	query := r.db.Model(&models.ProductionOrder{})

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("ItemGroup").Limit(limit).Offset(offset).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *productionOrderRepository) Update(order *models.ProductionOrder) error {
	return r.db.Save(order).Error
}

func (r *productionOrderRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductionOrder{}, "id = ?", id).Error
}

func (r *productionOrderRepository) FindByProductionOrderNumber(orderNo string) (*models.ProductionOrder, error) {
	var order models.ProductionOrder
	err := r.db.First(&order, "production_order_number = ?", orderNo).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
