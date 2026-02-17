package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ItemHandler struct {
	service  services.ItemService
	validate *validator.Validate
}

func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *ItemHandler) CreateItem(c *fiber.Ctx) error {
	var input input.CreateItemInput

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

	item, err := h.service.CreateItem(&input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ItemHandler) GetItem(c *fiber.Ctx) error {
	id := c.Params("id")

	item, err := h.service.GetItem(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

func (h *ItemHandler) GetAllItems(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	items, err := h.service.GetAllItems(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(items)
}

func (h *ItemHandler) UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var input input.UpdateItemInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	item, err := h.service.UpdateItem(id, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

func (h *ItemHandler) DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteItem(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item deleted successfully",
	})
}

func (h *ItemHandler) GetItemsByType(c *fiber.Ctx) error {
	itemType := c.Query("type")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	items, err := h.service.GetItemsByType(itemType, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(items)
}
