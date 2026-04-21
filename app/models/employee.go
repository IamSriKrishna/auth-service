package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Email         string         `gorm:"type:varchar(255)" json:"email"`
	Number        string         `gorm:"type:varchar(20)" json:"number"`
	Address       string         `gorm:"type:text" json:"address"`
	EmployeeType  string         `gorm:"type:varchar(50)" json:"employee_type"`
	MonthlySalary float64        `gorm:"type:decimal(10,2)" json:"monthly_salary"`
	DocumentURL   string         `gorm:"type:text" json:"document_url"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	CompanyID     uint           `gorm:"not null;index" json:"company_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (Employee) TableName() string {
	return "employees"
}
