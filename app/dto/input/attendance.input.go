package input

import "time"

type CreateAttendanceRequest struct {
	EmployeeID   uint       `json:"employee_id" validate:"required"`
	Date         string     `json:"date" validate:"required"`
	Status       string     `json:"status" validate:"required,oneof=on_time absent late holiday half_day leave"`
	Reason       string     `json:"reason"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	WorkingHours float64    `json:"working_hours"`
	Notes        string     `json:"notes"`
}

type BulkAttendanceRequest struct {
	EmployeeID uint                `json:"employee_id" validate:"required"`
	Attendance []BulkAttendanceDay `json:"attendance" validate:"required"`
}

type BulkAttendanceDay struct {
	Date         string     `json:"date" validate:"required"`
	Status       string     `json:"status" validate:"required,oneof=on_time absent late holiday half_day leave"`
	Reason       string     `json:"reason"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	WorkingHours float64    `json:"working_hours"`
	Notes        string     `json:"notes"`
}

type UpdateAttendanceRequest struct {
	Status       *string    `json:"status" validate:"omitempty,oneof=on_time absent late holiday half_day leave"`
	Reason       *string    `json:"reason"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	WorkingHours *float64   `json:"working_hours"`
	Notes        *string    `json:"notes"`
}

type GetAttendanceFilterRequest struct {
	EmployeeID *uint      `json:"employee_id"`
	Status     *string    `json:"status"`
	FromDate   *time.Time `json:"from_date"`
	ToDate     *time.Time `json:"to_date"`
}
