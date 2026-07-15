package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type PurchaseDispenseRepository interface {
	CreateTx(tx *gorm.DB, value *models.PurchaseDispense) error
	FindByIDAndCompany(id string, companyID uint) (*models.PurchaseDispense, error)
	FindByClaimAndCompany(claimID string, companyID uint) ([]models.PurchaseDispense, error)
	FindByClaimItemAndCompany(claimItemID uint, companyID uint) ([]models.PurchaseDispense, error)
	GetDB() *gorm.DB
}

type purchaseDispenseRepository struct {
	db *gorm.DB
}

func NewPurchaseDispenseRepository(db *gorm.DB) PurchaseDispenseRepository {
	return &purchaseDispenseRepository{db: db}
}

func (r *purchaseDispenseRepository) CreateTx(
	tx *gorm.DB,
	value *models.PurchaseDispense,
) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(value).Error
}

func (r *purchaseDispenseRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.PurchaseDispense, error) {
	var value models.PurchaseDispense
	err := r.db.
		Joins("JOIN purchase_claims pc ON pc.id = purchase_dispenses.purchase_claim_id").
		Where("purchase_dispenses.id = ? AND pc.company_id = ?", id, companyID).
		First(&value).Error
	return &value, err
}

func (r *purchaseDispenseRepository) FindByClaimAndCompany(
	claimID string,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	var values []models.PurchaseDispense
	err := r.db.
		Joins("JOIN purchase_claims pc ON pc.id = purchase_dispenses.purchase_claim_id").
		Where("purchase_dispenses.purchase_claim_id = ? AND pc.company_id = ?", claimID, companyID).
		Order("purchase_dispenses.created_at DESC").
		Find(&values).Error
	return values, err
}

func (r *purchaseDispenseRepository) FindByClaimItemAndCompany(
	claimItemID uint,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	var values []models.PurchaseDispense
	err := r.db.
		Joins("JOIN purchase_claims pc ON pc.id = purchase_dispenses.purchase_claim_id").
		Where("purchase_dispenses.purchase_claim_item_id = ? AND pc.company_id = ?", claimItemID, companyID).
		Order("purchase_dispenses.created_at DESC").
		Find(&values).Error
	return values, err
}

func (r *purchaseDispenseRepository) GetDB() *gorm.DB {
	return r.db
}
