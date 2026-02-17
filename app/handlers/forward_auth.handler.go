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

type ForwardAuthHandler struct {
	cache      map[string]*CacheEntry
	cacheMutex sync.RWMutex
}

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

func NewForwardAuthHandler() *ForwardAuthHandler {
	handler := &ForwardAuthHandler{
		cache: make(map[string]*CacheEntry),
	}

	go handler.startCacheCleanup()

	return handler
}

func (h *ForwardAuthHandler) startCacheCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		h.cleanupExpiredEntries()
	}
}

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

func (h *ForwardAuthHandler) ForwardAuth(c *fiber.Ctx) error {
	originalURI := c.Get("X-Forwarded-Uri")
	originalMethod := c.Get("X-Forwarded-Method")

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

	fmt.Printf("Forward Auth Request: %s %s from %s\n", originalMethod, originalURI, c.IP())
	fmt.Printf("Debug Headers - X-Forwarded-Uri: '%s', X-Forwarded-Method: '%s'\n", c.Get("X-Forwarded-Uri"), c.Get("X-Forwarded-Method"))
	fmt.Printf("All request headers: %+v\n", c.GetReqHeaders())

	if h.isPublicEndpoint(originalURI, originalMethod) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"authenticated": true,
			"public":        true,
			"message":       "Public endpoint access granted",
		})
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Authorization header is required",
			"code":    "AUTH_HEADER_MISSING",
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid authorization header format",
			"code":    "INVALID_AUTH_FORMAT",
		})
	}

	token := tokenParts[1]

	var claims *output.Claims

	if cachedEntry, found := h.getCachedAuth(token); found {
		fmt.Printf("CACHE HIT: Found cached auth for user %d (%s) - but still checking RBAC for URI: %s\n",
			cachedEntry.UserID, cachedEntry.UserType, originalURI)

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

	if !h.isValidUserType(claims.UserType) {
		fmt.Printf("RBAC: Invalid user type '%s' for URI: %s\n", claims.UserType, originalURI)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user type",
			"code":    "INVALID_USER_TYPE",
		})
	}

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
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		h.cacheMutex.Unlock()
	}

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

func (h *ForwardAuthHandler) isPublicEndpoint(uri, method string) bool {
	if strings.HasPrefix(uri, "/health") ||
		strings.HasPrefix(uri, "/swagger") ||
		strings.HasPrefix(uri, "/docs") ||
		strings.HasPrefix(uri, "/version") ||
		strings.HasPrefix(uri, "/info") ||
		strings.HasPrefix(uri, "/favicon.ico") ||
		uri == "/" {
		return true
	}

	if method == "GET" {
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

		if strings.HasPrefix(uri, "/public/products/") && strings.Contains(uri, "/images") {
			return true
		}
		if strings.HasPrefix(uri, "/media/public-url/") || strings.HasPrefix(uri, "/media/download-url/") {
			return true
		}
	}

	return false
}

func (h *ForwardAuthHandler) isValidUserType(userType string) bool {
	validTypes := []string{"mobile_user", "superadmin", "admin", "partner"}
	for _, validType := range validTypes {
		if userType == validType {
			return true
		}
	}
	return false
}

func (h *ForwardAuthHandler) isAuthorized(claims *output.Claims, uri, method string) bool {
	if strings.HasPrefix(uri, "/health") ||
		strings.HasPrefix(uri, "/swagger") ||
		strings.HasPrefix(uri, "/docs") ||
		strings.HasPrefix(uri, "/version") ||
		strings.HasPrefix(uri, "/info") ||
		strings.HasPrefix(uri, "/favicon.ico") ||
		uri == "/" {
		return true
	}

	if strings.HasPrefix(uri, "/public/") {
		return true
	}

	return h.checkPathBasedAccess(uri, method, claims)
}

