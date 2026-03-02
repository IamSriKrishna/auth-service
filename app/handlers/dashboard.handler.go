package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service services.DashboardService
}

func NewDashboardHandler(service services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// GetDashboard returns all dashboard metrics
// @Summary Get Dashboard Metrics
// @Description Retrieve all dashboard metrics including customers, vendors, items, shipments, invoices, orders
// @Tags Dashboard
// @Produce json
// @Success 200 {object} output.DashboardMetricsOutput
// @Router /dashboard [get]
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	metrics, err := h.service.GetDashboardMetrics()
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
	stock, err := h.service.GetStockSummary()
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
	activity, err := h.service.GetActivitySummary()
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
