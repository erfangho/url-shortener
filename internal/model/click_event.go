package model

import "gorm.io/gorm"

type ClickEvent struct {
	gorm.Model

	URLID     uint   `json:"url_id" gorm:"not null"`
	URL       *URL   `json:"url"`
	UserAgent string `json:"user_agent" gorm:"not null"`
	IPAddress string `json:"ip_address" gorm:"not null"`
}
