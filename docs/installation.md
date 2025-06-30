# Installation Guide

This guide covers all the ways to install Hawk TUI on your system.

## Quick Installation

### Install Script (Recommended)

The easiest way to install Hawk TUI is using our install scripts:

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/hawk-tui/hawk-tui/main/scripts/install.sh | bash

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/hawk-tui/hawk-tui/main/scripts/install.ps1 | iex
```

The install script will:
- Detect your platform and architecture
- Download the appropriate binary
- Install to the correct location
- Update your PATH if needed
- Verify the installation

## Package Managers

### Homebrew (macOS/Linux)
```bash
brew install hawk-tui/tap/hawk-tui
```

### NPM (Global)
```bash
npm install -g hawk-tui
```

### Python (PyPI)
```bash
pip install hawk-tui
```

### Docker
```bash
# Run directly
docker run -it --rm hawktui/hawk-tui:latest

# Or pull the image
docker pull hawktui/hawk-tui:latest
```

## Manual Installation

### Download Binary

1. Go to [GitHub Releases](https://github.com/hawk-tui/hawk-tui/releases)
2. Download the appropriate binary for your platform:
   - Linux: `hawk-tui-linux-amd64.tar.gz`
   - macOS: `hawk-tui-darwin-amd64.tar.gz` or `hawk-tui-darwin-arm64.tar.gz`
   - Windows: `hawk-tui-windows-amd64.zip`

3. Extract and install:

```bash
# Linux/macOS
tar -xzf hawk-tui-*.tar.gz
sudo mv hawk /usr/local/bin/

# Windows
# Extract the ZIP file and move hawk.exe to a directory in your PATH
```

### Build from Source

Requirements:
- Go 1.21 or later
- Git
- Make

```bash
git clone https://github.com/hawk-tui/hawk-tui.git
cd hawk-tui
make build
sudo make install
```

## Verification

After installation, verify that Hawk TUI is working:

```bash
hawk --version
hawk --help
```

You should see the version information and help text.

## Platform-Specific Notes

### Linux
- The install script will place the binary in `/usr/local/bin` (system-wide) or `~/.local/bin` (user-only)
- Make sure `~/.local/bin` is in your PATH for user installations

### macOS
- Same as Linux
- On Apple Silicon Macs, use the `arm64` version for better performance
- You may need to allow the binary in System Preferences > Security & Privacy

### Windows
- The PowerShell script will install to `%USERPROFILE%\.local\bin` or `%PROGRAMFILES%\Hawk TUI`
- The script will automatically add the installation directory to your PATH
- You may need to restart your terminal or run: `$env:PATH += ';C:\path\to\hawk'`

## Uninstallation

### Linux/macOS
```bash
# If installed system-wide
sudo rm /usr/local/bin/hawk

# If installed to user directory
rm ~/.local/bin/hawk
```

### Windows
```powershell
# Remove from user installation
Remove-Item "$env:USERPROFILE\.local\bin\hawk.exe"

# Remove from system installation (as admin)
Remove-Item "$env:PROGRAMFILES\Hawk TUI\hawk.exe"
```

### Package Managers
```bash
# Homebrew
brew uninstall hawk-tui

# NPM
npm uninstall -g hawk-tui

# Python
pip uninstall hawk-tui
```

## Troubleshooting

### Permission Denied
If you get permission errors, try:
- Using `sudo` for system-wide installation
- Installing to user directory instead
- Checking that the binary has execute permissions: `chmod +x hawk`

### Command Not Found
If `hawk` command is not found:
- Check that the installation directory is in your PATH
- Restart your terminal
- Verify the binary exists in the expected location

### Download Failures
If the install script fails to download:
- Check your internet connection
- Try downloading manually from GitHub Releases
- Use a VPN if you're behind a corporate firewall

### Antivirus Issues
Some antivirus software may flag the binary:
- Add an exception for the Hawk TUI installation directory
- Download from official sources only
- Verify the binary's signature if available

## Next Steps

After installation, check out the [Quick Start Guide](quickstart.md) to start using Hawk TUI with your applications.