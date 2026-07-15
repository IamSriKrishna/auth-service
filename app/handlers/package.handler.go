package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PackageHandler struct {
	service  services.PackageService
	validate *validator.Validate
}

func NewPackageHandler(service services.PackageService) *PackageHandler {
	return &PackageHandler{service: service, validate: validator.New()}
}

func packageContext(c *fiber.Ctx) (uint, string, error) {
	userID, err := packageLocalToUint(c.Locals("user_id"))
	if err != nil || userID == 0 {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "invalid authenticated user")
	}
	companyID, err := packageLocalToUint(c.Locals("company_id"))
	if err != nil || companyID == 0 {
		return 0, "", fiber.NewError(fiber.StatusForbidden, "invalid authenticated company")
	}
	return companyID, strconv.FormatUint(uint64(userID), 10), nil
}

func packagePagination(c *fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (h *PackageHandler) CreatePackage(c *fiber.Ctx) error {
	var req input.CreatePackageInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	companyID, userID, err := packageContext(c)
	if err != nil {
		return err
	}
	pkg, err := h.service.CreatePackageForCompany(&req, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Package created successfully", "data": pkg})
}

func (h *PackageHandler) GetPackage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Package ID is required"})
	}
	companyID, _, err := packageContext(c)
	if err != nil {
		return err
	}
	pkg, err := h.service.GetPackageForCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": pkg})
}

func (h *PackageHandler) GetAllPackages(c *fiber.Ctx) error {
	limit, offset := packagePagination(c)
	companyID, _, err := packageContext(c)
	if err != nil {
		return err
	}
	packages, total, err := h.service.GetAllPackagesForCompany(companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": packages, "total": total})
}

func (h *PackageHandler) GetPackagesByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Params("customer_id"), 10, 32)
	if err != nil || customerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid customer ID"})
	}
	limit, offset := packagePagination(c)
	companyID, _, err := packageContext(c)
	if err != nil {
		return err
	}
	packages, total, err := h.service.GetPackagesByCustomerForCompany(uint(customerID), companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": packages, "total": total})
}

func (h *PackageHandler) GetPackagesBySalesOrder(c *fiber.Ctx) error {
	salesOrderID := c.Params("sales_order_id")
	if salesOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Sales order ID is required"})
	}
	limit, offset := packagePagination(c)
	companyID, _, err := packageContext(c)
	if err != nil {
		return err
	}
	packages, total, err := h.service.GetPackagesBySalesOrderForCompany(salesOrderID, companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": packages, "total": total})
}

func (h *PackageHandler) GetPackagesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Package status is required"})
	}
	limit, offset := packagePagination(c)
	companyID, _, err := packageContext(c)
	if err != nil {
		return err
	}
	packages, total, err := h.service.GetPackagesByStatusForCompany(status, companyID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": packages, "total": total})
}

func (h *PackageHandler) UpdatePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Package ID is required"})
	}
	var req input.UpdatePackageInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	companyID, userID, err := packageContext(c)
	if err != nil {
		return err
	}
	pkg, err := h.service.UpdatePackageForCompany(id, &req, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Package updated successfully", "data": pkg})
}

func (h *PackageHandler) UpdatePackageStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Package ID is required"})
	}
	var req input.UpdatePackageStatusInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	companyID, userID, err := packageContext(c)
	if err != nil {
		return err
	}
	pkg, err := h.service.UpdatePackageStatusForCompany(id, req.Status, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Package status updated successfully", "data": pkg})
}

func (h *PackageHandler) DeletePackage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Package ID is required"})
	}
	companyID, _, err := packageContext(c)
	if err != nil {
		return err
	}
	if err := h.service.DeletePackageForCompany(id, companyID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Package deleted successfully"})
}

func packageLocalToUint(value interface{}) (uint, error) {
	if value == nil {
		return 0, fmt.Errorf("value is missing")
	}
	switch v := value.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case int64:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case float64:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case json.Number:
		n, err := strconv.ParseUint(v.String(), 10, 64)
		return uint(n), err
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		return uint(n), err
	default:
		n, err := strconv.ParseUint(fmt.Sprintf("%v", v), 10, 64)
		return uint(n), err
	}
}
