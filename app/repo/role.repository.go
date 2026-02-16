package repo

import (
	"github.com/bbapp-org/auth-service/app/models"

	"gorm.io/gorm"
)

// roleRepository implements RoleRepository interface
type roleRepository struct {
	db *gorm.DB // Database with dbresolver (handles read/write splitting automatically)
}

// NewRoleRepository creates a new role repository instance
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name
func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAll retrieves all roles
func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("is_active = ?", true).Find(&roles).Error
	return roles, err
}

// Create creates a new role
func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// Update updates a role
func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete soft deletes a role
func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
