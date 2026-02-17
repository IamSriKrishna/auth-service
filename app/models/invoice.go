package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type Invoice struct {
	ID            string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	InvoiceNumber string    `json:"invoice_number" gorm:"type:varchar(100);uniqueIndex;not null"`
	CustomerID    uint      `json:"customer_id" gorm:"not null;index"`
	Customer      *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	OrderNumber  string              `json:"order_number" gorm:"type:varchar(100)"`
	SalesOrderID *string             `json:"sales_order_id,omitempty" gorm:"type:varchar(255);index"`
	InvoiceDate  time.Time           `json:"invoice_date" gorm:"not null"`
	Terms        domain.PaymentTerms `json:"terms" gorm:"type:varchar(50);not null"`
	DueDate      time.Time           `json:"due_date" gorm:"not null"`

	SalespersonID *uint        `json:"salesperson_id,omitempty" gorm:"index"`
	Salesperson   *Salesperson `json:"salesperson,omitempty" gorm:"foreignKey:SalespersonID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Subject string `json:"subject" gorm:"type:text"`

	LineItems []InvoiceLineItem `json:"line_items" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`

	SubTotal        float64        `json:"sub_total" gorm:"not null;default:0"`
	ShippingCharges float64        `json:"shipping_charges" gorm:"default:0"`
	TaxType         domain.TaxType `json:"tax_type" gorm:"type:varchar(10)"`
	TaxID           *uint          `json:"tax_id,omitempty" gorm:"index"`
	Tax             *Tax           `json:"tax,omitempty" gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TaxAmount       float64        `json:"tax_amount" gorm:"default:0"`
	Adjustment      float64        `json:"adjustment" gorm:"default:0"`
	Total           float64        `json:"total" gorm:"not null;default:0"`

	CustomerNotes      string `json:"customer_notes" gorm:"type:text"`
	TermsAndConditions string `json:"terms_and_conditions" gorm:"type:text"`

	PaymentReceived bool           `json:"payment_received" gorm:"default:false"`
	Payments        []Payment      `json:"payments,omitempty" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`
	PaymentSplits   []PaymentSplit `json:"payment_splits,omitempty" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`

	EmailCommunications []EmailCommunication `json:"email_communications,omitempty" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`

	Status domain.InvoiceStatus `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`

	InventorySynced   bool       `json:"inventory_synced" gorm:"default:false;index"`
	InventorySyncDate *time.Time `json:"inventory_sync_date"`

	Attachments InvoiceAttachments `json:"attachments,omitempty" gorm:"type:json"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy string    `json:"updated_by" gorm:"type:varchar(255)"`
}

func (Invoice) TableName() string {
	return "invoices"
}

type InvoiceLineItem struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	InvoiceID string `gorm:"type:varchar(255);index;not null"`

	ItemID string `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item   *Item  `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	VariantID *uint    `json:"variant_id,omitempty" gorm:"index"`
	Variant   *Variant `json:"variant,omitempty" gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Description string  `json:"description" gorm:"type:text"`
	Quantity    float64 `json:"quantity" gorm:"type:decimal(18,2);not null"`
	Rate        float64 `json:"rate" gorm:"not null"`
	Amount      float64 `json:"amount" gorm:"not null"`

	VariantDetails VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`

	InventorySynced bool       `json:"inventory_synced" gorm:"default:false"`
	SyncedAt        *time.Time `json:"synced_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (InvoiceLineItem) TableName() string {
	return "invoice_line_items"
}

type VariantDetails map[string]string

func (v VariantDetails) Value() (driver.Value, error) {
	if len(v) == 0 {
		return nil, nil
	}
	return json.Marshal(v)
}

func (v *VariantDetails) Scan(value interface{}) error {
	if value == nil {
		*v = make(map[string]string)
		return nil
	}

	var bytes []byte
	switch val := value.(type) {
	case []byte:
		bytes = val
	case string:
		bytes = []byte(val)
	default:
		return fmt.Errorf("invalid type for VariantDetails.Scan: %T", value)
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	result := make(map[string]string)
	for k, v := range temp {
		switch val := v.(type) {
		case string:
			result[k] = val
		case float64:
			if val == float64(int(val)) {
				result[k] = fmt.Sprintf("%d", int(val))
			} else {
				result[k] = fmt.Sprintf("%v", val)
			}
		case bool:
			result[k] = fmt.Sprintf("%v", val)
		case uint:
			result[k] = fmt.Sprintf("%d", val)
		case int:
			result[k] = fmt.Sprintf("%d", val)
		case nil:
			result[k] = ""
		default:
			result[k] = fmt.Sprintf("%v", val)
		}
	}
	*v = result
	return nil
}

type InvoiceAttachments []string

func (a InvoiceAttachments) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *InvoiceAttachments) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal InvoiceAttachments value")
	}
	return json.Unmarshal(bytes, a)
}

type Salesperson struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(255);not null"`
	Email string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Salesperson) TableName() string {
	return "salespersons"
}

type Tax struct {
	ID      uint    `gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name" gorm:"type:varchar(255);not null"`
	TaxType string  `json:"tax_type" gorm:"type:varchar(50);not null"`
	Rate    float64 `json:"rate" gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tax) TableName() string {
	return "taxes"
}

type Payment struct {
	ID        uint     `gorm:"primaryKey;autoIncrement"`
	InvoiceID string   `gorm:"type:varchar(255);index;not null"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	PaymentDate time.Time `json:"payment_date" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	PaymentMode string    `json:"payment_mode" gorm:"type:varchar(50)"`
	Reference   string    `json:"reference" gorm:"type:varchar(255)"`
	Notes       string    `json:"notes" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:varchar(255)"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentSplit struct {
	ID        uint     `gorm:"primaryKey;autoIncrement"`
	InvoiceID string   `gorm:"type:varchar(255);index;not null"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	PaymentMode    string  `json:"payment_mode" gorm:"type:varchar(50);not null"`
	DepositTo      string  `json:"deposit_to" gorm:"type:varchar(255)"`
	AmountReceived float64 `json:"amount_received" gorm:"default:0"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by" gorm:"type:varchar(255)"`
}

func (PaymentSplit) TableName() string {
	return "payment_splits"
}

type EmailCommunication struct {
	ID        uint     `gorm:"primaryKey;autoIncrement"`
	InvoiceID string   `gorm:"type:varchar(255);index;not null"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	EmailAddress string `json:"email_address" gorm:"type:varchar(255);not null"`
	Subject      string `json:"subject" gorm:"type:varchar(255)"`
	Status       string `json:"status" gorm:"type:varchar(50);default:'pending'"`

	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by" gorm:"type:varchar(255)"`
	SentAt    *time.Time `json:"sent_at"`
}

func (EmailCommunication) TableName() string {
	return "email_communications"
}
