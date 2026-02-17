package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type SalesOrder struct {
	ID                   string                  `json:"id" gorm:"type:varchar(255);primaryKey"`
	SalesOrderNumber     string                  `json:"sales_order_no" gorm:"column:sales_order_no;type:varchar(100);uniqueIndex;not null"`
	CustomerID           uint                    `json:"customer_id" gorm:"not null;index"`
	Customer             *Customer               `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SalespersonID        *uint                   `json:"salesperson_id,omitempty" gorm:"index"`
	Salesperson          *Salesperson            `json:"salesperson,omitempty" gorm:"foreignKey:SalespersonID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ReferenceNo          string                  `json:"reference_no" gorm:"type:varchar(100)"`
	SODate               time.Time               `json:"sales_order_date" gorm:"not null"`
	ExpectedShipmentDate time.Time               `json:"expected_shipment_date" gorm:"not null"`
	PaymentTerms         domain.PaymentTerms     `json:"payment_terms" gorm:"type:varchar(50);not null"`
	DeliveryMethod       string                  `json:"delivery_method" gorm:"type:varchar(255)"`
	LineItems            []SalesOrderLineItem    `json:"line_items" gorm:"foreignKey:SalesOrderID;constraint:OnDelete:CASCADE"`
	SubTotal             float64                 `json:"sub_total" gorm:"not null;default:0"`
	ShippingCharges      float64                 `json:"shipping_charges" gorm:"default:0"`
	TaxType              *domain.TaxType         `json:"tax_type" gorm:"type:varchar(10)"`
	TaxID                *uint                   `json:"tax_id,omitempty" gorm:"index"`
	Tax                  *Tax                    `json:"tax,omitempty" gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TaxAmount            float64                 `json:"tax_amount" gorm:"default:0"`
	Adjustment           float64                 `json:"adjustment" gorm:"default:0"`
	Total                float64                 `json:"total" gorm:"not null;default:0"`
	CustomerNotes        string                  `json:"customer_notes" gorm:"type:text"`
	TermsAndConditions   string                  `json:"terms_and_conditions" gorm:"type:text"`
	Status               domain.SalesOrderStatus `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`

	InventoryReserved bool       `json:"inventory_reserved" gorm:"default:false;index"`
	InventoryDeducted bool       `json:"inventory_deducted" gorm:"default:false;index"`
	ReservedDate      *time.Time `json:"reserved_date"`
	DeductedDate      *time.Time `json:"deducted_date"`

	Attachments []string  `json:"attachments,omitempty" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy   string    `json:"updated_by" gorm:"type:varchar(255)"`
}

type SalesOrderLineItem struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	SalesOrderID     string         `json:"sales_order_id" gorm:"type:varchar(255);not null;index"`
	ItemID           string         `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item             *Item          `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantID        *uint          `json:"variant_id,omitempty" gorm:"index"`
	Variant          *Variant       `json:"variant,omitempty" gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Quantity         float64        `json:"quantity" gorm:"not null"`
	InvoicedQuantity float64        `json:"invoiced_quantity" gorm:"default:0"`
	Rate             float64        `json:"rate" gorm:"not null"`
	Amount           float64        `json:"amount" gorm:"not null"`
	VariantDetails   VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (SalesOrder) TableName() string {
	return "sales_orders"
}

func (SalesOrderLineItem) TableName() string {
	return "sales_order_line_items"
}
