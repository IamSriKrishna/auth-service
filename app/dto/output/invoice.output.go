package output

import (
	"time"

	"github.com/bbapp-org/auth-service/app/models"
)

type InvoiceOutput struct {
	ID                  string                     `json:"id"`
	InvoiceNumber       string                     `json:"invoice_number"`
	CustomerID          uint                       `json:"customer_id"`
	Customer            *CustomerInfo              `json:"customer,omitempty"`
	OrderNumber         string                     `json:"order_number,omitempty"`
	InvoiceDate         time.Time                  `json:"invoice_date"`
	Terms               string                     `json:"terms"`
	DueDate             time.Time                  `json:"due_date"`
	SalespersonID       *uint                      `json:"salesperson_id,omitempty"`
	Salesperson         *SalespersonInfo           `json:"salesperson,omitempty"`
	Subject             string                     `json:"subject,omitempty"`
	LineItems           []InvoiceLineItemOutput    `json:"line_items"`
	SubTotal            float64                    `json:"sub_total"`
	ShippingCharges     float64                    `json:"shipping_charges"`
	TaxType             *string                    `json:"tax_type,omitempty"`
	TaxID               *uint                      `json:"tax_id,omitempty"`
	Tax                 *TaxInfo                   `json:"tax,omitempty"`
	TaxAmount           float64                    `json:"tax_amount"`
	Adjustment          float64                    `json:"adjustment"`
	Total               float64                    `json:"total"`
	CustomerNotes       string                     `json:"customer_notes,omitempty"`
	TermsAndConditions  string                     `json:"terms_and_conditions,omitempty"`
	Status              string                     `json:"status"`
	Attachments         []string                   `json:"attachments,omitempty"`
	PaymentReceived     bool                       `json:"payment_received"`
	Payments            []PaymentOutput            `json:"payments,omitempty"`
	PaymentSplits       []PaymentSplitOutput       `json:"payment_splits,omitempty"`
	EmailCommunications []EmailCommunicationOutput `json:"email_communications,omitempty"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	CreatedBy           string                     `json:"created_by,omitempty"`
	UpdatedBy           string                     `json:"updated_by,omitempty"`
}

type InvoiceLineItemOutput struct {
	ID             uint              `json:"id"`
	ItemID         string            `json:"item_id"`
	Item           *ItemInfo         `json:"item,omitempty"`
	VariantSKU     *string           `json:"variant_sku,omitempty"`
	Variant        *VariantInfo      `json:"variant,omitempty"`
	Description    string            `json:"description,omitempty"`
	Quantity       float64           `json:"quantity"`
	Rate           float64           `json:"rate"`
	Amount         float64           `json:"amount"`
	VariantDetails map[string]string `json:"variant_details,omitempty"`
}

type InvoiceListOutput struct {
	Invoices []InvoiceOutput `json:"invoices"`
	Total    int64           `json:"total"`
}

type SalespersonOutput struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SalespersonListOutput struct {
	Salespersons []SalespersonOutput `json:"salespersons"`
	Total        int64               `json:"total"`
}

type TaxOutput struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	TaxType   string    `json:"tax_type"`
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaxListOutput struct {
	Taxes []TaxOutput `json:"taxes"`
	Total int64       `json:"total"`
}

type PaymentOutput struct {
	ID          uint      `json:"id"`
	InvoiceID   string    `json:"invoice_id"`
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	PaymentMode string    `json:"payment_mode"`
	Reference   string    `json:"reference,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by,omitempty"`
}

type PaymentListOutput struct {
	Payments []PaymentOutput `json:"payments"`
	Total    int64           `json:"total"`
}

func ToInvoiceOutput(invoice *models.Invoice) (*InvoiceOutput, error) {
	lineItems := make([]InvoiceLineItemOutput, len(invoice.LineItems))
	for i, item := range invoice.LineItems {
		lineItemOutput := InvoiceLineItemOutput{
			ID:             item.ID,
			ItemID:         item.ItemID,
			VariantSKU:     item.VariantSKU,
			Description:    item.Description,
			Quantity:       item.Quantity,
			Rate:           item.Rate,
			Amount:         item.Amount,
			VariantDetails: convertVariantDetails(item.VariantDetails),
		}

		if item.Item != nil {
			lineItemOutput.Item = &ItemInfo{
				ID:   item.Item.ID,
				Name: item.Item.Name,
				SKU:  item.Item.ItemDetails.SKU,
			}
		}

		if item.Variant != nil {
			attributeMap := make(map[string]string)
			for _, attr := range item.Variant.Attributes {
				attributeMap[attr.Key] = attr.Value
			}
			lineItemOutput.Variant = &VariantInfo{
				ID:           item.Variant.ID,
				SKU:          item.Variant.SKU,
				AttributeMap: attributeMap,
			}
		}

		lineItems[i] = lineItemOutput
	}

	output := &InvoiceOutput{
		ID:                 invoice.ID,
		InvoiceNumber:      invoice.InvoiceNumber,
		CustomerID:         invoice.CustomerID,
		OrderNumber:        invoice.OrderNumber,
		InvoiceDate:        invoice.InvoiceDate,
		Terms:              string(invoice.Terms),
		DueDate:            invoice.DueDate,
		SalespersonID:      invoice.SalespersonID,
		Subject:            invoice.Subject,
		LineItems:          lineItems,
		SubTotal:           invoice.SubTotal,
		ShippingCharges:    invoice.ShippingCharges,
		TaxType:            (*string)(nil),
		TaxID:              invoice.TaxID,
		TaxAmount:          invoice.TaxAmount,
		Adjustment:         invoice.Adjustment,
		Total:              invoice.Total,
		CustomerNotes:      invoice.CustomerNotes,
		TermsAndConditions: invoice.TermsAndConditions,
		Status:             string(invoice.Status),
		Attachments:        invoice.Attachments,
		PaymentReceived:    invoice.PaymentReceived,
		CreatedAt:          invoice.CreatedAt,
		UpdatedAt:          invoice.UpdatedAt,
		CreatedBy:          invoice.CreatedBy,
		UpdatedBy:          invoice.UpdatedBy,
	}

	if invoice.TaxType != "" {
		taxTypeStr := string(invoice.TaxType)
		output.TaxType = &taxTypeStr
	}

	if invoice.Customer != nil {
		output.Customer = &CustomerInfo{
			ID:          invoice.Customer.ID,
			DisplayName: invoice.Customer.DisplayName,
			CompanyName: invoice.Customer.CompanyName,
			Email:       invoice.Customer.EmailAddress,
			Phone:       invoice.Customer.WorkPhone,
		}
	}

	if invoice.Salesperson != nil {
		output.Salesperson = &SalespersonInfo{
			ID:    invoice.Salesperson.ID,
			Name:  invoice.Salesperson.Name,
			Email: invoice.Salesperson.Email,
		}
	}

	if invoice.Tax != nil {
		output.Tax = &TaxInfo{
			ID:      invoice.Tax.ID,
			Name:    invoice.Tax.Name,
			TaxType: string(invoice.Tax.TaxType),
			Rate:    invoice.Tax.Rate,
		}
	}

	if len(invoice.PaymentSplits) > 0 {
		paymentSplits := make([]PaymentSplitOutput, len(invoice.PaymentSplits))
		for i, split := range invoice.PaymentSplits {
			paymentSplits[i] = *ToPaymentSplitOutput(&split)
		}
		output.PaymentSplits = paymentSplits
	}

	if len(invoice.Payments) > 0 {
		payments := make([]PaymentOutput, len(invoice.Payments))
		for i, payment := range invoice.Payments {
			payments[i] = *ToPaymentOutput(&payment)
		}
		output.Payments = payments
	}

	if len(invoice.EmailCommunications) > 0 {
		emails := make([]EmailCommunicationOutput, len(invoice.EmailCommunications))
		for i, email := range invoice.EmailCommunications {
			emails[i] = *ToEmailCommunicationOutput(&email)
		}
		output.EmailCommunications = emails
	}

	return output, nil
}

func ToInvoiceListOutput(invoices []models.Invoice, total int64) (*InvoiceListOutput, error) {
	outputs := make([]InvoiceOutput, len(invoices))
	for i, invoice := range invoices {
		output, err := ToInvoiceOutput(&invoice)
		if err != nil {
			return nil, err
		}
		outputs[i] = *output
	}

	return &InvoiceListOutput{
		Invoices: outputs,
		Total:    total,
	}, nil
}

func ToSalespersonOutput(salesperson *models.Salesperson) *SalespersonOutput {
	return &SalespersonOutput{
		ID:        salesperson.ID,
		Name:      salesperson.Name,
		Email:     salesperson.Email,
		CreatedAt: salesperson.CreatedAt,
		UpdatedAt: salesperson.UpdatedAt,
	}
}

func ToSalespersonListOutput(salespersons []models.Salesperson, total int64) *SalespersonListOutput {
	outputs := make([]SalespersonOutput, len(salespersons))
	for i, sp := range salespersons {
		outputs[i] = *ToSalespersonOutput(&sp)
	}

	return &SalespersonListOutput{
		Salespersons: outputs,
		Total:        total,
	}
}

func ToTaxOutput(tax *models.Tax) *TaxOutput {
	return &TaxOutput{
		ID:        tax.ID,
		Name:      tax.Name,
		TaxType:   tax.TaxType,
		Rate:      tax.Rate,
		CreatedAt: tax.CreatedAt,
		UpdatedAt: tax.UpdatedAt,
	}
}

func ToTaxListOutput(taxes []models.Tax, total int64) *TaxListOutput {
	outputs := make([]TaxOutput, len(taxes))
	for i, tax := range taxes {
		outputs[i] = *ToTaxOutput(&tax)
	}

	return &TaxListOutput{
		Taxes: outputs,
		Total: total,
	}
}

func ToPaymentOutput(payment *models.Payment) *PaymentOutput {
	return &PaymentOutput{
		ID:          payment.ID,
		InvoiceID:   payment.InvoiceID,
		PaymentDate: payment.PaymentDate,
		Amount:      payment.Amount,
		PaymentMode: payment.PaymentMode,
		Reference:   payment.Reference,
		Notes:       payment.Notes,
		CreatedAt:   payment.CreatedAt,
		CreatedBy:   payment.CreatedBy,
	}
}

func ToPaymentListOutput(payments []models.Payment, total int64) *PaymentListOutput {
	outputs := make([]PaymentOutput, len(payments))
	for i, payment := range payments {
		outputs[i] = *ToPaymentOutput(&payment)
	}

	return &PaymentListOutput{
		Payments: outputs,
		Total:    total,
	}
}

type PaymentSplitOutput struct {
	ID             uint      `json:"id"`
	InvoiceID      string    `json:"invoice_id"`
	PaymentMode    string    `json:"payment_mode"`
	DepositTo      string    `json:"deposit_to"`
	AmountReceived float64   `json:"amount_received"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by,omitempty"`
}

type EmailCommunicationOutput struct {
	ID           uint       `json:"id"`
	InvoiceID    string     `json:"invoice_id"`
	EmailAddress string     `json:"email_address"`
	Subject      string     `json:"subject"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    string     `json:"created_by,omitempty"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
}

func ToPaymentSplitOutput(split *models.PaymentSplit) *PaymentSplitOutput {
	return &PaymentSplitOutput{
		ID:             split.ID,
		InvoiceID:      split.InvoiceID,
		PaymentMode:    split.PaymentMode,
		DepositTo:      split.DepositTo,
		AmountReceived: split.AmountReceived,
		CreatedAt:      split.CreatedAt,
		CreatedBy:      split.CreatedBy,
	}
}

func ToEmailCommunicationOutput(email *models.EmailCommunication) *EmailCommunicationOutput {
	return &EmailCommunicationOutput{
		ID:           email.ID,
		InvoiceID:    email.InvoiceID,
		EmailAddress: email.EmailAddress,
		Subject:      email.Subject,
		Status:       email.Status,
		CreatedAt:    email.CreatedAt,
		CreatedBy:    email.CreatedBy,
		SentAt:       email.SentAt,
	}
}
