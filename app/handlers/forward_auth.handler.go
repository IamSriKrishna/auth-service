package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/utils"

	"github.com/gofiber/fiber/v2"
)

// ForwardAuthHandler handles Traefik forward authentication
type ForwardAuthHandler struct {
	// Add caching for performance optimization
	cache      map[string]*CacheEntry
	cacheMutex sync.RWMutex
}

// CacheEntry represents a cached authentication result
type CacheEntry struct {
	Valid        bool
	UserID       uint
	UserType     string
	Role         string
	Email        string
	Phone        string
	IdentityType string
	ExpiresAt    time.Time
}

// NewForwardAuthHandler creates a new forward auth handler
func NewForwardAuthHandler() *ForwardAuthHandler {
	handler := &ForwardAuthHandler{
		cache: make(map[string]*CacheEntry),
	}

	// Start cache cleanup goroutine
	go handler.startCacheCleanup()

	return handler
}

// startCacheCleanup periodically cleans expired cache entries
func (h *ForwardAuthHandler) startCacheCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		h.cleanupExpiredEntries()
	}
}

// cleanupExpiredEntries removes expired entries from cache
func (h *ForwardAuthHandler) cleanupExpiredEntries() {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()

	now := time.Now()
	for key, entry := range h.cache {
		if now.After(entry.ExpiresAt) {
			delete(h.cache, key)
		}
	}
}

// getCachedAuth retrieves cached authentication result
func (h *ForwardAuthHandler) getCachedAuth(token string) (*CacheEntry, bool) {
	h.cacheMutex.RLock()
	defer h.cacheMutex.RUnlock()

	entry, exists := h.cache[token]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry, true
}

// ForwardAuth validates requests for Traefik forward authentication
// @Summary Forward authentication for Traefik
// @Description Validates JWT tokens for Traefik forward auth middleware
// @Tags Forward Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param X-Forwarded-Uri header string true "Original request URI"
// @Param X-Forwarded-Method header string true "Original request method"
// @Success 200 {object} map[string]interface{} "Authentication successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /forward-auth [get]
func (h *ForwardAuthHandler) ForwardAuth(c *fiber.Ctx) error {
	// Get original request details from Traefik headers
	originalURI := c.Get("X-Forwarded-Uri")
	originalMethod := c.Get("X-Forwarded-Method")

	// Debug: Try alternative header names if empty
	if originalURI == "" {
		originalURI = c.Get("X-Original-Url")
		if originalURI == "" {
			originalURI = c.Get("X-Forwarded-Prefix")
		}
		if originalURI == "" {
			originalURI = c.Get("X-Original-URI")
		}
	}

	if originalMethod == "" {
		originalMethod = c.Get("X-Original-Method")
	}

	// Enhanced logging for debugging
	fmt.Printf("Forward Auth Request: %s %s from %s\n", originalMethod, originalURI, c.IP())
	fmt.Printf("Debug Headers - X-Forwarded-Uri: '%s', X-Forwarded-Method: '%s'\n", c.Get("X-Forwarded-Uri"), c.Get("X-Forwarded-Method"))
	fmt.Printf("All request headers: %+v\n", c.GetReqHeaders())

	// Check if this is a public endpoint that doesn't require authentication
	if h.isPublicEndpoint(originalURI, originalMethod) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"authenticated": true,
			"public":        true,
			"message":       "Public endpoint access granted",
		})
	}

	// Get authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Authorization header is required",
			"code":    "AUTH_HEADER_MISSING",
		})
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid authorization header format",
			"code":    "INVALID_AUTH_FORMAT",
		})
	}

	token := tokenParts[1]

	// Check cache first - but still need to validate RBAC for protected endpoints
	var claims *output.Claims

	if cachedEntry, found := h.getCachedAuth(token); found {
		fmt.Printf("CACHE HIT: Found cached auth for user %d (%s) - but still checking RBAC for URI: %s\n",
			cachedEntry.UserID, cachedEntry.UserType, originalURI)

		// Create claims from cache for RBAC checking
		claims = &output.Claims{
			UserID:       cachedEntry.UserID,
			UserType:     cachedEntry.UserType,
			Role:         cachedEntry.Role,
			Email:        cachedEntry.Email,
			Phone:        cachedEntry.Phone,
			IdentityType: cachedEntry.IdentityType,
		}
	} else {
		fmt.Printf("CACHE MISS: Validating fresh JWT token for URI: %s\n", originalURI)

		// Validate token
		var err error
		claims, err = utils.ValidateJWT(token)
		if err != nil {
			fmt.Printf("Token validation failed for %s %s: %v\n", originalMethod, originalURI, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid or expired token",
				"code":    "TOKEN_INVALID",
			})
		}
	}

	// Validate user type according to BRD requirements
	if !h.isValidUserType(claims.UserType) {
		fmt.Printf("RBAC: Invalid user type '%s' for URI: %s\n", claims.UserType, originalURI)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user type",
			"code":    "INVALID_USER_TYPE",
		})
	}

	// Check authorization based on the original request - ALWAYS CHECK RBAC
	fmt.Printf("RBAC: Checking authorization for user %d (%s) accessing %s %s\n",
		claims.UserID, claims.UserType, originalMethod, originalURI)

	if !h.isAuthorized(claims, originalURI, originalMethod) {
		fmt.Printf("RBAC: Access DENIED for user type '%s' (ID: %d) to %s %s\n",
			claims.UserType, claims.UserID, originalMethod, originalURI)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Insufficient permissions for this resource",
			"code":    "INSUFFICIENT_PERMISSIONS",
		})
	}

	fmt.Printf("RBAC: Access GRANTED for user %d (%s) to %s %s\n",
		claims.UserID, claims.UserType, originalMethod, originalURI)

	// Cache the valid authentication result (only if not from cache)
	if _, found := h.getCachedAuth(token); !found {
		fmt.Printf("CACHE: Storing auth result for user %d\n", claims.UserID)
		h.cacheMutex.Lock()
		h.cache[token] = &CacheEntry{
			Valid:        true,
			UserID:       claims.UserID,
			UserType:     claims.UserType,
			Role:         claims.Role,
			Email:        claims.Email,
			Phone:        claims.Phone,
			IdentityType: claims.IdentityType,
			ExpiresAt:    time.Now().Add(24 * time.Hour), // Cache for 24 hours
		}
		h.cacheMutex.Unlock()
	}

	// Set user information in response headers for the backend service
	c.Set("X-User-Id", fmt.Sprintf("%d", claims.UserID))
	c.Set("X-User-Type", claims.UserType)
	c.Set("X-User-Role", claims.Role)
	if claims.Email != "" {
		c.Set("X-User-Email", claims.Email)
	}
	if claims.Phone != "" {
		c.Set("X-User-Phone", claims.Phone)
	}
	if claims.IdentityType != "" {
		c.Set("X-Identity-Type", claims.IdentityType)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"authenticated": true,
		"user_id":       claims.UserID,
		"user_type":     claims.UserType,
		"user_role":     claims.Role,
		"identity_type": claims.IdentityType,
	})
}

