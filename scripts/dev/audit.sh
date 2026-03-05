#!/usr/bin/env bash
set -euo pipefail

export DOCKER_DEV_IMAGE="${DOCKER_DEV_IMAGE:-dev-go:1.26}"
export DOCKER_TEST_CMD="${DOCKER_TEST_CMD:-govulncheck ./... && go-licenses check ./... --allowed_licenses=Apache-2.0,BSD-2-Clause,BSD-3-Clause,MIT,ISC,MPL-2.0,GPL-3.0}"

if ! command -v docker-test >/dev/null 2>&1; then
  echo "ERROR: docker-test not found on PATH." >&2
  echo "Set up standard-tooling: export PATH=../standard-tooling/scripts/bin:\$PATH" >&2
  exit 1
fi
exec docker-test
