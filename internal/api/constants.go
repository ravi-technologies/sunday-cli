package api

import "time"

const (
	// TokenExpiryBuffer is the time before actual expiry to trigger refresh.
	// Backend issues 5-minute tokens; we refresh at 4 minutes for safety.
	TokenExpiryBuffer = 4 * time.Minute
)

const (
	// API endpoint paths
	PathDeviceCode   = "/api/v1/auth/device/"
	PathDeviceToken  = "/api/v1/auth/device/token/"
	PathTokenRefresh = "/api/v1/auth/token/refresh/"
	PathInbox        = "/api/v1/inbox/"
	PathEmailInbox   = "/api/v1/email-inbox/"
	PathSMSInbox     = "/api/v1/sms-inbox/"
)
