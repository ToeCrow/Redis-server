# Progress

Running log for the Redis clone challenge: what works, what proves it, and open gaps.

## Environment

- **Go:** 1.26.1 (see [go.mod](go.mod); update if you change the toolchain)
- **redis-cli:** installed? (yes/no — install from [Redis download](https://redis.io/download) or OS package)
- **Port 6379:** note if a system Redis already uses it (use another port or stop the other server).

## Milestones

### Step 0 — Environment

- [ ] Dev environment ready for network programming and TDD.
- [ ] Git hooks enabled: `git config core.hooksPath tools/hooks` (see [readme.md](readme.md) Step 0).

### Step 1 — RESP

- [x] Encode/decode implemented under `internal/resp` with tests (`ReadValue`, `WriteValue`; RESP2: `+`, `-`, `:`, `$`, `*`).
- **Proves it:** `go test ./internal/resp/...` — readme examples (`*1\r\n$4\r\nping\r\n`, `+OK\r\n`, `$-1\r\n`), empty bulk `$0\r\n\r\n`, round-trip for all kinds, nested array.
- **Limits:** RESP2 only (no RESP3 types). Simple/error lines assume well-formed CRLF-terminated payloads (as in spec).

### Step 2 — Server, PING, ECHO

- [x] TCP server: [`internal/server`](internal/server) — `Run` listens on the given address (default [`main.go`](main.go) uses `:6379`); `Serve` accepts one goroutine per client; requests are RESP2 arrays parsed with [`resp.ReadValueFrom`](internal/resp/decode.go) on a **single** `bufio.Reader` per connection so pipelined commands do not lose buffered bytes.
- [x] `PING` (no args) → `+PONG\r\n`; `PING` with one bulk argument echoes that string as a bulk reply (Redis-compatible). `ECHO` → bulk string reply; unknown commands → `-ERR unknown command ...\r\n`.
- **Proves it:** `go test ./internal/server/...` — raw RESP over loopback `:0` for `PING`, `ECHO "Hello World"`, and unknown command. `go test ./internal/resp/...` includes `ReadValueFrom` sequential decode.
- **Port 6379:** if a system Redis already binds `:6379`, run with another address (e.g. change `main` or `go run .` with a fork that passes a different `addr`) or stop the other process; tests use `:0` to avoid collisions.

### Step 3 — GET, SET

- [x] In-memory string store: [`KVStore`](internal/server/store.go) — `map[string]string` with `sync.RWMutex` for safe concurrent access from multiple client goroutines.
- [x] `SET key value` (three bulk arguments) → `+OK\r\n`. `GET key` (two bulk arguments) → bulk string value, or `$-1\r\n` if the key is missing (Redis-compatible null bulk).
- [x] Wrong arity → `-ERR wrong number of arguments for 'set'|'get' command\r\n`; non-bulk key/value where required → same bulk-string error pattern as `ECHO`.
- **Proves it:** `go test ./internal/server/...` — `TestSetReturnsOK`, `TestSetThenGetReturnsValue`, `TestGetMissingKeyReturnsNullBulk`, `TestSetWrongArityReturnsError`, `TestGetWrongArityReturnsError`, plus existing PING/ECHO tests. **Manual:** `redis-cli -p 6379 SET foo bar` then `GET foo` (stop system Redis first or use another port if `:6379` is taken).

### Later steps

- Concurrency benchmarks (`redis-benchmark`), expiry, persistence, benchmarks vs official Redis.

## Design notes (ADR-light)

- **KV store (2026-04):** `map[string]string` + `sync.RWMutex` — O(1) average get/set; multiple readers can proceed in parallel. Alternatives later: `sync.Map` or sharding if profiling warrants it.
