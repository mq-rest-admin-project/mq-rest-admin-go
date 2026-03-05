#!/usr/bin/env bash
set -euo pipefail

export DOCKER_DEV_IMAGE="${DOCKER_DEV_IMAGE:-dev-go:1.26}"
export DOCKER_TEST_CMD="${DOCKER_TEST_CMD:-go vet ./... && go test -race -count=1 -coverprofile=coverage.out ./... && go-test-coverage --config .testcoverage.yml}"

if ! command -v docker-test >/dev/null 2>&1; then
  echo "ERROR: docker-test not found on PATH." >&2
  echo "Set up standard-tooling: export PATH=../standard-tooling/scripts/bin:\$PATH" >&2
  exit 1
fi
exec docker-test
