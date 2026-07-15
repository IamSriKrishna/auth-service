package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseClaimService interface {
	GetPurchaseOrderClaimSource(
		purchaseOrderID string,
		companyID uint,
	) (*output.PurchaseOrderClaimSourceOutput, error)

	CreateClaim(
		request *input.CreatePurchaseClaimInput,
		userID string,
		companyID uint,
	) (*output.PurchaseClaimOutput, error)

	GetClaimByID(
		id string,
		companyID uint,
	) (*output.PurchaseClaimOutput, error)

	GetClaimsByPurchaseOrder(
		purchaseOrderID string,
		companyID uint,
	) ([]output.PurchaseClaimOutput, error)

	ReceiveReplacement(
		claimID string,
		request *input.ReceivePurchaseClaimReplacementInput,
		userID string,
		companyID uint,
	) (*output.PurchaseClaimOutput, error)

	GetReplacementReceipts(
		claimItemID uint,
		companyID uint,
	) ([]models.PurchaseClaimReceipt, error)

	GetNetReceivableBaseQuantity(
		purchaseOrderItemID uint,
		orderedBaseQuantity float64,
	) (float64, error)
}

type purchaseClaimService struct {
	repository  repo.PurchaseClaimRepository
	poRepo      repo.PurchaseOrderRepository
	productRepo repo.ProductRepository
}

func NewPurchaseClaimService(
	repository repo.PurchaseClaimRepository,
	poRepo repo.PurchaseOrderRepository,
	productRepo repo.ProductRepository,
) PurchaseClaimService {
	return &purchaseClaimService{
		repository:  repository,
		poRepo:      poRepo,
		productRepo: productRepo,
	}
}

func normalizeRawMaterialQuantity(quantity float64, unit string) (float64, string, error) {
	switch strings.ToLower(strings.TrimSpace(unit)) {
	case "kg", "kilogram", "kilograms":
		return quantity * 1000, "gram", nil
	case "g", "gram", "grams":
		return quantity, "gram", nil
	case "mg", "milligram", "milligrams":
		return quantity / 1000, "gram", nil
	case "ton", "tons", "tonne", "tonnes":
		return quantity * 1_000_000, "gram", nil
	default:
		return 0, "", fmt.Errorf("unsupported raw-material unit: %s", unit)
	}
}

func normalizeClaimQuantity(
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	unit string,
) (float64, string, error) {
	if lineItem.IsRawMaterial {
		return normalizeRawMaterialQuantity(quantity, unit)
	}

	inputUnit := strings.ToLower(strings.TrimSpace(unit))
	baseUnit := strings.ToLower(strings.TrimSpace(lineItem.StockUnit))
	if baseUnit == "" {
		baseUnit = strings.ToLower(strings.TrimSpace(lineItem.PurchaseUnit))
	}
	if baseUnit == "" {
		baseUnit = "pieces"
	}

	if inputUnit != baseUnit {
		return 0, "", fmt.Errorf(
			"unit must be %s for finished product %s",
			baseUnit,
			lineItem.ProductName,
		)
	}
	return quantity, baseUnit, nil
}

func getOrderedBaseQuantity(
	lineItem *models.PurchaseOrderLineItem,
) (float64, string, error) {
	if lineItem.StockQuantity > 0 && lineItem.StockUnit != "" {
		return lineItem.StockQuantity, strings.ToLower(strings.TrimSpace(lineItem.StockUnit)), nil
	}
	return normalizeClaimQuantity(lineItem, lineItem.Quantity, lineItem.PurchaseUnit)
}

func getReceivedBaseQuantity(
	lineItem *models.PurchaseOrderLineItem,
	po *models.PurchaseOrder,
) (float64, error) {
	if lineItem.ReceivedQuantity > 0 {
		value, _, err := normalizeClaimQuantity(
			lineItem,
			lineItem.ReceivedQuantity,
			lineItem.PurchaseUnit,
		)
		return value, err
	}

	if po.InventorySynced {
		value, _, err := getOrderedBaseQuantity(lineItem)
		return value, err
	}

	return 0, nil
}

