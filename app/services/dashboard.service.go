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
	GetDashboardMetrics() (*output.DashboardMetricsOutput, error)

	GetDashboardMetricsWithUserContext(
		userID uint,
		userType string,
		email string,
		companyID uint,
	) (*output.DashboardMetricsOutput, error)

	RefreshDashboardMetrics() error

	AddShipmentTracking(
		shipmentID string,
		status string,
		location string,
		latitude float64,
		longitude float64,
		notes string,
		companyID uint,
		userType string,
	) error

	GetShipmentTracking(
		shipmentID string,
		limit int,
		companyID uint,
		userType string,
	) (*output.ShipmentTrackingListOutput, error)

	GetStockSummary() (*output.StockListOutput, error)

	GetStockSummaryWithUserContext(
		userID uint,
		userType string,
		companyID uint,
	) (*output.StockListOutput, error)

	GetEntityTrends(
		entityType string,
		days int,
		companyID uint,
		userType string,
	) (*output.EntityTrendOutput, error)

	GetActivitySummary() (*output.ActivitySummaryOutput, error)

	GetActivitySummaryWithUserContext(
		userID uint,
		userType string,
		companyID uint,
	) (*output.ActivitySummaryOutput, error)

	GetDiagnosticReport(
		companyID uint,
		userType string,
	) (*output.DiagnosticReportOutput, error)
}

type dashboardService struct {
	repo        repo.DashboardRepository
	userRepo    repo.UserRepository
	companyRepo repo.CompanyRepository
}

func NewDashboardService(
	dashRepo repo.DashboardRepository,
	userRepo repo.UserRepository,
	companyRepo repo.CompanyRepository,
) DashboardService {
	return &dashboardService{
		repo:        dashRepo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
	}
}

func shouldFilterByCompany(userType string, companyID uint) bool {
	return userType != "superadmin" && companyID > 0
}

// GetDashboardMetrics returns global metrics.
// Keep this only for public or superadmin usage.
func (s *dashboardService) GetDashboardMetrics() (*output.DashboardMetricsOutput, error) {
	return s.GetDashboardMetricsWithUserContext(
		0,
		"superadmin",
		"",
		0,
	)
}

