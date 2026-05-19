package services

import (
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ManufacturerService interface {
	Create(req *input.CreateManufacturerInput, userID uint, companyID uint) (*output.ManufacturerOutput, error)
	GetByID(id string) (*output.ManufacturerOutput, error)
	GetAll(limit, offset int) (*output.ListManufacturersOutput, error)
	GetAllWithFilter(limit, offset int, companyID *uint, productGroupID *string) (*output.ListManufacturersOutput, error)
	GetByProductGroupID(productGroupID string) ([]output.ManufacturerOutput, error)
	Update(id string, req *input.UpdateManufacturerInput) (*output.ManufacturerOutput, error)
	Delete(id string) error
}

type manufacturerService struct {
	manufacturerRepo repo.ManufacturerRepository
	productGroupRepo repo.ProductGroupRepository
	employeeRepo     repo.EmployeeRepository
	productStockRepo repo.ProductStockRepository
	stockMgmtService StockManagementService
}

func NewManufacturerService(repo repo.ManufacturerRepository) ManufacturerService {
	return &manufacturerService{manufacturerRepo: repo}
}

// NewManufacturerServiceWithDependencies creates a new manufacturer service with all required dependencies
func NewManufacturerServiceWithDependencies(
	manufacturerRepo repo.ManufacturerRepository,
	productGroupRepo repo.ProductGroupRepository,
	employeeRepo repo.EmployeeRepository,
	productStockRepo repo.ProductStockRepository,
	stockMgmtService StockManagementService,
) ManufacturerService {
	return &manufacturerService{
		manufacturerRepo: manufacturerRepo,
		productGroupRepo: productGroupRepo,
		employeeRepo:     employeeRepo,
		productStockRepo: productStockRepo,
		stockMgmtService: stockMgmtService,
	}
}

// Create creates a new manufacturer with product group and employee assignments
func (s *manufacturerService) Create(req *input.CreateManufacturerInput, userID uint, companyID uint) (*output.ManufacturerOutput, error) {
	// Validate product group exists
	productGroup, err := s.productGroupRepo.FindByID(req.ProductGroupID)
	if err != nil {
		return nil, fmt.Errorf("product group not found: %v", err)
	}

	if len(productGroup.Components) == 0 {
		return nil, fmt.Errorf("product group has no components for manufacturing")
	}

	// Validate all employees exist
	for _, empInput := range req.Employees {
		emp, err := s.employeeRepo.GetByID(empInput.EmployeeID)
		if err != nil {
			return nil, fmt.Errorf("employee with ID %d not found: %v", empInput.EmployeeID, err)
		}
		if emp == nil {
			return nil, fmt.Errorf("employee with ID %d not found", empInput.EmployeeID)
		}
	}

	// Check inventory availability for all components
	availabilityIssues := []string{}
	for _, component := range productGroup.Components {
		requiredQty := component.Quantity * req.Quantity

		if component.Product != nil && !component.Product.IsResource {
			// Check product stock (skip if no tracking enabled)
			productStock, err := s.productStockRepo.GetByProductID(component.ProductID)
			if err != nil {
				// Stock record not found - skip inventory check if tracking not enabled
				// This product might not have inventory tracking
				continue
			}

			if productStock != nil && productStock.AvailableStock < requiredQty {
				availabilityIssues = append(availabilityIssues, fmt.Sprintf(
					"Insufficient stock for %s: need %.2f, available %.2f",
					component.Product.Name,
					requiredQty,
					productStock.AvailableStock,
				))
			}
		}
	}

	if len(availabilityIssues) > 0 {
		return nil, fmt.Errorf("inventory check failed: %v", availabilityIssues)
	}

	// Create manufacturer with employee assignments as JSON
	manufacturerID := "mfg_" + uuid.New().String()[:12]

	// Convert input employees to model employee assignments
	employees := make(models.EmployeeAssignments, len(req.Employees))
	for i, empInput := range req.Employees {
		employees[i] = models.ManufacturerEmployeeAssignment{
			EmployeeID:  empInput.EmployeeID,
			ServiceCost: empInput.ServiceCost,
			CostType:    empInput.CostType,
		}
	}

	manufacturer := &models.Manufacturer{
		ID:             manufacturerID,
		Name:           req.Name,
		ProductGroupID: req.ProductGroupID,
		Quantity:       req.Quantity,
		Status:         "pending",
		Description:    req.Description,
		Employees:      employees,
		CompanyID:      companyID,
		UserID:         userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Deduct stock for all components BEFORE creating manufacturer
	// This ensures transactional consistency - if stock deduction fails, manufacturing isn't created
	for _, component := range productGroup.Components {
		if component.Product != nil && !component.Product.IsResource {
			requiredQty := float64(component.Quantity * req.Quantity)

			notes := fmt.Sprintf("Manufacturing %s: %.0f units of %s", req.Name, req.Quantity, component.Product.Name)

			// Deduct stock - use variant-aware deduction if variant SKU is present
			var err error
			if component.VariantSku != nil && *component.VariantSku != "" {
				// Deduct from variant stock
				err = s.stockMgmtService.RecordOutboundMovementWithVariant(
					component.ProductID,
					*component.VariantSku,
					"PRODUCTION_USAGE",
					manufacturerID,
					req.Name,
					requiredQty,
					notes,
					fmt.Sprintf("%d", userID),
				)
			} else {
				// Deduct from product stock
				err = s.stockMgmtService.RecordOutboundMovement(
					component.ProductID,
					"PRODUCTION_USAGE",
					manufacturerID,
					req.Name,
					requiredQty,
					notes,
					fmt.Sprintf("%d", userID),
				)
			}

			if err != nil {
				return nil, fmt.Errorf("failed to deduct stock for component %s (%s): %v", component.Product.Name, component.ProductID, err)
			}
		}
	}

	// Create manufacturer after stock deduction succeeds
	if err := s.manufacturerRepo.Create(manufacturer); err != nil {
		return nil, fmt.Errorf("failed to create manufacturer: %v", err)
	}

	// Fetch the created manufacturer with all relationships
	createdManufacturer, err := s.manufacturerRepo.FindByStringID(manufacturerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created manufacturer: %v", err)
	}

	return convertManufacturerToOutput(createdManufacturer), nil
}

func (s *manufacturerService) GetByID(id string) (*output.ManufacturerOutput, error) {
	manufacturer, err := s.manufacturerRepo.FindByStringID(id)
	if err != nil {
		return nil, err
	}
	return convertManufacturerToOutput(manufacturer), nil
}

func (s *manufacturerService) GetAll(limit, offset int) (*output.ListManufacturersOutput, error) {
	manufacturers, count, err := s.manufacturerRepo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}
	return convertManufacturersToOutput(manufacturers, int(count)), nil
}

func (s *manufacturerService) GetAllWithFilter(limit, offset int, companyID *uint, productGroupID *string) (*output.ListManufacturersOutput, error) {
	manufacturers, count, err := s.manufacturerRepo.FindAllWithFilter(limit, offset, companyID, productGroupID)
	if err != nil {
		return nil, err
	}
	return convertManufacturersToOutput(manufacturers, int(count)), nil
}

func (s *manufacturerService) GetByProductGroupID(productGroupID string) ([]output.ManufacturerOutput, error) {
	manufacturers, err := s.manufacturerRepo.FindByProductGroupID(productGroupID)
	if err != nil {
		return nil, err
	}
	outputs := make([]output.ManufacturerOutput, len(manufacturers))
	for i, m := range manufacturers {
		outputs[i] = *convertManufacturerToOutput(&m)
	}
	return outputs, nil
}

func (s *manufacturerService) Update(id string, req *input.UpdateManufacturerInput) (*output.ManufacturerOutput, error) {
	manufacturer, err := s.manufacturerRepo.FindByStringID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		manufacturer.Name = *req.Name
	}
	if req.Quantity != nil && *req.Quantity > 0 {
		manufacturer.Quantity = *req.Quantity
	}
	if req.Status != nil {
		manufacturer.Status = *req.Status
	}
	if req.Description != nil {
		manufacturer.Description = *req.Description
	}

	manufacturer.UpdatedAt = time.Now()
	if err := s.manufacturerRepo.Update(manufacturer); err != nil {
		return nil, err
	}

	// Fetch updated manufacturer
	updated, err := s.manufacturerRepo.FindByStringID(id)
	if err != nil {
		return nil, err
	}

	return convertManufacturerToOutput(updated), nil
}

