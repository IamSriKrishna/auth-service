package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service  services.ProductService
	validate *validator.Validate
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var input input.CreateProductInput

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

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	product, err := h.service.CreateProduct(&input, createdBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.service.GetProduct(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	products, err := h.service.GetAllProducts(limit, offset, createdBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var input input.UpdateProductInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	product, err := h.service.UpdateProduct(id, &input, createdBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	if err := h.service.DeleteProduct(id, createdBy); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

func (h *ProductHandler) GetProductVariants(c *fiber.Ctx) error {
	productID := c.Params("id")

	variants, err := h.service.GetProductVariants(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product_id": productID,
		"variants":   variants,
		"total":      len(variants),
	})
}
