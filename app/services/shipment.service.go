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
	// Existing methods retained.
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

	// Company-scoped methods.
	CreateShipmentForCompany(shipInput *input.CreateShipmentInput, userID string, companyID uint) (*output.ShipmentOutput, error)
	GetShipmentForCompany(id string, companyID uint) (*output.ShipmentOutput, error)
	GetAllShipmentsForCompany(companyID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByCustomerForCompany(customerID, companyID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByPackageForCompany(packageID string, companyID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsBySalesOrderForCompany(salesOrderID string, companyID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	GetShipmentsByStatusForCompany(status string, companyID uint, limit, offset int) ([]output.ShipmentOutput, int64, error)
	UpdateShipmentForCompany(id string, shipInput *input.UpdateShipmentInput, userID string, companyID uint) (*output.ShipmentOutput, error)
	UpdateShipmentStatusForCompany(id, status, userID string, companyID uint) (*output.ShipmentOutput, error)
	DeleteShipmentForCompany(id string, companyID uint) error
}

type shipmentService struct {
	shipRepo     repo.ShipmentRepository
	pkgRepo      repo.PackageRepository
	soRepo       repo.SalesOrderRepository
	customerRepo repo.CustomerRepository
}

func NewShipmentService(
	shipRepo repo.ShipmentRepository,
	pkgRepo repo.PackageRepository,
	soRepo repo.SalesOrderRepository,
	customerRepo repo.CustomerRepository,
) ShipmentService {
	return &shipmentService{
		shipRepo:     shipRepo,
		pkgRepo:      pkgRepo,
		soRepo:       soRepo,
		customerRepo: customerRepo,
	}
}

func shipmentsToOutput(
	shipments []models.Shipment,
) ([]output.ShipmentOutput, error) {
	result := make([]output.ShipmentOutput, 0, len(shipments))

	for index := range shipments {
		shipmentOutput, err :=
			output.ToShipmentOutput(&shipments[index])
		if err != nil {
			return nil, err
		}

		result = append(result, *shipmentOutput)
	}

	return result, nil
}

// -----------------------------------------------------------------------------
// Existing methods retained.
// -----------------------------------------------------------------------------

func (s *shipmentService) CreateShipment(
	shipInput *input.CreateShipmentInput,
	userID string,
) (*output.ShipmentOutput, error) {
	if shipInput == nil {
		return nil, errors.New("shipment input cannot be nil")
	}

	pkg, err := s.pkgRepo.FindByID(shipInput.PackageID)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	salesOrder, err :=
		s.soRepo.FindByID(shipInput.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	customer, err :=
		s.customerRepo.FindByID(shipInput.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	return s.createShipmentFromValidatedRecords(
		shipInput,
		userID,
		pkg,
		salesOrder,
		customer,
	)
}

func (s *shipmentService) GetShipment(
	id string,
) (*output.ShipmentOutput, error) {
	shipment, err := s.shipRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	return output.ToShipmentOutput(shipment)
}

func (s *shipmentService) GetAllShipments(
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByCustomer(
	customerID uint,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindByCustomer(
			customerID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByPackage(
	packageID string,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindByPackage(
			packageID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsBySalesOrder(
	salesOrderID string,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindBySalesOrder(
			salesOrderID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByStatus(
	status string,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindByStatus(
			status,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) UpdateShipment(
	id string,
	shipInput *input.UpdateShipmentInput,
	userID string,
) (*output.ShipmentOutput, error) {
	if shipInput == nil {
		return nil, errors.New("shipment input cannot be nil")
	}

	shipment, err := s.shipRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	s.applyShipmentUpdate(shipment, shipInput, userID)

	updatedShipment, err :=
		s.shipRepo.Update(id, shipment)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to update shipment: %w",
			err,
		)
	}

	return output.ToShipmentOutput(updatedShipment)
}

func (s *shipmentService) UpdateShipmentStatus(
	id string,
	status string,
	userID string,
) (*output.ShipmentOutput, error) {
	shipment, err := s.shipRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	return s.updateShipmentStatusRecord(
		shipment,
		status,
		userID,
		func(shipment *models.Shipment) (*models.Shipment, error) {
			return s.shipRepo.Update(id, shipment)
		},
		func(salesOrder *models.SalesOrder) error {
			_, err := s.soRepo.Update(
				salesOrder.ID,
				salesOrder,
			)
			return err
		},
	)
}

func (s *shipmentService) DeleteShipment(id string) error {
	return s.shipRepo.Delete(id)
}

// -----------------------------------------------------------------------------
// Company-scoped methods.
// -----------------------------------------------------------------------------

func (s *shipmentService) CreateShipmentForCompany(
	shipInput *input.CreateShipmentInput,
	userID string,
	companyID uint,
) (*output.ShipmentOutput, error) {
	if shipInput == nil {
		return nil, errors.New("shipment input cannot be nil")
	}
	if companyID == 0 {
		return nil, errors.New("invalid company")
	}

	pkg, err := s.pkgRepo.FindByIDAndCompany(
		shipInput.PackageID,
		companyID,
	)
	if err != nil {
		return nil, errors.New("package not found in your company")
	}

	salesOrder, err := s.soRepo.FindByIDAndCompany(
		shipInput.SalesOrderID,
		companyID,
	)
	if err != nil {
		return nil, errors.New(
			"sales order not found in your company",
		)
	}

	customer, err := s.customerRepo.FindByIDAndCompany(
		shipInput.CustomerID,
		companyID,
	)
	if err != nil {
		return nil, errors.New(
			"customer not found in your company",
		)
	}

	return s.createShipmentFromValidatedRecords(
		shipInput,
		userID,
		pkg,
		salesOrder,
		customer,
	)
}

func (s *shipmentService) GetShipmentForCompany(
	id string,
	companyID uint,
) (*output.ShipmentOutput, error) {
	shipment, err :=
		s.shipRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	return output.ToShipmentOutput(shipment)
}

func (s *shipmentService) GetAllShipmentsForCompany(
	companyID uint,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindAllByCompany(
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByCustomerForCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	if _, err := s.customerRepo.FindByIDAndCompany(
		customerID,
		companyID,
	); err != nil {
		return nil, 0, errors.New(
			"customer not found in your company",
		)
	}

	shipments, total, err :=
		s.shipRepo.FindByCustomerAndCompany(
			customerID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByPackageForCompany(
	packageID string,
	companyID uint,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	if _, err := s.pkgRepo.FindByIDAndCompany(
		packageID,
		companyID,
	); err != nil {
		return nil, 0, errors.New(
			"package not found in your company",
		)
	}

	shipments, total, err :=
		s.shipRepo.FindByPackageAndCompany(
			packageID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsBySalesOrderForCompany(
	salesOrderID string,
	companyID uint,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	if _, err := s.soRepo.FindByIDAndCompany(
		salesOrderID,
		companyID,
	); err != nil {
		return nil, 0, errors.New(
			"sales order not found in your company",
		)
	}

	shipments, total, err :=
		s.shipRepo.FindBySalesOrderAndCompany(
			salesOrderID,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) GetShipmentsByStatusForCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]output.ShipmentOutput, int64, error) {
	shipments, total, err :=
		s.shipRepo.FindByStatusAndCompany(
			status,
			companyID,
			limit,
			offset,
		)
	if err != nil {
		return nil, 0, err
	}

	outputs, err := shipmentsToOutput(shipments)
	if err != nil {
		return nil, 0, err
	}

	return outputs, total, nil
}

func (s *shipmentService) UpdateShipmentForCompany(
	id string,
	shipInput *input.UpdateShipmentInput,
	userID string,
	companyID uint,
) (*output.ShipmentOutput, error) {
	if shipInput == nil {
		return nil, errors.New("shipment input cannot be nil")
	}

	shipment, err :=
		s.shipRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	s.applyShipmentUpdate(shipment, shipInput, userID)

	updatedShipment, err :=
		s.shipRepo.UpdateByCompany(
			id,
			companyID,
			shipment,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to update shipment: %w",
			err,
		)
	}

	return output.ToShipmentOutput(updatedShipment)
}

func (s *shipmentService) UpdateShipmentStatusForCompany(
	id string,
	status string,
	userID string,
	companyID uint,
) (*output.ShipmentOutput, error) {
	shipment, err :=
		s.shipRepo.FindByIDAndCompany(id, companyID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	return s.updateShipmentStatusRecord(
		shipment,
		status,
		userID,
		func(shipment *models.Shipment) (*models.Shipment, error) {
			return s.shipRepo.UpdateByCompany(
				id,
				companyID,
				shipment,
			)
		},
		func(salesOrder *models.SalesOrder) error {
			_, err := s.soRepo.UpdateByCompany(
				salesOrder.ID,
				companyID,
				salesOrder,
			)
			return err
		},
	)
}

func (s *shipmentService) DeleteShipmentForCompany(
	id string,
	companyID uint,
) error {
	if _, err := s.shipRepo.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return errors.New("shipment not found")
	}

	return s.shipRepo.DeleteByCompany(id, companyID)
}

// -----------------------------------------------------------------------------
// Shared logic.
// -----------------------------------------------------------------------------

func (s *shipmentService) createShipmentFromValidatedRecords(
	shipInput *input.CreateShipmentInput,
	userID string,
	pkg *models.Package,
	salesOrder *models.SalesOrder,
	customer *models.Customer,
) (*output.ShipmentOutput, error) {
	if pkg == nil || salesOrder == nil || customer == nil {
		return nil, errors.New(
			"package, sales order, and customer are required",
		)
	}

	if pkg.CustomerID != shipInput.CustomerID ||
		salesOrder.CustomerID != shipInput.CustomerID {
		return nil, errors.New(
			"customer does not match package or sales order",
		)
	}

	if pkg.SalesOrderID != shipInput.SalesOrderID {
		return nil, errors.New(
			"package does not belong to the selected sales order",
		)
	}

	if pkg.Status != domain.PackageStatus("packed") &&
		pkg.Status != domain.PackageStatus("shipped") {
		return nil, errors.New(
			"package must be packed before creating a shipment",
		)
	}

	shipmentNo, err := s.shipRepo.GetNextShipmentNo()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate shipment number: %w",
			err,
		)
	}

	shipment := &models.Shipment{
		ID:              uuid.New().String(),
		ShipmentNo:      shipmentNo,
		PackageID:       shipInput.PackageID,
		Package:         pkg,
		SalesOrderID:    shipInput.SalesOrderID,
		SalesOrder:      salesOrder,
		CustomerID:      shipInput.CustomerID,
		Customer:        customer,
		ShipDate:        shipInput.ShipDate,
		Carrier:         shipInput.Carrier,
		TrackingNo:      shipInput.TrackingNo,
		TrackingURL:     shipInput.TrackingURL,
		ShippingCharges: shipInput.ShippingCharges,
		Status:          domain.ShipmentStatusCreated,
		Notes:           shipInput.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       userID,
		UpdatedBy:       userID,
	}

	createdShipment, err := s.shipRepo.Create(shipment)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create shipment: %w",
			err,
		)
	}

	// Do not deduct stock here.
	// Product stock is already deducted when the package status becomes packed.
	// Deducting again during shipment creation would reduce stock twice.
	log.Printf(
		"[SHIPMENT] Shipment %s created for packed package %s without duplicate stock deduction",
		createdShipment.ShipmentNo,
		pkg.PackageSlipNo,
	)

	return output.ToShipmentOutput(createdShipment)
}

func (s *shipmentService) applyShipmentUpdate(
	shipment *models.Shipment,
	shipInput *input.UpdateShipmentInput,
	userID string,
) {
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
		shipment.ShippingCharges =
			*shipInput.ShippingCharges
	}
	if shipInput.Notes != nil {
		shipment.Notes = *shipInput.Notes
	}
	if shipInput.Status != nil {
		shipment.Status =
			domain.ShipmentStatus(*shipInput.Status)
	}

	shipment.UpdatedBy = userID
	shipment.UpdatedAt = time.Now()
}

func (s *shipmentService) updateShipmentStatusRecord(
	shipment *models.Shipment,
	status string,
	userID string,
	saveShipment func(*models.Shipment) (*models.Shipment, error),
	saveSalesOrder func(*models.SalesOrder) error,
) (*output.ShipmentOutput, error) {
	switch status {
	case "created", "shipped", "in_transit", "delivered", "cancelled":
	default:
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	if status == "shipped" &&
		shipment.Status != domain.ShipmentStatus("shipped") {
		log.Printf(
			"[SHIPMENT] Shipment outbound recorded - Shipment: %s, Carrier: %s, Tracking: %s",
			shipment.ShipmentNo,
			shipment.Carrier,
			shipment.TrackingNo,
		)
	}

	if status == "delivered" &&
		shipment.Status != domain.ShipmentStatus("delivered") {
		salesOrder, err :=
			s.soRepo.FindByID(shipment.SalesOrderID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get sales order: %w",
				err,
			)
		}

		salesOrder.Status = domain.SalesOrderStatus("delivered")
		salesOrder.UpdatedBy = userID
		salesOrder.UpdatedAt = time.Now()

		if err := saveSalesOrder(salesOrder); err != nil {
			return nil, fmt.Errorf(
				"failed to update sales order delivery status: %w",
				err,
			)
		}
	}

	shipment.Status = domain.ShipmentStatus(status)
	shipment.UpdatedBy = userID
	shipment.UpdatedAt = time.Now()

	updatedShipment, err := saveShipment(shipment)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to update shipment status: %w",
			err,
		)
	}

	return output.ToShipmentOutput(updatedShipment)
}
