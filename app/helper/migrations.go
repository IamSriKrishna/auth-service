package helper

import (
	"log"
	"os"

	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/utils"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// Temporarily disable foreign key checks during migration
	log.Println("Disabling foreign key checks...")
	if err := db.Exec("SET FOREIGN_KEY_CHECKS=0").Error; err != nil {
		log.Printf("Warning: Failed to disable foreign key checks: %v", err)
	}
	defer func() {
		log.Println("Re-enabling foreign key checks...")
		if err := db.Exec("SET FOREIGN_KEY_CHECKS=1").Error; err != nil {
			log.Printf("Warning: Failed to re-enable foreign key checks: %v", err)
		}
	}()

	// Drop problematic foreign key constraints if they exist
	if db.Migrator().HasConstraint("inventory_balances", "fk_inventory_balances_variant") {
		log.Println("Removing problematic foreign key constraint fk_inventory_balances_variant...")
		if err := db.Migrator().DropConstraint("inventory_balances", "fk_inventory_balances_variant"); err != nil {
			log.Printf("Warning: Failed to drop constraint: %v", err)
		}
	}

	if db.Migrator().HasConstraint("inventory_aggregations", "fk_inventory_aggregations_variant") {
		log.Println("Removing problematic foreign key constraint fk_inventory_aggregations_variant...")
		if err := db.Migrator().DropConstraint("inventory_aggregations", "fk_inventory_aggregations_variant"); err != nil {
			log.Printf("Warning: Failed to drop constraint: %v", err)
		}
	}

	if db.Migrator().HasConstraint("sales_order_line_items", "fk_sales_order_line_items_variant") {
		log.Println("Removing problematic foreign key constraint fk_sales_order_line_items_variant...")
		if err := db.Migrator().DropConstraint("sales_order_line_items", "fk_sales_order_line_items_variant"); err != nil {
			log.Printf("Warning: Failed to drop constraint: %v", err)
		}
	}

	if os.Getenv("DROP_SALES_ORDER_TABLES") == "true" {
		log.Println("DROP_SALES_ORDER_TABLES=true detected, dropping sales order tables...")
		if err := DropSalesOrderTables(db); err != nil {
			log.Printf("Warning: Failed to drop sales order tables: %v", err)
		}
	}

	if os.Getenv("DROP_PURCHASE_ORDER_TABLES") == "true" {
		log.Println("DROP_PURCHASE_ORDER_TABLES=true detected, dropping purchase order tables...")
		if err := DropPurchaseOrderTables(db); err != nil {
			log.Printf("Warning: Failed to drop purchase order tables: %v", err)
		}
	}

	if os.Getenv("DROP_ORDER_FULFILLMENT_TABLES") == "true" {
		log.Println("DROP_ORDER_FULFILLMENT_TABLES=true detected, dropping order and fulfillment tables...")
		if err := DropOrderFulfillmentTables(db); err != nil {
			log.Printf("Warning: Failed to drop order fulfillment tables: %v", err)
		}
	}

	if os.Getenv("DROP_ALL_EXCEPT_USER") == "true" {
		log.Println("DROP_ALL_EXCEPT_USER=true detected, dropping all tables except user-related...")
		if err := DropAllTablesExceptUser(db); err != nil {
			log.Printf("Warning: Failed to drop tables: %v", err)
		}
	}

	if os.Getenv("DROP_ITEM_TABLES") == "true" {
		log.Println("DROP_ITEM_TABLES=true detected, dropping item tables...")
		if err := DropItemTables(db); err != nil {
			log.Printf("Warning: Failed to drop existing tables: %v", err)
		}
	}

	if os.Getenv("DROP_PRODUCT_DATA") == "true" {
		log.Println("DROP_PRODUCT_DATA=true detected, dropping all product data...")
		if err := DropProductData(db); err != nil {
			log.Printf("Warning: Failed to drop product data: %v", err)
		}
	}

	if os.Getenv("DROP_ITEM_MODEL_ONLY") == "true" {
		log.Println("DROP_ITEM_MODEL_ONLY=true detected, dropping only Item model...")
		if err := DropItemModelOnly(db); err != nil {
			log.Printf("Warning: Failed to drop Item model: %v", err)
		}
	}

	if os.Getenv("DROP_ALL_TABLES") == "true" {
		log.Println("DROP_ALL_TABLES=true detected, dropping ALL tables...")
		if err := DropAllTables(db); err != nil {
			log.Printf("Warning: Failed to drop all tables: %v", err)
		}
	}

	// Clean up orphaned sales order line items before adding foreign key constraints
	log.Println("Cleaning up orphaned sales order line items...")
	if err := cleanupOrphanedSalesOrderLineItems(db); err != nil {
		log.Printf("Warning: Failed to cleanup orphaned sales order line items: %v", err)
	}

	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
		&models.UserSession{},
		&models.Support{},
		&models.Employee{},
		&models.EmployeeAttendance{},

		&models.BusinessType{},
		&models.Country{},
		&models.State{},
		&models.TaxType{},
		&models.Bank{},

		&models.Company{},
		&models.CompanyContact{},
		&models.CompanyAddress{},
		&models.CompanyBankDetail{},
		&models.CompanyUPIDetail{},
		&models.CompanyInvoiceSetting{},
		&models.CompanyTaxSetting{},
		&models.CompanyRegionalSetting{},

		&models.Vendor{},
		&models.Customer{},
		&models.EntityOtherDetails{},
		&models.EntityAddress{},
		&models.EntityContactPerson{},
		&models.VendorBankDetail{},
		&models.EntityDocument{},

		&models.Item{},
		&models.ItemDetails{},
		&models.Variant{},
		&models.VariantAttribute{},
		&models.SalesInfo{},
		&models.PurchaseInfo{},
		&models.Inventory{},
		&models.ReturnPolicy{},
		&models.OpeningStock{},
		&models.VariantOpeningStock{},
		&models.StockMovement{},
		&models.Manufacturer{},

		&models.Product{},
		&models.ProductGroup{},
		&models.ProductGroupComponent{},
		&models.ProductGroupResource{},
		&models.ProductGroupInventory{},
		&models.ComponentInventory{},
		&models.ProductDetails{},
		&models.ProductVariant{},
		&models.ProductVariantAttribute{},

		&models.ItemGroup{},
		&models.ItemGroupComponent{},
		&models.ProductionOrder{},
		&models.ProductionOrderItem{},

		&models.ProductStock{},
		&models.StockLedger{},

		&models.VariantStock{},
		&models.VariantStockMovement{},
		&models.StockReservation{},

		&models.InventoryBalance{},
		&models.InventoryAggregation{},
		&models.InventoryJournal{},
		&models.SupplyChainSummary{},

		&models.Invoice{},
		&models.InvoiceLineItem{},
		&models.Salesperson{},
		&models.Tax{},
		&models.Payment{},
		&models.PaymentSplit{},
		&models.EmailCommunication{},

		&models.PurchaseOrder{},
		&models.PurchaseOrderLineItem{},
		&models.VendorPayment{},

		&models.SalesOrder{},
		&models.SalesOrderLineItem{},

		&models.Package{},
		&models.PackageItem{},

		&models.Shipment{},

		&models.Bill{},
		&models.BillLineItem{},
	)

	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Database migration completed successfully!")

	if err := utils.SeedInitialData(db); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
	}

	if err := utils.SeedDefaultCompany(db); err != nil {
		log.Printf("Warning: Failed to seed default company: %v", err)
	}

	return nil
}

