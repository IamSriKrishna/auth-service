package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ProductService interface {
	// Basic CRUD Operations
	CreateProduct(input *input.CreateProductInput, createdBy string) (*output.ProductOutput, error)
	GetProduct(id string, createdBy string) (*output.ProductOutput, error)
	GetAllProducts(limit, offset int, createdBy string) (*output.ProductListOutput, error)
	UpdateProduct(id string, input *input.UpdateProductInput, createdBy string) (*output.ProductOutput, error)
	DeleteProduct(id string, createdBy string) error

	// Variant Operations
	GetProductVariants(productID string) ([]models.ProductVariant, error)

	// Inventory Operations
	EnableInventoryTracking(productID string) error
	IsInventoryTrackingEnabled(productID string) (bool, error)
	GetProductWithInventoryStatus(productID string) (map[string]interface{}, error)
}

type productService struct {
	repo          repo.ProductRepository
	vendorRepo    repo.VendorRepository
	inventoryRepo repo.InventoryBalanceRepository
	userRepo      repo.UserRepository
	companyRepo   repo.CompanyRepository
}

func NewProductService(
	productRepo repo.ProductRepository,
	vendorRepo repo.VendorRepository,
	inventoryRepo repo.InventoryBalanceRepository,
	userRepo repo.UserRepository,
	companyRepo repo.CompanyRepository,
) ProductService {
	return &productService{
		repo:          productRepo,
		vendorRepo:    vendorRepo,
		inventoryRepo: inventoryRepo,
		userRepo:      userRepo,
		companyRepo:   companyRepo,
	}
}

