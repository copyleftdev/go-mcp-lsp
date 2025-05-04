#!/bin/bash
set -e

# Directory containing this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Ensure mcplsp is built
cd "${PROJECT_ROOT}"
go build -o "${PROJECT_ROOT}/bin/mcplsp" "${PROJECT_ROOT}/cmd/mcplsp/main.go"

# Test directories and their corresponding rules
declare -A TEST_DIRS=(
    ["testdata/error_handling"]="error_handling"
    ["testdata/api_design"]="api_design"
    ["testdata/concurrency"]="concurrent_map_access"
    ["testdata/security"]="secure_coding"
    ["testdata/organization"]="org_coding_standards"
    ["testdata/ast_analysis"]="error_handling,api_design,concurrent_map_access,secure_coding"
)

# Test modes
MODES=("standard" "deep")

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

TOTAL_TESTS=0
FAILED_TESTS=0
PASSED_TESTS=0

for TEST_DIR in "${!TEST_DIRS[@]}"; do
    RULES="${TEST_DIRS[$TEST_DIR]}"
    
    echo "====================================================="
    echo "Testing directory: ${TEST_DIR}"
    echo "Rules: ${RULES}"
    echo "-----------------------------------------------------"
    
    # Find all Go files in the directory
    if [ -d "${PROJECT_ROOT}/${TEST_DIR}" ]; then
        for MODE in "${MODES[@]}"; do
            echo "Analysis mode: ${MODE}"
            
            DEEP_FLAG=""
            if [ "$MODE" == "deep" ]; then
                DEEP_FLAG="-deep"
            fi
            
            for TEST_FILE in $(find "${PROJECT_ROOT}/${TEST_DIR}" -name "*.go"); do
                TOTAL_TESTS=$((TOTAL_TESTS + 1))
                BASENAME=$(basename "$TEST_FILE")
                
                echo "Testing ${BASENAME}..."
                # Using positional arguments as per CLI design: validate <file> [rule1,rule2,...]
                if "${PROJECT_ROOT}/bin/mcplsp" ${DEEP_FLAG} validate "${TEST_FILE}" "${RULES}" > /dev/null 2>&1; then
                    echo "✅ ${BASENAME} passed"
                    PASSED_TESTS=$((PASSED_TESTS + 1))
                else
                    echo "❌ ${BASENAME} failed"
                    FAILED_TESTS=$((FAILED_TESTS + 1))
                    
                    # Show detailed output for failed tests
                    echo "Detailed output:"
                    "${PROJECT_ROOT}/bin/mcplsp" ${DEEP_FLAG} validate "${TEST_FILE}" "${RULES}"
                fi
            done
            
            echo ""
        done
    else
        echo "Directory ${TEST_DIR} not found, skipping..."
    fi
    
    echo "====================================================="
    echo ""
done

# Kill the server if we started it
if [ -n "$SERVER_PID" ]; then
    echo "Stopping MCP server (PID: $SERVER_PID)..."
    kill $SERVER_PID
fi

echo "Governance testing complete"
echo "Total tests: ${TOTAL_TESTS}"
echo "Passed: ${PASSED_TESTS}"
echo "Failed: ${FAILED_TESTS}"

if [ $FAILED_TESTS -gt 0 ]; then
    exit 1
fi
