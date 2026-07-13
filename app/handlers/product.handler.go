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

func productLocalUint(c *fiber.Ctx, key string) uint {
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

func productAuthContext(c *fiber.Ctx) (string, uint, error) {
	userID := productLocalUint(c, "user_id")
	companyID := productLocalUint(c, "company_id")

	if userID == 0 {
		return "", 0, fmt.Errorf("invalid authenticated user")
	}

	if companyID == 0 {
		return "", 0, fmt.Errorf("user is not assigned to a company")
	}

	return strconv.FormatUint(uint64(userID), 10), companyID, nil
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req input.CreateProductInput

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

	createdBy, companyID, err := productAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	product, err := h.service.CreateProductForCompany(
		&req,
		createdBy,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	_, companyID, err := productAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	product, err := h.service.GetProductByCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
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

	_, companyID, err := productAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	products, err := h.service.GetAllProductsByCompany(
		limit,
		offset,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var req input.UpdateProductInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, companyID, err := productAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	product, err := h.service.UpdateProductForCompany(
		id,
		&req,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	_, companyID, err := productAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.DeleteProductForCompany(
		id,
		companyID,
	); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

func (h *ProductHandler) GetProductVariants(c *fiber.Ctx) error {
	productID := c.Params("id")

	_, companyID, err := productAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	variants, err := h.service.GetProductVariantsByCompany(
		productID,
		companyID,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product_id": productID,
		"variants":   variants,
		"total":      len(variants),
	})
}