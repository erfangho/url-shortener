package repository

import (
	"context"

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

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)

	return result.Error
}

func (r *UserRepository) FindByUserName(ctx context.Context, username string) (*model.User, error) {
	var user *model.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindAll(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	offset := (page - 1) * limit

	var users []model.User
	var total int64

	result := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(limit).Find(&users)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	count := r.db.WithContext(ctx).Model(&model.User{}).Count(&total)

	if count.Error != nil {
		return nil, 0, count.Error
	}

	return users, total, nil
}
