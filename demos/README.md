# HawkTUI Demos

This directory contains VHS tape files for recording animated GIF demos of HawkTUI examples.

## About VHS

[VHS](https://github.com/charmbracelet/vhs) is a tool for generating terminal GIFs from code. It allows you to write "tape files" that script terminal interactions and automatically record them as GIFs.

## Installing VHS

```bash
# Using Go
go install github.com/charmbracelet/vhs@latest

# Using Homebrew (macOS/Linux)
brew install vhs

# Using apt (Debian/Ubuntu)
sudo apt install vhs
```

VHS also requires [ttyd](https://github.com/tsl0922/ttyd) and [ffmpeg](https://ffmpeg.org/) to be installed.

## Pre-generated Demos

The demo GIF files are committed to the repository so they display in the main README on GitHub. You can view them directly in the [main README](../README.md#demos).

## Regenerating Demos

If you want to regenerate the demo GIFs (e.g., after modifying the examples or tape files):

```bash
# Generate individual demos
vhs demos/simple_demo.tape
vhs demos/dashboard_demo.tape
vhs demos/form_demo.tape
vhs demos/list_demo.tape

# Or use the provided script
./demos/generate_all.sh
```

## Available Demos

### Simple Demo (`simple_demo.tape`)
Demonstrates basic HawkTUI components including button, input field, and spinner. Shows focus management and component interaction.

**Recording:** `demos/simple_demo.gif`

### Dashboard Demo (`dashboard_demo.tape`)
Showcases the dashboard template with multiple widgets including metrics, status indicators, and charts using the Dracula theme.

**Recording:** `demos/dashboard_demo.gif`

### Form Demo (`form_demo.tape`)
Illustrates form validation, field navigation, and submit handling with the Nord theme.

**Recording:** `demos/form_demo.gif`

### List Demo (`list_demo.tape`)
Displays list component features including navigation, filtering, and item selection.

**Recording:** `demos/list_demo.gif`

## Customizing Demos

Each `.tape` file can be customized with different settings:

- `Set FontSize` - Adjust font size (default: 18-20)
- `Set Width` - Terminal width in pixels (default: 1400-1600)
- `Set Height` - Terminal height in pixels (default: 900-1000)
- `Set Theme` - Color theme (e.g., "Dracula", "Nord", "Catppuccin Mocha")
- `Set TypingSpeed` - How fast to "type" commands
- `Set PlaybackSpeed` - Overall playback speed

See the [VHS documentation](https://github.com/charmbracelet/vhs) for more options.
