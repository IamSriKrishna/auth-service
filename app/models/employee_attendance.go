package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceStatus string

const (
	AttendanceOnTime  AttendanceStatus = "on_time"
	AttendanceAbsent  AttendanceStatus = "absent"
	AttendanceLate    AttendanceStatus = "late"
	AttendanceHoliday AttendanceStatus = "holiday"
	AttendanceHalfDay AttendanceStatus = "half_day"
	AttendanceLeave   AttendanceStatus = "leave"
)

type EmployeeAttendance struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	EmployeeID   uint             `gorm:"not null;index" json:"employee_id"`
	CompanyID    uint             `gorm:"not null;index" json:"company_id"`
	Date         time.Time        `gorm:"type:date;not null;index" json:"date"`
	Status       AttendanceStatus `gorm:"type:varchar(50);not null" json:"status"`
	Reason       string           `gorm:"type:text" json:"reason"`
	CheckInTime  *time.Time       `json:"check_in_time"`
	CheckOutTime *time.Time       `json:"check_out_time"`
	WorkingHours float64          `json:"working_hours"`
	Notes        string           `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relations
	Employee *Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Company  *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (EmployeeAttendance) TableName() string {
	return "employee_attendance"
}
