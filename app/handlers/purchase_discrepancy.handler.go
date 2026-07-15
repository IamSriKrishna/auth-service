package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PurchaseDispenseHandler struct {
	service  services.PurchaseDispenseService
	validate *validator.Validate
}

func NewPurchaseDispenseHandler(
	service services.PurchaseDispenseService,
) *PurchaseDispenseHandler {
	return &PurchaseDispenseHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *PurchaseDispenseHandler) CreateDispense(c *fiber.Ctx) error {
	claimID := c.Params("claimId")

	var request input.CreatePurchaseDispenseInput
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

	result, err := h.service.CreateDispense(
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "purchase dispense created and replacement added to stock",
		"data":    result,
	})
}

func (h *PurchaseDispenseHandler) GetDispenseByID(c *fiber.Ctx) error {
	_, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.GetDispenseByID(c.Params("id"), companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "purchase dispense not found",
		})
	}

	return c.JSON(fiber.Map{"success": true, "data": result})
}

func (h *PurchaseDispenseHandler) GetDispensesByClaim(c *fiber.Ctx) error {
	_, companyID, err := purchaseOrderContext(c)
	if err != nil {
		return purchaseOrderContextError(c, err)
	}

	result, err := h.service.GetDispensesByClaim(c.Params("claimId"), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"total":   len(result),
		"data":    result,
	})
}

func (h *PurchaseDispenseHandler) GetDispensesByClaimItem(c *fiber.Ctx) error {
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

	result, err := h.service.GetDispensesByClaimItem(uint(itemID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"total":   len(result),
		"data":    result,
	})
}
