package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ManufacturerEmployeeAssignment represents an employee assignment with cost details
type ManufacturerEmployeeAssignment struct {
	EmployeeID  uint    `json:"employee_id"`
	ServiceCost float64 `json:"service_cost"` // Cost per unit manufactured or fixed cost
	CostType    string  `json:"cost_type"`    // 'fixed' or 'per_unit'
}

// EmployeeAssignments is a custom type for storing JSON array in database
type EmployeeAssignments []ManufacturerEmployeeAssignment

// Value implements the driver.Valuer interface
func (e EmployeeAssignments) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Scan implements the sql.Scanner interface
func (e *EmployeeAssignments) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &e)
}

// Manufacturer represents a manufacturing batch/job linked to a product group
type Manufacturer struct {
	ID             string              `json:"id" gorm:"type:varchar(255);primaryKey"`
	Name           string              `json:"name" gorm:"type:varchar(255);not null"`
	ProductGroupID string              `json:"product_group_id" gorm:"type:varchar(255);index;not null"`
	ProductGroup   *ProductGroup       `json:"product_group,omitempty" gorm:"foreignKey:ProductGroupID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity       float64             `json:"quantity" gorm:"not null"`                         // Quantity to manufacture
	Status         string              `json:"status" gorm:"type:varchar(50);default:'pending'"` // pending, in_progress, completed, cancelled
	Description    string              `json:"description" gorm:"type:text"`                     // Manufacturing notes
	Employees      EmployeeAssignments `json:"employees" gorm:"type:json"`                       // JSON array of employee assignments
	CompanyID      uint                `json:"company_id" gorm:"not null;index"`                 // Company ownership
	Company        *Company            `json:"company,omitempty" gorm:"foreignKey:CompanyID"`    // Company reference
	UserID         uint                `json:"user_id" gorm:"not null;index"`                    // User who created the record
	User           *User               `json:"user,omitempty" gorm:"foreignKey:UserID"`          // User reference
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

func (Manufacturer) TableName() string {
	return "manufacturers"
}
