package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type CustomerPaymentOutput struct {
	ID                   uint               `json:"id"`
	PaymentNumber        string             `json:"payment_number"`
	SalesOrderID         string             `json:"sales_order_id"`
	SalesOrder           *models.SalesOrder `json:"sales_order,omitempty"`
	CustomerID           uint               `json:"customer_id"`
	Customer             *CustomerInfo      `json:"customer,omitempty"`
	PaymentMode          string             `json:"payment_mode"`
	Amount               float64            `json:"amount"`
	ReceivedAmount       float64            `json:"received_amount"`
	RemainingAmount      float64            `json:"remaining_amount"`
	PaymentStatus        string             `json:"payment_status"`
	PaymentDate          time.Time          `json:"payment_date"`
	ReferenceNumber      string             `json:"reference_number"`
	Notes                string             `json:"notes"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
	CreatedByUserName    string             `json:"created_by_user_name"`
	CreatedByCompanyName string             `json:"created_by_company_name"`
}

type CustomerPaymentListResponse struct {
	Data  []CustomerPaymentOutput `json:"customer_payments"`
	Total int64                   `json:"total"`
}

func ConvertCustomerPaymentToOutput(payment *models.CustomerPayment) *CustomerPaymentOutput {
	output := &CustomerPaymentOutput{
		ID:                   payment.ID,
		PaymentNumber:        payment.PaymentNumber,
		SalesOrderID:         payment.SalesOrderID,
		CustomerID:           payment.CustomerID,
		PaymentMode:          string(payment.PaymentMode),
		Amount:               payment.Amount,
		ReceivedAmount:       payment.ReceivedAmount,
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

	if payment.Customer != nil {
		output.Customer = &CustomerInfo{
			ID:           payment.Customer.ID,
			DisplayName:  payment.Customer.DisplayName,
			CompanyName:  payment.Customer.CompanyName,
			EmailAddress: payment.Customer.EmailAddress,
		}
	}

	return output
}

func ConvertCustomerPaymentsToOutput(payments []models.CustomerPayment) []CustomerPaymentOutput {
	var outputs []CustomerPaymentOutput
	for _, payment := range payments {
		outputs = append(outputs, *ConvertCustomerPaymentToOutput(&payment))
	}
	return outputs
}
