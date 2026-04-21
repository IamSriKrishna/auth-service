package input

type CreateEmployeeRequest struct {
	Name          string  `form:"name" validate:"required"`
	Email         string  `form:"email" validate:"omitempty,email"`
	Number        string  `form:"number" validate:"required"`
	Address       string  `form:"address" validate:"required"`
	EmployeeType  string  `form:"employee_type" validate:"required,oneof=part-time full-time"`
	MonthlySalary float64 `form:"monthly_salary" validate:"required,gt=0"`
}

type UpdateEmployeeRequest struct {
	Name          *string  `json:"name"`
	Email         *string  `json:"email" validate:"omitempty,email"`
	Number        *string  `json:"number"`
	Address       *string  `json:"address"`
	EmployeeType  *string  `json:"employee_type" validate:"omitempty,oneof=part-time full-time"`
	MonthlySalary *float64 `json:"monthly_salary" validate:"omitempty,gt=0"`
}
