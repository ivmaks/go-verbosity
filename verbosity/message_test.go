package verbosity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMessage(t *testing.T) {
	// Mock server to simulate API responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and URL
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		expectedURL := "/msg/post/123/456"
		if r.URL.Path != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, r.URL.Path)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content type application/json, got %s", contentType)
		}

		// Verify API token header
		token := r.Header.Get("X-APIToken")
		if token != "test_token_1234567890123456789012" {
			t.Errorf("Expected API token 'test_token_1234567890123456789012', got '%s'", token)
		}

		// Parse request body
		var req UpdateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify request data
		if req.Text != "Updated message" {
			t.Errorf("Expected text 'Updated message', got '%s'", req.Text)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UpdateMessageResponse{
			UUID:    "test-uuid-123",
			ChatID:  123,
			PostNo:  456,
			Version: intPtr(2),
		})
	}))
	defer server.Close()

	// Create client with mock server
	config := &Config{
		APIURL:   server.URL,
		FileURL:  server.URL,
		APIToken: "test_token_1234567890123456789012",
	}
	client := NewClient(config)

	// Test successful update
	updateReq := &UpdateMessageRequest{
		Text: "Updated message",
	}

	response, err := client.UpdateMessage(123, 456, updateReq)
	if err != nil {
		t.Errorf("UpdateMessage should not return error: %v", err)
	}

	if response == nil {
		t.Error("UpdateMessage should return response")
	}

	if response.UUID != "test-uuid-123" {
		t.Errorf("Expected UUID 'test-uuid-123', got '%s'", response.UUID)
	}

	if response.ChatID != 123 {
		t.Errorf("Expected ChatID 123, got %d", response.ChatID)
	}

	if response.PostNo != 456 {
		t.Errorf("Expected PostNo 456, got %d", response.PostNo)
	}

	if response.Version == nil || *response.Version != 2 {
		t.Errorf("Expected Version 2, got %v", response.Version)
	}
}

func TestUpdateMessageValidation(t *testing.T) {
	config := &Config{
		APIURL:   "https://api.test.com",
		FileURL:  "https://file.test.com",
		APIToken: "test_token_1234567890123456789012",
	}
	client := NewClient(config)

	// Test with zero chat ID
	_, err := client.UpdateMessage(0, 456, &UpdateMessageRequest{Text: "test"})
	if err == nil {
		t.Error("UpdateMessage should return error for zero chat ID")
	}

	// Test with zero post number
	_, err = client.UpdateMessage(123, 0, &UpdateMessageRequest{Text: "test"})
	if err == nil {
		t.Error("UpdateMessage should return error for zero post number")
	}

	// Test with nil request
	_, err = client.UpdateMessage(123, 456, nil)
	if err == nil {
		t.Error("UpdateMessage should return error for nil request")
	}

	// Test with empty text
	_, err = client.UpdateMessage(123, 456, &UpdateMessageRequest{Text: ""})
	if err == nil {
		t.Error("UpdateMessage should return error for empty text")
	}
}

func TestUpdateMessageWithAttachments(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UpdateMessageResponse{
			UUID:    "test-uuid-attachments",
			ChatID:  123,
			PostNo:  789,
			Version: intPtr(3),
		})
	}))
	defer server.Close()

	config := &Config{
		APIURL:   server.URL,
		FileURL:  server.URL,
		APIToken: "test_token_1234567890123456789012",
	}
	client := NewClient(config)

	// Test update with attachments
	attachments := []string{"guid1", "guid2", "guid3"}
	response, err := client.UpdateMessageWithAttachments(123, 789, "Message with attachments", attachments)
	if err != nil {
		t.Errorf("UpdateMessageWithAttachments should not return error: %v", err)
	}

	if response == nil {
		t.Error("UpdateMessageWithAttachments should return response")
	}

	if response.UUID != "test-uuid-attachments" {
		t.Errorf("Expected UUID 'test-uuid-attachments', got '%s'", response.UUID)
	}
}

