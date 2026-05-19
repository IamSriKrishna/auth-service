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

	// Get attendance records for the month
	startDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	attendanceRecords, err := s.attendanceRepo.GetByEmployeeAndDateRangeNoLimit(req.EmployeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate working days and attendance breakdown
	totalWorkingDays, presentDays, absentDays, halfDays, holidayDays, leaveDays := s.countAttendanceDays(startDate, endDate, attendanceRecords)

	// Calculate daily/weekly rate
	var dailyRate float64
	if employee.SalaryType == "weekly" && employee.WeeklySalary > 0 {
		// Assuming 5 working days per week
		dailyRate = employee.WeeklySalary / 5
	} else {
		// Monthly salary: divided by working days in month
		if totalWorkingDays > 0 {
			dailyRate = employee.MonthlySalary / float64(totalWorkingDays)
		}
	}

	// Calculate earnings
	earningAmount := (float64(presentDays) * dailyRate) + (float64(halfDays) * (dailyRate / 2))

	// Calculate deductions
	deductionAmount := float64(absentDays) * dailyRate

	// Calculate net salary
	netSalary := earningAmount - deductionAmount

	// Create salary calculation record
	salaryCalc := &models.SalaryCalculation{
		EmployeeID:       req.EmployeeID,
		CompanyID:        companyID,
		Month:            req.Month,
		Year:             req.Year,
		BaseSalary:       employee.MonthlySalary,
		SalaryType:       employee.SalaryType,
		TotalWorkingDays: totalWorkingDays,
		PresentDays:      presentDays,
		AbsentDays:       absentDays,
		HalfDays:         halfDays,
		HolidayDays:      holidayDays,
		LeaveDays:        leaveDays,
		DailyRate:        dailyRate,
		EarningAmount:    earningAmount,
		DeductionAmount:  deductionAmount,
		NetSalary:        netSalary,
		Status:           "pending",
		Notes:            fmt.Sprintf("Calculated on %s", time.Now().Format("2006-01-02")),
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

func (s *salaryService) countAttendanceDays(startDate, endDate time.Time, records []models.EmployeeAttendance) (total, present, absent, halfDay, holiday, leave int) {
	attendanceMap := make(map[string]models.AttendanceStatus)
	for _, record := range records {
		attendanceMap[record.Date.Format("2006-01-02")] = record.Status
	}

	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		dayOfWeek := current.Weekday()

		// Skip Sundays
		if dayOfWeek == time.Sunday {
			current = current.AddDate(0, 0, 1)
			continue
		}

		total++

		status, exists := attendanceMap[current.Format("2006-01-02")]
		if !exists {
			// No record means absent
			absent++
		} else {
			switch status {
			case models.AttendanceOnTime, models.AttendanceLate:
				present++
			case models.AttendanceAbsent:
				absent++
			case models.AttendanceHalfDay:
				halfDay++
			case models.AttendanceHoliday:
				holiday++
			case models.AttendanceLeave:
				leave++
			default:
				absent++
			}
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
