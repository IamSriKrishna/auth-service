package services

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type DashboardService interface {
	// Get dashboard metrics
	GetDashboardMetrics() (*output.DashboardMetricsOutput, error)
	GetDashboardMetricsWithUserContext(userID uint, userType string, email string, companyID uint) (*output.DashboardMetricsOutput, error)
	RefreshDashboardMetrics() error

	// Shipment tracking
	AddShipmentTracking(shipmentID, status, location string, latitude, longitude float64, notes string) error
	GetShipmentTracking(shipmentID string, limit int) (*output.ShipmentTrackingListOutput, error)

	// Stock information
	GetStockSummary() (*output.StockListOutput, error)
	GetStockSummaryWithUserContext(userID uint, userType string) (*output.StockListOutput, error)

	// Trends
	GetEntityTrends(entityType string, days int) (*output.EntityTrendOutput, error)

	// Activity
	GetActivitySummary() (*output.ActivitySummaryOutput, error)
	GetActivitySummaryWithUserContext(userID uint, userType string) (*output.ActivitySummaryOutput, error)

	// Diagnostics
	GetDiagnosticReport() (*output.DiagnosticReportOutput, error)
}

type dashboardService struct {
	repo        repo.DashboardRepository
	userRepo    repo.UserRepository
	companyRepo repo.CompanyRepository
}

func NewDashboardService(dashRepo repo.DashboardRepository, userRepo repo.UserRepository, companyRepo repo.CompanyRepository) DashboardService {
	return &dashboardService{
		repo:        dashRepo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
	}
}

