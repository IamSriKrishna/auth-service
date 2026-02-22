package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetByAppleID(appleID string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int, search string) ([]models.User, int64, error)
	UpdateLastLogin(id uint) error
	UpdatePasswordChangedAt(id uint) error
	GetDashboardStats(customerType *string, fromDate, toDate *time.Time) (map[string]interface{}, error)
}

type RoleRepository interface {
	GetByID(id uint) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	GetAll() ([]models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
}

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	GetByTokenID(tokenID string) (*models.RefreshToken, error)
	GetByUserID(userID uint) ([]models.RefreshToken, error)
	Delete(tokenID string) error
	DeleteByUserID(userID uint) error
	DeleteExpired() error
}

type UserSessionRepository interface {
	Create(session *models.UserSession) error
	GetBySessionID(sessionID string) (*models.UserSession, error)
	GetByUserID(userID uint) ([]models.UserSession, error)
	Delete(sessionID string) error
	DeleteByUserID(userID uint) error
	DeleteExpired() error
}

type SupportRepository interface {
	Create(support *models.Support) error
	GetByID(id uint) (*models.Support, error)
	List(offset, limit int) ([]models.Support, int64, error)
	Update(support *models.Support) error
	Delete(id uint) error
}

type CompanyRepository interface {
	Create(company *models.Company) error
	FindByID(id uint) (*models.Company, error)
	FindAll(businessTypeID *uint, search *string, isActive *bool, page, pageSize int) ([]models.Company, int64, error)
	Update(company *models.Company) error
	Delete(id uint) error

	UpsertContact(contact *models.CompanyContact) error
	GetContact(companyID uint) (*models.CompanyContact, error)

	UpsertAddress(address *models.CompanyAddress) error
	GetAddress(companyID uint) (*models.CompanyAddress, error)

	CreateBankDetail(bankDetail *models.CompanyBankDetail) error
	GetBankDetails(companyID uint) ([]models.CompanyBankDetail, error)
	GetBankDetailByID(id uint) (*models.CompanyBankDetail, error)
	UpdateBankDetail(bankDetail *models.CompanyBankDetail) error
	DeleteBankDetail(id uint) error

	UpsertUPIDetail(upiDetail *models.CompanyUPIDetail) error
	GetUPIDetail(companyID uint) (*models.CompanyUPIDetail, error)

	UpsertInvoiceSettings(settings *models.CompanyInvoiceSetting) error
	GetInvoiceSettings(companyID uint) (*models.CompanyInvoiceSetting, error)

	UpsertTaxSettings(settings *models.CompanyTaxSetting) error
	GetTaxSettings(companyID uint) (*models.CompanyTaxSetting, error)

	UpsertRegionalSettings(settings *models.CompanyRegionalSetting) error
	GetRegionalSettings(companyID uint) (*models.CompanyRegionalSetting, error)

	GetCompleteProfile(companyID uint) (*models.Company, error)
}

type BusinessTypeRepository interface {
	FindAll() ([]models.BusinessType, error)
	FindByID(id uint) (*models.BusinessType, error)
	Create(businessType *models.BusinessType) error
	Update(businessType *models.BusinessType) error
	Delete(id uint) error
}

type TaxTypeRepository interface {
	FindAll() ([]models.TaxType, error)
	FindByID(id uint) (*models.TaxType, error)
	Create(taxType *models.TaxType) error
	Update(taxType *models.TaxType) error
	Delete(id uint) error
}
type LocationRepository interface {
	GetAllCountries() ([]models.Country, error)
	GetCountryByID(id uint) (*models.Country, error)
	GetStatesByCountry(countryID uint) ([]models.State, error)
	GetStateByID(id uint) (*models.State, error)
}

type VendorRepository interface {
	Create(vendor *models.Vendor) error
	Update(vendor *models.Vendor) error
	FindByID(id uint) (*models.Vendor, error)
	FindAll(page, limit int) ([]models.Vendor, int64, error)
	Delete(id uint) error
	FindByMobile(mobile string) (*models.Vendor, error)
}

type CustomerRepository interface {
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	FindByID(id uint) (*models.Customer, error)
	FindAll(page, limit int) ([]models.Customer, int64, error)
	Delete(customer *models.Customer) error
	FindByMobile(mobile string) (*models.Customer, error)
}

