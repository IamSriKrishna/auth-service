package models

import (
	"time"

	"gorm.io/gorm"
)

type SalaryCalculation struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	EmployeeID       uint           `gorm:"not null;index" json:"employee_id"`
	CompanyID        uint           `gorm:"not null;index" json:"company_id"`
	Month            int            `gorm:"type:int;not null" json:"month"`
	Year             int            `gorm:"type:int;not null" json:"year"`
	BaseSalary       float64        `gorm:"type:decimal(10,2);not null" json:"base_salary"`
	SalaryType       string         `gorm:"type:varchar(20)" json:"salary_type"`
	TotalWorkingDays int            `gorm:"type:int" json:"total_working_days"`
	PresentDays      int            `gorm:"type:int" json:"present_days"`
	AbsentDays       int            `gorm:"type:int" json:"absent_days"`
	HalfDays         int            `gorm:"type:int" json:"half_days"`
	HolidayDays      int            `gorm:"type:int" json:"holiday_days"`
	LeaveDays        int            `gorm:"type:int" json:"leave_days"`
	DailyRate        float64        `gorm:"type:decimal(10,2)" json:"daily_rate"`
	EarningAmount    float64        `gorm:"type:decimal(10,2)" json:"earning_amount"`
	DeductionAmount  float64        `gorm:"type:decimal(10,2)" json:"deduction_amount"`
	NetSalary        float64        `gorm:"type:decimal(10,2)" json:"net_salary"`
	Status           string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
	Notes            string         `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Employee *Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Company  *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (SalaryCalculation) TableName() string {
	return "salary_calculations"
}
