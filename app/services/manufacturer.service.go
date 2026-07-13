package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ManufacturerService interface {
	// Existing methods retained for internal/backward compatibility.
	Create(
		req *input.CreateManufacturerInput,
		userID uint,
		companyID uint,
	) (*output.ManufacturerOutput, error)

	GetByID(
		id string,
	) (*output.ManufacturerOutput, error)

	GetAll(
		limit int,
		offset int,
	) (*output.ListManufacturersOutput, error)

	GetAllWithFilter(
		limit int,
		offset int,
		companyID *uint,
		productGroupID *string,
	) (*output.ListManufacturersOutput, error)

	GetByProductGroupID(
		productGroupID string,
	) ([]output.ManufacturerOutput, error)

	Update(
		id string,
		req *input.UpdateManufacturerInput,
	) (*output.ManufacturerOutput, error)

	Delete(
		id string,
	) error

	// Company-scoped methods used by authenticated routes.
	CreateForCompany(
		req *input.CreateManufacturerInput,
		userID uint,
		companyID uint,
	) (*output.ManufacturerOutput, error)

	GetByIDForCompany(
		id string,
		companyID uint,
	) (*output.ManufacturerOutput, error)

	GetAllForCompany(
		companyID uint,
		limit int,
		offset int,
	) (*output.ListManufacturersOutput, error)

	GetAllWithFilterForCompany(
		companyID uint,
		productGroupID *string,
		limit int,
		offset int,
	) (*output.ListManufacturersOutput, error)

	GetByProductGroupIDForCompany(
		productGroupID string,
		companyID uint,
	) ([]output.ManufacturerOutput, error)

	UpdateForCompany(
		id string,
		req *input.UpdateManufacturerInput,
		userID uint,
		companyID uint,
	) (*output.ManufacturerOutput, error)

	DeleteForCompany(
		id string,
		userID uint,
		companyID uint,
	) error
}

type manufacturerService struct {
	manufacturerRepo repo.ManufacturerRepository
	productGroupRepo repo.ProductGroupRepository
	employeeRepo     repo.EmployeeRepository
	productStockRepo repo.ProductStockRepository
	productRepo      repo.ProductRepository
	userRepo         repo.UserRepository
	stockMgmtService StockManagementService
}

func NewManufacturerService(
	manufacturerRepo repo.ManufacturerRepository,
) ManufacturerService {
	return &manufacturerService{
		manufacturerRepo: manufacturerRepo,
	}
}

func NewManufacturerServiceWithDependencies(
	manufacturerRepo repo.ManufacturerRepository,
	productGroupRepo repo.ProductGroupRepository,
	employeeRepo repo.EmployeeRepository,
	productStockRepo repo.ProductStockRepository,
	userRepo repo.UserRepository,
	productRepo repo.ProductRepository,
	stockMgmtService StockManagementService,
) ManufacturerService {
	return &manufacturerService{
		manufacturerRepo: manufacturerRepo,
		productGroupRepo: productGroupRepo,
		employeeRepo:     employeeRepo,
		productStockRepo: productStockRepo,
		productRepo:      productRepo,
		userRepo:         userRepo,
		stockMgmtService: stockMgmtService,
	}
}

func (s *manufacturerService) Create(
	req *input.CreateManufacturerInput,
	userID uint,
	companyID uint,
) (*output.ManufacturerOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	productGroup, err :=
		s.productGroupRepo.FindByID(req.ProductGroupID)
	if err != nil {
		return nil, fmt.Errorf(
			"product group not found: %v",
			err,
		)
	}

	return s.createWithValidatedProductGroup(
		req,
		productGroup,
		userID,
		companyID,
		false,
	)
}

