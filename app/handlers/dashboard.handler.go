package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service  services.DashboardService
	userRepo repo.UserRepository
}

func NewDashboardHandler(service services.DashboardService, userRepo repo.UserRepository) *DashboardHandler {
	return &DashboardHandler{
		service:  service,
		userRepo: userRepo,
	}
}

// extractUserContext extracts and validates user context from JWT and query parameters
func (h *DashboardHandler) extractUserContext(c *fiber.Ctx) (userID uint, userType string, email string, companyID uint, err error) {
	authenticatedUserID := uint(0)
	authenticatedUserType := "superadmin"
	authenticatedEmail := ""
	authenticatedCompanyID := uint(0)

	// Extract authenticated user information from context (from JWT)
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

	if e := c.Locals("user_email"); e != nil {
		if emailStr, ok := e.(string); ok {
			authenticatedEmail = emailStr
		}
	}

	if cid := c.Locals("company_id"); cid != nil {
		switch v := cid.(type) {
		case uint:
			authenticatedCompanyID = v
		case int:
			authenticatedCompanyID = uint(v)
		case float64:
			authenticatedCompanyID = uint(v)
		}
	}

	userID = authenticatedUserID
	userType = authenticatedUserType
	email = authenticatedEmail
	companyID = authenticatedCompanyID

	// Check if view_user_id is provided in query parameters
	viewUserIDStr := c.Query("view_user_id")
	if viewUserIDStr != "" {
		var viewUserID uint64
		_, parseErr := fmt.Sscanf(viewUserIDStr, "%d", &viewUserID)
		if parseErr != nil || viewUserID <= 0 {
			return 0, "", "", 0, fmt.Errorf("invalid view_user_id parameter")
		}

		// Permission check: Only superadmin can view any user's dashboard, others can only view their own
		if authenticatedUserType != "superadmin" && uint(viewUserID) != authenticatedUserID {
			return 0, "", "", 0, fmt.Errorf("unauthorized: cannot view another user's dashboard")
		}

		// Validate that the user exists
		if h.userRepo != nil {
			user, getErr := h.userRepo.GetByID(uint(viewUserID))
			if getErr != nil || user == nil {
				return 0, "", "", 0, fmt.Errorf("user not found with id: %d", viewUserID)
			}
		}

		// Set the userID to the requested view_user_id
		userID = uint(viewUserID)
		// If viewing another user and authenticated user is superadmin, treat context as admin
		if uint(viewUserID) != authenticatedUserID && authenticatedUserType == "superadmin" {
			userType = "admin"
		}
	}

	return userID, userType, email, companyID, nil
}

