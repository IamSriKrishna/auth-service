package models

import "time"

type Country struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CountryName string    `gorm:"size:100;not null;uniqueIndex" json:"country_name"`
	CountryCode string    `gorm:"size:3;not null;uniqueIndex" json:"country_code"`
	PhoneCode   string    `gorm:"size:10" json:"phone_code,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	States []State `gorm:"foreignKey:CountryID" json:"states,omitempty"`
}

func (Country) TableName() string {
	return "countries"
}

type State struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CountryID uint      `gorm:"not null;index" json:"country_id"`
	StateName string    `gorm:"size:100;not null" json:"state_name"`
	StateCode string    `gorm:"size:10" json:"state_code,omitempty"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Country Country `gorm:"foreignKey:CountryID" json:"country,omitempty"`
}

func (State) TableName() string {
	return "states"
}

type CompanyContact struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CompanyID       uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	Mobile          string    `gorm:"size:15;not null;index" json:"mobile"`
	AlternateMobile string    `gorm:"size:15" json:"alternate_mobile,omitempty"`
	Email           string    `gorm:"size:255;not null;index" json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (CompanyContact) TableName() string {
	return "company_contacts"
}

type CompanyAddress struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CompanyID    uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	AddressLine1 string    `gorm:"size:255;not null" json:"address_line1"`
	AddressLine2 string    `gorm:"size:255" json:"address_line2,omitempty"`
	City         string    `gorm:"size:100;not null" json:"city"`
	StateID      uint      `gorm:"not null;index" json:"state_id"`
	CountryID    uint      `gorm:"not null;index" json:"country_id"`
	Pincode      string    `gorm:"size:10;not null;index" json:"pincode"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
	State   State   `gorm:"foreignKey:StateID" json:"state,omitempty"`
	Country Country `gorm:"foreignKey:CountryID" json:"country,omitempty"`
}

func (CompanyAddress) TableName() string {
	return "company_addresses"
}
