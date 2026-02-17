package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type ProductionOrder struct {
	ID                    string                       `json:"id" gorm:"type:varchar(255);primaryKey"`
	ProductionOrderNumber string                       `json:"production_order_no" gorm:"type:varchar(100);uniqueIndex;not null"`
	ItemGroupID           string                       `json:"item_group_id" gorm:"type:varchar(255);not null;index"`
	ItemGroup             *ItemGroup                   `json:"item_group,omitempty" gorm:"foreignKey:ItemGroupID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	QuantityToManufacture float64                      `json:"quantity_to_manufacture" gorm:"not null"`
	QuantityManufactured  float64                      `json:"quantity_manufactured" gorm:"default:0"`
	Status                domain.ProductionOrderStatus `json:"status" gorm:"type:varchar(50);not null;default:'planned'"`
	PlannedStartDate      time.Time                    `json:"planned_start_date"`
	PlannedEndDate        time.Time                    `json:"planned_end_date"`
	ActualStartDate       *time.Time                   `json:"actual_start_date"`
	ActualEndDate         *time.Time                   `json:"actual_end_date"`
	ManufacturedDate      *time.Time                   `json:"manufactured_date"`
	InventorySynced       bool                         `json:"inventory_synced" gorm:"default:false;index"`
	InventorySyncDate     *time.Time                   `json:"inventory_sync_date"`
	Notes                 string                       `json:"notes" gorm:"type:text"`
	ProductionOrderItems  []ProductionOrderItem        `json:"production_order_items" gorm:"foreignKey:ProductionOrderID;constraint:OnDelete:CASCADE"`
	CreatedAt             time.Time                    `json:"created_at"`
	UpdatedAt             time.Time                    `json:"updated_at"`
	CreatedBy             string                       `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy             string                       `json:"updated_by" gorm:"type:varchar(255)"`
}

func (ProductionOrder) TableName() string {
	return "production_orders"
}

type ProductionOrderItem struct {
	ID                   uint                `gorm:"primaryKey;autoIncrement"`
	ProductionOrderID    string              `json:"production_order_id" gorm:"type:varchar(255);index;not null"`
	ItemGroupComponentID uint                `json:"item_group_component_id" gorm:"index;not null"`
	ItemGroupComponent   *ItemGroupComponent `json:"item_group_component,omitempty" gorm:"foreignKey:ItemGroupComponentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	QuantityRequired     float64             `json:"quantity_required" gorm:"not null"`
	QuantityConsumed     float64             `json:"quantity_consumed" gorm:"default:0"`
	InventorySynced      bool                `json:"inventory_synced" gorm:"default:false"`
	SyncedAt             *time.Time          `json:"synced_at"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

func (ProductionOrderItem) TableName() string {
	return "production_order_items"
}
