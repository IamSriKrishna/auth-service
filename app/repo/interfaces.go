package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
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
	ListWithFilters(offset, limit int, search, role string) ([]models.User, int64, error)
	UpdateLastLogin(id uint) error
	UpdatePasswordChangedAt(id uint) error
	GetDashboardStats(customerType *string, fromDate, toDate *time.Time) (map[string]interface{}, error)
	GetByIDAndCompanyID(id uint, companyID uint) (*models.User, error)
	ListByCompanyID(
		companyID uint,
	) ([]models.User, error)
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
	// Existing methods retained.
	Create(vendor *models.Vendor) error
	Update(vendor *models.Vendor) error
	FindByID(id uint) (*models.Vendor, error)
	FindAll(page, limit int) ([]models.Vendor, int64, error)
	FindByUserID(
		userID uint,
		companyID uint,
		page int,
		limit int,
	) ([]models.Vendor, int64, error)
	FindByIDAndUser(
		id uint,
		userID uint,
	) (*models.Vendor, error)
	Delete(id uint) error
	FindByMobile(mobile string) (*models.Vendor, error)

	// Company-scoped methods added.
	FindByIDAndCompany(
		id uint,
		companyID uint,
	) (*models.Vendor, error)

	FindByCompanyID(
		companyID uint,
		page int,
		limit int,
	) ([]models.Vendor, int64, error)

	UpdateByCompanyID(
		vendor *models.Vendor,
		companyID uint,
	) error

	DeleteByIDAndCompany(
		id uint,
		companyID uint,
	) error

	FindByMobileAndCompany(
		mobile string,
		companyID uint,
	) (*models.Vendor, error)
}
type CustomerRepository interface {
	// Existing methods retained.
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	FindByID(id uint) (*models.Customer, error)
	FindAll(page, limit int) ([]models.Customer, int64, error)
	FindByUserID(
		userID uint,
		companyID uint,
		page int,
		limit int,
	) ([]models.Customer, int64, error)
	FindByIDAndUser(
		id uint,
		userID uint,
	) (*models.Customer, error)
	Delete(customer *models.Customer) error
	FindByMobile(mobile string) (*models.Customer, error)

	// Company-scoped methods added.
	FindByIDAndCompany(
		id uint,
		companyID uint,
	) (*models.Customer, error)

	FindByCompanyID(
		companyID uint,
		page int,
		limit int,
	) ([]models.Customer, int64, error)

	UpdateByCompanyID(
		customer *models.Customer,
		companyID uint,
	) error

	DeleteByIDAndCompany(
		id uint,
		companyID uint,
	) error

	FindByMobileAndCompany(
		mobile string,
		companyID uint,
	) (*models.Customer, error)
}
type ItemRepository interface {
	Create(item *models.Item) error
	FindByID(id string) (*models.Item, error)
	FindAll(limit, offset int) ([]models.Item, int64, error)
	FindByCreatedBy(createdBy string, limit, offset int) ([]models.Item, int64, error)
	Update(item *models.Item) error
	Delete(id string) error
	DeleteByCreatedBy(id string, createdBy string) error
	FindByType(itemType string, limit, offset int) ([]models.Item, int64, error)
	FindByTypeAndCreatedBy(itemType string, createdBy string, limit, offset int) ([]models.Item, int64, error)
	DeductStockQuantity(itemID string, variantSKU *string, quantity float64) error
	CheckReorderPoint(itemID string, variantSKU *string) (*models.Variant, error)
	GetVariantBySKU(sku string) (*models.Variant, error)
	UpdateVariantStock(variantID uint, newQuantity float64) error
}

