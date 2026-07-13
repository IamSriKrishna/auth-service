package models

import "time"

type ProductGroup struct {
	ID          string `json:"id" gorm:"type:varchar(255);primaryKey"`
	CompanyID   uint   `json:"company_id" gorm:"not null;index;uniqueIndex:idx_product_group_company_name"`
	Name        string `json:"name" gorm:"type:varchar(255);not null;uniqueIndex:idx_product_group_company_name"`
	Description string `json:"description" gorm:"type:text"`
	IsActive    bool   `json:"is_active" gorm:"default:true;index"`

	Components []ProductGroupComponent `json:"components" gorm:"foreignKey:ProductGroupID;constraint:OnDelete:CASCADE"`
	Resources  []ProductGroupResource  `json:"resources" gorm:"foreignKey:ProductGroupID;constraint:OnDelete:CASCADE"`

	UserID    uint      `json:"user_id" gorm:"not null;index"`
	CreatedBy uint      `json:"created_by" gorm:"not null;index"`
	UpdatedBy uint      `json:"updated_by" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductGroup) TableName() string { return "product_groups" }

type ProductGroupComponent struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
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

func (ProductGroupComponent) TableName() string { return "product_group_components" }

type ProductGroupResource struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductGroupID string    `json:"product_group_id" gorm:"type:varchar(255);index;not null"`
	ResourceType   string    `json:"resource_type" gorm:"type:varchar(100);not null"`
	Unit           string    `json:"unit" gorm:"type:varchar(50);not null"`
	Quantity       float64   `json:"quantity" gorm:"not null"`
	Cost           float64   `json:"cost" gorm:"not null"`
	Position       int       `json:"position,omitempty" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ProductGroupResource) TableName() string { return "product_group_resources" }
