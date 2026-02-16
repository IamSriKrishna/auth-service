package input

type CreateCompanyInput struct {
	CompanyName    string `json:"company_name" validate:"required,min=1,max=255"`
	BusinessTypeID uint   `json:"business_type_id" validate:"required"`
	GSTNumber      string `json:"gst_number,omitempty" validate:"omitempty,len=15"`
	PANNumber      string `json:"pan_number,omitempty" validate:"omitempty,len=10"`
}

type UpdateCompanyInput struct {
	CompanyName    *string `json:"company_name,omitempty" validate:"omitempty,min=1,max=255"`
	BusinessTypeID *uint   `json:"business_type_id,omitempty"`
	GSTNumber      *string `json:"gst_number,omitempty" validate:"omitempty,len=15"`
	PANNumber      *string `json:"pan_number,omitempty" validate:"omitempty,len=10"`
}

type UpsertCompanyContactInput struct {
	Mobile          string `json:"mobile" validate:"required,min=10,max=15"`
	AlternateMobile string `json:"alternate_mobile,omitempty" validate:"omitempty,min=10,max=15"`
	Email           string `json:"email" validate:"required,email,max=255"`
}

type UpsertCompanyAddressInput struct {
	AddressLine1 string `json:"address_line1" validate:"required,min=1,max=255"`
	AddressLine2 string `json:"address_line2,omitempty" validate:"omitempty,max=255"`
	City         string `json:"city" validate:"required,min=1,max=100"`
	StateID      uint   `json:"state_id" validate:"required"`
	CountryID    uint   `json:"country_id" validate:"required"`
	Pincode      string `json:"pincode" validate:"required,min=4,max=10"`
}

type CreateBankDetailInput struct {
	BankName          string `json:"bank_name" validate:"required,min=1,max=255"`
	AccountHolderName string `json:"account_holder_name" validate:"required,min=1,max=255"`
	AccountNumber     string `json:"account_number" validate:"required,min=1,max=50"`
	IFSCCode          string `json:"ifsc_code" validate:"required,len=11"`
	BranchName        string `json:"branch_name,omitempty" validate:"omitempty,max=255"`
	IsPrimary         bool   `json:"is_primary"`
}

type UpdateBankDetailInput struct {
	BankName          *string `json:"bank_name,omitempty" validate:"omitempty,min=1,max=255"`
	AccountHolderName *string `json:"account_holder_name,omitempty" validate:"omitempty,min=1,max=255"`
	AccountNumber     *string `json:"account_number,omitempty" validate:"omitempty,min=1,max=50"`
	IFSCCode          *string `json:"ifsc_code,omitempty" validate:"omitempty,len=11"`
	BranchName        *string `json:"branch_name,omitempty" validate:"omitempty,max=255"`
	IsPrimary         *bool   `json:"is_primary,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

type UpsertUPIDetailInput struct {
	UPIID    string `json:"upi_id" validate:"required,min=1,max=255"`
	UPIQRURL string `json:"upi_qr_url,omitempty"`
}

type UpsertInvoiceSettingsInput struct {
	InvoicePrefix      string `json:"invoice_prefix" validate:"required,min=1,max=10"`
	InvoiceStartNumber int    `json:"invoice_start_number" validate:"required,min=1"`
	ShowLogo           bool   `json:"show_logo"`
	ShowSignature      bool   `json:"show_signature"`
	RoundOffTotal      bool   `json:"round_off_total"`
}

type UpsertTaxSettingsInput struct {
	GSTEnabled bool `json:"gst_enabled"`
	TaxTypeID  uint `json:"tax_type_id" validate:"required"`
}

type UpsertRegionalSettingsInput struct {
	Timezone       string `json:"timezone" validate:"required,max=50"`
	DateFormat     string `json:"date_format" validate:"required,max=20"`
	TimeFormat     string `json:"time_format" validate:"required,max=20"`
	CurrencyCode   string `json:"currency_code" validate:"required,len=3"`
	CurrencySymbol string `json:"currency_symbol" validate:"required,max=10"`
	LanguageCode   string `json:"language_code" validate:"required,max=5"`
}

type CompleteCompanySetupInput struct {
	Company          CreateCompanyInput           `json:"company" validate:"required"`
	Contact          UpsertCompanyContactInput    `json:"contact" validate:"required"`
	Address          UpsertCompanyAddressInput    `json:"address" validate:"required"`
	BankDetails      *CreateBankDetailInput       `json:"bank_details,omitempty"`
	UPIDetails       *UpsertUPIDetailInput        `json:"upi_details,omitempty"`
	InvoiceSettings  *UpsertInvoiceSettingsInput  `json:"invoice_settings,omitempty"`
	TaxSettings      *UpsertTaxSettingsInput      `json:"tax_settings,omitempty"`
	RegionalSettings *UpsertRegionalSettingsInput `json:"regional_settings,omitempty"`
}
