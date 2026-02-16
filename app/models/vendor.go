package models

import (
	"time"

	"gorm.io/gorm"
)

type Vendor struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Salutation     string         `gorm:"type:varchar(10)" json:"salutation"`
	FirstName      string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName       string         `gorm:"type:varchar(100)" json:"last_name"`
	CompanyName    string         `gorm:"type:varchar(255)" json:"company_name"`
	DisplayName    string         `gorm:"type:varchar(255);not null" json:"display_name"`
	EmailAddress   string         `gorm:"type:varchar(255)" json:"email_address"`
	WorkPhone      string         `gorm:"type:varchar(20)" json:"work_phone"`
	WorkPhoneCode  string         `gorm:"type:varchar(5);default:'+91'" json:"work_phone_code"`
	Mobile         string         `gorm:"type:varchar(20)" json:"mobile"`
	MobileCode     string         `gorm:"type:varchar(5);default:'+91'" json:"mobile_code"`
	VendorLanguage string         `gorm:"type:varchar(50);default:'English'" json:"vendor_language"`
	GSTIN          string         `gorm:"type:varchar(15)" json:"gstin"` // For GST portal prefill
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships using polymorphic associations
	OtherDetails    *EntityOtherDetails   `gorm:"polymorphic:Entity;polymorphicValue:vendor" json:"other_details,omitempty"`
	BillingAddress  *EntityAddress        `gorm:"polymorphic:Entity;polymorphicValue:vendor" json:"billing_address,omitempty"`
	ShippingAddress *EntityAddress        `gorm:"polymorphic:Entity;polymorphicValue:vendor" json:"shipping_address,omitempty"`
	ContactPersons  []EntityContactPerson `gorm:"polymorphic:Entity;polymorphicValue:vendor" json:"contact_persons,omitempty"`
	BankDetails     []VendorBankDetail    `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bank_details,omitempty"`
	Documents       []EntityDocument      `gorm:"polymorphic:Entity;polymorphicValue:vendor" json:"documents,omitempty"`
}

func (Vendor) TableName() string {
	return "vendors"
}