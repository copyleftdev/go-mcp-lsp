#!/bin/bash
# Test MVP functionality of the MCP integration

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "[test] Building CLI tool..."
cd "$PROJECT_ROOT"
go build -o ./bin/mcplsp ./cmd/mcplsp

echo "[test] Starting MCP server in background..."
cd "$PROJECT_ROOT"
./scripts/start-mcpserver.sh &
SERVER_PID=$!

# Give the server a moment to start
sleep 1

echo "[test] Testing connection to MCP server..."
./bin/mcplsp test

echo "[test] Validating sample file..."
./bin/mcplsp validate ./pkg/mcpclient/test/sample.go

echo "[test] Running integration tests..."
cd "$PROJECT_ROOT/pkg/mcpclient/test"
go test -v

echo "[test] Cleaning up..."
kill $SERVER_PID

echo "[test] MVP validation complete"
