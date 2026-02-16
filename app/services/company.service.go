package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"gorm.io/gorm"
)

type CompanyService interface {
	CreateCompany(input *input.CreateCompanyInput, userID uint) (*output.CompleteCompanyProfileOutput, error)
	GetCompany(id uint) (*output.CompleteCompanyProfileOutput, error)
	GetAllCompanies(query *output.ListCompaniesQuery) (*output.CompanyPaginatedResponse, error)
	UpdateCompany(id uint, input *input.UpdateCompanyInput, userID uint) (*output.CompanyOutput, error)
	DeleteCompany(id uint) error
	CompleteCompanySetup(input *input.CompleteCompanySetupInput, userID uint) (*output.CompleteCompanyProfileOutput, error)
	UpsertContact(companyID uint, input *input.UpsertCompanyContactInput) (*output.CompanyContactOutput, error)
	GetContact(companyID uint) (*output.CompanyContactOutput, error)
	UpsertAddress(companyID uint, input *input.UpsertCompanyAddressInput) (*output.CompanyAddressOutput, error)
	GetAddress(companyID uint) (*output.CompanyAddressOutput, error)
	CreateBankDetail(companyID uint, input *input.CreateBankDetailInput) (*output.CompanyBankDetailOutput, error)
	GetBankDetails(companyID uint) ([]output.CompanyBankDetailOutput, error)
	UpdateBankDetail(id uint, input *input.UpdateBankDetailInput) (*output.CompanyBankDetailOutput, error)
	DeleteBankDetail(id uint) error
	UpsertUPIDetail(companyID uint, input *input.UpsertUPIDetailInput) (*output.CompanyUPIDetailOutput, error)
	GetUPIDetail(companyID uint) (*output.CompanyUPIDetailOutput, error)
	UpsertInvoiceSettings(companyID uint, input *input.UpsertInvoiceSettingsInput) (*output.CompanyInvoiceSettingsOutput, error)
	GetInvoiceSettings(companyID uint) (*output.CompanyInvoiceSettingsOutput, error)
	UpsertTaxSettings(companyID uint, input *input.UpsertTaxSettingsInput) (*output.CompanyTaxSettingsOutput, error)
	GetTaxSettings(companyID uint) (*output.CompanyTaxSettingsOutput, error)
	UpsertRegionalSettings(companyID uint, input *input.UpsertRegionalSettingsInput) (*output.CompanyRegionalSettingsOutput, error)
	GetRegionalSettings(companyID uint) (*output.CompanyRegionalSettingsOutput, error)
}

type companyService struct {
	companyRepo      repo.CompanyRepository
	businessTypeRepo repo.BusinessTypeRepository
	locationRepo     repo.LocationRepository
	taxTypeRepo      repo.TaxTypeRepository
	db               *gorm.DB
}

func NewCompanyService(
	companyRepo repo.CompanyRepository,
	businessTypeRepo repo.BusinessTypeRepository,
	locationRepo repo.LocationRepository,
	taxTypeRepo repo.TaxTypeRepository,
	db *gorm.DB,
) CompanyService {
	return &companyService{
		companyRepo:      companyRepo,
		businessTypeRepo: businessTypeRepo,
		locationRepo:     locationRepo,
		taxTypeRepo:      taxTypeRepo,
		db:               db,
	}
}

func (s *companyService) CreateCompany(input *input.CreateCompanyInput, userID uint) (*output.CompleteCompanyProfileOutput, error) {

	_, err := s.businessTypeRepo.FindByID(input.BusinessTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("business type with ID %d does not exist", input.BusinessTypeID)
		}
		return nil, fmt.Errorf("failed to validate business type: %v", err)
	}

	company := &models.Company{
		CompanyName:    input.CompanyName,
		BusinessTypeID: input.BusinessTypeID,
		GSTNumber:      input.GSTNumber,
		PANNumber:      input.PANNumber,
		CreatedBy:      &userID,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, fmt.Errorf("failed to create company: %v", err)
	}

	return s.GetCompany(company.ID)
}

func (s *companyService) GetCompany(id uint) (*output.CompleteCompanyProfileOutput, error) {
	company, err := s.companyRepo.GetCompleteProfile(id)
	if err != nil {
		return nil, err
	}
	return s.toCompleteProfileOutput(company), nil
}

