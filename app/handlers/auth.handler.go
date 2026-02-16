package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/bbapp-org/auth-service/app/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// handleError handles different types of errors and returns appropriate HTTP responses
func (h *AuthHandler) handleError(c *fiber.Ctx, err error) error {
	// Check if it's a custom HTTP error
	if httpErr, ok := err.(*utils.HTTPError); ok {
		return c.Status(httpErr.Code).JSON(output.ErrorResponse{
			Error:   true,
			Message: httpErr.Message,
			Code:    httpErr.Code,
		})
	}

	// Default to bad request for unknown errors
	return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
		Error:   true,
		Message: err.Error(),
	})
}

// RegisterEmail handles email registration
// @Summary Register with email
// @Description Register a new mobile user with email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.RegisterEmailRequest true "Registration request"
// @Success 200 {object} output.OTPResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/register/email [post]
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

// RegisterPhone handles phone registration
// @Summary Register with phone
// @Description Register a new mobile user with phone
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.RegisterPhoneRequest true "Registration request"
// @Success 200 {object} output.OTPResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/register/phone [post]
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

// RegisterGoogle handles Google OIDC registration
// @Summary Register with Google
// @Description Register a new mobile user with Google OIDC
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.RegisterGoogleRequest true "Registration request"
// @Success 200 {object} output.AuthResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/register/google [post]
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

// LoginEmail handles email login
// @Summary Login with email
// @Description Login mobile user with email (OTP)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.LoginEmailRequest true "Login request"
// @Success 200 {object} output.OTPResponse
// @Failure 400 {object} output.ErrorResponse
// @Failure 403 {object} output.ErrorResponse
// @Failure 404 {object} output.ErrorResponse
// @Failure 500 {object} output.ErrorResponse
// @Router /auth/login/email [post]
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

// LoginPhone handles phone login
// @Summary Login with phone
// @Description Login mobile user with phone (OTP)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.LoginPhoneRequest true "Login request"
// @Success 200 {object} output.OTPResponse
// @Failure 400 {object} output.ErrorResponse
// @Failure 403 {object} output.ErrorResponse
// @Failure 404 {object} output.ErrorResponse
// @Failure 500 {object} output.ErrorResponse
// @Router /auth/login/phone [post]
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

// LoginGoogle handles Google OIDC login
// @Summary Login with Google
// @Description Login mobile user with Google OIDC
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.LoginGoogleRequest true "Login request"
// @Success 200 {object} output.AuthResponse
// @Failure 400 {object} output.ErrorResponse
// @Failure 401 {object} output.ErrorResponse
// @Failure 403 {object} output.ErrorResponse
// @Failure 404 {object} output.ErrorResponse
// @Router /auth/login/google [post]
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

// LoginApple handles Apple login via Firebase ID token
// @Summary Login with Apple (Firebase)
// @Description Login user with Firebase ID token where provider is Apple
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.LoginAppleRequest true "Login request"
// @Success 200 {object} output.AuthResponse
// @Failure 400 {object} output.ErrorResponse
// @Failure 401 {object} output.ErrorResponse
// @Failure 403 {object} output.ErrorResponse
// @Router /auth/login/apple [post]
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

// LoginPassword handles password-based login
// @Summary Login with password
// @Description Login admin/partner user with password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.LoginPasswordRequest true "Login request"
// @Success 200 {object} output.AuthResponse
// @Failure 400 {object} output.ErrorResponse
// @Failure 401 {object} output.ErrorResponse
// @Failure 403 {object} output.ErrorResponse
// @Router /auth/login/password [post]
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

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body input.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} output.AuthResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/refresh-token [post]
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

// GetUserInfo handles user info retrieval
// @Summary Get user info
// @Description Get authenticated user information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} output.UserInfo
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/user-info [get]
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

// ChangePassword handles password change
// @Summary Change password
// @Description Change user password
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body input.ChangePasswordRequest true "Change password request"
// @Success 200 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/change-password [post]
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

// ValidateToken handles token validation
// @Summary Validate token
// @Description Validate JWT token (internal use)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body input.TokenValidationRequest true "Token validation request"
// @Success 200 {object} output.TokenValidationResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/validate-token [post]
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

// Logout handles user logout
// @Summary Logout
// @Description Logout user and invalidate tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// For now, we'll use a placeholder token ID
	// In a real implementation, you'd extract this from the token
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
