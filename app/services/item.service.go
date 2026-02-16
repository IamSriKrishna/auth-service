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
	CreateItem(input *input.CreateItemInput) (*output.ItemOutput, error)
	GetItem(id string) (*output.ItemOutput, error)
	GetAllItems(limit, offset int) (*output.ItemListOutput, error)
	UpdateItem(id string, input *input.UpdateItemInput) (*output.ItemOutput, error)
	DeleteItem(id string) error
	GetItemsByType(itemType string, limit, offset int) (*output.ItemListOutput, error)
}

type itemService struct {
	repo             repo.ItemRepository
	vendorRepo       repo.VendorRepository
	ManufacturerRepo repo.ManufacturerRepository
}

func NewItemService(itemRepo repo.ItemRepository, vendorRepo repo.VendorRepository, manufacturerRepo repo.ManufacturerRepository) ItemService {
	return &itemService{
		repo:             itemRepo,
		vendorRepo:       vendorRepo,
		ManufacturerRepo: manufacturerRepo,
	}
}

func (s *itemService) CreateItem(input *input.CreateItemInput) (*output.ItemOutput, error) {
	// Generate item ID
	id := fmt.Sprintf("item_%s", uuid.New().String()[:8])

	// Validate preferred vendor if provided
	if input.PurchaseInfo != nil && input.PurchaseInfo.PreferredVendorID != nil {
		_, err := s.vendorRepo.FindByID(*input.PurchaseInfo.PreferredVendorID)
		if err != nil {
			return nil, fmt.Errorf("preferred vendor not found")
		}
	}
	if input.ManufacturerID != nil {
		_, err := s.ManufacturerRepo.FindByID(*input.ManufacturerID)
		if err != nil {
			return nil, fmt.Errorf("manufacturer not found")
		}
	}

	if input.PurchaseInfo != nil && input.PurchaseInfo.PreferredVendorID != nil {
		_, err := s.vendorRepo.FindByID(*input.PurchaseInfo.PreferredVendorID)
		if err != nil {
			return nil, fmt.Errorf("preferred vendor not found")
		}
	}

	// Validate structure-specific requirements
	if input.ItemDetails.Structure == "single" {
		if len(input.ItemDetails.Variants) > 0 {
			return nil, fmt.Errorf("single items cannot have variants")
		}
	} else if input.ItemDetails.Structure == "variants" {
		if len(input.ItemDetails.Variants) == 0 {
			return nil, fmt.Errorf("variant items must have at least one variant")
		}
		if len(input.ItemDetails.Attributes) == 0 {
			return nil, fmt.Errorf("variant items must define attributes")
		}
	}

	// Build ItemDetails
	itemDetails := buildItemDetails(id, input)

	// Build SalesInfo
	salesInfo := models.SalesInfo{
		ItemID:       id,
		Account:      input.SalesInfo.Account,
		SellingPrice: input.SalesInfo.SellingPrice,
		Currency:     input.SalesInfo.Currency,
		Description:  input.SalesInfo.Description,
	}

	// Build PurchaseInfo (optional)
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

	// Build Inventory (optional, but usually present)
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

	// Build ReturnPolicy
	returnPolicy := models.ReturnPolicy{
		ItemID:     id,
		Returnable: false,
	}
	if input.ReturnPolicy != nil {
		returnPolicy.Returnable = input.ReturnPolicy.Returnable
	}

	// Create the item
	item := &models.Item{
		ID:             id,
		Name:           input.Name,
		Type:           input.Type,
		Brand:          input.Brand,
		ManufacturerID: input.ManufacturerID,
		ItemDetails:    itemDetails,
		SalesInfo:      salesInfo,
		PurchaseInfo:   purchaseInfo,
		Inventory:      inventory,
		ReturnPolicy:   returnPolicy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	// Fetch the created item with all associations
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

	// Add attribute definitions for variant items
	if input.ItemDetails.Structure == "variants" && len(input.ItemDetails.Attributes) > 0 {
		itemDetails.AttributeDefinitions = make(models.AttributeDefinitions, len(input.ItemDetails.Attributes))
		for i, attr := range input.ItemDetails.Attributes {
			itemDetails.AttributeDefinitions[i] = models.AttributeDefinition{
				Key:     attr.Key,
				Options: attr.Options,
			}
		}
	}

	// Add variants
	if len(input.ItemDetails.Variants) > 0 {
		itemDetails.Variants = make([]models.Variant, len(input.ItemDetails.Variants))
		for i, v := range input.ItemDetails.Variants {
			// Convert attribute map to attribute array
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

func (s *itemService) GetAllItems(limit, offset int) (*output.ItemListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToItemListOutput(items, total)
}

func (s *itemService) UpdateItem(id string, input *input.UpdateItemInput) (*output.ItemOutput, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update basic fields
	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Type != nil {
		item.Type = *input.Type
	}
	if input.Brand != nil {
		item.Brand = *input.Brand
	}
	if input.ManufacturerID != nil {
		// Validate manufacturer exists
		_, err := s.ManufacturerRepo.FindByID(*input.ManufacturerID)
		if err != nil {
			return nil, fmt.Errorf("manufacturer not found")
		}
		item.ManufacturerID = input.ManufacturerID
	}

	// Update SalesInfo
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

	// Update PurchaseInfo
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
			// Validate vendor exists
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

	// Update Inventory
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

	// Update ReturnPolicy
	if input.ReturnPolicy != nil {
		item.ReturnPolicy.Returnable = input.ReturnPolicy.Returnable
	}

	// Update ItemDetails
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

		// Update variants if provided
		if len(input.ItemDetails.Variants) > 0 {
			// Convert variants
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

		// Update attribute definitions
		if len(input.ItemDetails.Attributes) > 0 {
			item.ItemDetails.AttributeDefinitions = make(models.AttributeDefinitions, len(input.ItemDetails.Attributes))
			for i, attr := range input.ItemDetails.Attributes {
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

	// Fetch updated item
	updatedItem, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToItemOutput(updatedItem)
}

func (s *itemService) DeleteItem(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("item not found")
	}

	return s.repo.Delete(id)
}

func (s *itemService) GetItemsByType(itemType string, limit, offset int) (*output.ItemListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.FindByType(itemType, limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToItemListOutput(items, total)
}
