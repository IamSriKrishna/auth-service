package repo

import (
	"fmt"
	"strconv"
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
	GetTotalCustomersWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetActiveCustomersWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetCustomersCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Vendor metrics
	GetTotalVendors() (int64, error)
	GetActiveVendors() (int64, error)
	GetVendorsCreatedToday() (int64, error)
	GetTotalVendorsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetActiveVendorsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetVendorsCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Item metrics
	GetTotalItems() (int64, error)
	GetTotalItemGroups() (int64, error)
	GetTotalStock() (int64, error)
	GetLowStockItems(threshold int64) (int64, error)
	GetOutOfStockItems() (int64, error)
	GetItemsCreatedToday() (int64, error)
	GetItemStockDetails() ([]map[string]interface{}, error)
	GetItemStockDetailsWithFilter(shouldFilter bool, companyID uint) ([]map[string]interface{}, error)
	GetTotalItemsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetTotalItemGroupsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetTotalStockWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetLowStockItemsWithFilter(threshold int64, shouldFilter bool, companyID uint) (int64, error)
	GetOutOfStockItemsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetItemsCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Shipment metrics
	GetTotalShipments() (int64, error)
	GetShippedCount() (int64, error)
	GetShipmentsByStatus(status string) (int64, error)
	GetShippedToday() (int64, error)
	GetTotalShipmentsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetShippedCountWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetShippedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetShipmentsByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error)

	// Invoice metrics
	GetTotalInvoices() (int64, error)
	GetTotalInvoiceAmount() (float64, error)
	GetOutstandingInvoices() (float64, error)
	GetInvoicesByStatus(status string) (int64, error)
	GetOverdueInvoices() (int64, error)
	GetTotalInvoicesWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetTotalInvoiceAmountWithFilter(shouldFilter bool, companyID uint) (float64, error)
	GetOutstandingInvoicesWithFilter(shouldFilter bool, companyID uint) (float64, error)
	GetInvoicesByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error)
	GetOverdueInvoicesWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Sales order metrics
	GetTotalSalesOrders() (int64, error)
	GetTotalSalesOrderAmount() (float64, error)
	GetSalesOrdersByStatus(status string) (int64, error)
	GetSalesOrdersCreatedToday() (int64, error)
	GetTotalSalesOrdersWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetTotalSalesOrderAmountWithFilter(shouldFilter bool, companyID uint) (float64, error)
	GetSalesOrdersByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error)
	GetSalesOrdersCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Purchase order metrics
	GetTotalPurchaseOrders() (int64, error)
	GetTotalPurchaseOrderAmount() (float64, error)
	GetPurchaseOrdersByStatus(status string) (int64, error)
	GetPurchaseOrdersCreatedToday() (int64, error)
	GetTotalPurchaseOrdersWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetTotalPurchaseOrderAmountWithFilter(shouldFilter bool, companyID uint) (float64, error)
	GetPurchaseOrdersByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error)
	GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Package metrics
	GetTotalPackages() (int64, error)
	GetPackagesByStatus(status string) (int64, error)
	GetPackagesCreatedToday() (int64, error)
	GetTotalPackagesWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetPackagesByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error)
	GetPackagesCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error)

	// Product stock metrics
	GetTotalProducts() (int64, error)
	GetInStockProducts() (int64, error)
	GetLowStockProducts(threshold float64) (int64, error)
	GetOutOfStockProducts() (int64, error)
	GetTotalProductStock() (float64, error)
	GetProductStockDetails() ([]map[string]interface{}, error)
	GetProductStockDetailsWithFilter(shouldFilter bool, userID uint, companyID uint) ([]map[string]interface{}, error)
	GetTotalProductsWithFilter(shouldFilter bool, companyID uint) (int64, error)
	GetInStockProductsWithFilter(shouldFilter bool, userID uint, companyID uint) (int64, error)
	GetLowStockProductsWithFilter(threshold float64, shouldFilter bool, userID uint, companyID uint) (int64, error)
	GetOutOfStockProductsWithFilter(shouldFilter bool, userID uint, companyID uint) (int64, error)
	GetTotalProductStockWithFilter(shouldFilter bool, userID uint, companyID uint) (float64, error)

	// Shipment tracking
	AddShipmentTracking(tracking *models.ShipmentTracking) error
	GetShipmentTracking(shipmentID string, limit int) ([]models.ShipmentTracking, error)
	GetLatestShipmentTracking(shipmentID string) (*models.ShipmentTracking, error)
	AddShipmentTrackingWithFilter(tracking *models.ShipmentTracking, shouldFilter bool, companyID uint) error
	GetShipmentTrackingWithFilter(shipmentID string, limit int, shouldFilter bool, companyID uint) ([]models.ShipmentTracking, error)

	// History
	GetEntityCountHistory(entityType string, days int) ([]models.EntityCountHistory, error)
	GetEntityCountHistoryWithFilter(entityType string, days int, shouldFilter bool, companyID uint) ([]models.EntityCountHistory, error)
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

