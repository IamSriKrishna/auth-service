package models

import "time"

// DashboardMetrics represents aggregated business metrics
type DashboardMetrics struct {
	ID                  string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	TotalCustomers      int       `json:"total_customers"`
	ActiveCustomers     int       `json:"active_customers"`
	TotalVendors        int       `json:"total_vendors"`
	ActiveVendors       int       `json:"active_vendors"`
	TotalItems          int       `json:"total_items"`
	TotalItemGroups     int       `json:"total_item_groups"`
	TotalStock          int64     `json:"total_stock"`
	LowStockItems       int       `json:"low_stock_items"`
	TotalShipments      int       `json:"total_shipments"`
	ShippedCount        int       `json:"shipped_count"`
	PendingShipments    int       `json:"pending_shipments"`
	TotalInvoices       int       `json:"total_invoices"`
	InvoiceAmount       float64   `json:"invoice_amount"`
	TotalSalesOrders    int       `json:"total_sales_orders"`
	SalesOrderAmount    float64   `json:"sales_order_amount"`
	TotalPurchaseOrders int       `json:"total_purchase_orders"`
	PurchaseOrderAmount float64   `json:"purchase_order_amount"`
	TotalPackages       int       `json:"total_packages"`
	PackagesShipped     int       `json:"packages_shipped"`
	PendingPackages     int       `json:"pending_packages"`
	LastUpdatedAt       time.Time `json:"last_updated_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (DashboardMetrics) TableName() string {
	return "dashboard_metrics"
}

// ShipmentTracking represents shipment tracking information
type ShipmentTracking struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	ShipmentID string    `json:"shipment_id" gorm:"index"`
	Status     string    `json:"status" gorm:"index"`
	Location   string    `json:"location"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	UpdatedBy  string    `json:"updated_by"`
	Notes      string    `json:"notes" gorm:"type:text"`
	Timestamp  time.Time `json:"timestamp" gorm:"index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ShipmentTracking) TableName() string {
	return "shipment_tracking"
}

// EntityCountHistory stores historical counts for trend analysis
type EntityCountHistory struct {
	ID           string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	EntityType   string    `json:"entity_type" gorm:"index"`
	Date         time.Time `json:"date" gorm:"index;type:date"`
	Count        int       `json:"count"`
	ActiveCount  int       `json:"active_count"`
	CreatedToday int       `json:"created_today"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (EntityCountHistory) TableName() string {
	return "entity_count_history"
}
