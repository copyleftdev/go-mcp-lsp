#!/bin/bash
# Test our governance rules against the comprehensive test data

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

echo "=== Testing MCP Governance Rules ==="
echo ""

# Error Handling Rules
echo "=== Error Handling Rules ==="
echo "Testing file with missing error checks (should fail):"
./bin/mcplsp validate ./testdata/error_handling/missing_check.go
echo ""

echo "Testing file with proper error handling (should pass):"
./bin/mcplsp validate ./testdata/error_handling/proper_check.go
echo ""

# API Design Rules
echo "=== API Design Rules ==="
echo "Testing file missing context parameters (should fail):"
./bin/mcplsp validate ./testdata/api_design/missing_context.go
echo ""

echo "Testing file with proper context usage (should pass):"
./bin/mcplsp validate ./testdata/api_design/proper_context.go
echo ""

# Security Rules
echo "=== Security Rules ==="
echo "Testing file with insecure practices (should fail):"
./bin/mcplsp validate ./testdata/security/insecure_practices.go
echo ""

echo "Testing file with secure practices (should pass):"
./bin/mcplsp validate ./testdata/security/secure_practices.go
echo ""

# Organization Rules
echo "=== Organization Standards ==="
echo "Testing file violating organizational standards (should fail):"
./bin/mcplsp validate ./testdata/organization/non_compliant.go
echo ""

echo "Testing file following organizational standards (should pass):"
./bin/mcplsp validate ./testdata/organization/compliant.go
echo ""

# Concurrency Rules
echo "=== Concurrency Rules ==="
echo "Testing file with race conditions (should fail):"
./bin/mcplsp validate ./testdata/concurrency/race_condition.go
echo ""

echo "Testing file with proper synchronization (should pass):"
./bin/mcplsp validate ./testdata/concurrency/safe_access.go
echo ""

echo "=== Testing Complete ==="
