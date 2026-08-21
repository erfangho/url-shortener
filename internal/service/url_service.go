package service

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/erfangho/url-shortener/internal/model"
	"gorm.io/gorm"
)

type URLRepositoryInterface interface {
	Create(url *model.URL) error
	FindByShortCode(shortCode string) (*model.URL, error)
	FindByOriginalURL(originalUrl string) (*model.URL, error)
	IncrementClickCount(shortCode string) error
}

type URLService struct {
	repo URLRepositoryInterface
}

func NewURLService(repo URLRepositoryInterface) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateURL(originalURL string) (*model.URL, error) {
	randomExist := true
	var randomChars string

	for randomExist {
		b := make([]byte, 6)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}

		randomChars = fmt.Sprintf("%x", b)[:6]

		_, err = s.repo.FindByShortCode(randomChars)

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

	err := s.repo.Create(newURL)
	if err != nil {
		return nil, err
	}

	return newURL, nil
}

func (s *URLService) GetURL(shortCode string) (*model.URL, error) {
	return s.repo.FindByShortCode(shortCode)
}

func (s *URLService) Redirect(shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(shortCode)

	if err != nil {
		return "", err
	}

	go s.repo.IncrementClickCount(shortCode)

	return url.OriginalURL, nil
}
