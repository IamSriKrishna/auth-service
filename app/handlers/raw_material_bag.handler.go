package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type RawMaterialBagHandler struct {
	service services.RawMaterialBagService
}

func NewRawMaterialBagHandler(service services.RawMaterialBagService) *RawMaterialBagHandler {
	return &RawMaterialBagHandler{service: service}
}

func (h *RawMaterialBagHandler) ReceiveBags(c *fiber.Ctx) error {
	var req input.ReceiveRawMaterialBagsInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	res, err := h.service.ReceiveBags(&req, createdBy)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Raw material bags received successfully",
		"data":    res,
	})
}

func (h *RawMaterialBagHandler) GetByProduct(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	res, err := h.service.GetBagsByProduct(productID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *RawMaterialBagHandler) GetByPurchaseOrder(c *fiber.Ctx) error {
	poID := c.Params("po_id")

	res, err := h.service.GetBagsByPurchaseOrder(poID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}
func (h *RawMaterialBagHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	res, err := h.service.GetAll(limit, offset)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *RawMaterialBagHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "bag not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}
