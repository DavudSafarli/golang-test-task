package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"twitch_chat_analysis/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const messages_queue_name = "messages"

type RabbitMQ struct {
	channel *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{}
}

func (mq *RabbitMQ) InitConnection(connStr string) error {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return fmt.Errorf("failed to dial rabbitmq server: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel to rabbitmq server: %w", err)
	}

	mq.channel = ch
	return nil
}

func (mq *RabbitMQ) CreateQueues() error {
	_, err := mq.channel.QueueDeclare(
		messages_queue_name, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	return err
}

func (mq *RabbitMQ) PublishMessage(ctx context.Context, message domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = mq.channel.PublishWithContext(ctx,
		"",                  // exchange
		messages_queue_name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}
	return nil
}

type MessageProcessor interface {
	Handle(ctx context.Context, message domain.Message) error
}

func (mq *RabbitMQ) ConsumeMessage(ctx context.Context, consumer MessageProcessor) error {
	delivery, err := mq.channel.Consume(
		messages_queue_name, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("faield to start rabbitmq consumer: %w", err)
	}

	go func() {
		for {
			select {
			case d := <-delivery:
				msg := domain.Message{}
				json.Unmarshal(d.Body, &msg)
				log.Printf("Received a message: %s", d.Body)
				err := consumer.Handle(ctx, msg)
				if err != nil {
					fmt.Println("failed to acknowledge")
					// do not acknowledge if handler failed
					continue
				}
				// TODO(davud): check what happens unackowledged messages, they should be retried until configured attempts
				err = d.Ack(false)
				if err != nil {
					fmt.Errorf("failed to acknowledge a rabbitmq message: %w", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
