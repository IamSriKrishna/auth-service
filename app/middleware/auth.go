package middleware

import (
	"strings"

	"github.com/bbapp-org/auth-service/app/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Allow OPTIONS preflight to pass through without auth
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Authorization header is required",
			})
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid authorization header format",
			})
		}

		token := tokenParts[1]

		// Validate token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid or expired token",
			})
		}

		// Store user claims in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_type", claims.UserType)
		c.Locals("user_role", claims.Role)
		c.Locals("user_email", claims.Email)
		c.Locals("user_phone", claims.Phone)
		c.Locals("user_claims", claims)

		return c.Next()
	}
}

// SuperAdminMiddleware validates super admin role
func SuperAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "superadmin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Super admin access required",
			})
		}
		return c.Next()
	}
}

// AdminMiddleware validates admin or super admin role
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "admin" && userType != "superadmin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Admin access required",
			})
		}
		return c.Next()
	}
}

// PartnerMiddleware validates partner or admin role
func PartnerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "partner" && userType != "admin" && userType != "superadmin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Partner access required",
			})
		}
		return c.Next()
	}
}

// MobileUserMiddleware validates mobile user role
func MobileUserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "mobile_user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Mobile user access required",
			})
		}
		return c.Next()
	}
}
