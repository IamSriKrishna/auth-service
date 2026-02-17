package input

type ContactPersonInput struct {
	// Primary field names
	Salutation    string `json:"salutation"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmailAddress  string `json:"email_address" validate:"omitempty,email"`
	WorkPhone     string `json:"work_phone"`
	WorkPhoneCode string `json:"work_phone_code"`
	Mobile        string `json:"mobile"`
	MobileCode    string `json:"mobile_code"`

	// Alternative field names (for flexibility)
	Title     string `json:"title"`
	Email     string `json:"email" validate:"omitempty,email"`
	Phone     string `json:"phone"`
	PhoneCode string `json:"phone_code"`
}

// Normalize maps alternative field names to primary fields
func (c *ContactPersonInput) Normalize() {
	// Map title to salutation if salutation is empty
	if c.Salutation == "" && c.Title != "" {
		c.Salutation = c.Title
	}

	// Map email to email_address if email_address is empty
	if c.EmailAddress == "" && c.Email != "" {
		c.EmailAddress = c.Email
	}

	// Map phone to mobile if mobile is empty
	if c.Mobile == "" && c.Phone != "" {
		c.Mobile = c.Phone
	}

	// Map phone_code to mobile_code if mobile_code is empty or default
	if (c.MobileCode == "" || c.MobileCode == "+91") && c.PhoneCode != "" && c.PhoneCode != "+91" {
		c.MobileCode = c.PhoneCode
	}
}
