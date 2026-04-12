# ForgeCTL

A local, event-driven function execution daemon — written in Go as a learning project for concurrency and channels. Inspired by tools like [Inngest](https://www.inngest.com/) and [Hatchet](https://hatchet.run/), but stripped down to the essentials and meant to run on a single machine.

---

## What it does

`forgectl` runs as a long-lived **daemon** that listens on a Unix domain socket. You send it **events** from a separate **client** process, and the daemon dispatches each event to a pool of worker goroutines that execute the associated work concurrently.

The same binary plays both roles depending on the subcommand:

```
forgectl serve                        # start the daemon (blocks)
forgectl send user.created '{"email":"x@y.com"}'   # send an event
```

The eventual goal is a tool where you register Go functions against event names, and the daemon handles concurrent execution, retries, idempotent replays, debouncing, and persistence — but right now the focus has been on getting the **plumbing** right.

If you're interested on how it works, check out [HOW_IT_WORKS.md](./docs/HOW_IT_WORKS.md).

---

## Project layout

```
forgectl/
├── main.go                       # composition root: context, pool, server wiring
├── main_test.go                  # smoke tests for command dispatch
├── internal/
│   ├── server/server.go          # Unix socket listener, accept loop, handleConn
│   ├── client/client.go          # net.Dial + JSON encode/decode
│   ├── engine/pool.go            # buffered channel + N worker goroutines
│   └── function/                 # function registry (in progress, Phase 3+)
├── Makefile
├── go.mod
└── README.md
```

---

## Build

Requires Go 1.21+.

```bash
make build           # produces ./forgectl
# or directly:
go build -o forgectl .
```

---

## Dev

### Manual end-to-end test

Because `forgectl serve` blocks the terminal, you need **two terminals** (or a `tmux` split).

**Terminal 1 — start the daemon:**

```bash
./forgectl serve
```

You should see something like `daemon listening on /tmp/forgectl.sock` (adjust to whatever your server actually logs). Leave it running.

**Terminal 2 — send some events:**

```bash
./forgectl send user.created '{"email":"alice@example.com"}'
./forgectl send user.created '{"email":"bob@example.com"}'
```

Each call should print the daemon's JSON ack and exit immediately. In Terminal 1, you should see the daemon log each received event.

### Testing concurrency by eye

The whole reason the worker pool exists is to run jobs in parallel. To actually *see* that happening, send a burst of events from a shell loop:

```bash
for i in $(seq 1 10); do
  ./forgectl send user.created "{}" abc$i &
done
wait
```

With 3 workers, you should see jobs being picked up in batches of (up to) 3. If your handler logs timestamps, the parallelism will be visible in the daemon output. If everything looks strictly sequential, that's the fan-out/fan-in bug — go check that the result-collection loop runs *after* the submission loop, not interleaved with it.

## Why this project exists

This is a learning exercise, not a production tool (for now). The goal is to internalize Go's concurrency primitives — goroutines, channels, `sync.WaitGroup`, `context.Context`, `select`, mutexes — by building something where getting them wrong is immediately, viscerally obvious. Every phase is designed so that the bug you'll hit teaches the lesson better than any blog post could.
