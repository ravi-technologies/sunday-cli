# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

Sunday CLI is a Go command-line client for the Sunday backend service. It provides programmatic access to inbox (emails and SMS) for AI agents.

**Purpose:** AI agents get their own dedicated phone number and email address through Sunday. They can use these to:
- Sign up for websites and services (use the email/phone in registration forms)
- Receive OTPs and verification codes
- Complete 2FA authentication flows
- Read confirmation emails and SMS messages

## Using Sunday CLI as an AI Agent

### Quick Reference

```bash
# Get your assigned Sunday email and phone
sunday get email --json
sunday get phone --json

# Check authentication status
sunday auth status --json

# Check for new messages (always use --json for parsing)
sunday inbox list --unread --json

# Check for SMS messages (e.g., OTP codes)
sunday inbox sms --unread --json

# Check for emails (e.g., verification links)
sunday inbox email --unread --json

# View specific SMS conversation
sunday inbox sms <conversation_id> --json

# View specific email thread
sunday inbox email <thread_id> --json

# Access individual messages (flat list, not grouped)
sunday message sms --json              # List all SMS messages
sunday message sms <message_id> --json # Get specific SMS by ID
sunday message email --json            # List all email messages
sunday message email <message_id> --json # Get specific email by ID
```

### Workflow: Signing Up for a Service

1. Get your Sunday email: `sunday get email --json | jq -r '.email'`
2. Get your Sunday phone: `sunday get phone --json | jq -r '.phone_number'`
3. Fill out the signup form using these credentials
4. Wait for verification: `sleep 5 && sunday inbox list --unread --json`
5. Extract OTP or verification link from the message
6. Complete the verification

### Workflow: Receiving 2FA Codes

```bash
# After triggering 2FA, wait and check inbox
sleep 5
sunday inbox list --type sms --unread --json  # For SMS-based 2FA
sunday inbox list --type email --unread --json  # For email-based 2FA
```

See `.claude/skills/sunday-cli.md` for detailed usage instructions.

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
