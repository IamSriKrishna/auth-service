package helper

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
)

func MapCreateVendorInput(input *input.CreateVendorInput) *models.Vendor {
	vendor := &models.Vendor{
		Salutation:     input.Salutation,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		CompanyName:    input.CompanyName,
		DisplayName:    input.DisplayName,
		EmailAddress:   input.EmailAddress,
		WorkPhone:      input.WorkPhone,
		WorkPhoneCode:  input.WorkPhoneCode,
		Mobile:         input.Mobile,
		MobileCode:     input.MobileCode,
		VendorLanguage: input.VendorLanguage,
	}

	if input.OtherDetails != nil {
		vendor.OtherDetails = &models.EntityOtherDetails{
			PAN:              input.OtherDetails.PAN,
			IsMSMERegistered: input.OtherDetails.IsMSMERegistered,
			Currency:         input.OtherDetails.Currency,
			PaymentTerms:     input.OtherDetails.PaymentTerms,
			TDS:              input.OtherDetails.TDS,
			EnablePortal:     input.OtherDetails.EnablePortal,
			WebsiteURL:       input.OtherDetails.WebsiteURL,
			Department:       input.OtherDetails.Department,
			Designation:      input.OtherDetails.Designation,
			Twitter:          input.OtherDetails.Twitter,
			SkypeName:        input.OtherDetails.SkypeName,
			Facebook:         input.OtherDetails.Facebook,
		}
	}

	if input.BillingAddress != nil {
		billingAddr := input.BillingAddress
		billingAddr.Normalize()

		vendor.BillingAddress = &models.EntityAddress{
			AddressType:   "billing",
			Attention:     billingAddr.Attention,
			CountryRegion: billingAddr.CountryRegion,
			AddressLine1:  billingAddr.AddressLine1,
			AddressLine2:  billingAddr.AddressLine2,
			City:          billingAddr.City,
			State:         billingAddr.State,
			PinCode:       billingAddr.PinCode,
			Phone:         billingAddr.Phone,
			PhoneCode:     billingAddr.PhoneCode,
			FaxNumber:     billingAddr.FaxNumber,
		}
	}

	if input.ShippingAddress != nil {
		shippingAddr := input.ShippingAddress
		shippingAddr.Normalize()

		vendor.ShippingAddress = &models.EntityAddress{
			AddressType:   "shipping",
			Attention:     shippingAddr.Attention,
			CountryRegion: shippingAddr.CountryRegion,
			AddressLine1:  shippingAddr.AddressLine1,
			AddressLine2:  shippingAddr.AddressLine2,
			City:          shippingAddr.City,
			State:         shippingAddr.State,
			PinCode:       shippingAddr.PinCode,
			Phone:         shippingAddr.Phone,
			PhoneCode:     shippingAddr.PhoneCode,
			FaxNumber:     shippingAddr.FaxNumber,
		}
	}

	for _, cp := range input.ContactPersons {
		// Normalize alternative field names to primary fields
		cpCopy := cp
		cpCopy.Normalize()

		vendor.ContactPersons = append(vendor.ContactPersons, models.EntityContactPerson{
			Salutation:    cpCopy.Salutation,
			FirstName:     cpCopy.FirstName,
			LastName:      cpCopy.LastName,
			EmailAddress:  cpCopy.EmailAddress,
			WorkPhone:     cpCopy.WorkPhone,
			WorkPhoneCode: cpCopy.WorkPhoneCode,
			Mobile:        cpCopy.Mobile,
			MobileCode:    cpCopy.MobileCode,
		})
	}

	for _, bd := range input.BankDetails {
		vendor.BankDetails = append(vendor.BankDetails, models.VendorBankDetail{
			BankID:            bd.BankID,
			AccountHolderName: bd.AccountHolderName,
			AccountNumber:     bd.AccountNumber,
		})
	}

	return vendor
}

