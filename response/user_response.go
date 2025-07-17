package response

type UserPaginationResponse struct {
	Users      []UserResponse `json:"users,omitempty"`
	Pagination Pagination     `json:"pagination,omitempty"`
}

type UserResponse struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	Age     int           `json:"age"`
	Message string        `json:"message,omitempty"`
	Data    *UserResponse `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type Pagination struct {
	PerPage      int64 `json:"limit"`
	CurrentPage  int64 `json:"page"`
	LastPage     int64 `json:"lastPage"`
	TotalResults int64 `json:"total"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
