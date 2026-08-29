package model

import "gorm.io/gorm"

type URL struct {
	gorm.Model

	OriginalURL string       `json:"original_url" gorm:"not null"`
	ShortCode   string       `json:"short_code" gorm:"uniqueIndex;not null"`
	ClickCount  int          `json:"click_count" gorm:"default:0"`
	ClickEvents []ClickEvent `json:"click_events" gorm:"foreignKey:URLID;constraint:OnDelete:CASCADE"`
}

type CreateURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type GetAllURLsResponse struct {
	URLs       []URL      `json:"data"`
	Pagination Pagination `json:"pagination"`
}
