#!/bin/bash
# Test governance rules against our test data

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== Testing Intent-Based Governance Rules ==="
echo ""

test_file() {
    local file=$1
    local rule=$2
    local expected=$3
    
    echo "Testing $file against rule: $rule"
    echo "Expected result: $expected"
    
    $PROJECT_ROOT/bin/mcplsp validate "$PROJECT_ROOT/$file" "$rule" > /tmp/validation_result.txt 2>&1
    exit_code=$?
    
    if grep -q "Validation passed" /tmp/validation_result.txt; then
        if [ "$expected" == "pass" ]; then
            echo "✅ PASS: File passed validation as expected"
        else
            echo "❌ FAIL: File unexpectedly passed validation"
        fi
    else
        if [ "$expected" == "fail" ]; then
            echo "✅ PASS: File failed validation as expected"
            grep -A5 "Validation failed:" /tmp/validation_result.txt || true
        else
            echo "❌ FAIL: File unexpectedly failed validation"
            cat /tmp/validation_result.txt
        fi
    fi
    echo ""
}

# Error Handling Tests
echo "=== Error Handling Governance Tests ==="
test_file "testdata/error_handling/missing_check.go" "error_handling" "fail"
test_file "testdata/error_handling/proper_check.go" "error_handling" "pass"

# API Design Tests
echo "=== API Design Governance Tests ==="
test_file "testdata/api_design/missing_context.go" "api_design" "fail"
test_file "testdata/api_design/proper_context.go" "api_design" "pass"

# Concurrency Tests
echo "=== Concurrency Governance Tests ==="
test_file "testdata/concurrency/race_condition.go" "concurrent_map_access" "fail"
test_file "testdata/concurrency/safe_access.go" "concurrent_map_access" "pass"

# Security Tests
echo "=== Security Governance Tests ==="
test_file "testdata/security/insecure_practices.go" "secure_coding" "fail"
test_file "testdata/security/secure_practices.go" "secure_coding" "pass"

# Organization Standards Tests
echo "=== Organization Standards Tests ==="
test_file "testdata/organization/non_compliant.go" "org_coding_standards" "fail"
test_file "testdata/organization/compliant.go" "org_coding_standards" "pass"

echo "=== All Tests Completed ==="
