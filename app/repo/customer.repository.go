package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(
	customer *models.Customer,
) error {
	if customer == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Session(&gorm.Session{
				FullSaveAssociations: true,
			}).
			Create(customer).
			Error
	})
}

func (r *customerRepository) Update(
	customer *models.Customer,
) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) FindByID(
	id uint,
) (*models.Customer, error) {
	var customer models.Customer

	err := r.db.
		Preload("User").
		Preload("Company").
		Preload("OtherDetails").
		Preload("Addresses").
		Preload("ContactPersons").
		Where("id = ?", id).
		First(&customer).
		Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindAll(
	page int,
	limit int,
) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := r.db.Model(&models.Customer{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&customers).
		Error

	return customers, total, err
}

func (r *customerRepository) FindByUserID(
	userID uint,
	companyID uint,
	page int,
	limit int,
) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := r.db.
		Model(&models.Customer{}).
		Where(
			"user_id = ? AND company_id = ?",
			userID,
			companyID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&customers).
		Error

	return customers, total, err
}

func (r *customerRepository) FindByIDAndUser(
	id uint,
	userID uint,
) (*models.Customer, error) {
	var customer models.Customer

	err := r.db.
		Preload("User").
		Preload("Company").
		Preload("OtherDetails").
		Preload("Addresses").
		Preload("ContactPersons").
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		First(&customer).
		Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) Delete(
	customer *models.Customer,
) error {
	return r.db.Delete(customer).Error
}

func (r *customerRepository) FindByMobile(
	mobile string,
) (*models.Customer, error) {
	var customer models.Customer

	err := r.db.
		Where("mobile = ?", mobile).
		First(&customer).
		Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindByIDAndCompany(
	id uint,
	companyID uint,
) (*models.Customer, error) {
	var customer models.Customer

	err := r.db.
		Preload("User").
		Preload("Company").
		Preload("OtherDetails").
		Preload("Addresses").
		Preload("ContactPersons").
		Where(
			"customers.id = ? AND customers.company_id = ?",
			id,
			companyID,
		).
		First(&customer).
		Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindByCompanyID(
	companyID uint,
	page int,
	limit int,
) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := r.db.
		Model(&models.Customer{}).
		Where("company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&customers).
		Error

	return customers, total, err
}

func (r *customerRepository) UpdateByCompanyID(
	customer *models.Customer,
	companyID uint,
) error {
	if customer == nil {
		return gorm.ErrInvalidData
	}

	result := r.db.
		Model(&models.Customer{}).
		Where(
			"id = ? AND company_id = ?",
			customer.ID,
			companyID,
		).
		Updates(customer)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *customerRepository) DeleteByIDAndCompany(
	id uint,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND company_id = ?",
			id,
			companyID,
		).
		Delete(&models.Customer{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *customerRepository) FindByMobileAndCompany(
	mobile string,
	companyID uint,
) (*models.Customer, error) {
	var customer models.Customer

	err := r.db.
		Where(
			"mobile = ? AND company_id = ?",
			mobile,
			companyID,
		).
		First(&customer).
		Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
