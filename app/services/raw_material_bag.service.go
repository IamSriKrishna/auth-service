package services

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type RawMaterialBagService interface {
	// Existing methods retained.
	ReceiveBags(
		req *input.ReceiveRawMaterialBagsInput,
		createdBy string,
	) (*output.ReceiveRawMaterialBagsOutput, error)

	GetAll(
		limit int,
		offset int,
	) (*output.RawMaterialBagListOutput, error)

	GetByID(
		id string,
	) (*output.RawMaterialBagOutput, error)

	GetBagsByProduct(
		productID string,
	) ([]output.RawMaterialBagOutput, error)

	GetBagsByPurchaseOrder(
		purchaseOrderID string,
	) ([]output.RawMaterialBagOutput, error)

	UseBags(
		productID string,
		bags []input.UseRawMaterialBagInput,
	) (float64, error)

	// Company-scoped methods.
	ReceiveBagsForCompany(
		req *input.ReceiveRawMaterialBagsInput,
		createdBy string,
		companyID uint,
	) (*output.ReceiveRawMaterialBagsOutput, error)

	GetAllForCompany(
		companyID uint,
		limit int,
		offset int,
	) (*output.RawMaterialBagListOutput, error)

	GetByIDForCompany(
		id string,
		companyID uint,
	) (*output.RawMaterialBagOutput, error)

	GetBagsByProductForCompany(
		productID string,
		companyID uint,
	) ([]output.RawMaterialBagOutput, error)

	GetBagsByPurchaseOrderForCompany(
		purchaseOrderID string,
		companyID uint,
	) ([]output.RawMaterialBagOutput, error)

	UseBagsForCompany(
		productID string,
		bags []input.UseRawMaterialBagInput,
		companyID uint,
	) (float64, error)
}

type rawMaterialBagService struct {
	bagRepo     repo.RawMaterialBagRepository
	claimRepo   repo.VendorShortageClaimRepository
	poRepo      repo.PurchaseOrderRepository
	productRepo repo.ProductRepository
	userRepo    repo.UserRepository
	companyRepo repo.CompanyRepository
}

func NewRawMaterialBagService(
	bagRepo repo.RawMaterialBagRepository,
	claimRepo repo.VendorShortageClaimRepository,
	poRepo repo.PurchaseOrderRepository,
	productRepo repo.ProductRepository,
	userRepo repo.UserRepository,
	companyRepo repo.CompanyRepository,
) RawMaterialBagService {
	return &rawMaterialBagService{
		bagRepo:     bagRepo,
		claimRepo:   claimRepo,
		poRepo:      poRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
	}
}

func rawMaterialBagOutputs(
	bags []models.RawMaterialBag,
) []output.RawMaterialBagOutput {
	results := make(
		[]output.RawMaterialBagOutput,
		len(bags),
	)

	for index := range bags {
		results[index] =
			output.ToRawMaterialBagOutput(&bags[index])
	}

	return results
}

