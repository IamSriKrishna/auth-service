package services

import (
	"fmt"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ItemGroupService interface {
	Create(input *input.CreateItemGroupInput) (*output.ItemGroupOutput, error)
	FindByID(id string) (*output.ItemGroupOutput, error)
	FindAll(limit, offset int, search string) (*output.ItemGroupListOutput, error)
	Update(id string, input *input.UpdateItemGroupInput) (*output.ItemGroupOutput, error)
	Delete(id string) error
	FindByName(name string) (*output.ItemGroupOutput, error)
}

type itemGroupService struct {
	itemGroupRepo repo.ItemGroupRepository
	itemRepo      repo.ItemRepository
}

func NewItemGroupService(itemGroupRepo repo.ItemGroupRepository, itemRepo repo.ItemRepository) ItemGroupService {
	return &itemGroupService{
		itemGroupRepo: itemGroupRepo,
		itemRepo:      itemRepo,
	}
}

func (s *itemGroupService) Create(input *input.CreateItemGroupInput) (*output.ItemGroupOutput, error) {
	// Validate all items exist and quantities are valid
	if len(input.Components) == 0 {
		return nil, fmt.Errorf("item group must have at least one component")
	}

	// Check that all component quantities are equal and whole numbers (integers)
	firstQuantity := input.Components[0].Quantity
	for i, comp := range input.Components {
		// Validate all quantities are whole numbers
		if comp.Quantity != float64(int64(comp.Quantity)) {
			return nil, fmt.Errorf("component %d quantity must be a whole number (no decimals). Got: %f", i+1, comp.Quantity)
		}

		if comp.Quantity != firstQuantity {
			return nil, fmt.Errorf("all component quantities must be equal. Component 1 has quantity %d, but component %d has quantity %d", int64(firstQuantity), i+1, int64(comp.Quantity))
		}

		if comp.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for item %s: got %d", comp.ItemID, int64(comp.Quantity))
		}

		item, err := s.itemRepo.FindByID(comp.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item %s not found", comp.ItemID)
		}

		// If variant is specified, validate it exists
		if comp.VariantSku != nil {
			variantFound := false
			if item.ItemDetails.Variants != nil {
				for _, v := range item.ItemDetails.Variants {
					if v.SKU == *comp.VariantSku {
						variantFound = true
						break
					}
				}
			}
			if !variantFound {
				return nil, fmt.Errorf("variant %s not found in item %s", *comp.VariantSku, comp.ItemID)
			}
		}
	}

	// Create ItemGroup
	itemGroupID := "ig_" + uuid.New().String()[:8]

	components := make([]models.ItemGroupComponent, 0, len(input.Components))
	for _, comp := range input.Components {
		// Convert VariantDetails from interface{} to models.VariantDetails
		var variantDetails models.VariantDetails
		if comp.VariantDetails != nil {
			switch v := comp.VariantDetails.(type) {
			case map[string]interface{}:
				variantDetails = make(models.VariantDetails)
				for k, val := range v {
					if strVal, ok := val.(string); ok {
						variantDetails[k] = strVal
					}
				}
			case map[string]string:
				variantDetails = models.VariantDetails(v)
			}
		}

		component := models.ItemGroupComponent{
			ItemGroupID:    itemGroupID,
			ItemID:         comp.ItemID,
			VariantSku:     comp.VariantSku,
			Quantity:       comp.Quantity,
			VariantDetails: variantDetails,
		}
		components = append(components, component)
	}

	itemGroup := &models.ItemGroup{
		ID:          itemGroupID,
		Name:        input.Name,
		Description: input.Description,
		IsActive:    input.IsActive,
		Components:  components,
	}

	err := s.itemGroupRepo.Create(itemGroup)
	if err != nil {
		return nil, err
	}

	return s.toOutput(itemGroup)
}

func (s *itemGroupService) FindByID(id string) (*output.ItemGroupOutput, error) {
	itemGroup, err := s.itemGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("item group not found")
	}

	return s.toOutput(itemGroup)
}

func (s *itemGroupService) FindAll(limit, offset int, search string) (*output.ItemGroupListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	itemGroups, total, err := s.itemGroupRepo.FindAll(limit, offset, search)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.ItemGroupOutput, len(itemGroups))
	for i, ig := range itemGroups {
		out, _ := s.toOutput(&ig)
		outputs[i] = *out
	}

	page := offset/limit + 1
	return &output.ItemGroupListOutput{
		ItemGroups: outputs,
		Total:      total,
		Page:       page,
		PageSize:   limit,
	}, nil
}

