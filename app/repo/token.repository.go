package repo

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

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

func (r *refreshTokenRepository) GetByUserID(userID uint) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := r.db.Where("user_id = ? AND is_revoked = ?", userID, false).Find(&tokens).Error
	return tokens, err
}

func (r *refreshTokenRepository) Delete(tokenID string) error {
	return r.db.Where("token_id = ?", tokenID).Delete(&models.RefreshToken{}).Error
}

func (r *refreshTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

type userSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(session *models.UserSession) error {
	return r.db.Create(session).Error
}

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

func (r *userSessionRepository) GetByUserID(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := r.db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

func (r *userSessionRepository) Delete(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.UserSession{}).Error
}

func (r *userSessionRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

func (r *userSessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.UserSession{}).Error
}