func TestUpdateMessageWithReply(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UpdateMessageResponse{
			UUID:    "test-uuid-reply",
			ChatID:  123,
			PostNo:  999,
			Version: intPtr(1),
		})
	}))
	defer server.Close()

	config := &Config{
		APIURL:   server.URL,
		FileURL:  server.URL,
		APIToken: "test_token_1234567890123456789012",
	}
	client := NewClient(config)

	// Test update with reply
	response, err := client.UpdateMessageWithReply(123, 999, 100, "Reply message")
	if err != nil {
		t.Errorf("UpdateMessageWithReply should not return error: %v", err)
	}

	if response == nil {
		t.Error("UpdateMessageWithReply should return response")
	}

	if response.PostNo != 999 {
		t.Errorf("Expected PostNo 999, got %d", response.PostNo)
	}
}

func TestUpdateMessageE2E(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UpdateMessageResponse{
			UUID:    "test-uuid-e2e",
			ChatID:  123,
			PostNo:  111,
			Version: intPtr(4),
		})
	}))
	defer server.Close()

	config := &Config{
		APIURL:   server.URL,
		FileURL:  server.URL,
		APIToken: "test_token_1234567890123456789012",
	}
	client := NewClient(config)

	// Test update with E2E encryption
	response, err := client.UpdateMessageE2E(123, 111, "E2E message", true)
	if err != nil {
		t.Errorf("UpdateMessageE2E should not return error: %v", err)
	}

	if response == nil {
		t.Error("UpdateMessageE2E should return response")
	}

	if response.UUID != "test-uuid-e2e" {
		t.Errorf("Expected UUID 'test-uuid-e2e', got '%s'", response.UUID)
	}
}

func TestUpdateMessageRequestMarshaling(t *testing.T) {
	// Test complete request structure
	req := &UpdateMessageRequest{
		Text:        "Test message",
		E2E:         boolPtr(true),
		ReplyNo:     int64Ptr(123),
		Quote:       stringPtr("Original message"),
		Attachments: []string{"guid1", "guid2"},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal UpdateMessageRequest: %v", err)
	}

	var unmarshaled UpdateMessageRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal UpdateMessageRequest: %v", err)
	}

	if unmarshaled.Text != "Test message" {
		t.Errorf("Expected text 'Test message', got '%s'", unmarshaled.Text)
	}

	if unmarshaled.E2E == nil || !*unmarshaled.E2E {
		t.Error("Expected E2E to be true")
	}

	if unmarshaled.ReplyNo == nil || *unmarshaled.ReplyNo != 123 {
		t.Error("Expected ReplyNo to be 123")
	}

	if unmarshaled.Quote == nil || *unmarshaled.Quote != "Original message" {
		t.Error("Expected Quote to be 'Original message'")
	}

	if len(unmarshaled.Attachments) != 2 {
		t.Errorf("Expected 2 attachments, got %d", len(unmarshaled.Attachments))
	}

	if unmarshaled.Attachments[0] != "guid1" || unmarshaled.Attachments[1] != "guid2" {
		t.Error("Attachments don't match expected values")
	}
}

func TestUpdateMessageResponseMarshaling(t *testing.T) {
	// Test response structure
	resp := &UpdateMessageResponse{
		UUID:    "test-uuid",
		ChatID:  123,
		PostNo:  456,
		Version: intPtr(2),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("Failed to marshal UpdateMessageResponse: %v", err)
	}

	var unmarshaled UpdateMessageResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal UpdateMessageResponse: %v", err)
	}

	if unmarshaled.UUID != "test-uuid" {
		t.Errorf("Expected UUID 'test-uuid', got '%s'", unmarshaled.UUID)
	}

	if unmarshaled.ChatID != 123 {
		t.Errorf("Expected ChatID 123, got %d", unmarshaled.ChatID)
	}

	if unmarshaled.PostNo != 456 {
		t.Errorf("Expected PostNo 456, got %d", unmarshaled.PostNo)
	}

	if unmarshaled.Version == nil || *unmarshaled.Version != 2 {
		t.Errorf("Expected Version 2, got %v", unmarshaled.Version)
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
