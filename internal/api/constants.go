package api

import "time"

const (
	// TokenExpiryBuffer is the time before actual expiry to trigger refresh.
	// Backend issues 5-minute tokens; we refresh at 4 minutes for safety.
	TokenExpiryBuffer = 4 * time.Minute
)

const (
	// API endpoint paths
	PathDeviceCode    = "/api/auth/device/"
	PathDeviceToken   = "/api/auth/device/token/"
	PathTokenRefresh  = "/api/auth/token/refresh/"
	PathInbox         = "/api/inbox/"
	PathEmailInbox    = "/api/email-inbox/"
	PathSMSInbox      = "/api/sms-inbox/"
	PathPhone         = "/api/phone/"
	PathEmail         = "/api/email/"
	PathMessages      = "/api/messages/"
	PathEmailMessages = "/api/email-messages/"
	PathEncryption    = "/api/encryption/"
	PathOwner         = "/api/me/"
	PathPasswords     = "/api/passwords/"
	PathIdentities    = "/api/identities/"
	PathBindIdentity  = "/api/auth/bind-identity/"
)
