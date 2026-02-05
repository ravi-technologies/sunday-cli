# Sunday CLI

Command-line interface for AI agents to access their inbox (email and SMS).

## Overview

Sunday CLI enables AI agents to receive and read communications on dedicated phone numbers and email addresses. This allows agents to:

- **Receive OTPs and verification codes** to authenticate with websites and services
- **Sign up for services** using the assigned phone number and email address
- **Read incoming messages** from services, notifications, and confirmations
- **Automate workflows** that require email/SMS verification

Each agent gets their own dedicated inbox with:
- A unique phone number for SMS
- A unique email address for email

## Use Cases

### Receiving OTPs for Website Login
```bash
# Check for recent SMS messages containing verification codes
sunday inbox sms --unread --json | jq '.[0].messages[].body'

# Get the latest email with OTP
sunday inbox email --unread
```

### Signing Up for Services
When filling out registration forms:
1. Use `sunday auth status --json` to get your assigned email
2. Use the phone number from your Sunday account for SMS verification
3. Monitor `sunday inbox` for the verification code
4. Complete the signup process

### Automated Verification Flows
```bash
# Poll for new messages in JSON format (ideal for automation)
sunday inbox --unread --json

# Filter for SMS only
sunday inbox --type sms --unread --json

# Filter for email only
sunday inbox --type email --unread --json
```

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
   sunday inbox list
   ```

3. **View only unread messages:**
   ```bash
   sunday inbox list --unread
   ```

4. **Get messages in JSON format (for automation):**
   ```bash
   sunday inbox list --json
   ```

## Commands

### Authentication

| Command | Description |
|---------|-------------|
| `sunday auth login` | Authenticate via browser OAuth flow |
| `sunday auth logout` | Clear stored credentials |
| `sunday auth status` | Show current authentication status and user email |

### Inbox

| Command | Description |
|---------|-------------|
| `sunday inbox list` | List all inbox messages (combined SMS + email) |
| `sunday inbox list --type email` | Filter by message type (email/sms) |
| `sunday inbox list --type sms` | Filter to SMS messages only |
| `sunday inbox list --direction incoming` | Filter by direction (incoming/outgoing) |
| `sunday inbox list --unread` | Show only unread messages |
| `sunday inbox email` | List email threads |
| `sunday inbox email <thread-id>` | View specific email thread with all messages |
| `sunday inbox sms` | List SMS conversations |
| `sunday inbox sms <conversation-id>` | View specific SMS conversation with all messages |

### Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output in JSON format (recommended for AI agents) |
| `--help` | Show help for any command |
| `--version` | Show version information |

## JSON Output for AI Agents

All commands support the `--json` flag, which outputs structured JSON ideal for programmatic parsing:

```bash
# List all unread messages as JSON
sunday inbox list --unread --json

# Parse with jq to extract OTP from SMS
sunday inbox sms --json | jq -r '.[0].messages[] | select(.body | test("[0-9]{6}")) | .body'

# Get the most recent email subject
sunday inbox email --json | jq -r '.[0].subject'
```

### JSON Response Structure

**Inbox List:**
```json
[
  {
    "type": "sms",
    "from": "+1234567890",
    "preview": "Your verification code is 123456",
    "date": "2024-01-15T10:30:00Z",
    "is_read": false
  }
]
```

**SMS Conversation Detail:**
```json
{
  "conversation_id": "conv_123",
  "from_number": "+1234567890",
  "sunday_number": "+0987654321",
  "messages": [
    {
      "direction": "incoming",
      "body": "Your verification code is 123456",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## Configuration

Credentials are stored in `~/.sunday/config.json` with secure file permissions (0600).

The config file contains:
- Access token (auto-refreshes when expired)
- Refresh token
- User email address

## Development

### Prerequisites

- Go 1.21+
- Make

### Building

```bash
# Build with API URL (required)
make build API_URL=https://api.sunday.example.com

# Build for all platforms
make build-all API_URL=https://api.sunday.example.com

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
