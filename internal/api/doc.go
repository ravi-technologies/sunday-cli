// Package api provides the HTTP client and type definitions for communicating
// with the Sunday backend API.
//
// The package includes:
//   - Client: HTTP client with automatic token refresh and authentication
//   - Auth functions: Device code flow for OAuth 2.0 authentication
//   - Inbox functions: Fetching messages, emails, and SMS conversations
//   - Type definitions: Request/response structures for all API endpoints
//
// Example usage:
//
//	client, err := api.NewClient(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	messages, err := client.ListInbox(api.InboxFilters{})
package api
