package controller

import (
	"context"
	"log"
	"net/http"
	"twitch_chat_analysis/domain"

	"github.com/gin-gonic/gin"
)

type Storage interface {
	GetMessagesSortedDesc(ctx context.Context, sender, receiver string) ([]domain.Message, error)
}
type ReportingController struct {
	storage Storage
}

func NewReportingController(storage Storage) ReportingController {
	return ReportingController{
		storage: storage,
	}
}

func (ctrl ReportingController) Get(c *gin.Context) {

	req := domain.Message{}
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind query"})
		log.Println("failed to bind query: %w", err)
		return
	}

	msgs, err := ctrl.storage.GetMessagesSortedDesc(c, req.Sender, req.Receiver)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get messages"})
		log.Println("failed to get messages: %w", err)
		return
	}

	c.IndentedJSON(http.StatusOK, msgs)
}
