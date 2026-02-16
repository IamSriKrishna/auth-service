package services

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/utils"
)

// adminService implements AdminService interface
type adminService struct {
	userRepo repo.UserRepository
	roleRepo repo.RoleRepository
}

// NewAdminService creates a new admin service instance
func NewAdminService(
	userRepo repo.UserRepository,
	roleRepo repo.RoleRepository,
) AdminService {
	return &adminService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateUser creates a new admin or partner user
func (s *adminService) CreateUser(ctx context.Context, createdBy uint, req *input.CreateUserRequest) (*output.UserInfo, error) {
	// Validate user type
	if req.UserType != "admin" && req.UserType != "partner" {
		return nil, errors.New("invalid user type")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Get role by name
	role, err := s.roleRepo.GetByName(req.RoleName)
	if err != nil {
		return nil, errors.New("invalid role name")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        &req.Email,
		Username:     &req.Username,
		PasswordHash: &passwordHash,
		UserType:     models.UserType(req.UserType),
		RoleID:       role.ID,
		Status:       models.UserStatusActive,
		Phone:        &req.Phone,
		CreatedBy:    &createdBy,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Return user info
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

// ResetPassword resets user password
func (s *adminService) ResetPassword(ctx context.Context, req *input.ResetPasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Only allow password reset for admin and partner users
	if user.UserType == models.UserTypeMobile {
		return errors.New("cannot reset password for mobile users")
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = &passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Update password changed timestamp
	s.userRepo.UpdatePasswordChangedAt(req.UserID)

	return nil
}

// ResetPassword resets user password
func (s *adminService) ResetUserPassword(ctx context.Context, req *input.ResetUserPasswordRequest, userID uint64) error {
	// Get user
	user, err := s.userRepo.GetByID(uint(userID))
	if err != nil {
		return errors.New("user not found")
	}

	// Only allow password reset for admin and partner users
	if user.UserType == models.UserTypeMobile {
		return errors.New("cannot reset password for mobile users")
	}

	// Verify old password
	if !utils.CheckPassword(*user.PasswordHash, req.OldPassword) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = &passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Update password changed timestamp
	s.userRepo.UpdatePasswordChangedAt(uint(userID))

	return nil
}

// GetUsers retrieves users with pagination and search
func (s *adminService) GetUsers(ctx context.Context, page, limit int, search string) (*output.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, total, err := s.userRepo.List(offset, limit, search)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	userList := make([]output.UserListResponse, len(users))
	for i, user := range users {
		userList[i] = output.UserListResponse{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			Phone:       user.Phone,
			UserType:    string(user.UserType),
			Role:        user.Role.RoleName,
			Status:      string(user.Status),
			CreatedAt:   user.CreatedAt,
			CreatedBy:   user.CreatedBy,
			LastLoginAt: user.LastLoginAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &output.PaginatedResponse{
		Success: true,
		Data:    userList,
		Meta: output.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       int(total),
			TotalPages:  totalPages,
		},
	}, nil
}

// GetUser retrieves a specific user
func (s *adminService) GetUser(ctx context.Context, userID uint) (*output.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &output.UserInfo{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		Username:    user.Username,
		UserType:    string(user.UserType),
		Role:        user.Role.RoleName,
		Status:      string(user.Status),
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
	}, nil
}

// UpdateUser updates user information
func (s *adminService) UpdateUser(ctx context.Context, userID uint, req *input.UpdateUserRequest) (*output.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Username != nil {
		user.Username = req.Username
	}
	if req.RoleName != nil {
		// Get role by name
		role, err := s.roleRepo.GetByName(*req.RoleName)
		if err != nil {
			return nil, errors.New("invalid role name")
		}
		user.RoleID = role.ID
	}
	if req.Status != nil {
		user.Status = models.UserStatus(*req.Status)
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Get updated user with role information
	updatedUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &output.UserInfo{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		Phone:       updatedUser.Phone,
		Username:    updatedUser.Username,
		UserType:    string(updatedUser.UserType),
		Role:        updatedUser.Role.RoleName,
		Status:      string(updatedUser.Status),
		CreatedAt:   updatedUser.CreatedAt,
		LastLoginAt: updatedUser.LastLoginAt,
	}, nil
}

// DeleteUser soft deletes a user
func (s *adminService) DeleteUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Prevent deletion of superadmin users
	if user.UserType == models.UserTypeSuperAdmin {
		return errors.New("cannot delete superadmin user")
	}

	return s.userRepo.Delete(userID)
}

// UpdateUserStatus updates user status
func (s *adminService) UpdateUserStatus(ctx context.Context, userID uint, status string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Status = models.UserStatus(status)
	return s.userRepo.Update(user)
}

// UpdateUserRole updates user role
func (s *adminService) UpdateUserRole(ctx context.Context, userID uint, roleName string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Get role by name
	role, err := s.roleRepo.GetByName(roleName)
	if err != nil {
		return errors.New("invalid role name")
	}

	user.RoleID = role.ID
	return s.userRepo.Update(user)
}

// GetDashboardStats retrieves dashboard statistics with filters
func (s *adminService) GetDashboardStats(ctx context.Context, filter *input.DashboardStatsFilter) (*output.DashboardStatsResponse, error) {
	var fromDate, toDate *time.Time

	// Parse date strings if provided
	if filter.FromDate != nil && *filter.FromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", *filter.FromDate)
		if err != nil {
			return nil, errors.New("invalid from_date format. Use YYYY-MM-DD")
		}
		fromDate = &parsedFrom
	}

	if filter.ToDate != nil && *filter.ToDate != "" {
		parsedTo, err := time.Parse("2006-01-02", *filter.ToDate)
		if err != nil {
			return nil, errors.New("invalid to_date format. Use YYYY-MM-DD")
		}
		// Set to end of day
		endOfDay := parsedTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		toDate = &endOfDay
	}

	// Get stats from repository
	stats, err := s.userRepo.GetDashboardStats(
		filter.CustomerType,
		fromDate,
		toDate,
	)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &output.DashboardStatsResponse{
		TotalUsers:              stats["total_users"].(int),
		ActiveUsers:             stats["active_users"].(int),
		MembershipUsers:         stats["membership_users"].(int),
		MembershipLastUpdatedAt: stats["membership_last_updated_at"].(string),
		Filters: output.DashboardStatsFilterApplied{
			CustomerType: filter.CustomerType,
			FromDate:     filter.FromDate,
			ToDate:       filter.ToDate,
		},
	}

	return response, nil
}