func (s *companyService) GetAllCompanies(query *output.ListCompaniesQuery) (*output.CompanyPaginatedResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	companies, total, err := s.companyRepo.FindAll(
		query.BusinessTypeID,
		query.Search,
		query.IsActive,
		query.Page,
		query.PageSize,
	)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.CompleteCompanyProfileOutput, len(companies))
	for i, c := range companies {
		fullCompany, _ := s.companyRepo.GetCompleteProfile(c.ID)
		if fullCompany != nil {
			outputs[i] = *s.toCompleteProfileOutput(fullCompany)
		}
	}

	return &output.CompanyPaginatedResponse{
		Data:       outputs,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalCount: total,
		TotalPages: int(math.Ceil(float64(total) / float64(query.PageSize))),
	}, nil
}

func (s *companyService) UpdateCompany(id uint, input *input.UpdateCompanyInput, userID uint) (*output.CompanyOutput, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.BusinessTypeID != nil {
		_, err := s.businessTypeRepo.FindByID(*input.BusinessTypeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("business type with ID %d does not exist", *input.BusinessTypeID)
			}
			return nil, fmt.Errorf("failed to validate business type: %v", err)
		}
		company.BusinessTypeID = *input.BusinessTypeID
	}

	if input.CompanyName != nil {
		company.CompanyName = *input.CompanyName
	}
	if input.GSTNumber != nil {
		company.GSTNumber = *input.GSTNumber
	}
	if input.PANNumber != nil {
		company.PANNumber = *input.PANNumber
	}

	company.UpdatedBy = &userID

	if err := s.companyRepo.Update(company); err != nil {
		return nil, fmt.Errorf("failed to update company: %v", err)
	}

	updatedCompany, _ := s.companyRepo.FindByID(id)
	return s.toCompanyOutput(updatedCompany), nil
}

func (s *companyService) DeleteCompany(id uint) error {
	return s.companyRepo.Delete(id)
}

