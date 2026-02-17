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

type authService struct {
	userRepo         repo.UserRepository
	roleRepo         repo.RoleRepository
	refreshTokenRepo repo.RefreshTokenRepository
	sessionRepo      repo.UserSessionRepository
	oauthConfig      input.OAuthConfig
	firebaseAuth     *fbAuth.Client
}

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

func (s *authService) getFirebaseAuth(ctx context.Context) (*fbAuth.Client, error) {
	if s.firebaseAuth != nil {
		return s.firebaseAuth, nil
	}
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

func (s *authService) RegisterEmail(ctx context.Context, req *input.RegisterEmailRequest) (*output.OTPResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.sendOTPEmail(context.Background(), req.Email, otp); err != nil {
			log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to email successfully",
		ExpiresIn: 300,
	}, nil
}

func (s *authService) RegisterPhone(ctx context.Context, req *input.RegisterPhoneRequest) (*output.OTPResponse, error) {
	existingUser, err := s.userRepo.GetByPhone(req.Phone)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this phone")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.sendOTPSMS(context.Background(), req.Phone, otp); err != nil {
			log.Printf("Failed to send OTP SMS to %s: %v", req.Phone, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to phone successfully",
		ExpiresIn: 300,
	}, nil
}

type GoogleUserInfo struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func (s *authService) validateGoogleToken(ctx context.Context, tokenString string) (*GoogleUserInfo, error) {
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

	userInfo := &GoogleUserInfo{
		ID:            payload.Subject,
		Email:         payload.Claims["email"].(string),
		EmailVerified: payload.Claims["email_verified"].(bool),
		Name:          payload.Claims["name"].(string),
		Picture:       payload.Claims["picture"].(string),
	}
	return userInfo, nil
}

func (s *authService) RegisterGoogle(ctx context.Context, req *input.RegisterGoogleRequest) (*output.AuthResponse, error) {
	googleUserInfo, err := s.validateGoogleToken(ctx, req.GoogleToken)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByGoogleID(googleUserInfo.ID)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists with this Google account")
	}

	existingEmailUser, err := s.userRepo.GetByEmail(googleUserInfo.Email)
	if err == nil && existingEmailUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	role, err := s.roleRepo.GetByName("mobile_user")
	if err != nil {
		return nil, err
	}

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

	s.userRepo.UpdateLastLogin(user.ID)

	return s.generateTokens(user)
}

func (s *authService) LoginEmail(ctx context.Context, req *input.LoginEmailRequest) (*output.OTPResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, utils.NewInternalServerError("failed to generate OTP")
	}

	go func() {
		if err := s.sendOTPEmail(context.Background(), req.Email, otp); err != nil {
			log.Printf("Failed to send OTP email to %s: %v", req.Email, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to email successfully",
		ExpiresIn: 300,
	}, nil
}

func (s *authService) LoginPhone(ctx context.Context, req *input.LoginPhoneRequest) (*output.OTPResponse, error) {
	user, err := s.userRepo.GetByPhone(req.Phone)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, utils.NewInternalServerError("failed to generate OTP")
	}

	go func() {
		if err := s.sendOTPSMS(context.Background(), req.Phone, otp); err != nil {
			log.Printf("Failed to send OTP SMS to %s: %v", req.Phone, err)
		}
	}()

	return &output.OTPResponse{
		Message:   "OTP sent to phone successfully",
		ExpiresIn: 300,
	}, nil
}

func (s *authService) LoginGoogle(ctx context.Context, req *input.LoginGoogleRequest) (*output.AuthResponse, error) {
	googleUserInfo, err := s.validateGoogleToken(ctx, req.GoogleToken)
	if err != nil {
		return nil, utils.NewUnauthorizedError("invalid Google token")
	}

	user, err := s.userRepo.GetByEmail(googleUserInfo.Email)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	s.userRepo.UpdateLastLogin(user.ID)

	return s.generateTokens(user)
}

