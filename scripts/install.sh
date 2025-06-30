#!/bin/bash

# Hawk TUI Installation Script
# Universal installer for Linux/macOS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/hawk-tui/hawk-tui"
BINARY_NAME="hawk"
INSTALL_DIR="/usr/local/bin"
VERSION="latest"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for system-wide install
check_permissions() {
    if [[ $EUID -eq 0 ]]; then
        print_status "Installing system-wide as root"
        INSTALL_DIR="/usr/local/bin"
    else
        print_status "Installing to user directory"
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        
        # Add to PATH if not already there
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            print_warning "Adding $INSTALL_DIR to PATH in ~/.bashrc"
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> ~/.bashrc
            print_warning "Please run 'source ~/.bashrc' or restart your terminal"
        fi
    fi
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
        *) 
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    case $OS in
        linux|darwin) ;;
        *)
            print_error "Unsupported operating system: $OS"
            print_error "Please use the manual installation method"
            exit 1
            ;;
    esac
    
    print_status "Detected platform: $OS-$ARCH"
}

# Check dependencies
check_dependencies() {
    print_status "Checking dependencies..."
    
    if ! command -v curl >/dev/null 2>&1; then
        print_error "curl is required but not installed"
        print_error "Please install curl and try again"
        exit 1
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        print_error "tar is required but not installed"
        exit 1
    fi
    
    print_success "All dependencies found"
}

# Get latest release version
get_latest_version() {
    if [[ "$VERSION" == "latest" ]]; then
        print_status "Fetching latest release version..."
        VERSION=$(curl -s "https://api.github.com/repos/hawk-tui/hawk-tui/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [[ -z "$VERSION" ]]; then
            print_warning "Could not fetch latest version, using v1.0.0"
            VERSION="v1.0.0"
        fi
    fi
    print_status "Installing version: $VERSION"
}

# Download and install binary
install_binary() {
    DOWNLOAD_URL="https://github.com/hawk-tui/hawk-tui/releases/download/$VERSION/hawk-tui-$VERSION-$OS-$ARCH.tar.gz"
    TEMP_DIR=$(mktemp -d)
    
    print_status "Downloading Hawk TUI from $DOWNLOAD_URL"
    
    if ! curl -L "$DOWNLOAD_URL" -o "$TEMP_DIR/hawk-tui.tar.gz"; then
        print_error "Failed to download Hawk TUI"
        print_error "URL: $DOWNLOAD_URL"
        rm -rf "$TEMP_DIR"
        exit 1
    fi
    
    print_status "Extracting archive..."
    tar -xzf "$TEMP_DIR/hawk-tui.tar.gz" -C "$TEMP_DIR"
    
    print_status "Installing binary to $INSTALL_DIR"
    if [[ $EUID -eq 0 ]]; then
        cp "$TEMP_DIR/hawk" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/hawk"
    else
        cp "$TEMP_DIR/hawk" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/hawk"
    fi
    
    rm -rf "$TEMP_DIR"
    print_success "Binary installed successfully"
}

# Install from source (fallback)
install_from_source() {
    print_status "Installing from source..."
    
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is required for source installation"
        print_error "Please install Go 1.21+ or use a binary release"
        exit 1
    fi
    
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    print_status "Cloning repository..."
    git clone "$REPO_URL.git" .
    
    print_status "Building binary..."
    go build -o hawk ./cmd/hawk
    
    print_status "Installing binary to $INSTALL_DIR"
    if [[ $EUID -eq 0 ]]; then
        cp hawk "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/hawk"
    else
        cp hawk "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/hawk"
    fi
    
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
    print_success "Source installation completed"
}

# Verify installation
verify_installation() {
    if command -v hawk >/dev/null 2>&1; then
        INSTALLED_VERSION=$(hawk --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' || echo "unknown")
        print_success "Hawk TUI installed successfully!"
        print_success "Version: $INSTALLED_VERSION"
        print_success "Location: $(which hawk)"
    else
        print_error "Installation verification failed"
        print_error "hawk command not found in PATH"
        exit 1
    fi
}

# Show usage examples
show_examples() {
    echo ""
    print_success "Installation complete! Here are some quick examples:"
    echo ""
    echo "  # Basic usage with any application:"
    echo "  your-app | hawk"
    echo ""
    echo "  # With Node.js application:"
    echo "  node app.js | hawk"
    echo ""
    echo "  # With Python application:"
    echo "  python app.py | hawk"
    echo ""
    echo "  # Show help:"
    echo "  hawk --help"
    echo ""
    echo "  # Show version:"
    echo "  hawk --version"
    echo ""
    print_status "Visit https://hawktui.dev for documentation and examples"
}

# Main installation function
main() {
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                              🦅 Hawk TUI Installer                          ║"
    echo "║                     Universal TUI Framework for Any Language                ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo ""
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --source)
                SOURCE_INSTALL=true
                shift
                ;;
            --help)
                echo "Hawk TUI Installer"
                echo ""
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --version VERSION    Install specific version (default: latest)"
                echo "  --install-dir DIR    Installation directory (default: /usr/local/bin or ~/.local/bin)"
                echo "  --source             Install from source instead of binary"
                echo "  --help               Show this help message"
                echo ""
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                print_error "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    check_permissions
    detect_platform
    check_dependencies
    get_latest_version
    
    if [[ "$SOURCE_INSTALL" == "true" ]]; then
        install_from_source
    else
        install_binary || {
            print_warning "Binary installation failed, trying source installation..."
            install_from_source
        }
    fi
    
    verify_installation
    show_examples
    
    echo ""
    print_success "🎉 Happy monitoring with Hawk TUI!"
}

# Run main function
main "$@"