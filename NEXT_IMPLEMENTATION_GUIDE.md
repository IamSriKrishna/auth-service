# Next Steps: Repository Interfaces & DTO Templates

This file provides templates for the next implementation phase.

---

## Repository Interfaces

Create these in `app/repo/interfaces.go` or separate interface files:

### ItemGroupRepository Interface

```go
package repo

import (
	"context"
	"github.com/bbapp-org/auth-service/app/models"
)

type ItemGroupRepository interface {
	// Create new ItemGroup
	Create(ctx context.Context, itemGroup *models.ItemGroup) error
	
	// Get ItemGroup by ID with components
	GetByID(ctx context.Context, id string) (*models.ItemGroup, error)
	
	// List all ItemGroups with pagination
	List(ctx context.Context, page, pageSize int) ([]models.ItemGroup, int64, error)
	
	// List active ItemGroups only
	ListActive(ctx context.Context), ([]models.ItemGroup, error)
	
	// Search ItemGroups by name
	Search(ctx context.Context, keyword string) ([]models.ItemGroup, error)
	
	// Update ItemGroup
	Update(ctx context.Context, itemGroup *models.ItemGroup) error
	
	// Soft delete ItemGroup
	Delete(ctx context.Context, id string) error
	
	// Get ItemGroup with all components loaded
	GetWithComponents(ctx context.Context, id string) (*models.ItemGroup, error)
}
```

### ProductionOrderRepository Interface

```go
package repo

type ProductionOrderRepository interface {
	// Create new ProductionOrder
	Create(ctx context.Context, order *models.ProductionOrder) error
	
	// Get ProductionOrder by ID with items
	GetByID(ctx context.Context, id string) (*models.ProductionOrder, error)
	
	// Get by ProductionOrderNumber
	GetByNumber(ctx context.Context, number string) (*models.ProductionOrder, error)
	
	// List all ProductionOrders
	List(ctx context.Context, page, pageSize int) ([]models.ProductionOrder, int64, error)
	
	// Filter by status
	ListByStatus(ctx context.Context, status string) ([]models.ProductionOrder, error)
	
	// Filter by ItemGroup
	ListByItemGroup(ctx context.Context, itemGroupID string) ([]models.ProductionOrder, error)
	
	// Filter by date range
	ListByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.ProductionOrder, error)
	
	// Update status
	UpdateStatus(ctx context.Context, id string, status string) error
	
	// Update manufactured quantity
	UpdateManufacturedQuantity(ctx context.Context, id string, quantity float64) error
	
	// Update with completion date
	MarkCompleted(ctx context.Context, id string) error
	
	// Get with all related items
	GetWithItems(ctx context.Context, id string) (*models.ProductionOrder, error)
}
```

### InventoryRepository Interface

```go
package repo

type InventoryRepository interface {
	// Balance Operations
	GetBalance(ctx context.Context, itemID string, variantID *uint) (*models.InventoryBalance, error)
	UpdateBalance(ctx context.Context, balance *models.InventoryBalance) error
	CreateBalance(ctx context.Context, balance *models.InventoryBalance) error
	
	// Aggregation Operations
	GetAggregation(ctx context.Context, itemID string, variantID *uint) (*models.InventoryAggregation, error)
	UpdateAggregation(ctx context.Context, agg *models.InventoryAggregation) error
	RecalculateAggregation(ctx context.Context, itemID string, variantID *uint) error
	
	// Journal Operations (Audit Trail)
	LogTransaction(ctx context.Context, journal *models.InventoryJournal) error
	GetJournal(ctx context.Context, itemID string, startDate, endDate time.Time) ([]models.InventoryJournal, error)
	GetJournalByReference(ctx context.Context, refType, refID string) ([]models.InventoryJournal, error)
	
	// Supply Chain Summary
	GetSupplyChainSummary(ctx context.Context, itemID string) (*models.SupplyChainSummary, error)
	UpdateSupplyChainSummary(ctx context.Context, summary *models.SupplyChainSummary) error
	RecalculateSupplyChainSummary(ctx context.Context, itemID string) error
}
```

---

## DTO Input/Output Templates

### Create `app/dto/input/item_group.input.go`