func (s *companyService) CompleteCompanySetup(input *input.CompleteCompanySetupInput, userID uint) (*output.CompleteCompanyProfileOutput, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	company := &models.Company{
		CompanyName:    input.Company.CompanyName,
		BusinessTypeID: input.Company.BusinessTypeID,
		GSTNumber:      input.Company.GSTNumber,
		PANNumber:      input.Company.PANNumber,
		CreatedBy:      &userID,
	}

	if err := tx.Create(company).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create company: %v", err)
	}

	contact := &models.CompanyContact{
		CompanyID:       company.ID,
		Mobile:          input.Contact.Mobile,
		AlternateMobile: input.Contact.AlternateMobile,
		Email:           input.Contact.Email,
	}
	if err := tx.Create(contact).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create contact: %v", err)
	}

	address := &models.CompanyAddress{
		CompanyID:    company.ID,
		AddressLine1: input.Address.AddressLine1,
		AddressLine2: input.Address.AddressLine2,
		City:         input.Address.City,
		StateID:      input.Address.StateID,
		CountryID:    input.Address.CountryID,
		Pincode:      input.Address.Pincode,
	}
	if err := tx.Create(address).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create address: %v", err)
	}

	if input.BankDetails != nil {
		bankDetail := &models.CompanyBankDetail{
			CompanyID:         company.ID,
			BankName:          input.BankDetails.BankName,
			AccountHolderName: input.BankDetails.AccountHolderName,
			AccountNumber:     input.BankDetails.AccountNumber,
			IFSCCode:          input.BankDetails.IFSCCode,
			BranchName:        input.BankDetails.BranchName,
			IsPrimary:         input.BankDetails.IsPrimary,
		}
		if err := tx.Create(bankDetail).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create bank details: %v", err)
		}
	}

	if input.UPIDetails != nil {
		upiDetail := &models.CompanyUPIDetail{
			CompanyID: company.ID,
			UPIID:     input.UPIDetails.UPIID,
			UPIQRURL:  input.UPIDetails.UPIQRURL,
		}
		if err := tx.Create(upiDetail).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create UPI details: %v", err)
		}
	}

	invoiceSettings := &models.CompanyInvoiceSetting{
		CompanyID:            company.ID,
		InvoicePrefix:        "INV",
		InvoiceStartNumber:   1,
		CurrentInvoiceNumber: 1,
		ShowLogo:             true,
		ShowSignature:        false,
		RoundOffTotal:        true,
	}
	if input.InvoiceSettings != nil {
		invoiceSettings.InvoicePrefix = input.InvoiceSettings.InvoicePrefix
		invoiceSettings.InvoiceStartNumber = input.InvoiceSettings.InvoiceStartNumber
		invoiceSettings.CurrentInvoiceNumber = input.InvoiceSettings.InvoiceStartNumber
		invoiceSettings.ShowLogo = input.InvoiceSettings.ShowLogo
		invoiceSettings.ShowSignature = input.InvoiceSettings.ShowSignature
		invoiceSettings.RoundOffTotal = input.InvoiceSettings.RoundOffTotal
	}
	if err := tx.Create(invoiceSettings).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create invoice settings: %v", err)
	}

	taxSettings := &models.CompanyTaxSetting{
		CompanyID:  company.ID,
		GSTEnabled: true,
		TaxTypeID:  1,
	}
	if input.TaxSettings != nil {
		taxSettings.GSTEnabled = input.TaxSettings.GSTEnabled
		taxSettings.TaxTypeID = input.TaxSettings.TaxTypeID
	}
	if err := tx.Create(taxSettings).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create tax settings: %v", err)
	}

	regionalSettings := &models.CompanyRegionalSetting{
		CompanyID:      company.ID,
		Timezone:       "Asia/Kolkata",
		DateFormat:     "DD/MM/YYYY",
		TimeFormat:     "24h",
		CurrencyCode:   "INR",
		CurrencySymbol: "₹",
		LanguageCode:   "en",
	}
	if input.RegionalSettings != nil {
		regionalSettings.Timezone = input.RegionalSettings.Timezone
		regionalSettings.DateFormat = input.RegionalSettings.DateFormat
		regionalSettings.TimeFormat = input.RegionalSettings.TimeFormat
		regionalSettings.CurrencyCode = input.RegionalSettings.CurrencyCode
		regionalSettings.CurrencySymbol = input.RegionalSettings.CurrencySymbol
		regionalSettings.LanguageCode = input.RegionalSettings.LanguageCode
	}
	if err := tx.Create(regionalSettings).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create regional settings: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return s.GetCompany(company.ID)
}

func (s *companyService) UpsertContact(companyID uint, input *input.UpsertCompanyContactInput) (*output.CompanyContactOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	contact := &models.CompanyContact{
		CompanyID:       companyID,
		Mobile:          input.Mobile,
		AlternateMobile: input.AlternateMobile,
		Email:           input.Email,
	}

	if err := s.companyRepo.UpsertContact(contact); err != nil {
		return nil, fmt.Errorf("failed to upsert contact: %v", err)
	}

	return s.GetContact(companyID)
}

func (s *companyService) GetContact(companyID uint) (*output.CompanyContactOutput, error) {
	contact, err := s.companyRepo.GetContact(companyID)
	if err != nil {
		return nil, err
	}
	return s.toContactOutput(contact), nil
}

func (s *companyService) UpsertAddress(companyID uint, input *input.UpsertCompanyAddressInput) (*output.CompanyAddressOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	if _, err := s.locationRepo.GetStateByID(input.StateID); err != nil {
		return nil, fmt.Errorf("state not found: %v", err)
	}
	if _, err := s.locationRepo.GetCountryByID(input.CountryID); err != nil {
		return nil, fmt.Errorf("country not found: %v", err)
	}

	address := &models.CompanyAddress{
		CompanyID:    companyID,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		StateID:      input.StateID,
		CountryID:    input.CountryID,
		Pincode:      input.Pincode,
	}

	if err := s.companyRepo.UpsertAddress(address); err != nil {
		return nil, fmt.Errorf("failed to upsert address: %v", err)
	}

	return s.GetAddress(companyID)
}

func (s *companyService) GetAddress(companyID uint) (*output.CompanyAddressOutput, error) {
	address, err := s.companyRepo.GetAddress(companyID)
	if err != nil {
		return nil, err
	}
	return s.toAddressOutput(address), nil
}

