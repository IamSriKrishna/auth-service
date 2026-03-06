package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type employeeAttendanceRepository struct {
	db *gorm.DB
}

func NewEmployeeAttendanceRepository(db *gorm.DB) EmployeeAttendanceRepository {
	return &employeeAttendanceRepository{db: db}
}

func (r *employeeAttendanceRepository) Create(attendance *models.EmployeeAttendance) error {
	return r.db.Create(attendance).Error
}

func (r *employeeAttendanceRepository) GetByID(id uint) (*models.EmployeeAttendance, error) {
	var attendance models.EmployeeAttendance
	err := r.db.Where("id = ?", id).First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *employeeAttendanceRepository) GetByEmployeeAndDate(employeeID uint, date time.Time) (*models.EmployeeAttendance, error) {
	var attendance models.EmployeeAttendance
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	err := r.db.Where("employee_id = ? AND date >= ? AND date < ?", employeeID, startOfDay, endOfDay).
		First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *employeeAttendanceRepository) GetByEmployeeID(employeeID, companyID uint, offset, limit int) ([]models.EmployeeAttendance, int64, error) {
	var attendance []models.EmployeeAttendance
	var count int64

	err := r.db.Where("employee_id = ? AND company_id = ?", employeeID, companyID).
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendance).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("employee_id = ? AND company_id = ?", employeeID, companyID).
		Count(&count)

	return attendance, count, nil
}

func (r *employeeAttendanceRepository) GetByCompanyID(companyID uint, offset, limit int) ([]models.EmployeeAttendance, int64, error) {
	var attendance []models.EmployeeAttendance
	var count int64

	err := r.db.Where("company_id = ?", companyID).
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendance).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("company_id = ?", companyID).
		Count(&count)

	return attendance, count, nil
}

func (r *employeeAttendanceRepository) GetByDateRange(companyID uint, fromDate, toDate time.Time, offset, limit int) ([]models.EmployeeAttendance, int64, error) {
	var attendance []models.EmployeeAttendance
	var count int64

	err := r.db.Where("company_id = ? AND date >= ? AND date <= ?", companyID, fromDate, toDate).
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendance).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("company_id = ? AND date >= ? AND date <= ?", companyID, fromDate, toDate).
		Count(&count)

	return attendance, count, nil
}

func (r *employeeAttendanceRepository) GetByEmployeeAndDateRange(employeeID, companyID uint, fromDate, toDate time.Time, offset, limit int) ([]models.EmployeeAttendance, int64, error) {
	var attendance []models.EmployeeAttendance
	var count int64

	err := r.db.Where("employee_id = ? AND company_id = ? AND date >= ? AND date <= ?", employeeID, companyID, fromDate, toDate).
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendance).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("employee_id = ? AND company_id = ? AND date >= ? AND date <= ?", employeeID, companyID, fromDate, toDate).
		Count(&count)

	return attendance, count, nil
}

func (r *employeeAttendanceRepository) Update(attendance *models.EmployeeAttendance) error {
	return r.db.Save(attendance).Error
}

func (r *employeeAttendanceRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeAttendance{}, id).Error
}

func (r *employeeAttendanceRepository) GetAttendanceStats(companyID uint, fromDate, toDate time.Time) (map[string]interface{}, error) {
	var stats struct {
		OnTime  int64
		Absent  int64
		Late    int64
		Holiday int64
		HalfDay int64
		Leave   int64
		Total   int64
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("company_id = ? AND date >= ? AND date <= ?", companyID, fromDate, toDate).
		Select("COUNT(*) as total," +
			"SUM(CASE WHEN status = 'on_time' THEN 1 ELSE 0 END) as on_time," +
			"SUM(CASE WHEN status = 'absent' THEN 1 ELSE 0 END) as absent," +
			"SUM(CASE WHEN status = 'late' THEN 1 ELSE 0 END) as late," +
			"SUM(CASE WHEN status = 'holiday' THEN 1 ELSE 0 END) as holiday," +
			"SUM(CASE WHEN status = 'half_day' THEN 1 ELSE 0 END) as half_day," +
			"SUM(CASE WHEN status = 'leave' THEN 1 ELSE 0 END) as leave").
		Scan(&stats)

	return map[string]interface{}{
		"total":    stats.Total,
		"on_time":  stats.OnTime,
		"absent":   stats.Absent,
		"late":     stats.Late,
		"holiday":  stats.Holiday,
		"half_day": stats.HalfDay,
		"leave":    stats.Leave,
	}, nil
}
