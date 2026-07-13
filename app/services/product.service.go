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
	// Existing methods retained for compatibility.
	CreateProduct(input *input.CreateProductInput, createdBy string) (*output.ProductOutput, error)
	GetProduct(id string, createdBy string) (*output.ProductOutput, error)
	GetAllProducts(limit, offset int, createdBy string) (*output.ProductListOutput, error)
	UpdateProduct(id string, input *input.UpdateProductInput, createdBy string) (*output.ProductOutput, error)
	DeleteProduct(id string, createdBy string) error
	GetProductVariants(productID string) ([]models.ProductVariant, error)

	EnableInventoryTracking(productID string) error
	IsInventoryTrackingEnabled(productID string) (bool, error)
	GetProductWithInventoryStatus(productID string) (map[string]interface{}, error)

	// Company-scoped methods added.
	CreateProductForCompany(
		input *input.CreateProductInput,
		createdBy string,
		companyID uint,
	) (*output.ProductOutput, error)

	GetProductByCompany(
		id string,
		companyID uint,
	) (*output.ProductOutput, error)

	GetAllProductsByCompany(
		limit int,
		offset int,
		companyID uint,
	) (*output.ProductListOutput, error)

	UpdateProductForCompany(
		id string,
		input *input.UpdateProductInput,
		companyID uint,
	) (*output.ProductOutput, error)

	DeleteProductForCompany(
		id string,
		companyID uint,
	) error

	GetProductVariantsByCompany(
		productID string,
		companyID uint,
	) ([]models.ProductVariant, error)
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

