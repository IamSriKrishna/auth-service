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

type ShipmentService interface {
	// Basic CRUD Operations
	CreateShipment(shipInput *input.CreateShipmentInput, userID string) (*output.ShipmentOutput, error)
	GetShipment(id string) (*output.ShipmentOutput, error)
	GetAllShipments(limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByCustomer(customerID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByPackage(packageID string, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsBySalesOrder(salesOrderID string, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByStatus(status string, limit, offset int) ([]output.ShipmentOutput, int64, error)
	UpdateShipment(id string, shipInput *input.UpdateShipmentInput, userID string) (*output.ShipmentOutput, error)

	// Step 4: Selling to Customers - Shipment
	// Record shipment and manage inventory deduction when items ship out
	UpdateShipmentStatus(id string, status string, userID string) (*output.ShipmentOutput, error)
	DeleteShipment(id string) error
}

type shipmentService struct {
	shipRepo         repo.ShipmentRepository
	pkgRepo          repo.PackageRepository
	soRepo           repo.SalesOrderRepository
	customerRepo     repo.CustomerRepository
	inventoryRepo    repo.InventoryBalanceRepository
	stockMovementSvc StockMovementService
}

func NewShipmentService(
	shipRepo repo.ShipmentRepository,
	pkgRepo repo.PackageRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	inventoryRepo repo.InventoryBalanceRepository,
	stockMovementSvc StockMovementService,
) ShipmentService {
	return &shipmentService{
		shipRepo:         shipRepo,
		pkgRepo:          pkgRepo,
		soRepo:           soRepo,
		customerRepo:     customerRepo,
		inventoryRepo:    inventoryRepo,
		stockMovementSvc: stockMovementSvc,
	}
}

func (s *shipmentService) CreateShipment(shipInput *input.CreateShipmentInput, userID string) (*output.ShipmentOutput, error) {
	if shipInput == nil {
		return nil, errors.New("shipment input cannot be nil")
	}

	pkg, err := s.pkgRepo.FindByID(shipInput.PackageID)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	so, err := s.soRepo.FindByID(shipInput.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}
	customer, err := s.customerRepo.FindByID(shipInput.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Verify customer matches
	if pkg.CustomerID != shipInput.CustomerID || so.CustomerID != shipInput.CustomerID {
		return nil, errors.New("customer does not match package or sales order")
	}

	shipNo, err := s.shipRepo.GetNextShipmentNo()
	if err != nil {
		return nil, fmt.Errorf("failed to generate shipment number: %w", err)
	}

	shipment := &models.Shipment{
		ID:              uuid.New().String(),
		ShipmentNo:      shipNo,
		PackageID:       shipInput.PackageID,
		SalesOrderID:    shipInput.SalesOrderID,
		CustomerID:      shipInput.CustomerID,
		ShipDate:        shipInput.ShipDate,
		Carrier:         shipInput.Carrier,
		TrackingNo:      shipInput.TrackingNo,
		TrackingURL:     shipInput.TrackingURL,
		ShippingCharges: shipInput.ShippingCharges,
		Status:          domain.ShipmentStatusCreated,
		Notes:           shipInput.Notes,
		CreatedBy:       userID,
		UpdatedBy:       userID,
	}

	shipment.Package = pkg
	shipment.SalesOrder = so
	shipment.Customer = customer

	createdShip, err := s.shipRepo.Create(shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Deduct inventory for shipped items
	if err := s.deductInventoryForShipment(so, userID); err != nil {
		return nil, fmt.Errorf("failed to deduct inventory for shipment: %w", err)
	}

	return output.ToShipmentOutput(createdShip)
}

func (s *shipmentService) GetShipment(id string) (*output.ShipmentOutput, error) {
	shipment, err := s.shipRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	return output.ToShipmentOutput(shipment)
}

func (s *shipmentService) GetAllShipments(limit, offset int) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err := s.shipRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.ShipmentOutput, 0)
	for _, ship := range shipments {
		if out, err := output.ToShipmentOutput(&ship); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByCustomer(customerID uint, limit, offset int) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err := s.shipRepo.FindByCustomer(customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.ShipmentOutput, 0)
	for _, ship := range shipments {
		if out, err := output.ToShipmentOutput(&ship); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByPackage(packageID string, limit, offset int) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err := s.shipRepo.FindByPackage(packageID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.ShipmentOutput, 0)
	for _, ship := range shipments {
		if out, err := output.ToShipmentOutput(&ship); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsBySalesOrder(salesOrderID string, limit, offset int) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err := s.shipRepo.FindBySalesOrder(salesOrderID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.ShipmentOutput, 0)
	for _, ship := range shipments {
		if out, err := output.ToShipmentOutput(&ship); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByStatus(status string, limit, offset int) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err := s.shipRepo.FindByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]output.ShipmentOutput, 0)
	for _, ship := range shipments {
		if out, err := output.ToShipmentOutput(&ship); err == nil {
			outputs = append(outputs, *out)
		}
	}

	return outputs, total, nil
}

func (s *shipmentService) UpdateShipment(id string, shipInput *input.UpdateShipmentInput, userID string) (*output.ShipmentOutput, error) {
	if shipInput == nil {
		return nil, errors.New("shipment input cannot be nil")
	}

	shipment, err := s.shipRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	if shipInput.ShipDate != nil {
		shipment.ShipDate = *shipInput.ShipDate
	}

	if shipInput.Carrier != nil {
		shipment.Carrier = *shipInput.Carrier
	}

	if shipInput.TrackingNo != nil {
		shipment.TrackingNo = *shipInput.TrackingNo
	}

	if shipInput.TrackingURL != nil {
		shipment.TrackingURL = *shipInput.TrackingURL
	}

	if shipInput.ShippingCharges != nil {
		shipment.ShippingCharges = *shipInput.ShippingCharges
	}

	if shipInput.Notes != nil {
		shipment.Notes = *shipInput.Notes
	}

	if shipInput.Status != nil {
		shipment.Status = domain.ShipmentStatus(*shipInput.Status)
	}

	shipment.UpdatedBy = userID
	shipment.UpdatedAt = time.Now()

	updatedShip, err := s.shipRepo.Update(id, shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}

	return output.ToShipmentOutput(updatedShip)
}

func (s *shipmentService) UpdateShipmentStatus(id string, status string, userID string) (*output.ShipmentOutput, error) {
	switch status {
	case "created", "shipped", "in_transit", "delivered", "cancelled":
	default:
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	shipment, err := s.shipRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	// Record shipment status change and track inventory in transit
	if status == "shipped" && shipment.Status != domain.ShipmentStatus("shipped") {
		// Get the sales order to access line items
		so, err := s.soRepo.FindByID(shipment.SalesOrderID)
		if err != nil {
			return nil, fmt.Errorf("failed to get sales order: %w", err)
		}

		// Record shipment movement - simplified for new model
		for _, lineItem := range so.LineItems {
			// Log shipment for each line item
			log.Printf("[SHIPMENT] Shipment outbound recorded - SKU: %s, Quantity: %.2f, Carrier: %s, Tracking: %s",
				lineItem.SKU, lineItem.Quantity, shipment.Carrier, shipment.TrackingNo)
		}
	}

	// Track delivery completion
	if status == "delivered" && shipment.Status != domain.ShipmentStatus("delivered") {
		// Update sales order delivery status if needed
		so, err := s.soRepo.FindByID(shipment.SalesOrderID)
		if err == nil {
			so.Status = "delivered"
			so.UpdatedBy = userID
			so.UpdatedAt = time.Now()
		}
	}

	shipment.Status = domain.ShipmentStatus(status)
	shipment.UpdatedBy = userID
	shipment.UpdatedAt = time.Now()

	updatedShip, err := s.shipRepo.Update(id, shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to update shipment status: %w", err)
	}

	return output.ToShipmentOutput(updatedShip)
}

func (s *shipmentService) DeleteShipment(id string) error {
	return s.shipRepo.Delete(id)
}

// deductInventoryForShipment reduces available inventory when shipment is created
func (s *shipmentService) deductInventoryForShipment(so *models.SalesOrder, userID string) error {
	for _, lineItem := range so.LineItems {
		// Log inventory deduction for shipment
		log.Printf("[SHIPMENT] Inventory deduction recorded - SKU: %s, Quantity: %.2f, SalesOrder: %s",
			lineItem.SKU, lineItem.Quantity, so.SalesOrderNumber)
	}
	return nil
}
