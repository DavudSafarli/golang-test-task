package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"twitch_chat_analysis/domain"

	"github.com/redis/go-redis/v9"
)

const messages_key = "messages"

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Redis{
		client: rdb,
	}
}

func (r Redis) getMessageKey(sender, receiver string) string {
	// message-key_sender-name_receiver-name
	return fmt.Sprintf("%s_%s_%s", messages_key, sender, receiver)
}

func (r Redis) Store(ctx context.Context, message domain.Message) error {
	val, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// score := message.SentAt.UnixMilli()
	// err = r.client.ZAdd(ctx, messages_key, redis.Z{
	// 	Score:  float64(score),
	// 	Member: val,
	// }).Err()

	key := r.getMessageKey(message.Sender, message.Receiver)
	err = r.client.LPush(ctx, key, val).Err()
	if err != nil {
		return fmt.Errorf("Failed to store message to redis: %w", err)
	}
	return nil
}

func (r Redis) GetMessagesSortedDesc(ctx context.Context, sender, receiver string) ([]domain.Message, error) {
	key := r.getMessageKey(sender, receiver)
	cmd := r.client.LRange(ctx, key, 0, -1)

	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	values := cmd.Val()
	messages := make([]domain.Message, 0, len(values))

	for _, val := range values {
		m := domain.Message{}
		json.Unmarshal([]byte(val), &m)
		messages = append(messages, m)
	}
	return messages, nil
}
