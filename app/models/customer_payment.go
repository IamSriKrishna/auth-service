package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type CustomerPayment struct {
	ID                   uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentNumber        string               `json:"payment_number" gorm:"type:varchar(100);uniqueIndex;not null"`
	SalesOrderID         string               `json:"sales_order_id" gorm:"type:varchar(255);index;not null"`
	SalesOrder           *SalesOrder          `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CustomerID           uint                 `json:"customer_id" gorm:"not null;index"`
	Customer             *Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PaymentMode          domain.PaymentMode   `json:"payment_mode" gorm:"type:varchar(50);not null"`                     // cash, online
	Amount               float64              `json:"amount" gorm:"not null;default:0"`                                  // Initial amount being paid in this payment record
	ReceivedAmount       float64              `json:"received_amount" gorm:"default:0"`                                  // Amount actually recorded/confirmed for this payment
	RemainingAmount      float64              `json:"remaining_amount" gorm:"not null;default:0"`                        // DERIVED: Remaining on SO after all payments
	PaymentStatus        domain.PaymentStatus `json:"payment_status" gorm:"type:varchar(50);not null;default:'pending'"` // pending, received, cancelled
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

func (CustomerPayment) TableName() string {
	return "customer_payments"
}
