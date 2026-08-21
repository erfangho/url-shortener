package repository

import (
	"github.com/erfangho/url-shortener/internal/model"
	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Create(url *model.URL) error {
	result := r.db.Create(url)

	return result.Error
}

func (r *URLRepository) FindByShortCode(shortCode string) (*model.URL, error) {
	var url model.URL

	result := r.db.First(&url, "short_code = ?", shortCode)

	if result.Error != nil {
		return nil, result.Error
	}

	return &url, nil
}

func (r *URLRepository) FindByOriginalURL(originalUrl string) (*model.URL, error) {
	var url model.URL

	result := r.db.First(&url, "original_url = ?", originalUrl)

	if result.Error != nil {
		return nil, result.Error
	}

	return &url, nil
}

func (r *URLRepository) IncrementClickCount(shortCode string) error {
	result := r.db.Model(&model.URL{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1))

	return result.Error
}