// GetDashboardMetricsWithUserContext returns global data for superadmin
// and company-scoped data for every other user.
func (s *dashboardService) GetDashboardMetricsWithUserContext(
	userID uint,
	userType string,
	email string,
	companyID uint,
) (*output.DashboardMetricsOutput, error) {
	shouldFilter := shouldFilterByCompany(userType, companyID)

	totalCustomers, _ := s.repo.GetTotalCustomersWithFilter(shouldFilter, companyID)
	activeCustomers, _ := s.repo.GetActiveCustomersWithFilter(shouldFilter, companyID)
	customersCreatedToday, _ := s.repo.GetCustomersCreatedTodayWithFilter(shouldFilter, companyID)

	totalVendors, _ := s.repo.GetTotalVendorsWithFilter(shouldFilter, companyID)
	activeVendors, _ := s.repo.GetActiveVendorsWithFilter(shouldFilter, companyID)
	vendorsCreatedToday, _ := s.repo.GetVendorsCreatedTodayWithFilter(shouldFilter, companyID)

	totalItems, _ := s.repo.GetTotalItemsWithFilter(shouldFilter, companyID)
	totalItemGroups, _ := s.repo.GetTotalItemGroupsWithFilter(shouldFilter, companyID)
	totalStock, _ := s.repo.GetTotalStockWithFilter(shouldFilter, companyID)
	lowStockItems, _ := s.repo.GetLowStockItemsWithFilter(100, shouldFilter, companyID)
	outOfStockItems, _ := s.repo.GetOutOfStockItemsWithFilter(shouldFilter, companyID)
	itemsCreatedToday, _ := s.repo.GetItemsCreatedTodayWithFilter(shouldFilter, companyID)

	totalShipments, _ := s.repo.GetTotalShipmentsWithFilter(shouldFilter, companyID)
	shippedCount, _ := s.repo.GetShippedCountWithFilter(shouldFilter, companyID)
	pendingShipments, _ := s.repo.GetShipmentsByStatusWithFilter("pending", shouldFilter, companyID)
	deliveredShipments, _ := s.repo.GetShipmentsByStatusWithFilter("delivered", shouldFilter, companyID)
	inTransitShipments, _ := s.repo.GetShipmentsByStatusWithFilter("in_transit", shouldFilter, companyID)

	totalInvoices, _ := s.repo.GetTotalInvoicesWithFilter(shouldFilter, companyID)
	totalInvoiceAmount, _ := s.repo.GetTotalInvoiceAmountWithFilter(shouldFilter, companyID)
	outstandingInvoices, _ := s.repo.GetOutstandingInvoicesWithFilter(shouldFilter, companyID)
	paidInvoices, _ := s.repo.GetInvoicesByStatusWithFilter("paid", shouldFilter, companyID)
	pendingInvoices, _ := s.repo.GetInvoicesByStatusWithFilter("pending", shouldFilter, companyID)
	overdueInvoices, _ := s.repo.GetOverdueInvoicesWithFilter(shouldFilter, companyID)

	totalSalesOrders, _ := s.repo.GetTotalSalesOrdersWithFilter(shouldFilter, companyID)
	totalSalesOrderAmount, _ := s.repo.GetTotalSalesOrderAmountWithFilter(shouldFilter, companyID)
	completedSalesOrders, _ := s.repo.GetSalesOrdersByStatusWithFilter("completed", shouldFilter, companyID)
	pendingSalesOrders, _ := s.repo.GetSalesOrdersByStatusWithFilter("pending", shouldFilter, companyID)
	cancelledSalesOrders, _ := s.repo.GetSalesOrdersByStatusWithFilter("cancelled", shouldFilter, companyID)
	salesOrdersCreatedToday, _ := s.repo.GetSalesOrdersCreatedTodayWithFilter(shouldFilter, companyID)

	totalPurchaseOrders, _ := s.repo.GetTotalPurchaseOrdersWithFilter(shouldFilter, companyID)
	totalPurchaseOrderAmount, _ := s.repo.GetTotalPurchaseOrderAmountWithFilter(shouldFilter, companyID)
	completedPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatusWithFilter("completed", shouldFilter, companyID)
	pendingPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatusWithFilter("pending", shouldFilter, companyID)
	cancelledPurchaseOrders, _ := s.repo.GetPurchaseOrdersByStatusWithFilter("cancelled", shouldFilter, companyID)
	purchaseOrdersCreatedToday, _ := s.repo.GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter, companyID)

	totalPackages, _ := s.repo.GetTotalPackagesWithFilter(shouldFilter, companyID)
	shippedPackages, _ := s.repo.GetPackagesByStatusWithFilter("shipped", shouldFilter, companyID)
	pendingPackages, _ := s.repo.GetPackagesByStatusWithFilter("pending", shouldFilter, companyID)
	inTransitPackages, _ := s.repo.GetPackagesByStatusWithFilter("in_transit", shouldFilter, companyID)
	deliveredPackages, _ := s.repo.GetPackagesByStatusWithFilter("delivered", shouldFilter, companyID)
	packagesCreatedToday, _ := s.repo.GetPackagesCreatedTodayWithFilter(shouldFilter, companyID)

	userInfo := output.UserInfoOutput{
		UserID:      userID,
		UserRole:    userType,
		CompanyID:   companyID,
		CompanyName: "",
		UserName:    "",
		Email:       email,
	}

	if userID > 0 && s.userRepo != nil {
		user, err := s.userRepo.GetByID(userID)
		if err == nil && user != nil {
			if user.Username != nil {
				userInfo.UserName = *user.Username
			}

			if user.Email != nil {
				userInfo.Email = *user.Email
			}
		}
	}

	if companyID > 0 && s.companyRepo != nil {
		company, err := s.companyRepo.FindByID(companyID)
		if err == nil && company != nil {
			userInfo.CompanyName = company.CompanyName
		}
	}

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

// RefreshDashboardMetrics remains global.
// Apply SuperAdminMiddleware to its route.
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

func (s *dashboardService) AddShipmentTracking(
	shipmentID string,
	status string,
	location string,
	latitude float64,
	longitude float64,
	notes string,
	companyID uint,
	userType string,
) error {
	shouldFilter := shouldFilterByCompany(userType, companyID)

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

	return s.repo.AddShipmentTrackingWithFilter(
		tracking,
		shouldFilter,
		companyID,
	)
}

