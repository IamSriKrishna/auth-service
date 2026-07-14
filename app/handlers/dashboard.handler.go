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

func NewDashboardHandler(
	service services.DashboardService,
	userRepo repo.UserRepository,
) *DashboardHandler {
	return &DashboardHandler{
		service:  service,
		userRepo: userRepo,
	}
}

func (h *DashboardHandler) extractUserContext(
	c *fiber.Ctx,
) (
	userID uint,
	userType string,
	email string,
	companyID uint,
	err error,
) {
	userID = getUintLocal(c, "user_id")
	companyID = getUintLocal(c, "company_id")

	if value := c.Locals("user_type"); value != nil {
		userType, _ = value.(string)
	}

	if value := c.Locals("user_email"); value != nil {
		email, _ = value.(string)
	}

	if userID == 0 {
		return 0, "", "", 0, fmt.Errorf("invalid authenticated user")
	}

	if userType == "" {
		return 0, "", "", 0, fmt.Errorf("invalid authenticated user type")
	}

	if userType != "superadmin" && companyID == 0 {
		return 0, "", "", 0, fmt.Errorf("user is not assigned to a company")
	}

	viewUserIDString := c.Query("view_user_id")
	if viewUserIDString == "" {
		return userID, userType, email, companyID, nil
	}

	viewUserID64, parseErr := strconv.ParseUint(viewUserIDString, 10, 64)
	if parseErr != nil || viewUserID64 == 0 {
		return 0, "", "", 0, fmt.Errorf("invalid view_user_id parameter")
	}

	if h.userRepo == nil {
		return 0, "", "", 0, fmt.Errorf("user repository is unavailable")
	}

	viewUserID := uint(viewUserID64)

	if userType == "superadmin" {
		selectedUser, getErr := h.userRepo.GetByID(viewUserID)
		if getErr != nil || selectedUser == nil {
			return 0, "", "", 0, fmt.Errorf("user not found with id: %d", viewUserID)
		}

		userID = selectedUser.ID
		userType = string(selectedUser.UserType)

		if selectedUser.CompanyID == nil || *selectedUser.CompanyID == 0 {
			return 0, "", "", 0, fmt.Errorf(
				"selected user is not assigned to a company",
			)
		}

		companyID = *selectedUser.CompanyID

		if selectedUser.Email != nil {
			email = *selectedUser.Email
		}

		if companyID == 0 {
			return 0, "", "", 0, fmt.Errorf("selected user is not assigned to a company")
		}

		return userID, userType, email, companyID, nil
	}

	selectedUser, getErr := h.userRepo.GetByIDAndCompanyID(
		viewUserID,
		companyID,
	)
	if getErr != nil || selectedUser == nil {
		return 0, "", "", 0, fmt.Errorf("user not found in your company")
	}

	userID = selectedUser.ID
	userType = string(selectedUser.UserType)

	if selectedUser.Email != nil {
		email = *selectedUser.Email
	}

	return userID, userType, email, companyID, nil
}

func getUintLocal(c *fiber.Ctx, key string) uint {
	value := c.Locals(key)

	switch typedValue := value.(type) {
	case uint:
		return typedValue
	case uint64:
		return uint(typedValue)
	case int:
		if typedValue > 0 {
			return uint(typedValue)
		}
	case int64:
		if typedValue > 0 {
			return uint(typedValue)
		}
	case float64:
		if typedValue > 0 {
			return uint(typedValue)
		}
	case string:
		parsedValue, err := strconv.ParseUint(typedValue, 10, 64)
		if err == nil {
			return uint(parsedValue)
		}
	}

	return 0
}

func dashboardContextError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}

// GetDashboard returns company-scoped dashboard metrics.
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	userID, userType, email, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	metrics, err := h.service.GetDashboardMetricsWithUserContext(
		userID,
		userType,
		email,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(metrics)
}

// GetShipmentTracking returns company-scoped shipment tracking.
func (h *DashboardHandler) GetShipmentTracking(c *fiber.Ctx) error {
	_, userType, _, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	shipmentID := c.Params("shipment_id")
	if shipmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "shipment_id is required",
		})
	}

	limit, parseErr := strconv.Atoi(c.Query("limit", "10"))
	if parseErr != nil || limit < 1 {
		limit = 10
	}

	tracking, err := h.service.GetShipmentTracking(
		shipmentID,
		limit,
		companyID,
		userType,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(tracking)
}

// AddShipmentTracking adds tracking only to the authenticated company.
func (h *DashboardHandler) AddShipmentTracking(c *fiber.Ctx) error {
	_, userType, _, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	shipmentID := c.Params("shipment_id")
	if shipmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "shipment_id is required",
		})
	}

	var request struct {
		Status    string  `json:"status"`
		Location  string  `json:"location"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Notes     string  `json:"notes"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	err = h.service.AddShipmentTracking(
		shipmentID,
		request.Status,
		request.Location,
		request.Latitude,
		request.Longitude,
		request.Notes,
		companyID,
		userType,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Tracking record added successfully",
	})
}

// GetStockSummary returns company-scoped stock.
func (h *DashboardHandler) GetStockSummary(c *fiber.Ctx) error {
	userID, userType, _, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	stock, err := h.service.GetStockSummaryWithUserContext(
		userID,
		userType,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(stock)
}

// GetEntityTrends returns company-scoped trends.
func (h *DashboardHandler) GetEntityTrends(c *fiber.Ctx) error {
	_, userType, _, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	entityType := c.Params("entity_type")

	days, parseErr := strconv.Atoi(c.Query("days", "30"))
	if parseErr != nil || days < 1 {
		days = 30
	}

	trends, err := h.service.GetEntityTrends(
		entityType,
		days,
		companyID,
		userType,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(trends)
}

// GetActivitySummary returns company-scoped activity.
func (h *DashboardHandler) GetActivitySummary(c *fiber.Ctx) error {
	userID, userType, _, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	activity, err := h.service.GetActivitySummaryWithUserContext(
		userID,
		userType,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(activity)
}

// RefreshMetrics keeps the existing global refresh method.
// Protect this route using SuperAdminMiddleware.
func (h *DashboardHandler) RefreshMetrics(c *fiber.Ctx) error {
	if err := h.service.RefreshDashboardMetrics(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Dashboard metrics refreshed successfully",
	})
}

// GetDiagnosticReport returns company-scoped diagnostics.
func (h *DashboardHandler) GetDiagnosticReport(c *fiber.Ctx) error {
	_, userType, _, companyID, err := h.extractUserContext(c)
	if err != nil {
		return dashboardContextError(c, err)
	}

	report, err := h.service.GetDiagnosticReport(companyID, userType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate diagnostic report",
		})
	}

	return c.Status(fiber.StatusOK).JSON(report)
}

// GetPublicLiveStatus remains global because this route is public.
func (h *DashboardHandler) GetPublicLiveStatus(c *fiber.Ctx) error {
	metrics, err := h.service.GetDashboardMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

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
