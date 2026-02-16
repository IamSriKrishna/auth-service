package output

import "time"

type AddressOutput struct {
	ID            uint      `json:"id"`
	VendorID      uint      `json:"vendor_id"`
	AddressType   string    `json:"address_type"`
	Attention     string    `json:"attention"`
	CountryRegion string    `json:"country_region"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PinCode       string    `json:"pin_code"`
	Phone         string    `json:"phone"`
	PhoneCode     string    `json:"phone_code"`
	FaxNumber     string    `json:"fax_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CustomerAddressOutput struct {
	ID            uint      `json:"id"`
	CustomerID    uint      `json:"customer_id"`
	AddressType   string    `json:"address_type"`
	Attention     string    `json:"attention"`
	CountryRegion string    `json:"country_region"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PinCode       string    `json:"pin_code"`
	Phone         string    `json:"phone"`
	PhoneCode     string    `json:"phone_code"`
	FaxNumber     string    `json:"fax_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
