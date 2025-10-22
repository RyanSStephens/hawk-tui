#!/bin/bash

# Script to generate all HawkTUI demo GIFs using VHS

set -e

echo "🎬 Generating HawkTUI demos..."
echo ""

# Check if vhs is installed
if ! command -v vhs &> /dev/null; then
    echo "❌ VHS is not installed. Please install it first:"
    echo "   go install github.com/charmbracelet/vhs@latest"
    echo "   or visit: https://github.com/charmbracelet/vhs"
    exit 1
fi

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR/.."

# Build the examples first
echo "🔨 Building examples..."
go build -o /tmp/simple_demo examples/hawktui/simple_demo.go
go build -o /tmp/dashboard_demo examples/hawktui/dashboard_demo.go
go build -o /tmp/form_demo examples/hawktui/form_demo.go
go build -o /tmp/list_demo examples/hawktui/list_demo.go
echo "✅ Examples built"
echo ""

# Generate each demo
demos=(
    "simple_demo"
    "dashboard_demo"
    "form_demo"
    "list_demo"
)

for demo in "${demos[@]}"; do
    echo "📹 Recording $demo..."
    vhs "$SCRIPT_DIR/${demo}.tape"
    if [ -f "$SCRIPT_DIR/${demo}.gif" ]; then
        echo "✅ $demo.gif created"
    else
        echo "❌ Failed to create $demo.gif"
    fi
    echo ""
done

echo "🎉 All demos generated successfully!"
echo ""
echo "Demo files created in: $SCRIPT_DIR/"
ls -lh "$SCRIPT_DIR"/*.gif 2>/dev/null || echo "No GIF files found"
