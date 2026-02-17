package repo

import (
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type bankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &bankRepository{db: db}
}

func (r *bankRepository) Create(bank *models.Bank) error {
	return r.db.Create(bank).Error
}

func (r *bankRepository) FindByID(id uint) (*models.Bank, error) {
	var bank models.Bank
	err := r.db.First(&bank, id).Error
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *bankRepository) FindByIFSCCode(ifscCode string) (*models.Bank, error) {
	var bank models.Bank
	err := r.db.Where("ifsc_code = ?", ifscCode).First(&bank).Error
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *bankRepository) FindAll(limit, offset int) ([]models.Bank, int64, error) {
	var banks []models.Bank
	var count int64
	err := r.db.Model(&models.Bank{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Limit(limit).Offset(offset).Find(&banks).Error
	if err != nil {
		return nil, 0, err
	}
	return banks, count, nil
}

func (r *bankRepository) Update(bank *models.Bank) error {
	return r.db.Save(bank).Error
}

func (r *bankRepository) Delete(id uint) error {
	return r.db.Delete(&models.Bank{}, id).Error
}
