package services

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
)

type ProductGroupInventoryService interface {
	// Inventory initialization and management
	InitializeProductGroupInventory(productGroupID string) error
	GetInventoryStatus(productGroupID string) (*ProductGroupInventoryStatus, error)

	// Stock operations
	AddStock(productGroupID string, quantity float64, reason string, referenceID *string) error
	AllocateStock(productGroupID string, quantity float64, salesOrderID string) error
	DeductStock(productGroupID string, quantity float64, reason string, referenceID *string) error
	ReleaseAllocatedStock(productGroupID string, quantity float64, salesOrderID string) error
	ReverseTransaction(transactionID uint) error

	// Validation
	HasSufficientStock(productGroupID string, requiredQuantity float64) (bool, error)
	CanFulfillOrder(productGroupID string, quantity float64) (bool, string, error)
	GetMinimumStockLevel(productGroupID string) (float64, error)

	// Reporting
	GetInventoryHistory(productGroupID string, limit, offset int) ([]models.ProductGroupTransaction, int64, error)
	GetLowStockProducts() ([]models.ProductGroup, error)
	GetInventoryReport(productGroupID string) (*InventoryReport, error)
}

type ProductGroupInventoryStatus struct {
	ProductGroupID  string                 `json:"product_group_id"`
	CurrentStock    float64                `json:"current_stock"`
	AllocatedStock  float64                `json:"allocated_stock"`
	AvailableStock  float64                `json:"available_stock"`
	TotalReceived   float64                `json:"total_received"`
	TotalSold       float64                `json:"total_sold"`
	DamagedStock    float64                `json:"damaged_stock"`
	ComponentStatus []ComponentStockStatus `json:"component_status"`
}

type ComponentStockStatus struct {
	ComponentProductID   string  `json:"component_product_id"`
	ComponentProductName string  `json:"component_product_name"`
	VariantSku           *string `json:"variant_sku,omitempty"`
	QuantityPerGroup     float64 `json:"quantity_per_group"`
	CurrentStock         float64 `json:"current_stock"`
	AllocatedStock       float64 `json:"allocated_stock"`
	AvailableStock       float64 `json:"available_stock"`
	CanFulfill           bool    `json:"can_fulfill"` // Can this component fulfill one product group
}

type InventoryReport struct {
	ProductGroupID     string                           `json:"product_group_id"`
	ProductGroupName   string                           `json:"product_group_name"`
	CurrentStock       float64                          `json:"current_stock"`
	AllocatedStock     float64                          `json:"allocated_stock"`
	AvailableStock     float64                          `json:"available_stock"`
	TotalReceived      float64                          `json:"total_received"`
	TotalSold          float64                          `json:"total_sold"`
	DamagedStock       float64                          `json:"damaged_stock"`
	LastPurchaseDate   *time.Time                       `json:"last_purchase_date"`
	LastSaleDate       *time.Time                       `json:"last_sale_date"`
	Components         []ComponentInventoryReport       `json:"components"`
	RecentTransactions []models.ProductGroupTransaction `json:"recent_transactions"`
}

type ComponentInventoryReport struct {
	ComponentProductID   string  `json:"component_product_id"`
	ComponentProductName string  `json:"component_product_name"`
	VariantSku           *string `json:"variant_sku,omitempty"`
	QuantityPerGroup     float64 `json:"quantity_per_group"`
	CurrentStock         float64 `json:"current_stock"`
	RequiredForGroup     float64 `json:"required_for_current_stock"` // current_stock * quantity_per_group
}

type productGroupInventoryService struct {
	pgInventoryRepo   repo.ProductGroupInventoryRepository
	compInventoryRepo repo.ComponentInventoryRepository
	pgTransactionRepo repo.ProductGroupTransactionRepository
	productGroupRepo  repo.ProductGroupRepository
	productRepo       repo.ProductRepository
}

func NewProductGroupInventoryService(
	pgInventoryRepo repo.ProductGroupInventoryRepository,
	compInventoryRepo repo.ComponentInventoryRepository,
	pgTransactionRepo repo.ProductGroupTransactionRepository,
	productGroupRepo repo.ProductGroupRepository,
	productRepo repo.ProductRepository,
) ProductGroupInventoryService {
	return &productGroupInventoryService{
		pgInventoryRepo:   pgInventoryRepo,
		compInventoryRepo: compInventoryRepo,
		pgTransactionRepo: pgTransactionRepo,
		productGroupRepo:  productGroupRepo,
		productRepo:       productRepo,
	}
}