// GetItemStockDetailsWithFilter retrieves stock details for items filtered by company
func (r *dashboardRepository) GetItemStockDetailsWithFilter(shouldFilter bool, companyID uint) ([]map[string]interface{}, error) {
	// First, let's get items with optional filtering
	var items []models.Item
	query := r.db

	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}

	err := query.Find(&items).Error

	if err != nil {
		fmt.Printf("GetItemStockDetailsWithFilter Error fetching items: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetItemStockDetailsWithFilter: Found %d items (shouldFilter=%v, companyID=%d)\n", len(items), shouldFilter, companyID)

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

func (r *dashboardRepository) GetShippedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Shipment{}).
		Where("status = ? AND updated_at >= ?", "shipped", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
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

// AddShipmentTrackingWithFilter verifies that the shipment belongs to the
// authenticated company before inserting a tracking record.
func (r *dashboardRepository) AddShipmentTrackingWithFilter(
	tracking *models.ShipmentTracking,
	shouldFilter bool,
	companyID uint,
) error {
	shipmentQuery := r.db.Model(&models.Shipment{}).
		Where("id = ?", tracking.ShipmentID)

	if shouldFilter {
		shipmentQuery = shipmentQuery.Where("company_id = ?", companyID)
	}

	var count int64
	if err := shipmentQuery.Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.Create(tracking).Error
}

// GetShipmentTrackingWithFilter returns tracking only when the shipment belongs
// to the authenticated company.
func (r *dashboardRepository) GetShipmentTrackingWithFilter(
	shipmentID string,
	limit int,
	shouldFilter bool,
	companyID uint,
) ([]models.ShipmentTracking, error) {
	if limit < 1 {
		limit = 10
	}

	shipmentQuery := r.db.Model(&models.Shipment{}).
		Where("id = ?", shipmentID)

	if shouldFilter {
		shipmentQuery = shipmentQuery.Where("company_id = ?", companyID)
	}

	var count int64
	if err := shipmentQuery.Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var tracking []models.ShipmentTracking
	err := r.db.
		Where("shipment_id = ?", shipmentID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&tracking).Error

	return tracking, err
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

func (r *dashboardRepository) GetEntityCountHistoryWithFilter(
	entityType string,
	days int,
	shouldFilter bool,
	companyID uint,
) ([]models.EntityCountHistory, error) {
	var history []models.EntityCountHistory
	startDate := time.Now().AddDate(0, 0, -days)

	query := r.db.Model(&models.EntityCountHistory{}).
		Where("entity_type = ? AND date >= ?", entityType, startDate)

	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}

	err := query.
		Order("date ASC").
		Find(&history).Error

	return history, err
}

func (r *dashboardRepository) SaveEntityCountHistory(history *models.EntityCountHistory) error {
	return r.db.Save(history).Error
}

// ==================== Filter Methods for Company-based Dashboard ====================

// Customer metrics with filter
func (r *dashboardRepository) GetTotalCustomersWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Customer{})
	if shouldFilter && companyID > 0 {
		fmt.Printf("Filtering customers by company_id=%d\n", companyID)
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	fmt.Printf("GetTotalCustomersWithFilter result: count=%d, filter=%v, companyID=%d\n", count, shouldFilter, companyID)
	return count, err
}

