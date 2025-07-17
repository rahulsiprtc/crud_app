package request

type PaginationRequest struct {
	ID           string `json:"id" query:"id" validate:"required,min=2"`
	Page         int64  `query:"page" validate:"gte=1"`
	Limit        int64  `query:"limit" validate:"gte=1,lte=100"`
	MinAge       int64  `query:"minAge" validate:"gte=18,lte=80"`
	NameContains string `query:"nameContains" validate:"omitempty,min=2"`
}
