package output

type BrandOutput struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ListBrandOutput struct{
	Brands []BrandOutput `json:"brands"`
	Total  int64         `json:"total"`
}
