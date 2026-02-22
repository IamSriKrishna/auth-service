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

	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
		&models.UserSession{},
		&models.Support{},

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

		&models.ItemGroup{},
		&models.ItemGroupComponent{},
		&models.ProductionOrder{},
		&models.ProductionOrderItem{},

		&models.InventoryBalance{},
		&models.InventoryAggregation{},
		&models.InventoryJournal{},
		&models.SupplyChainSummary{},

		&models.Invoice{},
		&models.InvoiceLineItem{},
		&models.Salesperson{},
		&models.Tax{},
		&models.Payment{},

		&models.PurchaseOrder{},
		&models.PurchaseOrderLineItem{},

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

func DropItemTables(db *gorm.DB) error {
	log.Println("Dropping item-related tables...")

	tables := []interface{}{
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

		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},
		&models.ProductionOrderItem{},
		&models.ProductionOrder{},

		&models.InventoryAggregation{},
		&models.InventoryBalance{},
		&models.InventoryJournal{},
		&models.SupplyChainSummary{},

		&models.VariantOpeningStock{},
		&models.OpeningStock{},
		&models.StockMovement{},
		&models.ItemGroupComponent{},
		&models.ItemGroup{},
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
	log.Println("WARNING: Dropping ALL tables...")

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

		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},
		&models.ProductionOrderItem{},
		&models.ProductionOrder{},

		&models.InventoryAggregation{},
		&models.InventoryBalance{},
		&models.InventoryJournal{},
		&models.SupplyChainSummary{},

		&models.VariantOpeningStock{},
		&models.OpeningStock{},
		&models.StockMovement{},
		&models.ItemGroupComponent{},
		&models.ItemGroup{},
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

		&models.Support{},
		&models.UserSession{},
		&models.RefreshToken{},
		&models.User{},
		&models.Role{},
	}

	for _, table := range allTables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("All tables dropped!")
	return nil
}

func ResetDatabase(db *gorm.DB) error {
	log.Println("WARNING: Resetting database...")

	if err := DropAllTables(db); err != nil {
		return err
	}

	return RunMigrations(db)
}
