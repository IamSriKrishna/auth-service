package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type StockManagementHandler struct {
	service          services.StockManagementService
	variantStockMgmt services.VariantStockManagementService
	userRepo         repo.UserRepository
	productRepo      repo.ProductRepository
}

func NewStockManagementHandler(service services.StockManagementService, variantStockMgmt services.VariantStockManagementService, productRepo repo.ProductRepository) *StockManagementHandler {
	return &StockManagementHandler{
		service:          service,
		variantStockMgmt: variantStockMgmt,
		productRepo:      productRepo,
	}
}

func NewStockManagementHandlerWithUserRepo(service services.StockManagementService, variantStockMgmt services.VariantStockManagementService, userRepo repo.UserRepository, productRepo repo.ProductRepository) *StockManagementHandler {
	return &StockManagementHandler{
		service:          service,
		variantStockMgmt: variantStockMgmt,
		userRepo:         userRepo,
		productRepo:      productRepo,
	}
}

func (h *StockManagementHandler) shouldIncludeProduct(productID string) bool {
	if h.productRepo == nil {
		return true
	}

	if productID == "" {
		return false
	}

	if len(productID) >= 3 && productID[:3] == "pg_" {
		return true
	}

	product, err := h.productRepo.FindByID(productID)
	if err != nil {
		return false
	}

	return !product.IsRaw
}

// extractViewUserID validates the view_user_id parameter and checks permissions
// Always returns authenticated user ID by default, or view_user_id if provided and authorized
func (h *StockManagementHandler) extractViewUserID(c *fiber.Ctx) (uint, error) {
	// Get authenticated user info
	authenticatedUserID := uint(0)
	authenticatedUserType := ""

	if id := c.Locals("user_id"); id != nil {
		switch v := id.(type) {
		case uint:
			authenticatedUserID = v
		case int:
			authenticatedUserID = uint(v)
		case float64:
			authenticatedUserID = uint(v)
		case string:
			fmt.Sscanf(v, "%d", &authenticatedUserID)
		}
	}

	if ut := c.Locals("user_type"); ut != nil {
		if userTypeStr, ok := ut.(string); ok {
			authenticatedUserType = userTypeStr
		}
	}

	// Default to authenticated user
	userID := authenticatedUserID

	// Check if view_user_id is provided in query parameters
	viewUserIDStr := c.Query("view_user_id")
	if viewUserIDStr != "" {
		var viewUserID uint64
		_, parseErr := fmt.Sscanf(viewUserIDStr, "%d", &viewUserID)
		if parseErr != nil || viewUserID <= 0 {
			return 0, fmt.Errorf("invalid view_user_id parameter")
		}

		// Permission check: Only superadmin can view any user's stock, others can only view their own
		if authenticatedUserType != "superadmin" && uint(viewUserID) != authenticatedUserID {
			return 0, fmt.Errorf("unauthorized: cannot view another user's stock")
		}

		// Validate that the user exists
		if h.userRepo != nil {
			user, getErr := h.userRepo.GetByID(uint(viewUserID))
			if getErr != nil || user == nil {
				return 0, fmt.Errorf("user not found with id: %d", viewUserID)
			}
		}

		userID = uint(viewUserID)
	}

	return userID, nil
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

// GET /api/stock/summary/raw-materials
// Get stock summary for raw-material products only
func (h *StockManagementHandler) GetRawMaterialStocksSummary(c *fiber.Ctx) error {
	// Validate view_user_id if provided
	viewUserID, err := h.extractViewUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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

	stocks, total, getAllErr := h.service.GetAllRawMaterialProductStockByUser(viewUserID, offset, limit)
	if getAllErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": getAllErr.Error(),
		})
	}

	stocksResponse := make([]fiber.Map, 0)
	damagedResponse := make([]fiber.Map, 0)
	totalStockValue := 0.0
	totalSoldProductValue := 0.0
	totalDamagedValue := 0.0
	damagedCount := 0

	for _, stock := range stocks {
		if stock.DamagedStock > 0 {
			damagedValue := stock.DamagedStock * stock.AverageCost
			totalDamagedValue += damagedValue
			damagedCount++

			damagedResponse = append(damagedResponse, fiber.Map{
				"product_id":        stock.ProductID,
				"product_name":      stock.ProductName,
				"damaged_stock":     stock.DamagedStock,
				"average_cost":      stock.AverageCost,
				"damaged_value":     damagedValue,
				"damage_reason":     stock.DamageReason,
				"damaged_at":        stock.DamagedAt,
				"damaged_by":        stock.DamagedBy,
				"raw_material_unit": stock.RawMaterialUnit,
				"required_gram_per_unit": func() float64 {
					if stock.Product != nil {
						return stock.Product.RequiredGramPerUnit
					}
					return 0
				}(),
				"type": "product",
			})
		}

		stockValue := stock.CurrentStock * stock.AverageCost
		totalStockValue += stockValue

		soldProductValue := stock.SoldStock * stock.AverageCost
		totalSoldProductValue += soldProductValue

		totalPieces := computeRawMaterialPieces(&stock)

		stocksResponse = append(stocksResponse, fiber.Map{
			"product_id":        stock.ProductID,
			"product_name":      stock.ProductName,
			"current_stock":     stock.CurrentStock,
			"raw_material_unit": stock.RawMaterialUnit,
			"required_gram_per_unit": func() float64 {
				if stock.Product != nil {
					return stock.Product.RequiredGramPerUnit
				}
				return 0
			}(),
			"purchased_total": stock.PurchasedStock,
			"sold_total":      stock.SoldStock,
			"reserved_stock":  stock.ReservedStock,
			"available_stock": stock.AvailableStock,
			"damaged_stock":   stock.DamagedStock,
			"average_cost":    stock.AverageCost,
			"stock_value":     stockValue,
			"total_pieces":    totalPieces,
			"last_purchased":  stock.LastPurchasedDate,
			"last_sold":       stock.LastSoldDate,
			"type":            "product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stocks":                   stocksResponse,
		"damaged_products":         damagedResponse,
		"damaged_count":            damagedCount,
		"total":                    total,
		"total_stock_value":        totalStockValue,
		"total_sold_product_value": totalSoldProductValue,
		"total_damaged_value":      totalDamagedValue,
	})
}

