package main

import (
	"log"
	"time"
	"twitch_chat_analysis/adapters/streaming"
	"twitch_chat_analysis/api/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Todo(davud): read from env
	mqConnStr := "amqp://guest:guest@rabbitmq:5672"

	rabbitMQ := streaming.NewRabbitMQ()
	if err := rabbitMQ.InitConnection(mqConnStr); err != nil {
		log.Fatal("failed to init rabbitmq connection: ", err)
	}

	if err := rabbitMQ.CreateQueues(); err != nil {
		log.Fatal("failed create rabbitmq queues")
	}

	messageController := controller.NewMessageController(rabbitMQ)

	r.POST("/message", messageController.Create)

	r.Run(":8080")
	time.Sleep(time.Minute)
}