func (s *productService) CreateProduct(input *input.CreateProductInput, createdBy string) (*output.ProductOutput, error) {
	// For resources, skip variant and SKU validation
	if !input.IsResource {
		// Validate variant attributes first
		if err := input.ValidateVariantAttributes(); err != nil {
			return nil, err
		}

		// Validate SKU uniqueness
		if err := input.ValidateSKUUniqueness(); err != nil {
			return nil, err
		}
	}

	id := fmt.Sprintf("prod_%s", uuid.New().String()[:8])

	idPtr := id
	product := &models.Product{
		ID:                  id,
		Name:                input.Name,
		IsResource:          input.IsResource,
		ResourceName:        input.ResourceName,
		ResourceUnit:        input.ResourceUnit,
		ResourceCostPerUnit: input.ResourceCostPerUnit,
		CreatedBy:           createdBy,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Only add product details, sales info, purchase info, inventory, and return policy for non-resource products
	if !input.IsResource {
		productDetails := buildProductDetails(id, input)
		product.ProductDetails = productDetails

		salesInfo := models.SalesInfo{}
		if input.SalesInfo != nil {
			salesInfo = models.SalesInfo{
				ProductID:    &idPtr,
				Account:      input.SalesInfo.Account,
				SellingPrice: input.SalesInfo.SellingPrice,
				Currency:     input.SalesInfo.Currency,
				Description:  input.SalesInfo.Description,
			}
		}
		product.SalesInfo = salesInfo

		purchaseInfo := models.PurchaseInfo{}
		if input.PurchaseInfo != nil {
			purchaseInfo = models.PurchaseInfo{
				ProductID:   &idPtr,
				Account:     input.PurchaseInfo.Account,
				CostPrice:   input.PurchaseInfo.CostPrice,
				Currency:    input.PurchaseInfo.Currency,
				Description: input.PurchaseInfo.Description,
			}
		}
		product.PurchaseInfo = purchaseInfo

		inventory := models.Inventory{
			ProductID:      &idPtr,
			TrackInventory: false,
		}
		if input.Inventory != nil {
			inventory.TrackInventory = input.Inventory.TrackInventory
			inventory.InventoryAccount = input.Inventory.InventoryAccount
			inventory.InventoryValuationMethod = input.Inventory.InventoryValuationMethod
			inventory.ReorderPoint = input.Inventory.ReorderPoint
		}
		product.Inventory = inventory

		returnPolicy := models.ReturnPolicy{
			ProductID:  &idPtr,
			Returnable: false,
		}
		if input.ReturnPolicy != nil {
			returnPolicy.Returnable = input.ReturnPolicy.Returnable
		}
		product.ReturnPolicy = returnPolicy
	}

	// Fetch user details if createdBy is provided
	if createdBy != "" {
		var userID uint
		_, err := fmt.Sscanf(createdBy, "%d", &userID)
		if err == nil {
			user, err := s.userRepo.GetByID(userID)
			if err == nil && user != nil {
				// Set user name from email or username
				if user.Email != nil {
					product.CreatedByUserName = *user.Email
				} else if user.Username != nil {
					product.CreatedByUserName = *user.Username
				}

				// Set company details if user has a company
				if user.CompanyID != nil {
					product.CreatedByCompanyID = *user.CompanyID
					company, err := s.companyRepo.FindByID(*user.CompanyID)
					if err == nil && company != nil {
						product.CreatedByCompanyName = company.CompanyName
					}
				}
			}
		}
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	createdProduct, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToProductOutput(createdProduct)
}

func buildProductDetails(productID string, input *input.CreateProductInput) models.ProductDetails {
	productDetails := models.ProductDetails{
		ProductID:   productID,
		Unit:        input.ProductDetails.Unit,
		BaseSKU:     input.ProductDetails.BaseSKU,
		UPC:         input.ProductDetails.UPC,
		EAN:         input.ProductDetails.EAN,
		MPN:         input.ProductDetails.MPN,
		ISBN:        input.ProductDetails.ISBN,
		Description: input.ProductDetails.Description,
	}

	if len(input.ProductDetails.AttributeDefinitions) > 0 {
		productDetails.AttributeDefinitions = make(models.ProductAttributeDefinitions, len(input.ProductDetails.AttributeDefinitions))
		for i, attr := range input.ProductDetails.AttributeDefinitions {
			productDetails.AttributeDefinitions[i] = models.ProductAttributeDefinition{
				Key:     attr.Key,
				Options: attr.Options,
			}
		}
	}

	// Handle variants: either use provided variants or auto-create a default variant
	if len(input.ProductDetails.Variants) > 0 {
		productDetails.ProductVariants = make([]models.ProductVariant, len(input.ProductDetails.Variants))
		for i, v := range input.ProductDetails.Variants {
			attributes := make([]models.ProductVariantAttribute, 0, len(v.AttributeMap))
			for key, value := range v.AttributeMap {
				attributes = append(attributes, models.ProductVariantAttribute{
					Key:   key,
					Value: value,
				})
			}

			productDetails.ProductVariants[i] = models.ProductVariant{
				SKU:           v.SKU,
				VariantName:   v.VariantName,
				Attributes:    attributes,
				SellingPrice:  v.SellingPrice,
				CostPrice:     v.CostPrice,
				StockQuantity: v.StockQuantity,
				IsActive:      v.IsActive,
			}
		}
	} else {
		// No variants provided - auto-create a default variant from base product info
		sku := input.ProductDetails.BaseSKU
		if sku == "" {
			// If no base SKU provided, use product name as base
			sku = input.Name
		}

		defaultVariant := models.ProductVariant{
			SKU:           sku,
			VariantName:   input.Name,
			Attributes:    []models.ProductVariantAttribute{}, // No attributes for single variant
			SellingPrice:  input.SalesInfo.SellingPrice,
			CostPrice:     0.0, // Will be set from PurchaseInfo if available
			StockQuantity: 0.0,
			IsActive:      true,
		}

		// Set cost price from purchase info if available
		if input.PurchaseInfo != nil {
			defaultVariant.CostPrice = input.PurchaseInfo.CostPrice
		}

		productDetails.ProductVariants = []models.ProductVariant{defaultVariant}
	}

	return productDetails
}

func (s *productService) GetProduct(id string, createdBy string) (*output.ProductOutput, error) {
	// Parse user ID from createdBy
	var userID uint
	_, err := fmt.Sscanf(createdBy, "%d", &userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	// Get user to fetch company
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get product
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if product belongs to user's company
	if product.CreatedByCompanyID != *user.CompanyID {
		return nil, fmt.Errorf("unauthorized: product does not belong to your company")
	}

	return output.ToProductOutput(product)
}

func (s *productService) GetAllProducts(limit, offset int, createdBy string) (*output.ProductListOutput, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Parse user ID from createdBy
	var userID uint
	_, err := fmt.Sscanf(createdBy, "%d", &userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	// Get user to fetch company
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get products only from user's company
	products, total, err := s.repo.FindByCreatedByAndCompany(createdBy, *user.CompanyID, limit, offset)
	if err != nil {
		return nil, err
	}

	return output.ToProductListOutput(products, total)
}

func (s *productService) UpdateProduct(id string, input *input.UpdateProductInput, createdBy string) (*output.ProductOutput, error) {
	// Parse user ID from createdBy
	var userID uint
	_, err := fmt.Sscanf(createdBy, "%d", &userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	// Get user to fetch company
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user created the product
	if product.CreatedBy != createdBy {
		return nil, fmt.Errorf("unauthorized: you can only update products you created")
	}

	// Check if product belongs to user's company
	if product.CreatedByCompanyID != *user.CompanyID {
		return nil, fmt.Errorf("unauthorized: product does not belong to your company")
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	// Handle resource fields
	if input.IsResource != nil {
		product.IsResource = *input.IsResource
	}
	if input.ResourceName != nil {
		product.ResourceName = *input.ResourceName
	}
	if input.ResourceUnit != nil {
		product.ResourceUnit = *input.ResourceUnit
	}
	if input.ResourceCostPerUnit != nil {
		product.ResourceCostPerUnit = *input.ResourceCostPerUnit
	}

	if input.SalesInfo != nil {
		if input.SalesInfo.Account != "" {
			product.SalesInfo.Account = input.SalesInfo.Account
		}
		if input.SalesInfo.SellingPrice > 0 {
			product.SalesInfo.SellingPrice = input.SalesInfo.SellingPrice
		}
		if input.SalesInfo.Currency != "" {
			product.SalesInfo.Currency = input.SalesInfo.Currency
		}
		if input.SalesInfo.Description != "" {
			product.SalesInfo.Description = input.SalesInfo.Description
		}
	}

	if input.PurchaseInfo != nil {
		if input.PurchaseInfo.Account != "" {
			product.PurchaseInfo.Account = input.PurchaseInfo.Account
		}
		if input.PurchaseInfo.CostPrice > 0 {
			product.PurchaseInfo.CostPrice = input.PurchaseInfo.CostPrice
		}
		if input.PurchaseInfo.Currency != "" {
			product.PurchaseInfo.Currency = input.PurchaseInfo.Currency
		}
		if input.PurchaseInfo.Description != "" {
			product.PurchaseInfo.Description = input.PurchaseInfo.Description
		}
	}

	if input.Inventory != nil {
		product.Inventory.TrackInventory = input.Inventory.TrackInventory
		if input.Inventory.InventoryAccount != "" {
			product.Inventory.InventoryAccount = input.Inventory.InventoryAccount
		}
		if input.Inventory.InventoryValuationMethod != "" {
			product.Inventory.InventoryValuationMethod = input.Inventory.InventoryValuationMethod
		}
		if input.Inventory.ReorderPoint >= 0 {
			product.Inventory.ReorderPoint = input.Inventory.ReorderPoint
		}
	}

	if input.ReturnPolicy != nil {
		product.ReturnPolicy.Returnable = input.ReturnPolicy.Returnable
	}

	if input.ProductDetails != nil {
		if input.ProductDetails.Unit != "" {
			product.ProductDetails.Unit = input.ProductDetails.Unit
		}
		if input.ProductDetails.BaseSKU != "" {
			product.ProductDetails.BaseSKU = input.ProductDetails.BaseSKU
		}
		if input.ProductDetails.Description != "" {
			product.ProductDetails.Description = input.ProductDetails.Description
		}

		if len(input.ProductDetails.Variants) > 0 {
			product.ProductDetails.ProductVariants = make([]models.ProductVariant, len(input.ProductDetails.Variants))
			for i, v := range input.ProductDetails.Variants {
				attributes := make([]models.ProductVariantAttribute, 0, len(v.AttributeMap))
				for key, value := range v.AttributeMap {
					attributes = append(attributes, models.ProductVariantAttribute{
						Key:   key,
						Value: value,
					})
				}

				product.ProductDetails.ProductVariants[i] = models.ProductVariant{
					SKU:           v.SKU,
					VariantName:   v.VariantName,
					Attributes:    attributes,
					SellingPrice:  v.SellingPrice,
					CostPrice:     v.CostPrice,
					StockQuantity: v.StockQuantity,
					IsActive:      v.IsActive,
				}
			}
		}

		if len(input.ProductDetails.AttributeDefinitions) > 0 {
			product.ProductDetails.AttributeDefinitions = make(models.ProductAttributeDefinitions, len(input.ProductDetails.AttributeDefinitions))
			for i, attr := range input.ProductDetails.AttributeDefinitions {
				product.ProductDetails.AttributeDefinitions[i] = models.ProductAttributeDefinition{
					Key:     attr.Key,
					Options: attr.Options,
				}
			}
		}
	}

	product.UpdatedAt = time.Now()

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	updatedProduct, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return output.ToProductOutput(updatedProduct)
}

func (s *productService) DeleteProduct(id string, createdBy string) error {
	// Parse user ID from createdBy
	var userID uint
	_, err := fmt.Sscanf(createdBy, "%d", &userID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}

	// Get user to fetch company
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	product, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	// Check if user created the product
	if product.CreatedBy != createdBy {
		return errors.New("unauthorized: you can only delete products you created")
	}

	// Check if product belongs to user's company
	if product.CreatedByCompanyID != *user.CompanyID {
		return errors.New("unauthorized: product does not belong to your company")
	}

	return s.repo.DeleteByCreatedBy(id, createdBy)
}

// GetProductVariants retrieves all variants for a product
func (s *productService) GetProductVariants(productID string) ([]models.ProductVariant, error) {
	return s.repo.GetProductVariantsByProductID(productID)
}

// EnableInventoryTracking enables inventory tracking for a product
func (s *productService) EnableInventoryTracking(productID string) error {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	product.Inventory.TrackInventory = true
	product.UpdatedAt = time.Now()

	return s.repo.Update(product)
}

// IsInventoryTrackingEnabled checks if inventory tracking is enabled for a product
func (s *productService) IsInventoryTrackingEnabled(productID string) (bool, error) {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return false, fmt.Errorf("product not found: %w", err)
	}

	return product.Inventory.TrackInventory, nil
}

// GetProductWithInventoryStatus returns product details along with current inventory balance
func (s *productService) GetProductWithInventoryStatus(productID string) (map[string]interface{}, error) {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	variants := make([]map[string]interface{}, len(product.ProductDetails.ProductVariants))
	for i, v := range product.ProductDetails.ProductVariants {
		variants[i] = map[string]interface{}{
			"sku":            v.SKU,
			"variant_name":   v.VariantName,
			"selling_price":  v.SellingPrice,
			"cost_price":     v.CostPrice,
			"stock_quantity": v.StockQuantity,
			"reorder_level":  v.ReorderLevel,
			"is_active":      v.IsActive,
		}
	}

	result := map[string]interface{}{
		"product_id":         productID,
		"name":               product.Name,
		"inventory_tracking": product.Inventory.TrackInventory,
		"reorder_point":      product.Inventory.ReorderPoint,
		"purchase_price":     product.PurchaseInfo.CostPrice,
		"selling_price":      product.SalesInfo.SellingPrice,
		"variants":           variants,
	}

	return result, nil
}
