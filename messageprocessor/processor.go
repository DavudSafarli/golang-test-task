package messageprocessor

import (
	"context"
	"fmt"
	"log"
	"twitch_chat_analysis/domain"
)

type Storage interface {
	Store(ctx context.Context, message domain.Message) error
}

type MessageProcessor struct {
	storage Storage
}

func NewMessageProcessor(storage Storage) MessageProcessor {
	return MessageProcessor{
		storage: storage,
	}
}
func (p MessageProcessor) Handle(ctx context.Context, message domain.Message) error {
	err := p.storage.Store(ctx, message)
	if err != nil {
		log.Println("failed to store the message in storage: %w", err)
		return fmt.Errorf("failed to store the message in storage: %w", err)
	}
	return nil
}
