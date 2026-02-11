# Sunday CLI - AI Agent Inbox Access

This skill teaches you how to use the Sunday CLI to access your dedicated email and SMS inbox.

## What is Sunday CLI?

Sunday CLI provides AI agents with their own phone number and email address. You can:
- Receive SMS messages (including OTPs and verification codes)
- Receive emails (including signup confirmations and verification links)
- Use these credentials to sign up for websites and services
- Authenticate to services that require phone/email verification
- Store and retrieve E2E-encrypted passwords per identity

## Prerequisites

Before using Sunday CLI, ensure you are authenticated:
```bash
sunday auth status
```

If not authenticated, run:
```bash
sunday auth login
```

## Common Workflows

### 1. Checking for OTP/Verification Codes

When you need to receive a verification code (e.g., after signing up for a service):

```bash
# Check for new SMS messages
sunday inbox sms --unread --json

# Check for new emails
sunday inbox email --unread --json

# Check unified inbox for all new messages
sunday inbox list --unread --json
```

**Extracting OTP from SMS:**
```bash
# Get all unread SMS and look for 6-digit codes
sunday inbox sms --unread --json | jq -r '.[] | .messages[] | select(.body | test("[0-9]{6}")) | .body'
```

### 2. Signing Up for Services

When filling out a registration form that requires email/phone:

1. **Get your Sunday email address:**
   ```bash
   sunday get email --json | jq -r '.email'
   ```

2. **Get your Sunday phone number:**
   ```bash
   sunday get phone --json | jq -r '.phone_number'
   ```

3. **Use these credentials in the registration form**

4. **After submitting the form, check for verification:**
   ```bash
   # Wait a few seconds, then check for the verification message
   sunday inbox list --unread --json
   ```

### 3. Two-Factor Authentication (2FA)

When logging into a website that sends 2FA codes:

```bash
# For SMS-based 2FA
sunday inbox sms --unread --json

# For email-based 2FA
sunday inbox email --unread --json
```

### 4. Viewing Message Details

**View a specific SMS conversation:**
```bash
# List conversations first
sunday inbox sms --json

# Then view specific conversation
sunday inbox sms <conversation_id> --json
```

**View a specific email thread:**
```bash
# List threads first
sunday inbox email --json

# Then view specific thread
sunday inbox email <thread_id> --json
```

## Command Reference

### Authentication Commands
| Command | Description |
|---------|-------------|
| `sunday auth login` | Authenticate (opens browser) |
| `sunday auth logout` | Clear credentials |
| `sunday auth status` | Show auth status and account email |
| `sunday auth status --json` | Get auth info as JSON |

### Resource Commands
| Command | Description |
|---------|-------------|
| `sunday get phone` | Get your assigned Sunday phone number |
| `sunday get email` | Get your assigned Sunday email address |
| `sunday get phone --json` | Get phone as JSON |
| `sunday get email --json` | Get email as JSON |

### Inbox Commands (grouped by conversation/thread)
| Command | Description |
|---------|-------------|
| `sunday inbox list` | List all messages (SMS + email) |
| `sunday inbox list --unread` | Only unread messages |
| `sunday inbox list --type sms` | Only SMS messages |
| `sunday inbox list --type email` | Only email messages |
| `sunday inbox sms` | List SMS conversations |
| `sunday inbox sms <id>` | View SMS conversation |
| `sunday inbox email` | List email threads |
| `sunday inbox email <id>` | View email thread |

### Message Commands (individual messages)
| Command | Description |
|---------|-------------|
| `sunday message sms` | List all SMS messages (flat) |
| `sunday message sms <id>` | View specific SMS message by ID |
| `sunday message sms --unread` | Only unread SMS messages |
| `sunday message email` | List all email messages (flat) |
| `sunday message email <id>` | View specific email message by ID |
| `sunday message email --unread` | Only unread email messages |

### Important Flags
- `--json` - Output as JSON (always use this for parsing)
- `--unread` - Filter to unread messages only

## Best Practices

1. **Always use `--json` flag** when you need to parse the output programmatically

2. **Poll for new messages** after triggering a verification:
   ```bash
   # Wait a moment, then check
   sleep 5 && sunday inbox list --unread --json
   ```

3. **Use specific filters** to reduce noise:
   ```bash
   # If expecting SMS OTP, filter to SMS only
   sunday inbox list --type sms --unread --json
   ```

4. **Check both SMS and email** - some services send to either:
   ```bash
   sunday inbox list --unread --json
   ```

## Example: Complete Signup Flow

```bash
# 1. Get your Sunday email and phone
EMAIL=$(sunday get email --json | jq -r '.email')
PHONE=$(sunday get phone --json | jq -r '.phone_number')
echo "Use this email for signup: $EMAIL"
echo "Use this phone for signup: $PHONE"

# 2. [Fill out the signup form with these credentials]

# 3. Wait for verification email/SMS
sleep 10

# 4. Check for verification
sunday inbox list --unread --json

# 5. Extract verification link or code from the email
sunday inbox email <thread_id> --json | jq -r '.messages[].text_content'

# Or extract OTP code from SMS
sunday inbox sms <conversation_id> --json | jq -r '.messages[].body'
```

### 5. Managing Passwords

Store credentials for services you've signed up for:

```bash
# After signing up for a service, store the credentials
sunday passwords create example.com --username "$EMAIL" --password 'the-password-used'

# Or auto-generate a password during signup
sunday passwords create example.com
# Outputs: Generated password: xK9#mL2...  (use this in the signup form)

# Retrieve stored credentials later
sunday passwords list --json
sunday passwords get <uuid> --json

# Update a password
sunday passwords edit <uuid> --password 'new-password'

# Generate a password without storing it
sunday passwords generate --length 24 --json | jq -r '.password'
```

**Note:** URL inputs are automatically cleaned to domains (e.g. `https://mail.google.com/inbox` → `google.com`). Username defaults to your identity email if not specified. Password is auto-generated if not provided.

### Password Commands
| Command | Description |
|---------|-------------|
| `sunday passwords list` | List all stored passwords |
| `sunday passwords get <uuid>` | Show a stored password (decrypted) |
| `sunday passwords create <domain>` | Create a new password entry |
| `sunday passwords edit <uuid>` | Edit a stored password entry |
| `sunday passwords delete <uuid>` | Delete a stored password entry |
| `sunday passwords generate` | Generate a random password |

**Create flags:** `--username`, `--password`, `--generate`, `--length` (default: 16), `--no-special`, `--no-digits`, `--exclude-chars`, `--notes`

## Troubleshooting

**Not authenticated:**
```bash
sunday auth login
```

**No messages appearing:**
- Verify the correct email/phone was used
- Wait a few more seconds for delivery
- Check spam filters on the service side

**Token expired:**
The CLI automatically refreshes tokens. If issues persist:
```bash
sunday auth logout
sunday auth login
```
