package repo

import (
	"fmt"
	
	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *models.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("LineItems", "Customer", "Salesperson", "Tax").Create(invoice).Error; err != nil {
			return err
		}

		if len(invoice.LineItems) > 0 {
			for i := range invoice.LineItems {
				invoice.LineItems[i].InvoiceID = invoice.ID
			}
			if err := tx.Omit("Item", "Variant").Create(&invoice.LineItems).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *invoiceRepository) FindByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.
		Preload("LineItems.Item.ItemDetails").
		Preload("LineItems.Variant.Attributes").
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Where("id = ?", id).
		First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) FindAll(limit, offset int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	if err := r.db.Model(&models.Invoice{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Preload("LineItems.Item.ItemDetails").
		Preload("LineItems.Variant.Attributes").
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *invoiceRepository) Update(invoice *models.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(invoice).
			Omit("LineItems", "Customer", "Salesperson", "Tax").
			Updates(invoice).Error; err != nil {
			return err
		}

		if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceLineItem{}).Error; err != nil {
			return err
		}

		if len(invoice.LineItems) > 0 {
			for i := range invoice.LineItems {
				invoice.LineItems[i].InvoiceID = invoice.ID
			}
			if err := tx.Omit("Item", "Variant").Create(&invoice.LineItems).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *invoiceRepository) Delete(id string) error {
	return r.db.Delete(&models.Invoice{}, "id = ?", id).Error
}

func (r *invoiceRepository) FindByCustomerID(customerID string, limit, offset int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.Model(&models.Invoice{}).Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("LineItems.Item.ItemDetails").
		Preload("LineItems.Variant.Attributes").
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *invoiceRepository) FindByStatus(status string, limit, offset int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.Model(&models.Invoice{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("LineItems.Item.ItemDetails").
		Preload("LineItems.Variant.Attributes").
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *invoiceRepository) GetNextInvoiceNumber() (string, error) {
	var lastInvoice models.Invoice
	err := r.db.Order("created_at DESC").First(&lastInvoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "INV-000001", nil
		}
		return "", err
	}

	var number int
	_, err = fmt.Sscanf(lastInvoice.InvoiceNumber, "INV-%d", &number)
	if err != nil {
		return "INV-000001", nil
	}

	return fmt.Sprintf("INV-%06d", number+1), nil
}

type salespersonRepository struct {
	db *gorm.DB
}

func NewSalespersonRepository(db *gorm.DB) SalespersonRepository {
	return &salespersonRepository{db: db}
}

func (r *salespersonRepository) Create(salesperson *models.Salesperson) error {
	return r.db.Create(salesperson).Error
}

func (r *salespersonRepository) FindByID(id uint) (*models.Salesperson, error) {
	var salesperson models.Salesperson
	err := r.db.Where("id = ?", id).First(&salesperson).Error
	if err != nil {
		return nil, err
	}
	return &salesperson, nil
}

func (r *salespersonRepository) FindAll(limit, offset int) ([]models.Salesperson, int64, error) {
	var salespersons []models.Salesperson
	var total int64

	if err := r.db.Model(&models.Salesperson{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&salespersons).Error
	if err != nil {
		return nil, 0, err
	}

	return salespersons, total, nil
}

func (r *salespersonRepository) Update(salesperson *models.Salesperson) error {
	return r.db.Model(salesperson).Updates(salesperson).Error
}

func (r *salespersonRepository) Delete(id uint) error {
	return r.db.Delete(&models.Salesperson{}, "id = ?", id).Error
}

func (r *salespersonRepository) FindByEmail(email string) (*models.Salesperson, error) {
	var salesperson models.Salesperson
	err := r.db.Where("email = ?", email).First(&salesperson).Error
	if err != nil {
		return nil, err
	}
	return &salesperson, nil
}

type taxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) Create(tax *models.Tax) error {
	return r.db.Create(tax).Error
}

func (r *taxRepository) FindByID(id uint) (*models.Tax, error) {
	var tax models.Tax
	err := r.db.Where("id = ?", id).First(&tax).Error
	if err != nil {
		return nil, err
	}
	return &tax, nil
}

func (r *taxRepository) FindAll(limit, offset int) ([]models.Tax, int64, error) {
	var taxes []models.Tax
	var total int64

	if err := r.db.Model(&models.Tax{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&taxes).Error
	if err != nil {
		return nil, 0, err
	}

	return taxes, total, nil
}

func (r *taxRepository) Update(tax *models.Tax) error {
	return r.db.Model(tax).Updates(tax).Error
}

func (r *taxRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tax{}, "id = ?", id).Error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Invoice").Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByInvoiceID(invoiceID string) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("invoice_id = ?", invoiceID).Order("payment_date DESC").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payment{}, "id = ?", id).Error
}
