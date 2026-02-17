package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type BankHandler struct {
	service services.BankService
}

func NewBankHandler(service services.BankService) *BankHandler {
	return &BankHandler{service: service}
}

func (h *BankHandler) CreateBank(c *fiber.Ctx) error {
	var req input.CreateBankInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	bank, err := h.service.Create(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    bank,
	})
}

func (h *BankHandler) GetBankByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid bank id")
	}

	bank, err := h.service.FindByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bank,
	})
}

func (h *BankHandler) GetAllBanks(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	banks, total, err := h.service.FindAll(limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"banks": banks,
			"total": total,
		},
	})
}

func (h *BankHandler) UpdateBank(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid bank id")
	}

	var req input.UpdateBankInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	bank, err := h.service.Update(uint(id), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bank,
	})
}

func (h *BankHandler) DeleteBank(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid bank id")
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Bank deleted successfully",
	})
}
