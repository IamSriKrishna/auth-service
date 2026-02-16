package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CustomerType     string         `gorm:"type:varchar(20);not null;default:'Business'" json:"customer_type"` 
	Salutation       string         `gorm:"type:varchar(10)" json:"salutation"`
	FirstName        string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName         string         `gorm:"type:varchar(100)" json:"last_name"`
	CompanyName      string         `gorm:"type:varchar(255)" json:"company_name"`
	DisplayName      string         `gorm:"type:varchar(255);not null" json:"display_name"`
	EmailAddress     string         `gorm:"type:varchar(255)" json:"email_address"`
	WorkPhone        string         `gorm:"type:varchar(20)" json:"work_phone"`
	WorkPhoneCode    string         `gorm:"type:varchar(5);default:'+91'" json:"work_phone_code"`
	Mobile           string         `gorm:"type:varchar(20)" json:"mobile"`
	MobileCode       string         `gorm:"type:varchar(5);default:'+91'" json:"mobile_code"`
	CustomerLanguage string         `gorm:"type:varchar(50);default:'English'" json:"customer_language"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	OtherDetails   *EntityOtherDetails   `gorm:"polymorphic:Entity;polymorphicValue:customer" json:"other_details,omitempty"`
	Addresses      []EntityAddress       `gorm:"polymorphic:Entity;polymorphicValue:customer" json:"addresses,omitempty"`
	ContactPersons []EntityContactPerson `gorm:"polymorphic:Entity;polymorphicValue:customer" json:"contact_persons,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

func (c *Customer) GetBillingAddress() *EntityAddress {
	for i := range c.Addresses {
		if c.Addresses[i].AddressType == "billing" {
			return &c.Addresses[i]
		}
	}
	return nil
}

func (c *Customer) GetShippingAddress() *EntityAddress {
	for i := range c.Addresses {
		if c.Addresses[i].AddressType == "shipping" {
			return &c.Addresses[i]
		}
	}
	return nil
}

func (c *Customer) SetBillingAddress(address *EntityAddress) {
	if address == nil {
		return
	}
	
	address.AddressType = "billing"
	
	for i := range c.Addresses {
		if c.Addresses[i].AddressType == "billing" {
			c.Addresses[i] = *address
			return
		}
	}
	
	c.Addresses = append(c.Addresses, *address)
}

func (c *Customer) SetShippingAddress(address *EntityAddress) {
	if address == nil {
		return
	}
	
	address.AddressType = "shipping"
	
	for i := range c.Addresses {
		if c.Addresses[i].AddressType == "shipping" {
			c.Addresses[i] = *address
			return
		}
	}
	
	c.Addresses = append(c.Addresses, *address)
}