// GetDashboardMetrics returns current dashboard metrics
func (s *dashboardService) GetDashboardMetrics() (*output.DashboardMetricsOutput, error) {
	// Customer metrics
	totalCustomers, _ := s.repo.GetTotalCustomers()
	activeCustomers, _ := s.repo.GetActiveCustomers()
	customersCreatedToday, _ := s.repo.GetCustomersCreatedToday()

	// Vendor metrics
	totalVendors, _ := s.repo.GetTotalVendors()
	activeVendors, _ := s.repo.GetActiveVendors()
	vendorsCreatedToday, _ := s.repo.GetVendorsCreatedToday()

	// Item metrics
	totalItems, _ := s.repo.GetTotalItems()
	totalItemGroups, _ := s.repo.GetTotalItemGroups()
	totalStock, _ := s.repo.GetTotalStock()
	lowStockItems, _ := s.repo.GetLowStockItems(100)
	outOfStockItems, _ := s.repo.GetOutOfStockItems()
	itemsCreatedToday, _ := s.repo.GetItemsCreatedToday()

	// Shipment metrics
	totalShipments, _ := s.repo.GetTotalShipments()
	shippedCount, _ := s.repo.GetShippedCount()
	pendingShipments, _ := s.repo.GetShipmentsByStatus("pending")
	deliveredShipments, _ := s.repo.GetShipmentsByStatus("delivered")
	inTransitShipments, _ := s.repo.GetShipmentsByStatus("in_transit")

	// Invoice metrics
	totalInvoices, _ := s.repo.GetTotalInvoices()
	totalInvoiceAmount, _ := s.repo.GetTotalInvoiceAmount()
	outstandingInvoices, _ := s.repo.GetOutstandingInvoices()
	paidInvoices, _ := s.repo.GetInvoicesByStatus("paid")
	pendingInvoices, _ := s.repo.GetInvoicesByStatus("pending")
	overdueInvoices, _ := s.repo.GetOverdueInvoices()

	// Sales order metrics
	totalSalesOrders, _ := s.repo.GetTotalSalesOrders()
	totalSalesOrderAmount, _ := s.repo.GetTotalSalesOrderAmount()
	completedSalesOrders, _ := s.repo.GetSalesOrdersByStatus("completed")
	pendingSalesOrders, _ := s.repo.GetSalesOrdersByStatus("pending")
	cancelledSalesOrders, _ := s.repo.GetSalesOrdersByStatus("cancelled")
	salesOrdersCreatedToday, _ := s.repo.GetSalesOrdersCreatedToday()

	// Purchase order metrics
	totalPurchaseOrders, _ := s.repo.GetTotalPurchaseOrders()
	totalPurchaseOrderAmount, _ := s.repo.GetTotalPurchaseOrderAmount()
	completedPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatus("completed")
	pendingPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatus("pending")
	cancelledPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatus("cancelled")
	purchaseOrdersCreatedToday, _ := s.repo.GetPurchaseOrdersCreatedToday()

	// Package metrics
	totalPackages, _ := s.repo.GetTotalPackages()
	shippedPackages, _ := s.repo.GetPackagesByStatus("shipped")
	pendingPackages, _ := s.repo.GetPackagesByStatus("pending")
	inTransitPackages, _ := s.repo.GetPackagesByStatus("in_transit")
	deliveredPackages, _ := s.repo.GetPackagesByStatus("delivered")
	packagesCreatedToday, _ := s.repo.GetPackagesCreatedToday()

	return &output.DashboardMetricsOutput{
		CustomerMetrics: output.CustomerMetricsOutput{
			Total:        int(totalCustomers),
			Active:       int(activeCustomers),
			Inactive:     int(totalCustomers - activeCustomers),
			CreatedToday: int(customersCreatedToday),
		},
		VendorMetrics: output.VendorMetricsOutput{
			Total:        int(totalVendors),
			Active:       int(activeVendors),
			Inactive:     int(totalVendors - activeVendors),
			CreatedToday: int(vendorsCreatedToday),
		},
		ItemMetrics: output.ItemMetricsOutput{
			Total:          int(totalItems),
			TotalStock:     totalStock,
			LowStockItems:  int(lowStockItems),
			ItemGroups:     int(totalItemGroups),
			CreatedToday:   int(itemsCreatedToday),
			OutOfStockItem: int(outOfStockItems),
		},
		ShipmentMetrics: output.ShipmentMetricsOutput{
			Total:            int(totalShipments),
			Shipped:          int(shippedCount),
			Pending:          int(pendingShipments),
			InTransit:        int(inTransitShipments),
			Delivered:        int(deliveredShipments),
			CancelledShipped: 0,
			AverageTime:      0,
		},
		InvoiceMetrics: output.InvoiceMetricsOutput{
			Total:       int(totalInvoices),
			TotalAmount: totalInvoiceAmount,
			Outstanding: outstandingInvoices,
			Paid:        int(paidInvoices),
			Pending:     int(pendingInvoices),
			Overdue:     int(overdueInvoices),
		},
		SalesOrderMetrics: output.SalesOrderMetricsOutput{
			Total:        int(totalSalesOrders),
			TotalAmount:  totalSalesOrderAmount,
			Completed:    int(completedSalesOrders),
			Pending:      int(pendingSalesOrders),
			Cancelled:    int(cancelledSalesOrders),
			CreatedToday: int(salesOrdersCreatedToday),
		},
		PurchaseOrderMetrics: output.PurchaseOrderMetricsOutput{
			Total:        int(totalPurchaseOrders),
			TotalAmount:  totalPurchaseOrderAmount,
			Completed:    int(completedPurchaseOrders),
			Pending:      int(pendingPurchaseOrders),
			Cancelled:    int(cancelledPurchaseOrders),
			CreatedToday: int(purchaseOrdersCreatedToday),
		},
		PackageMetrics: output.PackageMetricsOutput{
			Total:        int(totalPackages),
			Shipped:      int(shippedPackages),
			Pending:      int(pendingPackages),
			InTransit:    int(inTransitPackages),
			Delivered:    int(deliveredPackages),
			CreatedToday: int(packagesCreatedToday),
		},
		LastUpdatedAt: time.Now(),
		GeneratedAt:   time.Now(),
	}, nil
}

