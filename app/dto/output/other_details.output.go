package output

import "time"

type OtherDetailsOutput struct {
	ID               uint      `json:"id"`
	VendorID         uint      `json:"vendor_id"`
	PAN              string    `json:"pan"`
	IsMSMERegistered bool      `json:"is_msme_registered"`
	Currency         string    `json:"currency"`
	PaymentTerms     string    `json:"payment_terms"`
	TDS              string    `json:"tds"`
	EnablePortal     bool      `json:"enable_portal"`
	WebsiteURL       string    `json:"website_url"`
	Department       string    `json:"department"`
	Designation      string    `json:"designation"`
	Twitter          string    `json:"twitter"`
	SkypeName        string    `json:"skype_name"`
	Facebook         string    `json:"facebook"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CustomerOtherDetailsOutput struct {
	ID           uint      `json:"id"`
	CustomerID   uint      `json:"customer_id"`
	PAN          string    `json:"pan"`
	Currency     string    `json:"currency"`
	PaymentTerms string    `json:"payment_terms"`
	EnablePortal bool      `json:"enable_portal"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
