package input

type CreateManufacturerInput struct {
	Name           string                    `json:"name" validate:"required"`
	ProductGroupID string                    `json:"product_group_id" validate:"required"`
	Quantity       float64                   `json:"quantity" validate:"required,gt=0"`
	Description    string                    `json:"description"`
	Employees      []EmployeeAssignmentInput `json:"employees" validate:"required,min=1"`
}

type EmployeeAssignmentInput struct {
	EmployeeID  uint    `json:"employee_id" validate:"required"`
	ServiceCost float64 `json:"service_cost" validate:"required,gt=0"`
	CostType    string  `json:"cost_type" validate:"required,oneof=fixed per_unit"` // 'fixed' or 'per_unit'
}

type UpdateManufacturerInput struct {
	Name        *string  `json:"name"`
	Quantity    *float64 `json:"quantity"`
	Status      *string  `json:"status"`
	Description *string  `json:"description"`
}
