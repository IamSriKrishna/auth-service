package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/utils"

	firebase "firebase.google.com/go/v4"
	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// AuthService interface defines authentication service methods
type AuthService interface {
	RegisterEmail(ctx context.Context, req *input.RegisterEmailRequest) (*output.OTPResponse, error)
	RegisterPhone(ctx context.Context, req *input.RegisterPhoneRequest) (*output.OTPResponse, error)
	RegisterGoogle(ctx context.Context, req *input.RegisterGoogleRequest) (*output.AuthResponse, error)
	LoginEmail(ctx context.Context, req *input.LoginEmailRequest) (*output.OTPResponse, error)
	LoginPhone(ctx context.Context, req *input.LoginPhoneRequest) (*output.OTPResponse, error)
	LoginGoogle(ctx context.Context, req *input.LoginGoogleRequest) (*output.AuthResponse, error)
	LoginApple(ctx context.Context, req *input.LoginAppleRequest) (*output.AuthResponse, error)
	LoginPassword(ctx context.Context, req *input.LoginPasswordRequest) (*output.AuthResponse, error)
	RefreshToken(ctx context.Context, req *input.RefreshTokenRequest) (*output.AuthResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *input.ChangePasswordRequest) error
	GetUserInfo(ctx context.Context, userID uint) (*output.UserInfo, error)
	ValidateToken(ctx context.Context, tokenString string) (*output.TokenValidationResponse, error)
	Logout(ctx context.Context, userID uint, tokenID string) error
}

// AdminService interface defines admin service methods
type AdminService interface {
	CreateUser(ctx context.Context, createdBy uint, req *input.CreateUserRequest) (*output.UserInfo, error)
	ResetPassword(ctx context.Context, req *input.ResetPasswordRequest) error
	ResetUserPassword(ctx context.Context, req *input.ResetUserPasswordRequest, userID uint64) error
	GetUsers(ctx context.Context, page, limit int, search string) (*output.PaginatedResponse, error)
	GetUser(ctx context.Context, userID uint) (*output.UserInfo, error)
	UpdateUser(ctx context.Context, userID uint, req *input.UpdateUserRequest) (*output.UserInfo, error)
	DeleteUser(ctx context.Context, userID uint) error
	UpdateUserStatus(ctx context.Context, userID uint, status string) error
	UpdateUserRole(ctx context.Context, userID uint, roleName string) error
	GetDashboardStats(ctx context.Context, filter *input.DashboardStatsFilter) (*output.DashboardStatsResponse, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo         repo.UserRepository
	roleRepo         repo.RoleRepository
	refreshTokenRepo repo.RefreshTokenRepository
	sessionRepo      repo.UserSessionRepository
	oauthConfig      input.OAuthConfig
	firebaseAuth     *fbAuth.Client
}

// NewAuthService creates a new auth service instance
func NewAuthService(
	userRepo repo.UserRepository,
	roleRepo repo.RoleRepository,
	refreshTokenRepo repo.RefreshTokenRepository,
	sessionRepo repo.UserSessionRepository,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
		sessionRepo:      sessionRepo,
	}
}

// getFirebaseAuth lazily initializes and returns Firebase Auth client
func (s *authService) getFirebaseAuth(ctx context.Context) (*fbAuth.Client, error) {
	if s.firebaseAuth != nil {
		return s.firebaseAuth, nil
	}
	// Prefer explicit project ID when available from config
	fbProjectID := s.oauthConfig.FirebaseProjectID
	var fbCfg *firebase.Config
	if fbProjectID != "" {
		fbCfg = &firebase.Config{ProjectID: fbProjectID}
	}

	app, err := firebase.NewApp(ctx, fbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init firebase app: %w", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init firebase auth client: %w", err)
	}
	s.firebaseAuth = client
	return client, nil
}

// RegisterEmail registers a user with email
func (s *authService) RegisterEmail(ctx context.Context, req *input.RegisterEmailRequest) (*output.OTPResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, err
	}

	// Send OTP via email asynchronously using notification service
	go func() {
		if err := s.sendOTPEmail(context.Background(), req.Email, otp); err != nil {
			log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to email successfully",
		ExpiresIn: 300, // 5 minutes
	}, nil
}

