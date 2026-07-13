package handlers

import (
	"fmt"
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type RawMaterialBagHandler struct {
	service services.RawMaterialBagService
}

func NewRawMaterialBagHandler(
	service services.RawMaterialBagService,
) *RawMaterialBagHandler {
	return &RawMaterialBagHandler{
		service: service,
	}
}

func rawMaterialBagLocalUint(
	c *fiber.Ctx,
	key string,
) uint {
	value := c.Locals(key)

	switch typed := value.(type) {
	case uint:
		return typed
	case uint64:
		return uint(typed)
	case int:
		if typed > 0 {
			return uint(typed)
		}
	case int64:
		if typed > 0 {
			return uint(typed)
		}
	case float64:
		if typed > 0 {
			return uint(typed)
		}
	case string:
		parsed, err := strconv.ParseUint(typed, 10, 64)
		if err == nil {
			return uint(parsed)
		}
	}

	return 0
}

func rawMaterialBagAuthContext(
	c *fiber.Ctx,
) (
	userID string,
	companyID uint,
	err error,
) {
	authenticatedUserID := rawMaterialBagLocalUint(
		c,
		"user_id",
	)
	companyID = rawMaterialBagLocalUint(
		c,
		"company_id",
	)

	if authenticatedUserID == 0 {
		return "", 0, fmt.Errorf(
			"invalid authenticated user",
		)
	}

	if companyID == 0 {
		return "", 0, fmt.Errorf(
			"user is not assigned to a company",
		)
	}

	return strconv.FormatUint(
		uint64(authenticatedUserID),
		10,
	), companyID, nil
}

func rawMaterialBagPagination(
	c *fiber.Ctx,
) (
	limit int,
	offset int,
) {
	limit, err := strconv.Atoi(
		c.Query("limit", "10"),
	)
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset, err = strconv.Atoi(
		c.Query("offset", "0"),
	)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

func rawMaterialBagContextError(
	c *fiber.Ctx,
	err error,
) error {
	return c.Status(fiber.StatusForbidden).JSON(
		fiber.Map{
			"success": false,
			"message": err.Error(),
		},
	)
}

func (h *RawMaterialBagHandler) ReceiveBags(
	c *fiber.Ctx,
) error {
	var req input.ReceiveRawMaterialBagsInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	userID, companyID, err :=
		rawMaterialBagAuthContext(c)
	if err != nil {
		return rawMaterialBagContextError(c, err)
	}

	response, err :=
		h.service.ReceiveBagsForCompany(
			&req,
			userID,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"success": true,
			"message": "Raw material bags received successfully",
			"data":    response,
		},
	)
}

func (h *RawMaterialBagHandler) GetByProduct(
	c *fiber.Ctx,
) error {
	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Product ID is required",
			},
		)
	}

	_, companyID, err :=
		rawMaterialBagAuthContext(c)
	if err != nil {
		return rawMaterialBagContextError(c, err)
	}

	response, err :=
		h.service.GetBagsByProductForCompany(
			productID,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data":    response,
		},
	)
}

func (h *RawMaterialBagHandler) GetByPurchaseOrder(
	c *fiber.Ctx,
) error {
	purchaseOrderID := c.Params("po_id")
	if purchaseOrderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Purchase order ID is required",
			},
		)
	}

	_, companyID, err :=
		rawMaterialBagAuthContext(c)
	if err != nil {
		return rawMaterialBagContextError(c, err)
	}

	response, err :=
		h.service.GetBagsByPurchaseOrderForCompany(
			purchaseOrderID,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data":    response,
		},
	)
}

func (h *RawMaterialBagHandler) GetAll(
	c *fiber.Ctx,
) error {
	limit, offset := rawMaterialBagPagination(c)

	_, companyID, err :=
		rawMaterialBagAuthContext(c)
	if err != nil {
		return rawMaterialBagContextError(c, err)
	}

	response, err :=
		h.service.GetAllForCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(
				fiber.Map{
					"success": false,
					"message": err.Error(),
				},
			)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data":    response,
		},
	)
}

func (h *RawMaterialBagHandler) GetByID(
	c *fiber.Ctx,
) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Bag ID is required",
			},
		)
	}

	_, companyID, err :=
		rawMaterialBagAuthContext(c)
	if err != nil {
		return rawMaterialBagContextError(c, err)
	}

	response, err :=
		h.service.GetByIDForCompany(
			id,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "bag not found",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"data":    response,
		},
	)
}

// UseBags consumes raw material from one or more bags.
//
// POST /raw-material-bags/use/:product_id
func (h *RawMaterialBagHandler) UseBags(
	c *fiber.Ctx,
) error {
	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Product ID is required",
			},
		)
	}

	var req struct {
		Bags []input.UseRawMaterialBagInput `json:"bags"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid request body",
			},
		)
	}

	if len(req.Bags) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "At least one bag is required",
			},
		)
	}

	_, companyID, err :=
		rawMaterialBagAuthContext(c)
	if err != nil {
		return rawMaterialBagContextError(c, err)
	}

	totalUsedKG, err :=
		h.service.UseBagsForCompany(
			productID,
			req.Bags,
			companyID,
		)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success":       true,
			"message":       "Raw material bags used successfully",
			"total_used_kg": totalUsedKG,
		},
	)
}