// GetDashboardMetricsWithUserContext returns dashboard metrics filtered by user role
func (s *dashboardService) GetDashboardMetricsWithUserContext(userID uint, userType string, email string, companyID uint) (*output.DashboardMetricsOutput, error) {
	// For superadmin, show all data; for admin, filter by user ID
	shouldFilter := userType == "admin" && userID > 0

	fmt.Printf("Service Debug - UserID: %d, UserType: %s, ShouldFilter: %v, CompanyID: %d\n", userID, userType, shouldFilter, companyID)

	// Customer metrics
	totalCustomers, _ := s.repo.GetTotalCustomersWithFilter(shouldFilter, userID)
	activeCustomers, _ := s.repo.GetActiveCustomersWithFilter(shouldFilter, userID)
	customersCreatedToday, _ := s.repo.GetCustomersCreatedTodayWithFilter(shouldFilter, userID)

	fmt.Printf("Customer metrics - Total: %d (should filter:%v)\n", totalCustomers, shouldFilter)

	// Vendor metrics
	totalVendors, _ := s.repo.GetTotalVendorsWithFilter(shouldFilter, userID)
	activeVendors, _ := s.repo.GetActiveVendorsWithFilter(shouldFilter, userID)
	vendorsCreatedToday, _ := s.repo.GetVendorsCreatedTodayWithFilter(shouldFilter, userID)

	// Item metrics
	totalItems, _ := s.repo.GetTotalItemsWithFilter(shouldFilter, userID)
	totalItemGroups, _ := s.repo.GetTotalItemGroupsWithFilter(shouldFilter, userID)
	totalStock, _ := s.repo.GetTotalStockWithFilter(shouldFilter, userID)
	lowStockItems, _ := s.repo.GetLowStockItemsWithFilter(100, shouldFilter, userID)
	outOfStockItems, _ := s.repo.GetOutOfStockItemsWithFilter(shouldFilter, userID)
	itemsCreatedToday, _ := s.repo.GetItemsCreatedTodayWithFilter(shouldFilter, userID)

	// Shipment metrics
	totalShipments, _ := s.repo.GetTotalShipmentsWithFilter(shouldFilter, userID)
	shippedCount, _ := s.repo.GetShippedCountWithFilter(shouldFilter, userID)
	pendingShipments, _ := s.repo.GetShipmentsByStatusWithFilter("pending", shouldFilter, userID)
	deliveredShipments, _ := s.repo.GetShipmentsByStatusWithFilter("delivered", shouldFilter, userID)
	inTransitShipments, _ := s.repo.GetShipmentsByStatusWithFilter("in_transit", shouldFilter, userID)

	// Invoice metrics
	totalInvoices, _ := s.repo.GetTotalInvoicesWithFilter(shouldFilter, userID)
	totalInvoiceAmount, _ := s.repo.GetTotalInvoiceAmountWithFilter(shouldFilter, userID)
	outstandingInvoices, _ := s.repo.GetOutstandingInvoicesWithFilter(shouldFilter, userID)
	paidInvoices, _ := s.repo.GetInvoicesByStatusWithFilter("paid", shouldFilter, userID)
	pendingInvoices, _ := s.repo.GetInvoicesByStatusWithFilter("pending", shouldFilter, userID)
	overdueInvoices, _ := s.repo.GetOverdueInvoicesWithFilter(shouldFilter, userID)

	// Sales order metrics
	totalSalesOrders, _ := s.repo.GetTotalSalesOrdersWithFilter(shouldFilter, userID)
	totalSalesOrderAmount, _ := s.repo.GetTotalSalesOrderAmountWithFilter(shouldFilter, userID)
	completedSalesOrders, _ := s.repo.GetSalesOrdersByStatusWithFilter("completed", shouldFilter, userID)
	pendingSalesOrders, _ := s.repo.GetSalesOrdersByStatusWithFilter("pending", shouldFilter, userID)
	cancelledSalesOrders, _ := s.repo.GetSalesOrdersByStatusWithFilter("cancelled", shouldFilter, userID)
	salesOrdersCreatedToday, _ := s.repo.GetSalesOrdersCreatedTodayWithFilter(shouldFilter, userID)

	// Purchase order metrics
	totalPurchaseOrders, _ := s.repo.GetTotalPurchaseOrdersWithFilter(shouldFilter, userID)
	totalPurchaseOrderAmount, _ := s.repo.GetTotalPurchaseOrderAmountWithFilter(shouldFilter, userID)
	completedPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatusWithFilter("completed", shouldFilter, userID)
	pendingPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatusWithFilter("pending", shouldFilter, userID)
	cancelledPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatusWithFilter("cancelled", shouldFilter, userID)
	purchaseOrdersCreatedToday, _ := s.repo.GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter, userID)

	// Package metrics
	totalPackages, _ := s.repo.GetTotalPackagesWithFilter(shouldFilter, userID)
	shippedPackages, _ := s.repo.GetPackagesByStatusWithFilter("shipped", shouldFilter, userID)
	pendingPackages, _ := s.repo.GetPackagesByStatusWithFilter("pending", shouldFilter, userID)
	inTransitPackages, _ := s.repo.GetPackagesByStatusWithFilter("in_transit", shouldFilter, userID)
	deliveredPackages, _ := s.repo.GetPackagesByStatusWithFilter("delivered", shouldFilter, userID)
	packagesCreatedToday, _ := s.repo.GetPackagesCreatedTodayWithFilter(shouldFilter, userID)

	// Prepare user info by fetching from database
	userInfo := output.UserInfoOutput{
		UserID:      userID,
		UserRole:    userType,
		CompanyID:   companyID,
		CompanyName: "",
		UserName:    "",
		Email:       email,
	}

	fmt.Printf("\n=== Service: Fetching User Info ===\n")
	fmt.Printf("UserID: %d, UserType: %s, CompanyID: %d\n", userID, userType, companyID)

	// Fetch actual user data from database
	if userID > 0 && s.userRepo != nil {
		fmt.Printf("Attempting to fetch user %d from database\n", userID)
		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			fmt.Printf("Error fetching user: %v\n", err)
		} else if user != nil {
			fmt.Printf("User found: %+v\n", user)
			if user.Username != nil {
				userInfo.UserName = *user.Username
				fmt.Printf("Set username: %s\n", userInfo.UserName)
			}
			if user.Email != nil {
				userInfo.Email = *user.Email
				fmt.Printf("Set email: %s\n", userInfo.Email)
			}
		} else {
			fmt.Println("User is nil")
		}
	} else {
		fmt.Printf("Skipping user fetch: userID=%d, userRepo=%v\n", userID, s.userRepo != nil)
	}

	// Fetch company data from database
	if companyID > 0 && s.companyRepo != nil {
		fmt.Printf("Attempting to fetch company %d from database\n", companyID)
		company, err := s.companyRepo.FindByID(companyID)
		if err != nil {
			fmt.Printf("Error fetching company: %v\n", err)
		} else if company != nil {
			fmt.Printf("Company found: %+v\n", company)
			userInfo.CompanyName = company.CompanyName
			fmt.Printf("Set company name: %s\n", userInfo.CompanyName)
		} else {
			fmt.Println("Company is nil")
		}
	} else {
		fmt.Printf("Skipping company fetch: companyID=%d, companyRepo=%v\n", companyID, s.companyRepo != nil)
	}

	fmt.Printf("=== Final User Info ===\n")
	fmt.Printf("UserID: %d, UserName: %s, UserRole: %s, CompanyID: %d, CompanyName: %s, Email: %s\n\n",
		userInfo.UserID, userInfo.UserName, userInfo.UserRole, userInfo.CompanyID, userInfo.CompanyName, userInfo.Email)

	return &output.DashboardMetricsOutput{
		UserInfo: userInfo,
		CustomerMetrics: output.CustomerMetricsOutput{
			Total:        int(totalCustomers),
			Active:       int(activeCustomers),
			Inactive:     int(totalCustomers - activeCustomers),
			CreatedToday: int(customersCreatedToday),
		},
		VendorMetrics: output.VendorMetricsOutput{
			Total:        int(totalVendors),
			Active:       int(activeVendors),
			Inactive:     int(totalVendors - activeVendors),
			CreatedToday: int(vendorsCreatedToday),
		},
		ItemMetrics: output.ItemMetricsOutput{
			Total:          int(totalItems),
			TotalStock:     totalStock,
			LowStockItems:  int(lowStockItems),
			ItemGroups:     int(totalItemGroups),
			CreatedToday:   int(itemsCreatedToday),
			OutOfStockItem: int(outOfStockItems),
		},
		ShipmentMetrics: output.ShipmentMetricsOutput{
			Total:            int(totalShipments),
			Shipped:          int(shippedCount),
			Pending:          int(pendingShipments),
			InTransit:        int(inTransitShipments),
			Delivered:        int(deliveredShipments),
			CancelledShipped: 0,
			AverageTime:      0,
		},
		InvoiceMetrics: output.InvoiceMetricsOutput{
			Total:       int(totalInvoices),
			TotalAmount: totalInvoiceAmount,
			Outstanding: outstandingInvoices,
			Paid:        int(paidInvoices),
			Pending:     int(pendingInvoices),
			Overdue:     int(overdueInvoices),
		},
		SalesOrderMetrics: output.SalesOrderMetricsOutput{
			Total:        int(totalSalesOrders),
			TotalAmount:  totalSalesOrderAmount,
			Completed:    int(completedSalesOrders),
			Pending:      int(pendingSalesOrders),
			Cancelled:    int(cancelledSalesOrders),
			CreatedToday: int(salesOrdersCreatedToday),
		},
		PurchaseOrderMetrics: output.PurchaseOrderMetricsOutput{
			Total:        int(totalPurchaseOrders),
			TotalAmount:  totalPurchaseOrderAmount,
			Completed:    int(completedPurchaseOrders),
			Pending:      int(pendingPurchaseOrders),
			Cancelled:    int(cancelledPurchaseOrders),
			CreatedToday: int(purchaseOrdersCreatedToday),
		},
		PackageMetrics: output.PackageMetricsOutput{
			Total:        int(totalPackages),
			Shipped:      int(shippedPackages),
			Pending:      int(pendingPackages),
			InTransit:    int(inTransitPackages),
			Delivered:    int(deliveredPackages),
			CreatedToday: int(packagesCreatedToday),
		},
		LastUpdatedAt: time.Now(),
		GeneratedAt:   time.Now(),
	}, nil
}

