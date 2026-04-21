package output

import "time"

type EmployeeOutput struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Number        string    `json:"number"`
	Address       string    `json:"address"`
	EmployeeType  string    `json:"employee_type"`
	MonthlySalary float64   `json:"monthly_salary"`
	DocumentURL   string    `json:"document_url"`
	UserID        uint      `json:"user_id"`
	CompanyID     uint      `json:"company_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmployeeListOutput struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Number        string    `json:"number"`
	Address       string    `json:"address"`
	EmployeeType  string    `json:"employee_type"`
	MonthlySalary float64   `json:"monthly_salary"`
	DocumentURL   string    `json:"document_url"`
	UserID        uint      `json:"user_id"`
	CompanyID     uint      `json:"company_id"`
	CreatedAt     time.Time `json:"created_at"`
}
