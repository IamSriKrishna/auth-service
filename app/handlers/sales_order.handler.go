package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SalesOrderHandler struct {
	service services.SalesOrderService
}

func NewSalesOrderHandler(service services.SalesOrderService) *SalesOrderHandler {
	return &SalesOrderHandler{service: service}
}

func (h *SalesOrderHandler) CreateSalesOrder(c *fiber.Ctx) error {
	var soInput input.CreateSalesOrderInput

	if err := c.BodyParser(&soInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(soInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	so, err := h.service.CreateSalesOrder(&soInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Sales order created successfully",
		"data":    so,
	})
}

func (h *SalesOrderHandler) GetSalesOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	so, err := h.service.GetSalesOrder(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Sales order not found",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    so,
	})
}

func (h *SalesOrderHandler) GetAllSalesOrders(c *fiber.Ctx) error {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	sos, total, err := h.service.GetAllSalesOrders(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    sos,
		"total":   total,
	})
}

func (h *SalesOrderHandler) GetSalesOrdersByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerId")
	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid customer ID",
			"success": false,
		})
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	sos, total, err := h.service.GetSalesOrdersByCustomer(uint(id), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    sos,
		"total":   total,
	})
}

func (h *SalesOrderHandler) GetSalesOrdersByStatus(c *fiber.Ctx) error {
	status := c.Params("status")

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	sos, total, err := h.service.GetSalesOrdersByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    sos,
		"total":   total,
	})
}

func (h *SalesOrderHandler) UpdateSalesOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	var soInput input.UpdateSalesOrderInput
	if err := c.BodyParser(&soInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	so, err := h.service.UpdateSalesOrder(id, &soInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Sales order updated successfully",
		"data":    so,
	})
}

func (h *SalesOrderHandler) UpdateSalesOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var statusInput input.UpdateSalesOrderStatusInput
	if err := c.BodyParser(&statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	so, err := h.service.UpdateSalesOrderStatus(id, statusInput.Status, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Sales order status updated successfully",
		"data":    so,
	})
}

func (h *SalesOrderHandler) DeleteSalesOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteSalesOrder(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Sales order deleted successfully",
	})
}
