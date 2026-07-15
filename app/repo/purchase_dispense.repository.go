package repo

import (
	"errors"
	"strings"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type PurchaseDispenseRepository interface {
	CreateTx(
		tx *gorm.DB,
		value *models.PurchaseDispense,
	) error

	FindByIDAndCompany(
		id string,
		companyID uint,
	) (*models.PurchaseDispense, error)

	FindByClaimAndCompany(
		claimID string,
		companyID uint,
	) ([]models.PurchaseDispense, error)

	FindByClaimItemAndCompany(
		claimItemID uint,
		companyID uint,
	) ([]models.PurchaseDispense, error)

	GetDB() *gorm.DB
}

type purchaseDispenseRepository struct {
	db *gorm.DB
}

func NewPurchaseDispenseRepository(
	db *gorm.DB,
) PurchaseDispenseRepository {
	return &purchaseDispenseRepository{
		db: db,
	}
}

func (r *purchaseDispenseRepository) CreateTx(
	tx *gorm.DB,
	value *models.PurchaseDispense,
) error {
	if value == nil {
		return gorm.ErrInvalidData
	}

	if tx == nil {
		tx = r.db
	}

	return tx.Create(value).Error
}

func (r *purchaseDispenseRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.PurchaseDispense, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New(
			"purchase dispense ID is required",
		)
	}

	var value models.PurchaseDispense

	err := r.db.
		Model(&models.PurchaseDispense{}).
		Joins(`
			JOIN purchase_claims pc
				ON pc.id =
					purchase_dispenses.purchase_claim_id
		`).
		Where(
			"purchase_dispenses.id = ? AND pc.company_id = ?",
			id,
			companyID,
		).
		Select("purchase_dispenses.*").
		First(&value).
		Error
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (r *purchaseDispenseRepository) FindByClaimAndCompany(
	claimID string,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	var values []models.PurchaseDispense

	err := r.db.
		Model(&models.PurchaseDispense{}).
		Joins(`
			JOIN purchase_claims pc
				ON pc.id =
					purchase_dispenses.purchase_claim_id
		`).
		Where(
			"purchase_dispenses.purchase_claim_id = ? AND pc.company_id = ?",
			claimID,
			companyID,
		).
		Select("purchase_dispenses.*").
		Order(
			"purchase_dispenses.created_at DESC",
		).
		Find(&values).
		Error

	return values, err
}

func (r *purchaseDispenseRepository) FindByClaimItemAndCompany(
	claimItemID uint,
	companyID uint,
) ([]models.PurchaseDispense, error) {
	var values []models.PurchaseDispense

	err := r.db.
		Model(&models.PurchaseDispense{}).
		Joins(`
			JOIN purchase_claims pc
				ON pc.id =
					purchase_dispenses.purchase_claim_id
		`).
		Where(
			"purchase_dispenses.purchase_claim_item_id = ? AND pc.company_id = ?",
			claimItemID,
			companyID,
		).
		Select("purchase_dispenses.*").
		Order(
			"purchase_dispenses.created_at DESC",
		).
		Find(&values).
		Error

	return values, err
}

func (r *purchaseDispenseRepository) GetDB() *gorm.DB {
	return r.db
}