```go
package input

import "github.com/bbapp-org/auth-service/app/dto/output"

type CreateItemGroupInput struct {
	Name        string                      `json:"name" binding:"required"`
	Description string                      `json:"description"`
	IsActive    bool                        `json:"is_active" binding:"required"`
	Components  []ItemGroupComponentInput   `json:"components" binding:"required,min=1"`
}

type UpdateItemGroupInput struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	IsActive    *bool                       `json:"is_active"`
	Components  []ItemGroupComponentInput   `json:"components"`
}

type ItemGroupComponentInput struct {
	ItemID         string                 `json:"item_id" binding:"required"`
	VariantID      *uint                  `json:"variant_id"`
	Quantity       float64                `json:"quantity" binding:"required,gt=0"`
	VariantDetails map[string]string      `json:"variant_details"`
}

type ListItemGroupInput struct {
	Page     int    `query:"page" binding:"required,min=1"`
	PageSize int    `query:"page_size" binding:"required,min=1,max=100"`
	Search   string `query:"search"`
	IsActive *bool  `query:"is_active"`
	Sort     string `query:"sort"`
}
```

### Create `app/dto/output/item_group.output.go`

```go
package output

import "time"

type ItemGroupOutput struct {
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	Description string                          `json:"description"`
	IsActive    bool                            `json:"is_active"`
	Components  []ItemGroupComponentOutput      `json:"components"`
	CreatedAt   time.Time                       `json:"created_at"`
	UpdatedAt   time.Time                       `json:"updated_at"`
}

type ItemGroupComponentOutput struct {
	ID             uint                   `json:"id"`
	ItemGroupID    string                 `json:"item_group_id"`
	ItemID         string                 `json:"item_id"`
	Item           *ItemOutput            `json:"item,omitempty"`
	VariantID      *uint                  `json:"variant_id"`
	Variant        *VariantOutput         `json:"variant,omitempty"`
	Quantity       float64                `json:"quantity"`
	VariantDetails map[string]string      `json:"variant_details"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type ListItemGroupOutput struct {
	Data      []ItemGroupOutput `json:"data"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}
```

### Create `app/dto/input/production_order.input.go`

```go
package input

import "time"

type CreateProductionOrderInput struct {
	ItemGroupID           string    `json:"item_group_id" binding:"required"`
	QuantityToManufacture float64   `json:"quantity_to_manufacture" binding:"required,gt=0"`
	PlannedStartDate      time.Time `json:"planned_start_date" binding:"required"`
	PlannedEndDate        time.Time `json:"planned_end_date" binding:"required"`
	Notes                 string    `json:"notes"`
}

type UpdateProductionOrderInput struct {
	QuantityToManufacture *float64  `json:"quantity_to_manufacture"`
	PlannedStartDate      *time.Time `json:"planned_start_date"`
	PlannedEndDate        *time.Time `json:"planned_end_date"`
	ActualStartDate       *time.Time `json:"actual_start_date"`
	ActualEndDate         *time.Time `json:"actual_end_date"`
	Notes                 string    `json:"notes"`
}

type UpdateProductionStatusInput struct {
	Status                string     `json:"status" binding:"required,oneof=planned in_progress completed cancelled"`
	QuantityManufactured  *float64   `json:"quantity_manufactured"`
	ActualEndDate         *time.Time `json:"actual_end_date"`
}

type ListProductionOrderInput struct {
	Page      int       `query:"page" binding:"required,min=1"`
	PageSize  int       `query:"page_size" binding:"required,min=1,max=100"`
	Status    string    `query:"status"`
	ItemGroup string    `query:"item_group_id"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
	Sort      string    `query:"sort"`
}
```

### Create `app/dto/output/production_order.output.go`

```go
package output

import "time"

type ProductionOrderOutput struct {
	ID                      string                          `json:"id"`
	ProductionOrderNumber   string                          `json:"production_order_no"`
	ItemGroupID             string                          `json:"item_group_id"`
	ItemGroup               *ItemGroupOutput                `json:"item_group,omitempty"`
	QuantityToManufacture   float64                         `json:"quantity_to_manufacture"`
	QuantityManufactured    float64                         `json:"quantity_manufactured"`
	Status                  string                          `json:"status"`
	PlannedStartDate        time.Time                       `json:"planned_start_date"`
	PlannedEndDate          time.Time                       `json:"planned_end_date"`
	ActualStartDate         *time.Time                      `json:"actual_start_date"`
	ActualEndDate           *time.Time                      `json:"actual_end_date"`
	Notes                   string                          `json:"notes"`
	ProductionOrderItems    []ProductionOrderItemOutput     `json:"production_order_items"`
	CreatedAt               time.Time                       `json:"created_at"`
	UpdatedAt               time.Time                       `json:"updated_at"`
	CreatedBy               string                          `json:"created_by"`
	UpdatedBy               string                          `json:"updated_by"`
}

