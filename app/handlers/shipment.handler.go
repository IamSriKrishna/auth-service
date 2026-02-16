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

// CreateShipment creates a new shipment
// @Summary Create a new shipment
// @Description Create a new shipment for delivery
// @Tags Shipment
// @Accept json
// @Produce json
// @Param shipment body input.CreateShipmentInput true "Shipment input"
// @Success 201 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /shipments [post]
func (h *ShipmentHandler) CreateShipment(c *fiber.Ctx) error {
	var shipInput input.CreateShipmentInput

	if err := c.BodyParser(&shipInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	// Validate input
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

// GetShipment retrieves a shipment by ID
// @Summary Get shipment by ID
// @Description Get a specific shipment by its ID
// @Tags Shipment
// @Produce json
// @Param id path string true "Shipment ID"
// @Success 200 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /shipments/{id} [get]
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

// GetAllShipments retrieves all shipments with pagination
// @Summary Get all shipments
// @Description Get all shipments with pagination support
// @Tags Shipment
// @Produce json
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /shipments [get]
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

// GetShipmentsByCustomer retrieves shipments for a specific customer
// @Summary Get shipments by customer
// @Description Get all shipments for a specific customer
// @Tags Shipment
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /shipments/customer/{customer_id} [get]
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

// GetShipmentsByPackage retrieves shipments for a specific package
// @Summary Get shipments by package
// @Description Get all shipments for a specific package
// @Tags Shipment
// @Produce json
// @Param package_id path string true "Package ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /shipments/package/{package_id} [get]
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

// GetShipmentsBySalesOrder retrieves shipments for a specific sales order
// @Summary Get shipments by sales order
// @Description Get all shipments for a specific sales order
// @Tags Shipment
// @Produce json
// @Param sales_order_id path string true "Sales Order ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /shipments/sales-order/{sales_order_id} [get]
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

// GetShipmentsByStatus retrieves shipments by status
// @Summary Get shipments by status
// @Description Get all shipments with a specific status
// @Tags Shipment
// @Produce json
// @Param status query string true "Status (created, shipped, in_transit, delivered, cancelled)"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /shipments/status/{status} [get]
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

// UpdateShipment updates a shipment
// @Summary Update a shipment
// @Description Update shipment details
// @Tags Shipment
// @Accept json
// @Produce json
// @Param id path string true "Shipment ID"
// @Param shipment body input.UpdateShipmentInput true "Shipment input"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /shipments/{id} [put]
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

// UpdateShipmentStatus updates shipment status
// @Summary Update shipment status
// @Description Update the status of a shipment
// @Tags Shipment
// @Accept json
// @Produce json
// @Param id path string true "Shipment ID"
// @Param status body input.UpdateShipmentStatusInput true "Status input"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /shipments/{id}/status [patch]
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

// DeleteShipment deletes a shipment
// @Summary Delete a shipment
// @Description Delete a shipment by ID
// @Tags Shipment
// @Produce json
// @Param id path string true "Shipment ID"
// @Success 200 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /shipments/{id} [delete]
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
