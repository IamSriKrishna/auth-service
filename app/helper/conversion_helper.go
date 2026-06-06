package helper

import (
	"github.com/bbapp-org/auth-service/app/models"
)

func ConvertToBaseUnit(quantity float64, purchaseUnit string, product *models.Product) (float64, string) {
	if product == nil {
		return quantity, purchaseUnit
	}

	baseUnit := product.BaseUnit
	if baseUnit == "" {
		baseUnit = product.ProductDetails.Unit
	}

	if product.PurchaseUnit != "" &&
		product.BaseUnit != "" &&
		purchaseUnit == product.PurchaseUnit &&
		product.ConversionFactor > 0 {
		return quantity * product.ConversionFactor, product.BaseUnit
	}

	return quantity, baseUnit
}
