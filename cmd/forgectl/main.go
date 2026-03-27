package main

import (
	"errors"
	"fmt"
	"os"
)

func handleArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("Usage: forgectl <command>")
	}
	switch args[0] {
	case "serve":
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

	cmd, err := handleArgs(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(cmd)
}
