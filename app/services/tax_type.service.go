package services

import (
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/repo"
)
type TaxTypeService interface {
	GetAll() ([]output.TaxTypeOutput, error)
	GetByID(id uint) (*output.TaxTypeOutput, error)
}

type taxTypeService struct {
	repo repo.TaxTypeRepository
}

func NewTaxTypeService(repo repo.TaxTypeRepository) TaxTypeService {
	return &taxTypeService{repo: repo}
}

func (s *taxTypeService) GetAll() ([]output.TaxTypeOutput, error) {
	taxTypes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	outputs := make([]output.TaxTypeOutput, len(taxTypes))
	for i, tt := range taxTypes {
		outputs[i] = output.TaxTypeOutput{
			ID:          tt.ID,
			TaxName:     tt.TaxName,
			TaxCode:     tt.TaxCode,
			Description: tt.Description,
		}
	}

	return outputs, nil
}

func (s *taxTypeService) GetByID(id uint) (*output.TaxTypeOutput, error) {
	tt, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &output.TaxTypeOutput{
		ID:          tt.ID,
		TaxName:     tt.TaxName,
		TaxCode:     tt.TaxCode,
		Description: tt.Description,
	}, nil
}
