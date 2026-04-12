# How it works (current state)

There are three layers, each owning one concern:

**1. Transport — `internal/server` and `internal/client`**

The daemon opens a Unix domain socket with `net.Listen("unix", path)` and enters an accept loop. Each incoming connection is handed to a goroutine that decodes a JSON event with `json.NewDecoder(conn).Decode(...)` and writes back a JSON ack. The client side mirrors this with `net.Dial` and `json.NewEncoder`.

If you're coming from TypeScript: think of the socket like a tiny local HTTP server, except instead of `express.listen()` you call `net.Listen`, and instead of middleware you just have a `handleConn(ctx, conn)` function spawned per connection — the goroutine-per-connection pattern is Go's idiomatic equivalent of `server.on('connection', handler)`.

**2. Worker pool — `internal/engine/pool.go`**

The pool owns a buffered `chan Job` and launches N workers (currently 3) with `go p.worker(...)`. Each worker is a `for job := range p.jobs` loop, which is the channel equivalent of `for await (const job of asyncIterable)` in TS. Closing the channel is how you signal "no more work" — the range loop exits naturally and the workers return. A `sync.WaitGroup` tracks them so the daemon can wait for in-flight work to drain on shutdown (the Go analog of `await Promise.all(workerPromises)`).

The submission pattern is **fan-out, then fan-in**: the handler submits *all* jobs to the channel first, then collects *all* results. Submitting one job and immediately blocking on its result in the same loop iteration would collapse the whole thing back to sequential execution — a bug I hit and fixed earlier.

**3. Lifecycle — `main.go`**

The root `context.Context` is created here using `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)` and threaded down into `NewPool(ctx, ...)` and `Serve(ctx, ...)`. Everything that needs to know about shutdown gets the context from above rather than creating its own.

One subtlety worth knowing: `listener.Accept()` is **not** context-aware — it blocks on a syscall and ignores `ctx.Done()` entirely. So the server starts a small "shutdown bridge" goroutine that watches `ctx.Done()` and calls `s.listener.Close()` when the context is canceled. Closing the listener makes `Accept()` return an error, the accept loop sees the canceled context, and shutdown proceeds cleanly. In TS terms, this is roughly like wiring an `AbortController` to a resource that doesn't natively support `AbortSignal` — you have to bridge them manually.

### Mental model recap

| Go concept used here   | TypeScript analog                                |
| ---------------------- | ------------------------------------------------ |
| `goroutine`            | `async` function scheduled on the event loop     |
| `chan Job`             | An async iterable / queue you can `for await` on |
| `sync.WaitGroup`       | `Promise.all([...])`                             |
| `context.Context` tree | Linked `AbortController` chain                   |
| `signal.NotifyContext` | `process.on('SIGINT', () => controller.abort())` |

