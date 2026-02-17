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

func (h *PackageHandler) CreatePackage(c *fiber.Ctx) error {
	var pkgInput input.CreatePackageInput

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
