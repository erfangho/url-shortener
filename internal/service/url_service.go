package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/erfangho/url-shortener/internal/model"
	"gorm.io/gorm"
)

type URLRepositoryInterface interface {
	Create(ctx context.Context, url *model.URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	FindByOriginalURL(ctx context.Context, originalUrl string) (*model.URL, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
	FindAll(ctx context.Context, page, limit int) ([]model.URL, int64, error)
}

type URLService struct {
	repo URLRepositoryInterface
}

func NewURLService(repo URLRepositoryInterface) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateURL(ctx context.Context, originalURL string) (*model.URL, error) {
	randomExist := true
	var randomChars string

	for randomExist {
		b := make([]byte, 6)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}

		randomChars = fmt.Sprintf("%x", b)[:6]

		_, err = s.repo.FindByShortCode(ctx, randomChars)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			randomExist = false
		} else if err != nil {
			return nil, err
		}
	}

	newURL := &model.URL{
		OriginalURL: originalURL,
		ShortCode:   randomChars,
	}

	err := s.repo.Create(ctx, newURL)
	if err != nil {
		return nil, err
	}

	return newURL, nil
}

func (s *URLService) GetURL(ctx context.Context, shortCode string) (*model.URL, error) {
	return s.repo.FindByShortCode(ctx, shortCode)
}

func (s *URLService) Redirect(ctx context.Context, shortCode string) (*model.URL, error) {
	url, err := s.repo.FindByShortCode(ctx, shortCode)

	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetAllURLs(ctx context.Context, page, perPage int) ([]model.URL, int64, int, error) {
	if page == 0 {
		page = 1
	}

	if perPage > 100 {
		perPage = 100
	}

	urls, total, err := s.repo.FindAll(ctx, page, perPage)

	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return urls, total, totalPages, nil
}
