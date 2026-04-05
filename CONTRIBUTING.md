# Contributing to TUI

Thank you for your interest in contributing to TUI! This document provides guidelines
and information for contributors.

## Reporting Bugs

- Check existing issues to avoid duplicates
- Include Go version, terminal emulator, and OS
- Provide a minimal reproducible example
- Include terminal screenshots if the issue is visual

## Suggesting Features

- Open an issue describing the feature and its use case
- For widget proposals, include a sketch of the expected API
- For KGP-related features, reference the relevant protocol specification

## Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run all checks (see below)
5. Submit a pull request with a clear description

## Development Guidelines

### Code Style

- Run `gofmt` on all code before committing
- Run `go vet ./...` to catch common issues
- Follow existing conventions in the codebase:
  - Short receiver names (1-2 characters): `(a *App)`, `(b *Buffer)`, `(s Style)`
  - Godoc comments on all exported types and functions
  - Godoc comments start with the name being documented

### Testing

- Add tests for new functionality
- Use table-driven tests where appropriate
- Run `go test ./...` before submitting
- Aim for high coverage on core packages

### Commit Messages

- Use present tense ("Add feature" not "Added feature")
- Start with a verb ("Fix bug" not "Bug fix")
- Reference issues where applicable ("Fix #42")
- Keep the first line under 72 characters

## Project Structure

```
tui/
├── *.go              # Core framework (app, buffer, events, layout, etc.)
├── image.go          # KGP image integration
├── animation.go      # KGP animation support
├── widget/           # Built-in widget library
│   ├── *.go          # Individual widgets
│   └── doc.go        # Package documentation
├── examples/         # Example applications
│   ├── hello/        # Simple getting started example
│   ├── demo/         # Comprehensive widget showcase
│   ├── image/        # KGP image rendering
│   └── animation/    # KGP animation demo
└── *_test.go         # Tests
```
