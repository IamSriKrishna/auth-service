package output

import (
	"github.com/bbapp-org/auth-service/app/models"
)

// convertVariantDetails safely converts VariantDetails map values to strings
func convertVariantDetails(details models.VariantDetails) map[string]string {
	if len(details) == 0 {
		return nil
	}
	result := make(map[string]string)
	for k, v := range details {
		result[k] = v
	}
	return result
}

// CustomerInfo represents basic customer information
type CustomerInfo struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	CompanyName string `json:"company_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

// VendorInfo represents basic vendor information
type VendorInfo struct {
	ID           uint   `json:"id"`
	DisplayName  string `json:"display_name"`
	CompanyName  string `json:"company_name,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	WorkPhone    string `json:"work_phone,omitempty"`
}

// ItemInfo represents basic item information
type ItemInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SKU  string `json:"sku,omitempty"`
}

// VariantInfo represents basic variant information
type VariantInfo struct {
	ID           uint              `json:"id"`
	SKU          string            `json:"sku"`
	AttributeMap map[string]string `json:"attribute_map"`
}

// TaxInfo represents basic tax information
type TaxInfo struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	TaxType string  `json:"tax_type"`
	Rate    float64 `json:"rate"`
}

// SalespersonInfo represents basic salesperson information
type SalespersonInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
