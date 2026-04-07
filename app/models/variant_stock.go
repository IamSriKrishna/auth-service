package models

import "time"

// VariantStock tracks inventory for product variants
// This enables granular stock tracking at the SKU/variant level
type VariantStock struct {
	ID          string   `gorm:"type:varchar(255);primaryKey" json:"id"`
	ProductID   string   `gorm:"type:varchar(255);index;not null" json:"product_id"`
	Product     *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	VariantSKU  string   `gorm:"type:varchar(100);uniqueIndex;not null" json:"variant_sku"`
	VariantName string   `gorm:"type:varchar(255)" json:"variant_name"`
	ProductName string   `gorm:"type:varchar(255)" json:"product_name"`

	// Stock quantities
	CurrentStock   float64 `gorm:"default:0;not null" json:"current_stock"`    // Total units on hand
	PurchasedStock float64 `gorm:"default:0;not null" json:"purchased_stock"`  // Total units purchased (all time)
	SoldStock      float64 `gorm:"default:0;not null" json:"sold_stock"`       // Total units sold (all time)
	ReservedStock  float64 `gorm:"default:0;not null" json:"reserved_stock"`   // Units reserved for sales orders
	AvailableStock float64 `gorm:"default:0;not null" json:"available_stock"`  // Available = Current - Reserved
	InTransitStock float64 `gorm:"default:0;not null" json:"in_transit_stock"` // Stock in shipment process

	// Reorder management
	ReorderLevel float64 `json:"reorder_level" gorm:"default:0"`    // Minimum stock level
	ReorderQty   float64 `json:"reorder_qty" gorm:"default:0"`      // Quantity to order when below reorder level
	IsLowStock   bool    `json:"is_low_stock" gorm:"default:false"` // Flag for low stock alert

	// Accounting
	AverageCost       float64 `gorm:"default:0" json:"average_cost"`       // Weighted average purchase cost
	SellingPrice      float64 `gorm:"default:0" json:"selling_price"`      // Current selling price
	RevaluationAmount float64 `gorm:"default:0" json:"revaluation_amount"` // Stock revaluation adjustments

	// Tracking
	LastPurchasedDate *time.Time `json:"last_purchased_date,omitempty"`
	LastSoldDate      *time.Time `json:"last_sold_date,omitempty"`
	LastStockSyncAt   *time.Time `json:"last_stock_sync_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (VariantStock) TableName() string {
	return "variant_stocks"
}

// VariantStockMovement tracks detailed movements for variants
type VariantStockMovement struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	VariantID  string `gorm:"type:varchar(255);index;not null" json:"variant_id"`
	ProductID  string `gorm:"type:varchar(255);index;not null" json:"product_id"`
	VariantSKU string `gorm:"type:varchar(100);index;not null" json:"variant_sku"`

	// Movement details
	MovementType string  `gorm:"type:varchar(50);index;not null" json:"movement_type"` // PURCHASE, SALES, SHIPMENT, etc.
	Quantity     float64 `gorm:"not null" json:"quantity"`                             // Can be negative for outbound
	Rate         float64 `gorm:"default:0" json:"rate"`                                // Cost per unit
	Amount       float64 `gorm:"not null" json:"amount"`                               // Qty * Rate

	// Reference information
	ReferenceType   string `gorm:"type:varchar(100)" json:"reference_type"`         // purchase_order, sales_order, bill, invoice, shipment, package
	ReferenceID     string `gorm:"type:varchar(255);index" json:"reference_id"`     // ID of the source document
	ReferenceNumber string `gorm:"type:varchar(100);index" json:"reference_number"` // PO-001, SO-001, INV-001, etc.

	// Balance tracking
	BalanceBeforeQty float64 `json:"balance_before_qty"`
	BalanceAfterQty  float64 `json:"balance_after_qty"`

	// Document stages
	Stage string `gorm:"type:varchar(50)" json:"stage"` // draft, confirmed, shipped, delivered, invoiced

	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(255)" json:"created_by"`
}

func (VariantStockMovement) TableName() string {
	return "variant_stock_movements"
}

// StockReservation tracks reserved stock for sales orders
type StockReservation struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	SalesOrderID    string     `gorm:"type:varchar(255);index;not null" json:"sales_order_id"`
	SalesOrderNo    string     `gorm:"type:varchar(100)" json:"sales_order_no"`
	ProductID       string     `gorm:"type:varchar(255);index;not null" json:"product_id"`
	VariantSKU      string     `gorm:"type:varchar(100);index;not null" json:"variant_sku"`
	VariantStockID  string     `gorm:"type:varchar(255);index" json:"variant_stock_id"`
	ReservedQty     float64    `gorm:"not null" json:"reserved_qty"`
	ShippedQty      float64    `gorm:"default:0" json:"shipped_qty"`
	InvoicedQty     float64    `gorm:"default:0" json:"invoiced_qty"`
	Status          string     `gorm:"type:varchar(50)" json:"status"` // reserved, partial_shipped, fully_shipped, invoiced, cancelled
	ReservationDate time.Time  `json:"reservation_date"`
	ShipByDate      *time.Time `json:"ship_by_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedBy       string     `gorm:"type:varchar(255)" json:"created_by"`
}

func (StockReservation) TableName() string {
	return "stock_reservations"
}
