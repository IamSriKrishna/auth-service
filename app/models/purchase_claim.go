package models

import "time"

type PurchaseClaimType string
type PurchaseClaimAction string
type PurchaseClaimStatus string

const (
	PurchaseClaimMissing PurchaseClaimType = "missing"
	PurchaseClaimDamaged PurchaseClaimType = "damaged"

	PurchaseClaimActionReplacement PurchaseClaimAction = "replacement"
	PurchaseClaimActionCreditNote  PurchaseClaimAction = "credit_note"
	PurchaseClaimActionReturn      PurchaseClaimAction = "return_to_vendor"
	PurchaseClaimActionScrap       PurchaseClaimAction = "scrap"
	PurchaseClaimActionAdjustment  PurchaseClaimAction = "adjustment_only"

	PurchaseClaimStatusOpen      PurchaseClaimStatus = "open"
	PurchaseClaimStatusPartial   PurchaseClaimStatus = "partial"
	PurchaseClaimStatusResolved  PurchaseClaimStatus = "resolved"
	PurchaseClaimStatusCancelled PurchaseClaimStatus = "cancelled"
)

type PurchaseClaim struct {
	ID                  string              `gorm:"type:varchar(255);primaryKey" json:"id"`
	ClaimNumber         string              `gorm:"type:varchar(100);uniqueIndex;not null" json:"claim_number"`
	PurchaseOrderID     string              `gorm:"type:varchar(255);not null;index" json:"purchase_order_id"`
	PurchaseOrderNumber string              `gorm:"type:varchar(100);index" json:"purchase_order_number"`
	VendorID            uint                `gorm:"not null;index" json:"vendor_id"`
	CompanyID           uint                `gorm:"not null;index" json:"company_id"`
	ClaimDate           time.Time           `gorm:"not null;index" json:"claim_date"`
	Status              PurchaseClaimStatus `gorm:"type:varchar(30);not null;default:'open';index" json:"status"`
	Notes               string              `gorm:"type:text" json:"notes"`
	CreatedBy           string              `gorm:"type:varchar(255);not null" json:"created_by"`
	Items               []PurchaseClaimItem `gorm:"foreignKey:PurchaseClaimID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

func (PurchaseClaim) TableName() string {
	return "purchase_claims"
}

type PurchaseClaimItem struct {
	ID                  uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseClaimID     string            `gorm:"type:varchar(255);not null;index" json:"purchase_claim_id"`
	PurchaseOrderItemID uint              `gorm:"not null;index" json:"purchase_order_item_id"`
	ProductID           string            `gorm:"type:varchar(255);not null;index" json:"product_id"`
	ProductName         string            `gorm:"type:varchar(255)" json:"product_name"`
	SKU                 string            `gorm:"type:varchar(100);index" json:"sku"`
	IsRawMaterial       bool              `gorm:"not null;default:false;index" json:"is_raw_material"`
	Type                PurchaseClaimType `gorm:"type:varchar(30);not null;index" json:"type"`

	Quantity float64 `gorm:"type:decimal(18,6);not null" json:"quantity"`
	Unit     string  `gorm:"type:varchar(50);not null" json:"unit"`

	BaseQuantity float64 `gorm:"type:decimal(18,6);not null" json:"base_quantity"`
	BaseUnit     string  `gorm:"type:varchar(50);not null" json:"base_unit"`

	Rate   float64 `gorm:"type:decimal(18,6);default:0" json:"rate"`
	Amount float64 `gorm:"type:decimal(18,6);default:0" json:"amount"`

	Reason string              `gorm:"type:text;not null" json:"reason"`
	Action PurchaseClaimAction `gorm:"type:varchar(50);not null" json:"action"`

	StockAdjusted bool `gorm:"not null;default:false" json:"stock_adjusted"`

	// Replacement tracking.
	ReplacementPendingBase  float64    `gorm:"type:decimal(18,6);default:0" json:"replacement_pending_base"`
	ReplacementReceivedBase float64    `gorm:"type:decimal(18,6);default:0" json:"replacement_received_base"`
	ReplacementCompleted    bool       `gorm:"not null;default:false" json:"replacement_completed"`
	ReplacementCompletedAt  *time.Time `json:"replacement_completed_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PurchaseClaimItem) TableName() string {
	return "purchase_claim_items"
}

type PurchaseClaimReceipt struct {
	ID                   string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	PurchaseClaimID      string    `gorm:"type:varchar(255);not null;index" json:"purchase_claim_id"`
	PurchaseClaimItemID  uint      `gorm:"not null;index" json:"purchase_claim_item_id"`
	ProductID            string    `gorm:"type:varchar(255);not null;index" json:"product_id"`
	ReceivedQuantity     float64   `gorm:"type:decimal(18,6);not null" json:"received_quantity"`
	ReceivedUnit         string    `gorm:"type:varchar(50);not null" json:"received_unit"`
	ReceivedBaseQuantity float64   `gorm:"type:decimal(18,6);not null" json:"received_base_quantity"`
	BaseUnit             string    `gorm:"type:varchar(50);not null" json:"base_unit"`
	ReceivedDate         time.Time `gorm:"not null;index" json:"received_date"`
	Notes                string    `gorm:"type:text" json:"notes"`
	ReceivedBy           string    `gorm:"type:varchar(255);not null" json:"received_by"`
	CreatedAt            time.Time `json:"created_at"`
}

func (PurchaseClaimReceipt) TableName() string {
	return "purchase_claim_receipts"
}
