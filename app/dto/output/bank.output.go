package output

import "time"

// BankResponse represents a single bank response
type BankResponse struct {
	ID        uint      `json:"id"`
	BankCode  *string   `json:"bank_code,omitempty"`
	BankName  string    `json:"bank_name"`
	ShortName *string   `json:"short_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// BankListItem is an item in bank list responses
type BankListItem = BankResponse
