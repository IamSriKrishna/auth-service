package models

import "time"

type ItemGroup struct {
	ID          string               `json:"id" gorm:"type:varchar(255);primaryKey"`
	Name        string               `json:"name" gorm:"type:varchar(255);not null;unique"`
	Description string               `json:"description" gorm:"type:text"`
	IsActive    bool                 `json:"is_active" gorm:"default:true"`
	Components  []ItemGroupComponent `json:"components" gorm:"foreignKey:ItemGroupID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func (ItemGroup) TableName() string {
	return "item_groups"
}

type ItemGroupComponent struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	ItemGroupID    string         `json:"item_group_id" gorm:"type:varchar(255);index;not null"`
	ItemID         string         `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item           *Item          `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantSku     *string        `json:"variant_sku,omitempty" gorm:"type:varchar(255);index"`
	Quantity       float64        `json:"quantity" gorm:"not null"`
	VariantDetails VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (ItemGroupComponent) TableName() string {
	return "item_group_components"
}
