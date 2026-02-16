package handlers

import (
	"strconv"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct {
	companyService      services.CompanyService
	businessTypeService services.BusinessTypeService
	locationService     services.LocationService
	taxTypeService      services.TaxTypeService
}

func NewCompanyHandler(
	companyService services.CompanyService,
	businessTypeService services.BusinessTypeService,
	locationService services.LocationService,
	taxTypeService services.TaxTypeService,
) *CompanyHandler {
	return &CompanyHandler{
		companyService:      companyService,
		businessTypeService: businessTypeService,
		locationService:     locationService,
		taxTypeService:      taxTypeService,
	}
}

// ===========================
// COMPANY CRUD
// ===========================

// CreateCompany godoc
// @Summary Create a new company
// @Tags Companies
// @Accept json
// @Produce json
// @Param input body input.CreateCompanyInput true "Company details"
// @Success 201 {object} output.CompleteCompanyProfileOutput
// @Router /companies [post]
func (h *CompanyHandler) CreateCompany(c *fiber.Ctx) error {
	var input input.CreateCompanyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("user_id").(uint)

	output, err := h.companyService.CreateCompany(&input, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

// GetCompany godoc
// @Summary Get company by ID
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompleteCompanyProfileOutput
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetCompany(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}

	return c.JSON(output)
}

// GetAllCompanies godoc
// @Summary Get all companies with filters
// @Tags Companies
// @Produce json
// @Param business_type_id query int false "Filter by business type ID"
// @Param search query string false "Search by company name, GST, or PAN"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} output.CompanyPaginatedResponse
// @Router /companies [get]
func (h *CompanyHandler) GetAllCompanies(c *fiber.Ctx) error {
	query := &output.ListCompaniesQuery{
		Page:     1,
		PageSize: 10,
	}

	if businessTypeID := c.Query("business_type_id"); businessTypeID != "" {
		id, _ := strconv.ParseUint(businessTypeID, 10, 32)
		val := uint(id)
		query.BusinessTypeID = &val
	}
	if search := c.Query("search"); search != "" {
		query.Search = &search
	}
	if isActive := c.Query("is_active"); isActive != "" {
		val := isActive == "true"
		query.IsActive = &val
	}
	if page := c.Query("page"); page != "" {
		query.Page, _ = strconv.Atoi(page)
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		query.PageSize, _ = strconv.Atoi(pageSize)
	}

	response, err := h.companyService.GetAllCompanies(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}

// UpdateCompany godoc
// @Summary Update company
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpdateCompanyInput true "Update details"
// @Success 200 {object} output.CompanyOutput
// @Router /companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpdateCompanyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(uint)

	output, err := h.companyService.UpdateCompany(uint(id), &input, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// DeleteCompany godoc
// @Summary Delete company
// @Tags Companies
// @Param id path int true "Company ID"
// @Success 204
// @Router /companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.companyService.DeleteCompany(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CompleteCompanySetup godoc
// @Summary Complete company setup in one request
// @Tags Companies
// @Accept json
// @Produce json
// @Param input body input.CompleteCompanySetupInput true "Complete company setup"
// @Success 201 {object} output.CompleteCompanyProfileOutput
// @Router /companies/setup [post]
func (h *CompanyHandler) CompleteCompanySetup(c *fiber.Ctx) error {
	var input input.CompleteCompanySetupInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(uint)

	output, err := h.companyService.CompleteCompanySetup(&input, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

// ===========================
// CONTACT OPERATIONS
// ===========================

// UpsertContact godoc
// @Summary Create or update company contact
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpsertCompanyContactInput true "Contact details"
// @Success 200 {object} output.CompanyContactOutput
// @Router /companies/{id}/contact [put]
func (h *CompanyHandler) UpsertContact(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpsertCompanyContactInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpsertContact(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// GetContact godoc
// @Summary Get company contact
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompanyContactOutput
// @Router /companies/{id}/contact [get]
func (h *CompanyHandler) GetContact(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetContact(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Contact not found"})
	}

	return c.JSON(output)
}

// ===========================
// ADDRESS OPERATIONS
// ===========================

// UpsertAddress godoc
// @Summary Create or update company address
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpsertCompanyAddressInput true "Address details"
// @Success 200 {object} output.CompanyAddressOutput
// @Router /companies/{id}/address [put]
func (h *CompanyHandler) UpsertAddress(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpsertCompanyAddressInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpsertAddress(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// GetAddress godoc
// @Summary Get company address
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompanyAddressOutput
// @Router /companies/{id}/address [get]
func (h *CompanyHandler) GetAddress(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetAddress(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Address not found"})
	}

	return c.JSON(output)
}

// ===========================
// BANK DETAILS OPERATIONS
// ===========================

// CreateBankDetail godoc
// @Summary Add bank detail to company
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.CreateBankDetailInput true "Bank details"
// @Success 201 {object} output.CompanyBankDetailOutput
// @Router /companies/{id}/bank-details [post]
func (h *CompanyHandler) CreateBankDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.CreateBankDetailInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.CreateBankDetail(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

// GetBankDetails godoc
// @Summary Get all bank details for a company
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {array} output.CompanyBankDetailOutput
// @Router /companies/{id}/bank-details [get]
func (h *CompanyHandler) GetBankDetails(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetBankDetails(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// UpdateBankDetail godoc
// @Summary Update bank detail
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Bank Detail ID"
// @Param input body input.UpdateBankDetailInput true "Update details"
// @Success 200 {object} output.CompanyBankDetailOutput
// @Router /companies/bank-details/{id} [put]
func (h *CompanyHandler) UpdateBankDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpdateBankDetailInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpdateBankDetail(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// DeleteBankDetail godoc
// @Summary Delete bank detail
// @Tags Companies
// @Param id path int true "Bank Detail ID"
// @Success 204
// @Router /companies/bank-details/{id} [delete]
func (h *CompanyHandler) DeleteBankDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.companyService.DeleteBankDetail(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ===========================
// UPI DETAILS OPERATIONS
// ===========================

// UpsertUPIDetail godoc
// @Summary Create or update company UPI details
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpsertUPIDetailInput true "UPI details"
// @Success 200 {object} output.CompanyUPIDetailOutput
// @Router /companies/{id}/upi-details [put]
func (h *CompanyHandler) UpsertUPIDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpsertUPIDetailInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpsertUPIDetail(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// GetUPIDetail godoc
// @Summary Get company UPI details
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompanyUPIDetailOutput
// @Router /companies/{id}/upi-details [get]
func (h *CompanyHandler) GetUPIDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetUPIDetail(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "UPI details not found"})
	}

	return c.JSON(output)
}

// ===========================
// INVOICE SETTINGS OPERATIONS
// ===========================

// UpsertInvoiceSettings godoc
// @Summary Create or update invoice settings
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpsertInvoiceSettingsInput true "Invoice settings"
// @Success 200 {object} output.CompanyInvoiceSettingsOutput
// @Router /companies/{id}/invoice-settings [put]
func (h *CompanyHandler) UpsertInvoiceSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpsertInvoiceSettingsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpsertInvoiceSettings(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// GetInvoiceSettings godoc
// @Summary Get invoice settings
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompanyInvoiceSettingsOutput
// @Router /companies/{id}/invoice-settings [get]
func (h *CompanyHandler) GetInvoiceSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetInvoiceSettings(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice settings not found"})
	}

	return c.JSON(output)
}

// ===========================
// TAX SETTINGS OPERATIONS
// ===========================

// UpsertTaxSettings godoc
// @Summary Create or update tax settings
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpsertTaxSettingsInput true "Tax settings"
// @Success 200 {object} output.CompanyTaxSettingsOutput
// @Router /companies/{id}/tax-settings [put]
func (h *CompanyHandler) UpsertTaxSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpsertTaxSettingsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpsertTaxSettings(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// GetTaxSettings godoc
// @Summary Get tax settings
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompanyTaxSettingsOutput
// @Router /companies/{id}/tax-settings [get]
func (h *CompanyHandler) GetTaxSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetTaxSettings(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tax settings not found"})
	}

	return c.JSON(output)
}

// ===========================
// REGIONAL SETTINGS OPERATIONS
// ===========================

// UpsertRegionalSettings godoc
// @Summary Create or update regional settings
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param input body input.UpsertRegionalSettingsInput true "Regional settings"
// @Success 200 {object} output.CompanyRegionalSettingsOutput
// @Router /companies/{id}/regional-settings [put]
func (h *CompanyHandler) UpsertRegionalSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input input.UpsertRegionalSettingsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	output, err := h.companyService.UpsertRegionalSettings(uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(output)
}

// GetRegionalSettings godoc
// @Summary Get regional settings
// @Tags Companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} output.CompanyRegionalSettingsOutput
// @Router /companies/{id}/regional-settings [get]
func (h *CompanyHandler) GetRegionalSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	output, err := h.companyService.GetRegionalSettings(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Regional settings not found"})
	}

	return c.JSON(output)
}

// ===========================
// HELPER DATA HANDLERS
// ===========================

type HelperHandler struct {
	businessTypeService services.BusinessTypeService
	locationService     services.LocationService
	taxTypeService      services.TaxTypeService
}

func NewHelperHandler(
	businessTypeService services.BusinessTypeService,
	locationService services.LocationService,
	taxTypeService services.TaxTypeService,
) *HelperHandler {
	return &HelperHandler{
		businessTypeService: businessTypeService,
		locationService:     locationService,
		taxTypeService:      taxTypeService,
	}
}

// GetBusinessTypes godoc
// @Summary Get all business types
// @Tags Helpers
// @Produce json
// @Success 200 {array} output.BusinessTypeOutput
// @Router /helpers/business-types [get]
func (h *HelperHandler) GetBusinessTypes(c *fiber.Ctx) error {
	output, err := h.businessTypeService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}

// GetCountries godoc
// @Summary Get all countries
// @Tags Helpers
// @Produce json
// @Success 200 {array} output.CountryOutput
// @Router /helpers/countries [get]
func (h *HelperHandler) GetCountries(c *fiber.Ctx) error {
	output, err := h.locationService.GetAllCountries()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}

// GetStatesByCountry godoc
// @Summary Get states by country
// @Tags Helpers
// @Produce json
// @Param country_id path int true "Country ID"
// @Success 200 {array} output.StateOutput
// @Router /helpers/countries/{country_id}/states [get]
func (h *HelperHandler) GetStatesByCountry(c *fiber.Ctx) error {
	countryID, err := strconv.ParseUint(c.Params("country_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid country ID"})
	}

	output, err := h.locationService.GetStatesByCountry(uint(countryID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}

// GetTaxTypes godoc
// @Summary Get all tax types
// @Tags Helpers
// @Produce json
// @Success 200 {array} output.TaxTypeOutput
// @Router /helpers/tax-types [get]
func (h *HelperHandler) GetTaxTypes(c *fiber.Ctx) error {
	output, err := h.taxTypeService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}
