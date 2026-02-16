package output

type ListBottlesQuery struct {
	SizeID     *uint   `query:"size_id"`
	NeckSizeMM *int    `query:"neck_size_mm"`
	ThreadType *string `query:"thread_type"`
	Page       int     `query:"page" validate:"min=1"`
	PageSize   int     `query:"page_size" validate:"min=1,max=100"`
}

type ListCapsQuery struct {
	NeckSizeMM *int    `query:"neck_size_mm"`
	ThreadType *string `query:"thread_type"`
	Material   *string `query:"material"`
	Color      *string `query:"color"`
	Page       int     `query:"page" validate:"min=1"`
	PageSize   int     `query:"page_size" validate:"min=1,max=100"`
}

type ListProductsQuery struct {
	BottleID *uint    `query:"bottle_id"`
	CapID    *uint    `query:"cap_id"`
	MinMRP   *float64 `query:"min_mrp"`
	MaxMRP   *float64 `query:"max_mrp"`
	Page     int      `query:"page" validate:"min=1"`
	PageSize int      `query:"page_size" validate:"min=1,max=100"`
}

type ProductPaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}