func (s *productService) CreateProduct(
	req *input.CreateProductInput,
	createdBy string,
) (*output.ProductOutput, error) {
	if createdBy == "" {
		return nil, errors.New("created_by is required")
	}

	var userID uint
	if _, err := fmt.Sscanf(createdBy, "%d", &userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if user.CompanyID == nil || *user.CompanyID == 0 {
		return nil, errors.New("user is not assigned to a company")
	}

	return s.CreateProductForCompany(
		req,
		createdBy,
		*user.CompanyID,
	)
}

func (s *productService) CreateProductForCompany(
	req *input.CreateProductInput,
	createdBy string,
	companyID uint,
) (*output.ProductOutput, error) {
	if req == nil {
		return nil, errors.New("input cannot be nil")
	}

	if companyID == 0 {
		return nil, errors.New("invalid company")
	}

	if !req.IsResource && !req.IsRaw {
		if req.ProductDetails == nil {
			return nil, fmt.Errorf("product_details is required for regular products")
		}

		if err := req.ValidateVariantAttributes(); err != nil {
			return nil, err
		}

		if err := req.ValidateSKUUniqueness(); err != nil {
			return nil, err
		}
	}

	company, err := s.companyRepo.FindByID(companyID)
	if err != nil || company == nil {
		return nil, errors.New("company not found")
	}

	createdByUserName := ""
	if createdBy != "" {
		var userID uint
		if _, err := fmt.Sscanf(createdBy, "%d", &userID); err == nil {
			user, err := s.userRepo.GetByIDAndCompanyID(userID, companyID)
			if err != nil || user == nil {
				return nil, errors.New("user does not belong to the company")
			}

			if user.Email != nil {
				createdByUserName = *user.Email
			} else if user.Username != nil {
				createdByUserName = *user.Username
			}
		}
	}

	id := fmt.Sprintf("prod_%s", uuid.New().String()[:8])
	idPointer := id

	product := &models.Product{
		ID:   id,
		Name: req.Name,

		BaseUnit:         req.BaseUnit,
		PurchaseUnit:     req.PurchaseUnit,
		ConversionFactor: req.ConversionFactor,

		IsResource:          req.IsResource,
		ResourceName:        req.ResourceName,
		ResourceUnit:        req.ResourceUnit,
		ResourceCostPerUnit: req.ResourceCostPerUnit,

		IsRaw:               req.IsRaw,
		RawName:             req.RawName,
		RawSpecification:    req.RawSpecification,
		RawUnit:             req.RawUnit,
		RawCostPerUnit:      req.RawCostPerUnit,
		RequiredGramPerUnit: req.RequiredGramPerUnit,
		ConsumptionPerUnit:  req.ConsumptionPerUnit,

		CreatedBy:            createdBy,
		CreatedByUserName:    createdByUserName,
		CreatedByCompanyID:   companyID,
		CreatedByCompanyName: company.CompanyName,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if product.ConversionFactor <= 0 {
		product.ConversionFactor = 1
	}

	if !req.IsResource && !req.IsRaw && req.ProductDetails != nil {
		product.ProductDetails = buildProductDetails(id, req)

		if req.SalesInfo != nil {
			product.SalesInfo = models.SalesInfo{
				ProductID:    &idPointer,
				Account:      req.SalesInfo.Account,
				SellingPrice: req.SalesInfo.SellingPrice,
				Currency:     req.SalesInfo.Currency,
				Description:  req.SalesInfo.Description,
			}
		}

		if req.PurchaseInfo != nil {
			product.PurchaseInfo = models.PurchaseInfo{
				ProductID:   &idPointer,
				Account:     req.PurchaseInfo.Account,
				CostPrice:   req.PurchaseInfo.CostPrice,
				Currency:    req.PurchaseInfo.Currency,
				Description: req.PurchaseInfo.Description,
			}
		}

		product.Inventory = models.Inventory{
			ProductID:      &idPointer,
			TrackInventory: false,
		}

		if req.Inventory != nil {
			product.Inventory.TrackInventory = req.Inventory.TrackInventory
			product.Inventory.InventoryAccount = req.Inventory.InventoryAccount
			product.Inventory.InventoryValuationMethod = req.Inventory.InventoryValuationMethod
			product.Inventory.ReorderPoint = req.Inventory.ReorderPoint
		}

		product.ReturnPolicy = models.ReturnPolicy{
			ProductID:  &idPointer,
			Returnable: false,
		}

		if req.ReturnPolicy != nil {
			product.ReturnPolicy.Returnable = req.ReturnPolicy.Returnable
		}
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	createdProduct, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	return output.ToProductOutput(createdProduct)
}

func buildProductDetails(
	productID string,
	req *input.CreateProductInput,
) models.ProductDetails {
	if req.ProductDetails == nil {
		return models.ProductDetails{
			ProductID: productID,
		}
	}

	productDetails := models.ProductDetails{
		ProductID:   productID,
		Unit:        req.ProductDetails.Unit,
		BaseSKU:     req.ProductDetails.BaseSKU,
		UPC:         req.ProductDetails.UPC,
		EAN:         req.ProductDetails.EAN,
		MPN:         req.ProductDetails.MPN,
		ISBN:        req.ProductDetails.ISBN,
		Description: req.ProductDetails.Description,
	}

	if len(req.ProductDetails.AttributeDefinitions) > 0 {
		productDetails.AttributeDefinitions = make(
			models.ProductAttributeDefinitions,
			len(req.ProductDetails.AttributeDefinitions),
		)

		for index, attribute := range req.ProductDetails.AttributeDefinitions {
			productDetails.AttributeDefinitions[index] =
				models.ProductAttributeDefinition{
					Key:     attribute.Key,
					Options: attribute.Options,
				}
		}
	}

	if len(req.ProductDetails.Variants) > 0 {
		productDetails.ProductVariants = make(
			[]models.ProductVariant,
			len(req.ProductDetails.Variants),
		)

		for index, variant := range req.ProductDetails.Variants {
			attributes := make(
				[]models.ProductVariantAttribute,
				0,
				len(variant.AttributeMap),
			)

			for key, value := range variant.AttributeMap {
				attributes = append(
					attributes,
					models.ProductVariantAttribute{
						Key:   key,
						Value: value,
					},
				)
			}

			productDetails.ProductVariants[index] = models.ProductVariant{
				SKU:           variant.SKU,
				VariantName:   variant.VariantName,
				Attributes:    attributes,
				SellingPrice:  variant.SellingPrice,
				CostPrice:     variant.CostPrice,
				StockQuantity: variant.StockQuantity,
				IsActive:      variant.IsActive,
			}
		}
	} else {
		sku := req.ProductDetails.BaseSKU
		if sku == "" {
			sku = req.Name
		}

		costPrice := 0.0
		if req.PurchaseInfo != nil {
			costPrice = req.PurchaseInfo.CostPrice
		}

		sellingPrice := 0.0
		if req.SalesInfo != nil {
			sellingPrice = req.SalesInfo.SellingPrice
		}

		productDetails.ProductVariants = []models.ProductVariant{
			{
				SKU:           sku,
				VariantName:   req.Name,
				Attributes:    []models.ProductVariantAttribute{},
				SellingPrice:  sellingPrice,
				CostPrice:     costPrice,
				StockQuantity: 0,
				IsActive:      true,
			},
		}
	}

	return productDetails
}

func (s *productService) GetProduct(
	id string,
	createdBy string,
) (*output.ProductOutput, error) {
	var userID uint
	if _, err := fmt.Sscanf(createdBy, "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.CompanyID == nil {
		return nil, fmt.Errorf("user is not assigned to a company")
	}

	return s.GetProductByCompany(id, *user.CompanyID)
}

func (s *productService) GetProductByCompany(
	id string,
	companyID uint,
) (*output.ProductOutput, error) {
	product, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	return output.ToProductOutput(product)
}

func (s *productService) GetAllProducts(
	limit int,
	offset int,
	createdBy string,
) (*output.ProductListOutput, error) {
	var userID uint
	if _, err := fmt.Sscanf(createdBy, "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.CompanyID == nil {
		return nil, fmt.Errorf("user is not assigned to a company")
	}

	return s.GetAllProductsByCompany(
		limit,
		offset,
		*user.CompanyID,
	)
}

func (s *productService) GetAllProductsByCompany(
	limit int,
	offset int,
	companyID uint,
) (*output.ProductListOutput, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	products, total, err := s.repo.FindByCompany(
		companyID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return output.ToProductListOutput(products, total)
}

func (s *productService) UpdateProduct(
	id string,
	req *input.UpdateProductInput,
	createdBy string,
) (*output.ProductOutput, error) {
	var userID uint
	if _, err := fmt.Sscanf(createdBy, "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.CompanyID == nil {
		return nil, fmt.Errorf("user is not assigned to a company")
	}

	return s.UpdateProductForCompany(
		id,
		req,
		*user.CompanyID,
	)
}

func (s *productService) UpdateProductForCompany(
	id string,
	req *input.UpdateProductInput,
	companyID uint,
) (*output.ProductOutput, error) {
	if req == nil {
		return nil, errors.New("input cannot be nil")
	}

	product, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.IsResource != nil {
		product.IsResource = *req.IsResource
	}
	if req.ResourceName != nil {
		product.ResourceName = *req.ResourceName
	}
	if req.ResourceUnit != nil {
		product.ResourceUnit = *req.ResourceUnit
	}
	if req.ResourceCostPerUnit != nil {
		product.ResourceCostPerUnit = *req.ResourceCostPerUnit
	}

	if req.IsRaw != nil {
		product.IsRaw = *req.IsRaw
	}
	if req.RawName != nil {
		product.RawName = *req.RawName
	}
	if req.RawSpecification != nil {
		product.RawSpecification = *req.RawSpecification
	}
	if req.RawUnit != nil {
		product.RawUnit = *req.RawUnit
	}
	if req.RawCostPerUnit != nil {
		product.RawCostPerUnit = *req.RawCostPerUnit
	}
	if req.RequiredGramPerUnit != nil {
		product.RequiredGramPerUnit = *req.RequiredGramPerUnit
	}
	if req.ConsumptionPerUnit != nil {
		product.ConsumptionPerUnit = *req.ConsumptionPerUnit
	}

	if req.SalesInfo != nil {
		if req.SalesInfo.Account != "" {
			product.SalesInfo.Account = req.SalesInfo.Account
		}
		if req.SalesInfo.SellingPrice > 0 {
			product.SalesInfo.SellingPrice = req.SalesInfo.SellingPrice
		}
		if req.SalesInfo.Currency != "" {
			product.SalesInfo.Currency = req.SalesInfo.Currency
		}
		if req.SalesInfo.Description != "" {
			product.SalesInfo.Description = req.SalesInfo.Description
		}
	}

	if req.PurchaseInfo != nil {
		if req.PurchaseInfo.Account != "" {
			product.PurchaseInfo.Account = req.PurchaseInfo.Account
		}
		if req.PurchaseInfo.CostPrice > 0 {
			product.PurchaseInfo.CostPrice = req.PurchaseInfo.CostPrice
		}
		if req.PurchaseInfo.Currency != "" {
			product.PurchaseInfo.Currency = req.PurchaseInfo.Currency
		}
		if req.PurchaseInfo.Description != "" {
			product.PurchaseInfo.Description = req.PurchaseInfo.Description
		}
	}

	if req.Inventory != nil {
		product.Inventory.TrackInventory = req.Inventory.TrackInventory

		if req.Inventory.InventoryAccount != "" {
			product.Inventory.InventoryAccount =
				req.Inventory.InventoryAccount
		}

		if req.Inventory.InventoryValuationMethod != "" {
			product.Inventory.InventoryValuationMethod =
				req.Inventory.InventoryValuationMethod
		}

		if req.Inventory.ReorderPoint >= 0 {
			product.Inventory.ReorderPoint =
				req.Inventory.ReorderPoint
		}
	}

	if req.ReturnPolicy != nil {
		product.ReturnPolicy.Returnable =
			req.ReturnPolicy.Returnable
	}

	if req.ProductDetails != nil {
		if req.ProductDetails.Unit != "" {
			product.ProductDetails.Unit = req.ProductDetails.Unit
		}
		if req.ProductDetails.BaseSKU != "" {
			product.ProductDetails.BaseSKU = req.ProductDetails.BaseSKU
		}
		if req.ProductDetails.Description != "" {
			product.ProductDetails.Description =
				req.ProductDetails.Description
		}

		if len(req.ProductDetails.Variants) > 0 {
			product.ProductDetails.ProductVariants = make(
				[]models.ProductVariant,
				len(req.ProductDetails.Variants),
			)

			for index, variant := range req.ProductDetails.Variants {
				attributes := make(
					[]models.ProductVariantAttribute,
					0,
					len(variant.AttributeMap),
				)

				for key, value := range variant.AttributeMap {
					attributes = append(
						attributes,
						models.ProductVariantAttribute{
							Key:   key,
							Value: value,
						},
					)
				}

				product.ProductDetails.ProductVariants[index] =
					models.ProductVariant{
						SKU:           variant.SKU,
						VariantName:   variant.VariantName,
						Attributes:    attributes,
						SellingPrice:  variant.SellingPrice,
						CostPrice:     variant.CostPrice,
						StockQuantity: variant.StockQuantity,
						IsActive:      variant.IsActive,
					}
			}
		}

		if len(req.ProductDetails.AttributeDefinitions) > 0 {
			product.ProductDetails.AttributeDefinitions = make(
				models.ProductAttributeDefinitions,
				len(req.ProductDetails.AttributeDefinitions),
			)

			for index, attribute :=
				range req.ProductDetails.AttributeDefinitions {
				product.ProductDetails.AttributeDefinitions[index] =
					models.ProductAttributeDefinition{
						Key:     attribute.Key,
						Options: attribute.Options,
					}
			}
		}
	}

	product.UpdatedAt = time.Now()

	if err := s.repo.UpdateByCompany(
		product,
		companyID,
	); err != nil {
		return nil, err
	}

	updatedProduct, err := s.repo.FindByIDAndCompany(
		id,
		companyID,
	)
	if err != nil {
		return nil, err
	}

	return output.ToProductOutput(updatedProduct)
}

func (s *productService) DeleteProduct(
	id string,
	createdBy string,
) error {
	var userID uint
	if _, err := fmt.Sscanf(createdBy, "%d", &userID); err != nil {
		return fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	if user.CompanyID == nil {
		return fmt.Errorf("user is not assigned to a company")
	}

	return s.DeleteProductForCompany(
		id,
		*user.CompanyID,
	)
}

func (s *productService) DeleteProductForCompany(
	id string,
	companyID uint,
) error {
	return s.repo.DeleteByCompany(
		id,
		companyID,
	)
}

func (s *productService) GetProductVariants(
	productID string,
) ([]models.ProductVariant, error) {
	return s.repo.GetProductVariantsByProductID(productID)
}

func (s *productService) GetProductVariantsByCompany(
	productID string,
	companyID uint,
) ([]models.ProductVariant, error) {
	if _, err := s.repo.FindByIDAndCompany(
		productID,
		companyID,
	); err != nil {
		return nil, fmt.Errorf("product not found")
	}

	return s.repo.GetProductVariantsByProductIDAndCompany(
		productID,
		companyID,
	)
}

func (s *productService) EnableInventoryTracking(
	productID string,
) error {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	product.Inventory.TrackInventory = true
	product.UpdatedAt = time.Now()

	return s.repo.Update(product)
}

func (s *productService) IsInventoryTrackingEnabled(
	productID string,
) (bool, error) {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return false, fmt.Errorf("product not found: %w", err)
	}

	return product.Inventory.TrackInventory, nil
}

func (s *productService) GetProductWithInventoryStatus(
	productID string,
) (map[string]interface{}, error) {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	variants := make(
		[]map[string]interface{},
		len(product.ProductDetails.ProductVariants),
	)

	for index, variant := range product.ProductDetails.ProductVariants {
		variants[index] = map[string]interface{}{
			"sku":            variant.SKU,
			"variant_name":   variant.VariantName,
			"selling_price":  variant.SellingPrice,
			"cost_price":     variant.CostPrice,
			"stock_quantity": variant.StockQuantity,
			"reorder_level":  variant.ReorderLevel,
			"is_active":      variant.IsActive,
		}
	}

	return map[string]interface{}{
		"product_id":         productID,
		"name":               product.Name,
		"inventory_tracking": product.Inventory.TrackInventory,
		"reorder_point":      product.Inventory.ReorderPoint,
		"purchase_price":     product.PurchaseInfo.CostPrice,
		"selling_price":      product.SalesInfo.SellingPrice,
		"variants":           variants,
	}, nil
}