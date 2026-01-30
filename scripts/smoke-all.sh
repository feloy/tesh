#!/bin/sh

set -e

./scripts/smoke-run.sh > /dev/null 2>&1
./scripts/smoke-scenarios.sh > /dev/null 2>&1
./scripts/smoke-expects.sh > /dev/null 2>&1
./scripts/smoke-calls.sh > /dev/null 2>&1

echo "All smoke tests passed"
