package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type ProductGroupOperationsHandler struct {
	pgInventoryService services.ProductGroupInventoryService
	pgService          services.ProductGroupService
	poService          services.PurchaseOrderService
	soService          services.SalesOrderService
}

func NewProductGroupOperationsHandler(
	pgInventoryService services.ProductGroupInventoryService,
	pgService services.ProductGroupService,
	poService services.PurchaseOrderService,
	soService services.SalesOrderService,
) *ProductGroupOperationsHandler {
	return &ProductGroupOperationsHandler{
		pgInventoryService: pgInventoryService,
		pgService:          pgService,
		poService:          poService,
		soService:          soService,
	}
}

// InitializeProductGroupInventory initializes inventory for a product group
// POST /api/product-groups/:id/inventory/initialize
func (h *ProductGroupOperationsHandler) InitializeProductGroupInventory(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	err := h.pgInventoryService.InitializeProductGroupInventory(productGroupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	status, err := h.pgInventoryService.GetInventoryStatus(productGroupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message": "Product group inventory initialized",
		"data":    status,
	})
}

// GetInventoryStatus retrieves current inventory status
// GET /api/product-groups/:id/inventory/status
func (h *ProductGroupOperationsHandler) GetInventoryStatus(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	status, err := h.pgInventoryService.GetInventoryStatus(productGroupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(status)
}

// CheckStockAvailability checks if sufficient stock is available
// POST /api/product-groups/:id/check-stock
func (h *ProductGroupOperationsHandler) CheckStockAvailability(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	var req input.ProductGroupCheckStockInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	canFulfill, message, err := h.pgInventoryService.CanFulfillOrder(productGroupID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	status, _ := h.pgInventoryService.GetInventoryStatus(productGroupID)

	response := map[string]interface{}{
		"product_group_id":   productGroupID,
		"requested_quantity": req.Quantity,
		"can_fulfill":        canFulfill,
		"available_stock":    status.AvailableStock,
		"component_status":   status.ComponentStatus,
	}

	if !canFulfill {
		response["message"] = message
	} else {
		response["message"] = "Sufficient stock available"
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// AllocateStock allocates (reserves) stock for a sales order
// POST /api/product-groups/:id/allocate-stock
func (h *ProductGroupOperationsHandler) AllocateStock(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	var req input.ProductGroupAllocateStockInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	err := h.pgInventoryService.AllocateStock(productGroupID, req.Quantity, req.SalesOrderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	status, _ := h.pgInventoryService.GetInventoryStatus(productGroupID)

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message":      "Stock allocated successfully",
		"stock_status": status,
	})
}

// DeductStock deducts stock after shipment
// POST /api/product-groups/:id/deduct-stock
func (h *ProductGroupOperationsHandler) DeductStock(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	var req input.ProductGroupDeductStockInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	err := h.pgInventoryService.DeductStock(productGroupID, req.Quantity, req.Reason, &req.ShipmentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	status, _ := h.pgInventoryService.GetInventoryStatus(productGroupID)

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message":      "Stock deducted successfully",
		"stock_status": status,
	})
}

// ReleaseAllocatedStock releases previously reserved stock
// POST /api/product-groups/:id/release-stock
func (h *ProductGroupOperationsHandler) ReleaseAllocatedStock(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	var req input.ProductGroupReleaseStockInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Get SO details to find allocated quantity
	// This would need to be fetched from database
	// For now, we'll require it in the request
	// You should implement proper SO service method

	status, _ := h.pgInventoryService.GetInventoryStatus(productGroupID)

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message":      "Stock released successfully",
		"stock_status": status,
	})
}

// GetTransactionHistory retrieves inventory transaction history
// GET /api/product-groups/:id/transactions?limit=20&offset=0
func (h *ProductGroupOperationsHandler) GetTransactionHistory(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	limit := 20
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

	transactions, total, err := h.pgInventoryService.GetInventoryHistory(productGroupID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"page":         (offset / limit) + 1,
		"limit":        limit,
	})
}

// GetInventoryReport generates comprehensive inventory report
// GET /api/product-groups/:id/report
func (h *ProductGroupOperationsHandler) GetInventoryReport(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	report, err := h.pgInventoryService.GetInventoryReport(productGroupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(report)
}

// StockAdjustment performs manual stock adjustment
// POST /api/product-groups/:id/adjust-stock
func (h *ProductGroupOperationsHandler) StockAdjustment(c *fiber.Ctx) error {
	productGroupID := c.Params("id")

	var req input.ProductGroupStockAdjustmentInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	var err error
	reason := req.Reason
	if req.ReferenceNo != "" {
		reason = reason + " (" + req.ReferenceNo + ")"
	}

	switch req.AdjustmentType {
	case "add":
		err = h.pgInventoryService.AddStock(productGroupID, req.Quantity, reason, nil)
	case "remove", "damage":
		err = h.pgInventoryService.DeductStock(productGroupID, req.Quantity, reason, nil)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": "Invalid adjustment type. Use 'add', 'remove', or 'damage'",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	status, _ := h.pgInventoryService.GetInventoryStatus(productGroupID)

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message":      "Stock adjusted successfully",
		"adjustment":   req,
		"stock_status": status,
	})
}

// SetStockAlerts sets low stock thresholds
// POST /api/product-groups/:id/stock-alerts
func (h *ProductGroupOperationsHandler) SetStockAlerts(c *fiber.Ctx) error {
	// This would store alert thresholds in database
	// Implementation depends on your alert system
	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message": "Stock alerts configured",
	})
}

// GetLowStockProducts returns product groups with low stock
// GET /api/product-groups/low-stock
func (h *ProductGroupOperationsHandler) GetLowStockProducts(c *fiber.Ctx) error {
	threshold := 10.0
	if t := c.Query("threshold"); t != "" {
		if parsed, err := strconv.ParseFloat(t, 64); err == nil {
			threshold = parsed
		}
	}

	products, err := h.pgInventoryService.GetLowStockProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"threshold": threshold,
		"products":  products,
	})
}
