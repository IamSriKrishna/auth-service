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
type stockAccessContext struct {
	UserID    uint
	UserType  string
	CompanyID uint
}

func (h *StockManagementHandler) extractStockAccessContext(c *fiber.Ctx) (*stockAccessContext, error) {
	ctx := &stockAccessContext{
		UserID:    stockLocalUint(c, "user_id"),
		UserType:  stockLocalString(c, "user_type"),
		CompanyID: stockLocalUint(c, "company_id"),
	}

	if ctx.UserID == 0 {
		return nil, fmt.Errorf("invalid authenticated user")
	}
	if ctx.UserType == "" {
		return nil, fmt.Errorf("invalid authenticated user type")
	}
	if ctx.UserType != "superadmin" && ctx.CompanyID == 0 {
		return nil, fmt.Errorf("user is not assigned to a company")
	}

	viewUserIDStr := c.Query("view_user_id")
	if viewUserIDStr == "" {
		return ctx, nil
	}

	viewUserID64, err := strconv.ParseUint(viewUserIDStr, 10, 64)
	if err != nil || viewUserID64 == 0 {
		return nil, fmt.Errorf("invalid view_user_id parameter")
	}
	if h.userRepo == nil {
		return nil, fmt.Errorf("user repository is unavailable")
	}

	viewUserID := uint(viewUserID64)

	if ctx.UserType == "superadmin" {
		selectedUser, err := h.userRepo.GetByID(viewUserID)
		if err != nil || selectedUser == nil {
			return nil, fmt.Errorf("user not found with id: %d", viewUserID)
		}
		if selectedUser.CompanyID == nil || *selectedUser.CompanyID == 0 {
			return nil, fmt.Errorf("selected user is not assigned to a company")
		}

		ctx.UserID = selectedUser.ID
		ctx.UserType = string(selectedUser.UserType)
		ctx.CompanyID = *selectedUser.CompanyID
		return ctx, nil
	}

	selectedUser, err := h.userRepo.GetByIDAndCompanyID(viewUserID, ctx.CompanyID)
	if err != nil || selectedUser == nil {
		return nil, fmt.Errorf("user not found in your company")
	}

	// Stock visibility is company-level, so keep the authenticated company.
	ctx.UserID = selectedUser.ID
	ctx.UserType = string(selectedUser.UserType)
	return ctx, nil
}

func stockLocalUint(c *fiber.Ctx, key string) uint {
	value := c.Locals(key)
	switch v := value.(type) {
	case uint:
		return v
	case uint64:
		return uint(v)
	case int:
		if v > 0 {
			return uint(v)
		}
	case int64:
		if v > 0 {
			return uint(v)
		}
	case float64:
		if v > 0 {
			return uint(v)
		}
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			return uint(parsed)
		}
	}
	return 0
}

func stockLocalString(c *fiber.Ctx, key string) string {
	value, _ := c.Locals(key).(string)
	return value
}

func stockUserIDString(userID uint) string {
	return strconv.FormatUint(uint64(userID), 10)
}

