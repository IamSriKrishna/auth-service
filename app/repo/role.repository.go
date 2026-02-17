package repo

import (
	"github.com/bbapp-org/auth-service/app/models"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("is_active = ?", true).Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