// cleanupOrphanedSalesOrderLineItems removes sales order line items with invalid product_group_id references
func cleanupOrphanedSalesOrderLineItems(db *gorm.DB) error {
	// Check if the table exists before trying to clean it
	if !db.Migrator().HasTable("sales_order_line_items") {
		log.Println("Table 'sales_order_line_items' does not exist yet, skipping cleanup")
		return nil
	}

	if !db.Migrator().HasTable("product_groups") {
		log.Println("Table 'product_groups' does not exist yet, skipping cleanup")
		return nil
	}

	// Delete sales order line items where product_group_id is empty or null
	result := db.Where("product_group_id IS NULL OR product_group_id = ''").Delete(&models.SalesOrderLineItem{})
	if result.Error != nil {
		log.Printf("Warning: Failed to delete NULL/empty product_group_id rows: %v", result.Error)
		return nil
	}
	if result.RowsAffected > 0 {
		log.Printf("Deleted %d sales order line items with NULL/empty product_group_id", result.RowsAffected)
	}

	// Delete sales order line items with orphaned product_group_id (not in product_groups table)
	result = db.Where("product_group_id NOT IN (SELECT id FROM product_groups)").Delete(&models.SalesOrderLineItem{})
	if result.Error != nil {
		log.Printf("Warning: Failed to delete orphaned product_group_id rows: %v", result.Error)
		return nil
	}
	if result.RowsAffected > 0 {
		log.Printf("Deleted %d sales order line items with orphaned product_group_id", result.RowsAffected)
	}

	log.Println("Orphaned sales order line items cleanup completed")
	return nil
}

