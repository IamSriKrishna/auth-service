package services

import (
	"context"
	"fmt"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type SupportService interface {
	CreateSupport(ctx context.Context, req *input.CreateSupportRequest) (*output.SuccessResponse, error)
}

type supportService struct {
	supportRepo repo.SupportRepository
}

func NewSupportService(supportRepo repo.SupportRepository) SupportService {
	return &supportService{
		supportRepo: supportRepo,
	}
}

func (s *supportService) CreateSupport(ctx context.Context, req *input.CreateSupportRequest) (*output.SuccessResponse, error) {
	support := &models.Support{
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		IssueType:   req.IssueType,
		Description: req.Description,
		Status:      models.SupportStatusOpen,
	}

	if err := s.supportRepo.Create(support); err != nil {
		return nil, fmt.Errorf("failed to create support ticket: %w", err)
	}

	response := &output.SuccessResponse{
		Success: true,
		Message: "Support ticket created successfully",
		Data: output.SupportResponse{
			ID:          support.ID,
			Name:        support.Name,
			Phone:       support.Phone,
			Email:       support.Email,
			IssueType:   support.IssueType,
			Description: support.Description,
			Status:      string(support.Status),
			CreatedAt:   support.CreatedAt,
		},
	}

	return response, nil
}
