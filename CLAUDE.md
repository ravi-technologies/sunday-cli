# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

Sunday CLI is a Go command-line client for the Sunday backend service. It provides programmatic access to inbox (emails and SMS) for AI agents.

## Commands

```bash
# Development
make build API_URL=https://api.sunday.app   # Build binary (API_URL required)
make test                                    # Run tests
make lint                                    # Check with golangci-lint
make lint-fix                                # Auto-fix lint issues
make clean                                   # Remove build artifacts

# Cross-compilation
make build-all API_URL=https://api.sunday.app  # Build for all platforms
```

## Architecture

```
cmd/sunday/           # Entry point
internal/
├── api/              # HTTP client and API types
├── auth/             # Device code flow orchestration
├── config/           # Token/config file management
├── output/           # Human/JSON formatters
└── version/          # Build-time version info
pkg/cli/              # Cobra commands
```

### Key Patterns

- **Output formatting**: All commands support `--json` flag for AI agent consumption
- **Token refresh**: API client automatically refreshes expired tokens
- **Build-time config**: API URL injected via ldflags (no runtime config needed)

## Code Style

- Use `gofmt` formatting
- Follow Go idioms and effective Go guidelines
- Error wrapping with `fmt.Errorf("context: %w", err)`
- Conventional commits: `feat(scope):`, `fix(scope):`, `refactor(scope):`
