package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)


type vendorShortageClaimRepository struct {
	db *gorm.DB
}

func NewVendorShortageClaimRepository(db *gorm.DB) VendorShortageClaimRepository {
	return &vendorShortageClaimRepository{db: db}
}

func (r *vendorShortageClaimRepository) Create(claim *models.VendorShortageClaim) error {
	return r.db.Create(claim).Error
}

func (r *vendorShortageClaimRepository) FindByPurchaseOrderID(poID string) ([]models.VendorShortageClaim, error) {
	var claims []models.VendorShortageClaim
	err := r.db.Where("purchase_order_id = ?", poID).Find(&claims).Error
	return claims, err
}