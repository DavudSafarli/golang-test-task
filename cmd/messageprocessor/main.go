package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"twitch_chat_analysis/adapters/storage"
	"twitch_chat_analysis/adapters/streaming"
	"twitch_chat_analysis/messageprocessor"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// wait for exit signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("interrupt")
		cancel()
	}()

	// todo(Davud): read from to env var
	mqConnStr := "amqp://guest:guest@rabbitmq:5672/"

	rabbitMQ := streaming.NewRabbitMQ()
	if err := rabbitMQ.InitConnection(mqConnStr); err != nil {
		log.Fatal("failed to init rabbitmq connection")
	}

	if err := rabbitMQ.CreateQueues(); err != nil {
		log.Fatal("failed create rabbitmq queues")
	}

	redisConnStr := "redis:6379"
	redis := storage.NewRedis(redisConnStr)
	co := messageprocessor.NewMessageProcessor(redis)
	rabbitMQ.ConsumeMessage(ctx, co)

	<-ctx.Done()
}
