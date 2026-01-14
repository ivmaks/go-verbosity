package verbosity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// newRequest creates a new HTTP request for the API.
func (c *Client) newRequest(method, path string, params url.Values, body io.Reader) (*http.Request, error) {
	requestURL := c.config.APIURL + path

	if params != nil && len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.Header.Set("X-APIToken", c.config.APIToken)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// newFileRequest creates a new HTTP request for file uploads.
func (c *Client) newFileRequest(path string, body *bytes.Buffer) (*http.Request, error) {
	requestURL := c.config.FileURL + path

	req, err := http.NewRequest(http.MethodPost, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.Header.Set("X-APIToken", c.config.APIToken)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// do executes the request and handles the response.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return c.handleError(resp.StatusCode, body)
	}

	// Try to parse the response
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			// Check if it's an API error in JSON format
			var errorResp ErrorResponse
			if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Code != "" {
				return fmt.Errorf("API error: %s - %s", errorResp.Code, errorResp.Message)
			}

			// Check for validation errors
			var validationResp ValidationErrorResponse
			if err := json.Unmarshal(body, &validationResp); err == nil && validationResp.TamtamResponseAPI {
				return fmt.Errorf("validation error: %s", validationResp.Error)
			}

			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// handleError handles error responses from the API.
func (c *Client) handleError(statusCode int, body []byte) error {
	// Try to parse as JSON error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Code != "" {
		return fmt.Errorf("API error (code=%s): %s (status=%d)", errorResp.Code, errorResp.Message, statusCode)
	}

	// Check for validation errors
	var validationResp ValidationErrorResponse
	if err := json.Unmarshal(body, &validationResp); err == nil && validationResp.TamtamResponseAPI {
		return fmt.Errorf("validation error: %s (status=%d)", validationResp.Error, statusCode)
	}

	// Return generic error
	return fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
}

// IsAccessDeniedError checks if the error is an access denied error.
func IsAccessDeniedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "access_deny")
}

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "validation error") || strings.Contains(err.Error(), "tamtam_response_api")
}

// IsNotFoundError checks if the error is a "not found" error.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not found")
}
