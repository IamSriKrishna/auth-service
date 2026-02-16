package input

type ContactPersonInput struct {
	Salutation    string `json:"salutation"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmailAddress  string `json:"email_address" validate:"omitempty,email"`
	WorkPhone     string `json:"work_phone"`
	WorkPhoneCode string `json:"work_phone_code"`
	Mobile        string `json:"mobile"`
	MobileCode    string `json:"mobile_code"`
}
