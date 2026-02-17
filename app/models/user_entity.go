package models

import (
	"time"

	"gorm.io/gorm"
)

type EntityType string

const (
	EntityTypeVendor   EntityType = "vendor"
	EntityTypeCustomer EntityType = "customer"
)


type EntityOtherDetails struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	EntityID         uint           `gorm:"not null;index:idx_entity" json:"entity_id"`
	EntityType       string         `gorm:"type:varchar(20);not null;index:idx_entity" json:"entity_type"` 
	PAN              string         `gorm:"type:varchar(10)" json:"pan"`
	IsMSMERegistered bool           `gorm:"default:false" json:"is_msme_registered"`
	Currency         string         `gorm:"type:varchar(50);default:'INR- Indian Rupee'" json:"currency"`
	PaymentTerms     string         `gorm:"type:varchar(100);default:'Due on Receipt'" json:"payment_terms"`
	TDS              string         `gorm:"type:varchar(100)" json:"tds"`
	EnablePortal     bool           `gorm:"default:false" json:"enable_portal"`
	WebsiteURL       string         `gorm:"type:varchar(255)" json:"website_url"`
	Department       string         `gorm:"type:varchar(100)" json:"department"`
	Designation      string         `gorm:"type:varchar(100)" json:"designation"`
	Twitter          string         `gorm:"type:varchar(255)" json:"twitter"`
	SkypeName        string         `gorm:"type:varchar(100)" json:"skype_name"`
	Facebook         string         `gorm:"type:varchar(255)" json:"facebook"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityOtherDetails) TableName() string {
	return "entity_other_details"
}

type EntityAddress struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	EntityID      uint           `gorm:"not null;index:idx_entity_address" json:"entity_id"`
	EntityType    string         `gorm:"type:varchar(20);not null;index:idx_entity_address" json:"entity_type"`
	AddressType   string         `gorm:"type:varchar(20);not null" json:"address_type"`
	Attention     string         `gorm:"type:varchar(255)" json:"attention"`
	CountryRegion string         `gorm:"type:varchar(100)" json:"country_region"`
	AddressLine1  string         `gorm:"type:text" json:"address_line1"`
	AddressLine2  string         `gorm:"type:text" json:"address_line2"`
	City          string         `gorm:"type:varchar(100)" json:"city"`
	State         string         `gorm:"type:varchar(100)" json:"state"`
	PinCode       string         `gorm:"type:varchar(10)" json:"pin_code"`
	Phone         string         `gorm:"type:varchar(20)" json:"phone"`
	PhoneCode     string         `gorm:"type:varchar(5);default:'+91'" json:"phone_code"`
	FaxNumber     string         `gorm:"type:varchar(20)" json:"fax_number"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityAddress) TableName() string {
	return "entity_addresses"
}

type EntityContactPerson struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	EntityID      uint           `gorm:"not null;index:idx_entity_contact" json:"entity_id"`
	EntityType    string         `gorm:"type:varchar(20);not null;index:idx_entity_contact" json:"entity_type"`
	Salutation    string         `gorm:"type:varchar(10)" json:"salutation"`
	FirstName     string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName      string         `gorm:"type:varchar(100)" json:"last_name"`
	EmailAddress  string         `gorm:"type:varchar(255)" json:"email_address"`
	WorkPhone     string         `gorm:"type:varchar(20)" json:"work_phone"`
	WorkPhoneCode string         `gorm:"type:varchar(5);default:'+91'" json:"work_phone_code"`
	Mobile        string         `gorm:"type:varchar(20)" json:"mobile"`
	MobileCode    string         `gorm:"type:varchar(5);default:'+91'" json:"mobile_code"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityContactPerson) TableName() string {
	return "entity_contact_persons"
}

type EntityBankDetail struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	EntityID          uint           `gorm:"not null;index:idx_entity_bank" json:"entity_id"`
	EntityType        string         `gorm:"type:varchar(20);not null;index:idx_entity_bank" json:"entity_type"` 
	AccountHolderName string         `gorm:"type:varchar(255)" json:"account_holder_name"`
	BankName          string         `gorm:"type:varchar(255)" json:"bank_name"`
	AccountNumber     string         `gorm:"type:varchar(50);not null" json:"account_number"`
	IFSC              string         `gorm:"type:varchar(11);not null" json:"ifsc"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityBankDetail) TableName() string {
	return "entity_bank_details"
}

type EntityDocument struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	EntityID   uint           `gorm:"not null;index:idx_entity_document" json:"entity_id"`
	EntityType string         `gorm:"type:varchar(20);not null;index:idx_entity_document" json:"entity_type"` 
	FileName   string         `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath   string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize   int64          `gorm:"not null" json:"file_size"`
	MimeType   string         `gorm:"type:varchar(100)" json:"mime_type"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EntityDocument) TableName() string {
	return "entity_documents"
}
