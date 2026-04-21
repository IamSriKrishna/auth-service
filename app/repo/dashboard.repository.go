package repo

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	// Metrics
	GetDashboardMetrics() (*models.DashboardMetrics, error)
	SaveDashboardMetrics(metrics *models.DashboardMetrics) error

	// Customer metrics
	GetTotalCustomers() (int64, error)
	GetActiveCustomers() (int64, error)
	GetCustomersCreatedToday() (int64, error)
	GetTotalCustomersWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetActiveCustomersWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetCustomersCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Vendor metrics
	GetTotalVendors() (int64, error)
	GetActiveVendors() (int64, error)
	GetVendorsCreatedToday() (int64, error)
	GetTotalVendorsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetActiveVendorsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetVendorsCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Item metrics
	GetTotalItems() (int64, error)
	GetTotalItemGroups() (int64, error)
	GetTotalStock() (int64, error)
	GetLowStockItems(threshold int64) (int64, error)
	GetOutOfStockItems() (int64, error)
	GetItemsCreatedToday() (int64, error)
	GetItemStockDetails() ([]map[string]interface{}, error)
	GetItemStockDetailsWithFilter(shouldFilter bool, userID uint) ([]map[string]interface{}, error)
	GetTotalItemsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetTotalItemGroupsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetTotalStockWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetLowStockItemsWithFilter(threshold int64, shouldFilter bool, userID uint) (int64, error)
	GetOutOfStockItemsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetItemsCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Shipment metrics
	GetTotalShipments() (int64, error)
	GetShippedCount() (int64, error)
	GetShipmentsByStatus(status string) (int64, error)
	GetShippedToday() (int64, error)
	GetTotalShipmentsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetShippedCountWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetShippedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetShipmentsByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error)

	// Invoice metrics
	GetTotalInvoices() (int64, error)
	GetTotalInvoiceAmount() (float64, error)
	GetOutstandingInvoices() (float64, error)
	GetInvoicesByStatus(status string) (int64, error)
	GetOverdueInvoices() (int64, error)
	GetTotalInvoicesWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetTotalInvoiceAmountWithFilter(shouldFilter bool, userID uint) (float64, error)
	GetOutstandingInvoicesWithFilter(shouldFilter bool, userID uint) (float64, error)
	GetInvoicesByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error)
	GetOverdueInvoicesWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Sales order metrics
	GetTotalSalesOrders() (int64, error)
	GetTotalSalesOrderAmount() (float64, error)
	GetSalesOrdersByStatus(status string) (int64, error)
	GetSalesOrdersCreatedToday() (int64, error)
	GetTotalSalesOrdersWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetTotalSalesOrderAmountWithFilter(shouldFilter bool, userID uint) (float64, error)
	GetSalesOrdersByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error)
	GetSalesOrdersCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Purchase order metrics
	GetTotalPurchaseOrders() (int64, error)
	GetTotalPurchaseOrderAmount() (float64, error)
	GetPurchaseOrdersByStatus(status string) (int64, error)
	GetPurchaseOrdersCreatedToday() (int64, error)
	GetTotalPurchaseOrdersWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetTotalPurchaseOrderAmountWithFilter(shouldFilter bool, userID uint) (float64, error)
	GetPurchaseOrdersByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error)
	GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Package metrics
	GetTotalPackages() (int64, error)
	GetPackagesByStatus(status string) (int64, error)
	GetPackagesCreatedToday() (int64, error)
	GetTotalPackagesWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetPackagesByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error)
	GetPackagesCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error)

	// Product stock metrics
	GetTotalProducts() (int64, error)
	GetInStockProducts() (int64, error)
	GetLowStockProducts(threshold float64) (int64, error)
	GetOutOfStockProducts() (int64, error)
	GetTotalProductStock() (float64, error)
	GetProductStockDetails() ([]map[string]interface{}, error)
	GetProductStockDetailsWithFilter(shouldFilter bool, userID uint) ([]map[string]interface{}, error)
	GetTotalProductsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetInStockProductsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetLowStockProductsWithFilter(threshold float64, shouldFilter bool, userID uint) (int64, error)
	GetOutOfStockProductsWithFilter(shouldFilter bool, userID uint) (int64, error)
	GetTotalProductStockWithFilter(shouldFilter bool, userID uint) (float64, error)

	// Shipment tracking
	AddShipmentTracking(tracking *models.ShipmentTracking) error
	GetShipmentTracking(shipmentID string, limit int) ([]models.ShipmentTracking, error)
	GetLatestShipmentTracking(shipmentID string) (*models.ShipmentTracking, error)

	// History
	GetEntityCountHistory(entityType string, days int) ([]models.EntityCountHistory, error)
	SaveEntityCountHistory(history *models.EntityCountHistory) error
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetDashboardMetrics retrieves the latest dashboard metrics
func (r *dashboardRepository) GetDashboardMetrics() (*models.DashboardMetrics, error) {
	var metrics models.DashboardMetrics
	err := r.db.Order("created_at DESC").First(&metrics).Error
	return &metrics, err
}

