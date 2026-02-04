package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ravi-technologies/sunday-cli/internal/config"
)

// newTestClient creates a Client configured to use the test server.
// It sets up minimal config with a valid access token and future expiry.
func newTestClient(serverURL string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		baseURL:    serverURL,
		config: &config.Config{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(time.Hour), // Token won't expire during test
		},
	}
}

// TestListInbox_NoFilters verifies that ListInbox returns all messages without filters.
func TestListInbox_NoFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path
		if r.URL.Path != PathInbox {
			t.Errorf("Expected path %s, got %s", PathInbox, r.URL.Path)
		}

		// Verify no query params
		if r.URL.RawQuery != "" {
			t.Errorf("Expected no query params, got %s", r.URL.RawQuery)
		}

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			t.Errorf("Expected Bearer token in Authorization header, got %s", authHeader)
		}

		// Return mock response
		messages := []InboxMessage{
			{
				ID:          1,
				Type:        "email",
				FromAddress: "sender@example.com",
				ToAddress:   "user@sunday.app",
				Subject:     "Test Subject",
				Body:        "Test body content",
				Direction:   "inbound",
				IsRead:      false,
				CreatedDt:   time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			},
			{
				ID:          2,
				Type:        "sms",
				FromAddress: "+15551234567",
				ToAddress:   "+15559876543",
				Subject:     "",
				Body:        "Hello from SMS",
				Direction:   "inbound",
				IsRead:      true,
				CreatedDt:   time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	messages, err := client.ListInbox("", "", false)

	if err != nil {
		t.Fatalf("ListInbox() error = %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	// Verify first message
	if messages[0].ID != 1 {
		t.Errorf("First message ID = %d, want 1", messages[0].ID)
	}
	if messages[0].Type != "email" {
		t.Errorf("First message Type = %s, want email", messages[0].Type)
	}
	if messages[0].FromAddress != "sender@example.com" {
		t.Errorf("First message FromAddress = %s, want sender@example.com", messages[0].FromAddress)
	}

	// Verify second message
	if messages[1].ID != 2 {
		t.Errorf("Second message ID = %d, want 2", messages[1].ID)
	}
	if messages[1].Type != "sms" {
		t.Errorf("Second message Type = %s, want sms", messages[1].Type)
	}
}

// TestListInbox_TypeFilter verifies that ListInbox correctly filters by type (sms/email).
func TestListInbox_TypeFilter(t *testing.T) {
	testCases := []struct {
		name         string
		msgType      string
		expectedType string
	}{
		{
			name:         "filter by sms",
			msgType:      "sms",
			expectedType: "sms",
		},
		{
			name:         "filter by email",
			msgType:      "email",
			expectedType: "email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify type query parameter
				typeParam := r.URL.Query().Get("type")
				if typeParam != tc.expectedType {
					t.Errorf("Expected type=%s, got %s", tc.expectedType, typeParam)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]InboxMessage{})
			}))
			defer server.Close()

			client := newTestClient(server.URL)
			_, err := client.ListInbox(tc.msgType, "", false)

			if err != nil {
				t.Fatalf("ListInbox() error = %v", err)
			}
		})
	}
}

// TestListInbox_DirectionFilter verifies that ListInbox correctly filters by direction.
func TestListInbox_DirectionFilter(t *testing.T) {
	testCases := []struct {
		name              string
		direction         string
		expectedDirection string
	}{
		{
			name:              "filter by inbound",
			direction:         "inbound",
			expectedDirection: "inbound",
		},
		{
			name:              "filter by outbound",
			direction:         "outbound",
			expectedDirection: "outbound",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify direction query parameter
				directionParam := r.URL.Query().Get("direction")
				if directionParam != tc.expectedDirection {
					t.Errorf("Expected direction=%s, got %s", tc.expectedDirection, directionParam)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]InboxMessage{})
			}))
			defer server.Close()

			client := newTestClient(server.URL)
			_, err := client.ListInbox("", tc.direction, false)

			if err != nil {
				t.Fatalf("ListInbox() error = %v", err)
			}
		})
	}
}

