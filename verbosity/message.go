package verbosity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendMessageRequest represents a request to send a message to a chat.
type SendMessageRequest struct {
	Key        string        `json:"key"`
	ChatID     int64         `json:"chat_id"`
	Text       string        `json:"text"`
	TextParsed []interface{} `json:"text_parsed,omitempty"`
	ReplyNo    *int64        `json:"reply_no,omitempty"`
}

// SendMessage sends a message to a non-private chat.
//
// API: POST /bot/message
func (c *Client) SendMessage(chatID int64, text string, replyNo *int64) (*MessageResponse, error) {
	if chatID == 0 {
		return nil, fmt.Errorf("chat_id cannot be zero")
	}
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	reqBody := SendMessageRequest{
		Key:     c.BotToken(),
		ChatID:  chatID,
		Text:    text,
		ReplyNo: replyNo,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, "/bot/message", nil, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response MessageResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendReply sends a reply to a specific message.
func (c *Client) SendReply(chatID, postNo int64, text string) (*MessageResponse, error) {
	return c.SendMessage(chatID, text, &postNo)
}

// PrivateMessageRequest represents a request to send a private message.
type PrivateMessageRequest struct {
	Text           string `json:"text"`
	UserID         *int64 `json:"user_id,omitempty"`
	UserEmail      string `json:"user_email,omitempty"`
	UserUniqueName string `json:"user_unique_name,omitempty"`
	ReplyNo        *int64 `json:"reply_no,omitempty"`
}

// SendPrivateMessageByID sends a private message to a user by their ID.
//
// API: POST /msg/post/private
func (c *Client) SendPrivateMessageByID(userID int64, text string, replyNo *int64) (*PrivateMessageResponse, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user_id cannot be zero")
	}
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	reqBody := PrivateMessageRequest{
		Text:    text,
		UserID:  &userID,
		ReplyNo: replyNo,
	}

	return c.sendPrivateMessage(reqBody)
}

// SendPrivateMessageByEmail sends a private message to a user by their email.
//
// API: POST /msg/post/private
func (c *Client) SendPrivateMessageByEmail(email, text string, replyNo *int64) (*PrivateMessageResponse, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	reqBody := PrivateMessageRequest{
		Text:      text,
		UserEmail: email,
		ReplyNo:   replyNo,
	}

	return c.sendPrivateMessage(reqBody)
}

// SendPrivateMessageByUniqueName sends a private message to a user by their unique name.
//
// API: POST /msg/post/private
func (c *Client) SendPrivateMessageByUniqueName(uniqueName, text string, replyNo *int64) (*PrivateMessageResponse, error) {
	if uniqueName == "" {
		return nil, fmt.Errorf("unique_name cannot be empty")
	}
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	reqBody := PrivateMessageRequest{
		Text:           text,
		UserUniqueName: uniqueName,
		ReplyNo:        replyNo,
	}

	return c.sendPrivateMessage(reqBody)
}

// sendPrivateMessage is a helper function to send private messages.
func (c *Client) sendPrivateMessage(reqBody PrivateMessageRequest) (*PrivateMessageResponse, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, "/msg/post/private", nil, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response PrivateMessageResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendPrivateReply sends a private reply to a specific message.
func (c *Client) SendPrivateReply(userID, replyPostNo int64, text string) (*PrivateMessageResponse, error) {
	return c.SendPrivateMessageByID(userID, text, &replyPostNo)
}

// BroadcastMessage sends a message to multiple chats.
func (c *Client) BroadcastMessage(chatIDs []int64, text string) ([]MessageResponse, error) {
	if len(chatIDs) == 0 {
		return nil, fmt.Errorf("chat_ids slice cannot be empty")
	}

	var responses []MessageResponse
	for _, chatID := range chatIDs {
		response, err := c.SendMessage(chatID, text, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to send message to chat %d: %w", chatID, err)
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// SendMentionMessage sends a message with a mention to all members in a chat.
func (c *Client) SendMentionMessage(chatID int64, text string) (*MessageResponse, error) {
	return c.SendMessage(chatID, "@all "+text, nil)
}

// SendMessageToAllMyChats sends a message to all chats where the bot is a member.
func (c *Client) SendMessageToAllMyChats(text string) ([]MessageResponse, error) {
	chats, err := c.GetMyChats()
	if err != nil {
		return nil, err
	}

	chatIDs := make([]int64, len(chats.Chats))
	for i, chat := range chats.Chats {
		chatIDs[i] = chat.ID
	}

	return c.BroadcastMessage(chatIDs, text)
}

// UpdateMessage updates an existing message in a chat.
//
// API: PUT /msg/post/{chat_id}/{post_no}
func (c *Client) UpdateMessage(chatID, postNo int64, updateReq *UpdateMessageRequest) (*UpdateMessageResponse, error) {
	if chatID == 0 {
		return nil, fmt.Errorf("chat_id cannot be zero")
	}
	if postNo == 0 {
		return nil, fmt.Errorf("post_no cannot be zero")
	}
	if updateReq == nil {
		return nil, fmt.Errorf("update_request cannot be nil")
	}
	if updateReq.Text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	body, err := json.Marshal(updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("/msg/post/%d/%d", chatID, postNo)
	req, err := c.newRequest(http.MethodPut, url, nil, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response UpdateMessageResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateMessageWithAttachments updates a message with new attachments.
func (c *Client) UpdateMessageWithAttachments(chatID, postNo int64, text string, attachments []string) (*UpdateMessageResponse, error) {
	updateReq := &UpdateMessageRequest{
		Text:        text,
		Attachments: attachments,
	}
	return c.UpdateMessage(chatID, postNo, updateReq)
}

// UpdateMessageWithReply updates a message and sets reply reference.
func (c *Client) UpdateMessageWithReply(chatID, postNo, replyPostNo int64, text string) (*UpdateMessageResponse, error) {
	updateReq := &UpdateMessageRequest{
		Text:    text,
		ReplyNo: &replyPostNo,
	}
	return c.UpdateMessage(chatID, postNo, updateReq)
}

// UpdateMessageE2E updates a message with E2E encryption flag.
func (c *Client) UpdateMessageE2E(chatID, postNo int64, text string, e2e bool) (*UpdateMessageResponse, error) {
	e2eFlag := e2e
	updateReq := &UpdateMessageRequest{
		Text: text,
		E2E:  &e2eFlag,
	}
	return c.UpdateMessage(chatID, postNo, updateReq)
}
