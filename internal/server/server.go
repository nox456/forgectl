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
	"github.com/nox456/forgectl/internal/function"
)

const SocketPath = "/tmp/forgectl.sock"

type Response struct {
	Status   string   `json:"status"`
	Messages []string `json:"messages,omitempty"`
}

type Server struct {
	listener net.Listener
	registry *function.Registry
}

func NewServer(registry *function.Registry) *Server {
	return &Server{
		registry: registry,
	}
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
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	var evt event.Event
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&evt); err != nil {
		fmt.Println("failed to decode event:", err)
		return
	}

	fmt.Printf("received event: %s | data: %v\n", evt.Name, evt.Data)

	functions := s.registry.Lookup(evt.Name)

	resp := Response{Status: "accepted"}
	if len(functions) == 0 {
		fmt.Println("no functions found for event")
		resp = Response{Status: "no_functions"}
	}

	for _, fn := range functions {
		fmt.Printf("invoking function: Name: %s | ID: %s\n", fn.Name, fn.ID)
		err := fn.Handler(ctx, evt)
		if err != nil {
			resp.Status = "function_error"
			resp.Messages = append(resp.Messages, err.Error())
		}
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(resp); err != nil {
		fmt.Println("failed to encode response:", err)
		return
	}
}
