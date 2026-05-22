package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductConversionHandler struct {
	service  services.ProductConversionService
	validate *validator.Validate
}

func NewProductConversionHandler(service services.ProductConversionService) *ProductConversionHandler {
	return &ProductConversionHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateConversion creates a new conversion rule
func (h *ProductConversionHandler) CreateConversion(c *fiber.Ctx) error {
	var conversionInput input.CreateProductConversionInput

	if err := c.BodyParser(&conversionInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(conversionInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	userName := ""
	if un := c.Locals("user_name"); un != nil {
		userName = fmt.Sprintf("%v", un)
	}

	companyID := uint(0)
	if cid := c.Locals("company_id"); cid != nil {
		companyID = cid.(uint)
	}

	companyName := ""
	if cn := c.Locals("company_name"); cn != nil {
		companyName = fmt.Sprintf("%v", cn)
	}

	conversion, err := h.service.CreateConversion(&conversionInput, userID, companyID, userName, companyName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(conversion)
}

// UpdateConversion updates a conversion rule
func (h *ProductConversionHandler) UpdateConversion(c *fiber.Ctx) error {
	conversionID := c.Params("id")

	var updateInput input.UpdateProductConversionInput
	if err := c.BodyParser(&updateInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	userName := ""
	if un := c.Locals("user_name"); un != nil {
		userName = fmt.Sprintf("%v", un)
	}

	conversion, err := h.service.UpdateConversion(conversionID, &updateInput, userID, userName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversion)
}

// DeleteConversion deletes a conversion rule
func (h *ProductConversionHandler) DeleteConversion(c *fiber.Ctx) error {
	conversionID := c.Params("id")

	if err := h.service.DeleteConversion(conversionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Conversion deleted successfully",
	})
}

// GetConversion retrieves a conversion rule
func (h *ProductConversionHandler) GetConversion(c *fiber.Ctx) error {
	conversionID := c.Params("id")

	conversion, err := h.service.GetConversion(conversionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if conversion == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Conversion not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversion)
}

// ListConversions lists all conversions
func (h *ProductConversionHandler) ListConversions(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	conversions, err := h.service.ListConversions(offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversions)
}

// ListActiveConversions lists all active conversions
func (h *ProductConversionHandler) ListActiveConversions(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	conversions, err := h.service.ListActiveConversions(offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversions)
}

// GetConversionsByRawProduct lists conversions for a raw product
func (h *ProductConversionHandler) GetConversionsByRawProduct(c *fiber.Ctx) error {
	rawProductID := c.Query("raw_product_id")
	if rawProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "raw_product_id is required",
		})
	}

	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	conversions, err := h.service.GetConversionsByRawProduct(rawProductID, offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversions)
}

// GetConversionsByFinishedProduct lists conversions for a finished product
func (h *ProductConversionHandler) GetConversionsByFinishedProduct(c *fiber.Ctx) error {
	finishedProductID := c.Query("finished_product_id")
	if finishedProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "finished_product_id is required",
		})
	}

	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	conversions, err := h.service.GetConversionsByFinishedProduct(finishedProductID, offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversions)
}

// ExecuteConversion executes a product conversion
func (h *ProductConversionHandler) ExecuteConversion(c *fiber.Ctx) error {
	var conversionRecordInput input.CreateProductConversionRecordInput

	if err := c.BodyParser(&conversionRecordInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(conversionRecordInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		userID = fmt.Sprintf("%v", uid)
	}

	userName := ""
	if un := c.Locals("user_name"); un != nil {
		userName = fmt.Sprintf("%v", un)
	}

	companyID := uint(0)
	if cid := c.Locals("company_id"); cid != nil {
		companyID = cid.(uint)
	}

	companyName := ""
	if cn := c.Locals("company_name"); cn != nil {
		companyName = fmt.Sprintf("%v", cn)
	}

	result, err := h.service.ExecuteConversion(&conversionRecordInput, userID, companyID, userName, companyName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetConversionRecord retrieves a conversion record
func (h *ProductConversionHandler) GetConversionRecord(c *fiber.Ctx) error {
	recordID := c.Params("record_id")

	record, err := h.service.GetConversionRecord(recordID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if record == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Record not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(record)
}

// ListConversionRecords lists all conversion records
func (h *ProductConversionHandler) ListConversionRecords(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	records, err := h.service.ListConversionRecords(offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(records)
}

// ListConversionRecordsByRule lists conversion records for a specific rule
func (h *ProductConversionHandler) ListConversionRecordsByRule(c *fiber.Ctx) error {
	conversionID := c.Query("conversion_id")
	if conversionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "conversion_id is required",
		})
	}

	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	records, err := h.service.ListConversionRecordsByRule(conversionID, offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(records)
}

// ListConversionRecordsByDateRange lists conversion records within a date range
func (h *ProductConversionHandler) ListConversionRecordsByDateRange(c *fiber.Ctx) error {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "from_date and to_date are required",
		})
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid from_date format (use YYYY-MM-DD)",
		})
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid to_date format (use YYYY-MM-DD)",
		})
	}

	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	records, err := h.service.ListConversionRecordsByDateRange(fromDate, toDate.AddDate(0, 0, 1), offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(records)
}
