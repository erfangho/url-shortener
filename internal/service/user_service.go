package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/pkg/cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepositoryInterface interface {
	Create(user *model.User) error
	FindByUserName(username string) (*model.User, error)
	FindAll(page, limit int) ([]model.User, int64, error)
}

type UserService struct {
	repo       UserRepositoryInterface
	redisCache *cache.RedisCache
}

func NewUserService(repo UserRepositoryInterface, redisCache *cache.RedisCache) *UserService {
	return &UserService{
		repo:       repo,
		redisCache: redisCache,
	}
}

func (s *UserService) CreateUser(name, username, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindByUserName(username string) (*model.User, error) {
	user, err := s.repo.FindByUserName(username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(user *model.User, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
}

func (s *UserService) GetAllUsers(ctx context.Context, page, perPage int) ([]model.User, int64, int, error) {
	if page == 0 {
		page = 1
	}

	if perPage > 100 {
		perPage = 100
	}

	cacheKey := fmt.Sprintf("users:page:%d:per_page:%d", page, perPage)

	type cachedResponse struct {
		Users      []model.User `json:"users"`
		Total      int64        `json:"total"`
		TotalPages int          `json:"total_pages"`
	}

	var cached cachedResponse
	found, err := s.redisCache.Get(ctx, cacheKey, &cached)
	if err == nil && found {
		return cached.Users, cached.Total, cached.TotalPages, nil
	}

	users, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	s.redisCache.Set(ctx, cacheKey, &cachedResponse{
		Users:      users,
		Total:      total,
		TotalPages: totalPages,
	})

	return users, total, totalPages, nil
}
