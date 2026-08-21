package model

import "gorm.io/gorm"

type URL struct {
	gorm.Model

	OriginalURL string `json:"original_url" gorm:"not null"`
	ShortCode   string `json:"short_code" gorm:"unique;not null"`
	ClickCount  int    `json:"click_count" gorm:"default:0"`
}

type CreateURLRequest struct {
	URL string `json:"url" binding:"required"`
}
