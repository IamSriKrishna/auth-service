package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type VendorPaymentOutput struct {
	ID                   uint                  `json:"id"`
	PaymentNumber        string                `json:"payment_number"`
	PurchaseOrderID      string                `json:"purchase_order_id"`
	PurchaseOrder        *models.PurchaseOrder `json:"purchase_order,omitempty"`
	VendorID             uint                  `json:"vendor_id"`
	Vendor               *VendorInfo           `json:"vendor,omitempty"`
	PaymentMode          string                `json:"payment_mode"`
	Amount               float64               `json:"amount"`
	PaidAmount           float64               `json:"paid_amount"`
	RemainingAmount      float64               `json:"remaining_amount"`
	PaymentStatus        string                `json:"payment_status"`
	PaymentDate          time.Time             `json:"payment_date"`
	ReferenceNumber      string                `json:"reference_number"`
	Notes                string                `json:"notes"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	CreatedByUserName    string                `json:"created_by_user_name"`
	CreatedByCompanyName string                `json:"created_by_company_name"`
}

type VendorPaymentListResponse struct {
	Data  []VendorPaymentOutput `json:"vendor_payments"`
	Total int64                 `json:"total"`
}

func ConvertVendorPaymentToOutput(payment *models.VendorPayment) *VendorPaymentOutput {
	output := &VendorPaymentOutput{
		ID:                   payment.ID,
		PaymentNumber:        payment.PaymentNumber,
		PurchaseOrderID:      payment.PurchaseOrderID,
		VendorID:             payment.VendorID,
		PaymentMode:          string(payment.PaymentMode),
		Amount:               payment.Amount,
		PaidAmount:           payment.PaidAmount,
		RemainingAmount:      payment.RemainingAmount,
		PaymentStatus:        string(payment.PaymentStatus),
		PaymentDate:          payment.PaymentDate,
		ReferenceNumber:      payment.ReferenceNumber,
		Notes:                payment.Notes,
		CreatedAt:            payment.CreatedAt,
		UpdatedAt:            payment.UpdatedAt,
		CreatedByUserName:    payment.CreatedByUserName,
		CreatedByCompanyName: payment.CreatedByCompanyName,
	}

	if payment.Vendor != nil {
		output.Vendor = &VendorInfo{
			ID:           payment.Vendor.ID,
			DisplayName:  payment.Vendor.DisplayName,
			CompanyName:  payment.Vendor.CompanyName,
			EmailAddress: payment.Vendor.EmailAddress,
		}
	}

	return output
}

func ConvertVendorPaymentsToOutput(payments []models.VendorPayment) []VendorPaymentOutput {
	var outputs []VendorPaymentOutput
	for _, payment := range payments {
		outputs = append(outputs, *ConvertVendorPaymentToOutput(&payment))
	}
	return outputs
}
