package handlers

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"

	"github.com/gofiber/fiber/v2"
)

type SupportHandler struct {
	supportService services.SupportService
}

func NewSupportHandler(supportService services.SupportService) *SupportHandler {
	return &SupportHandler{
		supportService: supportService,
	}
}

func (h *SupportHandler) CreateSupport(c *fiber.Ctx) error {
	var req input.CreateSupportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.supportService.CreateSupport(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
