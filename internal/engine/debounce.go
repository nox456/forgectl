package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/nox456/forgectl/internal/event"
	"github.com/nox456/forgectl/internal/function"
)

type DebounceEntry struct {
	function function.Function
	event    event.Event
	timer    *time.Timer
}

type Debouncer struct {
	mu      sync.Mutex
	entries map[string]*DebounceEntry
	pool    *Pool
}

func NewDebouncer(pool *Pool) *Debouncer {
	return &Debouncer{
		entries: make(map[string]*DebounceEntry),
		pool:    pool,
	}
}

func (d *Debouncer) Debounce(fn function.Function, evt event.Event) {
	debounceKey, exists := evt.Data[fn.DebounceConfig.Key]

	if !exists {
		fmt.Printf("debounce key not found in event data: %s\n", fn.DebounceConfig.Key)
		d.pool.Run(Job{
			Function: fn,
			Event:    evt,
		})
		return
	}

	composedKey := fmt.Sprintf("%s-%s", fn.ID, debounceKey)

	d.mu.Lock()
	defer d.mu.Unlock()

	entry, exists := d.entries[composedKey]

	if exists {
		entry.timer.Stop()
		entry.event = evt
	} else {
		entry = &DebounceEntry{
			function: fn,
			event:    evt,
		}
		d.entries[composedKey] = entry
	}

	var newTimer *time.Timer

	newTimer = time.AfterFunc(fn.DebounceConfig.Period, func() {
		d.mu.Lock()
		defer d.mu.Unlock()

		current, exists := d.entries[composedKey]
		if !exists {
			return
		}

		if current.timer != newTimer {
			return
		}

		delete(d.entries, composedKey)
		d.pool.Run(Job{
			Function: current.function,
			Event:    current.event,
		})
	})

	entry.timer = newTimer
}
