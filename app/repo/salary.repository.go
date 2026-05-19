package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type salaryRepository struct {
	db *gorm.DB
}

func NewSalaryRepository(db *gorm.DB) SalaryRepository {
	return &salaryRepository{db: db}
}

func (r *salaryRepository) Create(salary *models.SalaryCalculation) error {
	return r.db.Create(salary).Error
}

func (r *salaryRepository) GetByID(id uint) (*models.SalaryCalculation, error) {
	var salary models.SalaryCalculation
	err := r.db.Where("id = ?", id).First(&salary).Error
	if err != nil {
		return nil, err
	}
	return &salary, nil
}

func (r *salaryRepository) GetByEmployee(employeeID uint) ([]models.SalaryCalculation, error) {
	var salaries []models.SalaryCalculation
	err := r.db.Where("employee_id = ?", employeeID).
		Order("year DESC, month DESC").
		Find(&salaries).Error
	if err != nil {
		return nil, err
	}
	return salaries, nil
}

func (r *salaryRepository) GetByEmployeeAndMonth(employeeID uint, month, year int) (*models.SalaryCalculation, error) {
	var salary models.SalaryCalculation
	err := r.db.Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		First(&salary).Error
	if err != nil {
		return nil, err
	}
	return &salary, nil
}

func (r *salaryRepository) GetByCompany(companyID uint, limit, offset int) ([]models.SalaryCalculation, int64, error) {
	var salaries []models.SalaryCalculation
	var count int64

	err := r.db.Where("company_id = ?", companyID).
		Order("year DESC, month DESC").
		Offset(offset).
		Limit(limit).
		Find(&salaries).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.SalaryCalculation{}).
		Where("company_id = ?", companyID).
		Count(&count)

	return salaries, count, nil
}

func (r *salaryRepository) Update(salary *models.SalaryCalculation) error {
	return r.db.Save(salary).Error
}

func (r *salaryRepository) Delete(id uint) error {
	return r.db.Delete(&models.SalaryCalculation{}, id).Error
}
