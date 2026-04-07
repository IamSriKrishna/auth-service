package output

import (
	"time"
)

// ProductGroupInventoryOutput represents the inventory status of a product group
type ProductGroupInventoryOutput struct {
	ProductGroupID   string                       `json:"product_group_id"`
	ProductGroupName string                       `json:"product_group_name"`
	CurrentStock     float64                      `json:"current_stock"`
	AllocatedStock   float64                      `json:"allocated_stock"`
	AvailableStock   float64                      `json:"available_stock"`
	TotalReceived    float64                      `json:"total_received"`
	TotalSold        float64                      `json:"total_sold"`
	DamagedStock     float64                      `json:"damaged_stock"`
	ComponentStatus  []ComponentStockStatusOutput `json:"component_status"`
	LastUpdated      time.Time                    `json:"last_updated"`
}

type ComponentStockStatusOutput struct {
	ComponentProductID   string  `json:"component_product_id"`
	ComponentProductName string  `json:"component_product_name"`
	VariantSku           *string `json:"variant_sku,omitempty"`
	QuantityPerGroup     float64 `json:"quantity_per_group"`
	CurrentStock         float64 `json:"current_stock"`
	AllocatedStock       float64 `json:"allocated_stock"`
	AvailableStock       float64 `json:"available_stock"`
	CanFulfillGroup      bool    `json:"can_fulfill_group"` // Whether this component can fulfill one product group
	Percentage           float64 `json:"percentage"`        // Percentage of required stock available
}

