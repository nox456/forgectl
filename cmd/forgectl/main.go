package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nox456/forgectl/internal/client"
	"github.com/nox456/forgectl/internal/engine"
	"github.com/nox456/forgectl/internal/event"
	"github.com/nox456/forgectl/internal/function"
	"github.com/nox456/forgectl/internal/server"
)

func handleArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("Usage: forgectl <command>")
	}
	switch args[0] {
	case "serve":
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		pool := engine.NewPool(3, 3, ctx)
		debouncer := engine.NewDebouncer(pool)
		registry := function.NewRegistry()

		s := server.NewServer(registry, pool, ctx, debouncer)

		if err := s.Serve(); err != nil {
			return "", err
		}

		return "serve", nil
	case "send":
		if len(args) < 2 {
			return "", errors.New("Usage: forgectl send <event name> [data]")
		}

		data := make(map[string]any)
		if len(args) > 2 {
			if err := json.Unmarshal([]byte(args[2]), &data); err != nil {
				return "", err
			}
		}

		var idempotencyKey string
		if args[3] == "" {
			idempotencyKey = args[1]
		} else {
			idempotencyKey = args[3]
		}

		evt := event.NewEvent(args[1], data, idempotencyKey)

		if err := client.Send(*evt); err != nil {
			return "", err
		}

		return "send", nil
	case "list":
		return "list", nil
	default:
		return "", errors.New("Unknown command: " + args[0])
	}
}

func main() {
	args := os.Args[1:]

	_, err := handleArgs(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