func findPurchaseOrderLineItem(
	po *models.PurchaseOrder,
	itemID uint,
) (*models.PurchaseOrderLineItem, error) {
	for index := range po.LineItems {
		if po.LineItems[index].ID == itemID {
			return &po.LineItems[index], nil
		}
	}
	return nil, fmt.Errorf("purchase order line item %d not found", itemID)
}

func getVendorName(po *models.PurchaseOrder) string {
	if po.Vendor == nil {
		return ""
	}
	if po.Vendor.DisplayName != "" {
		return po.Vendor.DisplayName
	}
	return po.Vendor.CompanyName
}

func (s *purchaseClaimService) GetPurchaseOrderClaimSource(
	purchaseOrderID string,
	companyID uint,
) (*output.PurchaseOrderClaimSourceOutput, error) {
	po, err := s.poRepo.FindByIDAndCompany(purchaseOrderID, companyID)
	if err != nil {
		return nil, errors.New("purchase order not found in your company")
	}

	result := &output.PurchaseOrderClaimSourceOutput{
		PurchaseOrderID:     po.ID,
		PurchaseOrderNumber: po.PurchaseOrderNumber,
		VendorID:            po.VendorID,
		VendorName:          getVendorName(po),
		Status:              string(po.Status),
		InventorySynced:     po.InventorySynced,
		Items:               make([]output.PurchaseOrderClaimSourceItemOutput, 0),
	}

	for index := range po.LineItems {
		item := &po.LineItems[index]
		if item.ProductID == nil || *item.ProductID == "" {
			continue
		}

		orderedBase, baseUnit, err := getOrderedBaseQuantity(item)
		if err != nil {
			return nil, err
		}
		receivedBase, err := getReceivedBaseQuantity(item, po)
		if err != nil {
			return nil, err
		}

		missing, err := s.repository.SumClaimedByPOItem(
			nil,
			item.ID,
			models.PurchaseClaimMissing,
		)
		if err != nil {
			return nil, err
		}

		damaged, err := s.repository.SumClaimedByPOItem(
			nil,
			item.ID,
			models.PurchaseClaimDamaged,
		)
		if err != nil {
			return nil, err
		}

		replacementPending, err := s.repository.SumReplacementPendingByPOItem(nil, item.ID)
		if err != nil {
			return nil, err
		}

		result.Items = append(result.Items, output.PurchaseOrderClaimSourceItemOutput{
			PurchaseOrderItemID: item.ID,
			ProductID:           *item.ProductID,
			ProductName:         item.ProductName,
			SKU:                 item.SKU,
			IsRawMaterial:       item.IsRawMaterial,

			OrderedQuantity:     item.Quantity,
			OrderedUnit:         item.PurchaseUnit,
			OrderedBaseQuantity: orderedBase,
			BaseUnit:            baseUnit,

			ReceivedQuantity:     item.ReceivedQuantity,
			ReceivedBaseQuantity: receivedBase,

			MissingReportedBase: missing,
			DamagedReportedBase: damaged,

			MissingRemainingBase: math.Max(orderedBase-missing-damaged, 0),
			DamagedRemainingBase: math.Max(receivedBase-missing-damaged, 0),

			ReplacementPendingBase: replacementPending,

			NumberOfPacks:   item.NumberOfPacks,
			QuantityPerPack: item.QuantityPerPack,
			ReceivedPacks:   item.ReceivedPacks,
			Rate:            item.Rate,
		})
	}

	return result, nil
}

