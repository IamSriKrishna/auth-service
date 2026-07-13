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

func productConversionLocalUint(c *fiber.Ctx, key string) uint {
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

func productConversionAuthContext(c *fiber.Ctx) (string, uint, string, string, error) {
	userID := productConversionLocalUint(c, "user_id")
	companyID := productConversionLocalUint(c, "company_id")

	if userID == 0 {
		return "", 0, "", "", fmt.Errorf("invalid authenticated user")
	}
	if companyID == 0 {
		return "", 0, "", "", fmt.Errorf("user is not assigned to a company")
	}

	userName := ""
	if value := c.Locals("user_name"); value != nil {
		userName = fmt.Sprintf("%v", value)
	} else if value := c.Locals("user_email"); value != nil {
		userName = fmt.Sprintf("%v", value)
	}

	companyName := ""
	if value := c.Locals("company_name"); value != nil {
		companyName = fmt.Sprintf("%v", value)
	}

	return strconv.FormatUint(uint64(userID), 10), companyID, userName, companyName, nil
}

func productConversionContextError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
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

	userID, companyID, userName, companyName, err := productConversionAuthContext(c)
	if err != nil {
		return productConversionContextError(c, err)
	}

	conversion, err := h.service.CreateConversionForCompany(&conversionInput, userID, companyID, userName, companyName)
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

	userID, companyID, userName, _, err := productConversionAuthContext(c)
	if err != nil {
		return productConversionContextError(c, err)
	}

	conversion, err := h.service.UpdateConversionForCompany(conversionID, &updateInput, userID, userName, companyID)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	if err := h.service.DeleteConversionForCompany(conversionID, companyID); err != nil {
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	conversion, err := h.service.GetConversionForCompany(conversionID, companyID)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	conversions, err := h.service.ListConversionsForCompany(companyID, offset, limitNum)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	conversions, err := h.service.ListActiveConversionsForCompany(companyID, offset, limitNum)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	conversions, err := h.service.GetConversionsByRawProductForCompany(rawProductID, companyID, offset, limitNum)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	conversions, err := h.service.GetConversionsByFinishedProductForCompany(finishedProductID, companyID, offset, limitNum)
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

	userID, companyID, userName, companyName, err := productConversionAuthContext(c)
	if err != nil {
		return productConversionContextError(c, err)
	}

	result, err := h.service.ExecuteConversionForCompany(&conversionRecordInput, userID, companyID, userName, companyName)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	record, err := h.service.GetConversionRecordForCompany(recordID, companyID)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	records, err := h.service.ListConversionRecordsForCompany(companyID, offset, limitNum)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	records, err := h.service.ListConversionRecordsByRuleForCompany(conversionID, companyID, offset, limitNum)
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

	_, companyID, _, _, authErr := productConversionAuthContext(c)
	if authErr != nil {
		return productConversionContextError(c, authErr)
	}

	records, err := h.service.ListConversionRecordsByDateRangeForCompany(fromDate, toDate.AddDate(0, 0, 1), companyID, offset, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(records)
}