// RegisterPhone registers a user with phone
func (s *authService) RegisterPhone(ctx context.Context, req *input.RegisterPhoneRequest) (*output.OTPResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByPhone(req.Phone)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this phone")
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, err
	}

	// Send OTP via SMS asynchronously using notification service
	go func() {
		if err := s.sendOTPSMS(context.Background(), req.Phone, otp); err != nil {
			log.Printf("Failed to send OTP SMS to %s: %v", req.Phone, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to phone successfully",
		ExpiresIn: 300, // 5 minutes
	}, nil
}

// GoogleUserInfo represents user info from Google token
type GoogleUserInfo struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// validateGoogleToken validates Google ID token and extracts user info
func (s *authService) validateGoogleToken(ctx context.Context, tokenString string) (*GoogleUserInfo, error) {
	// Get Google OAuth client IDs from config
	var validClientIDs []string

	if s.oauthConfig.GoogleClientID != "" {
		validClientIDs = append(validClientIDs, s.oauthConfig.GoogleClientID)
	}
	if s.oauthConfig.GoogleIOSClientID != "" {
		validClientIDs = append(validClientIDs, s.oauthConfig.GoogleIOSClientID)
	}
	if s.oauthConfig.GoogleAndroidClientID != "" {
		validClientIDs = append(validClientIDs, s.oauthConfig.GoogleAndroidClientID)
	}

	if len(validClientIDs) == 0 {
		return nil, errors.New("no google oauth client ids configured")
	}

	var payload *idtoken.Payload
	var lastErr error

	// Try validating against each client ID
	for _, clientID := range validClientIDs {
		p, err := idtoken.Validate(ctx, tokenString, clientID)
		if err == nil {
			payload = p
			break
		}
		lastErr = err
	}

	if payload == nil {
		return nil, fmt.Errorf("failed to validate Google token: %v", lastErr)
	}

	// Extract user information
	userInfo := &GoogleUserInfo{
		ID:            payload.Subject,
		Email:         payload.Claims["email"].(string),
		EmailVerified: payload.Claims["email_verified"].(bool),
		Name:          payload.Claims["name"].(string),
		Picture:       payload.Claims["picture"].(string),
	}
	return userInfo, nil
}

// RegisterGoogle registers a user with Google OIDC
func (s *authService) RegisterGoogle(ctx context.Context, req *input.RegisterGoogleRequest) (*output.AuthResponse, error) {
	// Validate Google token and extract user info
	googleUserInfo, err := s.validateGoogleToken(ctx, req.GoogleToken)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByGoogleID(googleUserInfo.ID)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this Google account")
	}

	// Check if user already exists with this email
	existingEmailUser, err := s.userRepo.GetByEmail(googleUserInfo.Email)
	if err == nil && existingEmailUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Get mobile user role
	role, err := s.roleRepo.GetByName("mobile_user")
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Email:         &googleUserInfo.Email,
		GoogleID:      &googleUserInfo.ID,
		UserType:      models.UserTypeMobile,
		RoleID:        role.ID,
		Status:        models.UserStatusActive,
		EmailVerified: googleUserInfo.EmailVerified,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens
	return s.generateTokens(user)
}

// LoginEmail logs in a user with email
func (s *authService) LoginEmail(ctx context.Context, req *input.LoginEmailRequest) (*output.OTPResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, utils.NewInternalServerError("failed to generate OTP")
	}

	// Send OTP via email asynchronously using notification service
	go func() {
		if err := s.sendOTPEmail(context.Background(), req.Email, otp); err != nil {
			log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to email successfully",
		ExpiresIn: 300, // 5 minutes
	}, nil
}

