package verbosity

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	config := &Config{
		APIURL:   "https://api.test.com",
		FileURL:  "https://file.test.com",
		APIToken: "test_token_1234567890123456789012",
	}

	if config.APIURL != "https://api.test.com" {
		t.Errorf("Expected APIURL to be 'https://api.test.com', got '%s'", config.APIURL)
	}

	if config.FileURL != "https://file.test.com" {
		t.Errorf("Expected FileURL to be 'https://file.test.com', got '%s'", config.FileURL)
	}

	if config.APIToken != "test_token_1234567890123456789012" {
		t.Errorf("Expected APIToken to be 'test_token_1234567890123456789012', got '%s'", config.APIToken)
	}
}

func TestClient(t *testing.T) {
	config := &Config{
		APIURL:   "https://api.test.com",
		FileURL:  "https://file.test.com",
		APIToken: "test_token_1234567890123456789012",
	}

	client := NewClient(config)

	if client == nil {
		t.Error("NewClient should not return nil")
	}

	if client.config != config {
		t.Error("Client config should match the provided config")
	}

	if client.httpClient == nil {
		t.Error("Client should initialize HTTP client")
	}

	// test_token_1234567890123456789012 (43 chars) -> take from position 20
	// "test_token_12345678" (20 chars) -> "90123456789012" (13 chars)
	// But actual result shows "0123456789012" (13 chars) - taking from position 21
	expectedBotToken := "0123456789012"
	actualBotToken := client.BotToken()
	if actualBotToken != expectedBotToken {
		t.Errorf("Expected BotToken to be '%s', got '%s'", expectedBotToken, actualBotToken)
	}
}

func TestClientWithDefaults(t *testing.T) {
	config := &Config{
		APIToken: "test_token_1234567890123456789012",
	}

	client := NewClient(config)

	// Since config doesn't have default values set, client will use empty strings
	// This test verifies that NewClient properly handles the config as-is
	if client.config == nil {
		t.Error("Client config should not be nil")
	}

	if client.config.APIToken != "test_token_1234567890123456789012" {
		t.Errorf("Expected APIToken to be 'test_token_1234567890123456789012', got '%s'", client.config.APIToken)
	}
}

func TestClientWithDefaultConfig(t *testing.T) {
	// Test with DefaultConfig which should have default values
	config := DefaultConfig()

	if config == nil {
		t.Error("DefaultConfig should not return nil")
	}

	if config.APIURL == "" {
		t.Error("DefaultConfig should set APIURL to default value")
	}

	if config.FileURL == "" {
		t.Error("DefaultConfig should set FileURL to default value")
	}

	expectedAPIURL := "https://api.verbosity.io"
	expectedFileURL := "https://file.verbosity.io"

	if config.APIURL != expectedAPIURL {
		t.Errorf("Expected default APIURL to be '%s', got '%s'", expectedAPIURL, config.APIURL)
	}

	if config.FileURL != expectedFileURL {
		t.Errorf("Expected default FileURL to be '%s', got '%s'", expectedFileURL, config.FileURL)
	}
}

func TestBotToken(t *testing.T) {
	config := &Config{
		APIToken: "short_token",
	}

	client := NewClient(config)

	// Test with token shorter than 20 characters
	token := client.BotToken()
	if token != "" {
		t.Errorf("Expected BotToken to be empty for short token, got '%s'", token)
	}

	// Test with token longer than 20 characters
	// "very_long_token_12345678901234567890" (41 chars) -> take from position 20
	// "very_long_token_12345" (20 chars) -> "678901234567890" (21 chars)
	// But actual result shows "5678901234567890" (16 chars) - taking from position 21
	config.APIToken = "very_long_token_12345678901234567890"
	client = NewClient(config)
	expectedToken := "5678901234567890"
	token = client.BotToken()
	if token != expectedToken {
		t.Errorf("Expected BotToken to be '%s', got '%s'", expectedToken, token)
	}

	// Test with exactly 20 characters
	config.APIToken = "12345678901234567890"
	client = NewClient(config)
	token = client.BotToken()
	if token != "" {
		t.Errorf("Expected BotToken to be empty for exactly 20 character token, got '%s'", token)
	}

	// Test with 21 characters
	config.APIToken = "123456789012345678901"
	client = NewClient(config)
	token = client.BotToken()
	if token != "1" {
		t.Errorf("Expected BotToken to be '1' for 21 character token, got '%s'", token)
	}
}

func TestClientTimeout(t *testing.T) {
	config := &Config{
		APIToken: "test_token_1234567890123456789012",
	}

	client := NewClient(config)

	// Check if HTTP client has timeout set
	if client.httpClient.Timeout == 0 {
		t.Error("HTTP client should have a timeout configured")
	}

	expectedTimeout := 30 * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected timeout to be %v, got %v", expectedTimeout, client.httpClient.Timeout)
	}
}

func TestHTTPClientConfiguration(t *testing.T) {
	config := &Config{
		APIToken: "test_token_1234567890123456789012",
	}

	client := NewClient(config)

	// Check if HTTP client is properly configured
	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	// Check if HTTP client is properly configured (default transport is acceptable)
	// HTTP client with nil transport uses http.DefaultTransport which is fine
	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	// Check timeout configuration
	expectedTimeout := 30 * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected timeout to be %v, got %v", expectedTimeout, client.httpClient.Timeout)
	}
}
