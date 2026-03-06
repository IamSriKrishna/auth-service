package services

import (
	"context"
	"errors"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type EmployeeAttendanceService interface {
	CreateAttendance(ctx context.Context, companyID uint, req *input.CreateAttendanceRequest) (*output.AttendanceOutput, error)
	BulkCreateAttendance(ctx context.Context, companyID uint, req *input.BulkAttendanceRequest) (*output.BulkAttendanceResponse, error)
	GetAttendanceByID(ctx context.Context, attendanceID, companyID uint) (*output.AttendanceOutput, error)
	GetAttendanceByEmployeeID(ctx context.Context, employeeID, companyID uint, page, limit int) (*output.AttendanceListOutput, error)
	GetAttendanceByCompanyID(ctx context.Context, companyID uint, page, limit int) (*output.AttendanceListOutput, error)
	GetAttendanceByDateRange(ctx context.Context, companyID uint, fromDate, toDate time.Time, page, limit int) (*output.AttendanceListOutput, error)
	GetAttendanceByEmployeeAndDateRange(ctx context.Context, employeeID, companyID uint, fromDate, toDate time.Time, page, limit int) (*output.AttendanceListOutput, error)
	GetAttendanceCalendarView(ctx context.Context, employeeID, companyID uint, fromDate, toDate time.Time) (*output.AttendanceCalendarOutput, error)
	GetCompanyAttendanceWeekView(ctx context.Context, companyID uint, fromDate, toDate time.Time) (*output.CompanyAttendanceWeekOutput, error)
	UpdateAttendance(ctx context.Context, attendanceID, companyID uint, req *input.UpdateAttendanceRequest) (*output.AttendanceOutput, error)
	DeleteAttendance(ctx context.Context, attendanceID, companyID uint) error
	GetAttendanceStats(ctx context.Context, companyID uint, fromDate, toDate time.Time) (map[string]interface{}, error)
	CheckInEmployee(ctx context.Context, employeeID, companyID uint) (*output.AttendanceOutput, error)
	CheckOutEmployee(ctx context.Context, employeeID, companyID uint) (*output.AttendanceOutput, error)
}

type employeeAttendanceService struct {
	attendanceRepo repo.EmployeeAttendanceRepository
	employeeRepo   repo.EmployeeRepository
}

func NewEmployeeAttendanceService(
	attendanceRepo repo.EmployeeAttendanceRepository,
	employeeRepo repo.EmployeeRepository,
) EmployeeAttendanceService {
	return &employeeAttendanceService{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
	}
}

func (s *employeeAttendanceService) CreateAttendance(ctx context.Context, companyID uint, req *input.CreateAttendanceRequest) (*output.AttendanceOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Parse date string
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Verify employee exists
	employee, err := s.employeeRepo.GetByID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("employee does not belong to this company")
	}

	// Check if attendance already exists for this date
	existingAttendance, _ := s.attendanceRepo.GetByEmployeeAndDate(req.EmployeeID, parsedDate)
	if existingAttendance != nil {
		return nil, errors.New("attendance already recorded for this date")
	}

	attendance := &models.EmployeeAttendance{
		EmployeeID:   req.EmployeeID,
		CompanyID:    companyID,
		Date:         parsedDate,
		Status:       models.AttendanceStatus(req.Status),
		Reason:       req.Reason,
		CheckInTime:  req.CheckInTime,
		CheckOutTime: req.CheckOutTime,
		WorkingHours: req.WorkingHours,
		Notes:        req.Notes,
	}

	if err := s.attendanceRepo.Create(attendance); err != nil {
		return nil, err
	}

	return s.mapAttendanceToOutput(attendance), nil
}

func (s *employeeAttendanceService) BulkCreateAttendance(ctx context.Context, companyID uint, req *input.BulkAttendanceRequest) (*output.BulkAttendanceResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Verify employee exists
	employee, err := s.employeeRepo.GetByID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("employee does not belong to this company")
	}

	response := &output.BulkAttendanceResponse{
		Success:         true,
		Message:         "Attendance records processed",
		SuccessfulDates: []string{},
		FailedDates:     []string{},
	}

	for _, day := range req.Attendance {
		parsedDate, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			response.FailedDates = append(response.FailedDates, day.Date)
			response.FailedCount++
			continue
		}

		// Check if attendance already exists for this date
		existingAttendance, _ := s.attendanceRepo.GetByEmployeeAndDate(req.EmployeeID, parsedDate)
		if existingAttendance != nil {
			response.FailedDates = append(response.FailedDates, day.Date)
			response.FailedCount++
			continue
		}

		attendance := &models.EmployeeAttendance{
			EmployeeID:   req.EmployeeID,
			CompanyID:    companyID,
			Date:         parsedDate,
			Status:       models.AttendanceStatus(day.Status),
			Reason:       day.Reason,
			CheckInTime:  day.CheckInTime,
			CheckOutTime: day.CheckOutTime,
			WorkingHours: day.WorkingHours,
			Notes:        day.Notes,
		}

		if err := s.attendanceRepo.Create(attendance); err != nil {
			response.FailedDates = append(response.FailedDates, day.Date)
			response.FailedCount++
			continue
		}

		response.SuccessfulDates = append(response.SuccessfulDates, day.Date)
		response.CreatedCount++
	}

	return response, nil
}

