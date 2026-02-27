package models

import (
	"time"

	"github.com/bbapp-org/auth-service/app/domain"
)

type Package struct {
	ID            string               `json:"id" gorm:"type:varchar(255);primaryKey"`
	PackageSlipNo string               `json:"package_slip_no" gorm:"column:package_slip_no;type:varchar(100);uniqueIndex;not null"`
	SalesOrderID  string               `json:"sales_order_id" gorm:"type:varchar(255);not null;index"`
	SalesOrder    *SalesOrder          `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CustomerID    uint                 `json:"customer_id" gorm:"not null;index"`
	Customer      *Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PackageDate   time.Time            `json:"package_date" gorm:"not null"`
	Items         []PackageItem        `json:"items" gorm:"foreignKey:PackageID;constraint:OnDelete:CASCADE"`
	Status        domain.PackageStatus `json:"status" gorm:"type:varchar(50);not null;default:'created'"`
	InternalNotes string               `json:"internal_notes" gorm:"type:text"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	CreatedBy     string               `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedBy     string               `json:"updated_by" gorm:"type:varchar(255)"`
}

type PackageItem struct {
	ID               uint                `json:"id" gorm:"primaryKey"`
	PackageID        string              `json:"package_id" gorm:"type:varchar(255);not null;index"`
	SalesOrderItemID uint                `json:"sales_order_item_id" gorm:"not null;index"`
	SalesOrderItem   *SalesOrderLineItem `json:"sales_order_item,omitempty" gorm:"foreignKey:SalesOrderItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ItemID           string              `json:"item_id" gorm:"type:varchar(255);not null;index"`
	Item             *Item               `json:"item,omitempty" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VariantSKU       *string             `json:"variant_sku,omitempty" gorm:"type:varchar(255);index"`
	Variant          *Variant            `json:"variant,omitempty" gorm:"foreignKey:VariantSKU;references:SKU;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	OrderedQty       float64             `json:"ordered_qty" gorm:"not null"`
	PackedQty        float64             `json:"packed_qty" gorm:"not null;default:0"`
	VariantDetails   VariantDetails      `json:"variant_details,omitempty" gorm:"type:json"`
}

func (Package) TableName() string {
	return "packages"
}

func (PackageItem) TableName() string {
	return "package_items"
}
