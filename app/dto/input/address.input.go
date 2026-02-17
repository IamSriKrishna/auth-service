package input

type AddressInput struct {
	// Primary field names
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

	// Alternative field names (for flexibility)
	Street     string `json:"street"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// Normalize maps alternative field names to primary fields
func (a *AddressInput) Normalize() {
	// Map street to address_line1 if address_line1 is empty
	if a.AddressLine1 == "" && a.Street != "" {
		a.AddressLine1 = a.Street
	}

	// Map country to country_region if country_region is empty
	if a.CountryRegion == "" && a.Country != "" {
		a.CountryRegion = a.Country
	}

	// Map postal_code to pin_code if pin_code is empty
	if a.PinCode == "" && a.PostalCode != "" {
		a.PinCode = a.PostalCode
	}
}
