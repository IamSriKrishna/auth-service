package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PurchaseClaimHandler struct {
	service  services.PurchaseClaimService
	validate *validator.Validate
}

func NewPurchaseClaimHandler(
	service services.PurchaseClaimService,
) *PurchaseClaimHandler {
	return &PurchaseClaimHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *PurchaseClaimHandler) GetPurchaseOrderClaimSource(c *fiber.Ctx) error {
	purchaseOrderID := c.Params("purchaseOrderId")

	_, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.GetPurchaseOrderClaimSource(
		purchaseOrderID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *PurchaseClaimHandler) CreateClaim(c *fiber.Ctx) error {
	var request input.CreatePurchaseClaimInput

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if err := h.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.CreateClaim(
		&request,
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
		"message": "purchase claim created successfully",
		"data":    result,
	})
}

func (h *PurchaseClaimHandler) GetClaimByID(c *fiber.Ctx) error {
	id := c.Params("id")

	_, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.GetClaimByID(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *PurchaseClaimHandler) GetClaimsByPurchaseOrder(c *fiber.Ctx) error {
	purchaseOrderID := c.Params("purchaseOrderId")

	_, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.GetClaimsByPurchaseOrder(
		purchaseOrderID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"total":   len(result),
		"data":    result,
	})
}

func (h *PurchaseClaimHandler) ReceiveReplacement(c *fiber.Ctx) error {
	claimID := c.Params("id")

	var request input.ReceivePurchaseClaimReplacementInput
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if err := h.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	userID, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.ReceiveReplacement(
		claimID,
		&request,
		userID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "vendor replacement received and added to stock",
		"data":    result,
	})
}

func (h *PurchaseClaimHandler) GetReplacementReceipts(c *fiber.Ctx) error {
	itemID, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil || itemID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid claim item ID",
		})
	}

	_, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.GetReplacementReceipts(uint(itemID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"total":   len(result),
		"data":    result,
	})
}
