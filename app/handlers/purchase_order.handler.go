package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PurchaseOrderHandler struct {
	service  services.PurchaseOrderService
	validate *validator.Validate
}

func NewPurchaseOrderHandler(
	service services.PurchaseOrderService,
) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		service:  service,
		validate: validator.New(),
	}
}

func purchaseOrderContext(
	c *fiber.Ctx,
) (string, uint, error) {
	var userID uint

	switch value := c.Locals("user_id").(type) {
	case uint:
		userID = value
	case int:
		if value > 0 {
			userID = uint(value)
		}
	case float64:
		if value > 0 {
			userID = uint(value)
		}
	case string:
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			userID = uint(parsed)
		}
	}

	var companyID uint
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
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			companyID = uint(parsed)
		}
	}

	if userID == 0 {
		return "", 0, fmt.Errorf(
			"invalid authenticated user",
		)
	}
	if companyID == 0 {
		return "", 0, fmt.Errorf(
			"user is not assigned to a company",
		)
	}

	return strconv.FormatUint(
		uint64(userID),
		10,
	), companyID, nil
}

func purchaseOrderPagination(
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

func purchaseOrderContextError(
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

func (h *PurchaseOrderHandler) CreatePurchaseOrder(
	c *fiber.Ctx,
) error {
	var poInput input.CreatePurchaseOrderInput

	if err := c.BodyParser(&poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			},
		)
	}

	if err := h.validate.Struct(poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userID, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrder, err :=
		h.service.CreatePurchaseOrderForCompany(
			&poInput,
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

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"success": true,
			"message": "Purchase order created successfully",
			"data":    purchaseOrder,
		},
	)
}

func (h *PurchaseOrderHandler) GetPurchaseOrder(
	c *fiber.Ctx,
) error {
	id := c.Params("id")

	_, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrder, err :=
		h.service.GetPurchaseOrderByCompany(
			id,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"error":   "Purchase order not found",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data":    purchaseOrder,
		},
	)
}

func (h *PurchaseOrderHandler) GetAllPurchaseOrders(
	c *fiber.Ctx,
) error {
	limit, offset := purchaseOrderPagination(c)

	_, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrders, err :=
		h.service.GetAllPurchaseOrdersByCompany(
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
			"data":    purchaseOrders,
		},
	)
}

func (h *PurchaseOrderHandler) UpdatePurchaseOrder(
	c *fiber.Ctx,
) error {
	id := c.Params("id")
	var poInput input.UpdatePurchaseOrderInput

	if err := c.BodyParser(&poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			},
		)
	}

	if err := h.validate.Struct(poInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userID, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrder, err :=
		h.service.UpdatePurchaseOrderForCompany(
			id,
			&poInput,
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
			"message": "Purchase order updated successfully",
			"data":    purchaseOrder,
		},
	)
}

func (h *PurchaseOrderHandler) DeletePurchaseOrder(
	c *fiber.Ctx,
) error {
	id := c.Params("id")

	_, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	if err := h.service.DeletePurchaseOrderForCompany(
		id,
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
			"message": "Purchase order deleted successfully",
		},
	)
}

func (h *PurchaseOrderHandler) GetPurchaseOrdersByVendor(
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

	limit, offset := purchaseOrderPagination(c)

	_, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrders, err :=
		h.service.GetPurchaseOrdersByVendorAndCompany(
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
			"data":    purchaseOrders,
		},
	)
}

func (h *PurchaseOrderHandler) GetPurchaseOrdersByCustomer(
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

	limit, offset := purchaseOrderPagination(c)

	_, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrders, err :=
		h.service.GetPurchaseOrdersByCustomerAndCompany(
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
			"data":    purchaseOrders,
		},
	)
}

func (h *PurchaseOrderHandler) GetPurchaseOrdersByStatus(
	c *fiber.Ctx,
) error {
	status := c.Params("status")
	limit, offset := purchaseOrderPagination(c)

	_, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrders, err :=
		h.service.GetPurchaseOrdersByStatusAndCompany(
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
			"data":    purchaseOrders,
		},
	)
}

func (h *PurchaseOrderHandler) UpdatePurchaseOrderStatus(
	c *fiber.Ctx,
) error {
	id := c.Params("id")
	var statusInput input.UpdatePurchaseOrderStatusInput

	if err := c.BodyParser(&statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			},
		)
	}

	if err := h.validate.Struct(statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		)
	}

	userID, companyID, err :=
		purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	purchaseOrder, err :=
		h.service.UpdatePurchaseOrderStatusForCompany(
			id,
			domain.PurchaseOrderStatus(statusInput.Status),
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
			"message": "Purchase order status updated successfully",
			"data":    purchaseOrder,
		},
	)
}
