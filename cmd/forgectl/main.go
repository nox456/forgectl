package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/nox456/forgectl/internal/server"
)

func handleArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("Usage: forgectl <command>")
	}
	switch args[0] {
	case "serve":
		s := server.NewServer()

		if err := s.Serve(); err != nil {
			return "", err
		}

		return "serve", nil
	case "send":
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
