package domain

type PaymentTerms string

const (
	PaymentTermsDueOnReceipt    PaymentTerms = "due_on_receipt"
	PaymentTermsNet15           PaymentTerms = "net_15"
	PaymentTermsNet30           PaymentTerms = "net_30"
	PaymentTermsNet45           PaymentTerms = "net_45"
	PaymentTermsNet60           PaymentTerms = "net_60"
	PaymentTermsDueEndOfMonth   PaymentTerms = "due_end_of_month"
	PaymentTermsDueEndNextMonth PaymentTerms = "due_end_next_month"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft   InvoiceStatus = "draft"
	InvoiceStatusSent    InvoiceStatus = "sent"
	InvoiceStatusPartial InvoiceStatus = "partial"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusOverdue InvoiceStatus = "overdue"
	InvoiceStatusVoid    InvoiceStatus = "void"
)

type TaxType string

const (
	TaxTypeTDS TaxType = "tds"
	TaxTypeTCS TaxType = "tcs"
)

type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft             PurchaseOrderStatus = "draft"
	PurchaseOrderStatusSent              PurchaseOrderStatus = "sent"
	PurchaseOrderStatusPartiallyReceived PurchaseOrderStatus = "partially_received"
	PurchaseOrderStatusReceived          PurchaseOrderStatus = "received"
	PurchaseOrderStatusCancelled         PurchaseOrderStatus = "cancelled"
)

type SalesOrderStatus string

const (
	SalesOrderStatusDraft       SalesOrderStatus = "draft"
	SalesOrderStatusSent        SalesOrderStatus = "sent"
	SalesOrderStatusConfirmed   SalesOrderStatus = "confirmed"
	SalesOrderStatusPartialShip SalesOrderStatus = "partial_shipped"
	SalesOrderStatusShipped     SalesOrderStatus = "shipped"
	SalesOrderStatusDelivered   SalesOrderStatus = "delivered"
	SalesOrderStatusCancelled   SalesOrderStatus = "cancelled"
)

type PackageStatus string

const (
	PackageStatusCreated   PackageStatus = "created"
	PackageStatusPacked    PackageStatus = "packed"
	PackageStatusShipped   PackageStatus = "shipped"
	PackageStatusDelivered PackageStatus = "delivered"
	PackageStatusCancelled PackageStatus = "cancelled"
)

type ShipmentStatus string

const (
	ShipmentStatusCreated   ShipmentStatus = "created"
	ShipmentStatusShipped   ShipmentStatus = "shipped"
	ShipmentStatusInTransit ShipmentStatus = "in_transit"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusCancelled ShipmentStatus = "cancelled"
)

type BillStatus string

const (
	BillStatusDraft   BillStatus = "draft"
	BillStatusSent    BillStatus = "sent"
	BillStatusPartial BillStatus = "partial"
	BillStatusPaid    BillStatus = "paid"
	BillStatusOverdue BillStatus = "overdue"
	BillStatusVoid    BillStatus = "void"
)

type ProductionOrderStatus string

const (
	ProductionOrderStatusPlanned    ProductionOrderStatus = "planned"
	ProductionOrderStatusInProgress ProductionOrderStatus = "in_progress"
	ProductionOrderStatusCompleted  ProductionOrderStatus = "completed"
	ProductionOrderStatusCancelled  ProductionOrderStatus = "cancelled"
)