// TestListInbox_UnreadFilter verifies that ListInbox correctly filters by unread status.
func TestListInbox_UnreadFilter(t *testing.T) {
	testCases := []struct {
		name           string
		unreadOnly     bool
		expectIsRead   bool
		expectNoParam  bool
	}{
		{
			name:          "unread only filter",
			unreadOnly:    true,
			expectIsRead:  false,
			expectNoParam: false,
		},
		{
			name:          "no unread filter",
			unreadOnly:    false,
			expectNoParam: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isReadParam := r.URL.Query().Get("is_read")

				if tc.expectNoParam {
					if isReadParam != "" {
						t.Errorf("Expected no is_read param, got %s", isReadParam)
					}
				} else {
					if isReadParam != "false" {
						t.Errorf("Expected is_read=false, got %s", isReadParam)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]InboxMessage{})
			}))
			defer server.Close()

			client := newTestClient(server.URL)
			_, err := client.ListInbox("", "", tc.unreadOnly)

			if err != nil {
				t.Fatalf("ListInbox() error = %v", err)
			}
		})
	}
}

// TestListInbox_MultipleFilters verifies that ListInbox correctly combines multiple filters.
func TestListInbox_MultipleFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Verify all three parameters are present
		typeParam := query.Get("type")
		if typeParam != "email" {
			t.Errorf("Expected type=email, got %s", typeParam)
		}

		directionParam := query.Get("direction")
		if directionParam != "inbound" {
			t.Errorf("Expected direction=inbound, got %s", directionParam)
		}

		isReadParam := query.Get("is_read")
		if isReadParam != "false" {
			t.Errorf("Expected is_read=false, got %s", isReadParam)
		}

		// Return filtered messages
		messages := []InboxMessage{
			{
				ID:          1,
				Type:        "email",
				FromAddress: "sender@example.com",
				ToAddress:   "user@sunday.app",
				Subject:     "Unread inbound email",
				Body:        "This matches all filters",
				Direction:   "inbound",
				IsRead:      false,
				CreatedDt:   time.Now(),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	messages, err := client.ListInbox("email", "inbound", true)

	if err != nil {
		t.Fatalf("ListInbox() error = %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(messages))
	}

	if messages[0].Type != "email" {
		t.Errorf("Message Type = %s, want email", messages[0].Type)
	}
	if messages[0].Direction != "inbound" {
		t.Errorf("Message Direction = %s, want inbound", messages[0].Direction)
	}
	if messages[0].IsRead != false {
		t.Errorf("Message IsRead = %v, want false", messages[0].IsRead)
	}
}

// TestListEmailThreads_Success verifies that ListEmailThreads returns email threads.
func TestListEmailThreads_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path
		if r.URL.Path != PathEmailInbox {
			t.Errorf("Expected path %s, got %s", PathEmailInbox, r.URL.Path)
		}

		// Verify no has_unread param when not filtering
		if r.URL.Query().Get("has_unread") != "" {
			t.Errorf("Expected no has_unread param, got %s", r.URL.Query().Get("has_unread"))
		}

		threads := []EmailThread{
			{
				ThreadID:        "<thread-1@example.com>",
				Subject:         "First thread subject",
				Preview:         "Preview of first email...",
				FromEmail:       "alice@example.com",
				SundayEmail:     "user@sunday.app",
				MessageCount:    3,
				UnreadCount:     1,
				LatestMessageDt: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
				OldestMessageDt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
			},
			{
				ThreadID:        "<thread-2@example.com>",
				Subject:         "Second thread subject",
				Preview:         "Preview of second email...",
				FromEmail:       "bob@example.com",
				SundayEmail:     "user@sunday.app",
				MessageCount:    1,
				UnreadCount:     0,
				LatestMessageDt: time.Date(2024, 1, 14, 10, 0, 0, 0, time.UTC),
				OldestMessageDt: time.Date(2024, 1, 14, 10, 0, 0, 0, time.UTC),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(threads)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	threads, err := client.ListEmailThreads(false)

	if err != nil {
		t.Fatalf("ListEmailThreads() error = %v", err)
	}

	if len(threads) != 2 {
		t.Fatalf("Expected 2 threads, got %d", len(threads))
	}

	// Verify first thread
	if threads[0].ThreadID != "<thread-1@example.com>" {
		t.Errorf("First thread ID = %s, want <thread-1@example.com>", threads[0].ThreadID)
	}
	if threads[0].Subject != "First thread subject" {
		t.Errorf("First thread Subject = %s, want 'First thread subject'", threads[0].Subject)
	}
	if threads[0].MessageCount != 3 {
		t.Errorf("First thread MessageCount = %d, want 3", threads[0].MessageCount)
	}
	if threads[0].UnreadCount != 1 {
		t.Errorf("First thread UnreadCount = %d, want 1", threads[0].UnreadCount)
	}

	// Verify second thread
	if threads[1].ThreadID != "<thread-2@example.com>" {
		t.Errorf("Second thread ID = %s, want <thread-2@example.com>", threads[1].ThreadID)
	}
	if threads[1].UnreadCount != 0 {
		t.Errorf("Second thread UnreadCount = %d, want 0", threads[1].UnreadCount)
	}
}

// TestListEmailThreads_UnreadOnly verifies that ListEmailThreads filters unread threads.
func TestListEmailThreads_UnreadOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify has_unread param is set to true
		hasUnreadParam := r.URL.Query().Get("has_unread")
		if hasUnreadParam != "true" {
			t.Errorf("Expected has_unread=true, got %s", hasUnreadParam)
		}

		threads := []EmailThread{
			{
				ThreadID:     "<unread-thread@example.com>",
				Subject:      "Unread thread",
				Preview:      "This thread has unread messages",
				FromEmail:    "sender@example.com",
				SundayEmail:  "user@sunday.app",
				MessageCount: 2,
				UnreadCount:  2,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(threads)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	threads, err := client.ListEmailThreads(true)

	if err != nil {
		t.Fatalf("ListEmailThreads() error = %v", err)
	}

	if len(threads) != 1 {
		t.Fatalf("Expected 1 thread, got %d", len(threads))
	}

	if threads[0].UnreadCount == 0 {
		t.Error("Expected thread with unread messages")
	}
}

// TestGetEmailThread_Success verifies that GetEmailThread returns thread detail.
func TestGetEmailThread_Success(t *testing.T) {
	threadID := "<test-thread-123@example.com>"
	encodedID := url.PathEscape(threadID)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path contains encoded thread ID using EscapedPath()
		expectedPath := PathEmailInbox + encodedID + "/"
		if r.URL.EscapedPath() != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.EscapedPath())
		}

		threadDetail := EmailThreadDetail{
			ThreadID:     threadID,
			Subject:      "Test Thread Subject",
			MessageCount: 2,
			Messages: []EmailMessage{
				{
					ID:          1,
					FromEmail:   "alice@example.com",
					ToEmail:     "user@sunday.app",
					CC:          "",
					Subject:     "Test Thread Subject",
					TextContent: "First message in thread",
					HTMLContent: "<p>First message in thread</p>",
					Direction:   "inbound",
					IsRead:      true,
					CreatedDt:   time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
				},
				{
					ID:          2,
					FromEmail:   "user@sunday.app",
					ToEmail:     "alice@example.com",
					CC:          "",
					Subject:     "Re: Test Thread Subject",
					TextContent: "Reply to first message",
					HTMLContent: "<p>Reply to first message</p>",
					Direction:   "outbound",
					IsRead:      true,
					CreatedDt:   time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(threadDetail)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	detail, err := client.GetEmailThread(threadID)

	if err != nil {
		t.Fatalf("GetEmailThread() error = %v", err)
	}

	if detail == nil {
		t.Fatal("GetEmailThread() returned nil")
	}

	if detail.ThreadID != threadID {
		t.Errorf("ThreadID = %s, want %s", detail.ThreadID, threadID)
	}
	if detail.Subject != "Test Thread Subject" {
		t.Errorf("Subject = %s, want 'Test Thread Subject'", detail.Subject)
	}
	if detail.MessageCount != 2 {
		t.Errorf("MessageCount = %d, want 2", detail.MessageCount)
	}
	if len(detail.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(detail.Messages))
	}

	// Verify first message
	if detail.Messages[0].Direction != "inbound" {
		t.Errorf("First message Direction = %s, want inbound", detail.Messages[0].Direction)
	}

	// Verify second message
	if detail.Messages[1].Direction != "outbound" {
		t.Errorf("Second message Direction = %s, want outbound", detail.Messages[1].Direction)
	}
}

// TestGetEmailThread_NotFound verifies that GetEmailThread handles 404 error.
func TestGetEmailThread_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIError{Detail: "Thread not found"})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetEmailThread("<nonexistent-thread@example.com>")

	if err == nil {
		t.Fatal("GetEmailThread() expected error for 404, got nil")
	}

	if !strings.Contains(err.Error(), "Thread not found") {
		t.Errorf("Error should contain 'Thread not found', got: %v", err)
	}
}

// TestListSMSConversations_Success verifies that ListSMSConversations returns SMS conversations.
func TestListSMSConversations_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path
		if r.URL.Path != PathSMSInbox {
			t.Errorf("Expected path %s, got %s", PathSMSInbox, r.URL.Path)
		}

		// Verify no has_unread param when not filtering
		if r.URL.Query().Get("has_unread") != "" {
			t.Errorf("Expected no has_unread param, got %s", r.URL.Query().Get("has_unread"))
		}

		conversations := []SMSConversation{
			{
				ConversationID:    "+15551234567:+15559876543",
				FromNumber:        "+15551234567",
				SundayPhone:       "My Phone",
				SundayPhoneNumber: "+15559876543",
				Preview:           "Latest message preview...",
				MessageCount:      5,
				UnreadCount:       2,
				LatestMessageDt:   time.Date(2024, 1, 15, 16, 0, 0, 0, time.UTC),
			},
			{
				ConversationID:    "+15552223333:+15559876543",
				FromNumber:        "+15552223333",
				SundayPhone:       "My Phone",
				SundayPhoneNumber: "+15559876543",
				Preview:           "Another conversation...",
				MessageCount:      10,
				UnreadCount:       0,
				LatestMessageDt:   time.Date(2024, 1, 14, 12, 0, 0, 0, time.UTC),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversations)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	conversations, err := client.ListSMSConversations(false)

	if err != nil {
		t.Fatalf("ListSMSConversations() error = %v", err)
	}

	if len(conversations) != 2 {
		t.Fatalf("Expected 2 conversations, got %d", len(conversations))
	}

	// Verify first conversation
	if conversations[0].ConversationID != "+15551234567:+15559876543" {
		t.Errorf("First conversation ID = %s, want '+15551234567:+15559876543'", conversations[0].ConversationID)
	}
	if conversations[0].FromNumber != "+15551234567" {
		t.Errorf("First conversation FromNumber = %s, want '+15551234567'", conversations[0].FromNumber)
	}
	if conversations[0].MessageCount != 5 {
		t.Errorf("First conversation MessageCount = %d, want 5", conversations[0].MessageCount)
	}
	if conversations[0].UnreadCount != 2 {
		t.Errorf("First conversation UnreadCount = %d, want 2", conversations[0].UnreadCount)
	}

	// Verify second conversation
	if conversations[1].UnreadCount != 0 {
		t.Errorf("Second conversation UnreadCount = %d, want 0", conversations[1].UnreadCount)
	}
}

// TestListSMSConversations_UnreadOnly verifies that ListSMSConversations filters unread conversations.
func TestListSMSConversations_UnreadOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify has_unread param is set to true
		hasUnreadParam := r.URL.Query().Get("has_unread")
		if hasUnreadParam != "true" {
			t.Errorf("Expected has_unread=true, got %s", hasUnreadParam)
		}

		conversations := []SMSConversation{
			{
				ConversationID:    "+15551234567:+15559876543",
				FromNumber:        "+15551234567",
				SundayPhone:       "My Phone",
				SundayPhoneNumber: "+15559876543",
				Preview:           "Unread message...",
				MessageCount:      3,
				UnreadCount:       1,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversations)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	conversations, err := client.ListSMSConversations(true)

	if err != nil {
		t.Fatalf("ListSMSConversations() error = %v", err)
	}

	if len(conversations) != 1 {
		t.Fatalf("Expected 1 conversation, got %d", len(conversations))
	}

	if conversations[0].UnreadCount == 0 {
		t.Error("Expected conversation with unread messages")
	}
}

// TestGetSMSConversation_Success verifies that GetSMSConversation returns conversation detail.
func TestGetSMSConversation_Success(t *testing.T) {
	conversationID := "+15551234567:+15559876543"
	encodedID := url.PathEscape(conversationID)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path contains encoded conversation ID
		expectedPath := PathSMSInbox + encodedID + "/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		conversationDetail := SMSConversationDetail{
			ConversationID: conversationID,
			FromNumber:     "+15551234567",
			SundayPhone:    "My Phone",
			MessageCount:   3,
			Messages: []SMSMessage{
				{
					ID:        1,
					Body:      "Hello!",
					Direction: "inbound",
					IsRead:    true,
					CreatedDt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					Body:      "Hi there!",
					Direction: "outbound",
					IsRead:    true,
					CreatedDt: time.Date(2024, 1, 10, 9, 5, 0, 0, time.UTC),
				},
				{
					ID:        3,
					Body:      "How are you?",
					Direction: "inbound",
					IsRead:    false,
					CreatedDt: time.Date(2024, 1, 10, 9, 10, 0, 0, time.UTC),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversationDetail)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	detail, err := client.GetSMSConversation(conversationID)

	if err != nil {
		t.Fatalf("GetSMSConversation() error = %v", err)
	}

	if detail == nil {
		t.Fatal("GetSMSConversation() returned nil")
	}

	if detail.ConversationID != conversationID {
		t.Errorf("ConversationID = %s, want %s", detail.ConversationID, conversationID)
	}
	if detail.FromNumber != "+15551234567" {
		t.Errorf("FromNumber = %s, want '+15551234567'", detail.FromNumber)
	}
	if detail.MessageCount != 3 {
		t.Errorf("MessageCount = %d, want 3", detail.MessageCount)
	}
	if len(detail.Messages) != 3 {
		t.Fatalf("Expected 3 messages, got %d", len(detail.Messages))
	}

	// Verify first message
	if detail.Messages[0].Body != "Hello!" {
		t.Errorf("First message Body = %s, want 'Hello!'", detail.Messages[0].Body)
	}
	if detail.Messages[0].Direction != "inbound" {
		t.Errorf("First message Direction = %s, want inbound", detail.Messages[0].Direction)
	}

	// Verify last message is unread
	if detail.Messages[2].IsRead != false {
		t.Errorf("Last message IsRead = %v, want false", detail.Messages[2].IsRead)
	}
}

// TestGetSMSConversation_NotFound verifies that GetSMSConversation handles 404 error.
func TestGetSMSConversation_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIError{Detail: "Conversation not found"})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetSMSConversation("+15550000000:+15559999999")

	if err == nil {
		t.Fatal("GetSMSConversation() expected error for 404, got nil")
	}

	if !strings.Contains(err.Error(), "Conversation not found") {
		t.Errorf("Error should contain 'Conversation not found', got: %v", err)
	}
}

// TestGetEmailThread_URLEncoding verifies that special characters in thread ID are properly URL encoded.
func TestGetEmailThread_URLEncoding(t *testing.T) {
	testCases := []struct {
		name      string
		threadID  string
		encodedID string
	}{
		{
			name:      "thread ID with angle brackets",
			threadID:  "<CABx+y@mail.example.com>",
			encodedID: url.PathEscape("<CABx+y@mail.example.com>"),
		},
		{
			name:      "thread ID with spaces",
			threadID:  "thread with spaces@example.com",
			encodedID: url.PathEscape("thread with spaces@example.com"),
		},
		{
			name:      "simple thread ID",
			threadID:  "simple-thread-123",
			encodedID: "simple-thread-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := PathEmailInbox + tc.encodedID + "/"
				// Use EscapedPath() to check the encoded URL path
				if r.URL.EscapedPath() != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.EscapedPath())
				}

				threadDetail := EmailThreadDetail{
					ThreadID:     tc.threadID,
					Subject:      "Test",
					MessageCount: 0,
					Messages:     []EmailMessage{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(threadDetail)
			}))
			defer server.Close()

			client := newTestClient(server.URL)
			_, err := client.GetEmailThread(tc.threadID)

			if err != nil {
				t.Fatalf("GetEmailThread() error = %v", err)
			}
		})
	}
}

