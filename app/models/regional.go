package models

import (
	"time"
)

type CompanyRegionalSetting struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CompanyID      uint      `gorm:"not null;uniqueIndex" json:"company_id"`
	Timezone       string    `gorm:"size:50;default:'Asia/Kolkata'" json:"timezone"`
	DateFormat     string    `gorm:"size:20;default:'DD/MM/YYYY'" json:"date_format"`
	TimeFormat     string    `gorm:"size:20;default:'24h'" json:"time_format"`
	CurrencyCode   string    `gorm:"size:3;default:'INR'" json:"currency_code"`
	CurrencySymbol string    `gorm:"size:10;default:'₹'" json:"currency_symbol"`
	LanguageCode   string    `gorm:"size:5;default:'en'" json:"language_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (CompanyRegionalSetting) TableName() string {
	return "company_regional_settings"
}
