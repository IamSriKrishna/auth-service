package input

type CreateEmployeeRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"omitempty,email"`
	Number  string `json:"number" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type UpdateEmployeeRequest struct {
	Name    *string `json:"name"`
	Email   *string `json:"email" validate:"omitempty,email"`
	Number  *string `json:"number"`
	Address *string `json:"address"`
}