func (s *rawMaterialBagService) GetAll(
	limit int,
	offset int,
) (*output.RawMaterialBagListOutput, error) {
	bags, total, err := s.bagRepo.GetAll(
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return &output.RawMaterialBagListOutput{
		Bags:  rawMaterialBagOutputs(bags),
		Total: total,
	}, nil
}

func (s *rawMaterialBagService) GetByID(
	id string,
) (*output.RawMaterialBagOutput, error) {
	bag, err := s.bagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := output.ToRawMaterialBagOutput(bag)
	return &result, nil
}

func (s *rawMaterialBagService) ReceiveBags(
	req *input.ReceiveRawMaterialBagsInput,
	createdBy string,
) (*output.ReceiveRawMaterialBagsOutput, error) {
	return s.receiveBags(
		req,
		createdBy,
		0,
		false,
	)
}

func (s *rawMaterialBagService) ReceiveBagsForCompany(
	req *input.ReceiveRawMaterialBagsInput,
	createdBy string,
	companyID uint,
) (*output.ReceiveRawMaterialBagsOutput, error) {
	if companyID == 0 {
		return nil, errors.New("invalid company")
	}

	userID, err := strconv.ParseUint(
		createdBy,
		10,
		64,
	)
	if err != nil || userID == 0 {
		return nil, errors.New(
			"invalid authenticated user",
		)
	}

	user, err := s.userRepo.GetByIDAndCompanyID(
		uint(userID),
		companyID,
	)
	if err != nil || user == nil {
		return nil, errors.New(
			"user does not belong to the company",
		)
	}

	return s.receiveBags(
		req,
		createdBy,
		companyID,
		true,
	)
}

func (s *rawMaterialBagService) receiveBags(
	req *input.ReceiveRawMaterialBagsInput,
	createdBy string,
	companyID uint,
	useCompanyFilter bool,
) (*output.ReceiveRawMaterialBagsOutput, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	if req.PurchaseOrderID == "" {
		return nil, errors.New(
			"purchase_order_id is required",
		)
	}

	if req.ProductID == "" {
		return nil, errors.New(
			"product_id is required",
		)
	}

	if len(req.Bags) == 0 {
		return nil, errors.New(
			"at least one bag is required",
		)
	}

	if req.ExpectedKgPerBag <= 0 {
		return nil, errors.New(
			"expected kg per bag must be greater than zero",
		)
	}

	var purchaseOrder *models.PurchaseOrder
	var product *models.Product
	var err error

	if useCompanyFilter {
		purchaseOrder, err =
			s.poRepo.FindByIDAndCompany(
				req.PurchaseOrderID,
				companyID,
			)
	} else {
		purchaseOrder, err =
			s.poRepo.FindByID(
				req.PurchaseOrderID,
			)
	}

	if err != nil {
		return nil, errors.New(
			"purchase order not found",
		)
	}

	if useCompanyFilter {
		product, err =
			s.productRepo.FindByIDAndCompany(
				req.ProductID,
				companyID,
			)
	} else {
		product, err =
			s.productRepo.FindByID(
				req.ProductID,
			)
	}

	if err != nil {
		return nil, errors.New("product not found")
	}

	productExistsInPO := false
	for _, lineItem := range purchaseOrder.LineItems {
		if lineItem.ProductID != nil &&
			*lineItem.ProductID == req.ProductID {
			productExistsInPO = true
			break
		}
	}

	if !productExistsInPO {
		return nil, errors.New(
			"product does not belong to the purchase order",
		)
	}

	createdByUserName := ""
	createdByCompanyID := uint(0)
	createdByCompanyName := ""

	if createdBy != "" {
		userID, parseErr := strconv.ParseUint(
			createdBy,
			10,
			64,
		)
		if parseErr == nil && userID > 0 {
			var user *models.User

			if useCompanyFilter {
				user, err =
					s.userRepo.GetByIDAndCompanyID(
						uint(userID),
						companyID,
					)
			} else {
				user, err =
					s.userRepo.GetByID(
						uint(userID),
					)
			}

			if err == nil && user != nil {
				if user.Email != nil {
					createdByUserName = *user.Email
				} else if user.Username != nil {
					createdByUserName = *user.Username
				}

				if user.CompanyID != nil {
					createdByCompanyID = *user.CompanyID

					company, companyErr :=
						s.companyRepo.FindByID(
							*user.CompanyID,
						)
					if companyErr == nil &&
						company != nil {
						createdByCompanyName =
							company.CompanyName
					}
				}
			}
		}
	}

	if useCompanyFilter &&
		createdByCompanyID != companyID {
		return nil, errors.New(
			"authenticated user company mismatch",
		)
	}

	expectedTotalKG :=
		float64(len(req.Bags)) *
			req.ExpectedKgPerBag

	actualTotalKG := 0.0
	bags := make(
		[]models.RawMaterialBag,
		0,
		len(req.Bags),
	)

	vendorName := ""
	if purchaseOrder.Vendor != nil {
		vendorName =
			purchaseOrder.Vendor.DisplayName
	}

	now := time.Now()

	for _, bagInput := range req.Bags {
		if bagInput.ActualKg < 0 {
			return nil, fmt.Errorf(
				"actual kg cannot be negative for bag %d",
				bagInput.BagNumber,
			)
		}

		actualTotalKG += bagInput.ActualKg

		bags = append(
			bags,
			models.RawMaterialBag{
				ID: "bag_" +
					uuid.New().String()[:12],
				PurchaseOrderID:
					purchaseOrder.ID,
				PurchaseOrderNo:
					purchaseOrder.PurchaseOrderNumber,
				VendorID:
					purchaseOrder.VendorID,
				VendorName:
					vendorName,
				ProductID:
					product.ID,
				ProductName:
					product.Name,
				CreatedBy:
					createdBy,
				CreatedByUserName:
					createdByUserName,
				CreatedByCompanyID:
					createdByCompanyID,
				CreatedByCompanyName:
					createdByCompanyName,
				BagNumber:
					bagInput.BagNumber,
				ExpectedKg:
					req.ExpectedKgPerBag,
				ActualKg:
					bagInput.ActualKg,
				RemainingKg:
					bagInput.ActualKg,
				Status:
					"available",
				CreatedAt:
					now,
				UpdatedAt:
					now,
			},
		)
	}

	if useCompanyFilter {
		err = s.bagRepo.CreateManyForCompany(
			bags,
			companyID,
		)
	} else {
		err = s.bagRepo.CreateMany(bags)
	}

	if err != nil {
		return nil, err
	}

	shortageKG :=
		expectedTotalKG - actualTotalKG
	shortageKG =
		math.Round(shortageKG*10000) / 10000

	if shortageKG > 0 {
		ratePerKG := 0.0

		for _, lineItem :=
			range purchaseOrder.LineItems {
			if lineItem.ProductID != nil &&
				*lineItem.ProductID == req.ProductID {
				ratePerKG = lineItem.Rate
				break
			}
		}

		claim := &models.VendorShortageClaim{
			ID: "vsc_" +
				uuid.New().String()[:12],
			PurchaseOrderID:
				purchaseOrder.ID,
			PurchaseOrderNo:
				purchaseOrder.PurchaseOrderNumber,
			VendorID:
				purchaseOrder.VendorID,
			VendorName:
				vendorName,
			ProductID:
				product.ID,
			ProductName:
				product.Name,
			ExpectedKg:
				expectedTotalKG,
			ReceivedKg:
				actualTotalKG,
			ShortageKg:
				shortageKG,
			ShortageGrams:
				shortageKG * 1000,
			RatePerKg:
				ratePerKG,
			ClaimAmount:
				shortageKG * ratePerKG,
			Status:
				"pending",
			Notes:
				"Auto-created from actual bag shortage",
			CreatedAt:
				now,
			UpdatedAt:
				now,
		}

		if useCompanyFilter {
			_ = s.claimRepo.CreateForCompany(
				claim,
				companyID,
			)
		} else {
			_ = s.claimRepo.Create(claim)
		}
	}

	return &output.ReceiveRawMaterialBagsOutput{
		PurchaseOrderID: purchaseOrder.ID,
		ProductID:       product.ID,
		ProductName:     product.Name,
		ExpectedKg:      expectedTotalKG,
		ActualKg:        actualTotalKG,
		ShortageKg:      shortageKG,
		ShortageGrams:   shortageKG * 1000,
		Bags:            rawMaterialBagOutputs(bags),
	}, nil
}

func (s *rawMaterialBagService) GetBagsByProduct(
	productID string,
) ([]output.RawMaterialBagOutput, error) {
	bags, err := s.bagRepo.GetByProductID(
		productID,
	)
	if err != nil {
		return nil, err
	}

	return rawMaterialBagOutputs(bags), nil
}

func (s *rawMaterialBagService) GetBagsByPurchaseOrder(
	purchaseOrderID string,
) ([]output.RawMaterialBagOutput, error) {
	bags, err :=
		s.bagRepo.GetByPurchaseOrderID(
			purchaseOrderID,
		)
	if err != nil {
		return nil, err
	}

	return rawMaterialBagOutputs(bags), nil
}

func (s *rawMaterialBagService) UseBags(
	productID string,
	bags []input.UseRawMaterialBagInput,
) (float64, error) {
	return s.useBags(
		productID,
		bags,
		0,
		false,
	)
}

func (s *rawMaterialBagService) GetAllForCompany(
	companyID uint,
	limit int,
	offset int,
) (*output.RawMaterialBagListOutput, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	bags, total, err :=
		s.bagRepo.GetAllByCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, err
	}

	return &output.RawMaterialBagListOutput{
		Bags:  rawMaterialBagOutputs(bags),
		Total: total,
	}, nil
}