func (s *purchaseClaimService) CreateClaim(
	request *input.CreatePurchaseClaimInput,
	userID string,
	companyID uint,
) (*output.PurchaseClaimOutput, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	po, err := s.poRepo.FindByIDAndCompany(request.PurchaseOrderID, companyID)
	if err != nil {
		return nil, errors.New("purchase order not found in your company")
	}

	now := time.Now()
	claim := &models.PurchaseClaim{
		ID:                  uuid.New().String(),
		ClaimNumber:         fmt.Sprintf("PC-%s-%d", now.Format("20060102"), now.UnixNano()),
		PurchaseOrderID:     po.ID,
		PurchaseOrderNumber: po.PurchaseOrderNumber,
		VendorID:            po.VendorID,
		CompanyID:           companyID,
		ClaimDate:           request.Date,
		Status:              models.PurchaseClaimStatusOpen,
		Notes:               request.Notes,
		CreatedBy:           userID,
		CreatedAt:           now,
		UpdatedAt:           now,
		Items:               make([]models.PurchaseClaimItem, 0, len(request.Items)),
	}

	err = s.repository.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, itemInput := range request.Items {
			lineItem, err := findPurchaseOrderLineItem(po, itemInput.PurchaseOrderItemID)
			if err != nil {
				return err
			}
			if lineItem.ProductID == nil || *lineItem.ProductID == "" {
				return errors.New("purchase order line item has no product")
			}

			if _, err := s.productRepo.FindByIDAndCompany(*lineItem.ProductID, companyID); err != nil {
				return fmt.Errorf("product %s does not belong to your company", *lineItem.ProductID)
			}

			baseQuantity, baseUnit, err := normalizeClaimQuantity(
				lineItem,
				itemInput.Quantity,
				itemInput.Unit,
			)
			if err != nil {
				return err
			}

			orderedBase, _, err := getOrderedBaseQuantity(lineItem)
			if err != nil {
				return err
			}
			receivedBase, err := getReceivedBaseQuantity(lineItem, po)
			if err != nil {
				return err
			}

			existingMissing, err := s.repository.SumClaimedByPOItem(
				tx,
				lineItem.ID,
				models.PurchaseClaimMissing,
			)
			if err != nil {
				return err
			}
			existingDamaged, err := s.repository.SumClaimedByPOItem(
				tx,
				lineItem.ID,
				models.PurchaseClaimDamaged,
			)
			if err != nil {
				return err
			}

			switch itemInput.Type {
			case models.PurchaseClaimMissing:
				maxMissing := orderedBase - existingMissing - existingDamaged
				if baseQuantity > maxMissing+0.000001 {
					return fmt.Errorf(
						"missing quantity %.3f %s exceeds remaining PO quantity %.3f %s for %s",
						baseQuantity,
						baseUnit,
						math.Max(maxMissing, 0),
						baseUnit,
						lineItem.ProductName,
					)
				}

			case models.PurchaseClaimDamaged:
				if !po.InventorySynced {
					return fmt.Errorf(
						"purchase order %s is not received; receive it before reporting damaged stock",
						po.PurchaseOrderNumber,
					)
				}
				maxDamage := receivedBase - existingMissing - existingDamaged
				if baseQuantity > maxDamage+0.000001 {
					return fmt.Errorf(
						"damaged quantity %.3f %s exceeds remaining received quantity %.3f %s for %s",
						baseQuantity,
						baseUnit,
						math.Max(maxDamage, 0),
						baseUnit,
						lineItem.ProductName,
					)
				}

			default:
				return fmt.Errorf("unsupported claim type: %s", itemInput.Type)
			}

			stockAdjusted := false

			// Missing quantity:
			// - If PO is not yet synced, do not add missing quantity during PO receive.
			// - If PO is already synced, remove missing quantity from stock.
			//
			// Damaged quantity:
			// - PO must already be synced.
			// - Move quantity from available stock to damaged stock.
			if po.InventorySynced {
				if err := applyPurchaseClaimToCorrectStock(
					tx,
					lineItem,
					baseQuantity,
					itemInput.Type,
					itemInput.Reason,
					claim.ID,
					claim.ClaimNumber,
					userID,
				); err != nil {
					return err
				}
				stockAdjusted = true
			}

			replacementPending := 0.0
			if itemInput.Action == models.PurchaseClaimActionReplacement {
				replacementPending = baseQuantity
			}

			claim.Items = append(claim.Items, models.PurchaseClaimItem{
				PurchaseClaimID:         claim.ID,
				PurchaseOrderItemID:     lineItem.ID,
				ProductID:               *lineItem.ProductID,
				ProductName:             lineItem.ProductName,
				SKU:                     lineItem.SKU,
				IsRawMaterial:           lineItem.IsRawMaterial,
				Type:                    itemInput.Type,
				Quantity:                itemInput.Quantity,
				Unit:                    strings.ToLower(strings.TrimSpace(itemInput.Unit)),
				BaseQuantity:            baseQuantity,
				BaseUnit:                baseUnit,
				Rate:                    lineItem.Rate,
				Amount:                  baseQuantity * lineItem.Rate,
				Reason:                  itemInput.Reason,
				Action:                  itemInput.Action,
				StockAdjusted:           stockAdjusted,
				ReplacementPendingBase:  replacementPending,
				ReplacementReceivedBase: 0,
				ReplacementCompleted:    false,
				CreatedAt:               now,
				UpdatedAt:               now,
			})
		}

		return s.repository.CreateTx(tx, claim)
	})
	if err != nil {
		return nil, err
	}

	created, err := s.repository.FindByIDAndCompany(claim.ID, companyID)
	if err != nil {
		return nil, err
	}
	return output.ToPurchaseClaimOutput(created), nil
}

