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

func NewVendorPaymentHandler(
	service services.VendorPaymentService,
) *VendorPaymentHandler {
	return &VendorPaymentHandler{
		service:  service,
		validate: validator.New(),
	}
}

func vendorPaymentContext(
	c *fiber.Ctx,
) (
	userID string,
	userName string,
	companyID uint,
	companyName string,
	err error,
) {
	var parsedUserID uint

	switch value := c.Locals("user_id").(type) {
	case uint:
		parsedUserID = value
	case int:
		if value > 0 {
			parsedUserID = uint(value)
		}
	case float64:
		if value > 0 {
			parsedUserID = uint(value)
		}
	case string:
		parsed, parseErr := strconv.ParseUint(
			value,
			10,
			64,
		)
		if parseErr == nil {
			parsedUserID = uint(parsed)
		}
	}

	switch value := c.Locals("company_id").(type) {
	case uint:
		companyID = value
	case int:
		if value > 0 {
			companyID = uint(value)
		}
	case float64:
		if value > 0 {
			companyID = uint(value)
		}
	case string:
		parsed, parseErr := strconv.ParseUint(
			value,
			10,
			64,
		)
		if parseErr == nil {
			companyID = uint(parsed)
		}
	}

	if email := c.Locals("user_email"); email != nil {
		userName = fmt.Sprintf("%v", email)
	}

	if name := c.Locals("company_name"); name != nil {
		companyName = fmt.Sprintf("%v", name)
	}

	if parsedUserID == 0 {
		return "", "", 0, "", fmt.Errorf(
			"invalid authenticated user",
		)
	}

	if companyID == 0 {
		return "", "", 0, "", fmt.Errorf(
			"user is not assigned to a company",
		)
	}

	userID = strconv.FormatUint(
		uint64(parsedUserID),
		10,
	)

	return userID, userName, companyID, companyName, nil
}

func vendorPaymentPagination(
	c *fiber.Ctx,
) (int, int) {
	limit, err := strconv.Atoi(
		c.Query("limit", "10"),
	)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(
		c.Query("offset", "0"),
	)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

func vendorPaymentContextError(
	c *fiber.Ctx,
	err error,
) error {
	return c.Status(fiber.StatusForbidden).JSON(
		fiber.Map{
			"success": false,
			"error":   err.Error(),
		},
	)
}

func (h *VendorPaymentHandler) CreateVendorPayment(
	c *fiber.Ctx,
) error {
	var paymentInput input.CreateVendorPaymentInput

	if err := c.BodyParser(&paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			},
		)
	}

	if err := h.validate.Struct(paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userID, userName, companyID, companyName, err :=
		vendorPaymentContext(c)
	if err != nil {
		return vendorPaymentContextError(c, err)
	}

	payment, err :=
		h.service.CreateVendorPaymentForCompany(
			&paymentInput,
			userID,
			userName,
			companyName,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"success": true,
			"message": "Vendor payment created successfully",
			"data":    payment,
		},
	)
}

func (h *VendorPaymentHandler) GetVendorPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		32,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	_, _, companyID, _, contextErr :=
		vendorPaymentContext(c)
	if contextErr != nil {
		return vendorPaymentContextError(c, contextErr)
	}

	payment, err :=
		h.service.GetVendorPaymentByCompany(
			uint(id),
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data":    payment,
		},
	)
}

func (h *VendorPaymentHandler) GetAllVendorPayments(
	c *fiber.Ctx,
) error {
	limit, offset := vendorPaymentPagination(c)

	_, _, companyID, _, err :=
		vendorPaymentContext(c)
	if err != nil {
		return vendorPaymentContextError(c, err)
	}

	payments, total, err :=
		h.service.GetAllVendorPaymentsByCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(
				fiber.Map{
					"success": false,
					"error":   err.Error(),
				},
			)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data": fiber.Map{
				"vendor_payments": payments,
				"total":           total,
			},
		},
	)
}

func (h *VendorPaymentHandler) GetVendorPaymentsByPurchaseOrder(
	c *fiber.Ctx,
) error {
	purchaseOrderID := c.Params("purchaseOrderId")
	if purchaseOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Purchase order ID is required",
			},
		)
	}

	limit, offset := vendorPaymentPagination(c)

	_, _, companyID, _, err :=
		vendorPaymentContext(c)
	if err != nil {
		return vendorPaymentContextError(c, err)
	}

	payments, total, err :=
		h.service.GetVendorPaymentsByPurchaseOrderAndCompany(
			purchaseOrderID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data": fiber.Map{
				"vendor_payments": payments,
				"total":           total,
			},
		},
	)
}

func (h *VendorPaymentHandler) GetVendorPaymentsByVendor(
	c *fiber.Ctx,
) error {
	vendorID, err := strconv.ParseUint(
		c.Params("vendorId"),
		10,
		32,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid vendor ID",
			},
		)
	}

	limit, offset := vendorPaymentPagination(c)

	_, _, companyID, _, contextErr :=
		vendorPaymentContext(c)
	if contextErr != nil {
		return vendorPaymentContextError(
			c,
			contextErr,
		)
	}

	payments, total, err :=
		h.service.GetVendorPaymentsByVendorAndCompany(
			uint(vendorID),
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data": fiber.Map{
				"vendor_payments": payments,
				"total":           total,
			},
		},
	)
}

func (h *VendorPaymentHandler) GetVendorPaymentsByStatus(
	c *fiber.Ctx,
) error {
	status := c.Query("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Status parameter is required",
			},
		)
	}

	limit, offset := vendorPaymentPagination(c)

	_, _, companyID, _, err :=
		vendorPaymentContext(c)
	if err != nil {
		return vendorPaymentContextError(c, err)
	}

	payments, total, err :=
		h.service.GetVendorPaymentsByStatusAndCompany(
			status,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(
				fiber.Map{
					"success": false,
					"error":   err.Error(),
				},
			)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data": fiber.Map{
				"vendor_payments": payments,
				"total":           total,
			},
		},
	)
}

func (h *VendorPaymentHandler) UpdateVendorPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		32,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	var paymentInput input.UpdateVendorPaymentInput

	if err := c.BodyParser(&paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			},
		)
	}

	if err := h.validate.Struct(paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userID, _, companyID, _, contextErr :=
		vendorPaymentContext(c)
	if contextErr != nil {
		return vendorPaymentContextError(
			c,
			contextErr,
		)
	}

	payment, err :=
		h.service.UpdateVendorPaymentForCompany(
			uint(id),
			&paymentInput,
			userID,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Vendor payment updated successfully",
			"data":    payment,
		},
	)
}

func (h *VendorPaymentHandler) RecordPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		32,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	var paymentInput input.RecordVendorPaymentInput

	if err := c.BodyParser(&paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			},
		)
	}

	if err := h.validate.Struct(paymentInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userID, userName, companyID, companyName, contextErr :=
		vendorPaymentContext(c)
	if contextErr != nil {
		return vendorPaymentContextError(
			c,
			contextErr,
		)
	}

	payment, err :=
		h.service.RecordPaymentForCompany(
			uint(id),
			&paymentInput,
			userID,
			userName,
			companyName,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Payment recorded successfully",
			"data":    payment,
		},
	)
}

func (h *VendorPaymentHandler) DeleteVendorPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		32,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	_, _, companyID, _, contextErr :=
		vendorPaymentContext(c)
	if contextErr != nil {
		return vendorPaymentContextError(
			c,
			contextErr,
		)
	}

	if err := h.service.DeleteVendorPaymentForCompany(
		uint(id),
		companyID,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Vendor payment deleted successfully",
		},
	)
}
