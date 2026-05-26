# Developer Setup

## Prerequisites

| Tool | Version | Purpose |
| --- | --- | --- |
| Go | 1.25+ | Build and test |
| git | Latest | Version control |
| Docker | Latest | Local MQ containers (integration tests) |

No external build tools, task runners, or code generators are required.
The Go toolchain provides everything needed.

## Required repositories

mqrestadmin depends on two sibling repositories:

| Repository | Purpose |
| --- | --- |
| [mq-rest-admin-go](https://github.com/mq-rest-admin-project/mq-rest-admin-go) | This project |
| [vergil-tooling](https://github.com/vergil-project/vergil-tooling) | Canonical project standards (referenced by `AGENTS.md` and git hooks) |
| [mq-rest-admin-dev-environment](https://github.com/mq-rest-admin-project/mq-rest-admin-dev-environment) | Dockerized MQ test infrastructure (local and CI) |

## Recommended directory layout

Clone all three repositories as siblings:

```text
~/dev/
├── mq-rest-admin-go/
├── vergil-tooling/
└── mq-rest-admin-dev-environment/
```

```bash
cd ~/dev
git clone https://github.com/mq-rest-admin-project/mq-rest-admin-go.git
git clone https://github.com/vergil-project/vergil-tooling.git
git clone https://github.com/mq-rest-admin-project/mq-rest-admin-dev-environment.git
```

## Building

```bash
go build ./...          # Compile all packages
go vet ./...            # Run static analysis
```

There is no separate clean step -- Go caches builds automatically and
invalidates when source files change.

## Testing

```bash
go test ./...                          # Run all unit tests
go test -v ./...                       # Verbose output
go test -race ./...                    # Race detector enabled
go test -cover ./...                   # Show coverage summary
go test -coverprofile=coverage.out ./...  # Generate coverage profile
go tool cover -html=coverage.out       # View coverage in browser
```

## Validation

Run all quality checks:

```bash
go build ./...          # Compile
go vet ./...            # Static analysis
go test -race ./...     # Tests with race detector
```

## Claude Code hook guard

The `.claude/hooks/guard.sh` PreToolUse hook blocks raw `git` and
`gh` commands in AI agent sessions — all operations must go through
the `vrg-git` / `vrg-gh` wrappers.

## Documentation

### Local setup

```bash
# Set up shared fragments symlink
scripts/dev/docs-setup.sh

# Install MkDocs
pip install mkdocs-material

# Build the documentation site
mkdocs build -f docs/site/mkdocs.yml

# Serve locally with live reload
mkdocs serve -f docs/site/mkdocs.yml
```
