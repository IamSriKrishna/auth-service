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

type adminService struct {
	userRepo repo.UserRepository
	roleRepo repo.RoleRepository
}

func NewAdminService(
	userRepo repo.UserRepository,
	roleRepo repo.RoleRepository,
) AdminService {
	return &adminService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *adminService) CreateUser(ctx context.Context, createdBy uint, req *input.CreateUserRequest) (*output.UserInfo, error) {
	if req.UserType != "admin" && req.UserType != "partner" {
		return nil, errors.New("invalid user type")
	}

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	role, err := s.roleRepo.GetByName(req.RoleName)
	if err != nil {
		return nil, errors.New("invalid role name")
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

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
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.UserType == models.UserTypeMobile {
		return errors.New("cannot reset password for mobile users")
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = &passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.userRepo.UpdatePasswordChangedAt(req.UserID)

	return nil
}

func (s *adminService) ResetUserPassword(ctx context.Context, req *input.ResetUserPasswordRequest, userID uint64) error {
	user, err := s.userRepo.GetByID(uint(userID))
	if err != nil {
		return errors.New("user not found")
	}

	if user.UserType == models.UserTypeMobile {
		return errors.New("cannot reset password for mobile users")
	}

	if !utils.CheckPassword(*user.PasswordHash, req.OldPassword) {
		return errors.New("current password is incorrect")
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = &passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.userRepo.UpdatePasswordChangedAt(uint(userID))

	return nil
}

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

func (s *adminService) UpdateUser(ctx context.Context, userID uint, req *input.UpdateUserRequest) (*output.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Username != nil {
		user.Username = req.Username
	}
	if req.RoleName != nil {
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

func (s *adminService) DeleteUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.UserType == models.UserTypeSuperAdmin {
		return errors.New("cannot delete superadmin user")
	}

	return s.userRepo.Delete(userID)
}

func (s *adminService) UpdateUserStatus(ctx context.Context, userID uint, status string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Status = models.UserStatus(status)
	return s.userRepo.Update(user)
}

func (s *adminService) UpdateUserRole(ctx context.Context, userID uint, roleName string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	role, err := s.roleRepo.GetByName(roleName)
	if err != nil {
		return errors.New("invalid role name")
	}

	user.RoleID = role.ID
	return s.userRepo.Update(user)
}

func (s *adminService) GetDashboardStats(ctx context.Context, filter *input.DashboardStatsFilter) (*output.DashboardStatsResponse, error) {
	var fromDate, toDate *time.Time

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
		endOfDay := parsedTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		toDate = &endOfDay
	}

	stats, err := s.userRepo.GetDashboardStats(
		filter.CustomerType,
		fromDate,
		toDate,
	)
	if err != nil {
		return nil, err
	}

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
