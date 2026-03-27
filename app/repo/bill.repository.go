package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (r *billRepository) Create(bill *models.Bill) (*models.Bill, error) {
	if err := r.db.Create(bill).Error; err != nil {
		return nil, err
	}
	return bill, nil
}

func (r *billRepository) FindByID(id string) (*models.Bill, error) {
	var bill models.Bill
	if err := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant").
		Where("id = ?", id).
		First(&bill).Error; err != nil {
		return nil, err
	}

	// Populate user information from joined query
	if bill.CreatedBy != "" {
		var userInfo struct {
			UserName    string
			CompanyID   uint
			CompanyName string
		}

		if err := r.db.
			Select("u.email as user_name, c.id as company_id, c.company_name as company_name").
			Table("bills").
			Joins("LEFT JOIN users u ON u.id = CAST(bills.created_by AS UNSIGNED)").
			Joins("LEFT JOIN companies c ON c.id = u.company_id").
			Where("bills.id = ?", id).
			Scan(&userInfo).Error; err == nil {
			bill.CreatedByUserName = userInfo.UserName
			bill.CreatedByCompanyID = userInfo.CompanyID
			bill.CreatedByCompanyName = userInfo.CompanyName
		}
	}

	return &bill, nil
}

func (r *billRepository) FindAll(limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)
	return bills, total, nil
}

func (r *billRepository) FindByCreatedBy(createdBy string, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("created_by = ?", createdBy).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("created_by = ?", createdBy).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)
	return bills, total, nil
}

func (r *billRepository) FindByVendor(vendorID uint, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("vendor_id = ?", vendorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("vendor_id = ?", vendorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)
	return bills, total, nil
}

func (r *billRepository) FindByVendorAndCreatedBy(vendorID uint, createdBy string, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("vendor_id = ? AND created_by = ?", vendorID, createdBy).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("vendor_id = ? AND created_by = ?", vendorID, createdBy).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)
	return bills, total, nil
}

func (r *billRepository) FindByStatus(status string, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)
	return bills, total, nil
}

func (r *billRepository) FindByStatusAndCreatedBy(status string, createdBy string, limit, offset int) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")

	if err := query.Model(&models.Bill{}).Where("status = ? AND created_by = ?", status, createdBy).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Where("status = ? AND created_by = ?", status, createdBy).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)
	return bills, total, nil
}

// populateUserInfoForBills fetches and populates user information for all bills
func (r *billRepository) populateUserInfoForBills(bills []models.Bill) {
	if len(bills) == 0 {
		return
	}

	billIDMap := make(map[string]*models.Bill)
	var billIDs []string
	for i := range bills {
		if bills[i].CreatedBy != "" {
			billIDMap[bills[i].ID] = &bills[i]
			billIDs = append(billIDs, bills[i].ID)
		}
	}

	if len(billIDs) == 0 {
		return
	}

	var userInfos []struct {
		BillID      string
		UserName    string
		CompanyID   uint
		CompanyName string
	}

	r.db.
		Select("bills.id as bill_id, u.email as user_name, c.id as company_id, c.company_name as company_name").
		Table("bills").
		Joins("LEFT JOIN users u ON u.id = CAST(bills.created_by AS UNSIGNED)").
		Joins("LEFT JOIN companies c ON c.id = u.company_id").
		Where("bills.id IN ?", billIDs).
		Scan(&userInfos)

	for _, info := range userInfos {
		if bill, ok := billIDMap[info.BillID]; ok {
			bill.CreatedByUserName = info.UserName
			bill.CreatedByCompanyID = info.CompanyID
			bill.CreatedByCompanyName = info.CompanyName
		}
	}
}

func (r *billRepository) Update(id string, bill *models.Bill) (*models.Bill, error) {
	if err := r.db.Model(&models.Bill{}).Where("id = ?", id).Updates(bill).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *billRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Bill{}).Error
}

func (r *billRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Bill{}).Where("id = ?", id).Update("status", status).Error
}

func (r *billRepository) GetDB() *gorm.DB {
	return r.db
}