func (s *companyService) CreateBankDetail(companyID uint, input *input.CreateBankDetailInput) (*output.CompanyBankDetailOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	bankDetail := &models.CompanyBankDetail{
		CompanyID:         companyID,
		BankName:          input.BankName,
		AccountHolderName: input.AccountHolderName,
		AccountNumber:     input.AccountNumber,
		IFSCCode:          input.IFSCCode,
		BranchName:        input.BranchName,
		IsPrimary:         input.IsPrimary,
	}

	if err := s.companyRepo.CreateBankDetail(bankDetail); err != nil {
		return nil, fmt.Errorf("failed to create bank detail: %v", err)
	}

	created, _ := s.companyRepo.GetBankDetailByID(bankDetail.ID)
	return s.toBankDetailOutput(created), nil
}

func (s *companyService) GetBankDetails(companyID uint) ([]output.CompanyBankDetailOutput, error) {
	bankDetails, err := s.companyRepo.GetBankDetails(companyID)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.CompanyBankDetailOutput, len(bankDetails))
	for i, bd := range bankDetails {
		outputs[i] = *s.toBankDetailOutput(&bd)
	}
	return outputs, nil
}

func (s *companyService) UpdateBankDetail(id uint, input *input.UpdateBankDetailInput) (*output.CompanyBankDetailOutput, error) {
	bankDetail, err := s.companyRepo.GetBankDetailByID(id)
	if err != nil {
		return nil, err
	}

	if input.BankName != nil {
		bankDetail.BankName = *input.BankName
	}
	if input.AccountHolderName != nil {
		bankDetail.AccountHolderName = *input.AccountHolderName
	}
	if input.AccountNumber != nil {
		bankDetail.AccountNumber = *input.AccountNumber
	}
	if input.IFSCCode != nil {
		bankDetail.IFSCCode = *input.IFSCCode
	}
	if input.BranchName != nil {
		bankDetail.BranchName = *input.BranchName
	}
	if input.IsPrimary != nil {
		bankDetail.IsPrimary = *input.IsPrimary
	}
	if input.IsActive != nil {
		bankDetail.IsActive = *input.IsActive
	}

	if err := s.companyRepo.UpdateBankDetail(bankDetail); err != nil {
		return nil, fmt.Errorf("failed to update bank detail: %v", err)
	}

	updated, _ := s.companyRepo.GetBankDetailByID(id)
	return s.toBankDetailOutput(updated), nil
}

func (s *companyService) DeleteBankDetail(id uint) error {
	return s.companyRepo.DeleteBankDetail(id)
}

func (s *companyService) UpsertUPIDetail(companyID uint, input *input.UpsertUPIDetailInput) (*output.CompanyUPIDetailOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	upiDetail := &models.CompanyUPIDetail{
		CompanyID: companyID,
		UPIID:     input.UPIID,
		UPIQRURL:  input.UPIQRURL,
	}

	if err := s.companyRepo.UpsertUPIDetail(upiDetail); err != nil {
		return nil, fmt.Errorf("failed to upsert UPI detail: %v", err)
	}

	return s.GetUPIDetail(companyID)
}

func (s *companyService) GetUPIDetail(companyID uint) (*output.CompanyUPIDetailOutput, error) {
	upiDetail, err := s.companyRepo.GetUPIDetail(companyID)
	if err != nil {
		return nil, err
	}
	return s.toUPIDetailOutput(upiDetail), nil
}

func (s *companyService) UpsertInvoiceSettings(companyID uint, input *input.UpsertInvoiceSettingsInput) (*output.CompanyInvoiceSettingsOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	settings := &models.CompanyInvoiceSetting{
		CompanyID:            companyID,
		InvoicePrefix:        input.InvoicePrefix,
		InvoiceStartNumber:   input.InvoiceStartNumber,
		CurrentInvoiceNumber: input.InvoiceStartNumber,
		ShowLogo:             input.ShowLogo,
		ShowSignature:        input.ShowSignature,
		RoundOffTotal:        input.RoundOffTotal,
	}

	if err := s.companyRepo.UpsertInvoiceSettings(settings); err != nil {
		return nil, fmt.Errorf("failed to upsert invoice settings: %v", err)
	}

	return s.GetInvoiceSettings(companyID)
}

