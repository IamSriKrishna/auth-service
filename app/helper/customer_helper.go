package helper

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
)

func MapCreateCustomerInput(input *input.CreateCustomerInput) *models.Customer {
	customer := &models.Customer{
		Salutation:       input.Salutation,
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		CompanyName:      input.CompanyName,
		DisplayName:      input.DisplayName,
		EmailAddress:     input.EmailAddress,
		WorkPhone:        input.WorkPhone,
		WorkPhoneCode:    input.WorkPhoneCode,
		Mobile:           input.Mobile,
		MobileCode:       input.MobileCode,
		CustomerLanguage: input.CustomerLanguage,
	}

	customer.OtherDetails = mapOtherDetails(input.OtherDetails)

	if billingAddr := mapAddress(input.BillingAddress, "billing"); billingAddr != nil {
		customer.Addresses = append(customer.Addresses, *billingAddr)
	}
	if shippingAddr := mapAddress(input.ShippingAddress, "shipping"); shippingAddr != nil {
		customer.Addresses = append(customer.Addresses, *shippingAddr)
	}

	for _, cp := range input.ContactPersons {
		// Normalize alternative field names to primary fields
		cpCopy := cp
		cpCopy.Normalize()

		customer.ContactPersons = append(customer.ContactPersons, models.EntityContactPerson{
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

	return customer
}

func ApplyUpdateCustomerInput(customer *models.Customer, input *input.UpdateCustomerInput) {
	if input.Salutation != nil {
		customer.Salutation = *input.Salutation
	}
	if input.FirstName != nil {
		customer.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		customer.LastName = *input.LastName
	}
	if input.CompanyName != nil {
		customer.CompanyName = *input.CompanyName
	}
	if input.DisplayName != nil {
		customer.DisplayName = *input.DisplayName
	}
	if input.EmailAddress != nil {
		customer.EmailAddress = *input.EmailAddress
	}
	if input.WorkPhone != nil {
		customer.WorkPhone = *input.WorkPhone
	}
	if input.WorkPhoneCode != nil {
		customer.WorkPhoneCode = *input.WorkPhoneCode
	}
	if input.Mobile != nil {
		customer.Mobile = *input.Mobile
	}
	if input.MobileCode != nil {
		customer.MobileCode = *input.MobileCode
	}
	if input.CustomerLanguage != nil {
		customer.CustomerLanguage = *input.CustomerLanguage
	}

	if input.OtherDetails != nil {
		if customer.OtherDetails == nil {
			customer.OtherDetails = &models.EntityOtherDetails{}
		}
		customer.OtherDetails.PAN = input.OtherDetails.PAN
		customer.OtherDetails.Currency = input.OtherDetails.Currency
		customer.OtherDetails.PaymentTerms = input.OtherDetails.PaymentTerms
		customer.OtherDetails.EnablePortal = input.OtherDetails.EnablePortal
	}

	if input.BillingAddress != nil {
		updateAddressInSlice(&customer.Addresses, input.BillingAddress, "billing")
	}
	if input.ShippingAddress != nil {
		updateAddressInSlice(&customer.Addresses, input.ShippingAddress, "shipping")
	}

	if len(input.ContactPersons) > 0 {
		customer.ContactPersons = []models.EntityContactPerson{}

		for _, cp := range input.ContactPersons {
			// Normalize alternative field names to primary fields
			cpCopy := cp
			cpCopy.Normalize()

			customer.ContactPersons = append(customer.ContactPersons, models.EntityContactPerson{
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
	}
}
func mapAddress(input *input.AddressInput, addressType string) *models.EntityAddress {
	if input == nil {
		return nil
	}

	// Normalize alternative field names to primary fields
	input.Normalize()

	return &models.EntityAddress{
		AddressType:   addressType,
		Attention:     input.Attention,
		CountryRegion: input.CountryRegion,
		AddressLine1:  input.AddressLine1,
		AddressLine2:  input.AddressLine2,
		City:          input.City,
		State:         input.State,
		PinCode:       input.PinCode,
		Phone:         input.Phone,
		PhoneCode:     input.PhoneCode,
		FaxNumber:     input.FaxNumber,
	}
}

func updateAddressInSlice(addresses *[]models.EntityAddress, input *input.AddressInput, addressType string) {
	if input == nil {
		return
	}

	// Normalize alternative field names to primary fields
	input.Normalize()

	for i := range *addresses {
		if (*addresses)[i].AddressType == addressType {
			(*addresses)[i].Attention = input.Attention
			(*addresses)[i].CountryRegion = input.CountryRegion
			(*addresses)[i].AddressLine1 = input.AddressLine1
			(*addresses)[i].AddressLine2 = input.AddressLine2
			(*addresses)[i].City = input.City
			(*addresses)[i].State = input.State
			(*addresses)[i].PinCode = input.PinCode
			(*addresses)[i].Phone = input.Phone
			(*addresses)[i].PhoneCode = input.PhoneCode
			(*addresses)[i].FaxNumber = input.FaxNumber
			return
		}
	}
	newAddr := mapAddress(input, addressType)
	if newAddr != nil {
		*addresses = append(*addresses, *newAddr)
	}
}

func mapOtherDetails(input *input.OtherDetailsCustomerInput) *models.EntityOtherDetails {
	if input == nil {
		return nil
	}

	return &models.EntityOtherDetails{
		PAN:          input.PAN,
		Currency:     input.Currency,
		PaymentTerms: input.PaymentTerms,
		EnablePortal: input.EnablePortal,
	}
}

func MapCustomerToOutput(customer *models.Customer) *output.CustomerOutput {
	out := &output.CustomerOutput{
		ID:               customer.ID,
		Salutation:       customer.Salutation,
		FirstName:        customer.FirstName,
		LastName:         customer.LastName,
		CompanyName:      customer.CompanyName,
		DisplayName:      customer.DisplayName,
		EmailAddress:     customer.EmailAddress,
		WorkPhone:        customer.WorkPhone,
		WorkPhoneCode:    customer.WorkPhoneCode,
		Mobile:           customer.Mobile,
		MobileCode:       customer.MobileCode,
		CustomerLanguage: customer.CustomerLanguage,
		CreatedAt:        customer.CreatedAt,
		UpdatedAt:        customer.UpdatedAt,
	}

	if customer.OtherDetails != nil {
		out.OtherDetails = &output.CustomerOtherDetailsOutput{
			ID:           customer.OtherDetails.ID,
			CustomerID:   customer.OtherDetails.EntityID,
			PAN:          customer.OtherDetails.PAN,
			Currency:     customer.OtherDetails.Currency,
			PaymentTerms: customer.OtherDetails.PaymentTerms,
			EnablePortal: customer.OtherDetails.EnablePortal,
			CreatedAt:    customer.OtherDetails.CreatedAt,
			UpdatedAt:    customer.OtherDetails.UpdatedAt,
		}
	}

	billingAddr := customer.GetBillingAddress()
	if billingAddr != nil {
		out.BillingAddress = mapAddressOutput(billingAddr)
	}

	shippingAddr := customer.GetShippingAddress()
	if shippingAddr != nil {
		out.ShippingAddress = mapAddressOutput(shippingAddr)
	}

	for _, cp := range customer.ContactPersons {
		out.ContactPersons = append(out.ContactPersons, output.CustomerContactPersonOutput{
			ID:            cp.ID,
			CustomerID:    cp.EntityID,
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

	return out
}

func mapAddressOutput(addr *models.EntityAddress) *output.CustomerAddressOutput {
	return &output.CustomerAddressOutput{
		ID:            addr.ID,
		CustomerID:    addr.EntityID,
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
