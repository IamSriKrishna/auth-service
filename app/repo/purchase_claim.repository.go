package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type PurchaseClaimRepository interface {
	CreateTx(tx *gorm.DB, claim *models.PurchaseClaim) error
	CreateReceiptTx(tx *gorm.DB, receipt *models.PurchaseClaimReceipt) error

	FindByIDAndCompany(id string, companyID uint) (*models.PurchaseClaim, error)
	FindByPurchaseOrderAndCompany(purchaseOrderID string, companyID uint) ([]models.PurchaseClaim, error)
	FindItemByIDAndCompany(itemID uint, companyID uint) (*models.PurchaseClaimItem, *models.PurchaseClaim, error)
	FindReceiptsByItemAndCompany(itemID uint, companyID uint) ([]models.PurchaseClaimReceipt, error)

	SumClaimedByPOItem(
		tx *gorm.DB,
		purchaseOrderItemID uint,
		claimType models.PurchaseClaimType,
	) (float64, error)

	SumReplacementPendingByPOItem(
		tx *gorm.DB,
		purchaseOrderItemID uint,
	) (float64, error)

	UpdateClaimStatusTx(
		tx *gorm.DB,
		claimID string,
		status models.PurchaseClaimStatus,
	) error

	GetDB() *gorm.DB
}

type purchaseClaimRepository struct {
	db *gorm.DB
}

func NewPurchaseClaimRepository(db *gorm.DB) PurchaseClaimRepository {
	return &purchaseClaimRepository{db: db}
}

func (r *purchaseClaimRepository) CreateTx(
	tx *gorm.DB,
	claim *models.PurchaseClaim,
) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(claim).Error
}

func (r *purchaseClaimRepository) CreateReceiptTx(
	tx *gorm.DB,
	receipt *models.PurchaseClaimReceipt,
) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(receipt).Error
}

func (r *purchaseClaimRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.PurchaseClaim, error) {
	var claim models.PurchaseClaim
	err := r.db.
		Preload("Items").
		Where("id = ? AND company_id = ?", id, companyID).
		First(&claim).Error
	return &claim, err
}

func (r *purchaseClaimRepository) FindByPurchaseOrderAndCompany(
	purchaseOrderID string,
	companyID uint,
) ([]models.PurchaseClaim, error) {
	var claims []models.PurchaseClaim
	err := r.db.
		Preload("Items").
		Where("purchase_order_id = ? AND company_id = ?", purchaseOrderID, companyID).
		Order("created_at DESC").
		Find(&claims).Error
	return claims, err
}

func (r *purchaseClaimRepository) FindItemByIDAndCompany(
	itemID uint,
	companyID uint,
) (*models.PurchaseClaimItem, *models.PurchaseClaim, error) {
	var item models.PurchaseClaimItem
	err := r.db.
		Joins("JOIN purchase_claims pc ON pc.id = purchase_claim_items.purchase_claim_id").
		Where("purchase_claim_items.id = ? AND pc.company_id = ?", itemID, companyID).
		First(&item).Error
	if err != nil {
		return nil, nil, err
	}

	var claim models.PurchaseClaim
	err = r.db.
		Where("id = ? AND company_id = ?", item.PurchaseClaimID, companyID).
		First(&claim).Error
	if err != nil {
		return nil, nil, err
	}

	return &item, &claim, nil
}

func (r *purchaseClaimRepository) FindReceiptsByItemAndCompany(
	itemID uint,
	companyID uint,
) ([]models.PurchaseClaimReceipt, error) {
	var receipts []models.PurchaseClaimReceipt
	err := r.db.
		Joins("JOIN purchase_claims pc ON pc.id = purchase_claim_receipts.purchase_claim_id").
		Where("purchase_claim_receipts.purchase_claim_item_id = ? AND pc.company_id = ?", itemID, companyID).
		Order("purchase_claim_receipts.created_at DESC").
		Find(&receipts).Error
	return receipts, err
}

func (r *purchaseClaimRepository) SumClaimedByPOItem(
	tx *gorm.DB,
	purchaseOrderItemID uint,
	claimType models.PurchaseClaimType,
) (float64, error) {
	if tx == nil {
		tx = r.db
	}

	var total float64
	err := tx.
		Model(&models.PurchaseClaimItem{}).
		Joins("JOIN purchase_claims pc ON pc.id = purchase_claim_items.purchase_claim_id").
		Where(
			"purchase_claim_items.purchase_order_item_id = ? AND purchase_claim_items.type = ? AND pc.status <> ?",
			purchaseOrderItemID,
			claimType,
			models.PurchaseClaimStatusCancelled,
		).
		Select("COALESCE(SUM(purchase_claim_items.base_quantity), 0)").
		Scan(&total).Error
	return total, err
}

func (r *purchaseClaimRepository) SumReplacementPendingByPOItem(
	tx *gorm.DB,
	purchaseOrderItemID uint,
) (float64, error) {
	if tx == nil {
		tx = r.db
	}

	var total float64
	err := tx.
		Model(&models.PurchaseClaimItem{}).
		Joins("JOIN purchase_claims pc ON pc.id = purchase_claim_items.purchase_claim_id").
		Where(
			"purchase_claim_items.purchase_order_item_id = ? AND purchase_claim_items.action = ? AND pc.status <> ?",
			purchaseOrderItemID,
			models.PurchaseClaimActionReplacement,
			models.PurchaseClaimStatusCancelled,
		).
		Select("COALESCE(SUM(purchase_claim_items.replacement_pending_base), 0)").
		Scan(&total).Error
	return total, err
}

func (r *purchaseClaimRepository) UpdateClaimStatusTx(
	tx *gorm.DB,
	claimID string,
	status models.PurchaseClaimStatus,
) error {
	if tx == nil {
		tx = r.db
	}
	return tx.
		Model(&models.PurchaseClaim{}).
		Where("id = ?", claimID).
		Update("status", status).Error
}

func (r *purchaseClaimRepository) GetDB() *gorm.DB {
	return r.db
}