func stockAccessError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func (h *StockManagementHandler) getCompanyUsers(companyID uint) ([]models.User, error) {
	if h.userRepo == nil {
		return nil, fmt.Errorf("user repository is unavailable")
	}
	return h.userRepo.ListByCompanyID(companyID)
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
	ctx, err := h.extractStockAccessContext(c)
	if err != nil {
		return stockAccessError(c, err)
	}

	limit := 100
	offset := 0
	if value := c.Query("limit"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed > 0 {
			limit = parsed
		}
	}
	if value := c.Query("offset"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed >= 0 {
			offset = parsed
		}
	}

	stocks, total, err := h.service.GetAllRawMaterialProductStockByCompany(
		ctx.CompanyID,
		ctx.UserType,
		offset,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
				"product_id": stock.ProductID, "product_name": stock.ProductName,
				"damaged_stock": stock.DamagedStock, "average_cost": stock.AverageCost,
				"damaged_value": damagedValue, "damage_reason": stock.DamageReason,
				"damaged_at": stock.DamagedAt, "damaged_by": stock.DamagedBy,
				"raw_material_unit": stock.RawMaterialUnit,
				"required_gram_per_unit": func() float64 {
					if stock.Product != nil { return stock.Product.RequiredGramPerUnit }
					return 0
				}(),
				"type": "product",
			})
		}

		stockValue := stock.CurrentStock * stock.AverageCost
		totalStockValue += stockValue
		totalSoldProductValue += stock.SoldStock * stock.AverageCost

		stocksResponse = append(stocksResponse, fiber.Map{
			"product_id": stock.ProductID, "product_name": stock.ProductName,
			"current_stock": stock.CurrentStock, "raw_material_unit": stock.RawMaterialUnit,
			"required_gram_per_unit": func() float64 {
				if stock.Product != nil { return stock.Product.RequiredGramPerUnit }
				return 0
			}(),
			"purchased_total": stock.PurchasedStock, "sold_total": stock.SoldStock,
			"reserved_stock": stock.ReservedStock, "available_stock": stock.AvailableStock,
			"damaged_stock": stock.DamagedStock, "average_cost": stock.AverageCost,
			"stock_value": stockValue, "total_pieces": computeRawMaterialPieces(&stock),
			"last_purchased": stock.LastPurchasedDate, "last_sold": stock.LastSoldDate,
			"type": "product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stocks": stocksResponse, "damaged_products": damagedResponse,
		"damaged_count": damagedCount, "total": total,
		"total_stock_value": totalStockValue,
		"total_sold_product_value": totalSoldProductValue,
		"total_damaged_value": totalDamagedValue,
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
	case "g", "gram", "grams", "piece", "pieces", "pcs", "pc", "":
		// Keep the current stock as-is for gram/piece-based values.
	default:
		// Unknown unit: fall back to the current stock so the API still returns a value.
	}

	if stock.Product.RequiredGramPerUnit <= 0 {
		return 0
	}

	return int64(stockInGrams / stock.Product.RequiredGramPerUnit)
}

