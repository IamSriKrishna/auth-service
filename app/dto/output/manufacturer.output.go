package output

import "time"

type ManufacturerOutput struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListManufacturersOutput struct {
	Manufacturers []ManufacturerOutput `json:"manufacturers"`
	TotalCount    int                `json:"total_count"`
}