func (s *companyService) GetInvoiceSettings(companyID uint) (*output.CompanyInvoiceSettingsOutput, error) {
	settings, err := s.companyRepo.GetInvoiceSettings(companyID)
	if err != nil {
		return nil, err
	}
	return s.toInvoiceSettingsOutput(settings), nil
}

func (s *companyService) UpsertTaxSettings(companyID uint, input *input.UpsertTaxSettingsInput) (*output.CompanyTaxSettingsOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	if _, err := s.taxTypeRepo.FindByID(input.TaxTypeID); err != nil {
		return nil, fmt.Errorf("tax type not found: %v", err)
	}

	settings := &models.CompanyTaxSetting{
		CompanyID:  companyID,
		GSTEnabled: input.GSTEnabled,
		TaxTypeID:  input.TaxTypeID,
	}

	if err := s.companyRepo.UpsertTaxSettings(settings); err != nil {
		return nil, fmt.Errorf("failed to upsert tax settings: %v", err)
	}

	return s.GetTaxSettings(companyID)
}

func (s *companyService) GetTaxSettings(companyID uint) (*output.CompanyTaxSettingsOutput, error) {
	settings, err := s.companyRepo.GetTaxSettings(companyID)
	if err != nil {
		return nil, err
	}
	return s.toTaxSettingsOutput(settings), nil
}

func (s *companyService) UpsertRegionalSettings(companyID uint, input *input.UpsertRegionalSettingsInput) (*output.CompanyRegionalSettingsOutput, error) {

	if _, err := s.companyRepo.FindByID(companyID); err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}

	settings := &models.CompanyRegionalSetting{
		CompanyID:      companyID,
		Timezone:       input.Timezone,
		DateFormat:     input.DateFormat,
		TimeFormat:     input.TimeFormat,
		CurrencyCode:   input.CurrencyCode,
		CurrencySymbol: input.CurrencySymbol,
		LanguageCode:   input.LanguageCode,
	}

	if err := s.companyRepo.UpsertRegionalSettings(settings); err != nil {
		return nil, fmt.Errorf("failed to upsert regional settings: %v", err)
	}

	return s.GetRegionalSettings(companyID)
}

func (s *companyService) GetRegionalSettings(companyID uint) (*output.CompanyRegionalSettingsOutput, error) {
	settings, err := s.companyRepo.GetRegionalSettings(companyID)
	if err != nil {
		return nil, err
	}
	return s.toRegionalSettingsOutput(settings), nil
}

