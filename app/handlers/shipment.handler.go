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

type ShipmentHandler struct {
	service  services.ShipmentService
	validate *validator.Validate
}

func NewShipmentHandler(service services.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{
		service:  service,
		validate: validator.New(),
	}
}

func shipmentContext(c *fiber.Ctx) (uint, string, error) {
	userID, err := shipmentLocalToUint(c.Locals("user_id"))
	if err != nil || userID == 0 {
		return 0, "", fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid authenticated user",
		)
	}

	companyID, err := shipmentLocalToUint(c.Locals("company_id"))
	if err != nil || companyID == 0 {
		return 0, "", fiber.NewError(
			fiber.StatusForbidden,
			"invalid authenticated company",
		)
	}

	return companyID, strconv.FormatUint(uint64(userID), 10), nil
}

func shipmentPagination(c *fiber.Ctx) (int, int) {
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

func (h *ShipmentHandler) CreateShipment(c *fiber.Ctx) error {
	var req input.CreateShipmentInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	companyID, userID, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipment, err := h.service.CreateShipmentForCompany(
		&req,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Shipment created successfully",
		"data":    shipment,
	})
}

func (h *ShipmentHandler) GetShipment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Shipment ID is required",
		})
	}

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipment, err := h.service.GetShipmentForCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipment,
	})
}

func (h *ShipmentHandler) GetAllShipments(c *fiber.Ctx) error {
	limit, offset := shipmentPagination(c)

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipments, total, err := h.service.GetAllShipmentsForCompany(
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(
		c.Params("customer_id"),
		10,
		32,
	)
	if err != nil || customerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid customer ID",
		})
	}

	limit, offset := shipmentPagination(c)

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipments, total, err :=
		h.service.GetShipmentsByCustomerForCompany(
			uint(customerID),
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsByPackage(c *fiber.Ctx) error {
	packageID := c.Params("package_id")
	if packageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Package ID is required",
		})
	}

	limit, offset := shipmentPagination(c)

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipments, total, err :=
		h.service.GetShipmentsByPackageForCompany(
			packageID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsBySalesOrder(c *fiber.Ctx) error {
	salesOrderID := c.Params("sales_order_id")
	if salesOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Sales order ID is required",
		})
	}

	limit, offset := shipmentPagination(c)

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipments, total, err :=
		h.service.GetShipmentsBySalesOrderForCompany(
			salesOrderID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) GetShipmentsByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Shipment status is required",
		})
	}

	limit, offset := shipmentPagination(c)

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipments, total, err :=
		h.service.GetShipmentsByStatusForCompany(
			status,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    shipments,
		"total":   total,
	})
}

func (h *ShipmentHandler) UpdateShipment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Shipment ID is required",
		})
	}

	var req input.UpdateShipmentInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	companyID, userID, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipment, err := h.service.UpdateShipmentForCompany(
		id,
		&req,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Shipment updated successfully",
		"data":    shipment,
	})
}

func (h *ShipmentHandler) UpdateShipmentStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Shipment ID is required",
		})
	}

	var req input.UpdateShipmentStatusInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	companyID, userID, err := shipmentContext(c)
	if err != nil {
		return err
	}

	shipment, err :=
		h.service.UpdateShipmentStatusForCompany(
			id,
			req.Status,
			userID,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Shipment status updated successfully",
		"data":    shipment,
	})
}

func (h *ShipmentHandler) DeleteShipment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Shipment ID is required",
		})
	}

	companyID, _, err := shipmentContext(c)
	if err != nil {
		return err
	}

	if err := h.service.DeleteShipmentForCompany(
		id,
		companyID,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Shipment deleted successfully",
	})
}

// Retained route/function name.
// It now creates the actual company-scoped shipment instead of returning
// a temporary response that is not saved in the database.
func (h *ShipmentHandler) CreateShipmentWithVariants(
	c *fiber.Ctx,
) error {
	return h.CreateShipment(c)
}

func shipmentLocalToUint(value interface{}) (uint, error) {
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
	case int8:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case int16:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case int32:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case int64:
		if v <= 0 {
			return 0, fmt.Errorf("invalid value")
		}
		return uint(v), nil
	case float32:
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
		n, err := strconv.ParseUint(
			fmt.Sprintf("%v", v),
			10,
			64,
		)
		return uint(n), err
	}
}