// TestGetSMSConversation_URLEncoding verifies that phone numbers with + are properly URL encoded.
func TestGetSMSConversation_URLEncoding(t *testing.T) {
	testCases := []struct {
		name           string
		conversationID string
		encodedID      string
	}{
		{
			name:           "conversation ID with plus signs",
			conversationID: "+15551234567:+15559876543",
			encodedID:      url.PathEscape("+15551234567:+15559876543"),
		},
		{
			name:           "conversation ID without plus signs",
			conversationID: "15551234567:15559876543",
			encodedID:      "15551234567:15559876543",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := PathSMSInbox + tc.encodedID + "/"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				conversationDetail := SMSConversationDetail{
					ConversationID: tc.conversationID,
					FromNumber:     "+15551234567",
					SundayPhone:    "My Phone",
					MessageCount:   0,
					Messages:       []SMSMessage{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(conversationDetail)
			}))
			defer server.Close()

			client := newTestClient(server.URL)
			_, err := client.GetSMSConversation(tc.conversationID)

			if err != nil {
				t.Fatalf("GetSMSConversation() error = %v", err)
			}
		})
	}
}

// TestListInbox_EmptyResponse verifies handling of empty inbox.
func TestListInbox_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]InboxMessage{})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	messages, err := client.ListInbox("", "", false)

	if err != nil {
		t.Fatalf("ListInbox() error = %v", err)
	}

	if messages == nil {
		t.Fatal("ListInbox() returned nil, expected empty slice")
	}

	if len(messages) != 0 {
		t.Errorf("Expected 0 messages, got %d", len(messages))
	}
}

// TestListEmailThreads_EmptyResponse verifies handling of empty email threads list.
func TestListEmailThreads_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]EmailThread{})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	threads, err := client.ListEmailThreads(false)

	if err != nil {
		t.Fatalf("ListEmailThreads() error = %v", err)
	}

	if threads == nil {
		t.Fatal("ListEmailThreads() returned nil, expected empty slice")
	}

	if len(threads) != 0 {
		t.Errorf("Expected 0 threads, got %d", len(threads))
	}
}

// TestListSMSConversations_EmptyResponse verifies handling of empty SMS conversations list.
func TestListSMSConversations_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SMSConversation{})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	conversations, err := client.ListSMSConversations(false)

	if err != nil {
		t.Fatalf("ListSMSConversations() error = %v", err)
	}

	if conversations == nil {
		t.Fatal("ListSMSConversations() returned nil, expected empty slice")
	}

	if len(conversations) != 0 {
		t.Errorf("Expected 0 conversations, got %d", len(conversations))
	}
}
