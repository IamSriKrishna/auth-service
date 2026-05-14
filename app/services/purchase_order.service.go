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

type PurchaseOrderService interface {
	// Basic CRUD Operations
	CreatePurchaseOrder(poInput *input.CreatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error)
	GetPurchaseOrder(id string) (*output.PurchaseOrderOutput, error)
	GetAllPurchaseOrders(limit, offset int) (*output.PurchaseOrderListOutput, error)
	UpdatePurchaseOrder(id string, poInput *input.UpdatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error)
	DeletePurchaseOrder(id string) error
	GetPurchaseOrdersByVendor(vendorID uint, limit, offset int) (*output.PurchaseOrderListOutput, error)
	GetPurchaseOrdersByCustomer(customerID uint, limit, offset int) (*output.PurchaseOrderListOutput, error)
	GetPurchaseOrdersByStatus(status string, limit, offset int) (*output.PurchaseOrderListOutput, error)

	// Step 3: Purchasing Stock (Inbound Operations)
	// Update PO status and trigger inventory sync when stock is received
	UpdatePurchaseOrderStatus(id string, status domain.PurchaseOrderStatus, userID string) (*output.PurchaseOrderOutput, error)

	// Reorder Operations
	// ReorderProductGroup allows reordering a product group with custom component quantities
	ReorderProductGroup(reorderInput *input.ReorderProductGroupInput, userID string) (*output.PurchaseOrderOutput, error)
}

type purchaseOrderService struct {
	poRepo             repo.PurchaseOrderRepository
	vendorRepo         repo.VendorRepository
	customerRepo       repo.CustomerRepository
	productRepo        repo.ProductRepository
	productGroupRepo   repo.ProductGroupRepository
	taxRepo            repo.TaxRepository
	userRepo           repo.UserRepository
	companyRepo        repo.CompanyRepository
	stockManagementSvc StockManagementService
	stockLedgerRepo    repo.StockLedgerRepository
	variantStockMgmt   VariantStockManagementService
}

func NewPurchaseOrderService(
	poRepo repo.PurchaseOrderRepository,
	vendorRepo repo.VendorRepository,
	customerRepo repo.CustomerRepository,
	productRepo repo.ProductRepository,
	productGroupRepo repo.ProductGroupRepository,
	taxRepo repo.TaxRepository,
	userRepo repo.UserRepository,
	companyRepo repo.CompanyRepository,
	stockManagementSvc StockManagementService,
	stockLedgerRepo repo.StockLedgerRepository,
	variantStockMgmt VariantStockManagementService,
) PurchaseOrderService {
	return &purchaseOrderService{
		poRepo:             poRepo,
		vendorRepo:         vendorRepo,
		customerRepo:       customerRepo,
		productRepo:        productRepo,
		productGroupRepo:   productGroupRepo,
		taxRepo:            taxRepo,
		userRepo:           userRepo,
		companyRepo:        companyRepo,
		stockManagementSvc: stockManagementSvc,
		stockLedgerRepo:    stockLedgerRepo,
		variantStockMgmt:   variantStockMgmt,
	}
}

