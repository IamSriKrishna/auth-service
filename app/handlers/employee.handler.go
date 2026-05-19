package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type EmployeeHandler struct {
	employeeService services.EmployeeService
}

func NewEmployeeHandler(employeeService services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req input.CreateEmployeeRequest

	// Parse form data
	if err := c.FormValue("name"); err == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Name is required",
		})
	}

	req.Name = c.FormValue("name")
	req.Email = c.FormValue("email")
	req.Number = c.FormValue("number")
	req.Address = c.FormValue("address")
	req.EmployeeType = c.FormValue("employee_type")
	req.SalaryType = c.FormValue("salary_type")

	// Validate salary_type
	if req.SalaryType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Salary type is required (monthly or weekly)",
		})
	}

	if req.SalaryType != "monthly" && req.SalaryType != "weekly" {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Salary type must be either 'monthly' or 'weekly'",
		})
	}

	// Parse salary based on salary_type
	if req.SalaryType == "monthly" {
		monthlySalaryStr := c.FormValue("monthly_salary")
		if monthlySalaryStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Monthly salary is required",
			})
		}

		var monthlySalary float64
		if _, err := fmt.Sscanf(monthlySalaryStr, "%f", &monthlySalary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Invalid monthly salary format",
			})
		}
		req.MonthlySalary = monthlySalary
	} else if req.SalaryType == "weekly" {
		weeklySalaryStr := c.FormValue("weekly_salary")
		if weeklySalaryStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Weekly salary is required",
			})
		}

		var weeklySalary float64
		if _, err := fmt.Sscanf(weeklySalaryStr, "%f", &weeklySalary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Invalid weekly salary format",
			})
		}
		req.WeeklySalary = weeklySalary
	}

	// Get file from form
	file, err := c.FormFile("document")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Document file is required",
		})
	}

	createdByID := c.Locals("user_id").(uint)
	companyID := c.Locals("company_id").(uint)

	resp, err := h.employeeService.CreateEmployeeWithFile(c.Context(), createdByID, companyID, &req, file)
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

func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	createdByID := c.Locals("user_id").(uint)

	resp, err := h.employeeService.GetEmployeeByID(c.Context(), uint(employeeID), createdByID)
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

func (h *EmployeeHandler) GetEmployees(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	createdByID := c.Locals("user_id").(uint)
	companyID := c.Locals("company_id").(uint)

	resp, err := h.employeeService.GetEmployeesByUser(c.Context(), createdByID, companyID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	var req input.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	createdByID := c.Locals("user_id").(uint)

	resp, err := h.employeeService.UpdateEmployee(c.Context(), uint(employeeID), createdByID, &req)
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

func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	createdByID := c.Locals("user_id").(uint)

	err = h.employeeService.DeleteEmployee(c.Context(), uint(employeeID), createdByID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "Employee deleted successfully",
	})
}
