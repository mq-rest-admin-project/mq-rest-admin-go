# Repository Standards

## Table of Contents

- [Pre-flight checklist](#pre-flight-checklist)
- [Local validation](#local-validation)
- [Tooling requirement](#tooling-requirement)
- [Merge strategy override](#merge-strategy-override)
- [Approved domain abbreviations](#approved-domain-abbreviations)
- [Accepted naming deviations](#accepted-naming-deviations)

## Pre-flight checklist

- Before modifying any files, check the current branch with `git status -sb`.
- If on `develop`, create a short-lived `feature/*` branch or ask for explicit approval to proceed on `develop`.
- If approval is granted to work on `develop`, call it out in the response and proceed only for that user-approved scope.
- Enable repository git hooks before committing: `git config core.hooksPath .githooks`.

## Local validation

Canonical local validation command:

```bash
vrg-docker-run -- vrg-validate
```

## Tooling requirement

Required for daily workflow:

- Go 1.22+ (`brew install go` or <https://go.dev/dl/>)
- `golangci-lint` (`brew install golangci-lint`)
- `gocyclo` (`go install github.com/fzipp/gocyclo/cmd/gocyclo@latest`)
- `govulncheck` (`go install golang.org/x/vuln/cmd/govulncheck@latest`)
- `vrg-docker-run -- vrg-validate` (canonical validation)

## Merge strategy override

- Feature, bugfix, and chore PRs targeting `develop` use squash merges (`--squash`).
- Release PRs targeting `main` use regular merges (`--merge`) to preserve shared
  ancestry between `main` and `develop`.
- Auto-merge commands:
  - Feature PRs: `gh pr merge --auto --squash --delete-branch`
  - Release PRs: `gh pr merge --auto --merge --delete-branch`

## Commit and PR scripts

AI agents **must** use the wrapper scripts for commits and PR
submission. Do not construct commit messages or PR bodies manually.

### Committing

```bash
vrg-commit \
  --type TYPE --message MESSAGE \
  [--scope SCOPE] [--body BODY]
```

- `--type` (required): one of
  `feat|fix|docs|style|refactor|test|chore|ci|build`
- `--message` (required): commit description
- `--scope` (optional): conventional commit scope
- `--body` (optional): detailed commit body

The script resolves the correct `Co-Authored-By` identity from
`vergil.toml` and the git hooks validate the result.

### Submitting PRs

```bash
vrg-submit-pr \
  --issue NUMBER --summary TEXT \
  [--linkage KEYWORD] [--title TEXT] \
  [--notes TEXT] [--dry-run]
```

- `--issue` (required): GitHub issue number (just the number)
- `--summary` (required): one-line PR summary
- `--linkage` (optional, default: `Ref`):
  `Fixes|Closes|Resolves|Ref`
- `--title` (optional): PR title (default: most recent commit
  subject)
- `--notes` (optional): additional notes
- `--dry-run` (optional): print generated PR without executing

The script detects the target branch and merge strategy
automatically.

## Approved domain abbreviations

Domain-specific abbreviations that are well-established in the IBM MQ ecosystem
and may be used in identifiers without expansion:

- `qmgr` — queue manager (established MQSC domain term)

## Accepted naming deviations

None yet.