func (s *manufacturerService) Delete(id string) error {
	return s.manufacturerRepo.DeleteByStringID(id)
}

// Helper functions

func convertManufacturerToOutput(manufacturer *models.Manufacturer) *output.ManufacturerOutput {
	if manufacturer == nil {
		return nil
	}

	var pgOutput *output.ProductGroupOutput
	if manufacturer.ProductGroup != nil {
		pgOutput = convertProductGroupToOutput(manufacturer.ProductGroup)
	}

	employeeOutputs := make([]output.EmployeeAssignmentOutput, len(manufacturer.Employees))
	for i, emp := range manufacturer.Employees {
		employeeOutputs[i] = output.EmployeeAssignmentOutput{
			EmployeeID:  emp.EmployeeID,
			ServiceCost: emp.ServiceCost,
			CostType:    emp.CostType,
			Employee:    nil, // Employee details can be fetched separately if needed
		}
	}

	return &output.ManufacturerOutput{
		ID:             manufacturer.ID,
		Name:           manufacturer.Name,
		ProductGroupID: manufacturer.ProductGroupID,
		ProductGroup:   pgOutput,
		Quantity:       manufacturer.Quantity,
		Status:         manufacturer.Status,
		Description:    manufacturer.Description,
		Employees:      employeeOutputs,
		CreatedAt:      manufacturer.CreatedAt,
		UpdatedAt:      manufacturer.UpdatedAt,
	}
}

func convertManufacturersToOutput(manufacturers []models.Manufacturer, count int) *output.ListManufacturersOutput {
	outputs := make([]output.ManufacturerOutput, len(manufacturers))
	for i, m := range manufacturers {
		output := convertManufacturerToOutput(&m)
		if output != nil {
			outputs[i] = *output
		}
	}
	return &output.ListManufacturersOutput{
		Manufacturers: outputs,
		TotalCount:    count,
	}
}

func convertProductGroupToOutput(pg *models.ProductGroup) *output.ProductGroupOutput {
	if pg == nil {
		return nil
	}

	// Use the proper conversion function that includes all fields and calculates pricing
	pgOutput, err := output.ToProductGroupOutput(pg)
	if err != nil {
		// Fallback to basic conversion if there's an error
		return &output.ProductGroupOutput{
			ID:          pg.ID,
			Name:        pg.Name,
			Description: pg.Description,
			IsActive:    pg.IsActive,
		}
	}
	return pgOutput
}
