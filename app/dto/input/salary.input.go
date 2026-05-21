package input

type CalculateSalaryRequest struct {
	EmployeeID uint   `json:"employee_id" validate:"required"`
	FromDate   string `json:"from_date" validate:"required,datetime=2006-01-02"`
	ToDate     string `json:"to_date" validate:"required,datetime=2006-01-02"`
}

type ApproveSalaryRequest struct {
	SalaryCalculationID uint `json:"salary_calculation_id" validate:"required"`
}
