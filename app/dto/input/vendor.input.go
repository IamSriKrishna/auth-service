package input

type CreateVendorInput struct {
	Salutation      string               `json:"salutation" validate:"required"`
	FirstName       string               `json:"first_name" validate:"required"`
	LastName        string               `json:"last_name"`
	CompanyName     string               `json:"company_name"`
	DisplayName     string               `json:"display_name" validate:"required"`
	EmailAddress    string               `json:"email_address" validate:"omitempty,email"`
	WorkPhone       string               `json:"work_phone"`
	WorkPhoneCode   string               `json:"work_phone_code"`
	Mobile          string               `json:"mobile"`
	MobileCode      string               `json:"mobile_code"`
	VendorLanguage  string               `json:"vendor_language"`
	OtherDetails    *OtherDetailsInput   `json:"other_details"`
	BillingAddress  *AddressInput        `json:"billing_address"`
	ShippingAddress *AddressInput        `json:"shipping_address"`
	ContactPersons  []ContactPersonInput `json:"contact_persons"`
	BankDetails     []BankDetailInput    `json:"bank_details"`
}

type UpdateVendorInput struct {
	Salutation      *string              `json:"salutation"`
	FirstName       *string              `json:"first_name"`
	LastName        *string              `json:"last_name"`
	CompanyName     *string              `json:"company_name"`
	DisplayName     *string              `json:"display_name"`
	EmailAddress    *string              `json:"email_address" validate:"omitempty,email"`
	WorkPhone       *string              `json:"work_phone"`
	WorkPhoneCode   *string              `json:"work_phone_code"`
	Mobile          *string              `json:"mobile"`
	MobileCode      *string              `json:"mobile_code"`
	VendorLanguage  *string              `json:"vendor_language"`
	OtherDetails    *OtherDetailsInput   `json:"other_details"`
	BillingAddress  *AddressInput        `json:"billing_address"`
	ShippingAddress *AddressInput        `json:"shipping_address"`
	ContactPersons  []ContactPersonInput `json:"contact_persons"`
	BankDetails     []BankDetailInput    `json:"bank_details"`
}
