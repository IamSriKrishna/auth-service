package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type SalaryHandler struct {
	salaryService services.SalaryService
}

func NewSalaryHandler(salaryService services.SalaryService) *SalaryHandler {
	return &SalaryHandler{
		salaryService: salaryService,
	}
}

func (h *SalaryHandler) CalculateSalary(c *fiber.Ctx) error {
	var req input.CalculateSalaryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.salaryService.CalculateSalary(c.Context(), 0, companyID, &req)
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

func (h *SalaryHandler) GetSalaryCalculation(c *fiber.Ctx) error {
	salaryCalculationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid salary calculation ID",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.salaryService.GetSalaryCalculation(c.Context(), uint(salaryCalculationID), companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

func (h *SalaryHandler) GetSalaryCalculationsByEmployee(c *fiber.Ctx) error {
	employeeID, err := strconv.ParseUint(c.Params("employee_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid employee ID",
		})
	}

	companyID := c.Locals("company_id").(uint)

	resp, err := h.salaryService.GetSalaryCalculationsByEmployee(c.Context(), uint(employeeID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(output.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

func (h *SalaryHandler) ApproveSalary(c *fiber.Ctx) error {
	var req input.ApproveSalaryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	companyID := c.Locals("company_id").(uint)

	err := h.salaryService.ApproveSalary(c.Context(), req.SalaryCalculationID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(output.SuccessResponse{
		Success: true,
		Message: "Salary approved successfully",
	})
}