// GET /api/stock/summary
// Get updated stock summary for all products
func (h *StockManagementHandler) GetAllStocksSummary(c *fiber.Ctx) error {
	ctx, err := h.extractStockAccessContext(c)
	if err != nil {
		return stockAccessError(c, err)
	}

	limit := 100
	offset := 0
	if value := c.Query("limit"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed > 0 {
			limit = parsed
		}
	}
	if value := c.Query("offset"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed >= 0 {
			offset = parsed
		}
	}

	stocks, total, err := h.service.GetAllProductStockByCompany(
		ctx.CompanyID,
		ctx.UserType,
		offset,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	stocksResponse := make([]fiber.Map, 0)
	damagedResponse := make([]fiber.Map, 0)
	totalStockValue := 0.0
	totalSoldProductValue := 0.0
	totalDamagedValue := 0.0
	damagedCount := 0

	for _, stock := range stocks {
		if !h.shouldIncludeProduct(stock.ProductID) {
			continue
		}

		if stock.DamagedStock > 0 {
			damagedValue := stock.DamagedStock * stock.AverageCost
			totalDamagedValue += damagedValue
			damagedCount++
			damagedResponse = append(damagedResponse, fiber.Map{
				"product_id": stock.ProductID, "product_name": stock.ProductName,
				"sku": stock.SKU, "damaged_stock": stock.DamagedStock,
				"average_cost": stock.AverageCost, "damaged_value": damagedValue,
				"damage_reason": stock.DamageReason, "damaged_at": stock.DamagedAt,
				"damaged_by": stock.DamagedBy, "type": "product",
			})
		}

		stockValue := stock.CurrentStock * stock.AverageCost
		totalStockValue += stockValue
		totalSoldProductValue += stock.SoldStock * stock.AverageCost

		stocksResponse = append(stocksResponse, fiber.Map{
			"product_id": stock.ProductID, "product_name": stock.ProductName,
			"sku": stock.SKU, "current_stock": stock.CurrentStock,
			"purchased_total": stock.PurchasedStock, "sold_total": stock.SoldStock,
			"reserved_stock": stock.ReservedStock, "available_stock": stock.AvailableStock,
			"damaged_stock": stock.DamagedStock, "average_cost": stock.AverageCost,
			"stock_value": stockValue, "last_purchased": stock.LastPurchasedDate,
			"last_sold": stock.LastSoldDate, "type": "product",
		})
	}

	// Preserve the existing variant service. Aggregate all users in the same company.
	variantSeen := make(map[string]bool)
	if ctx.UserType == "superadmin" && ctx.CompanyID == 0 {
		log.Printf("[STOCK_SUMMARY] Superadmin global variant listing needs view_user_id to select a company")
	} else {
		users, usersErr := h.getCompanyUsers(ctx.CompanyID)
		if usersErr != nil {
			log.Printf("[STOCK_SUMMARY] Error fetching company users: %v", usersErr)
		} else {
			for _, user := range users {
				variantStocks, _, variantErr := h.variantStockMgmt.GetAllVariantStocksByUser(user.ID, 0, 10000)
				if variantErr != nil {
					log.Printf("[STOCK_SUMMARY] Error fetching variants for user %d: %v", user.ID, variantErr)
					continue
				}

				for _, stock := range variantStocks {
					if variantSeen[stock.ID] || !h.shouldIncludeProduct(stock.ProductID) {
						continue
					}
					variantSeen[stock.ID] = true

					if stock.DamagedStock > 0 {
						damagedValue := stock.DamagedStock * stock.AverageCost
						totalDamagedValue += damagedValue
						damagedCount++
						damagedResponse = append(damagedResponse, fiber.Map{
							"product_id": stock.ProductID, "product_name": stock.ProductName,
							"variant_sku": stock.VariantSKU, "variant_name": stock.VariantName,
							"damaged_stock": stock.DamagedStock, "average_cost": stock.AverageCost,
							"damaged_value": damagedValue, "damage_reason": stock.DamageReason,
							"damaged_at": stock.DamagedAt, "damaged_by": stock.DamagedBy,
							"type": "variant",
						})
					}

					stockValue := stock.CurrentStock * stock.AverageCost
					totalStockValue += stockValue
					totalSoldProductValue += stock.SoldStock * stock.AverageCost

					stocksResponse = append(stocksResponse, fiber.Map{
						"product_id": stock.ProductID, "product_name": stock.ProductName,
						"sku": stock.VariantSKU, "variant_name": stock.VariantName,
						"current_stock": stock.CurrentStock, "purchased_total": stock.PurchasedStock,
						"sold_total": stock.SoldStock, "reserved_stock": stock.ReservedStock,
						"available_stock": stock.AvailableStock, "damaged_stock": stock.DamagedStock,
						"average_cost": stock.AverageCost, "stock_value": stockValue,
						"last_purchased": stock.LastPurchasedDate, "last_sold": stock.LastSoldDate,
						"type": "variant",
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stocks": stocksResponse, "damaged_products": damagedResponse,
		"damaged_count": damagedCount, "total": total,
		"total_stock_value": totalStockValue,
		"total_sold_product_value": totalSoldProductValue,
		"total_damaged_value": totalDamagedValue,
	})
}

// GET /api/stock/product/:product_id/movements
// Get product movement history
func (h *StockManagementHandler) GetProductMovements(c *fiber.Ctx) error {
	ctx, err := h.extractStockAccessContext(c)
	if err != nil {
		return stockAccessError(c, err)
	}

	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}

	limit := 10
	offset := 0
	if value := c.Query("limit"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed > 0 { limit = parsed }
	}
	if value := c.Query("offset"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed >= 0 { offset = parsed }
	}

	movements, total, err := h.service.GetProductMovementHistoryByCompany(
		productID, ctx.CompanyID, ctx.UserType, offset, limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	movementsResponse := make([]fiber.Map, len(movements))
	for i, movement := range movements {
		movementsResponse[i] = fiber.Map{
			"id": movement.ID, "product_id": movement.ProductID,
			"movement_type": movement.MovementType, "quantity": movement.Quantity,
			"rate": movement.Rate, "amount": movement.Amount,
			"reference_type": movement.ReferenceType, "reference_id": movement.ReferenceID,
			"reference_number": movement.ReferenceNumber,
			"balance_before_qty": movement.BalanceBeforeQty,
			"balance_after_qty": movement.BalanceAfterQty,
			"notes": movement.Notes, "created_at": movement.CreatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"movements": movementsResponse, "total": total})
}

// GET /api/stock/debug/product/:product_id
// Get detailed stock breakdown and audit trail
func (h *StockManagementHandler) GetProductStockDebug(c *fiber.Ctx) error {
	ctx, err := h.extractStockAccessContext(c)
	if err != nil {
		return stockAccessError(c, err)
	}

	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID is required"})
	}

	stock, err := h.service.GetProductStockByCompany(productID, ctx.CompanyID, ctx.UserType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	movements, total, err := h.service.GetProductMovementHistoryByCompany(
		productID, ctx.CompanyID, ctx.UserType, 0, 1000,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	movementsBreakdown := make([]fiber.Map, 0)
	inboundTotal := 0.0
	outboundTotal := 0.0
	for _, movement := range movements {
		if movement.Quantity > 0 { inboundTotal += movement.Quantity } else { outboundTotal += movement.Quantity }
		movementsBreakdown = append(movementsBreakdown, fiber.Map{
			"id": movement.ID, "movement_type": movement.MovementType,
			"quantity": movement.Quantity, "rate": movement.Rate, "amount": movement.Amount,
			"reference_type": movement.ReferenceType, "reference_id": movement.ReferenceID,
			"reference_number": movement.ReferenceNumber,
			"balance_before_qty": movement.BalanceBeforeQty,
			"balance_after_qty": movement.BalanceAfterQty,
			"notes": movement.Notes, "created_at": movement.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product_id": productID, "current_stock": stock.CurrentStock,
		"purchased_total": stock.PurchasedStock, "sold_total": stock.SoldStock,
		"reserved_stock": stock.ReservedStock, "available_stock": stock.AvailableStock,
		"average_cost": stock.AverageCost, "inbound_total": inboundTotal,
		"outbound_total": outboundTotal, "net_movement": inboundTotal + outboundTotal,
		"expected_current": inboundTotal + outboundTotal,
		"movements_count": total, "movements": movementsBreakdown,
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
	ctx, err := h.extractStockAccessContext(c)
	if err != nil {
		return stockAccessError(c, err)
	}

	var input MarkDamagedInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "error": "Invalid request body",
		})
	}

	userID := stockUserIDString(ctx.UserID)

	if input.VariantSKU != nil && *input.VariantSKU != "" {
		// Validate variant's product belongs to the authenticated company before mutation.
		currentVariant, err := h.variantStockMgmt.GetVariantStockSummary(*input.VariantSKU)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "error": "Variant stock not found",
			})
		}

		if _, err := h.service.GetProductStockByCompany(
			currentVariant.ProductID, ctx.CompanyID, ctx.UserType,
		); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false, "error": "Variant does not belong to your company",
			})
		}

		if err := h.variantStockMgmt.MarkVariantAsDamaged(
			*input.VariantSKU, input.Quantity, input.Reason, userID,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false, "error": err.Error(),
			})
		}

		updatedStock, err := h.variantStockMgmt.GetVariantStockSummary(*input.VariantSKU)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false, "error": "Failed to fetch updated stock",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true, "message": "Variant marked as damaged successfully",
			"type": "variant", "variant_sku": updatedStock.VariantSKU,
			"damaged_stock": updatedStock.DamagedStock,
			"available_stock": updatedStock.AvailableStock,
			"damage_reason": updatedStock.DamageReason,
			"damaged_at": updatedStock.DamagedAt,
		})
	}

	if err := h.service.MarkProductAsDamagedByCompany(
		input.ProductID, input.Quantity, input.Reason, userID,
		ctx.CompanyID, ctx.UserType,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "error": err.Error(),
		})
	}

	updatedStock, err := h.service.GetProductStockByCompany(
		input.ProductID, ctx.CompanyID, ctx.UserType,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "error": "Failed to fetch updated stock",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true, "message": "Product marked as damaged successfully",
		"type": "product", "product_id": updatedStock.ProductID,
		"damaged_stock": updatedStock.DamagedStock,
		"available_stock": updatedStock.AvailableStock,
		"damage_reason": updatedStock.DamageReason,
		"damaged_at": updatedStock.DamagedAt,
	})
}

