package repo

import (
	"fmt"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type PackageRepository interface {
	Create(pkg *models.Package) (*models.Package, error)
	FindByID(id string) (*models.Package, error)
	FindAll(limit, offset int) ([]models.Package, int64, error)
	FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Package, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.Package, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.Package, int64, error)
	Update(id string, pkg *models.Package) (*models.Package, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetNextPackageSlipNo() (string, error)
	GetDB() *gorm.DB
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) Create(pkg *models.Package) (*models.Package, error) {
	if err := r.db.Create(pkg).Error; err != nil {
		return nil, err
	}
	return pkg, nil
}

func (r *packageRepository) FindByID(id string) (*models.Package, error) {
	var pkg models.Package
	if err := r.db.
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.Variant").
		Preload("Items.SalesOrderItem").
		Where("id = ?", id).
		First(&pkg).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) FindAll(limit, offset int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	query := r.db.
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.Variant")

	if err := query.Model(&models.Package{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *packageRepository) FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	query := r.db.
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.Variant").
		Where("sales_order_id = ?", salesOrderID)

	if err := query.Model(&models.Package{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *packageRepository) FindByCustomer(customerID uint, limit, offset int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	query := r.db.
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.Variant").
		Where("customer_id = ?", customerID)

	if err := query.Model(&models.Package{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *packageRepository) FindByStatus(status string, limit, offset int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64

	query := r.db.
		Preload("SalesOrder").
		Preload("SalesOrder.Customer").
		Preload("Customer").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.Variant").
		Where("status = ?", status)

	if err := query.Model(&models.Package{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *packageRepository) Update(id string, pkg *models.Package) (*models.Package, error) {
	if err := r.db.Model(&models.Package{}).Where("id = ?", id).Updates(pkg).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *packageRepository) Delete(id string) error {
	return r.db.Delete(&models.Package{}, "id = ?", id).Error
}

func (r *packageRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Package{}).Where("id = ?", id).Update("status", status).Error
}

func (r *packageRepository) GetNextPackageSlipNo() (string, error) {
	var count int64
	if err := r.db.Model(&models.Package{}).Count(&count).Error; err != nil {
		return "", err
	}

	slipNo := generatePackageSlipNo(count + 1)
	return slipNo, nil
}

func (r *packageRepository) GetDB() *gorm.DB {
	return r.db
}

func generatePackageSlipNo(number int64) string {
	return fmt.Sprintf("PKG-%05d", number)
}
