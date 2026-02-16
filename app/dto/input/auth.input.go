package input

// RegisterEmailRequest represents email registration request
type RegisterEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// RegisterPhoneRequest represents phone registration request
type RegisterPhoneRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

// RegisterGoogleRequest represents Google OIDC registration request
type RegisterGoogleRequest struct {
	GoogleToken string `json:"google_token" validate:"required"`
}

// LoginEmailRequest represents email login request
type LoginEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// LoginPhoneRequest represents phone login request
type LoginPhoneRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

// LoginGoogleRequest represents Google OIDC login request
type LoginGoogleRequest struct {
	GoogleToken string `json:"google_token" validate:"required"`
}

// LoginAppleRequest represents Apple (via Firebase) login request
type LoginAppleRequest struct {
	AppleToken string `json:"apple_token" validate:"required"`
}

// LoginPasswordRequest represents password-based login request
type LoginPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ResetPasswordRequest represents password reset request by super admin
type ResetPasswordRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetUserPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// CreateUserRequest represents user creation request by super admin
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username,omitempty"`
	Password string `json:"password" validate:"required,min=8"`
	UserType string `json:"user_type" validate:"required,oneof=admin partner"`
	RoleName string `json:"role_name" validate:"required,oneof=admin partner"`
	Phone    string `json:"phone,omitempty"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Username *string `json:"username,omitempty"`
	RoleName *string `json:"role_name,omitempty" validate:"omitempty,oneof=admin partner mobile_user"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive pending"`
	Phone    *string `json:"phone,omitempty"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenValidationRequest represents token validation request
type TokenValidationRequest struct {
	Token string `json:"token" validate:"required"`
}