func (s *rawMaterialBagService) GetByIDForCompany(
	id string,
	companyID uint,
) (*output.RawMaterialBagOutput, error) {
	bag, err :=
		s.bagRepo.GetByIDAndCompany(
			id,
			companyID,
		)
	if err != nil {
		return nil, err
	}

	result := output.ToRawMaterialBagOutput(bag)
	return &result, nil
}

func (s *rawMaterialBagService) GetBagsByProductForCompany(
	productID string,
	companyID uint,
) ([]output.RawMaterialBagOutput, error) {
	if _, err :=
		s.productRepo.FindByIDAndCompany(
			productID,
			companyID,
		); err != nil {
		return nil, errors.New(
			"product not found in your company",
		)
	}

	bags, err :=
		s.bagRepo.GetByProductIDAndCompany(
			productID,
			companyID,
		)
	if err != nil {
		return nil, err
	}

	return rawMaterialBagOutputs(bags), nil
}

func (s *rawMaterialBagService) GetBagsByPurchaseOrderForCompany(
	purchaseOrderID string,
	companyID uint,
) ([]output.RawMaterialBagOutput, error) {
	if _, err :=
		s.poRepo.FindByIDAndCompany(
			purchaseOrderID,
			companyID,
		); err != nil {
		return nil, errors.New(
			"purchase order not found in your company",
		)
	}

	bags, err :=
		s.bagRepo.GetByPurchaseOrderIDAndCompany(
			purchaseOrderID,
			companyID,
		)
	if err != nil {
		return nil, err
	}

	return rawMaterialBagOutputs(bags), nil
}

