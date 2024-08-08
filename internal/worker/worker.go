package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/internal/dto"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
	"github.com/lunarnuts/inboxsuite-test/internal/publisher"
	"github.com/lunarnuts/inboxsuite-test/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"sync"
	"sync/atomic"
)

func New(
	id int,
	cache interfaces.CacheLoader,
	mappers interfaces.Mappers,
	logger *slog.Logger,
	wg *sync.WaitGroup,
	jobCh chan *dto.Job,
	errorCh chan error,
	publishCh chan *publisher.Publishing,
	stats *atomic.Int64,
	cfg *config.RabbitMQ,
) *Worker {
	return &Worker{
		id:         id,
		cache:      cache,
		mappers:    mappers,
		logger:     logger,
		cfg:        cfg,
		statistics: stats,
		wg:         wg,
		publishCh:  publishCh,
		jobCh:      jobCh,
		errorCh:    errorCh,
	}
}

type Worker struct {
	cache      interfaces.CacheLoader
	mappers    interfaces.Mappers
	logger     *slog.Logger
	cfg        *config.RabbitMQ
	statistics *atomic.Int64

	id        int
	wg        *sync.WaitGroup
	publishCh chan *publisher.Publishing
	jobCh     chan *dto.Job
	errorCh   chan error
}

func (w *Worker) Start(ctx context.Context) {
	defer w.wg.Done()
	w.logger.Info("worker starting", "id", w.id)
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker shutting down", "workerId", w.id)
			return
		case job := <-w.jobCh:
			w.logger.Info("worker received job", "job", job)
			roadmapId, err := w.cache.Get(w.mappers.ClassID().ToEntity(job.ClassID))
			if err != nil {
				w.errorCh <- errors.Wrap(err, fmt.Sprintf("worker #%d: failed to load cache", w.id))
				continue
			}

			result := dto.Result{
				ProfileID: job.ProfileID,
				RoadmapID: w.mappers.RoadmapID().ToDTO(roadmapId),
			}
			w.PublishResult(ctx, result)

			w.statistics.Add(1)
			stat := w.statistics.Load()
			if stat%10 == 0 {
				w.PublishStatistics(ctx, stat)
			}
		}
	}
}

func (w *Worker) PublishResult(ctx context.Context, result dto.Result) {
	body, err := json.Marshal(result)
	if err != nil {
		w.errorCh <- errors.Wrap(err, fmt.Sprintf("worker #%d: error marshalling result", w.id))
	}
	w.publishCh <- &publisher.Publishing{
		Ctx:       ctx,
		Exchange:  w.cfg.ResultExchange,
		Key:       w.cfg.ResultExchange,
		Mandatory: false,
		Immediate: false,
		Msg: amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	}
	w.logger.Info("sent result", "result", string(body))
}

func (w *Worker) PublishStatistics(ctx context.Context, count int64) {
	body, err := json.Marshal(&dto.Statistics{Count: count})
	if err != nil {
		w.errorCh <- errors.Wrap(err, fmt.Sprintf("worker #%d: error marshalling result", w.id))
	}
	w.publishCh <- &publisher.Publishing{
		Ctx:       ctx,
		Exchange:  w.cfg.StatisticsExchange,
		Key:       w.cfg.StatisticsExchange,
		Mandatory: false,
		Immediate: false,
		Msg: amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	}
	w.logger.Info("sent statistics", "statistics", string(body))
}