func (r *dashboardRepository) GetActiveCustomersWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Customer{}).Where("is_active = ?", true)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetCustomersCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Customer{}).Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Vendor metrics with filter
func (r *dashboardRepository) GetTotalVendorsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Vendor{})
	if shouldFilter && companyID > 0 {
		fmt.Printf("Filtering vendors by company_id=%d\n", companyID)
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	fmt.Printf("GetTotalVendorsWithFilter result: count=%d, filter=%v, companyID=%d\n", count, shouldFilter, companyID)
	return count, err
}

func (r *dashboardRepository) GetActiveVendorsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Vendor{}).Where("is_active = ?", true)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetVendorsCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Vendor{}).Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Item metrics with filter
func (r *dashboardRepository) GetTotalItemsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Item{})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalItemGroupsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.ItemGroup{})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalStockWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var total int64
	query := r.db.Model(&models.InventoryBalance{})
	if shouldFilter {
		// For inventory, we may need to join with items to filter by company
		query = query.Joins("LEFT JOIN items ON items.id = inventory_balances.item_id").
			Where("items.company_id = ?", companyID)
	}
	err := query.Select("COALESCE(SUM(quantity), 0)").Row().Scan(&total)
	return total, err
}

func (r *dashboardRepository) GetLowStockItemsWithFilter(threshold int64, shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.InventoryBalance{}).
		Where("quantity > 0 AND quantity <= ?", threshold)
	if shouldFilter {
		query = query.Joins("LEFT JOIN items ON items.id = inventory_balances.item_id").
			Where("items.company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOutOfStockItemsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.InventoryBalance{}).
		Where("quantity <= 0")
	if shouldFilter {
		query = query.Joins("LEFT JOIN items ON items.id = inventory_balances.item_id").
			Where("items.company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetItemsCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Item{}).Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Shipment metrics with filter
func (r *dashboardRepository) GetTotalShipmentsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Shipment{})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShippedCountWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Shipment{}).
		Where("status IN ?", []string{"shipped", "delivered"})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetShipmentsByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Shipment{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Invoice metrics with filter
func (r *dashboardRepository) GetTotalInvoicesWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Invoice{})
	if shouldFilter && companyID > 0 {
		companyIDStr := companyID
		fmt.Printf("Filtering invoices by company_id=%s\n", companyIDStr)
		query = query.Where("company_id = ?", companyIDStr)
	}
	err := query.Count(&count).Error
	fmt.Printf("GetTotalInvoicesWithFilter result: count=%d, filter=%v, companyID=%d\n", count, shouldFilter, companyID)
	return count, err
}

func (r *dashboardRepository) GetTotalInvoiceAmountWithFilter(shouldFilter bool, companyID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.Invoice{})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Select("COALESCE(SUM(total), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetOutstandingInvoicesWithFilter(shouldFilter bool, companyID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.Invoice{}).
		Where("status IN ?", []string{"pending", "partially_paid"})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Select("COALESCE(SUM(total), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetInvoicesByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Invoice{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOverdueInvoicesWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Invoice{}).
		Where("status IN ? AND due_date < ?", []string{"pending", "partially_paid"}, time.Now())
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Sales order metrics with filter
func (r *dashboardRepository) GetTotalSalesOrdersWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.SalesOrder{})
	if shouldFilter {
		query = query.
			Joins("JOIN customers ON customers.id = sales_orders.customer_id").
			Where("customers.company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalSalesOrderAmountWithFilter(shouldFilter bool, companyID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.SalesOrder{})
	if shouldFilter {
		query = query.
			Joins("JOIN customers ON customers.id = sales_orders.customer_id").
			Where("customers.company_id = ?", companyID)
	}
	err := query.Select("COALESCE(SUM(sales_orders.total), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetSalesOrdersByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.SalesOrder{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetSalesOrdersCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.SalesOrder{}).
		Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Purchase order metrics with filter
func (r *dashboardRepository) GetTotalPurchaseOrdersWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.PurchaseOrder{})
	if shouldFilter {
		query = query.
			Joins("JOIN vendors ON vendors.id = purchase_orders.vendor_id").
			Where("vendors.company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalPurchaseOrderAmountWithFilter(shouldFilter bool, companyID uint) (float64, error) {
	var amount float64
	query := r.db.Model(&models.PurchaseOrder{})
	if shouldFilter {
		query = query.
			Joins("JOIN vendors ON vendors.id = purchase_orders.vendor_id").
			Where("vendors.company_id = ?", companyID)
	}
	err := query.Select("COALESCE(SUM(purchase_orders.total), 0)").Row().Scan(&amount)
	return amount, err
}

func (r *dashboardRepository) GetPurchaseOrdersByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.PurchaseOrder{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.PurchaseOrder{}).
		Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

// Package metrics with filter
func (r *dashboardRepository) GetTotalPackagesWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Package{})
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPackagesByStatusWithFilter(status string, shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Package{}).
		Where("status = ?", status)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetPackagesCreatedTodayWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.Model(&models.Package{}).
		Where("created_at >= ?", today)
	if shouldFilter {
		query = query.Where("company_id = ?", companyID)
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

func applyVariantStockOwnershipFilter(
	query *gorm.DB,
	shouldFilter bool,
	userID uint,
	companyID uint,
) *gorm.DB {
	if !shouldFilter {
		return query
	}

	return query.
		Joins("INNER JOIN products ON products.id = variant_stocks.product_id").
		Where(
			"products.created_by = ? AND products.created_by_company_id = ?",
			strconv.FormatUint(uint64(userID), 10),
			companyID,
		)
}

func (r *dashboardRepository) GetProductStockDetailsWithFilter(
	shouldFilter bool,
	userID uint,
	companyID uint,
) ([]map[string]interface{}, error) {
	type ProductStockAggregate struct {
		ProductID         string     `gorm:"column:product_id"`
		ProductName       string     `gorm:"column:product_name"`
		CurrentStock      float64    `gorm:"column:current_stock"`
		AvailableStock    float64    `gorm:"column:available_stock"`
		ReservedStock     float64    `gorm:"column:reserved_stock"`
		PurchasedStock    float64    `gorm:"column:purchased_stock"`
		SoldStock         float64    `gorm:"column:sold_stock"`
		AverageCost       float64    `gorm:"column:average_cost"`
		RevaluationAmount float64    `gorm:"column:revaluation_amount"`
		LastPurchasedDate *time.Time `gorm:"column:last_purchased_date"`
		LastSoldDate      *time.Time `gorm:"column:last_sold_date"`
		RawMaterialUnit   string     `gorm:"column:raw_material_unit"`
		IsRaw             bool       `gorm:"column:is_raw"`
		RawName           string     `gorm:"column:raw_name"`
		RawSpecification  string     `gorm:"column:raw_specification"`
	}

	var products []ProductStockAggregate

	query := r.db.Table("variant_stocks")
	query = applyVariantStockOwnershipFilter(query, shouldFilter, userID, companyID)

	err := query.
		Select(`
			variant_stocks.product_id AS product_id,
			variant_stocks.product_name AS product_name,
			COALESCE(SUM(variant_stocks.current_stock), 0) AS current_stock,
			COALESCE(SUM(variant_stocks.available_stock), 0) AS available_stock,
			COALESCE(SUM(variant_stocks.reserved_stock), 0) AS reserved_stock,
			COALESCE(SUM(variant_stocks.purchased_stock), 0) AS purchased_stock,
			COALESCE(SUM(variant_stocks.sold_stock), 0) AS sold_stock,
			COALESCE(AVG(variant_stocks.average_cost), 0) AS average_cost,
			COALESCE(SUM(variant_stocks.revaluation_amount), 0) AS revaluation_amount,
			MAX(variant_stocks.last_purchased_date) AS last_purchased_date,
			MAX(variant_stocks.last_sold_date) AS last_sold_date,
			COALESCE(products.raw_unit, '') AS raw_material_unit,
			COALESCE(products.is_raw, false) AS is_raw,
			COALESCE(products.raw_name, '') AS raw_name,
			COALESCE(products.raw_specification, '') AS raw_specification
		`).
		Group(`
			variant_stocks.product_id,
			variant_stocks.product_name,
			products.raw_unit,
			products.is_raw,
			products.raw_name,
			products.raw_specification
		`).
		Scan(&products).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(products))
	for _, product := range products {
		status := "in_stock"
		if product.AvailableStock <= 0 {
			status = "out_of_stock"
		} else if product.AvailableStock <= 100 {
			status = "low_stock"
		}

		result = append(result, map[string]interface{}{
			"product_id":          product.ProductID,
			"product_name":        product.ProductName,
			"sku":                 "",
			"current_stock":       product.CurrentStock,
			"available_stock":     product.AvailableStock,
			"reserved_stock":      product.ReservedStock,
			"purchased_stock":     product.PurchasedStock,
			"sold_stock":          product.SoldStock,
			"average_cost":        product.AverageCost,
			"revaluation_amount":  product.RevaluationAmount,
			"last_purchased_date": product.LastPurchasedDate,
			"last_sold_date":      product.LastSoldDate,
			"raw_material_unit":   product.RawMaterialUnit,
			"is_raw":              product.IsRaw,
			"raw_name":            product.RawName,
			"raw_specification":   product.RawSpecification,
			"status":              status,
		})
	}

	return result, nil
}

func (r *dashboardRepository) GetTotalProductsWithFilter(shouldFilter bool, companyID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.VariantStock{}).Distinct("product_id")

	if shouldFilter {
		query = query.
			Joins("JOIN products ON variant_stocks.product_id = products.id").
			Where("products.created_by_company_id = ?", companyID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetInStockProductsWithFilter(
	shouldFilter bool,
	userID uint,
	companyID uint,
) (int64, error) {
	var count int64
	query := r.db.Table("variant_stocks").Where("variant_stocks.available_stock > 0")
	query = applyVariantStockOwnershipFilter(query, shouldFilter, userID, companyID)
	err := query.Distinct("variant_stocks.product_id").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetLowStockProductsWithFilter(
	threshold float64,
	shouldFilter bool,
	userID uint,
	companyID uint,
) (int64, error) {
	var count int64
	query := r.db.Table("variant_stocks").
		Where("variant_stocks.available_stock > 0 AND variant_stocks.available_stock <= ?", threshold)
	query = applyVariantStockOwnershipFilter(query, shouldFilter, userID, companyID)
	err := query.Distinct("variant_stocks.product_id").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetOutOfStockProductsWithFilter(
	shouldFilter bool,
	userID uint,
	companyID uint,
) (int64, error) {
	var count int64
	query := r.db.Table("variant_stocks").Where("variant_stocks.available_stock <= 0")
	query = applyVariantStockOwnershipFilter(query, shouldFilter, userID, companyID)
	err := query.Distinct("variant_stocks.product_id").Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetTotalProductStockWithFilter(
	shouldFilter bool,
	userID uint,
	companyID uint,
) (float64, error) {
	var total float64
	query := r.db.Table("variant_stocks")
	query = applyVariantStockOwnershipFilter(query, shouldFilter, userID, companyID)
	err := query.Select("COALESCE(SUM(variant_stocks.current_stock), 0)").Scan(&total).Error
	return total, err
}
