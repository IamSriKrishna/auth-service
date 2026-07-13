package services

import (
	"errors"
	"fmt"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/helper"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"gorm.io/gorm"
)

type VendorService interface {
	// Existing methods retained for compatibility.
	CreateVendor(input *input.CreateVendorInput) (*output.VendorOutput, error)
	UpdateVendor(id uint, input *input.UpdateVendorInput) (*output.VendorOutput, error)
	GetVendorByID(id uint) (*output.VendorOutput, error)
	GetAllVendors(page, limit int) ([]output.VendorListOutput, int64, error)
	DeleteVendor(id uint) error

	CreateVendorForUser(
		userID uint,
		companyID uint,
		input *input.CreateVendorInput,
	) (*output.VendorOutput, error)

	UpdateVendorForUser(
		id uint,
		userID uint,
		input *input.UpdateVendorInput,
	) (*output.VendorOutput, error)

	GetVendorByIDAndUser(
		id uint,
		userID uint,
	) (*output.VendorOutput, error)

	GetVendorsByUser(
		userID uint,
		companyID uint,
		page int,
		limit int,
	) ([]output.VendorListOutput, int64, error)

	DeleteVendorForUser(
		id uint,
		userID uint,
	) error

	// Company-scoped methods added.
	GetVendorByIDAndCompany(
		id uint,
		companyID uint,
	) (*output.VendorOutput, error)

	GetVendorsByCompany(
		companyID uint,
		page int,
		limit int,
	) ([]output.VendorListOutput, int64, error)

	UpdateVendorForCompany(
		id uint,
		companyID uint,
		input *input.UpdateVendorInput,
	) (*output.VendorOutput, error)

	DeleteVendorForCompany(
		id uint,
		companyID uint,
	) error
}

type vendorService struct {
	repo        repo.VendorRepository
	companyRepo repo.CompanyRepository
}

func NewVendorService(
	repo repo.VendorRepository,
	companyRepo repo.CompanyRepository,
) VendorService {
	return &vendorService{
		repo:        repo,
		companyRepo: companyRepo,
	}
}

func vendorListOutputs(
	vendors []models.Vendor,
) []output.VendorListOutput {
	outputs := make([]output.VendorListOutput, 0, len(vendors))

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

	return outputs
}

func (s *vendorService) CreateVendor(
	input *input.CreateVendorInput,
) (*output.VendorOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if input.Mobile != "" {
		existingVendor, err := s.repo.FindByMobile(input.Mobile)
		if err == nil && existingVendor != nil {
			return nil, fmt.Errorf(
				"mobile number %s already exists with vendor: %s",
				input.Mobile,
				existingVendor.DisplayName,
			)
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

func (s *vendorService) UpdateVendor(
	id uint,
	input *input.UpdateVendorInput,
) (*output.VendorOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	vendor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	if input.Mobile != nil &&
		*input.Mobile != "" &&
		*input.Mobile != vendor.Mobile {
		existingVendor, err := s.repo.FindByMobile(*input.Mobile)
		if err == nil &&
			existingVendor != nil &&
			existingVendor.ID != id {
			return nil, fmt.Errorf(
				"mobile number %s already exists with vendor: %s",
				*input.Mobile,
				existingVendor.DisplayName,
			)
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

func (s *vendorService) GetVendorByID(
	id uint,
) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return helper.MapVendorToOutput(vendor), nil
}

func (s *vendorService) GetAllVendors(
	page int,
	limit int,
) ([]output.VendorListOutput, int64, error) {
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

	return vendorListOutputs(vendors), total, nil
}

func (s *vendorService) DeleteVendor(
	id uint,
) error {
	return s.repo.Delete(id)
}

func (s *vendorService) CreateVendorForUser(
	userID uint,
	companyID uint,
	input *input.CreateVendorInput,
) (*output.VendorOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if userID == 0 {
		return nil, errors.New("invalid user")
	}

	if companyID == 0 {
		return nil, errors.New("invalid company")
	}

	company, err := s.companyRepo.FindByID(companyID)
	if err != nil || company == nil {
		return nil, errors.New("company not found")
	}

	if input.Mobile != "" {
		existingVendor, err := s.repo.FindByMobileAndCompany(
			input.Mobile,
			companyID,
		)

		if err == nil && existingVendor != nil {
			return nil, fmt.Errorf(
				"mobile number %s already exists with vendor: %s",
				input.Mobile,
				existingVendor.DisplayName,
			)
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

	return s.GetVendorByIDAndCompany(
		vendor.ID,
		companyID,
	)
}

// Existing compatibility method retained.
func (s *vendorService) GetVendorByIDAndUser(
	id uint,
	userID uint,
) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("vendor not found")
	}

	return helper.MapVendorToOutput(vendor), nil
}

// Existing compatibility method retained.
func (s *vendorService) GetVendorsByUser(
	userID uint,
	companyID uint,
	page int,
	limit int,
) ([]output.VendorListOutput, int64, error) {
	_ = userID

	return s.GetVendorsByCompany(
		companyID,
		page,
		limit,
	)
}

// Existing compatibility method retained.
func (s *vendorService) UpdateVendorForUser(
	id uint,
	userID uint,
	input *input.UpdateVendorInput,
) (*output.VendorOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	vendor, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("vendor not found")
	}

	helper.ApplyUpdateVendorInput(vendor, input)

	if err := s.repo.Update(vendor); err != nil {
		return nil, err
	}

	return helper.MapVendorToOutput(vendor), nil
}

// Existing compatibility method retained.
func (s *vendorService) DeleteVendorForUser(
	id uint,
	userID uint,
) error {
	vendor, err := s.repo.FindByIDAndUser(id, userID)
	if err != nil {
		return fmt.Errorf("vendor not found")
	}

	return s.repo.Delete(vendor.ID)
}

func (s *vendorService) GetVendorByIDAndCompany(
	id uint,
	companyID uint,
) (*output.VendorOutput, error) {
	vendor, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("vendor not found")
	}

	return helper.MapVendorToOutput(vendor), nil
}

func (s *vendorService) GetVendorsByCompany(
	companyID uint,
	page int,
	limit int,
) ([]output.VendorListOutput, int64, error) {
	if companyID == 0 {
		return nil, 0, errors.New("invalid company")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	vendors, total, err := s.repo.FindByCompanyID(
		companyID,
		page,
		limit,
	)
	if err != nil {
		return nil, 0, err
	}

	return vendorListOutputs(vendors), total, nil
}

func (s *vendorService) UpdateVendorForCompany(
	id uint,
	companyID uint,
	input *input.UpdateVendorInput,
) (*output.VendorOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	vendor, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("vendor not found")
	}

	if input.Mobile != nil &&
		*input.Mobile != "" &&
		*input.Mobile != vendor.Mobile {
		existingVendor, err := s.repo.FindByMobileAndCompany(
			*input.Mobile,
			companyID,
		)

		if err == nil &&
			existingVendor != nil &&
			existingVendor.ID != id {
			return nil, fmt.Errorf(
				"mobile number %s already exists with vendor: %s",
				*input.Mobile,
				existingVendor.DisplayName,
			)
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	helper.ApplyUpdateVendorInput(vendor, input)

	if err := s.repo.UpdateByCompanyID(
		vendor,
		companyID,
	); err != nil {
		return nil, err
	}

	return s.GetVendorByIDAndCompany(
		vendor.ID,
		companyID,
	)
}

func (s *vendorService) DeleteVendorForCompany(
	id uint,
	companyID uint,
) error {
	return s.repo.DeleteByIDAndCompany(
		id,
		companyID,
	)
}