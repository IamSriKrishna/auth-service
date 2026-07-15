package repo

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) packagePreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("SalesOrder.LineItems").
		Preload("Customer").
		Preload("Items").
		Preload("Items.SalesOrderItem").
		Preload("Items.Product").
		Preload("Items.Product.ProductDetails").
		Preload("Items.Variant")
}

func (r *packageRepository) Create(pkg *models.Package) (*models.Package, error) {
	if pkg == nil {
		return nil, fmt.Errorf("package cannot be nil")
	}
	if err := r.db.Omit("SalesOrder", "Customer").Create(pkg).Error; err != nil {
		return nil, err
	}
	return r.FindByID(pkg.ID)
}

func (r *packageRepository) FindByID(id string) (*models.Package, error) {
	var pkg models.Package
	if err := r.packagePreloads(r.db).
		Where("packages.id = ?", id).
		First(&pkg).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) FindAll(limit, offset int) ([]models.Package, int64, error) {
	return r.findMany(
		r.db.Model(&models.Package{}),
		limit,
		offset,
	)
}

func (r *packageRepository) FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Package, int64, error) {
	return r.findMany(
		r.db.Model(&models.Package{}).
			Where("packages.sales_order_id = ?", salesOrderID),
		limit,
		offset,
	)
}

func (r *packageRepository) FindByCustomer(customerID uint, limit, offset int) ([]models.Package, int64, error) {
	return r.findMany(
		r.db.Model(&models.Package{}).
			Where("packages.customer_id = ?", customerID),
		limit,
		offset,
	)
}

func (r *packageRepository) FindByStatus(status string, limit, offset int) ([]models.Package, int64, error) {
	return r.findMany(
		r.db.Model(&models.Package{}).
			Where("packages.status = ?", status),
		limit,
		offset,
	)
}

func (r *packageRepository) findMany(query *gorm.DB, limit, offset int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	if err := query.Distinct("packages.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.packagePreloads(query).
		Select("packages.*").
		Limit(limit).
		Offset(offset).
		Order("packages.created_at DESC").
		Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *packageRepository) Update(id string, pkg *models.Package) (*models.Package, error) {
	if pkg == nil {
		return nil, fmt.Errorf("package cannot be nil")
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Package
		if err := tx.Where("id = ?", id).First(&existing).Error; err != nil {
			return err
		}

		if err := tx.Model(&existing).
			Omit("Items", "SalesOrder", "Customer").
			Updates(pkg).Error; err != nil {
			return err
		}

		for i := range pkg.Items {
			if err := tx.Model(&models.PackageItem{}).
				Where(
					"package_id = ? AND sales_order_item_id = ?",
					id,
					pkg.Items[i].SalesOrderItemID,
				).
				Updates(map[string]interface{}{
					"packed_qty": pkg.Items[i].PackedQty,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *packageRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("package_id = ?", id).
			Delete(&models.PackageItem{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).
			Delete(&models.Package{}).Error
	})
}

func (r *packageRepository) UpdateStatus(id, status string) error {
	return r.db.Model(&models.Package{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *packageRepository) GetNextPackageSlipNo() (string, error) {
	prefix := "PKG-" + time.Now().Format("20060102")
	var count int64
	if err := r.db.Model(&models.Package{}).
		Where("package_slip_no LIKE ?", prefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%04d", prefix, count+1), nil
}

func (r *packageRepository) FindByIDAndCompany(id string, companyID uint) (*models.Package, error) {
	var pkg models.Package
	if err := r.packagePreloads(r.db).
		Joins("INNER JOIN customers ON customers.id = packages.customer_id").
		Where(
			"packages.id = ? AND customers.company_id = ?",
			id,
			companyID,
		).
		First(&pkg).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) FindAllByCompany(companyID uint, limit, offset int) ([]models.Package, int64, error) {
	return r.findMany(
		r.companyQuery(companyID),
		limit,
		offset,
	)
}

func (r *packageRepository) FindBySalesOrderAndCompany(
	salesOrderID string,
	companyID uint,
	limit,
	offset int,
) ([]models.Package, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where("packages.sales_order_id = ?", salesOrderID),
		limit,
		offset,
	)
}

func (r *packageRepository) FindByCustomerAndCompany(
	customerID,
	companyID uint,
	limit,
	offset int,
) ([]models.Package, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where("packages.customer_id = ?", customerID),
		limit,
		offset,
	)
}

func (r *packageRepository) FindByStatusAndCompany(
	status string,
	companyID uint,
	limit,
	offset int,
) ([]models.Package, int64, error) {
	return r.findMany(
		r.companyQuery(companyID).
			Where("packages.status = ?", status),
		limit,
		offset,
	)
}

func (r *packageRepository) companyQuery(companyID uint) *gorm.DB {
	return r.db.Model(&models.Package{}).
		Joins("INNER JOIN customers ON customers.id = packages.customer_id").
		Where("customers.company_id = ?", companyID)
}

func (r *packageRepository) UpdateByCompany(
	id string,
	companyID uint,
	pkg *models.Package,
) (*models.Package, error) {
	if _, err := r.FindByIDAndCompany(id, companyID); err != nil {
		return nil, err
	}
	if _, err := r.Update(id, pkg); err != nil {
		return nil, err
	}
	return r.FindByIDAndCompany(id, companyID)
}

func (r *packageRepository) DeleteByCompany(id string, companyID uint) error {
	if _, err := r.FindByIDAndCompany(id, companyID); err != nil {
		return err
	}
	return r.Delete(id)
}

func (r *packageRepository) UpdateStatusByCompany(id string, companyID uint, status string) error {
	if _, err := r.FindByIDAndCompany(id, companyID); err != nil {
		return err
	}
	return r.UpdateStatus(id, status)
}
