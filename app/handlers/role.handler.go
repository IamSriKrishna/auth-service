package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	service services.RoleService
}

func NewRoleHandler(service services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req input.CreateRoleRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	role, err := h.service.CreateRole(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    role,
	})
}

func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid role id")
	}

	role, err := h.service.GetRole(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    role,
	})
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    roles,
	})
}
