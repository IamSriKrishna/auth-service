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


func (h *CompanyHandler) CreateCompany(c *fiber.Ctx) error {
	var input input.CreateCompanyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(uint)

	output, err := h.companyService.CreateCompany(&input, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

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

func (h *HelperHandler) GetBusinessTypes(c *fiber.Ctx) error {
	output, err := h.businessTypeService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}

func (h *HelperHandler) GetCountries(c *fiber.Ctx) error {
	output, err := h.locationService.GetAllCountries()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}

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

func (h *HelperHandler) GetTaxTypes(c *fiber.Ctx) error {
	output, err := h.taxTypeService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(output)
}
