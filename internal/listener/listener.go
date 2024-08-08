package listener

import (
	"context"
	"encoding/json"
	"github.com/lunarnuts/inboxsuite-test/internal/dto"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type Listener struct {
	delivery <-chan amqp091.Delivery
	jobCh    chan *dto.Job
	errChan  chan error
	logger   *slog.Logger
}

func New(delivery <-chan amqp091.Delivery, jobCh chan *dto.Job, errChan chan error, logger *slog.Logger) *Listener {
	return &Listener{delivery: delivery, jobCh: jobCh, errChan: errChan, logger: logger}
}

func (p *Listener) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(p.jobCh)

		case d := <-p.delivery:
			var err error
			var job dto.Job
			err = json.Unmarshal(d.Body, &job)
			if err != nil {
				p.errChan <- err
				continue
			}
			p.jobCh <- &job
			p.logger.Info("sent job", "body", job)
		}
	}
}
