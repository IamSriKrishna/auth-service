package models

import "time"

type InventoryBalance struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement"`
	ItemID              string     `gorm:"type:varchar(255);index;not null"`
	Item                *Item      `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantSKU          *string    `gorm:"type:varchar(255);index"`
	CurrentQuantity     float64    `json:"current_quantity" gorm:"type:decimal(18,2);default:0"`
	ReservedQuantity    float64    `json:"reserved_quantity" gorm:"type:decimal(18,2);default:0"`
	AvailableQuantity   float64    `json:"available_quantity" gorm:"type:decimal(18,2);default:0"`
	InTransitQuantity   float64    `json:"in_transit_quantity" gorm:"type:decimal(18,2);default:0"`
	AverageRate         float64    `json:"average_rate" gorm:"default:0"`
	LastReceivedDate    *time.Time `json:"last_received_date"`
	LastConsumedDate    *time.Time `json:"last_consumed_date"`
	LastSoldDate        *time.Time `json:"last_sold_date"`
	LastInventorySyncAt time.Time  `json:"last_inventory_sync_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (InventoryBalance) TableName() string {
	return "inventory_balances"
}

type InventoryAggregation struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	ItemID             string    `gorm:"type:varchar(255);index;not null"`
	Item               *Item     `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantSKU         *string   `gorm:"type:varchar(255);index"`
	TotalPurchased     float64   `json:"total_purchased" gorm:"default:0"`
	TotalManufactured  float64   `json:"total_manufactured" gorm:"default:0"`
	TotalConsumedInMfg float64   `json:"total_consumed_in_mfg" gorm:"default:0"`
	TotalSold          float64   `json:"total_sold" gorm:"default:0"`
	AverageRate        float64   `json:"average_rate" gorm:"default:0"`
	CalculatedAt       time.Time `json:"calculated_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (InventoryAggregation) TableName() string {
	return "inventory_aggregations"
}

type InventoryJournal struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	ItemID          string    `gorm:"type:varchar(255);index;not null"`
	VariantSKU      *string   `gorm:"type:varchar(255);index"`
	TransactionType string    `json:"transaction_type" gorm:"type:varchar(50);not null"`
	Quantity        float64   `json:"quantity" gorm:"not null"`
	ReferenceType   string    `json:"reference_type" gorm:"type:varchar(50)"`
	ReferenceID     string    `json:"reference_id" gorm:"type:varchar(255);index"`
	ReferenceNo     string    `json:"reference_no" gorm:"type:varchar(100)"`
	Notes           string    `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by" gorm:"type:varchar(255)"`
}

func (InventoryJournal) TableName() string {
	return "inventory_journals"
}

type SupplyChainSummary struct {
	ID                           uint      `gorm:"primaryKey;autoIncrement"`
	ItemID                       string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Item                         *Item     `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantSKU                   *string   `gorm:"type:varchar(255);uniqueIndex:,composite:variant_item"`
	Variant                      *Variant  `json:"variant,omitempty" gorm:"foreignKey:VariantSKU;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	OpeningStock                 float64   `json:"opening_stock" gorm:"default:0"`
	TotalPurchaseOrderQuantity   float64   `json:"total_po_quantity" gorm:"default:0"`
	TotalPurchaseOrderAmount     float64   `json:"total_po_amount" gorm:"default:0"`
	AveragePurchaseRate          float64   `json:"avg_purchase_rate" gorm:"default:0"`
	TotalProductionOrderQuantity float64   `json:"total_prod_qty" gorm:"default:0"`
	TotalManufacturedQuantity    float64   `json:"total_mfg_qty" gorm:"default:0"`
	TotalConsumedInProduction    float64   `json:"total_consumed_in_mfg" gorm:"default:0"`
	TotalSalesOrderQuantity      float64   `json:"total_so_quantity" gorm:"default:0"`
	TotalSalesOrderAmount        float64   `json:"total_so_amount" gorm:"default:0"`
	AverageSalesRate             float64   `json:"avg_sales_rate" gorm:"default:0"`
	TotalInvoicedQuantity        float64   `json:"total_invoiced_qty" gorm:"default:0"`
	CurrentQuantity              float64   `json:"current_qty" gorm:"default:0"`
	UpdatedAt                    time.Time `json:"updated_at"`
}

func (SupplyChainSummary) TableName() string {
	return "supply_chain_summary"
}