// SaveDashboardMetrics saves dashboard metrics
func (r *dashboardRepository) SaveDashboardMetrics(metrics *models.DashboardMetrics) error {
	return r.db.Save(metrics).Error
}

// Customer metrics
func (r *dashboardRepository) GetTotalCustomers() (int64, error) {
	var count int64
	err := r.db.Model(&models.Customer{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetActiveCustomers() (int64, error) {
	var count int64
	err := r.db.Model(&models.Customer{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetCustomersCreatedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.Customer{}).
		Where("created_at >= ?", today).
		Count(&count).Error
	return count, err
}

// Vendor metrics
func (r *dashboardRepository) GetTotalVendors() (int64, error) {
	var count int64
	err := r.db.Model(&models.Vendor{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetActiveVendors() (int64, error) {
	var count int64
	err := r.db.Model(&models.Vendor{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetVendorsCreatedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.Vendor{}).
		Where("created_at >= ?", today).
		Count(&count).Error
	return count, err
}

// Item metrics
func (r *dashboardRepository) GetTotalItems() (int64, error) {
	var count int64
	err := r.db.Model(&models.Item{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalItemGroups() (int64, error) {
	var count int64
	err := r.db.Model(&models.ItemGroup{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalStock() (int64, error) {
	var total int64
	err := r.db.Model(&models.InventoryBalance{}).
		Select("COALESCE(SUM(quantity), 0)").
		Row().
		Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetLowStockItems(threshold int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.InventoryBalance{}).
		Where("quantity > 0 AND quantity <= ?", threshold).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOutOfStockItems() (int64, error) {
	var count int64
	err := r.db.Model(&models.InventoryBalance{}).
		Where("quantity <= 0").
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetItemsCreatedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.Item{}).
		Where("created_at >= ?", today).
		Count(&count).Error
	return count, err
}

// GetItemStockDetails retrieves stock details for all items
func (r *dashboardRepository) GetItemStockDetails() ([]map[string]interface{}, error) {
	// First, let's get all items with their basic info
	var items []models.Item
	err := r.db.Find(&items).Error

	if err != nil {
		fmt.Printf("GetItemStockDetails Error fetching items: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetItemStockDetails: Found %d items from items table\n", len(items))

	result := make([]map[string]interface{}, 0)

	// For each item, try to get its inventory balance
	for _, item := range items {
		fmt.Printf("Processing item: ID=%s, Name=%s\n", item.ID, item.Name)

		var inventory models.InventoryBalance
		invErr := r.db.Where("item_id = ?", item.ID).First(&inventory).Error

		currentQty := 0.0
		availableQty := 0.0
		reservedQty := 0.0
		inTransitQty := 0.0

		if invErr == nil {
			currentQty = inventory.CurrentQuantity
			availableQty = inventory.AvailableQuantity
			reservedQty = inventory.ReservedQuantity
			inTransitQty = inventory.InTransitQuantity
			fmt.Printf("  Found inventory: CurrentQty=%f\n", currentQty)
		} else if invErr.Error() != "record not found" {
			fmt.Printf("  Inventory query error: %v\n", invErr)
		} else {
			fmt.Printf("  No inventory record found\n")
		}

		stockStatus := "in_stock"
		if currentQty == 0 {
			stockStatus = "out_of_stock"
		} else if currentQty <= 100 { // default threshold
			stockStatus = "low_stock"
		}

		result = append(result, map[string]interface{}{
			"item_id":             item.ID,
			"item_name":           item.Name,
			"current_quantity":    currentQty,
			"available_quantity":  availableQty,
			"reserved_quantity":   reservedQty,
			"in_transit_quantity": inTransitQty,
			"status":              stockStatus,
		})
	}

	fmt.Printf("GetItemStockDetails: Returning %d results\n", len(result))
	return result, nil
}

// GetItemStockDetailsWithFilter retrieves stock details for items filtered by user
func (r *dashboardRepository) GetItemStockDetailsWithFilter(shouldFilter bool, userID uint) ([]map[string]interface{}, error) {
	// First, let's get items with optional filtering
	var items []models.Item
	query := r.db

	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}

	err := query.Find(&items).Error

	if err != nil {
		fmt.Printf("GetItemStockDetailsWithFilter Error fetching items: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetItemStockDetailsWithFilter: Found %d items (shouldFilter=%v, userID=%d)\n", len(items), shouldFilter, userID)

	result := make([]map[string]interface{}, 0)

	// For each item, try to get its inventory balance
	for _, item := range items {
		fmt.Printf("Processing item: ID=%s, Name=%s, CreatedBy=%s\n", item.ID, item.Name, item.CreatedBy)

		var inventory models.InventoryBalance
		invErr := r.db.Where("item_id = ?", item.ID).First(&inventory).Error

		currentQty := 0.0
		availableQty := 0.0
		reservedQty := 0.0
		inTransitQty := 0.0

		if invErr == nil {
			currentQty = inventory.CurrentQuantity
			availableQty = inventory.AvailableQuantity
			reservedQty = inventory.ReservedQuantity
			inTransitQty = inventory.InTransitQuantity
			fmt.Printf("  Found inventory: CurrentQty=%f\n", currentQty)
		} else if invErr.Error() != "record not found" {
			fmt.Printf("  Inventory query error: %v\n", invErr)
		} else {
			fmt.Printf("  No inventory record found\n")
		}

		stockStatus := "in_stock"
		if currentQty == 0 {
			stockStatus = "out_of_stock"
		} else if currentQty <= 100 { // default threshold
			stockStatus = "low_stock"
		}

		result = append(result, map[string]interface{}{
			"item_id":             item.ID,
			"item_name":           item.Name,
			"current_quantity":    currentQty,
			"available_quantity":  availableQty,
			"reserved_quantity":   reservedQty,
			"in_transit_quantity": inTransitQty,
			"status":              stockStatus,
		})
	}

	fmt.Printf("GetItemStockDetailsWithFilter: Returning %d results\n", len(result))
	return result, nil
}

// Shipment metrics
func (r *dashboardRepository) GetTotalShipments() (int64, error) {
	var count int64
	err := r.db.Model(&models.Shipment{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShippedCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Shipment{}).
		Where("status IN ?", []string{"shipped", "delivered"}).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShipmentsByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Shipment{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShippedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.Shipment{}).
		Where("status = ? AND updated_at >= ?", "shipped", today).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShippedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Shipment{}).
		Where("status = ? AND updated_at >= ?", "shipped", today)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Invoice metrics
func (r *dashboardRepository) GetTotalInvoices() (int64, error) {
	var count int64
	err := r.db.Model(&models.Invoice{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalInvoiceAmount() (float64, error) {
	var amount float64
	err := r.db.Model(&models.Invoice{}).
		Select("COALESCE(SUM(total), 0)").
		Row().
		Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetOutstandingInvoices() (float64, error) {
	var amount float64
	err := r.db.Model(&models.Invoice{}).
		Where("status IN ?", []string{"pending", "partially_paid"}).
		Select("COALESCE(SUM(total), 0)").
		Row().
		Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetInvoicesByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Invoice{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOverdueInvoices() (int64, error) {
	var count int64
	err := r.db.Model(&models.Invoice{}).
		Where("status IN ? AND due_date < NOW()", []string{"pending", "partially_paid"}).
		Count(&count).Error
	return count, err
}

// Sales order metrics
func (r *dashboardRepository) GetTotalSalesOrders() (int64, error) {
	var count int64
	err := r.db.Model(&models.SalesOrder{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalSalesOrderAmount() (float64, error) {
	var amount float64
	err := r.db.Model(&models.SalesOrder{}).
		Select("COALESCE(SUM(total), 0)").
		Row().
		Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetSalesOrdersByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.SalesOrder{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetSalesOrdersCreatedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.SalesOrder{}).
		Where("created_at >= ?", today).
		Count(&count).Error
	return count, err
}

// Purchase order metrics
func (r *dashboardRepository) GetTotalPurchaseOrders() (int64, error) {
	var count int64
	err := r.db.Model(&models.PurchaseOrder{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalPurchaseOrderAmount() (float64, error) {
	var amount float64
	err := r.db.Model(&models.PurchaseOrder{}).
		Select("COALESCE(SUM(total), 0)").
		Row().
		Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetPurchaseOrdersByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.PurchaseOrder{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPurchaseOrdersCreatedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.PurchaseOrder{}).
		Where("created_at >= ?", today).
		Count(&count).Error
	return count, err
}

// Package metrics
func (r *dashboardRepository) GetTotalPackages() (int64, error) {
	var count int64
	err := r.db.Model(&models.Package{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPackagesByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Package{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPackagesCreatedToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&models.Package{}).
		Where("created_at >= ?", today).
		Count(&count).Error
	return count, err
}

// Shipment tracking
func (r *dashboardRepository) AddShipmentTracking(tracking *models.ShipmentTracking) error {
	return r.db.Create(tracking).Error
}

func (r *dashboardRepository) GetShipmentTracking(shipmentID string, limit int) ([]models.ShipmentTracking, error) {
	var tracking []models.ShipmentTracking
	err := r.db.Where("shipment_id = ?", shipmentID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&tracking).Error
	return tracking, err
}

func (r *dashboardRepository) GetLatestShipmentTracking(shipmentID string) (*models.ShipmentTracking, error) {
	var tracking models.ShipmentTracking
	err := r.db.Where("shipment_id = ?", shipmentID).
		Order("timestamp DESC").
		First(&tracking).Error
	return &tracking, err
}

// History
func (r *dashboardRepository) GetEntityCountHistory(entityType string, days int) ([]models.EntityCountHistory, error) {
	var history []models.EntityCountHistory
	startDate := time.Now().AddDate(0, 0, -days)
	err := r.db.Where("entity_type = ? AND date >= ?", entityType, startDate).
		Order("date ASC").
		Find(&history).Error
	return history, err
}

func (r *dashboardRepository) SaveEntityCountHistory(history *models.EntityCountHistory) error {
	return r.db.Save(history).Error
}

// ==================== Filter Methods for User-based Dashboard ====================

// Customer metrics with filter
func (r *dashboardRepository) GetTotalCustomersWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Customer{})
	if shouldFilter && userID > 0 {
		fmt.Printf("Filtering customers by user_id=%d\n", userID)
		query = query.Where("user_id = ?", userID)
	}
	err := query.Count(&count).Error
	fmt.Printf("GetTotalCustomersWithFilter result: count=%d, filter=%v, userID=%d\n", count, shouldFilter, userID)
	return count, err
}

func (r *dashboardRepository) GetActiveCustomersWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Customer{}).Where("is_active = ?", true)
	if shouldFilter {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetCustomersCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Customer{}).Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Vendor metrics with filter
func (r *dashboardRepository) GetTotalVendorsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Vendor{})
	if shouldFilter && userID > 0 {
		fmt.Printf("Filtering vendors by user_id=%d\n", userID)
		query = query.Where("user_id = ?", userID)
	}
	err := query.Count(&count).Error
	fmt.Printf("GetTotalVendorsWithFilter result: count=%d, filter=%v, userID=%d\n", count, shouldFilter, userID)
	return count, err
}

func (r *dashboardRepository) GetActiveVendorsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Vendor{}).Where("is_active = ?", true)
	if shouldFilter {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetVendorsCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Vendor{}).Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Item metrics with filter
func (r *dashboardRepository) GetTotalItemsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Item{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalItemGroupsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.ItemGroup{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalStockWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var total int64
	query := r.db.Model(&models.InventoryBalance{})
	if shouldFilter {
		// For inventory, we may need to join with items to filter by user
		query = query.Joins("LEFT JOIN items ON items.id = inventory_balances.item_id").
			Where("items.created_by = ?", userID)
	}
	err := query.Select("COALESCE(SUM(quantity), 0)").Row().Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetLowStockItemsWithFilter(threshold int64, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.InventoryBalance{}).
		Where("quantity > 0 AND quantity <= ?", threshold)
	if shouldFilter {
		query = query.Joins("LEFT JOIN items ON items.id = inventory_balances.item_id").
			Where("items.created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOutOfStockItemsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.InventoryBalance{}).
		Where("quantity <= 0")
	if shouldFilter {
		query = query.Joins("LEFT JOIN items ON items.id = inventory_balances.item_id").
			Where("items.created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetItemsCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Item{}).Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Shipment metrics with filter
func (r *dashboardRepository) GetTotalShipmentsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Shipment{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShippedCountWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Shipment{}).
		Where("status IN ?", []string{"shipped", "delivered"})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShipmentsByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Shipment{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Invoice metrics with filter
func (r *dashboardRepository) GetTotalInvoicesWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Invoice{})
	if shouldFilter && userID > 0 {
		userIDStr := fmt.Sprintf("%d", userID)
		fmt.Printf("Filtering invoices by created_by=%s\n", userIDStr)
		query = query.Where("created_by = ?", userIDStr)
	}
	err := query.Count(&count).Error
	fmt.Printf("GetTotalInvoicesWithFilter result: count=%d, filter=%v, userID=%d\n", count, shouldFilter, userID)
	return count, err
}

func (r *dashboardRepository) GetTotalInvoiceAmountWithFilter(shouldFilter bool, userID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.Invoice{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Select("COALESCE(SUM(total), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetOutstandingInvoicesWithFilter(shouldFilter bool, userID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.Invoice{}).
		Where("status IN ?", []string{"pending", "partially_paid"})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Select("COALESCE(SUM(total), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetInvoicesByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Invoice{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOverdueInvoicesWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Invoice{}).
		Where("status IN ? AND due_date < ?", []string{"pending", "partially_paid"}, time.Now())
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Sales order metrics with filter
func (r *dashboardRepository) GetTotalSalesOrdersWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.SalesOrder{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalSalesOrderAmountWithFilter(shouldFilter bool, userID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.SalesOrder{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetSalesOrdersByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.SalesOrder{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetSalesOrdersCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.SalesOrder{}).
		Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Purchase order metrics with filter
func (r *dashboardRepository) GetTotalPurchaseOrdersWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.PurchaseOrder{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalPurchaseOrderAmountWithFilter(shouldFilter bool, userID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.PurchaseOrder{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetPurchaseOrdersByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.PurchaseOrder{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.PurchaseOrder{}).
		Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Package metrics with filter
func (r *dashboardRepository) GetTotalPackagesWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Package{})
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPackagesByStatusWithFilter(status string, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Package{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPackagesCreatedTodayWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Package{}).
		Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("created_by = ?", userID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Product stock metrics - Aggregated from variant stocks
func (r *dashboardRepository) GetTotalProducts() (int64, error) {
	var count int64
	err := r.db.Model(&models.VariantStock{}).
		Distinct("product_id").
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetInStockProducts() (int64, error) {
	var count int64
	err := r.db.Model(&models.VariantStock{}).
		Select("DISTINCT product_id").
		Where("available_stock > 0").
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetLowStockProducts(threshold float64) (int64, error) {
	var count int64
	err := r.db.Model(&models.VariantStock{}).
		Select("DISTINCT product_id").
		Where("available_stock > 0 AND available_stock <= ?", threshold).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOutOfStockProducts() (int64, error) {
	var count int64
	err := r.db.Model(&models.VariantStock{}).
		Select("DISTINCT product_id").
		Where("available_stock <= 0").
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalProductStock() (float64, error) {
	var total float64
	err := r.db.Model(&models.VariantStock{}).
		Select("COALESCE(SUM(current_stock), 0)").
		Row().
		Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetProductStockDetails() ([]map[string]interface{}, error) {
	// Aggregate variant stocks by product
	type ProductStockAggregate struct {
		ProductID         string `gorm:"column:product_id"`
		ProductName       string `gorm:"column:product_name"`
		CurrentStock      float64
		AvailableStock    float64
		ReservedStock     float64
		PurchasedStock    float64
		SoldStock         float64
		AverageCost       float64
		RevaluationAmount float64
		LastPurchasedDate *time.Time
		LastSoldDate      *time.Time
	}

	var products []ProductStockAggregate
	// Use LEFT JOIN to include product groups (pg_xxx) which may not have corresponding products
	err := r.db.Model(&models.VariantStock{}).
		Joins("LEFT JOIN products ON variant_stocks.product_id = products.id").
		Select(`
			product_id,
			product_name,
			COALESCE(SUM(current_stock), 0) AS current_stock,
			COALESCE(SUM(available_stock), 0) AS available_stock,
			COALESCE(SUM(reserved_stock), 0) AS reserved_stock,
			COALESCE(SUM(purchased_stock), 0) AS purchased_stock,
			COALESCE(SUM(sold_stock), 0) AS sold_stock,
			COALESCE(AVG(average_cost), 0) AS average_cost,
			COALESCE(SUM(revaluation_amount), 0) AS revaluation_amount,
			MAX(last_purchased_date) AS last_purchased_date,
			MAX(last_sold_date) AS last_sold_date
		`).
		Group("product_id, product_name").
		Scan(&products).Error

	if err != nil {
		fmt.Printf("GetProductStockDetails Error fetching products: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetProductStockDetails: Found %d products\n", len(products))

	result := make([]map[string]interface{}, 0)

	for _, product := range products {
		stockStatus := "in_stock"
		if product.AvailableStock <= 0 {
			stockStatus = "out_of_stock"
		} else if product.AvailableStock <= 100 { // default threshold
			stockStatus = "low_stock"
		}

		result = append(result, map[string]interface{}{
			"product_id":          product.ProductID,
			"product_name":        product.ProductName,
			"sku":                 "", // Not available at product level
			"current_stock":       product.CurrentStock,
			"available_stock":     product.AvailableStock,
			"reserved_stock":      product.ReservedStock,
			"purchased_stock":     product.PurchasedStock,
			"sold_stock":          product.SoldStock,
			"average_cost":        product.AverageCost,
			"revaluation_amount":  product.RevaluationAmount,
			"last_purchased_date": product.LastPurchasedDate,
			"last_sold_date":      product.LastSoldDate,
			"status":              stockStatus,
		})
	}

	fmt.Printf("GetProductStockDetails: Returning %d results\n", len(result))
	return result, nil
}

func (r *dashboardRepository) GetProductStockDetailsWithFilter(shouldFilter bool, userID uint) ([]map[string]interface{}, error) {
	// Aggregate variant stocks by product, filtered by user
	type ProductStockAggregate struct {
		ProductID         string `gorm:"column:product_id"`
		ProductName       string `gorm:"column:product_name"`
		CurrentStock      float64
		AvailableStock    float64
		ReservedStock     float64
		PurchasedStock    float64
		SoldStock         float64
		AverageCost       float64
		RevaluationAmount float64
		LastPurchasedDate *time.Time
		LastSoldDate      *time.Time
	}

	var products []ProductStockAggregate
	query := r.db.Model(&models.VariantStock{})

	if shouldFilter {
		// Use LEFT JOIN to include product groups (pg_xxx) which don't have corresponding products
		// Filter for: (real products created by user) OR (product groups)
		query = query.
			Joins("LEFT JOIN products ON variant_stocks.product_id = products.id").
			Where("(products.created_by = ?) OR (variant_stocks.product_id LIKE 'pg_%')", fmt.Sprintf("%d", userID))
	}

	err := query.
		Select(`
			product_id,
			product_name,
			COALESCE(SUM(current_stock), 0) AS current_stock,
			COALESCE(SUM(available_stock), 0) AS available_stock,
			COALESCE(SUM(reserved_stock), 0) AS reserved_stock,
			COALESCE(SUM(purchased_stock), 0) AS purchased_stock,
			COALESCE(SUM(sold_stock), 0) AS sold_stock,
			COALESCE(AVG(average_cost), 0) AS average_cost,
			COALESCE(SUM(revaluation_amount), 0) AS revaluation_amount,
			MAX(last_purchased_date) AS last_purchased_date,
			MAX(last_sold_date) AS last_sold_date
		`).
		Group("product_id, product_name").
		Scan(&products).Error

	if err != nil {
		fmt.Printf("GetProductStockDetailsWithFilter Error fetching products: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetProductStockDetailsWithFilter (shouldFilter=%v, userID=%d): Found %d products\n", shouldFilter, userID, len(products))

	result := make([]map[string]interface{}, 0)

	for _, product := range products {
		stockStatus := "in_stock"
		if product.AvailableStock <= 0 {
			stockStatus = "out_of_stock"
		} else if product.AvailableStock <= 100 { // default threshold
			stockStatus = "low_stock"
		}

		result = append(result, map[string]interface{}{
			"product_id":          product.ProductID,
			"product_name":        product.ProductName,
			"sku":                 "", // Not available at product level
			"current_stock":       product.CurrentStock,
			"available_stock":     product.AvailableStock,
			"reserved_stock":      product.ReservedStock,
			"purchased_stock":     product.PurchasedStock,
			"sold_stock":          product.SoldStock,
			"average_cost":        product.AverageCost,
			"revaluation_amount":  product.RevaluationAmount,
			"last_purchased_date": product.LastPurchasedDate,
			"last_sold_date":      product.LastSoldDate,
			"status":              stockStatus,
		})
	}

	fmt.Printf("GetProductStockDetailsWithFilter: Returning %d results\n", len(result))
	return result, nil
}

func (r *dashboardRepository) GetTotalProductsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.VariantStock{}).Distinct("product_id")

	if shouldFilter {
		query = query.
			Joins("JOIN products ON variant_stocks.product_id = products.id").
			Where("products.created_by = ?", fmt.Sprintf("%d", userID))
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetInStockProductsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.VariantStock{}).
		Select("DISTINCT product_id").
		Where("available_stock > 0")

	if shouldFilter {
		query = query.
			Joins("JOIN products ON variant_stocks.product_id = products.id").
			Where("products.created_by = ?", fmt.Sprintf("%d", userID))
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetLowStockProductsWithFilter(threshold float64, shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.VariantStock{}).
		Select("DISTINCT product_id").
		Where("available_stock > 0 AND available_stock <= ?", threshold)

	if shouldFilter {
		query = query.
			Joins("JOIN products ON variant_stocks.product_id = products.id").
			Where("products.created_by = ?", fmt.Sprintf("%d", userID))
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOutOfStockProductsWithFilter(shouldFilter bool, userID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.VariantStock{}).
		Select("DISTINCT product_id").
		Where("available_stock <= 0")

	if shouldFilter {
		query = query.
			Joins("JOIN products ON variant_stocks.product_id = products.id").
			Where("products.created_by = ?", fmt.Sprintf("%d", userID))
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalProductStockWithFilter(shouldFilter bool, userID uint) (float64, error) {
	var total float64
	query := r.db.Model(&models.VariantStock{})

	if shouldFilter {
		query = query.
			Joins("JOIN products ON variant_stocks.product_id = products.id").
			Where("products.created_by = ?", fmt.Sprintf("%d", userID))
	}

	err := query.
		Select("COALESCE(SUM(current_stock), 0)").
		Row().
		Scan(&total)
	return total, err
}
