package output

import "time"

type BankOutput struct {
	ID         uint      `json:"id"`
	BankName   string    `json:"bank_name"`
	Address    string    `json:"address,omitempty"`
	City       string    `json:"city,omitempty"`
	State      string    `json:"state,omitempty"`
	PostalCode string    `json:"postal_code,omitempty"`
	Country    string    `json:"country,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BankResponse struct {
	ID        uint      `json:"id"`
	BankCode  *string   `json:"bank_code,omitempty"`
	BankName  string    `json:"bank_name"`
	ShortName *string   `json:"short_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ListBankOutput struct {
	Banks []BankOutput `json:"banks"`
	Total int64        `json:"total"`
}

type BankListItem = BankResponse
