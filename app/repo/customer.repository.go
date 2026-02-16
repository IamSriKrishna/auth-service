package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(customer).Updates(customer).Error; err != nil {
			return err
		}

		if customer.OtherDetails != nil {
			customer.OtherDetails.EntityID = customer.ID
			customer.OtherDetails.EntityType = "customer"
			if customer.OtherDetails.ID == 0 {
				if err := tx.Create(customer.OtherDetails).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(customer.OtherDetails).Updates(customer.OtherDetails).Error; err != nil {
					return err
				}
			}
		}

		for i := range customer.Addresses {
			customer.Addresses[i].EntityID = customer.ID
			customer.Addresses[i].EntityType = "customer"

			var existingAddr models.EntityAddress
			err := tx.Where("entity_id = ? AND entity_type = ? AND address_type = ?",
				customer.ID, "customer", customer.Addresses[i].AddressType).
				First(&existingAddr).Error

			switch err {
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&customer.Addresses[i]).Error; err != nil {
					return err
				}
			case nil:
				customer.Addresses[i].ID = existingAddr.ID
				if err := tx.Model(&customer.Addresses[i]).Updates(&customer.Addresses[i]).Error; err != nil {
					return err
				}
			default:
				return err
			}
		}

		if err := tx.Where("entity_id = ? AND entity_type = ?", customer.ID, "customer").
			Delete(&models.EntityContactPerson{}).Error; err != nil {
			return err
		}

		for i := range customer.ContactPersons {
			customer.ContactPersons[i].EntityID = customer.ID
			customer.ContactPersons[i].EntityType = "customer"
			customer.ContactPersons[i].ID = 0
			if err := tx.Create(&customer.ContactPersons[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *customerRepository) FindByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.
		Preload("OtherDetails").
		Preload("Addresses").
		Preload("ContactPersons").
		First(&customer, id).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindAll(page, limit int) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&customers).Error

	return customers, total, err
}

func (r *customerRepository) Delete(customer *models.Customer) error {
	return r.db.Select(clause.Associations).Delete(customer).Error
}

func (r *customerRepository) FindByMobile(mobile string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("mobile = ?", mobile).First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
