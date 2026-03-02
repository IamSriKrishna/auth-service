package output

import "time"

// DashboardMetricsOutput represents the main dashboard view
type DashboardMetricsOutput struct {
	CustomerMetrics      CustomerMetricsOutput      `json:"customer_metrics"`
	VendorMetrics        VendorMetricsOutput        `json:"vendor_metrics"`
	ItemMetrics          ItemMetricsOutput          `json:"item_metrics"`
	ShipmentMetrics      ShipmentMetricsOutput      `json:"shipment_metrics"`
	InvoiceMetrics       InvoiceMetricsOutput       `json:"invoice_metrics"`
	SalesOrderMetrics    SalesOrderMetricsOutput    `json:"sales_order_metrics"`
	PurchaseOrderMetrics PurchaseOrderMetricsOutput `json:"purchase_order_metrics"`
	PackageMetrics       PackageMetricsOutput       `json:"package_metrics"`
	LastUpdatedAt        time.Time                  `json:"last_updated_at"`
	GeneratedAt          time.Time                  `json:"generated_at"`
}

type CustomerMetricsOutput struct {
	Total        int `json:"total"`
	Active       int `json:"active"`
	Inactive     int `json:"inactive"`
	CreatedToday int `json:"created_today"`
}

type VendorMetricsOutput struct {
	Total        int `json:"total"`
	Active       int `json:"active"`
	Inactive     int `json:"inactive"`
	CreatedToday int `json:"created_today"`
}

type ItemMetricsOutput struct {
	Total          int   `json:"total"`
	TotalStock     int64 `json:"total_stock"`
	LowStockItems  int   `json:"low_stock_items"`
	ItemGroups     int   `json:"item_groups"`
	CreatedToday   int   `json:"created_today"`
	OutOfStockItem int   `json:"out_of_stock_items"`
}

type ShipmentMetricsOutput struct {
	Total            int     `json:"total"`
	Shipped          int     `json:"shipped"`
	Pending          int     `json:"pending"`
	InTransit        int     `json:"in_transit"`
	Delivered        int     `json:"delivered"`
	CancelledShipped int     `json:"cancelled_shipped"`
	AverageTime      float64 `json:"average_delivery_time_days"`
}

type InvoiceMetricsOutput struct {
	Total       int     `json:"total"`
	TotalAmount float64 `json:"total_amount"`
	Outstanding float64 `json:"outstanding_amount"`
	Paid        int     `json:"paid_count"`
	Pending     int     `json:"pending_count"`
	Overdue     int     `json:"overdue_count"`
}

type SalesOrderMetricsOutput struct {
	Total        int     `json:"total"`
	TotalAmount  float64 `json:"total_amount"`
	Completed    int     `json:"completed_count"`
	Pending      int     `json:"pending_count"`
	Cancelled    int     `json:"cancelled_count"`
	CreatedToday int     `json:"created_today"`
}

type PurchaseOrderMetricsOutput struct {
	Total        int     `json:"total"`
	TotalAmount  float64 `json:"total_amount"`
	Completed    int     `json:"completed_count"`
	Pending      int     `json:"pending_count"`
	Cancelled    int     `json:"cancelled_count"`
	CreatedToday int     `json:"created_today"`
}

type PackageMetricsOutput struct {
	Total        int `json:"total"`
	Shipped      int `json:"shipped_count"`
	Pending      int `json:"pending_count"`
	InTransit    int `json:"in_transit_count"`
	Delivered    int `json:"delivered_count"`
	CreatedToday int `json:"created_today"`
}

// ShipmentTrackingOutput represents shipment tracking details
type ShipmentTrackingOutput struct {
	ID         string    `json:"id"`
	ShipmentID string    `json:"shipment_id"`
	Status     string    `json:"status"`
	Location   string    `json:"location"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Notes      string    `json:"notes"`
	Timestamp  time.Time `json:"timestamp"`
}

// ShipmentTrackingListOutput represents a list of shipment tracking details
type ShipmentTrackingListOutput struct {
	Data  []ShipmentTrackingOutput `json:"data"`
	Total int                      `json:"total"`
}

// EntityTrendOutput represents trend data for an entity
type EntityTrendOutput struct {
	EntityType string       `json:"entity_type"`
	Data       []TrendPoint `json:"data"`
}

type TrendPoint struct {
	Date        time.Time `json:"date"`
	Count       int       `json:"count"`
	ActiveCount int       `json:"active_count"`
	NewToday    int       `json:"created_today"`
}

// EntityTrendsListOutput represents trends for multiple entities
type EntityTrendsListOutput struct {
	Trends    []EntityTrendOutput `json:"trends"`
	Timeframe string              `json:"timeframe"`
}

// StockDetailOutput represents detailed stock information
type StockDetailOutput struct {
	ItemID            string  `json:"item_id"`
	ItemName          string  `json:"item_name"`
	CurrentQuantity   float64 `json:"current_quantity"`
	AvailableQuantity float64 `json:"available_quantity"`
	ReservedQuantity  float64 `json:"reserved_quantity"`
	InTransitQuantity float64 `json:"in_transit_quantity"`
	Status            string  `json:"status"` // in_stock, low_stock, out_of_stock
}

// StockListOutput represents list of stock items
type StockListOutput struct {
	Data          []StockDetailOutput `json:"data"`
	TotalItems    int                 `json:"total_items"`
	InStock       int                 `json:"in_stock_count"`
	LowStock      int                 `json:"low_stock_count"`
	OutOfStock    int                 `json:"out_of_stock_count"`
	TotalQuantity int64               `json:"total_quantity"`
}

// ActivitySummaryOutput represents recent activity
type ActivitySummaryOutput struct {
	CreatedCustomersToday      int `json:"created_customers_today"`
	CreatedVendorsToday        int `json:"created_vendors_today"`
	CreatedItemsToday          int `json:"created_items_today"`
	CreatedSalesOrdersToday    int `json:"created_sales_orders_today"`
	CreatedPurchaseOrdersToday int `json:"created_purchase_orders_today"`
	ShippedToday               int `json:"shipped_today"`
	DeliveredToday             int `json:"delivered_today"`
}

// DiagnosticReportOutput for identifying data issues
type DiagnosticReportOutput struct {
	DataIssues  []string                  `json:"data_issues"`
	Diagnostics map[string]DiagnosticItem `json:"diagnostics"`
	Summary     string                    `json:"summary"`
}

type DiagnosticItem struct {
	Label       string      `json:"label"`
	Value       interface{} `json:"value"`
	Status      string      `json:"status"` // "ok", "warning", "error"
	Description string      `json:"description"`
}
