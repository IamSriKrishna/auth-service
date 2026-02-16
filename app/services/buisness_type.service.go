package services

import (
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/repo"
)
type BusinessTypeService interface {
	GetAll() ([]output.BusinessTypeOutput, error)
	GetByID(id uint) (*output.BusinessTypeOutput, error)
}

type businessTypeService struct {
	repo repo.BusinessTypeRepository
}

func NewBusinessTypeService(repo repo.BusinessTypeRepository) BusinessTypeService {
	return &businessTypeService{repo: repo}
}

func (s *businessTypeService) GetAll() ([]output.BusinessTypeOutput, error) {
	businessTypes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	outputs := make([]output.BusinessTypeOutput, len(businessTypes))
	for i, bt := range businessTypes {
		outputs[i] = output.BusinessTypeOutput{
			ID:          bt.ID,
			TypeName:    bt.TypeName,
			Description: bt.Description,
			IsActive:    bt.IsActive,
			CreatedAt:   bt.CreatedAt,
		}
	}

	return outputs, nil
}

func (s *businessTypeService) GetByID(id uint) (*output.BusinessTypeOutput, error) {
	bt, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &output.BusinessTypeOutput{
		ID:          bt.ID,
		TypeName:    bt.TypeName,
		Description: bt.Description,
		IsActive:    bt.IsActive,
		CreatedAt:   bt.CreatedAt,
	}, nil
}