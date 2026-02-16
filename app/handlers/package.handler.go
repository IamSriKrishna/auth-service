package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PackageHandler struct {
	service services.PackageService
}

func NewPackageHandler(service services.PackageService) *PackageHandler {
	return &PackageHandler{service: service}
}

// CreatePackage creates a new package
// @Summary Create a new package
// @Description Create a new package with items for shipment
// @Tags Package
// @Accept json
// @Produce json
// @Param package body input.CreatePackageInput true "Package input"
// @Success 201 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /packages [post]
func (h *PackageHandler) CreatePackage(c *fiber.Ctx) error {
	var pkgInput input.CreatePackageInput

	if err := c.BodyParser(&pkgInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(pkgInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	pkg, err := h.service.CreatePackage(&pkgInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Package created successfully",
		"data":    pkg,
	})
}

// GetPackage retrieves a package by ID
// @Summary Get package by ID
// @Description Get a specific package by its ID
// @Tags Package
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /packages/{id} [get]
func (h *PackageHandler) GetPackage(c *fiber.Ctx) error {
	id := c.Params("id")

	pkg, err := h.service.GetPackage(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Package not found",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    pkg,
	})
}

// GetAllPackages retrieves all packages with pagination
// @Summary Get all packages
// @Description Get all packages with pagination support
// @Tags Package
// @Produce json
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /packages [get]
func (h *PackageHandler) GetAllPackages(c *fiber.Ctx) error {
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

	packages, total, err := h.service.GetAllPackages(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    packages,
		"total":   total,
	})
}

// GetPackagesByCustomer retrieves packages for a specific customer
// @Summary Get packages by customer
// @Description Get all packages for a specific customer
// @Tags Package
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /packages/customer/{customer_id} [get]
func (h *PackageHandler) GetPackagesByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Params("customer_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid customer ID",
			"success": false,
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

	packages, total, err := h.service.GetPackagesByCustomer(uint(customerID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    packages,
		"total":   total,
	})
}

// GetPackagesBySalesOrder retrieves packages for a specific sales order
// @Summary Get packages by sales order
// @Description Get all packages for a specific sales order
// @Tags Package
// @Produce json
// @Param sales_order_id path string true "Sales Order ID"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /packages/sales-order/{sales_order_id} [get]
func (h *PackageHandler) GetPackagesBySalesOrder(c *fiber.Ctx) error {
	salesOrderID := c.Params("sales_order_id")

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

	packages, total, err := h.service.GetPackagesBySalesOrder(salesOrderID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    packages,
		"total":   total,
	})
}

// GetPackagesByStatus retrieves packages by status
// @Summary Get packages by status
// @Description Get all packages with a specific status
// @Tags Package
// @Produce json
// @Param status query string true "Status (created, packed, shipped, delivered, cancelled)"
// @Param limit query int false "Limit (default: 10)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /packages/status/{status} [get]
func (h *PackageHandler) GetPackagesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")

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

	packages, total, err := h.service.GetPackagesByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    packages,
		"total":   total,
	})
}

// UpdatePackage updates a package
// @Summary Update a package
// @Description Update package details
// @Tags Package
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Param package body input.UpdatePackageInput true "Package input"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /packages/{id} [put]
func (h *PackageHandler) UpdatePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	var pkgInput input.UpdatePackageInput

	if err := c.BodyParser(&pkgInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(pkgInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	pkg, err := h.service.UpdatePackage(id, &pkgInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Package updated successfully",
		"data":    pkg,
	})
}

// UpdatePackageStatus updates package status
// @Summary Update package status
// @Description Update the status of a package
// @Tags Package
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Param status body input.UpdatePackageStatusInput true "Status input"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /packages/{id}/status [patch]
func (h *PackageHandler) UpdatePackageStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var statusInput input.UpdatePackageStatusInput

	if err := c.BodyParser(&statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	pkg, err := h.service.UpdatePackageStatus(id, statusInput.Status, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Package status updated successfully",
		"data":    pkg,
	})
}

// DeletePackage deletes a package
// @Summary Delete a package
// @Description Delete a package by ID
// @Tags Package
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /packages/{id} [delete]
func (h *PackageHandler) DeletePackage(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeletePackage(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Package deleted successfully",
	})
}