func (s *purchaseOrderService) CreatePurchaseOrder(poInput *input.CreatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error) {
	vendor, err := s.vendorRepo.FindByID(poInput.VendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	if poInput.DeliveryAddressType != "organization" && poInput.DeliveryAddressType != "customer" {
		return nil, errors.New("invalid delivery address type")
	}

	var customer *models.Customer
	if poInput.DeliveryAddressType == "customer" {
		if poInput.CustomerID == nil {
			return nil, errors.New("customer_id is required for customer delivery")
		}
		var err error
		customer, err = s.customerRepo.FindByID(*poInput.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}

	var tax *models.Tax
	if poInput.TaxID != nil {
		tax, err = s.taxRepo.FindByID(*poInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
	}

	lineItems := make([]models.PurchaseOrderLineItem, 0)
	subTotal := 0.0

	for _, itemInput := range poInput.LineItems {
		// Validate product exists
		if itemInput.ProductID == nil {
			return nil, errors.New("product_id is required for each line item")
		}

		product, err := s.productRepo.FindByID(*itemInput.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", *itemInput.ProductID)
		}

		amount := itemInput.Quantity * itemInput.Rate
		subTotal += amount

		lineItem := models.PurchaseOrderLineItem{
			ProductID:   itemInput.ProductID,
			Product:     product,
			ProductName: itemInput.ProductName,
			SKU:         itemInput.SKU,
			Account:     itemInput.Account,
			Quantity:    itemInput.Quantity,
			Rate:        itemInput.Rate,
			Amount:      amount,
		}

		lineItems = append(lineItems, lineItem)
	}

	discount := poInput.Discount
	if poInput.DiscountType == "percentage" {
		discount = (subTotal * poInput.Discount) / 100
	}

	taxAmount := 0.0
	if tax != nil {
		taxAmount = ((subTotal - discount) * tax.Rate) / 100
	}

	total := subTotal - discount + taxAmount + poInput.Adjustment

	poNumber := fmt.Sprintf("PO-%s-%04d", time.Now().Format("20060102"), s.generatePOSequence())

	var taxType *domain.TaxType
	if poInput.TaxType != nil {
		tt := domain.TaxType(*poInput.TaxType)
		taxType = &tt
	}

	po := &models.PurchaseOrder{
		ID:                  uuid.New().String(),
		PurchaseOrderNumber: poNumber,
		VendorID:            poInput.VendorID,
		Vendor:              vendor,
		DeliveryAddressType: poInput.DeliveryAddressType,
		DeliveryAddressID:   poInput.DeliveryAddressID,
		OrganizationName:    poInput.OrganizationName,
		OrganizationAddress: poInput.OrganizationAddress,
		CustomerID:          poInput.CustomerID,
		Customer:            customer,
		ReferenceNo:         poInput.ReferenceNo,
		PODate:              poInput.Date,
		DeliveryDate:        poInput.DeliveryDate,
		PaymentTerms:        domain.PaymentTerms(poInput.PaymentTerms),
		ShipmentPreference:  poInput.ShipmentPreference,
		LineItems:           lineItems,
		SubTotal:            subTotal,
		Discount:            discount,
		DiscountType:        poInput.DiscountType,
		TaxType:             taxType,
		TaxID:               poInput.TaxID,
		Tax:                 tax,
		TaxAmount:           taxAmount,
		Adjustment:          poInput.Adjustment,
		Total:               total,
		Notes:               poInput.Notes,
		TermsAndConditions:  poInput.TermsAndConditions,
		Status:              domain.PurchaseOrderStatusDraft,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CreatedBy:           userID,
	}

	// Fetch user details if userID is provided
	if userID != "" {
		var userIDUint uint
		_, err := fmt.Sscanf(userID, "%d", &userIDUint)
		if err == nil {
			user, err := s.userRepo.GetByID(userIDUint)
			if err == nil && user != nil {
				// Set user name from email or username
				if user.Email != nil {
					po.CreatedByUserName = *user.Email
				} else if user.Username != nil {
					po.CreatedByUserName = *user.Username
				}

				// Set company details if user has a company
				if user.CompanyID != nil {
					po.CreatedByCompanyID = *user.CompanyID
					company, err := s.companyRepo.FindByID(*user.CompanyID)
					if err == nil && company != nil {
						po.CreatedByCompanyName = company.CompanyName
					}
				}
			}
		}
	}

	if len(poInput.Attachments) > 0 {
		po.Attachments = poInput.Attachments
	}

	createdPO, err := s.poRepo.Create(po)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return output.ToPurchaseOrderOutput(createdPO)
}

func (s *purchaseOrderService) GetPurchaseOrder(id string) (*output.PurchaseOrderOutput, error) {
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	return output.ToPurchaseOrderOutput(po)
}

func (s *purchaseOrderService) GetAllPurchaseOrders(limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) UpdatePurchaseOrder(id string, poInput *input.UpdatePurchaseOrderInput, userID string) (*output.PurchaseOrderOutput, error) {
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	if poInput.VendorID != nil {
		vendor, err := s.vendorRepo.FindByID(*poInput.VendorID)
		if err != nil {
			return nil, errors.New("vendor not found")
		}
		po.VendorID = *poInput.VendorID
		po.Vendor = vendor
	}

	if poInput.DeliveryAddressType != nil {
		po.DeliveryAddressType = *poInput.DeliveryAddressType
	}

	if poInput.OrganizationName != nil {
		po.OrganizationName = *poInput.OrganizationName
	}

	if poInput.OrganizationAddress != nil {
		po.OrganizationAddress = *poInput.OrganizationAddress
	}

	if poInput.CustomerID != nil {
		customer, err := s.customerRepo.FindByID(*poInput.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
		po.CustomerID = poInput.CustomerID
		po.Customer = customer
	}

	if poInput.ReferenceNo != nil {
		po.ReferenceNo = *poInput.ReferenceNo
	}

	if poInput.Date != nil {
		po.PODate = *poInput.Date
	}

	if poInput.DeliveryDate != nil {
		po.DeliveryDate = *poInput.DeliveryDate
	}

	if poInput.PaymentTerms != nil {
		po.PaymentTerms = domain.PaymentTerms(*poInput.PaymentTerms)
	}

	if poInput.ShipmentPreference != nil {
		po.ShipmentPreference = *poInput.ShipmentPreference
	}

	if len(poInput.LineItems) > 0 {
		lineItems := make([]models.PurchaseOrderLineItem, 0)
		subTotal := 0.0

		for _, itemInput := range poInput.LineItems {
			// Validate product exists
			if itemInput.ProductID == nil {
				return nil, errors.New("product_id is required for each line item")
			}

			product, err := s.productRepo.FindByID(*itemInput.ProductID)
			if err != nil {
				return nil, fmt.Errorf("product %s not found", *itemInput.ProductID)
			}

			amount := itemInput.Quantity * itemInput.Rate
			subTotal += amount

			lineItem := models.PurchaseOrderLineItem{
				ProductID:   itemInput.ProductID,
				Product:     product,
				ProductName: itemInput.ProductName,
				SKU:         itemInput.SKU,
				Account:     itemInput.Account,
				Quantity:    itemInput.Quantity,
				Rate:        itemInput.Rate,
				Amount:      amount,
			}

			lineItems = append(lineItems, lineItem)
		}

		po.LineItems = lineItems
		po.SubTotal = subTotal
	}

	if poInput.Discount != nil {
		po.Discount = *poInput.Discount
	}

	if poInput.DiscountType != nil {
		po.DiscountType = *poInput.DiscountType
	}

	if poInput.TaxID != nil {
		tax, err := s.taxRepo.FindByID(*poInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
		po.TaxID = poInput.TaxID
		po.Tax = tax
	}

	if poInput.TaxType != nil {
		taxType := domain.TaxType(*poInput.TaxType)
		po.TaxType = &taxType
	}

	if poInput.Adjustment != nil {
		po.Adjustment = *poInput.Adjustment
	}

	discount := po.Discount
	if po.DiscountType == "percentage" {
		discount = (po.SubTotal * po.Discount) / 100
	}

	taxAmount := 0.0
	if po.Tax != nil {
		taxAmount = ((po.SubTotal - discount) * po.Tax.Rate) / 100
	}

	po.TaxAmount = taxAmount
	po.Total = po.SubTotal - discount + taxAmount + po.Adjustment

	if poInput.Notes != nil {
		po.Notes = *poInput.Notes
	}

	if poInput.TermsAndConditions != nil {
		po.TermsAndConditions = *poInput.TermsAndConditions
	}

	if len(poInput.Attachments) > 0 {
		po.Attachments = poInput.Attachments
	}

	updatedPO, err := s.poRepo.Update(id, po)
	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return output.ToPurchaseOrderOutput(updatedPO)
}

func (s *purchaseOrderService) DeletePurchaseOrder(id string) error {
	// First, get the PO to access its line items and inventory sync status
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return errors.New("purchase order not found")
	}

	// Only reverse stock if it was synced to inventory
	if po.InventorySynced {
		for _, lineItem := range po.LineItems {
			if lineItem.ProductID == nil || *lineItem.ProductID == "" {
				continue
			}

			log.Printf("[PO_DELETE] Reversing stock for PO %s: %s (Qty: %.2f, SKU: %s)", po.PurchaseOrderNumber, *lineItem.ProductID, lineItem.Quantity, lineItem.SKU)

			// MUTUALLY EXCLUSIVE: Reverse from either product or variant level
			if lineItem.SKU != "" {
				// This is a variant - use variant service to reverse
				err := s.variantStockMgmt.RecordStockAdjustment(
					lineItem.SKU,
					lineItem.Quantity,
					"out", // Out because we're reversing an "in"
					fmt.Sprintf("Purchase Order %s deleted", po.PurchaseOrderNumber),
					"system",
				)
				if err != nil {
					log.Printf("[PO_DELETE] Warning: Failed to reverse variant stock for SKU %s: %v", lineItem.SKU, err)
				}
			} else {
				// This is a base product - use stock management service
				err := s.stockManagementSvc.RecordStockAdjustment(
					*lineItem.ProductID,
					lineItem.Quantity,
					"out", // Out because we're reversing an "in"
					fmt.Sprintf("Purchase Order %s deleted", po.PurchaseOrderNumber),
					"system",
				)
				if err != nil {
					log.Printf("[PO_DELETE] Error reversing stock for product %s: %v", *lineItem.ProductID, err)
				}
			}
		}
	}

	// Delete associated stock ledger entries for this purchase order
	if err := s.stockLedgerRepo.DeleteByReferenceID(id); err != nil {
		log.Printf("[PO_DELETE] Warning: Failed to delete stock ledger entries: %v", err)
	}

	// Delete the purchase order itself
	log.Printf("[PO_DELETE] Deleted PO: %s", po.PurchaseOrderNumber)
	return s.poRepo.Delete(id)
}

func (s *purchaseOrderService) GetPurchaseOrdersByVendor(vendorID uint, limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindByVendor(vendorID, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) GetPurchaseOrdersByCustomer(customerID uint, limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindByCustomer(customerID, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) GetPurchaseOrdersByStatus(status string, limit, offset int) (*output.PurchaseOrderListOutput, error) {
	pos, total, err := s.poRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, err
	}

	outputs := make([]output.PurchaseOrderOutput, len(pos))
	for i, po := range pos {
		out, err := output.ToPurchaseOrderOutput(&po)
		if err != nil {
			return nil, err
		}
		outputs[i] = *out
	}

	return &output.PurchaseOrderListOutput{
		PurchaseOrders: outputs,
		Total:          total,
	}, nil
}

func (s *purchaseOrderService) UpdatePurchaseOrderStatus(id string, status domain.PurchaseOrderStatus, userID string) (*output.PurchaseOrderOutput, error) {
	po, err := s.poRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	// When PO is received, update inventory using stock management service
	// Only record inventory once, even if status changes back and forth
	if status == domain.PurchaseOrderStatusReceived && !po.InventorySynced {
		for _, lineItem := range po.LineItems {
			// Only process if ProductID is provided
			if lineItem.ProductID == nil || *lineItem.ProductID == "" {
				log.Printf("[PO_STATUS] Warning: No ProductID for line item in PO %s", po.ID)
				continue
			}

			// Record inbound stock movement - MUTUALLY EXCLUSIVE:
			// If line item has SKU (variant), record only at variant level
			// If no SKU (base product), record only at product level
			if lineItem.SKU != "" {
				// This is a variant purchase - record at variant level only
				err := s.variantStockMgmt.RecordPurchaseInbound(
					*lineItem.ProductID,
					lineItem.SKU,
					lineItem.Quantity,
					lineItem.Rate,
					"purchase_order",
					po.ID,
					po.PurchaseOrderNumber,
					userID,
				)

				if err != nil {
					log.Printf("[PO_STATUS] Error recording variant stock for SKU %s (Qty: %.2f): %v", lineItem.SKU, lineItem.Quantity, err)
				} else {
					log.Printf("[PO_STATUS] Successfully recorded variant stock: %s +%.2f units", lineItem.SKU, lineItem.Quantity)
				}
			} else {
				// This is a base product (no SKU) - record at product level only
				err := s.stockManagementSvc.RecordInboundMovement(
					*lineItem.ProductID,    // productID
					"purchase_order",       // referenceType
					po.ID,                  // referenceID
					po.PurchaseOrderNumber, // referenceNo
					lineItem.Quantity,      // quantity
					lineItem.Rate,          // rate
					fmt.Sprintf("Received from vendor %s", po.Vendor.DisplayName), // notes
					userID, // userID
				)

				if err != nil {
					log.Printf("[PO_STATUS] Error recording stock for product %s (Qty: %.2f): %v", *lineItem.ProductID, lineItem.Quantity, err)
				} else {
					log.Printf("[PO_STATUS] Successfully recorded stock: %s +%.2f units", *lineItem.ProductID, lineItem.Quantity)
				}
			}
		}

		// Mark inventory as synced
		po.InventorySynced = true
		now := time.Now()
		po.InventorySyncDate = &now
	}

	po.Status = status
	po.UpdatedAt = time.Now()

	// Save all fields, including InventorySynced and InventorySyncDate
	_, err = s.poRepo.Update(id, po)
	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return s.GetPurchaseOrder(id)
}

// ReorderProductGroup creates a new purchase order by reordering a product group with custom quantities
func (s *purchaseOrderService) ReorderProductGroup(reorderInput *input.ReorderProductGroupInput, userID string) (*output.PurchaseOrderOutput, error) {
	// Validate vendor
	vendor, err := s.vendorRepo.FindByID(reorderInput.VendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	// Validate delivery address type
	if reorderInput.DeliveryAddressType != "organization" && reorderInput.DeliveryAddressType != "customer" {
		return nil, errors.New("invalid delivery address type")
	}

	var customer *models.Customer
	if reorderInput.DeliveryAddressType == "customer" {
		if reorderInput.CustomerID == nil {
			return nil, errors.New("customer_id is required for customer delivery")
		}
		customer, err = s.customerRepo.FindByID(*reorderInput.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}

	// Validate product group exists
	productGroup, err := s.productGroupRepo.FindByID(reorderInput.ProductGroupID)
	if err != nil {
		return nil, errors.New("product group not found")
	}

	// Validate tax if provided
	var tax *models.Tax
	if reorderInput.TaxID != nil {
		tax, err = s.taxRepo.FindByID(*reorderInput.TaxID)
		if err != nil {
			return nil, errors.New("tax not found")
		}
	}

	// Get user for tracking
	var createdByUserName string
	var createdByCompanyID uint
	var createdByCompanyName string

	if userID != "" {
		var userIDUint uint
		_, err := fmt.Sscanf(userID, "%d", &userIDUint)
		if err == nil {
			user, err := s.userRepo.GetByID(userIDUint)
			if err == nil && user != nil {
				// Set user name from email or username
				if user.Email != nil {
					createdByUserName = *user.Email
				} else if user.Username != nil {
					createdByUserName = *user.Username
				}

				// Set company details if user has a company
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

	// Convert reorder components to line items
	lineItems := make([]models.PurchaseOrderLineItem, len(reorderInput.Components))
	subTotal := float64(0)

	for i, component := range reorderInput.Components {
		amount := component.Quantity * component.Rate
		lineItems[i] = models.PurchaseOrderLineItem{
			ProductID:      &component.ProductID,
			ProductName:    component.ProductName,
			SKU:            component.VariantSku,
			ProductGroupID: &reorderInput.ProductGroupID,
			IsProductGroup: false, // Individual component line items
			Account:        "PURCHASE_EXPENSE",
			Quantity:       component.Quantity,
			Rate:           component.Rate,
			Amount:         amount,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		subTotal += amount
	}

	// Calculate discount
	discountAmount := float64(0)
	if reorderInput.Discount > 0 {
		if reorderInput.DiscountType == "percentage" {
			discountAmount = subTotal * (reorderInput.Discount / 100)
		} else {
			discountAmount = reorderInput.Discount
		}
	}

	// Calculate tax amount
	taxAmount := float64(0)
	if tax != nil {
		taxAmount = (subTotal - discountAmount) * (tax.Rate / 100)
	}

	// Calculate totals
	total := subTotal - discountAmount + taxAmount + reorderInput.Adjustment

	// Generate unique PO number
	poNumber := fmt.Sprintf("PO-%s-%04d", time.Now().Format("20060102"), s.generatePOSequence())

	// Create new purchase order
	po := &models.PurchaseOrder{
		ID:                   uuid.New().String(),
		PurchaseOrderNumber:  poNumber,
		VendorID:             reorderInput.VendorID,
		Vendor:               vendor,
		DeliveryAddressType:  reorderInput.DeliveryAddressType,
		DeliveryAddressID:    reorderInput.DeliveryAddressID,
		OrganizationName:     reorderInput.OrganizationName,
		OrganizationAddress:  reorderInput.OrganizationAddress,
		CustomerID:           reorderInput.CustomerID,
		Customer:             customer,
		PODate:               reorderInput.Date,
		DeliveryDate:         reorderInput.DeliveryDate,
		PaymentTerms:         domain.PaymentTerms(reorderInput.PaymentTerms),
		ShipmentPreference:   reorderInput.ShipmentPreference,
		LineItems:            lineItems,
		SubTotal:             subTotal,
		Discount:             reorderInput.Discount,
		DiscountType:         reorderInput.DiscountType,
		TaxID:                reorderInput.TaxID,
		Tax:                  tax,
		TaxAmount:            taxAmount,
		Adjustment:           reorderInput.Adjustment,
		Total:                total,
		Notes:                reorderInput.Notes,
		TermsAndConditions:   reorderInput.TermsAndConditions,
		Status:               domain.PurchaseOrderStatusDraft,
		Attachments:          reorderInput.Attachments,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CreatedBy:            userID,
		CreatedByUserName:    createdByUserName,
		CreatedByCompanyID:   createdByCompanyID,
		CreatedByCompanyName: createdByCompanyName,
		POPaymentStatus:      domain.PaymentStatusPending,
		PaidAmount:           0,
		RemainingAmount:      total,
	}

	// Save to database
	createdPO, err := s.poRepo.Create(po)
	if err != nil {
		log.Printf("[REORDER] Error creating purchase order for product group %s: %v", reorderInput.ProductGroupID, err)
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	log.Printf("[REORDER] Successfully created purchase order %s for product group %s with %d components", createdPO.PurchaseOrderNumber, productGroup.Name, len(lineItems))

	// Convert to output DTO
	return output.ToPurchaseOrderOutput(createdPO)
}

func (s *purchaseOrderService) generatePOSequence() int {
	var count int64
	today := time.Now().Format("2006-01-02")

	s.poRepo.GetDB().Where("DATE(created_at) = ?", today).Model(&models.PurchaseOrder{}).Count(&count)

	return int(count) + 1
}

// validateNoDuplicateItems checks that no item_id appears more than once in the line items
// This prevents adding the same item with different quantities in the same purchase order
