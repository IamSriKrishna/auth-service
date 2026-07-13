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

	// Company-scoped methods.
	CreateConversionForCompany(*input.CreateProductConversionInput, string, uint, string, string) (*output.ProductConversionOutput, error)
	UpdateConversionForCompany(string, *input.UpdateProductConversionInput, string, string, uint) (*output.ProductConversionOutput, error)
	DeleteConversionForCompany(string, uint) error
	GetConversionForCompany(string, uint) (*output.ProductConversionOutput, error)
	ListConversionsForCompany(uint, int, int) (*output.ProductConversionListOutput, error)
	ListActiveConversionsForCompany(uint, int, int) (*output.ProductConversionListOutput, error)
	GetConversionsByRawProductForCompany(string, uint, int, int) (*output.ProductConversionListOutput, error)
	GetConversionsByFinishedProductForCompany(string, uint, int, int) (*output.ProductConversionListOutput, error)
	ExecuteConversionForCompany(*input.CreateProductConversionRecordInput, string, uint, string, string) (*output.ConversionExecutionOutput, error)
	GetConversionRecordForCompany(string, uint) (*output.ProductConversionRecordOutput, error)
	ListConversionRecordsForCompany(uint, int, int) (*output.ProductConversionRecordListOutput, error)
	ListConversionRecordsByRuleForCompany(string, uint, int, int) (*output.ProductConversionRecordListOutput, error)
	ListConversionRecordsByDateRangeForCompany(time.Time, time.Time, uint, int, int) (*output.ProductConversionRecordListOutput, error)
}

type productConversionService struct {
	conversionRepo     repo.ProductConversionRepository
	recordRepo         repo.ProductConversionRecordRepository
	bagUsageRepo       repo.ConversionRecordBagUsageRepository
	productRepo        repo.ProductRepository
	stockManagementSvc StockManagementService
	variantStockMgmt   VariantStockManagementService
	rawMaterialBagSvc  RawMaterialBagService
	userRepo           repo.UserRepository
}

