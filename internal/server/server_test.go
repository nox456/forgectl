package server

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/nox456/forgectl/internal/event"
)

func TestHandleConn(t *testing.T) {
	clientConn, serverConn := net.Pipe()

	srv := NewServer()

	go srv.handleConn(serverConn)

	evt := event.Event{
		Name: "test",
		Data: map[string]any{
			"test": "test",
		},
		Timestamp:      time.Now(),
		IdempotencyKey: "test",
	}

	encoder := json.NewEncoder(clientConn)
	if err := encoder.Encode(evt); err != nil {
		t.Fatalf("failed to send event: %v", err)
	}

	var resp Response
	decoder := json.NewDecoder(clientConn)
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "accepted" {
		t.Errorf("expected response status to be 'accepted', got %s", resp.Status)
	}
}

func TestInvalidEventData(t *testing.T) {
	clientConn, serverConn := net.Pipe()

	srv := NewServer()

	go srv.handleConn(serverConn)

	clientConn.Write([]byte("invalid json"))
	clientConn.Close()
}
