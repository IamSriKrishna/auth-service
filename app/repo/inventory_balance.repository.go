package repo

import (
	"fmt"
	"log"
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

func (r *inventoryBalanceRepository) GetBalance(itemID string, variantSKU *string) (*models.InventoryBalance, error) {
	var balance models.InventoryBalance
	query := r.db.Where("item_id = ?", itemID)

	variantDesc := "nil"
	if variantSKU != nil {
		variantDesc = *variantSKU
		query = query.Where("variant_sku = ?", *variantSKU)
	} else {
		query = query.Where("variant_sku IS NULL")
	}

	log.Printf("[INVENTORY_BALANCE] GetBalance - itemID: %s, variant: %s", itemID, variantDesc)
	err := query.First(&balance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[INVENTORY_BALANCE] Record not found, creating new balance - itemID: %s, variant: %s", itemID, variantDesc)
			// Create a new balance record if it doesn't exist
			balance = models.InventoryBalance{
				ItemID:              itemID,
				VariantSKU:          variantSKU,
				CurrentQuantity:     0,
				ReservedQuantity:    0,
				AvailableQuantity:   0,
				LastInventorySyncAt: time.Now(),
				UpdatedAt:           time.Now(),
			}
			if err := r.db.Create(&balance).Error; err != nil {
				log.Printf("[INVENTORY_BALANCE] Error creating balance - itemID: %s, variant: %s, err: %v", itemID, variantDesc, err)
				return nil, err
			}
			log.Printf("[INVENTORY_BALANCE] New balance created - ID: %d, itemID: %s, variant: %s", balance.ID, itemID, variantDesc)
			return &balance, nil
		}
		log.Printf("[INVENTORY_BALANCE] Error retrieving balance - itemID: %s, variant: %s, err: %v", itemID, variantDesc, err)
		return nil, err
	}

	log.Printf("[INVENTORY_BALANCE] Existing balance found - ID: %d, Available: %.2f", balance.ID, balance.AvailableQuantity)
	return &balance, nil
}

func (r *inventoryBalanceRepository) GetBalances(itemID string) ([]models.InventoryBalance, error) {
	var balances []models.InventoryBalance
	err := r.db.
		Preload("Item").
		Where("item_id = ?", itemID).
		Find(&balances).Error

	return balances, err
}

func (r *inventoryBalanceRepository) UpdateBalance(balance *models.InventoryBalance) error {
	balance.UpdatedAt = time.Now()
	log.Printf("[INVENTORY_BALANCE] UpdateBalance - ID: %d, Available: %.2f, Current: %.2f", balance.ID, balance.AvailableQuantity, balance.CurrentQuantity)
	result := r.db.Save(balance)
	if result.Error != nil {
		log.Printf("[INVENTORY_BALANCE] Error updating balance - ID: %d, err: %v", balance.ID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("[INVENTORY_BALANCE] WARNING: No rows affected when updating balance - ID: %d", balance.ID)
	} else {
		log.Printf("[INVENTORY_BALANCE] Balance updated successfully - ID: %d, Rows affected: %d", balance.ID, result.RowsAffected)
	}
	return nil
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

func (r *inventoryBalanceRepository) ReserveInventory(itemID string, variantSKU *string, quantity float64, referenceID, referenceNo string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current balance
		balance, err := r.GetBalance(itemID, variantSKU)
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
			VariantSKU:      variantSKU,
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

func (r *inventoryBalanceRepository) ReleaseReservation(itemID string, variantSKU *string, quantity float64, referenceID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current balance
		balance, err := r.GetBalance(itemID, variantSKU)
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
			VariantSKU:      variantSKU,
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
