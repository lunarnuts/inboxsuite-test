package pool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/internal/dto"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
	"github.com/lunarnuts/inboxsuite-test/internal/listener"
	"github.com/lunarnuts/inboxsuite-test/internal/mappers"
	"github.com/lunarnuts/inboxsuite-test/internal/publisher"
	"github.com/lunarnuts/inboxsuite-test/internal/worker"
	"github.com/lunarnuts/inboxsuite-test/pkg/errors"
	"github.com/lunarnuts/inboxsuite-test/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	exchangeKind = "fanout"
)

type Pool struct {
	ctx       context.Context
	workerNum int
	logger    *slog.Logger
	mappers   interfaces.Mappers
	cache     interfaces.CacheLoader
	cfg       *config.RabbitMQ

	channel *amqp091.Channel
	stats   atomic.Int64
	errChan chan error
	jobCh   <-chan amqp091.Delivery
}

func (p *Pool) Run(ctx context.Context) {
	const op = "Pool.Run"
	jobs := make(chan *dto.Job)
	pubs := make(chan *publisher.Publishing)

	l := listener.New(p.jobCh, jobs, p.errChan, p.logger)
	go l.Listen(ctx)

	pub := publisher.New(pubs, p.errChan, p.channel)
	go pub.Run(ctx)

	var wg sync.WaitGroup
	wg.Add(p.workerNum)

	for i := 0; i < p.workerNum; i++ {
		id := i
		w := worker.New(id, p.cache, p.mappers, p.logger, &wg, jobs, p.errChan, pubs, &p.stats, p.cfg)
		go w.Start(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			p.Stop(ctx)
		case err := <-p.errChan:
			p.logger.Error(op, "error", err.Error())
		}
	}

}

func (p *Pool) Stop(ctx context.Context) {
	const op = "Pool.Stop"
	defer close(p.errChan)
	err := p.PublishStatistics(ctx)
	if err != nil {
		p.logger.Error(op, "error", err.Error())
		return
	}
	err = p.channel.Close()
	if err != nil {
		p.logger.Error(op, "error", err.Error())
		return
	}
}

func (p *Pool) PublishStatistics(ctx context.Context) error {
	body, err := json.Marshal(&dto.Statistics{Count: p.stats.Load()})
	if err != nil {
		return errors.Wrap(err, fmt.Sprint("PublishStatistics: error marshalling result"))
	}
	err = p.channel.PublishWithContext(ctx,
		p.cfg.StatisticsExchange,
		p.cfg.StatisticsExchange,
		true,
		true,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprint("PublishStatistics: error marshalling result"))
	}
	p.logger.Info("sent statistics", "statistics", string(body))
	return nil
}

func (p *Pool) Notify() <-chan error {
	return p.errChan
}

func New(ctx context.Context, cfg *config.Config, cache interfaces.CacheLoader, log *slog.Logger, rabbit *rabbitmq.Rabbit) (interfaces.RunnerNotifier, error) {
	ch, err := rabbit.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open channel")
	}

	err = declareAndBindQueue(ch, cfg.RabbitMQ.JobQueue, "")
	if err != nil {
		return nil, err
	}

	msgs, err := consumeQueue(ch, cfg.RabbitMQ.JobQueue)
	if err != nil {
		return nil, err
	}

	err = declareAndBindExchangeQueue(ch, cfg.RabbitMQ.ResultExchange, exchangeKind)
	if err != nil {
		return nil, err
	}

	err = declareAndBindExchangeQueue(ch, cfg.RabbitMQ.StatisticsExchange, exchangeKind)
	if err != nil {
		return nil, err
	}

	return &Pool{
		ctx:       ctx,
		workerNum: cfg.Worker,
		logger:    log,
		cache:     cache,
		cfg:       &cfg.RabbitMQ,
		channel:   ch,
		mappers:   mappers.New(),
		errChan:   make(chan error),
		jobCh:     msgs,
	}, nil
}

func declareAndBindQueue(ch *amqp091.Channel, queueName, exchangeName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare queue")
	}

	if exchangeName != "" {
		err = ch.QueueBind(queueName, "", exchangeName, false, nil)
		if err != nil {
			return errors.Wrap(err, "failed to bind queue")
		}
	}

	return nil
}

func consumeQueue(ch *amqp091.Channel, queueName string) (<-chan amqp091.Delivery, error) {
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to consume queue")
	}
	return msgs, nil
}

func declareAndBindExchangeQueue(ch *amqp091.Channel, exchangeName, exchangeKind string) error {
	_, err := ch.QueueDeclare(
		exchangeName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare exchange queue")
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare exchange")
	}

	err = ch.QueueBind(exchangeName, "", exchangeName, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to bind exchange queue")
	}

	return nil
}
