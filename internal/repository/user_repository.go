package repository

import (
	"github.com/erfangho/url-shortener/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *model.User) error {
	result := r.db.Create(user)

	return result.Error
}
