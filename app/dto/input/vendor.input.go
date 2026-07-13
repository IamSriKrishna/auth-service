package input

import (
	"encoding/json"
	"fmt"
)

type CreateVendorInput struct {
	Salutation      string               `json:"salutation" validate:"required"`
	FirstName       string               `json:"first_name" validate:"required"`
	LastName        string               `json:"last_name"`
	DisplayName     string               `json:"display_name" validate:"required"`
	EmailAddress    string               `json:"email_address" validate:"omitempty,email"`
	WorkPhone       string               `json:"work_phone"`
	WorkPhoneCode   string               `json:"work_phone_code"`
	Mobile          string               `json:"mobile"`
	MobileCode      string               `json:"mobile_code"`
	VendorLanguage  string               `json:"vendor_language"`
	GSTIN           string               `json:"gstin"`
	OtherDetails    *OtherDetailsInput   `json:"other_details"`
	BillingAddress  *AddressInput        `json:"billing_address"`
	ShippingAddress *AddressInput        `json:"shipping_address"`
	ContactPersons  []ContactPersonInput `json:"contact_persons"`
	BankDetails     []BankDetailInput    `json:"bank_details"`
}

func (v *CreateVendorInput) UnmarshalJSON(data []byte) error {
	type vendorJSON struct {
		Salutation      string               `json:"salutation"`
		FirstName       string               `json:"first_name"`
		LastName        string               `json:"last_name"`
		DisplayName     string               `json:"display_name"`
		EmailAddress    string               `json:"email_address"`
		WorkPhone       json.RawMessage      `json:"work_phone"`
		WorkPhoneCode   string               `json:"work_phone_code"`
		Mobile          json.RawMessage      `json:"mobile"`
		MobileCode      string               `json:"mobile_code"`
		VendorLanguage  string               `json:"vendor_language"`
		GSTIN           string               `json:"gstin"`
		OtherDetails    *OtherDetailsInput   `json:"other_details"`
		BillingAddress  *AddressInput        `json:"billing_address"`
		ShippingAddress *AddressInput        `json:"shipping_address"`
		ContactPersons  []ContactPersonInput `json:"contact_persons"`
		BankDetails     []BankDetailInput    `json:"bank_details"`
	}

	var raw vendorJSON

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

	v.Salutation = raw.Salutation
	v.FirstName = raw.FirstName
	v.LastName = raw.LastName
	v.DisplayName = raw.DisplayName
	v.EmailAddress = raw.EmailAddress
	v.WorkPhone = workPhone
	v.WorkPhoneCode = raw.WorkPhoneCode
	v.Mobile = mobile
	v.MobileCode = raw.MobileCode
	v.VendorLanguage = raw.VendorLanguage
	v.GSTIN = raw.GSTIN
	v.OtherDetails = raw.OtherDetails
	v.BillingAddress = raw.BillingAddress
	v.ShippingAddress = raw.ShippingAddress
	v.ContactPersons = raw.ContactPersons
	v.BankDetails = raw.BankDetails

	return nil
}

type UpdateVendorInput struct {
	Salutation      *string              `json:"salutation"`
	FirstName       *string              `json:"first_name"`
	LastName        *string              `json:"last_name"`
	DisplayName     *string              `json:"display_name"`
	EmailAddress    *string              `json:"email_address" validate:"omitempty,email"`
	WorkPhone       *string              `json:"work_phone"`
	WorkPhoneCode   *string              `json:"work_phone_code"`
	Mobile          *string              `json:"mobile"`
	MobileCode      *string              `json:"mobile_code"`
	VendorLanguage  *string              `json:"vendor_language"`
	GSTIN           *string              `json:"gstin"`
	OtherDetails    *OtherDetailsInput   `json:"other_details"`
	BillingAddress  *AddressInput        `json:"billing_address"`
	ShippingAddress *AddressInput        `json:"shipping_address"`
	ContactPersons  []ContactPersonInput `json:"contact_persons"`
	BankDetails     []BankDetailInput    `json:"bank_details"`
}

func (v *UpdateVendorInput) UnmarshalJSON(data []byte) error {
	type updateVendorJSON struct {
		Salutation      *string              `json:"salutation"`
		FirstName       *string              `json:"first_name"`
		LastName        *string              `json:"last_name"`
		DisplayName     *string              `json:"display_name"`
		EmailAddress    *string              `json:"email_address"`
		WorkPhone       json.RawMessage      `json:"work_phone"`
		WorkPhoneCode   *string              `json:"work_phone_code"`
		Mobile          json.RawMessage      `json:"mobile"`
		MobileCode      *string              `json:"mobile_code"`
		VendorLanguage  *string              `json:"vendor_language"`
		GSTIN           *string              `json:"gstin"`
		OtherDetails    *OtherDetailsInput   `json:"other_details"`
		BillingAddress  *AddressInput        `json:"billing_address"`
		ShippingAddress *AddressInput        `json:"shipping_address"`
		ContactPersons  []ContactPersonInput `json:"contact_persons"`
		BankDetails     []BankDetailInput    `json:"bank_details"`
	}

	var raw updateVendorJSON

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Salutation = raw.Salutation
	v.FirstName = raw.FirstName
	v.LastName = raw.LastName
	v.DisplayName = raw.DisplayName
	v.EmailAddress = raw.EmailAddress
	v.WorkPhoneCode = raw.WorkPhoneCode
	v.MobileCode = raw.MobileCode
	v.VendorLanguage = raw.VendorLanguage
	v.GSTIN = raw.GSTIN
	v.OtherDetails = raw.OtherDetails
	v.BillingAddress = raw.BillingAddress
	v.ShippingAddress = raw.ShippingAddress
	v.ContactPersons = raw.ContactPersons
	v.BankDetails = raw.BankDetails

	if len(raw.WorkPhone) > 0 && string(raw.WorkPhone) != "null" {
		workPhone, err := parseStringOrNumber(raw.WorkPhone)
		if err != nil {
			return fmt.Errorf("invalid work_phone: %w", err)
		}

		v.WorkPhone = &workPhone
	}

	if len(raw.Mobile) > 0 && string(raw.Mobile) != "null" {
		mobile, err := parseStringOrNumber(raw.Mobile)
		if err != nil {
			return fmt.Errorf("invalid mobile: %w", err)
		}

		v.Mobile = &mobile
	}

	return nil
}
