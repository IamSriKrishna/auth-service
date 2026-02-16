package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin endpoints
type AdminHandler struct {
	adminService services.AdminService
}

// NewAdminHandler creates a new admin handler instance
func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// CreateUser handles user creation by super admin
// @Summary Create user
// @Description Create a new admin or partner user
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body input.CreateUserRequest true "Create user request"
// @Success 200 {object} output.UserInfo
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/create-user [post]
func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	var req input.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	createdBy := c.Locals("user_id").(uint)

	resp, err := h.adminService.CreateUser(c.Context(), createdBy, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

// ResetPassword handles password reset by super admin
// @Summary Reset password
// @Description Reset user password (super admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body input.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/reset-password [post]
func (h *AdminHandler) ResetAdminPassword(c *fiber.Ctx) error {
	var req input.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	err := h.adminService.ResetPassword(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}

func (h *AdminHandler) ResetUserPassword(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("partner_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid partner ID",
		})
	}
	var req input.ResetUserPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	err = h.adminService.ResetUserPassword(c.Context(), &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}

// GetUsers handles user list retrieval
// @Summary Get users
// @Description Get list of users with pagination and search
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by email or phone"
// @Success 200 {object} output.PaginatedResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/users [get]
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	search := c.Query("search", "")

	resp, err := h.adminService.GetUsers(c.Context(), page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

// GetUser handles specific user retrieval
// @Summary Get user
// @Description Get specific user details
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} output.UserInfo
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	resp, err := h.adminService.GetUser(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

// UpdateUser handles user update
// @Summary Update user
// @Description Update user information
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param request body input.UpdateUserRequest true "Update user request"
// @Success 200 {object} output.UserInfo
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	var req input.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	resp, err := h.adminService.UpdateUser(c.Context(), uint(userID), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(resp)
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Delete user account
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	err = h.adminService.DeleteUser(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

// UpdateUserStatus handles user status update
// @Summary Update user status
// @Description Update user status
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param status body object{status=string} true "Status update"
// @Success 200 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	err = h.adminService.UpdateUserStatus(c.Context(), uint(userID), req.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "User status updated successfully",
	})
}

// UpdateUserRole handles user role update
// @Summary Update user role
// @Description Update user role
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param role body object{role_name=string} true "Role update"
// @Success 200 {object} output.SuccessResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/users/{id}/role [put]
func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid user ID",
		})
	}

	var req struct {
		RoleName string `json:"role_name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: "Invalid request body",
		})
	}

	err = h.adminService.UpdateUserRole(c.Context(), uint(userID), req.RoleName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(output.SuccessResponse{
		Success: true,
		Message: "User role updated successfully",
	})
}

// GetDashboardStats handles dashboard statistics retrieval
// @Summary Get dashboard statistics
// @Description Get dashboard statistics with filters (customer_type can be: mobile_user, partner, vendor, admin, superadmin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param customer_type query string false "Filter by user type (mobile_user, partner, vendor, admin, superadmin)"
// @Param from_date query string false "From date (YYYY-MM-DD)"
// @Param to_date query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} output.DashboardStatsResponse
// @Failure 400 {object} output.ErrorResponse
// @Router /auth/admin/dashboard/stats [get]
func (h *AdminHandler) GetDashboardStats(c *fiber.Ctx) error {
	// Parse query parameters
	filter := &input.DashboardStatsFilter{}

	if customerType := c.Query("customer_type"); customerType != "" {
		if customerType == "vendor" {
			customerType = "admin"
		}
		filter.CustomerType = &customerType
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		filter.FromDate = &fromDate
	}

	if toDate := c.Query("to_date"); toDate != "" {
		filter.ToDate = &toDate
	}

	// Get stats from service
	stats, err := h.adminService.GetDashboardStats(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(stats)
}
