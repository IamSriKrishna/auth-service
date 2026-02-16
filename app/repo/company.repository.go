package repo

import (
	"errors"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) FindByID(id uint) (*models.Company, error) {
	var company models.Company
	err := r.db.Preload("BusinessType").First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) FindAll(businessTypeID *uint, search *string, isActive *bool, page, pageSize int) ([]models.Company, int64, error) {
	var companies []models.Company
	var total int64

	query := r.db.Model(&models.Company{})

	if businessTypeID != nil {
		query = query.Where("business_type_id = ?", *businessTypeID)
	}

	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		query = query.Where("company_name LIKE ? OR gst_number LIKE ? OR pan_number LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("BusinessType").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&companies).Error

	return companies, total, err
}

func (r *companyRepository) Update(company *models.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Company{}, id).Error
}

func (r *companyRepository) UpsertContact(contact *models.CompanyContact) error {
	var existing models.CompanyContact
	err := r.db.Where("company_id = ?", contact.CompanyID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(contact).Error
		}
		return err
	}

	contact.ID = existing.ID
	return r.db.Save(contact).Error
}

func (r *companyRepository) GetContact(companyID uint) (*models.CompanyContact, error) {
	var contact models.CompanyContact
	err := r.db.Where("company_id = ?", companyID).First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *companyRepository) UpsertAddress(address *models.CompanyAddress) error {
	var existing models.CompanyAddress
	err := r.db.Where("company_id = ?", address.CompanyID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(address).Error
		}
		return err
	}

	address.ID = existing.ID
	return r.db.Save(address).Error
}

func (r *companyRepository) GetAddress(companyID uint) (*models.CompanyAddress, error) {
	var address models.CompanyAddress
	err := r.db.Preload("State").Preload("Country").
		Where("company_id = ?", companyID).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *companyRepository) CreateBankDetail(bankDetail *models.CompanyBankDetail) error {
	return r.db.Create(bankDetail).Error
}

func (r *companyRepository) GetBankDetails(companyID uint) ([]models.CompanyBankDetail, error) {
	var bankDetails []models.CompanyBankDetail
	err := r.db.Where("company_id = ?", companyID).
		Order("is_primary DESC, created_at DESC").
		Find(&bankDetails).Error
	return bankDetails, err
}

func (r *companyRepository) GetBankDetailByID(id uint) (*models.CompanyBankDetail, error) {
	var bankDetail models.CompanyBankDetail
	err := r.db.First(&bankDetail, id).Error
	if err != nil {
		return nil, err
	}
	return &bankDetail, nil
}

func (r *companyRepository) UpdateBankDetail(bankDetail *models.CompanyBankDetail) error {
	return r.db.Save(bankDetail).Error
}

func (r *companyRepository) DeleteBankDetail(id uint) error {
	return r.db.Delete(&models.CompanyBankDetail{}, id).Error
}

func (r *companyRepository) UpsertUPIDetail(upiDetail *models.CompanyUPIDetail) error {
	var existing models.CompanyUPIDetail
	err := r.db.Where("company_id = ?", upiDetail.CompanyID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(upiDetail).Error
		}
		return err
	}

	upiDetail.ID = existing.ID
	return r.db.Save(upiDetail).Error
}

func (r *companyRepository) GetUPIDetail(companyID uint) (*models.CompanyUPIDetail, error) {
	var upiDetail models.CompanyUPIDetail
	err := r.db.Where("company_id = ?", companyID).First(&upiDetail).Error
	if err != nil {
		return nil, err
	}
	return &upiDetail, nil
}

func (r *companyRepository) UpsertInvoiceSettings(settings *models.CompanyInvoiceSetting) error {
	var existing models.CompanyInvoiceSetting
	err := r.db.Where("company_id = ?", settings.CompanyID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(settings).Error
		}
		return err
	}

	settings.ID = existing.ID
	return r.db.Save(settings).Error
}

func (r *companyRepository) GetInvoiceSettings(companyID uint) (*models.CompanyInvoiceSetting, error) {
	var settings models.CompanyInvoiceSetting
	err := r.db.Where("company_id = ?", companyID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *companyRepository) UpsertTaxSettings(settings *models.CompanyTaxSetting) error {
	var existing models.CompanyTaxSetting
	err := r.db.Where("company_id = ?", settings.CompanyID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(settings).Error
		}
		return err
	}

	settings.ID = existing.ID
	return r.db.Save(settings).Error
}

func (r *companyRepository) GetTaxSettings(companyID uint) (*models.CompanyTaxSetting, error) {
	var settings models.CompanyTaxSetting
	err := r.db.Preload("TaxType").Where("company_id = ?", companyID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *companyRepository) UpsertRegionalSettings(settings *models.CompanyRegionalSetting) error {
	var existing models.CompanyRegionalSetting
	err := r.db.Where("company_id = ?", settings.CompanyID).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(settings).Error
		}
		return err
	}

	settings.ID = existing.ID
	return r.db.Save(settings).Error
}

func (r *companyRepository) GetRegionalSettings(companyID uint) (*models.CompanyRegionalSetting, error) {
	var settings models.CompanyRegionalSetting
	err := r.db.Where("company_id = ?", companyID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *companyRepository) GetCompleteProfile(companyID uint) (*models.Company, error) {
	var company models.Company
	err := r.db.Preload("BusinessType").
		Preload("Contact").
		Preload("Address.State").
		Preload("Address.Country").
		Preload("BankDetails", "deleted_at IS NULL").
		Preload("UPIDetails").
		Preload("InvoiceSettings").
		Preload("TaxSettings.TaxType").
		Preload("RegionalSettings").
		First(&company, companyID).Error

	if err != nil {
		return nil, err
	}
	return &company, nil
}
