package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type ProductGroupHandler struct {
	service services.ProductGroupService
}

func NewProductGroupHandler(service services.ProductGroupService) *ProductGroupHandler {
	return &ProductGroupHandler{service: service}
}

func productGroupContext(c *fiber.Ctx) (uint, uint, error) {
	userID, err := productGroupLocalToUint(c.Locals("user_id"))
	if err != nil || userID == 0 {
		return 0, 0, fiber.NewError(fiber.StatusUnauthorized, "invalid authenticated user")
	}
	companyID, err := productGroupLocalToUint(c.Locals("company_id"))
	if err != nil || companyID == 0 {
		return 0, 0, fiber.NewError(fiber.StatusForbidden, "invalid authenticated company")
	}
	return userID, companyID, nil
}

func productGroupPagination(c *fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (h *ProductGroupHandler) CreateProductGroup(c *fiber.Ctx) error {
	var req input.CreateProductGroupInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	userID, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	productGroup, err := h.service.CreateForCompany(&req, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    productGroup,
		"message": "Product Group created successfully",
	})
}

func (h *ProductGroupHandler) GetProductGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "product group id is required")
	}
	_, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	productGroup, err := h.service.FindByIDForCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": productGroup})
}

func (h *ProductGroupHandler) GetAllProductGroups(c *fiber.Ctx) error {
	limit, offset := productGroupPagination(c)
	search := c.Query("search", "")
	_, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	result, err := h.service.FindAllForCompany(companyID, limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    result.ProductGroups,
		"total":   result.Total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *ProductGroupHandler) UpdateProductGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "product group id is required")
	}
	var req input.UpdateProductGroupInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	userID, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	productGroup, err := h.service.UpdateForCompany(id, &req, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": productGroup, "message": "Product Group updated successfully"})
}

func (h *ProductGroupHandler) DeleteProductGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "product group id is required")
	}
	userID, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	if err := h.service.DeleteForCompany(id, userID, companyID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Product Group deleted successfully"})
}

func (h *ProductGroupHandler) GetProductGroupByName(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name query parameter is required")
	}
	_, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	productGroup, err := h.service.FindByNameForCompany(name, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": productGroup})
}

func (h *ProductGroupHandler) ReorderProductGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "product group id is required")
	}
	var req input.UpdateProductGroupInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	userID, companyID, err := productGroupContext(c)
	if err != nil {
		return err
	}
	response, err := h.service.ReorderWithSummaryForCompany(id, &req, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"code":    200,
		"message": "Product group reordered successfully",
		"data":    response,
	})
}

func productGroupLocalToUint(value interface{}) (uint, error) {
	if value == nil {
		return 0, fmt.Errorf("value is missing")
	}
	switch v := value.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case int8:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case int16:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case int32:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case int64:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case float32:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case float64:
		if v <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(v), nil
	case json.Number:
		parsed, err := strconv.ParseUint(v.String(), 10, 64)
		return uint(parsed), err
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		return uint(parsed), err
	default:
		parsed, err := strconv.ParseUint(fmt.Sprintf("%v", v), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unsupported numeric type %T", value)
		}
		return uint(parsed), nil
	}
}