// InitializeProductGroupInventory creates initial inventory records for a product group and its components
func (s *productGroupInventoryService) InitializeProductGroupInventory(productGroupID string) error {
	// Check if product group exists
	productGroup, err := s.productGroupRepo.FindByID(productGroupID)
	if err != nil {
		return fmt.Errorf("product group not found: %w", err)
	}

	// Create product group inventory
	pgInventory := &models.ProductGroupInventory{
		ProductGroupID: productGroupID,
		CurrentStock:   0,
		AllocatedStock: 0,
		AvailableStock: 0,
		TotalReceived:  0,
		TotalSold:      0,
		DamagedStock:   0,
	}

	err = s.pgInventoryRepo.Create(pgInventory)
	if err != nil {
		return fmt.Errorf("failed to create product group inventory: %w", err)
	}

	// Create component inventories
	for _, component := range productGroup.Components {
		compInventory := &models.ComponentInventory{
			ProductGroupID:      productGroupID,
			ComponentProductID:  component.ProductID,
			ComponentVariantSku: component.VariantSku,
			QuantityPerGroup:    component.Quantity,
			CurrentStock:        0,
			AllocatedStock:      0,
			AvailableStock:      0,
		}

		err = s.compInventoryRepo.Create(compInventory)
		if err != nil {
			return fmt.Errorf("failed to create component inventory: %w", err)
		}
	}

	return nil
}

// GetInventoryStatus retrieves current inventory status with component details
func (s *productGroupInventoryService) GetInventoryStatus(productGroupID string) (*ProductGroupInventoryStatus, error) {
	pgInventory, err := s.pgInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return nil, fmt.Errorf("product group inventory not found: %w", err)
	}

	compInventories, err := s.compInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch component inventories: %w", err)
	}

	status := &ProductGroupInventoryStatus{
		ProductGroupID:  productGroupID,
		CurrentStock:    pgInventory.CurrentStock,
		AllocatedStock:  pgInventory.AllocatedStock,
		AvailableStock:  pgInventory.AvailableStock,
		TotalReceived:   pgInventory.TotalReceived,
		TotalSold:       pgInventory.TotalSold,
		DamagedStock:    pgInventory.DamagedStock,
		ComponentStatus: make([]ComponentStockStatus, 0, len(compInventories)),
	}

	for _, ci := range compInventories {
		var componentName string
		if ci.ComponentProduct != nil {
			componentName = ci.ComponentProduct.Name
		} else if product, err := s.productRepo.FindByID(ci.ComponentProductID); err == nil {
			componentName = product.Name
		}

		// A component can fulfil one product group if its available stock >= quantity_per_group
		canFulfill := ci.AvailableStock >= ci.QuantityPerGroup

		compStatus := ComponentStockStatus{
			ComponentProductID:   ci.ComponentProductID,
			ComponentProductName: componentName,
			VariantSku:           ci.ComponentVariantSku,
			QuantityPerGroup:     ci.QuantityPerGroup,
			CurrentStock:         ci.CurrentStock,
			AllocatedStock:       ci.AllocatedStock,
			AvailableStock:       ci.AvailableStock,
			CanFulfill:           canFulfill,
		}
		status.ComponentStatus = append(status.ComponentStatus, compStatus)
	}

	return status, nil
}

// AddStock adds stock to the product group (received from purchase)
func (s *productGroupInventoryService) AddStock(productGroupID string, quantity float64, reason string, referenceID *string) error {
	pgInventory, err := s.pgInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return fmt.Errorf("product group inventory not found: %w", err)
	}

	// Update stock
	pgInventory.CurrentStock += quantity
	pgInventory.AvailableStock = pgInventory.CurrentStock - pgInventory.AllocatedStock
	pgInventory.TotalReceived += quantity
	pgInventory.UpdatedAt = time.Now()

	err = s.pgInventoryRepo.Update(pgInventory)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Log transaction
	transaction := &models.ProductGroupTransaction{
		ProductGroupID:  productGroupID,
		TransactionType: "purchase",
		Quantity:        quantity,
		Notes:           reason,
		CreatedAt:       time.Now(),
	}

	if referenceID != nil {
		transaction.PurchaseOrderID = referenceID
	}

	err = s.pgTransactionRepo.Create(transaction)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %w", err)
	}

	return nil
}

// AllocateStock reserves stock for a sales order (inventory reserved but not deducted)
func (s *productGroupInventoryService) AllocateStock(productGroupID string, quantity float64, salesOrderID string) error {
	// Check if sufficient stock is available
	hasStock, err := s.HasSufficientStock(productGroupID, quantity)
	if err != nil {
		return err
	}
	if !hasStock {
		return fmt.Errorf("insufficient stock available for allocation")
	}

	pgInventory, err := s.pgInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return fmt.Errorf("product group inventory not found: %w", err)
	}

	pgInventory.AllocatedStock += quantity
	pgInventory.AvailableStock = pgInventory.CurrentStock - pgInventory.AllocatedStock
	pgInventory.UpdatedAt = time.Now()

	err = s.pgInventoryRepo.Update(pgInventory)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}

// DeductStock removes stock from inventory (after shipment)
func (s *productGroupInventoryService) DeductStock(productGroupID string, quantity float64, reason string, referenceID *string) error {
	pgInventory, err := s.pgInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return fmt.Errorf("product group inventory not found: %w", err)
	}

	if pgInventory.CurrentStock < quantity {
		return fmt.Errorf("insufficient current stock to deduct")
	}

	pgInventory.CurrentStock -= quantity
	pgInventory.AvailableStock = pgInventory.CurrentStock - pgInventory.AllocatedStock
	pgInventory.TotalSold += quantity
	pgInventory.UpdatedAt = time.Now()

	err = s.pgInventoryRepo.Update(pgInventory)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Log transaction
	transaction := &models.ProductGroupTransaction{
		ProductGroupID:  productGroupID,
		TransactionType: "sales",
		Quantity:        quantity,
		Notes:           reason,
		CreatedAt:       time.Now(),
	}

	if referenceID != nil {
		transaction.ShipmentID = referenceID
	}

	err = s.pgTransactionRepo.Create(transaction)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %w", err)
	}

	return nil
}

