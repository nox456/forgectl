package client

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/nox456/forgectl/internal/event"
	"github.com/nox456/forgectl/internal/server"
)

func Send(evt event.Event) error {
	conn, err := net.Dial("unix", server.SocketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to socket: %w", err)
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(evt); err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	decoder := json.NewDecoder(conn)
	var resp server.Response
	if err := decoder.Decode(&resp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("received response: %s\n\n  messages: %s\n", resp.Status, resp.Messages)

	return nil
}
