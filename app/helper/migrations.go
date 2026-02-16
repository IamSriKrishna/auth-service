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

	if os.Getenv("DROP_ITEM_TABLES") == "true" {
		log.Println("DROP_ITEM_TABLES=true detected, dropping item tables...")
		if err := DropItemTables(db); err != nil {
			log.Printf("Warning: Failed to drop existing tables: %v", err)
		}
	}

	// Run migrations in order of dependencies
	err := db.AutoMigrate(
		// Auth & User Management
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
		&models.UserSession{},
		&models.Support{},

		// Master Data
		&models.BusinessType{},
		&models.Country{},
		&models.State{},
		&models.TaxType{},

		// Company Management
		&models.Company{},
		&models.CompanyContact{},
		&models.CompanyAddress{},
		&models.CompanyBankDetail{},
		&models.CompanyUPIDetail{},
		&models.CompanyInvoiceSetting{},
		&models.CompanyTaxSetting{},
		&models.CompanyRegionalSetting{},

		// Vendor & Customer Management
		&models.Vendor{},
		&models.Customer{},
		&models.EntityOtherDetails{},
		&models.EntityAddress{},
		&models.EntityContactPerson{},
		&models.VendorBankDetail{},
		&models.EntityDocument{},

		// Item Management (in order of foreign key dependencies)
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

		//invoice management
		&models.Invoice{},
		&models.InvoiceLineItem{},
		&models.Salesperson{},
		&models.Tax{},
		&models.Payment{},

		// Purchase Order management
		&models.PurchaseOrder{},
		&models.PurchaseOrderLineItem{},

		// Sales Order management
		&models.SalesOrder{},
		&models.SalesOrderLineItem{},

		// Package management
		&models.Package{},
		&models.PackageItem{},

		// Shipment management
		&models.Shipment{},

		// Bill management
		&models.Bill{},
		&models.BillLineItem{},
	)

	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return err
	}

	log.Println("Database migration completed successfully!")

	// Seed initial data
	if err := utils.SeedInitialData(db); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
	}

	if err := utils.SeedDefaultCompany(db); err != nil {
		log.Printf("Warning: Failed to seed default company: %v", err)
	}

	return nil
}

// DropItemTables drops all item-related tables in the correct order (reverse of creation)
func DropItemTables(db *gorm.DB) error {
	log.Println("Dropping item-related tables...")

	// Drop in reverse order of dependencies
	tables := []interface{}{
		&models.VariantAttribute{},
		&models.Variant{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.ItemDetails{},
		&models.Item{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
			// Continue dropping other tables even if one fails
		}
	}

	log.Println("Item tables dropped successfully!")
	return nil
}

// DropAllTables drops all tables in the database (useful for development/testing)
func DropAllTables(db *gorm.DB) error {
	log.Println("WARNING: Dropping ALL tables...")

	allTables := []interface{}{
		// Sales Order Management (drop first due to foreign keys)
		&models.SalesOrderLineItem{},
		&models.SalesOrder{},

		// Purchase Order Management (drop first due to foreign keys)
		&models.PurchaseOrderLineItem{},
		&models.PurchaseOrder{},

		// Item Management (drop first due to foreign keys)
		&models.VariantAttribute{},
		&models.Variant{},
		&models.ReturnPolicy{},
		&models.Inventory{},
		&models.PurchaseInfo{},
		&models.SalesInfo{},
		&models.ItemDetails{},
		&models.Item{},

		// Vendor & Customer Management
		&models.EntityDocument{},
		&models.VendorBankDetail{},
		&models.EntityContactPerson{},
		&models.EntityAddress{},
		&models.EntityOtherDetails{},
		&models.Customer{},
		&models.Vendor{},

		// Company Management
		&models.CompanyRegionalSetting{},
		&models.CompanyTaxSetting{},
		&models.CompanyInvoiceSetting{},
		&models.CompanyUPIDetail{},
		&models.CompanyBankDetail{},
		&models.CompanyAddress{},
		&models.CompanyContact{},
		&models.Company{},

		// Master Data
		&models.TaxType{},
		&models.State{},
		&models.Country{},
		&models.BusinessType{},

		// Auth & User Management
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

// ResetDatabase drops all tables and re-runs migrations (useful for development)
func ResetDatabase(db *gorm.DB) error {
	log.Println("WARNING: Resetting database...")

	if err := DropAllTables(db); err != nil {
		return err
	}

	return RunMigrations(db)
}
