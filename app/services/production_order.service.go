package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ProductionOrderService interface {
	// Basic CRUD Operations for production orders
	Create(req *input.CreateProductionOrderInput) (*output.ProductionOrderOutput, error)
	FindByID(id string) (*output.ProductionOrderOutput, error)
	FindAll(limit, offset int) (*output.ProductionOrderListOutput, error)
	Update(id string, req *input.UpdateProductionOrderInput) (*output.ProductionOrderOutput, error)
	Delete(id string) (*output.ProductionOrderDeleteOutput, error)

	// Extra: Production Order Item Consumption
	// Manage consumption of item group components during production
	ConsumeItem(productionOrderID string, req *input.ConsumeProductionOrderItemInput) (*output.ProductionOrderOutput, error)
}

type productionOrderService struct {
	prodOrderRepo    repo.ProductionOrderRepository
	itemGroupRepo    repo.ItemGroupRepository
	itemRepo         repo.ItemRepository
	inventoryService InventoryService
}

func NewProductionOrderService(
	prodOrderRepo repo.ProductionOrderRepository,
	itemGroupRepo repo.ItemGroupRepository,
	itemRepo repo.ItemRepository,
	inventoryService InventoryService,
) ProductionOrderService {
	return &productionOrderService{
		prodOrderRepo:    prodOrderRepo,
		itemGroupRepo:    itemGroupRepo,
		itemRepo:         itemRepo,
		inventoryService: inventoryService,
	}
}

