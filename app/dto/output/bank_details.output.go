package output

import "time"

type BankDetailOutput struct {
	ID                uint      `json:"id"`
	VendorID          uint      `json:"vendor_id"`
	AccountHolderName string    `json:"account_holder_name"`
	BankName          string    `json:"bank_name"`
	AccountNumber     string    `json:"account_number"`
	IFSC              string    `json:"ifsc"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
