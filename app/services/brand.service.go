package services

import (
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type BrandService interface {
	Create(input *input.CreateBrandInput) (*output.BrandOutput, error)
	FindByID(id uint) (*output.BrandOutput, error)
	FindAll(limit, offset int) ([]output.BrandOutput, int64, error)
	Update(id uint, input *input.UpdateBrandInput) (*output.BrandOutput, error)
	Delete(id uint) error
}

type brandService struct {
	repo repo.BrandRepository
}

func NewBrandService(repo repo.BrandRepository) BrandService {
	return &brandService{repo: repo}
}

func (s *brandService) Create(input *input.CreateBrandInput) (*output.BrandOutput, error) {
	brand := &models.Brand{
		Name: input.Name,
	}

	err := s.repo.Create(brand)
	if err != nil {
		return nil, err
	}

	return &output.BrandOutput{
		ID:   brand.ID,
		Name: brand.Name,
	}, nil
}

func (s *brandService) FindByID(id uint) (*output.BrandOutput, error) {
	brand, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &output.BrandOutput{
		ID:   brand.ID,
		Name: brand.Name,
	}, nil
}

func (s *brandService) FindAll(limit, offset int) ([]output.BrandOutput, int64, error) {
	brands, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var brandOutputs []output.BrandOutput
	for _, brand := range brands {
		brandOutputs = append(brandOutputs, output.BrandOutput{
			ID:   brand.ID,
			Name: brand.Name,
		})
	}

	return brandOutputs, total, nil
}

func (s *brandService) Update(id uint, input *input.UpdateBrandInput) (*output.BrandOutput, error) {
	brand, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		brand.Name = *input.Name
	}

	err = s.repo.Update(brand)
	if err != nil {
		return nil, err
	}

	return &output.BrandOutput{
		ID:   brand.ID,
		Name: brand.Name,
	}, nil
}

func (s *brandService) Delete(id uint) error {
	return s.repo.Delete(id)
}