type ProductRepository interface {
	// Existing methods retained.
	Create(product *models.Product) error
	FindByID(id string) (*models.Product, error)
	FindAll(limit, offset int) ([]models.Product, int64, error)
	FindByCreatedBy(
		createdBy string,
		limit int,
		offset int,
	) ([]models.Product, int64, error)
	FindByCreatedByAndCompany(
		createdBy string,
		companyID uint,
		limit int,
		offset int,
	) ([]models.Product, int64, error)
	Update(product *models.Product) error
	Delete(id string) error
	DeleteByCreatedBy(id string, createdBy string) error
	DeductProductVariantStock(
		productID string,
		variantSKU string,
		quantity float64,
	) error
	CheckProductVariantReorderPoint(
		productID string,
		variantSKU string,
	) (*models.ProductVariant, error)
	GetProductVariantBySKU(
		sku string,
	) (*models.ProductVariant, error)
	UpdateProductVariantStock(
		variantID uint,
		newQuantity float64,
	) error
	GetProductVariantsByProductID(
		productID string,
	) ([]models.ProductVariant, error)

	// Company-scoped methods added.
	FindByIDAndCompany(
		id string,
		companyID uint,
	) (*models.Product, error)

	FindByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]models.Product, int64, error)

	UpdateByCompany(
		product *models.Product,
		companyID uint,
	) error

	DeleteByCompany(
		id string,
		companyID uint,
	) error

	GetProductVariantsByProductIDAndCompany(
		productID string,
		companyID uint,
	) ([]models.ProductVariant, error)
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
	// Existing methods retained.
	Create(invoice *models.Invoice) error
	FindByID(id string) (*models.Invoice, error)
	FindAll(limit, offset int) ([]models.Invoice, int64, error)
	Update(invoice *models.Invoice) error
	Delete(id string) error
	FindByCustomerID(
		customerID string,
		limit int,
		offset int,
	) ([]models.Invoice, int64, error)
	FindByStatus(
		status string,
		limit int,
		offset int,
	) ([]models.Invoice, int64, error)
	GetNextInvoiceNumber() (string, error)

	// Company-scoped methods.
	FindByIDAndCompany(
		id string,
		companyID uint,
	) (*models.Invoice, error)

	FindAllByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]models.Invoice, int64, error)

	FindByCustomerIDAndCompany(
		customerID uint,
		companyID uint,
		limit int,
		offset int,
	) ([]models.Invoice, int64, error)

	FindByStatusAndCompany(
		status string,
		companyID uint,
		limit int,
		offset int,
	) ([]models.Invoice, int64, error)

	UpdateByCompany(
		invoice *models.Invoice,
		companyID uint,
	) error

	DeleteByCompany(
		id string,
		companyID uint,
	) error
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
	FindByStringID(id string) (*models.Manufacturer, error)
	FindAll(limit, offset int) ([]models.Manufacturer, int64, error)
	FindAllWithFilter(limit, offset int, companyID *uint, productGroupID *string) ([]models.Manufacturer, int64, error)
	FindByProductGroupID(productGroupID string) ([]models.Manufacturer, error)
	Update(manufacturer *models.Manufacturer) error
	Delete(id uint) error
	DeleteByStringID(id string) error
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

type ItemGroupRepository interface {
	Create(itemGroup *models.ItemGroup) error
	FindByID(id string) (*models.ItemGroup, error)
	FindAll(limit, offset int, search string) ([]models.ItemGroup, int64, error)
	Update(itemGroup *models.ItemGroup) error
	Delete(id string) error
	FindByName(name string) (*models.ItemGroup, error)
	FindActiveGroups(limit, offset int) ([]models.ItemGroup, int64, error)
}

type ProductGroupRepository interface {
	Create(productGroup *models.ProductGroup) error
	FindByID(id string) (*models.ProductGroup, error)
	FindAll(limit, offset int, search string) ([]models.ProductGroup, int64, error)
	Update(productGroup *models.ProductGroup) error
	Delete(id string) error
	FindByName(name string) (*models.ProductGroup, error)
	FindActiveGroups(limit, offset int) ([]models.ProductGroup, int64, error)
}

