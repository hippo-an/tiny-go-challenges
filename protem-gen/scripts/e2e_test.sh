#!/bin/bash

# E2E Test Script for protem-gen
# This script tests the full project generation flow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test directory
TEST_DIR="/tmp/protem-gen-e2e-test"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Counters
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((TESTS_FAILED++))
}

cleanup() {
    log_info "Cleaning up test directory..."
    rm -rf "$TEST_DIR"
}

# Build protem-gen binary
build_binary() {
    log_info "Building protem-gen binary..."
    cd "$PROJECT_ROOT"
    go build -o "$TEST_DIR/protem-gen" .
    if [ $? -eq 0 ]; then
        log_success "Binary built successfully"
    else
        log_error "Failed to build binary"
        exit 1
    fi
}

# Test: Version command
test_version() {
    log_info "Testing version command..."
    if "$TEST_DIR/protem-gen" version | grep -q "protem-gen"; then
        log_success "Version command works"
    else
        log_error "Version command failed"
    fi
}

# Test: Generate project with postgres
test_generate_postgres() {
    local project_name="test-postgres-app"
    local project_dir="$TEST_DIR/$project_name"

    log_info "Testing project generation with PostgreSQL..."

    # Generate project using flags
    cd "$TEST_DIR"
    if "$TEST_DIR/protem-gen" create \
        --name "$project_name" \
        --module "github.com/test/$project_name" \
        --database postgres \
        --no-interactive 2>&1; then
        log_success "Project generated successfully"
    else
        log_error "Project generation failed"
        return 1
    fi

    # Verify directory structure
    local required_dirs=(
        "cmd/server"
        "internal/domain"
        "internal/application"
        "internal/infrastructure/database"
        "internal/infrastructure/http"
        "internal/interfaces/http"
        "web/templates/layouts"
        "web/templates/pages"
        "web/static/css"
        "migrations"
        "sqlc/queries"
    )

    for dir in "${required_dirs[@]}"; do
        if [ -d "$project_dir/$dir" ]; then
            log_success "Directory exists: $dir"
        else
            log_error "Missing directory: $dir"
        fi
    done

    # Verify key files
    local required_files=(
        "go.mod"
        "Makefile"
        ".gitignore"
        "package.json"
        ".air.toml"
        "cmd/server/main.go"
    )

    for file in "${required_files[@]}"; do
        if [ -f "$project_dir/$file" ]; then
            log_success "File exists: $file"
        else
            log_error "Missing file: $file"
        fi
    done

    # Verify go.mod contains correct module path
    if grep -q "github.com/test/$project_name" "$project_dir/go.mod"; then
        log_success "go.mod has correct module path"
    else
        log_error "go.mod has incorrect module path"
    fi

    # Test build
    log_info "Testing project build..."
    cd "$project_dir"
    if go build ./...; then
        log_success "Project builds successfully"
    else
        log_error "Project build failed"
    fi
}

# Test: Generate project with SQLite
test_generate_sqlite() {
    local project_name="test-sqlite-app"
    local project_dir="$TEST_DIR/$project_name"

    log_info "Testing project generation with SQLite..."

    cd "$TEST_DIR"
    if "$TEST_DIR/protem-gen" create \
        --name "$project_name" \
        --module "github.com/test/$project_name" \
        --database sqlite \
        --no-interactive 2>&1; then
        log_success "SQLite project generated successfully"
    else
        log_error "SQLite project generation failed"
        return 1
    fi

    # Verify SQLite-specific config
    if [ -f "$project_dir/sqlc.yaml" ]; then
        if grep -q "sqlite" "$project_dir/sqlc.yaml"; then
            log_success "sqlc.yaml configured for SQLite"
        else
            log_error "sqlc.yaml not configured for SQLite"
        fi
    fi

    # Test build
    cd "$project_dir"
    if go build ./...; then
        log_success "SQLite project builds successfully"
    else
        log_error "SQLite project build failed"
    fi
}

# Test: Generate project with all options
test_generate_all_options() {
    local project_name="test-full-app"
    local project_dir="$TEST_DIR/$project_name"

    log_info "Testing project generation with all options..."

    cd "$TEST_DIR"
    if "$TEST_DIR/protem-gen" create \
        --name "$project_name" \
        --module "github.com/test/$project_name" \
        --database postgres \
        --grpc \
        --auth \
        --ai \
        --no-interactive 2>&1; then
        log_success "Full-featured project generated successfully"
    else
        log_error "Full-featured project generation failed"
        return 1
    fi

    # Verify optional directories
    local optional_dirs=(
        "internal/interfaces/grpc"
        "proto"
        "internal/infrastructure/auth"
        "internal/infrastructure/llm"
        "internal/infrastructure/prompt"
        "internal/infrastructure/stream"
    )

    for dir in "${optional_dirs[@]}"; do
        if [ -d "$project_dir/$dir" ]; then
            log_success "Optional directory exists: $dir"
        else
            log_error "Missing optional directory: $dir"
        fi
    done

    # Test build
    cd "$project_dir"
    if go build ./...; then
        log_success "Full-featured project builds successfully"
    else
        log_error "Full-featured project build failed"
    fi
}

# Main test runner
main() {
    echo ""
    echo "=========================================="
    echo "  protem-gen E2E Test Suite"
    echo "=========================================="
    echo ""

    # Setup
    cleanup
    mkdir -p "$TEST_DIR"

    # Build
    build_binary

    # Run tests
    test_version
    test_generate_postgres
    test_generate_sqlite
    test_generate_all_options

    # Cleanup
    cleanup

    # Summary
    echo ""
    echo "=========================================="
    echo "  Test Summary"
    echo "=========================================="
    echo -e "  ${GREEN}Passed:${NC} $TESTS_PASSED"
    echo -e "  ${RED}Failed:${NC} $TESTS_FAILED"
    echo "=========================================="
    echo ""

    if [ $TESTS_FAILED -gt 0 ]; then
        exit 1
    fi

    exit 0
}

# Run main
main "$@"
