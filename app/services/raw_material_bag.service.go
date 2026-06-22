package services

import (
	"fmt"
	"math"
	"time"

	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type RawMaterialBagService interface {
	ReceiveBags(req *input.ReceiveRawMaterialBagsInput, createdBy string) (*output.ReceiveRawMaterialBagsOutput, error)
	GetAll(limit, offset int) (*output.RawMaterialBagListOutput, error)
	GetByID(id string) (*output.RawMaterialBagOutput, error)
	GetBagsByProduct(productID string) ([]output.RawMaterialBagOutput, error)
	GetBagsByPurchaseOrder(poID string) ([]output.RawMaterialBagOutput, error)
	UseBags(productID string, bags []input.UseRawMaterialBagInput) (float64, error)
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

func (s *rawMaterialBagService) GetAll(limit, offset int) (*output.RawMaterialBagListOutput, error) {
	bags, total, err := s.bagRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	res := make([]output.RawMaterialBagOutput, len(bags))
	for i := range bags {
		res[i] = output.ToRawMaterialBagOutput(&bags[i])
	}

	return &output.RawMaterialBagListOutput{
		Bags:  res,
		Total: total,
	}, nil
}

func (s *rawMaterialBagService) GetByID(id string) (*output.RawMaterialBagOutput, error) {
	bag, err := s.bagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	res := output.ToRawMaterialBagOutput(bag)
	return &res, nil
}

func (s *rawMaterialBagService) ReceiveBags(req *input.ReceiveRawMaterialBagsInput, createdBy string) (*output.ReceiveRawMaterialBagsOutput, error) {
	po, err := s.poRepo.FindByID(req.PurchaseOrderID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found")
	}

	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	expectedTotalKg := float64(len(req.Bags)) * req.ExpectedKgPerBag
	actualTotalKg := 0.0

	bags := make([]models.RawMaterialBag, 0, len(req.Bags))

	var createdByUserName string
	var createdByCompanyID uint
	var createdByCompanyName string

	if createdBy != "" {
		var userID uint
		_, err := fmt.Sscanf(createdBy, "%d", &userID)
		if err == nil {
			user, err := s.userRepo.GetByID(userID)
			if err == nil && user != nil {
				if user.Email != nil {
					createdByUserName = *user.Email
				} else if user.Username != nil {
					createdByUserName = *user.Username
				}

				if user.CompanyID != nil {
					createdByCompanyID = *user.CompanyID
					company, err := s.companyRepo.FindByID(*user.CompanyID)
					if err == nil && company != nil {
						createdByCompanyName = company.CompanyName
					}
				}
			}
		}
	}

	for _, b := range req.Bags {
		actualTotalKg += b.ActualKg

		vendorName := ""
		if po.Vendor != nil {
			vendorName = po.Vendor.DisplayName
		}

		bag := models.RawMaterialBag{
			ID:                   "bag_" + uuid.New().String()[:12],
			PurchaseOrderID:      po.ID,
			PurchaseOrderNo:      po.PurchaseOrderNumber,
			VendorID:             po.VendorID,
			VendorName:           vendorName,
			ProductID:            product.ID,
			ProductName:          product.Name,
			CreatedBy:            createdBy,
			CreatedByUserName:    createdByUserName,
			CreatedByCompanyID:   createdByCompanyID,
			CreatedByCompanyName: createdByCompanyName,
			BagNumber:            b.BagNumber,
			ExpectedKg:           req.ExpectedKgPerBag,
			ActualKg:             b.ActualKg,
			RemainingKg:          b.ActualKg,
			Status:               "available",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		bags = append(bags, bag)
	}

	if err := s.bagRepo.CreateMany(bags); err != nil {
		return nil, err
	}

	shortageKg := expectedTotalKg - actualTotalKg
	shortageKg = math.Round(shortageKg*10000) / 10000

	if shortageKg > 0 {
		ratePerKg := 0.0

		for _, item := range po.LineItems {
			if item.ProductID != nil && *item.ProductID == req.ProductID {
				ratePerKg = item.Rate
				break
			}
		}

		vendorName := ""
		if po.Vendor != nil {
			vendorName = po.Vendor.DisplayName
		}

		claim := &models.VendorShortageClaim{
			ID:              "vsc_" + uuid.New().String()[:12],
			PurchaseOrderID: po.ID,
			PurchaseOrderNo: po.PurchaseOrderNumber,
			VendorID:        po.VendorID,
			VendorName:      vendorName,
			ProductID:       product.ID,
			ProductName:     product.Name,
			ExpectedKg:      expectedTotalKg,
			ReceivedKg:      actualTotalKg,
			ShortageKg:      shortageKg,
			ShortageGrams:   shortageKg * 1000,
			RatePerKg:       ratePerKg,
			ClaimAmount:     shortageKg * ratePerKg,
			Status:          "pending",
			Notes:           "Auto-created from actual bag shortage",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		_ = s.claimRepo.Create(claim)
	}

	bagOutputs := make([]output.RawMaterialBagOutput, len(bags))
	for i := range bags {
		bagOutputs[i] = output.ToRawMaterialBagOutput(&bags[i])
	}

	return &output.ReceiveRawMaterialBagsOutput{
		PurchaseOrderID: po.ID,
		ProductID:       product.ID,
		ProductName:     product.Name,
		ExpectedKg:      expectedTotalKg,
		ActualKg:        actualTotalKg,
		ShortageKg:      shortageKg,
		ShortageGrams:   shortageKg * 1000,
		Bags:            bagOutputs,
	}, nil
}

func (s *rawMaterialBagService) GetBagsByProduct(productID string) ([]output.RawMaterialBagOutput, error) {
	bags, err := s.bagRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	res := make([]output.RawMaterialBagOutput, len(bags))
	for i := range bags {
		res[i] = output.ToRawMaterialBagOutput(&bags[i])
	}

	return res, nil
}

func (s *rawMaterialBagService) GetBagsByPurchaseOrder(poID string) ([]output.RawMaterialBagOutput, error) {
	bags, err := s.bagRepo.GetByPurchaseOrderID(poID)
	if err != nil {
		return nil, err
	}

	res := make([]output.RawMaterialBagOutput, len(bags))
	for i := range bags {
		res[i] = output.ToRawMaterialBagOutput(&bags[i])
	}

	return res, nil
}

func (s *rawMaterialBagService) UseBags(productID string, bags []input.UseRawMaterialBagInput) (float64, error) {
	totalUsedKg := 0.0

	for _, b := range bags {
		bag, err := s.bagRepo.GetByID(b.BagID)
		if err != nil {
			return 0, fmt.Errorf("bag not found: %s", b.BagID)
		}

		if bag.ProductID != productID {
			return 0, fmt.Errorf("bag %s does not belong to product %s", b.BagID, productID)
		}

		if bag.RemainingKg < b.QuantityKg {
			return 0, fmt.Errorf("insufficient quantity in bag %d: available %.4f kg, requested %.4f kg",
				bag.BagNumber,
				bag.RemainingKg,
				b.QuantityKg,
			)
		}

		bag.RemainingKg -= b.QuantityKg

		if bag.RemainingKg <= 0 {
			bag.RemainingKg = 0
			bag.Status = "used"
		} else {
			bag.Status = "partial"
		}

		bag.UpdatedAt = time.Now()

		if err := s.bagRepo.Update(bag); err != nil {
			return 0, err
		}

		totalUsedKg += b.QuantityKg
	}

	return totalUsedKg, nil
}
