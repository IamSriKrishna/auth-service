package routes

import (
	"time"

	"github.com/bbapp-org/auth-service/app/config"
	"github.com/bbapp-org/auth-service/app/config/database"
	"github.com/bbapp-org/auth-service/app/handlers"
	"github.com/bbapp-org/auth-service/app/middleware"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/services"
	"github.com/bbapp-org/auth-service/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize repositories with dbresolver-enabled DB
	// dbresolver automatically handles read/write splitting
	db := database.GetDB()

	// Initialize HTTP client for calling customer service
	httpClient := utils.NewHTTPClient(cfg.Service.CustomerServiceURL, 10*time.Second)

	userRepo := repo.NewUserRepository(db, httpClient)
	roleRepo := repo.NewRoleRepository(db)
	refreshTokenRepo := repo.NewRefreshTokenRepository(db)
	sessionRepo := repo.NewUserSessionRepository(db)
	supportRepo := repo.NewSupportRepository(db)
	vendorRepo := repo.NewVendorRepository(db)
	companyRepo := repo.NewCompanyRepository(db)
	businessTypeRepo := repo.NewBusinessTypeRepository(db)
	locationRepo := repo.NewLocationRepository(db)
	taxTypeRepo := repo.NewTaxTypeRepository(db)
	itemRepo := repo.NewItemRepository(db)
	customerRepo := repo.NewCustomerRepository(db)
	openStockRepo := repo.NewOpeningStockRepository(db)
	manufacturerRepo := repo.NewManufacturerRepository(db)
	brandRepo := repo.NewBrandRepository(db)
	invoiceRepo := repo.NewInvoiceRepository(db)
	salespersonRepo := repo.NewSalespersonRepository(db)
	taxRepo := repo.NewTaxRepository(db)
	paymentRepo := repo.NewPaymentRepository(db)
	purchaseOrderRepo := repo.NewPurchaseOrderRepository(db)
	salesOrderRepo := repo.NewSalesOrderRepository(db)
	packageRepo := repo.NewPackageRepository(db)
	shipmentRepo := repo.NewShipmentRepository(db)
	billRepo := repo.NewBillRepository(db)
	// Initialize services
	authService := services.NewAuthService(userRepo, roleRepo, refreshTokenRepo, sessionRepo)
	adminService := services.NewAdminService(userRepo, roleRepo)
	supportService := services.NewSupportService(supportRepo)
	businessTypeService := services.NewBusinessTypeService(businessTypeRepo)
	locationService := services.NewLocationService(locationRepo)
	taxTypeService := services.NewTaxTypeService(taxTypeRepo)
	companyService := services.NewCompanyService(companyRepo, businessTypeRepo, locationRepo, taxTypeRepo, db)
	itemService := services.NewItemService(itemRepo, vendorRepo, manufacturerRepo)
	vendorService := services.NewVendorService(vendorRepo)
	customerService := services.NewCustomerService(customerRepo)
	openStockService := services.NewOpeningStockService(openStockRepo, itemRepo)
	manufacturerService := services.NewManufacturerService(manufacturerRepo)
	brandService := services.NewBrandService(brandRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo, itemRepo, customerRepo, salespersonRepo, taxRepo, paymentRepo, "./pdf_outputs")
	salespersonService := services.NewSalespersonService(salespersonRepo)
	taxService := services.NewTaxService(taxRepo)
	paymentService := services.NewPaymentService(paymentRepo, invoiceRepo)
	purchaseOrderService := services.NewPurchaseOrderService(purchaseOrderRepo, vendorRepo, customerRepo, itemRepo, taxRepo)
	salesOrderService := services.NewSalesOrderService(salesOrderRepo, customerRepo, itemRepo, taxRepo, salespersonRepo)
	packageService := services.NewPackageService(packageRepo, salesOrderRepo, customerRepo, itemRepo)
	shipmentService := services.NewShipmentService(shipmentRepo, packageRepo, salesOrderRepo, customerRepo)
	billService := services.NewBillService(billRepo, vendorRepo, itemRepo, taxRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	supportHandler := handlers.NewSupportHandler(supportService)
	forwardAuthHandler := handlers.NewForwardAuthHandler()
	vendorHandler := handlers.NewVendorHandler(vendorService)
	companyHandler := handlers.NewCompanyHandler(companyService, businessTypeService, locationService, taxTypeService)
	helperHandler := handlers.NewHelperHandler(businessTypeService, locationService, taxTypeService)
	itemHandler := handlers.NewItemHandler(itemService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	openStockHandler := handlers.NewOpeningStockHandler(openStockService)
	manufacturerHandler := handlers.NewManufacturerHandler(manufacturerService)
	brandHandler := handlers.NewBrandHandler(brandService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)
	salespersonHandler := handlers.NewSalespersonHandler(salespersonService)
	taxHandler := handlers.NewTaxHandler(taxService)
	purchaseOrderHandler := handlers.NewPurchaseOrderHandler(purchaseOrderService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	salesOrderHandler := handlers.NewSalesOrderHandler(salesOrderService)
	packageHandler := handlers.NewPackageHandler(packageService)
	shipmentHandler := handlers.NewShipmentHandler(shipmentService)
	billHandler := handlers.NewBillHandler(billService)

	// Swagger documentation
	app.Get("/docs/*", swagger.HandlerDefault)

	// Public authentication routes
	authGroup := app.Group("/auth")
	{
		// Registration routes
		authGroup.Post("/register/email", authHandler.RegisterEmail)
		authGroup.Post("/register/phone", authHandler.RegisterPhone)
		authGroup.Post("/register/google", authHandler.RegisterGoogle)

		// Login routes
		authGroup.Post("/login/email", authHandler.LoginEmail)
		authGroup.Post("/login/phone", authHandler.LoginPhone)
		authGroup.Post("/login/google", authHandler.LoginGoogle)
		authGroup.Post("/login/apple", authHandler.LoginApple)
		authGroup.Post("/login/password", authHandler.LoginPassword)

		// Internal/Service endpoints
		authGroup.Post("/validate-token", authHandler.ValidateToken)
	}

	// Protected authentication routes (require valid JWT)
	protectedAuthGroup := app.Group("/auth")
	protectedAuthGroup.Use(middleware.AuthMiddleware())
	{
		protectedAuthGroup.Post("/refresh-token", authHandler.RefreshToken)
		protectedAuthGroup.Get("/user-info", authHandler.GetUserInfo)
		protectedAuthGroup.Post("/change-password", authHandler.ChangePassword) // should not be used
		protectedAuthGroup.Post("/logout", authHandler.Logout)
	}

	// Manufacturer routes
	manufacturerGroup := app.Group("/manufacturers")
	{
		manufacturerGroup.Get("/", manufacturerHandler.GetAllManufacturers)
		manufacturerGroup.Get("/:id", manufacturerHandler.GetManufacturerByID)
		manufacturerGroup.Post("/", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), manufacturerHandler.CreateManufacturer)
		manufacturerGroup.Put("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), manufacturerHandler.UpdateManufacturer)
		manufacturerGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), manufacturerHandler.DeleteManufacturer)
	}

	// Brand routes
	brandGroup := app.Group("/brands")
	{
		brandGroup.Get("/", brandHandler.GetAllBrands)
		brandGroup.Get("/:id", brandHandler.GetBrandByID)
		brandGroup.Post("/", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), brandHandler.CreateBrand)
		brandGroup.Put("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), brandHandler.UpdateBrand)
		brandGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), brandHandler.DeleteBrand)
	}

	// Super Admin only routes
	superAdminGroup := app.Group("/auth/admin")
	superAdminGroup.Use(middleware.AuthMiddleware())
	superAdminGroup.Use(middleware.SuperAdminMiddleware())
	{
		superAdminGroup.Post("/create-user", adminHandler.CreateUser)
		superAdminGroup.Post("/reset-password", adminHandler.ResetAdminPassword)
		superAdminGroup.Get("/users", adminHandler.GetUsers)
		superAdminGroup.Get("/users/:id", adminHandler.GetUser)
		superAdminGroup.Put("/users/:id", adminHandler.UpdateUser)
		superAdminGroup.Delete("/users/:id", adminHandler.DeleteUser)
		superAdminGroup.Put("/users/:id/status", adminHandler.UpdateUserStatus)
		superAdminGroup.Put("/users/:id/role", adminHandler.UpdateUserRole)
		superAdminGroup.Get("/dashboard/stats", adminHandler.GetDashboardStats)
	}

	vendorGroup := app.Group("/vendors")
	{
		vendorGroup.Get("/", vendorHandler.GetAllVendors)
		vendorGroup.Get("/:id", vendorHandler.GetVendor)
		vendorGroup.Post("/", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), vendorHandler.CreateVendor)
		vendorGroup.Put("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), vendorHandler.UpdateVendor)
		vendorGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), vendorHandler.DeleteVendor)
	}

	customerGroup := app.Group("/customers")
	{
		customerGroup.Get("/", customerHandler.GetAllCustomers)
		customerGroup.Get("/:id", customerHandler.GetCustomerByID)
		customerGroup.Post("/", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), customerHandler.CreateCustomer)
		customerGroup.Put("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), customerHandler.UpdateCustomer)
		customerGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), customerHandler.DeleteCustomer)
	}

	partners := app.Group("/partners")
	{
		partners.Patch("/:partner_id/reset-password", adminHandler.ResetUserPassword)
	}

	// Admin routes (admin and superadmin access) - limited to partner management
	adminGroup := app.Group("/auth/manage")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.AdminMiddleware())
	{
		adminGroup.Post("/create-partner", adminHandler.CreateUser) // Reuse CreateUser but limit to partners
		adminGroup.Get("/partners", adminHandler.GetUsers)          // Reuse GetUsers but filter partners
	}

	// Forward auth routes for Traefik middleware
	forwardAuthGroup := app.Group("/forward-auth")
	{
		forwardAuthGroup.Get("/", forwardAuthHandler.ForwardAuth)
		forwardAuthGroup.Get("/product", forwardAuthHandler.ProductAuth)
		forwardAuthGroup.Get("/customer", forwardAuthHandler.CustomerAuth)
	}

	// Public support route
	app.Post("/public/support", supportHandler.CreateSupport)

	// Helper/Lookup data routes (public)
	helperRoutes := app.Group("/helpers")
	{
		helperRoutes.Get("/business-types", helperHandler.GetBusinessTypes)
		helperRoutes.Get("/countries", helperHandler.GetCountries)
		helperRoutes.Get("/countries/:country_id/states", helperHandler.GetStatesByCountry)
		helperRoutes.Get("/tax-types", helperHandler.GetTaxTypes)
	}

	// Company routes
	companyRoutes := app.Group("/companies")
	companyRoutes.Use(middleware.AuthMiddleware()) // All company routes require authentication
	{
		// Complete company setup (one-time wizard)
		companyRoutes.Post("/setup", companyHandler.CompleteCompanySetup)

		// Company CRUD
		companyRoutes.Get("/", companyHandler.GetAllCompanies)
		companyRoutes.Post("/", companyHandler.CreateCompany)
		companyRoutes.Get("/:id", companyHandler.GetCompany)
		companyRoutes.Put("/:id", companyHandler.UpdateCompany)
		companyRoutes.Delete("/:id", middleware.SuperAdminMiddleware(), companyHandler.DeleteCompany)

		// Contact management
		companyRoutes.Put("/:id/contact", companyHandler.UpsertContact)
		companyRoutes.Get("/:id/contact", companyHandler.GetContact)

		// Address management
		companyRoutes.Put("/:id/address", companyHandler.UpsertAddress)
		companyRoutes.Get("/:id/address", companyHandler.GetAddress)

		// Bank details management
		companyRoutes.Post("/:id/bank-details", companyHandler.CreateBankDetail)
		companyRoutes.Get("/:id/bank-details", companyHandler.GetBankDetails)
		companyRoutes.Put("/bank-details/:id", companyHandler.UpdateBankDetail)
		companyRoutes.Delete("/bank-details/:id", companyHandler.DeleteBankDetail)

		// UPI details management
		companyRoutes.Put("/:id/upi-details", companyHandler.UpsertUPIDetail)
		companyRoutes.Get("/:id/upi-details", companyHandler.GetUPIDetail)

		// Invoice settings management
		companyRoutes.Put("/:id/invoice-settings", companyHandler.UpsertInvoiceSettings)
		companyRoutes.Get("/:id/invoice-settings", companyHandler.GetInvoiceSettings)

		// Tax settings management
		companyRoutes.Put("/:id/tax-settings", companyHandler.UpsertTaxSettings)
		companyRoutes.Get("/:id/tax-settings", companyHandler.GetTaxSettings)

		// Regional settings management
		companyRoutes.Put("/:id/regional-settings", companyHandler.UpsertRegionalSettings)
		companyRoutes.Get("/:id/regional-settings", companyHandler.GetRegionalSettings)
	}

	itemRoutes := app.Group("/items")
	{
		itemRoutes.Get("/", itemHandler.GetAllItems)
		itemRoutes.Get("/:id", itemHandler.GetItem)

		itemRoutes.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemHandler.CreateItem)
		itemRoutes.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemHandler.UpdateItem)
		itemRoutes.Delete("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), itemHandler.DeleteItem)

		itemRoutes.Put("/:id/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.UpdateOpeningStock)
		itemRoutes.Get("/:id/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.GetOpeningStock)

		// Variant opening stock
		itemRoutes.Put("/:id/variants/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.UpdateVariantsOpeningStock)
		itemRoutes.Get("/:id/variants/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.GetVariantsOpeningStock)
		// Stock summary
		itemRoutes.Get("/:id/stock-summary", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.GetStockSummary)
	}

	// Invoice routes
	invoiceRoutes := app.Group("/invoices")
	invoiceRoutes.Use(middleware.AuthMiddleware()) // All invoice routes require authentication
	{
		// Invoice CRUD
		invoiceRoutes.Post("/", middleware.AdminMiddleware(), invoiceHandler.CreateInvoice)
		invoiceRoutes.Get("/", invoiceHandler.GetAllInvoices)
		invoiceRoutes.Get("/:id", invoiceHandler.GetInvoice)
		invoiceRoutes.Put("/:id", middleware.AdminMiddleware(), invoiceHandler.UpdateInvoice)
		invoiceRoutes.Delete("/:id", middleware.AdminMiddleware(), invoiceHandler.DeleteInvoice)

		// Invoice status update
		invoiceRoutes.Patch("/:id/status", middleware.AdminMiddleware(), invoiceHandler.UpdateInvoiceStatus)

		// Invoice-specific payments
		invoiceRoutes.Get("/:invoiceId/payments", paymentHandler.GetPaymentsByInvoice)
	}

	// Status-based invoice filtering
	app.Get("/invoices/status/:status", middleware.AuthMiddleware(), invoiceHandler.GetInvoicesByStatus)

	// Customer-specific invoices
	customerGroup.Get("/:customerId/invoices", middleware.AuthMiddleware(), invoiceHandler.GetInvoicesByCustomer)

	// Salesperson routes
	salespersonRoutes := app.Group("/salespersons")
	salespersonRoutes.Use(middleware.AuthMiddleware()) // All salesperson routes require authentication
	{
		salespersonRoutes.Post("/", middleware.AdminMiddleware(), salespersonHandler.CreateSalesperson)
		salespersonRoutes.Get("/", salespersonHandler.GetAllSalespersons)
		salespersonRoutes.Get("/:id", salespersonHandler.GetSalesperson)
		salespersonRoutes.Put("/:id", middleware.AdminMiddleware(), salespersonHandler.UpdateSalesperson)
		salespersonRoutes.Delete("/:id", middleware.SuperAdminMiddleware(), salespersonHandler.DeleteSalesperson)
	}

	// Tax routes
	taxRoutes := app.Group("/taxes")
	taxRoutes.Use(middleware.AuthMiddleware()) // All tax routes require authentication
	{
		taxRoutes.Post("/", middleware.AdminMiddleware(), taxHandler.CreateTax)
		taxRoutes.Get("/", taxHandler.GetAllTaxes)
		taxRoutes.Get("/:id", taxHandler.GetTax)
		taxRoutes.Put("/:id", middleware.AdminMiddleware(), taxHandler.UpdateTax)
		taxRoutes.Delete("/:id", middleware.SuperAdminMiddleware(), taxHandler.DeleteTax)
	}

	// Payment routes
	paymentRoutes := app.Group("/payments")
	paymentRoutes.Use(middleware.AuthMiddleware()) // All payment routes require authentication
	{
		paymentRoutes.Post("/", middleware.AdminMiddleware(), paymentHandler.CreatePayment)
		paymentRoutes.Get("/:id", paymentHandler.GetPayment)
		paymentRoutes.Delete("/:id", middleware.AdminMiddleware(), paymentHandler.DeletePayment)
	}
	purchaseOrderRoutes := app.Group("/purchase-orders")
	purchaseOrderRoutes.Use(middleware.AuthMiddleware()) // All purchase order routes require authentication
	{
		// Purchase Order CRUD
		purchaseOrderRoutes.Post("/", middleware.AdminMiddleware(), purchaseOrderHandler.CreatePurchaseOrder)
		purchaseOrderRoutes.Get("/", purchaseOrderHandler.GetAllPurchaseOrders)
		purchaseOrderRoutes.Get("/:id", purchaseOrderHandler.GetPurchaseOrder)
		purchaseOrderRoutes.Put("/:id", middleware.AdminMiddleware(), purchaseOrderHandler.UpdatePurchaseOrder)
		purchaseOrderRoutes.Delete("/:id", middleware.AdminMiddleware(), purchaseOrderHandler.DeletePurchaseOrder)

		// Purchase Order status update
		purchaseOrderRoutes.Patch("/:id/status", middleware.AdminMiddleware(), purchaseOrderHandler.UpdatePurchaseOrderStatus)

		// Purchase Order-specific routes
		purchaseOrderRoutes.Get("/vendor/:vendorId", purchaseOrderHandler.GetPurchaseOrdersByVendor)
		purchaseOrderRoutes.Get("/customer/:customerId", purchaseOrderHandler.GetPurchaseOrdersByCustomer)
		purchaseOrderRoutes.Get("/status/:status", purchaseOrderHandler.GetPurchaseOrdersByStatus)
	}

	// Sales Order routes
	salesOrderRoutes := app.Group("/sales-orders")
	salesOrderRoutes.Use(middleware.AuthMiddleware()) // All sales order routes require authentication
	{
		// Sales Order CRUD
		salesOrderRoutes.Post("/", middleware.AdminMiddleware(), salesOrderHandler.CreateSalesOrder)
		salesOrderRoutes.Get("/", salesOrderHandler.GetAllSalesOrders)
		salesOrderRoutes.Get("/:id", salesOrderHandler.GetSalesOrder)
		salesOrderRoutes.Put("/:id", middleware.AdminMiddleware(), salesOrderHandler.UpdateSalesOrder)
		salesOrderRoutes.Delete("/:id", middleware.AdminMiddleware(), salesOrderHandler.DeleteSalesOrder)

		// Sales Order status update
		salesOrderRoutes.Patch("/:id/status", middleware.AdminMiddleware(), salesOrderHandler.UpdateSalesOrderStatus)

		// Sales Order-specific routes
		salesOrderRoutes.Get("/customer/:customerId", salesOrderHandler.GetSalesOrdersByCustomer)
		salesOrderRoutes.Get("/status/:status", salesOrderHandler.GetSalesOrdersByStatus)
	}

	// Package routes
	packageRoutes := app.Group("/packages")
	packageRoutes.Use(middleware.AuthMiddleware()) // All package routes require authentication
	{
		// Package CRUD
		packageRoutes.Post("/", middleware.AdminMiddleware(), packageHandler.CreatePackage)
		packageRoutes.Get("/", packageHandler.GetAllPackages)
		packageRoutes.Get("/:id", packageHandler.GetPackage)
		packageRoutes.Put("/:id", middleware.AdminMiddleware(), packageHandler.UpdatePackage)
		packageRoutes.Delete("/:id", middleware.AdminMiddleware(), packageHandler.DeletePackage)

		// Package status update
		packageRoutes.Patch("/:id/status", middleware.AdminMiddleware(), packageHandler.UpdatePackageStatus)

		// Package-specific routes
		packageRoutes.Get("/customer/:customer_id", packageHandler.GetPackagesByCustomer)
		packageRoutes.Get("/sales-order/:sales_order_id", packageHandler.GetPackagesBySalesOrder)
		packageRoutes.Get("/status/:status", packageHandler.GetPackagesByStatus)
	}

	// Shipment routes
	shipmentRoutes := app.Group("/shipments")
	shipmentRoutes.Use(middleware.AuthMiddleware()) // All shipment routes require authentication
	{
		// Shipment CRUD
		shipmentRoutes.Post("/", middleware.AdminMiddleware(), shipmentHandler.CreateShipment)
		shipmentRoutes.Get("/", shipmentHandler.GetAllShipments)
		shipmentRoutes.Get("/:id", shipmentHandler.GetShipment)
		shipmentRoutes.Put("/:id", middleware.AdminMiddleware(), shipmentHandler.UpdateShipment)
		shipmentRoutes.Delete("/:id", middleware.AdminMiddleware(), shipmentHandler.DeleteShipment)

		// Shipment status update
		shipmentRoutes.Patch("/:id/status", middleware.AdminMiddleware(), shipmentHandler.UpdateShipmentStatus)

		// Shipment-specific routes
		shipmentRoutes.Get("/customer/:customer_id", shipmentHandler.GetShipmentsByCustomer)
		shipmentRoutes.Get("/package/:package_id", shipmentHandler.GetShipmentsByPackage)
		shipmentRoutes.Get("/sales-order/:sales_order_id", shipmentHandler.GetShipmentsBySalesOrder)
		shipmentRoutes.Get("/status/:status", shipmentHandler.GetShipmentsByStatus)
	}

	// Bill routes
	billRoutes := app.Group("/bills")
	billRoutes.Use(middleware.AuthMiddleware()) // All bill routes require authentication
	{
		// Bill CRUD
		billRoutes.Post("/", middleware.AdminMiddleware(), billHandler.CreateBill)
		billRoutes.Get("/", billHandler.GetAllBills)
		billRoutes.Get("/:id", billHandler.GetBill)
		billRoutes.Put("/:id", middleware.AdminMiddleware(), billHandler.UpdateBill)
		billRoutes.Delete("/:id", middleware.AdminMiddleware(), billHandler.DeleteBill)

		// Bill status update
		billRoutes.Patch("/:id/status", middleware.AdminMiddleware(), billHandler.UpdateBillStatus)

		// Bill-specific routes
		billRoutes.Get("/vendor/:vendorId", billHandler.GetBillsByVendor)
		billRoutes.Get("/status/:status", billHandler.GetBillsByStatus)
	}

	//
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "github.com/bbapp-org/auth-service",
			"version": "1.0.0",
		})
	})

}