func DropSalesOrderTables(db *gorm.DB) error {
	log.Println("Dropping sales order tables...")

	tables := []interface{}{
		&models.Shipment{},
		&models.Package{},
		&models.PackageItem{},
		&models.SalesOrderLineItem{},
		&models.SalesOrder{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("Sales order tables dropped successfully!")
	return nil
}

func DropPurchaseOrderTables(db *gorm.DB) error {
	log.Println("Dropping purchase order related tables...")

	tables := []interface{}{
		&models.VendorPayment{},
		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("Purchase order tables dropped successfully!")
	return nil
}

func DropItemTables(db *gorm.DB) error {
	log.Println("Dropping item-related tables...")

	tables := []interface{}{
		&models.StockReservation{},
		&models.VariantStockMovement{},
		&models.VariantStock{},
		&models.InventoryJournal{},
		&models.InventoryAggregation{},
		&models.InventoryBalance{},
		&models.VariantOpeningStock{},
		&models.OpeningStock{},
		&models.StockMovement{},
		&models.ProductionOrderItem{},
		&models.ProductionOrder{},
		&models.ItemGroupComponent{},
		&models.ItemGroup{},
		&models.ProductVariantAttribute{},
		&models.ProductVariant{},
		&models.ProductDetails{},
		&models.Product{},
		&models.ProductGroup{},
		&models.ProductGroupComponent{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.VariantAttribute{},
		&models.Variant{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.ItemDetails{},
		&models.Manufacturer{},
		&models.Item{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("Item tables dropped successfully!")
	return nil
}

func DropItemModelOnly(db *gorm.DB) error {
	log.Println("Dropping Item model only...")

	tables := []interface{}{
		&models.StockReservation{},
		&models.VariantStockMovement{},
		&models.VariantStock{},
		&models.ProductVariantAttribute{},
		&models.ProductVariant{},
		&models.ProductDetails{},
		&models.Product{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.VariantAttribute{},
		&models.Variant{},
		&models.ItemDetails{},
		&models.Item{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("Item model dropped successfully!")
	return nil
}

func DropAllTablesExceptUser(db *gorm.DB) error {
	log.Println("WARNING: Dropping ALL tables except User-related tables...")

	// Tables to drop (all except Role, User, RefreshToken, UserSession)
	tablesToDrop := []interface{}{
		&models.BillLineItem{},
		&models.Bill{},

		&models.Shipment{},

		&models.PackageItem{},
		&models.Package{},

		&models.InvoiceLineItem{},
		&models.Invoice{},
		&models.Payment{},
		&models.Salesperson{},

		&models.SalesOrderLineItem{},
		&models.SalesOrder{},

		&models.VendorPayment{},
		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},
		&models.ProductionOrderItem{},
		&models.ProductionOrder{},

		&models.InventoryAggregation{},
		&models.InventoryBalance{},
		&models.InventoryJournal{},
		&models.SupplyChainSummary{},

		&models.StockReservation{},
		&models.VariantStockMovement{},
		&models.VariantStock{},

		&models.VariantOpeningStock{},
		&models.OpeningStock{},
		&models.StockMovement{},
		&models.ItemGroupComponent{},
		&models.ItemGroup{},
		&models.ProductVariantAttribute{},
		&models.ProductVariant{},
		&models.ProductDetails{},
		&models.Product{},
		&models.VariantAttribute{},
		&models.Variant{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.ItemDetails{},
		&models.Manufacturer{},
		&models.Item{},

		&models.Tax{},

		&models.EntityDocument{},
		&models.VendorBankDetail{},
		&models.EntityContactPerson{},
		&models.EntityAddress{},
		&models.EntityOtherDetails{},
		&models.Customer{},
		&models.Vendor{},

		&models.CompanyRegionalSetting{},
		&models.CompanyTaxSetting{},
		&models.CompanyInvoiceSetting{},
		&models.CompanyUPIDetail{},
		&models.CompanyBankDetail{},
		&models.CompanyAddress{},
		&models.CompanyContact{},
		&models.Company{},

		&models.TaxType{},
		&models.State{},
		&models.Country{},
		&models.BusinessType{},
		&models.Bank{},

		&models.Support{},
	}

	for _, table := range tablesToDrop {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("All tables except User-related tables dropped successfully!")
	return nil
}

func DropAllTables(db *gorm.DB) error {
	log.Println("WARNING: Dropping ALL tables completely...")

	// Disable foreign key constraints temporarily
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		log.Printf("Warning: Could not disable foreign key checks: %v", err)
	}

	// First, try dropping all known model tables
	allTables := []interface{}{
		&models.BillLineItem{},
		&models.Bill{},

		&models.Shipment{},

		&models.PackageItem{},
		&models.Package{},

		&models.InvoiceLineItem{},
		&models.Invoice{},
		&models.Payment{},
		&models.Salesperson{},

		&models.SalesOrderLineItem{},
		&models.SalesOrder{},

		&models.VendorPayment{},
		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},
		&models.ProductionOrderItem{},
		&models.ProductionOrder{},

		&models.InventoryAggregation{},
		&models.InventoryBalance{},
		&models.InventoryJournal{},
		&models.SupplyChainSummary{},

		&models.StockReservation{},
		&models.VariantStockMovement{},
		&models.VariantStock{},
		&models.StockLedger{},
		&models.ProductStock{},

		&models.VariantOpeningStock{},
		&models.OpeningStock{},
		&models.StockMovement{},
		&models.ItemGroupComponent{},
		&models.ItemGroup{},
		&models.ProductGroupComponent{},
		&models.ComponentInventory{},
		&models.ProductGroupInventory{},
		&models.ProductGroup{},
		&models.ProductVariantAttribute{},
		&models.ProductVariant{},
		&models.ProductDetails{},
		&models.Product{},
		&models.VariantAttribute{},
		&models.Variant{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.ItemDetails{},
		&models.Manufacturer{},
		&models.Item{},

		&models.Tax{},

		&models.EntityDocument{},
		&models.VendorBankDetail{},
		&models.EntityContactPerson{},
		&models.EntityAddress{},
		&models.EntityOtherDetails{},
		&models.Customer{},
		&models.Vendor{},

		&models.CompanyRegionalSetting{},
		&models.CompanyTaxSetting{},
		&models.CompanyInvoiceSetting{},
		&models.CompanyUPIDetail{},
		&models.CompanyBankDetail{},
		&models.CompanyAddress{},
		&models.CompanyContact{},
		&models.Company{},

		&models.TaxType{},
		&models.State{},
		&models.Country{},
		&models.BusinessType{},
		&models.Bank{},

		&models.Support{},
		&models.UserSession{},
		&models.RefreshToken{},
		&models.User{},
		&models.Role{},
	}

	for _, table := range allTables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Info: Attempted to drop table %T: %v", table, err)
		}
	}

	// Get all remaining tables from database and drop them
	var tables []string
	if err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tables).Error; err != nil {
		log.Printf("Warning: Could not query remaining tables: %v", err)
	} else {
		log.Printf("Remaining tables after initial drop: %v\n", tables)
		for _, table := range tables {
			if table != "" {
				if err := db.Exec("DROP TABLE IF EXISTS `" + table + "`").Error; err != nil {
					log.Printf("Warning: Failed to drop remaining table %s: %v", table, err)
				} else {
					log.Printf("Dropped remaining table: %s", table)
				}
			}
		}
	}

	// Re-enable foreign key constraints
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		log.Printf("Warning: Could not re-enable foreign key checks: %v", err)
	}

	log.Println("All tables dropped completely!")
	return nil
}

func DropOrderFulfillmentTables(db *gorm.DB) error {
	log.Println("Dropping order and fulfillment related tables...")

	// Order matters due to foreign key constraints
	// Drop in reverse order of creation dependencies
	tables := []interface{}{
		// Fulfillment
		&models.Shipment{},
		&models.Package{},
		&models.PackageItem{},

		// Invoicing and Billing
		&models.Invoice{},
		&models.InvoiceLineItem{},
		&models.Bill{},
		&models.BillLineItem{},

		// Sales Orders
		&models.SalesOrderLineItem{},
		&models.SalesOrder{},

		// Purchase Orders
		&models.VendorPayment{},
		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},

		// Stock and Inventory Management
		&models.StockReservation{},
		&models.StockLedger{},
		&models.VariantStockMovement{},
		&models.VariantStock{},
		&models.ProductStock{},
		&models.InventoryJournal{},
		&models.InventoryAggregation{},
		&models.InventoryBalance{},
		&models.VariantOpeningStock{},
		&models.OpeningStock{},
		&models.StockMovement{},

		// Production
		&models.ProductionOrderItem{},
		&models.ProductionOrder{},

		// Products and Variants
		&models.ItemGroupComponent{},
		&models.ProductGroupComponent{},
		&models.ItemGroup{},
		&models.ProductVariantAttribute{},
		&models.ProductVariant{},
		&models.ProductDetails{},
		&models.ProductGroup{},
		&models.ProductGroupInventory{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.VariantAttribute{},
		&models.Variant{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.ItemDetails{},
		&models.Manufacturer{},
		&models.Item{},
		&models.Product{},
	}

	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			if err := db.Migrator().DropTable(table); err != nil {
				log.Printf("Warning: Failed to drop table %T: %v", table, err)
			} else {
				log.Printf("Dropped table: %T", table)
			}
		}
	}

	log.Println("Order and fulfillment tables dropped successfully!")
	return nil
}

func ResetDatabase(db *gorm.DB) error {
	log.Println("WARNING: Resetting database...")

	if err := DropAllTables(db); err != nil {
		return err
	}

	return RunMigrations(db)
}

// DropProductData drops all product-related data: products, product groups, purchase orders, sales orders, packages, shipments, and stock
func DropProductData(db *gorm.DB) error {
	log.Println("WARNING: Dropping all product data (products, product groups, POs, SOs, packages, shipments, stock)...")

	// Disable foreign key checks temporarily
	if err := db.Exec("SET FOREIGN_KEY_CHECKS=0").Error; err != nil {
		log.Printf("Warning: Failed to disable foreign key checks: %v", err)
	}
	defer func() {
		if err := db.Exec("SET FOREIGN_KEY_CHECKS=1").Error; err != nil {
			log.Printf("Warning: Failed to re-enable foreign key checks: %v", err)
		}
	}()

	// Clear tables in order of dependencies (reverse order of creation)
	tablesToClear := []struct {
		name  string
		model interface{}
	}{
		{"shipments", &models.Shipment{}},
		{"packages", &models.Package{}},
		{"sales_orders", &models.SalesOrder{}},
		{"sales_order_line_items", &models.SalesOrderLineItem{}},
		{"purchase_orders", &models.PurchaseOrder{}},
		{"purchase_order_line_items", &models.PurchaseOrderLineItem{}},
		{"product_stock", &models.ProductStock{}},
		{"variant_stock", &models.VariantStock{}},
		{"variant_stock_movements", &models.VariantStockMovement{}},
		{"stock_ledgers", &models.StockLedger{}},
		{"stock_movements", &models.StockMovement{}},
		{"stock_reservations", &models.StockReservation{}},
		{"product_group_transactions", &models.ProductGroupTransaction{}},
		{"component_inventory", &models.ComponentInventory{}},
		{"product_group_inventory", &models.ProductGroupInventory{}},
		{"product_groups", &models.ProductGroup{}},
		{"product_groups_components", &models.ProductGroupComponent{}},
		{"product_details", &models.ProductDetails{}},
		{"products", &models.Product{}},
	}

	for _, table := range tablesToClear {
		if db.Migrator().HasTable(table.name) {
			result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table.model)
			if result.Error != nil {
				log.Printf("Warning: Failed to clear table %s: %v", table.name, result.Error)
			} else {
				log.Printf("✓ Cleared table %s (%d rows deleted)", table.name, result.RowsAffected)
			}
		} else {
			log.Printf("⊘ Table %s does not exist, skipping", table.name)
		}
	}

	log.Println("✓ Product data dropped successfully!")
	return nil
}
