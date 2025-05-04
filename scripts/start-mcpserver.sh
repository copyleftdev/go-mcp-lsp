#!/bin/bash
# Start the MCP server for testing

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT/server/mcpserver"

# Make sure rules and templates directories exist
mkdir -p rules templates

echo "[info] Starting MCP Server on localhost:9000..."
go run cmd/main.go --address=localhost:9000 --rules=./rules --templates=./templates
