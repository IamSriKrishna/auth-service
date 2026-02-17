package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type PurchaseOrder struct {
	ID                  string  `json:"id" gorm:"type:varchar(255);primaryKey"`
	PurchaseOrderNumber string  `json:"purchase_order_no" gorm:"column:purchase_order_no;type:varchar(100);uniqueIndex;not null"`
	VendorID            uint    `json:"vendor_id" gorm:"not null;index"`
	Vendor              *Vendor `json:"vendor,omitempty" gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	DeliveryAddressType string         `json:"delivery_address_type" gorm:"type:varchar(50);not null"`
	DeliveryAddressID   *uint          `json:"delivery_address_id,omitempty" gorm:"index"`
	DeliveryAddress     *EntityAddress `json:"delivery_address,omitempty" gorm:"foreignKey:DeliveryAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	OrganizationName    string `json:"organization_name" gorm:"type:varchar(255)"`
	OrganizationAddress string `json:"organization_address" gorm:"type:text"`

	CustomerID *uint     `json:"customer_id,omitempty" gorm:"index"`
	Customer   *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	ReferenceNo        string              `json:"reference_no" gorm:"type:varchar(100)"`
	PODate             time.Time           `json:"date" gorm:"not null"`
	DeliveryDate       time.Time           `json:"delivery_date" gorm:"not null"`
	DeliveryDateActual *time.Time          `json:"delivery_date_actual,omitempty"`
	PaymentTerms       domain.PaymentTerms `json:"payment_terms" gorm:"type:varchar(50);not null"`
	ShipmentPreference string              `json:"shipment_preference" gorm:"type:varchar(255)"`

	InventorySynced   bool       `json:"inventory_synced" gorm:"default:false"`
	InventorySyncDate *time.Time `json:"inventory_sync_date"`

	LineItems []PurchaseOrderLineItem `json:"line_items" gorm:"foreignKey:PurchaseOrderID;constraint:OnDelete:CASCADE"`

	SubTotal     float64         `json:"sub_total" gorm:"not null;default:0"`
	Discount     float64         `json:"discount" gorm:"default:0"`
	DiscountType string          `json:"discount_type" gorm:"type:varchar(50)"`
	TaxType      *domain.TaxType `json:"tax_type" gorm:"type:varchar(10)"`
	TaxID        *uint           `json:"tax_id,omitempty" gorm:"index"`
	Tax          *Tax            `json:"tax,omitempty" gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TaxAmount    float64         `json:"tax_amount" gorm:"default:0"`
	Adjustment   float64         `json:"adjustment" gorm:"default:0"`
	Total        float64         `json:"total" gorm:"not null;default:0"`

	Notes              string `json:"notes" gorm:"type:text"`
	TermsAndConditions string `json:"terms_and_conditions" gorm:"type:text"`

	Status domain.PurchaseOrderStatus `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`

	Attachments []string `json:"attachments,omitempty" gorm:"type:json"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy string    `json:"updated_by" gorm:"type:varchar(255)"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

type PurchaseOrderLineItem struct {
	ID              uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseOrderID string `gorm:"type:varchar(255);index;not null" json:"purchase_order_id"`

	ItemID string `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item   *Item  `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	VariantID        *uint          `json:"variant_id,omitempty" gorm:"index"`
	Variant          *Variant       `json:"variant,omitempty" gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Account          string         `json:"account" gorm:"type:varchar(100)"`
	Quantity         float64        `json:"quantity" gorm:"not null"`
	ReceivedQuantity float64        `json:"received_quantity" gorm:"default:0"`
	Rate             float64        `json:"rate" gorm:"not null"`
	Amount           float64        `json:"amount" gorm:"not null"`
	VariantDetails   VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (PurchaseOrderLineItem) TableName() string {
	return "purchase_order_line_items"
}
