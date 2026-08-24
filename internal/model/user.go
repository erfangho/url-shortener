package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"not null;uniqueIndex"`
	Password string `json:"-" gorm:"not null"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required,min=6,max=10"`
	Password string `json:"password" binding:"required,min=8,max=12"`
}
