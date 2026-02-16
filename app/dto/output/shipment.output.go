package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type ShipmentOutput struct {
	ID              string          `json:"id"`
	ShipmentNo      string          `json:"shipment_no"`
	PackageID       string          `json:"package_id"`
	Package         *PackageInfo    `json:"package,omitempty"`
	SalesOrderID    string          `json:"sales_order_id"`
	SalesOrder      *SalesOrderInfo `json:"sales_order,omitempty"`
	CustomerID      uint            `json:"customer_id"`
	Customer        *CustomerInfo   `json:"customer,omitempty"`
	ShipDate        time.Time       `json:"ship_date"`
	Carrier         string          `json:"carrier,omitempty"`
	TrackingNo      string          `json:"tracking_no,omitempty"`
	TrackingURL     string          `json:"tracking_url,omitempty"`
	ShippingCharges float64         `json:"shipping_charges"`
	Status          string          `json:"status"`
	Notes           string          `json:"notes,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedBy       string          `json:"created_by,omitempty"`
	UpdatedBy       string          `json:"updated_by,omitempty"`
}

type PackageInfo struct {
	ID            string `json:"id"`
	PackageSlipNo string `json:"package_slip_no"`
	Status        string `json:"status"`
}

func ToShipmentOutput(shipment *models.Shipment) (*ShipmentOutput, error) {
	// Build package info
	var packageInfo *PackageInfo
	if shipment.Package != nil {
		packageInfo = &PackageInfo{
			ID:            shipment.Package.ID,
			PackageSlipNo: shipment.Package.PackageSlipNo,
			Status:        string(shipment.Package.Status),
		}
	}

	// Build sales order info
	var soInfo *SalesOrderInfo
	if shipment.SalesOrder != nil {
		soInfo = &SalesOrderInfo{
			ID:               shipment.SalesOrder.ID,
			SalesOrderNo:     shipment.SalesOrder.SalesOrderNumber,
			CustomerID:       shipment.SalesOrder.CustomerID,
			ReferenceNo:      shipment.SalesOrder.ReferenceNo,
			SODate:           shipment.SalesOrder.SODate,
			ExpectedShipDate: shipment.SalesOrder.ExpectedShipmentDate,
			Status:           string(shipment.SalesOrder.Status),
		}
	}

	// Build customer info
	var customerInfo *CustomerInfo
	if shipment.Customer != nil {
		customerInfo = &CustomerInfo{
			ID:          shipment.Customer.ID,
			DisplayName: shipment.Customer.DisplayName,
			CompanyName: shipment.Customer.CompanyName,
			Email:       shipment.Customer.EmailAddress,
			Phone:       shipment.Customer.Mobile,
		}
	}

	output := &ShipmentOutput{
		ID:              shipment.ID,
		ShipmentNo:      shipment.ShipmentNo,
		PackageID:       shipment.PackageID,
		Package:         packageInfo,
		SalesOrderID:    shipment.SalesOrderID,
		SalesOrder:      soInfo,
		CustomerID:      shipment.CustomerID,
		Customer:        customerInfo,
		ShipDate:        shipment.ShipDate,
		Carrier:         shipment.Carrier,
		TrackingNo:      shipment.TrackingNo,
		TrackingURL:     shipment.TrackingURL,
		ShippingCharges: shipment.ShippingCharges,
		Status:          string(shipment.Status),
		Notes:           shipment.Notes,
		CreatedAt:       shipment.CreatedAt,
		UpdatedAt:       shipment.UpdatedAt,
		CreatedBy:       shipment.CreatedBy,
		UpdatedBy:       shipment.UpdatedBy,
	}

	return output, nil
}
