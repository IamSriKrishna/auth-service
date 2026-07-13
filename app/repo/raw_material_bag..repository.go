package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)


type rawMaterialBagRepository struct {
	db *gorm.DB
}

func NewRawMaterialBagRepository(
	db *gorm.DB,
) RawMaterialBagRepository {
	return &rawMaterialBagRepository{
		db: db,
	}
}

// Raw-material bag company ownership is resolved through:
//
// raw_material_bags.created_by
//     -> users.id
//     -> users.company_id
//
// No new company column is required in raw_material_bags.
func rawMaterialBagCompanyQuery(
	db *gorm.DB,
	companyID uint,
) *gorm.DB {
	return db.
		Joins(`
			JOIN users raw_bag_creator
				ON raw_bag_creator.id =
					CAST(raw_material_bags.created_by AS UNSIGNED)
		`).
		Where(
			"raw_bag_creator.company_id = ?",
			companyID,
		)
}

func (r *rawMaterialBagRepository) CreateMany(
	bags []models.RawMaterialBag,
) error {
	if len(bags) == 0 {
		return nil
	}

	return r.db.Create(&bags).Error
}

func (r *rawMaterialBagRepository) CreateManyForCompany(
	bags []models.RawMaterialBag,
	companyID uint,
) error {
	if len(bags) == 0 {
		return nil
	}

	for index := range bags {
		if bags[index].CreatedBy == "" {
			return gorm.ErrInvalidData
		}

		var count int64
		err := r.db.
			Table("users").
			Where(
				"id = CAST(? AS UNSIGNED) AND company_id = ?",
				bags[index].CreatedBy,
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
	}

	return r.db.Create(&bags).Error
}

func (r *rawMaterialBagRepository) GetAll(
	limit int,
	offset int,
) ([]models.RawMaterialBag, int64, error) {
	var bags []models.RawMaterialBag
	var total int64

	if err := r.db.
		Model(&models.RawMaterialBag{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bags).
		Error

	return bags, total, err
}

func (r *rawMaterialBagRepository) GetAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.RawMaterialBag, int64, error) {
	var bags []models.RawMaterialBag
	var total int64

	countQuery := rawMaterialBagCompanyQuery(
		r.db.Model(&models.RawMaterialBag{}),
		companyID,
	)

	if err := countQuery.
		Distinct("raw_material_bags.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := rawMaterialBagCompanyQuery(
		r.db.Model(&models.RawMaterialBag{}),
		companyID,
	)

	err := findQuery.
		Select("raw_material_bags.*").
		Order("raw_material_bags.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bags).
		Error

	return bags, total, err
}

func (r *rawMaterialBagRepository) GetByID(
	id string,
) (*models.RawMaterialBag, error) {
	var bag models.RawMaterialBag

	if err := r.db.
		Where("id = ?", id).
		First(&bag).
		Error; err != nil {
		return nil, err
	}

	return &bag, nil
}

func (r *rawMaterialBagRepository) GetByIDAndCompany(
	id string,
	companyID uint,
) (*models.RawMaterialBag, error) {
	var bag models.RawMaterialBag

	query := rawMaterialBagCompanyQuery(
		r.db.Model(&models.RawMaterialBag{}),
		companyID,
	)

	if err := query.
		Select("raw_material_bags.*").
		Where("raw_material_bags.id = ?", id).
		First(&bag).
		Error; err != nil {
		return nil, err
	}

	return &bag, nil
}

func (r *rawMaterialBagRepository) GetByProductID(
	productID string,
) ([]models.RawMaterialBag, error) {
	var bags []models.RawMaterialBag

	err := r.db.
		Where("product_id = ?", productID).
		Order("created_at ASC").
		Find(&bags).
		Error

	return bags, err
}

func (r *rawMaterialBagRepository) GetByProductIDAndCompany(
	productID string,
	companyID uint,
) ([]models.RawMaterialBag, error) {
	var bags []models.RawMaterialBag

	query := rawMaterialBagCompanyQuery(
		r.db.Model(&models.RawMaterialBag{}),
		companyID,
	)

	err := query.
		Select("raw_material_bags.*").
		Where(
			"raw_material_bags.product_id = ?",
			productID,
		).
		Order("raw_material_bags.created_at ASC").
		Find(&bags).
		Error

	return bags, err
}

func (r *rawMaterialBagRepository) GetByPurchaseOrderID(
	purchaseOrderID string,
) ([]models.RawMaterialBag, error) {
	var bags []models.RawMaterialBag

	err := r.db.
		Where(
			"purchase_order_id = ?",
			purchaseOrderID,
		).
		Order("created_at ASC").
		Find(&bags).
		Error

	return bags, err
}

func (r *rawMaterialBagRepository) GetByPurchaseOrderIDAndCompany(
	purchaseOrderID string,
	companyID uint,
) ([]models.RawMaterialBag, error) {
	var bags []models.RawMaterialBag

	query := rawMaterialBagCompanyQuery(
		r.db.Model(&models.RawMaterialBag{}),
		companyID,
	)

	err := query.
		Select("raw_material_bags.*").
		Where(
			"raw_material_bags.purchase_order_id = ?",
			purchaseOrderID,
		).
		Order("raw_material_bags.created_at ASC").
		Find(&bags).
		Error

	return bags, err
}

func (r *rawMaterialBagRepository) Update(
	bag *models.RawMaterialBag,
) error {
	if bag == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Save(bag).Error
}

func (r *rawMaterialBagRepository) UpdateByCompany(
	bag *models.RawMaterialBag,
	companyID uint,
) error {
	if bag == nil {
		return gorm.ErrInvalidData
	}

	if _, err := r.GetByIDAndCompany(
		bag.ID,
		companyID,
	); err != nil {
		return err
	}

	result := r.db.
		Model(&models.RawMaterialBag{}).
		Where("id = ?", bag.ID).
		Updates(bag)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
