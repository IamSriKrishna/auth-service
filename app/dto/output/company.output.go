package output

import "time"

type BusinessTypeOutput struct {
	ID          uint      `json:"id"`
	TypeName    string    `json:"type_name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CompanyOutput struct {
	ID             uint               `json:"id"`
	CompanyName    string             `json:"company_name"`
	BusinessTypeID uint               `json:"business_type_id"`
	BusinessType   BusinessTypeOutput `json:"business_type"`
	GSTNumber      string             `json:"gst_number,omitempty"`
	PANNumber      string             `json:"pan_number,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type CountryOutput struct {
	ID          uint   `json:"id"`
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
	PhoneCode   string `json:"phone_code,omitempty"`
}

type StateOutput struct {
	ID        uint          `json:"id"`
	CountryID uint          `json:"country_id"`
	StateName string        `json:"state_name"`
	StateCode string        `json:"state_code,omitempty"`
	Country   CountryOutput `json:"country,omitempty"`
}

type CompanyContactOutput struct {
	ID              uint      `json:"id"`
	CompanyID       uint      `json:"company_id"`
	Mobile          string    `json:"mobile"`
	AlternateMobile string    `json:"alternate_mobile,omitempty"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CompanyAddressOutput struct {
	ID           uint          `json:"id"`
	CompanyID    uint          `json:"company_id"`
	AddressLine1 string        `json:"address_line1"`
	AddressLine2 string        `json:"address_line2,omitempty"`
	City         string        `json:"city"`
	StateID      uint          `json:"state_id"`
	State        StateOutput   `json:"state"`
	CountryID    uint          `json:"country_id"`
	Country      CountryOutput `json:"country"`
	Pincode      string        `json:"pincode"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type CompanyBankDetailOutput struct {
	ID                uint                   `json:"id"`
	CompanyID         uint                   `json:"company_id"`
	BankID            uint                   `json:"bank_id"`
	Bank              map[string]interface{} `json:"bank,omitempty"`
	AccountHolderName string                 `json:"account_holder_name"`
	AccountNumber     string                 `json:"account_number"`
	IsPrimary         bool                   `json:"is_primary"`
	IsActive          bool                   `json:"is_active"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type CompanyUPIDetailOutput struct {
	ID        uint      `json:"id"`
	CompanyID uint      `json:"company_id"`
	UPIID     string    `json:"upi_id"`
	UPIQRURL  string    `json:"upi_qr_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaxTypeOutput struct {
	ID          uint   `json:"id"`
	TaxName     string `json:"tax_name"`
	TaxCode     string `json:"tax_code"`
	Description string `json:"description,omitempty"`
}

type CompanyInvoiceSettingsOutput struct {
	ID                   uint      `json:"id"`
	CompanyID            uint      `json:"company_id"`
	InvoicePrefix        string    `json:"invoice_prefix"`
	InvoiceStartNumber   int       `json:"invoice_start_number"`
	CurrentInvoiceNumber int       `json:"current_invoice_number"`
	ShowLogo             bool      `json:"show_logo"`
	ShowSignature        bool      `json:"show_signature"`
	RoundOffTotal        bool      `json:"round_off_total"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CompanyTaxSettingsOutput struct {
	ID         uint          `json:"id"`
	CompanyID  uint          `json:"company_id"`
	GSTEnabled bool          `json:"gst_enabled"`
	TaxTypeID  uint          `json:"tax_type_id"`
	TaxType    TaxTypeOutput `json:"tax_type"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type CompanyRegionalSettingsOutput struct {
	ID             uint      `json:"id"`
	CompanyID      uint      `json:"company_id"`
	Timezone       string    `json:"timezone"`
	DateFormat     string    `json:"date_format"`
	TimeFormat     string    `json:"time_format"`
	CurrencyCode   string    `json:"currency_code"`
	CurrencySymbol string    `json:"currency_symbol"`
	LanguageCode   string    `json:"language_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CompleteCompanyProfileOutput struct {
	Company          CompanyOutput                  `json:"company"`
	Contact          *CompanyContactOutput          `json:"contact,omitempty"`
	Address          *CompanyAddressOutput          `json:"address,omitempty"`
	BankDetails      []CompanyBankDetailOutput      `json:"bank_details,omitempty"`
	UPIDetails       *CompanyUPIDetailOutput        `json:"upi_details,omitempty"`
	InvoiceSettings  *CompanyInvoiceSettingsOutput  `json:"invoice_settings,omitempty"`
	TaxSettings      *CompanyTaxSettingsOutput      `json:"tax_settings,omitempty"`
	RegionalSettings *CompanyRegionalSettingsOutput `json:"regional_settings,omitempty"`
}

type CompanyPaginatedResponse struct {
	Data       []CompleteCompanyProfileOutput `json:"data"`
	Page       int                            `json:"page"`
	PageSize   int                            `json:"page_size"`
	TotalCount int64                          `json:"total_count"`
	TotalPages int                            `json:"total_pages"`
}

type ListCompaniesQuery struct {
	BusinessTypeID *uint   `json:"business_type_id,omitempty"`
	Search         *string `json:"search,omitempty"`
	IsActive       *bool   `json:"is_active,omitempty"`
	Page           int     `json:"page"`
	PageSize       int     `json:"page_size"`
}
