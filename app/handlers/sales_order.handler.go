package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SalesOrderHandler struct {
	service  services.SalesOrderService
	validate *validator.Validate
}

func NewSalesOrderHandler(service services.SalesOrderService) *SalesOrderHandler {
	return &SalesOrderHandler{
		service:  service,
		validate: validator.New(),
	}
}

func salesOrderContext(c *fiber.Ctx) (uint, uint, string, error) {
	userID, err := salesOrderLocalToUint(c.Locals("user_id"))
	if err != nil || userID == 0 {
		return 0, 0, "", fiber.NewError(fiber.StatusUnauthorized, "invalid authenticated user")
	}

	companyID, err := salesOrderLocalToUint(c.Locals("company_id"))
	if err != nil || companyID == 0 {
		return 0, 0, "", fiber.NewError(fiber.StatusForbidden, "invalid authenticated company")
	}

	return userID, companyID, strconv.FormatUint(uint64(userID), 10), nil
}

func salesOrderPagination(c *fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

func (h *SalesOrderHandler) CreateSalesOrder(c *fiber.Ctx) error {
	var soInput input.CreateSalesOrderInput
	if err := c.BodyParser(&soInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "success": false})
	}
	if err := h.validate.Struct(soInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	_, companyID, userID, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	so, err := h.service.CreateSalesOrderForCompany(&soInput, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Sales order created successfully",
		"data":    so,
	})
}

func (h *SalesOrderHandler) GetSalesOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Sales order ID is required", "success": false})
	}

	_, companyID, _, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	so, err := h.service.GetSalesOrderForCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	return c.JSON(fiber.Map{"success": true, "data": so})
}

func (h *SalesOrderHandler) GetAllSalesOrders(c *fiber.Ctx) error {
	limit, offset := salesOrderPagination(c)

	_, companyID, _, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	sos, total, err := h.service.GetAllSalesOrdersForCompany(companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	totalAmount := 0.0
	for _, so := range sos {
		totalAmount += so.Total
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"data":         sos,
		"total":        total,
		"total_amount": totalAmount,
	})
}

func (h *SalesOrderHandler) GetSalesOrdersByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Params("customerId"), 10, 32)
	if err != nil || customerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid customer ID", "success": false})
	}

	limit, offset := salesOrderPagination(c)
	_, companyID, _, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	sos, total, err := h.service.GetSalesOrdersByCustomerForCompany(uint(customerID), companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	totalAmount := 0.0
	for _, so := range sos {
		totalAmount += so.Total
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"data":         sos,
		"total":        total,
		"total_amount": totalAmount,
	})
}

func (h *SalesOrderHandler) GetSalesOrdersByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status is required", "success": false})
	}

	limit, offset := salesOrderPagination(c)
	_, companyID, _, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	sos, total, err := h.service.GetSalesOrdersByStatusForCompany(status, companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	totalAmount := 0.0
	for _, so := range sos {
		totalAmount += so.Total
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"data":         sos,
		"total":        total,
		"total_amount": totalAmount,
	})
}

func (h *SalesOrderHandler) UpdateSalesOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Sales order ID is required", "success": false})
	}

	var soInput input.UpdateSalesOrderInput
	if err := c.BodyParser(&soInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "success": false})
	}

	_, companyID, userID, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	so, err := h.service.UpdateSalesOrderForCompany(id, &soInput, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Sales order updated successfully", "data": so})
}

func (h *SalesOrderHandler) UpdateSalesOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Sales order ID is required", "success": false})
	}

	var statusInput input.UpdateSalesOrderStatusInput
	if err := c.BodyParser(&statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "success": false})
	}
	if err := h.validate.Struct(statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	_, companyID, userID, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	so, err := h.service.UpdateSalesOrderStatusForCompany(id, statusInput.Status, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Sales order status updated successfully", "data": so})
}

func (h *SalesOrderHandler) DeleteSalesOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Sales order ID is required", "success": false})
	}

	_, companyID, _, err := salesOrderContext(c)
	if err != nil {
		return err
	}

	if err := h.service.DeleteSalesOrderForCompany(id, companyID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "success": false})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Sales order deleted successfully"})
}

func salesOrderLocalToUint(value interface{}) (uint, error) {
	if value == nil {
		return 0, fmt.Errorf("value is missing")
	}

	switch v := value.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case int8:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case int16:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case int32:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case int64:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case float32:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case float64:
		if v <= 0 {
			return 0, fmt.Errorf("value must be positive")
		}
		return uint(v), nil
	case json.Number:
		parsed, err := strconv.ParseUint(v.String(), 10, 64)
		return uint(parsed), err
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		return uint(parsed), err
	default:
		parsed, err := strconv.ParseUint(fmt.Sprintf("%v", v), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unsupported numeric type %T with value %v", value, value)
		}
		return uint(parsed), nil
	}
}