func applyPurchaseClaimToCorrectStock(
	tx *gorm.DB,
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	claimType models.PurchaseClaimType,
	reason string,
	referenceID string,
	referenceNumber string,
	userID string,
) error {
	if lineItem == nil || lineItem.ProductID == nil || *lineItem.ProductID == "" {
		return errors.New("valid purchase order line item is required")
	}

	if lineItem.IsRawMaterial {
		return applyRawMaterialClaimToProductStock(tx, lineItem, quantity, claimType, reason, referenceID, referenceNumber, userID)
	}

	return applyVariantClaimToVariantStock(tx, lineItem, quantity, claimType, reason, referenceID, referenceNumber, userID)
}

func applyRawMaterialClaimToProductStock(
	tx *gorm.DB,
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	claimType models.PurchaseClaimType,
	reason string,
	referenceID string,
	referenceNumber string,
	userID string,
) error {
	productID := *lineItem.ProductID
	var stock models.ProductStock

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ?", productID).
		First(&stock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("raw-material product stock not found for %s", productID)
	}
	if err != nil {
		return err
	}
	if quantity <= 0 {
		return errors.New("claim quantity must be greater than zero")
	}
	if stock.AvailableStock+0.000001 < quantity {
		return fmt.Errorf("insufficient raw-material available stock for %s: available %.3f, requested %.3f", productID, stock.AvailableStock, quantity)
	}

	before := stock.AvailableStock
	now := time.Now()
	switch claimType {
	case models.PurchaseClaimMissing:
		if stock.CurrentStock+0.000001 < quantity || stock.PurchasedStock+0.000001 < quantity {
			return fmt.Errorf("cannot remove missing raw-material quantity %.3f from stock", quantity)
		}
		stock.CurrentStock -= quantity
		stock.PurchasedStock -= quantity
		stock.AvailableStock -= quantity
	case models.PurchaseClaimDamaged:
		stock.AvailableStock -= quantity
		stock.DamagedStock += quantity
		stock.DamageReason = reason
		stock.DamagedAt = &now
		stock.DamagedBy = userID
	default:
		return fmt.Errorf("unsupported raw-material claim type: %s", claimType)
	}

	stock.LastStockSyncAt = now
	stock.UpdatedAt = now
	if err := tx.Save(&stock).Error; err != nil {
		return fmt.Errorf("failed to update raw-material stock: %w", err)
	}

	cost := lineItem.Rate
	if cost <= 0 {
		cost = stock.AverageCost
	}
	movementType := "PURCHASE_MISSING"
	if claimType == models.PurchaseClaimDamaged {
		movementType = "PURCHASE_DAMAGED"
	}
	ledger := &models.StockLedger{
		ProductID: productID, MovementType: movementType, Quantity: -quantity, Rate: cost,
		Amount: -(quantity * cost), ReferenceType: "purchase_claim", ReferenceID: referenceID,
		ReferenceNumber: referenceNumber, BalanceBeforeQty: before, BalanceAfterQty: stock.AvailableStock,
		CostBeforeAmount: before * cost, CostAfterAmount: stock.AvailableStock * cost,
		Notes: reason, CreatedAt: now, CreatedBy: userID,
	}
	if err := tx.Create(ledger).Error; err != nil {
		return fmt.Errorf("failed to create raw-material stock ledger: %w", err)
	}
	return nil
}

