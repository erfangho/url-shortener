package repository

import (
	"context"
	"errors"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/pkg/cache"
	"gorm.io/gorm"
)

type URLRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewURLRepository(db *gorm.DB, cache *cache.Cache) *URLRepository {
	return &URLRepository{
		db:    db,
		cache: cache,
	}
}

func (r *URLRepository) Create(ctx context.Context, url *model.URL) error {
	result := r.db.WithContext(ctx).Create(url)

	return result.Error
}

func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	cacheValue, exists := r.cache.Get(shortCode)

	if exists {
		url, ok := cacheValue.(*model.URL)

		if !ok {
			return nil, errors.New("cache retrieve failed")
		}

		return url, nil
	}

	var url model.URL

	result := r.db.WithContext(ctx).First(&url, "short_code = ?", shortCode)

	if result.Error != nil {
		return nil, result.Error
	}

	r.cache.Set(shortCode, &url)

	return &url, nil
}

func (r *URLRepository) FindByOriginalURL(ctx context.Context, originalUrl string) (*model.URL, error) {
	var url model.URL

	result := r.db.WithContext(ctx).First(&url, "original_url = ?", originalUrl)

	if result.Error != nil {
		return nil, result.Error
	}

	return &url, nil
}

func (r *URLRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	result := r.db.WithContext(ctx).Model(&model.URL{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1))

	return result.Error
}

func (r *URLRepository) FindAll(ctx context.Context, page, limit int) ([]model.URL, int64, error) {
	offset := (page - 1) * limit

	var urls []model.URL
	var total int64

	result := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(limit).Find(&urls)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	count := r.db.WithContext(ctx).Model(&model.URL{}).Count(&total)

	if count.Error != nil {
		return nil, 0, count.Error
	}

	return urls, total, nil
}

func (r *URLRepository) SaveClickEventsBatch(events []model.ClickEvent) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&events).Error; err != nil {
			return err
		}

		counts := make(map[uint]int)
		for _, clickEvent := range events {
			counts[clickEvent.URLID]++
		}

		for urlID, clickCount := range counts {
			if err := tx.Model(&model.URL{}).
				Where("id = ?", urlID).
				UpdateColumn("click_count", gorm.Expr("click_count + ?", clickCount)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
