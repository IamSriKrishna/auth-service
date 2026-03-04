package output

import "time"

type EmployeeOutput struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Number    string    `json:"number"`
	Address   string    `json:"address"`
	UserID    uint      `json:"user_id"`
	CompanyID uint      `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmployeeListOutput struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Number    string    `json:"number"`
	Address   string    `json:"address"`
	UserID    uint      `json:"user_id"`
	CompanyID uint      `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
}
