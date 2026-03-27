package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ItemService interface {
	// Basic CRUD Operations
	CreateItem(input *input.CreateItemInput, createdBy string) (*output.ItemOutput, error)
	GetItem(id string) (*output.ItemOutput, error)
	GetAllItems(limit, offset int, createdBy string) (*output.ItemListOutput, error)
	UpdateItem(id string, input *input.UpdateItemInput, createdBy string) (*output.ItemOutput, error)
	DeleteItem(id string, createdBy string) error
	GetItemsByType(itemType string, limit, offset int, createdBy string) (*output.ItemListOutput, error)

	// Step 2: Inventory Tracking Operations
	// Enable/disable inventory tracking for items to track stock movements
	EnableInventoryTracking(itemID string) error
	IsInventoryTrackingEnabled(itemID string) (bool, error)

	// Get item with current inventory balance
	// This is used to check current stock levels for purchasing and selling decisions
	GetItemWithInventoryStatus(itemID string) (map[string]interface{}, error)
}

type itemService struct {
	repo             repo.ItemRepository
	vendorRepo       repo.VendorRepository
	ManufacturerRepo repo.ManufacturerRepository
	inventoryRepo    repo.InventoryBalanceRepository
	userRepo         repo.UserRepository
	companyRepo      repo.CompanyRepository
}

func NewItemService(itemRepo repo.ItemRepository, vendorRepo repo.VendorRepository, manufacturerRepo repo.ManufacturerRepository, inventoryRepo repo.InventoryBalanceRepository, userRepo repo.UserRepository, companyRepo repo.CompanyRepository) ItemService {
	return &itemService{
		repo:             itemRepo,
		vendorRepo:       vendorRepo,
		ManufacturerRepo: manufacturerRepo,
		inventoryRepo:    inventoryRepo,
		userRepo:         userRepo,
		companyRepo:      companyRepo,
	}
}