// GetDashboard returns all dashboard metrics
// @Summary Get Dashboard Metrics
// @Description Retrieve all dashboard metrics including customers, vendors, items, shipments, invoices, orders
// @Tags Dashboard
// @Produce json
// @Success 200 {object} output.DashboardMetricsOutput
// @Router /dashboard [get]
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	userID, userType, email, companyID, err := h.extractUserContext(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	metrics, err := h.service.GetDashboardMetricsWithUserContext(userID, userType, email, companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(metrics)
}

// GetShipmentTracking returns shipment tracking details
// @Summary Get Shipment Tracking
// @Description Retrieve tracking details for a specific shipment
// @Tags Dashboard
// @Produce json
// @Param shipment_id path string true "Shipment ID"
// @Param limit query int false "Limit (default: 10)"
// @Success 200 {object} output.ShipmentTrackingListOutput
// @Router /dashboard/shipment/{shipment_id}/tracking [get]
func (h *DashboardHandler) GetShipmentTracking(c *fiber.Ctx) error {
	shipmentID := c.Params("shipment_id")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	tracking, err := h.service.GetShipmentTracking(shipmentID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(tracking)
}

// AddShipmentTracking adds a new shipment tracking record
// @Summary Add Shipment Tracking
// @Description Add a new tracking record for a shipment
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param shipment_id path string true "Shipment ID"
// @Param request body map[string]interface{} true "Tracking details"
// @Success 201 {object} map[string]string
// @Router /dashboard/shipment/{shipment_id}/tracking [post]
func (h *DashboardHandler) AddShipmentTracking(c *fiber.Ctx) error {
	shipmentID := c.Params("shipment_id")

	var req struct {
		Status    string  `json:"status"`
		Location  string  `json:"location"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Notes     string  `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.service.AddShipmentTracking(shipmentID, req.Status, req.Location, req.Latitude, req.Longitude, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tracking record added successfully",
	})
}

// GetStockSummary returns stock summary
// @Summary Get Stock Summary
// @Description Retrieve stock information including in-stock, low stock, and out of stock items
// @Tags Dashboard
// @Produce json
// @Success 200 {object} output.StockListOutput
// @Router /dashboard/stock [get]
func (h *DashboardHandler) GetStockSummary(c *fiber.Ctx) error {
	userID, userType, _, _, err := h.extractUserContext(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	stock, err := h.service.GetStockSummaryWithUserContext(userID, userType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(stock)
}

// GetEntityTrends returns trend data for an entity
// @Summary Get Entity Trends
// @Description Retrieve historical trend data for a specific entity
// @Tags Dashboard
// @Produce json
// @Param entity_type path string true "Entity type (customer, vendor, item, etc.)"
// @Param days query int false "Number of days (default: 30)"
// @Success 200 {object} output.EntityTrendOutput
// @Router /dashboard/trends/{entity_type} [get]
func (h *DashboardHandler) GetEntityTrends(c *fiber.Ctx) error {
	entityType := c.Params("entity_type")
	days, _ := strconv.Atoi(c.Query("days", "30"))

	trends, err := h.service.GetEntityTrends(entityType, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(trends)
}

// GetActivitySummary returns today's activity summary
// @Summary Get Activity Summary
// @Description Retrieve today's activity including created items, shipments, orders
// @Tags Dashboard
// @Produce json
// @Success 200 {object} output.ActivitySummaryOutput
// @Router /dashboard/activity [get]
func (h *DashboardHandler) GetActivitySummary(c *fiber.Ctx) error {
	userID, userType, _, _, err := h.extractUserContext(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	activity, err := h.service.GetActivitySummaryWithUserContext(userID, userType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(activity)
}

// RefreshMetrics triggers a refresh of dashboard metrics
// @Summary Refresh Dashboard Metrics
// @Description Manually trigger a refresh of all dashboard metrics
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]string
// @Router /dashboard/refresh [post]
func (h *DashboardHandler) RefreshMetrics(c *fiber.Ctx) error {
	err := h.service.RefreshDashboardMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Dashboard metrics refreshed successfully",
	})
}

// GetDiagnosticReport handler - GET /dashboard/diagnose
func (h *DashboardHandler) GetDiagnosticReport(c *fiber.Ctx) error {
	report, err := h.service.GetDiagnosticReport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate diagnostic report",
		})
	}

	return c.Status(fiber.StatusOK).JSON(report)
}

// GetPublicLiveStatus returns live status of all data (customers, vendors, products, stock) - PUBLIC endpoint
// @Summary Get Public Live Status
// @Description Retrieve live status about all customers, vendors, products, and stock (public access, no authentication required)
// @Tags Dashboard
// @Produce json
// @Success 200 {object} output.PublicLiveStatusOutput
// @Router /public/live-status [get]
func (h *DashboardHandler) GetPublicLiveStatus(c *fiber.Ctx) error {
	metrics, err := h.service.GetDashboardMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return simplified public status focused on customers, vendors, products, and stock
	publicStatus := output.PublicLiveStatusOutput{
		Customers: output.PublicCustomerStatus{
			Total:  metrics.CustomerMetrics.Total,
			Active: metrics.CustomerMetrics.Active,
		},
		Vendors: output.PublicVendorStatus{
			Total:  metrics.VendorMetrics.Total,
			Active: metrics.VendorMetrics.Active,
		},
		Products: output.PublicProductStatus{
			Total:           metrics.ItemMetrics.Total,
			TotalStock:      metrics.ItemMetrics.TotalStock,
			LowStockItems:   metrics.ItemMetrics.LowStockItems,
			OutOfStockItems: metrics.ItemMetrics.OutOfStockItem,
		},
		Stock: output.PublicStockStatus{
			TotalItems:    metrics.ItemMetrics.Total,
			TotalQuantity: metrics.ItemMetrics.TotalStock,
			LowStock:      metrics.ItemMetrics.LowStockItems,
			OutOfStock:    metrics.ItemMetrics.OutOfStockItem,
		},
		LastUpdatedAt: metrics.LastUpdatedAt,
		GeneratedAt:   metrics.GeneratedAt,
	}

	return c.Status(fiber.StatusOK).JSON(publicStatus)
}