// RefreshDashboardMetrics refreshes the cached metrics
func (s *dashboardService) RefreshDashboardMetrics() error {
	totalCustomers, _ := s.repo.GetTotalCustomers()
	activeCustomers, _ := s.repo.GetActiveCustomers()
	totalVendors, _ := s.repo.GetTotalVendors()
	activeVendors, _ := s.repo.GetActiveVendors()
	totalItems, _ := s.repo.GetTotalItems()
	totalItemGroups, _ := s.repo.GetTotalItemGroups()
	totalStock, _ := s.repo.GetTotalStock()
	lowStockItems, _ := s.repo.GetLowStockItems(100)
	totalShipments, _ := s.repo.GetTotalShipments()
	shippedCount, _ := s.repo.GetShippedCount()
	pendingShipments, _ := s.repo.GetShipmentsByStatus("pending")
	totalInvoices, _ := s.repo.GetTotalInvoices()
	totalInvoiceAmount, _ := s.repo.GetTotalInvoiceAmount()
	totalSalesOrders, _ := s.repo.GetTotalSalesOrders()
	totalSalesOrderAmount, _ := s.repo.GetTotalSalesOrderAmount()
	totalPurchaseOrders, _ := s.repo.GetTotalPurchaseOrders()
	totalPurchaseOrderAmount, _ := s.repo.GetTotalPurchaseOrderAmount()
	totalPackages, _ := s.repo.GetTotalPackages()
	shippedPackages, _ := s.repo.GetPackagesByStatus("shipped")
	pendingPackages, _ := s.repo.GetPackagesByStatus("pending")

	metrics := &models.DashboardMetrics{
		ID:                  uuid.New().String(),
		TotalCustomers:      int(totalCustomers),
		ActiveCustomers:     int(activeCustomers),
		TotalVendors:        int(totalVendors),
		ActiveVendors:       int(activeVendors),
		TotalItems:          int(totalItems),
		TotalItemGroups:     int(totalItemGroups),
		TotalStock:          totalStock,
		LowStockItems:       int(lowStockItems),
		TotalShipments:      int(totalShipments),
		ShippedCount:        int(shippedCount),
		PendingShipments:    int(pendingShipments),
		TotalInvoices:       int(totalInvoices),
		InvoiceAmount:       totalInvoiceAmount,
		TotalSalesOrders:    int(totalSalesOrders),
		SalesOrderAmount:    totalSalesOrderAmount,
		TotalPurchaseOrders: int(totalPurchaseOrders),
		PurchaseOrderAmount: totalPurchaseOrderAmount,
		TotalPackages:       int(totalPackages),
		PackagesShipped:     int(shippedPackages),
		PendingPackages:     int(pendingPackages),
		LastUpdatedAt:       time.Now(),
	}

	return s.repo.SaveDashboardMetrics(metrics)
}

