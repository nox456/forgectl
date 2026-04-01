package event

import "time"

type Event struct {
	Name           string         `json:"name"`
	Data           map[string]any `json:"data"`
	Timestamp      time.Time      `json:"timestamp"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
}

func NewEvent(name string, data map[string]any, idempotencyKey string) *Event {
	return &Event{
		Name:           name,
		Data:           data,
		Timestamp:      time.Now(),
		IdempotencyKey: idempotencyKey,
	}
}
