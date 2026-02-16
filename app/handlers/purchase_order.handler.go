package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PurchaseOrderHandler handles purchase order-related HTTP requests
type PurchaseOrderHandler struct {
	service  services.PurchaseOrderService
	validate *validator.Validate
}

func NewPurchaseOrderHandler(service services.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreatePurchaseOrder creates a new purchase order
// POST /api/purchase-orders
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c *fiber.Ctx) error {
	var poInput input.CreatePurchaseOrderInput

	if err := c.BodyParser(&poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get user ID from context (assuming auth middleware sets this)
	userID := ""
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	po, err := h.service.CreatePurchaseOrder(&poInput, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Purchase order created successfully",
		"data":    po,
	})
}

// GetPurchaseOrder retrieves a purchase order by ID
// GET /api/purchase-orders/:id
func (h *PurchaseOrderHandler) GetPurchaseOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	po, err := h.service.GetPurchaseOrder(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Purchase order not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    po,
	})
}

// GetAllPurchaseOrders retrieves all purchase orders with pagination
// GET /api/purchase-orders
func (h *PurchaseOrderHandler) GetAllPurchaseOrders(c *fiber.Ctx) error {
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

	pos, err := h.service.GetAllPurchaseOrders(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pos,
	})
}

// UpdatePurchaseOrder updates an existing purchase order
// PUT /api/purchase-orders/:id
func (h *PurchaseOrderHandler) UpdatePurchaseOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var poInput input.UpdatePurchaseOrderInput

	if err := c.BodyParser(&poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	po, err := h.service.UpdatePurchaseOrder(id, &poInput, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Purchase order updated successfully",
		"data":    po,
	})
}

// DeletePurchaseOrder deletes a purchase order
// DELETE /api/purchase-orders/:id
func (h *PurchaseOrderHandler) DeletePurchaseOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeletePurchaseOrder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Purchase order deleted successfully",
	})
}

// GetPurchaseOrdersByVendor retrieves purchase orders for a specific vendor
// GET /api/purchase-orders/vendor/:vendorId
func (h *PurchaseOrderHandler) GetPurchaseOrdersByVendor(c *fiber.Ctx) error {
	vendorID, err := strconv.ParseUint(c.Params("vendorId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vendor ID",
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

	pos, err := h.service.GetPurchaseOrdersByVendor(uint(vendorID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pos,
	})
}

// GetPurchaseOrdersByCustomer retrieves purchase orders for a specific customer
// GET /api/purchase-orders/customer/:customerId
func (h *PurchaseOrderHandler) GetPurchaseOrdersByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Params("customerId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid customer ID",
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

	pos, err := h.service.GetPurchaseOrdersByCustomer(uint(customerID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pos,
	})
}

// GetPurchaseOrdersByStatus retrieves purchase orders filtered by status
// GET /api/purchase-orders/status/:status
func (h *PurchaseOrderHandler) GetPurchaseOrdersByStatus(c *fiber.Ctx) error {
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

	pos, err := h.service.GetPurchaseOrdersByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pos,
	})
}

// UpdatePurchaseOrderStatus updates the status of a purchase order
// PATCH /api/purchase-orders/:id/status
func (h *PurchaseOrderHandler) UpdatePurchaseOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var statusInput input.UpdatePurchaseOrderStatusInput

	if err := c.BodyParser(&statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	po, err := h.service.UpdatePurchaseOrderStatus(id, statusInput.Status, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Purchase order status updated successfully",
		"data":    po,
	})
}
