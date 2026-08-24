package model

import "gorm.io/gorm"

type URL struct {
	gorm.Model

	OriginalURL string `json:"original_url" gorm:"not null"`
	ShortCode   string `json:"short_code" gorm:"uniqueIndex;not null"`
	ClickCount  int    `json:"click_count" gorm:"default:0"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type CreateURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type GetAllURLsResponse struct {
	URLs       []URL      `json:"data"`
	Pagination Pagination `json:"pagination"`
}
