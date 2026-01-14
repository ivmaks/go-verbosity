package verbosity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
)

// BotInfo represents information about the current bot.
type BotInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	ReceiverURL string `json:"receiver_url"`
	SendingURL  string `json:"sending_url"`
	FullKey     string `json:"key"`
	URLKey      string `json:"url_key"`
	APIKey      string `json:"api_key"`
}

// BotRequest represents an incoming request to the bot.
type BotRequest struct {
	UserID         int64       `json:"user_id"`
	UserUniqueName *string     `json:"user_unique_name,omitempty"`
	PostNo         int64       `json:"post_no"`
	ChatID         int64       `json:"chat_id"`
	OrganizationID *int64      `json:"organization_id,omitempty"`
	Text           string      `json:"text"`
	TextParsed     []TextBlock `json:"text_parsed"`
	ReplyNo        *int64      `json:"reply_no,omitempty"`
	ReplyText      *string     `json:"reply_text,omitempty"`
	FileGUID       *string     `json:"file_guid,omitempty"`
	FileName       *string     `json:"file_name,omitempty"`
	Attachments    []string    `json:"attachments,omitempty"`
}

// TextBlock represents a parsed text block.
type TextBlock struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ActionRequest represents an action callback request.
type ActionRequest struct {
	UserID         int64             `json:"user_id"`
	ChatID         int64             `json:"chat_id"`
	PostNo         int64             `json:"post_no"`
	OrganizationID *int64            `json:"organization_id,omitempty"`
	Action         string            `json:"action"`
	Params         map[string]string `json:"params"`
}

// GetBotInfo retrieves information about the current bot.
// Note: This method uses the InfoBot API through the messaging system.
func (c *Client) GetBotInfo() (*User, error) {
	// Get bot info from the users endpoint
	// Since we can't directly get bot info, we use a workaround
	// The bot's user ID can be determined from its API key

	// For now, return nil with explanation
	// In a real implementation, you would need to query InfoBot directly
	return nil, fmt.Errorf("GetBotInfo requires special InfoBot API access. " +
		"Use direct messaging with @infobot for bot management operations")
}

// GetBotInfoFromAPI retrieves bot information from the core API if available.
func (c *Client) GetBotInfoFromAPI() (*User, error) {
	// Try to get bot info from the users endpoint
	// This is a placeholder that should be implemented based on actual API behavior

	// For now, return an error indicating this requires special handling
	return nil, fmt.Errorf("bot info retrieval requires InfoBot integration")
}

// VerifySignature verifies the X-Signature header for incoming bot requests.
//
// The signature is computed as:
//
//	base64(hmac_sha256(bkey, body))
//
// where bkey = to_bytes(int(key, 16))
func (c *Client) VerifySignature(body, signature string) (bool, error) {
	if len(c.config.APIToken) < 20 {
		return false, fmt.Errorf("API token is too short")
	}

	// Convert hex key to bytes
	bkey, ok := new(big.Int).SetString(c.config.APIToken[:20], 16)
	if !ok {
		return false, fmt.Errorf("failed to parse API token key")
	}

	h := hmac.New(sha256.New, bkey.Bytes())
	h.Write([]byte(body))
	expected := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature == expected, nil
}

// ParseBotRequest parses a JSON bot request.
func ParseBotRequest(data []byte) (*BotRequest, error) {
	var request BotRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("failed to parse bot request: %w", err)
	}
	return &request, nil
}

// ParseActionRequest parses a JSON action request.
func ParseActionRequest(data []byte) (*ActionRequest, error) {
	var request ActionRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("failed to parse action request: %w", err)
	}
	return &request, nil
}

// CreateActionURL creates a bot action URL for interactive buttons/links.
//
// Format: bot://{action}?title={title}[&{key}={value}]
func CreateActionURL(action, title string, params map[string]string) string {
	// URL-encode the action and title
	encodedAction := url.PathEscape(action)
	encodedTitle := url.PathEscape(title)

	result := fmt.Sprintf("bot://%s?title=%s", encodedAction, encodedTitle)
	for key, value := range params {
		encodedKey := url.PathEscape(key)
		encodedValue := url.PathEscape(value)
		result += fmt.Sprintf("&%s=%s", encodedKey, encodedValue)
	}
	return result
}

// IsUserMentioned checks if the bot is mentioned in the message text.
func (r *BotRequest) IsUserMentioned() bool {
	for _, block := range r.TextParsed {
		if block.Type == "mention" && block.Value == "bot" {
			return true
		}
	}
	return false
}

// IsCommand checks if the message starts with a command prefix.
func (r *BotRequest) IsCommand() bool {
	return len(r.Text) > 0 && r.Text[0] == '/'
}

// GetCommand returns the command and arguments from the message.
func (r *BotRequest) GetCommand() (string, []string) {
	if !r.IsCommand() {
		return "", nil
	}

	parts := splitCommand(r.Text)
	if len(parts) == 0 {
		return "", nil
	}

	return parts[0], parts[1:]
}

func splitCommand(text string) []string {
	var result []string
	var current string
	inWord := false

	for _, c := range text {
		if c == ' ' || c == '\n' || c == '\t' {
			if inWord {
				result = append(result, current)
				current = ""
				inWord = false
			}
		} else {
			current += string(c)
			inWord = true
		}
	}

	if inWord {
		result = append(result, current)
	}

	return result
}

// MessageIsEmpty checks if the message has no content.
func (r *BotRequest) MessageIsEmpty() bool {
	return r.Text == "" && r.FileGUID == nil && len(r.Attachments) == 0
}

// HasFile checks if the message contains a file.
func (r *BotRequest) HasFile() bool {
	return r.FileGUID != nil
}

// HasReply checks if the message is a reply to another message.
func (r *BotRequest) HasReply() bool {
	return r.ReplyNo != nil
}

// GetReplyText returns the text of the message being replied to.
func (r *BotRequest) GetReplyText() string {
	if r.ReplyText != nil {
		return *r.ReplyText
	}
	return ""
}

// GetFileGUID returns the file GUID if present.
func (r *BotRequest) GetFileGUID() string {
	if r.FileGUID != nil {
		return *r.FileGUID
	}
	return ""
}

// GetFileName returns the file name if present.
func (r *BotRequest) GetFileName() string {
	if r.FileName != nil {
		return *r.FileName
	}
	return ""
}