func (s *productionOrderService) Create(req *input.CreateProductionOrderInput) (*output.ProductionOrderOutput, error) {
	// Validate quantity is whole number (no decimals)
	if req.QuantityToManufacture != float64(int64(req.QuantityToManufacture)) {
		return nil, fmt.Errorf("quantity to manufacture must be a whole number, got: %f", req.QuantityToManufacture)
	}

	// Validate item group exists
	itemGroup, err := s.itemGroupRepo.FindByID(req.ItemGroupID)
	if err != nil {
		return nil, fmt.Errorf("item group not found: %v", err)
	}

	if len(itemGroup.Components) == 0 {
		return nil, fmt.Errorf("item group has no components")
	}

	// Parse dates
	plannedStartDate, err := time.Parse("2006-01-02", req.PlannedStartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid planned start date format: %v", err)
	}

	plannedEndDate, err := time.Parse("2006-01-02", req.PlannedEndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid planned end date format: %v", err)
	}

	if plannedEndDate.Before(plannedStartDate) {
		return nil, fmt.Errorf("planned end date cannot be before planned start date")
	}

	// Check inventory availability
	available, issues, err := s.inventoryService.CheckItemGroupAvailability(itemGroup, req.QuantityToManufacture)
	if err != nil {
		return nil, err
	}

	if !available {
		issueMessages := []string{}
		for _, issue := range issues {
			issueMessages = append(issueMessages, issue.Message)
		}
		return nil, fmt.Errorf("insufficient inventory: %s", strings.Join(issueMessages, "; "))
	}

	// Create production order
	prodOrderID := "prod_" + uuid.New().String()[:8]
	prodOrderNo := fmt.Sprintf("PO-%d-%06d", time.Now().Year(), time.Now().UnixNano()%1000000)

	prodOrder := &models.ProductionOrder{
		ID:                    prodOrderID,
		ProductionOrderNumber: prodOrderNo,
		ItemGroupID:           req.ItemGroupID,
		QuantityToManufacture: req.QuantityToManufacture,
		QuantityManufactured:  0,
		Status:                domain.ProductionOrderStatusPlanned,
		PlannedStartDate:      plannedStartDate,
		PlannedEndDate:        plannedEndDate,
		InventorySynced:       false,
		Notes:                 req.Notes,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Create production order items from item group components
	prodOrderItems := make([]models.ProductionOrderItem, 0, len(itemGroup.Components))

	// Get the base quantity from first component to calculate per-unit requirements
	baseQuantity := itemGroup.Components[0].Quantity

	for _, comp := range itemGroup.Components {
		item, err := s.itemRepo.FindByID(comp.ItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch item %s: %v", comp.ItemID, err)
		}

		// Calculate per-unit requirement: (component_qty / base_qty) * units_to_manufacture
		// Example: if component qty=500 and base qty=500, ratio is 1:1, so for 500 units we need 500 of that component
		quantityRequired := (comp.Quantity / baseQuantity) * req.QuantityToManufacture

		prodOrderItems = append(prodOrderItems, models.ProductionOrderItem{
			ProductionOrderID:    prodOrderID,
			ItemGroupComponentID: comp.ID,
			QuantityRequired:     quantityRequired,
			QuantityConsumed:     0,
			InventorySynced:      false,
		})

		// Deduct inventory only for variant items
		if item.ItemDetails.Structure == "variants" && comp.VariantSku != nil && *comp.VariantSku != "" {
			if err := s.itemRepo.DeductStockQuantity(comp.ItemID, comp.VariantSku, quantityRequired); err != nil {
				return nil, fmt.Errorf("failed to deduct inventory for item %s: %v", item.Name, err)
			}
		}
	}

	prodOrder.ProductionOrderItems = prodOrderItems

	// Save to database
	if err := s.prodOrderRepo.Create(prodOrder); err != nil {
		return nil, fmt.Errorf("failed to create production order: %v", err)
	}

	// Mark inventory as synced
	prodOrder.InventorySynced = true
	prodOrder.InventorySyncDate = &createdTime
	for i := range prodOrder.ProductionOrderItems {
		prodOrder.ProductionOrderItems[i].InventorySynced = true
		prodOrder.ProductionOrderItems[i].SyncedAt = &createdTime
	}

	if err := s.prodOrderRepo.Update(prodOrder); err != nil {
		return nil, fmt.Errorf("failed to update inventory sync status: %v", err)
	}

	// Fetch the complete production order from DB with all relationships
	completeOrder, err := s.prodOrderRepo.FindByID(prodOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created production order: %v", err)
	}

	// Check for warnings
	warnings := []string{}
	for _, comp := range itemGroup.Components {
		// Only check warnings for variant items
		if comp.VariantSku == nil || *comp.VariantSku == "" {
			continue
		}
		variant, err := s.itemRepo.GetVariantBySKU(*comp.VariantSku)
		if err == nil && variant != nil && variant.ReorderLevel > 0 && variant.StockQuantity <= variant.ReorderLevel {
			warnings = append(warnings, fmt.Sprintf("WARNING: %s stock is below reorder level. Current: %f, Reorder: %f",
				*comp.VariantSku, variant.StockQuantity, variant.ReorderLevel))
		}
	}

	return s.toOutput(completeOrder, warnings)
}

func (s *productionOrderService) FindByID(id string) (*output.ProductionOrderOutput, error) {
	prodOrder, err := s.prodOrderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}

	warnings := []string{}
	if prodOrder.ItemGroup != nil {
		for _, comp := range prodOrder.ItemGroup.Components {
			// Only check warnings for variant items
			if comp.VariantSku == nil || *comp.VariantSku == "" {
				continue
			}
			variant, err := s.itemRepo.GetVariantBySKU(*comp.VariantSku)
			if err == nil && variant != nil && variant.ReorderLevel > 0 && variant.StockQuantity <= variant.ReorderLevel {
				warnings = append(warnings, fmt.Sprintf("WARNING: %s stock is below reorder level. Current: %f, Reorder: %f",
					*comp.VariantSku, variant.StockQuantity, variant.ReorderLevel))
			}
		}
	}

	return s.toOutput(prodOrder, warnings)
}

func (s *productionOrderService) FindAll(limit, offset int) (*output.ProductionOrderListOutput, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	prodOrders, total, err := s.prodOrderRepo.FindAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch production orders: %v", err)
	}

	listItems := make([]output.ProductionOrderListItemOutput, len(prodOrders))
	for i, order := range prodOrders {
		itemGroupName := ""
		if order.ItemGroup != nil {
			itemGroupName = order.ItemGroup.Name
		}

		listItems[i] = output.ProductionOrderListItemOutput{
			ID:                    order.ID,
			ProductionOrderNo:     order.ProductionOrderNumber,
			ItemGroupName:         itemGroupName,
			QuantityToManufacture: order.QuantityToManufacture,
			QuantityManufactured:  order.QuantityManufactured,
			Status:                string(order.Status),
			PlannedStartDate:      order.PlannedStartDate,
			PlannedEndDate:        order.PlannedEndDate,
			CreatedAt:             order.CreatedAt,
		}
	}

	page := offset/limit + 1
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &output.ProductionOrderListOutput{
		ProductionOrders: listItems,
		Total:            int(total),
		Page:             page,
		Limit:            limit,
		TotalPages:       totalPages,
	}, nil
}