func (s *manufacturerService) createWithValidatedProductGroup(
	req *input.CreateManufacturerInput,
	productGroup *models.ProductGroup,
	userID uint,
	companyID uint,
	companyScoped bool,
) (*output.ManufacturerOutput, error) {
	if productGroup == nil {
		return nil, errors.New("product group not found")
	}

	if req.Quantity <= 0 {
		return nil, errors.New(
			"manufacturing quantity must be greater than zero",
		)
	}

	if len(productGroup.Components) == 0 {
		return nil, errors.New(
			"product group has no components for manufacturing",
		)
	}

	for _, employeeInput := range req.Employees {
		var employee *models.Employee
		var err error

		if companyScoped {
			employee, err =
				s.employeeRepo.GetByIDAndCompanyID(
					employeeInput.EmployeeID,
					companyID,
				)
		} else {
			employee, err =
				s.employeeRepo.GetByID(
					employeeInput.EmployeeID,
				)
		}

		if err != nil || employee == nil {
			return nil, fmt.Errorf(
				"employee with ID %d not found in your company",
				employeeInput.EmployeeID,
			)
		}
	}

	availabilityIssues := make([]string, 0)

	for _, component := range productGroup.Components {
		requiredQuantity :=
			component.Quantity * req.Quantity

		if component.Product == nil ||
			component.Product.IsResource {
			continue
		}

		var productStock *models.ProductStock
		var err error

		if companyScoped {
			productStock, err =
				s.productStockRepo.GetByProductIDAndCompany(
					component.ProductID,
					companyID,
					true,
				)
		} else {
			productStock, err =
				s.productStockRepo.GetByProductID(
					component.ProductID,
				)
		}

		if err != nil {
			// Preserve the original behavior: a missing stock row may mean
			// inventory tracking is disabled for the component.
			continue
		}

		if productStock != nil &&
			productStock.AvailableStock < requiredQuantity {
			availabilityIssues = append(
				availabilityIssues,
				fmt.Sprintf(
					"insufficient stock for %s: need %.2f, available %.2f",
					component.Product.Name,
					requiredQuantity,
					productStock.AvailableStock,
				),
			)
		}
	}

	if len(availabilityIssues) > 0 {
		return nil, fmt.Errorf(
			"inventory check failed: %v",
			availabilityIssues,
		)
	}

	manufacturerID :=
		"mfg_" + uuid.New().String()[:12]

	employees := make(
		models.EmployeeAssignments,
		len(req.Employees),
	)

	for index, employeeInput := range req.Employees {
		employees[index] =
			models.ManufacturerEmployeeAssignment{
				EmployeeID:  employeeInput.EmployeeID,
				ServiceCost: employeeInput.ServiceCost,
				CostType:    employeeInput.CostType,
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

	// Validate every component before any stock deduction starts.
	if companyScoped {
		for _, component := range productGroup.Components {
			if component.Product == nil ||
				component.Product.IsResource {
				continue
			}

			if err := s.validateCompanyScopedComponentProduct(
				component.ProductID,
				companyID,
			); err != nil {
				return nil, err
			}
		}
	}

	if err := s.consumeComponentStockForManufacturer(
		productGroup,
		manufacturerID,
		req.Name,
		req.Quantity,
		userID,
		companyID,
		companyScoped,
	); err != nil {
		return nil, err
	}

	if err := s.manufacturerRepo.Create(
		manufacturer,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create manufacturer: %v",
			err,
		)
	}

	var (
		createdManufacturer *models.Manufacturer
		fetchErr            error
	)

	if companyScoped {
		createdManufacturer, fetchErr =
			s.manufacturerRepo.FindByStringIDAndCompany(
				manufacturerID,
				companyID,
			)
	} else {
		createdManufacturer, fetchErr =
			s.manufacturerRepo.FindByStringID(
				manufacturerID,
			)
	}

	if fetchErr != nil {
		return nil, fmt.Errorf(
			"failed to fetch created manufacturer: %w",
			fetchErr,
		)
	}

	return convertManufacturerToOutput(
		createdManufacturer,
	), nil
}

func (s *manufacturerService) GetByID(
	id string,
) (*output.ManufacturerOutput, error) {
	manufacturer, err :=
		s.manufacturerRepo.FindByStringID(id)
	if err != nil {
		return nil, err
	}

	return convertManufacturerToOutput(
		manufacturer,
	), nil
}

func (s *manufacturerService) GetAll(
	limit int,
	offset int,
) (*output.ListManufacturersOutput, error) {
	manufacturers, count, err :=
		s.manufacturerRepo.FindAll(
			limit,
			offset,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturersToOutput(
		manufacturers,
		int(count),
	), nil
}

func (s *manufacturerService) GetAllWithFilter(
	limit int,
	offset int,
	companyID *uint,
	productGroupID *string,
) (*output.ListManufacturersOutput, error) {
	manufacturers, count, err :=
		s.manufacturerRepo.FindAllWithFilter(
			limit,
			offset,
			companyID,
			productGroupID,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturersToOutput(
		manufacturers,
		int(count),
	), nil
}

func (s *manufacturerService) GetByProductGroupID(
	productGroupID string,
) ([]output.ManufacturerOutput, error) {
	manufacturers, err :=
		s.manufacturerRepo.FindByProductGroupID(
			productGroupID,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturerSlice(manufacturers), nil
}

func (s *manufacturerService) Update(
	id string,
	req *input.UpdateManufacturerInput,
) (*output.ManufacturerOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	manufacturer, err :=
		s.manufacturerRepo.FindByStringID(id)
	if err != nil {
		return nil, err
	}

	oldStatus := manufacturer.Status
	applyManufacturerUpdate(manufacturer, req)

	if shouldConsumeStockOnCompletion(oldStatus, manufacturer.Status) {
		if err := s.consumeComponentStockForManufacturer(
			manufacturer.ProductGroup,
			manufacturer.ID,
			manufacturer.Name,
			manufacturer.Quantity,
			0,
			0,
			false,
		); err != nil {
			return nil, err
		}
	}

	if err := s.manufacturerRepo.Update(
		manufacturer,
	); err != nil {
		return nil, err
	}

	updated, err :=
		s.manufacturerRepo.FindByStringID(id)
	if err != nil {
		return nil, err
	}

	return convertManufacturerToOutput(updated), nil
}

func (s *manufacturerService) Delete(
	id string,
) error {
	return s.manufacturerRepo.DeleteByStringID(id)
}

func (s *manufacturerService) validateCompanyScopedComponentProduct(
	productID string,
	companyID uint,
) error {
	if productID == "" {
		return nil
	}

	if companyID == 0 {
		return errors.New("invalid company")
	}

	if s.productRepo == nil {
		return errors.New(
			"product repository is required for company validation",
		)
	}

	product, err := s.productRepo.FindByIDAndCompany(
		productID,
		companyID,
	)
	if err != nil || product == nil {
		return fmt.Errorf(
			"component product %s does not belong to your company",
			productID,
		)
	}

	return nil
}

func (s *manufacturerService) validateUserCompany(
	userID uint,
	companyID uint,
) error {
	if userID == 0 {
		return errors.New("invalid authenticated user")
	}
	if companyID == 0 {
		return errors.New("invalid company")
	}

	if s.userRepo == nil {
		return errors.New(
			"user repository is required for company validation",
		)
	}

	user, err := s.userRepo.GetByIDAndCompanyID(
		userID,
		companyID,
	)
	if err != nil || user == nil {
		return errors.New(
			"user does not belong to the company",
		)
	}

	return nil
}

func (s *manufacturerService) CreateForCompany(
	req *input.CreateManufacturerInput,
	userID uint,
	companyID uint,
) (*output.ManufacturerOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if err := s.validateUserCompany(
		userID,
		companyID,
	); err != nil {
		return nil, err
	}

	productGroup, err :=
		s.productGroupRepo.FindByIDAndCompany(
			req.ProductGroupID,
			companyID,
		)
	if err != nil || productGroup == nil {
		return nil, errors.New(
			"product group not found in your company",
		)
	}

	return s.createWithValidatedProductGroup(
		req,
		productGroup,
		userID,
		companyID,
		true,
	)
}

func (s *manufacturerService) GetByIDForCompany(
	id string,
	companyID uint,
) (*output.ManufacturerOutput, error) {
	manufacturer, err :=
		s.manufacturerRepo.FindByStringIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"manufacturer not found",
		)
	}

	return convertManufacturerToOutput(
		manufacturer,
	), nil
}

func (s *manufacturerService) GetAllForCompany(
	companyID uint,
	limit int,
	offset int,
) (*output.ListManufacturersOutput, error) {
	manufacturers, count, err :=
		s.manufacturerRepo.FindAllByCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturersToOutput(
		manufacturers,
		int(count),
	), nil
}

func (s *manufacturerService) GetAllWithFilterForCompany(
	companyID uint,
	productGroupID *string,
	limit int,
	offset int,
) (*output.ListManufacturersOutput, error) {
	if productGroupID != nil &&
		*productGroupID != "" {
		if _, err :=
			s.productGroupRepo.FindByIDAndCompany(
				*productGroupID,
				companyID,
			); err != nil {
			return nil, errors.New(
				"product group not found in your company",
			)
		}
	}

	manufacturers, count, err :=
		s.manufacturerRepo.FindAllWithFilter(
			limit,
			offset,
			&companyID,
			productGroupID,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturersToOutput(
		manufacturers,
		int(count),
	), nil
}

func (s *manufacturerService) GetByProductGroupIDForCompany(
	productGroupID string,
	companyID uint,
) ([]output.ManufacturerOutput, error) {
	if _, err := s.productGroupRepo.FindByIDAndCompany(
		productGroupID,
		companyID,
	); err != nil {
		return nil, errors.New(
			"product group not found in your company",
		)
	}

	manufacturers, err :=
		s.manufacturerRepo.FindByProductGroupIDAndCompany(
			productGroupID,
			companyID,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturerSlice(manufacturers), nil
}

func (s *manufacturerService) UpdateForCompany(
	id string,
	req *input.UpdateManufacturerInput,
	userID uint,
	companyID uint,
) (*output.ManufacturerOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if err := s.validateUserCompany(
		userID,
		companyID,
	); err != nil {
		return nil, err
	}

	manufacturer, err :=
		s.manufacturerRepo.FindByStringIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, errors.New(
			"manufacturer not found",
		)
	}

	oldStatus := manufacturer.Status
	applyManufacturerUpdate(manufacturer, req)

	if shouldConsumeStockOnCompletion(oldStatus, manufacturer.Status) {
		if err := s.consumeComponentStockForManufacturer(
			manufacturer.ProductGroup,
			manufacturer.ID,
			manufacturer.Name,
			manufacturer.Quantity,
			userID,
			companyID,
			true,
		); err != nil {
			return nil, err
		}
	}

	if err := s.manufacturerRepo.UpdateByCompany(
		manufacturer,
		companyID,
	); err != nil {
		return nil, err
	}

	updated, err :=
		s.manufacturerRepo.FindByStringIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, err
	}

	return convertManufacturerToOutput(updated), nil
}

func (s *manufacturerService) DeleteForCompany(
	id string,
	userID uint,
	companyID uint,
) error {
	if err := s.validateUserCompany(
		userID,
		companyID,
	); err != nil {
		return err
	}

	if _, err :=
		s.manufacturerRepo.FindByStringIDAndCompany(
			id,
			companyID,
		); err != nil {
		return errors.New("manufacturer not found")
	}

	// This preserves the existing behavior and does not restore stock.
	return s.manufacturerRepo.DeleteByStringIDAndCompany(
		id,
		companyID,
	)
}

func shouldConsumeStockOnCompletion(oldStatus, newStatus string) bool {
	oldNormalized := strings.ToLower(strings.TrimSpace(oldStatus))
	newNormalized := strings.ToLower(strings.TrimSpace(newStatus))

	return newNormalized == "completed" && oldNormalized != "completed"
}

func (s *manufacturerService) consumeComponentStockForManufacturer(
	productGroup *models.ProductGroup,
	manufacturerID, manufacturerName string,
	quantity float64,
	userID uint,
	companyID uint,
	companyScoped bool,
) error {
	if productGroup == nil {
		return nil
	}

	if s.stockMgmtService == nil {
		return errors.New("stock management service is required")
	}

	if companyScoped {
		for _, component := range productGroup.Components {
			if component.Product == nil || component.Product.IsResource {
				continue
			}

			if err := s.validateCompanyScopedComponentProduct(
				component.ProductID,
				companyID,
			); err != nil {
				return err
			}
		}
	}

	for _, component := range productGroup.Components {
		if component.Product == nil || component.Product.IsResource {
			continue
		}

		requiredQuantity := float64(component.Quantity) * quantity
		notes := fmt.Sprintf(
			"Manufacturing %s: %.0f units of %s",
			manufacturerName,
			quantity,
			component.Product.Name,
		)

		var deductionErr error
		if component.VariantSku != nil && *component.VariantSku != "" {
			deductionErr = s.stockMgmtService.RecordOutboundMovementWithVariant(
				component.ProductID,
				*component.VariantSku,
				"PRODUCTION_USAGE",
				manufacturerID,
				manufacturerName,
				requiredQuantity,
				notes,
				fmt.Sprintf("%d", userID),
			)
		} else {
			deductionErr = s.stockMgmtService.RecordOutboundMovement(
				component.ProductID,
				"PRODUCTION_USAGE",
				manufacturerID,
				manufacturerName,
				requiredQuantity,
				notes,
				fmt.Sprintf("%d", userID),
			)
		}

		if deductionErr != nil {
			return fmt.Errorf(
				"failed to deduct stock for component %s (%s): %v",
				component.Product.Name,
				component.ProductID,
				deductionErr,
			)
		}
	}

	return nil
}

func applyManufacturerUpdate(
	manufacturer *models.Manufacturer,
	req *input.UpdateManufacturerInput,
) {
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
}

func convertManufacturerSlice(
	manufacturers []models.Manufacturer,
) []output.ManufacturerOutput {
	outputs := make(
		[]output.ManufacturerOutput,
		len(manufacturers),
	)

	for index := range manufacturers {
		mapped := convertManufacturerToOutput(
			&manufacturers[index],
		)
		if mapped != nil {
			outputs[index] = *mapped
		}
	}

	return outputs
}

func convertManufacturerToOutput(
	manufacturer *models.Manufacturer,
) *output.ManufacturerOutput {
	if manufacturer == nil {
		return nil
	}

	var productGroupOutput *output.ProductGroupOutput

	if manufacturer.ProductGroup != nil {
		productGroupOutput =
			convertProductGroupToOutput(
				manufacturer.ProductGroup,
			)
	}

	employeeOutputs := make(
		[]output.EmployeeAssignmentOutput,
		len(manufacturer.Employees),
	)

	for index, employee := range manufacturer.Employees {
		employeeOutputs[index] =
			output.EmployeeAssignmentOutput{
				EmployeeID:  employee.EmployeeID,
				ServiceCost: employee.ServiceCost,
				CostType:    employee.CostType,
				Employee:    nil,
			}
	}

	return &output.ManufacturerOutput{
		ID:             manufacturer.ID,
		Name:           manufacturer.Name,
		ProductGroupID: manufacturer.ProductGroupID,
		ProductGroup:   productGroupOutput,
		Quantity:       manufacturer.Quantity,
		Status:         manufacturer.Status,
		Description:    manufacturer.Description,
		Employees:      employeeOutputs,
		CreatedAt:      manufacturer.CreatedAt,
		UpdatedAt:      manufacturer.UpdatedAt,
	}
}

func convertManufacturersToOutput(
	manufacturers []models.Manufacturer,
	count int,
) *output.ListManufacturersOutput {
	return &output.ListManufacturersOutput{
		Manufacturers: convertManufacturerSlice(
			manufacturers,
		),
		TotalCount: count,
	}
}

func convertProductGroupToOutput(
	productGroup *models.ProductGroup,
) *output.ProductGroupOutput {
	if productGroup == nil {
		return nil
	}

	productGroupOutput, err :=
		output.ToProductGroupOutput(productGroup)
	if err != nil {
		return &output.ProductGroupOutput{
			ID:          productGroup.ID,
			Name:        productGroup.Name,
			Description: productGroup.Description,
			IsActive:    productGroup.IsActive,
		}
	}

	return productGroupOutput
}