func (s *rawMaterialBagService) UseBagsForCompany(
	productID string,
	bags []input.UseRawMaterialBagInput,
	companyID uint,
) (float64, error) {
	if _, err :=
		s.productRepo.FindByIDAndCompany(
			productID,
			companyID,
		); err != nil {
		return 0, errors.New(
			"product not found in your company",
		)
	}

	return s.useBags(
		productID,
		bags,
		companyID,
		true,
	)
}

func (s *rawMaterialBagService) useBags(
	productID string,
	bags []input.UseRawMaterialBagInput,
	companyID uint,
	useCompanyFilter bool,
) (float64, error) {
	totalUsedKG := 0.0

	for _, bagInput := range bags {
		if bagInput.QuantityKg <= 0 {
			return 0, fmt.Errorf(
				"quantity must be greater than zero for bag %s",
				bagInput.BagID,
			)
		}

		var bag *models.RawMaterialBag
		var err error

		if useCompanyFilter {
			bag, err =
				s.bagRepo.GetByIDAndCompany(
					bagInput.BagID,
					companyID,
				)
		} else {
			bag, err =
				s.bagRepo.GetByID(
					bagInput.BagID,
				)
		}

		if err != nil {
			return 0, fmt.Errorf(
				"bag not found: %s",
				bagInput.BagID,
			)
		}

		if bag.ProductID != productID {
			return 0, fmt.Errorf(
				"bag %s does not belong to product %s",
				bagInput.BagID,
				productID,
			)
		}

		if bag.RemainingKg < bagInput.QuantityKg {
			return 0, fmt.Errorf(
				"insufficient quantity in bag %d: available %.4f kg, requested %.4f kg",
				bag.BagNumber,
				bag.RemainingKg,
				bagInput.QuantityKg,
			)
		}

		bag.RemainingKg -= bagInput.QuantityKg

		if bag.RemainingKg <= 0 {
			bag.RemainingKg = 0
			bag.Status = "used"
		} else {
			bag.Status = "partial"
		}

		bag.UpdatedAt = time.Now()

		if useCompanyFilter {
			err = s.bagRepo.UpdateByCompany(
				bag,
				companyID,
			)
		} else {
			err = s.bagRepo.Update(bag)
		}

		if err != nil {
			return 0, err
		}

		totalUsedKG += bagInput.QuantityKg
	}

	return totalUsedKG, nil
}