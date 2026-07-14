package repo

import (
	"errors"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type salesOrderRepository struct {
	db *gorm.DB
}

func NewSalesOrderRepository(db *gorm.DB) SalesOrderRepository {
	return &salesOrderRepository{db: db}
}

func (r *salesOrderRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *salesOrderRepository) salesOrderPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Customer").
		Preload("Salesperson").
		Preload("LineItems").
		Preload("LineItems.Manufacturer")
}

func (r *salesOrderRepository) Create(so *models.SalesOrder) (*models.SalesOrder, error) {
	if so == nil {
		return nil, errors.New("sales order cannot be nil")
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Customer", "Salesperson", "LineItems.Manufacturer").Create(so).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.FindByID(so.ID)
}

func (r *salesOrderRepository) FindByID(id string) (*models.SalesOrder, error) {
	var so models.SalesOrder
	err := r.salesOrderPreloads(r.db).
		Where("sales_orders.id = ?", id).
		First(&so).Error
	if err != nil {
		return nil, err
	}
	return &so, nil
}

func (r *salesOrderRepository) FindAll(limit, offset int) ([]models.SalesOrder, int64, error) {
	var salesOrders []models.SalesOrder
	var total int64

	query := r.db.Model(&models.SalesOrder{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.salesOrderPreloads(r.db).
		Order("sales_orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&salesOrders).Error
	if err != nil {
		return nil, 0, err
	}

	return salesOrders, total, nil
}

func (r *salesOrderRepository) FindByCustomer(customerID uint, limit, offset int) ([]models.SalesOrder, int64, error) {
	var salesOrders []models.SalesOrder
	var total int64

	query := r.db.Model(&models.SalesOrder{}).
		Where("sales_orders.customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.salesOrderPreloads(query).
		Order("sales_orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&salesOrders).Error
	if err != nil {
		return nil, 0, err
	}

	return salesOrders, total, nil
}

func (r *salesOrderRepository) FindByStatus(status string, limit, offset int) ([]models.SalesOrder, int64, error) {
	var salesOrders []models.SalesOrder
	var total int64

	query := r.db.Model(&models.SalesOrder{}).
		Where("sales_orders.status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.salesOrderPreloads(query).
		Order("sales_orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&salesOrders).Error
	if err != nil {
		return nil, 0, err
	}

	return salesOrders, total, nil
}

func (r *salesOrderRepository) Update(id string, so *models.SalesOrder) (*models.SalesOrder, error) {
	if so == nil {
		return nil, errors.New("sales order cannot be nil")
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.SalesOrder
		if err := tx.Where("id = ?", id).First(&existing).Error; err != nil {
			return err
		}

		if err := tx.Model(&existing).
			Omit("Customer", "Salesperson", "LineItems").
			Updates(so).Error; err != nil {
			return err
		}

		if so.LineItems != nil {
			if err := tx.Where("sales_order_id = ?", id).
				Delete(&models.SalesOrderLineItem{}).Error; err != nil {
				return err
			}

			for i := range so.LineItems {
				so.LineItems[i].ID = 0
				so.LineItems[i].SalesOrderID = id
			}

			if len(so.LineItems) > 0 {
				if err := tx.Omit("Manufacturer").Create(&so.LineItems).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *salesOrderRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("sales_order_id = ?", id).
			Delete(&models.SalesOrderLineItem{}).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", id).Delete(&models.SalesOrder{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *salesOrderRepository) UpdateStatus(id string, status string) error {
	result := r.db.Model(&models.SalesOrder{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Company ownership is derived from the linked customer. This avoids requiring
// a company_id column on sales_orders while still isolating tenant data.
func (r *salesOrderRepository) companyQuery(companyID uint) *gorm.DB {
	return r.db.Model(&models.SalesOrder{}).
		Joins("JOIN customers ON customers.id = sales_orders.customer_id").
		Where("customers.company_id = ?", companyID)
}

func (r *salesOrderRepository) FindByIDAndCompany(id string, companyID uint) (*models.SalesOrder, error) {
	var so models.SalesOrder

	err := r.salesOrderPreloads(r.db).
		Joins("JOIN customers ON customers.id = sales_orders.customer_id").
		Where("sales_orders.id = ? AND customers.company_id = ?", id, companyID).
		First(&so).Error
	if err != nil {
		return nil, err
	}

	return &so, nil
}

func (r *salesOrderRepository) FindAllByCompany(companyID uint, limit, offset int) ([]models.SalesOrder, int64, error) {
	var salesOrders []models.SalesOrder
	var total int64

	query := r.companyQuery(companyID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.salesOrderPreloads(query).
		Select("sales_orders.*").
		Order("sales_orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&salesOrders).Error
	if err != nil {
		return nil, 0, err
	}

	return salesOrders, total, nil
}

func (r *salesOrderRepository) FindByCustomerAndCompany(customerID, companyID uint, limit, offset int) ([]models.SalesOrder, int64, error) {
	var salesOrders []models.SalesOrder
	var total int64

	query := r.companyQuery(companyID).
		Where("sales_orders.customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.salesOrderPreloads(query).
		Select("sales_orders.*").
		Order("sales_orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&salesOrders).Error
	if err != nil {
		return nil, 0, err
	}

	return salesOrders, total, nil
}

func (r *salesOrderRepository) FindByStatusAndCompany(status string, companyID uint, limit, offset int) ([]models.SalesOrder, int64, error) {
	var salesOrders []models.SalesOrder
	var total int64

	query := r.companyQuery(companyID).
		Where("sales_orders.status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.salesOrderPreloads(query).
		Select("sales_orders.*").
		Order("sales_orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&salesOrders).Error
	if err != nil {
		return nil, 0, err
	}

	return salesOrders, total, nil
}

func (r *salesOrderRepository) UpdateByCompany(id string, companyID uint, so *models.SalesOrder) (*models.SalesOrder, error) {
	if _, err := r.FindByIDAndCompany(id, companyID); err != nil {
		return nil, err
	}

	updated, err := r.Update(id, so)
	if err != nil {
		return nil, err
	}

	return r.FindByIDAndCompany(updated.ID, companyID)
}

func (r *salesOrderRepository) DeleteByCompany(id string, companyID uint) error {
	if _, err := r.FindByIDAndCompany(id, companyID); err != nil {
		return err
	}
	return r.Delete(id)
}

func (r *salesOrderRepository) UpdateStatusByCompany(id string, companyID uint, status string) error {
	if _, err := r.FindByIDAndCompany(id, companyID); err != nil {
		return err
	}
	return r.UpdateStatus(id, status)
}
