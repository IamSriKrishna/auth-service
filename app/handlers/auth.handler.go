package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/bbapp-org/auth-service/app/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) handleError(c *fiber.Ctx, err error) error {
	if httpErr, ok := err.(*utils.HTTPError); ok {
		return c.Status(httpErr.Code).JSON(output.ErrorResponse{
			Error:   true,
			Message: httpErr.Message,
			Code:    httpErr.Code,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
		Error:   true,
		Message: err.Error(),
	})
}

func (h *AuthHandler) RegisterEmail(c *fiber.Ctx) error {
	var req input.RegisterEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.RegisterEmail(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *AuthHandler) RegisterPhone(c *fiber.Ctx) error {
	var req input.RegisterPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.RegisterPhone(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *AuthHandler) RegisterGoogle(c *fiber.Ctx) error {
	var req input.RegisterGoogleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.RegisterGoogle(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *AuthHandler) LoginEmail(c *fiber.Ctx) error {
	var req input.LoginEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.LoginEmail(c.Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) LoginPhone(c *fiber.Ctx) error {
	var req input.LoginPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.LoginPhone(c.Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) LoginGoogle(c *fiber.Ctx) error {
	var req input.LoginGoogleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.LoginGoogle(c.Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) LoginApple(c *fiber.Ctx) error {
	var req input.LoginAppleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.LoginApple(c.Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) LoginPassword(c *fiber.Ctx) error {
	var req input.LoginPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.LoginPassword(c.Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(resp)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req input.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.RefreshToken(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	resp, err := h.authService.GetUserInfo(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req input.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	err = h.authService.ChangePassword(c.Context(), uint(userID), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	var req input.TokenValidationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.authService.ValidateToken(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	tokenID := "placeholder_token_id"

	err := h.authService.Logout(c.Context(), userID, tokenID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}
