package input

type OtherDetailsInput struct {
	PAN              string `json:"pan" validate:"omitempty,len=10"`
	IsMSMERegistered bool   `json:"is_msme_registered"`
	Currency         string `json:"currency"`
	PaymentTerms     string `json:"payment_terms"`
	TDS              string `json:"tds"`
	EnablePortal     bool   `json:"enable_portal"`
	WebsiteURL       string `json:"website_url" validate:"omitempty,url"`
	Department       string `json:"department"`
	Designation      string `json:"designation"`
	Twitter          string `json:"twitter"`
	SkypeName        string `json:"skype_name"`
	Facebook         string `json:"facebook"`
}

type OtherDetailsCustomerInput struct {
	PAN              string `json:"pan" validate:"omitempty,len=10"`
	Currency         string `json:"currency"`
	PaymentTerms     string `json:"payment_terms"`
	EnablePortal     bool   `json:"enable_portal"`
}
