# Contributing to Hawk TUI

Thank you for your interest in contributing to Hawk TUI! We welcome contributions from developers of all skill levels.

## 🚀 Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/yourusername/hawk.git
   cd hawk
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/awesome-feature
   ```
5. **Make your changes** and test them
6. **Commit and push**:
   ```bash
   git commit -m "Add awesome feature"
   git push origin feature/awesome-feature
   ```
7. **Create a Pull Request** on GitHub

## 🛠️ Development Environment

### Prerequisites
- Go 1.21 or higher
- Git
- Terminal that supports ANSI colors

### Building from Source
```bash
# Build the main TUI application
go build ./cmd/hawk

# Run tests
go test ./...

# Run with example data
go run examples/go/simple_demo.go | ./hawk
```

### Project Structure
```
hawk-tui/
├── cmd/hawk/              # Main TUI application
├── internal/
│   ├── tui/              # TUI components and models
│   ├── protocol/         # Communication protocol
│   └── config/           # Configuration management
├── pkg/
│   ├── types/            # Shared types and interfaces
│   └── client/           # Client libraries
├── examples/
│   ├── python/           # Python examples and client
│   ├── nodejs/           # Node.js examples and client
│   └── go/               # Go examples and client
└── docs/                 # Documentation
```

## 🧪 Testing

### Go Tests
```bash
# Run all Go tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/tui
go test ./pkg/types
```

### Python Client Tests
```bash
cd examples/python
python -m pytest test_hawk.py -v
```

### Integration Tests
```bash
# Test TUI with real application
go run examples/go/simple_demo.go | ./hawk

# Test Python integration
cd examples/python
python flask_demo.py | ../../hawk
```

## 📝 Code Style

### Go Code Style
- Follow standard Go formatting with `gofmt`
- Use `golint` and `go vet` for code quality
- Write meaningful commit messages
- Add comments for exported functions

### Python Code Style
- Follow PEP 8 style guide
- Use type hints where appropriate
- Add docstrings for public functions
- Keep line length under 88 characters

### Commit Messages
Use conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

Examples:
- `feat(tui): add metric chart visualization`
- `fix(protocol): handle malformed JSON messages`
- `docs(readme): update installation instructions`

## 🎯 Areas for Contribution

### 🔥 High Priority
- **Language Clients**: Node.js, Rust, Java client implementations
- **UI Components**: New widget types (tables, forms, trees)
- **Performance**: Optimization for high-frequency updates
- **Documentation**: Tutorials, examples, API documentation

### 🚀 Features
- **Themes**: Additional color schemes and styling options
- **Export**: Save dashboards and metrics to files
- **Plugins**: Extension system for custom widgets
- **Remote**: Web interface for remote monitoring

### 🐛 Bug Fixes
- **Terminal Compatibility**: Support for various terminal types
- **Memory Leaks**: Optimize long-running sessions
- **Protocol**: Edge cases in message handling
- **Error Handling**: Graceful degradation improvements

### 📖 Documentation
- **Tutorials**: Step-by-step integration guides
- **Examples**: Real-world application examples
- **API Docs**: Comprehensive API documentation
- **Videos**: Demo videos and screencasts

## 🔧 Adding New Features

### Adding a New UI Component
1. Create component in `internal/tui/components/`
2. Implement the `Component` interface
3. Add to the main TUI model
4. Write tests and examples
5. Update documentation

### Adding Language Support
1. Create client library in `examples/{language}/`
2. Implement the protocol specification
3. Add examples and tests
4. Update README and documentation
5. Create integration tests

### Adding Protocol Features
1. Update protocol spec in `docs/PROTOCOL.md`
2. Add types in `pkg/types/protocol.go`
3. Update protocol handler
4. Implement in TUI components
5. Update client libraries

## 🚦 Pull Request Process

### Before Submitting
- [ ] All tests pass
- [ ] Code follows style guidelines
- [ ] Documentation is updated
- [ ] Examples work correctly
- [ ] No breaking changes (or clearly documented)

### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Examples updated
```

### Review Process
1. **Automated Checks**: All CI checks must pass
2. **Code Review**: At least one maintainer review
3. **Testing**: Manual testing by reviewer
4. **Documentation**: Ensure docs are updated
5. **Merge**: Squash and merge with clean commit message

## 🏗️ Architecture Guidelines

### Protocol Design
- **Backward Compatibility**: Don't break existing clients
- **Performance**: Minimize message overhead
- **Extensibility**: Easy to add new message types
- **Error Handling**: Graceful degradation

### TUI Design
- **Responsiveness**: 60 FPS target for smooth updates
- **Memory Efficiency**: Bounded data structures
- **Accessibility**: Keyboard navigation and screen readers
- **Themes**: Consistent styling system

### Client Libraries
- **Simplicity**: One-line integration for basic use
- **Performance**: Minimal overhead when TUI not present
- **Thread Safety**: Safe for concurrent use
- **Language Idioms**: Follow language-specific patterns

## 📋 Issue Guidelines

### Bug Reports
Use the bug report template and include:
- **Environment**: OS, terminal, Go/Python version
- **Steps to Reproduce**: Minimal example
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Screenshots**: If applicable

### Feature Requests
Use the feature request template and include:
- **Problem**: What problem does this solve?
- **Solution**: Proposed implementation
- **Alternatives**: Other approaches considered
- **Examples**: How would it be used?

### Questions
For questions and discussions:
- **Search**: Check existing issues first
- **Context**: Provide relevant background
- **Examples**: Include code samples
- **Environment**: Relevant system details

## 🎖️ Recognition

Contributors will be:
- **Listed**: In the README contributors section
- **Credited**: In release notes for their contributions
- **Invited**: To join the maintainer team for significant contributions
- **Featured**: In blog posts and social media

## 📞 Getting Help

### Communication Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and general discussion
- **Discord**: Real-time chat (invite in README)
- **Email**: maintainers@hawk-tui.dev

### Mentorship
New contributors can:
- **Find Issues**: Look for "good first issue" labels
- **Ask Questions**: Don't hesitate to ask for help
- **Pair Program**: Arrange sessions with maintainers
- **Get Guidance**: Request code review and feedback

## 📜 Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold this code.

### Our Pledge
We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

---

**Thank you for contributing to Hawk TUI!** 🦅

Together, we're making command-line tools more beautiful and accessible for everyone.