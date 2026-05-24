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

type ProductConversionService interface {
	// Conversion Rule Management
	CreateConversion(conversionInput *input.CreateProductConversionInput, userID string, companyID uint, userName string, companyName string) (*output.ProductConversionOutput, error)
	UpdateConversion(conversionID string, updateInput *input.UpdateProductConversionInput, userID string, userName string) (*output.ProductConversionOutput, error)
	DeleteConversion(conversionID string) error
	GetConversion(conversionID string) (*output.ProductConversionOutput, error)
	ListConversions(offset, limit int) (*output.ProductConversionListOutput, error)
	ListActiveConversions(offset, limit int) (*output.ProductConversionListOutput, error)
	GetConversionsByRawProduct(rawProductID string, offset, limit int) (*output.ProductConversionListOutput, error)
	GetConversionsByFinishedProduct(finishedProductID string, offset, limit int) (*output.ProductConversionListOutput, error)

	// Conversion Execution
	ExecuteConversion(conversionInput *input.CreateProductConversionRecordInput, userID string, companyID uint, userName string, companyName string) (*output.ConversionExecutionOutput, error)

	// Conversion Record Management
	GetConversionRecord(recordID string) (*output.ProductConversionRecordOutput, error)
	ListConversionRecords(offset, limit int) (*output.ProductConversionRecordListOutput, error)
	ListConversionRecordsByRule(conversionID string, offset, limit int) (*output.ProductConversionRecordListOutput, error)
	ListConversionRecordsByDateRange(fromDate, toDate time.Time, offset, limit int) (*output.ProductConversionRecordListOutput, error)
}

type productConversionService struct {
	conversionRepo     repo.ProductConversionRepository
	recordRepo         repo.ProductConversionRecordRepository
	productRepo        repo.ProductRepository
	stockManagementSvc StockManagementService
	variantStockMgmt   VariantStockManagementService
}

func NewProductConversionService(
	conversionRepo repo.ProductConversionRepository,
	recordRepo repo.ProductConversionRecordRepository,
	productRepo repo.ProductRepository,
	stockManagementSvc StockManagementService,
	variantStockMgmt VariantStockManagementService,
) ProductConversionService {
	return &productConversionService{
		conversionRepo:     conversionRepo,
		recordRepo:         recordRepo,
		productRepo:        productRepo,
		stockManagementSvc: stockManagementSvc,
		variantStockMgmt:   variantStockMgmt,
	}
}

