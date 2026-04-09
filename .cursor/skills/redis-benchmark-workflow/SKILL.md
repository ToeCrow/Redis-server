---
name: redis-benchmark-workflow
description: >-
  Run redis-benchmark against this clone, compare roughly to official Redis,
  and record results in PROGRESS.md. Use near readme steps 4 and 7.
---

# redis-benchmark workflow

## Prerequisites

- Server listening on the same host/port the benchmark targets (default **6379**).
- If official Redis is already on 6379, stop it or run this server on another port and pass `-p` to `redis-benchmark`.

## Example commands

```bash
redis-benchmark -t SET,GET -q
```

Add `-p <port>` if not using 6379.

## What to record in PROGRESS.md

- Date, command line, and key options (e.g. `-t`, `-q`, `-p`).
- Rough ops/sec or p50/p99 if you capture them; note hardware/load caveats.
- Optional: one-line comparison vs stock Redis on the same machine (informal).

## Interpretation

- Large gaps vs Redis are expected for a learning implementation; focus on regressions when you change concurrency or storage.
