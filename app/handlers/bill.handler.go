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
	service services.BillService
}

// Helper function to extract user information from context
func extractUserInfo(c *fiber.Ctx) UserInfo {
	userInfo := UserInfo{}

	if uid := c.Locals("user_id"); uid != nil {
		userInfo.UserID = fmt.Sprintf("%v", uid)
	}

	if email := c.Locals("user_email"); email != nil {
		userInfo.UserName = fmt.Sprintf("%v", email)
	}

	if compID := c.Locals("company_id"); compID != nil {
		if cid, ok := compID.(uint); ok {
			userInfo.CompanyID = cid
		}
	}

	return userInfo
}

func NewBillHandler(service services.BillService) *BillHandler {
	return &BillHandler{service: service}
}

func (h *BillHandler) CreateBill(c *fiber.Ctx) error {
	var billInput input.CreateBillInput

	if err := c.BodyParser(&billInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(billInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	bill, err := h.service.CreateBill(&billInput, userID)
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

	bill, err := h.service.GetBill(id)
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
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	bills, total, err := h.service.GetAllBills(limit, offset, createdBy)
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
	vendorID := c.Params("vendorId")
	vendorIDUint := uint(0)

	if vid, err := strconv.ParseUint(vendorID, 10, 32); err == nil {
		vendorIDUint = uint(vid)
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	bills, total, err := h.service.GetBillsByVendor(vendorIDUint, limit, offset, createdBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get bills for vendor",
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

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	createdBy := ""
	if uid := c.Locals("user_id"); uid != nil {
		createdBy = fmt.Sprintf("%v", uid)
	}

	bills, total, err := h.service.GetBillsByStatus(status, limit, offset, createdBy)
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
	var billInput input.UpdateBillInput

	if err := c.BodyParser(&billInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(billInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	bill, err := h.service.UpdateBill(id, &billInput, userID)
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
	var statusInput input.UpdateBillStatusInput

	if err := c.BodyParser(&statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(statusInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	bill, err := h.service.UpdateBillStatus(id, statusInput.Status, userID)
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

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	err := h.service.DeleteBill(id, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to delete bill",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Bill deleted successfully",
		"success": true,
	})
}

// CreateBillFromVariants creates a bill from purchase order variants
func (h *BillHandler) CreateBillFromVariants(c *fiber.Ctx) error {
	var billInput input.CreateBillInput

	if err := c.BodyParser(&billInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(billInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	// Generate bill number (or use provided one)
	billNo := billInput.BillNumber
	if billNo == "" {
		importTime := time.Now()
		billNo = fmt.Sprintf("BILL-%d-%05d", importTime.Year(), c.Locals("request_id"))
	}

	// Build line items
	var subTotal float64
	lineItems := make([]output.BillLineItemOutput, len(billInput.LineItems))

	for i, item := range billInput.LineItems {
		amount := item.Quantity * item.Rate
		subTotal += amount

		lineItems[i] = output.BillLineItemOutput{
			VariantSKU:  &item.SKU,
			Description: item.ProductName,
			Account:     item.Account,
			Quantity:    item.Quantity,
			Rate:        item.Rate,
			Amount:      amount,
		}
	}

	// Calculate totals
	total := subTotal + billInput.Discount + billInput.Adjustment

	// Create bill output
	billOutput := output.CreateBillVariantOutput(
		billNo,
		billInput.PurchaseOrderID,
		billInput.VendorID,
		"", // Fetch vendor name from DB
		lineItems,
		subTotal,
		0, // totalTax
		total,
	)

	billOutput.UserName = userID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    billOutput,
		"message": "Bill created successfully from purchase order variants",
	})
}
