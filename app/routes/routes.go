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

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()

	httpClient := utils.NewHTTPClient(cfg.Service.CustomerServiceURL, 10*time.Second)

	// Initialize Cloudinary
	cloudinaryClient, err := utils.InitCloudinary(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		app.Static("/", "./public", fiber.Static{
			Compress:  true,
			ByteRange: true,
		})
		// Log error but continue - Cloudinary is optional
		app.Use(func(c *fiber.Ctx) error {
			// You can log here if needed
			return c.Next()
		})
	}

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
	productRepo := repo.NewProductRepository(db)
	rawBagRepo := repo.NewRawMaterialBagRepository(db)
	claimRepo := repo.NewVendorShortageClaimRepository(db)
	customerRepo := repo.NewCustomerRepository(db)
	openStockRepo := repo.NewOpeningStockRepository(db)
	manufacturerRepo := repo.NewManufacturerRepository(db)
	invoiceRepo := repo.NewInvoiceRepository(db)
	salespersonRepo := repo.NewSalespersonRepository(db)
	taxRepo := repo.NewTaxRepository(db)
	paymentRepo := repo.NewPaymentRepository(db)
	purchaseOrderRepo := repo.NewPurchaseOrderRepository(db)
	vendorPaymentRepo := repo.NewVendorPaymentRepository(db)
	customerPaymentRepo := repo.NewCustomerPaymentRepository(db)
	salesOrderRepo := repo.NewSalesOrderRepository(db)
	packageRepo := repo.NewPackageRepository(db)
	shipmentRepo := repo.NewShipmentRepository(db)
	billRepo := repo.NewBillRepository(db)
	bankRepo := repo.NewBankRepository(db)
	inventoryBalanceRepo := repo.NewInventoryBalanceRepository(db)
	itemGroupRepo := repo.NewItemGroupRepository(db)
	productGroupRepo := repo.NewProductGroupRepository(db)
	dashboardRepo := repo.NewDashboardRepository(db)
	employeeRepo := repo.NewEmployeeRepository(db)
	attendanceRepo := repo.NewEmployeeAttendanceRepository(db)
	salaryRepo := repo.NewSalaryRepository(db)
	productStockRepo := repo.NewProductStockRepository(db)
	stockLedgerRepo := repo.NewStockLedgerRepository(db)
	variantStockRepo := repo.NewVariantStockRepository(db)
	variantMovementRepo := repo.NewVariantStockMovementRepository(db)
	reservationRepo := repo.NewStockReservationRepository(db)
	pgInventoryRepo := repo.NewProductGroupInventoryRepository(db)
	compInventoryRepo := repo.NewComponentInventoryRepository(db)
	pgTransactionRepo := repo.NewProductGroupTransactionRepository(db)
	customerPricingRepo := repo.NewCustomerPricingRepository(db)
	productConversionRepo := repo.NewProductConversionRepository(db)
	productConversionRecordRepo := repo.NewProductConversionRecordRepository(db)
	conversionRecordBagUsageRepo := repo.NewConversionRecordBagUsageRepository(db)
	purchaseClaimRepo := repo.NewPurchaseClaimRepository(db)
	purchaseDispenseRepo := repo.NewPurchaseDispenseRepository(db)

	authService := services.NewAuthService(userRepo, roleRepo, refreshTokenRepo, sessionRepo, companyRepo)
	adminService := services.NewAdminService(userRepo, roleRepo, companyRepo)

	rawBagService := services.NewRawMaterialBagService(
		rawBagRepo,
		claimRepo,
		purchaseOrderRepo,
		productRepo,
		userRepo,
		companyRepo,
	)

	purchaseClaimService := services.NewPurchaseClaimService(
		purchaseClaimRepo,
		purchaseOrderRepo,
		productRepo,
	)
	purchaseDispenseService :=
		services.NewPurchaseDispenseService(
			purchaseDispenseRepo,
			purchaseClaimRepo,
			purchaseOrderRepo,
			rawBagRepo,
		)

	roleService := services.NewRoleService(roleRepo)
	supportService := services.NewSupportService(supportRepo)
	businessTypeService := services.NewBusinessTypeService(businessTypeRepo)
	locationService := services.NewLocationService(locationRepo)
	taxTypeService := services.NewTaxTypeService(taxTypeRepo)
	companyService := services.NewCompanyService(companyRepo, businessTypeRepo, locationRepo, taxTypeRepo, db)
	itemService := services.NewItemService(itemRepo, vendorRepo, manufacturerRepo, inventoryBalanceRepo, userRepo, companyRepo)
	productService := services.NewProductService(productRepo, vendorRepo, inventoryBalanceRepo, userRepo, companyRepo)
	vendorService := services.NewVendorService(vendorRepo, companyRepo)
	customerService := services.NewCustomerService(customerRepo)
	openStockService := services.NewOpeningStockService(openStockRepo, itemRepo, inventoryBalanceRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo, itemRepo, customerRepo, salespersonRepo, taxRepo, paymentRepo, productRepo, productStockRepo, stockLedgerRepo, userRepo, "./pdf_outputs")
	salespersonService := services.NewSalespersonService(salespersonRepo)
	taxService := services.NewTaxService(taxRepo)
	paymentService := services.NewPaymentService(paymentRepo, invoiceRepo)

	// Stock Movement Service for old services (will be migrated to StockManagementService)
	stockMovementService := services.NewStockMovementService(inventoryBalanceRepo, itemRepo)

	// Stock Management Service for inventory tracking
	stockManagementService := services.NewStockManagementServiceWithVariant(productStockRepo, stockLedgerRepo, productRepo, variantStockRepo)

	// Variant Stock Management Service for SKU-level tracking
	variantStockManagementService := services.NewVariantStockManagementService(variantStockRepo, variantMovementRepo, reservationRepo, productStockRepo, stockLedgerRepo, productRepo, db)

	// Product Group Inventory Service for tracking product group stock
	productGroupInventoryService := services.NewProductGroupInventoryService(pgInventoryRepo, compInventoryRepo, pgTransactionRepo, productGroupRepo, productRepo)

	// Manufacturer Service with dependencies
	manufacturerService := services.NewManufacturerServiceWithDependencies(manufacturerRepo, productGroupRepo, employeeRepo, productStockRepo, userRepo, productRepo, stockManagementService)

	purchaseOrderService := services.NewPurchaseOrderService(purchaseOrderRepo, vendorRepo, customerRepo, productRepo, productGroupRepo, taxRepo, userRepo, companyRepo, stockManagementService, stockLedgerRepo, variantStockManagementService)
	vendorPaymentService := services.NewVendorPaymentService(vendorPaymentRepo, purchaseOrderRepo, vendorRepo, stockManagementService, variantStockManagementService, userRepo)
	customerPaymentService := services.NewCustomerPaymentService(customerPaymentRepo, salesOrderRepo, customerRepo, userRepo)
	salesOrderService := services.NewSalesOrderService(salesOrderRepo, customerRepo, itemRepo, taxRepo, salespersonRepo, inventoryBalanceRepo, stockMovementService, productStockRepo, stockLedgerRepo, variantStockManagementService, stockManagementService, productGroupInventoryService)
	packageService := services.NewPackageService(packageRepo, salesOrderRepo, customerRepo, productRepo, stockManagementService)
	shipmentService := services.NewShipmentService(shipmentRepo, packageRepo, salesOrderRepo, customerRepo)
	billService := services.NewBillService(billRepo, vendorRepo, productRepo, taxRepo, userRepo)
	bankService := services.NewBankService(bankRepo)
	itemGroupService := services.NewItemGroupService(itemGroupRepo, itemRepo)
	productGroupService := services.NewProductGroupServiceWithDependencies(
		productGroupRepo,
		productRepo,
		variantStockManagementService,
		productGroupInventoryService,
		stockManagementService,
		productStockRepo,
		userRepo,
	)
	dashboardService := services.NewDashboardService(dashboardRepo, userRepo, companyRepo)
	employeeService := services.NewEmployeeService(employeeRepo, userRepo, cloudinaryClient)
	attendanceService := services.NewEmployeeAttendanceService(attendanceRepo, employeeRepo)
	salaryService := services.NewSalaryService(salaryRepo, employeeRepo, attendanceRepo)
	customerPricingService := services.NewCustomerPricingServiceWithDependencies(customerPricingRepo, customerRepo, productRepo, userRepo)
	productConversionService := services.NewProductConversionService(
		productConversionRepo, productConversionRecordRepo, conversionRecordBagUsageRepo, productRepo,
		stockManagementService, variantStockManagementService, rawBagService, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	roleHandler := handlers.NewRoleHandler(roleService)
	supportHandler := handlers.NewSupportHandler(supportService)
	forwardAuthHandler := handlers.NewForwardAuthHandler()

	rawBagHandler := handlers.NewRawMaterialBagHandler(rawBagService)

	vendorHandler := handlers.NewVendorHandler(vendorService)
	companyHandler := handlers.NewCompanyHandler(companyService, businessTypeService, locationService, taxTypeService)
	helperHandler := handlers.NewHelperHandler(businessTypeService, locationService, taxTypeService)
	itemHandler := handlers.NewItemHandler(itemService)
	productHandler := handlers.NewProductHandler(productService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	openStockHandler := handlers.NewOpeningStockHandler(openStockService)
	manufacturerHandler := handlers.NewManufacturerHandler(manufacturerService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)
	salespersonHandler := handlers.NewSalespersonHandler(salespersonService)
	taxHandler := handlers.NewTaxHandler(taxService)
	purchaseOrderHandler := handlers.NewPurchaseOrderHandler(purchaseOrderService)
	vendorPaymentHandler := handlers.NewVendorPaymentHandler(vendorPaymentService)
	customerPaymentHandler := handlers.NewCustomerPaymentHandler(customerPaymentService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	salesOrderHandler := handlers.NewSalesOrderHandler(salesOrderService)
	packageHandler := handlers.NewPackageHandler(packageService)
	shipmentHandler := handlers.NewShipmentHandler(shipmentService)
	billHandler := handlers.NewBillHandler(billService)
	bankHandler := handlers.NewBankHandler(bankService)
	itemGroupHandler := handlers.NewItemGroupHandler(itemGroupService)
	productGroupHandler := handlers.NewProductGroupHandler(productGroupService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, userRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	salaryHandler := handlers.NewSalaryHandler(salaryService)
	stockManagementHandler := handlers.NewStockManagementHandlerWithUserRepo(stockManagementService, variantStockManagementService, userRepo, productRepo)
	customerPricingHandler := handlers.NewCustomerPricingHandler(customerPricingService)
	productConversionHandler := handlers.NewProductConversionHandler(productConversionService)
	purchaseClaimHandler := handlers.NewPurchaseClaimHandler(
		purchaseClaimService,
	)
	purchaseDispenseHandler := handlers.NewPurchaseDispenseHandler(
		purchaseDispenseService,
	)

	// Swagger documentation endpoint
	app.Get("/docs/*", swagger.HandlerDefault)

	authGroup := app.Group("/auth")
	{
		authGroup.Post("/register/email", authHandler.RegisterEmail)
		authGroup.Post("/register/phone", authHandler.RegisterPhone)
		authGroup.Post("/register/google", authHandler.RegisterGoogle)

		authGroup.Post("/login/email", authHandler.LoginEmail)
		authGroup.Post("/login/phone", authHandler.LoginPhone)
		authGroup.Post("/login/google", authHandler.LoginGoogle)
		authGroup.Post("/login/apple", authHandler.LoginApple)
		authGroup.Post("/login/password", authHandler.LoginPassword)

		authGroup.Post("/validate-token", authHandler.ValidateToken)
		authGroup.Post("/create-super-admin", adminHandler.CreateSuperAdmin)
		authGroup.Post("/admin/create-user", adminHandler.CreateUser)
	}

	protectedAuthGroup := app.Group("/auth")
	protectedAuthGroup.Use(middleware.AuthMiddleware())
	{
		protectedAuthGroup.Post("/refresh-token", authHandler.RefreshToken)
		protectedAuthGroup.Get("/user-info", authHandler.GetUserInfo)
		protectedAuthGroup.Post("/change-password", authHandler.ChangePassword)
		protectedAuthGroup.Post("/logout", authHandler.Logout)
	}

	manufacturerGroup := app.Group("/manufacturers")
	{
		manufacturerGroup.Get("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), manufacturerHandler.GetAllManufacturers)
		manufacturerGroup.Get("/product-group/:product_group_id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), manufacturerHandler.GetManufacturersByProductGroup)
		manufacturerGroup.Get("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), manufacturerHandler.GetManufacturerByID)
		manufacturerGroup.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), manufacturerHandler.CreateManufacturer)
		manufacturerGroup.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), manufacturerHandler.UpdateManufacturer)
		manufacturerGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), manufacturerHandler.DeleteManufacturer)
	}

	bankGroup := app.Group("/banks")
	{
		bankGroup.Get("/", bankHandler.GetAllBanks)
		bankGroup.Get("/:id", bankHandler.GetBankByID)
		bankGroup.Post("/", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), bankHandler.CreateBank)
		bankGroup.Put("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), bankHandler.UpdateBank)
		bankGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.SuperAdminMiddleware(), bankHandler.DeleteBank)
	}

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

	purchaseClaimGroup := app.Group("/purchase-claims")
	{
		purchaseClaimGroup.Get(
			"/purchase-orders/:purchaseOrderId/items",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseClaimHandler.GetPurchaseOrderClaimSource,
		)

		purchaseClaimGroup.Get(
			"/purchase-orders/:purchaseOrderId",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseClaimHandler.GetClaimsByPurchaseOrder,
		)

		purchaseClaimGroup.Get(
			"/items/:itemId/replacement-receipts",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseClaimHandler.GetReplacementReceipts,
		)

		purchaseClaimGroup.Get(
			"/:id",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseClaimHandler.GetClaimByID,
		)

		purchaseClaimGroup.Post(
			"/",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseClaimHandler.CreateClaim,
		)

		purchaseClaimGroup.Post(
			"/:id/receive-replacement",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseClaimHandler.ReceiveReplacement,
		)
	}

	purchaseDispenseGroup := app.Group("/purchase-dispenses")
	{
		purchaseDispenseGroup.Post(
			"/claims/:claimId",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseDispenseHandler.CreateDispense,
		)

		purchaseDispenseGroup.Get(
			"/claims/:claimId",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseDispenseHandler.GetDispensesByClaim,
		)

		purchaseDispenseGroup.Get(
			"/claim-items/:itemId",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseDispenseHandler.GetDispensesByClaimItem,
		)

		purchaseDispenseGroup.Get(
			"/:id",
			middleware.AuthMiddleware(),
			middleware.AdminMiddleware(),
			purchaseDispenseHandler.GetDispenseByID,
		)
	}

	rawMaterialBagGroup := app.Group("/raw-material-bags")
	{
		rawMaterialBagGroup.Post("/receive", middleware.AuthMiddleware(), middleware.AdminMiddleware(), rawBagHandler.ReceiveBags)

		rawMaterialBagGroup.Get("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), rawBagHandler.GetAll)
		rawMaterialBagGroup.Get("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), rawBagHandler.GetByID)
		rawMaterialBagGroup.Get("/product/:product_id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), rawBagHandler.GetByProduct)
		rawMaterialBagGroup.Get("/purchase-order/:po_id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), rawBagHandler.GetByPurchaseOrder)

	}

	customerGroup := app.Group("/customers")
	{
		customerGroup.Get("/", customerHandler.GetAllCustomers)
		customerGroup.Get("/:id", customerHandler.GetCustomerByID)
		customerGroup.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), customerHandler.CreateCustomer)
		customerGroup.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), customerHandler.UpdateCustomer)
		customerGroup.Delete("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), customerHandler.DeleteCustomer)
	}

	roleGroup := app.Group("/roles")
	roleGroup.Use(middleware.AuthMiddleware())
	roleGroup.Use(middleware.SuperAdminMiddleware())
	{
		roleGroup.Post("/", roleHandler.CreateRole)
		roleGroup.Get("/:id", roleHandler.GetRole)
		roleGroup.Get("/", roleHandler.GetAllRoles)
	}

	partners := app.Group("/partners")
	{
		partners.Patch("/:partner_id/reset-password", adminHandler.ResetUserPassword)
	}

	adminGroup := app.Group("/auth/manage")
	adminGroup.Use(middleware.AuthMiddleware())
	adminGroup.Use(middleware.AdminMiddleware())
	{
		adminGroup.Post("/create-partner", adminHandler.CreateUser)
		adminGroup.Get("/partners", adminHandler.GetUsers)

		adminGroup.Post("/employees", employeeHandler.CreateEmployee)
		adminGroup.Get("/employees", employeeHandler.GetEmployees)
		adminGroup.Get("/employees/:id", employeeHandler.GetEmployee)
		adminGroup.Put("/employees/:id", employeeHandler.UpdateEmployee)
		adminGroup.Patch("/employees/:id", employeeHandler.UpdateEmployee)
		adminGroup.Delete("/employees/:id", employeeHandler.DeleteEmployee)

		// Attendance Routes
		adminGroup.Post("/attendance", attendanceHandler.CreateAttendance)
		adminGroup.Post("/attendance/bulk", attendanceHandler.BulkCreateAttendance)
		adminGroup.Get("/attendance/:id", attendanceHandler.GetAttendanceByID)
		adminGroup.Get("/attendance/employee/:employee_id", attendanceHandler.GetAttendanceByEmployeeID)
		adminGroup.Get("/attendance/employee/:employee_id/calendar", attendanceHandler.GetAttendanceCalendarView)
		adminGroup.Get("/attendance/company/week-view", attendanceHandler.GetCompanyAttendanceWeekView)
		adminGroup.Get("/attendance", attendanceHandler.GetAttendanceByCompanyID)
		adminGroup.Get("/attendance/date-range", attendanceHandler.GetAttendanceByDateRange)
		adminGroup.Get("/attendance/employee/:employee_id/date-range", attendanceHandler.GetAttendanceByEmployeeAndDateRange)
		adminGroup.Put("/attendance/:id", attendanceHandler.UpdateAttendance)
		adminGroup.Delete("/attendance/:id", attendanceHandler.DeleteAttendance)
		adminGroup.Get("/attendance/stats/report", attendanceHandler.GetAttendanceStats)
		adminGroup.Post("/attendance/checkin/:employee_id", attendanceHandler.CheckInEmployee)
		adminGroup.Post("/attendance/checkout/:employee_id", attendanceHandler.CheckOutEmployee)

		// Salary Routes
		adminGroup.Post("/salary/calculate", salaryHandler.CalculateSalary)
		adminGroup.Get("/salary/:id", salaryHandler.GetSalaryCalculation)
		adminGroup.Get("/salary/employee/:employee_id", salaryHandler.GetSalaryCalculationsByEmployee)
		adminGroup.Post("/salary/:id/approve", salaryHandler.ApproveSalary)

		adminGroup.Post("/vendors", vendorHandler.CreateVendor)
		adminGroup.Get("/vendors", vendorHandler.GetUserVendors)
		adminGroup.Get("/vendors/:id", vendorHandler.GetUserVendor)
		adminGroup.Put("/vendors/:id", vendorHandler.UpdateUserVendor)
		adminGroup.Delete("/vendors/:id", vendorHandler.DeleteUserVendor)

		adminGroup.Post("/customers", customerHandler.CreateCustomer)
		adminGroup.Get("/customers", customerHandler.GetUserCustomers)
		adminGroup.Get("/customers/:id", customerHandler.GetUserCustomerByID)
		adminGroup.Put("/customers/:id", customerHandler.UpdateUserCustomer)
		adminGroup.Delete("/customers/:id", customerHandler.DeleteUserCustomer)

		// Items Routes
		adminGroup.Post("/items", itemHandler.CreateItem)
		adminGroup.Get("/items", itemHandler.GetAllItems)
		adminGroup.Get("/items/:id", itemHandler.GetItem)
		adminGroup.Put("/items/:id", itemHandler.UpdateItem)
		adminGroup.Delete("/items/:id", itemHandler.DeleteItem)
		adminGroup.Get("/items/by-type/:type", itemHandler.GetItemsByType)

		// Item Group Routes
		adminGroup.Post("/item-groups", itemGroupHandler.CreateItemGroup)
		adminGroup.Get("/item-groups", itemGroupHandler.GetAllItemGroups)
		adminGroup.Get("/item-groups/:id", itemGroupHandler.GetItemGroupByID)
		adminGroup.Put("/item-groups/:id", itemGroupHandler.UpdateItemGroup)
		adminGroup.Delete("/item-groups/:id", itemGroupHandler.DeleteItemGroup)
		adminGroup.Get("/item-groups/search/by-name", itemGroupHandler.GetItemGroupByName)

		// Bill Routes
		adminGroup.Post("/bills", billHandler.CreateBill)
		adminGroup.Get("/bills", billHandler.GetAllBills)
		adminGroup.Get("/bills/:id", billHandler.GetBill)
		adminGroup.Put("/bills/:id", billHandler.UpdateBill)
		adminGroup.Delete("/bills/:id", billHandler.DeleteBill)
		adminGroup.Patch("/bills/:id/status", billHandler.UpdateBillStatus)
		adminGroup.Get("/bills/vendor/:vendorId", billHandler.GetBillsByVendor)

		// Invoice Routes
		adminGroup.Post("/invoices", invoiceHandler.CreateInvoice)
		adminGroup.Get("/invoices", invoiceHandler.GetAllInvoices)
		adminGroup.Get("/invoices/:id", invoiceHandler.GetInvoice)
		adminGroup.Put("/invoices/:id", invoiceHandler.UpdateInvoice)
		adminGroup.Delete("/invoices/:id", invoiceHandler.DeleteInvoice)
		adminGroup.Patch("/invoices/:id/status", invoiceHandler.UpdateInvoiceStatus)
		adminGroup.Get("/invoices/:invoiceId/payments", paymentHandler.GetPaymentsByInvoice)

		// Purchase Order Routes
		adminGroup.Post("/purchase-orders", purchaseOrderHandler.CreatePurchaseOrder)
		adminGroup.Get("/purchase-orders", purchaseOrderHandler.GetAllPurchaseOrders)
		adminGroup.Get("/purchase-orders/:id", purchaseOrderHandler.GetPurchaseOrder)
		adminGroup.Put("/purchase-orders/:id", purchaseOrderHandler.UpdatePurchaseOrder)
		adminGroup.Delete("/purchase-orders/:id", purchaseOrderHandler.DeletePurchaseOrder)
		adminGroup.Patch("/purchase-orders/:id/status", purchaseOrderHandler.UpdatePurchaseOrderStatus)
		adminGroup.Get("/purchase-orders/vendor/:vendorId", purchaseOrderHandler.GetPurchaseOrdersByVendor)

		// Vendor Payment Routes
		adminGroup.Post("/vendor-payments", vendorPaymentHandler.CreateVendorPayment)
		adminGroup.Get("/vendor-payments", vendorPaymentHandler.GetAllVendorPayments)
		adminGroup.Get("/vendor-payments/:id", vendorPaymentHandler.GetVendorPayment)
		adminGroup.Put("/vendor-payments/:id", vendorPaymentHandler.UpdateVendorPayment)
		adminGroup.Delete("/vendor-payments/:id", vendorPaymentHandler.DeleteVendorPayment)
		adminGroup.Post("/vendor-payments/:id/record-payment", vendorPaymentHandler.RecordPayment)
		adminGroup.Get("/vendor-payments/purchase-order/:purchaseOrderId", vendorPaymentHandler.GetVendorPaymentsByPurchaseOrder)
		adminGroup.Get("/vendor-payments/vendor/:vendorId", vendorPaymentHandler.GetVendorPaymentsByVendor)

		// Customer Payment Routes
		adminGroup.Post("/customer-payments", customerPaymentHandler.CreateCustomerPayment)
		adminGroup.Get("/customer-payments", customerPaymentHandler.GetAllCustomerPayments)
		adminGroup.Get("/customer-payments/:id", customerPaymentHandler.GetCustomerPayment)
		adminGroup.Put("/customer-payments/:id", customerPaymentHandler.UpdateCustomerPayment)
		adminGroup.Delete("/customer-payments/:id", customerPaymentHandler.DeleteCustomerPayment)
		adminGroup.Post("/customer-payments/:id/record-payment", customerPaymentHandler.RecordPayment)
		adminGroup.Get("/customer-payments/sales-order/:salesOrderId", customerPaymentHandler.GetCustomerPaymentsBySalesOrder)
		adminGroup.Get("/customer-payments/customer/:customerId", customerPaymentHandler.GetCustomerPaymentsByCustomer)

		// Sales Order Routes
		adminGroup.Post("/sales-orders", middleware.AdminMiddleware(), salesOrderHandler.CreateSalesOrder)
		adminGroup.Get("/sales-orders", middleware.AdminMiddleware(), salesOrderHandler.GetAllSalesOrders)
		adminGroup.Get("/sales-orders/:id", middleware.AdminMiddleware(), salesOrderHandler.GetSalesOrder)
		adminGroup.Put("/sales-orders/:id", middleware.AdminMiddleware(), salesOrderHandler.UpdateSalesOrder)
		adminGroup.Delete("/sales-orders/:id", middleware.AdminMiddleware(), salesOrderHandler.DeleteSalesOrder)
		adminGroup.Patch("/sales-orders/:id/status", middleware.AdminMiddleware(), salesOrderHandler.UpdateSalesOrderStatus)
		adminGroup.Get("/sales-orders/customer/:customerId", middleware.AdminMiddleware(), salesOrderHandler.GetSalesOrdersByCustomer)

		// Package Routes
		adminGroup.Post("/packages", packageHandler.CreatePackage)
		adminGroup.Get("/packages", packageHandler.GetAllPackages)
		adminGroup.Get("/packages/:id", packageHandler.GetPackage)
		adminGroup.Put("/packages/:id", packageHandler.UpdatePackage)
		adminGroup.Delete("/packages/:id", packageHandler.DeletePackage)
		adminGroup.Patch("/packages/:id/status", packageHandler.UpdatePackageStatus)
		adminGroup.Get("/packages/customer/:customer_id", packageHandler.GetPackagesByCustomer)

		// Shipment Routes
		adminGroup.Post("/shipments", shipmentHandler.CreateShipment)
		adminGroup.Get("/shipments", shipmentHandler.GetAllShipments)
		adminGroup.Get("/shipments/:id", shipmentHandler.GetShipment)
		adminGroup.Put("/shipments/:id", shipmentHandler.UpdateShipment)
		adminGroup.Delete("/shipments/:id", shipmentHandler.DeleteShipment)
		adminGroup.Patch("/shipments/:id/status", shipmentHandler.UpdateShipmentStatus)
		adminGroup.Get("/shipments/customer/:customer_id", shipmentHandler.GetShipmentsByCustomer)
	}

	// Admin group for /admin prefix routes
	adminPrefixGroup := app.Group("/admin")
	adminPrefixGroup.Use(middleware.AuthMiddleware())
	adminPrefixGroup.Use(middleware.AdminMiddleware())
	{
		// Salary Routes
		adminPrefixGroup.Post("/salary/calculate", salaryHandler.CalculateSalary)
		adminPrefixGroup.Get("/salary/:id", salaryHandler.GetSalaryCalculation)
		adminPrefixGroup.Get("/salary/employee/:employee_id", salaryHandler.GetSalaryCalculationsByEmployee)
		adminPrefixGroup.Post("/salary/:id/approve", salaryHandler.ApproveSalary)
	}

	forwardAuthGroup := app.Group("/forward-auth")
	{
		forwardAuthGroup.Get("/", forwardAuthHandler.ForwardAuth)
		forwardAuthGroup.Get("/product", forwardAuthHandler.ProductAuth)
		forwardAuthGroup.Get("/customer", forwardAuthHandler.CustomerAuth)
	}

	app.Post("/public/support", supportHandler.CreateSupport)

	helperRoutes := app.Group("/helpers")
	{
		helperRoutes.Get("/business-types", helperHandler.GetBusinessTypes)
		helperRoutes.Get("/countries", helperHandler.GetCountries)
		helperRoutes.Get("/countries/:country_id/states", helperHandler.GetStatesByCountry)
		helperRoutes.Get("/tax-types", helperHandler.GetTaxTypes)
	}

	companyRoutes := app.Group("/companies")
	companyRoutes.Use(middleware.AuthMiddleware())
	{
		companyRoutes.Post("/setup", companyHandler.CompleteCompanySetup)
		companyRoutes.Get("/me", companyHandler.GetMyCompanyProfile)

		companyRoutes.Get("/", companyHandler.GetAllCompanies)
		companyRoutes.Post("/", companyHandler.CreateCompany)
		companyRoutes.Get("/:id", companyHandler.GetCompany)
		companyRoutes.Put("/:id", companyHandler.UpdateCompany)
		companyRoutes.Delete("/:id", middleware.SuperAdminMiddleware(), companyHandler.DeleteCompany)

		companyRoutes.Put("/:id/contact", companyHandler.UpsertContact)
		companyRoutes.Get("/:id/contact", companyHandler.GetContact)

		companyRoutes.Put("/:id/address", companyHandler.UpsertAddress)
		companyRoutes.Get("/:id/address", companyHandler.GetAddress)

		companyRoutes.Post("/:id/bank-details", companyHandler.CreateBankDetail)
		companyRoutes.Get("/:id/bank-details", companyHandler.GetBankDetails)
		companyRoutes.Put("/bank-details/:id", companyHandler.UpdateBankDetail)
		companyRoutes.Delete("/bank-details/:id", companyHandler.DeleteBankDetail)

		companyRoutes.Put("/:id/upi-details", companyHandler.UpsertUPIDetail)
		companyRoutes.Get("/:id/upi-details", companyHandler.GetUPIDetail)

		companyRoutes.Put("/:id/invoice-settings", companyHandler.UpsertInvoiceSettings)
		companyRoutes.Get("/:id/invoice-settings", companyHandler.GetInvoiceSettings)

		companyRoutes.Put("/:id/tax-settings", companyHandler.UpsertTaxSettings)
		companyRoutes.Get("/:id/tax-settings", companyHandler.GetTaxSettings)

		companyRoutes.Put("/:id/regional-settings", companyHandler.UpsertRegionalSettings)
		companyRoutes.Get("/:id/regional-settings", companyHandler.GetRegionalSettings)
	}

	itemRoutes := app.Group("/items")
	{
		itemRoutes.Get("/:id", itemHandler.GetItem)

		itemRoutes.Get("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemHandler.GetAllItems)
		itemRoutes.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemHandler.CreateItem)
		itemRoutes.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemHandler.UpdateItem)
		itemRoutes.Delete("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemHandler.DeleteItem)

		itemRoutes.Put("/:id/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.UpdateOpeningStock)
		itemRoutes.Get("/:id/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.GetOpeningStock)

		itemRoutes.Put("/:id/variants/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.UpdateVariantsOpeningStock)
		itemRoutes.Get("/:id/variants/opening-stock", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.GetVariantsOpeningStock)
		itemRoutes.Get("/:id/stock-summary", middleware.AuthMiddleware(), middleware.AdminMiddleware(), openStockHandler.GetStockSummary)
	}

	productRoutes := app.Group("/products")
	{
		productRoutes.Get("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandler.GetProduct)

		productRoutes.Get("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandler.GetAllProducts)
		productRoutes.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandler.CreateProduct)
		productRoutes.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandler.UpdateProduct)
		productRoutes.Delete("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandler.DeleteProduct)

		productRoutes.Get("/:id/variants", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandler.GetProductVariants)
	}

	itemGroupRoutes := app.Group("/item-groups")
	{
		itemGroupRoutes.Get("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemGroupHandler.GetAllItemGroups)
		itemGroupRoutes.Get("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemGroupHandler.GetItemGroupByID)

		itemGroupRoutes.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemGroupHandler.CreateItemGroup)
		itemGroupRoutes.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemGroupHandler.UpdateItemGroup)
		itemGroupRoutes.Delete("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), itemGroupHandler.DeleteItemGroup)

		itemGroupRoutes.Get("/search/by-name", itemGroupHandler.GetItemGroupByName)
	}

	productGroupRoutes := app.Group("/product-groups")
	{
		// More specific routes first
		productGroupRoutes.Get("/search/by-name", productGroupHandler.GetProductGroupByName)
		productGroupRoutes.Post("/:id/reorder", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productGroupHandler.ReorderProductGroup)

		// Then generic routes
		productGroupRoutes.Get("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productGroupHandler.GetAllProductGroups)
		productGroupRoutes.Get("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productGroupHandler.GetProductGroupByID)

		productGroupRoutes.Post("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productGroupHandler.CreateProductGroup)
		productGroupRoutes.Put("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productGroupHandler.UpdateProductGroup)
		productGroupRoutes.Delete("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productGroupHandler.DeleteProductGroup)
	}

	invoiceRoutes := app.Group("/invoices")
	invoiceRoutes.Use(middleware.AuthMiddleware())
	{
		invoiceRoutes.Post("/", middleware.AdminMiddleware(), invoiceHandler.CreateInvoice)
		invoiceRoutes.Get("/", middleware.AdminMiddleware(), invoiceHandler.GetAllInvoices)

		// More specific routes must be registered before generic parameter routes
		invoiceRoutes.Patch("/:id/status", middleware.AdminMiddleware(), invoiceHandler.UpdateInvoiceStatus)
		invoiceRoutes.Get("/:invoiceId/payments", paymentHandler.GetPaymentsByInvoice)

		invoiceRoutes.Get("/:id", middleware.AdminMiddleware(), invoiceHandler.GetInvoice)
		invoiceRoutes.Put("/:id", middleware.AdminMiddleware(), invoiceHandler.UpdateInvoice)
		invoiceRoutes.Delete("/:id", middleware.AdminMiddleware(), invoiceHandler.DeleteInvoice)
	}

	app.Get("/invoices/status/:status", middleware.AuthMiddleware(), invoiceHandler.GetInvoicesByStatus)

	customerGroup.Get("/:customerId/invoices", middleware.AuthMiddleware(), invoiceHandler.GetInvoicesByCustomer)

	salespersonRoutes := app.Group("/salespersons")
	salespersonRoutes.Use(middleware.AuthMiddleware())
	{
		salespersonRoutes.Post("/", middleware.AdminMiddleware(), salespersonHandler.CreateSalesperson)
		salespersonRoutes.Get("/", salespersonHandler.GetAllSalespersons)
		salespersonRoutes.Get("/:id", salespersonHandler.GetSalesperson)
		salespersonRoutes.Put("/:id", middleware.AdminMiddleware(), salespersonHandler.UpdateSalesperson)
		salespersonRoutes.Delete("/:id", middleware.SuperAdminMiddleware(), salespersonHandler.DeleteSalesperson)
	}

	taxRoutes := app.Group("/taxes")
	taxRoutes.Use(middleware.AuthMiddleware())
	{
		taxRoutes.Post("/", middleware.AdminMiddleware(), taxHandler.CreateTax)
		taxRoutes.Get("/", taxHandler.GetAllTaxes)
		taxRoutes.Get("/:id", taxHandler.GetTax)
		taxRoutes.Put("/:id", middleware.AdminMiddleware(), taxHandler.UpdateTax)
		taxRoutes.Delete("/:id", middleware.SuperAdminMiddleware(), taxHandler.DeleteTax)
	}

	paymentRoutes := app.Group("/payments")
	paymentRoutes.Use(middleware.AuthMiddleware())
	{
		paymentRoutes.Post("/", middleware.AdminMiddleware(), paymentHandler.CreatePayment)
		paymentRoutes.Get("/:id", paymentHandler.GetPayment)
		paymentRoutes.Delete("/:id", middleware.AdminMiddleware(), paymentHandler.DeletePayment)
	}
	purchaseOrderRoutes := app.Group("/purchase-orders")
	purchaseOrderRoutes.Use(middleware.AuthMiddleware())
	{
		purchaseOrderRoutes.Post("/", middleware.AdminMiddleware(), purchaseOrderHandler.CreatePurchaseOrder)
		purchaseOrderRoutes.Get("/", middleware.AdminMiddleware(), purchaseOrderHandler.GetAllPurchaseOrders)

		// More specific routes must be registered before generic parameter routes
		purchaseOrderRoutes.Patch("/:id/status", middleware.AdminMiddleware(), purchaseOrderHandler.UpdatePurchaseOrderStatus)
		purchaseOrderRoutes.Get("/vendor/:vendorId", purchaseOrderHandler.GetPurchaseOrdersByVendor)
		purchaseOrderRoutes.Get("/customer/:customerId", purchaseOrderHandler.GetPurchaseOrdersByCustomer)
		purchaseOrderRoutes.Get("/status/:status", purchaseOrderHandler.GetPurchaseOrdersByStatus)

		purchaseOrderRoutes.Get("/:id", middleware.AdminMiddleware(), purchaseOrderHandler.GetPurchaseOrder)
		purchaseOrderRoutes.Put("/:id", middleware.AdminMiddleware(), purchaseOrderHandler.UpdatePurchaseOrder)
		purchaseOrderRoutes.Delete("/:id", middleware.AdminMiddleware(), purchaseOrderHandler.DeletePurchaseOrder)
	}

	vendorPaymentRoutes := app.Group("/vendor-payments")
	vendorPaymentRoutes.Use(middleware.AuthMiddleware())
	{
		vendorPaymentRoutes.Post("/", middleware.AdminMiddleware(), vendorPaymentHandler.CreateVendorPayment)
		vendorPaymentRoutes.Get("/", middleware.AdminMiddleware(), vendorPaymentHandler.GetAllVendorPayments)
		vendorPaymentRoutes.Get("/:id", middleware.AdminMiddleware(), vendorPaymentHandler.GetVendorPayment)
		vendorPaymentRoutes.Put("/:id", middleware.AdminMiddleware(), vendorPaymentHandler.UpdateVendorPayment)
		vendorPaymentRoutes.Delete("/:id", middleware.AdminMiddleware(), vendorPaymentHandler.DeleteVendorPayment)

		vendorPaymentRoutes.Post("/:id/record-payment", middleware.AdminMiddleware(), vendorPaymentHandler.RecordPayment)

		vendorPaymentRoutes.Get("/purchase-order/:purchaseOrderId", vendorPaymentHandler.GetVendorPaymentsByPurchaseOrder)
		vendorPaymentRoutes.Get("/vendor/:vendorId", vendorPaymentHandler.GetVendorPaymentsByVendor)
		vendorPaymentRoutes.Get("/status/:status", vendorPaymentHandler.GetVendorPaymentsByStatus)
	}

	customerPaymentRoutes := app.Group("/customer-payments")
	customerPaymentRoutes.Use(middleware.AuthMiddleware())
	{
		customerPaymentRoutes.Post("/", middleware.AdminMiddleware(), customerPaymentHandler.CreateCustomerPayment)
		customerPaymentRoutes.Get("/", middleware.AdminMiddleware(), customerPaymentHandler.GetAllCustomerPayments)

		// Specific routes before wildcard :id
		customerPaymentRoutes.Post("/:id/record-payment", middleware.AdminMiddleware(), customerPaymentHandler.RecordPayment)
		customerPaymentRoutes.Get("/sales-order/:salesOrderId", customerPaymentHandler.GetCustomerPaymentsBySalesOrder)
		customerPaymentRoutes.Get("/customer/:customerId", customerPaymentHandler.GetCustomerPaymentsByCustomer)
		customerPaymentRoutes.Get("/status/:status", customerPaymentHandler.GetCustomerPaymentsByStatus)

		// Wildcard :id routes last
		customerPaymentRoutes.Get("/:id", middleware.AdminMiddleware(), customerPaymentHandler.GetCustomerPayment)
		customerPaymentRoutes.Put("/:id", middleware.AdminMiddleware(), customerPaymentHandler.UpdateCustomerPayment)
		customerPaymentRoutes.Delete("/:id", middleware.AdminMiddleware(), customerPaymentHandler.DeleteCustomerPayment)
	}

	salesOrderRoutes := app.Group("/sales-orders")
	salesOrderRoutes.Use(middleware.AuthMiddleware())
	{
		salesOrderRoutes.Post("/", middleware.AdminMiddleware(), salesOrderHandler.CreateSalesOrder)
		salesOrderRoutes.Get("/", middleware.AdminMiddleware(), salesOrderHandler.GetAllSalesOrders)

		// More specific routes must be registered before generic parameter routes
		salesOrderRoutes.Patch("/:id/status", middleware.AdminMiddleware(), salesOrderHandler.UpdateSalesOrderStatus)
		salesOrderRoutes.Get("/customer/:customerId", middleware.AdminMiddleware(), salesOrderHandler.GetSalesOrdersByCustomer)
		salesOrderRoutes.Get("/status/:status", middleware.AdminMiddleware(), salesOrderHandler.GetSalesOrdersByStatus)

		salesOrderRoutes.Get("/:id", middleware.AdminMiddleware(), salesOrderHandler.GetSalesOrder)
		salesOrderRoutes.Put("/:id", middleware.AdminMiddleware(), salesOrderHandler.UpdateSalesOrder)
		salesOrderRoutes.Delete("/:id", middleware.AdminMiddleware(), salesOrderHandler.DeleteSalesOrder)
	}

	packageRoutes := app.Group("/packages")
	packageRoutes.Use(middleware.AuthMiddleware())
	{
		packageRoutes.Post("/", middleware.AdminMiddleware(), packageHandler.CreatePackage)
		packageRoutes.Get("/", middleware.AdminMiddleware(), packageHandler.GetAllPackages)

		// More specific routes must be registered before generic parameter routes
		packageRoutes.Patch("/:id/status", middleware.AdminMiddleware(), packageHandler.UpdatePackageStatus)
		packageRoutes.Get("/customer/:customer_id", packageHandler.GetPackagesByCustomer)
		packageRoutes.Get("/sales-order/:sales_order_id", packageHandler.GetPackagesBySalesOrder)
		packageRoutes.Get("/status/:status", packageHandler.GetPackagesByStatus)

		packageRoutes.Get("/:id", middleware.AdminMiddleware(), packageHandler.GetPackage)
		packageRoutes.Put("/:id", middleware.AdminMiddleware(), packageHandler.UpdatePackage)
		packageRoutes.Delete("/:id", middleware.AdminMiddleware(), packageHandler.DeletePackage)
	}

	shipmentRoutes := app.Group("/shipments")
	shipmentRoutes.Use(middleware.AuthMiddleware())
	{
		shipmentRoutes.Post("/", middleware.AdminMiddleware(), shipmentHandler.CreateShipment)
		shipmentRoutes.Get("/", shipmentHandler.GetAllShipments)

		// More specific routes must be registered before generic parameter routes
		shipmentRoutes.Patch("/:id/status", middleware.AdminMiddleware(), shipmentHandler.UpdateShipmentStatus)
		shipmentRoutes.Get("/customer/:customer_id", shipmentHandler.GetShipmentsByCustomer)
		shipmentRoutes.Get("/package/:package_id", shipmentHandler.GetShipmentsByPackage)
		shipmentRoutes.Get("/sales-order/:sales_order_id", shipmentHandler.GetShipmentsBySalesOrder)
		shipmentRoutes.Get("/status/:status", shipmentHandler.GetShipmentsByStatus)

		shipmentRoutes.Get("/:id", shipmentHandler.GetShipment)
		shipmentRoutes.Put("/:id", middleware.AdminMiddleware(), shipmentHandler.UpdateShipment)
		shipmentRoutes.Delete("/:id", middleware.AdminMiddleware(), shipmentHandler.DeleteShipment)
	}

	billRoutes := app.Group("/bills")
	billRoutes.Use(middleware.AuthMiddleware())
	{
		billRoutes.Post("/", middleware.AdminMiddleware(), billHandler.CreateBill)
		billRoutes.Get("/", middleware.AdminMiddleware(), billHandler.GetAllBills)

		// More specific routes must be registered before generic parameter routes
		billRoutes.Patch("/:id/status", middleware.AdminMiddleware(), billHandler.UpdateBillStatus)
		billRoutes.Get("/vendor/:vendorId", billHandler.GetBillsByVendor)
		billRoutes.Get("/status/:status", billHandler.GetBillsByStatus)

		billRoutes.Get("/:id", middleware.AdminMiddleware(), billHandler.GetBill)
		billRoutes.Put("/:id", middleware.AdminMiddleware(), billHandler.UpdateBill)
		billRoutes.Delete("/:id", middleware.AdminMiddleware(), billHandler.DeleteBill)
	}

	// Stock Management Routes
	stockRoutes := app.Group("/api/stock")
	stockRoutes.Use(middleware.AuthMiddleware())
	{
		stockRoutes.Get("/summary", middleware.AdminMiddleware(), stockManagementHandler.GetAllStocksSummary)
		stockRoutes.Get("/summary/raw-materials", middleware.AdminMiddleware(), stockManagementHandler.GetRawMaterialStocksSummary)
		stockRoutes.Get("/damaged", middleware.AdminMiddleware(), stockManagementHandler.GetDamagedProducts)
		stockRoutes.Patch("/mark-damaged", middleware.AdminMiddleware(), stockManagementHandler.MarkProductAsDamaged)
		stockRoutes.Get("/product/:product_id/movements", middleware.AdminMiddleware(), stockManagementHandler.GetProductMovements)
		stockRoutes.Get("/debug/product/:product_id", middleware.AdminMiddleware(), stockManagementHandler.GetProductStockDebug)
	}

	// Customer Pricing Routes
	customerPricingRoutes := app.Group("/customer-pricing")
	customerPricingRoutes.Use(middleware.AuthMiddleware())
	{
		customerPricingRoutes.Post("/", middleware.AdminMiddleware(), customerPricingHandler.CreateCustomerPricing)
		customerPricingRoutes.Get("/", middleware.AdminMiddleware(), customerPricingHandler.GetAllCustomerPricing)
		// Define specific routes BEFORE wildcard /:id
		customerPricingRoutes.Get("/customer/active", middleware.AdminMiddleware(), customerPricingHandler.GetActivePricingByCustomer)
		customerPricingRoutes.Get("/customer", middleware.AdminMiddleware(), customerPricingHandler.GetPricingByCustomer)
		// Wildcard routes must be defined LAST
		customerPricingRoutes.Get("/:id", middleware.AdminMiddleware(), customerPricingHandler.GetCustomerPricingByID)
		customerPricingRoutes.Put("/:id", middleware.AdminMiddleware(), customerPricingHandler.UpdateCustomerPricing)
		customerPricingRoutes.Delete("/:id", middleware.AdminMiddleware(), customerPricingHandler.DeleteCustomerPricing)
		customerPricingRoutes.Put("/:id/date-range", middleware.AdminMiddleware(), customerPricingHandler.SetEffectiveDateRange)
	}

	// Product Conversion Routes
	conversionRoutes := app.Group("/product-conversions")
	conversionRoutes.Use(middleware.AuthMiddleware())
	{
		// Conversion Rule Management - Define specific routes BEFORE wildcard /:id
		conversionRoutes.Post("/", middleware.AdminMiddleware(), productConversionHandler.CreateConversion)
		conversionRoutes.Get("/", middleware.AdminMiddleware(), productConversionHandler.ListConversions)
		conversionRoutes.Get("/active", middleware.AdminMiddleware(), productConversionHandler.ListActiveConversions)
		conversionRoutes.Get("/by-raw", middleware.AdminMiddleware(), productConversionHandler.GetConversionsByRawProduct)
		conversionRoutes.Get("/by-finished", middleware.AdminMiddleware(), productConversionHandler.GetConversionsByFinishedProduct)

		// Conversion Records
		conversionRoutes.Get("/records", middleware.AdminMiddleware(), productConversionHandler.ListConversionRecords)
		conversionRoutes.Get("/records/by-rule", middleware.AdminMiddleware(), productConversionHandler.ListConversionRecordsByRule)
		conversionRoutes.Get("/records/by-date-range", middleware.AdminMiddleware(), productConversionHandler.ListConversionRecordsByDateRange)
		conversionRoutes.Get("/records/:record_id", middleware.AdminMiddleware(), productConversionHandler.GetConversionRecord)

		conversionRoutes.Get("/:id", middleware.AdminMiddleware(), productConversionHandler.GetConversion)
		conversionRoutes.Put("/:id", middleware.AdminMiddleware(), productConversionHandler.UpdateConversion)
		conversionRoutes.Delete("/:id", middleware.AdminMiddleware(), productConversionHandler.DeleteConversion)

		// Conversion Execution
		conversionRoutes.Post("/execute", middleware.AdminMiddleware(), productConversionHandler.ExecuteConversion)
	}

	// Dashboard Routes
	dashboardRoutes := app.Group("/dashboard")
	dashboardRoutes.Use(middleware.AuthMiddleware())
	{
		dashboardRoutes.Get("/", middleware.AdminMiddleware(), dashboardHandler.GetDashboard)
		dashboardRoutes.Get("/activity", middleware.AdminMiddleware(), dashboardHandler.GetActivitySummary)
		dashboardRoutes.Get("/stock", middleware.AdminMiddleware(), dashboardHandler.GetStockSummary)
		dashboardRoutes.Get("/trends/:entity_type", middleware.AdminMiddleware(), dashboardHandler.GetEntityTrends)
		dashboardRoutes.Get("/shipment/:shipment_id/tracking", middleware.AdminMiddleware(), dashboardHandler.GetShipmentTracking)
		dashboardRoutes.Post("/shipment/:shipment_id/tracking", middleware.AdminMiddleware(), dashboardHandler.AddShipmentTracking)
		dashboardRoutes.Post("/refresh", middleware.AdminMiddleware(), dashboardHandler.RefreshMetrics)
		dashboardRoutes.Get("/diagnose", middleware.AdminMiddleware(), dashboardHandler.GetDiagnosticReport)
	}

	// Public routes for live status (no authentication required)
	publicRoutes := app.Group("/public")
	{
		publicRoutes.Get("/live-status", dashboardHandler.GetPublicLiveStatus)
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "github.com/bbapp-org/auth-service",
			"version": "1.0.0",
		})
	})

}