type PurchaseOrderRepository interface {
	Create(po *models.PurchaseOrder) (*models.PurchaseOrder, error)
	FindByID(id string) (*models.PurchaseOrder, error)
	FindAll(limit, offset int) ([]models.PurchaseOrder, int64, error)
	FindByVendor(vendorID uint, limit, offset int) ([]models.PurchaseOrder, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.PurchaseOrder, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.PurchaseOrder, int64, error)
	Update(id string, po *models.PurchaseOrder) (*models.PurchaseOrder, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetDB() *gorm.DB
	purchaseOrderPreloads(db *gorm.DB) *gorm.DB
	FindByIDAndCompany(
		id string,
		companyID uint,
	) (*models.PurchaseOrder, error)
	FindAllByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]models.PurchaseOrder, int64, error)
	FindByVendorAndCompany(
		vendorID uint,
		companyID uint,
		limit int,
		offset int,
	) ([]models.PurchaseOrder, int64, error)
	FindByCustomerAndCompany(
		customerID uint,
		companyID uint,
		limit int,
		offset int,
	) ([]models.PurchaseOrder, int64, error)
	FindByStatusAndCompany(
		status string,
		companyID uint,
		limit int,
		offset int,
	) ([]models.PurchaseOrder, int64, error)
	UpdateByCompany(
		id string,
		companyID uint,
		purchaseOrder *models.PurchaseOrder,
	) (*models.PurchaseOrder, error)
	DeleteByCompany(
		id string,
		companyID uint,
	) error
	UpdateStatusByCompany(
		id string,
		companyID uint,
		status string,
	) error
}

