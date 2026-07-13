package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type ManufacturerHandler struct {
	service services.ManufacturerService
}

func NewManufacturerHandler(
	service services.ManufacturerService,
) *ManufacturerHandler {
	return &ManufacturerHandler{
		service: service,
	}
}

func manufacturerContext(
	c *fiber.Ctx,
) (uint, uint, error) {
	rawUserID := c.Locals("user_id")
	rawCompanyID := c.Locals("company_id")

	userID, err := localToUint(rawUserID)
	if err != nil || userID == 0 {
		return 0, 0, fiber.NewError(
			fiber.StatusUnauthorized,
			fmt.Sprintf(
				"invalid authenticated user: value=%v type=%T",
				rawUserID,
				rawUserID,
			),
		)
	}

	companyID, err := localToUint(rawCompanyID)
	if err != nil || companyID == 0 {
		return 0, 0, fiber.NewError(
			fiber.StatusForbidden,
			fmt.Sprintf(
				"invalid authenticated company: value=%v type=%T",
				rawCompanyID,
				rawCompanyID,
			),
		)
	}

	return userID, companyID, nil
}

func manufacturerPagination(
	c *fiber.Ctx,
) (int, int) {
	limit, err := strconv.Atoi(
		c.Query("limit", "10"),
	)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(
		c.Query("offset", "0"),
	)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

// CreateManufacturer creates a manufacturer inside the JWT company.
func (h *ManufacturerHandler) CreateManufacturer(
	c *fiber.Ctx,
) error {
	var req input.CreateManufacturerInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	userID, companyID, err := manufacturerContext(c)
	if err != nil {
		return err
	}

	manufacturer, err :=
		h.service.CreateForCompany(
			&req,
			userID,
			companyID,
		)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"success": true,
			"data":    manufacturer,
		},
	)
}

// UpdateManufacturer updates a manufacturer only in the JWT company.
func (h *ManufacturerHandler) UpdateManufacturer(
	c *fiber.Ctx,
) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid manufacturer id",
		)
	}

	var req input.UpdateManufacturerInput
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	userID, companyID, err := manufacturerContext(c)
	if err != nil {
		return err
	}

	manufacturer, err :=
		h.service.UpdateForCompany(
			id,
			&req,
			userID,
			companyID,
		)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			err.Error(),
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    manufacturer,
		},
	)
}

// GetManufacturerByID retrieves a manufacturer only from the JWT company.
func (h *ManufacturerHandler) GetManufacturerByID(
	c *fiber.Ctx,
) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid manufacturer id",
		)
	}

	_, companyID, err := manufacturerContext(c)
	if err != nil {
		return err
	}

	manufacturer, err :=
		h.service.GetByIDForCompany(
			id,
			companyID,
		)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			err.Error(),
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    manufacturer,
		},
	)
}

// GetAllManufacturers retrieves manufacturers only from the JWT company.
func (h *ManufacturerHandler) GetAllManufacturers(
	c *fiber.Ctx,
) error {
	limit, offset := manufacturerPagination(c)

	_, companyID, err := manufacturerContext(c)
	if err != nil {
		return err
	}

	manufacturers, err :=
		h.service.GetAllForCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    manufacturers,
		},
	)
}

// GetManufacturersByProductGroup retrieves matching company manufacturers.
func (h *ManufacturerHandler) GetManufacturersByProductGroup(
	c *fiber.Ctx,
) error {
	productGroupID := c.Params("product_group_id")
	if productGroupID == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"product group id is required",
		)
	}

	_, companyID, err := manufacturerContext(c)
	if err != nil {
		return err
	}

	manufacturers, err :=
		h.service.GetByProductGroupIDForCompany(
			productGroupID,
			companyID,
		)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    manufacturers,
		},
	)
}

// DeleteManufacturer deletes only a manufacturer in the JWT company.
func (h *ManufacturerHandler) DeleteManufacturer(
	c *fiber.Ctx,
) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid manufacturer id",
		)
	}

	userID, companyID, err := manufacturerContext(c)
	if err != nil {
		return err
	}

	if err := h.service.DeleteForCompany(
		id,
		userID,
		companyID,
	); err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			err.Error(),
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "manufacturer deleted successfully",
		},
	)
}

func localToUint(value interface{}) (uint, error) {
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
		parsedValue, err := strconv.ParseUint(
			typedValue.String(),
			10,
			64,
		)
		if err != nil {
			return 0, err
		}
		return uint(parsedValue), nil

	case string:
		parsedValue, err := strconv.ParseUint(
			typedValue,
			10,
			64,
		)
		if err != nil {
			return 0, err
		}
		return uint(parsedValue), nil

	default:
		parsedValue, err := strconv.ParseUint(
			fmt.Sprintf("%v", typedValue),
			10,
			64,
		)
		if err != nil {
			return 0, fmt.Errorf(
				"unsupported numeric type %T with value %v",
				value,
				value,
			)
		}

		return uint(parsedValue), nil
	}
}
