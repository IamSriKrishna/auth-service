package services

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type ManufacturerService interface {
	Create(input *input.CreateManufacturerInput) (*output.ManufacturerOutput, error)
	GetByID(id uint) (*output.ManufacturerOutput, error)
	GetAll(limit, offset int) (*output.ListManufacturersOutput, error)
	Update(id uint, input *input.UpdateManufacturerInput) (*output.ManufacturerOutput, error)
	Delete(id uint) error
}

type manufacturerService struct {
	repo repo.ManufacturerRepository
}

func NewManufacturerService(repo repo.ManufacturerRepository) ManufacturerService {
	return &manufacturerService{repo: repo}
}

func (s *manufacturerService) Create(input *input.CreateManufacturerInput) (*output.ManufacturerOutput, error) {
	manufacturer := &models.Manufacturer{
		Name: input.Name,
	}
	err := s.repo.Create(manufacturer)
	if err != nil {
		return nil, err
	}
	return &output.ManufacturerOutput{
		ID:        manufacturer.ID,
		Name:      manufacturer.Name,
		CreatedAt: manufacturer.CreatedAt,
		UpdatedAt: manufacturer.UpdatedAt,
	}, nil
}

func (s *manufacturerService) GetByID(id uint) (*output.ManufacturerOutput, error) {
	manufacturer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err	
	}
	return &output.ManufacturerOutput{
		ID:        manufacturer.ID,
		Name:      manufacturer.Name,
		CreatedAt: manufacturer.CreatedAt,
		UpdatedAt: manufacturer.UpdatedAt,
	}, nil
}

func (s *manufacturerService) GetAll(limit, offset int) (*output.ListManufacturersOutput, error) {
	manufacturers, count, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}
	outputs := make([]output.ManufacturerOutput, len(manufacturers))
	for i, m := range manufacturers {
		outputs[i] = output.ManufacturerOutput{
			ID:        m.ID,
			Name:      m.Name,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return &output.ListManufacturersOutput{
		Manufacturers: outputs,
		TotalCount:    int(count),
	}, nil
}

func (s *manufacturerService) Update(id uint, input *input.UpdateManufacturerInput) (*output.ManufacturerOutput, error) {
	manufacturer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	manufacturer.Name = *input.Name
	err = s.repo.Update(manufacturer)
	if err != nil {
		return nil, err
	}
	return &output.ManufacturerOutput{
		ID:        manufacturer.ID,
		Name:      manufacturer.Name,
		CreatedAt: manufacturer.CreatedAt,
		UpdatedAt: manufacturer.UpdatedAt,
	}, nil
}

func (s *manufacturerService) Delete(id uint) error {
	return s.repo.Delete(id)
}