package input

import "time"

type CreateCustomerPaymentInput struct {
	SalesOrderID    string    `json:"sales_order_id" validate:"required"`
	CustomerID      uint      `json:"customer_id" validate:"required"`
	PaymentMode     string    `json:"payment_mode" validate:"required,oneof=cash online"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	PaymentDate     time.Time `json:"payment_date" validate:"required"`
	ReferenceNumber string    `json:"reference_number"`
	Notes           string    `json:"notes"`
}

type UpdateCustomerPaymentInput struct {
	PaymentMode     *string    `json:"payment_mode" validate:"omitempty,oneof=cash online"`
	Amount          *float64   `json:"amount" validate:"omitempty,gt=0"`
	PaymentDate     *time.Time `json:"payment_date"`
	ReferenceNumber *string    `json:"reference_number"`
	Notes           *string    `json:"notes"`
}

type RecordCustomerPaymentInput struct {
	ReceivedAmount  float64 `json:"received_amount" validate:"required,gt=0"`
	PaymentMode     string  `json:"payment_mode" validate:"required,oneof=cash online"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
}
