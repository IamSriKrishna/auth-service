package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type PackageService interface {
	// Basic CRUD Operations
	CreatePackage(pkgInput *input.CreatePackageInput, userID string) (*output.PackageOutput, error)
	GetPackage(id string) (*output.PackageOutput, error)
	GetAllPackages(limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesByCustomer(customerID uint, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesBySalesOrder(salesOrderID string, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesByStatus(status string, limit, offset int) ([]output.PackageOutput, int64, error)
	UpdatePackage(id string, pkgInput *input.UpdatePackageInput, userID string) (*output.PackageOutput, error)

	// Step 4: Selling to Customers - Package Prep
	// Prepare items for shipping and update package status
	UpdatePackageStatus(id string, status string, userID string) (*output.PackageOutput, error)
	DeletePackage(id string) error
}

type packageService struct {
	pkgRepo      repo.PackageRepository
	soRepo       repo.SalesOrderRepository
	customerRepo repo.CustomerRepository
	itemRepo     repo.ItemRepository
}

func NewPackageService(
	pkgRepo repo.PackageRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	itemRepo repo.ItemRepository,
) PackageService {
	return &packageService{
		pkgRepo:      pkgRepo,
		soRepo:       soRepo,
		customerRepo: customerRepo,
		itemRepo:     itemRepo,
	}
}

func (s *packageService) CreatePackage(pkgInput *input.CreatePackageInput, userID string) (*output.PackageOutput, error) {
	if pkgInput == nil {
		return nil, errors.New("package input cannot be nil")
	}

	// Fetch and validate sales order
	so, err := s.soRepo.FindByID(pkgInput.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	// Fetch and validate customer
	customer, err := s.customerRepo.FindByID(pkgInput.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Verify customer matches sales order
	if so.CustomerID != pkgInput.CustomerID {
		return nil, errors.New("customer does not match sales order")
	}

	// Generate or use provided package slip number
	var slipNo string
	if pkgInput.PackageSlipNo != nil && *pkgInput.PackageSlipNo != "" {
		slipNo = *pkgInput.PackageSlipNo
	} else {
		generatedSlip, err := s.pkgRepo.GetNextPackageSlipNo()
		if err != nil {
			return nil, fmt.Errorf("failed to generate package slip number: %w", err)
		}
		slipNo = generatedSlip
	}

	// Build a map of input items for quick lookup by SalesOrderItemID
	inputItemsMap := make(map[uint]*input.PackageLineItemInput)
	if len(pkgInput.Items) > 0 {
		for i := range pkgInput.Items {
			inputItemsMap[pkgInput.Items[i].SalesOrderItemID] = &pkgInput.Items[i]
		}
	}

	// Populate package items from sales order line items
	var packageItems []models.PackageItem
	for _, soLineItem := range so.LineItems {
		// Check if this item is in the input items map
		var packedQty float64 = 0
		if inputItem, exists := inputItemsMap[soLineItem.ID]; exists {
			packedQty = inputItem.PackedQty
		} else if len(pkgInput.Items) == 0 {
			// If no items specified in input, default packed qty to 0 (user will fill it manually)
			packedQty = 0
		} else {
			// If items were specified but this one wasn't, skip it
			continue
		}

		// Fetch item details
		item, err := s.itemRepo.FindByID(soLineItem.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item %s not found: %w", soLineItem.ItemID, err)
		}

		// Trim VariantSKU to remove whitespace and newlines
		var variantSKU *string
		if soLineItem.VariantSKU != nil {
			trimmed := strings.TrimSpace(*soLineItem.VariantSKU)
			variantSKU = &trimmed
		}

		packageItem := models.PackageItem{
			SalesOrderItemID: soLineItem.ID,
			ItemID:           soLineItem.ItemID,
			Item:             item,
			VariantSKU:       variantSKU,
			Variant:          soLineItem.Variant,
			OrderedQty:       soLineItem.Quantity,
			PackedQty:        packedQty,
			VariantDetails:   soLineItem.VariantDetails,
		}

		packageItems = append(packageItems, packageItem)
	}

	// Create package model
	pkg := &models.Package{
		ID:            uuid.New().String(),
		PackageSlipNo: slipNo,
		SalesOrderID:  pkgInput.SalesOrderID,
		CustomerID:    pkgInput.CustomerID,
		PackageDate:   pkgInput.PackageDate,
		Status:        domain.PackageStatusCreated,
		InternalNotes: pkgInput.InternalNotes,
		Items:         packageItems,
		SalesOrder:    so,
		Customer:      customer,
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}

	// Save to repository
	createdPkg, err := s.pkgRepo.Create(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	return output.ToPackageOutput(createdPkg)
}

func (s *packageService) GetPackage(id string) (*output.PackageOutput, error) {
	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	return output.ToPackageOutput(pkg)
}

func (s *packageService) GetAllPackages(limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.PackageOutput, 0)
	for _, pkg := range packages {
		if out, err := output.ToPackageOutput(&pkg); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *packageService) GetPackagesByCustomer(customerID uint, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindByCustomer(customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.PackageOutput, 0)
	for _, pkg := range packages {
		if out, err := output.ToPackageOutput(&pkg); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *packageService) GetPackagesBySalesOrder(salesOrderID string, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindBySalesOrder(salesOrderID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.PackageOutput, 0)
	for _, pkg := range packages {
		if out, err := output.ToPackageOutput(&pkg); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *packageService) GetPackagesByStatus(status string, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.PackageOutput, 0)
	for _, pkg := range packages {
		if out, err := output.ToPackageOutput(&pkg); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *packageService) UpdatePackage(id string, pkgInput *input.UpdatePackageInput, userID string) (*output.PackageOutput, error) {
	if pkgInput == nil {
		return nil, errors.New("package input cannot be nil")
	}

	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	if pkgInput.PackageDate != nil {
		pkg.PackageDate = *pkgInput.PackageDate
	}

	if pkgInput.InternalNotes != nil {
		pkg.InternalNotes = *pkgInput.InternalNotes
	}

	if pkgInput.Status != nil {
		pkg.Status = domain.PackageStatus(*pkgInput.Status)
	}

	// Update packed quantities for specified items
	if len(pkgInput.Items) > 0 {
		// Build a map of input items for quick lookup
		inputItemsMap := make(map[uint]float64)
		for _, itemInput := range pkgInput.Items {
			inputItemsMap[itemInput.SalesOrderItemID] = itemInput.PackedQty
		}

		// Update existing package items
		for i := range pkg.Items {
			if packedQty, exists := inputItemsMap[pkg.Items[i].SalesOrderItemID]; exists {
				pkg.Items[i].PackedQty = packedQty
			}
		}
	}

	pkg.UpdatedBy = userID
	pkg.UpdatedAt = time.Now()

	updatedPkg, err := s.pkgRepo.Update(id, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}

	return output.ToPackageOutput(updatedPkg)
}

func (s *packageService) UpdatePackageStatus(id string, status string, userID string) (*output.PackageOutput, error) {
	switch status {
	case "created", "packed", "shipped", "delivered", "cancelled":
	default:
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	pkg.Status = domain.PackageStatus(status)
	pkg.UpdatedBy = userID
	pkg.UpdatedAt = time.Now()

	updatedPkg, err := s.pkgRepo.Update(id, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package status: %w", err)
	}

	return output.ToPackageOutput(updatedPkg)
}

func (s *packageService) DeletePackage(id string) error {
	return s.pkgRepo.Delete(id)
}
