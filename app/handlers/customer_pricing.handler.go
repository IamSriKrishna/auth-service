package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

// CustomerPricingHandler handles HTTP requests for customer pricing
type CustomerPricingHandler struct {
	service services.CustomerPricingService
}

// NewCustomerPricingHandler creates a new instance of customer pricing handler
func NewCustomerPricingHandler(service services.CustomerPricingService) *CustomerPricingHandler {
	return &CustomerPricingHandler{
		service: service,
	}
}

// CreateCustomerPricing creates new customer pricing records from line items
// @Summary Create customer pricing
// @Description Create multiple customer pricing records from line items
// @Tags Customer Pricing
// @Accept json
// @Produce json
// @Param request body input.CreateCustomerPricingDTO true "Pricing details with line items"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing [post]
func (h *CustomerPricingHandler) CreateCustomerPricing(c *fiber.Ctx) error {
	var req input.CreateCustomerPricingDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
			"details": err.Error(),
		})
	}

	pricings, err := h.service.CreatePricing(req.CustomerID, req.LineItems)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    pricings,
		"count":   len(pricings),
		"message": "customer pricing created successfully",
	})
}

// UpdateCustomerPricing updates an existing customer pricing record
// @Summary Update customer pricing
// @Description Update customer pricing details
// @Tags Customer Pricing
// @Accept json
// @Produce json
// @Param id path string true "Pricing ID"
// @Param request body input.UpdateCustomerPricingDTO true "Updated pricing details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing/{id} [put]
func (h *CustomerPricingHandler) UpdateCustomerPricing(c *fiber.Ctx) error {
	id := c.Params("id")
	var req input.UpdateCustomerPricingDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
			"details": err.Error(),
		})
	}

	pricing, err := h.service.UpdatePricing(id, req.Rate, req.Account, req.Description, req.IsActive)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pricing,
	})
}

// DeleteCustomerPricing deletes a customer pricing record
// @Summary Delete customer pricing
// @Description Delete a customer pricing record
// @Tags Customer Pricing
// @Param id path string true "Pricing ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing/{id} [delete]
func (h *CustomerPricingHandler) DeleteCustomerPricing(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeletePricing(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "customer pricing deleted successfully",
	})
}

// GetCustomerPricingByID retrieves a pricing record by ID
// @Summary Get customer pricing by ID
// @Description Retrieve a specific customer pricing record
// @Tags Customer Pricing
// @Produce json
// @Param id path string true "Pricing ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing/{id} [get]
func (h *CustomerPricingHandler) GetCustomerPricingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	pricing, err := h.service.GetPricingByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "pricing not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pricing,
	})
}

// GetPricingByCustomer retrieves all pricing for a customer
// @Summary Get pricing by customer
// @Description Retrieve all pricing records for a specific customer
// @Tags Customer Pricing
// @Produce json
// @Param customer_id query uint true "Customer ID"
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing/customer [get]
func (h *CustomerPricingHandler) GetPricingByCustomer(c *fiber.Ctx) error {
	customerIDStr := c.Query("customer_id")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	customerID, err := strconv.ParseUint(customerIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid customer ID",
		})
	}

	pricings, total, err := h.service.GetPricingByCustomer(uint(customerID), offset, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"pricings":    pricings,
			"total_count": total,
		},
	})
}

// GetAllCustomerPricing retrieves all customer pricing records
// @Summary Get all customer pricing
// @Description Retrieve all customer pricing records with pagination
// @Tags Customer Pricing
// @Produce json
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing [get]
func (h *CustomerPricingHandler) GetAllCustomerPricing(c *fiber.Ctx) error {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pricings, total, err := h.service.GetAllPricing(offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"pricings":    pricings,
			"total_count": total,
		},
	})
}

// GetActivePricingByCustomer retrieves all active pricing for a customer
// @Summary Get active pricing by customer
// @Description Retrieve all active pricing records for a customer
// @Tags Customer Pricing
// @Produce json
// @Param customer_id query uint true "Customer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer-pricing/customer/active [get]
func (h *CustomerPricingHandler) GetActivePricingByCustomer(c *fiber.Ctx) error {
	customerIDStr := c.Query("customer_id")
	customerID, err := strconv.ParseUint(customerIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid customer ID",
		})
	}

	pricings, err := h.service.GetActivePricingByCustomer(uint(customerID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"pricings": pricings,
		},
	})
}

// SetEffectiveDateRange sets the effective date range for pricing
// @Summary Set effective date range
// @Description Set the effective date range for a customer pricing record
// @Tags Customer Pricing
// @Accept json
// @Produce json
// @Param id path string true "Pricing ID"
// @Param request body input.SetDateRangeDTO true "Date range details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/customer-pricing/{id}/date-range [put]
func (h *CustomerPricingHandler) SetEffectiveDateRange(c *fiber.Ctx) error {
	id := c.Params("id")
	var req input.SetDateRangeDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
			"details": err.Error(),
		})
	}

	if err := h.service.SetEffectiveDateRange(id, req.EffectiveFrom, req.EffectiveTo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "date range updated successfully",
	})
}
