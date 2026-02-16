package output

import (
	"time"
)

type VendorOutput struct {
	ID              uint                  `json:"id"`
	Salutation      string                `json:"salutation"`
	FirstName       string                `json:"first_name"`
	LastName        string                `json:"last_name"`
	CompanyName     string                `json:"company_name"`
	DisplayName     string                `json:"display_name"`
	EmailAddress    string                `json:"email_address"`
	WorkPhone       string                `json:"work_phone"`
	WorkPhoneCode   string                `json:"work_phone_code"`
	Mobile          string                `json:"mobile"`
	MobileCode      string                `json:"mobile_code"`
	VendorLanguage  string                `json:"vendor_language"`
	OtherDetails    *OtherDetailsOutput    `json:"other_details,omitempty"`
	BillingAddress  *AddressOutput         `json:"billing_address,omitempty"`
	ShippingAddress *AddressOutput         `json:"shipping_address,omitempty"`
	ContactPersons  []ContactPersonOutput `json:"contact_persons,omitempty"`
	BankDetails     []BankDetailOutput    `json:"bank_details,omitempty"`
	Documents       []DocumentOutput      `json:"documents,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

type VendorListOutput struct {
	ID             uint      `json:"id"`
	DisplayName    string    `json:"display_name"`
	CompanyName    string    `json:"company_name"`
	EmailAddress   string    `json:"email_address"`
	WorkPhone      string    `json:"work_phone"`
	Mobile         string    `json:"mobile"`
	VendorLanguage string    `json:"vendor_language"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
