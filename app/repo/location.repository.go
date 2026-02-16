package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetAllCountries() ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Where("is_active = ?", true).Order("country_name ASC").Find(&countries).Error
	return countries, err
}

func (r *locationRepository) GetCountryByID(id uint) (*models.Country, error) {
	var country models.Country
	err := r.db.First(&country, id).Error
	return &country, err
}

func (r *locationRepository) GetStatesByCountry(countryID uint) ([]models.State, error) {
	var states []models.State
	err := r.db.Where("country_id = ? AND is_active = ?", countryID, true).
		Order("state_name ASC").Find(&states).Error
	return states, err
}

func (r *locationRepository) GetStateByID(id uint) (*models.State, error) {
	var state models.State
	err := r.db.Preload("Country").First(&state, id).Error
	return &state, err
}
