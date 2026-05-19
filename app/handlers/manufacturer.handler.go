package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type ManufacturerHandler struct {
	service services.ManufacturerService
}

func NewManufacturerHandler(service services.ManufacturerService) *ManufacturerHandler {
	return &ManufacturerHandler{service: service}
}

// CreateManufacturer creates a new manufacturer with product group and employees
func (h *ManufacturerHandler) CreateManufacturer(c *fiber.Ctx) error {
	var req input.CreateManufacturerInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Extract user info from context
	userID := uint(0)
	if id := c.Locals("user_id"); id != nil {
		switch v := id.(type) {
		case uint:
			userID = v
		case int:
			userID = uint(v)
		case float64:
			userID = uint(v)
		}
	}

	companyID := uint(0)
	if id := c.Locals("company_id"); id != nil {
		switch v := id.(type) {
		case uint:
			companyID = v
		case int:
			companyID = uint(v)
		case float64:
			companyID = uint(v)
		}
	}

	if userID == 0 || companyID == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "user or company information missing")
	}

	manufacturer, err := h.service.Create(&req, userID, companyID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    manufacturer,
	})
}

// UpdateManufacturer updates an existing manufacturer
func (h *ManufacturerHandler) UpdateManufacturer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid manufacturer id")
	}

	var req input.UpdateManufacturerInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	manufacturer, err := h.service.Update(id, &req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturer,
	})
}

// GetManufacturerByID retrieves a manufacturer by ID
func (h *ManufacturerHandler) GetManufacturerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid manufacturer id")
	}

	manufacturer, err := h.service.GetByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturer,
	})
}

// GetAllManufacturers retrieves all manufacturers with pagination
func (h *ManufacturerHandler) GetAllManufacturers(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	manufacturers, err := h.service.GetAll(limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturers,
	})
}

// GetManufacturersByProductGroup retrieves manufacturers for a specific product group
func (h *ManufacturerHandler) GetManufacturersByProductGroup(c *fiber.Ctx) error {
	productGroupID := c.Params("product_group_id")
	if productGroupID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "product group id is required")
	}

	manufacturers, err := h.service.GetByProductGroupID(productGroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturers,
	})
}

// DeleteManufacturer deletes a manufacturer
func (h *ManufacturerHandler) DeleteManufacturer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid manufacturer id")
	}

	err := h.service.Delete(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "manufacturer deleted successfully",
	})
}
