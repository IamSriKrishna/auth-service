package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
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

func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var input input.CreateInvoiceInput

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

	userID := ""
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	invoice, err := h.service.CreateInvoice(&input, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(invoice)
}

func (h *InvoiceHandler) GetInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	invoice, err := h.service.GetInvoice(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (h *InvoiceHandler) GetAllInvoices(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	invoices, err := h.service.GetAllInvoices(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *InvoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	var input input.UpdateInvoiceInput

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

	userID := ""
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	invoice, err := h.service.UpdateInvoice(id, &input, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (h *InvoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteInvoice(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Invoice deleted successfully",
	})
}

func (h *InvoiceHandler) GetInvoicesByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerId")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	invoices, err := h.service.GetInvoicesByCustomer(customerID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *InvoiceHandler) GetInvoicesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	invoices, err := h.service.GetInvoicesByStatus(status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}

func (h *InvoiceHandler) UpdateInvoiceStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var input input.InvoiceStatusUpdateInput

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

	invoice, err := h.service.UpdateInvoiceStatus(id, input.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
	var input input.CreatePaymentInput

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

	userID := ""
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	payment, err := h.service.CreatePayment(&input, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	payment, err := h.service.GetPayment(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(payment)
}

func (h *PaymentHandler) GetPaymentsByInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("invoiceId")

	payments, err := h.service.GetPaymentsByInvoice(invoiceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(payments)
}

func (h *PaymentHandler) DeletePayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	if err := h.service.DeletePayment(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment deleted successfully",
	})
}
