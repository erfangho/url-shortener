package service

import (
	"context"
	"testing"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Create(ctx context.Context, url *model.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLRepository) FindByOriginalURL(ctx context.Context, originalUrl string) (*model.URL, error) {
	args := m.Called(ctx, originalUrl)
	return args.Get(0).(*model.URL), args.Error(1)
}

func (m *MockURLRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	args := m.Called(ctx, shortCode)
	return args.Error(0)
}

func (m *MockURLRepository) FindAll(ctx context.Context, page, limit int) ([]model.URL, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]model.URL), args.Get(1).(int64), args.Error(2)
}

func (m *MockURLRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	args := m.Called(ctx, shortCode)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.URL), args.Error(1)
}

func TestURLService_GetURL(t *testing.T) {
	repo := &MockURLRepository{}
	repo.On("FindByShortCode", context.Background(), "abc123").
		Return(&model.URL{ShortCode: "abc123", OriginalURL: "https://test.com"}, nil)

	service := NewURLService(repo)

	url, err := service.GetURL(context.Background(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://test.com", url.OriginalURL)
	assert.Equal(t, "abc123", url.ShortCode)
	repo.AssertExpectations(t)
}

func TestURLService_GetURL_NotFound(t *testing.T) {
	repo := &MockURLRepository{}
	repo.On("FindByShortCode", context.Background(), "abc123").
		Return(nil, gorm.ErrRecordNotFound)

	service := NewURLService(repo)

	_, err := service.GetURL(context.Background(), "abc123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	repo.AssertExpectations(t)
}

func TestURLService_CreateURL_Success(t *testing.T) {
	repo := &MockURLRepository{}
	repo.On("FindByShortCode", context.Background(), mock.Anything).
		Return(nil, gorm.ErrRecordNotFound)
	repo.On("Create", context.Background(), mock.AnythingOfType("*model.URL")).
		Return(nil)

	service := NewURLService(repo)

	url, err := service.CreateURL(context.Background(), "https://example.com")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.OriginalURL)
	assert.Len(t, url.ShortCode, 6)
	repo.AssertExpectations(t)
}

func TestURLService_CreateURL_RepositoryError(t *testing.T) {
	repo := &MockURLRepository{}
	repo.On("FindByShortCode", context.Background(), mock.Anything).
		Return(nil, gorm.ErrRecordNotFound)
	repo.On("Create", context.Background(), mock.AnythingOfType("*model.URL")).
		Return(gorm.ErrRecordNotFound)

	service := NewURLService(repo)

	_, err := service.CreateURL(context.Background(), "https://example.com")

	assert.Error(t, err)
	repo.AssertExpectations(t)
}
