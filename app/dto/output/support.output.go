package output

import "time"

// SupportResponse represents the response after creating a support ticket
type SupportResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Phone       *string   `json:"phone,omitempty"`
	Email       *string   `json:"email,omitempty"`
	IssueType   *string   `json:"issue_type,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
