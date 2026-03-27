package main

import "testing"

func TestCommands(t *testing.T) {
	test_commands := []struct {
		name string
		args []string
	}{
		{
			name: "serve",
			args: []string{"serve"}}, // WORKS
		{
			name: "send",
			args: []string{"send"}}, // WORKS
		{
			name: "list",
			args: []string{"list"}}, // WORKS
		{
			name: "unknown",
			args: []string{"unknown"}}, // FAILS WITH "Unknown command: unknown"
		{
			name: "empty",
			args: []string{}}, // FAILS WITH "Usage: forgectl <command>"
	}
	for _, tt := range test_commands {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := handleArgs(tt.args)

			if cmd == "" && err == nil {
				t.Errorf("Expected error, got nil")
			}

			if cmd != "" && err != nil {
				t.Errorf("Expected nil, got error: %s", err)
			}
		})
	}
}
