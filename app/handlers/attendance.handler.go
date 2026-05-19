package handlers

import (
	"strconv"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	attendanceService services.EmployeeAttendanceService
}

func NewAttendanceHandler(attendanceService services.EmployeeAttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendanceService,
	}
}

// CreateAttendance creates a new attendance record
func (h *AttendanceHandler) CreateAttendance(c *fiber.Ctx) error {
	var req input.CreateAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.CreateAttendance(c.Context(), companyID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// BulkCreateAttendance creates multiple attendance records at once
func (h *AttendanceHandler) BulkCreateAttendance(c *fiber.Ctx) error {
	var req input.BulkAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.BulkCreateAttendance(c.Context(), companyID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetAttendanceByID retrieves a specific attendance record
func (h *AttendanceHandler) GetAttendanceByID(c *fiber.Ctx) error {
	attendanceID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid attendance ID",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceByID(c.Context(), uint(attendanceID), companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetAttendanceByEmployeeID retrieves all attendance records for an employee
func (h *AttendanceHandler) GetAttendanceByEmployeeID(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("employee_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceByEmployeeID(c.Context(), uint(employeeID), companyID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetAttendanceByCompanyID retrieves all attendance records for a company
func (h *AttendanceHandler) GetAttendanceByCompanyID(c *fiber.Ctx) error {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceByCompanyID(c.Context(), companyID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetAttendanceByDateRange retrieves attendance records within a date range
func (h *AttendanceHandler) GetAttendanceByDateRange(c *fiber.Ctx) error {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "from_date and to_date query parameters are required",
		})
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid from_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid to_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceByDateRange(c.Context(), companyID, fromDate, toDate, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetAttendanceByEmployeeAndDateRange retrieves attendance records for an employee within a date range
func (h *AttendanceHandler) GetAttendanceByEmployeeAndDateRange(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("employee_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "from_date and to_date query parameters are required",
		})
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid from_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid to_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceByEmployeeAndDateRange(c.Context(), uint(employeeID), companyID, fromDate, toDate, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// UpdateAttendance updates an attendance record
func (h *AttendanceHandler) UpdateAttendance(c *fiber.Ctx) error {
	attendanceID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid attendance ID",
		})
	}

	var req input.UpdateAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.UpdateAttendance(c.Context(), uint(attendanceID), companyID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// DeleteAttendance deletes an attendance record
func (h *AttendanceHandler) DeleteAttendance(c *fiber.Ctx) error {
	attendanceID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid attendance ID",
		})
	}

	companyID := c.Locals("company_id").(uint)

	err = h.attendanceService.DeleteAttendance(c.Context(), uint(attendanceID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "Attendance record deleted successfully",
	})
}

// GetAttendanceStats retrieves attendance statistics
func (h *AttendanceHandler) GetAttendanceStats(c *fiber.Ctx) error {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "from_date and to_date query parameters are required",
		})
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid from_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid to_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local)

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceStats(c.Context(), companyID, fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetAttendanceCalendarView retrieves attendance in calendar view format
func (h *AttendanceHandler) GetAttendanceCalendarView(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("employee_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "from_date and to_date query parameters are required",
		})
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid from_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid to_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local)

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetAttendanceCalendarView(c.Context(), uint(employeeID), companyID, fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// GetCompanyAttendanceWeekView retrieves all employees with attendance for a week/date range
func (h *AttendanceHandler) GetCompanyAttendanceWeekView(c *fiber.Ctx) error {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "from_date and to_date query parameters are required",
		})
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid from_date format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.Local)

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid toDate format (use YYYY-MM-DD)",
		})
	}
	// Normalize to local timezone
	toDate = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.Local)

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.GetCompanyAttendanceWeekView(c.Context(), companyID, fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// CheckInEmployee marks employee as checked in
func (h *AttendanceHandler) CheckInEmployee(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("employee_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.CheckInEmployee(c.Context(), uint(employeeID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// CheckOutEmployee marks employee as checked out
func (h *AttendanceHandler) CheckOutEmployee(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("employee_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.attendanceService.CheckOutEmployee(c.Context(), uint(employeeID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}
