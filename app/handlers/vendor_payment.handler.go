package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type VendorPaymentHandler struct {
	service  services.VendorPaymentService
	validate *validator.Validate
}

func NewVendorPaymentHandler(service services.VendorPaymentService) *VendorPaymentHandler {
	return &VendorPaymentHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *VendorPaymentHandler) CreateVendorPayment(c *fiber.Ctx) error {
	var paymentInput input.CreateVendorPaymentInput

	if err := c.BodyParser(&paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	userName := ""
	if email := c.Locals("user_email"); email != nil {
		userName = fmt.Sprintf("%v", email)
	}

	companyID := uint(0)
	if cid := c.Locals("company_id"); cid != nil {
		if id, ok := cid.(uint); ok {
			companyID = id
		} else if id, ok := cid.(float64); ok {
			companyID = uint(id)
		}
	}

	companyName := ""
	if cname := c.Locals("company_name"); cname != nil {
		companyName = fmt.Sprintf("%v", cname)
	}

	payment, err := h.service.CreateVendorPayment(&paymentInput, userID, userName, companyName, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Vendor payment created successfully",
		"data":    payment,
	})
}

func (h *VendorPaymentHandler) GetVendorPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	payment, err := h.service.GetVendorPayment(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    payment,
	})
}

func (h *VendorPaymentHandler) GetAllVendorPayments(c *fiber.Ctx) error {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	payments, total, err := h.service.GetAllVendorPayments(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"vendor_payments": payments,
			"total":           total,
		},
	})
}

func (h *VendorPaymentHandler) GetVendorPaymentsByPurchaseOrder(c *fiber.Ctx) error {
	poID := c.Params("purchaseOrderId")
	if poID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Purchase order ID is required",
		})
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	payments, total, err := h.service.GetVendorPaymentsByPurchaseOrder(poID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"vendor_payments": payments,
			"total":           total,
		},
	})
}

func (h *VendorPaymentHandler) GetVendorPaymentsByVendor(c *fiber.Ctx) error {
	vendorID, err := strconv.ParseUint(c.Params("vendorId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vendor ID",
		})
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	payments, total, err := h.service.GetVendorPaymentsByVendor(uint(vendorID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"vendor_payments": payments,
			"total":           total,
		},
	})
}

func (h *VendorPaymentHandler) GetVendorPaymentsByStatus(c *fiber.Ctx) error {
	status := c.Query("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Status parameter is required",
		})
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	payments, total, err := h.service.GetVendorPaymentsByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"vendor_payments": payments,
			"total":           total,
		},
	})
}

func (h *VendorPaymentHandler) UpdateVendorPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	var paymentInput input.UpdateVendorPaymentInput

	if err := c.BodyParser(&paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	payment, err := h.service.UpdateVendorPayment(uint(id), &paymentInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor payment updated successfully",
		"data":    payment,
	})
}

func (h *VendorPaymentHandler) RecordPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	var paymentInput input.RecordVendorPaymentInput

	if err := c.BodyParser(&paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.validate.Struct(paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	userName := ""
	if email := c.Locals("user_email"); email != nil {
		userName = fmt.Sprintf("%v", email)
	}

	companyID := uint(0)
	if cid := c.Locals("company_id"); cid != nil {
		if id, ok := cid.(uint); ok {
			companyID = id
		} else if id, ok := cid.(float64); ok {
			companyID = uint(id)
		}
	}

	companyName := ""
	if cname := c.Locals("company_name"); cname != nil {
		companyName = fmt.Sprintf("%v", cname)
	}

	payment, err := h.service.RecordPayment(uint(id), &paymentInput, userID, userName, companyName, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Payment recorded successfully",
		"data":    payment,
	})
}

func (h *VendorPaymentHandler) DeleteVendorPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	if err := h.service.DeleteVendorPayment(uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vendor payment deleted successfully",
	})
}
