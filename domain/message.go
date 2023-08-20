package domain

import "time"

type Message struct {
	Message  string    `json:"message"`
	Sender   string    `json:"sender" form:"sender"`
	Receiver string    `json:"receiver" form:"receiver"`
	SentAt   time.Time `json:"sent_at"`
}
