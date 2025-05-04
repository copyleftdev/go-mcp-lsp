#!/bin/bash
set -e

# Directory containing this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Ensure the mcplsp tool is built
cd "${PROJECT_ROOT}"
go build -o "${PROJECT_ROOT}/bin/mcplsp" "${PROJECT_ROOT}/cmd/mcplsp/main.go"

# Test AST files
TEST_FILES=(
    "testdata/ast_analysis/error_handling_test.go"
    "testdata/ast_analysis/api_design_test.go"
    "testdata/ast_analysis/security_concurrency_test.go"
)

# Rule combinations to test
RULE_SETS=(
    "error_handling"
    "api_design"
    "secure_coding,concurrent_map_access"
)

# Ensure the MCP server is running
SERVER_RUNNING=$(ps aux | grep "mcpserver" | grep -v grep | wc -l)
if [ $SERVER_RUNNING -eq 0 ]; then
    echo "Starting MCP server..."
    "${PROJECT_ROOT}/bin/mcpserver" &
    SERVER_PID=$!
    # Give the server a moment to start
    sleep 2
    echo "MCP server started with PID: $SERVER_PID"
else
    echo "MCP server already running"
fi

# Run tests comparing direct AST analysis vs. server-based validation
for i in "${!TEST_FILES[@]}"; do
    TEST_FILE="${TEST_FILES[$i]}"
    RULE_SET="${RULE_SETS[$i]}"
    
    echo "======================================================"
    echo "Testing file: ${TEST_FILE}"
    echo "Rules: ${RULE_SET}"
    echo "------------------------------------------------------"
    
    echo "Running with direct AST analysis (-deep):"
    "${PROJECT_ROOT}/bin/mcplsp" -deep validate "${PROJECT_ROOT}/${TEST_FILE}" "${RULE_SET}" || true
    
    echo ""
    echo "Running with server-based validation:"
    "${PROJECT_ROOT}/bin/mcplsp" validate "${PROJECT_ROOT}/${TEST_FILE}" "${RULE_SET}" || true
    
    echo "======================================================"
    echo ""
done

# Kill the server if we started it
if [ -n "$SERVER_PID" ]; then
    echo "Stopping MCP server (PID: $SERVER_PID)..."
    kill $SERVER_PID
fi

echo "Validation comparison complete"
