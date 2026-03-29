package event

import "time"

type Event struct {
	Name           string         `json:"name"`
	Data           map[string]any `json:"data"`
	Timestamp      time.Time      `json:"timestamp"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
}
