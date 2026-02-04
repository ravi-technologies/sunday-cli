# Sunday CLI

Command-line interface for the Sunday 2FA backend service.

## Overview

Sunday CLI provides secure access to your Sunday 2FA account from the terminal. View your inbox messages, manage authentication, and integrate with scripts using JSON output.

## Installation

### From Source

```bash
git clone <repo-url>
cd sunday-cli
make build API_URL=https://api.sunday.example.com
```

### Pre-built Binaries

Download the latest release for your platform from the releases page.

## Quick Start

1. **Login to your account:**
   ```bash
   sunday auth login
   ```
   This opens your browser for OAuth authentication.

2. **Check your inbox:**
   ```bash
   sunday inbox
   ```

3. **View only unread messages:**
   ```bash
   sunday inbox --unread
   ```

## Commands

### Authentication

| Command | Description |
|---------|-------------|
| `sunday auth login` | Authenticate via browser OAuth flow |
| `sunday auth logout` | Clear stored credentials |
| `sunday auth status` | Show current authentication status |

### Inbox

| Command | Description |
|---------|-------------|
| `sunday inbox` | List all inbox messages |
| `sunday inbox --type email` | Filter by message type (email/sms) |
| `sunday inbox --direction inbound` | Filter by direction (inbound/outbound) |
| `sunday inbox --unread` | Show only unread messages |
| `sunday inbox email` | List email threads |
| `sunday inbox email <thread-id>` | View specific email thread |
| `sunday inbox sms` | List SMS conversations |
| `sunday inbox sms <conversation-id>` | View specific SMS conversation |

### Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output in JSON format for scripting |
| `--help` | Show help for any command |
| `--version` | Show version information |

## Configuration

Credentials are stored in `~/.sunday/config.json` with secure file permissions (0600).

## JSON Output

All commands support `--json` flag for script integration:

```bash
sunday inbox --json | jq '.[] | select(.is_read == false)'
```

## Development

### Prerequisites

- Go 1.21+
- Make

### Building

```bash
# Build with API URL
make build API_URL=https://api.sunday.example.com

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint
```

### Project Structure

```
sunday-cli/
├── cmd/sunday/         # Main entry point
├── internal/
│   ├── api/           # HTTP client and API types
│   ├── auth/          # OAuth device flow
│   ├── config/        # Credential storage
│   ├── output/        # Human/JSON formatters
│   └── version/       # Build-time version info
└── pkg/cli/           # Cobra command definitions
```

## License

[Add license information]
