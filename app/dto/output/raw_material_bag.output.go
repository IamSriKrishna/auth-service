package output

import "github.com/bbapp-org/auth-service/app/models"

type RawMaterialBagOutput struct {
	ID              string  `json:"id"`
	PurchaseOrderID string  `json:"purchase_order_id"`
	PurchaseOrderNo string  `json:"purchase_order_no"`
	VendorID        uint    `json:"vendor_id"`
	VendorName      string  `json:"vendor_name"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	BagNumber       int     `json:"bag_number"`
	ExpectedKg      float64 `json:"expected_kg"`
	ActualKg        float64 `json:"actual_kg"`
	RemainingKg     float64 `json:"remaining_kg"`
	Status          string  `json:"status"`
}

type ReceiveRawMaterialBagsOutput struct {
	PurchaseOrderID string                 `json:"purchase_order_id"`
	ProductID       string                 `json:"product_id"`
	ProductName     string                 `json:"product_name"`
	ExpectedKg      float64                `json:"expected_kg"`
	ActualKg        float64                `json:"actual_kg"`
	ShortageKg      float64                `json:"shortage_kg"`
	ShortageGrams   float64                `json:"shortage_grams"`
	Bags            []RawMaterialBagOutput `json:"bags"`
}

type RawMaterialBagListOutput struct {
	Bags  []RawMaterialBagOutput `json:"bags"`
	Total int64                  `json:"total"`
}

func ToRawMaterialBagOutput(bag *models.RawMaterialBag) RawMaterialBagOutput {
	return RawMaterialBagOutput{
		ID:              bag.ID,
		PurchaseOrderID: bag.PurchaseOrderID,
		PurchaseOrderNo: bag.PurchaseOrderNo,
		VendorID:        bag.VendorID,
		VendorName:      bag.VendorName,
		ProductID:       bag.ProductID,
		ProductName:     bag.ProductName,
		BagNumber:       bag.BagNumber,
		ExpectedKg:      bag.ExpectedKg,
		ActualKg:        bag.ActualKg,
		RemainingKg:     bag.RemainingKg,
		Status:          bag.Status,
	}
}
