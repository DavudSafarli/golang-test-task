package controller

import (
	"context"
	"log"
	"net/http"
	"time"
	"twitch_chat_analysis/domain"

	"github.com/gin-gonic/gin"
)

type EventPublisher interface {
	PublishMessage(ctx context.Context, message domain.Message) error
}
type MessageController struct {
	eventPublisher EventPublisher
}

func NewMessageController(eventPublisher EventPublisher) MessageController {
	return MessageController{
		eventPublisher: eventPublisher,
	}
}

func (ctrl MessageController) Create(c *gin.Context) {
	msg := domain.Message{}
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode message"})
		log.Println("failed to decode message: %w", err)
		return
	}

	// should there be validation?

	msg.SentAt = time.Now()
	if err := ctrl.eventPublisher.PublishMessage(c, msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode message"})
		log.Println("failed to publish: %w", err)
		return
	}

	c.Status(http.StatusOK)
}
