# Build Your Own Redis Server

This repository contains my implementation of a lite version of Redis, built as part of the [Coding Challenges](https://codingchallenges.fyi/challenges/intro) series.

Redis is an in-memory data structure server that supports strings, hashes, lists, sets, and more. It was originally created by Salvatore Sanfilippo in just over 300 lines of TCL (original version [available here](https://gist.github.com/antirez/601351)).

## 🚀 The Challenge

The goal is to build a functional Redis clone that supports the same core operations and protocols as the original version, following a step-by-step evolution.

-----

## 🛠 Progress & Roadmap

### [ ] Step 0: Environment Setup

Set up the development environment for network programming and **Test-Driven Development (TDD)**.

  - Install the official [Redis server](https://redis.io/download) to use `redis-cli` as a test client.
  - **Git hooks (recommended):** run tests and vet before push. From the repo root:

    ```bash
    git config core.hooksPath tools/hooks
    ```

    On Windows, use Git Bash or another environment where `sh` can run the hook. The [pre-push](tools/hooks/pre-push) hook runs `go test ./...` and `go vet ./...`.

  - Ongoing log: see [PROGRESS.md](PROGRESS.md).

### [ ] Step 1: RESP Protocol Implementation

Implement the **Redis Serialisation Protocol (RESP)**. Refer to the [RESP protocol specification](https://redis.io/docs/latest/develop/reference/protocol-spec/).

  - Handle Simple Strings (`+`), Errors (`-`), Integers (`:`), Bulk Strings (`$`), and Arrays (`*`).
  - Test serialization/deserialization for cases like:
      - `*1\r\n$4\r\nping\r\n`
      - `+OK\r\n`
      - `$-1\r\n` (Null value)

### [ ] Step 2: Basic Server & PING/ECHO

Create the server listening on port **6379**.

  - [ ] Implement `PING`: Respond with `PONG`.
  - [ ] Implement `ECHO`: Return the provided string.

<!-- end list -->

```bash
redis-cli ECHO "Hello World"
# Returns "Hello World"
```

### [ ] Step 3: GET & SET (The Core)

Implement the "Remote Dictionary" functionality.

  - [ ] `SET <key> <value>`: Store a value.
  - [ ] `GET <key>`: Retrieve a value.
  - Choose an efficient internal data structure (like a Hash Map).

### [ ] Step 4: Concurrency

Handle multiple concurrent clients simultaneously.

  - Options: One thread per client or asynchronous programming (Event loop/Goroutines).
  - Test with [redis-benchmark](https://redis.io/docs/latest/operate/oss_and_stack/management/optimization/benchmarks/):

<!-- end list -->

```bash
redis-benchmark -t SET,GET -q
```

### [ ] Step 5: Expiry Options

Extend the `SET` command to support expiry options: `EX`, `PX`, `EXAT`, and `PXAT`.

  - Implement low-overhead expiry (constant time operations). Refer to [EXPIRE command](https://redis.io/commands/expire/) docs for logic.

### [ ] Step 6: Advanced Commands & Persistence

Add support for common operations and disk persistence.

  - [ ] `EXISTS`, `DEL`, `INCR`, `DECR`.
  - [ ] `LPUSH`, `RPUSH` (Lists).
  - [ ] `SAVE`: Save state to disk and load it on startup.

### [ ] Step 7: Performance Testing

Benchmark the implementation against the official Redis server.

  - Focus on data structure efficiency and concurrent access management.

-----

## 📚 Resources

  - **Challenge Source:** [Coding Challenges by John Crickett](https://codingchallenges.fyi/challenges/challenge-redis)
  - **Course:** [Become a Better Software Developer by Building Your Own Redis Server](https://codingchallenges.fyi/courses/redis)
  - **Protocol:** [RESP Specification](https://redis.io/docs/latest/develop/reference/protocol-spec/)
  - **Community:** [Discord Server](https://discord.gg/codingchallenges) | [Twitter](https://twitter.com/johncrickett) | [LinkedIn](https://www.linkedin.com/in/johncrickett/)
  - **Shared solutions:** [Coding Challenges Shared Solutions (GitHub)](https://github.com/JohnCrickett/coding-challenges)

## 🤝 Sharing

If you find this implementation useful, feel free to check out other solutions in the [Coding Challenges Shared Solutions GitHub repo](https://github.com/JohnCrickett/coding-challenges).
