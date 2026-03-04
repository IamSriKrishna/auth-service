package input

type CreateRoleRequest struct {
	RoleName    string   `json:"role_name" validate:"required,min=1,max=255"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
	Description string   `json:"description" validate:"omitempty,max=1000"`
	IsActive    bool     `json:"is_active" validate:""`
}
