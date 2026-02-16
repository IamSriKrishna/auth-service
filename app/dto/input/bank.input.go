package input

type BankDetailInput struct {
	AccountHolderName    string `json:"account_holder_name"`
	BankName             string `json:"bank_name"`
	AccountNumber        string `json:"account_number" validate:"required"`
	ReenterAccountNumber string `json:"reenter_account_number" validate:"required,eqfield=AccountNumber"`
	IFSC                 string `json:"ifsc" validate:"required,len=11"`
}
