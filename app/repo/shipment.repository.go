package repo

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type shipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) shipmentPreloads(
	db *gorm.DB,
) *gorm.DB {
	return db.
		Preload("Package").
		Preload("Package.Items").
		Preload("Package.Items.Product").
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer")
}

func (r *shipmentRepository) Create(
	shipment *models.Shipment,
) (*models.Shipment, error) {
	if shipment == nil {
		return nil, fmt.Errorf("shipment cannot be nil")
	}

	if err := r.db.
		Omit("Package", "SalesOrder", "Customer").
		Create(shipment).
		Error; err != nil {
		return nil, err
	}

	return r.FindByID(shipment.ID)
}

func (r *shipmentRepository) FindByID(
	id string,
) (*models.Shipment, error) {
	var shipment models.Shipment

	if err := r.shipmentPreloads(r.db).
		Where("shipments.id = ?", id).
		First(&shipment).
		Error; err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *shipmentRepository) FindAll(
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.db.Model(&models.Shipment{}),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindByPackage(
	packageID string,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.db.
			Model(&models.Shipment{}).
			Where("shipments.package_id = ?", packageID),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindBySalesOrder(
	salesOrderID string,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.db.
			Model(&models.Shipment{}).
			Where(
				"shipments.sales_order_id = ?",
				salesOrderID,
			),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindByCustomer(
	customerID uint,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.db.
			Model(&models.Shipment{}).
			Where(
				"shipments.customer_id = ?",
				customerID,
			),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindByStatus(
	status string,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.db.
			Model(&models.Shipment{}).
			Where("shipments.status = ?", status),
		limit,
		offset,
	)
}

func (r *shipmentRepository) findMany(
	query *gorm.DB,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	if err := query.
		Distinct("shipments.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	if err := r.shipmentPreloads(query).
		Select("shipments.*").
		Limit(limit).
		Offset(offset).
		Order("shipments.created_at DESC").
		Find(&shipments).
		Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) Update(
	id string,
	shipment *models.Shipment,
) (*models.Shipment, error) {
	if shipment == nil {
		return nil, fmt.Errorf("shipment cannot be nil")
	}

	result := r.db.
		Model(&models.Shipment{}).
		Where("id = ?", id).
		Omit("Package", "SalesOrder", "Customer").
		Updates(shipment)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.FindByID(id)
}

func (r *shipmentRepository) Delete(id string) error {
	result := r.db.
		Where("id = ?", id).
		Delete(&models.Shipment{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *shipmentRepository) UpdateStatus(
	id string,
	status string,
) error {
	result := r.db.
		Model(&models.Shipment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *shipmentRepository) GetNextShipmentNo() (
	string,
	error,
) {
	prefix := "SHP-" + time.Now().Format("20060102")

	var count int64
	if err := r.db.
		Model(&models.Shipment{}).
		Where("shipment_no LIKE ?", prefix+"%").
		Count(&count).
		Error; err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s-%04d",
		prefix,
		count+1,
	), nil
}

func (r *shipmentRepository) GetDB() *gorm.DB {
	return r.db
}

// -----------------------------------------------------------------------------
// Company-scoped methods.
// Company ownership is derived through customers.company_id.
// -----------------------------------------------------------------------------

func (r *shipmentRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.Shipment, error) {
	var shipment models.Shipment

	if err := r.shipmentPreloads(r.db).
		Joins(
			"INNER JOIN customers ON customers.id = shipments.customer_id",
		).
		Where(
			"shipments.id = ? AND customers.company_id = ?",
			id,
			companyID,
		).
		First(&shipment).
		Error; err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *shipmentRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.companyQuery(companyID),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindByPackageAndCompany(
	packageID string,
	companyID uint,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where(
				"shipments.package_id = ?",
				packageID,
			),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindBySalesOrderAndCompany(
	salesOrderID string,
	companyID uint,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where(
				"shipments.sales_order_id = ?",
				salesOrderID,
			),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindByCustomerAndCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where(
				"shipments.customer_id = ?",
				customerID,
			),
		limit,
		offset,
	)
}

func (r *shipmentRepository) FindByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]models.Shipment, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where("shipments.status = ?", status),
		limit,
		offset,
	)
}

func (r *shipmentRepository) companyQuery(
	companyID uint,
) *gorm.DB {
	return r.db.
		Model(&models.Shipment{}).
		Joins(
			"INNER JOIN customers ON customers.id = shipments.customer_id",
		).
		Where("customers.company_id = ?", companyID)
}

func (r *shipmentRepository) UpdateByCompany(
	id string,
	companyID uint,
	shipment *models.Shipment,
) (*models.Shipment, error) {
	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return nil, err
	}

	if _, err := r.Update(id, shipment); err != nil {
		return nil, err
	}

	return r.FindByIDAndCompany(id, companyID)
}

func (r *shipmentRepository) DeleteByCompany(
	id string,
	companyID uint,
) error {
	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return err
	}

	return r.Delete(id)
}

func (r *shipmentRepository) UpdateStatusByCompany(
	id string,
	companyID uint,
	status string,
) error {
	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return err
	}

	return r.UpdateStatus(id, status)
}
