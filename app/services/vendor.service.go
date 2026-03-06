package services

import (
	"errors"
	"fmt"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/helper"
	"github.com/bbapp-org/auth-service/app/repo"
	"gorm.io/gorm"
)

type VendorService interface {
	// Step 1: Set Up Your Contacts - Vendors
	// Define your suppliers, their tax details, and currency information
	CreateVendor(input *input.CreateVendorInput) (*output.VendorOutput, error)
	UpdateVendor(id uint, input *input.UpdateVendorInput) (*output.VendorOutput, error)
	GetVendorByID(id uint) (*output.VendorOutput, error)
	GetAllVendors(page, limit int) ([]output.VendorListOutput, int64, error)
	DeleteVendor(id uint) error

	// User-specific vendor operations
	CreateVendorForUser(userID, companyID uint, input *input.CreateVendorInput) (*output.VendorOutput, error)
	UpdateVendorForUser(id, userID uint, input *input.UpdateVendorInput) (*output.VendorOutput, error)
	GetVendorByIDAndUser(id, userID uint) (*output.VendorOutput, error)
	GetVendorsByUser(userID, companyID uint, page, limit int) ([]output.VendorListOutput, int64, error)
	DeleteVendorForUser(id, userID uint) error
}

type vendorService struct {
	repo repo.VendorRepository
}

func NewVendorService(repo repo.VendorRepository) VendorService {
	return &vendorService{repo: repo}
}

func (s *vendorService) CreateVendor(input *input.CreateVendorInput) (*output.VendorOutput, error) {
	if input.Mobile != "" {
		existingVendor, err := s.repo.FindByMobile(input.Mobile)
		if err == nil && existingVendor != nil {
			return nil, fmt.Errorf("mobile number %s already exists with vendor: %s",
				input.Mobile, existingVendor.DisplayName)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	vendor := helper.MapCreateVendorInput(input)

	if err := s.repo.Create(vendor); err != nil {
		return nil, err
	}

	return s.GetVendorByID(vendor.ID)
}

func (s *vendorService) UpdateVendor(id uint, input *input.UpdateVendorInput) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	if input.Mobile != nil && *input.Mobile != "" && *input.Mobile != vendor.Mobile {
		existingVendor, err := s.repo.FindByMobile(*input.Mobile)
		if err == nil && existingVendor != nil && existingVendor.ID != id {
			return nil, fmt.Errorf("mobile number %s already exists with vendor: %s",
				*input.Mobile, existingVendor.DisplayName)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	helper.ApplyUpdateVendorInput(vendor, input)

	if err := s.repo.Update(vendor); err != nil {
		return nil, err
	}

	return s.GetVendorByID(vendor.ID)
}

func (s *vendorService) GetVendorByID(id uint) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return helper.MapVendorToOutput(vendor), nil
}

func (s *vendorService) GetAllVendors(page, limit int) ([]output.VendorListOutput, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	vendors, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var outputs []output.VendorListOutput
	for _, vendor := range vendors {
		outputs = append(outputs, output.VendorListOutput{
			ID:             vendor.ID,
			DisplayName:    vendor.DisplayName,
			EmailAddress:   vendor.EmailAddress,
			WorkPhone:      vendor.WorkPhone,
			Mobile:         vendor.Mobile,
			VendorLanguage: vendor.VendorLanguage,
			CreatedAt:      vendor.CreatedAt,
			UpdatedAt:      vendor.UpdatedAt,
		})
	}

	return outputs, total, nil
}

func (s *vendorService) DeleteVendor(id uint) error {
	return s.repo.Delete(id)
}

func (s *vendorService) CreateVendorForUser(userID, companyID uint, input *input.CreateVendorInput) (*output.VendorOutput, error) {
	if input.Mobile != "" {
		existingVendor, err := s.repo.FindByMobile(input.Mobile)
		if err == nil && existingVendor != nil {
			return nil, fmt.Errorf("mobile number %s already exists with vendor: %s",
				input.Mobile, existingVendor.DisplayName)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	vendor := helper.MapCreateVendorInput(input)
	vendor.UserID = userID
	vendor.CompanyID = companyID

	if err := s.repo.Create(vendor); err != nil {
		return nil, err
	}

	return s.GetVendorByID(vendor.ID)
}

func (s *vendorService) GetVendorByIDAndUser(id, userID uint) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("vendor not found")
	}

	return helper.MapVendorToOutput(vendor), nil
}

func (s *vendorService) GetVendorsByUser(userID, companyID uint, page, limit int) ([]output.VendorListOutput, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	vendors, total, err := s.repo.FindByUserID(userID, companyID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var outputs []output.VendorListOutput
	for _, vendor := range vendors {
		outputs = append(outputs, output.VendorListOutput{
			ID:             vendor.ID,
			DisplayName:    vendor.DisplayName,
			EmailAddress:   vendor.EmailAddress,
			WorkPhone:      vendor.WorkPhone,
			Mobile:         vendor.Mobile,
			VendorLanguage: vendor.VendorLanguage,
			CreatedAt:      vendor.CreatedAt,
			UpdatedAt:      vendor.UpdatedAt,
		})
	}

	return outputs, total, nil
}

func (s *vendorService) UpdateVendorForUser(id, userID uint, input *input.UpdateVendorInput) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("vendor not found")
	}

	if input.FirstName != nil {
		vendor.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		vendor.LastName = *input.LastName
	}
	if input.DisplayName != nil {
		vendor.DisplayName = *input.DisplayName
	}
	if input.EmailAddress != nil {
		vendor.EmailAddress = *input.EmailAddress
	}
	if input.Mobile != nil {
		vendor.Mobile = *input.Mobile
	}
	if input.WorkPhone != nil {
		vendor.WorkPhone = *input.WorkPhone
	}

	if err := s.repo.Update(vendor); err != nil {
		return nil, err
	}

	return s.GetVendorByID(vendor.ID)
}

func (s *vendorService) DeleteVendorForUser(id, userID uint) error {
	vendor, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return fmt.Errorf("vendor not found")
	}

	return s.repo.Delete(vendor.ID)
}
