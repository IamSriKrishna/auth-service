package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type BrandHandler struct {
	service services.BrandService
}

func NewBrandHandler(service services.BrandService) *BrandHandler {
	return &BrandHandler{service: service}
}

func (h *BrandHandler) CreateBrand(c *fiber.Ctx) error {
	var req input.CreateBrandInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	brand, err := h.service.Create(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    brand,
	})
}

func (h *BrandHandler) GetBrandByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid brand id")
	}

	brand, err := h.service.FindByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    brand,
	})
}

func (h *BrandHandler) GetAllBrands(c *fiber.Ctx) error	 {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid limit parameter")
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid offset parameter")
	}

	brands, total, err := h.service.FindAll(limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    brands,
		"total":   total,
	})
}

func (h *BrandHandler) UpdateBrand(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid brand id")
	}

	var req input.UpdateBrandInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	brand, err := h.service.Update(uint(id), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    brand,
	})
}

func (h *BrandHandler) DeleteBrand(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid brand id")
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "brand deleted successfully",
	})
}
