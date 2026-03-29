package function

import (
	"context"
	"fmt"

	"github.com/nox456/forgectl/internal/event"
)

type Handler func(ctx context.Context, e event.Event) error

type Function struct {
	ID      string
	Name    string
	Trigger string
	Handler Handler
}

type Registry struct {
	functions map[string][]Function
	byID      map[string]struct{}
}

func NewRegistry() *Registry {
	return &Registry{
		functions: make(map[string][]Function),
		byID:      make(map[string]struct{}),
	}
}

func (r *Registry) Register(fn Function) error {
	if _, exists := r.byID[fn.ID]; exists {
		return fmt.Errorf("function already registered: %s", fn.ID)
	}

	r.byID[fn.ID] = struct{}{}
	r.functions[fn.Trigger] = append(r.functions[fn.Trigger], fn)

	return nil
}

func (r *Registry) Lookup(trigger string) []Function {
	return r.functions[trigger]
}