// AddShipmentTracking adds a shipment tracking record
func (s *dashboardService) AddShipmentTracking(shipmentID, status, location string, latitude, longitude float64, notes string) error {
	tracking := &models.ShipmentTracking{
		ID:         uuid.New().String(),
		ShipmentID: shipmentID,
		Status:     status,
		Location:   location,
		Latitude:   latitude,
		Longitude:  longitude,
		Notes:      notes,
		Timestamp:  time.Now(),
	}
	return s.repo.AddShipmentTracking(tracking)
}

// GetShipmentTracking retrieves shipment tracking records
func (s *dashboardService) GetShipmentTracking(shipmentID string, limit int) (*output.ShipmentTrackingListOutput, error) {
	if limit == 0 {
		limit = 10
	}
	tracking, err := s.repo.GetShipmentTracking(shipmentID, limit)
	if err != nil {
		return nil, err
	}

	result := &output.ShipmentTrackingListOutput{
		Data:  make([]output.ShipmentTrackingOutput, len(tracking)),
		Total: len(tracking),
	}

	for i, t := range tracking {
		result.Data[i] = output.ShipmentTrackingOutput{
			ID:         t.ID,
			ShipmentID: t.ShipmentID,
			Status:     t.Status,
			Location:   t.Location,
			Latitude:   t.Latitude,
			Longitude:  t.Longitude,
			Notes:      t.Notes,
			Timestamp:  t.Timestamp,
		}
	}

	return result, nil
}

