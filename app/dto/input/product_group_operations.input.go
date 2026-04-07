package input

// ProductGroupPurchaseInput - Create purchase order for product groups
type ProductGroupPurchaseInput struct {
	VendorID            uint    `json:"vendor_id" validate:"required"`
	ProductGroupID      string  `json:"product_group_id" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required,gt=0"`
	RatePerProductGroup float64 `json:"rate_per_product_group" validate:"required,gte=0"`
	DeliveryDate        string  `json:"delivery_date" validate:"required"`
	PaymentTerms        string  `json:"payment_terms" validate:"required"`
	ReferenceNo         string  `json:"reference_no"`
	Notes               string  `json:"notes"`
	TermsAndConditions  string  `json:"terms_and_conditions"`
}

// ProductGroupSalesInput - Create sales order for product groups
type ProductGroupSalesInput struct {
	CustomerID           uint    `json:"customer_id" validate:"required"`
	ProductGroupID       string  `json:"product_group_id" validate:"required"`
	Quantity             float64 `json:"quantity" validate:"required,gt=0"`
	RatePerProductGroup  float64 `json:"rate_per_product_group" validate:"required,gte=0"`
	ExpectedShipmentDate string  `json:"expected_shipment_date" validate:"required"`
	PaymentTerms         string  `json:"payment_terms" validate:"required"`
	DeliveryMethod       string  `json:"delivery_method"`
	ReferenceNo          string  `json:"reference_no"`
	CustomerNotes        string  `json:"customer_notes"`
	TermsAndConditions   string  `json:"terms_and_conditions"`
	AutoAllocateStock    bool    `json:"auto_allocate_stock"` // If true, automatically allocate stock
}

// ProductGroupAllocateStockInput - Allocate stock for a sales order
type ProductGroupAllocateStockInput struct {
	ProductGroupID string  `json:"product_group_id" validate:"required"`
	SalesOrderID   string  `json:"sales_order_id" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
}

// ProductGroupDeductStockInput - Deduct stock after shipment
type ProductGroupDeductStockInput struct {
	ProductGroupID string  `json:"product_group_id" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
	ShipmentID     string  `json:"shipment_id"`
	Reason         string  `json:"reason"` // e.g., "shipment", "damaged", "return"
}

// ProductGroupReleaseStockInput - Release previously allocated stock
type ProductGroupReleaseStockInput struct {
	ProductGroupID string `json:"product_group_id" validate:"required"`
	SalesOrderID   string `json:"sales_order_id" validate:"required"`
}

// ProductGroupShipmentInput - Create shipment for product groups
type ProductGroupShipmentInput struct {
	SalesOrderID      string  `json:"sales_order_id" validate:"required"`
	ProductGroupID    string  `json:"product_group_id" validate:"required"`
	Quantity          float64 `json:"quantity" validate:"required,gt=0"`
	TrackingNumber    string  `json:"tracking_number"`
	ShippingAddress   string  `json:"shipping_address"`
	EstimatedDelivery string  `json:"estimated_delivery"`
	Notes             string  `json:"notes"`
	AutoDeductStock   bool    `json:"auto_deduct_stock"` // If true, automatically deduct stock
}

// ProductGroupInvoiceInput - Create invoice for product groups
type ProductGroupInvoiceInput struct {
	SalesOrderID        string  `json:"sales_order_id" validate:"required"`
	ProductGroupID      string  `json:"product_group_id" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required,gt=0"`
	RatePerProductGroup float64 `json:"rate_per_product_group" validate:"required,gte=0"`
	InvoiceDate         string  `json:"invoice_date"` // YYYY-MM-DD
	DueDate             string  `json:"due_date"`     // YYYY-MM-DD
	TaxType             string  `json:"tax_type"`     // e.g., "GST", "VAT"
	TaxPercentage       float64 `json:"tax_percentage"`
	Notes               string  `json:"notes"`
}

// ProductGroupBillInput - Create bill for product groups
type ProductGroupBillInput struct {
	PurchaseOrderID     string  `json:"purchase_order_id" validate:"required"`
	ProductGroupID      string  `json:"product_group_id" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required,gt=0"`
	RatePerProductGroup float64 `json:"rate_per_product_group" validate:"required,gte=0"`
	BillDate            string  `json:"bill_date"` // YYYY-MM-DD
	DueDate             string  `json:"due_date"`  // YYYY-MM-DD
	TaxType             string  `json:"tax_type"`  // e.g., "GST", "VAT"
	TaxPercentage       float64 `json:"tax_percentage"`
	Notes               string  `json:"notes"`
}

// ProductGroupCheckStockInput - Check if stock is available
type ProductGroupCheckStockInput struct {
	ProductGroupID string  `json:"product_group_id" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
}

// ProductGroupStockAdjustmentInput - Manual stock adjustment
type ProductGroupStockAdjustmentInput struct {
	ProductGroupID string  `json:"product_group_id" validate:"required"`
	AdjustmentType string  `json:"adjustment_type" validate:"required"` // "add", "remove", "damage"
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
	Reason         string  `json:"reason" validate:"required"`
	ReferenceNo    string  `json:"reference_no"`
}

// ProductGroupTransferStockInput - Transfer stock between groups
type ProductGroupTransferStockInput struct {
	FromProductGroupID string  `json:"from_product_group_id" validate:"required"`
	ToProductGroupID   string  `json:"to_product_group_id" validate:"required"`
	Quantity           float64 `json:"quantity" validate:"required,gt=0"`
	Reason             string  `json:"reason"`
}

// ProductGroupInventoryReportInput - Request inventory report
type ProductGroupInventoryReportInput struct {
	ProductGroupID string `json:"product_group_id" validate:"required"`
	IncludeHistory bool   `json:"include_history"` // Include transaction history
	HistoryLimit   int    `json:"history_limit"`   // Number of recent transactions
}

// ProductGroupStockAlertInput - Set low stock alerts
type ProductGroupStockAlertInput struct {
	ProductGroupID          string  `json:"product_group_id" validate:"required"`
	LowStockThreshold       float64 `json:"low_stock_threshold" validate:"required,gt=0"`
	CriticalStockThreshold  float64 `json:"critical_stock_threshold" validate:"required,gt=0"`
	DisableAutoRestockAlert bool    `json:"disable_auto_restock_alert"`
}
