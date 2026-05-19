package input

type CalculateSalaryRequest struct {
	EmployeeID uint `json:"employee_id" validate:"required"`
	Month      int  `json:"month" validate:"required,min=1,max=12"`
	Year       int  `json:"year" validate:"required,min=2000"`
}

type ApproveSalaryRequest struct {
	SalaryCalculationID uint `json:"salary_calculation_id" validate:"required"`
}
