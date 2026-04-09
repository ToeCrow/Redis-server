---
name: resp-protocol-go
description: >-
  Implement and test RESP types in Go for this repo (simple strings, errors,
  integers, bulk strings, arrays). Use when editing internal/resp or wiring
  protocol to the server.
---

# RESP in Go (this project)

## Spec

- Use the official spec: https://redis.io/docs/latest/develop/reference/protocol-spec/

## Types to cover

| Prefix | Type |
|--------|------|
| `+` | Simple string |
| `-` | Error |
| `:` | Integer |
| `$` | Bulk string |
| `*` | Array |

## Edge cases for tests

- Bulk string null: `$-1\r\n`
- Empty bulk: `$0\r\n\r\n`
- Arrays of arrays; nested parsing boundaries
- Sample frames: `*1\r\n$4\r\nping\r\n`, `+OK\r\n`

## Placement

- Implement encoding/decoding under [internal/resp](../../../internal/resp/); keep `internal/server` focused on I/O and command dispatch.