func computeRawMaterialPieces(stock *models.ProductStock) int64 {
	if stock == nil || stock.CurrentStock <= 0 || stock.Product == nil || stock.Product.RequiredGramPerUnit <= 0 {
		return 0
	}

	unit := strings.TrimSpace(strings.ToLower(stock.RawMaterialUnit))
	if unit == "" {
		unit = strings.TrimSpace(strings.ToLower(stock.Product.RawUnit))
	}

	stockInGrams := stock.CurrentStock
	switch unit {
	case "kg", "kilogram", "kilograms":
		stockInGrams *= 1000
	case "mg", "milligram", "milligrams":
		stockInGrams /= 1000
	case "g", "gram", "grams":
		// already grams
	default:
		// Unknown unit: cannot reliably convert to grams
		return 0
	}

	return int64(stockInGrams / stock.Product.RequiredGramPerUnit)
}

// GET /api/stock/summary
// Get updated stock summary for all products
func (h *StockManagementHandler) GetAllStocksSummary(c *fiber.Ctx) error {
	// Validate view_user_id if provided
	viewUserID, err := h.extractViewUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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

	// Get product stocks - always filtered by authenticated user
	var stocks []models.ProductStock
	var total int64
	var getAllErr error

	stocks, total, getAllErr = h.service.GetAllProductStockByUser(viewUserID, offset, limit)

	if getAllErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": getAllErr.Error(),
		})
	}

	// Build the response with stock value calculations
	stocksResponse := make([]fiber.Map, 0)
	damagedResponse := make([]fiber.Map, 0)
	totalStockValue := 0.0
	totalSoldProductValue := 0.0
	totalDamagedValue := 0.0
	damagedCount := 0

	// Add product stocks
	for _, stock := range stocks {
		if !h.shouldIncludeProduct(stock.ProductID) {
			continue
		}

		if stock.DamagedStock > 0 {
			// Add to damaged products list
			damagedValue := stock.DamagedStock * stock.AverageCost
			totalDamagedValue += damagedValue
			damagedCount++

			damagedResponse = append(damagedResponse, fiber.Map{
				"product_id":    stock.ProductID,
				"product_name":  stock.ProductName,
				"sku":           stock.SKU,
				"damaged_stock": stock.DamagedStock,
				"average_cost":  stock.AverageCost,
				"damaged_value": damagedValue,
				"damage_reason": stock.DamageReason,
				"damaged_at":    stock.DamagedAt,
				"damaged_by":    stock.DamagedBy,
				"type":          "product",
			})
		}

		// Add to regular stock (without damaged quantity)
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
			"damaged_stock":   stock.DamagedStock,
			"average_cost":    stock.AverageCost,
			"stock_value":     stockValue,
			"last_purchased":  stock.LastPurchasedDate,
			"last_sold":       stock.LastSoldDate,
			"type":            "product",
		})
	}

	// Get variant stocks - always filtered by authenticated user
	var variantStocks []models.VariantStock
	var variantGetErr error

	variantStocks, _, variantGetErr = h.variantStockMgmt.GetAllVariantStocksByUser(viewUserID, offset, limit)

	if variantGetErr != nil {
		// Log error but continue - variants are optional
		log.Printf("[STOCK_SUMMARY] Error fetching variant stocks: %v", variantGetErr)
	} else {
		// Add variant stocks
		for _, vStock := range variantStocks {
			if !h.shouldIncludeProduct(vStock.ProductID) {
				continue
			}

			if vStock.DamagedStock > 0 {
				// Add to damaged variants list
				damagedValue := vStock.DamagedStock * vStock.AverageCost
				totalDamagedValue += damagedValue
				damagedCount++

				damagedResponse = append(damagedResponse, fiber.Map{
					"product_id":    vStock.ProductID,
					"product_name":  vStock.ProductName,
					"variant_sku":   vStock.VariantSKU,
					"variant_name":  vStock.VariantName,
					"damaged_stock": vStock.DamagedStock,
					"average_cost":  vStock.AverageCost,
					"damaged_value": damagedValue,
					"damage_reason": vStock.DamageReason,
					"damaged_at":    vStock.DamagedAt,
					"damaged_by":    vStock.DamagedBy,
					"type":          "variant",
				})
			}

			// Add to regular stock
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
				"damaged_stock":   vStock.DamagedStock,
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
		"damaged_products":         damagedResponse,
		"damaged_count":            damagedCount,
		"total":                    total,
		"total_stock_value":        totalStockValue,
		"total_sold_product_value": totalSoldProductValue,
		"total_damaged_value":      totalDamagedValue,
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

// PATCH /api/stock/mark-damaged
// Mark product or variant as damaged
type MarkDamagedInput struct {
	ProductID  string  `json:"product_id" validate:"required"`
	VariantSKU *string `json:"variant_sku"`
	Quantity   float64 `json:"quantity" validate:"required,gt=0"`
	Reason     string  `json:"reason" validate:"required"`
}

func (h *StockManagementHandler) MarkProductAsDamaged(c *fiber.Ctx) error {
	var input MarkDamagedInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		// user_id can be uint or string depending on middleware
		switch v := uid.(type) {
		case string:
			userID = v
		case uint:
			userID = strconv.FormatUint(uint64(v), 10)
		case int:
			userID = strconv.Itoa(v)
		case float64:
			userID = fmt.Sprintf("%d", int(v))
		}
	}

	// Mark variant as damaged if VariantSKU is provided
	if input.VariantSKU != nil && *input.VariantSKU != "" {
		if err := h.variantStockMgmt.MarkVariantAsDamaged(*input.VariantSKU, input.Quantity, input.Reason, userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		// Get updated variant stock
		updatedStock, err := h.variantStockMgmt.GetVariantStockSummary(*input.VariantSKU)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch updated stock",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":         true,
			"message":         "Variant marked as damaged successfully",
			"type":            "variant",
			"variant_sku":     updatedStock.VariantSKU,
			"damaged_stock":   updatedStock.DamagedStock,
			"available_stock": updatedStock.AvailableStock,
			"damage_reason":   updatedStock.DamageReason,
			"damaged_at":      updatedStock.DamagedAt,
		})
	}

	// Mark product as damaged
	if err := h.service.MarkProductAsDamaged(input.ProductID, input.Quantity, input.Reason, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get updated product stock
	updatedStock, err := h.service.GetProductStock(input.ProductID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch updated stock",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":         true,
		"message":         "Product marked as damaged successfully",
		"type":            "product",
		"product_id":      updatedStock.ProductID,
		"damaged_stock":   updatedStock.DamagedStock,
		"available_stock": updatedStock.AvailableStock,
		"damage_reason":   updatedStock.DamageReason,
		"damaged_at":      updatedStock.DamagedAt,
	})
}

