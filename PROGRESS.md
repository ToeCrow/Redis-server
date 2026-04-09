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

- [ ] Listener and commands documented; manual `redis-cli` checks noted here if useful.

### Later steps

- Add dated bullets as you complete GET/SET, concurrency, expiry, persistence, benchmarks.

## Design notes (ADR-light)

- (e.g. **Concurrency:** goroutine per client vs event loop — decision + date.)
