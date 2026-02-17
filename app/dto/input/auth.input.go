package input

type RegisterEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RegisterPhoneRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

type RegisterGoogleRequest struct {
	GoogleToken string `json:"google_token" validate:"required"`
}

type LoginEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginPhoneRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

type LoginGoogleRequest struct {
	GoogleToken string `json:"google_token" validate:"required"`
}

type LoginAppleRequest struct {
	AppleToken string `json:"apple_token" validate:"required"`
}

type LoginPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type VerifyOTPRequest struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetUserPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username,omitempty"`
	Password string `json:"password" validate:"required,min=8"`
	UserType string `json:"user_type" validate:"required,oneof=admin partner"`
	RoleName string `json:"role_name" validate:"required,oneof=admin partner"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Username *string `json:"username,omitempty"`
	RoleName *string `json:"role_name,omitempty" validate:"omitempty,oneof=admin partner mobile_user"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive pending"`
	Phone    *string `json:"phone,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenValidationRequest struct {
	Token string `json:"token" validate:"required"`
}
