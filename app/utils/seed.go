package utils

import (
	"errors"
	"log"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

func SeedInitialData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	businessTypes := []models.BusinessType{
		{
			TypeName:    "Retail",
			Description: "Business selling goods directly to consumers",
			IsActive:    true,
		},
		{
			TypeName:    "Wholesale",
			Description: "Business selling goods in bulk to retailers or businesses",
			IsActive:    true,
		},
		{
			TypeName:    "Manufacturing",
			Description: "Business involved in producing goods from raw materials",
			IsActive:    true,
		},
		{
			TypeName:    "Service",
			Description: "Business providing services instead of physical products",
			IsActive:    true,
		},
		{
			TypeName:    "Distribution",
			Description: "Business involved in logistics and distribution of goods",
			IsActive:    true,
		},
		{
			TypeName:    "Others",
			Description: "Other types of businesses not categorized above",
			IsActive:    true,
		},
	}

	for _, bt := range businessTypes {
		var existing models.BusinessType
		if err := db.Where("type_name = ?", bt.TypeName).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&bt).Error; err != nil {
					log.Printf("Failed to create business type %s: %v", bt.TypeName, err)
				}
			}
		}
	}

	taxTypes := []models.TaxType{
		{TaxName: "GST", TaxCode: "GST", Description: "Goods and Services Tax", IsActive: true},
		{TaxName: "IGST", TaxCode: "IGST", Description: "Integrated Goods and Services Tax", IsActive: true},
		{TaxName: "VAT", TaxCode: "VAT", Description: "Value Added Tax", IsActive: true},
		{TaxName: "None", TaxCode: "NONE", Description: "No Tax Applied", IsActive: true},
	}

	for _, tt := range taxTypes {
		var existing models.TaxType
		if err := db.Where("tax_code = ?", tt.TaxCode).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&tt).Error; err != nil {
					log.Printf("Failed to create tax type %s: %v", tt.TaxName, err)
				}
			}
		}
	}

	countries := []models.Country{
		{CountryName: "India", CountryCode: "IN", PhoneCode: "+91", IsActive: true},
		{CountryName: "United States", CountryCode: "US", PhoneCode: "+1", IsActive: true},
		{CountryName: "United Kingdom", CountryCode: "GB", PhoneCode: "+44", IsActive: true},
	}

	for _, c := range countries {
		var existing models.Country
		if err := db.Where("country_code = ?", c.CountryCode).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&c).Error; err != nil {
					log.Printf("Failed to create country %s: %v", c.CountryName, err)
				}
			}
		}
	}

	var indiaCountry models.Country
	if err := db.Where("country_code = ?", "IN").First(&indiaCountry).Error; err == nil {
		states := []models.State{
			{CountryID: indiaCountry.ID, StateName: "Tamil Nadu", StateCode: "TN", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Karnataka", StateCode: "KA", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Maharashtra", StateCode: "MH", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Delhi", StateCode: "DL", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Gujarat", StateCode: "GJ", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Rajasthan", StateCode: "RJ", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Uttar Pradesh", StateCode: "UP", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "West Bengal", StateCode: "WB", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Kerala", StateCode: "KL", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Andhra Pradesh", StateCode: "AP", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Telangana", StateCode: "TS", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Punjab", StateCode: "PB", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Haryana", StateCode: "HR", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Madhya Pradesh", StateCode: "MP", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Bihar", StateCode: "BR", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Odisha", StateCode: "OR", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Assam", StateCode: "AS", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Jharkhand", StateCode: "JH", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Chhattisgarh", StateCode: "CG", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Uttarakhand", StateCode: "UK", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Himachal Pradesh", StateCode: "HP", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Goa", StateCode: "GA", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Jammu and Kashmir", StateCode: "JK", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Ladakh", StateCode: "LA", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Puducherry", StateCode: "PY", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Chandigarh", StateCode: "CH", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Dadra and Nagar Haveli and Daman and Diu", StateCode: "DD", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Lakshadweep", StateCode: "LD", IsActive: true},
			{CountryID: indiaCountry.ID, StateName: "Andaman and Nicobar Islands", StateCode: "AN", IsActive: true},
		}

		for _, s := range states {
			var existing models.State
			if err := db.Where("country_id = ? AND state_name = ?", s.CountryID, s.StateName).First(&existing).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					if err := db.Create(&s).Error; err != nil {
						log.Printf("Failed to create state %s: %v", s.StateName, err)
					}
				}
			}
		}
	}

	log.Println("Initial data seeding completed!")
	return nil
}

func SeedDefaultCompany(db *gorm.DB) error {
	log.Println("Seeding default company...")

	return db.Transaction(func(tx *gorm.DB) error {

		var existing models.Company
		if err := tx.Where("company_name = ?", "BB Cloud Technologies").
			First(&existing).Error; err == nil {
			log.Println("Company already exists, skipping seed")
			return nil
		}

		var businessType models.BusinessType
		if err := tx.Where("type_name = ?", "Private Limited Company").
			First(&businessType).Error; err != nil {
			return err
		}

		company := models.Company{
			CompanyName:    "BB Cloud Technologies",
			BusinessTypeID: businessType.ID,
			GSTNumber:      "33ABCDE1234F1Z5",
			PANNumber:      "ABCDE1234F",
		}

		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		contact := models.CompanyContact{
			CompanyID: company.ID,
			Mobile:    "9876543210",
			Email:     "info@bbcloud.app",
		}
		if err := tx.Create(&contact).Error; err != nil {
			return err
		}

		var country models.Country
		if err := tx.Where("country_code = ?", "IN").First(&country).Error; err != nil {
			return err
		}

		var state models.State
		if err := tx.Where("state_code = ?", "TN").First(&state).Error; err != nil {
			return err
		}

		address := models.CompanyAddress{
			CompanyID:    company.ID,
			AddressLine1: "No 10, Anna Nagar",
			City:         "Chennai",
			StateID:      state.ID,
			CountryID:    country.ID,
			Pincode:      "600040",
		}
		if err := tx.Create(&address).Error; err != nil {
			return err
		}

		invoiceSettings := models.CompanyInvoiceSetting{
			CompanyID:          company.ID,
			InvoicePrefix:      "INV",
			InvoiceStartNumber: 1,
		}
		if err := tx.Create(&invoiceSettings).Error; err != nil {
			return err
		}

		var taxType models.TaxType
		if err := tx.Where("tax_code = ?", "GST").First(&taxType).Error; err != nil {
			return err
		}

		taxSettings := models.CompanyTaxSetting{
			CompanyID:  company.ID,
			GSTEnabled: true,
			TaxTypeID:  taxType.ID,
		}
		if err := tx.Create(&taxSettings).Error; err != nil {
			return err
		}

		regionalSettings := models.CompanyRegionalSetting{
			CompanyID: company.ID,
		}
		if err := tx.Create(&regionalSettings).Error; err != nil {
			return err
		}

		bank := models.CompanyBankDetail{
			CompanyID:         company.ID,
			BankID:            1, // Will need to fetch actual bank ID from database
			AccountHolderName: "BB Cloud Technologies",
			AccountNumber:     "123456789012",
			IsPrimary:         true,
		}
		if err := tx.Create(&bank).Error; err != nil {
			return err
		}

		upi := models.CompanyUPIDetail{
			CompanyID: company.ID,
			UPIID:     "bbcloud@upi",
		}
		if err := tx.Create(&upi).Error; err != nil {
			return err
		}

		log.Println("Default company seeded successfully!")
		return nil
	})
}