func (s *authService) LoginApple(ctx context.Context, req *input.LoginAppleRequest) (*output.AuthResponse, error) {
	fbClient, err := s.getFirebaseAuth(ctx)
	if err != nil {
		return nil, utils.NewInternalServerError("apple login is not configured: " + err.Error())
	}

	idToken, err := fbClient.VerifyIDToken(ctx, req.AppleToken)
	if err != nil {
		log.Printf("LoginApple: VerifyIDToken failed: %v", err)
		return nil, utils.NewUnauthorizedError("invalid apple token")
	}

	uid := idToken.UID
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
		return nil, utils.NewBadRequestError("email not available from apple token")
	}

	user, err := s.userRepo.GetByEmail(*emailPtr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	if user != nil && user.AppleID == nil {
		user.AppleID = &uid
		if uErr := s.userRepo.Update(user); uErr != nil {
			log.Printf("LoginApple: failed to link AppleID to user %d: %v", user.ID, uErr)
		}
	}

	s.userRepo.UpdateLastLogin(user.ID)

	return s.generateTokens(user)
}

func (s *authService) LoginPassword(ctx context.Context, req *input.LoginPasswordRequest) (*output.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, utils.NewUnauthorizedError("invalid credentials")
	}

	if user.Status != models.UserStatusActive {
		return nil, utils.NewForbiddenError("user account is not active")
	}

	if user.PasswordHash == nil || !utils.CheckPassword(req.Password, *user.PasswordHash) {
		return nil, utils.NewUnauthorizedError("invalid credentials")
	}

	s.userRepo.UpdateLastLogin(user.ID)

	return s.generateTokens(user)
}

func (s *authService) RefreshToken(ctx context.Context, req *input.RefreshTokenRequest) (*output.AuthResponse, error) {
	userID, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != models.UserStatusActive {
		return nil, errors.New("user account is not active")
	}

	return s.generateTokens(user)
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, req *input.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.PasswordHash == nil || !utils.CheckPassword(req.CurrentPassword, *user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = &newPasswordHash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.userRepo.UpdatePasswordChangedAt(userID)

	return nil
}

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

func (s *authService) Logout(ctx context.Context, userID uint, tokenID string) error {
	s.refreshTokenRepo.Delete(tokenID)

	s.sessionRepo.DeleteByUserID(userID)

	return nil
}

func (s *authService) generateTokens(user *models.User) (*output.AuthResponse, error) {
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
	if user.AppleID != nil {
		claims.AppleID = *user.AppleID
	}
	if user.AppleID != nil {
		claims.FirebaseUID = *user.AppleID
	} else if user.GoogleID != nil {
		claims.FirebaseUID = *user.GoogleID
	}

	accessToken, err := utils.GenerateJWT(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	tokenRecord := &models.RefreshToken{
		TokenID:   uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 90),
	}
	s.refreshTokenRepo.Create(tokenRecord)

	return &output.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    7 * 24 * 60 * 60,
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

type EmailRequest struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	HtmlBody  string `json:"html_body"`
}

type SMSRequest struct {
	ToPhone string `json:"to_phone"`
	OTP     string `json:"otp"`
}

func (s *authService) sendOTPSMS(ctx context.Context, phone, otp string) error {
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notification-service"
	}

	smsReq := SMSRequest{
		ToPhone: phone,
		OTP:     otp,
	}

	jsonData, err := json.Marshal(smsReq)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %v", err)
	}

	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	req := client.Post(notificationServiceURL + "/sms/send")
	req.Body(jsonData)
	req.Set("Content-Type", "application/json")

	status, body, errs := req.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("failed to send SMS request: %v", errs[0])
	}

	if status != 200 {
		return fmt.Errorf("notification service returned status %d: %s", status, string(body))
	}

	return nil
}

func (s *authService) sendOTPEmail(ctx context.Context, email, otp string) error {
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notification-service"
	}

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

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %v", err)
	}

	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	req := client.Post(notificationServiceURL + "/email/send")
	req.Body(jsonData)
	req.Set("Content-Type", "application/json")

	status, body, errs := req.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("failed to send email request: %v", errs[0])
	}

	if status != 200 {
		return fmt.Errorf("notification service returned status %d: %s", status, string(body))
	}

	return nil
}
