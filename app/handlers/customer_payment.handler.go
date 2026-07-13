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

func NewCustomerPaymentHandler(
	service services.CustomerPaymentService,
) *CustomerPaymentHandler {
	return &CustomerPaymentHandler{
		service:  service,
		validate: validator.New(),
	}
}

func customerPaymentContext(
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
		parsed, parseErr := strconv.ParseUint(value, 10, 64)
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
		parsed, parseErr := strconv.ParseUint(value, 10, 64)
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
		return "", "", 0, "", fmt.Errorf("invalid authenticated user")
	}

	if companyID == 0 {
		return "", "", 0, "", fmt.Errorf(
			"user is not assigned to a company",
		)
	}

	userID = strconv.FormatUint(uint64(parsedUserID), 10)

	return userID, userName, companyID, companyName, nil
}

func customerPaymentPagination(
	c *fiber.Ctx,
) (int, int) {
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

func customerPaymentContextError(
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

func (h *CustomerPaymentHandler) CreateCustomerPayment(
	c *fiber.Ctx,
) error {
	var paymentInput input.CreateCustomerPaymentInput

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
		customerPaymentContext(c)
	if err != nil {
		return customerPaymentContextError(c, err)
	}

	payment, err :=
		h.service.CreateCustomerPaymentForCompany(
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
			"message": "Customer payment created successfully",
			"data":    payment,
		},
	)
}

func (h *CustomerPaymentHandler) GetCustomerPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	_, _, companyID, _, contextErr :=
		customerPaymentContext(c)
	if contextErr != nil {
		return customerPaymentContextError(c, contextErr)
	}

	payment, err :=
		h.service.GetCustomerPaymentByCompany(
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

func (h *CustomerPaymentHandler) GetAllCustomerPayments(
	c *fiber.Ctx,
) error {
	limit, offset := customerPaymentPagination(c)

	_, _, companyID, _, err :=
		customerPaymentContext(c)
	if err != nil {
		return customerPaymentContextError(c, err)
	}

	payments, total, err :=
		h.service.GetAllCustomerPaymentsByCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
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
				"customer_payments": payments,
				"total":             total,
			},
		},
	)
}

func (h *CustomerPaymentHandler) GetCustomerPaymentsBySalesOrder(
	c *fiber.Ctx,
) error {
	salesOrderID := c.Params("salesOrderId")
	if salesOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Sales order ID is required",
			},
		)
	}

	limit, offset := customerPaymentPagination(c)

	_, _, companyID, _, err :=
		customerPaymentContext(c)
	if err != nil {
		return customerPaymentContextError(c, err)
	}

	payments, total, err :=
		h.service.GetCustomerPaymentsBySalesOrderAndCompany(
			salesOrderID,
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
				"customer_payments": payments,
				"total":             total,
			},
		},
	)
}

func (h *CustomerPaymentHandler) GetCustomerPaymentsByCustomer(
	c *fiber.Ctx,
) error {
	customerID, err := strconv.ParseUint(
		c.Params("customerId"),
		10,
		32,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid customer ID",
			},
		)
	}

	limit, offset := customerPaymentPagination(c)

	_, _, companyID, _, contextErr :=
		customerPaymentContext(c)
	if contextErr != nil {
		return customerPaymentContextError(c, contextErr)
	}

	payments, total, err :=
		h.service.GetCustomerPaymentsByCustomerAndCompany(
			uint(customerID),
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
				"customer_payments": payments,
				"total":             total,
			},
		},
	)
}

func (h *CustomerPaymentHandler) GetCustomerPaymentsByStatus(
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

	limit, offset := customerPaymentPagination(c)

	_, _, companyID, _, err :=
		customerPaymentContext(c)
	if err != nil {
		return customerPaymentContextError(c, err)
	}

	payments, total, err :=
		h.service.GetCustomerPaymentsByStatusAndCompany(
			status,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
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
				"customer_payments": payments,
				"total":             total,
			},
		},
	)
}

func (h *CustomerPaymentHandler) UpdateCustomerPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	var paymentInput input.UpdateCustomerPaymentInput

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
		customerPaymentContext(c)
	if contextErr != nil {
		return customerPaymentContextError(c, contextErr)
	}

	payment, err :=
		h.service.UpdateCustomerPaymentForCompany(
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
			"message": "Customer payment updated successfully",
			"data":    payment,
		},
	)
}

func (h *CustomerPaymentHandler) RecordPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	var paymentInput input.RecordCustomerPaymentInput

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
		customerPaymentContext(c)
	if contextErr != nil {
		return customerPaymentContextError(c, contextErr)
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

func (h *CustomerPaymentHandler) DeleteCustomerPayment(
	c *fiber.Ctx,
) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid payment ID",
			},
		)
	}

	_, _, companyID, _, contextErr :=
		customerPaymentContext(c)
	if contextErr != nil {
		return customerPaymentContextError(c, contextErr)
	}

	if err := h.service.DeleteCustomerPaymentForCompany(
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
			"message": "Customer payment deleted successfully",
		},
	)
}
