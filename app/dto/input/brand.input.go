package input
type CreateBrandInput struct {
	Name    string `json:"name" validate:"required"`
}
type UpdateBrandInput struct {
	Name    *string `json:"name"`
}
