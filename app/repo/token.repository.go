package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"

	"gorm.io/gorm"
)

// refreshTokenRepository implements RefreshTokenRepository interface
type refreshTokenRepository struct {
	db *gorm.DB // Database with dbresolver (handles read/write splitting automatically)
}

// NewRefreshTokenRepository creates a new refresh token repository instance
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetByTokenID retrieves a refresh token by token ID
func (r *refreshTokenRepository) GetByTokenID(tokenID string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Preload("User").
		Where("token_id = ? AND is_revoked = ? AND expires_at > ?", tokenID, false, time.Now()).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *refreshTokenRepository) GetByUserID(userID uint) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := r.db.Where("user_id = ? AND is_revoked = ?", userID, false).Find(&tokens).Error
	return tokens, err
}

// Delete deletes a refresh token by token ID
func (r *refreshTokenRepository) Delete(tokenID string) error {
	return r.db.Where("token_id = ?", tokenID).Delete(&models.RefreshToken{}).Error
}

// DeleteByUserID deletes all refresh tokens for a user
func (r *refreshTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

// DeleteExpired deletes expired refresh tokens
func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

// userSessionRepository implements UserSessionRepository interface
type userSessionRepository struct {
	db *gorm.DB // Database with dbresolver (handles read/write splitting automatically)
}

// NewUserSessionRepository creates a new user session repository instance
func NewUserSessionRepository(db *gorm.DB) UserSessionRepository {
	return &userSessionRepository{db: db}
}

// Create creates a new user session
func (r *userSessionRepository) Create(session *models.UserSession) error {
	return r.db.Create(session).Error
}

// GetBySessionID retrieves a user session by session ID
func (r *userSessionRepository) GetBySessionID(sessionID string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.Preload("User").
		Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByUserID retrieves all user sessions for a user
func (r *userSessionRepository) GetByUserID(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := r.db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

// Delete deletes a user session by session ID
func (r *userSessionRepository) Delete(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.UserSession{}).Error
}

// DeleteByUserID deletes all user sessions for a user
func (r *userSessionRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

// DeleteExpired deletes expired user sessions
func (r *userSessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.UserSession{}).Error
}
