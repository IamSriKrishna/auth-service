package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type VendorPayment struct {
	ID                   uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentNumber        string               `json:"payment_number" gorm:"type:varchar(100);uniqueIndex;not null"`
	PurchaseOrderID      string               `json:"purchase_order_id" gorm:"type:varchar(255);index;not null"`
	PurchaseOrder        *PurchaseOrder       `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VendorID             uint                 `json:"vendor_id" gorm:"not null;index"`
	Vendor               *Vendor              `json:"vendor,omitempty" gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PaymentMode          domain.PaymentMode   `json:"payment_mode" gorm:"type:varchar(50);not null"`                     // cash, online
	Amount               float64              `json:"amount" gorm:"not null;default:0"`                                  // Initial amount being paid in this payment record
	PaidAmount           float64              `json:"paid_amount" gorm:"default:0"`                                      // Amount actually recorded/confirmed for this payment
	RemainingAmount      float64              `json:"remaining_amount" gorm:"not null;default:0"`                        // DERIVED: Remaining on PO after all payments (see PurchaseOrder.RemainingAmount)
	PaymentStatus        domain.PaymentStatus `json:"payment_status" gorm:"type:varchar(50);not null;default:'pending'"` // pending, recorded, cancelled (status of THIS payment record)
	PaymentDate          time.Time            `json:"payment_date" gorm:"not null"`
	ReferenceNumber      string               `json:"reference_number" gorm:"type:varchar(100)"` // Check number, online transaction ID, etc.
	Notes                string               `json:"notes" gorm:"type:text"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	CreatedByUserID      string               `json:"created_by_user_id" gorm:"type:varchar(255)"`
	CreatedByUserName    string               `json:"created_by_user_name" gorm:"type:varchar(255)"`
	CreatedByCompanyID   uint                 `json:"created_by_company_id"`
	CreatedByCompanyName string               `json:"created_by_company_name" gorm:"type:varchar(255)"`
}

func (VendorPayment) TableName() string {
	return "vendor_payments"
}
