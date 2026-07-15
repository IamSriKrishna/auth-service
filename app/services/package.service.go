package services

import (
	"errors"
	"fmt"
	"log"
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

	CreatePackageForCompany(pkgInput *input.CreatePackageInput, userID string, companyID uint) (*output.PackageOutput, error)
	GetPackageForCompany(id string, companyID uint) (*output.PackageOutput, error)
	GetAllPackagesForCompany(companyID uint, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesByCustomerForCompany(customerID, companyID uint, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesBySalesOrderForCompany(salesOrderID string, companyID uint, limit, offset int) ([]output.PackageOutput, int64, error)
	GetPackagesByStatusForCompany(status string, companyID uint, limit, offset int) ([]output.PackageOutput, int64, error)
	UpdatePackageForCompany(id string, pkgInput *input.UpdatePackageInput, userID string, companyID uint) (*output.PackageOutput, error)
	UpdatePackageStatusForCompany(id, status, userID string, companyID uint) (*output.PackageOutput, error)
	DeletePackageForCompany(id string, companyID uint) error
}

type packageService struct {
	pkgRepo            repo.PackageRepository
	soRepo             repo.SalesOrderRepository
	customerRepo       repo.CustomerRepository
	productRepo        repo.ProductRepository
	stockManagementSvc StockManagementService
}

func NewPackageService(
	pkgRepo repo.PackageRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	productRepo repo.ProductRepository,
	stockManagementSvc StockManagementService,
) PackageService {
	return &packageService{
		pkgRepo:            pkgRepo,
		soRepo:             soRepo,
		customerRepo:       customerRepo,
		productRepo:        productRepo,
		stockManagementSvc: stockManagementSvc,
	}
}

func packagesToOutput(packages []models.Package) ([]output.PackageOutput, error) {
	result := make([]output.PackageOutput, 0, len(packages))
	for i := range packages {
		out, err := output.ToPackageOutput(&packages[i])
		if err != nil {
			return nil, err
		}
		result = append(result, *out)
	}
	return result, nil
}

func (s *packageService) CreatePackage(req *input.CreatePackageInput, userID string) (*output.PackageOutput, error) {
	if req == nil {
		return nil, errors.New("package input cannot be nil")
	}
	so, err := s.soRepo.FindByID(req.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}
	customer, err := s.customerRepo.FindByID(req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return s.createPackage(req, userID, so, customer)
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
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) GetPackagesByCustomer(customerID uint, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindByCustomer(customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) GetPackagesBySalesOrder(salesOrderID string, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindBySalesOrder(salesOrderID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) GetPackagesByStatus(status string, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) UpdatePackage(id string, req *input.UpdatePackageInput, userID string) (*output.PackageOutput, error) {
	if req == nil {
		return nil, errors.New("package input cannot be nil")
	}
	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	s.applyUpdate(pkg, req, userID)
	updated, err := s.pkgRepo.Update(id, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}
	return output.ToPackageOutput(updated)
}

func (s *packageService) UpdatePackageStatus(id, status, userID string) (*output.PackageOutput, error) {
	pkg, err := s.pkgRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	return s.updateStatus(pkg, status, userID, func(p *models.Package) (*models.Package, error) {
		return s.pkgRepo.Update(id, p)
	})
}

func (s *packageService) DeletePackage(id string) error {
	return s.pkgRepo.Delete(id)
}

func (s *packageService) CreatePackageForCompany(req *input.CreatePackageInput, userID string, companyID uint) (*output.PackageOutput, error) {
	if req == nil {
		return nil, errors.New("package input cannot be nil")
	}
	if companyID == 0 {
		return nil, errors.New("invalid company")
	}
	so, err := s.soRepo.FindByIDAndCompany(req.SalesOrderID, companyID)
	if err != nil {
		return nil, errors.New("sales order not found in your company")
	}
	customer, err := s.customerRepo.FindByIDAndCompany(req.CustomerID, companyID)
	if err != nil {
		return nil, errors.New("customer not found in your company")
	}
	return s.createPackage(req, userID, so, customer)
}

func (s *packageService) GetPackageForCompany(id string, companyID uint) (*output.PackageOutput, error) {
	pkg, err := s.pkgRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("package not found")
	}
	return output.ToPackageOutput(pkg)
}

func (s *packageService) GetAllPackagesForCompany(companyID uint, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindAllByCompany(companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) GetPackagesByCustomerForCompany(customerID, companyID uint, limit, offset int) ([]output.PackageOutput, int64, error) {
	if _, err := s.customerRepo.FindByIDAndCompany(customerID, companyID); err != nil {
		return nil, 0, errors.New("customer not found in your company")
	}
	packages, total, err := s.pkgRepo.FindByCustomerAndCompany(customerID, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) GetPackagesBySalesOrderForCompany(salesOrderID string, companyID uint, limit, offset int) ([]output.PackageOutput, int64, error) {
	if _, err := s.soRepo.FindByIDAndCompany(salesOrderID, companyID); err != nil {
		return nil, 0, errors.New("sales order not found in your company")
	}
	packages, total, err := s.pkgRepo.FindBySalesOrderAndCompany(salesOrderID, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) GetPackagesByStatusForCompany(status string, companyID uint, limit, offset int) ([]output.PackageOutput, int64, error) {
	packages, total, err := s.pkgRepo.FindByStatusAndCompany(status, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out, err := packagesToOutput(packages)
	return out, total, err
}

func (s *packageService) UpdatePackageForCompany(id string, req *input.UpdatePackageInput, userID string, companyID uint) (*output.PackageOutput, error) {
	if req == nil {
		return nil, errors.New("package input cannot be nil")
	}
	pkg, err := s.pkgRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("package not found")
	}
	s.applyUpdate(pkg, req, userID)
	updated, err := s.pkgRepo.UpdateByCompany(id, companyID, pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}
	return output.ToPackageOutput(updated)
}

func (s *packageService) UpdatePackageStatusForCompany(id, status, userID string, companyID uint) (*output.PackageOutput, error) {
	pkg, err := s.pkgRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("package not found")
	}
	return s.updateStatus(pkg, status, userID, func(p *models.Package) (*models.Package, error) {
		return s.pkgRepo.UpdateByCompany(id, companyID, p)
	})
}

func (s *packageService) DeletePackageForCompany(id string, companyID uint) error {
	if _, err := s.pkgRepo.FindByIDAndCompany(id, companyID); err != nil {
		return errors.New("package not found")
	}
	return s.pkgRepo.DeleteByCompany(id, companyID)
}

func (s *packageService) createPackage(req *input.CreatePackageInput, userID string, so *models.SalesOrder, customer *models.Customer) (*output.PackageOutput, error) {
	if so == nil || customer == nil {
		return nil, errors.New("sales order and customer are required")
	}
	if so.CustomerID != req.CustomerID {
		return nil, errors.New("customer does not match sales order")
	}

	slipNo := ""
	if req.PackageSlipNo != nil && *req.PackageSlipNo != "" {
		slipNo = *req.PackageSlipNo
	} else {
		var err error
		slipNo, err = s.pkgRepo.GetNextPackageSlipNo()
		if err != nil {
			return nil, fmt.Errorf("failed to generate package slip number: %w", err)
		}
	}

	inputMap := make(map[uint]*input.PackageLineItemInput)
	for i := range req.Items {
		inputMap[req.Items[i].SalesOrderItemID] = &req.Items[i]
	}

	items := make([]models.PackageItem, 0)
	for _, soItem := range so.LineItems {
		packed := 0.0
		inputItem, exists := inputMap[soItem.ID]
		if exists {
			packed = inputItem.PackedQty
		} else if len(req.Items) > 0 {
			continue
		}
		if packed < 0 || packed > soItem.Quantity {
			return nil, fmt.Errorf("invalid packed quantity for sales order item %d", soItem.ID)
		}
		items = append(items, models.PackageItem{
			SalesOrderItemID: soItem.ID,
			OrderedQty:       soItem.Quantity,
			PackedQty:        packed,
		})
	}
	if len(items) == 0 {
		return nil, errors.New("package must contain at least one sales order item")
	}

	pkg := &models.Package{
		ID:            uuid.New().String(),
		PackageSlipNo: slipNo,
		SalesOrderID:  req.SalesOrderID,
		SalesOrder:    so,
		CustomerID:    req.CustomerID,
		Customer:      customer,
		PackageDate:   req.PackageDate,
		Items:         items,
		Status:        domain.PackageStatusCreated,
		InternalNotes: req.InternalNotes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}

	created, err := s.pkgRepo.Create(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}
	return output.ToPackageOutput(created)
}

func (s *packageService) applyUpdate(pkg *models.Package, req *input.UpdatePackageInput, userID string) {
	if req.PackageDate != nil {
		pkg.PackageDate = *req.PackageDate
	}
	if req.InternalNotes != nil {
		pkg.InternalNotes = *req.InternalNotes
	}
	if req.Status != nil {
		pkg.Status = domain.PackageStatus(*req.Status)
	}
	if len(req.Items) > 0 {
		updates := make(map[uint]float64)
		for _, item := range req.Items {
			updates[item.SalesOrderItemID] = item.PackedQty
		}
		for i := range pkg.Items {
			if packed, ok := updates[pkg.Items[i].SalesOrderItemID]; ok {
				pkg.Items[i].PackedQty = packed
			}
		}
	}
	pkg.UpdatedBy = userID
	pkg.UpdatedAt = time.Now()
}

func (s *packageService) updateStatus(
	pkg *models.Package,
	status string,
	userID string,
	save func(*models.Package) (*models.Package, error),
) (*output.PackageOutput, error) {
	switch status {
	case "created", "packed", "shipped", "delivered", "cancelled":
	default:
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	if status == "packed" && pkg.Status != domain.PackageStatus("packed") {
		for _, item := range pkg.Items {
			if item.ProductID == nil || *item.ProductID == "" || item.PackedQty <= 0 {
				continue
			}
			if s.productRepo != nil {
				if _, err := s.productRepo.FindByID(*item.ProductID); err != nil {
					return nil, fmt.Errorf("product %s not found: %w", *item.ProductID, err)
				}
			}

			var err error
			if item.VariantSKU != nil && *item.VariantSKU != "" {
				err = s.stockManagementSvc.RecordOutboundMovementWithVariant(
					*item.ProductID,
					*item.VariantSKU,
					"package",
					pkg.ID,
					pkg.PackageSlipNo,
					item.PackedQty,
					fmt.Sprintf("Product variant packed for package %s", pkg.PackageSlipNo),
					userID,
				)
			} else {
				err = s.stockManagementSvc.RecordOutboundMovement(
					*item.ProductID,
					"package",
					pkg.ID,
					pkg.PackageSlipNo,
					item.PackedQty,
					fmt.Sprintf("Product packed for package %s", pkg.PackageSlipNo),
					userID,
				)
			}
			if err != nil {
				log.Printf("[PACKAGE] stock deduction failed for product %s: %v", *item.ProductID, err)
				return nil, err
			}
		}
	}

	pkg.Status = domain.PackageStatus(status)
	pkg.UpdatedBy = userID
	pkg.UpdatedAt = time.Now()

	updated, err := save(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to update package status: %w", err)
	}
	return output.ToPackageOutput(updated)
}
