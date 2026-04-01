package client

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/nox456/forgectl/internal/event"
	"github.com/nox456/forgectl/internal/server"
)

func createTestServer(t *testing.T) chan event.Event {
	// Clean up stale test socket file from a previous crash.
	os.Remove(server.SocketPath)

	listener, err := net.Listen("unix", server.SocketPath)
	if err != nil {
		t.Fatalf("failed to listen on socket: %v", err)
	}

	received := make(chan event.Event, 1)

	go func() {
		defer listener.Close()
		defer os.Remove(server.SocketPath)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		var evt event.Event
		json.NewDecoder(conn).Decode(&evt)
		received <- evt

		resp := server.Response{
			Status: "accepted",
		}
		json.NewEncoder(conn).Encode(resp)
	}()

	time.Sleep(50 * time.Millisecond)

	return received
}

func TestSend(t *testing.T) {
	received := createTestServer(t)

	evt := event.Event{
		Name: "test",
		Data: map[string]any{
			"foo": "bar",
		},
	}

	err := Send(evt)
	if err != nil {
		t.Fatalf("failed to send event: %v", err)
	}

	got := <-received

	if got.Name != evt.Name {
		t.Errorf("expected name %q, got %q", evt.Name, got.Name)
	}
	if got.Data["foo"] != evt.Data["foo"] {
		t.Errorf("expected data %q, got %q", evt.Data["foo"], got.Data["foo"])
	}
}

func TestSendNoDaemon(t *testing.T) {
	// No server running — Send should return an error
	os.Remove(server.SocketPath)

	evt := event.Event{Name: "test.event"}
	err := Send(evt)
	if err == nil {
		t.Error("expected error when no daemon is running, got nil")
	}
}
