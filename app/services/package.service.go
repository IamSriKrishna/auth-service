package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type PackageService interface {
	CreatePackage(pkgInput *input.CreatePackageInput, userID string) (*output.PackageOutput, error)
	GetPackage(id string) (*output.PackageOutput, error)
	GetAllPackages(limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesByCustomer(customerID uint, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesBySalesOrder(salesOrderID string, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesByStatus(status string, limit, offset int) ([]output.PackageOutput, int64, error)
	UpdatePackage(id string, pkgInput *input.UpdatePackageInput, userID string) (*output.PackageOutput, error)
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

	// Verify sales order exists
	so, err := s.soRepo.FindByID(pkgInput.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	// Verify customer exists
	customer, err := s.customerRepo.FindByID(pkgInput.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Verify customer matches sales order
	if so.CustomerID != pkgInput.CustomerID {
		return nil, errors.New("customer does not match sales order")
	}

	// Generate package slip number
	slipNo, err := s.pkgRepo.GetNextPackageSlipNo()
	if err != nil {
		return nil, fmt.Errorf("failed to generate package slip number: %w", err)
	}

	// Create package items
	var packageItems []models.PackageItem
	for _, itemInput := range pkgInput.Items {
		// Verify item exists
		item, err := s.itemRepo.FindByID(itemInput.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item %s not found: %w", itemInput.ItemID, err)
		}

		packageItem := models.PackageItem{
			SalesOrderItemID: itemInput.SalesOrderItemID,
			ItemID:           itemInput.ItemID,
			VariantID:        itemInput.VariantID,
			OrderedQty:       itemInput.OrderedQty,
			PackedQty:        itemInput.PackedQty,
		}

		// Convert variant details if present
		if itemInput.VariantDetails != nil {
			packageItem.VariantDetails = models.VariantDetails(itemInput.VariantDetails)
		}

		// Set references to models
		packageItem.Item = item

		packageItems = append(packageItems, packageItem)
	}

	// Create package
	pkg := &models.Package{
		ID:            uuid.New().String(),
		PackageSlipNo: slipNo,
		SalesOrderID:  pkgInput.SalesOrderID,
		CustomerID:    pkgInput.CustomerID,
		PackageDate:   pkgInput.PackageDate,
		Status:        domain.PackageStatusCreated,
		InternalNotes: pkgInput.InternalNotes,
		Items:         packageItems,
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}

	// Set references
	pkg.SalesOrder = so
	pkg.Customer = customer

	// Save package
	createdPkg, err := s.pkgRepo.Create(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	// Convert to output
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

	// Get existing package
	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	// Update fields
	if pkgInput.PackageDate != nil {
		pkg.PackageDate = *pkgInput.PackageDate
	}

	if pkgInput.InternalNotes != nil {
		pkg.InternalNotes = *pkgInput.InternalNotes
	}

	if pkgInput.Status != nil {
		pkg.Status = domain.PackageStatus(*pkgInput.Status)
	}

	// Update items if provided
	if len(pkgInput.Items) > 0 {
		var packageItems []models.PackageItem
		for _, itemInput := range pkgInput.Items {
			item, err := s.itemRepo.FindByID(itemInput.ItemID)
			if err != nil {
				return nil, fmt.Errorf("item %s not found: %w", itemInput.ItemID, err)
			}

			packageItem := models.PackageItem{
				SalesOrderItemID: itemInput.SalesOrderItemID,
				ItemID:           itemInput.ItemID,
				VariantID:        itemInput.VariantID,
				OrderedQty:       itemInput.OrderedQty,
				PackedQty:        itemInput.PackedQty,
				Item:             item,
			}

			if itemInput.VariantDetails != nil {
				packageItem.VariantDetails = models.VariantDetails(itemInput.VariantDetails)
			}

			packageItems = append(packageItems, packageItem)
		}
		pkg.Items = packageItems
	}

	pkg.UpdatedBy = userID
	pkg.UpdatedAt = time.Now()

	// Update package
	updatedPkg, err := s.pkgRepo.Update(id, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}

	return output.ToPackageOutput(updatedPkg)
}

func (s *packageService) UpdatePackageStatus(id string, status string, userID string) (*output.PackageOutput, error) {
	// Verify status is valid
	switch status {
	case "created", "packed", "shipped", "delivered", "cancelled":
	default:
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// Get existing package
	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	// Update status
	pkg.Status = domain.PackageStatus(status)
	pkg.UpdatedBy = userID
	pkg.UpdatedAt = time.Now()

	// Update in database
	updatedPkg, err := s.pkgRepo.Update(id, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package status: %w", err)
	}

	return output.ToPackageOutput(updatedPkg)
}

func (s *packageService) DeletePackage(id string) error {
	return s.pkgRepo.Delete(id)
}