// isPublicEndpoint checks if the endpoint is publicly accessible without authentication
func (h *ForwardAuthHandler) isPublicEndpoint(uri, method string) bool {
	// Health, version, and documentation endpoints are always public
	if strings.HasPrefix(uri, "/health") ||
		strings.HasPrefix(uri, "/swagger") ||
		strings.HasPrefix(uri, "/docs") ||
		strings.HasPrefix(uri, "/version") ||
		strings.HasPrefix(uri, "/info") ||
		strings.HasPrefix(uri, "/favicon.ico") ||
		uri == "/" {
		return true
	}

	// Public GET endpoints for product information (mobile users and general access)
	if method == "GET" {
		// Product endpoints - public access to read data
		if strings.HasPrefix(uri, "/public/products") && !strings.Contains(uri, "/admin") {
			return true
		}
		if strings.HasPrefix(uri, "/public/categories") && !strings.Contains(uri, "/admin") {
			return true
		}
		if strings.HasPrefix(uri, "/public/tags") && !strings.Contains(uri, "/admin") {
			return true
		}
		if strings.HasPrefix(uri, "/public/feedback") && !strings.Contains(uri, "/admin") {
			return true
		}

		// Public access to product images and media downloads
		if strings.HasPrefix(uri, "/public/products/") && strings.Contains(uri, "/images") {
			return true
		}
		if strings.HasPrefix(uri, "/media/public-url/") || strings.HasPrefix(uri, "/media/download-url/") {
			return true
		}
	}

	return false
}

// isValidUserType validates user type according to BRD requirements
func (h *ForwardAuthHandler) isValidUserType(userType string) bool {
	validTypes := []string{"mobile_user", "superadmin", "admin", "partner"}
	for _, validType := range validTypes {
		if userType == validType {
			return true
		}
	}
	return false
}