// GetStockSummary retrieves the stock summary
func (s *dashboardService) GetStockSummary() (*output.StockListOutput, error) {
	inStock, _ := s.repo.GetTotalItems()
	lowStock, _ := s.repo.GetLowStockItems(100)
	outOfStock, _ := s.repo.GetOutOfStockItems()
	totalStock, _ := s.repo.GetTotalStock()

	// Get detailed stock information
	itemStockDetails, err := s.repo.GetItemStockDetails()
	if err != nil {
		fmt.Printf("Error getting item stock details: %v\n", err)
		return &output.StockListOutput{
			Data:          make([]output.StockDetailOutput, 0),
			TotalItems:    int(inStock),
			InStock:       int(inStock) - int(lowStock) - int(outOfStock),
			LowStock:      int(lowStock),
			OutOfStock:    int(outOfStock),
			TotalQuantity: totalStock,
		}, nil
	}

	stockDetails := make([]output.StockDetailOutput, 0)
	fmt.Printf("DEBUG: itemStockDetails count: %d, nil: %v\n", len(itemStockDetails), itemStockDetails == nil)

	if len(itemStockDetails) > 0 {
		for _, item := range itemStockDetails {
			fmt.Printf("DEBUG: Processing item: %v\n", item)
			stockDetails = append(stockDetails, output.StockDetailOutput{
				ItemID:            item["item_id"].(string),
				ItemName:          item["item_name"].(string),
				CurrentQuantity:   item["current_quantity"].(float64),
				AvailableQuantity: item["available_quantity"].(float64),
				ReservedQuantity:  item["reserved_quantity"].(float64),
				InTransitQuantity: item["in_transit_quantity"].(float64),
				Status:            item["status"].(string),
			})
		}
	}

	fmt.Printf("DEBUG: Final stockDetails count: %d\n", len(stockDetails))

	return &output.StockListOutput{
		Data:          stockDetails,
		TotalItems:    int(inStock),
		InStock:       int(inStock) - int(lowStock) - int(outOfStock),
		LowStock:      int(lowStock),
		OutOfStock:    int(outOfStock),
		TotalQuantity: totalStock,
	}, nil
}

// GetStockSummaryWithUserContext retrieves stock summary with user filtering
func (s *dashboardService) GetStockSummaryWithUserContext(userID uint, userType string) (*output.StockListOutput, error) {
	shouldFilter := userType == "admin" && userID > 0

	var inStock, lowStock, outOfStock, totalStock int64

	if shouldFilter {
		inStock, _ = s.repo.GetTotalItemsWithFilter(shouldFilter, userID)
		lowStock, _ = s.repo.GetLowStockItemsWithFilter(100, shouldFilter, userID)
		outOfStock, _ = s.repo.GetOutOfStockItemsWithFilter(shouldFilter, userID)
		totalStock, _ = s.repo.GetTotalStockWithFilter(shouldFilter, userID)
	} else {
		inStock, _ = s.repo.GetTotalItems()
		lowStock, _ = s.repo.GetLowStockItems(100)
		outOfStock, _ = s.repo.GetOutOfStockItems()
		totalStock, _ = s.repo.GetTotalStock()
	}

	// Get detailed stock information
	var itemStockDetails []map[string]interface{}
	var err error

	if shouldFilter {
		itemStockDetails, err = s.repo.GetItemStockDetailsWithFilter(shouldFilter, userID)
	} else {
		itemStockDetails, err = s.repo.GetItemStockDetails()
	}

	if err != nil {
		fmt.Printf("Error getting item stock details: %v\n", err)
		return &output.StockListOutput{
			Data:          make([]output.StockDetailOutput, 0),
			TotalItems:    int(inStock),
			InStock:       int(inStock) - int(lowStock) - int(outOfStock),
			LowStock:      int(lowStock),
			OutOfStock:    int(outOfStock),
			TotalQuantity: totalStock,
		}, nil
	}

	stockDetails := make([]output.StockDetailOutput, 0)
	fmt.Printf("DEBUG: itemStockDetails count: %d, nil: %v\n", len(itemStockDetails), itemStockDetails == nil)

	if len(itemStockDetails) > 0 {
		for _, item := range itemStockDetails {
			fmt.Printf("DEBUG: Processing item: %v\n", item)
			stockDetails = append(stockDetails, output.StockDetailOutput{
				ItemID:            item["item_id"].(string),
				ItemName:          item["item_name"].(string),
				CurrentQuantity:   item["current_quantity"].(float64),
				AvailableQuantity: item["available_quantity"].(float64),
				ReservedQuantity:  item["reserved_quantity"].(float64),
				InTransitQuantity: item["in_transit_quantity"].(float64),
				Status:            item["status"].(string),
			})
		}
	}

	fmt.Printf("DEBUG: Final stockDetails count: %d\n", len(stockDetails))

	return &output.StockListOutput{
		Data:          stockDetails,
		TotalItems:    int(inStock),
		InStock:       int(inStock) - int(lowStock) - int(outOfStock),
		LowStock:      int(lowStock),
		OutOfStock:    int(outOfStock),
		TotalQuantity: totalStock,
	}, nil
}

