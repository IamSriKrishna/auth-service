package output

import "time"

type ProductionOrderOutput struct {
	ID                    string                      `json:"id"`
	ProductionOrderNo     string                      `json:"production_order_no"`
	ItemGroupID           string                      `json:"item_group_id"`
	ItemGroupName         string                      `json:"item_group_name"`
	QuantityToManufacture float64                     `json:"quantity_to_manufacture"`
	QuantityManufactured  float64                     `json:"quantity_manufactured"`
	Status                string                      `json:"status"`
	PlannedStartDate      time.Time                   `json:"planned_start_date"`
	PlannedEndDate        time.Time                   `json:"planned_end_date"`
	ActualStartDate       *time.Time                  `json:"actual_start_date"`
	ActualEndDate         *time.Time                  `json:"actual_end_date"`
	ManufacturedDate      *time.Time                  `json:"manufactured_date"`
	InventorySynced       bool                        `json:"inventory_synced"`
	Notes                 string                      `json:"notes"`
	ProductionOrderItems  []ProductionOrderItemOutput `json:"production_order_items"`
	CreatedAt             time.Time                   `json:"created_at"`
	UpdatedAt             time.Time                   `json:"updated_at"`
	Warnings              []string                    `json:"warnings,omitempty"`
}

type ProductionOrderItemOutput struct {
	ID                   uint    `json:"id"`
	ItemGroupComponentID uint    `json:"item_group_component_id"`
	ItemID               string  `json:"item_id"`
	ItemName             string  `json:"item_name"`
	VariantSku           *string `json:"variant_sku,omitempty"`
	QuantityRequired     float64 `json:"quantity_required"`
	QuantityConsumed     float64 `json:"quantity_consumed"`
	InventorySynced      bool    `json:"inventory_synced"`
}

type ProductionOrderListOutput struct {
	ProductionOrders []ProductionOrderListItemOutput `json:"data"`
	Total            int                             `json:"total"`
	Page             int                             `json:"page"`
	Limit            int                             `json:"limit"`
	TotalPages       int                             `json:"total_pages"`
}

type ProductionOrderListItemOutput struct {
	ID                    string    `json:"id"`
	ProductionOrderNo     string    `json:"production_order_no"`
	ItemGroupName         string    `json:"item_group_name"`
	QuantityToManufacture float64   `json:"quantity_to_manufacture"`
	QuantityManufactured  float64   `json:"quantity_manufactured"`
	Status                string    `json:"status"`
	PlannedStartDate      time.Time `json:"planned_start_date"`
	PlannedEndDate        time.Time `json:"planned_end_date"`
	CreatedAt             time.Time `json:"created_at"`
}

type ProductionOrderDeleteOutput struct {
	ID                string    `json:"id"`
	ProductionOrderNo string    `json:"production_order_no"`
	DeletedAt         time.Time `json:"deleted_at"`
}