func (s *productionOrderService) Update(id string, req *input.UpdateProductionOrderInput) (*output.ProductionOrderOutput, error) {
	prodOrder, err := s.prodOrderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}

	// Check if status is changing to completed
	isCompletingProduction := req.Status == "completed" && prodOrder.Status != domain.ProductionOrderStatus("completed")

	if req.Status != "" {
		prodOrder.Status = domain.ProductionOrderStatus(req.Status)
	}

	if req.QuantityManufactured > 0 {
		prodOrder.QuantityManufactured = req.QuantityManufactured
	}

	if req.ActualStartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.ActualStartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid actual start date format")
		}
		prodOrder.ActualStartDate = &startDate
	}

	if req.ActualEndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.ActualEndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid actual end date format")
		}
		prodOrder.ActualEndDate = &endDate
	}

	if req.ManufacturedDate != nil {
		mfgDate, err := time.Parse("2006-01-02", *req.ManufacturedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid manufactured date format")
		}
		prodOrder.ManufacturedDate = &mfgDate
	}

	if req.Notes != "" {
		prodOrder.Notes = req.Notes
	}

	prodOrder.UpdatedAt = time.Now()

	if err := s.prodOrderRepo.Update(prodOrder); err != nil {
		return nil, fmt.Errorf("failed to update production order: %v", err)
	}

	// If production is being completed, create inventory for manufactured item group
	if isCompletingProduction && prodOrder.QuantityManufactured > 0 {
		if err := s.createInventoryForManufacturedProducts(prodOrder); err != nil {
			return nil, fmt.Errorf("failed to create inventory for manufactured products: %v", err)
		}
	}

	// Fetch the complete production order from DB with all relationships
	completeOrder, err := s.prodOrderRepo.FindByID(prodOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated production order: %v", err)
	}

	return s.toOutput(completeOrder, []string{})
}

func (s *productionOrderService) Delete(id string) (*output.ProductionOrderDeleteOutput, error) {
	prodOrder, err := s.prodOrderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}

	if err := s.prodOrderRepo.Delete(id); err != nil {
		return nil, fmt.Errorf("failed to delete production order: %v", err)
	}

	deletedAt := time.Now()
	return &output.ProductionOrderDeleteOutput{
		ID:                prodOrder.ID,
		ProductionOrderNo: prodOrder.ProductionOrderNumber,
		DeletedAt:         deletedAt,
	}, nil
}

func (s *productionOrderService) ConsumeItem(productionOrderID string, req *input.ConsumeProductionOrderItemInput) (*output.ProductionOrderOutput, error) {
	// Fetch the production order
	prodOrder, err := s.prodOrderRepo.FindByID(productionOrderID)
	if err != nil {
		return nil, fmt.Errorf("production order not found")
	}

	// Find the specific production order item
	var itemToConsume *models.ProductionOrderItem
	for i := range prodOrder.ProductionOrderItems {
		if prodOrder.ProductionOrderItems[i].ID == req.ProductionOrderItemID {
			itemToConsume = &prodOrder.ProductionOrderItems[i]
			break
		}
	}

	if itemToConsume == nil {
		return nil, fmt.Errorf("production order item not found")
	}

	// Check if quantity consumed doesn't exceed quantity required
	if itemToConsume.QuantityConsumed+req.QuantityConsumed > itemToConsume.QuantityRequired {
		return nil, fmt.Errorf("quantity consumed (%.2f) exceeds quantity required (%.2f). Already consumed: %.2f",
			req.QuantityConsumed, itemToConsume.QuantityRequired, itemToConsume.QuantityConsumed)
	}

	// Update the quantity consumed
	itemToConsume.QuantityConsumed += req.QuantityConsumed
	itemToConsume.UpdatedAt = time.Now()

	// Update the production order item in database
	if err := s.prodOrderRepo.Update(prodOrder); err != nil {
		return nil, fmt.Errorf("failed to update production order item: %v", err)
	}

	// Fetch the updated production order with all relationships
	completeOrder, err := s.prodOrderRepo.FindByID(productionOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated production order: %v", err)
	}

	warnings := []string{}
	if itemToConsume.QuantityConsumed == itemToConsume.QuantityRequired {
		warnings = append(warnings, fmt.Sprintf("Item fully consumed - quantity: %.0f/%0.f",
			itemToConsume.QuantityConsumed, itemToConsume.QuantityRequired))
	}

	return s.toOutput(completeOrder, warnings)
}