func (s *employeeAttendanceService) GetAttendanceByID(ctx context.Context, attendanceID, companyID uint) (*output.AttendanceOutput, error) {
	attendance, err := s.attendanceRepo.GetByID(attendanceID)
	if err != nil {
		return nil, errors.New("attendance not found")
	}

	if attendance.CompanyID != companyID {
		return nil, errors.New("unauthorized access to attendance record")
	}

	return s.mapAttendanceToOutput(attendance), nil
}

func (s *employeeAttendanceService) GetAttendanceByEmployeeID(ctx context.Context, employeeID, companyID uint, page, limit int) (*output.AttendanceListOutput, error) {
	// Verify employee exists and belongs to company
	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("employee does not belong to this company")
	}

	offset := (page - 1) * limit
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	attendances, count, err := s.attendanceRepo.GetByEmployeeID(employeeID, companyID, offset, limit)
	if err != nil {
		return nil, err
	}

	output := &output.AttendanceListOutput{
		Data: make([]output.AttendanceOutput, len(attendances)),
	}

	for i, a := range attendances {
		output.Data[i] = *s.mapAttendanceToOutput(&a)
	}

	output.Pagination.Total = count
	output.Pagination.Page = page
	output.Pagination.PageSize = limit
	output.Pagination.TotalPages = int((count + int64(limit) - 1) / int64(limit))

	return output, nil
}

func (s *employeeAttendanceService) GetAttendanceByCompanyID(ctx context.Context, companyID uint, page, limit int) (*output.AttendanceListOutput, error) {
	offset := (page - 1) * limit
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	attendances, count, err := s.attendanceRepo.GetByCompanyID(companyID, offset, limit)
	if err != nil {
		return nil, err
	}

	output := &output.AttendanceListOutput{
		Data: make([]output.AttendanceOutput, len(attendances)),
	}

	for i, a := range attendances {
		output.Data[i] = *s.mapAttendanceToOutput(&a)
	}

	output.Pagination.Total = count
	output.Pagination.Page = page
	output.Pagination.PageSize = limit
	output.Pagination.TotalPages = int((count + int64(limit) - 1) / int64(limit))

	return output, nil
}

func (s *employeeAttendanceService) GetAttendanceByDateRange(ctx context.Context, companyID uint, fromDate, toDate time.Time, page, limit int) (*output.AttendanceListOutput, error) {
	offset := (page - 1) * limit
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	attendances, count, err := s.attendanceRepo.GetByDateRange(companyID, fromDate, toDate, offset, limit)
	if err != nil {
		return nil, err
	}

	output := &output.AttendanceListOutput{
		Data: make([]output.AttendanceOutput, len(attendances)),
	}

	for i, a := range attendances {
		output.Data[i] = *s.mapAttendanceToOutput(&a)
	}

	output.Pagination.Total = count
	output.Pagination.Page = page
	output.Pagination.PageSize = limit
	output.Pagination.TotalPages = int((count + int64(limit) - 1) / int64(limit))

	return output, nil
}

func (s *employeeAttendanceService) GetAttendanceByEmployeeAndDateRange(ctx context.Context, employeeID, companyID uint, fromDate, toDate time.Time, page, limit int) (*output.AttendanceListOutput, error) {
	// Verify employee exists and belongs to company
	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("employee does not belong to this company")
	}

	offset := (page - 1) * limit
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	attendances, count, err := s.attendanceRepo.GetByEmployeeAndDateRange(employeeID, companyID, fromDate, toDate, offset, limit)
	if err != nil {
		return nil, err
	}

	output := &output.AttendanceListOutput{
		Data: make([]output.AttendanceOutput, len(attendances)),
	}

	for i, a := range attendances {
		output.Data[i] = *s.mapAttendanceToOutput(&a)
	}

	output.Pagination.Total = count
	output.Pagination.Page = page
	output.Pagination.PageSize = limit
	output.Pagination.TotalPages = int((count + int64(limit) - 1) / int64(limit))

	return output, nil
}

