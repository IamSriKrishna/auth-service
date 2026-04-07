package models

import "time"

// ProductStock tracks inventory for products (not items/variants)
// This is the master stock ledger for all products in the system
type ProductStock struct {
	ID          string   `gorm:"type:varchar(255);primaryKey" json:"id"`
	ProductID   string   `gorm:"type:varchar(255);uniqueIndex:idx_product_id;not null" json:"product_id"`
	Product     *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductName string   `gorm:"type:varchar(255)" json:"product_name"`
	SKU         string   `gorm:"type:varchar(100);index" json:"sku"`

	// Stock quantities
	CurrentStock   float64 `gorm:"default:0;not null" json:"current_stock"`   // Total units on hand
	PurchasedStock float64 `gorm:"default:0;not null" json:"purchased_stock"` // Total units purchased (all time)
	SoldStock      float64 `gorm:"default:0;not null" json:"sold_stock"`      // Total units sold (all time)
	ReservedStock  float64 `gorm:"default:0;not null" json:"reserved_stock"`  // Units reserved for sales orders
	AvailableStock float64 `gorm:"default:0;not null" json:"available_stock"` // Available = Current - Reserved

	// Accounting
	AverageCost       float64 `gorm:"default:0" json:"average_cost"`       // Weighted average purchase cost
	RevaluationAmount float64 `gorm:"default:0" json:"revaluation_amount"` // Stock revaluation adjustments

	// Tracking
	LastPurchasedDate *time.Time `json:"last_purchased_date,omitempty"`
	LastSoldDate      *time.Time `json:"last_sold_date,omitempty"`
	LastStockSyncAt   time.Time  `json:"last_stock_sync_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductStock) TableName() string {
	return "product_stocks"
}

// StockLedger tracks all stock movements (inbound/outbound) for audit trail
type StockLedger struct {
	ID        uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID string   `gorm:"type:varchar(255);index;not null" json:"product_id"`
	Product   *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Movement details
	MovementType string  `gorm:"type:varchar(50);index;not null" json:"movement_type"` // PURCHASE, SALES, ADJUSTMENT, etc.
	Quantity     float64 `gorm:"not null" json:"quantity"`                             // Can be negative for outbound
	Rate         float64 `gorm:"default:0" json:"rate"`                                // Cost per unit
	Amount       float64 `gorm:"not null" json:"amount"`                               // Qty * Rate

	// Reference information
	ReferenceType   string `gorm:"type:varchar(100)" json:"reference_type"`         // purchase_order, sales_order, bill, invoice, etc.
	ReferenceID     string `gorm:"type:varchar(255);index" json:"reference_id"`     // ID of the source document
	ReferenceNumber string `gorm:"type:varchar(100);index" json:"reference_number"` // PO-2024-001, SO-2024-001, etc.

	// Balance tracking
	BalanceBeforeQty float64 `json:"balance_before_qty"`
	BalanceAfterQty  float64 `json:"balance_after_qty"`
	CostBeforeAmount float64 `json:"cost_before_amount"`
	CostAfterAmount  float64 `json:"cost_after_amount"`

	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(255)" json:"created_by"`
}

func (StockLedger) TableName() string {
	return "stock_ledgers"
}
