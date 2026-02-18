package input

type CreateItemGroupInput struct {
	Name        string                    `json:"name" validate:"required"`
	Description string                    `json:"description"`
	IsActive    bool                      `json:"is_active"`
	Components  []ItemGroupComponentInput `json:"components" validate:"required,dive"`
}

type ItemGroupComponentInput struct {
	ItemID         string      `json:"item_id" validate:"required"`
	VariantSku     *string     `json:"variant_sku,omitempty"`
	Quantity       float64     `json:"quantity" validate:"required,gt=0"`
	VariantDetails interface{} `json:"variant_details,omitempty"`
}

type UpdateItemGroupInput struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	IsActive    *bool                     `json:"is_active"`
	Components  []ItemGroupComponentInput `json:"components"`
}
