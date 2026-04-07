package services

import (
	"context"
	"errors"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/utils"
)

type adminService struct {
	userRepo    repo.UserRepository
	roleRepo    repo.RoleRepository
	companyRepo repo.CompanyRepository
}

func NewAdminService(
	userRepo repo.UserRepository,
	roleRepo repo.RoleRepository,
	companyRepo repo.CompanyRepository,
) AdminService {
	return &adminService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		companyRepo: companyRepo,
	}
}

func (s *adminService) CreateUser(ctx context.Context, createdBy *uint, req *input.CreateUserRequest) (*output.UserInfo, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	email := req.Email
	phone := req.Number
	username := req.Username
	passwordHashPtr := &passwordHash
	userType := models.UserType(req.UserType)
	companyID := req.CompanyID

	user := &models.User{
		Email:        &email,
		Phone:        &phone,
		Username:     &username,
		PasswordHash: passwordHashPtr,
		Status:       models.UserStatusActive,
		UserType:     userType,
		CreatedBy:    createdBy,
		CompanyID:    &companyID,
	}

	var role *models.Role
	if req.RoleName != "" {
		var err2 error
		role, err2 = s.roleRepo.GetByName(req.RoleName)
		if err2 != nil || role == nil {
			return nil, errors.New("role does not exist")
		}
		user.RoleID = role.ID
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	roleName := ""
	if role != nil && role.RoleName != "" {
		roleName = role.RoleName
	}

	var companyName *string
	if user.CompanyID != nil && *user.CompanyID > 0 {
		company, err := s.companyRepo.FindByID(*user.CompanyID)
		if err == nil && company != nil {
			companyName = &company.CompanyName
		}
	}

	return &output.UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		Username:    user.Username,
		UserType:    string(user.UserType),
		Role:        roleName,
		Status:      string(user.Status),
		CompanyID:   user.CompanyID,
		CompanyName: companyName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *adminService) CreateSuperAdmin(ctx context.Context, createdBy *uint, req *input.CreateSuperAdminRequest) (*output.UserInfo, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	email := req.Email
	phone := req.Phone
	passwordHashPtr := &passwordHash

	user := &models.User{
		Email:        &email,
		Phone:        &phone,
		Username:     &req.Username,
		PasswordHash: passwordHashPtr,
		Status:       models.UserStatusActive,
		UserType:     models.UserTypeSuperAdmin,
		CreatedBy:    createdBy,
	}

	// Get super admin role
	role, err := s.roleRepo.GetByName("superadmin")
	if err != nil || role == nil {
		return nil, errors.New("superadmin role does not exist")
	}
	user.RoleID = role.ID

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &output.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Username:  user.Username,
		UserType:  string(user.UserType),
		Role:      role.RoleName,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *adminService) ResetPassword(ctx context.Context, req *input.ResetPasswordRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = &passwordHash
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *adminService) ResetUserPassword(ctx context.Context, req *input.ResetUserPasswordRequest, userID uint64) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	user, err := s.userRepo.GetByID(uint(userID))
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = &passwordHash
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *adminService) GetUsers(ctx context.Context, page, limit int, search, role string) (*output.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []models.User
	var err error

	// Filter to show only admin users (not superadmin)
	// Default to "admin" role filtering
	filterRole := "admin"
	if role != "" {
		filterRole = role
	}

	users, _, err = s.userRepo.ListWithFilters(offset, limit, search, filterRole)
	if err != nil {
		return nil, err
	}

	// Additional filter: exclude superadmin users, only show admin user type
	var filteredUsers []models.User
	for _, user := range users {
		// Only include users with userType "admin", exclude superadmin
		if user.UserType != models.UserTypeSuperAdmin {
			filteredUsers = append(filteredUsers, user)
		}
	}

	userInfos := make([]interface{}, len(filteredUsers))
	for i, user := range filteredUsers {
		var companyName *string
		if user.CompanyID != nil && *user.CompanyID > 0 {
			company, err := s.companyRepo.FindByID(*user.CompanyID)
			if err == nil && company != nil {
				companyName = &company.CompanyName
			}
		}
		userInfos[i] = &output.UserInfo{
			ID:          user.ID,
			Email:       user.Email,
			Phone:       user.Phone,
			Username:    user.Username,
			UserType:    string(user.UserType),
			Role:        user.Role.RoleName,
			Status:      string(user.Status),
			CompanyID:   user.CompanyID,
			CompanyName: companyName,
			CreatedAt:   user.CreatedAt,
		}
	}

	totalPages := int((int(len(filteredUsers)) + limit - 1) / limit)
	return &output.PaginatedResponse{
		Success: true,
		Data:    userInfos,
		Meta: output.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       len(filteredUsers),
			TotalPages:  totalPages,
		},
	}, nil
}

func (s *adminService) GetUser(ctx context.Context, userID uint) (*output.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	var companyName *string
	if user.CompanyID != nil && *user.CompanyID > 0 {
		company, err := s.companyRepo.FindByID(*user.CompanyID)
		if err == nil && company != nil {
			companyName = &company.CompanyName
		}
	}

	return &output.UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		Username:    user.Username,
		UserType:    string(user.UserType),
		Role:        user.Role.RoleName,
		Status:      string(user.Status),
		CompanyID:   user.CompanyID,
		CompanyName: companyName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *adminService) UpdateUser(ctx context.Context, userID uint, req *input.UpdateUserRequest) (*output.UserInfo, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Username != nil {
		user.Username = req.Username
	}

	if req.RoleName != nil {
		role, err := s.roleRepo.GetByName(*req.RoleName)
		if err != nil || role == nil {
			return nil, errors.New("role does not exist")
		}
		user.RoleID = role.ID
	}

	if req.Status != nil {
		user.Status = models.UserStatus(*req.Status)
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	var companyName *string
	if user.CompanyID != nil && *user.CompanyID > 0 {
		company, err := s.companyRepo.FindByID(*user.CompanyID)
		if err == nil && company != nil {
			companyName = &company.CompanyName
		}
	}

	return &output.UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		Username:    user.Username,
		UserType:    string(user.UserType),
		Role:        user.Role.RoleName,
		Status:      string(user.Status),
		CompanyID:   user.CompanyID,
		CompanyName: companyName,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *adminService) DeleteUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	err = s.userRepo.Delete(userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *adminService) UpdateUserStatus(ctx context.Context, userID uint, status string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"pending":  true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	user.Status = models.UserStatus(status)
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *adminService) UpdateUserRole(ctx context.Context, userID uint, roleName string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	role, err := s.roleRepo.GetByName(roleName)
	if err != nil || role == nil {
		return errors.New("role does not exist")
	}

	user.RoleID = role.ID
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *adminService) GetDashboardStats(ctx context.Context, filter *input.DashboardStatsFilter) (*output.DashboardStatsResponse, error) {
	stats := &output.DashboardStatsResponse{
		TotalUsers:  0,
		ActiveUsers: 0,
		Filters: output.DashboardStatsFilterApplied{
			CustomerType: filter.CustomerType,
			FromDate:     filter.FromDate,
			ToDate:       filter.ToDate,
		},
	}

	// TODO: Implement actual statistics calculation based on your requirements
	// This is a placeholder implementation

	return stats, nil
}
