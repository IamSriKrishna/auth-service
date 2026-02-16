package handlers

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type OpeningStockHandler struct {
	service services.OpeningStockService
}

func NewOpeningStockHandler(service services.OpeningStockService) *OpeningStockHandler {
	return &OpeningStockHandler{service: service}
}

func (h *OpeningStockHandler) UpdateOpeningStock(c *fiber.Ctx) error {
	itemID := c.Params("id")

	var input input.OpeningStockInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID := c.Locals("user_id")
	userIDStr, _ := userID.(string)

	result, err := h.service.UpdateOpeningStock(itemID, &input, userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Opening stock updated successfully",
		"data":    result,
	})
}

func (h *OpeningStockHandler) GetOpeningStock(c *fiber.Ctx) error {
	itemID := c.Params("id")

	result, err := h.service.GetOpeningStock(itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}

func (h *OpeningStockHandler) UpdateVariantsOpeningStock(c *fiber.Ctx) error {
	itemID := c.Params("id")

	var input input.UpdateVariantsOpeningStockInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID := c.Locals("user_id")
	userIDStr, _ := userID.(string)

	result, err := h.service.UpdateVariantsOpeningStock(itemID, &input, userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Variants opening stock updated successfully",
		"data":    result,
	})
}

func (h *OpeningStockHandler) GetVariantsOpeningStock(c *fiber.Ctx) error {
	itemID := c.Params("id")

	result, err := h.service.GetVariantsOpeningStock(itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}

func (h *OpeningStockHandler) GetStockSummary(c *fiber.Ctx) error {
	itemID := c.Params("id")

	result, err := h.service.GetStockSummary(itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}