func ApplyUpdateVendorInput(vendor *models.Vendor, input *input.UpdateVendorInput) {
	if input.Salutation != nil {
		vendor.Salutation = *input.Salutation
	}
	if input.FirstName != nil {
		vendor.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		vendor.LastName = *input.LastName
	}
	if input.CompanyName != nil {
		vendor.CompanyName = *input.CompanyName
	}
	if input.DisplayName != nil {
		vendor.DisplayName = *input.DisplayName
	}
	if input.EmailAddress != nil {
		vendor.EmailAddress = *input.EmailAddress
	}
	if input.WorkPhone != nil {
		vendor.WorkPhone = *input.WorkPhone
	}
	if input.WorkPhoneCode != nil {
		vendor.WorkPhoneCode = *input.WorkPhoneCode
	}
	if input.Mobile != nil {
		vendor.Mobile = *input.Mobile
	}
	if input.MobileCode != nil {
		vendor.MobileCode = *input.MobileCode
	}
	if input.VendorLanguage != nil {
		vendor.VendorLanguage = *input.VendorLanguage
	}

	if input.OtherDetails != nil {
		if vendor.OtherDetails == nil {
			vendor.OtherDetails = &models.EntityOtherDetails{}
		}
		vendor.OtherDetails.PAN = input.OtherDetails.PAN
		vendor.OtherDetails.IsMSMERegistered = input.OtherDetails.IsMSMERegistered
		vendor.OtherDetails.Currency = input.OtherDetails.Currency
		vendor.OtherDetails.PaymentTerms = input.OtherDetails.PaymentTerms
		vendor.OtherDetails.TDS = input.OtherDetails.TDS
		vendor.OtherDetails.EnablePortal = input.OtherDetails.EnablePortal
		vendor.OtherDetails.WebsiteURL = input.OtherDetails.WebsiteURL
		vendor.OtherDetails.Department = input.OtherDetails.Department
		vendor.OtherDetails.Designation = input.OtherDetails.Designation
		vendor.OtherDetails.Twitter = input.OtherDetails.Twitter
		vendor.OtherDetails.SkypeName = input.OtherDetails.SkypeName
		vendor.OtherDetails.Facebook = input.OtherDetails.Facebook
	}

	if input.BillingAddress != nil {
		if vendor.BillingAddress == nil {
			vendor.BillingAddress = &models.EntityAddress{AddressType: "billing"}
		}
		billingAddr := input.BillingAddress
		billingAddr.Normalize()

		vendor.BillingAddress.Attention = billingAddr.Attention
		vendor.BillingAddress.CountryRegion = billingAddr.CountryRegion
		vendor.BillingAddress.AddressLine1 = billingAddr.AddressLine1
		vendor.BillingAddress.AddressLine2 = billingAddr.AddressLine2
		vendor.BillingAddress.City = billingAddr.City
		vendor.BillingAddress.State = billingAddr.State
		vendor.BillingAddress.PinCode = billingAddr.PinCode
		vendor.BillingAddress.Phone = billingAddr.Phone
		vendor.BillingAddress.PhoneCode = billingAddr.PhoneCode
		vendor.BillingAddress.FaxNumber = billingAddr.FaxNumber
	}

	if input.ShippingAddress != nil {
		if vendor.ShippingAddress == nil {
			vendor.ShippingAddress = &models.EntityAddress{AddressType: "shipping"}
		}
		shippingAddr := input.ShippingAddress
		shippingAddr.Normalize()

		vendor.ShippingAddress.Attention = shippingAddr.Attention
		vendor.ShippingAddress.CountryRegion = shippingAddr.CountryRegion
		vendor.ShippingAddress.AddressLine1 = shippingAddr.AddressLine1
		vendor.ShippingAddress.AddressLine2 = shippingAddr.AddressLine2
		vendor.ShippingAddress.City = shippingAddr.City
		vendor.ShippingAddress.State = shippingAddr.State
		vendor.ShippingAddress.PinCode = shippingAddr.PinCode
		vendor.ShippingAddress.Phone = shippingAddr.Phone
		vendor.ShippingAddress.PhoneCode = shippingAddr.PhoneCode
		vendor.ShippingAddress.FaxNumber = shippingAddr.FaxNumber
	}
}

