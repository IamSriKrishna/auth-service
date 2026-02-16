package input

type AddressInput struct {
	Attention     string `json:"attention"`
	CountryRegion string `json:"country_region"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	PinCode       string `json:"pin_code"`
	Phone         string `json:"phone"`
	PhoneCode     string `json:"phone_code"`
	FaxNumber     string `json:"fax_number"`
}
