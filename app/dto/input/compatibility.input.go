package input

type CheckCompatibilityInput struct {
	BottleID uint `json:"bottle_id" validate:"required"`
	CapID    uint `json:"cap_id" validate:"required"`
}