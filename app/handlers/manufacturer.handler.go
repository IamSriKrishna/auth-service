package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type ManufacturerHandler struct {
	service services.ManufacturerService
}

func NewManufacturerHandler(service services.ManufacturerService) *ManufacturerHandler {
	return &ManufacturerHandler{service: service}
}

func (h *ManufacturerHandler) CreateManufacturer(c *fiber.Ctx) error {
	var req input.CreateManufacturerInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	manufacturer, err := h.service.Create(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    manufacturer,
	})
}

func (h *ManufacturerHandler) UpdateManufacturer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid manufacturer id")
	}

	var req input.UpdateManufacturerInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	manufacturer, err := h.service.Update(uint(id), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturer,
	})
}

func (h *ManufacturerHandler) GetManufacturerByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid manufacturer id")
	}

	manufacturer, err := h.service.GetByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturer,
	})
}

func (h *ManufacturerHandler) GetAllManufacturers(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid limit parameter")
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid offset parameter")
	}

	manufacturers, err := h.service.GetAll(limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    manufacturers,
	})
}	

func (h *ManufacturerHandler) DeleteManufacturer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid manufacturer id")
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "manufacturer deleted successfully",
	})
}