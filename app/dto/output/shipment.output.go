package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type ShipmentOutput struct {
	ID                string                   `json:"id"`
	ShipmentNo        string                   `json:"shipment_no"`
	PackageID         string                   `json:"package_id,omitempty"`
	Package           *PackageInfo             `json:"package,omitempty"`
	SalesOrderID      string                   `json:"sales_order_id"`
	SalesOrderNo      string                   `json:"sales_order_no,omitempty"`
	SalesOrder        *SalesOrderInfo          `json:"sales_order,omitempty"`
	CustomerID        uint                     `json:"customer_id"`
	CustomerName      string                   `json:"customer_name,omitempty"`
	Customer          *CustomerInfo            `json:"customer,omitempty"`
	ShipDate          time.Time                `json:"ship_date"`
	ShipmentType      string                   `json:"shipment_type,omitempty"`
	ShippingMethod    string                   `json:"shipping_method,omitempty"`
	Carrier           string                   `json:"carrier,omitempty"`
	TrackingNo        string                   `json:"tracking_no,omitempty"`
	TrackingURL       string                   `json:"tracking_url,omitempty"`
	ShippingAddress   string                   `json:"shipping_address,omitempty"`
	EstimatedDelivery *time.Time               `json:"estimated_delivery,omitempty"`
	ShippingCharges   float64                  `json:"shipping_charges"`
	LineItems         []ShipmentLineItemOutput `json:"line_items"`
	TotalItems        float64                  `json:"total_items"`
	StockDeducted     bool                     `json:"stock_deducted"`
	Status            string                   `json:"status"`
	Notes             string                   `json:"notes,omitempty"`
	Message           string                   `json:"message,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	CreatedBy         string                   `json:"created_by,omitempty"`
	UpdatedBy         string                   `json:"updated_by,omitempty"`
}

// ShipmentLineItemOutput is defined in product_group_operations.output.go to avoid duplication

type PackageInfo struct {
	ID            string `json:"id"`
	PackageSlipNo string `json:"package_slip_no"`
	Status        string `json:"status"`
}

func ToShipmentOutput(shipment *models.Shipment) (*ShipmentOutput, error) {
	var packageInfo *PackageInfo
	if shipment.Package != nil {
		packageInfo = &PackageInfo{
			ID:            shipment.Package.ID,
			PackageSlipNo: shipment.Package.PackageSlipNo,
			Status:        string(shipment.Package.Status),
		}
	}

	var soInfo *SalesOrderInfo
	if shipment.SalesOrder != nil {
		soInfo = &SalesOrderInfo{
			ID:                   shipment.SalesOrder.ID,
			SalesOrderNo:         shipment.SalesOrder.SalesOrderNumber,
			CustomerID:           shipment.SalesOrder.CustomerID,
			ReferenceNo:          shipment.SalesOrder.ReferenceNo,
			Date:                 shipment.SalesOrder.Date,
			ExpectedShipmentDate: shipment.SalesOrder.ExpectedShipmentDate,
			Status:               string(shipment.SalesOrder.Status),
		}
	}

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

// CreateShipmentVariantOutput creates a shipment output with variant-specific line items
func CreateShipmentVariantOutput(shipmentNo string, salesOrderID string, salesOrderNo string, customerID uint, customerName string, lineItems []ShipmentLineItemOutput, totalItems float64, stockDeducted bool, message string) *ShipmentOutput {
	return &ShipmentOutput{
		ShipmentNo:    shipmentNo,
		SalesOrderID:  salesOrderID,
		SalesOrderNo:  salesOrderNo,
		CustomerID:    customerID,
		CustomerName:  customerName,
		LineItems:     lineItems,
		TotalItems:    totalItems,
		StockDeducted: stockDeducted,
		Status:        "shipped",
		Message:       message,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