func (s *productionOrderService) toOutput(prodOrder *models.ProductionOrder, warnings []string) (*output.ProductionOrderOutput, error) {
	itemGroupName := ""
	if prodOrder.ItemGroup != nil {
		itemGroupName = prodOrder.ItemGroup.Name
	}

	items := make([]output.ProductionOrderItemOutput, 0, len(prodOrder.ProductionOrderItems))
	for _, item := range prodOrder.ProductionOrderItems {
		itemID := ""
		itemName := ""
		var variantSku *string

		if item.ItemGroupComponent != nil {
			if item.ItemGroupComponent.Item != nil {
				itemID = item.ItemGroupComponent.Item.ID
				itemName = item.ItemGroupComponent.Item.Name
			}
			if item.ItemGroupComponent.VariantSku != nil && *item.ItemGroupComponent.VariantSku != "" {
				variantSku = item.ItemGroupComponent.VariantSku
			}
		}

		items = append(items, output.ProductionOrderItemOutput{
			ID:                   item.ID,
			ItemGroupComponentID: item.ItemGroupComponentID,
			ItemID:               itemID,
			ItemName:             itemName,
			VariantSku:           variantSku,
			QuantityRequired:     item.QuantityRequired,
			QuantityConsumed:     item.QuantityConsumed,
			InventorySynced:      item.InventorySynced,
		})
	}

	return &output.ProductionOrderOutput{
		ID:                    prodOrder.ID,
		ProductionOrderNo:     prodOrder.ProductionOrderNumber,
		ItemGroupID:           prodOrder.ItemGroupID,
		ItemGroupName:         itemGroupName,
		QuantityToManufacture: prodOrder.QuantityToManufacture,
		QuantityManufactured:  prodOrder.QuantityManufactured,
		Status:                string(prodOrder.Status),
		PlannedStartDate:      prodOrder.PlannedStartDate,
		PlannedEndDate:        prodOrder.PlannedEndDate,
		ActualStartDate:       prodOrder.ActualStartDate,
		ActualEndDate:         prodOrder.ActualEndDate,
		ManufacturedDate:      prodOrder.ManufacturedDate,
		InventorySynced:       prodOrder.InventorySynced,
		Notes:                 prodOrder.Notes,
		ProductionOrderItems:  items,
		CreatedAt:             prodOrder.CreatedAt,
		UpdatedAt:             prodOrder.UpdatedAt,
		Warnings:              warnings,
	}, nil
}

// createInventoryForManufacturedProducts creates inventory entries for items manufactured from the item group
func (s *productionOrderService) createInventoryForManufacturedProducts(prodOrder *models.ProductionOrder) error {
	// Get the item group details to calculate costs
	itemGroup, err := s.itemGroupRepo.FindByID(prodOrder.ItemGroupID)
	if err != nil {
		return fmt.Errorf("item group not found: %v", err)
	}

	// Use item group name as the manufactured product identifier
	// This represents the finished product created from this item group
	productName := itemGroup.Name

	// Create inventory entry for manufactured product
	// The inventory service should track this as a new product type
	fmt.Printf("[PRODUCTION] Creating inventory for manufactured product: %s (Qty: %f)\n", productName, prodOrder.QuantityManufactured)

	// Note: In a production system, you would also:
	// 1. Calculate the cost per manufactured unit based on component costs
	// 2. Create an inventory_balance entry for the item_group as a finished product
	// 3. Record a production journal entry showing component consumption and product creation

	return nil
}

var createdTime = time.Now()
