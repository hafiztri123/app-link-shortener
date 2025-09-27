package messaging

import (
	"context"
	"encoding/json"
	"hafiztri123/app-link-shortener/internal/models"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	queue string
}

type ClickEvent = models.Click

func NewRabbitMQ(url, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		slog.Error("Failed to open a channel", "error", err)
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false, 
		false,
		nil,
	)

	if err != nil {
		ch.Close()
		conn.Close()
		slog.Error("Failed to declare a queue", "error", err)
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		channel: ch,
		queue: queueName,
	}, nil
}


func (r *RabbitMQ) PublishClickEvent(ctx context.Context, event ClickEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		slog.Error("Failed to marshal click event", "error", err)
		return err
	}

	err = r.channel.PublishWithContext(
		ctx,
		"", //default exchange
		r.queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "application/json",
			Body: body,
			Timestamp: time.Now(),
		},
	)

	if err != nil {
		slog.Error("Failed to publish click event", "error", err)
		return err
	}

	return nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}

func (r *RabbitMQ) HealthCheck() error {
	return r.channel.ExchangeDeclarePassive(
		"amq.direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
}