func applyVariantClaimToVariantStock(
	tx *gorm.DB,
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	claimType models.PurchaseClaimType,
	reason string,
	referenceID string,
	referenceNumber string,
	userID string,
) error {
	productID := *lineItem.ProductID
	variantSKU := strings.TrimSpace(lineItem.SKU)
	if variantSKU == "" {
		return fmt.Errorf("variant SKU is required for product %s", productID)
	}

	var stock models.VariantStock
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ? AND variant_sku = ?", productID, variantSKU).
		First(&stock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("variant stock not found for product %s and SKU %s", productID, variantSKU)
	}
	if err != nil {
		return err
	}
	if quantity <= 0 {
		return errors.New("claim quantity must be greater than zero")
	}
	if stock.AvailableStock+0.000001 < quantity {
		return fmt.Errorf("insufficient variant available stock for %s: available %.3f, requested %.3f", variantSKU, stock.AvailableStock, quantity)
	}

	before := stock.AvailableStock
	now := time.Now()
	switch claimType {
	case models.PurchaseClaimMissing:
		if stock.CurrentStock+0.000001 < quantity || stock.PurchasedStock+0.000001 < quantity {
			return fmt.Errorf("cannot remove missing variant quantity %.3f from stock", quantity)
		}
		stock.CurrentStock -= quantity
		stock.PurchasedStock -= quantity
		stock.AvailableStock -= quantity
	case models.PurchaseClaimDamaged:
		stock.AvailableStock -= quantity
		stock.DamagedStock += quantity
		stock.DamageReason = reason
		stock.DamagedAt = &now
		stock.DamagedBy = userID
	default:
		return fmt.Errorf("unsupported variant claim type: %s", claimType)
	}

	stock.LastStockSyncAt = &now
	stock.UpdatedAt = now
	if err := tx.Save(&stock).Error; err != nil {
		return fmt.Errorf("failed to update variant stock: %w", err)
	}

	cost := lineItem.Rate
	if cost <= 0 {
		cost = stock.AverageCost
	}
	movementType := "PURCHASE_MISSING"
	if claimType == models.PurchaseClaimDamaged {
		movementType = "PURCHASE_DAMAGED"
	}
	movement := &models.VariantStockMovement{
		VariantID: stock.ID, ProductID: stock.ProductID, VariantSKU: stock.VariantSKU,
		MovementType: movementType, Quantity: -quantity, Rate: cost, Amount: -(quantity * cost),
		ReferenceType: "purchase_claim", ReferenceID: referenceID, ReferenceNumber: referenceNumber,
		BalanceBeforeQty: before, BalanceAfterQty: stock.AvailableStock, Stage: "confirmed",
		Notes: reason, CreatedAt: now, CreatedBy: userID,
	}
	if err := tx.Create(movement).Error; err != nil {
		return fmt.Errorf("failed to create variant stock movement: %w", err)
	}
	return nil
}

