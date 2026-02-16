package input

type CreateManufacturerInput struct {
	Name string `json:"name" validate:"required"`
}

type UpdateManufacturerInput struct {
	Name *string `json:"name"`
}

