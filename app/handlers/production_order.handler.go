package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type ProductionOrderHandler struct {
	service services.ProductionOrderService
}

func NewProductionOrderHandler(service services.ProductionOrderService) *ProductionOrderHandler {
	return &ProductionOrderHandler{service: service}
}

// CreateProductionOrder godoc
// @Summary Create a new production order
// @Description Create a production order to manufacture products from components
// @Tags Production Orders
// @Accept json
// @Produce json
// @Param request body input.CreateProductionOrderInput true "Production Order Request"
// @Success 201 {object} map[string]interface{} "Success response with production order data"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /production-orders [post]
func (h *ProductionOrderHandler) CreateProductionOrder(c *fiber.Ctx) error {
	var req input.CreateProductionOrderInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	prodOrder, err := h.service.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    prodOrder,
		"message": "Production Order created successfully",
	})
}

// GetProductionOrderByID godoc
// @Summary Get production order by ID
// @Description Retrieve a production order by its ID
// @Tags Production Orders
// @Produce json
// @Param id path string true "Production Order ID"
// @Success 200 {object} map[string]interface{} "Success response with production order data"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Router /production-orders/{id} [get]
func (h *ProductionOrderHandler) GetProductionOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "production order id is required",
		})
	}

	prodOrder, err := h.service.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    prodOrder,
		"message": "Production Order retrieved successfully",
	})
}

// GetAllProductionOrders godoc
// @Summary Get all production orders
// @Description Retrieve all production orders with pagination
// @Tags Production Orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "Success response with production orders list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /production-orders [get]
func (h *ProductionOrderHandler) GetAllProductionOrders(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	result, err := h.service.FindAll(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result.ProductionOrders,
		"pagination": fiber.Map{
			"page":        result.Page,
			"limit":       result.Limit,
			"total":       result.Total,
			"total_pages": result.TotalPages,
		},
		"message": "Production Orders retrieved successfully",
	})
}

// UpdateProductionOrder godoc
// @Summary Update production order
// @Description Update production order status and details
// @Tags Production Orders
// @Accept json
// @Produce json
// @Param id path string true "Production Order ID"
// @Param request body input.UpdateProductionOrderInput true "Update Request"
// @Success 200 {object} map[string]interface{} "Success response with updated production order"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Router /production-orders/{id} [put]
func (h *ProductionOrderHandler) UpdateProductionOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "production order id is required",
		})
	}

	var req input.UpdateProductionOrderInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	prodOrder, err := h.service.Update(id, &req)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "production order not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    prodOrder,
		"message": "Production Order updated successfully",
	})
}

// DeleteProductionOrder godoc
// @Summary Delete production order
// @Description Delete a production order by ID
// @Tags Production Orders
// @Produce json
// @Param id path string true "Production Order ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Router /production-orders/{id} [delete]
func (h *ProductionOrderHandler) DeleteProductionOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "production order id is required",
		})
	}

	result, err := h.service.Delete(id)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "production order not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "Production Order deleted successfully",
	})
}

// ConsumeProductionOrderItem godoc
// @Summary Consume/use items in production order
// @Description Mark items as consumed during the production process
// @Tags Production Orders
// @Accept json
// @Produce json
// @Param id path string true "Production Order ID"
// @Param request body input.ConsumeProductionOrderItemInput true "Consume Item Request"
// @Success 200 {object} map[string]interface{} "Success response with updated production order"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Router /production-orders/{id}/consume-item [post]
func (h *ProductionOrderHandler) ConsumeProductionOrderItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "production order id is required",
		})
	}

	var req input.ConsumeProductionOrderItemInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	result, err := h.service.ConsumeItem(id, &req)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err.Error() == "production order not found" || err.Error() == "production order item not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "Item consumed successfully",
	})
}