func (s *itemService) CreateItem(input *input.CreateItemInput, createdBy string) (*output.ItemOutput, error) {
	// Validate variant attributes first
	if err := input.ValidateVariantAttributes(); err != nil {
		return nil, err
	}

	id := fmt.Sprintf("item_%s", uuid.New().String()[:8])

	if input.PurchaseInfo != nil && input.PurchaseInfo.PreferredVendorID != nil {
		_, err := s.vendorRepo.FindByID(*input.PurchaseInfo.PreferredVendorID)
		if err != nil {
			return nil, fmt.Errorf("preferred vendor not found")
		}
	}

	itemDetails := buildItemDetails(id, input)

	salesInfo := models.SalesInfo{
		ItemID:       id,
		Account:      input.SalesInfo.Account,
		SellingPrice: input.SalesInfo.SellingPrice,
		Currency:     input.SalesInfo.Currency,
		Description:  input.SalesInfo.Description,
	}

	purchaseInfo := models.PurchaseInfo{}
	if input.PurchaseInfo != nil {
		purchaseInfo = models.PurchaseInfo{
			ItemID:            id,
			Account:           input.PurchaseInfo.Account,
			CostPrice:         input.PurchaseInfo.CostPrice,
			Currency:          input.PurchaseInfo.Currency,
			PreferredVendorID: input.PurchaseInfo.PreferredVendorID,
			Description:       input.PurchaseInfo.Description,
		}
	}

	inventory := models.Inventory{
		ItemID:         id,
		TrackInventory: false,
	}
	if input.Inventory != nil {
		inventory.TrackInventory = input.Inventory.TrackInventory
		inventory.InventoryAccount = input.Inventory.InventoryAccount
		inventory.InventoryValuationMethod = input.Inventory.InventoryValuationMethod
		inventory.ReorderPoint = input.Inventory.ReorderPoint
	}

	returnPolicy := models.ReturnPolicy{
		ItemID:     id,
		Returnable: false,
	}
	if input.ReturnPolicy != nil {
		returnPolicy.Returnable = input.ReturnPolicy.Returnable
	}

	item := &models.Item{
		ID:           id,
		Name:         input.Name,
		Type:         input.Type,
		ItemDetails:  itemDetails,
		SalesInfo:    salesInfo,
		PurchaseInfo: purchaseInfo,
		Inventory:    inventory,
		ReturnPolicy: returnPolicy,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Fetch user details if createdBy is provided
	if createdBy != "" {
		var userID uint
		_, err := fmt.Sscanf(createdBy, "%d", &userID)
		if err == nil {
			user, err := s.userRepo.GetByID(userID)
			if err == nil && user != nil {
				// Set user name from email or username
				if user.Email != nil {
					item.CreatedByUserName = *user.Email
				} else if user.Username != nil {
					item.CreatedByUserName = *user.Username
				}

				// Set company details if user has a company
				if user.CompanyID != nil {
					item.CreatedByCompanyID = *user.CompanyID
					company, err := s.companyRepo.FindByID(*user.CompanyID)
					if err == nil && company != nil {
						item.CreatedByCompanyName = company.CompanyName
					}
				}
			}
		}
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	createdItem, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToItemOutput(createdItem)
}

func buildItemDetails(itemID string, input *input.CreateItemInput) models.ItemDetails {
	itemDetails := models.ItemDetails{
		ItemID:      itemID,
		Structure:   input.ItemDetails.Structure,
		Unit:        input.ItemDetails.Unit,
		SKU:         input.ItemDetails.SKU,
		UPC:         input.ItemDetails.UPC,
		EAN:         input.ItemDetails.EAN,
		MPN:         input.ItemDetails.MPN,
		ISBN:        input.ItemDetails.ISBN,
		Description: input.ItemDetails.Description,
	}

	if input.ItemDetails.Structure == "variants" && len(input.ItemDetails.AttributeDefinitions) > 0 {
		itemDetails.AttributeDefinitions = make(models.AttributeDefinitions, len(input.ItemDetails.AttributeDefinitions))
		for i, attr := range input.ItemDetails.AttributeDefinitions {
			itemDetails.AttributeDefinitions[i] = models.AttributeDefinition{
				Key:     attr.Key,
				Options: attr.Options,
			}
		}
	}

	if len(input.ItemDetails.Variants) > 0 {
		itemDetails.Variants = make([]models.Variant, len(input.ItemDetails.Variants))
		for i, v := range input.ItemDetails.Variants {
			attributes := make([]models.VariantAttribute, 0, len(v.AttributeMap))
			for key, value := range v.AttributeMap {
				attributes = append(attributes, models.VariantAttribute{
					Key:   key,
					Value: value,
				})
			}

			itemDetails.Variants[i] = models.Variant{
				SKU:           v.SKU,
				Attributes:    attributes,
				SellingPrice:  v.SellingPrice,
				CostPrice:     v.CostPrice,
				StockQuantity: v.StockQuantity,
			}
		}
	}

	return itemDetails
}

func (s *itemService) GetItem(id string) (*output.ItemOutput, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToItemOutput(item)
}

func (s *itemService) GetAllItems(limit, offset int, createdBy string) (*output.ItemListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.FindByCreatedBy(createdBy, limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToItemListOutput(items, total)
}

func (s *itemService) UpdateItem(id string, input *input.UpdateItemInput, createdBy string) (*output.ItemOutput, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if item.CreatedBy != createdBy {
		return nil, fmt.Errorf("unauthorized: you can only update items you created")
	}

	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Type != nil {
		item.Type = *input.Type
	}

	if input.SalesInfo != nil {
		if input.SalesInfo.Account != "" {
			item.SalesInfo.Account = input.SalesInfo.Account
		}
		if input.SalesInfo.SellingPrice > 0 {
			item.SalesInfo.SellingPrice = input.SalesInfo.SellingPrice
		}
		if input.SalesInfo.Currency != "" {
			item.SalesInfo.Currency = input.SalesInfo.Currency
		}
		if input.SalesInfo.Description != "" {
			item.SalesInfo.Description = input.SalesInfo.Description
		}
	}

	if input.PurchaseInfo != nil {
		if input.PurchaseInfo.Account != "" {
			item.PurchaseInfo.Account = input.PurchaseInfo.Account
		}
		if input.PurchaseInfo.CostPrice > 0 {
			item.PurchaseInfo.CostPrice = input.PurchaseInfo.CostPrice
		}
		if input.PurchaseInfo.Currency != "" {
			item.PurchaseInfo.Currency = input.PurchaseInfo.Currency
		}
		if input.PurchaseInfo.PreferredVendorID != nil {
			_, err := s.vendorRepo.FindByID(*input.PurchaseInfo.PreferredVendorID)
			if err != nil {
				return nil, fmt.Errorf("preferred vendor not found")
			}
			item.PurchaseInfo.PreferredVendorID = input.PurchaseInfo.PreferredVendorID
		}
		if input.PurchaseInfo.Description != "" {
			item.PurchaseInfo.Description = input.PurchaseInfo.Description
		}
	}

	if input.Inventory != nil {
		item.Inventory.TrackInventory = input.Inventory.TrackInventory
		if input.Inventory.InventoryAccount != "" {
			item.Inventory.InventoryAccount = input.Inventory.InventoryAccount
		}
		if input.Inventory.InventoryValuationMethod != "" {
			item.Inventory.InventoryValuationMethod = input.Inventory.InventoryValuationMethod
		}
		if input.Inventory.ReorderPoint >= 0 {
			item.Inventory.ReorderPoint = input.Inventory.ReorderPoint
		}
	}

	if input.ReturnPolicy != nil {
		item.ReturnPolicy.Returnable = input.ReturnPolicy.Returnable
	}

	if input.ItemDetails != nil {
		if input.ItemDetails.Unit != "" {
			item.ItemDetails.Unit = input.ItemDetails.Unit
		}
		if input.ItemDetails.SKU != "" {
			item.ItemDetails.SKU = input.ItemDetails.SKU
		}
		if input.ItemDetails.Description != "" {
			item.ItemDetails.Description = input.ItemDetails.Description
		}

		if len(input.ItemDetails.Variants) > 0 {
			item.ItemDetails.Variants = make([]models.Variant, len(input.ItemDetails.Variants))
			for i, v := range input.ItemDetails.Variants {
				attributes := make([]models.VariantAttribute, 0, len(v.AttributeMap))
				for key, value := range v.AttributeMap {
					attributes = append(attributes, models.VariantAttribute{
						Key:   key,
						Value: value,
					})
				}

				item.ItemDetails.Variants[i] = models.Variant{
					SKU:           v.SKU,
					Attributes:    attributes,
					SellingPrice:  v.SellingPrice,
					CostPrice:     v.CostPrice,
					StockQuantity: v.StockQuantity,
				}
			}
		}

		if len(input.ItemDetails.AttributeDefinitions) > 0 {
			item.ItemDetails.AttributeDefinitions = make(models.AttributeDefinitions, len(input.ItemDetails.AttributeDefinitions))
			for i, attr := range input.ItemDetails.AttributeDefinitions {
				item.ItemDetails.AttributeDefinitions[i] = models.AttributeDefinition{
					Key:     attr.Key,
					Options: attr.Options,
				}
			}
		}
	}

	item.UpdatedAt = time.Now()

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	updatedItem, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToItemOutput(updatedItem)
}

func (s *itemService) DeleteItem(id string, createdBy string) error {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("item not found")
	}

	if item.CreatedBy != createdBy {
		return errors.New("unauthorized: you can only delete items you created")
	}

	return s.repo.DeleteByCreatedBy(id, createdBy)
}

func (s *itemService) GetItemsByType(itemType string, limit, offset int, createdBy string) (*output.ItemListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.FindByTypeAndCreatedBy(itemType, createdBy, limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToItemListOutput(items, total)
}

// Step 2: Inventory Tracking Operations
// EnableInventoryTracking enables inventory tracking for an item
// Required for items that need stock movement tracking
func (s *itemService) EnableInventoryTracking(itemID string) error {
	item, err := s.repo.FindByID(itemID)
	if err != nil {
		return fmt.Errorf("item not found: %w", err)
	}

	item.Inventory.TrackInventory = true
	item.UpdatedAt = time.Now()

	return s.repo.Update(item)
}

// IsInventoryTrackingEnabled checks if inventory tracking is enabled for an item
func (s *itemService) IsInventoryTrackingEnabled(itemID string) (bool, error) {
	item, err := s.repo.FindByID(itemID)
	if err != nil {
		return false, fmt.Errorf("item not found: %w", err)
	}

	return item.Inventory.TrackInventory, nil
}

// GetItemWithInventoryStatus returns item details along with current inventory balance
// Used to check stock availability for purchasing and selling decisions
func (s *itemService) GetItemWithInventoryStatus(itemID string) (map[string]interface{}, error) {
	item, err := s.repo.FindByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	result := map[string]interface{}{
		"item_id":             itemID,
		"name":                item.Name,
		"type":                item.Type,
		"inventory_tracking":  item.Inventory.TrackInventory,
		"reorder_point":       item.Inventory.ReorderPoint,
		"purchase_price":      item.PurchaseInfo.CostPrice,
		"selling_price":       item.SalesInfo.SellingPrice,
		"preferred_vendor_id": item.PurchaseInfo.PreferredVendorID,
	}

	// If inventory tracking is enabled, get current balance
	if item.Inventory.TrackInventory {
		balance, err := s.inventoryRepo.GetBalance(itemID, nil)
		if err == nil && balance != nil {
			result["current_stock"] = balance.AvailableQuantity
			result["reserved_quantity"] = balance.ReservedQuantity
			result["total_quantity"] = balance.CurrentQuantity
		}
	}

	return result, nil
}