func (s *dashboardService) GetShipmentTracking(
	shipmentID string,
	limit int,
	companyID uint,
	userType string,
) (*output.ShipmentTrackingListOutput, error) {
	if limit < 1 {
		limit = 10
	}

	shouldFilter := shouldFilterByCompany(userType, companyID)

	tracking, err := s.repo.GetShipmentTrackingWithFilter(
		shipmentID,
		limit,
		shouldFilter,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	result := &output.ShipmentTrackingListOutput{
		Data:  make([]output.ShipmentTrackingOutput, len(tracking)),
		Total: len(tracking),
	}

	for index, item := range tracking {
		result.Data[index] = output.ShipmentTrackingOutput{
			ID:         item.ID,
			ShipmentID: item.ShipmentID,
			Status:     item.Status,
			Location:   item.Location,
			Latitude:   item.Latitude,
			Longitude:  item.Longitude,
			Notes:      item.Notes,
			Timestamp:  item.Timestamp,
		}
	}

	return result, nil
}

func (s *dashboardService) GetStockSummary() (*output.StockListOutput, error) {
	return s.GetStockSummaryWithUserContext(
		0,
		"superadmin",
		0,
	)
}

func (s *dashboardService) GetStockSummaryWithUserContext(
	userID uint,
	userType string,
	companyID uint,
) (*output.StockListOutput, error) {
	_ = userID

	shouldFilter := shouldFilterByCompany(userType, companyID)

	inStock, _ := s.repo.GetInStockProductsWithFilter(shouldFilter, companyID)
	lowStock, _ := s.repo.GetLowStockProductsWithFilter(100, shouldFilter, companyID)
	outOfStock, _ := s.repo.GetOutOfStockProductsWithFilter(shouldFilter, companyID)
	totalStock, _ := s.repo.GetTotalProductStockWithFilter(shouldFilter, companyID)

	productStockDetails, err := s.repo.GetProductStockDetailsWithFilter(
		shouldFilter,
		companyID,
	)
	if err != nil {
		return &output.StockListOutput{
			Data:            []output.ProductStockDetailOutput{},
			TotalProducts:   int(inStock),
			InStockCount:    int(inStock - lowStock - outOfStock),
			LowStockCount:   int(lowStock),
			OutOfStockCount: int(outOfStock),
			TotalQuantity:   totalStock,
		}, nil
	}

	stockDetails := make([]output.ProductStockDetailOutput, 0)

	for _, product := range productStockDetails {
		lastPurchasedDate, _ := product["last_purchased_date"].(*time.Time)
		lastSoldDate, _ := product["last_sold_date"].(*time.Time)

		rawMaterialUnit := ""
		if value, ok := product["raw_material_unit"]; ok && value != nil {
			rawMaterialUnit, _ = value.(string)
		}

		productName, _ := product["product_name"].(string)

		if isRaw, ok := product["is_raw"].(bool); ok && isRaw {
			rawName, _ := product["raw_name"].(string)
			rawSpecification, _ := product["raw_specification"].(string)

			if rawName != "" && rawSpecification != "" {
				productName = rawName + "_" + rawSpecification
			} else if rawName != "" {
				productName = rawName
			}
		}

		stockDetails = append(stockDetails, output.ProductStockDetailOutput{
			ProductID:         stringValue(product, "product_id"),
			ProductName:       productName,
			SKU:               stringValue(product, "sku"),
			CurrentStock:      floatValue(product, "current_stock"),
			AvailableStock:    floatValue(product, "available_stock"),
			ReservedStock:     floatValue(product, "reserved_stock"),
			PurchasedStock:    floatValue(product, "purchased_stock"),
			SoldStock:         floatValue(product, "sold_stock"),
			AverageCost:       floatValue(product, "average_cost"),
			RevaluationAmount: floatValue(product, "revaluation_amount"),
			LastPurchasedDate: lastPurchasedDate,
			LastSoldDate:      lastSoldDate,
			Status:            stringValue(product, "status"),
			RawMaterialUnit:   rawMaterialUnit,
		})
	}

	return &output.StockListOutput{
		Data:            stockDetails,
		TotalProducts:   int(inStock),
		InStockCount:    int(inStock - lowStock - outOfStock),
		LowStockCount:   int(lowStock),
		OutOfStockCount: int(outOfStock),
		TotalQuantity:   totalStock,
	}, nil
}

func (s *dashboardService) GetEntityTrends(
	entityType string,
	days int,
	companyID uint,
	userType string,
) (*output.EntityTrendOutput, error) {
	shouldFilter := shouldFilterByCompany(userType, companyID)

	history, err := s.repo.GetEntityCountHistoryWithFilter(
		entityType,
		days,
		shouldFilter,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	result := &output.EntityTrendOutput{
		EntityType: entityType,
		Data:       make([]output.TrendPoint, len(history)),
	}

	for index, item := range history {
		result.Data[index] = output.TrendPoint{
			Date:        item.Date,
			Count:       item.Count,
			ActiveCount: item.ActiveCount,
			NewToday:    item.CreatedToday,
		}
	}

	return result, nil
}

func (s *dashboardService) GetActivitySummary() (*output.ActivitySummaryOutput, error) {
	return s.GetActivitySummaryWithUserContext(
		0,
		"superadmin",
		0,
	)
}

func (s *dashboardService) GetActivitySummaryWithUserContext(
	userID uint,
	userType string,
	companyID uint,
) (*output.ActivitySummaryOutput, error) {
	_ = userID

	shouldFilter := shouldFilterByCompany(userType, companyID)

	customersCreatedToday, _ := s.repo.GetCustomersCreatedTodayWithFilter(shouldFilter, companyID)
	vendorsCreatedToday, _ := s.repo.GetVendorsCreatedTodayWithFilter(shouldFilter, companyID)
	itemsCreatedToday, _ := s.repo.GetItemsCreatedTodayWithFilter(shouldFilter, companyID)
	salesOrdersCreatedToday, _ := s.repo.GetSalesOrdersCreatedTodayWithFilter(shouldFilter, companyID)
	purchaseOrdersCreatedToday, _ := s.repo.GetPurchaseOrdersCreatedTodayWithFilter(shouldFilter, companyID)
	shippedToday, _ := s.repo.GetShippedTodayWithFilter(shouldFilter, companyID)
	deliveredToday, _ := s.repo.GetShipmentsByStatusWithFilter("delivered", shouldFilter, companyID)

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

func (s *dashboardService) GetDiagnosticReport(
	companyID uint,
	userType string,
) (*output.DiagnosticReportOutput, error) {
	shouldFilter := shouldFilterByCompany(userType, companyID)

	issues := []string{}
	diagnostics := make(map[string]output.DiagnosticItem)

	totalCustomers, _ := s.repo.GetTotalCustomersWithFilter(shouldFilter, companyID)
	activeCustomers, _ := s.repo.GetActiveCustomersWithFilter(shouldFilter, companyID)

	if totalCustomers > 0 && activeCustomers == 0 {
		issues = append(issues, "No active customers found - check is_active field")
	}

	diagnostics["customers"] = output.DiagnosticItem{
		Label: "Customers",
		Value: map[string]interface{}{
			"total":  totalCustomers,
			"active": activeCustomers,
		},
		Status:      statusCheck(activeCustomers, totalCustomers),
		Description: "Total vs Active customers",
	}

	totalVendors, _ := s.repo.GetTotalVendorsWithFilter(shouldFilter, companyID)
	activeVendors, _ := s.repo.GetActiveVendorsWithFilter(shouldFilter, companyID)

	if totalVendors > 0 && activeVendors == 0 {
		issues = append(issues, "No active vendors found - check is_active field")
	}

	diagnostics["vendors"] = output.DiagnosticItem{
		Label: "Vendors",
		Value: map[string]interface{}{
			"total":  totalVendors,
			"active": activeVendors,
		},
		Status:      statusCheck(activeVendors, totalVendors),
		Description: "Total vs Active vendors",
	}

	totalItems, _ := s.repo.GetTotalItemsWithFilter(shouldFilter, companyID)
	totalStock, _ := s.repo.GetTotalStockWithFilter(shouldFilter, companyID)

	if totalItems > 0 && totalStock == 0 {
		issues = append(issues, "No inventory balance records - inventory_balance table may be empty")
	}

	diagnostics["inventory"] = output.DiagnosticItem{
		Label: "Inventory",
		Value: map[string]interface{}{
			"total_items": totalItems,
			"total_stock": totalStock,
		},
		Status:      statusCheck(totalStock, totalItems),
		Description: "Items with inventory tracking",
	}

	totalInvoices, _ := s.repo.GetTotalInvoicesWithFilter(shouldFilter, companyID)
	totalInvoiceAmount, _ := s.repo.GetTotalInvoiceAmountWithFilter(shouldFilter, companyID)

	if totalInvoices > 0 && totalInvoiceAmount == 0 {
		issues = append(issues, "Invoice amounts are zero - check invoice.total_amount field")
	}

	diagnostics["invoices"] = output.DiagnosticItem{
		Label: "Invoices",
		Value: map[string]interface{}{
			"total":        totalInvoices,
			"total_amount": totalInvoiceAmount,
		},
		Status:      statusCheck(totalInvoiceAmount, totalInvoices),
		Description: "Total invoices with amounts",
	}

	totalSalesOrders, _ := s.repo.GetTotalSalesOrdersWithFilter(shouldFilter, companyID)
	totalSalesOrderAmount, _ := s.repo.GetTotalSalesOrderAmountWithFilter(shouldFilter, companyID)

	if totalSalesOrders > 0 && totalSalesOrderAmount == 0 {
		issues = append(issues, "Sales order amounts are zero - check sales_order.total_amount field")
	}

	diagnostics["sales_orders"] = output.DiagnosticItem{
		Label: "Sales Orders",
		Value: map[string]interface{}{
			"total":        totalSalesOrders,
			"total_amount": totalSalesOrderAmount,
		},
		Status:      statusCheck(totalSalesOrderAmount, totalSalesOrders),
		Description: "Total sales orders with amounts",
	}

	totalPurchaseOrders, _ := s.repo.GetTotalPurchaseOrdersWithFilter(shouldFilter, companyID)
	totalPurchaseOrderAmount, _ := s.repo.GetTotalPurchaseOrderAmountWithFilter(shouldFilter, companyID)

	if totalPurchaseOrders > 0 && totalPurchaseOrderAmount == 0 {
		issues = append(issues, "Purchase order amounts are zero - check purchase_order.total_amount field")
	}

	diagnostics["purchase_orders"] = output.DiagnosticItem{
		Label: "Purchase Orders",
		Value: map[string]interface{}{
			"total":        totalPurchaseOrders,
			"total_amount": totalPurchaseOrderAmount,
		},
		Status:      statusCheck(totalPurchaseOrderAmount, totalPurchaseOrders),
		Description: "Total purchase orders with amounts",
	}

	totalShipments, _ := s.repo.GetTotalShipmentsWithFilter(shouldFilter, companyID)

	if totalShipments == 0 {
		issues = append(issues, "No shipments found - shipment tracking not available")
	}

	diagnostics["shipments"] = output.DiagnosticItem{
		Label: "Shipments",
		Value: map[string]interface{}{
			"total": totalShipments,
		},
		Status:      statusCheck(totalShipments, 1),
		Description: "Total shipments in database",
	}

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

func stringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok && value != nil {
		if result, ok := value.(string); ok {
			return result
		}
	}

	return ""
}

func floatValue(data map[string]interface{}, key string) float64 {
	if value, ok := data[key]; ok && value != nil {
		switch typedValue := value.(type) {
		case float64:
			return typedValue
		case float32:
			return float64(typedValue)
		case int:
			return float64(typedValue)
		case int64:
			return float64(typedValue)
		}
	}

	return 0
}

func statusCheck(actual interface{}, expected interface{}) string {
	actualValue := getValue(actual)
	expectedValue := getValue(expected)

	if actualValue == expectedValue && actualValue > 0 {
		return "ok"
	}

	if actualValue > 0 {
		return "warning"
	}

	if expectedValue > 0 {
		return "error"
	}

	return "ok"
}

func getValue(value interface{}) int {
	switch typedValue := value.(type) {
	case int:
		return typedValue
	case int64:
		return int(typedValue)
	case float64:
		return int(typedValue)
	default:
		return 0
	}
}
