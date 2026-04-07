package models

import "time"

// ProductGroupInventory tracks stock levels for entire product groups
// When you have 100 units of a product group, it means you have enough
// of ALL components (bottle + cap + label) to make 100 complete kits
type ProductGroupInventory struct {
	ID             uint          `gorm:"primaryKey;autoIncrement"`
	ProductGroupID string        `json:"product_group_id" gorm:"type:varchar(255);index;not null;uniqueIndex"`
	ProductGroup   *ProductGroup `json:"product_group,omitempty" gorm:"foreignKey:ProductGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Current inventory level (quantity of complete product groups available)
	CurrentStock float64 `json:"current_stock" gorm:"default:0;not null"`

	// Allocated stock (reserved for sales orders but not yet shipped)
	AllocatedStock float64 `json:"allocated_stock" gorm:"default:0;not null"`

	// Available for sale (current - allocated)
	AvailableStock float64 `json:"available_stock" gorm:"default:0;not null"`

	// Total received from purchases
	TotalReceived float64 `json:"total_received" gorm:"default:0;not null"`

	// Total sold/shipped
	TotalSold float64 `json:"total_sold" gorm:"default:0;not null"`

	// For tracking damaged/expired stock
	DamagedStock float64 `json:"damaged_stock" gorm:"default:0;not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductGroupInventory) TableName() string {
	return "product_group_inventory"
}

// ComponentInventory represents inventory of individual components within a product group
type ComponentInventory struct {
	ID             uint          `gorm:"primaryKey;autoIncrement"`
	ProductGroupID string        `json:"product_group_id" gorm:"type:varchar(255);index;not null"`
	ProductGroup   *ProductGroup `json:"product_group,omitempty" gorm:"foreignKey:ProductGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ComponentProductID string   `json:"component_product_id" gorm:"type:varchar(255);not null;index"`
	ComponentProduct   *Product `json:"component_product,omitempty" gorm:"foreignKey:ComponentProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	ComponentVariantSku *string `json:"component_variant_sku,omitempty" gorm:"type:varchar(255);index"`

	// Quantity of this component needed per product group
	// e.g., 1 bottle per kit, 1 cap per kit, 1 label per kit
	QuantityPerGroup float64 `json:"quantity_per_group" gorm:"not null"`

	// Stock levels - tracked separately for each component
	CurrentStock   float64 `json:"current_stock" gorm:"default:0;not null"`
	AllocatedStock float64 `json:"allocated_stock" gorm:"default:0;not null"`
	AvailableStock float64 `json:"available_stock" gorm:"default:0;not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ComponentInventory) TableName() string {
	return "component_inventory"
}

// ProductGroupTransaction logs all stock movements
type ProductGroupTransaction struct {
	ID             uint          `gorm:"primaryKey;autoIncrement"`
	ProductGroupID string        `json:"product_group_id" gorm:"type:varchar(255);index;not null"`
	ProductGroup   *ProductGroup `json:"product_group,omitempty" gorm:"foreignKey:ProductGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	TransactionType string  `json:"transaction_type" gorm:"type:varchar(50);not null"` // "purchase", "sales", "return", "adjustment", "damage"
	Quantity        float64 `json:"quantity" gorm:"not null"`

	// Reference to source transaction
	PurchaseOrderID *string `json:"purchase_order_id,omitempty" gorm:"type:varchar(255);index"`
	SalesOrderID    *string `json:"sales_order_id,omitempty" gorm:"type:varchar(255);index"`
	ShipmentID      *string `json:"shipment_id,omitempty" gorm:"type:varchar(255);index"`
	BillID          *string `json:"bill_id,omitempty" gorm:"type:varchar(255);index"`
	InvoiceID       *string `json:"invoice_id,omitempty" gorm:"type:varchar(255);index"`

	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:varchar(255)"`
}

func (ProductGroupTransaction) TableName() string {
	return "product_group_transactions"
}