// ProductGroupTransactionOutput represents a single transaction
type ProductGroupTransactionOutput struct {
	ID               uint      `json:"id"`
	ProductGroupID   string    `json:"product_group_id"`
	ProductGroupName string    `json:"product_group_name"`
	TransactionType  string    `json:"transaction_type"` // purchase, sales, return, adjustment, damage
	Quantity         float64   `json:"quantity"`
	ReferenceID      string    `json:"reference_id,omitempty"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
}

type ProductGroupTransactionListOutput struct {
	Transactions []ProductGroupTransactionOutput `json:"transactions"`
	Total        int64                           `json:"total"`
	Page         int                             `json:"page"`
	Limit        int                             `json:"limit"`
}

// ProductGroupPurchaseOutput - Response when purchasing a product group
type ProductGroupPurchaseOutput struct {
	PurchaseOrderID     string                   `json:"purchase_order_id"`
	PurchaseOrderNumber string                   `json:"purchase_order_number"`
	ProductGroupID      string                   `json:"product_group_id"`
	ProductGroupName    string                   `json:"product_group_name"`
	VendorID            uint                     `json:"vendor_id"`
	VendorName          string                   `json:"vendor_name"`
	Quantity            float64                  `json:"quantity"` // Number of complete product groups
	LineItems           []PurchaseLineItemOutput `json:"line_items"`
	TotalAmount         float64                  `json:"total_amount"`
	Status              string                   `json:"status"`
	CreatedAt           time.Time                `json:"created_at"`
	Message             string                   `json:"message"`
}

// ProductGroupSalesOutput - Response when selling a product group
type ProductGroupSalesOutput struct {
	SalesOrderID     string                `json:"sales_order_id"`
	SalesOrderNumber string                `json:"sales_order_number"`
	ProductGroupID   string                `json:"product_group_id"`
	ProductGroupName string                `json:"product_group_name"`
	CustomerID       uint                  `json:"customer_id"`
	CustomerName     string                `json:"customer_name"`
	Quantity         float64               `json:"quantity"` // Number of complete product groups
	LineItems        []SalesLineItemOutput `json:"line_items"`
	TotalAmount      float64               `json:"total_amount"`
	StockAllocated   bool                  `json:"stock_allocated"`
	Status           string                `json:"status"`
	CreatedAt        time.Time             `json:"created_at"`
	Message          string                `json:"message"`
}

type PurchaseLineItemOutput struct {
	ID             uint    `json:"id"`
	ProductID      string  `json:"product_id,omitempty"`
	ProductName    string  `json:"product_name"`
	VariantSKU     string  `json:"variant_sku,omitempty"`
	Quantity       float64 `json:"quantity"`
	Rate           float64 `json:"rate"`
	Amount         float64 `json:"amount"`
	IsProductGroup bool    `json:"is_product_group"`
	ProductGroupID string  `json:"product_group_id,omitempty"`
}

type SalesLineItemOutput struct {
	ID             uint    `json:"id"`
	ProductID      string  `json:"product_id,omitempty"`
	ProductName    string  `json:"product_name"`
	VariantSKU     string  `json:"variant_sku,omitempty"`
	Quantity       float64 `json:"quantity"`
	Rate           float64 `json:"rate"`
	Amount         float64 `json:"amount"`
	IsProductGroup bool    `json:"is_product_group"`
	ProductGroupID string  `json:"product_group_id,omitempty"`
}

// ProductGroupShipmentOutput - Response when shipping product groups
type ProductGroupShipmentOutput struct {
	ShipmentID       string                   `json:"shipment_id"`
	SalesOrderID     string                   `json:"sales_order_id"`
	ProductGroupID   string                   `json:"product_group_id"`
	ProductGroupName string                   `json:"product_group_name"`
	Quantity         float64                  `json:"quantity"` // Number of complete product groups
	LineItems        []ShipmentLineItemOutput `json:"line_items"`
	ShipmentDate     time.Time                `json:"shipment_date"`
	TrackingNumber   string                   `json:"tracking_number,omitempty"`
	Status           string                   `json:"status"`
	Message          string                   `json:"message"`
}

type ShipmentLineItemOutput struct {
	ID              uint    `json:"id"`
	ProductName     string  `json:"product_name"`
	VariantSKU      string  `json:"variant_sku"`
	QuantityShipped float64 `json:"quantity_shipped"`
	Status          string  `json:"status,omitempty"`
}

// ProductGroupInvoiceOutput - Response when creating invoice for product groups
type ProductGroupInvoiceOutput struct {
	InvoiceID        string                  `json:"invoice_id"`
	InvoiceNumber    string                  `json:"invoice_number"`
	SalesOrderID     string                  `json:"sales_order_id"`
	ProductGroupID   string                  `json:"product_group_id"`
	ProductGroupName string                  `json:"product_group_name"`
	CustomerID       uint                    `json:"customer_id"`
	CustomerName     string                  `json:"customer_name"`
	Quantity         float64                 `json:"quantity"`
	LineItems        []InvoiceLineItemOutput `json:"line_items"`
	SubTotal         float64                 `json:"sub_total"`
	TaxAmount        float64                 `json:"tax_amount"`
	Total            float64                 `json:"total"`
	InvoiceDate      time.Time               `json:"invoice_date"`
	DueDate          time.Time               `json:"due_date"`
	Status           string                  `json:"status"`
	Message          string                  `json:"message"`
}

// ProductGroupBillOutput - Response when creating bill for product groups
type ProductGroupBillOutput struct {
	BillID           string               `json:"bill_id"`
	BillNumber       string               `json:"bill_number"`
	PurchaseOrderID  string               `json:"purchase_order_id"`
	ProductGroupID   string               `json:"product_group_id"`
	ProductGroupName string               `json:"product_group_name"`
	VendorID         uint                 `json:"vendor_id"`
	VendorName       string               `json:"vendor_name"`
	Quantity         float64              `json:"quantity"`
	LineItems        []BillLineItemOutput `json:"line_items"`
	SubTotal         float64              `json:"sub_total"`
	TaxAmount        float64              `json:"tax_amount"`
	Total            float64              `json:"total"`
	BillDate         time.Time            `json:"bill_date"`
	DueDate          time.Time            `json:"due_date"`
	Status           string               `json:"status"`
	Message          string               `json:"message"`
}

// ProductGroupStockManagementOutput - Comprehensive stock management response
type ProductGroupStockManagementOutput struct {
	ProductGroupID     string                          `json:"product_group_id"`
	ProductGroupName   string                          `json:"product_group_name"`
	StockStatus        ProductGroupInventoryOutput     `json:"stock_status"`
	RecentTransactions []ProductGroupTransactionOutput `json:"recent_transactions"`
	Message            string                          `json:"message"`
}

// ErrorResponseOutput for stock-related errors
type ProductGroupStockErrorOutput struct {
	ErrorCode         string  `json:"error_code"`
	Message           string  `json:"message"`
	Details           string  `json:"details"`
	ProductGroupID    string  `json:"product_group_id,omitempty"`
	RequiredQuantity  float64 `json:"required_quantity,omitempty"`
	AvailableQuantity float64 `json:"available_quantity,omitempty"`
}
