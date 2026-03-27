package output

import "time"

type CustomerUserOutput struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CustomerCompanyOutput struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CompanyCode string `json:"company_code"`
}

type CustomerOutput struct {
	ID               uint                          `json:"id"`
	Salutation       string                        `json:"salutation"`
	FirstName        string                        `json:"first_name"`
	LastName         string                        `json:"last_name"`
	DisplayName      string                        `json:"display_name"`
	EmailAddress     string                        `json:"email_address"`
	WorkPhone        string                        `json:"work_phone"`
	WorkPhoneCode    string                        `json:"work_phone_code"`
	Mobile           string                        `json:"mobile"`
	MobileCode       string                        `json:"mobile_code"`
	CustomerLanguage string                        `json:"customer_language"`
	User             *CustomerUserOutput           `json:"user"`
	Company          *CustomerCompanyOutput        `json:"company"`
	OtherDetails     *CustomerOtherDetailsOutput   `json:"other_details,omitempty"`
	BillingAddress   *CustomerAddressOutput        `json:"billing_address,omitempty"`
	ShippingAddress  *CustomerAddressOutput        `json:"shipping_address,omitempty"`
	ContactPersons   []CustomerContactPersonOutput `json:"contact_persons,omitempty"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
	UserID           uint                          `json:"user_id,omitempty"`
	UserName         string                        `json:"user_name,omitempty"`
	CompanyID        uint                          `json:"company_id,omitempty"`
	CompanyName      string                        `json:"company_name,omitempty"`
}

type CustomerListOutput struct {
	ID               uint      `json:"id"`
	DisplayName      string    `json:"display_name"`
	EmailAddress     string    `json:"email_address"`
	WorkPhone        string    `json:"work_phone"`
	Mobile           string    `json:"mobile"`
	CustomerLanguage string    `json:"customer_language"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
