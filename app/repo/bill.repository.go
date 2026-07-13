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

func billPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Vendor").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Product").
		Preload("LineItems.Item").
		Preload("LineItems.Variant")
}

func billCompanyQuery(
	db *gorm.DB,
	companyID uint,
) *gorm.DB {
	return db.
		Joins(`
			JOIN users bill_creator
				ON bill_creator.id =
				CAST(bills.created_by AS UNSIGNED)
		`).
		Where("bill_creator.company_id = ?", companyID)
}

func (r *billRepository) Create(
	bill *models.Bill,
) (*models.Bill, error) {
	if err := r.db.Create(bill).Error; err != nil {
		return nil, err
	}

	return bill, nil
}

func (r *billRepository) FindByID(
	id string,
) (*models.Bill, error) {
	var bill models.Bill

	if err := billPreloads(r.db).
		Where("bills.id = ?", id).
		First(&bill).
		Error; err != nil {
		return nil, err
	}

	r.populateUserInfoForBill(&bill)

	return &bill, nil
}

func (r *billRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.Bill, error) {
	var bill models.Bill

	query := billCompanyQuery(
		billPreloads(r.db).
			Model(&models.Bill{}),
		companyID,
	)

	if err := query.
		Select("bills.*").
		Where("bills.id = ?", id).
		First(&bill).
		Error; err != nil {
		return nil, err
	}

	r.populateUserInfoForBill(&bill)

	return &bill, nil
}

func (r *billRepository) FindAll(
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	if err := r.db.
		Model(&models.Bill{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	if err := billPreloads(r.db).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	countQuery := billCompanyQuery(
		r.db.Model(&models.Bill{}),
		companyID,
	)

	if err := countQuery.
		Distinct("bills.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := billCompanyQuery(
		billPreloads(r.db).
			Model(&models.Bill{}),
		companyID,
	)

	if err := findQuery.
		Select("bills.*").
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByCreatedBy(
	createdBy string,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Model(&models.Bill{}).
		Where("created_by = ?", createdBy)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := billPreloads(r.db).
		Where("created_by = ?", createdBy).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByVendor(
	vendorID uint,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Model(&models.Bill{}).
		Where("vendor_id = ?", vendorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := billPreloads(r.db).
		Where("vendor_id = ?", vendorID).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByVendorAndCreatedBy(
	vendorID uint,
	createdBy string,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Model(&models.Bill{}).
		Where(
			"vendor_id = ? AND created_by = ?",
			vendorID,
			createdBy,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := billPreloads(r.db).
		Where(
			"vendor_id = ? AND created_by = ?",
			vendorID,
			createdBy,
		).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByVendorAndCompany(
	vendorID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	countQuery := billCompanyQuery(
		r.db.Model(&models.Bill{}),
		companyID,
	).Where("bills.vendor_id = ?", vendorID)

	if err := countQuery.
		Distinct("bills.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := billCompanyQuery(
		billPreloads(r.db).
			Model(&models.Bill{}),
		companyID,
	)

	if err := findQuery.
		Select("bills.*").
		Where("bills.vendor_id = ?", vendorID).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByStatus(
	status string,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Model(&models.Bill{}).
		Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := billPreloads(r.db).
		Where("status = ?", status).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByStatusAndCreatedBy(
	status string,
	createdBy string,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	query := r.db.
		Model(&models.Bill{}).
		Where(
			"status = ? AND created_by = ?",
			status,
			createdBy,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := billPreloads(r.db).
		Where(
			"status = ? AND created_by = ?",
			status,
			createdBy,
		).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) FindByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	countQuery := billCompanyQuery(
		r.db.Model(&models.Bill{}),
		companyID,
	).Where("bills.status = ?", status)

	if err := countQuery.
		Distinct("bills.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := billCompanyQuery(
		billPreloads(r.db).
			Model(&models.Bill{}),
		companyID,
	)

	if err := findQuery.
		Select("bills.*").
		Where("bills.status = ?", status).
		Order("bills.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bills).
		Error; err != nil {
		return nil, 0, err
	}

	r.populateUserInfoForBills(bills)

	return bills, total, nil
}

func (r *billRepository) Update(
	id string,
	bill *models.Bill,
) (*models.Bill, error) {
	if bill == nil {
		return nil, gorm.ErrInvalidData
	}

	if err := r.db.
		Model(&models.Bill{}).
		Where("id = ?", id).
		Updates(bill).
		Error; err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *billRepository) UpdateByCompany(
	id string,
	companyID uint,
	bill *models.Bill,
) (*models.Bill, error) {
	if bill == nil {
		return nil, gorm.ErrInvalidData
	}

	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return nil, err
	}

	if err := r.db.
		Model(&models.Bill{}).
		Where("id = ?", id).
		Updates(bill).
		Error; err != nil {
		return nil, err
	}

	return r.FindByIDAndCompany(
		id,
		companyID,
	)
}

func (r *billRepository) Delete(
	id string,
) error {
	return r.db.
		Where("id = ?", id).
		Delete(&models.Bill{}).
		Error
}

func (r *billRepository) DeleteByCompany(
	id string,
	companyID uint,
) error {
	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return err
	}

	result := r.db.
		Where("id = ?", id).
		Delete(&models.Bill{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *billRepository) UpdateStatus(
	id string,
	status string,
) error {
	return r.db.
		Model(&models.Bill{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *billRepository) UpdateStatusByCompany(
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

	result := r.db.
		Model(&models.Bill{}).
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

func (r *billRepository) populateUserInfoForBill(
	bill *models.Bill,
) {
	if bill == nil || bill.CreatedBy == "" {
		return
	}

	var userInfo struct {
		UserName    string
		CompanyID   uint
		CompanyName string
	}

	err := r.db.
		Select(`
			u.email AS user_name,
			c.id AS company_id,
			c.company_name AS company_name
		`).
		Table("users AS u").
		Joins(
			"LEFT JOIN companies AS c ON c.id = u.company_id",
		).
		Where(
			"u.id = CAST(? AS UNSIGNED)",
			bill.CreatedBy,
		).
		Scan(&userInfo).
		Error

	if err != nil {
		return
	}

	bill.CreatedByUserName = userInfo.UserName
	bill.CreatedByCompanyID = userInfo.CompanyID
	bill.CreatedByCompanyName = userInfo.CompanyName
}

func (r *billRepository) populateUserInfoForBills(
	bills []models.Bill,
) {
	if len(bills) == 0 {
		return
	}

	billMap := make(map[string]*models.Bill)
	billIDs := make([]string, 0, len(bills))

	for index := range bills {
		if bills[index].CreatedBy == "" {
			continue
		}

		billMap[bills[index].ID] = &bills[index]
		billIDs = append(billIDs, bills[index].ID)
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

	if err := r.db.
		Select(`
			bills.id AS bill_id,
			u.email AS user_name,
			c.id AS company_id,
			c.company_name AS company_name
		`).
		Table("bills").
		Joins(`
			LEFT JOIN users AS u
				ON u.id = CAST(bills.created_by AS UNSIGNED)
		`).
		Joins(
			"LEFT JOIN companies AS c ON c.id = u.company_id",
		).
		Where("bills.id IN ?", billIDs).
		Scan(&userInfos).
		Error; err != nil {
		return
	}

	for _, info := range userInfos {
		bill, exists := billMap[info.BillID]
		if !exists {
			continue
		}

		bill.CreatedByUserName = info.UserName
		bill.CreatedByCompanyID = info.CompanyID
		bill.CreatedByCompanyName = info.CompanyName
	}
}

func (r *billRepository) GetDB() *gorm.DB {
	return r.db
}
