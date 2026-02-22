package services

import (
	"errors"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type BankService interface {
	Create(input *input.CreateBankInput) (*output.BankOutput, error)
	FindByID(id uint) (*output.BankOutput, error)
	FindAll(limit, offset int) ([]output.BankOutput, int64, error)
	Update(id uint, input *input.UpdateBankInput) (*output.BankOutput, error)
	Delete(id uint) error
}

type bankService struct {
	repo repo.BankRepository
}

func NewBankService(repo repo.BankRepository) BankService {
	return &bankService{repo: repo}
}

func (s *bankService) Create(input *input.CreateBankInput) (*output.BankOutput, error) {
	if input == nil {
		return nil, errors.New("invalid input")
	}

	bank := &models.Bank{
		BankName:   input.BankName,
		Address:    input.Address,
		City:       input.City,
		State:      input.State,
		PostalCode: input.PostalCode,
		Country:    input.Country,
		IsActive:   input.IsActive,
	}

	err := s.repo.Create(bank)
	if err != nil {
		return nil, err
	}

	return &output.BankOutput{
		ID:         bank.ID,
		BankName:   bank.BankName,
		Address:    bank.Address,
		City:       bank.City,
		State:      bank.State,
		PostalCode: bank.PostalCode,
		Country:    bank.Country,
		IsActive:   bank.IsActive,
		CreatedAt:  bank.CreatedAt,
		UpdatedAt:  bank.UpdatedAt,
	}, nil
}

func (s *bankService) FindByID(id uint) (*output.BankOutput, error) {
	bank, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &output.BankOutput{
		ID:         bank.ID,
		BankName:   bank.BankName,
		Address:    bank.Address,
		City:       bank.City,
		State:      bank.State,
		PostalCode: bank.PostalCode,
		Country:    bank.Country,
		IsActive:   bank.IsActive,
		CreatedAt:  bank.CreatedAt,
		UpdatedAt:  bank.UpdatedAt,
	}, nil
}

func (s *bankService) FindAll(limit, offset int) ([]output.BankOutput, int64, error) {
	banks, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var bankOutputs []output.BankOutput
	for _, bank := range banks {
		bankOutputs = append(bankOutputs, output.BankOutput{
			ID:         bank.ID,
			BankName:   bank.BankName,
			Address:    bank.Address,
			City:       bank.City,
			State:      bank.State,
			PostalCode: bank.PostalCode,
			Country:    bank.Country,
			IsActive:   bank.IsActive,
			CreatedAt:  bank.CreatedAt,
			UpdatedAt:  bank.UpdatedAt,
		})
	}
	return bankOutputs, total, nil
}

func (s *bankService) Update(id uint, input *input.UpdateBankInput) (*output.BankOutput, error) {
	bank, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.BankName != nil {
		bank.BankName = *input.BankName
	}
	if input.Address != nil {
		bank.Address = *input.Address
	}
	if input.City != nil {
		bank.City = *input.City
	}
	if input.State != nil {
		bank.State = *input.State
	}
	if input.PostalCode != nil {
		bank.PostalCode = *input.PostalCode
	}
	if input.Country != nil {
		bank.Country = *input.Country
	}
	if input.IsActive != nil {
		bank.IsActive = *input.IsActive
	}

	err = s.repo.Update(bank)
	if err != nil {
		return nil, err
	}

	return &output.BankOutput{
		ID:         bank.ID,
		BankName:   bank.BankName,
		Address:    bank.Address,
		City:       bank.City,
		State:      bank.State,
		PostalCode: bank.PostalCode,
		Country:    bank.Country,
		IsActive:   bank.IsActive,
		CreatedAt:  bank.CreatedAt,
		UpdatedAt:  bank.UpdatedAt,
	}, nil
}

func (s *bankService) Delete(id uint) error {
	return s.repo.Delete(id)
}
