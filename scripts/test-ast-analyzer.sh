#!/bin/bash
set -e

# Directory containing this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Build the AST analyzer test tool
echo "Building the AST analyzer test tool..."
cd "${PROJECT_ROOT}"
go build -o "${PROJECT_ROOT}/bin/ast-analyzer" "${PROJECT_ROOT}/cmd/ast-analyzer/main.go"

# Run tests on each test file in the testdata/ast_analysis directory
echo "Running AST analysis tests..."
echo "=============================="

# Error handling tests
echo "Testing error handling patterns..."
"${PROJECT_ROOT}/bin/ast-analyzer" analyze \
    --file "${PROJECT_ROOT}/testdata/ast_analysis/error_handling_test.go" \
    --rules error_handling

# API design tests
echo -e "\nTesting API design patterns..."
"${PROJECT_ROOT}/bin/ast-analyzer" analyze \
    --file "${PROJECT_ROOT}/testdata/ast_analysis/api_design_test.go" \
    --rules api_design

# Security and concurrency tests
echo -e "\nTesting security and concurrency patterns..."
"${PROJECT_ROOT}/bin/ast-analyzer" analyze \
    --file "${PROJECT_ROOT}/testdata/ast_analysis/security_concurrency_test.go" \
    --rules concurrent_map_access,secure_coding

echo -e "\nAST analysis testing complete"
