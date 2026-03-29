package function

import (
	"context"
	"testing"

	"github.com/nox456/forgectl/internal/event"
)

func testHandler(ctx context.Context, e event.Event) error {
	return nil
}

func TestRegistryAndLookup(t *testing.T) {
	t.Run("Duplicated function ID", func(t *testing.T) {
		r := NewRegistry()

		err := r.Register(Function{ID: "test-1", Name: "user/created", Trigger: "user.created", Handler: testHandler})
		if err != nil {
			t.Fatalf("Expected success, got error: %s", err)
		}

		err = r.Register(Function{ID: "test-1", Name: "user/updated", Trigger: "user.updated", Handler: testHandler})

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("With function registered", func(t *testing.T) {
		r := NewRegistry()

		r.Register(Function{ID: "test-1", Name: "user/created", Trigger: "user.created", Handler: testHandler})
		r.Register(Function{ID: "test-2", Name: "user/updated", Trigger: "user.updated", Handler: testHandler})
		r.Register(Function{ID: "test-3", Name: "user/updated", Trigger: "user.updated", Handler: testHandler})

		fns := r.Lookup("user.created")
		if len(fns) != 1 {
			t.Errorf("Expected 1 function, got %d", len(fns))
		}

		fns = r.Lookup("user.updated")
		if len(fns) != 2 {
			t.Errorf("Expected 2 functions, got %d", len(fns))
		}
	})

	t.Run("No functions registered", func(t *testing.T) {
		r := NewRegistry()

		fns := r.Lookup("user.created")
		if len(fns) != 0 {
			t.Errorf("Expected 0 functions, got %d", len(fns))
		}
	})
}
