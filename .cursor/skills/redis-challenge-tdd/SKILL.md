---
name: redis-challenge-tdd
description: >-
  TDD workflow and PROGRESS updates for this Redis clone (Coding Challenges).
  Use when implementing features, commands, or protocol steps from readme.md.
---

# Redis challenge TDD

## Workflow

1. Write a failing test that describes the next small behavior (RESP case, command, or server response).
2. Implement the minimum code to pass.
3. Refactor; keep tests green.
4. For commands (`PING`, `ECHO`, `GET`, `SET`, …), add a short `redis-cli` example under the relevant step in [readme.md](../../../readme.md) or note it in [PROGRESS.md](../../../PROGRESS.md).

## PROGRESS checklist (after each milestone)

- Which readme step is done (or partially done).
- What works; which tests or manual `redis-cli` checks prove it.
- Known gaps (e.g. only happy path).
- Port **6379** vs local Redis: document if you had to use another port or stop system Redis.

## When skipping tests temporarily

- Add one line to [PROGRESS.md](../../../PROGRESS.md) with reason and a follow-up to add tests.