// CreateConversion creates a new conversion rule
func (s *productConversionService) CreateConversion(
	conversionInput *input.CreateProductConversionInput,
	userID string,
	companyID uint,
	userName string,
	companyName string,
) (*output.ProductConversionOutput, error) {
	// Validate raw product exists
	rawProduct, err := s.productRepo.FindByID(conversionInput.RawProductID)
	if err != nil {
		return nil, fmt.Errorf("error fetching raw product: %w", err)
	}
	if rawProduct == nil {
		return nil, errors.New("raw product not found")
	}

	// Validate finished product exists
	finishedProduct, err := s.productRepo.FindByID(conversionInput.FinishedProductID)
	if err != nil {
		return nil, fmt.Errorf("error fetching finished product: %w", err)
	}
	if finishedProduct == nil {
		return nil, errors.New("finished product not found")
	}

	// Check if conversion already exists
	existing, err := s.conversionRepo.GetByProductPair(conversionInput.RawProductID, conversionInput.FinishedProductID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing conversion: %w", err)
	}
	if existing != nil {
		return nil, errors.New("conversion rule already exists for this product pair")
	}

	conversion := &models.ProductConversion{
		ID:                   uuid.New().String(),
		RawProductID:         conversionInput.RawProductID,
		RawProductName:       rawProduct.Name,
		RawProductSpec:       rawProduct.RawSpecification,
		FinishedProductID:    conversionInput.FinishedProductID,
		FinishedProductName:  finishedProduct.Name,
		FinishedVariantSKU:   conversionInput.FinishedVariantSKU,
		ConversionRatio:      conversionInput.ConversionRatio,
		LossPercentage:       conversionInput.LossPercentage,
		IsActive:             conversionInput.IsActive,
		Notes:                conversionInput.Notes,
		CreatedBy:            userID,
		CreatedByUserName:    userName,
		CreatedByCompanyID:   companyID,
		CreatedByCompanyName: companyName,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.conversionRepo.Create(conversion); err != nil {
		return nil, fmt.Errorf("error creating conversion: %w", err)
	}

	return s.mapConversionToOutput(conversion), nil
}

// UpdateConversion updates a conversion rule
func (s *productConversionService) UpdateConversion(
	conversionID string,
	updateInput *input.UpdateProductConversionInput,
	userID string,
	userName string,
) (*output.ProductConversionOutput, error) {
	conversion, err := s.conversionRepo.GetByID(conversionID)
	if err != nil {
		return nil, fmt.Errorf("error fetching conversion: %w", err)
	}
	if conversion == nil {
		return nil, errors.New("conversion not found")
	}

	if updateInput.ConversionRatio != nil {
		conversion.ConversionRatio = *updateInput.ConversionRatio
	}
	if updateInput.LossPercentage != nil {
		conversion.LossPercentage = *updateInput.LossPercentage
	}
	if updateInput.IsActive != nil {
		conversion.IsActive = *updateInput.IsActive
	}
	if updateInput.FinishedVariantSKU != nil {
		conversion.FinishedVariantSKU = *updateInput.FinishedVariantSKU
	}
	if updateInput.Notes != nil {
		conversion.Notes = *updateInput.Notes
	}

	conversion.UpdatedBy = userID
	conversion.UpdatedByUserName = userName
	conversion.UpdatedAt = time.Now()

	if err := s.conversionRepo.Update(conversion); err != nil {
		return nil, fmt.Errorf("error updating conversion: %w", err)
	}

	return s.mapConversionToOutput(conversion), nil
}

// DeleteConversion deletes a conversion rule
func (s *productConversionService) DeleteConversion(conversionID string) error {
	conversion, err := s.conversionRepo.GetByID(conversionID)
	if err != nil {
		return fmt.Errorf("error fetching conversion: %w", err)
	}
	if conversion == nil {
		return errors.New("conversion not found")
	}

	return s.conversionRepo.Delete(conversionID)
}

// GetConversion retrieves a conversion rule
func (s *productConversionService) GetConversion(conversionID string) (*output.ProductConversionOutput, error) {
	conversion, err := s.conversionRepo.GetByID(conversionID)
	if err != nil {
		return nil, fmt.Errorf("error fetching conversion: %w", err)
	}
	if conversion == nil {
		return nil, errors.New("conversion not found")
	}

	return s.mapConversionToOutput(conversion), nil
}

// ListConversions lists all conversions
func (s *productConversionService) ListConversions(offset, limit int) (*output.ProductConversionListOutput, error) {
	conversions, total, err := s.conversionRepo.GetAll(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching conversions: %w", err)
	}

	outputs := make([]output.ProductConversionOutput, len(conversions))
	for i, c := range conversions {
		outputs[i] = *s.mapConversionToOutput(&c)
	}

	return &output.ProductConversionListOutput{
		Conversions: outputs,
		Total:       total,
		Page:        offset/limit + 1,
		Limit:       limit,
	}, nil
}

// ListActiveConversions lists all active conversions
func (s *productConversionService) ListActiveConversions(offset, limit int) (*output.ProductConversionListOutput, error) {
	conversions, total, err := s.conversionRepo.GetActiveConversions(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching active conversions: %w", err)
	}

	outputs := make([]output.ProductConversionOutput, len(conversions))
	for i, c := range conversions {
		outputs[i] = *s.mapConversionToOutput(&c)
	}

	return &output.ProductConversionListOutput{
		Conversions: outputs,
		Total:       total,
		Page:        offset/limit + 1,
		Limit:       limit,
	}, nil
}

// GetConversionsByRawProduct lists conversions for a raw product
func (s *productConversionService) GetConversionsByRawProduct(rawProductID string, offset, limit int) (*output.ProductConversionListOutput, error) {
	conversions, total, err := s.conversionRepo.GetByRawProductID(rawProductID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching conversions: %w", err)
	}

	outputs := make([]output.ProductConversionOutput, len(conversions))
	for i, c := range conversions {
		outputs[i] = *s.mapConversionToOutput(&c)
	}

	return &output.ProductConversionListOutput{
		Conversions: outputs,
		Total:       total,
		Page:        offset/limit + 1,
		Limit:       limit,
	}, nil
}

// GetConversionsByFinishedProduct lists conversions for a finished product
func (s *productConversionService) GetConversionsByFinishedProduct(finishedProductID string, offset, limit int) (*output.ProductConversionListOutput, error) {
	conversions, total, err := s.conversionRepo.GetByFinishedProductID(finishedProductID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching conversions: %w", err)
	}

	outputs := make([]output.ProductConversionOutput, len(conversions))
	for i, c := range conversions {
		outputs[i] = *s.mapConversionToOutput(&c)
	}

	return &output.ProductConversionListOutput{
		Conversions: outputs,
		Total:       total,
		Page:        offset/limit + 1,
		Limit:       limit,
	}, nil
}

// ExecuteConversion executes a product conversion and updates stock
func (s *productConversionService) ExecuteConversion(
	conversionInput *input.CreateProductConversionRecordInput,
	userID string,
	companyID uint,
	userName string,
	companyName string,
) (*output.ConversionExecutionOutput, error) {
	// Get conversion rule
	conversion, err := s.conversionRepo.GetByID(conversionInput.ConversionID)
	if err != nil {
		return nil, fmt.Errorf("error fetching conversion: %w", err)
	}
	if conversion == nil {
		return nil, errors.New("conversion not found")
	}

	if !conversion.IsActive {
		return nil, errors.New("conversion rule is not active")
	}

	// Calculate produced quantity
	lossRatio := conversion.LossPercentage / 100.0
	producedQuantity := (conversionInput.RawQuantityUsed / conversion.ConversionRatio) * (1 - lossRatio)
	lossQuantity := conversionInput.RawQuantityUsed * lossRatio

	// Set conversion date
	conversionDate := time.Now()
	if conversionInput.ConversionDate != nil {
		conversionDate = *conversionInput.ConversionDate
	}

	// Determine which variant SKU to use (if any)
	// Priority: input variant > conversion rule variant > product's base/first variant (if has variants) > product-level only
	variantSKU := conversionInput.FinishedVariantSKU
	if variantSKU == "" {
		variantSKU = conversion.FinishedVariantSKU
	}

	// If still no variant SKU, check if finished product has variants
	// If it does, use the base SKU or first variant's SKU
	if variantSKU == "" {
		finishedProduct, err := s.productRepo.FindByID(conversion.FinishedProductID)
		if err == nil && finishedProduct != nil {
			// Check if product has variants
			if len(finishedProduct.ProductDetails.ProductVariants) > 0 {
				// Use base SKU if available, otherwise use first variant's SKU
				if finishedProduct.ProductDetails.BaseSKU != "" {
					variantSKU = finishedProduct.ProductDetails.BaseSKU
				} else {
					variantSKU = finishedProduct.ProductDetails.ProductVariants[0].SKU
				}
			}
		}
	}

	// Create conversion record
	record := &models.ProductConversionRecord{
		ID:                       uuid.New().String(),
		ConversionID:             conversionInput.ConversionID,
		RawProductID:             conversion.RawProductID,
		RawProductName:           conversion.RawProductName,
		RawQuantityUsed:          conversionInput.RawQuantityUsed,
		FinishedProductID:        conversion.FinishedProductID,
		FinishedProductName:      conversion.FinishedProductName,
		FinishedVariantSKU:       variantSKU,
		FinishedQuantityProduced: producedQuantity,
		LossQuantity:             lossQuantity,
		ConversionDate:           conversionDate,
		Status:                   "COMPLETED",
		Notes:                    conversionInput.Notes,
		CreatedBy:                userID,
		CreatedByUserName:        userName,
		CreatedByCompanyID:       companyID,
		CreatedByCompanyName:     companyName,
		CreatedAt:                time.Now(),
	}

	if err := s.recordRepo.Create(record); err != nil {
		return nil, fmt.Errorf("error creating conversion record: %w", err)
	}

	// Update stock if ExecuteConversion is true
	if conversionInput.ExecuteConversion {
		// Deduct from raw product stock using outbound movement
		if err := s.stockManagementSvc.RecordOutboundMovement(
			conversion.RawProductID,
			"PRODUCTION_USAGE",
			record.ID,
			"",
			conversionInput.RawQuantityUsed,
			fmt.Sprintf("Conversion: %s → %s", conversion.RawProductName, conversion.FinishedProductName),
			userID,
		); err != nil {
			// Mark record as failed and return error
			record.Status = "FAILED"
			s.recordRepo.Update(record)
			return nil, fmt.Errorf("error deducting raw product stock: %w", err)
		}

		// Add to finished product stock - use variant if specified, otherwise product-level
		if variantSKU != "" {
			// Get finished product to find variant details
			finishedProduct, err := s.productRepo.FindByID(conversion.FinishedProductID)
			if err == nil && finishedProduct != nil {
				// Try to find the variant definition in the product
				variantName := ""
				sellingPrice := 0.0
				costPrice := 0.0

				for _, v := range finishedProduct.ProductDetails.ProductVariants {
					if v.SKU == variantSKU {
						variantName = v.VariantName
						sellingPrice = v.SellingPrice
						costPrice = v.CostPrice
						break
					}
				}

				// Initialize variant stock if it doesn't exist
				if variantName != "" {
					s.variantStockMgmt.InitializeVariantStock(
						conversion.FinishedProductID,
						variantSKU,
						variantName,
						conversion.FinishedProductName,
						sellingPrice,
						costPrice,
					)
				}
			}

			// Add stock to specific variant using stock adjustment
			if err := s.variantStockMgmt.RecordStockAdjustment(
				variantSKU,
				producedQuantity,
				"in",
				fmt.Sprintf("Conversion from: %s (Qty: %v, Conversion ID: %s)", conversion.RawProductName, conversionInput.RawQuantityUsed, record.ID),
				userID,
			); err != nil {
				// Mark record as failed
				record.Status = "FAILED"
				s.recordRepo.Update(record)
				return nil, fmt.Errorf("error adding stock to variant %s: %w", variantSKU, err)
			}
		} else {
			// Add stock to product level using inbound movement
			if err := s.stockManagementSvc.RecordInboundMovement(
				conversion.FinishedProductID,
				"PRODUCTION_USAGE",
				record.ID,
				"",
				producedQuantity,
				0, // rate is not applicable for conversions
				fmt.Sprintf("Conversion from: %s (Qty: %v)", conversion.RawProductName, conversionInput.RawQuantityUsed),
				userID,
			); err != nil {
				// Mark record as failed
				record.Status = "FAILED"
				s.recordRepo.Update(record)
				return nil, fmt.Errorf("error adding finished product stock: %w", err)
			}
		}
	}

	return &output.ConversionExecutionOutput{
		RecordID:                 record.ID,
		Status:                   record.Status,
		RawProductName:           record.RawProductName,
		RawQuantityUsed:          record.RawQuantityUsed,
		FinishedProductName:      record.FinishedProductName,
		FinishedVariantSKU:       record.FinishedVariantSKU,
		FinishedQuantityProduced: record.FinishedQuantityProduced,
		LossQuantity:             record.LossQuantity,
		Message: fmt.Sprintf("Conversion executed successfully. Raw: %v used, Finished: %v produced to %s", conversionInput.RawQuantityUsed, producedQuantity, func() string {
			if record.FinishedVariantSKU != "" {
				return "variant " + record.FinishedVariantSKU
			} else {
				return "product level"
			}
		}()),
	}, nil
}

// GetConversionRecord retrieves a conversion record
func (s *productConversionService) GetConversionRecord(recordID string) (*output.ProductConversionRecordOutput, error) {
	record, err := s.recordRepo.GetByID(recordID)
	if err != nil {
		return nil, fmt.Errorf("error fetching record: %w", err)
	}
	if record == nil {
		return nil, errors.New("record not found")
	}

	return s.mapRecordToOutput(record), nil
}

// ListConversionRecords lists all conversion records
func (s *productConversionService) ListConversionRecords(offset, limit int) (*output.ProductConversionRecordListOutput, error) {
	records, total, err := s.recordRepo.GetAll(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching records: %w", err)
	}

	outputs := make([]output.ProductConversionRecordOutput, len(records))
	for i, r := range records {
		outputs[i] = *s.mapRecordToOutput(&r)
	}

	return &output.ProductConversionRecordListOutput{
		Records: outputs,
		Total:   total,
		Page:    offset/limit + 1,
		Limit:   limit,
	}, nil
}

// ListConversionRecordsByRule lists conversion records for a specific rule
func (s *productConversionService) ListConversionRecordsByRule(conversionID string, offset, limit int) (*output.ProductConversionRecordListOutput, error) {
	records, total, err := s.recordRepo.GetByConversionID(conversionID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching records: %w", err)
	}

	outputs := make([]output.ProductConversionRecordOutput, len(records))
	for i, r := range records {
		outputs[i] = *s.mapRecordToOutput(&r)
	}

	return &output.ProductConversionRecordListOutput{
		Records: outputs,
		Total:   total,
		Page:    offset/limit + 1,
		Limit:   limit,
	}, nil
}

// ListConversionRecordsByDateRange lists conversion records within a date range
func (s *productConversionService) ListConversionRecordsByDateRange(fromDate, toDate time.Time, offset, limit int) (*output.ProductConversionRecordListOutput, error) {
	records, total, err := s.recordRepo.GetByDateRange(fromDate, toDate, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching records: %w", err)
	}

	outputs := make([]output.ProductConversionRecordOutput, len(records))
	for i, r := range records {
		outputs[i] = *s.mapRecordToOutput(&r)
	}

	return &output.ProductConversionRecordListOutput{
		Records: outputs,
		Total:   total,
		Page:    offset/limit + 1,
		Limit:   limit,
	}, nil
}

// Helper methods
func (s *productConversionService) mapConversionToOutput(conversion *models.ProductConversion) *output.ProductConversionOutput {
	return &output.ProductConversionOutput{
		ID:                   conversion.ID,
		RawProductID:         conversion.RawProductID,
		RawProductName:       conversion.RawProductName,
		RawProductSpec:       conversion.RawProductSpec,
		FinishedProductID:    conversion.FinishedProductID,
		FinishedProductName:  conversion.FinishedProductName,
		FinishedProductSpec:  conversion.FinishedProductSpec,
		FinishedVariantSKU:   conversion.FinishedVariantSKU,
		ConversionRatio:      conversion.ConversionRatio,
		LossPercentage:       conversion.LossPercentage,
		IsActive:             conversion.IsActive,
		Notes:                conversion.Notes,
		CreatedBy:            conversion.CreatedBy,
		CreatedByUserName:    conversion.CreatedByUserName,
		CreatedByCompanyID:   conversion.CreatedByCompanyID,
		CreatedByCompanyName: conversion.CreatedByCompanyName,
		UpdatedBy:            conversion.UpdatedBy,
		UpdatedByUserName:    conversion.UpdatedByUserName,
		CreatedAt:            conversion.CreatedAt,
		UpdatedAt:            conversion.UpdatedAt,
	}
}

func (s *productConversionService) mapRecordToOutput(record *models.ProductConversionRecord) *output.ProductConversionRecordOutput {
	return &output.ProductConversionRecordOutput{
		ID:                       record.ID,
		ConversionID:             record.ConversionID,
		RawProductID:             record.RawProductID,
		RawProductName:           record.RawProductName,
		RawQuantityUsed:          record.RawQuantityUsed,
		FinishedProductID:        record.FinishedProductID,
		FinishedProductName:      record.FinishedProductName,
		FinishedVariantSKU:       record.FinishedVariantSKU,
		FinishedQuantityProduced: record.FinishedQuantityProduced,
		LossQuantity:             record.LossQuantity,
		ConversionDate:           record.ConversionDate,
		Status:                   record.Status,
		Notes:                    record.Notes,
		CreatedBy:                record.CreatedBy,
		CreatedByUserName:        record.CreatedByUserName,
		CreatedByCompanyID:       record.CreatedByCompanyID,
		CreatedByCompanyName:     record.CreatedByCompanyName,
		CreatedAt:                record.CreatedAt,
	}
}
