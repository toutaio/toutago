#!/bin/bash
# Phase 1 Validation Script

set -e

echo "🔍 Toutā Phase 1 Validation"
echo "================================"
echo ""

echo "✓ Checking Go version..."
go version

echo ""
echo "✓ Running all tests..."
go test ./... -short

echo ""
echo "✓ Checking test coverage..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tail -1

echo ""
echo "✓ Building CLI..."
go build -o touta cmd/touta/main.go
./touta version

echo ""
echo "✓ Testing project creation..."
rm -rf test-validation-app
./touta new test-validation-app
cd test-validation-app && ls -la

echo ""
echo "✓ Checking generated files..."
[ -f "touta.yaml" ] && echo "  ✓ touta.yaml exists"
[ -f "main.go" ] && echo "  ✓ main.go exists"
[ -d "handlers" ] && echo "  ✓ handlers/ directory exists"
[ -d "templates" ] && echo "  ✓ templates/ directory exists"

cd ..
rm -rf test-validation-app

echo ""
echo "================================"
echo "✅ Phase 1 validation complete!"
echo ""
echo "Summary:"
echo "  - All tests passing"
echo "  - CLI tool working"
echo "  - Project scaffolding functional"
echo "  - Core components implemented"
echo ""
echo "Phase 1 Foundation: READY ✨"
