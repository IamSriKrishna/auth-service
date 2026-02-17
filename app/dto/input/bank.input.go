package input

type CreateBankInput struct {
	BankName   string `json:"bank_name" validate:"required,min=1,max=255"`
	IFSCCode   string `json:"ifsc_code" validate:"required,min=1,max=11"`
	BranchName string `json:"branch_name" validate:"omitempty,max=255"`
	BranchCode string `json:"branch_code" validate:"omitempty,max=50"`
	Address    string `json:"address" validate:"omitempty"`
	City       string `json:"city" validate:"omitempty,max=100"`
	State      string `json:"state" validate:"omitempty,max=100"`
	PostalCode string `json:"postal_code" validate:"omitempty,max=10"`
	Country    string `json:"country" validate:"omitempty,max=100"`
	IsActive   bool   `json:"is_active"`
}

type UpdateBankInput struct {
	BankName   *string `json:"bank_name" validate:"omitempty,min=1,max=255"`
	IFSCCode   *string `json:"ifsc_code" validate:"omitempty,min=1,max=11"`
	BranchName *string `json:"branch_name" validate:"omitempty,max=255"`
	BranchCode *string `json:"branch_code" validate:"omitempty,max=50"`
	Address    *string `json:"address" validate:"omitempty"`
	City       *string `json:"city" validate:"omitempty,max=100"`
	State      *string `json:"state" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code" validate:"omitempty,max=10"`
	Country    *string `json:"country" validate:"omitempty,max=100"`
	IsActive   *bool   `json:"is_active"`
}

type BankDetailInput struct {
	BankID               uint   `json:"bank_id" validate:"required"`
	AccountHolderName    string `json:"account_holder_name"`
	AccountNumber        string `json:"account_number" validate:"required"`
	ReenterAccountNumber string `json:"reenter_account_number" validate:"required,eqfield=AccountNumber"`
}
