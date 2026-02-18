package output

import "time"

type ItemGroupOutput struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	IsActive    bool                       `json:"is_active"`
	Components  []ItemGroupComponentOutput `json:"components"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

type ItemGroupComponentOutput struct {
	ID             uint        `json:"id"`
	ItemGroupID    string      `json:"item_group_id"`
	ItemID         string      `json:"item_id"`
	Item           *ItemInfo   `json:"item,omitempty"`
	VariantSku     *string     `json:"variant_sku,omitempty"`
	Quantity       float64     `json:"quantity"`
	VariantDetails interface{} `json:"variant_details,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type ItemGroupListOutput struct {
	ItemGroups []ItemGroupOutput `json:"item_groups"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}
