package service

import (
	"sync"

	"github.com/erfangho/url-shortener/internal/model"
)

type URLClickRepositoryInterface interface {
	SaveClickEventsBatch(events []model.ClickEvent) error
}

type AnalyticsService struct {
	eventChan chan model.ClickEvent
	urlRepo   URLClickRepositoryInterface
	wg        sync.WaitGroup
}

func NewAnalyticsService(urlRepo URLClickRepositoryInterface, bufferSize int, workerCount int) *AnalyticsService {
	s := &AnalyticsService{
		urlRepo:   urlRepo,
		eventChan: make(chan model.ClickEvent, bufferSize),
	}

	s.wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go s.worker()
	}

	return s
}

func (s *AnalyticsService) worker() {

}