func MapVendorToOutput(vendor *models.Vendor) *output.VendorOutput {
	out := &output.VendorOutput{
		ID:             vendor.ID,
		Salutation:     vendor.Salutation,
		FirstName:      vendor.FirstName,
		LastName:       vendor.LastName,
		CompanyName:    vendor.CompanyName,
		DisplayName:    vendor.DisplayName,
		EmailAddress:   vendor.EmailAddress,
		WorkPhone:      vendor.WorkPhone,
		WorkPhoneCode:  vendor.WorkPhoneCode,
		Mobile:         vendor.Mobile,
		MobileCode:     vendor.MobileCode,
		VendorLanguage: vendor.VendorLanguage,
		CreatedAt:      vendor.CreatedAt,
		UpdatedAt:      vendor.UpdatedAt,
	}

	if vendor.OtherDetails != nil {
		out.OtherDetails = &output.OtherDetailsOutput{
			ID:               vendor.OtherDetails.ID,
			VendorID:         vendor.OtherDetails.EntityID,
			PAN:              vendor.OtherDetails.PAN,
			IsMSMERegistered: vendor.OtherDetails.IsMSMERegistered,
			Currency:         vendor.OtherDetails.Currency,
			PaymentTerms:     vendor.OtherDetails.PaymentTerms,
			TDS:              vendor.OtherDetails.TDS,
			EnablePortal:     vendor.OtherDetails.EnablePortal,
			WebsiteURL:       vendor.OtherDetails.WebsiteURL,
			Department:       vendor.OtherDetails.Department,
			Designation:      vendor.OtherDetails.Designation,
			Twitter:          vendor.OtherDetails.Twitter,
			SkypeName:        vendor.OtherDetails.SkypeName,
			Facebook:         vendor.OtherDetails.Facebook,
			CreatedAt:        vendor.OtherDetails.CreatedAt,
			UpdatedAt:        vendor.OtherDetails.UpdatedAt,
		}
	}

	if vendor.BillingAddress != nil {
		out.BillingAddress = mapVendorAddressOutput(vendor.BillingAddress)
	}

	if vendor.ShippingAddress != nil {
		out.ShippingAddress = mapVendorAddressOutput(vendor.ShippingAddress)
	}

	for _, cp := range vendor.ContactPersons {
		out.ContactPersons = append(out.ContactPersons, output.ContactPersonOutput{
			ID:            cp.ID,
			VendorID:      cp.EntityID,
			Salutation:    cp.Salutation,
			FirstName:     cp.FirstName,
			LastName:      cp.LastName,
			EmailAddress:  cp.EmailAddress,
			WorkPhone:     cp.WorkPhone,
			WorkPhoneCode: cp.WorkPhoneCode,
			Mobile:        cp.Mobile,
			MobileCode:    cp.MobileCode,
			CreatedAt:     cp.CreatedAt,
			UpdatedAt:     cp.UpdatedAt,
		})
	}

	for _, bd := range vendor.BankDetails {
		out.BankDetails = append(out.BankDetails, output.BankDetailOutput{
			ID:                bd.ID,
			VendorID:          bd.VendorID,
			BankID:            bd.BankID,
			AccountHolderName: bd.AccountHolderName,
			AccountNumber:     bd.AccountNumber,
			CreatedAt:         bd.CreatedAt,
			UpdatedAt:         bd.UpdatedAt,
		})
	}

	for _, doc := range vendor.Documents {
		out.Documents = append(out.Documents, output.DocumentOutput{
			ID:        doc.ID,
			VendorID:  doc.EntityID,
			FileName:  doc.FileName,
			FilePath:  doc.FilePath,
			FileSize:  doc.FileSize,
			MimeType:  doc.MimeType,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}

	return out
}

func mapVendorAddressOutput(addr *models.EntityAddress) *output.AddressOutput {
	return &output.AddressOutput{
		ID:            addr.ID,
		VendorID:      addr.EntityID,
		AddressType:   addr.AddressType,
		Attention:     addr.Attention,
		CountryRegion: addr.CountryRegion,
		AddressLine1:  addr.AddressLine1,
		AddressLine2:  addr.AddressLine2,
		City:          addr.City,
		State:         addr.State,
		PinCode:       addr.PinCode,
		Phone:         addr.Phone,
		PhoneCode:     addr.PhoneCode,
		FaxNumber:     addr.FaxNumber,
		CreatedAt:     addr.CreatedAt,
		UpdatedAt:     addr.UpdatedAt,
	}
}
