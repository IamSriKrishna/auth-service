package output

import "time"

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int      `json:"expires_in"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID          uint       `json:"id"`
	Email       *string    `json:"email,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Username    *string    `json:"username,omitempty"`
	UserType    string     `json:"user_type"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	VendorID    *uint      `json:"vendor_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type SocialUserData struct {
	UID          string  `json:"uid"`
	Email        *string `json:"email,omitempty"`
	Name         *string `json:"name,omitempty"`
	AuthProvider string  `json:"auth_provider"`
}

type AppleLoginResponse struct {
	Token        string         `json:"token"`
	RefreshToken string         `json:"refreshToken"`
	UserData     SocialUserData `json:"userData"`
	IsNewUser    bool           `json:"isNewUser"`
	UserID       uint           `json:"user_id"`
}

type OTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type TokenValidationResponse struct {
	Valid    bool   `json:"valid"`
	UserID   uint   `json:"user_id,omitempty"`
	UserType string `json:"user_type,omitempty"`
	Role     string `json:"role,omitempty"`
	Claims   Claims `json:"claims,omitempty"`
}

type Claims struct {
	UserID       uint   `json:"user_id"`
	UserType     string `json:"user_type"`
	Role         string `json:"role"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	GoogleID     string `json:"google_id,omitempty"`
	AppleID      string `json:"apple_id,omitempty"`
	FirebaseUID  string `json:"firebase_uid,omitempty"`
	IdentityType string `json:"identity_type"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}

type UserListResponse struct {
	ID          uint       `json:"id"`
	Email       *string    `json:"email,omitempty"`
	Username    *string    `json:"username,omitempty"`
	UserType    string     `json:"user_type"`
	Role        string     `json:"role"`
	Phone       *string    `json:"phone,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uint      `json:"created_by,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type DashboardStatsResponse struct {
	TotalUsers              int                         `json:"total_users"`
	ActiveUsers             int                         `json:"active_users"`
	MembershipUsers         int                         `json:"membership_users"`
	MembershipLastUpdatedAt string                      `json:"membership_last_updated_at,omitempty"`
	Filters                 DashboardStatsFilterApplied `json:"filters_applied"`
}

type DashboardStatsFilterApplied struct {
	CustomerType *string `json:"customer_type,omitempty"`
	FromDate     *string `json:"from_date,omitempty"`
	ToDate       *string `json:"to_date,omitempty"`
}
