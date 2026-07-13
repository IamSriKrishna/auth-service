package handlers

import (
	"fmt"
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

func NewVendorHandler(service services.VendorService) *VendorHandler {
	return &VendorHandler{
		service:  service,
		validate: validator.New(),
	}
}

func vendorLocalUint(c *fiber.Ctx, key string) uint {
	value := c.Locals(key)

	switch typed := value.(type) {
	case uint:
		return typed
	case uint64:
		return uint(typed)
	case int:
		if typed > 0 {
			return uint(typed)
		}
	case int64:
		if typed > 0 {
			return uint(typed)
		}
	case float64:
		if typed > 0 {
			return uint(typed)
		}
	case string:
		parsed, err := strconv.ParseUint(typed, 10, 64)
		if err == nil {
			return uint(parsed)
		}
	}

	return 0
}

func vendorAuthContext(c *fiber.Ctx) (uint, uint, error) {
	userID := vendorLocalUint(c, "user_id")
	companyID := vendorLocalUint(c, "company_id")

	if userID == 0 {
		return 0, 0, fmt.Errorf("invalid authenticated user")
	}

	if companyID == 0 {
		return 0, 0, fmt.Errorf("user is not assigned to a company")
	}

	return userID, companyID, nil
}

func normalizeVendorPagination(c *fiber.Ctx) (int, int) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return page, limit
}

func (h *VendorHandler) CreateVendor(c *fiber.Ctx) error {
	var req input.CreateVendorInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.CreateVendorForUser(
		userID,
		companyID,
		&req,
	)
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

func (h *VendorHandler) UpdateVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	var req input.UpdateVendorInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.UpdateVendorForCompany(
		uint(id),
		companyID,
		&req,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor updated successfully",
		"data":    vendor,
	})
}

func (h *VendorHandler) GetVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.GetVendorByIDAndCompany(
		uint(id),
		companyID,
	)
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

func (h *VendorHandler) GetAllVendors(c *fiber.Ctx) error {
	page, limit := normalizeVendorPagination(c)

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendors, total, err := h.service.GetVendorsByCompany(
		companyID,
		page,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

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

func (h *VendorHandler) DeleteVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.DeleteVendorForCompany(
		uint(id),
		companyID,
	); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor deleted successfully",
	})
}

// Kept for route compatibility.
// It now returns all vendors belonging to the authenticated company.
func (h *VendorHandler) GetUserVendors(c *fiber.Ctx) error {
	page, limit := normalizeVendorPagination(c)

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendors, total, err := h.service.GetVendorsByCompany(
		companyID,
		page,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

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

func (h *VendorHandler) GetUserVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.GetVendorByIDAndCompany(
		uint(id),
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vendor not found or unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    vendor,
	})
}

func (h *VendorHandler) UpdateUserVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	var req input.UpdateVendorInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	vendor, err := h.service.UpdateVendorForCompany(
		uint(id),
		companyID,
		&req,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor updated successfully",
		"data":    vendor,
	})
}

func (h *VendorHandler) DeleteUserVendor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid vendor ID",
		})
	}

	_, companyID, err := vendorAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.DeleteVendorForCompany(
		uint(id),
		companyID,
	); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor deleted successfully",
	})
}