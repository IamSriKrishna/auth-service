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
	err := r.db.Where("id = ?", id).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByUserID(userID uint, offset, limit int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var count int64

	err := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&employees).
		Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.Employee{}).
		Where("user_id = ?", userID).
		Count(&count)

	return employees, count, nil
}

func (r *employeeRepository) GetByCompanyAndUser(companyID, userID uint, offset, limit int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var count int64

	err := r.db.Where("company_id = ? AND user_id = ?", companyID, userID).
		Offset(offset).
		Limit(limit).
		Find(&employees).
		Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.Employee{}).
		Where("company_id = ? AND user_id = ?", companyID, userID).
		Count(&count)

	return employees, count, nil
}

func (r *employeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}
