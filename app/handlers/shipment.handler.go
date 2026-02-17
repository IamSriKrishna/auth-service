package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ShipmentHandler struct {
	service services.ShipmentService
}

func NewShipmentHandler(service services.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{service: service}
}

func (h *ShipmentHandler) CreateShipment(c *fiber.Ctx) error {
	var shipInput input.CreateShipmentInput

	if err := c.BodyParser(&shipInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(shipInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	shipment, err := h.service.CreateShipment(&shipInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Shipment created successfully",
		"data":    shipment,
	})
}

func (h *ShipmentHandler) GetShipment(c *fiber.Ctx) error {
	id := c.Params("id")

	shipment, err := h.service.GetShipment(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Shipment not found",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipment,
	})
}

func (h *ShipmentHandler) GetAllShipments(c *fiber.Ctx) error {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	shipments, total, err := h.service.GetAllShipments(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Params("customer_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid customer ID",
			"success": false,
		})
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	shipments, total, err := h.service.GetShipmentsByCustomer(uint(customerID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsByPackage(c *fiber.Ctx) error {
	packageID := c.Params("package_id")

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	shipments, total, err := h.service.GetShipmentsByPackage(packageID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsBySalesOrder(c *fiber.Ctx) error {
	salesOrderID := c.Params("sales_order_id")

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	shipments, total, err := h.service.GetShipmentsBySalesOrder(salesOrderID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsByStatus(c *fiber.Ctx) error {
	status := c.Params("status")

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	shipments, total, err := h.service.GetShipmentsByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) UpdateShipment(c *fiber.Ctx) error {
	id := c.Params("id")
	var shipInput input.UpdateShipmentInput

	if err := c.BodyParser(&shipInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(shipInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	shipment, err := h.service.UpdateShipment(id, &shipInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Shipment updated successfully",
		"data":    shipment,
	})
}

func (h *ShipmentHandler) UpdateShipmentStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var statusInput input.UpdateShipmentStatusInput

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

	shipment, err := h.service.UpdateShipmentStatus(id, statusInput.Status, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Shipment status updated successfully",
		"data":    shipment,
	})
}

func (h *ShipmentHandler) DeleteShipment(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteShipment(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Shipment deleted successfully",
	})
}
