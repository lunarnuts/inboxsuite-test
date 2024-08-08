package publisher

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type Publishing struct {
	Ctx       context.Context
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
	Msg       amqp091.Publishing
}

type Publisher struct {
	publishCh chan *Publishing
	errorCh   chan error
	channel   *amqp091.Channel
}

func New(publishCh chan *Publishing, errorCh chan error, channel *amqp091.Channel) *Publisher {
	return &Publisher{
		publishCh: publishCh,
		errorCh:   errorCh,
		channel:   channel,
	}
}

func (p *Publisher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			err := p.channel.Close()
			if err != nil {
				p.errorCh <- err
			}
		case publishing := <-p.publishCh:
			err := p.channel.Publish(publishing.Exchange, publishing.Key, publishing.Mandatory, publishing.Immediate, publishing.Msg)
			if err != nil {
				p.errorCh <- err
			}
		}
	}
}
