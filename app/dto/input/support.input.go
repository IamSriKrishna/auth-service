package input

type CreateSupportRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	IssueType   *string `json:"issue_type,omitempty" validate:"omitempty,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
