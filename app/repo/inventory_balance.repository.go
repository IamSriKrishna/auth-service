package repo

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type inventoryBalanceRepository struct {
	db *gorm.DB
}

func NewInventoryBalanceRepository(db *gorm.DB) InventoryBalanceRepository {
	return &inventoryBalanceRepository{db: db}
}

func (r *inventoryBalanceRepository) GetBalance(itemID string, variantID *uint) (*models.InventoryBalance, error) {
	var balance models.InventoryBalance
	query := r.db.Where("item_id = ?", itemID)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.First(&balance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new balance record if it doesn't exist
			balance = models.InventoryBalance{
				ItemID:              itemID,
				VariantID:           variantID,
				CurrentQuantity:     0,
				ReservedQuantity:    0,
				AvailableQuantity:   0,
				LastInventorySyncAt: time.Now(),
				UpdatedAt:           time.Now(),
			}
			if err := r.db.Create(&balance).Error; err != nil {
				return nil, err
			}
			return &balance, nil
		}
		return nil, err
	}

	return &balance, nil
}

func (r *inventoryBalanceRepository) GetBalances(itemID string) ([]models.InventoryBalance, error) {
	var balances []models.InventoryBalance
	err := r.db.
		Preload("Item").
		Preload("Variant").
		Where("item_id = ?", itemID).
		Find(&balances).Error

	return balances, err
}

func (r *inventoryBalanceRepository) UpdateBalance(balance *models.InventoryBalance) error {
	balance.UpdatedAt = time.Now()
	return r.db.Save(balance).Error
}

func (r *inventoryBalanceRepository) CreateJournalEntry(entry *models.InventoryJournal) error {
	entry.CreatedAt = time.Now()
	return r.db.Create(entry).Error
}

func (r *inventoryBalanceRepository) GetJournalEntries(itemID string, limit, offset int) ([]models.InventoryJournal, int64, error) {
	var entries []models.InventoryJournal
	var total int64

	query := r.db.Where("item_id = ?", itemID)

	if err := query.Model(&models.InventoryJournal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error

	return entries, total, err
}

func (r *inventoryBalanceRepository) ReserveInventory(itemID string, variantID *uint, quantity float64, referenceID, referenceNo string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current balance
		balance, err := r.GetBalance(itemID, variantID)
		if err != nil {
			return err
		}

		// Check if enough inventory is available
		if balance.AvailableQuantity < quantity {
			return fmt.Errorf("insufficient inventory for item %s. Available: %f, Requested: %f", itemID, balance.AvailableQuantity, quantity)
		}

		// Update balance
		balance.ReservedQuantity += quantity
		balance.AvailableQuantity -= quantity
		balance.UpdatedAt = time.Now()

		if err := tx.Save(balance).Error; err != nil {
			return err
		}

		// Create journal entry
		entry := &models.InventoryJournal{
			ItemID:          itemID,
			VariantID:       variantID,
			TransactionType: "SALES_ORDER_RESERVED",
			Quantity:        quantity,
			ReferenceType:   "SalesOrder",
			ReferenceID:     referenceID,
			ReferenceNo:     referenceNo,
			Notes:           fmt.Sprintf("Reserved for sales order %s", referenceNo),
			CreatedAt:       time.Now(),
		}

		if err := tx.Create(entry).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *inventoryBalanceRepository) ReleaseReservation(itemID string, variantID *uint, quantity float64, referenceID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current balance
		balance, err := r.GetBalance(itemID, variantID)
		if err != nil {
			return err
		}

		// Update balance
		if balance.ReservedQuantity >= quantity {
			balance.ReservedQuantity -= quantity
			balance.AvailableQuantity += quantity
			balance.UpdatedAt = time.Now()

			if err := tx.Save(balance).Error; err != nil {
				return err
			}
		}

		// Create journal entry
		entry := &models.InventoryJournal{
			ItemID:          itemID,
			VariantID:       variantID,
			TransactionType: "SALES_ORDER_CANCELLED",
			Quantity:        -quantity,
			ReferenceType:   "SalesOrder",
			ReferenceID:     referenceID,
			Notes:           fmt.Sprintf("Cancelled reservation for sales order %s", referenceID),
			CreatedAt:       time.Now(),
		}

		if err := tx.Create(entry).Error; err != nil {
			return err
		}

		return nil
	})
}
