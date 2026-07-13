package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserInfo struct {
	UserID      string
	UserName    string
	CompanyID   uint
	CompanyName string
}

type BillHandler struct {
	service  services.BillService
	validate *validator.Validate
}

func NewBillHandler(service services.BillService) *BillHandler {
	return &BillHandler{
		service:  service,
		validate: validator.New(),
	}
}

func billLocalUint(c *fiber.Ctx, key string) uint {
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

func extractUserInfo(c *fiber.Ctx) UserInfo {
	userInfo := UserInfo{
		UserID:    "",
		UserName:  "",
		CompanyID: billLocalUint(c, "company_id"),
	}

	if userID := billLocalUint(c, "user_id"); userID > 0 {
		userInfo.UserID = strconv.FormatUint(uint64(userID), 10)
	}

	if email := c.Locals("user_email"); email != nil {
		userInfo.UserName = fmt.Sprintf("%v", email)
	}

	return userInfo
}

func billAuthContext(c *fiber.Ctx) (UserInfo, error) {
	userInfo := extractUserInfo(c)

	if userInfo.UserID == "" {
		return UserInfo{}, fmt.Errorf("invalid authenticated user")
	}

	if userInfo.CompanyID == 0 {
		return UserInfo{}, fmt.Errorf("user is not assigned to a company")
	}

	return userInfo, nil
}

func billContextError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   err.Error(),
		"success": false,
	})
}

func billPagination(c *fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
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

func (h *BillHandler) CreateBill(c *fiber.Ctx) error {
	var req input.CreateBillInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bill, err := h.service.CreateBillForCompany(
		&req,
		userInfo.UserID,
		userInfo.CompanyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    bill,
		"message": "Bill created successfully",
		"success": true,
	})
}

func (h *BillHandler) GetBill(c *fiber.Ctx) error {
	id := c.Params("id")

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bill, err := h.service.GetBillByCompany(
		id,
		userInfo.CompanyID,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Bill not found",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    bill,
		"success": true,
	})
}

func (h *BillHandler) GetAllBills(c *fiber.Ctx) error {
	limit, offset := billPagination(c)

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bills, total, err := h.service.GetAllBillsByCompany(
		userInfo.CompanyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get bills",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    bills,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"success": true,
	})
}

func (h *BillHandler) GetBillsByVendor(c *fiber.Ctx) error {
	vendorID64, err := strconv.ParseUint(
		c.Params("vendorId"),
		10,
		32,
	)
	if err != nil || vendorID64 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid vendor ID",
			"success": false,
		})
	}

	limit, offset := billPagination(c)

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bills, total, err := h.service.GetBillsByVendorAndCompany(
		uint(vendorID64),
		userInfo.CompanyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    bills,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"success": true,
	})
}

func (h *BillHandler) GetBillsByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	limit, offset := billPagination(c)

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bills, total, err := h.service.GetBillsByStatusAndCompany(
		status,
		userInfo.CompanyID,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get bills by status",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    bills,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"success": true,
	})
}

func (h *BillHandler) UpdateBill(c *fiber.Ctx) error {
	id := c.Params("id")

	var req input.UpdateBillInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bill, err := h.service.UpdateBillForCompany(
		id,
		&req,
		userInfo.UserID,
		userInfo.CompanyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    bill,
		"message": "Bill updated successfully",
		"success": true,
	})
}

func (h *BillHandler) UpdateBillStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var req input.UpdateBillStatusInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	bill, err := h.service.UpdateBillStatusForCompany(
		id,
		req.Status,
		userInfo.UserID,
		userInfo.CompanyID,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    bill,
		"message": "Bill status updated successfully",
		"success": true,
	})
}

func (h *BillHandler) DeleteBill(c *fiber.Ctx) error {
	id := c.Params("id")

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	if err := h.service.DeleteBillForCompany(
		id,
		userInfo.CompanyID,
	); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Bill deleted successfully",
		"success": true,
	})
}

// This endpoint only builds a preview response.
// Company validation is still applied to vendor and products.
func (h *BillHandler) CreateBillFromVariants(c *fiber.Ctx) error {
	var req input.CreateBillInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userInfo, err := billAuthContext(c)
	if err != nil {
		return billContextError(c, err)
	}

	if err := h.service.ValidateBillInputForCompany(
		&req,
		userInfo.CompanyID,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	billNumber := req.BillNumber
	if billNumber == "" {
		billNumber = fmt.Sprintf(
			"BILL-%d-%d",
			time.Now().Year(),
			time.Now().UnixNano(),
		)
	}

	lineItems := make(
		[]output.BillLineItemOutput,
		len(req.LineItems),
	)

	subTotal := 0.0

	for index, item := range req.LineItems {
		amount := item.Quantity * item.Rate
		subTotal += amount

		lineItems[index] = output.BillLineItemOutput{
			VariantSKU:  &item.SKU,
			Description: item.ProductName,
			Account:     item.Account,
			Quantity:    item.Quantity,
			Rate:        item.Rate,
			Amount:      amount,
		}
	}

	total := subTotal - req.Discount + req.Adjustment

	billOutput := output.CreateBillVariantOutput(
		billNumber,
		req.PurchaseOrderID,
		req.VendorID,
		"",
		lineItems,
		subTotal,
		0,
		total,
	)

	billOutput.UserName = userInfo.UserID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    billOutput,
		"message": "Bill created successfully from purchase order variants",
	})
}