// LoginPhone logs in a user with phone
func (s *authService) LoginPhone(ctx context.Context, req *input.LoginPhoneRequest) (*output.OTPResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByPhone(req.Phone)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, utils.NewInternalServerError("failed to generate OTP")
	}

	// Send OTP via SMS asynchronously using notification service
	go func() {
		if err := s.sendOTPSMS(context.Background(), req.Phone, otp); err != nil {
			log.Printf("Failed to send OTP SMS to %s: %v", req.Phone, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to phone successfully",
		ExpiresIn: 300, // 5 minutes
	}, nil
}

// LoginGoogle logs in a user with Google OIDC
func (s *authService) LoginGoogle(ctx context.Context, req *input.LoginGoogleRequest) (*output.AuthResponse, error) {
	// Validate Google token and extract user info
	googleUserInfo, err := s.validateGoogleToken(ctx, req.GoogleToken)
	if err != nil {
		return nil, utils.NewUnauthorizedError("invalid Google token")
	}

	// Check if user exists
	user, err := s.userRepo.GetByEmail(googleUserInfo.Email)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens
	return s.generateTokens(user)
}

// LoginApple logs in a user with Apple (via Firebase ID token)
func (s *authService) LoginApple(ctx context.Context, req *input.LoginAppleRequest) (*output.AuthResponse, error) {
	// Verify Firebase ID token
	fbClient, err := s.getFirebaseAuth(ctx)
	if err != nil {
		return nil, utils.NewInternalServerError("apple login is not configured: " + err.Error())
	}

	idToken, err := fbClient.VerifyIDToken(ctx, req.AppleToken)
	if err != nil {
		// Log detailed error for troubleshooting (do not expose to clients)
		log.Printf("LoginApple: VerifyIDToken failed: %v", err)
		return nil, utils.NewUnauthorizedError("invalid apple token")
	}

	uid := idToken.UID
	// Optional: ensure the sign-in provider is Apple
	if firebaseObj, ok := idToken.Claims["firebase"].(map[string]interface{}); ok {
		if provider, ok2 := firebaseObj["sign_in_provider"].(string); ok2 && provider != "apple.com" {
			return nil, utils.NewUnauthorizedError("invalid apple token")
		}
	}
	var emailPtr *string
	var namePtr *string
	if v, ok := idToken.Claims["email"].(string); ok && v != "" {
		emailPtr = &v
	}
	if v, ok := idToken.Claims["name"].(string); ok && v != "" {
		namePtr = &v
	}

	if emailPtr == nil {
		// For our current model, we require an email to create/login
		return nil, utils.NewBadRequestError("email not available from apple token")
	}

	// Try by email first
	user, err := s.userRepo.GetByEmail(*emailPtr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new mobile user
			role, rErr := s.roleRepo.GetByName("mobile_user")
			if rErr != nil {
				return nil, rErr
			}
			user = &models.User{
				Email:         emailPtr,
				Username:      namePtr,
				AppleID:       &uid,
				UserType:      models.UserTypeMobile,
				RoleID:        role.ID,
				Status:        models.UserStatusActive,
				EmailVerified: true,
			}
			if cErr := s.userRepo.Create(user); cErr != nil {
				return nil, cErr
			}
		} else {
			return nil, utils.NewInternalServerError("failed to fetch user")
		}
	}

	// If existing user found and AppleID not set, link it now
	if user != nil && user.AppleID == nil {
		user.AppleID = &uid
		if uErr := s.userRepo.Update(user); uErr != nil {
			log.Printf("LoginApple: failed to link AppleID to user %d: %v", user.ID, uErr)
		}
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens (standard auth response)
	return s.generateTokens(user)
}

// LoginPassword logs in a user with password
func (s *authService) LoginPassword(ctx context.Context, req *input.LoginPasswordRequest) (*output.AuthResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, utils.NewUnauthorizedError("invalid credentials")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	// Verify password
	if user.PasswordHash == nil || !utils.CheckPassword(req.Password, *user.PasswordHash) {
		return nil, utils.NewUnauthorizedError("invalid credentials")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens
	return s.generateTokens(user)
}

// RefreshToken refreshes access token
func (s *authService) RefreshToken(ctx context.Context, req *input.RefreshTokenRequest) (*output.AuthResponse, error) {
	// Validate refresh token
	userID, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, errors.New("user account is not active")
	}

	// Generate new tokens
	return s.generateTokens(user)
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, userID uint, req *input.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if user.PasswordHash == nil || !utils.CheckPassword(req.CurrentPassword, *user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = &newPasswordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Update password changed timestamp
	s.userRepo.UpdatePasswordChangedAt(userID)

	return nil
}

// GetUserInfo retrieves user information
func (s *authService) GetUserInfo(ctx context.Context, userID uint) (*output.UserInfo, error) {
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

// ValidateToken validates JWT token
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*output.TokenValidationResponse, error) {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return &output.TokenValidationResponse{Valid: false}, nil
	}

	return &output.TokenValidationResponse{
		Valid:    true,
		UserID:   claims.UserID,
		UserType: claims.UserType,
		Role:     claims.Role,
		Claims:   *claims,
	}, nil
}

// Logout logs out a user
func (s *authService) Logout(ctx context.Context, userID uint, tokenID string) error {
	// Delete refresh token
	s.refreshTokenRepo.Delete(tokenID)

	// Delete user sessions
	s.sessionRepo.DeleteByUserID(userID)

	return nil
}

// generateTokens generates access and refresh tokens
func (s *authService) generateTokens(user *models.User) (*output.AuthResponse, error) {
	// Create JWT claims
	var email, phone, googleID string
	if user.Email != nil {
		email = *user.Email
	}
	if user.Phone != nil {
		phone = *user.Phone
	}
	if user.GoogleID != nil {
		googleID = *user.GoogleID
	}

	claims := output.Claims{
		UserID:       user.ID,
		UserType:     string(user.UserType),
		Role:         user.Role.RoleName,
		IdentityType: utils.GetIdentityType(email, phone, googleID),
	}

	if user.Email != nil {
		claims.Email = *user.Email
	}
	if user.Phone != nil {
		claims.Phone = *user.Phone
	}
	if user.GoogleID != nil {
		claims.GoogleID = *user.GoogleID
	}
	// Attach AppleID if model has it
	if user.AppleID != nil {
		claims.AppleID = *user.AppleID
	}
	// For social providers we treat GoogleID or AppleID (or future provider IDs) as FirebaseUID if available.
	// If Apple login created the user but we didn't store AppleID yet, prefer GoogleID fallback already set.
	if user.AppleID != nil {
		claims.FirebaseUID = *user.AppleID
	} else if user.GoogleID != nil {
		claims.FirebaseUID = *user.GoogleID
	}

	// Generate access token
	accessToken, err := utils.GenerateJWT(claims)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	tokenRecord := &models.RefreshToken{
		TokenID:   uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 90), // 90 days
	}
	s.refreshTokenRepo.Create(tokenRecord)

	return &output.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    7 * 24 * 60 * 60, // 7 days in seconds
		User: output.UserInfo{
			ID:          user.ID,
			Email:       user.Email,
			Phone:       user.Phone,
			Username:    user.Username,
			UserType:    string(user.UserType),
			Role:        user.Role.RoleName,
			Status:      string(user.Status),
			CreatedAt:   user.CreatedAt,
			LastLoginAt: user.LastLoginAt,
		},
	}, nil
}

