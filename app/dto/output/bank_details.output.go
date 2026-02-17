package output

import "time"

type BankDetailOutput struct {
	ID                uint      `json:"id"`
	VendorID          uint      `json:"vendor_id"`
	BankID            uint      `json:"bank_id"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountNumber     string    `json:"account_number"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
