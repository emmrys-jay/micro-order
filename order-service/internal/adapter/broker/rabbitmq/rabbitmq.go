package rabbitmq

import (
	"context"
	"fmt"
	"order-service/internal/adapter/config"
	"order-service/internal/core/domain"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Message struct {
	amqp.Delivery
}

type Broker struct {
	conn *amqp.Connection

	channel *sync.Pool
}

func New(ctx context.Context, conf *config.RabbitMqConfiguration) (*Broker, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s/", conf.User, conf.Password, conf.Host)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Broker{
		conn: conn,
		channel: &sync.Pool{
			New: func() any {
				ch, err := conn.Channel()
				if err != nil {
					zap.L().Error("Failed to create channel", zap.Error(err))
					return nil
				}
				return ch
			},
		},
	}, nil
}

func (b *Broker) declareQueue(queue string) (*amqp.Channel, amqp.Queue, error) {
	ch, ok := b.channel.Get().(*amqp.Channel)
	if !ok {
		return nil, amqp.Queue{}, fmt.Errorf("No channel found in the channel pool")
	}

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return ch, q, err
}

func (b *Broker) Publish(ctx context.Context, queue string, msg []byte, headers map[string]any) error {
	ch, q, err := b.declareQueue(queue)
	if err != nil {
		return fmt.Errorf("error declaring queue, %w", err)
	}
	defer b.channel.Put(ch)

	ch = b.channel.Get().(*amqp.Channel)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
			Headers:     headers,
		},
	)

	return err
}

func (b *Broker) Consume(ctx context.Context, queue string, handler func(*zap.Logger, []byte) error) {
	log := zap.L()
	ch, q, err := b.declareQueue(queue)
	if err != nil {
		log.Error("Could not declare queue", zap.Error(err))
		return
	}
	defer b.channel.Put(ch)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			correlationId := msg.Headers[string(domain.CorrelationIDCtxKey)] // prevent panic
			logger := log.With(zap.Any("correlation_id", correlationId))

			if err := handler(logger, msg.Body); err != nil {
				log.Error("Error executing message from the queue", zap.Error(err))
				msg.Nack(false, true) // requeue message
			} else {
				msg.Ack(false) // acknowledge
			}
		}
	}

}
