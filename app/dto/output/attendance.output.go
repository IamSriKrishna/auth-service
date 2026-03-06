package output

import "time"

type AttendanceOutput struct {
	ID           uint       `json:"id"`
	EmployeeID   uint       `json:"employee_id"`
	CompanyID    uint       `json:"company_id"`
	Date         time.Time  `json:"date"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	WorkingHours float64    `json:"working_hours"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type AttendanceListOutput struct {
	Data       []AttendanceOutput `json:"data"`
	Pagination struct {
		Total      int64 `json:"total"`
		Page       int   `json:"page"`
		PageSize   int   `json:"page_size"`
		TotalPages int   `json:"total_pages"`
	} `json:"pagination"`
}

// AttendanceCalendarOutput represents calendar view with daily attendance
type AttendanceCalendarOutput struct {
	EmployeeID   uint                       `json:"employee_id"`
	EmployeeName string                     `json:"employee_name"`
	CompanyID    uint                       `json:"company_id"`
	StartDate    string                     `json:"start_date"`
	EndDate      string                     `json:"end_date"`
	DailyStatus  map[string]DailyAttendance `json:"daily_status"` // Key: "2026-03-05"
	Statistics   AttendanceStats            `json:"statistics"`
}

type DailyAttendance struct {
	Date         string     `json:"date"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	WorkingHours float64    `json:"working_hours"`
	Notes        string     `json:"notes"`
}

type AttendanceStats struct {
	Total   int `json:"total"`
	OnTime  int `json:"on_time"`
	Absent  int `json:"absent"`
	Late    int `json:"late"`
	Holiday int `json:"holiday"`
	HalfDay int `json:"half_day"`
	Leave   int `json:"leave"`
}

// BulkAttendanceResponse for bulk create operations
type BulkAttendanceResponse struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message"`
	CreatedCount    int      `json:"created_count"`
	FailedCount     int      `json:"failed_count"`
	FailedDates     []string `json:"failed_dates,omitempty"`
	SuccessfulDates []string `json:"successful_dates,omitempty"`
}

// CompanyAttendanceWeekOutput - All employees with attendance for a date range
type CompanyAttendanceWeekOutput struct {
	CompanyID    uint                     `json:"company_id"`
	StartDate    string                   `json:"start_date"`
	EndDate      string                   `json:"end_date"`
	DateRange    []string                 `json:"date_range"` // ["2026-03-01", "2026-03-02", ...]
	Employees    []EmployeeWeekAttendance `json:"employees"`
	CompanyStats CompanyAttendanceStats   `json:"company_stats"`
}

type EmployeeWeekAttendance struct {
	EmployeeID      uint                       `json:"employee_id"`
	EmployeeName    string                     `json:"employee_name"`
	Email           string                     `json:"email"`
	EmployeeType    string                     `json:"employee_type"`
	DailyAttendance map[string]DailyAttendance `json:"daily_attendance"` // Key: "2026-03-01"
	WeekStats       AttendanceStats            `json:"week_stats"`
}

type CompanyAttendanceStats struct {
	TotalEmployees int `json:"total_employees"`
	TotalPresent   int `json:"total_present"`  // Sum of on_time across all employees
	TotalAbsent    int `json:"total_absent"`   // Sum of absent across all employees
	TotalLate      int `json:"total_late"`     // Sum of late across all employees
	TotalHoliday   int `json:"total_holiday"`  // Sum of holiday across all employees
	TotalHalfDay   int `json:"total_half_day"` // Sum of half_day across all employees
	TotalLeave     int `json:"total_leave"`    // Sum of leave across all employees
}
