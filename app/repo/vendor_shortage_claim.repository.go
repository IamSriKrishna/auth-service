package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type vendorShortageClaimRepository struct {
	db *gorm.DB
}

func NewVendorShortageClaimRepository(
	db *gorm.DB,
) VendorShortageClaimRepository {
	return &vendorShortageClaimRepository{
		db: db,
	}
}

func (r *vendorShortageClaimRepository) Create(
	claim *models.VendorShortageClaim,
) error {
	return r.db.Create(claim).Error
}

func (r *vendorShortageClaimRepository) FindByPurchaseOrderID(
	purchaseOrderID string,
) ([]models.VendorShortageClaim, error) {
	var claims []models.VendorShortageClaim

	err := r.db.
		Where(
			"purchase_order_id = ?",
			purchaseOrderID,
		).
		Order("created_at DESC").
		Find(&claims).
		Error

	return claims, err
}

func (r *vendorShortageClaimRepository) CreateForCompany(
	claim *models.VendorShortageClaim,
	companyID uint,
) error {
	if claim == nil {
		return gorm.ErrInvalidData
	}

	var count int64
	err := r.db.
		Table("purchase_orders").
		Joins(`
			JOIN users po_creator
				ON po_creator.id =
					CAST(purchase_orders.created_by AS UNSIGNED)
		`).
		Where(
			"purchase_orders.id = ? AND po_creator.company_id = ?",
			claim.PurchaseOrderID,
			companyID,
		).
		Count(&count).
		Error
	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.Create(claim).Error
}

func (r *vendorShortageClaimRepository) FindByPurchaseOrderIDAndCompany(
	purchaseOrderID string,
	companyID uint,
) ([]models.VendorShortageClaim, error) {
	var claims []models.VendorShortageClaim

	err := r.db.
		Model(&models.VendorShortageClaim{}).
		Joins(`
			JOIN purchase_orders
				ON purchase_orders.id =
					vendor_shortage_claims.purchase_order_id
		`).
		Joins(`
			JOIN users po_creator
				ON po_creator.id =
					CAST(purchase_orders.created_by AS UNSIGNED)
		`).
		Where(
			"vendor_shortage_claims.purchase_order_id = ? AND po_creator.company_id = ?",
			purchaseOrderID,
			companyID,
		).
		Select("vendor_shortage_claims.*").
		Order(
			"vendor_shortage_claims.created_at DESC",
		).
		Find(&claims).
		Error

	return claims, err
}
