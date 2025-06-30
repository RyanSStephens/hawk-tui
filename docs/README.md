# Hawk TUI Documentation

Welcome to the comprehensive documentation for Hawk TUI - the universal TUI framework for any programming language.

## Table of Contents

- [Quick Start Guide](quickstart.md) - Get up and running in 5 minutes
- [Installation Guide](installation.md) - All installation methods and options
- [API Reference](api-reference.md) - Complete API documentation
- [Language Clients](clients/) - Language-specific client libraries
  - [Python Client](clients/python.md)
  - [Node.js Client](clients/nodejs.md)
  - [Go Client](clients/go.md)
- [Examples](examples.md) - Real-world usage examples
- [Configuration](configuration.md) - Configuration options and customization
- [Architecture](architecture.md) - How Hawk TUI works internally
- [Contributing](../CONTRIBUTING.md) - How to contribute to the project
- [Troubleshooting](troubleshooting.md) - Common issues and solutions

## Overview

Hawk TUI transforms any command-line application into a beautiful, interactive TUI with zero architectural changes. Simply pipe your application's output through Hawk TUI and get instant visualization of logs, metrics, configuration, and more.

## Core Concepts

### Universal Protocol
Hawk TUI uses a simple JSON-RPC protocol over stdin/stdout that works with any programming language. Your application sends structured data, and Hawk TUI renders it beautifully.

### Zero Dependencies
Host applications don't need to install heavy TUI libraries. The thin client libraries handle all communication with the Hawk TUI binary.

### Drop-in Integration
Add monitoring to existing applications with minimal code changes. Most integrations require just one line of code.

## Getting Help

- **Issues**: Report bugs and request features on [GitHub Issues](https://github.com/hawk-tui/hawk-tui/issues)
- **Discussions**: Ask questions and share ideas on [GitHub Discussions](https://github.com/hawk-tui/hawk-tui/discussions)
- **Commercial Support**: For enterprise support, contact: support@hawktui.dev

## License

Hawk TUI is dual-licensed under AGPL-3.0 and commercial licenses. See the [LICENSE](../LICENSE) file for details.