func (s *companyService) toCompanyOutput(c *models.Company) *output.CompanyOutput {
	return &output.CompanyOutput{
		ID:             c.ID,
		CompanyName:    c.CompanyName,
		BusinessTypeID: c.BusinessTypeID,
		BusinessType: output.BusinessTypeOutput{
			ID:          c.BusinessType.ID,
			TypeName:    c.BusinessType.TypeName,
			Description: c.BusinessType.Description,
			IsActive:    c.BusinessType.IsActive,
			CreatedAt:   c.BusinessType.CreatedAt,
		},
		GSTNumber: c.GSTNumber,
		PANNumber: c.PANNumber,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (s *companyService) toContactOutput(c *models.CompanyContact) *output.CompanyContactOutput {
	return &output.CompanyContactOutput{
		ID:              c.ID,
		CompanyID:       c.CompanyID,
		Mobile:          c.Mobile,
		AlternateMobile: c.AlternateMobile,
		Email:           c.Email,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

func (s *companyService) toAddressOutput(a *models.CompanyAddress) *output.CompanyAddressOutput {
	return &output.CompanyAddressOutput{
		ID:           a.ID,
		CompanyID:    a.CompanyID,
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		City:         a.City,
		StateID:      a.StateID,
		State: output.StateOutput{
			ID:        a.State.ID,
			CountryID: a.State.CountryID,
			StateName: a.State.StateName,
			StateCode: a.State.StateCode,
			Country: output.CountryOutput{
				ID:          a.State.Country.ID,
				CountryName: a.State.Country.CountryName,
				CountryCode: a.State.Country.CountryCode,
				PhoneCode:   a.State.Country.PhoneCode,
			},
		},
		CountryID: a.CountryID,
		Country: output.CountryOutput{
			ID:          a.Country.ID,
			CountryName: a.Country.CountryName,
			CountryCode: a.Country.CountryCode,
			PhoneCode:   a.Country.PhoneCode,
		},
		Pincode:   a.Pincode,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (s *companyService) toBankDetailOutput(b *models.CompanyBankDetail) *output.CompanyBankDetailOutput {
	return &output.CompanyBankDetailOutput{
		ID:                b.ID,
		CompanyID:         b.CompanyID,
		BankName:          b.BankName,
		AccountHolderName: b.AccountHolderName,
		AccountNumber:     b.AccountNumber,
		IFSCCode:          b.IFSCCode,
		BranchName:        b.BranchName,
		IsPrimary:         b.IsPrimary,
		IsActive:          b.IsActive,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
}

func (s *companyService) toUPIDetailOutput(u *models.CompanyUPIDetail) *output.CompanyUPIDetailOutput {
	return &output.CompanyUPIDetailOutput{
		ID:        u.ID,
		CompanyID: u.CompanyID,
		UPIID:     u.UPIID,
		UPIQRURL:  u.UPIQRURL,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (s *companyService) toInvoiceSettingsOutput(i *models.CompanyInvoiceSetting) *output.CompanyInvoiceSettingsOutput {
	return &output.CompanyInvoiceSettingsOutput{
		ID:                   i.ID,
		CompanyID:            i.CompanyID,
		InvoicePrefix:        i.InvoicePrefix,
		InvoiceStartNumber:   i.InvoiceStartNumber,
		CurrentInvoiceNumber: i.CurrentInvoiceNumber,
		ShowLogo:             i.ShowLogo,
		ShowSignature:        i.ShowSignature,
		RoundOffTotal:        i.RoundOffTotal,
		CreatedAt:            i.CreatedAt,
		UpdatedAt:            i.UpdatedAt,
	}
}

func (s *companyService) toTaxSettingsOutput(t *models.CompanyTaxSetting) *output.CompanyTaxSettingsOutput {
	return &output.CompanyTaxSettingsOutput{
		ID:         t.ID,
		CompanyID:  t.CompanyID,
		GSTEnabled: t.GSTEnabled,
		TaxTypeID:  t.TaxTypeID,
		TaxType: output.TaxTypeOutput{
			ID:          t.TaxType.ID,
			TaxName:     t.TaxType.TaxName,
			TaxCode:     t.TaxType.TaxCode,
			Description: t.TaxType.Description,
		},
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (s *companyService) toRegionalSettingsOutput(r *models.CompanyRegionalSetting) *output.CompanyRegionalSettingsOutput {
	return &output.CompanyRegionalSettingsOutput{
		ID:             r.ID,
		CompanyID:      r.CompanyID,
		Timezone:       r.Timezone,
		DateFormat:     r.DateFormat,
		TimeFormat:     r.TimeFormat,
		CurrencyCode:   r.CurrencyCode,
		CurrencySymbol: r.CurrencySymbol,
		LanguageCode:   r.LanguageCode,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func (s *companyService) toCompleteProfileOutput(c *models.Company) *output.CompleteCompanyProfileOutput {
	profile := &output.CompleteCompanyProfileOutput{
		Company: *s.toCompanyOutput(c),
	}

	if c.Contact != nil {
		profile.Contact = s.toContactOutput(c.Contact)
	}

	if c.Address != nil {
		profile.Address = s.toAddressOutput(c.Address)
	}

	if len(c.BankDetails) > 0 {
		bankDetails := make([]output.CompanyBankDetailOutput, len(c.BankDetails))
		for i, bd := range c.BankDetails {
			bankDetails[i] = *s.toBankDetailOutput(&bd)
		}
		profile.BankDetails = bankDetails
	}

	if c.UPIDetails != nil {
		profile.UPIDetails = s.toUPIDetailOutput(c.UPIDetails)
	}

	if c.InvoiceSettings != nil {
		profile.InvoiceSettings = s.toInvoiceSettingsOutput(c.InvoiceSettings)
	}

	if c.TaxSettings != nil {
		profile.TaxSettings = s.toTaxSettingsOutput(c.TaxSettings)
	}

	if c.RegionalSettings != nil {
		profile.RegionalSettings = s.toRegionalSettingsOutput(c.RegionalSettings)
	}

	return profile
}
