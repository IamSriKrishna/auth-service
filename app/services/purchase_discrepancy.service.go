package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseDispenseService interface {
	CreateDispense(
		claimID string,
		request *input.CreatePurchaseDispenseInput,
		userID string,
		companyID uint,
	) (*models.PurchaseDispense, error)

	GetDispenseByID(
		id string,
		companyID uint,
	) (*models.PurchaseDispense, error)

	GetDispensesByClaim(
		claimID string,
		companyID uint,
	) ([]models.PurchaseDispense, error)

	GetDispensesByClaimItem(
		claimItemID uint,
		companyID uint,
	) ([]models.PurchaseDispense, error)
}

type purchaseDispenseService struct {
	repository      repo.PurchaseDispenseRepository
	claimRepository repo.PurchaseClaimRepository
	claimService    PurchaseClaimService
}

func NewPurchaseDispenseService(
	repository repo.PurchaseDispenseRepository,
	claimRepository repo.PurchaseClaimRepository,
	claimService PurchaseClaimService,
) PurchaseDispenseService {
	return &purchaseDispenseService{
		repository:      repository,
		claimRepository: claimRepository,
		claimService:    claimService,
	}
}

func (s *purchaseDispenseService) CreateDispense(
	claimID string,
	request *input.CreatePurchaseDispenseInput,
	userID string,
	companyID uint,
) (*models.PurchaseDispense, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	item, claim, err := s.claimRepository.FindItemByIDAndCompany(
		request.PurchaseClaimItemID,
		companyID,
	)
	if err != nil {
		return nil, errors.New("purchase claim item not found")
	}
	if claim.ID != claimID {
		return nil, errors.New("claim item does not belong to this claim")
	}
	if item.Action != models.PurchaseClaimActionReplacement {
		return nil, errors.New("claim action must be replacement")
	}

	// Reuse the stock-safe replacement logic from PurchaseClaimService.
	updatedClaim, err := s.claimService.ReceiveReplacement(
		claimID,
		&input.ReceivePurchaseClaimReplacementInput{
			PurchaseClaimItemID: request.PurchaseClaimItemID,
			Quantity:            request.Quantity,
			Unit:                request.Unit,
			ReceivedDate:        request.DispenseDate,
			Notes:               request.Notes,
		},
		userID,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	var updatedItem *models.PurchaseClaimItem
	for index := range updatedClaim.Items {
		if updatedClaim.Items[index].ID == request.PurchaseClaimItemID {
			updatedItem = &updatedClaim.Items[index]
			break
		}
	}
	if updatedItem == nil {
		return nil, errors.New("updated claim item not found")
	}

	baseQuantity := request.Quantity
	baseUnit := strings.ToLower(strings.TrimSpace(request.Unit))
	if item.IsRawMaterial {
		baseQuantity, baseUnit, err = normalizeRawMaterialQuantity(
			request.Quantity,
			request.Unit,
		)
		if err != nil {
			return nil, err
		}
	}

	value := &models.PurchaseDispense{
		ID:                  uuid.New().String(),
		PurchaseClaimID:     claim.ID,
		PurchaseClaimItemID: item.ID,
		PurchaseOrderID:     claim.PurchaseOrderID,
		ProductID:           item.ProductID,
		ProductName:         item.ProductName,
		IsRawMaterial:       item.IsRawMaterial,
		Quantity:            request.Quantity,
		Unit:                strings.ToLower(strings.TrimSpace(request.Unit)),
		BaseQuantity:        baseQuantity,
		BaseUnit:            baseUnit,
		DispenseDate:        request.DispenseDate,
		Notes:               request.Notes,
		CreatedBy:           userID,
		CreatedAt:           time.Now(),
	}

	err = s.repository.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := s.repository.CreateTx(tx, value); err != nil {
			return fmt.Errorf("failed to create purchase dispense: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *purchaseDispenseService) GetDispenseByID(
	id string,
	companyID uint,
) (*models.PurchaseDispense, error) {
	return s.repository.FindByIDAndCompany(id, companyID)
}

func (s *purchaseDispenseService) GetDispensesByClaim(
	claimID string,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	return s.repository.FindByClaimAndCompany(claimID, companyID)
}

func (s *purchaseDispenseService) GetDispensesByClaimItem(
	claimItemID uint,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	return s.repository.FindByClaimItemAndCompany(claimItemID, companyID)
}
