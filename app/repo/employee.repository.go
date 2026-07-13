package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) GetByID(id uint) (*models.Employee, error) {
	var employee models.Employee

	err := r.db.
		Where("id = ?", id).
		First(&employee).
		Error

	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *employeeRepository) GetByIDAndCompanyID(
	id uint,
	companyID uint,
) (*models.Employee, error) {
	var employee models.Employee

	err := r.db.
		Where(
			"id = ? AND company_id = ?",
			id,
			companyID,
		).
		First(&employee).
		Error

	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *employeeRepository) GetByUserID(
	userID uint,
	offset int,
	limit int,
) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.
		Model(&models.Employee{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&employees).
		Error

	return employees, total, err
}

func (r *employeeRepository) GetByCompany(
	companyID uint,
	offset int,
	limit int,
) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.
		Model(&models.Employee{}).
		Where("company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&employees).
		Error

	return employees, total, err
}

func (r *employeeRepository) GetByCompanyAndUser(
	companyID uint,
	userID uint,
	offset int,
	limit int,
) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.
		Model(&models.Employee{}).
		Where(
			"company_id = ? AND user_id = ?",
			companyID,
			userID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&employees).
		Error

	return employees, total, err
}

func (r *employeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) UpdateByCompanyID(
	employee *models.Employee,
	companyID uint,
) error {
	if employee == nil {
		return gorm.ErrInvalidData
	}

	result := r.db.
		Model(&models.Employee{}).
		Where(
			"id = ? AND company_id = ?",
			employee.ID,
			companyID,
		).
		Updates(employee)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

func (r *employeeRepository) DeleteByIDAndCompanyID(
	id uint,
	companyID uint,
) error {
	result := r.db.
		Where(
			"id = ? AND company_id = ?",
			id,
			companyID,
		).
		Delete(&models.Employee{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
