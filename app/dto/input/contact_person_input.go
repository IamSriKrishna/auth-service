package input

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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

	// Alternative field names
	Title     string `json:"title"`
	Email     string `json:"email" validate:"omitempty,email"`
	Phone     string `json:"phone"`
	PhoneCode string `json:"phone_code"`
}

func parseStringOrNumber(data json.RawMessage) (string, error) {
	data = bytes.TrimSpace(data)

	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return "", nil
	}

	if data[0] == '"' {
		var value string

		if err := json.Unmarshal(data, &value); err != nil {
			return "", err
		}

		return value, nil
	}

	var number json.Number

	if err := json.Unmarshal(data, &number); err != nil {
		return "", fmt.Errorf("value must be a string or number")
	}

	return number.String(), nil
}

func (c *ContactPersonInput) UnmarshalJSON(data []byte) error {
	type contactPersonAlias struct {
		Salutation    string          `json:"salutation"`
		FirstName     string          `json:"first_name"`
		LastName      string          `json:"last_name"`
		EmailAddress  string          `json:"email_address"`
		WorkPhone     json.RawMessage `json:"work_phone"`
		WorkPhoneCode string          `json:"work_phone_code"`
		Mobile        json.RawMessage `json:"mobile"`
		MobileCode    string          `json:"mobile_code"`

		Title     string          `json:"title"`
		Email     string          `json:"email"`
		Phone     json.RawMessage `json:"phone"`
		PhoneCode string          `json:"phone_code"`
	}

	var raw contactPersonAlias

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	workPhone, err := parseStringOrNumber(raw.WorkPhone)
	if err != nil {
		return fmt.Errorf("invalid work_phone: %w", err)
	}

	mobile, err := parseStringOrNumber(raw.Mobile)
	if err != nil {
		return fmt.Errorf("invalid mobile: %w", err)
	}

	phone, err := parseStringOrNumber(raw.Phone)
	if err != nil {
		return fmt.Errorf("invalid phone: %w", err)
	}

	c.Salutation = raw.Salutation
	c.FirstName = raw.FirstName
	c.LastName = raw.LastName
	c.EmailAddress = raw.EmailAddress
	c.WorkPhone = workPhone
	c.WorkPhoneCode = raw.WorkPhoneCode
	c.Mobile = mobile
	c.MobileCode = raw.MobileCode

	c.Title = raw.Title
	c.Email = raw.Email
	c.Phone = phone
	c.PhoneCode = raw.PhoneCode

	c.Normalize()

	return nil
}

func (c *ContactPersonInput) Normalize() {
	if c.Salutation == "" && c.Title != "" {
		c.Salutation = c.Title
	}

	if c.EmailAddress == "" && c.Email != "" {
		c.EmailAddress = c.Email
	}

	if c.Mobile == "" && c.Phone != "" {
		c.Mobile = c.Phone
	}

	if (c.MobileCode == "" || c.MobileCode == "+91") &&
		c.PhoneCode != "" &&
		c.PhoneCode != "+91" {
		c.MobileCode = c.PhoneCode
	}
}
