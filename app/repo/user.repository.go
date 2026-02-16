package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/utils"

	"gorm.io/gorm"
)

// userRepository implements UserRepository interface
type userRepository struct {
	db         *gorm.DB          // Database with dbresolver (handles read/write splitting automatically)
	httpClient *utils.HTTPClient // HTTP client for calling external services
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB, httpClient *utils.HTTPClient) UserRepository {
	return &userRepository{
		db:         db,
		httpClient: httpClient,
	}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone retrieves a user by phone
func (r *userRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByGoogleID retrieves a user by Google ID
func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByAppleID retrieves a user by Apple ID
func (r *userRepository) GetByAppleID(appleID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("apple_id = ?", appleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List retrieves users with pagination and search
func (r *userRepository) List(offset, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Add search conditions if search term is provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("email LIKE ? OR phone LIKE ?", searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Preload("Role").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

// UpdateLastLogin updates the last login time for a user
func (r *userRepository) UpdateLastLogin(id uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// UpdatePasswordChangedAt updates the password changed timestamp
func (r *userRepository) UpdatePasswordChangedAt(id uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("password_changed_at", now).Error
}

// GetDashboardStats retrieves dashboard statistics with filters
func (r *userRepository) GetDashboardStats(customerType *string, fromDate, toDate *time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Base query - filter by user_type if customer_type is specified, otherwise all users
	baseQuery := r.db.Model(&models.User{})

	// Apply customer_type filter
	if customerType != nil && *customerType != "" {
		baseQuery = baseQuery.Where("user_type = ?", *customerType)
	} else {
		// Default to mobile users only if no customer_type specified
		baseQuery = baseQuery.Where("user_type = ?", models.UserTypeMobile)
	}

	// Apply date filters
	if fromDate != nil {
		baseQuery = baseQuery.Where("created_at >= ?", *fromDate)
	}
	if toDate != nil {
		baseQuery = baseQuery.Where("created_at <= ?", *toDate)
	} // Get total users with filters
	var totalUsers int64
	if err := baseQuery.Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	stats["total_users"] = int(totalUsers)

	// Get active users count
	var activeUsers int64
	activeQuery := baseQuery.Where("status = ?", models.UserStatusActive)
	if err := activeQuery.Count(&activeUsers).Error; err != nil {
		return nil, err
	}
	stats["active_users"] = int(activeUsers)

	// Get membership users count from customer service
	membershipData, err := r.httpClient.GetMembershipCount()
	if err != nil {
		// Log error but don't fail the entire request, return 0 for membership users
		println("ERROR: Failed to fetch membership count from customer service:", err.Error())
		stats["membership_users"] = 0
		stats["membership_last_updated_at"] = ""
	} else {
		stats["membership_users"] = membershipData.Count
		stats["membership_last_updated_at"] = membershipData.LastUpdatedAt
	}

	return stats, nil
}
