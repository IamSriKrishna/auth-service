package output

import "time"

type ManufacturerOutput struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	ProductGroupID string                     `json:"product_group_id"`
	ProductGroup   *ProductGroupOutput        `json:"product_group,omitempty"`
	Quantity       float64                    `json:"quantity"`
	Status         string                     `json:"status"`
	Description    string                     `json:"description"`
	Employees      []EmployeeAssignmentOutput `json:"employees,omitempty"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

type EmployeeAssignmentOutput struct {
	EmployeeID  uint            `json:"employee_id"`
	Employee    *EmployeeOutput `json:"employee,omitempty"`
	ServiceCost float64         `json:"service_cost"`
	CostType    string          `json:"cost_type"`
}

type ListManufacturersOutput struct {
	Manufacturers []ManufacturerOutput `json:"manufacturers"`
	TotalCount    int                  `json:"total_count"`
}
