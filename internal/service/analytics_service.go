package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/erfangho/url-shortener/internal/grpc"
	"github.com/erfangho/url-shortener/internal/model"
)

type URLClickRepositoryInterface interface {
	SaveClickEventsBatch(events []model.ClickEvent) error
}

type AnalyticsService struct {
	eventChan  chan model.ClickEvent
	urlRepo    URLClickRepositoryInterface
	grpcClient *grpc.AnalyticsClient
	wg         sync.WaitGroup
}

func NewAnalyticsService(urlRepo URLClickRepositoryInterface, grpcClient *grpc.AnalyticsClient, bufferSize int, workerCount int) *AnalyticsService {
	s := &AnalyticsService{
		urlRepo:    urlRepo,
		grpcClient: grpcClient,
		eventChan:  make(chan model.ClickEvent, bufferSize),
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
				s.flush(context.Background(), batch)
				return
			}

			batch = append(batch, event)
		case <-ticker.C:
			s.flush(context.Background(), batch)
			batch = []model.ClickEvent{}
		}
	}
}

func (s *AnalyticsService) Publish(event model.ClickEvent) {
	s.eventChan <- event
}

func (s *AnalyticsService) flush(ctx context.Context, events []model.ClickEvent) {
	if len(events) == 0 {
		return
	}

	err := s.urlRepo.SaveClickEventsBatch(events)

	if err != nil {
		slog.Error(err.Error())
	}

	for _, event := range events {
		s.grpcClient.RecordClick(ctx, &event)
	}
	return
}

func (s *AnalyticsService) Close() {
	close(s.eventChan)
	s.wg.Wait()
}