// EmailRequest represents the notification service email request
type EmailRequest struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	HtmlBody  string `json:"html_body"`
}

// SMSRequest represents the notification service SMS request
type SMSRequest struct {
	ToPhone string `json:"to_phone"`
	OTP     string `json:"otp"`
}

// sendOTPSMS sends OTP via notification service SMS
func (s *authService) sendOTPSMS(ctx context.Context, phone, otp string) error {
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notification-service"
	}

	// Create SMS request
	smsReq := SMSRequest{
		ToPhone: phone,
		OTP:     otp,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(smsReq)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %v", err)
	}

	// Create Fiber client
	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	// Create request
	req := client.Post(notificationServiceURL + "/sms/send")
	req.Body(jsonData)
	req.Set("Content-Type", "application/json")

	// Send request
	status, body, errs := req.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("failed to send SMS request: %v", errs[0])
	}

	if status != 200 {
		return fmt.Errorf("notification service returned status %d: %s", status, string(body))
	}

	return nil
}

// sendOTPEmail sends OTP via notification service
func (s *authService) sendOTPEmail(ctx context.Context, email, otp string) error {
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notification-service"
	}

	// Create email request
	emailReq := EmailRequest{
		ToAddress: email,
		Subject:   "Varthagan OTP Verification",
		HtmlBody: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">Varthagan OTP Verification</h2>
				<p>Your OTP code is:</p>
				<div style="background-color: #f4f4f4; padding: 20px; text-align: center; margin: 20px 0;">
					<h1 style="color: #007bff; font-size: 32px; margin: 0;">%s</h1>
				</div>
				<p>This code will expire in 5 minutes.</p>
				<p>Best regards,<br/>The Varthagan Team</p>
			</div>
		`, otp),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %v", err)
	}

	// Create Fiber client
	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	// Create request
	req := client.Post(notificationServiceURL + "/email/send")
	req.Body(jsonData)
	req.Set("Content-Type", "application/json")

	// Send request
	status, body, errs := req.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("failed to send email request: %v", errs[0])
	}

	if status != 200 {
		return fmt.Errorf("notification service returned status %d: %s", status, string(body))
	}

	return nil
}
