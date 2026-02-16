package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type VendorHandler struct {
	service  services.VendorService
	validate *validator.Validate
}

func 	NewVendorHandler(service services.VendorService) *VendorHandler {
	return &VendorHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateVendor creates a new vendor
// @Summary Create a new vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Param vendor body input.CreateVendorInput true "Vendor data"
// @Success 201 {object} input.VendorOutput
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /vendors [post]
func (h *VendorHandler) CreateVendor(c *fiber.Ctx) error {
	var input input.CreateVendorInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.CreateVendor(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Vendor created successfully",
		"data":    vendor,
	})
}

// UpdateVendor updates an existing vendor
// @Summary Update a vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Param id path int true "Vendor ID"
// @Param vendor body input.UpdateVendorInput true "Vendor data"
// @Success 200 {object} input.VendorOutput
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /vendors/{id} [put]
func (h *VendorHandler) UpdateVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	var input input.UpdateVendorInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.UpdateVendor(uint(id), &input)
	if err != nil {
		if err.Error() == "vendor not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor updated successfully",
		"data":    vendor,
	})
}

// GetVendor retrieves a vendor by ID
// @Summary Get a vendor by ID
// @Tags vendors
// @Produce json
// @Param id path int true "Vendor ID"
// @Success 200 {object} input.VendorOutput
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vendors/{id} [get]
func (h *VendorHandler) GetVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	vendor, err := h.service.GetVendorByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vendor not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    vendor,
	})
}

// GetAllVendors retrieves all vendors with pagination
// @Summary Get all vendors
// @Tags vendors
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /vendors [get]
func (h *VendorHandler) GetAllVendors(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	vendors, total, err := h.service.GetAllVendors(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    vendors,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// DeleteVendor deletes a vendor by ID
// @Summary Delete a vendor
// @Tags vendors
// @Produce json
// @Param id path int true "Vendor ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /vendors/{id} [delete]
func (h *VendorHandler) DeleteVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	if err := h.service.DeleteVendor(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor deleted successfully",
	})
}
