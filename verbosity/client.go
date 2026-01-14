package verbosity

import (
	"net/http"
	"os"
	"strings"
	"time"
)

// Config holds the configuration for the Verbosity client.
type Config struct {
	// API URL for all endpoints (users, chats, bots, orgs)
	APIURL string
	// File upload URL
	FileURL string
	// Bot API token
	APIToken string
}

// DefaultConfig returns a Config with values from environment variables.
// Environment variables:
// - VERBOSITY_API_URL: API URL (default: https://api.verbosity.io)
// - VERBOSITY_FILE_URL: File upload URL (default: https://file.verbosity.io)
// - VERBOSITY_API_TOKEN: Bot API token
func DefaultConfig() *Config {
	return &Config{
		APIURL:   getEnv("VERBOSITY_API_URL", "https://api.verbosity.io"),
		FileURL:  getEnv("VERBOSITY_FILE_URL", "https://file.verbosity.io"),
		APIToken: os.Getenv("VERBOSITY_API_TOKEN"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.TrimRight(value, "/")
}

// Client is the main Verbosity API client.
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new Verbosity API client with the given configuration.
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientFromEnv creates a new Verbosity API client with configuration from environment variables.
func NewClientFromEnv() *Client {
	return NewClient(DefaultConfig())
}

// Config returns the client configuration.
func (c *Client) Config() *Config {
	return c.config
}

// APIURL returns the API URL.
func (c *Client) APIURL() string {
	return c.config.APIURL
}

// FileURL returns the file upload URL.
func (c *Client) FileURL() string {
	return c.config.FileURL
}

// APIToken returns the full API token.
func (c *Client) APIToken() string {
	return c.config.APIToken
}

// BotToken returns the API token for bot operations (without first 20 characters).
// According to the API documentation, the token for signature verification
// and API calls should not include the first 20 characters.
func (c *Client) BotToken() string {
	if len(c.config.APIToken) <= 20 {
		return ""
	}
	return c.config.APIToken[20:]
}