func (s *purchaseClaimService) ReceiveReplacement(
	claimID string,
	request *input.ReceivePurchaseClaimReplacementInput,
	userID string,
	companyID uint,
) (*output.PurchaseClaimOutput, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	item, claim, err := s.repository.FindItemByIDAndCompany(request.PurchaseClaimItemID, companyID)
	if err != nil {
		return nil, errors.New("purchase claim item not found")
	}
	if claim.ID != claimID {
		return nil, errors.New("claim item does not belong to this purchase claim")
	}
	if item.Action != models.PurchaseClaimActionReplacement {
		return nil, errors.New("this claim item is not marked for vendor replacement")
	}
	if item.ReplacementCompleted {
		return nil, errors.New("replacement is already completed")
	}

	po, err := s.poRepo.FindByIDAndCompany(claim.PurchaseOrderID, companyID)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}
	lineItem, err := findPurchaseOrderLineItem(po, item.PurchaseOrderItemID)
	if err != nil {
		return nil, err
	}

	receivedBase, baseUnit, err := normalizeClaimQuantity(lineItem, request.Quantity, request.Unit)
	if err != nil {
		return nil, err
	}
	if receivedBase > item.ReplacementPendingBase+0.000001 {
		return nil, fmt.Errorf("replacement received %.3f %s exceeds pending %.3f %s", receivedBase, baseUnit, item.ReplacementPendingBase, item.BaseUnit)
	}

	now := time.Now()
	err = s.repository.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := addReplacementToCorrectStock(tx, lineItem, receivedBase, item.Rate, claim.ID, claim.ClaimNumber, request.Notes, userID); err != nil {
			return err
		}

		item.ReplacementReceivedBase += receivedBase
		item.ReplacementPendingBase -= receivedBase
		if item.ReplacementPendingBase <= 0.000001 {
			item.ReplacementPendingBase = 0
			item.ReplacementCompleted = true
			item.ReplacementCompletedAt = &now
		}
		item.UpdatedAt = now
		if err := tx.Save(item).Error; err != nil {
			return fmt.Errorf("failed to update claim item: %w", err)
		}

		receipt := &models.PurchaseClaimReceipt{
			ID: uuid.New().String(), PurchaseClaimID: claim.ID, PurchaseClaimItemID: item.ID,
			ProductID: item.ProductID, ReceivedQuantity: request.Quantity,
			ReceivedUnit:         strings.ToLower(strings.TrimSpace(request.Unit)),
			ReceivedBaseQuantity: receivedBase, BaseUnit: baseUnit,
			ReceivedDate: request.ReceivedDate, Notes: request.Notes,
			ReceivedBy: userID, CreatedAt: now,
		}
		if err := s.repository.CreateReceiptTx(tx, receipt); err != nil {
			return fmt.Errorf("failed to save replacement receipt: %w", err)
		}

		var pendingCount int64
		if err := tx.Model(&models.PurchaseClaimItem{}).
			Where("purchase_claim_id = ? AND action = ? AND replacement_pending_base > 0", claim.ID, models.PurchaseClaimActionReplacement).
			Count(&pendingCount).Error; err != nil {
			return err
		}
		status := models.PurchaseClaimStatusResolved
		if pendingCount > 0 {
			status = models.PurchaseClaimStatusPartial
		}
		return s.repository.UpdateClaimStatusTx(tx, claim.ID, status)
	})
	if err != nil {
		return nil, err
	}

	updated, err := s.repository.FindByIDAndCompany(claim.ID, companyID)
	if err != nil {
		return nil, err
	}
	return output.ToPurchaseClaimOutput(updated), nil
}

func addReplacementToCorrectStock(
	tx *gorm.DB,
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	rate float64,
	referenceID string,
	referenceNumber string,
	notes string,
	userID string,
) error {
	if lineItem == nil || lineItem.ProductID == nil || *lineItem.ProductID == "" {
		return errors.New("valid purchase order line item is required")
	}
	if lineItem.IsRawMaterial {
		return addRawMaterialReplacementStock(tx, lineItem, quantity, rate, referenceID, referenceNumber, notes, userID)
	}
	return addVariantReplacementStock(tx, lineItem, quantity, rate, referenceID, referenceNumber, notes, userID)
}

func addRawMaterialReplacementStock(
	tx *gorm.DB,
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	rate float64,
	referenceID string,
	referenceNumber string,
	notes string,
	userID string,
) error {
	productID := *lineItem.ProductID
	now := time.Now()
	var stock models.ProductStock
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ?", productID).First(&stock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stock = models.ProductStock{ID: uuid.New().String(), ProductID: productID, ProductName: lineItem.ProductName, SKU: lineItem.SKU, RawMaterialUnit: "gram", AverageCost: rate, LastStockSyncAt: now, CreatedAt: now, UpdatedAt: now}
		if err := tx.Create(&stock).Error; err != nil {
			return fmt.Errorf("failed to create raw-material stock: %w", err)
		}
	} else if err != nil {
		return err
	}

	before := stock.AvailableStock
	cost := rate
	if cost <= 0 {
		cost = stock.AverageCost
	}
	stock.PurchasedStock += quantity
	stock.CurrentStock += quantity
	stock.AvailableStock += quantity
	stock.LastPurchasedDate = &now
	stock.LastStockSyncAt = now
	stock.UpdatedAt = now
	if stock.AverageCost <= 0 {
		stock.AverageCost = cost
	}
	if err := tx.Save(&stock).Error; err != nil {
		return fmt.Errorf("failed to add raw-material replacement stock: %w", err)
	}

	ledger := &models.StockLedger{ProductID: productID, MovementType: "VENDOR_REPLACEMENT", Quantity: quantity, Rate: cost, Amount: quantity * cost, ReferenceType: "purchase_claim", ReferenceID: referenceID, ReferenceNumber: referenceNumber, BalanceBeforeQty: before, BalanceAfterQty: stock.AvailableStock, CostBeforeAmount: before * cost, CostAfterAmount: stock.AvailableStock * cost, Notes: notes, CreatedAt: now, CreatedBy: userID}
	if err := tx.Create(ledger).Error; err != nil {
		return fmt.Errorf("failed to create raw-material replacement ledger: %w", err)
	}
	return nil
}

