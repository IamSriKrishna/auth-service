package input

type BagReceiveInput struct {
	BagNumber int     `json:"bag_number" validate:"required,gt=0"`
	ActualKg  float64 `json:"actual_kg" validate:"required,gt=0"`
}

type ReceiveRawMaterialBagsInput struct {
	PurchaseOrderID string            `json:"purchase_order_id" validate:"required"`
	ProductID       string            `json:"product_id" validate:"required"`
	ExpectedKgPerBag float64          `json:"expected_kg_per_bag" validate:"required,gt=0"`
	Bags            []BagReceiveInput `json:"bags" validate:"required,min=1"`
}

type UseRawMaterialBagInput struct {
	BagID            string  `json:"bag_id" validate:"required"`
	QuantityKg       float64 `json:"quantity_kg,omitempty"`
	FinishedQuantity float64 `json:"finished_quantity,omitempty"`
}