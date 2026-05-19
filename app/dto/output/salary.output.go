package output

import "time"

type SalaryCalculationOutput struct {
	ID               uint      `json:"id"`
	EmployeeID       uint      `json:"employee_id"`
	EmployeeName     string    `json:"employee_name"`
	CompanyID        uint      `json:"company_id"`
	Month            int       `json:"month"`
	Year             int       `json:"year"`
	BaseSalary       float64   `json:"base_salary"`
	SalaryType       string    `json:"salary_type"`
	TotalWorkingDays int       `json:"total_working_days"`
	PresentDays      int       `json:"present_days"`
	AbsentDays       int       `json:"absent_days"`
	HalfDays         int       `json:"half_days"`
	HolidayDays      int       `json:"holiday_days"`
	LeaveDays        int       `json:"leave_days"`
	DailyRate        float64   `json:"daily_rate"`
	EarningAmount    float64   `json:"earning_amount"`
	DeductionAmount  float64   `json:"deduction_amount"`
	NetSalary        float64   `json:"net_salary"`
	Status           string    `json:"status"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type SalaryCalculationListOutput struct {
	ID           uint      `json:"id"`
	EmployeeID   uint      `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	Month        int       `json:"month"`
	Year         int       `json:"year"`
	BaseSalary   float64   `json:"base_salary"`
	NetSalary    float64   `json:"net_salary"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