// GET /api/stock/damaged
// Get all damaged products and variants
func (h *StockManagementHandler) GetDamagedProducts(c *fiber.Ctx) error {
	ctx, err := h.extractStockAccessContext(c)
	if err != nil {
		return stockAccessError(c, err)
	}

	limit := 50
	offset := 0
	if value := c.Query("limit"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed > 0 { limit = parsed }
	}
	if value := c.Query("offset"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed >= 0 { offset = parsed }
	}

	products, productTotal, productErr := h.service.GetDamagedProductsByCompany(
		ctx.CompanyID, ctx.UserType, offset, limit,
	)
	if productErr != nil {
		log.Printf("[DAMAGE_TRACKING] Error fetching damaged products: %v", productErr)
	}

	variants := make([]models.VariantStock, 0)
	variantTotal := int64(0)
	variantSeen := make(map[string]bool)

	if ctx.UserType != "superadmin" || ctx.CompanyID > 0 {
		users, usersErr := h.getCompanyUsers(ctx.CompanyID)
		if usersErr != nil {
			log.Printf("[DAMAGE_TRACKING] Error fetching company users: %v", usersErr)
		} else {
			for _, user := range users {
				userVariants, _, variantErr := h.variantStockMgmt.GetDamagedVariantsByUser(user.ID, 0, 10000)
				if variantErr != nil {
					log.Printf("[DAMAGE_TRACKING] Error fetching damaged variants for user %d: %v", user.ID, variantErr)
					continue
				}
				for _, variant := range userVariants {
					if variantSeen[variant.ID] { continue }
					variantSeen[variant.ID] = true
					variants = append(variants, variant)
					variantTotal++
				}
			}
		}
	}

	damagedItems := make([]fiber.Map, 0)
	totalDamagedValue := 0.0

	for _, product := range products {
		damagedValue := product.DamagedStock * product.AverageCost
		totalDamagedValue += damagedValue
		damagedItems = append(damagedItems, fiber.Map{
			"type": "product", "product_id": product.ProductID,
			"product_name": product.ProductName, "sku": product.SKU,
			"damaged_stock": product.DamagedStock, "damage_reason": product.DamageReason,
			"damaged_at": product.DamagedAt, "damaged_by": product.DamagedBy,
			"average_cost": product.AverageCost, "damaged_value": damagedValue,
		})
	}

	for _, variant := range variants {
		damagedValue := variant.DamagedStock * variant.AverageCost
		totalDamagedValue += damagedValue
		damagedItems = append(damagedItems, fiber.Map{
			"type": "variant", "product_id": variant.ProductID,
			"product_name": variant.ProductName, "variant_sku": variant.VariantSKU,
			"variant_name": variant.VariantName, "damaged_stock": variant.DamagedStock,
			"damage_reason": variant.DamageReason, "damaged_at": variant.DamagedAt,
			"damaged_by": variant.DamagedBy, "average_cost": variant.AverageCost,
			"damaged_value": damagedValue,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true, "damaged_items": damagedItems,
		"damaged_products": productTotal, "damaged_variants": variantTotal,
		"total_damaged": productTotal + variantTotal,
		"total_damaged_value": totalDamagedValue,
		"limit": limit, "offset": offset,
	})
}