// GET /api/stock/damaged
// Get all damaged products and variants
func (h *StockManagementHandler) GetDamagedProducts(c *fiber.Ctx) error {
	limit := 50
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

	// Get damaged products
	products, productTotal, err := h.service.GetDamagedProducts(offset, limit)
	if err != nil {
		log.Printf("[DAMAGE_TRACKING] Error fetching damaged products: %v", err)
	}

	// Get damaged variants
	variants, variantTotal, variantErr := h.variantStockMgmt.GetDamagedVariants(offset, limit)
	if variantErr != nil {
		log.Printf("[DAMAGE_TRACKING] Error fetching damaged variants: %v", variantErr)
	}

	// Build response
	damagedItems := make([]fiber.Map, 0)
	totalDamagedValue := 0.0

	// Add damaged products
	for _, product := range products {
		damagedValue := product.DamagedStock * product.AverageCost
		totalDamagedValue += damagedValue

		damagedItems = append(damagedItems, fiber.Map{
			"type":          "product",
			"product_id":    product.ProductID,
			"product_name":  product.ProductName,
			"sku":           product.SKU,
			"damaged_stock": product.DamagedStock,
			"damage_reason": product.DamageReason,
			"damaged_at":    product.DamagedAt,
			"damaged_by":    product.DamagedBy,
			"average_cost":  product.AverageCost,
			"damaged_value": damagedValue,
		})
	}

	// Add damaged variants
	for _, variant := range variants {
		damagedValue := variant.DamagedStock * variant.AverageCost
		totalDamagedValue += damagedValue

		damagedItems = append(damagedItems, fiber.Map{
			"type":          "variant",
			"product_id":    variant.ProductID,
			"product_name":  variant.ProductName,
			"variant_sku":   variant.VariantSKU,
			"variant_name":  variant.VariantName,
			"damaged_stock": variant.DamagedStock,
			"damage_reason": variant.DamageReason,
			"damaged_at":    variant.DamagedAt,
			"damaged_by":    variant.DamagedBy,
			"average_cost":  variant.AverageCost,
			"damaged_value": damagedValue,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":             true,
		"damaged_items":       damagedItems,
		"damaged_products":    productTotal,
		"damaged_variants":    variantTotal,
		"total_damaged":       productTotal + variantTotal,
		"total_damaged_value": totalDamagedValue,
		"limit":               limit,
		"offset":              offset,
	})
}
