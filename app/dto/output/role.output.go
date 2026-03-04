package output

import "time"

type RoleResponse struct {
	ID          uint      `json:"id"`
	RoleName    string    `json:"role_name"`
	Permissions []string  `json:"permissions"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