func addVariantReplacementStock(
	tx *gorm.DB,
	lineItem *models.PurchaseOrderLineItem,
	quantity float64,
	rate float64,
	referenceID string,
	referenceNumber string,
	notes string,
	userID string,
) error {
	productID := *lineItem.ProductID
	variantSKU := strings.TrimSpace(lineItem.SKU)
	if variantSKU == "" {
		return fmt.Errorf("variant SKU is required for product %s", productID)
	}
	now := time.Now()
	var stock models.VariantStock
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ? AND variant_sku = ?", productID, variantSKU).First(&stock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stock = models.VariantStock{ID: uuid.New().String(), ProductID: productID, VariantSKU: variantSKU, VariantName: variantSKU, ProductName: lineItem.ProductName, AverageCost: rate, LastStockSyncAt: &now, CreatedAt: now, UpdatedAt: now}
		if err := tx.Create(&stock).Error; err != nil {
			return fmt.Errorf("failed to create variant stock: %w", err)
		}
	} else if err != nil {
		return err
	}

	before := stock.AvailableStock
	cost := rate
	if cost <= 0 {
		cost = stock.AverageCost
	}
	stock.PurchasedStock += quantity
	stock.CurrentStock += quantity
	stock.AvailableStock += quantity
	stock.LastPurchasedDate = &now
	stock.LastStockSyncAt = &now
	stock.UpdatedAt = now
	if stock.AverageCost <= 0 {
		stock.AverageCost = cost
	}
	if err := tx.Save(&stock).Error; err != nil {
		return fmt.Errorf("failed to add variant replacement stock: %w", err)
	}

	movement := &models.VariantStockMovement{VariantID: stock.ID, ProductID: stock.ProductID, VariantSKU: stock.VariantSKU, MovementType: "VENDOR_REPLACEMENT", Quantity: quantity, Rate: cost, Amount: quantity * cost, ReferenceType: "purchase_claim", ReferenceID: referenceID, ReferenceNumber: referenceNumber, BalanceBeforeQty: before, BalanceAfterQty: stock.AvailableStock, Stage: "confirmed", Notes: notes, CreatedAt: now, CreatedBy: userID}
	if err := tx.Create(movement).Error; err != nil {
		return fmt.Errorf("failed to create variant replacement movement: %w", err)
	}
	return nil
}

func (s *purchaseClaimService) GetClaimByID(
	id string,
	companyID uint,
) (*output.PurchaseClaimOutput, error) {
	claim, err := s.repository.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("purchase claim not found")
	}
	return output.ToPurchaseClaimOutput(claim), nil
}

func (s *purchaseClaimService) GetClaimsByPurchaseOrder(
	purchaseOrderID string,
	companyID uint,
) ([]output.PurchaseClaimOutput, error) {
	if _, err := s.poRepo.FindByIDAndCompany(purchaseOrderID, companyID); err != nil {
		return nil, errors.New("purchase order not found in your company")
	}

	claims, err := s.repository.FindByPurchaseOrderAndCompany(purchaseOrderID, companyID)
	if err != nil {
		return nil, err
	}

	result := make([]output.PurchaseClaimOutput, 0, len(claims))
	for index := range claims {
		result = append(result, *output.ToPurchaseClaimOutput(&claims[index]))
	}
	return result, nil
}

func (s *purchaseClaimService) GetReplacementReceipts(
	claimItemID uint,
	companyID uint,
) ([]models.PurchaseClaimReceipt, error) {
	return s.repository.FindReceiptsByItemAndCompany(claimItemID, companyID)
}

func (s *purchaseClaimService) GetNetReceivableBaseQuantity(
	purchaseOrderItemID uint,
	orderedBaseQuantity float64,
) (float64, error) {
	missing, err := s.repository.SumClaimedByPOItem(
		nil,
		purchaseOrderItemID,
		models.PurchaseClaimMissing,
	)
	if err != nil {
		return 0, err
	}
	return math.Max(orderedBaseQuantity-missing, 0), nil
}
