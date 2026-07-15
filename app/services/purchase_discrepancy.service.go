package services

import (
	"errors"
	"fmt"
	"math"
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
	poRepository    repo.PurchaseOrderRepository
	bagRepository   repo.RawMaterialBagRepository
}

func NewPurchaseDispenseService(
	repository repo.PurchaseDispenseRepository,
	claimRepository repo.PurchaseClaimRepository,
	poRepository repo.PurchaseOrderRepository,
	bagRepository repo.RawMaterialBagRepository,
) PurchaseDispenseService {
	return &purchaseDispenseService{
		repository:      repository,
		claimRepository: claimRepository,
		poRepository:    poRepository,
		bagRepository:   bagRepository,
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

	if strings.TrimSpace(claimID) == "" {
		return nil, errors.New("purchase claim ID is required")
	}

	if companyID == 0 {
		return nil, errors.New("company ID is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("authenticated user is required")
	}

	if request.PurchaseClaimItemID == 0 {
		return nil, errors.New("purchase claim item ID is required")
	}

	if request.Quantity <= 0 {
		return nil, errors.New("dispense quantity must be greater than zero")
	}

	if strings.TrimSpace(request.Unit) == "" {
		return nil, errors.New("dispense unit is required")
	}

	claimItem, claim, err :=
		s.claimRepository.FindItemByIDAndCompany(
			request.PurchaseClaimItemID,
			companyID,
		)
	if err != nil {
		return nil, errors.New("purchase claim item not found")
	}

	if claim == nil {
		return nil, errors.New("purchase claim not found")
	}

	if claimItem == nil {
		return nil, errors.New("purchase claim item not found")
	}

	if claim.ID != claimID {
		return nil, errors.New(
			"claim item does not belong to the selected purchase claim",
		)
	}

	if claimItem.Action != models.PurchaseClaimActionReplacement {
		return nil, errors.New(
			"purchase claim item action must be replacement",
		)
	}

	if claimItem.ReplacementCompleted {
		return nil, errors.New(
			"vendor replacement is already completed for this claim item",
		)
	}

	if claimItem.ReplacementPendingBase <= 0 {
		return nil, errors.New(
			"there is no pending replacement quantity for this claim item",
		)
	}

	purchaseOrder, err :=
		s.poRepository.FindByIDAndCompany(
			claim.PurchaseOrderID,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"purchase order not found in your company",
		)
	}

	purchaseOrderItem, err :=
		findDispensePurchaseOrderLineItem(
			purchaseOrder,
			claimItem.PurchaseOrderItemID,
		)
	if err != nil {
		return nil, err
	}

	if purchaseOrderItem.ProductID == nil ||
		strings.TrimSpace(*purchaseOrderItem.ProductID) == "" {
		return nil, errors.New(
			"purchase order line item has no product",
		)
	}

	if *purchaseOrderItem.ProductID != claimItem.ProductID {
		return nil, errors.New(
			"purchase order product does not match the claim item product",
		)
	}

	baseQuantity, baseUnit, err :=
		normalizeDispenseQuantity(
			purchaseOrderItem,
			request.Quantity,
			request.Unit,
		)
	if err != nil {
		return nil, err
	}

	if baseQuantity >
		claimItem.ReplacementPendingBase+0.000001 {
		return nil, fmt.Errorf(
			"dispense quantity %.3f %s exceeds pending replacement %.3f %s",
			baseQuantity,
			baseUnit,
			claimItem.ReplacementPendingBase,
			claimItem.BaseUnit,
		)
	}

	now := time.Now()

	dispenseDate := request.DispenseDate
	if dispenseDate.IsZero() {
		dispenseDate = now
	}

	dispense := &models.PurchaseDispense{
		ID: uuid.New().String(),

		PurchaseClaimID:     claim.ID,
		PurchaseClaimItemID: claimItem.ID,
		PurchaseOrderID:     claim.PurchaseOrderID,

		ProductID:     claimItem.ProductID,
		ProductName:   claimItem.ProductName,
		IsRawMaterial: claimItem.IsRawMaterial,

		Quantity: request.Quantity,
		Unit: strings.ToLower(
			strings.TrimSpace(request.Unit),
		),

		BaseQuantity: baseQuantity,
		BaseUnit:     baseUnit,

		DispenseDate: dispenseDate,
		Notes:        strings.TrimSpace(request.Notes),

		CreatedBy: userID,
		CreatedAt: now,
	}

	err = s.repository.GetDB().
		Transaction(func(tx *gorm.DB) error {
			/*
				1. Add the vendor replacement to the correct stock.

				Raw material:
					product_stocks
					stock_ledgers

				Finished product:
					variant_stocks
					variant_stock_movements
			*/
			if err := addReplacementToCorrectStock(
				tx,
				purchaseOrderItem,
				baseQuantity,
				claimItem.Rate,
				claim.ID,
				claim.ClaimNumber,
				request.Notes,
				userID,
			); err != nil {
				return fmt.Errorf(
					"failed to add replacement into stock: %w",
					err,
				)
			}

			/*
				2. For a raw material, automatically allocate the
				   replacement quantity to bags with shortages.

				Order:
					bag_number ASC

				Shortage:
					expected_kg - actual_kg
			*/
			if claimItem.IsRawMaterial {
				replacementKg := baseQuantity / 1000

				if err :=
					s.allocateReplacementToLowBagsTx(
						tx,
						claim.PurchaseOrderID,
						claimItem.ProductID,
						replacementKg,
						companyID,
					); err != nil {
					return fmt.Errorf(
						"failed to allocate replacement to raw-material bags: %w",
						err,
					)
				}
			}

			/*
				3. Update purchase claim replacement progress.
			*/
			claimItem.ReplacementReceivedBase +=
				baseQuantity

			claimItem.ReplacementPendingBase -=
				baseQuantity

			if claimItem.ReplacementPendingBase <= 0.000001 {
				claimItem.ReplacementPendingBase = 0
				claimItem.ReplacementCompleted = true
				claimItem.ReplacementCompletedAt = &now
			}

			claimItem.UpdatedAt = now

			if err := tx.Save(claimItem).Error; err != nil {
				return fmt.Errorf(
					"failed to update purchase claim item: %w",
					err,
				)
			}

			/*
				4. Create the purchase dispense history row.
			*/
			if err := s.repository.CreateTx(
				tx,
				dispense,
			); err != nil {
				return fmt.Errorf(
					"failed to create purchase dispense: %w",
					err,
				)
			}

			/*
				5. Create the vendor replacement receipt.
			*/
			receipt := &models.PurchaseClaimReceipt{
				ID: uuid.New().String(),

				PurchaseClaimID:     claim.ID,
				PurchaseClaimItemID: claimItem.ID,

				ProductID: claimItem.ProductID,

				ReceivedQuantity: request.Quantity,
				ReceivedUnit: strings.ToLower(
					strings.TrimSpace(request.Unit),
				),

				ReceivedBaseQuantity: baseQuantity,
				BaseUnit:             baseUnit,

				ReceivedDate: dispenseDate,
				Notes:        strings.TrimSpace(request.Notes),

				ReceivedBy: userID,
				CreatedAt:  now,
			}

			if err :=
				s.claimRepository.CreateReceiptTx(
					tx,
					receipt,
				); err != nil {
				return fmt.Errorf(
					"failed to create purchase claim receipt: %w",
					err,
				)
			}

			/*
				6. Update the parent purchase claim status.

				Any pending replacement:
					partial

				No pending replacement:
					resolved
			*/
			var pendingReplacementCount int64

			if err := tx.
				Model(&models.PurchaseClaimItem{}).
				Where(
					"purchase_claim_id = ? AND action = ? AND replacement_pending_base > ?",
					claim.ID,
					models.PurchaseClaimActionReplacement,
					0.000001,
				).
				Count(&pendingReplacementCount).
				Error; err != nil {
				return fmt.Errorf(
					"failed to calculate pending replacements: %w",
					err,
				)
			}

			newClaimStatus :=
				models.PurchaseClaimStatusResolved

			if pendingReplacementCount > 0 {
				newClaimStatus =
					models.PurchaseClaimStatusPartial
			}

			if err :=
				s.claimRepository.UpdateClaimStatusTx(
					tx,
					claim.ID,
					newClaimStatus,
				); err != nil {
				return fmt.Errorf(
					"failed to update purchase claim status: %w",
					err,
				)
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return s.repository.FindByIDAndCompany(
		dispense.ID,
		companyID,
	)
}

func (s *purchaseDispenseService) allocateReplacementToLowBagsTx(
	tx *gorm.DB,
	purchaseOrderID string,
	productID string,
	replacementKg float64,
	companyID uint,
) error {
	if tx == nil {
		return errors.New(
			"database transaction is required",
		)
	}

	if replacementKg <= 0 {
		return errors.New(
			"raw-material replacement quantity must be greater than zero",
		)
	}

	bags, err :=
		s.bagRepository.GetShortageBagsForUpdateTx(
			tx,
			purchaseOrderID,
			productID,
			companyID,
		)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve raw-material shortage bags: %w",
			err,
		)
	}

	if len(bags) == 0 {
		return errors.New(
			"no raw-material bags with missing quantity were found",
		)
	}

	totalShortageKg := 0.0

	for index := range bags {
		shortageKg :=
			bags[index].ExpectedKg -
				bags[index].ActualKg

		if shortageKg > 0 {
			totalShortageKg += shortageKg
		}
	}

	if replacementKg >
		totalShortageKg+0.000001 {
		return fmt.Errorf(
			"replacement quantity %.3f kg exceeds total bag shortage %.3f kg",
			replacementKg,
			totalShortageKg,
		)
	}

	remainingReplacementKg := replacementKg
	now := time.Now()

	for index := range bags {
		if remainingReplacementKg <= 0.000001 {
			break
		}

		bag := &bags[index]

		bagShortageKg :=
			bag.ExpectedKg -
				bag.ActualKg

		if bagShortageKg <= 0 {
			continue
		}

		allocationKg :=
			math.Min(
				bagShortageKg,
				remainingReplacementKg,
			)

		/*
			Preserve material already consumed from this bag.

			Example:

				ActualKg:    10
				RemainingKg: 4

				ConsumedKg:
					10 - 4 = 6
		*/
		consumedKg :=
			bag.ActualKg -
				bag.RemainingKg

		if consumedKg < 0 {
			consumedKg = 0
		}

		/*
			Add the vendor replacement into the original bag.
		*/
		bag.ActualKg += allocationKg

		/*
			Remaining stock should include the replacement,
			but must still account for previously consumed stock.
		*/
		bag.RemainingKg =
			bag.ActualKg -
				consumedKg

		if bag.RemainingKg < 0 {
			bag.RemainingKg = 0
		}

		switch {
		case bag.RemainingKg <= 0.000001:
			bag.RemainingKg = 0
			bag.Status = "used"

		case consumedKg > 0:
			bag.Status = "partial"

		default:
			bag.Status = "available"
		}

		bag.UpdatedAt = now

		if err :=
			s.bagRepository.UpdateReplacementBagTx(
				tx,
				bag,
			); err != nil {
			return fmt.Errorf(
				"failed to update raw-material bag %d: %w",
				bag.BagNumber,
				err,
			)
		}

		remainingReplacementKg -=
			allocationKg
	}

	if remainingReplacementKg >
		0.000001 {
		return fmt.Errorf(
			"unable to allocate remaining replacement quantity %.3f kg",
			remainingReplacementKg,
		)
	}

	return nil
}

func normalizeDispenseQuantity(
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	unit string,
) (float64, string, error) {
	if lineItem == nil {
		return 0, "", errors.New(
			"purchase order line item is required",
		)
	}

	if quantity <= 0 {
		return 0, "", errors.New(
			"dispense quantity must be greater than zero",
		)
	}

	normalizedUnit :=
		strings.ToLower(
			strings.TrimSpace(unit),
		)

	if lineItem.IsRawMaterial {
		return normalizeRawMaterialQuantity(
			quantity,
			normalizedUnit,
		)
	}

	expectedUnit :=
		strings.ToLower(
			strings.TrimSpace(
				lineItem.StockUnit,
			),
		)

	if expectedUnit == "" {
		expectedUnit =
			strings.ToLower(
				strings.TrimSpace(
					lineItem.PurchaseUnit,
				),
			)
	}

	if expectedUnit == "" {
		expectedUnit = "pieces"
	}

	if normalizedUnit != expectedUnit {
		return 0, "", fmt.Errorf(
			"unit must be %s for product %s",
			expectedUnit,
			lineItem.ProductName,
		)
	}

	return quantity, expectedUnit, nil
}

func findDispensePurchaseOrderLineItem(
	purchaseOrder *models.PurchaseOrder,
	lineItemID uint,
) (*models.PurchaseOrderLineItem, error) {
	if purchaseOrder == nil {
		return nil, errors.New(
			"purchase order is required",
		)
	}

	for index := range purchaseOrder.LineItems {
		if purchaseOrder.LineItems[index].ID ==
			lineItemID {
			return &purchaseOrder.LineItems[index], nil
		}
	}

	return nil, fmt.Errorf(
		"purchase order line item %d not found",
		lineItemID,
	)
}

func (s *purchaseDispenseService) GetDispenseByID(
	id string,
	companyID uint,
) (*models.PurchaseDispense, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New(
			"purchase dispense ID is required",
		)
	}

	return s.repository.FindByIDAndCompany(
		id,
		companyID,
	)
}

func (s *purchaseDispenseService) GetDispensesByClaim(
	claimID string,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	if strings.TrimSpace(claimID) == "" {
		return nil, errors.New(
			"purchase claim ID is required",
		)
	}

	return s.repository.FindByClaimAndCompany(
		claimID,
		companyID,
	)
}

func (s *purchaseDispenseService) GetDispensesByClaimItem(
	claimItemID uint,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	if claimItemID == 0 {
		return nil, errors.New(
			"purchase claim item ID is required",
		)
	}

	return s.repository.FindByClaimItemAndCompany(
		claimItemID,
		companyID,
	)
}
