package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type SalaryService interface {
	CalculateSalary(ctx context.Context, createdByID, companyID uint, req *input.CalculateSalaryRequest) (*output.SalaryCalculationOutput, error)
	GetSalaryCalculation(ctx context.Context, salaryCalculationID, companyID uint) (*output.SalaryCalculationOutput, error)
	GetSalaryCalculationsByEmployee(ctx context.Context, employeeID, companyID uint) ([]output.SalaryCalculationListOutput, error)
	ApproveSalary(ctx context.Context, salaryCalculationID, companyID uint) error
}

type salaryService struct {
	salaryRepo     repo.SalaryRepository
	employeeRepo   repo.EmployeeRepository
	attendanceRepo repo.EmployeeAttendanceRepository
}

func NewSalaryService(
	salaryRepo repo.SalaryRepository,
	employeeRepo repo.EmployeeRepository,
	attendanceRepo repo.EmployeeAttendanceRepository,
) SalaryService {
	return &salaryService{
		salaryRepo:     salaryRepo,
		employeeRepo:   employeeRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (s *salaryService) CalculateSalary(ctx context.Context, createdByID, companyID uint, req *input.CalculateSalaryRequest) (*output.SalaryCalculationOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Get employee
	employee, err := s.employeeRepo.GetByID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("unauthorized access")
	}

	// Parse date range from request
	startDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		return nil, errors.New("invalid from_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		return nil, errors.New("invalid to_date format, use YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("to_date must be after from_date")
	}

	// Normalize to calendar dates in the service timezone so the full requested range is counted consistently.
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.Local)

	attendanceRecords, err := s.attendanceRepo.GetByEmployeeAndDateRangeNoLimit(req.EmployeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate working days and attendance breakdown from the attendance rows that actually exist.
	totalWorkingDays, presentAndHolidayDays, absentDays, halfDays, holidayDays, leaveDays := s.countAttendanceDaysWithHolidayAsSame(startDate, endDate, attendanceRecords)

	// Calculate daily/weekly rate based on salary type
	var dailyRate float64
	if employee.SalaryType == "weekly" && employee.WeeklySalary > 0 {
		// Weekly salary: divided by 7 days per week
		dailyRate = employee.WeeklySalary / 7
	} else {
		// Monthly salary: divided by total days in month
		daysInMonth := startDate.AddDate(0, 1, -1).Day()
		if daysInMonth > 0 {
			dailyRate = employee.MonthlySalary / float64(daysInMonth)
		}
	}

	// Calculate earnings (both present and holiday days count as full working days)
	earningAmount := (float64(presentAndHolidayDays) * dailyRate) + (float64(halfDays) * (dailyRate / 2))

	// Calculate deductions (only absent days are deducted)
	deductionAmount := float64(absentDays) * dailyRate

	// Calculate net salary
	netSalary := earningAmount - deductionAmount

	// Determine base salary based on salary type
	var baseSalary float64
	if employee.SalaryType == "weekly" {
		baseSalary = employee.WeeklySalary
	} else {
		baseSalary = employee.MonthlySalary
	}

	// Create salary calculation record
	salaryCalc := &models.SalaryCalculation{
		EmployeeID:       req.EmployeeID,
		CompanyID:        companyID,
		Month:            int(startDate.Month()),
		Year:             startDate.Year(),
		BaseSalary:       baseSalary,
		SalaryType:       employee.SalaryType,
		TotalWorkingDays: totalWorkingDays,
		PresentDays:      presentAndHolidayDays,
		AbsentDays:       absentDays,
		HalfDays:         halfDays,
		HolidayDays:      holidayDays,
		LeaveDays:        leaveDays,
		DailyRate:        dailyRate,
		EarningAmount:    earningAmount,
		DeductionAmount:  deductionAmount,
		NetSalary:        netSalary,
		Status:           "pending",
		Notes:            fmt.Sprintf("Calculated based on %s salary on %s", employee.SalaryType, time.Now().Format("2006-01-02")),
	}

	if err := s.salaryRepo.Create(salaryCalc); err != nil {
		return nil, err
	}

	return &output.SalaryCalculationOutput{
		ID:               salaryCalc.ID,
		EmployeeID:       salaryCalc.EmployeeID,
		EmployeeName:     employee.Name,
		CompanyID:        salaryCalc.CompanyID,
		Month:            salaryCalc.Month,
		Year:             salaryCalc.Year,
		BaseSalary:       salaryCalc.BaseSalary,
		SalaryType:       salaryCalc.SalaryType,
		TotalWorkingDays: salaryCalc.TotalWorkingDays,
		PresentDays:      salaryCalc.PresentDays,
		AbsentDays:       salaryCalc.AbsentDays,
		HalfDays:         salaryCalc.HalfDays,
		DailyRate:        salaryCalc.DailyRate,
		EarningAmount:    salaryCalc.EarningAmount,
		DeductionAmount:  salaryCalc.DeductionAmount,
		NetSalary:        salaryCalc.NetSalary,
		Status:           salaryCalc.Status,
		Notes:            salaryCalc.Notes,
		CreatedAt:        salaryCalc.CreatedAt,
		UpdatedAt:        salaryCalc.UpdatedAt,
	}, nil
}

// countAttendanceDaysWithHolidayAsSame uses the attendance rows that actually exist for the selected range.
// Days without an explicit attendance record are ignored, so salary follows the attendance data directly.
func (s *salaryService) countAttendanceDaysWithHolidayAsSame(startDate, endDate time.Time, records []models.EmployeeAttendance) (total, presentAndHoliday, absent, halfDay, holidayDays, leaveDays int) {
	attendanceMap := make(map[string]models.AttendanceStatus, len(records))
	for _, record := range records {
		attendanceMap[record.Date.Format("2006-01-02")] = record.Status
	}

	current := startDate
	for !current.After(endDate) {
		status, exists := attendanceMap[current.Format("2006-01-02")]
		if !exists {
			current = current.AddDate(0, 0, 1)
			continue
		}

		total++

		switch status {
		case models.AttendanceOnTime, models.AttendanceLate:
			presentAndHoliday++
		case models.AttendanceHoliday:
			holidayDays++
			presentAndHoliday++
		case models.AttendanceLeave:
			leaveDays++
			presentAndHoliday++
		case models.AttendanceAbsent:
			absent++
		case models.AttendanceHalfDay:
			halfDay++
		default:
			absent++
		}

		current = current.AddDate(0, 0, 1)
	}

	return
}

func (s *salaryService) GetSalaryCalculation(ctx context.Context, salaryCalculationID, companyID uint) (*output.SalaryCalculationOutput, error) {
	salaryCalc, err := s.salaryRepo.GetByID(salaryCalculationID)
	if err != nil {
		return nil, errors.New("salary calculation not found")
	}

	if salaryCalc.CompanyID != companyID {
		return nil, errors.New("unauthorized access")
	}

	employee, err := s.employeeRepo.GetByID(salaryCalc.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	return &output.SalaryCalculationOutput{
		ID:               salaryCalc.ID,
		EmployeeID:       salaryCalc.EmployeeID,
		EmployeeName:     employee.Name,
		CompanyID:        salaryCalc.CompanyID,
		Month:            salaryCalc.Month,
		Year:             salaryCalc.Year,
		BaseSalary:       salaryCalc.BaseSalary,
		SalaryType:       salaryCalc.SalaryType,
		TotalWorkingDays: salaryCalc.TotalWorkingDays,
		PresentDays:      salaryCalc.PresentDays,
		AbsentDays:       salaryCalc.AbsentDays,
		HalfDays:         salaryCalc.HalfDays,
		HolidayDays:      salaryCalc.HolidayDays,
		LeaveDays:        salaryCalc.LeaveDays,
		DailyRate:        salaryCalc.DailyRate,
		EarningAmount:    salaryCalc.EarningAmount,
		DeductionAmount:  salaryCalc.DeductionAmount,
		NetSalary:        salaryCalc.NetSalary,
		Status:           salaryCalc.Status,
		Notes:            salaryCalc.Notes,
		CreatedAt:        salaryCalc.CreatedAt,
		UpdatedAt:        salaryCalc.UpdatedAt,
	}, nil
}

func (s *salaryService) GetSalaryCalculationsByEmployee(ctx context.Context, employeeID, companyID uint) ([]output.SalaryCalculationListOutput, error) {
	salaryCalcs, err := s.salaryRepo.GetByEmployee(employeeID)
	if err != nil {
		return nil, err
	}

	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("unauthorized access")
	}

	var outputs []output.SalaryCalculationListOutput
	for _, salaryCalc := range salaryCalcs {
		outputs = append(outputs, output.SalaryCalculationListOutput{
			ID:           salaryCalc.ID,
			EmployeeID:   salaryCalc.EmployeeID,
			EmployeeName: employee.Name,
			Month:        salaryCalc.Month,
			Year:         salaryCalc.Year,
			BaseSalary:   salaryCalc.BaseSalary,
			NetSalary:    salaryCalc.NetSalary,
			Status:       salaryCalc.Status,
			CreatedAt:    salaryCalc.CreatedAt,
		})
	}

	return outputs, nil
}

func (s *salaryService) ApproveSalary(ctx context.Context, salaryCalculationID, companyID uint) error {
	salaryCalc, err := s.salaryRepo.GetByID(salaryCalculationID)
	if err != nil {
		return errors.New("salary calculation not found")
	}

	if salaryCalc.CompanyID != companyID {
		return errors.New("unauthorized access")
	}

	salaryCalc.Status = "approved"
	return s.salaryRepo.Update(salaryCalc)
}
