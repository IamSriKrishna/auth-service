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

	// Vendor metrics
	GetTotalVendors() (int64, error)
	GetActiveVendors() (int64, error)
	GetVendorsCreatedToday() (int64, error)

	// Item metrics
	GetTotalItems() (int64, error)
	GetTotalItemGroups() (int64, error)
	GetTotalStock() (int64, error)
	GetLowStockItems(threshold int64) (int64, error)
	GetOutOfStockItems() (int64, error)
	GetItemsCreatedToday() (int64, error)
	GetItemStockDetails() ([]map[string]interface{}, error)

	// Shipment metrics
	GetTotalShipments() (int64, error)
	GetShippedCount() (int64, error)
	GetShipmentsByStatus(status string) (int64, error)
	GetShippedToday() (int64, error)

	// Invoice metrics
	GetTotalInvoices() (int64, error)
	GetTotalInvoiceAmount() (float64, error)
	GetOutstandingInvoices() (float64, error)
	GetInvoicesByStatus(status string) (int64, error)
	GetOverdueInvoices() (int64, error)

	// Sales order metrics
	GetTotalSalesOrders() (int64, error)
	GetTotalSalesOrderAmount() (float64, error)
	GetSalesOrdersByStatus(status string) (int64, error)
	GetSalesOrdersCreatedToday() (int64, error)

	// Purchase order metrics
	GetTotalPurchaseOrders() (int64, error)
	GetTotalPurchaseOrderAmount() (float64, error)
	GetPurchaseOrdersByStatus(status string) (int64, error)
	GetPurchaseOrdersCreatedToday() (int64, error)

	// Package metrics
	GetTotalPackages() (int64, error)
	GetPackagesByStatus(status string) (int64, error)
	GetPackagesCreatedToday() (int64, error)

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
