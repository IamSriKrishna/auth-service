package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CustomerPaymentHandler struct {
	service  services.CustomerPaymentService
	validate *validator.Validate
}

func NewCustomerPaymentHandler(service services.CustomerPaymentService) *CustomerPaymentHandler {
	return &CustomerPaymentHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *CustomerPaymentHandler) CreateCustomerPayment(c *fiber.Ctx) error {
	var paymentInput input.CreateCustomerPaymentInput

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

	payment, err := h.service.CreateCustomerPayment(&paymentInput, userID, userName, companyName, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Customer payment created successfully",
		"data":    payment,
	})
}

func (h *CustomerPaymentHandler) GetCustomerPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	payment, err := h.service.GetCustomerPayment(uint(id))
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

func (h *CustomerPaymentHandler) GetAllCustomerPayments(c *fiber.Ctx) error {
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

	payments, total, err := h.service.GetAllCustomerPayments(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"customer_payments": payments,
			"total":             total,
		},
	})
}

func (h *CustomerPaymentHandler) GetCustomerPaymentsBySalesOrder(c *fiber.Ctx) error {
	soID := c.Params("salesOrderId")
	if soID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Sales order ID is required",
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

	payments, total, err := h.service.GetCustomerPaymentsBySalesOrder(soID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"customer_payments": payments,
			"total":             total,
		},
	})
}

func (h *CustomerPaymentHandler) GetCustomerPaymentsByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Params("customerId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid customer ID",
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

	payments, total, err := h.service.GetCustomerPaymentsByCustomer(uint(customerID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"customer_payments": payments,
			"total":             total,
		},
	})
}

func (h *CustomerPaymentHandler) GetCustomerPaymentsByStatus(c *fiber.Ctx) error {
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

	payments, total, err := h.service.GetCustomerPaymentsByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"customer_payments": payments,
			"total":             total,
		},
	})
}

func (h *CustomerPaymentHandler) UpdateCustomerPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	var paymentInput input.UpdateCustomerPaymentInput

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

	payment, err := h.service.UpdateCustomerPayment(uint(id), &paymentInput, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Customer payment updated successfully",
		"data":    payment,
	})
}

func (h *CustomerPaymentHandler) RecordPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	var paymentInput input.RecordCustomerPaymentInput

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

func (h *CustomerPaymentHandler) DeleteCustomerPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid payment ID",
		})
	}

	if err := h.service.DeleteCustomerPayment(uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Customer payment deleted successfully",
	})
}
