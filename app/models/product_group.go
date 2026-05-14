package models

import "time"

type ProductGroup struct {
	ID          string                  `json:"id" gorm:"type:varchar(255);primaryKey"`
	Name        string                  `json:"name" gorm:"type:varchar(255);not null;unique"`
	Description string                  `json:"description" gorm:"type:text"`
	IsActive    bool                    `json:"is_active" gorm:"default:true"`
	Components  []ProductGroupComponent `json:"components" gorm:"foreignKey:ProductGroupID;constraint:OnDelete:CASCADE"`
	Resources   []ProductGroupResource  `json:"resources" gorm:"foreignKey:ProductGroupID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

func (ProductGroup) TableName() string {
	return "product_groups"
}

type ProductGroupComponent struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	ProductGroupID string         `json:"product_group_id" gorm:"type:varchar(255);index;not null"`
	ProductID      string         `json:"product_id" gorm:"type:varchar(255);not null;index"`
	Product        *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantSku     *string        `json:"variant_sku,omitempty" gorm:"type:varchar(255);index"`
	Quantity       float64        `json:"quantity" gorm:"not null"`
	Position       int            `json:"position,omitempty" gorm:"default:0"`
	VariantDetails VariantDetails `json:"variant_details,omitempty" gorm:"type:json"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (ProductGroupComponent) TableName() string {
	return "product_group_components"
}

type ProductGroupResource struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	ProductGroupID string    `json:"product_group_id" gorm:"type:varchar(255);index;not null"`
	ResourceType   string    `json:"resource_type" gorm:"type:varchar(100);not null"` // e.g., "electricity", "water"
	Unit           string    `json:"unit" gorm:"type:varchar(50);not null"`           // e.g., "watt", "liter", "kwh"
	Quantity       float64   `json:"quantity" gorm:"not null"`                        // e.g., 100 (watts), 5 (liters)
	Cost           float64   `json:"cost" gorm:"not null"`                            // Cost of the resource
	Position       int       `json:"position,omitempty" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ProductGroupResource) TableName() string {
	return "product_group_resources"
}
