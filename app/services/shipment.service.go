package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
	"github.com/bbapp-org/auth-service/app/dto/input"
	"github.com/bbapp-org/auth-service/app/dto/output"
	"github.com/bbapp-org/auth-service/app/models"
	"github.com/bbapp-org/auth-service/app/repo"
	"github.com/google/uuid"
)

type ShipmentService interface {
	CreateShipment(shipInput *input.CreateShipmentInput, userID string) (*output.ShipmentOutput, error)
	GetShipment(id string) (*output.ShipmentOutput, error)
	GetAllShipments(limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByCustomer(customerID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByPackage(packageID string, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsBySalesOrder(salesOrderID string, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByStatus(status string, limit, offset int) ([]output.ShipmentOutput, int64, error)
	UpdateShipment(id string, shipInput *input.UpdateShipmentInput, userID string) (*output.ShipmentOutput, error)
	UpdateShipmentStatus(id string, status string, userID string) (*output.ShipmentOutput, error)
	DeleteShipment(id string) error
}

type shipmentService struct {
	shipRepo      repo.ShipmentRepository
	pkgRepo       repo.PackageRepository
	soRepo        repo.SalesOrderRepository
	customerRepo  repo.CustomerRepository
	inventoryRepo repo.InventoryBalanceRepository
}

func NewShipmentService(
	shipRepo repo.ShipmentRepository,
	pkgRepo repo.PackageRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
	inventoryRepo repo.InventoryBalanceRepository,
) ShipmentService {
	return &shipmentService{
		shipRepo:      shipRepo,
		pkgRepo:       pkgRepo,
		soRepo:        soRepo,
		customerRepo:  customerRepo,
		inventoryRepo: inventoryRepo,
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
		// Get current inventory balance
		balance, err := s.inventoryRepo.GetBalance(lineItem.ItemID, lineItem.VariantSKU)
		if err != nil {
			return fmt.Errorf("failed to get inventory balance for item %s: %w", lineItem.ItemID, err)
		}

		// Deduct shipped quantity from available inventory
		if balance.AvailableQuantity < lineItem.Quantity {
			return fmt.Errorf("insufficient inventory for item %s. Required: %f, Available: %f", lineItem.ItemID, lineItem.Quantity, balance.AvailableQuantity)
		}

		balance.AvailableQuantity -= lineItem.Quantity
		balance.CurrentQuantity -= lineItem.Quantity
		balance.UpdatedAt = time.Now()

		if err := s.inventoryRepo.UpdateBalance(balance); err != nil {
			return fmt.Errorf("failed to update inventory balance: %w", err)
		}

		// Create inventory journal entry for shipment
		entry := &models.InventoryJournal{
			ItemID:          lineItem.ItemID,
			VariantSKU:      lineItem.VariantSKU,
			TransactionType: "SHIPMENT_DEDUCTION",
			Quantity:        -lineItem.Quantity,
			ReferenceType:   "SalesOrder",
			ReferenceID:     so.ID,
			ReferenceNo:     so.SalesOrderNumber,
			Notes:           fmt.Sprintf("Inventory deducted for shipment - SO: %s", so.SalesOrderNumber),
			CreatedBy:       userID,
		}

		if err := s.inventoryRepo.CreateJournalEntry(entry); err != nil {
			return fmt.Errorf("failed to create inventory journal: %w", err)
		}
	}

	return nil
}