func NewProductConversionService(
	conversionRepo repo.ProductConversionRepository,
	recordRepo repo.ProductConversionRecordRepository,
	bagUsageRepo repo.ConversionRecordBagUsageRepository,
	productRepo repo.ProductRepository,
	stockManagementSvc StockManagementService,
	variantStockMgmt VariantStockManagementService,
	rawMaterialBagSvc RawMaterialBagService,
	userRepo repo.UserRepository,
) ProductConversionService {
	return &productConversionService{
		conversionRepo:     conversionRepo,
		recordRepo:         recordRepo,
		bagUsageRepo:       bagUsageRepo,
		productRepo:        productRepo,
		stockManagementSvc: stockManagementSvc,
		variantStockMgmt:   variantStockMgmt,
		rawMaterialBagSvc:  rawMaterialBagSvc,
		userRepo:           userRepo,
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
func (s *productConversionService) ExecuteConversion(
	conversionInput *input.CreateProductConversionRecordInput,
	userID string,
	companyID uint,
	userName string,
	companyName string,
) (*output.ConversionExecutionOutput, error) {
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

	rawProduct, err := s.productRepo.FindByID(conversion.RawProductID)
	if err != nil {
		return nil, fmt.Errorf("raw product not found: %w", err)
	}

	if rawProduct.RequiredGramPerUnit <= 0 {
		return nil, errors.New("required_gram_per_unit must be greater than 0")
	}

	var rawQuantityUsed float64
	var producedQuantity float64
	var bagUsageDetails []models.ConversionRecordBagUsage

	if len(conversionInput.RawMaterialBags) > 0 {
		if s.rawMaterialBagSvc == nil {
			return nil, errors.New("raw material bag service not configured")
		}

		for _, bagInput := range conversionInput.RawMaterialBags {
			if bagInput.FinishedQuantity <= 0 {
				return nil, errors.New("finished_quantity is required when using raw material bags")
			}

			requiredGrams := bagInput.FinishedQuantity * rawProduct.RequiredGramPerUnit
			requiredKg := requiredGrams / 1000

			// Get bag details before using it
			bagRepo := s.recordRepo.GetDB().Model(&models.RawMaterialBag{})
			var bag models.RawMaterialBag
			if err := bagRepo.Where("id = ?", bagInput.BagID).First(&bag).Error; err != nil {
				return nil, fmt.Errorf("failed to fetch bag details: %w", err)
			}

			_, err := s.rawMaterialBagSvc.UseBags(
				conversion.RawProductID,
				[]input.UseRawMaterialBagInput{
					{
						BagID:      bagInput.BagID,
						QuantityKg: requiredKg,
					},
				},
			)
			if err != nil {
				return nil, err
			}

			// Create bag usage detail for tracking
			bagUsageDetail := models.ConversionRecordBagUsage{
				ID:             uuid.New().String(),
				BagID:          bagInput.BagID,
				BagNumber:      bag.BagNumber,
				ProductID:      bag.ProductID,
				ProductName:    bag.ProductName,
				QuantityUsedKg: requiredKg,
				CreatedAt:      time.Now(),
			}
			bagUsageDetails = append(bagUsageDetails, bagUsageDetail)

			rawQuantityUsed += requiredGrams
			producedQuantity += bagInput.FinishedQuantity
		}
	} else {
		rawQuantityUsed = conversionInput.RawQuantityUsed

		if rawQuantityUsed <= 0 {
			return nil, errors.New("raw quantity used is required")
		}

		producedQuantity = rawQuantityUsed / rawProduct.RequiredGramPerUnit
	}

	lossRatio := conversion.LossPercentage / 100
	lossQuantity := rawQuantityUsed * lossRatio
	producedQuantity = producedQuantity * (1 - lossRatio)

	conversionDate := time.Now()
	if conversionInput.ConversionDate != nil {
		conversionDate = *conversionInput.ConversionDate
	}

	variantSKU := conversionInput.FinishedVariantSKU
	if variantSKU == "" {
		variantSKU = conversion.FinishedVariantSKU
	}

	if variantSKU == "" {
		finishedProduct, err := s.productRepo.FindByID(conversion.FinishedProductID)
		if err == nil && finishedProduct != nil {
			if len(finishedProduct.ProductDetails.ProductVariants) > 0 {
				if finishedProduct.ProductDetails.BaseSKU != "" {
					variantSKU = finishedProduct.ProductDetails.BaseSKU
				} else {
					variantSKU = finishedProduct.ProductDetails.ProductVariants[0].SKU
				}
			}
		}
	}

	record := &models.ProductConversionRecord{
		ID:                       uuid.New().String(),
		ConversionID:             conversionInput.ConversionID,
		RawProductID:             conversion.RawProductID,
		RawProductName:           conversion.RawProductName,
		RawQuantityUsed:          rawQuantityUsed,
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

	// Save bag usage details
	for i := range bagUsageDetails {
		bagUsageDetails[i].ConversionRecordID = record.ID
		if err := s.bagUsageRepo.Create(&bagUsageDetails[i]); err != nil {
			// Log error but don't fail the conversion - tracking is secondary
			fmt.Printf("warning: failed to save bag usage detail: %v\n", err)
		}
	}

	if conversionInput.ExecuteConversion {
		if err := s.stockManagementSvc.RecordOutboundMovement(
			conversion.RawProductID,
			"PRODUCTION_USAGE",
			record.ID,
			"",
			rawQuantityUsed,
			fmt.Sprintf("Conversion: %s → %s", conversion.RawProductName, conversion.FinishedProductName),
			userID,
		); err != nil {
			record.Status = "FAILED"
			s.recordRepo.Update(record)
			return nil, fmt.Errorf("error deducting raw product stock: %w", err)
		}

		if variantSKU != "" {
			finishedProduct, err := s.productRepo.FindByID(conversion.FinishedProductID)
			if err == nil && finishedProduct != nil {
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

			if err := s.variantStockMgmt.RecordStockAdjustment(
				variantSKU,
				producedQuantity,
				"in",
				fmt.Sprintf("Conversion from: %s (Qty: %v grams, Conversion ID: %s)", conversion.RawProductName, rawQuantityUsed, record.ID),
				userID,
			); err != nil {
				record.Status = "FAILED"
				s.recordRepo.Update(record)
				return nil, fmt.Errorf("error adding stock to variant %s: %w", variantSKU, err)
			}
		} else {
			if err := s.stockManagementSvc.RecordInboundMovement(
				conversion.FinishedProductID,
				"PRODUCTION_USAGE",
				record.ID,
				"",
				producedQuantity,
				0,
				fmt.Sprintf("Conversion from: %s (Qty: %v grams)", conversion.RawProductName, rawQuantityUsed),
				userID,
			); err != nil {
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
		Message: fmt.Sprintf(
			"Conversion executed successfully. Raw: %.2f grams used, Finished: %.2f produced to %s",
			rawQuantityUsed,
			producedQuantity,
			func() string {
				if record.FinishedVariantSKU != "" {
					return "variant " + record.FinishedVariantSKU
				}
				return "product level"
			}(),
		),
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
	// Retrieve bag usage details
	bagsUsed, _ := s.bagUsageRepo.GetByConversionRecordID(record.ID)
	bagOutputs := make([]output.ConversionRecordBagUsageOutput, len(bagsUsed))
	for i, bagUsage := range bagsUsed {
		bagOutputs[i] = output.ConversionRecordBagUsageOutput{
			ID:             bagUsage.ID,
			BagID:          bagUsage.BagID,
			BagNumber:      bagUsage.BagNumber,
			ProductID:      bagUsage.ProductID,
			ProductName:    bagUsage.ProductName,
			QuantityUsedKg: bagUsage.QuantityUsedKg,
			CreatedAt:      bagUsage.CreatedAt,
		}
	}

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
		BagsUsed:                 bagOutputs,
		CreatedAt:                record.CreatedAt,
	}
}

func (s *productConversionService) validateConversionUser(userID string, companyID uint) error {
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil || id == 0 {
		return errors.New("invalid authenticated user")
	}
	user, err := s.userRepo.GetByIDAndCompanyID(id, companyID)
	if err != nil || user == nil {
		return errors.New("user does not belong to the company")
	}
	return nil
}

func conversionListOutput(items []models.ProductConversion, total int64, offset, limit int, mapper func(*models.ProductConversion) *output.ProductConversionOutput) *output.ProductConversionListOutput {
	results := make([]output.ProductConversionOutput, len(items))
	for i := range items {
		results[i] = *mapper(&items[i])
	}
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}
	return &output.ProductConversionListOutput{Conversions: results, Total: total, Page: page, Limit: limit}
}

func recordListOutput(items []models.ProductConversionRecord, total int64, offset, limit int, mapper func(*models.ProductConversionRecord) *output.ProductConversionRecordOutput) *output.ProductConversionRecordListOutput {
	results := make([]output.ProductConversionRecordOutput, len(items))
	for i := range items {
		results[i] = *mapper(&items[i])
	}
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}
	return &output.ProductConversionRecordListOutput{Records: results, Total: total, Page: page, Limit: limit}
}

func (s *productConversionService) CreateConversionForCompany(in *input.CreateProductConversionInput, userID string, companyID uint, userName, companyName string) (*output.ProductConversionOutput, error) {
	if in == nil {
		return nil, errors.New("input cannot be nil")
	}
	if err := s.validateConversionUser(userID, companyID); err != nil {
		return nil, err
	}
	raw, err := s.productRepo.FindByIDAndCompany(in.RawProductID, companyID)
	if err != nil || raw == nil {
		return nil, errors.New("raw product not found in your company")
	}
	finished, err := s.productRepo.FindByIDAndCompany(in.FinishedProductID, companyID)
	if err != nil || finished == nil {
		return nil, errors.New("finished product not found in your company")
	}
	existing, err := s.conversionRepo.GetByProductPairAndCompany(in.RawProductID, in.FinishedProductID, companyID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("conversion rule already exists for this product pair in your company")
	}
	conversion := &models.ProductConversion{
		ID: uuid.New().String(), RawProductID: in.RawProductID, RawProductName: raw.Name,
		RawProductSpec: raw.RawSpecification, FinishedProductID: in.FinishedProductID,
		FinishedProductName: finished.Name, FinishedVariantSKU: in.FinishedVariantSKU,
		ConversionRatio: in.ConversionRatio, LossPercentage: in.LossPercentage, IsActive: in.IsActive,
		Notes: in.Notes, CreatedBy: userID, CreatedByUserName: userName,
		CreatedByCompanyID: companyID, CreatedByCompanyName: companyName,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := s.conversionRepo.CreateForCompany(conversion, companyID); err != nil {
		return nil, err
	}
	return s.mapConversionToOutput(conversion), nil
}

func (s *productConversionService) UpdateConversionForCompany(id string, in *input.UpdateProductConversionInput, userID, userName string, companyID uint) (*output.ProductConversionOutput, error) {
	if err := s.validateConversionUser(userID, companyID); err != nil {
		return nil, err
	}
	c, err := s.conversionRepo.GetByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("conversion not found")
	}
	if in.ConversionRatio != nil {
		c.ConversionRatio = *in.ConversionRatio
	}
	if in.LossPercentage != nil {
		c.LossPercentage = *in.LossPercentage
	}
	if in.IsActive != nil {
		c.IsActive = *in.IsActive
	}
	if in.FinishedVariantSKU != nil {
		c.FinishedVariantSKU = *in.FinishedVariantSKU
	}
	if in.Notes != nil {
		c.Notes = *in.Notes
	}
	c.UpdatedBy, c.UpdatedByUserName, c.UpdatedAt = userID, userName, time.Now()
	if err := s.conversionRepo.UpdateByCompany(c, companyID); err != nil {
		return nil, err
	}
	return s.mapConversionToOutput(c), nil
}

func (s *productConversionService) DeleteConversionForCompany(id string, companyID uint) error {
	return s.conversionRepo.DeleteByCompany(id, companyID)
}
func (s *productConversionService) GetConversionForCompany(id string, companyID uint) (*output.ProductConversionOutput, error) {
	c, err := s.conversionRepo.GetByIDAndCompany(id, companyID)
	if err != nil {
		return nil, err
	}
	return s.mapConversionToOutput(c), nil
}
func (s *productConversionService) ListConversionsForCompany(companyID uint, offset, limit int) (*output.ProductConversionListOutput, error) {
	items, total, err := s.conversionRepo.GetAllByCompany(companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return conversionListOutput(items, total, offset, limit, s.mapConversionToOutput), nil
}
func (s *productConversionService) ListActiveConversionsForCompany(companyID uint, offset, limit int) (*output.ProductConversionListOutput, error) {
	items, total, err := s.conversionRepo.GetActiveConversionsByCompany(companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return conversionListOutput(items, total, offset, limit, s.mapConversionToOutput), nil
}
func (s *productConversionService) GetConversionsByRawProductForCompany(productID string, companyID uint, offset, limit int) (*output.ProductConversionListOutput, error) {
	if _, err := s.productRepo.FindByIDAndCompany(productID, companyID); err != nil {
		return nil, errors.New("raw product not found in your company")
	}
	items, total, err := s.conversionRepo.GetByRawProductIDAndCompany(productID, companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return conversionListOutput(items, total, offset, limit, s.mapConversionToOutput), nil
}
func (s *productConversionService) GetConversionsByFinishedProductForCompany(productID string, companyID uint, offset, limit int) (*output.ProductConversionListOutput, error) {
	if _, err := s.productRepo.FindByIDAndCompany(productID, companyID); err != nil {
		return nil, errors.New("finished product not found in your company")
	}
	items, total, err := s.conversionRepo.GetByFinishedProductIDAndCompany(productID, companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return conversionListOutput(items, total, offset, limit, s.mapConversionToOutput), nil
}

func (s *productConversionService) ExecuteConversionForCompany(in *input.CreateProductConversionRecordInput, userID string, companyID uint, userName, companyName string) (*output.ConversionExecutionOutput, error) {
	if in == nil {
		return nil, errors.New("input cannot be nil")
	}
	if err := s.validateConversionUser(userID, companyID); err != nil {
		return nil, err
	}
	conversion, err := s.conversionRepo.GetByIDAndCompany(in.ConversionID, companyID)
	if err != nil || conversion == nil {
		return nil, errors.New("conversion not found")
	}
	if _, err = s.productRepo.FindByIDAndCompany(conversion.RawProductID, companyID); err != nil {
		return nil, errors.New("raw product not found in your company")
	}
	if _, err = s.productRepo.FindByIDAndCompany(conversion.FinishedProductID, companyID); err != nil {
		return nil, errors.New("finished product not found in your company")
	}
	for _, bag := range in.RawMaterialBags {
		if _, err := s.rawMaterialBagSvc.GetByIDForCompany(bag.BagID, companyID); err != nil {
			return nil, fmt.Errorf("bag %s not found in your company", bag.BagID)
		}
	}
	// Existing execution logic is retained after all company ownership checks.
	return s.ExecuteConversion(in, userID, companyID, userName, companyName)
}

func (s *productConversionService) GetConversionRecordForCompany(id string, companyID uint) (*output.ProductConversionRecordOutput, error) {
	r, err := s.recordRepo.GetByIDAndCompany(id, companyID)
	if err != nil {
		return nil, err
	}
	return s.mapRecordToOutput(r), nil
}
func (s *productConversionService) ListConversionRecordsForCompany(companyID uint, offset, limit int) (*output.ProductConversionRecordListOutput, error) {
	items, total, err := s.recordRepo.GetAllByCompany(companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return recordListOutput(items, total, offset, limit, s.mapRecordToOutput), nil
}
func (s *productConversionService) ListConversionRecordsByRuleForCompany(conversionID string, companyID uint, offset, limit int) (*output.ProductConversionRecordListOutput, error) {
	if _, err := s.conversionRepo.GetByIDAndCompany(conversionID, companyID); err != nil {
		return nil, errors.New("conversion not found")
	}
	items, total, err := s.recordRepo.GetByConversionIDAndCompany(conversionID, companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return recordListOutput(items, total, offset, limit, s.mapRecordToOutput), nil
}
func (s *productConversionService) ListConversionRecordsByDateRangeForCompany(from, to time.Time, companyID uint, offset, limit int) (*output.ProductConversionRecordListOutput, error) {
	items, total, err := s.recordRepo.GetByDateRangeAndCompany(from, to, companyID, offset, limit)
	if err != nil {
		return nil, err
	}
	return recordListOutput(items, total, offset, limit, s.mapRecordToOutput), nil
}