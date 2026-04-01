package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nox456/forgectl/internal/event"
)

const SocketPath = "/tmp/forgectl.sock"

type Response struct {
	Status string `json:"status"`
}

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve() error {
	// Clean up stale socket file from a previous crash.
	// If the daemon crashed without cleaning up, the file
	// still exists and net.Listen will fail with
	// "address already in use."
	os.Remove(SocketPath)

	var err error
	s.listener, err = net.Listen("unix", SocketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	defer os.Remove(SocketPath)
	defer s.listener.Close()

	// Create a context that cancels on SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("forgectl daemon listening on", SocketPath)

	go func() {
		<-ctx.Done()
		s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("shutting down...")
				return nil
			default:
				return fmt.Errorf("accept error: %w", err)
			}
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	var evt event.Event
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&evt); err != nil {
		fmt.Println("failed to decode event:", err)
		return
	}

	fmt.Printf("received event: %s | data: %v\n", evt.Name, evt.Data)

	resp := Response{Status: "accepted"}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(resp); err != nil {
		fmt.Println("failed to encode response:", err)
		return
	}
}
