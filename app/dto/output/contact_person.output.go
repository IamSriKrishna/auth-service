package output

import "time"

type ContactPersonOutput struct {
	ID            uint      `json:"id"`
	VendorID      uint      `json:"vendor_id"`
	Salutation    string    `json:"salutation"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	EmailAddress  string    `json:"email_address"`
	WorkPhone     string    `json:"work_phone"`
	WorkPhoneCode string    `json:"work_phone_code"`
	Mobile        string    `json:"mobile"`
	MobileCode    string    `json:"mobile_code"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CustomerContactPersonOutput struct {
	ID            uint      `json:"id"`
	CustomerID    uint      `json:"customer_id"`
	Salutation    string    `json:"salutation"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	EmailAddress  string    `json:"email_address"`
	WorkPhone     string    `json:"work_phone"`
	WorkPhoneCode string    `json:"work_phone_code"`
	Mobile        string    `json:"mobile"`
	MobileCode    string    `json:"mobile_code"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
