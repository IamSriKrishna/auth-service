package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/bbapp-org/auth-service/app/utils"
)

type EmployeeService interface {
	CreateEmployeeWithFile(
		ctx context.Context,
		createdByID uint,
		companyID uint,
		req *input.CreateEmployeeRequest,
		file *multipart.FileHeader,
	) (*output.EmployeeOutput, error)

	GetEmployeeByID(
		ctx context.Context,
		employeeID uint,
		companyID uint,
	) (*output.EmployeeOutput, error)

	GetEmployeesByCompany(
		ctx context.Context,
		companyID uint,
		page int,
		limit int,
	) (*output.PaginatedResponse, error)

	UpdateEmployee(
		ctx context.Context,
		employeeID uint,
		companyID uint,
		req *input.UpdateEmployeeRequest,
	) (*output.EmployeeOutput, error)

	DeleteEmployee(
		ctx context.Context,
		employeeID uint,
		companyID uint,
	) error

	// Existing compatibility method retained for other services.
	GetEmployeesByUser(
		ctx context.Context,
		createdByID uint,
		companyID uint,
		page int,
		limit int,
	) (*output.PaginatedResponse, error)
}

type employeeService struct {
	employeeRepo repo.EmployeeRepository
	userRepo     repo.UserRepository
	cloudinary   *utils.CloudinaryUploader
}

func NewEmployeeService(
	employeeRepo repo.EmployeeRepository,
	userRepo repo.UserRepository,
	cloudinary *utils.CloudinaryUploader,
) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
		cloudinary:   cloudinary,
	}
}

func employeeToOutput(employee *models.Employee) *output.EmployeeOutput {
	return &output.EmployeeOutput{
		ID:            employee.ID,
		Name:          employee.Name,
		Email:         employee.Email,
		Number:        employee.Number,
		Address:       employee.Address,
		EmployeeType:  employee.EmployeeType,
		MonthlySalary: employee.MonthlySalary,
		WeeklySalary:  employee.WeeklySalary,
		SalaryType:    employee.SalaryType,
		DocumentURL:   employee.DocumentURL,
		UserID:        employee.UserID,
		CompanyID:     employee.CompanyID,
		CreatedAt:     employee.CreatedAt,
		UpdatedAt:     employee.UpdatedAt,
	}
}

func (s *employeeService) CreateEmployeeWithFile(
	ctx context.Context,
	createdByID uint,
	companyID uint,
	req *input.CreateEmployeeRequest,
	file *multipart.FileHeader,
) (*output.EmployeeOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if file == nil {
		return nil, errors.New("document file is required")
	}

	if createdByID == 0 {
		return nil, errors.New("invalid creating user")
	}

	if companyID == 0 {
		return nil, errors.New("invalid company")
	}

	employee := &models.Employee{
		Name:          req.Name,
		Email:         req.Email,
		Number:        req.Number,
		Address:       req.Address,
		EmployeeType:  req.EmployeeType,
		MonthlySalary: req.MonthlySalary,
		WeeklySalary:  req.WeeklySalary,
		SalaryType:    req.SalaryType,
		UserID:        createdByID,
		CompanyID:     companyID,
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	if s.cloudinary != nil {
		fileReader, err := file.Open()
		if err != nil {
			_ = s.employeeRepo.DeleteByIDAndCompanyID(employee.ID, companyID)
			return nil, errors.New("failed to open document file: " + err.Error())
		}
		defer fileReader.Close()

		documentURL, err := s.cloudinary.UploadEmployeeDocumentFromReader(
			ctx,
			fileReader,
			file.Filename,
			employee.ID,
		)
		if err != nil {
			_ = s.employeeRepo.DeleteByIDAndCompanyID(employee.ID, companyID)
			return nil, errors.New("failed to upload document: " + err.Error())
		}

		employee.DocumentURL = documentURL

		if err := s.employeeRepo.UpdateByCompanyID(employee, companyID); err != nil {
			return nil, err
		}
	}

	return employeeToOutput(employee), nil
}

func (s *employeeService) GetEmployeeByID(
	ctx context.Context,
	employeeID uint,
	companyID uint,
) (*output.EmployeeOutput, error) {
	_ = ctx

	employee, err := s.employeeRepo.GetByIDAndCompanyID(
		employeeID,
		companyID,
	)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	return employeeToOutput(employee), nil
}

func (s *employeeService) GetEmployeesByCompany(
	ctx context.Context,
	companyID uint,
	page int,
	limit int,
) (*output.PaginatedResponse, error) {
	_ = ctx

	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	employees, total, err := s.employeeRepo.GetByCompany(
		companyID,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	employeeOutputs := make([]output.EmployeeListOutput, 0, len(employees))

	for _, employee := range employees {
		employeeOutputs = append(employeeOutputs, output.EmployeeListOutput{
			ID:            employee.ID,
			Name:          employee.Name,
			Email:         employee.Email,
			Number:        employee.Number,
			Address:       employee.Address,
			EmployeeType:  employee.EmployeeType,
			MonthlySalary: employee.MonthlySalary,
			WeeklySalary:  employee.WeeklySalary,
			SalaryType:    employee.SalaryType,
			DocumentURL:   employee.DocumentURL,
			UserID:        employee.UserID,
			CompanyID:     employee.CompanyID,
			CreatedAt:     employee.CreatedAt,
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &output.PaginatedResponse{
		Success: true,
		Data:    employeeOutputs,
		Meta: output.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       int(total),
			TotalPages:  totalPages,
		},
	}, nil
}

// Compatibility method retained.
// It now returns all employees in the company instead of only one user's employees.
func (s *employeeService) GetEmployeesByUser(
	ctx context.Context,
	createdByID uint,
	companyID uint,
	page int,
	limit int,
) (*output.PaginatedResponse, error) {
	_ = createdByID

	return s.GetEmployeesByCompany(
		ctx,
		companyID,
		page,
		limit,
	)
}

func (s *employeeService) UpdateEmployee(
	ctx context.Context,
	employeeID uint,
	companyID uint,
	req *input.UpdateEmployeeRequest,
) (*output.EmployeeOutput, error) {
	_ = ctx

	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	employee, err := s.employeeRepo.GetByIDAndCompanyID(
		employeeID,
		companyID,
	)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if req.Name != nil {
		employee.Name = *req.Name
	}
	if req.Email != nil {
		employee.Email = *req.Email
	}
	if req.Number != nil {
		employee.Number = *req.Number
	}
	if req.Address != nil {
		employee.Address = *req.Address
	}
	if req.EmployeeType != nil {
		employee.EmployeeType = *req.EmployeeType
	}
	if req.MonthlySalary != nil {
		employee.MonthlySalary = *req.MonthlySalary
	}
	if req.WeeklySalary != nil {
		employee.WeeklySalary = *req.WeeklySalary
	}
	if req.SalaryType != nil {
		employee.SalaryType = *req.SalaryType
	}

	if err := s.employeeRepo.UpdateByCompanyID(
		employee,
		companyID,
	); err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	return employeeToOutput(employee), nil
}

func (s *employeeService) DeleteEmployee(
	ctx context.Context,
	employeeID uint,
	companyID uint,
) error {
	_ = ctx

	if err := s.employeeRepo.DeleteByIDAndCompanyID(
		employeeID,
		companyID,
	); err != nil {
		return errors.New("employee not found")
	}

	return nil
}