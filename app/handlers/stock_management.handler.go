package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type StockManagementHandler struct {
	service          services.StockManagementService
	variantStockMgmt services.VariantStockManagementService
}

func NewStockManagementHandler(service services.StockManagementService, variantStockMgmt services.VariantStockManagementService) *StockManagementHandler {
	return &StockManagementHandler{
		service:          service,
		variantStockMgmt: variantStockMgmt,
	}
}

// GET /api/stock/products/:productId
// Get current stock for a product
func (h *StockManagementHandler) GetProductStock(c *fiber.Ctx) error {
	productID := c.Params("productId")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Product ID is required",
		})
	}

	stock, err := h.service.GetProductStock(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stock,
	})
}

// GET /api/stock/summary/:productId
// Get comprehensive stock summary for a product
func (h *StockManagementHandler) GetStockSummary(c *fiber.Ctx) error {
	productID := c.Params("productId")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Product ID is required",
		})
	}

	summary, err := h.service.GetStockSummary(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    summary,
	})
}

// GET /api/stock/all
// Get all product stocks with pagination
func (h *StockManagementHandler) GetAllStocks(c *fiber.Ctx) error {
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

	stocks, total, err := h.service.GetAllProductStock(offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stocks,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GET /api/stock/low-stock
// Get products with stock below threshold
func (h *StockManagementHandler) GetLowStockProducts(c *fiber.Ctx) error {
	threshold := 10.0
	limit := 10
	offset := 0

	if t := c.Query("threshold"); t != "" {
		if parsed, err := strconv.ParseFloat(t, 64); err == nil {
			threshold = parsed
		}
	}

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

	stocks, total, err := h.service.GetLowStockProducts(threshold, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stocks,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GET /api/stock/history/:productId
// Get movement history for a product
func (h *StockManagementHandler) GetProductHistory(c *fiber.Ctx) error {
	productID := c.Params("productId")
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

	movements, total, err := h.service.GetProductMovementHistory(productID, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    movements,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GET /api/stock/report/date-range
// Get stock movements within a date range
func (h *StockManagementHandler) GetMovementsByDateRange(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")
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

	fromDate, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid 'from' date format (use YYYY-MM-DD)",
		})
	}

	toDate, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid 'to' date format (use YYYY-MM-DD)",
		})
	}

	movements, total, err := h.service.GetMovementsByDateRange(fromDate, toDate.AddDate(0, 0, 1), offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    movements,
		"total":   total,
		"from":    fromDate,
		"to":      toDate,
		"limit":   limit,
		"offset":  offset,
	})
}