// isAuthorized checks if the user has permission to access the requested resource
// Based on AUTH_SYSTEM_BRD.md requirements and path-based access control
func (h *ForwardAuthHandler) isAuthorized(claims *output.Claims, uri, method string) bool {
	// Health, version, and documentation endpoints are always public
	if strings.HasPrefix(uri, "/health") ||
		strings.HasPrefix(uri, "/swagger") ||
		strings.HasPrefix(uri, "/docs") ||
		strings.HasPrefix(uri, "/version") ||
		strings.HasPrefix(uri, "/info") ||
		strings.HasPrefix(uri, "/favicon.ico") ||
		uri == "/" {
		return true
	}

	// All /public/* endpoints are accessible without authentication
	if strings.HasPrefix(uri, "/public/") {
		return true
	}

	// Path-based access control for authenticated endpoints
	return h.checkPathBasedAccess(uri, method, claims)
}

// checkPathBasedAccess validates access based on path patterns and JWT claims
func (h *ForwardAuthHandler) checkPathBasedAccess(uri, method string, claims *output.Claims) bool {
	// /admin/* endpoints - only superadmin and admin access
	if strings.HasPrefix(uri, "/admin/") {
		return claims.UserType == "superadmin" || claims.UserType == "admin"
	}

	// /user/:userid/* endpoints - user can only access their own resources
	if strings.HasPrefix(uri, "/user/") {
		// Extract user ID from path: /user/123/feedback -> userID should be 123
		pathParts := strings.Split(strings.TrimPrefix(uri, "/user/"), "/")
		if len(pathParts) > 0 {
			requestedUserID := pathParts[0]
			currentUserID := fmt.Sprintf("%d", claims.UserID)

			// User can only access their own resources
			if requestedUserID == currentUserID {
				return claims.UserType == "mobile_user" || claims.UserType == "admin" || claims.UserType == "superadmin"
			}

			// Admins and superadmins can access any user's resources
			return claims.UserType == "admin" || claims.UserType == "superadmin"
		}
		return false
	}

	// /partner/:partner-id/* endpoints - partner can only access their own resources
	if strings.HasPrefix(uri, "/partner/") {
		// Extract partner ID from path: /partner/abc123/products -> partnerID should be abc123
		pathParts := strings.Split(strings.TrimPrefix(uri, "/partner/"), "/")
		if len(pathParts) > 0 {
			// requestedPartnerID := pathParts[0] - Will be used when partner claims are implemented

			// For partner access, we need to check if user is a partner and owns the resource
			// For now, only admin and superadmin can access partner endpoints until partner system is implemented
			if claims.UserType == "partner" {
				// TODO: Add partner ID validation when partner claims are added to JWT
				// For now, allow any partner to access any partner endpoint
				return true
			}

			// Admins and superadmins can access any partner's resources
			return claims.UserType == "admin" || claims.UserType == "superadmin"
		}
		return false
	}

	// /partners/:partner_id/* endpoints - partner can only access their own resources (plural form)
	if strings.HasPrefix(uri, "/partners/") {
		// Extract partner ID from path: /partners/123/orders -> partnerID should be 123
		pathParts := strings.Split(strings.TrimPrefix(uri, "/partners/"), "/")
		if len(pathParts) > 0 {
			requestedPartnerID := pathParts[0]
			currentUserID := fmt.Sprintf("%d", claims.UserID)

			// Enhanced logging for debugging partner access
			fmt.Printf("RBAC Partner Check - URI: %s, UserType: %s, UserID: %d, RequestedPartnerID: %s, CurrentUserID: %s\n",
				uri, claims.UserType, claims.UserID, requestedPartnerID, currentUserID)

			// Partners can only access their own resources (partner_id == user_id for partners)
			if claims.UserType == "partner" && requestedPartnerID == currentUserID {
				fmt.Printf("RBAC: Access GRANTED - Partner ID matches user ID\n")
				return true
			}

			// Admins and superadmins can access any partner's resources
			if claims.UserType == "admin" || claims.UserType == "superadmin" {
				fmt.Printf("RBAC: Access GRANTED - Admin/SuperAdmin access\n")
				return true
			}

			fmt.Printf("RBAC: Access DENIED - Partner ID mismatch or insufficient privileges\n")
			return false
		}
		return false
	}

	// /vendors/:vendor_id/* endpoints - vendor admin can only access their assigned vendor's resources
	if strings.HasPrefix(uri, "/vendors/") {
		// Extract vendor ID from path: /vendors/VND123/orders -> vendorID should be VND123
		pathParts := strings.Split(strings.TrimPrefix(uri, "/vendors/"), "/")
		if len(pathParts) > 0 {
			requestedVendorID := pathParts[0]

			// Enhanced logging for debugging vendor access
			fmt.Printf("RBAC Vendor Check - URI: %s, UserType: %s, UserID: %d, RequestedVendorID: %s\n",
				uri, claims.UserType, claims.UserID, requestedVendorID)

			// Vendor admins can only access their assigned vendor's resources
			// Note: Vendor admin validation will be handled by the vendor service
			// For now, we allow admin type users to access vendor endpoints
			if claims.UserType == "admin" || claims.UserType == "superadmin" {
				fmt.Printf("RBAC: Access GRANTED - Admin/SuperAdmin access to vendor endpoint\n")
				return true
			}

			fmt.Printf("RBAC: Access DENIED - Only admin/superadmin can access vendor endpoints\n")
			return false
		}
		return false
	}

	// /customers/:customer-id/* endpoints - customer can only access their own resources
	if strings.HasPrefix(uri, "/customers/") {
		// Extract customer ID from path: /customers/123/orders -> customerID should be 123
		pathParts := strings.Split(strings.TrimPrefix(uri, "/customers/"), "/")
		if len(pathParts) > 0 {
			requestedCustomerID := pathParts[0]
			currentUserID := fmt.Sprintf("%d", claims.UserID)

			// Enhanced logging for debugging customer access
			fmt.Printf("RBAC Customer Check - URI: %s, UserType: %s, UserID: %d, RequestedCustomerID: %s, CurrentUserID: %s\n",
				uri, claims.UserType, claims.UserID, requestedCustomerID, currentUserID)

			// Mobile users can only access their own customer resources
			// Accept either internal numeric ID or external IDs (Firebase/Google/Apple) present in claims
			if requestedCustomerID == currentUserID ||
				(claims.FirebaseUID != "" && requestedCustomerID == claims.FirebaseUID) ||
				(claims.GoogleID != "" && requestedCustomerID == claims.GoogleID) ||
				(claims.AppleID != "" && requestedCustomerID == claims.AppleID) {
				fmt.Printf("RBAC: Access GRANTED - Customer identifier matches claims (numeric or external)\n")
				return claims.UserType == "mobile_user" || claims.UserType == "admin" || claims.UserType == "superadmin"
			}

			// Admins and superadmins can access any customer's resources
			if claims.UserType == "admin" || claims.UserType == "superadmin" {
				fmt.Printf("RBAC: Access GRANTED - Admin/Superadmin access\n")
				return true
			}

			fmt.Printf("RBAC: Access DENIED - Customer ID %s does not match numeric ID %s nor external IDs\n", requestedCustomerID, currentUserID)
			return false
		}
		return false
	}

	// Default deny for any other patterns
	return false
}

