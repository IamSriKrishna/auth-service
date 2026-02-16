package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

// Bill represents a bill to vendor (similar to invoice for customers)
type Bill struct {
	ID                string              `json:"id" gorm:"type:varchar(255);primaryKey"`
	BillNumber        string              `json:"bill_number" gorm:"column:bill_number;type:varchar(100);uniqueIndex;not null"`
	VendorID          uint                `json:"vendor_id" gorm:"not null;index"`
	Vendor            *Vendor             `json:"vendor,omitempty" gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BillingAddress    string              `json:"billing_address" gorm:"type:text"`
	OrderNumber       string              `json:"order_number" gorm:"type:varchar(100)"`
	BillDate          time.Time           `json:"bill_date" gorm:"not null"`
	DueDate           time.Time           `json:"due_date" gorm:"not null"`
	PaymentTerms      domain.PaymentTerms `json:"payment_terms" gorm:"type:varchar(50);not null"`
	Subject           string              `json:"subject" gorm:"type:text"`
	LineItems         []BillLineItem      `json:"line_items" gorm:"foreignKey:BillID;constraint:OnDelete:CASCADE"`
	SubTotal          float64             `json:"sub_total" gorm:"not null;default:0"`
	Discount          float64             `json:"discount" gorm:"default:0"`
	TaxType           *domain.TaxType     `json:"tax_type" gorm:"type:varchar(10)"`
	TaxID             *uint               `json:"tax_id,omitempty" gorm:"index"`
	Tax               *Tax                `json:"tax,omitempty" gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TaxAmount         float64             `json:"tax_amount" gorm:"default:0"`
	Adjustment        float64             `json:"adjustment" gorm:"default:0"`
	Total             float64             `json:"total" gorm:"not null;default:0"`
	Notes             string              `json:"notes" gorm:"type:text"`
	Status            domain.BillStatus   `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`
	InventorySynced   bool                `json:"inventory_synced" gorm:"default:false;index"`
	InventorySyncDate *time.Time          `json:"inventory_sync_date"`
	PurchaseOrderID   *string             `json:"purchase_order_id" gorm:"type:varchar(255);index"`
	Attachments       []string            `json:"attachments,omitempty" gorm:"type:json"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	CreatedBy         string              `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy         string              `json:"updated_by" gorm:"type:varchar(255)"`
}

// BillLineItem represents a line item in a bill
type BillLineItem struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	BillID          string         `json:"bill_id" gorm:"type:varchar(255);not null;index"`
	ItemID          string         `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item            *Item          `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantID       *uint          `json:"variant_id,omitempty" gorm:"index"`
	Variant         *Variant       `json:"variant,omitempty" gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Description     string         `json:"description" gorm:"type:text"`
	Account         string         `json:"account" gorm:"type:varchar(255)"`
	Quantity        float64        `json:"quantity" gorm:"not null"`
	Rate            float64        `json:"rate" gorm:"not null"`
	Amount          float64        `json:"amount" gorm:"not null"`
	InventorySynced bool           `json:"inventory_synced" gorm:"default:false"`
	SyncedAt        *time.Time     `json:"synced_at"`
	VariantDetails  VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Bill) TableName() string {
	return "bills"
}

func (BillLineItem) TableName() string {
	return "bill_line_items"
}
