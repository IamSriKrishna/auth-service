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

func employeeLocalUint(c *fiber.Ctx, key string) uint {
	value := c.Locals(key)

	switch typed := value.(type) {
	case uint:
		return typed
	case uint64:
		return uint(typed)
	case int:
		if typed > 0 {
			return uint(typed)
		}
	case int64:
		if typed > 0 {
			return uint(typed)
		}
	case float64:
		if typed > 0 {
			return uint(typed)
		}
	case string:
		parsed, err := strconv.ParseUint(typed, 10, 64)
		if err == nil {
			return uint(parsed)
		}
	}

	return 0
}

func employeeAuthContext(c *fiber.Ctx) (uint, uint, error) {
	userID := employeeLocalUint(c, "user_id")
	companyID := employeeLocalUint(c, "company_id")

	if userID == 0 {
		return 0, 0, fmt.Errorf("invalid authenticated user")
	}

	if companyID == 0 {
		return 0, 0, fmt.Errorf("user is not assigned to a company")
	}

	return userID, companyID, nil
}

func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req input.CreateEmployeeRequest

	if c.FormValue("name") == "" {
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

	if req.SalaryType == "monthly" {
		monthlySalaryString := c.FormValue("monthly_salary")
		if monthlySalaryString == "" {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Monthly salary is required",
			})
		}

		monthlySalary, err := strconv.ParseFloat(monthlySalaryString, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Invalid monthly salary format",
			})
		}

		req.MonthlySalary = monthlySalary
	}

	if req.SalaryType == "weekly" {
		weeklySalaryString := c.FormValue("weekly_salary")
		if weeklySalaryString == "" {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Weekly salary is required",
			})
		}

		weeklySalary, err := strconv.ParseFloat(weeklySalaryString, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
				Error:   true,
				Message: "Invalid weekly salary format",
			})
		}

		req.WeeklySalary = weeklySalary
	}

	file, err := c.FormFile("document")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Document file is required",
		})
	}

	createdByID, companyID, err := employeeAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	response, err := h.employeeService.CreateEmployeeWithFile(
		c.Context(),
		createdByID,
		companyID,
		&req,
		file,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(output.SuccessResponse{
		Success: true,
		Data:    response,
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

	_, companyID, err := employeeAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	response, err := h.employeeService.GetEmployeeByID(
		c.Context(),
		uint(employeeID),
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    response,
	})
}

func (h *EmployeeHandler) GetEmployees(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	_, companyID, err := employeeAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	response, err := h.employeeService.GetEmployeesByCompany(
		c.Context(),
		companyID,
		page,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(response)
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

	_, companyID, err := employeeAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	response, err := h.employeeService.UpdateEmployee(
		c.Context(),
		uint(employeeID),
		companyID,
		&req,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Data:    response,
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

	_, companyID, err := employeeAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	if err := h.employeeService.DeleteEmployee(
		c.Context(),
		uint(employeeID),
		companyID,
	); err != nil {
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