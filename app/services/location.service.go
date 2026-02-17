package services

import (
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/repo"
)
type LocationService interface {
	GetAllCountries() ([]output.CountryOutput, error)
	GetCountryByID(id uint) (*output.CountryOutput, error)
	GetStatesByCountry(countryID uint) ([]output.StateOutput, error)
	GetStateByID(id uint) (*output.StateOutput, error)
}

type locationService struct {
	repo repo.LocationRepository
}

func NewLocationService(repo repo.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

func (s *locationService) GetAllCountries() ([]output.CountryOutput, error) {
	countries, err := s.repo.GetAllCountries()
	if err != nil {
		return nil, err
	}

	outputs := make([]output.CountryOutput, len(countries))
	for i, c := range countries {
		outputs[i] = output.CountryOutput{
			ID:          c.ID,
			CountryName: c.CountryName,
			CountryCode: c.CountryCode,
			PhoneCode:   c.PhoneCode,
		}
	}

	return outputs, nil
}

func (s *locationService) GetCountryByID(id uint) (*output.CountryOutput, error) {
	country, err := s.repo.GetCountryByID(id)
	if err != nil {
		return nil, err
	}

	return &output.CountryOutput{
		ID:          country.ID,
		CountryName: country.CountryName,
		CountryCode: country.CountryCode,
		PhoneCode:   country.PhoneCode,
	}, nil
}

func (s *locationService) GetStatesByCountry(countryID uint) ([]output.StateOutput, error) {
	states, err := s.repo.GetStatesByCountry(countryID)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.StateOutput, len(states))
	for i, s := range states {
		outputs[i] = output.StateOutput{
			ID:        s.ID,
			CountryID: s.CountryID,
			StateName: s.StateName,
			StateCode: s.StateCode,
		}
	}

	return outputs, nil
}

func (s *locationService) GetStateByID(id uint) (*output.StateOutput, error) {
	state, err := s.repo.GetStateByID(id)
	if err != nil {
		return nil, err
	}

	return &output.StateOutput{
		ID:        state.ID,
		CountryID: state.CountryID,
		StateName: state.StateName,
		StateCode: state.StateCode,
		Country: output.CountryOutput{
			ID:          state.Country.ID,
			CountryName: state.Country.CountryName,
			CountryCode: state.Country.CountryCode,
			PhoneCode:   state.Country.PhoneCode,
		},
	}, nil
}