type ItemRepository interface {
	Create(item *models.Item) error
	FindByID(id string) (*models.Item, error)
	FindAll(limit, offset int) ([]models.Item, int64, error)
	Update(item *models.Item) error
	Delete(id string) error
	FindByType(itemType string, limit, offset int) ([]models.Item, int64, error)
	DeductStockQuantity(itemID string, variantSKU *string, quantity float64) error
	CheckReorderPoint(itemID string, variantSKU *string) (*models.Variant, error)
	GetVariantBySKU(sku string) (*models.Variant, error)
	UpdateVariantStock(variantID uint, newQuantity float64) error
}

type OpeningStockRepository interface {
	CreateOrUpdateOpeningStock(itemID string, openingStock, ratePerUnit float64) error
	GetOpeningStock(itemID string) (*models.OpeningStock, error)
	CreateOrUpdateVariantOpeningStock(variantSKU string, openingStock, ratePerUnit float64) error
	GetVariantOpeningStock(variantSKU string) (*models.VariantOpeningStock, error)
	GetAllVariantOpeningStocks(itemID string) ([]models.VariantOpeningStock, error)
	RecordStockMovement(movement *models.StockMovement) error
	GetStockMovements(itemID string) ([]models.StockMovement, error)
}

type InvoiceRepository interface {
	Create(invoice *models.Invoice) error
	FindByID(id string) (*models.Invoice, error)
	FindAll(limit, offset int) ([]models.Invoice, int64, error)
	Update(invoice *models.Invoice) error
	Delete(id string) error
	FindByCustomerID(customerID string, limit, offset int) ([]models.Invoice, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.Invoice, int64, error)
	GetNextInvoiceNumber() (string, error)
}

type SalespersonRepository interface {
	Create(salesperson *models.Salesperson) error
	FindByID(id uint) (*models.Salesperson, error)
	FindAll(limit, offset int) ([]models.Salesperson, int64, error)
	Update(salesperson *models.Salesperson) error
	Delete(id uint) error
	FindByEmail(email string) (*models.Salesperson, error)
}

type TaxRepository interface {
	Create(tax *models.Tax) error
	FindByID(id uint) (*models.Tax, error)
	FindAll(limit, offset int) ([]models.Tax, int64, error)
	Update(tax *models.Tax) error
	Delete(id uint) error
}

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uint) (*models.Payment, error)
	FindByInvoiceID(invoiceID string) ([]models.Payment, error)
	Delete(id uint) error
}

type ManufacturerRepository interface {
	Create(manufacturer *models.Manufacturer) error
	FindByID(id uint) (*models.Manufacturer, error)
	FindAll(limit, offset int) ([]models.Manufacturer, int64, error)
	Update(manufacturer *models.Manufacturer) error
	Delete(id uint) error
}

type BrandRepository interface {
	Create(brand *models.Brand) error
	FindByID(id uint) (*models.Brand, error)
	FindAll(limit, offset int) ([]models.Brand, int64, error)
	Update(brand *models.Brand) error
	Delete(id uint) error
}

type BankRepository interface {
	Create(bank *models.Bank) error
	FindByID(id uint) (*models.Bank, error)
	FindByIFSCCode(ifscCode string) (*models.Bank, error)
	FindAll(limit, offset int) ([]models.Bank, int64, error)
	Update(bank *models.Bank) error
	Delete(id uint) error
}

type InventoryBalanceRepository interface {
	GetBalance(itemID string, variantSKU *string) (*models.InventoryBalance, error)
	GetBalances(itemID string) ([]models.InventoryBalance, error)
	UpdateBalance(balance *models.InventoryBalance) error
	CreateJournalEntry(entry *models.InventoryJournal) error
	GetJournalEntries(itemID string, limit, offset int) ([]models.InventoryJournal, int64, error)
	ReserveInventory(itemID string, variantSKU *string, quantity float64, referenceID, referenceNo string) error
	ReleaseReservation(itemID string, variantSKU *string, quantity float64, referenceID string) error
}
type ProductionOrderRepository interface {
	Create(order *models.ProductionOrder) error
	FindByID(id string) (*models.ProductionOrder, error)
	FindAll(limit, offset int) ([]models.ProductionOrder, int64, error)
	Update(order *models.ProductionOrder) error
	Delete(id string) error
	FindByProductionOrderNumber(orderNo string) (*models.ProductionOrder, error)
}
