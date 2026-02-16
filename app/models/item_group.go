package models

import "time"

// ItemGroup represents a group/bundle of items (Bill of Materials)
// Example: A 300ml Water Bottle ItemGroup contains 100 bottles + 100 caps
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

// ItemGroupComponent represents a single item in an ItemGroup
// Example: 100 units of "Bottle" variant as a component of "300ml Water Bottle" ItemGroup
type ItemGroupComponent struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	ItemGroupID    string         `json:"item_group_id" gorm:"type:varchar(255);index;not null"`
	ItemID         string         `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item           *Item          `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantID      *uint          `json:"variant_id,omitempty" gorm:"index"`
	Variant        *Variant       `json:"variant,omitempty" gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Quantity       float64        `json:"quantity" gorm:"not null"`
	VariantDetails VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (ItemGroupComponent) TableName() string {
	return "item_group_components"
}
