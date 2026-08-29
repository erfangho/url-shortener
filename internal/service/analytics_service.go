package service

import (
	"log/slog"
	"sync"
	"time"

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
	defer s.wg.Done()

	batch := []model.ClickEvent{}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-s.eventChan:
			if !ok {
				s.flush(batch)
				return
			}

			batch = append(batch, event)
		case <-ticker.C:
			s.flush(batch)
			batch = []model.ClickEvent{}
		}
	}
}

func (s *AnalyticsService) Publish(event model.ClickEvent) {
	s.eventChan <- event
}

func (s *AnalyticsService) flush(events []model.ClickEvent) {
	if len(events) == 0 {
		return
	}

	err := s.urlRepo.SaveClickEventsBatch(events)

	if err != nil {
		slog.Error(err.Error())
	}
	return
}

func (s *AnalyticsService) Close() {
	close(s.eventChan)
	s.wg.Wait()
}
