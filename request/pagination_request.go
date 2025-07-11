package request

type PaginationRequest struct {
	Page         int64  `query:"page"`
	Limit        int64  `query:"limit"`
	MinAge       int64  `query:"minAge"`
	NameContains string `query:"nameContains"`
}
