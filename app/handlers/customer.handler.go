package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	service services.CustomerService
}

func NewCustomerHandler(service services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func customerLocalUint(c *fiber.Ctx, key string) uint {
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

func customerAuthContext(c *fiber.Ctx) (uint, uint, error) {
	userID := customerLocalUint(c, "user_id")
	companyID := customerLocalUint(c, "company_id")

	if userID == 0 {
		return 0, 0, fmt.Errorf("invalid authenticated user")
	}

	if companyID == 0 {
		return 0, 0, fmt.Errorf("user is not assigned to a company")
	}

	return userID, companyID, nil
}

func normalizeCustomerPagination(c *fiber.Ctx) (int, int) {
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

	return page, limit
}

func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req input.CreateCustomerInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userID, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customer, err := h.service.CreateCustomerForUser(
		userID,
		companyID,
		&req,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid customer id")
	}

	var req input.UpdateCustomerInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customer, err := h.service.UpdateCustomerForCompany(
		uint(id),
		companyID,
		&req,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (h *CustomerHandler) GetCustomerByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid customer id")
	}

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customer, err := h.service.GetCustomerByIDAndCompany(
		uint(id),
		companyID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "customer not found")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (h *CustomerHandler) GetAllCustomers(c *fiber.Ctx) error {
	page, limit := normalizeCustomerPagination(c)

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customers, total, err := h.service.GetCustomersByCompany(
		companyID,
		page,
		limit,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customers,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid customer id")
	}

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	if err := h.service.DeleteCustomerForCompany(
		uint(id),
		companyID,
	); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "customer deleted successfully",
	})
}

// Kept for route compatibility.
// It now returns all customers belonging to the authenticated company.
func (h *CustomerHandler) GetUserCustomers(c *fiber.Ctx) error {
	page, limit := normalizeCustomerPagination(c)

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customers, total, err := h.service.GetCustomersByCompany(
		companyID,
		page,
		limit,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customers,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *CustomerHandler) GetUserCustomerByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid customer id")
	}

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customer, err := h.service.GetCustomerByIDAndCompany(
		uint(id),
		companyID,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			"customer not found or unauthorized",
		)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (h *CustomerHandler) UpdateUserCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid customer id")
	}

	var req input.UpdateCustomerInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	customer, err := h.service.UpdateCustomerForCompany(
		uint(id),
		companyID,
		&req,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    customer,
	})
}

func (h *CustomerHandler) DeleteUserCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid customer id")
	}

	_, companyID, err := customerAuthContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	if err := h.service.DeleteCustomerForCompany(
		uint(id),
		companyID,
	); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "customer deleted successfully",
	})
}

// This keeps the models import used in projects where DeleteCustomer
// still constructs a model directly elsewhere.
var _ = models.Customer{}