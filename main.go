package main

import (
	"log"

	"github.com/bbapp-org/auth-service/app/config"
	"github.com/bbapp-org/auth-service/app/config/database"
	"github.com/bbapp-org/auth-service/app/helper"
	"github.com/bbapp-org/auth-service/app/routes"

	_ "github.com/bbapp-org/auth-service/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// Build-time variables for version info
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
	gitRef    = "unknown"
)

func main() {
	// Initialize application configuration
	cfg := config.LoadConfig()

	// Initialize database connections
	database.ConnectDatabase(cfg)

	// Run database migrations
	db := database.GetDB()
	if err := helper.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		// Skip logging for health check endpoints to reduce log noise
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/health" || c.Path() == "/v1/health"
		},
	}))
	app.Use(recover.New())

	// CORS - use configured origins (comma-separated)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.App.AllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: false,
	}))

	// Respond to preflight requests with 200 universally (so auth middleware isn't triggered)
	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Next()
	})

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "github.com/bbapp-org/auth-service",
			"version": version,
		})
	})

	// Version endpoint
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":   "github.com/bbapp-org/auth-service",
			"version":   version,
			"buildTime": buildTime,
			"gitCommit": gitCommit,
			"gitRef":    gitRef,
		})
	})

	// Routes
	routes.SetupRoutes(app, cfg)

	// Start server
	port := cfg.App.ServerPort
	log.Printf("Auth Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
