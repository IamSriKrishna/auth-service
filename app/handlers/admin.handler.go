package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

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

func (h *AdminHandler) GetDashboardStats(c *fiber.Ctx) error {
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

	stats, err := h.adminService.GetDashboardStats(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(output.ErrorResponse{
			Error:   true,
			Message: err.Error(),
		})
	}

	return c.JSON(stats)
}
