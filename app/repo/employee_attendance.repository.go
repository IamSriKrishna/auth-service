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
	// Use local timezone to match how dates are stored
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	err := r.db.Where("employee_id = ? AND date >= ? AND date < ? AND deleted_at IS NULL", employeeID, startOfDay, endOfDay).
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

	// Add one day to toDate to ensure all records for that date are included
	toDateInclusive := toDate.AddDate(0, 0, 1)

	err := r.db.Where("company_id = ? AND date >= ? AND date < ?", companyID, fromDate, toDateInclusive).
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendance).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("company_id = ? AND date >= ? AND date < ?", companyID, fromDate, toDateInclusive).
		Count(&count)

	return attendance, count, nil
}

func (r *employeeAttendanceRepository) GetByEmployeeAndDateRange(employeeID, companyID uint, fromDate, toDate time.Time, offset, limit int) ([]models.EmployeeAttendance, int64, error) {
	var attendance []models.EmployeeAttendance
	var count int64

	// Add one day to toDate to ensure all records for that date are included
	// Use local timezone to match how dates are stored
	fromDateLocal := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)
	toDateInclusive := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)

	err := r.db.Where("employee_id = ? AND company_id = ? AND date >= ? AND date < ?", employeeID, companyID, fromDateLocal, toDateInclusive).
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&attendance).Error

	if err != nil {
		return nil, 0, err
	}

	r.db.Model(&models.EmployeeAttendance{}).
		Where("employee_id = ? AND company_id = ? AND date >= ? AND date < ?", employeeID, companyID, fromDateLocal, toDateInclusive).
		Count(&count)

	return attendance, count, nil
}

func (r *employeeAttendanceRepository) GetByEmployeeAndDateRangeNoLimit(employeeID uint, fromDate, toDate time.Time) ([]models.EmployeeAttendance, error) {
	var attendance []models.EmployeeAttendance

	// Normalize to local calendar dates and include the full end date in the range.
	fromDateLocal := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)
	toDateInclusive := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)

	err := r.db.Where("employee_id = ? AND date >= ? AND date < ?", employeeID, fromDateLocal, toDateInclusive).
		Order("date ASC").
		Find(&attendance).Error

	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func (r *employeeAttendanceRepository) Update(attendance *models.EmployeeAttendance) error {
	return r.db.Save(attendance).Error
}

func (r *employeeAttendanceRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeAttendance{}, id).Error
}

func (r *employeeAttendanceRepository) DeleteByEmployeeAndDate(employeeID uint, date time.Time) error {
	// Use local timezone to match how dates are stored
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// Hard delete (unscoped) to remove even soft-deleted records
	return r.db.Unscoped().Where("employee_id = ? AND date >= ? AND date < ?", employeeID, startOfDay, endOfDay).Delete(&models.EmployeeAttendance{}).Error
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

	// Add one day to toDate to ensure all records for that date are included
	toDateInclusive := toDate.AddDate(0, 0, 1)

	r.db.Model(&models.EmployeeAttendance{}).
		Where("company_id = ? AND date >= ? AND date < ?", companyID, fromDate, toDateInclusive).
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
