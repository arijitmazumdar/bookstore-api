#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

./scripts/run-tests.sh
./scripts/run-integration-tests.sh
