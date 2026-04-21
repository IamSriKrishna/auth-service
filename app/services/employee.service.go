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
	CreateEmployeeWithFile(ctx context.Context, createdByID, companyID uint, req *input.CreateEmployeeRequest, file *multipart.FileHeader) (*output.EmployeeOutput, error)
	GetEmployeeByID(ctx context.Context, employeeID, createdByID uint) (*output.EmployeeOutput, error)
	GetEmployeesByUser(ctx context.Context, createdByID, companyID uint, page, limit int) (*output.PaginatedResponse, error)
	UpdateEmployee(ctx context.Context, employeeID, createdByID uint, req *input.UpdateEmployeeRequest) (*output.EmployeeOutput, error)
	DeleteEmployee(ctx context.Context, employeeID, createdByID uint) error
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

func (s *employeeService) CreateEmployeeWithFile(ctx context.Context, createdByID, companyID uint, req *input.CreateEmployeeRequest, file *multipart.FileHeader) (*output.EmployeeOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if file == nil {
		return nil, errors.New("document file is required")
	}

	// Create employee first without document URL
	employee := &models.Employee{
		Name:          req.Name,
		Email:         req.Email,
		Number:        req.Number,
		Address:       req.Address,
		EmployeeType:  req.EmployeeType,
		MonthlySalary: req.MonthlySalary,
		UserID:        createdByID,
		CompanyID:     companyID,
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	// Upload document to Cloudinary
	if s.cloudinary != nil {
		fileReader, err := file.Open()
		if err != nil {
			return nil, errors.New("failed to open document file: " + err.Error())
		}
		defer fileReader.Close()

		fmt.Printf("DEBUG: Uploading document %s for employee %d\n", file.Filename, employee.ID)
		documentURL, err := s.cloudinary.UploadEmployeeDocumentFromReader(ctx, fileReader, file.Filename, employee.ID)
		if err != nil {
			return nil, errors.New("failed to upload document: " + err.Error())
		}
		fmt.Printf("DEBUG: Document uploaded successfully, URL: %s\n", documentURL)
		employee.DocumentURL = documentURL

		if err := s.employeeRepo.Update(employee); err != nil {
			return nil, err
		}
	}

	return &output.EmployeeOutput{
		ID:            employee.ID,
		Name:          employee.Name,
		Email:         employee.Email,
		Number:        employee.Number,
		Address:       employee.Address,
		EmployeeType:  employee.EmployeeType,
		MonthlySalary: employee.MonthlySalary,
		DocumentURL:   employee.DocumentURL,
		UserID:        employee.UserID,
		CompanyID:     employee.CompanyID,
		CreatedAt:     employee.CreatedAt,
		UpdatedAt:     employee.UpdatedAt,
	}, nil
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, employeeID, createdByID uint) (*output.EmployeeOutput, error) {
	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Verify that the employee belongs to the user
	if employee.UserID != createdByID {
		return nil, errors.New("unauthorized access to employee")
	}

	return &output.EmployeeOutput{
		ID:            employee.ID,
		Name:          employee.Name,
		Email:         employee.Email,
		Number:        employee.Number,
		Address:       employee.Address,
		EmployeeType:  employee.EmployeeType,
		MonthlySalary: employee.MonthlySalary,
		DocumentURL:   employee.DocumentURL,
		UserID:        employee.UserID,
		CompanyID:     employee.CompanyID,
		CreatedAt:     employee.CreatedAt,
		UpdatedAt:     employee.UpdatedAt,
	}, nil
}

func (s *employeeService) GetEmployeesByUser(ctx context.Context, createdByID, companyID uint, page, limit int) (*output.PaginatedResponse, error) {
	offset := (page - 1) * limit
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	employees, total, err := s.employeeRepo.GetByCompanyAndUser(companyID, createdByID, offset, limit)
	if err != nil {
		return nil, err
	}

	var employeeOutputs []output.EmployeeListOutput
	for _, emp := range employees {
		employeeOutputs = append(employeeOutputs, output.EmployeeListOutput{
			ID:            emp.ID,
			Name:          emp.Name,
			Email:         emp.Email,
			Number:        emp.Number,
			Address:       emp.Address,
			EmployeeType:  emp.EmployeeType,
			MonthlySalary: emp.MonthlySalary,
			DocumentURL:   emp.DocumentURL,
			UserID:        emp.UserID,
			CompanyID:     emp.CompanyID,
			CreatedAt:     emp.CreatedAt,
		})
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &output.PaginatedResponse{
		Success: true,
		Data:    employeeOutputs,
		Meta: output.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       int(total),
			TotalPages:  int(totalPages),
		},
	}, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, employeeID, createdByID uint, req *input.UpdateEmployeeRequest) (*output.EmployeeOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Verify that the employee belongs to the user
	if employee.UserID != createdByID {
		return nil, errors.New("unauthorized access to employee")
	}

	// Update fields if provided
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

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, err
	}

	return &output.EmployeeOutput{
		ID:            employee.ID,
		Name:          employee.Name,
		Email:         employee.Email,
		Number:        employee.Number,
		Address:       employee.Address,
		EmployeeType:  employee.EmployeeType,
		MonthlySalary: employee.MonthlySalary,
		DocumentURL:   employee.DocumentURL,
		UserID:        employee.UserID,
		CompanyID:     employee.CompanyID,
		CreatedAt:     employee.CreatedAt,
		UpdatedAt:     employee.UpdatedAt,
	}, nil
}

func (s *employeeService) DeleteEmployee(ctx context.Context, employeeID, createdByID uint) error {
	employee, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	// Verify that the employee belongs to the user
	if employee.UserID != createdByID {
		return errors.New("unauthorized access to employee")
	}

	return s.employeeRepo.Delete(employeeID)
}
