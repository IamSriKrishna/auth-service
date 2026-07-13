package repo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bbapp-org/auth-service/app/models"
	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func invoicePreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Customer").
		Preload("Salesperson").
		Preload("Tax").
		Preload("LineItems").
		Preload("LineItems.Product")
}

// Invoice company ownership is resolved through:
//
// invoices.created_by -> users.id -> users.company_id
//
// This requires no company_id column in invoices.
func invoiceCompanyQuery(
	db *gorm.DB,
	companyID uint,
) *gorm.DB {
	return db.
		Joins(`
			JOIN users invoice_creator
				ON invoice_creator.id =
					CAST(invoices.created_by AS UNSIGNED)
		`).
		Where("invoice_creator.company_id = ?", companyID)
}

func (r *invoiceRepository) Create(
	invoice *models.Invoice,
) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) FindByID(
	id string,
) (*models.Invoice, error) {
	var invoice models.Invoice

	err := invoicePreloads(r.db).
		Where("invoices.id = ?", id).
		First(&invoice).
		Error

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) FindByIDAndCompany(
	id string,
	companyID uint,
) (*models.Invoice, error) {
	var invoice models.Invoice

	query := invoiceCompanyQuery(
		invoicePreloads(r.db).
			Model(&models.Invoice{}),
		companyID,
	)

	err := query.
		Select("invoices.*").
		Where("invoices.id = ?", id).
		First(&invoice).
		Error

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) FindAll(
	limit int,
	offset int,
) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	if err := r.db.
		Model(&models.Invoice{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	err := invoicePreloads(r.db).
		Order("invoices.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).
		Error

	return invoices, total, err
}

func (r *invoiceRepository) FindAllByCompany(
	companyID uint,
	limit int,
	offset int,
) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	countQuery := invoiceCompanyQuery(
		r.db.Model(&models.Invoice{}),
		companyID,
	)

	if err := countQuery.
		Distinct("invoices.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := invoiceCompanyQuery(
		invoicePreloads(r.db).
			Model(&models.Invoice{}),
		companyID,
	)

	err := findQuery.
		Select("invoices.*").
		Order("invoices.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).
		Error

	return invoices, total, err
}

func (r *invoiceRepository) Update(
	invoice *models.Invoice,
) error {
	if invoice == nil {
		return gorm.ErrInvalidData
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&models.Invoice{}).
			Where("id = ?", invoice.ID).
			Omit("LineItems").
			Updates(invoice).
			Error; err != nil {
			return err
		}

		if invoice.LineItems != nil {
			if err := tx.
				Where("invoice_id = ?", invoice.ID).
				Delete(&models.InvoiceLineItem{}).
				Error; err != nil {
				return err
			}

			if len(invoice.LineItems) > 0 {
				for index := range invoice.LineItems {
					invoice.LineItems[index].InvoiceID = invoice.ID
				}

				if err := tx.
					Create(&invoice.LineItems).
					Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *invoiceRepository) UpdateByCompany(
	invoice *models.Invoice,
	companyID uint,
) error {
	if invoice == nil {
		return gorm.ErrInvalidData
	}

	if _, err := r.FindByIDAndCompany(
		invoice.ID,
		companyID,
	); err != nil {
		return err
	}

	return r.Update(invoice)
}

func (r *invoiceRepository) Delete(
	id string,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("invoice_id = ?", id).
			Delete(&models.Payment{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			Where("invoice_id = ?", id).
			Delete(&models.InvoiceLineItem{}).
			Error; err != nil {
			return err
		}

		result := tx.
			Where("id = ?", id).
			Delete(&models.Invoice{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func (r *invoiceRepository) DeleteByCompany(
	id string,
	companyID uint,
) error {
	if _, err := r.FindByIDAndCompany(
		id,
		companyID,
	); err != nil {
		return err
	}

	return r.Delete(id)
}

func (r *invoiceRepository) FindByCustomerID(
	customerID string,
	limit int,
	offset int,
) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.
		Model(&models.Invoice{}).
		Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := invoicePreloads(r.db).
		Where("customer_id = ?", customerID).
		Order("invoices.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).
		Error

	return invoices, total, err
}

func (r *invoiceRepository) FindByCustomerIDAndCompany(
	customerID uint,
	companyID uint,
	limit int,
	offset int,
) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	countQuery := invoiceCompanyQuery(
		r.db.Model(&models.Invoice{}),
		companyID,
	).Where("invoices.customer_id = ?", customerID)

	if err := countQuery.
		Distinct("invoices.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := invoiceCompanyQuery(
		invoicePreloads(r.db).
			Model(&models.Invoice{}),
		companyID,
	)

	err := findQuery.
		Select("invoices.*").
		Where("invoices.customer_id = ?", customerID).
		Order("invoices.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).
		Error

	return invoices, total, err
}

func (r *invoiceRepository) FindByStatus(
	status string,
	limit int,
	offset int,
) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.
		Model(&models.Invoice{}).
		Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := invoicePreloads(r.db).
		Where("status = ?", status).
		Order("invoices.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).
		Error

	return invoices, total, err
}

func (r *invoiceRepository) FindByStatusAndCompany(
	status string,
	companyID uint,
	limit int,
	offset int,
) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	countQuery := invoiceCompanyQuery(
		r.db.Model(&models.Invoice{}),
		companyID,
	).Where("invoices.status = ?", status)

	if err := countQuery.
		Distinct("invoices.id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	findQuery := invoiceCompanyQuery(
		invoicePreloads(r.db).
			Model(&models.Invoice{}),
		companyID,
	)

	err := findQuery.
		Select("invoices.*").
		Where("invoices.status = ?", status).
		Order("invoices.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).
		Error

	return invoices, total, err
}

func (r *invoiceRepository) GetNextInvoiceNumber() (
	string,
	error,
) {
	var lastInvoice models.Invoice

	err := r.db.
		Select("invoice_number").
		Where("invoice_number LIKE ?", "INV-%").
		Order("created_at DESC").
		First(&lastInvoice).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "INV-00001", nil
		}
		return "", err
	}

	number := 0
	parts := strings.Split(lastInvoice.InvoiceNumber, "-")
	if len(parts) > 0 {
		parsed, parseErr := strconv.Atoi(parts[len(parts)-1])
		if parseErr == nil {
			number = parsed
		}
	}

	return fmt.Sprintf("INV-%05d", number+1), nil
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
