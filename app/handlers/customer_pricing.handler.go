package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type CustomerPricingHandler struct {
	service services.CustomerPricingService
}

func NewCustomerPricingHandler(service services.CustomerPricingService) *CustomerPricingHandler {
	return &CustomerPricingHandler{service: service}
}

func customerPricingLocalToUint(value interface{}) (uint, error) {
	if value == nil {
		return 0, fmt.Errorf("value is missing")
	}

	switch typedValue := value.(type) {
	case uint:
		return typedValue, nil
	case uint8:
		return uint(typedValue), nil
	case uint16:
		return uint(typedValue), nil
	case uint32:
		return uint(typedValue), nil
	case uint64:
		return uint(typedValue), nil
	case int:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case int8:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case int16:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case int32:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case int64:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case float32:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case float64:
		if typedValue <= 0 {
			return 0, fmt.Errorf("value must be greater than zero")
		}
		return uint(typedValue), nil
	case json.Number:
		parsed, err := strconv.ParseUint(typedValue.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(parsed), nil
	case string:
		parsed, err := strconv.ParseUint(typedValue, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(parsed), nil
	default:
		parsed, err := strconv.ParseUint(fmt.Sprintf("%v", typedValue), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unsupported numeric type %T with value %v", value, value)
		}
		return uint(parsed), nil
	}
}

func customerPricingContext(c *fiber.Ctx) (uint, uint, error) {
	rawUserID := c.Locals("user_id")
	rawCompanyID := c.Locals("company_id")

	userID, err := customerPricingLocalToUint(rawUserID)
	if err != nil || userID == 0 {
		return 0, 0, fiber.NewError(
			fiber.StatusUnauthorized,
			fmt.Sprintf("invalid authenticated user: value=%v type=%T", rawUserID, rawUserID),
		)
	}

	companyID, err := customerPricingLocalToUint(rawCompanyID)
	if err != nil || companyID == 0 {
		return 0, 0, fiber.NewError(
			fiber.StatusForbidden,
			fmt.Sprintf("invalid authenticated company: value=%v type=%T", rawCompanyID, rawCompanyID),
		)
	}

	return userID, companyID, nil
}

func customerPricingPagination(c *fiber.Ctx) (int, int) {
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return offset, limit
}

func (h *CustomerPricingHandler) CreateCustomerPricing(c *fiber.Ctx) error {
	var req input.CreateCustomerPricingDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body", "details": err.Error()})
	}

	userID, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	pricings, err := h.service.CreatePricingForCompany(req.CustomerID, req.LineItems, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": pricings, "count": len(pricings), "message": "customer pricing created successfully"})
}

func (h *CustomerPricingHandler) UpdateCustomerPricing(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "pricing ID is required"})
	}

	var req input.UpdateCustomerPricingDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body", "details": err.Error()})
	}

	userID, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	pricing, err := h.service.UpdatePricingForCompany(id, req.Rate, req.Account, req.Description, req.IsActive, userID, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": pricing})
}

func (h *CustomerPricingHandler) DeleteCustomerPricing(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "pricing ID is required"})
	}

	userID, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	if err := h.service.DeletePricingForCompany(id, userID, companyID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "customer pricing deleted successfully"})
}

func (h *CustomerPricingHandler) GetCustomerPricingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "pricing ID is required"})
	}

	_, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	pricing, err := h.service.GetPricingByIDForCompany(id, companyID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "pricing not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": pricing})
}

func (h *CustomerPricingHandler) GetPricingByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Query("customer_id"), 10, 32)
	if err != nil || customerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid customer ID"})
	}

	offset, limit := customerPricingPagination(c)
	_, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	pricings, total, err := h.service.GetPricingByCustomerForCompany(uint(customerID), companyID, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": fiber.Map{"pricings": pricings, "total_count": total}})
}

func (h *CustomerPricingHandler) GetAllCustomerPricing(c *fiber.Ctx) error {
	offset, limit := customerPricingPagination(c)
	_, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	pricings, total, err := h.service.GetAllPricingForCompany(companyID, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": fiber.Map{"pricings": pricings, "total_count": total}})
}

func (h *CustomerPricingHandler) GetActivePricingByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseUint(c.Query("customer_id"), 10, 32)
	if err != nil || customerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid customer ID"})
	}

	_, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	pricings, err := h.service.GetActivePricingByCustomerForCompany(uint(customerID), companyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": fiber.Map{"pricings": pricings}})
}

func (h *CustomerPricingHandler) SetEffectiveDateRange(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "pricing ID is required"})
	}

	var req input.SetDateRangeDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body", "details": err.Error()})
	}

	userID, companyID, err := customerPricingContext(c)
	if err != nil {
		return err
	}

	if err := h.service.SetEffectiveDateRangeForCompany(id, req.EffectiveFrom, req.EffectiveTo, userID, companyID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "date range updated successfully"})
}
