package handlers

import (
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
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	createdByID := c.Locals("user_id").(uint)
	companyID := c.Locals("company_id").(uint)

	resp, err := h.employeeService.CreateEmployee(c.Context(), createdByID, companyID, &req)
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