func (h *ForwardAuthHandler) checkPathBasedAccess(uri, method string, claims *output.Claims) bool {
	if strings.HasPrefix(uri, "/admin/") {
		return claims.UserType == "superadmin" || claims.UserType == "admin"
	}

	if strings.HasPrefix(uri, "/user/") {
		pathParts := strings.Split(strings.TrimPrefix(uri, "/user/"), "/")
		if len(pathParts) > 0 {
			requestedUserID := pathParts[0]
			currentUserID := fmt.Sprintf("%d", claims.UserID)

			if requestedUserID == currentUserID {
				return claims.UserType == "mobile_user" || claims.UserType == "admin" || claims.UserType == "superadmin"
			}

			return claims.UserType == "admin" || claims.UserType == "superadmin"
		}
		return false
	}

	if strings.HasPrefix(uri, "/partner/") {
		pathParts := strings.Split(strings.TrimPrefix(uri, "/partner/"), "/")
		if len(pathParts) > 0 {

			if claims.UserType == "partner" {
				return true
			}

			return claims.UserType == "admin" || claims.UserType == "superadmin"
		}
		return false
	}

	if strings.HasPrefix(uri, "/partners/") {
		pathParts := strings.Split(strings.TrimPrefix(uri, "/partners/"), "/")
		if len(pathParts) > 0 {
			requestedPartnerID := pathParts[0]
			currentUserID := fmt.Sprintf("%d", claims.UserID)

			fmt.Printf("RBAC Partner Check - URI: %s, UserType: %s, UserID: %d, RequestedPartnerID: %s, CurrentUserID: %s\n",
				uri, claims.UserType, claims.UserID, requestedPartnerID, currentUserID)

			if claims.UserType == "partner" && requestedPartnerID == currentUserID {
				fmt.Printf("RBAC: Access GRANTED - Partner ID matches user ID\n")
				return true
			}

			if claims.UserType == "admin" || claims.UserType == "superadmin" {
				fmt.Printf("RBAC: Access GRANTED - Admin/SuperAdmin access\n")
				return true
			}

			fmt.Printf("RBAC: Access DENIED - Partner ID mismatch or insufficient privileges\n")
			return false
		}
		return false
	}

	if strings.HasPrefix(uri, "/vendors/") {
		pathParts := strings.Split(strings.TrimPrefix(uri, "/vendors/"), "/")
		if len(pathParts) > 0 {
			requestedVendorID := pathParts[0]

			fmt.Printf("RBAC Vendor Check - URI: %s, UserType: %s, UserID: %d, RequestedVendorID: %s\n",
				uri, claims.UserType, claims.UserID, requestedVendorID)

			if claims.UserType == "admin" || claims.UserType == "superadmin" {
				fmt.Printf("RBAC: Access GRANTED - Admin/SuperAdmin access to vendor endpoint\n")
				return true
			}

			fmt.Printf("RBAC: Access DENIED - Only admin/superadmin can access vendor endpoints\n")
			return false
		}
		return false
	}

	if strings.HasPrefix(uri, "/customers/") {
		pathParts := strings.Split(strings.TrimPrefix(uri, "/customers/"), "/")
		if len(pathParts) > 0 {
			requestedCustomerID := pathParts[0]
			currentUserID := fmt.Sprintf("%d", claims.UserID)

			fmt.Printf("RBAC Customer Check - URI: %s, UserType: %s, UserID: %d, RequestedCustomerID: %s, CurrentUserID: %s\n",
				uri, claims.UserType, claims.UserID, requestedCustomerID, currentUserID)

			if requestedCustomerID == currentUserID ||
				(claims.FirebaseUID != "" && requestedCustomerID == claims.FirebaseUID) ||
				(claims.GoogleID != "" && requestedCustomerID == claims.GoogleID) ||
				(claims.AppleID != "" && requestedCustomerID == claims.AppleID) {
				fmt.Printf("RBAC: Access GRANTED - Customer identifier matches claims (numeric or external)\n")
				return claims.UserType == "mobile_user" || claims.UserType == "admin" || claims.UserType == "superadmin"
			}

			if claims.UserType == "admin" || claims.UserType == "superadmin" {
				fmt.Printf("RBAC: Access GRANTED - Admin/Superadmin access\n")
				return true
			}

			fmt.Printf("RBAC: Access DENIED - Customer ID %s does not match numeric ID %s nor external IDs\n", requestedCustomerID, currentUserID)
			return false
		}
		return false
	}

	return false
}

func (h *ForwardAuthHandler) ProductAuth(c *fiber.Ctx) error {
	return h.ForwardAuth(c)
}

func (h *ForwardAuthHandler) CustomerAuth(c *fiber.Ctx) error {
	originalURI := c.Get("X-Forwarded-Uri")
	originalMethod := c.Get("X-Forwarded-Method")

	fmt.Printf("CustomerAuth: Processing request %s %s\n", originalMethod, originalURI)

	result := h.ForwardAuth(c)

	if result != nil {
		if fiberErr, ok := result.(*fiber.Error); ok {
			fmt.Printf("CustomerAuth: Request DENIED - Status: %d\n", fiberErr.Code)
		}
	} else {
		fmt.Printf("CustomerAuth: Request GRANTED\n")
	}

	return result
}