// ProductAuth specifically handles product service authentication
// @Summary Product service forward authentication
// @Description Validates JWT tokens specifically for product service endpoints
// @Tags Forward Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param X-Forwarded-Uri header string true "Original request URI"
// @Param X-Forwarded-Method header string true "Original request method"
// @Success 200 {object} map[string]interface{} "Authentication successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /forward-auth/product [get]
func (h *ForwardAuthHandler) ProductAuth(c *fiber.Ctx) error {
	return h.ForwardAuth(c)
}

// CustomerAuth specifically handles customer service authentication
// @Summary Customer service forward authentication
// @Description Validates JWT tokens specifically for customer service endpoints with enhanced RBAC logging
// @Tags Forward Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param X-Forwarded-Uri header string true "Original request URI"
// @Param X-Forwarded-Method header string true "Original request method"
// @Success 200 {object} map[string]interface{} "Authentication successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /forward-auth/customer [get]
func (h *ForwardAuthHandler) CustomerAuth(c *fiber.Ctx) error {
	// Enhanced logging for customer service RBAC debugging
	originalURI := c.Get("X-Forwarded-Uri")
	originalMethod := c.Get("X-Forwarded-Method")

	fmt.Printf("CustomerAuth: Processing request %s %s\n", originalMethod, originalURI)

	result := h.ForwardAuth(c)

	// Log the result for debugging
	if result != nil {
		if fiberErr, ok := result.(*fiber.Error); ok {
			fmt.Printf("CustomerAuth: Request DENIED - Status: %d\n", fiberErr.Code)
		}
	} else {
		fmt.Printf("CustomerAuth: Request GRANTED\n")
	}

	return result
}
