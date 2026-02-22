package input

type CreateProductionOrderInput struct {
	ItemGroupID           string  `json:"item_group_id" validate:"required"`
	QuantityToManufacture float64 `json:"quantity_to_manufacture" validate:"required,gt=0"`
	PlannedStartDate      string  `json:"planned_start_date" validate:"required"`
	PlannedEndDate        string  `json:"planned_end_date" validate:"required"`
	Notes                 string  `json:"notes"`
}

type UpdateProductionOrderInput struct {
	Status               string  `json:"status"`
	QuantityManufactured float64 `json:"quantity_manufactured"`
	ActualStartDate      *string `json:"actual_start_date,omitempty"`
	ActualEndDate        *string `json:"actual_end_date,omitempty"`
	ManufacturedDate     *string `json:"manufactured_date,omitempty"`
	Notes                string  `json:"notes,omitempty"`
}
type ConsumeProductionOrderItemInput struct {
	ProductionOrderItemID uint    `json:"production_order_item_id" validate:"required"`
	QuantityConsumed      float64 `json:"quantity_consumed" validate:"required,gt=0"`
	Notes                 string  `json:"notes,omitempty"`
}