// GetEntityTrends retrieves trend data for an entity
func (s *dashboardService) GetEntityTrends(entityType string, days int) (*output.EntityTrendOutput, error) {
	history, err := s.repo.GetEntityCountHistory(entityType, days)
	if err != nil {
		return nil, err
	}

	result := &output.EntityTrendOutput{
		EntityType: entityType,
		Data:       make([]output.TrendPoint, len(history)),
	}

	for i, h := range history {
		result.Data[i] = output.TrendPoint{
			Date:        h.Date,
			Count:       h.Count,
			ActiveCount: h.ActiveCount,
			NewToday:    h.CreatedToday,
		}
	}

	return result, nil
}

// GetActivitySummary retrieves activity summary for today
func (s *dashboardService) GetActivitySummary() (*output.ActivitySummaryOutput, error) {
	customersCreatedToday, _ := s.repo.GetCustomersCreatedToday()
	vendorsCreatedToday, _ := s.repo.GetVendorsCreatedToday()
	itemsCreatedToday, _ := s.repo.GetItemsCreatedToday()
	salesOrdersCreatedToday, _ := s.repo.GetSalesOrdersCreatedToday()
	purchaseOrdersCreatedToday, _ := s.repo.GetPurchaseOrdersCreatedToday()
	shippedToday, _ := s.repo.GetShippedToday()
	deliveredToday, _ := s.repo.GetShipmentsByStatus("delivered")

	return &output.ActivitySummaryOutput{
		CreatedCustomersToday:      int(customersCreatedToday),
		CreatedVendorsToday:        int(vendorsCreatedToday),
		CreatedItemsToday:          int(itemsCreatedToday),
		CreatedSalesOrdersToday:    int(salesOrdersCreatedToday),
		CreatedPurchaseOrdersToday: int(purchaseOrdersCreatedToday),
		ShippedToday:               int(shippedToday),
		DeliveredToday:             int(deliveredToday),
	}, nil
}

// GetActivitySummaryWithUserContext retrieves activity summary for today with user filtering
func (s *dashboardService) GetActivitySummaryWithUserContext(userID uint, userType string) (*output.ActivitySummaryOutput, error) {
	shouldFilter := userType == "admin" && userID > 0

	var customersCreatedToday, vendorsCreatedToday, itemsCreatedToday int64
	var salesOrdersCreatedToday, purchaseOrdersCreatedToday, shippedToday, deliveredToday int64

	if shouldFilter {
		customersCreatedToday, _ = s.repo.GetCustomersCreatedTodayWithFilter(shouldFilter, userID)
		vendorsCreatedToday, _ = s.repo.GetVendorsCreatedTodayWithFilter(shouldFilter, userID)
		itemsCreatedToday, _ = s.repo.GetItemsCreatedTodayWithFilter(shouldFilter, userID)
		salesOrdersCreatedToday, _ = s.repo.GetSalesOrdersCreatedTodayWithFilter(shouldFilter, userID)
		purchaseOrdersCreatedToday, _ = s.repo.GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter, userID)
	} else {
		customersCreatedToday, _ = s.repo.GetCustomersCreatedToday()
		vendorsCreatedToday, _ = s.repo.GetVendorsCreatedToday()
		itemsCreatedToday, _ = s.repo.GetItemsCreatedToday()
		salesOrdersCreatedToday, _ = s.repo.GetSalesOrdersCreatedToday()
		purchaseOrdersCreatedToday, _ = s.repo.GetPurchaseOrdersCreatedToday()
	}

	shippedToday, _ = s.repo.GetShippedToday()
	deliveredToday, _ = s.repo.GetShipmentsByStatus("delivered")

	return &output.ActivitySummaryOutput{
		CreatedCustomersToday:      int(customersCreatedToday),
		CreatedVendorsToday:        int(vendorsCreatedToday),
		CreatedItemsToday:          int(itemsCreatedToday),
		CreatedSalesOrdersToday:    int(salesOrdersCreatedToday),
		CreatedPurchaseOrdersToday: int(purchaseOrdersCreatedToday),
		ShippedToday:               int(shippedToday),
		DeliveredToday:             int(deliveredToday),
	}, nil
}

