package verbosity

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// UploadFile uploads a file to the server and returns its GUID.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadFile(chatID int64, filePath string) (*FileUploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return c.UploadFileData(chatID, file, fileInfo.Size(), fileInfo.Name())
}

// UploadFileData uploads file data directly.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadFileData(chatID int64, reader io.Reader, size int64, filename string) (*FileUploadResponse, error) {
	if chatID == 0 {
		return nil, fmt.Errorf("chat_id cannot be zero")
	}
	if size <= 0 {
		return nil, fmt.Errorf("file size must be positive")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add chat_id field
	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", chatID)); err != nil {
		return nil, fmt.Errorf("failed to write chat_id field: %w", err)
	}

	// Add size field
	if err := writer.WriteField("size", fmt.Sprintf("%d", size)); err != nil {
		return nil, fmt.Errorf("failed to write size field: %w", err)
	}

	// Add file field
	part, err := writer.CreateFormFile("data", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, reader); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := c.newFileRequest("/new/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var response FileUploadResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UploadFileFromBytes uploads a file from byte slice.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadFileFromBytes(chatID int64, data []byte, filename string) (*FileUploadResponse, error) {
	if chatID == 0 {
		return nil, fmt.Errorf("chat_id cannot be zero")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("file data cannot be empty")
	}

	return c.UploadFileData(chatID, bytes.NewReader(data), int64(len(data)), filename)
}

// UploadTextFile uploads a text file.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadTextFile(chatID int64, content, filename string) (*FileUploadResponse, error) {
	if chatID == 0 {
		return nil, fmt.Errorf("chat_id cannot be zero")
	}
	if filename == "" {
		filename = "file.txt"
	}

	return c.UploadFileFromBytes(chatID, []byte(content), filename)
}

// UploadImage uploads an image file.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadImage(chatID int64, imagePath string) (*FileUploadResponse, error) {
	return c.UploadFile(chatID, imagePath)
}

// UploadDocument uploads a document file.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadDocument(chatID int64, docPath string) (*FileUploadResponse, error) {
	return c.UploadFile(chatID, docPath)
}

// UploadAudio uploads an audio file.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadAudio(chatID int64, audioPath string) (*FileUploadResponse, error) {
	return c.UploadFile(chatID, audioPath)
}

// UploadVideo uploads a video file.
//
// API: POST https://file.verbosity.io/new/upload
func (c *Client) UploadVideo(chatID int64, videoPath string) (*FileUploadResponse, error) {
	return c.UploadFile(chatID, videoPath)
}

// UploadToMultipleChats uploads a file and sends it to multiple chats.
//
// Returns the file GUID and any error.
func (c *Client) UploadToMultipleChats(filePath string, chatIDs []int64) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	uploadResponse, err := c.UploadFileData(0, file, fileInfo.Size(), fileInfo.Name())
	if err != nil {
		return "", err
	}

	for _, chatID := range chatIDs {
		_, err := c.UploadFileData(chatID, bytes.NewReader([]byte{}), fileInfo.Size(), fileInfo.Name())
		if err != nil {
			return uploadResponse.GUID, fmt.Errorf("failed to upload to chat %d: %w", chatID, err)
		}
	}

	return uploadResponse.GUID, nil
}
