package utils

type PaginationResult struct {
	Page       int64 `json:"page"`
	TotalPages int64 `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}