// ReleaseAllocatedStock releases previously allocated stock (order cancelled)
func (s *productGroupInventoryService) ReleaseAllocatedStock(productGroupID string, quantity float64, salesOrderID string) error {
	pgInventory, err := s.pgInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return fmt.Errorf("product group inventory not found: %w", err)
	}

	if pgInventory.AllocatedStock < quantity {
		return fmt.Errorf("cannot release more than allocated stock")
	}

	pgInventory.AllocatedStock -= quantity
	pgInventory.AvailableStock = pgInventory.CurrentStock - pgInventory.AllocatedStock
	pgInventory.UpdatedAt = time.Now()

	err = s.pgInventoryRepo.Update(pgInventory)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}

// HasSufficientStock checks if enough stock is available
func (s *productGroupInventoryService) HasSufficientStock(productGroupID string, requiredQuantity float64) (bool, error) {
	pgInventory, err := s.pgInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return false, fmt.Errorf("product group inventory not found: %w", err)
	}

	return pgInventory.AvailableStock >= requiredQuantity, nil
}

// CanFulfillOrder checks if all components have sufficient stock
func (s *productGroupInventoryService) CanFulfillOrder(productGroupID string, quantity float64) (bool, string, error) {
	compInventories, err := s.compInventoryRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch component inventories: %w", err)
	}

	for _, ci := range compInventories {
		required := quantity * ci.QuantityPerGroup
		if ci.AvailableStock < required {
			return false, fmt.Sprintf("Insufficient stock for component %s. Required: %f, Available: %f",
				ci.ComponentProductID, required, ci.AvailableStock), nil
		}
	}

	return true, "", nil
}

// ReverseTransaction rolls back a transaction
func (s *productGroupInventoryService) ReverseTransaction(transactionID uint) error {
	transaction, err := s.pgTransactionRepo.FindByID(transactionID)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Reverse based on transaction type
	switch transaction.TransactionType {
	case "purchase":
		return s.DeductStock(transaction.ProductGroupID, transaction.Quantity, "Reverse purchase", transaction.PurchaseOrderID)
	case "sales":
		return s.AddStock(transaction.ProductGroupID, transaction.Quantity, "Reverse sales", transaction.ShipmentID)
	default:
		return fmt.Errorf("unknown transaction type: %s", transaction.TransactionType)
	}
}

// GetInventoryHistory retrieves transaction history
func (s *productGroupInventoryService) GetInventoryHistory(productGroupID string, limit, offset int) ([]models.ProductGroupTransaction, int64, error) {
	return s.pgTransactionRepo.FindByProductGroupID(productGroupID, limit, offset)
}

// GetLowStockProducts returns product groups with low stock
func (s *productGroupInventoryService) GetLowStockProducts() ([]models.ProductGroup, error) {
	// This would require a threshold setting
	// For now, returning empty - implement based on your business rule
	return []models.ProductGroup{}, nil
}

// GetInventoryReport generates comprehensive inventory report
func (s *productGroupInventoryService) GetInventoryReport(productGroupID string) (*InventoryReport, error) {
	productGroup, err := s.productGroupRepo.FindByID(productGroupID)
	if err != nil {
		return nil, fmt.Errorf("product group not found: %w", err)
	}

	status, err := s.GetInventoryStatus(productGroupID)
	if err != nil {
		return nil, err
	}

	compInventories, _ := s.compInventoryRepo.FindByProductGroupID(productGroupID)
	transactions, _, _ := s.GetInventoryHistory(productGroupID, 20, 0)

	report := &InventoryReport{
		ProductGroupID:     productGroupID,
		ProductGroupName:   productGroup.Name,
		CurrentStock:       status.CurrentStock,
		AllocatedStock:     status.AllocatedStock,
		AvailableStock:     status.AvailableStock,
		TotalReceived:      status.TotalReceived,
		TotalSold:          status.TotalSold,
		DamagedStock:       status.DamagedStock,
		Components:         make([]ComponentInventoryReport, len(compInventories)),
		RecentTransactions: transactions,
	}

	for i, ci := range compInventories {
		report.Components[i] = ComponentInventoryReport{
			ComponentProductID: ci.ComponentProductID,
			VariantSku:         ci.ComponentVariantSku,
			QuantityPerGroup:   ci.QuantityPerGroup,
			CurrentStock:       ci.CurrentStock,
			RequiredForGroup:   status.CurrentStock * ci.QuantityPerGroup,
		}
	}

	return report, nil
}

func (s *productGroupInventoryService) GetMinimumStockLevel(productGroupID string) (float64, error) {
	// Implement based on your business rules
	// For now returning 10 as default minimum
	return 10, nil
}
