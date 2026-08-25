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

func (r *UserRepository) FindByUserName(username string) (*model.User, error) {
	var user *model.User

	err := r.db.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}