type ProductionOrderItemOutput struct {
	ID                   uint                       `json:"id"`
	ProductionOrderID    string                     `json:"production_order_id"`
	ItemGroupComponentID uint                       `json:"item_group_component_id"`
	QuantityRequired     float64                    `json:"quantity_required"`
	QuantityConsumed     float64                    `json:"quantity_consumed"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
}

type ListProductionOrderOutput struct {
	Data      []ProductionOrderOutput `json:"data"`
	Total     int64                   `json:"total"`
	Page      int                     `json:"page"`
	PageSize  int                     `json:"page_size"`
	TotalPage int                     `json:"total_page"`
}
```

### Create `app/dto/output/inventory.output.go`

```go
package output

import "time"

type InventoryBalanceOutput struct {
	ItemID           string     `json:"item_id"`
	VariantID        *uint      `json:"variant_id"`
	Item             *ItemOutput `json:"item,omitempty"`
	Variant          *VariantOutput `json:"variant,omitempty"`
	CurrentQuantity  float64    `json:"current_quantity"`
	ReservedQuantity float64    `json:"reserved_quantity"`
	AvailableQuantity float64   `json:"available_quantity"`
	LastReceivedDate *time.Time `json:"last_received_date"`
	LastConsumedDate *time.Time `json:"last_consumed_date"`
	LastSoldDate     *time.Time `json:"last_sold_date"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type InventoryAggregationOutput struct {
	ItemID               string     `json:"item_id"`
	VariantID            *uint      `json:"variant_id"`
	Item                 *ItemOutput `json:"item,omitempty"`
	Variant              *VariantOutput `json:"variant,omitempty"`
	TotalPurchased       float64    `json:"total_purchased"`
	TotalManufactured    float64    `json:"total_manufactured"`
	TotalConsumedInMfg   float64    `json:"total_consumed_in_mfg"`
	TotalSold            float64    `json:"total_sold"`
	CalculatedAt         time.Time  `json:"calculated_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type InventoryJournalOutput struct {
	ID              uint      `json:"id"`
	ItemID          string    `json:"item_id"`
	VariantID       *uint     `json:"variant_id"`
	TransactionType string    `json:"transaction_type"`
	Quantity        float64   `json:"quantity"`
	ReferenceType   string    `json:"reference_type"`
	ReferenceID     string    `json:"reference_id"`
	ReferenceNo     string    `json:"reference_no"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
}

type SupplyChainSummaryOutput struct {
	ItemID                     string     `json:"item_id"`
	VariantID                  *uint      `json:"variant_id"`
	Item                       *ItemOutput `json:"item,omitempty"`
	OpeningStock               float64    `json:"opening_stock"`
	TotalPurchaseOrderQuantity float64    `json:"total_po_quantity"`
	TotalPurchaseOrderAmount   float64    `json:"total_po_amount"`
	AveragePurchaseRate        float64    `json:"avg_purchase_rate"`
	TotalProductionOrderQuantity float64  `json:"total_prod_qty"`
	TotalManufacturedQuantity  float64    `json:"total_mfg_qty"`
	TotalConsumedInProduction  float64    `json:"total_consumed_in_mfg"`
	TotalSalesOrderQuantity    float64    `json:"total_so_quantity"`
	TotalSalesOrderAmount      float64    `json:"total_so_amount"`
	AverageSalesRate           float64    `json:"avg_sales_rate"`
	TotalInvoicedQuantity      float64    `json:"total_invoiced_qty"`
	CurrentQuantity            float64    `json:"current_qty"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}

type ListInventoryJournalOutput struct {
	Data      []InventoryJournalOutput `json:"data"`
	Total     int64                    `json:"total"`
	StartDate time.Time                `json:"start_date"`
	EndDate   time.Time                `json:"end_date"`
}
```

---

## Service Implementation Templates

### ItemGroupService Skeleton

```go
package services

import (
	"context"
	"errors"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
)

type ItemGroupService interface {
	Create(ctx context.Context, req input.CreateItemGroupInput) (*output.ItemGroupOutput, error)
	GetByID(ctx context.Context, id string) (*output.ItemGroupOutput, error)
	List(ctx context.Context, req input.ListItemGroupInput) (*output.ListItemGroupOutput, error)
	Update(ctx context.Context, id string, req input.UpdateItemGroupInput) (*output.ItemGroupOutput, error)
	Delete(ctx context.Context, id string) error
	ValidateComponents(ctx context.Context, components []input.ItemGroupComponentInput) error
	GetComponentRequirements(ctx context.Context, itemGroupID string, quantity float64) map[string]float64
}

type itemGroupService struct {
	repo repo.ItemGroupRepository
}

func NewItemGroupService(repo repo.ItemGroupRepository) ItemGroupService {
	return &itemGroupService{repo: repo}
}

func (s *itemGroupService) Create(ctx context.Context, req input.CreateItemGroupInput) (*output.ItemGroupOutput, error) {
	// Validate components exist
	// TODO: Call ValidateComponents
	
	// Create ItemGroup with components
	// TODO: Implementation
	
	return nil, nil
}

func (s *itemGroupService) ValidateComponents(ctx context.Context, components []input.ItemGroupComponentInput) error {
	// For each component:
	// 1. Verify Item exists
	// 2. If VariantID provided, verify Variant exists
	// 3. Verify Variant belongs to Item
	// TODO: Implementation
	
	return nil
}

func (s *itemGroupService) GetComponentRequirements(ctx context.Context, itemGroupID string, quantity float64) map[string]float64 {
	// itemGroup := get itemGroup by ID
	// requirements := make map
	// for each component:
	//     requirements[itemID|variantID] = component.quantity * quantity
	// return requirements
	// TODO: Implementation
	
	return nil
}
```

### ProductionOrderService Skeleton

```go
package services

type ProductionOrderService interface {
	Create(ctx context.Context, req input.CreateProductionOrderInput) (*output.ProductionOrderOutput, error)
	GetByID(ctx context.Context, id string) (*output.ProductionOrderOutput, error)
	List(ctx context.Context, req input.ListProductionOrderInput) (*output.ListProductionOrderOutput, error)
	UpdateStatus(ctx context.Context, id string, req input.UpdateProductionStatusInput) (*output.ProductionOrderOutput, error)
	MarkInProgress(ctx context.Context, id string) (*output.ProductionOrderOutput, error)
	MarkCompleted(ctx context.Context, id string, actualQuantity float64) (*output.ProductionOrderOutput, error)
	GetRequiredComponents(ctx context.Context, id string) ([]ComponentRequirement, error)
	CheckInventoryAvailability(ctx context.Context, id string) (bool, error)
	ConsumeComponentInventory(ctx context.Context, id string) error
	CreateFinishedInventory(ctx context.Context, id string) error
}

// TODO: Implement interface with business logic
```

### InventoryService Skeleton

```go
package services

type InventoryService interface {
	// Balance operations
	GetBalance(ctx context.Context, itemID string, variantID *uint) (*output.InventoryBalanceOutput, error)
	UpdateBalance(ctx context.Context, itemID string, variantID *uint, quantity float64) error
	ReserveInventory(ctx context.Context, itemID string, variantID *uint, quantity float64) error
	ReleaseReservation(ctx context.Context, itemID string, variantID *uint, quantity float64) error
	
	// Transaction logging
	LogPurchase(ctx context.Context, itemID string, variantID *uint, quantity float64, refID, refNo string) error
	LogManufacture(ctx context.Context, itemID string, variantID *uint, quantity float64, refID, refNo string) error
	LogConsumption(ctx context.Context, itemID string, variantID *uint, quantity float64, refID, refNo string) error
	LogSale(ctx context.Context, itemID string, variantID *uint, quantity float64, refID, refNo string) error
	
	// Aggregation
	GetAggregation(ctx context.Context, itemID string, variantID *uint) (*output.InventoryAggregationOutput, error)
	RecalculateAggregation(ctx context.Context, itemID string, variantID *uint) error
	
	// Supply chain summary
	GetSupplyChainSummary(ctx context.Context, itemID string) (*output.SupplyChainSummaryOutput, error)
	RecalculateSupplyChainSummary(ctx context.Context, itemID string) error
	
	// Journal
	GetJournal(ctx context.Context, itemID string, startDate, endDate time.Time, page, pageSize int) (*output.ListInventoryJournalOutput, error)
}

// TODO: Implement interface with business logic
```

---

## Handler Route Templates

### ItemGroup Routes

```go
// In routes/routes.go

func SetupItemGroupRoutes(router *gin.Engine, handler *handlers.ItemGroupHandler) {
	group := router.Group("/api/item-groups")
	
	group.POST("", handler.Create)              // Create ItemGroup
	group.GET("", handler.List)                 // List ItemGroups
	group.GET("/:id", handler.GetByID)          // Get ItemGroup
	group.PUT("/:id", handler.Update)           // Update ItemGroup
	group.DELETE("/:id", handler.Delete)        // Delete ItemGroup
	group.GET("/:id/components", handler.GetComponents)  // Get components
}
```

### ProductionOrder Routes

```go
func SetupProductionOrderRoutes(router *gin.Engine, handler *handlers.ProductionOrderHandler) {
	group := router.Group("/api/production-orders")
	
	group.POST("", handler.Create)              // Create ProductionOrder
	group.GET("", handler.List)                 // List ProductionOrders
	group.GET("/:id", handler.GetByID)          // Get ProductionOrder
	group.PUT("/:id", handler.Update)           // Update ProductionOrder
	group.PUT("/:id/status", handler.UpdateStatus) // Update status
	group.GET("/:id/components", handler.GetComponents) // Get required components
	group.POST("/:id/start", handler.MarkInProgress)   // Start production
	group.POST("/:id/complete", handler.Complete)      // Complete production
}
```

### Inventory Routes

```go
func SetupInventoryRoutes(router *gin.Engine, handler *handlers.InventoryHandler) {
	group := router.Group("/api/inventory")
	
	// Balance
	group.GET("/balance/:item_id", handler.GetBalance)
	group.PUT("/balance/:item_id/reserve", handler.ReserveInventory)
	group.PUT("/balance/:item_id/release", handler.ReleaseReservation)
	
	// Aggregation
	group.GET("/aggregation/:item_id", handler.GetAggregation)
	
	// Journal
	group.GET("/journal/:item_id", handler.GetJournal)
	
	// Supply Chain
	group.GET("/supply-chain/:item_id", handler.GetSupplyChainSummary)
}
```

---

## Quick Implementation Checklist

### Phase 1: Repository Layer
- [ ] Create `app/repo/item_group.repository.go`
- [ ] Create `app/repo/production_order.repository.go`
- [ ] Create `app/repo/inventory.repository.go`
- [ ] Implement all methods for each repository

### Phase 2: DTO Layer
- [ ] Create all input DTOs in `app/dto/input/`
- [ ] Create all output DTOs in `app/dto/output/`
- [ ] Add validation tags to input DTOs

### Phase 3: Service Layer
- [ ] Create `app/services/item_group.service.go`
- [ ] Create `app/services/production_order.service.go`
- [ ] Create `app/services/inventory.service.go`
- [ ] Implement business logic for all services

### Phase 4: Handler Layer
- [ ] Create `app/handlers/item_group.handler.go`
- [ ] Create `app/handlers/production_order.handler.go`
- [ ] Create `app/handlers/inventory.handler.go`
- [ ] Implement REST endpoints for all handlers

### Phase 5: Routes
- [ ] Update `app/routes/routes.go`
- [ ] Register all new route groups
- [ ] Add middleware for authentication/authorization

### Phase 6: Testing
- [ ] Unit tests for repositories
- [ ] Unit tests for services
- [ ] Integration tests for handlers
- [ ] API tests with Postman/curl

---

## Notes

- All services should accept `context.Context` for cancellation and timeout support
- All methods should include proper error handling
- Use dependency injection for repositories in services
- Use dependency injection for services in handlers
- All responses should follow the standard response wrapper format
- All errors should be properly logged
- Add proper validation in input DTOs
- Add pagination for list endpoints
- Consider caching for frequently accessed data (ItemGroups, aggregations)
