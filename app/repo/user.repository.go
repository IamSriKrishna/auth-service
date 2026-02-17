package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/utils"

	"gorm.io/gorm"
)

type userRepository struct {
	db         *gorm.DB
	httpClient *utils.HTTPClient
}

func NewUserRepository(db *gorm.DB, httpClient *utils.HTTPClient) UserRepository {
	return &userRepository{
		db:         db,
		httpClient: httpClient,
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByAppleID(appleID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("apple_id = ?", appleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(offset, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("email LIKE ? OR phone LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Role").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

func (r *userRepository) UpdatePasswordChangedAt(id uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("password_changed_at", now).Error
}

func (r *userRepository) GetDashboardStats(customerType *string, fromDate, toDate *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	baseQuery := r.db.Model(&models.User{})

	if customerType != nil && *customerType != "" {
		baseQuery = baseQuery.Where("user_type = ?", *customerType)
	} else {
		baseQuery = baseQuery.Where("user_type = ?", models.UserTypeMobile)
	}

	if fromDate != nil {
		baseQuery = baseQuery.Where("created_at >= ?", *fromDate)
	}
	if toDate != nil {
		baseQuery = baseQuery.Where("created_at <= ?", *toDate)
	}
	var totalUsers int64
	if err := baseQuery.Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	stats["total_users"] = int(totalUsers)

	var activeUsers int64
	activeQuery := baseQuery.Where("status = ?", models.UserStatusActive)
	if err := activeQuery.Count(&activeUsers).Error; err != nil {
		return nil, err
	}
	stats["active_users"] = int(activeUsers)

	membershipData, err := r.httpClient.GetMembershipCount()
	if err != nil {
		println("ERROR: Failed to fetch membership count from customer service:", err.Error())
		stats["membership_users"] = 0
		stats["membership_last_updated_at"] = ""
	} else {
		stats["membership_users"] = membershipData.Count
		stats["membership_last_updated_at"] = membershipData.LastUpdatedAt
	}

	return stats, nil
}
