package input

type CreateEmployeeRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"omitempty,email"`
	Number       string `json:"number" validate:"required"`
	Address      string `json:"address" validate:"required"`
	EmployeeType string `json:"employee_type" validate:"required,oneof=part-time full-time"`
}

type UpdateEmployeeRequest struct {
	Name         *string `json:"name"`
	Email        *string `json:"email" validate:"omitempty,email"`
	Number       *string `json:"number"`
	Address      *string `json:"address"`
	EmployeeType *string `json:"employee_type" validate:"omitempty,oneof=part-time full-time"`
}
