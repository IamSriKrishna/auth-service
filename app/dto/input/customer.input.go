package input

type CreateCustomerInput struct {
	Salutation       string                     `json:"salutation" validate:"required"`
	FirstName        string                     `json:"first_name" validate:"required"`
	LastName         string                     `json:"last_name"`
	CompanyName      string                     `json:"company_name"`
	DisplayName      string                     `json:"display_name" validate:"required"`
	EmailAddress     string                     `json:"email_address" validate:"omitempty,email"`
	WorkPhone        string                     `json:"work_phone"`
	WorkPhoneCode    string                     `json:"work_phone_code"`
	Mobile           string                     `json:"mobile"`
	MobileCode       string                     `json:"mobile_code"`
	CustomerLanguage string                     `json:"customer_language"`
	OtherDetails     *OtherDetailsCustomerInput `json:"other_details"`
	BillingAddress   *AddressInput              `json:"billing_address"`
	ShippingAddress  *AddressInput              `json:"shipping_address"`
	ContactPersons   []ContactPersonInput       `json:"contact_persons"`
}

type UpdateCustomerInput struct {
	Salutation      *string                    `json:"salutation"`
	FirstName       *string                    `json:"first_name"`
	LastName        *string                    `json:"last_name"`
	CompanyName     *string                    `json:"company_name"`
	DisplayName     *string                    `json:"display_name"`
	EmailAddress    *string                    `json:"email_address" validate:"omitempty,email"`
	WorkPhone       *string                    `json:"work_phone"`
	WorkPhoneCode   *string                    `json:"work_phone_code"`
	Mobile          *string                    `json:"mobile"`
	MobileCode      *string                    `json:"mobile_code"`
	CustomerLanguage *string                   `json:"customer_language"`
	OtherDetails    *OtherDetailsCustomerInput `json:"other_details"`
	BillingAddress  *AddressInput              `json:"billing_address"`
	ShippingAddress *AddressInput              `json:"shipping_address"`
	ContactPersons  []ContactPersonInput       `json:"contact_persons"`
}