// GetDiagnosticReport generates a diagnostic report for data quality
func (s *dashboardService) GetDiagnosticReport() (*output.DiagnosticReportOutput, error) {
	issues := []string{}
	diagnostics := make(map[string]output.DiagnosticItem)

	// Check customer data
	totalCustomers, _ := s.repo.GetTotalCustomers()
	activeCustomers, _ := s.repo.GetActiveCustomers()
	if totalCustomers > 0 && activeCustomers == 0 {
		issues = append(issues, "No active customers found - check is_active field")
	}
	diagnostics["customers"] = output.DiagnosticItem{
		Label:       "Customers",
		Value:       map[string]interface{}{"total": totalCustomers, "active": activeCustomers},
		Status:      statusCheck(activeCustomers, totalCustomers),
		Description: "Total vs Active customers",
	}

	// Check vendor data
	totalVendors, _ := s.repo.GetTotalVendors()
	activeVendors, _ := s.repo.GetActiveVendors()
	if totalVendors > 0 && activeVendors == 0 {
		issues = append(issues, "No active vendors found - check is_active field")
	}
	diagnostics["vendors"] = output.DiagnosticItem{
		Label:       "Vendors",
		Value:       map[string]interface{}{"total": totalVendors, "active": activeVendors},
		Status:      statusCheck(activeVendors, totalVendors),
		Description: "Total vs Active vendors",
	}

	// Check inventory balance
	totalItems, _ := s.repo.GetTotalItems()
	totalStock, _ := s.repo.GetTotalStock()
	if totalItems > 0 && totalStock == 0 {
		issues = append(issues, "No inventory balance records - inventory_balance table may be empty")
	}
	diagnostics["inventory"] = output.DiagnosticItem{
		Label:       "Inventory",
		Value:       map[string]interface{}{"total_items": totalItems, "total_stock": totalStock},
		Status:      statusCheck(totalStock, totalItems),
		Description: "Items with inventory tracking",
	}

	// Check invoice amounts
	totalInvoices, _ := s.repo.GetTotalInvoices()
	totalInvoiceAmount, _ := s.repo.GetTotalInvoiceAmount()
	if totalInvoices > 0 && totalInvoiceAmount == 0 {
		issues = append(issues, "Invoice amounts are zero - check invoice.total_amount field")
	}
	diagnostics["invoices"] = output.DiagnosticItem{
		Label:       "Invoices",
		Value:       map[string]interface{}{"total": totalInvoices, "total_amount": totalInvoiceAmount},
		Status:      statusCheck(totalInvoiceAmount, totalInvoices),
		Description: "Total invoices with amounts",
	}

	// Check sales order amounts
	totalSalesOrders, _ := s.repo.GetTotalSalesOrders()
	totalSOAmount, _ := s.repo.GetTotalSalesOrderAmount()
	if totalSalesOrders > 0 && totalSOAmount == 0 {
		issues = append(issues, "Sales order amounts are zero - check sales_order.total_amount field")
	}
	diagnostics["sales_orders"] = output.DiagnosticItem{
		Label:       "Sales Orders",
		Value:       map[string]interface{}{"total": totalSalesOrders, "total_amount": totalSOAmount},
		Status:      statusCheck(totalSOAmount, totalSalesOrders),
		Description: "Total sales orders with amounts",
	}

	// Check purchase order amounts
	totalPurchaseOrders, _ := s.repo.GetTotalPurchaseOrders()
	totalPOAmount, _ := s.repo.GetTotalPurchaseOrderAmount()
	if totalPurchaseOrders > 0 && totalPOAmount == 0 {
		issues = append(issues, "Purchase order amounts are zero - check purchase_order.total_amount field")
	}
	diagnostics["purchase_orders"] = output.DiagnosticItem{
		Label:       "Purchase Orders",
		Value:       map[string]interface{}{"total": totalPurchaseOrders, "total_amount": totalPOAmount},
		Status:      statusCheck(totalPOAmount, totalPurchaseOrders),
		Description: "Total purchase orders with amounts",
	}

	// Check shipments
	totalShipments, _ := s.repo.GetTotalShipments()
	diagnostics["shipments"] = output.DiagnosticItem{
		Label:       "Shipments",
		Value:       map[string]interface{}{"total": totalShipments},
		Status:      statusCheck(totalShipments, 1),
		Description: "Total shipments in database",
	}
	if totalShipments == 0 {
		issues = append(issues, "No shipments found - shipment tracking not available")
	}

	// Create summary
	summary := "✓ All systems operational"
	if len(issues) > 0 {
		summary = fmt.Sprintf("⚠️ %d data quality issues detected", len(issues))
	}

	return &output.DiagnosticReportOutput{
		DataIssues:  issues,
		Diagnostics: diagnostics,
		Summary:     summary,
	}, nil
}

func statusCheck(actual, expected interface{}) string {
	actualVal := getValue(actual)
	expectedVal := getValue(expected)

	if actualVal == expectedVal && actualVal > 0 {
		return "ok"
	}
	if actualVal > 0 {
		return "warning"
	}
	if expectedVal > 0 {
		return "error"
	}
	return "ok"
}

func getValue(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}
