package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	service  services.InvoiceService
	validate *validator.Validate
}

func NewInvoiceHandler(service services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		service:  service,
		validate: validator.New(),
	}
}

func invoiceLocalUint(c *fiber.Ctx, key string) uint {
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

func invoiceAuthContext(c *fiber.Ctx) (string, uint, error) {
	userID := invoiceLocalUint(c, "user_id")
	companyID := invoiceLocalUint(c, "company_id")

	if userID == 0 {
		return "", 0, fmt.Errorf("invalid authenticated user")
	}

	if companyID == 0 {
		return "", 0, fmt.Errorf("user is not assigned to a company")
	}

	return strconv.FormatUint(uint64(userID), 10), companyID, nil
}

func invoicePagination(c *fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
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

func invoiceContextError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var req input.CreateInvoiceInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoice, err := h.service.CreateInvoiceForCompany(
		&req,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(invoice)
}

func (h *InvoiceHandler) GetInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoice, err := h.service.GetInvoiceByCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (h *InvoiceHandler) GetAllInvoices(c *fiber.Ctx) error {
	limit, offset := invoicePagination(c)

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoices, err := h.service.GetAllInvoicesByCompany(
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *InvoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	var req input.UpdateInvoiceInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoice, err := h.service.UpdateInvoiceForCompany(
		id,
		&req,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (h *InvoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	if err := h.service.DeleteInvoiceForCompany(id, companyID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Invoice deleted successfully",
	})
}

func (h *InvoiceHandler) GetInvoicesByCustomer(c *fiber.Ctx) error {
	customerID64, err := strconv.ParseUint(
		c.Params("customerId"),
		10,
		32,
	)
	if err != nil || customerID64 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	limit, offset := invoicePagination(c)

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoices, err := h.service.GetInvoicesByCustomerAndCompany(
		uint(customerID64),
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *InvoiceHandler) GetInvoicesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	limit, offset := invoicePagination(c)

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoices, err := h.service.GetInvoicesByStatusAndCompany(
		status,
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *InvoiceHandler) UpdateInvoiceStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var req input.InvoiceStatusUpdateInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	invoice, err := h.service.UpdateInvoiceStatusForCompany(
		id,
		req.Status,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

type SalespersonHandler struct {
	service  services.SalespersonService
	validate *validator.Validate
}

func NewSalespersonHandler(service services.SalespersonService) *SalespersonHandler {
	return &SalespersonHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *SalespersonHandler) CreateSalesperson(c *fiber.Ctx) error {
	var input input.CreateSalespersonInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	salesperson, err := h.service.CreateSalesperson(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(salesperson)
}

func (h *SalespersonHandler) GetSalesperson(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid salesperson ID",
		})
	}

	salesperson, err := h.service.GetSalesperson(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Salesperson not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(salesperson)
}

func (h *SalespersonHandler) GetAllSalespersons(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	salespersons, err := h.service.GetAllSalespersons(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(salespersons)
}

func (h *SalespersonHandler) UpdateSalesperson(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid salesperson ID",
		})
	}

	var input input.UpdateSalespersonInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	salesperson, err := h.service.UpdateSalesperson(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(salesperson)
}

func (h *SalespersonHandler) DeleteSalesperson(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid salesperson ID",
		})
	}

	if err := h.service.DeleteSalesperson(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Salesperson deleted successfully",
	})
}

type TaxHandler struct {
	service  services.TaxService
	validate *validator.Validate
}

func NewTaxHandler(service services.TaxService) *TaxHandler {
	return &TaxHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *TaxHandler) CreateTax(c *fiber.Ctx) error {
	var input input.CreateTaxInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tax, err := h.service.CreateTax(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tax)
}

func (h *TaxHandler) GetTax(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tax ID",
		})
	}

	tax, err := h.service.GetTax(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tax not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tax)
}

func (h *TaxHandler) GetAllTaxes(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	taxes, err := h.service.GetAllTaxes(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(taxes)
}

func (h *TaxHandler) UpdateTax(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tax ID",
		})
	}

	var input input.UpdateTaxInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tax, err := h.service.UpdateTax(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(tax)
}

func (h *TaxHandler) DeleteTax(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tax ID",
		})
	}

	if err := h.service.DeleteTax(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tax deleted successfully",
	})
}

type PaymentHandler struct {
	service  services.PaymentService
	validate *validator.Validate
}

func NewPaymentHandler(service services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var req input.CreatePaymentInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	payment, err := h.service.CreatePaymentForCompany(
		&req,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	payment, err := h.service.GetPaymentForCompany(
		uint(id),
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(payment)
}

func (h *PaymentHandler) GetPaymentsByInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("invoiceId")

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	payments, err := h.service.GetPaymentsByInvoiceForCompany(
		invoiceID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(payments)
}

func (h *PaymentHandler) DeletePayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	if err := h.service.DeletePaymentForCompany(
		uint(id),
		companyID,
	); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment deleted successfully",
	})
}

// CreateInvoiceFromVariants creates an invoice preview from variants.
// It validates that the customer and products belong to the authenticated company.
func (h *InvoiceHandler) CreateInvoiceFromVariants(c *fiber.Ctx) error {
	var invoiceInput input.CreateInvoiceInput

	if err := c.BodyParser(&invoiceInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(invoiceInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, companyID, err := invoiceAuthContext(c)
	if err != nil {
		return invoiceContextError(c, err)
	}

	if err := h.service.ValidateInvoiceInputForCompany(
		&invoiceInput,
		companyID,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invoiceNo := fmt.Sprintf(
		"INV-%d-%d",
		time.Now().Year(),
		time.Now().UnixNano(),
	)

	var subTotal float64
	lineItems := make(
		[]output.InvoiceLineItemOutput,
		len(invoiceInput.LineItems),
	)

	for i, item := range invoiceInput.LineItems {
		if item.ProductID == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "product_id is required",
			})
		}

		amount := item.Quantity * item.Rate
		subTotal += amount

		lineItems[i] = output.InvoiceLineItemOutput{
			ProductID:   *item.ProductID,
			ProductName: item.ProductName,
			SKU:         item.SKU,
			Account:     item.Account,
			Quantity:    item.Quantity,
			Rate:        item.Rate,
			Amount:      amount,
		}
	}

	total := subTotal +
		invoiceInput.ShippingCharges +
		invoiceInput.Adjustment

	invoiceOutput := output.CreateInvoiceVariantOutput(
		invoiceNo,
		invoiceInput.SalesOrderID,
		invoiceInput.CustomerID,
		"",
		lineItems,
		subTotal,
		invoiceInput.ShippingCharges,
		0,
		total,
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    invoiceOutput,
		"message": "Invoice created successfully from variants",
	})
}