// POST /api/stock/adjustment
// Record manual stock adjustment
type StockAdjustmentInput struct {
	ProductID      string  `json:"product_id" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
	AdjustmentType string  `json:"adjustment_type" validate:"required,oneof=in out"`
	Reason         string  `json:"reason" validate:"required"`
}

func (h *StockManagementHandler) RecordAdjustment(c *fiber.Ctx) error {
	var input StockAdjustmentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = uid.(string)
	}

	if err := h.service.RecordStockAdjustment(input.ProductID, input.Quantity, input.AdjustmentType, input.Reason, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Stock adjustment recorded successfully",
	})
}

// GET /api/stock/summary
// Get updated stock summary for all products
func (h *StockManagementHandler) GetAllStocksSummary(c *fiber.Ctx) error {
	limit := 100
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

	// Get product stocks
	stocks, _, err := h.service.GetAllProductStock(offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Build the response with stock value calculations
	stocksResponse := make([]fiber.Map, 0)
	totalStockValue := 0.0
	totalSoldProductValue := 0.0

	// Add product stocks
	for _, stock := range stocks {
		stockValue := stock.CurrentStock * stock.AverageCost
		totalStockValue += stockValue

		soldProductValue := stock.SoldStock * stock.AverageCost
		totalSoldProductValue += soldProductValue

		stocksResponse = append(stocksResponse, fiber.Map{
			"product_id":      stock.ProductID,
			"product_name":    stock.ProductName,
			"sku":             stock.SKU,
			"current_stock":   stock.CurrentStock,
			"purchased_total": stock.PurchasedStock,
			"sold_total":      stock.SoldStock,
			"reserved_stock":  stock.ReservedStock,
			"available_stock": stock.AvailableStock,
			"average_cost":    stock.AverageCost,
			"stock_value":     stockValue,
			"last_purchased":  stock.LastPurchasedDate,
			"last_sold":       stock.LastSoldDate,
			"type":            "product",
		})
	}

	// Get variant stocks
	variantStocks, _, err := h.variantStockMgmt.GetAllVariantStocks(offset, limit)
	if err != nil {
		// Log error but continue - variants are optional
		log.Printf("[STOCK_SUMMARY] Error fetching variant stocks: %v", err)
	} else {
		// Add variant stocks
		for _, vStock := range variantStocks {
			stockValue := vStock.CurrentStock * vStock.AverageCost
			totalStockValue += stockValue

			soldProductValue := vStock.SoldStock * vStock.AverageCost
			totalSoldProductValue += soldProductValue

			stocksResponse = append(stocksResponse, fiber.Map{
				"product_id":      vStock.ProductID,
				"product_name":    vStock.ProductName,
				"sku":             vStock.VariantSKU,
				"variant_name":    vStock.VariantName,
				"current_stock":   vStock.CurrentStock,
				"purchased_total": vStock.PurchasedStock,
				"sold_total":      vStock.SoldStock,
				"reserved_stock":  vStock.ReservedStock,
				"available_stock": vStock.AvailableStock,
				"average_cost":    vStock.AverageCost,
				"stock_value":     stockValue,
				"last_purchased":  vStock.LastPurchasedDate,
				"last_sold":       vStock.LastSoldDate,
				"type":            "variant",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stocks":                   stocksResponse,
		"total_stock_value":        totalStockValue,
		"total_sold_product_value": totalSoldProductValue,
	})
}

// GET /api/stock/product/:product_id/movements
// Get product movement history
func (h *StockManagementHandler) GetProductMovements(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
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

	movements, total, err := h.service.GetProductMovementHistory(productID, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Build movements response with correct structure
	movementsResponse := make([]fiber.Map, len(movements))
	for i, movement := range movements {
		movementsResponse[i] = fiber.Map{
			"id":                 movement.ID,
			"product_id":         movement.ProductID,
			"movement_type":      movement.MovementType,
			"quantity":           movement.Quantity,
			"rate":               movement.Rate,
			"amount":             movement.Amount,
			"reference_type":     movement.ReferenceType,
			"reference_id":       movement.ReferenceID,
			"reference_number":   movement.ReferenceNumber,
			"balance_before_qty": movement.BalanceBeforeQty,
			"balance_after_qty":  movement.BalanceAfterQty,
			"notes":              movement.Notes,
			"created_at":         movement.CreatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"movements": movementsResponse,
		"total":     total,
	})
}

// GET /api/stock/debug/product/:product_id
// Get detailed stock breakdown and audit trail
func (h *StockManagementHandler) GetProductStockDebug(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	// Get current stock
	stock, err := h.service.GetProductStock(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get all movements
	movements, total, err := h.service.GetProductMovementHistory(productID, 0, 1000)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Build detailed audit trail
	movementsBreakdown := make([]fiber.Map, 0)
	inboundTotal := 0.0
	outboundTotal := 0.0

	for _, movement := range movements {
		if movement.Quantity > 0 {
			inboundTotal += movement.Quantity
		} else {
			outboundTotal += movement.Quantity // This is already negative
		}

		movementsBreakdown = append(movementsBreakdown, fiber.Map{
			"id":                 movement.ID,
			"movement_type":      movement.MovementType,
			"quantity":           movement.Quantity,
			"rate":               movement.Rate,
			"amount":             movement.Amount,
			"reference_type":     movement.ReferenceType,
			"reference_id":       movement.ReferenceID,
			"reference_number":   movement.ReferenceNumber,
			"balance_before_qty": movement.BalanceBeforeQty,
			"balance_after_qty":  movement.BalanceAfterQty,
			"notes":              movement.Notes,
			"created_at":         movement.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product_id":       productID,
		"current_stock":    stock.CurrentStock,
		"purchased_total":  stock.PurchasedStock,
		"sold_total":       stock.SoldStock,
		"reserved_stock":   stock.ReservedStock,
		"available_stock":  stock.AvailableStock,
		"average_cost":     stock.AverageCost,
		"inbound_total":    inboundTotal,
		"outbound_total":   outboundTotal,
		"net_movement":     inboundTotal + outboundTotal,
		"expected_current": inboundTotal + outboundTotal,
		"movements_count":  total,
		"movements":        movementsBreakdown,
	})
}