func (s *employeeAttendanceService) UpdateAttendance(ctx context.Context, attendanceID, companyID uint, req *input.UpdateAttendanceRequest) (*output.AttendanceOutput, error) {
	attendance, err := s.attendanceRepo.GetByID(attendanceID)
	if err != nil {
		return nil, errors.New("attendance not found")
	}

	if attendance.CompanyID != companyID {
		return nil, errors.New("unauthorized access to attendance record")
	}

	if req.Status != nil {
		attendance.Status = models.AttendanceStatus(*req.Status)
	}
	if req.Reason != nil {
		attendance.Reason = *req.Reason
	}
	if req.CheckInTime != nil {
		attendance.CheckInTime = req.CheckInTime
	}
	if req.CheckOutTime != nil {
		attendance.CheckOutTime = req.CheckOutTime
	}
	if req.WorkingHours != nil {
		attendance.WorkingHours = *req.WorkingHours
	}
	if req.Notes != nil {
		attendance.Notes = *req.Notes
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return s.mapAttendanceToOutput(attendance), nil
}

func (s *employeeAttendanceService) DeleteAttendance(ctx context.Context, attendanceID, companyID uint) error {
	attendance, err := s.attendanceRepo.GetByID(attendanceID)
	if err != nil {
		return errors.New("attendance not found")
	}

	if attendance.CompanyID != companyID {
		return errors.New("unauthorized access to attendance record")
	}

	return s.attendanceRepo.Delete(attendanceID)
}

func (s *employeeAttendanceService) GetAttendanceStats(ctx context.Context, companyID uint, fromDate, toDate time.Time) (map[string]interface{}, error) {
	return s.attendanceRepo.GetAttendanceStats(companyID, fromDate, toDate)
}

func (s *employeeAttendanceService) CheckInEmployee(ctx context.Context, employeeID, companyID uint) (*output.AttendanceOutput, error) {
	// Get or create today's attendance record
	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	attendance, _ := s.attendanceRepo.GetByEmployeeAndDate(employeeID, startOfDay)

	now := time.Now().UTC()

	if attendance == nil {
		// Create new attendance record
		attendance = &models.EmployeeAttendance{
			EmployeeID:  employeeID,
			CompanyID:   companyID,
			Date:        startOfDay,
			Status:      models.AttendanceOnTime,
			CheckInTime: &now,
		}

		if err := s.attendanceRepo.Create(attendance); err != nil {
			return nil, err
		}
	} else {
		// Update existing record with check-in time
		attendance.CheckInTime = &now
		attendance.Status = models.AttendanceOnTime

		if err := s.attendanceRepo.Update(attendance); err != nil {
			return nil, err
		}
	}

	return s.mapAttendanceToOutput(attendance), nil
}

func (s *employeeAttendanceService) CheckOutEmployee(ctx context.Context, employeeID, companyID uint) (*output.AttendanceOutput, error) {
	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	attendance, err := s.attendanceRepo.GetByEmployeeAndDate(employeeID, startOfDay)
	if err != nil || attendance == nil {
		return nil, errors.New("no check-in found for today")
	}

	now := time.Now().UTC()
	attendance.CheckOutTime = &now

	// Calculate working hours
	if attendance.CheckInTime != nil {
		duration := now.Sub(*attendance.CheckInTime)
		attendance.WorkingHours = duration.Hours()
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return s.mapAttendanceToOutput(attendance), nil
}

func (s *employeeAttendanceService) mapAttendanceToOutput(attendance *models.EmployeeAttendance) *output.AttendanceOutput {
	return &output.AttendanceOutput{
		ID:           attendance.ID,
		EmployeeID:   attendance.EmployeeID,
		CompanyID:    attendance.CompanyID,
		Date:         attendance.Date,
		Status:       string(attendance.Status),
		Reason:       attendance.Reason,
		CheckInTime:  attendance.CheckInTime,
		CheckOutTime: attendance.CheckOutTime,
		WorkingHours: attendance.WorkingHours,
		Notes:        attendance.Notes,
		CreatedAt:    attendance.CreatedAt,
		UpdatedAt:    attendance.UpdatedAt,
	}
}

func (s *employeeAttendanceService) GetAttendanceCalendarView(ctx context.Context, employeeID, companyID uint, fromDate, toDate time.Time) (*output.AttendanceCalendarOutput, error) {
	// Verify employee exists and belongs to company
	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return nil, errors.New("employee does not belong to this company")
	}

	// Get attendance records for the date range
	attendances, _, err := s.attendanceRepo.GetByEmployeeAndDateRange(employeeID, companyID, fromDate, toDate, 0, 1000)
	if err != nil {
		return nil, err
	}

	// Build calendar view
	calendarOutput := &output.AttendanceCalendarOutput{
		EmployeeID:   employeeID,
		EmployeeName: employee.Name,
		CompanyID:    companyID,
		StartDate:    fromDate.Format("2006-01-02"),
		EndDate:      toDate.Format("2006-01-02"),
		DailyStatus:  make(map[string]output.DailyAttendance),
		Statistics: output.AttendanceStats{
			Total:   0,
			OnTime:  0,
			Absent:  0,
			Late:    0,
			Holiday: 0,
			HalfDay: 0,
			Leave:   0,
		},
	}

	// Process attendance records
	for _, att := range attendances {
		dateKey := att.Date.Format("2006-01-02")
		calendarOutput.DailyStatus[dateKey] = output.DailyAttendance{
			Date:         dateKey,
			Status:       string(att.Status),
			Reason:       att.Reason,
			CheckInTime:  att.CheckInTime,
			CheckOutTime: att.CheckOutTime,
			WorkingHours: att.WorkingHours,
			Notes:        att.Notes,
		}

		// Update statistics
		calendarOutput.Statistics.Total++
		switch att.Status {
		case models.AttendanceOnTime:
			calendarOutput.Statistics.OnTime++
		case models.AttendanceAbsent:
			calendarOutput.Statistics.Absent++
		case models.AttendanceLate:
			calendarOutput.Statistics.Late++
		case models.AttendanceHoliday:
			calendarOutput.Statistics.Holiday++
		case models.AttendanceHalfDay:
			calendarOutput.Statistics.HalfDay++
		case models.AttendanceLeave:
			calendarOutput.Statistics.Leave++
		}
	}

	return calendarOutput, nil
}

func (s *employeeAttendanceService) GetCompanyAttendanceWeekView(ctx context.Context, companyID uint, fromDate, toDate time.Time) (*output.CompanyAttendanceWeekOutput, error) {
	// Get all employees in the company
	employees, _, err := s.employeeRepo.GetByCompany(companyID, 0, 1000)
	if err != nil {
		return nil, errors.New("failed to fetch employees")
	}

	weekOutput := &output.CompanyAttendanceWeekOutput{
		CompanyID: companyID,
		StartDate: fromDate.Format("2006-01-02"),
		EndDate:   toDate.Format("2006-01-02"),
		DateRange: generateDateRange(fromDate, toDate),
		Employees: make([]output.EmployeeWeekAttendance, 0),
		CompanyStats: output.CompanyAttendanceStats{
			TotalEmployees: len(employees),
		},
	}

	// Process each employee
	for _, emp := range employees {
		// Get attendance for this employee in the date range
		attendances, _, err := s.attendanceRepo.GetByEmployeeAndDateRange(emp.ID, companyID, fromDate, toDate, 0, 1000)
		if err != nil {
			continue
		}

		empWeekAttendance := output.EmployeeWeekAttendance{
			EmployeeID:      emp.ID,
			EmployeeName:    emp.Name,
			Email:           emp.Email,
			EmployeeType:    emp.EmployeeType,
			DailyAttendance: make(map[string]output.DailyAttendance),
			WeekStats: output.AttendanceStats{
				Total:   0,
				OnTime:  0,
				Absent:  0,
				Late:    0,
				Holiday: 0,
				HalfDay: 0,
				Leave:   0,
			},
		}

		// Process attendance records
		for _, att := range attendances {
			dateKey := att.Date.Format("2006-01-02")
			empWeekAttendance.DailyAttendance[dateKey] = output.DailyAttendance{
				Date:         dateKey,
				Status:       string(att.Status),
				Reason:       att.Reason,
				CheckInTime:  att.CheckInTime,
				CheckOutTime: att.CheckOutTime,
				WorkingHours: att.WorkingHours,
				Notes:        att.Notes,
			}

			// Update employee statistics
			empWeekAttendance.WeekStats.Total++
			switch att.Status {
			case models.AttendanceOnTime:
				empWeekAttendance.WeekStats.OnTime++
				weekOutput.CompanyStats.TotalPresent++
			case models.AttendanceAbsent:
				empWeekAttendance.WeekStats.Absent++
				weekOutput.CompanyStats.TotalAbsent++
			case models.AttendanceLate:
				empWeekAttendance.WeekStats.Late++
				weekOutput.CompanyStats.TotalLate++
			case models.AttendanceHoliday:
				empWeekAttendance.WeekStats.Holiday++
				weekOutput.CompanyStats.TotalHoliday++
			case models.AttendanceHalfDay:
				empWeekAttendance.WeekStats.HalfDay++
				weekOutput.CompanyStats.TotalHalfDay++
			case models.AttendanceLeave:
				empWeekAttendance.WeekStats.Leave++
				weekOutput.CompanyStats.TotalLeave++
			}
		}

		weekOutput.Employees = append(weekOutput.Employees, empWeekAttendance)
	}

	return weekOutput, nil
}

// Helper function to generate date range
func generateDateRange(fromDate, toDate time.Time) []string {
	var dates []string
	for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}
