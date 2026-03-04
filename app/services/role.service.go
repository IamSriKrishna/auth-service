package services

import (
	"errors"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type RoleService interface {
	CreateRole(req *input.CreateRoleRequest) (*output.RoleResponse, error)
	GetRole(id uint) (*output.RoleResponse, error)
	GetAllRoles() ([]output.RoleResponse, error)
}

type roleService struct {
	roleRepo repo.RoleRepository
}

func NewRoleService(roleRepo repo.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) CreateRole(req *input.CreateRoleRequest) (*output.RoleResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Check if role already exists
	existing, err := s.roleRepo.GetByName(req.RoleName)
	if err == nil && existing != nil {
		return nil, errors.New("role already exists with this name")
	}

	role := &models.Role{
		RoleName:    req.RoleName,
		Permissions: models.StringArray(req.Permissions),
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	err = s.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	return &output.RoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Permissions: []string(role.Permissions),
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

func (s *roleService) GetRole(id uint) (*output.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &output.RoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Permissions: []string(role.Permissions),
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

func (s *roleService) GetAllRoles() ([]output.RoleResponse, error) {
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []output.RoleResponse
	for _, role := range roles {
		responses = append(responses, output.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Permissions: []string(role.Permissions),
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	return responses, nil
}
