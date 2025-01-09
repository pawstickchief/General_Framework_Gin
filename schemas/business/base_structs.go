package business

type PaginationParams struct {
	Page  int `json:"page" binding:"required,min=1"`  // 页码，必填，最小值为 1
	Limit int `json:"limit" binding:"required,min=1"` // 每页条数，必填，最小值为 1
}