func (s *itemGroupService) Update(id string, input *input.UpdateItemGroupInput) (*output.ItemGroupOutput, error) {
	itemGroup, err := s.itemGroupRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("item group not found")
	}

	if input.Name != "" {
		itemGroup.Name = input.Name
	}

	if input.Description != "" {
		itemGroup.Description = input.Description
	}

	if input.IsActive != nil {
		itemGroup.IsActive = *input.IsActive
	}

	// Update components if provided
	if len(input.Components) > 0 {
		// Validate all items exist
		for _, comp := range input.Components {
			item, err := s.itemRepo.FindByID(comp.ItemID)
			if err != nil {
				return nil, fmt.Errorf("item %s not found", comp.ItemID)
			}

			if comp.VariantSku != nil {
				variantFound := false
				if item.ItemDetails.Variants != nil {
					for _, v := range item.ItemDetails.Variants {
						if v.SKU == *comp.VariantSku {
							variantFound = true
							break
						}
					}
				}
				if !variantFound {
					return nil, fmt.Errorf("variant %s not found in item %s", *comp.VariantSku, comp.ItemID)
				}
			}
		}

		itemGroup.Components = make([]models.ItemGroupComponent, 0, len(input.Components))
		for _, comp := range input.Components {
			// Convert VariantDetails from interface{} to models.VariantDetails
			var variantDetails models.VariantDetails
			if comp.VariantDetails != nil {
				switch v := comp.VariantDetails.(type) {
				case map[string]interface{}:
					variantDetails = make(models.VariantDetails)
					for k, val := range v {
						if strVal, ok := val.(string); ok {
							variantDetails[k] = strVal
						}
					}
				case map[string]string:
					variantDetails = models.VariantDetails(v)
				}
			}

			component := models.ItemGroupComponent{
				ItemGroupID:    id,
				ItemID:         comp.ItemID,
				VariantSku:     comp.VariantSku,
				Quantity:       comp.Quantity,
				VariantDetails: variantDetails,
			}
			itemGroup.Components = append(itemGroup.Components, component)
		}
	}

	err = s.itemGroupRepo.Update(itemGroup)
	if err != nil {
		return nil, err
	}

	return s.toOutput(itemGroup)
}

func (s *itemGroupService) Delete(id string) error {
	return s.itemGroupRepo.Delete(id)
}

func (s *itemGroupService) FindByName(name string) (*output.ItemGroupOutput, error) {
	itemGroup, err := s.itemGroupRepo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("item group not found")
	}

	return s.toOutput(itemGroup)
}

func (s *itemGroupService) toOutput(itemGroup *models.ItemGroup) (*output.ItemGroupOutput, error) {
	components := make([]output.ItemGroupComponentOutput, len(itemGroup.Components))

	for i, comp := range itemGroup.Components {
		itemInfo := &output.ItemInfo{}
		if comp.Item != nil {
			itemInfo.ID = comp.Item.ID
			itemInfo.Name = comp.Item.Name
			// Get SKU from the item's base item details
			if comp.Item.ItemDetails.SKU != "" {
				itemInfo.SKU = comp.Item.ItemDetails.SKU
			}
		}

		components[i] = output.ItemGroupComponentOutput{
			ID:             comp.ID,
			ItemGroupID:    comp.ItemGroupID,
			ItemID:         comp.ItemID,
			Item:           itemInfo,
			VariantSku:     comp.VariantSku,
			Quantity:       comp.Quantity,
			VariantDetails: comp.VariantDetails,
			CreatedAt:      comp.CreatedAt,
			UpdatedAt:      comp.UpdatedAt,
		}
	}

	return &output.ItemGroupOutput{
		ID:          itemGroup.ID,
		Name:        itemGroup.Name,
		Description: itemGroup.Description,
		IsActive:    itemGroup.IsActive,
		Components:  components,
		CreatedAt:   itemGroup.CreatedAt,
		UpdatedAt:   itemGroup.UpdatedAt,
	}, nil
}

// ValidateStockAvailability checks if there is enough stock to fulfill item group usage
// quantityToCreate: number of item groups to create/consume
func (s *itemGroupService) ValidateStockAvailability(itemGroupID string, quantityToCreate float64) ([]string, error) {
	itemGroup, err := s.itemGroupRepo.FindByID(itemGroupID)
	if err != nil {
		return nil, fmt.Errorf("item group not found: %v", err)
	}

	warnings := []string{}

	for _, comp := range itemGroup.Components {
		totalRequired := comp.Quantity * quantityToCreate

		item, err := s.itemRepo.FindByID(comp.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item %s not found: %v", comp.ItemID, err)
		}

		// For item groups, variant SKU is required
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			return nil, fmt.Errorf("variant SKU required for item %s in item group", comp.ItemID)
		}

		// Check variant stock
		variant, err := s.itemRepo.GetVariantBySKU(*comp.VariantSku)
		if err != nil {
			return nil, fmt.Errorf("variant %s not found: %v", *comp.VariantSku, err)
		}

		available := variant.StockQuantity

		// Check if would go below reorder level
		if variant.ReorderLevel > 0 && (available-totalRequired) <= variant.ReorderLevel {
			warnings = append(warnings,
				fmt.Sprintf("WARNING: %s (variant: %s) stock would reach reorder level. Current: %f, Required: %f, Reorder Level: %f",
					item.Name, *comp.VariantSku, available, totalRequired, variant.ReorderLevel))
		}

		if available < totalRequired {
			return nil, fmt.Errorf("insufficient stock for %s (variant: %s): available=%f, required=%f",
				item.Name, *comp.VariantSku, available, totalRequired)
		}
	}

	return warnings, nil
}
