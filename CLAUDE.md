# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Standards reference**: <https://github.com/vergil-project/vergil-tooling>
— active standards documentation lives in the vergil-tooling repository under `docs/`.
Repository profile: `vergil.toml`.

## Memory management

Memory is allowed with human approval. The authoritative policy is in
the user's global `~/.claude/CLAUDE.md` — agents must propose memory
writes and suggest a destination (repo memory, global CLAUDE.md, or
plugin/skill issue) before writing. See that file for the full
workflow.

Available skills:
- `/vergil:memory-init` — set up or update the policy header
  in a project's `MEMORY.md`.
- `/vergil:memory-audit` — structured collaborative review
  of memory files.

## Parallel AI agent development

This repository supports running multiple Claude Code agents in parallel via
git worktrees. The convention keeps parallel agents' working trees isolated
while preserving shared project memory (which Claude Code derives from the
session's starting CWD).

**Canonical spec:**
[`vergil-tooling/docs/specs/worktree-convention.md`](https://github.com/vergil-project/vergil-tooling/blob/develop/docs/specs/worktree-convention.md)
— full rationale, trust model, failure modes, and memory-path implications.
The canonical text lives in `vergil-tooling`; this section is the local
on-ramp.

### Structure

```text
<project-root>/                              ← sessions ALWAYS start here
  .git/
  CLAUDE.md, …                               ← main worktree (usually `develop`)
  .worktrees/                                ← container for parallel worktrees
    issue-<N>-<short-slug>/                  ← worktree on feature/<N>-<short-slug>
    …
```

### Rules

1. **Sessions always start at the project root.**
   Never start Claude from inside `.worktrees/<name>/`. This keeps the
   memory-path slug stable and shared.
2. **Each parallel agent is assigned exactly one worktree.** The session
   prompt names the worktree (see Agent prompt contract below).
   - For Read / Edit / Write tools: use the worktree's absolute path.
   - For Bash commands that touch files: `cd` into the worktree first,
     or use absolute paths.
3. **The main worktree is read-only.** All edits flow through a worktree
   on a feature branch — the logical endpoint of the standing
   "no direct commits to develop" policy.
4. **One worktree per issue.** Don't stack in-flight issues. When a
   branch lands, remove the worktree before starting the next.
5. **Naming: `issue-<N>-<short-slug>`.** `<N>` is the GitHub issue
   number; `<short-slug>` is 2–4 kebab-case tokens.

### Agent prompt contract

When launching a parallel-agent session, use this template (fill in the
placeholders):

```text
You are working on issue #<N>: <issue title>.

Your worktree is: <project-root>/.worktrees/issue-<N>-<slug>/
Your branch is:   feature/<N>-<slug>

Rules for this session:
- Do all git operations from inside your worktree:
    cd <absolute-worktree-path> && vrg-git <command>
- For Read / Edit / Write tools, use the absolute worktree path.
- For Bash commands that touch files, cd into the worktree first
  or use absolute paths.
- Do not edit files at the project root. The main worktree is
  read-only — all changes flow through your worktree on your
  feature branch.
- When you need to run validation, run it from inside your worktree
  (vrg-container-run mounts the current directory).
```

All fields are required.

## Shell command policy

Use `vrg-git` instead of `git` for all git operations. Use `vrg-gh`
instead of `gh` for all GitHub CLI operations. These wrappers enforce
subcommand allowlists, flag deny lists, and credential selection.

Raw `git` and `gh` are denied by the permission model. If a command
is not available through the wrappers, explain the situation to the
human who can run it directly via `! <command>` in the prompt.

## Validation

```bash
vrg-container-run -- vrg-validate
```

This is the **only** validation command. Do not run individual linters,
formatters, or other tools outside of `vrg-validate`. If a tool is not
invoked by `vrg-validate`, it is not part of the validation pipeline.

## Project Overview

This is a Go port of `pymqrest`, providing a Go wrapper for the IBM MQ administrative REST API. The project provides typed Go functions for every MQSC command exposed by the `runCommandJSON` REST endpoint, with automatic attribute name translation between Go `snake_case` and native MQSC parameter names.

**Project name**: mq-rest-admin-go

**Status**: Pre-alpha (initial setup)

**Module path**: `github.com/mq-rest-admin-project/mq-rest-admin-go`

**Package name**: `mqrestadmin`

**Canonical Standards**: This repository follows standards at <https://github.com/vergil-project/vergil-tooling> (local path: `../vergil-tooling` if available)

## Development Commands

### Environment Setup

- **Go**: 1.25+ (CI tests 1.25 and 1.26; go.mod declares 1.26)
- **golangci-lint**: `brew install golangci-lint` (not in `tools.go` per project recommendation)
- **Dev tools** (pinned in `tools.go`): `go install golang.org/x/vuln/cmd/govulncheck && go install github.com/vladopajic/go-test-coverage/v2 && go install github.com/fzipp/gocyclo/cmd/gocyclo`
- **Hook guard**: The Claude Code PreToolUse hook guard (`.claude/hooks/guard.sh`) blocks raw `git`/`gh` — use `vrg-git`/`vrg-gh`
- **VERGIL CLI tools** (`vrg-commit`, `vrg-validate`, etc.) are pre-installed in the dev container images

### Build

```bash
go build ./...          # Compile all packages
go vet ./...            # Static analysis
```

### Validation

```bash
vrg-container-run -- vrg-validate   # Canonical validation (runs in dev container)
```

### CI

PR CI (`.github/workflows/ci.yml`) uses vergil-actions v2.0 reusable
workflows for quality (lint, typecheck), unit tests (Go 1.25/1.26
matrix), dependency audit, security (CodeQL, Trivy, Semgrep, standards),
and release gates. A bespoke job handles integration tests (MQ containers).

### Testing

```bash
go test ./...                                   # Unit tests
go test -race -count=1 ./...                    # Unit tests with race detection
go test -run TestFunctionName ./mqrestadmin/...  # Run a single test
go test -race -count=1 -tags=integration ./...  # Integration tests (needs MQ env)
go test -coverprofile=coverage.out ./... && go-test-coverage --config .testcoverage.yml  # Coverage gate
go tool cover -html=coverage.out                # View coverage in browser
```

- **Framework**: stdlib `testing` package
- **Coverage**: Target 100% line coverage (file, package, and total), enforced in
  CI via `go-test-coverage` with `.testcoverage.yml`. Structurally untestable
  lines (e.g., `json.Marshal` on `map[string]any`, embedded JSON parse errors) are
  annotated with `// coverage-ignore -- <reason>` on the **preceding line** (the
  `{` line) and excluded from measurement.
- **Integration tests**: Require `MQ_REST_ADMIN_RUN_INTEGRATION=1` and a running
  MQ container. CI uses the `mq-rest-admin-project/mq-rest-admin-dev-environment` action.
  Wrapper scripts in `scripts/dev/mq_*.sh` manage the MQ lifecycle with
  `COMPOSE_PROJECT_NAME=mqrest-go` and Go-specific port allocation
  (REST: 9463/9464, MQ: 1434/1435).

### Local MQ Container

The MQ development environment is owned by the
[mq-rest-admin-dev-environment](https://github.com/mq-rest-admin-project/mq-rest-admin-dev-environment)
repository. Clone it as a sibling directory before running lifecycle
scripts:

```bash
# Prerequisite (one-time)
git clone https://github.com/mq-rest-admin-project/mq-rest-admin-dev-environment.git ../mq-rest-admin-dev-environment

# Start the containerized MQ queue managers
./scripts/dev/mq_start.sh

# Seed deterministic test objects (DEV.* prefix)
./scripts/dev/mq_seed.sh

# Verify REST-based MQSC responses
./scripts/dev/mq_verify.sh

# Stop the queue managers
./scripts/dev/mq_stop.sh

# Reset to clean state (removes data volumes)
./scripts/dev/mq_reset.sh
```

The lifecycle scripts are thin wrappers that delegate to
`../mq-rest-admin-dev-environment`. Override the path with `MQ_DEV_ENV_PATH`.

Integration tests are gated by the `MQ_REST_ADMIN_RUN_INTEGRATION`
environment variable. When unset, integration tests are skipped. For local
runs:

```bash
./scripts/dev/mq_start.sh
./scripts/dev/mq_seed.sh
export MQ_REST_ADMIN_RUN_INTEGRATION=true
go test -race -count=1 -tags=integration ./...
```

Container details:
- Queue managers: `QM1` and `QM2`
- QM1 ports: `1434` (MQ listener), `9463` (REST API)
- QM2 ports: `1435` (MQ listener), `9464` (REST API)
- Admin credentials: `mqadmin` / `mqadmin`
- Object prefix: `DEV.*`

## Architecture

Direct port of `pymqrest`'s architecture, adapted to Go idioms. Uses the
`mq-rest-admin` Java port as a secondary reference.

### Technology Stack

- **HTTP client**: `net/http` (stdlib, zero runtime dependencies)
- **JSON library**: `encoding/json` (stdlib)
- **TLS/mTLS**: `crypto/tls` (stdlib)
- **Mapping data**: `go:embed` (stdlib)
- **Runtime dependencies**: 0 (stdlib only)

### API Surface

- **Style**: Method-per-command mirroring pymqrest (`DisplayQueue()`, `DefineQlocal()`, etc.)
- **Parameters/results**: `map[string]any` for MQ attributes (dynamic attribute sets)
- **Fixed-schema types**: Go structs for `TransportResponse`, `SyncConfig`, `SyncResult`, `EnsureResult`, `MappingIssue`, credential types
- **Credentials**: Interface with unexported method (closed set: `BasicAuth`, `LTPAAuth`, `CertificateAuth`)
- **Session**: Single `Session` struct with functional options constructor
- **Context**: All I/O methods accept `context.Context` as first parameter

### Transport Layer

- `Transport` interface — enables mock-based testing
- `HTTPTransport` implementation using `net/http`
- Supports TLS/mTLS via `crypto/tls`, timeouts via `time.Duration`

### Attribute Mapping

- Direct port of pymqrest's 3-layer pipeline: key map, value map, key-value map
- Two directions: request and response
- Strict/permissive modes with `MappingIssue` tracking
- Mapping data embedded via `go:embed` from JSON file
- Override mechanism with merge/replace modes

### Error Types

Go typed errors (not exception hierarchy), used with `errors.As()`:

```text
TransportError   — network/connection failures
ResponseError    — malformed JSON, unexpected structure
AuthError        — authentication/authorization failures
CommandError     — MQSC command returned error codes
TimeoutError     — polling timeout exceeded
MappingError     — attribute translation failures (strict mode)
```

### Package Structure

Single flat package under `mqrestadmin/`. Key files:

- `session.go` — `Session` struct, `NewSession`, functional options, core `mqscCommand` dispatch
- `session_commands.go` — All MQSC command methods (DISPLAY, DEFINE, ALTER, DELETE, START, STOP, etc.)
- `session_ensure.go` / `session_sync.go` — Idempotent ensure and synchronous polling methods
- `auth.go` — `Credentials` interface (closed set via unexported method), `BasicAuth`, `LTPAAuth`, `CertificateAuth`
- `transport.go` — `Transport` interface, `HTTPTransport`
- `mapping.go` + `mapping_data.go` + `mapping_data.json` — 3-layer attribute translation pipeline
- `errors.go` — All typed error types

### Testing Patterns

Tests use a `mockTransport` (in `mock_test.go`) that records calls and returns
pre-configured responses. Helper constructors:

- `newTestSession(transport)` — mapping disabled, for simple assertions
- `newTestSessionWithMapping(transport)` — mapping enabled
- `newTestSessionWithClock(transport, clock)` — mock clock for sync/polling tests
- `generateSelfSignedCert(t)` — creates test TLS certificates (in `testhelpers_test.go`)

## Key References

**Reference implementation**: `../pymqrest` (Python version)

**Java port reference**: `../mq-rest-admin` (Java version)

**External Documentation**:

- IBM MQ 9.4 administrative REST API
- MQSC command reference
