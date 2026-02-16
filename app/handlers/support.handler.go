package handlers

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"

	"github.com/gofiber/fiber/v2"
)

// SupportHandler handles support endpoints
type SupportHandler struct {
	supportService services.SupportService
}

// NewSupportHandler creates a new support handler instance
func NewSupportHandler(supportService services.SupportService) *SupportHandler {
	return &SupportHandler{
		supportService: supportService,
	}
}

// CreateSupport handles support ticket creation
// @Summary Create support ticket
// @Description Create a new support ticket (public endpoint)
// @Tags Support
// @Accept json
// @Produce json
// @Param request body input.CreateSupportRequest true "Support ticket request"
// @Success 201 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /public/support [post]
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