type SalesOrderRepository interface {
	Create(so *models.SalesOrder) (*models.SalesOrder, error)
	FindByID(id string) (*models.SalesOrder, error)
	FindAll(limit, offset int) ([]models.SalesOrder, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.SalesOrder, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.SalesOrder, int64, error)
	Update(id string, so *models.SalesOrder) (*models.SalesOrder, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetDB() *gorm.DB
}

type BillRepository interface {
	// Existing methods retained.
	Create(bill *models.Bill) (*models.Bill, error)
	FindByID(id string) (*models.Bill, error)
	FindAll(limit, offset int) ([]models.Bill, int64, error)
	FindByCreatedBy(
		createdBy string,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	FindByVendor(
		vendorID uint,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	FindByVendorAndCreatedBy(
		vendorID uint,
		createdBy string,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	FindByStatus(
		status string,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	FindByStatusAndCreatedBy(
		status string,
		createdBy string,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	Update(
		id string,
		bill *models.Bill,
	) (*models.Bill, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetDB() *gorm.DB

	// Company-scoped methods.
	FindByIDAndCompany(
		id string,
		companyID uint,
	) (*models.Bill, error)
	FindAllByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	FindByVendorAndCompany(
		vendorID uint,
		companyID uint,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	FindByStatusAndCompany(
		status string,
		companyID uint,
		limit int,
		offset int,
	) ([]models.Bill, int64, error)
	UpdateByCompany(
		id string,
		companyID uint,
		bill *models.Bill,
	) (*models.Bill, error)
	DeleteByCompany(
		id string,
		companyID uint,
	) error
	UpdateStatusByCompany(
		id string,
		companyID uint,
		status string,
	) error
}

type PackageRepository interface {
	Create(pkg *models.Package) (*models.Package, error)
	FindByID(id string) (*models.Package, error)
	FindAll(limit, offset int) ([]models.Package, int64, error)
	FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Package, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.Package, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.Package, int64, error)
	Update(id string, pkg *models.Package) (*models.Package, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetNextPackageSlipNo() (string, error)
}

type ShipmentRepository interface {
	Create(shipment *models.Shipment) (*models.Shipment, error)
	FindByID(id string) (*models.Shipment, error)
	FindAll(limit, offset int) ([]models.Shipment, int64, error)
	FindByPackage(packageID string, limit, offset int) ([]models.Shipment, int64, error)
	FindBySalesOrder(salesOrderID string, limit, offset int) ([]models.Shipment, int64, error)
	FindByCustomer(customerID uint, limit, offset int) ([]models.Shipment, int64, error)
	FindByStatus(status string, limit, offset int) ([]models.Shipment, int64, error)
	Update(id string, shipment *models.Shipment) (*models.Shipment, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	GetNextShipmentNo() (string, error)
}

type EmployeeRepository interface {
	Create(employee *models.Employee) error

	// Existing methods retained.
	GetByID(id uint) (*models.Employee, error)
	GetByUserID(
		userID uint,
		offset int,
		limit int,
	) ([]models.Employee, int64, error)
	GetByCompany(
		companyID uint,
		offset int,
		limit int,
	) ([]models.Employee, int64, error)
	GetByCompanyAndUser(
		companyID uint,
		userID uint,
		offset int,
		limit int,
	) ([]models.Employee, int64, error)
	Update(employee *models.Employee) error
	Delete(id uint) error

	// Company-safe methods added.
	GetByIDAndCompanyID(
		id uint,
		companyID uint,
	) (*models.Employee, error)

	UpdateByCompanyID(
		employee *models.Employee,
		companyID uint,
	) error

	DeleteByIDAndCompanyID(
		id uint,
		companyID uint,
	) error
}
type EmployeeAttendanceRepository interface {
	Create(attendance *models.EmployeeAttendance) error
	GetByID(id uint) (*models.EmployeeAttendance, error)
	GetByEmployeeAndDate(employeeID uint, date time.Time) (*models.EmployeeAttendance, error)
	GetByEmployeeID(employeeID, companyID uint, offset, limit int) ([]models.EmployeeAttendance, int64, error)
	GetByCompanyID(companyID uint, offset, limit int) ([]models.EmployeeAttendance, int64, error)
	GetByDateRange(companyID uint, fromDate, toDate time.Time, offset, limit int) ([]models.EmployeeAttendance, int64, error)
	GetByEmployeeAndDateRange(employeeID, companyID uint, fromDate, toDate time.Time, offset, limit int) ([]models.EmployeeAttendance, int64, error)
	GetByEmployeeAndDateRangeNoLimit(employeeID uint, fromDate, toDate time.Time) ([]models.EmployeeAttendance, error)
	Update(attendance *models.EmployeeAttendance) error
	Delete(id uint) error
	DeleteByEmployeeAndDate(employeeID uint, date time.Time) error
	GetAttendanceStats(companyID uint, fromDate, toDate time.Time) (map[string]interface{}, error)
}

// ProductStock repository for managing product-level inventory
type ProductStockRepository interface {
	Create(stock *models.ProductStock) error
	GetByID(id string) (*models.ProductStock, error)
	GetByProductID(productID string) (*models.ProductStock, error)
	Update(stock *models.ProductStock) error
	Delete(id string) error
	GetAll(offset, limit int) ([]models.ProductStock, int64, error)
	GetAllByUser(userID uint, offset, limit int) ([]models.ProductStock, int64, error)
	GetAllByUserWithRawFilter(userID uint, offset, limit int, includeRaw bool) ([]models.ProductStock, int64, error)
	GetByProductIDs(productIDs []string) ([]models.ProductStock, error)
	GetLowStockProducts(threshold float64, offset, limit int) ([]models.ProductStock, int64, error)
	GetLowStockProductsByUser(userID uint, threshold float64, offset, limit int) ([]models.ProductStock, int64, error)
	GetDamagedProducts(offset, limit int) ([]models.ProductStock, int64, error)
	GetDamagedProductsByUser(userID uint, offset, limit int) ([]models.ProductStock, int64, error)
	GetByProductIDAndCompany(
		productID string,
		companyID uint,
		shouldFilter bool,
	) (*models.ProductStock, error)
	GetAllByCompanyWithRawFilter(
		companyID uint,
		shouldFilter bool,
		offset int,
		limit int,
		includeRaw bool,
	) ([]models.ProductStock, int64, error)
	UpdateByCompany(
		stock *models.ProductStock,
		companyID uint,
		shouldFilter bool,
	) error
	GetDamagedProductsByCompany(
		companyID uint,
		shouldFilter bool,
		offset int,
		limit int,
	) ([]models.ProductStock, int64, error)
}

// StockLedger repository for tracking all stock movements
type StockLedgerRepository interface {
	Create(ledger *models.StockLedger) error
	GetByID(id uint) (*models.StockLedger, error)
	GetByProductID(productID string, offset, limit int) ([]models.StockLedger, int64, error)
	GetByReferenceID(referenceID string) ([]models.StockLedger, error)
	DeleteByReferenceID(referenceID string) error
	GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.StockLedger, int64, error)
	GetProductMovementHistory(productID string, offset, limit int) ([]models.StockLedger, int64, error)
	CreateForCompany(
		ledger *models.StockLedger,
		companyID uint,
		shouldFilter bool,
	) error
	GetProductMovementHistoryByCompany(
		productID string,
		companyID uint,
		shouldFilter bool,
		offset int,
		limit int,
	) ([]models.StockLedger, int64, error)
}

// VariantStock repository for managing variant-level inventory
type VariantStockRepository interface {
	Create(stock *models.VariantStock) error
	GetByID(id string) (*models.VariantStock, error)
	GetBySKU(sku string) (*models.VariantStock, error)
	GetByProductID(productID string, offset, limit int) ([]models.VariantStock, int64, error)
	Update(stock *models.VariantStock) error
	Delete(id string) error
	GetAll(offset, limit int) ([]models.VariantStock, int64, error)
	GetAllByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error)
	GetAllByUserWithRawFilter(userID uint, offset, limit int, includeRaw bool) ([]models.VariantStock, int64, error)
	GetBySKUs(skus []string) ([]models.VariantStock, error)
	GetLowStockVariants(threshold float64, offset, limit int) ([]models.VariantStock, int64, error)
	GetDamagedVariants(offset, limit int) ([]models.VariantStock, int64, error)
	GetDamagedVariantsByUser(userID uint, offset, limit int) ([]models.VariantStock, int64, error)
}

// VariantStockMovement repository for tracking variant stock movements
type VariantStockMovementRepository interface {
	Create(movement *models.VariantStockMovement) error
	GetByID(id uint) (*models.VariantStockMovement, error)
	GetByVariantSKU(sku string, offset, limit int) ([]models.VariantStockMovement, int64, error)
	GetByReferenceID(referenceID string) ([]models.VariantStockMovement, error)
	GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.VariantStockMovement, int64, error)
	DeleteByReferenceID(referenceID string) error
}

// StockReservation repository for managing reserved stock
type StockReservationRepository interface {
	Create(reservation *models.StockReservation) error
	GetByID(id uint) (*models.StockReservation, error)
	GetBySalesOrderID(salesOrderID string) ([]models.StockReservation, error)
	GetByVariantSKU(sku string, offset, limit int) ([]models.StockReservation, int64, error)
	GetByStatus(status string, offset, limit int) ([]models.StockReservation, int64, error)
	Update(reservation *models.StockReservation) error
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
}

// ProductGroupInventory repository for managing product group stock levels
type ProductGroupInventoryRepository interface {
	Create(inventory *models.ProductGroupInventory) error
	FindByID(id uint) (*models.ProductGroupInventory, error)
	FindByProductGroupID(productGroupID string) (*models.ProductGroupInventory, error)
	Update(inventory *models.ProductGroupInventory) error
	Delete(id uint) error
	GetLowStockGroups(threshold float64) ([]models.ProductGroupInventory, error)
}

// ComponentInventory repository for managing component stock within product groups
type ComponentInventoryRepository interface {
	Create(inventory *models.ComponentInventory) error
	FindByID(id uint) (*models.ComponentInventory, error)
	FindByProductGroupID(productGroupID string) ([]models.ComponentInventory, error)
	FindByComponentProductID(productID string) ([]models.ComponentInventory, error)
	Update(inventory *models.ComponentInventory) error
	Delete(id uint) error
	UpdateBatch(inventories []models.ComponentInventory) error
}

// ProductGroupTransaction repository for tracking inventory movements
type ProductGroupTransactionRepository interface {
	Create(transaction *models.ProductGroupTransaction) error
	FindByID(id uint) (*models.ProductGroupTransaction, error)
	FindByProductGroupID(productGroupID string, limit, offset int) ([]models.ProductGroupTransaction, int64, error)
	FindByReferenceID(referenceID string) ([]models.ProductGroupTransaction, error)
	GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.ProductGroupTransaction, int64, error)
	Delete(id uint) error
}

// VendorPaymentRepository interface for vendor payment operations
type VendorPaymentRepository interface {
	Create(vp *models.VendorPayment) (*models.VendorPayment, error)
	FindByID(id uint) (*models.VendorPayment, error)
	FindByPaymentNumber(paymentNumber string) (*models.VendorPayment, error)
	FindByPurchaseOrderID(poID string, limit, offset int) ([]models.VendorPayment, int64, error)
	FindByVendorID(vendorID uint, limit, offset int) ([]models.VendorPayment, int64, error)
	FindAll(limit, offset int) ([]models.VendorPayment, int64, error)
	FindByPaymentStatus(status string, limit, offset int) ([]models.VendorPayment, int64, error)
	Update(id uint, vp *models.VendorPayment) (*models.VendorPayment, error)
	UpdatePaymentStatus(id uint, status string, paidAmount, remainingAmount float64) error
	Delete(id uint) error
	GetDB() *gorm.DB
}

// CustomerPaymentRepository interface for customer payment operations
type CustomerPaymentRepository interface {
	Create(cp *models.CustomerPayment) (*models.CustomerPayment, error)
	FindByID(id uint) (*models.CustomerPayment, error)
	FindByPaymentNumber(paymentNumber string) (*models.CustomerPayment, error)
	FindBySalesOrderID(soID string, limit, offset int) ([]models.CustomerPayment, int64, error)
	FindByCustomerID(customerID uint, limit, offset int) ([]models.CustomerPayment, int64, error)
	FindAll(limit, offset int) ([]models.CustomerPayment, int64, error)
	FindByPaymentStatus(status string, limit, offset int) ([]models.CustomerPayment, int64, error)
	Update(id uint, cp *models.CustomerPayment) (*models.CustomerPayment, error)
	UpdatePaymentStatus(id uint, status string, receivedAmount, remainingAmount float64) error
	Delete(id uint) error
	GetDB() *gorm.DB
}

// SalaryRepository interface for salary calculation operations
type SalaryRepository interface {
	Create(salary *models.SalaryCalculation) error
	GetByID(id uint) (*models.SalaryCalculation, error)
	GetByEmployee(employeeID uint) ([]models.SalaryCalculation, error)
	GetByEmployeeAndMonth(employeeID uint, month, year int) (*models.SalaryCalculation, error)
	GetByCompany(companyID uint, limit, offset int) ([]models.SalaryCalculation, int64, error)
	Update(salary *models.SalaryCalculation) error
	Delete(id uint) error
}

// ProductConversionRepository interface for managing product conversions
type ProductConversionRepository interface {
	Create(conversion *models.ProductConversion) error
	GetByID(id string) (*models.ProductConversion, error)
	GetAll(offset, limit int) ([]models.ProductConversion, int64, error)
	GetActiveConversions(offset, limit int) ([]models.ProductConversion, int64, error)
	GetByRawProductID(rawProductID string, offset, limit int) ([]models.ProductConversion, int64, error)
	GetByFinishedProductID(finishedProductID string, offset, limit int) ([]models.ProductConversion, int64, error)
	GetByProductPair(rawProductID, finishedProductID string) (*models.ProductConversion, error)
	Update(conversion *models.ProductConversion) error
	Delete(id string) error
	GetDB() *gorm.DB

	CreateForCompany(*models.ProductConversion, uint) error
	GetByIDAndCompany(string, uint) (*models.ProductConversion, error)
	GetAllByCompany(uint, int, int) ([]models.ProductConversion, int64, error)
	GetActiveConversionsByCompany(uint, int, int) ([]models.ProductConversion, int64, error)
	GetByRawProductIDAndCompany(string, uint, int, int) ([]models.ProductConversion, int64, error)
	GetByFinishedProductIDAndCompany(string, uint, int, int) ([]models.ProductConversion, int64, error)
	GetByProductPairAndCompany(string, string, uint) (*models.ProductConversion, error)
	UpdateByCompany(*models.ProductConversion, uint) error
	DeleteByCompany(string, uint) error
}

// ProductConversionRecordRepository interface for managing conversion records
type ProductConversionRecordRepository interface {
	Create(record *models.ProductConversionRecord) error
	GetByID(id string) (*models.ProductConversionRecord, error)
	GetByConversionID(conversionID string, offset, limit int) ([]models.ProductConversionRecord, int64, error)
	GetAll(offset, limit int) ([]models.ProductConversionRecord, int64, error)
	GetByDateRange(fromDate, toDate time.Time, offset, limit int) ([]models.ProductConversionRecord, int64, error)
	GetByStatus(status string, offset, limit int) ([]models.ProductConversionRecord, int64, error)
	Update(record *models.ProductConversionRecord) error
	Delete(id string) error
	GetDB() *gorm.DB

	CreateForCompany(*models.ProductConversionRecord, uint) error
	GetByIDAndCompany(string, uint) (*models.ProductConversionRecord, error)
	GetByConversionIDAndCompany(string, uint, int, int) ([]models.ProductConversionRecord, int64, error)
	GetAllByCompany(uint, int, int) ([]models.ProductConversionRecord, int64, error)
	GetByDateRangeAndCompany(time.Time, time.Time, uint, int, int) ([]models.ProductConversionRecord, int64, error)
	GetByStatusAndCompany(string, uint, int, int) ([]models.ProductConversionRecord, int64, error)
	UpdateByCompany(*models.ProductConversionRecord, uint) error
	DeleteByCompany(string, uint) error
}

type ConversionRecordBagUsageRepository interface {
	Create(bagUsage *models.ConversionRecordBagUsage) error
	GetByConversionRecordID(recordID string) ([]models.ConversionRecordBagUsage, error)
	GetByBagID(bagID string) ([]models.ConversionRecordBagUsage, error)
	GetDB() *gorm.DB

	CreateForCompany(*models.ConversionRecordBagUsage, uint) error
	GetByConversionRecordIDAndCompany(string, uint) ([]models.ConversionRecordBagUsage, error)
	GetByBagIDAndCompany(string, uint) ([]models.ConversionRecordBagUsage, error)
}
type RawMaterialBagRepository interface {
	// Existing methods retained.
	CreateMany(bags []models.RawMaterialBag) error
	GetAll(
		limit int,
		offset int,
	) ([]models.RawMaterialBag, int64, error)
	GetByID(
		id string,
	) (*models.RawMaterialBag, error)
	GetByProductID(
		productID string,
	) ([]models.RawMaterialBag, error)
	GetByPurchaseOrderID(
		purchaseOrderID string,
	) ([]models.RawMaterialBag, error)
	Update(
		bag *models.RawMaterialBag,
	) error

	// Company-scoped methods.
	CreateManyForCompany(
		bags []models.RawMaterialBag,
		companyID uint,
	) error

	GetAllByCompany(
		companyID uint,
		limit int,
		offset int,
	) ([]models.RawMaterialBag, int64, error)

	GetByIDAndCompany(
		id string,
		companyID uint,
	) (*models.RawMaterialBag, error)

	GetByProductIDAndCompany(
		productID string,
		companyID uint,
	) ([]models.RawMaterialBag, error)

	GetByPurchaseOrderIDAndCompany(
		purchaseOrderID string,
		companyID uint,
	) ([]models.RawMaterialBag, error)

	UpdateByCompany(
		bag *models.RawMaterialBag,
		companyID uint,
	) error
}


type VendorShortageClaimRepository interface {
	// Existing methods retained.
	Create(
		claim *models.VendorShortageClaim,
	) error

	FindByPurchaseOrderID(
		purchaseOrderID string,
	) ([]models.VendorShortageClaim, error)

	// Company-scoped methods.
	CreateForCompany(
		claim *models.VendorShortageClaim,
		companyID uint,
	) error

	FindByPurchaseOrderIDAndCompany(
		purchaseOrderID string,
		companyID uint,
	) ([]models.VendorShortageClaim, error)
}